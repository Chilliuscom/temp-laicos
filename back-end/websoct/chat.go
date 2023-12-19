package websoct

import (
	"SocialNetwork/DB"
	"SocialNetwork/models"
	"log"
	"time"
)

var Database = DB.ConnectDB()

func GetRoomUsers(roomID int) []string {
	var roomUsers []models.Chatroom
	Database.Find(&roomUsers, "id = ?", roomID)
	result := []string{}
	for i := range roomUsers {
		result = append(result, GetUsernameByID(roomUsers[i].First_user))
	}
	return result
}

func AddMessageToDB(message, from string, roomID int) {
	userID := GetIDByUsername(from)
	now := time.Now().Format("2006-01-02 15:04:05")
	MSG := models.Message{Message: message, User_id: userID, Chatroom_id: roomID, Timestamp: now}
	Database.Create(&MSG)
}

func GetMessageFromDB(messageID, amount, roomID int) []models.Message {
	var data []models.Message
	if messageID == 0 {
		messageID = 999999999999
	}
	if err := Database.Raw(`
	SELECT
	messages.id,
	messages.message,
	messages.user_id,
	messages.chatroom_id,
	messages.timestamp 
	FROM messages 
	WHERE messages.chatroom_id = ?
	AND messages.id < ? 
	ORDER BY messages.id DESC
	LIMIT ?`, roomID, messageID, amount).Scan(&data); err != nil {
		//log.Println(err, "CHAT_DB_ERR")
	}
	return data
}

func GetIDByUsername(username string) int {
	var user models.User
	Database.First(&user, "email = ?", username)
	return int(user.ID)
}

func GetUsernameByID(userID int) string {
	var username string
	Database.Raw(` 
	SELECT
	users.email

	FROM users 
	WHERE users.id = ?`, userID).Scan(&username)
	return username
}

func GetRoomsContainingUser(username string) []models.Chatroom {
	var chatrooms []models.Chatroom
	Database.Find(&chatrooms, "username = ?", username)
	return chatrooms
}

func CreateChatroom(user1, user2 int) {
	var entry1 models.Chatroom
	var entry2 models.Chatroom
	var lastChatroomID int
	if err := Database.Raw(`
	SELECT
	chatrooms.chatroom_id
	FROM chatrooms
	ORDER BY chatrooms.id DESC
	LIMIT 1
	`).Scan(&lastChatroomID); err != nil {
		log.Println(err)
		//lastChatroomID = 0
	}
	lastChatroomID++
	log.Println(lastChatroomID, entry1, entry2, user1, user2)
	entry1.Chatroom_id = lastChatroomID
	entry2.Chatroom_id = lastChatroomID
	entry1.First_user = user1
	entry2.First_user = user2
	if !CheckIfChatroomExists(user1, user2) {
		Database.Create(&entry1)
		Database.Create(&entry2)
	}
}

func CheckIfChatroomExists(user1, user2 int) bool {
	var chatroom_id int
	//exists := false
	Database.Raw(`SELECT chatroom_id
	FROM chatrooms
	WHERE (
		(first_user = ? AND EXISTS (SELECT 1 FROM chatrooms WHERE chatroom_id = chatrooms.chatroom_id AND first_user = ?))
		OR
		(first_user = ? AND EXISTS (SELECT 1 FROM chatrooms WHERE chatroom_id = chatrooms.chatroom_id AND first_user = ?))
	)
	GROUP BY chatroom_id
	HAVING COUNT(*) = 2;
	`, user1, user2, user2, user1).Scan(&chatroom_id)
	log.Println(chatroom_id)
	if chatroom_id == 0 {
		return false
	} else {
		return true
	}
}
