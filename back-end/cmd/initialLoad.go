package cmd

import (
	"SocialNetwork/middleware"
	"SocialNetwork/models"
	"SocialNetwork/websoct"

	"net/http"
)

func initialLoad(w http.ResponseWriter, r *http.Request) {
	username := middleware.GetUsernameByUserID(getUserBySession(r))
	user := middleware.GetProfile(getUserBySession(r))
	chatrooms := websoct.GetRoomsContainingUser(username)
	var messages []models.Message
	for roomID := range chatrooms {
		messages = append(messages, websoct.GetMessageFromDB(0, 1, roomID)[0])
	}
	data := make(map[string]interface{})
	data["data"] = models.InitialData{UserData: user, Chatrooms: chatrooms, Messages: messages}

	respond(w, data)
}
