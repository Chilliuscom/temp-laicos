package websoct

import (
	"bytes"
	"encoding/json"
	"log"
	"strconv"
	"time"

	//"SocialNetwork/middleware"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

// CreateNewSocketUser creates a new socket user
func CreateNewSocketUser(hub *Hub, connection *websocket.Conn, username string) {
	uniqueID := uuid.New()
	client := &Client{
		hub:                 hub,
		webSocketConnection: connection,
		send:                make(chan SocketEventStruct),
		username:            username,
		userID:              uniqueID.String(),
	}

	go client.writePump()
	go client.readPump()

	client.hub.register <- client
}

// HandleUserRegisterEvent will handle the Join event for New socket users
func HandleUserRegisterEvent(hub *Hub, client *Client) {
	hub.clients[client] = true
	handleSocketPayloadEvents(client, SocketEventStruct{
		EventName:    "join",
		EventPayload: client.userID,
	})
}

// HandleUserDisconnectEvent will handle the Disconnect event for socket users
func HandleUserDisconnectEvent(hub *Hub, client *Client) {
	_, ok := hub.clients[client]
	if ok {
		delete(hub.clients, client)
		close(client.send)

		handleSocketPayloadEvents(client, SocketEventStruct{
			EventName:    "disconnect",
			EventPayload: client.userID,
		})
	}
}

// EmitToSpecificClient will emit the socket event to specific socket user
func EmitToSpecificClient(hub *Hub, payload SocketEventStruct, userID string) {
	for client := range hub.clients {
		if client.userID == userID {
			select {
			case client.send <- payload:
				/*default:
				close(client.send)
				delete(hub.clients, client)*/ //DISABLED FOR DIAGNOSTIC PURPOSES
			}
		}
	}
}

// BroadcastSocketEventToAllClient will emit the socket events to all socket users
func BroadcastSocketEventToAllClient(hub *Hub, payload SocketEventStruct) {
	for client := range hub.clients {
		select {
		case client.send <- payload:
			/*default:
			close(client.send)
			delete(hub.clients, client)*/ //DISABLED FOR DIAGNOSTIC PURPOSES
		}
	}
}

func BroadcastJoinEvent(hub *Hub, userID string) {
	for client := range hub.clients {
		EmitToSpecificClient(client.hub, SocketEventStruct{
			EventName: "join",
			EventPayload: JoinDisconnectPayload{
				UserID: userID,
				Users:  getAllConnectedUsers(client.hub),
			},
		}, client.userID)
	}
}

func handleSocketPayloadEvents(client *Client, socketEventPayload SocketEventStruct) {
	var socketEventResponse SocketEventStruct
	switch socketEventPayload.EventName {
	case "join":
		log.Printf("Join Event triggered")
		BroadcastSocketEventToAllClient(client.hub, SocketEventStruct{
			EventName: socketEventPayload.EventName,
			EventPayload: JoinDisconnectPayload{
				UserID: client.userID,
				Users:  getAllConnectedUsers(client.hub),
			},
		})

	case "disconnect":
		log.Printf("Disconnect Event triggered")
		BroadcastSocketEventToAllClient(client.hub, SocketEventStruct{
			EventName: socketEventPayload.EventName,
			EventPayload: JoinDisconnectPayload{
				UserID: client.userID,
				Users:  getAllConnectedUsers(client.hub),
			},
		})

	case "message":
		log.Printf("Message Event triggered")
		selectedRoomID, _ := strconv.Atoi(socketEventPayload.EventPayload.(map[string]interface{})["roomID"].(string))
		socketEventResponse.EventName = "message response"
		socketEventResponse.EventPayload = map[string]interface{}{
			"message": socketEventPayload.EventPayload.(map[string]interface{})["message"],
			"from":    client.username,
			"room":    selectedRoomID,
		}
		AddMessageToDB(socketEventPayload.EventPayload.(map[string]interface{})["message"].(string), client.username, selectedRoomID)
		log.Println(selectedRoomID)
		EmitToRoom(client.hub, socketEventResponse, selectedRoomID)

	case "notification":
		log.Printf("Notification Event triggered")
		socketEventResponse.EventName = "notification"
		socketEventResponse.EventPayload = map[string]interface{}{
			"name":  socketEventPayload.EventPayload.(map[string]interface{})["name"],
			"from":  socketEventPayload.EventPayload.(map[string]interface{})["from"],
			"where": socketEventPayload.EventPayload.(map[string]interface{})["where"],
			"link":  socketEventPayload.EventPayload.(map[string]interface{})["link"],
		}
		recipient := socketEventPayload.EventPayload.(map[string]interface{})["to"].(string)
		EmitToSpecificClient(client.hub, socketEventResponse, getUserIDByUsername(client.hub, recipient))

	}
}

func EmitToRoom(hub *Hub, payload SocketEventStruct, roomID int) {
	roomUsers := GetRoomUsers(roomID)

	for recipient := range roomUsers {
		EmitToSpecificClient(hub, payload, getUserIDByUsername(hub, roomUsers[recipient]))
	}
}

func getUserIDByUsername(hub *Hub, username string) string {
	var userID string
	for client := range hub.clients {
		if client.username == username {
			userID = client.userID
		}
	}
	return userID
}

func getUsernameByUserID(hub *Hub, userID string) string {
	var username string
	for client := range hub.clients {
		if client.userID == userID {
			username = client.username
		}
	}
	return username
}

func getAllConnectedUsers(hub *Hub) []UserStruct {
	var users []UserStruct
	for singleClient := range hub.clients {
		users = append(users, UserStruct{
			Username: singleClient.username,
			UserID:   singleClient.userID,
		})
	}
	return users
}

func (c *Client) readPump() {
	var socketEventPayload SocketEventStruct

	defer unRegisterAndCloseConnection(c)

	setSocketPayloadReadConfig(c)

	for {
		_, payload, err := c.webSocketConnection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error ===: %v", err)
			}
			break
		}
		decoder := json.NewDecoder(bytes.NewReader(payload))
		decoderErr := decoder.Decode(&socketEventPayload)

		if decoderErr != nil {

			log.Printf("error: %v", decoderErr)
			break

		}
		handleSocketPayloadEvents(c, socketEventPayload)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.webSocketConnection.Close()
	}()
	for {
		select {
		case payload, ok := <-c.send:
			reqBodyBytes := new(bytes.Buffer)
			json.NewEncoder(reqBodyBytes).Encode(payload)
			finalPayload := reqBodyBytes.Bytes()

			c.webSocketConnection.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.webSocketConnection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.webSocketConnection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(finalPayload)

			n := len(c.send)
			for i := 0; i < n; i++ {
				json.NewEncoder(reqBodyBytes).Encode(<-c.send)
				w.Write(reqBodyBytes.Bytes())
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.webSocketConnection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.webSocketConnection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func unRegisterAndCloseConnection(c *Client) {
	c.hub.unregister <- c
	c.webSocketConnection.Close()
}

func setSocketPayloadReadConfig(c *Client) {
	c.webSocketConnection.SetReadLimit(maxMessageSize)
	c.webSocketConnection.SetReadDeadline(time.Now().Add(pongWait))
	c.webSocketConnection.SetPongHandler(func(string) error { c.webSocketConnection.SetReadDeadline(time.Now().Add(pongWait)); return nil })
}
