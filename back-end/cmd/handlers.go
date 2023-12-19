package cmd

import (
	"SocialNetwork/DB"
	"SocialNetwork/middleware"
	"SocialNetwork/models"
	"SocialNetwork/websoct"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

const imgDir = "./pkg/img"

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if validateSession(w, r) { //if already logged in - should not be needed, router shouldnt allow
			fmt.Println("falsee2")
			getUserBySession(r)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

		var input models.Login
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		authErr := middleware.FindUserAndComparePassword(DB.ConnectDB(), input.Email, input.Password)
		switch authErr {
		case "OK":
			id := middleware.GetUserIDByEmail(input.Email)
			session, _ := store.Get(r, "session-id")
			session.Values["email"] = input.Email
			session.Values["id"] = id
			session.Values["expires"] = time.Now().Add(time.Second * time.Duration(3600))
			err = session.Save(r, w)
			if err != nil {
				fmt.Println(err)
			}
			response := models.Response{
				Status:  http.StatusOK,
				Message: "SUCCESSFUL LOGIN",
				Id:      id,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			//respond(w, message(200, "SUCCESSFUL LOGIN"))
			//create new session, assign cookie, return user details n stuff
		case "INVALIDPW":
			println("invalidpw")
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		case "NOUSER":
			println("nouser")
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func register(w http.ResponseWriter, r *http.Request) { //ADD PASSWORD HASHING
	if r.Method == "POST" {
		if validateSession(w, r) { //if already logged in - should not be needed, router shouldnt allow
			getUserBySession(r)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var input models.User

		err := r.ParseMultipartForm(10 << 20) // Max 10MB in memory, remaining to disk
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		input.Email = r.PostFormValue("email")
		input.Password = r.PostFormValue("password")
		input.FirstName = r.PostFormValue("firstName")
		input.LastName = r.PostFormValue("lastName")
		input.Username = r.PostFormValue("username")
		input.DateOfBirth = r.PostFormValue("dateOfBirth")
		input.Gender = r.PostFormValue("gender")
		input.Country = r.PostFormValue("country")
		input.Phone = r.PostFormValue("phone")
		input.Privacy = r.PostFormValue("privacy")
		input.About_Me = r.PostFormValue("aboutMe")
		input.Username = r.PostFormValue("username")
		if input.Username == "" {
			input.Username = "NULL"
		}
		if input.About_Me == "" {
			input.About_Me = "NULL"
		}

		_, _, err = r.FormFile("avatar")
		var avatar string
		if err != nil {
			avatar = "./pkg/img/non-avatar.png"
		}
		if err == nil {
			avatar = saveAvatar(w, r)
		}
		hashpassword, err := HashPassword(input.Password)
		user := models.User{FirstName: input.FirstName, LastName: input.LastName, Username: input.Username, Email: input.Email,
			Phone: input.Phone, Password: hashpassword, Avatar: avatar, Privacy: input.Privacy,
			DateOfBirth: input.DateOfBirth, Gender: input.Gender, About_Me: input.About_Me, Country: input.Country, Created_at: time.Now().Format("2006-01-02 15:04:05")}
		var users []int
		DB.ConnectDB().Table("users").Select("id").Scan(&users)
		fmt.Println(users)
		if result := DB.ConnectDB().Create(&user); result != nil {
			if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
				http.Error(w, "Email or Username already exists", http.StatusConflict)
				//deletes avatar if fails
				if avatar != "./pkg/img/non-avatar.png" {
					err := os.Remove(avatar)
					if err != nil {
						fmt.Println("Error deleting file:", err)
						return
					}
				}
				return
			}
			/* 			http.Error(w, "Internal server error. Try again later or contact the administration", http.StatusInternalServerError)
			   			return */
		}
		for _, userId := range users {
			follower := models.Followers{First_user: int(user.ID), Second_user: userId}
			DB.ConnectDB().Table("following").Create(&follower)
			follow := models.Followers{First_user: userId, Second_user: int(user.ID)}
			DB.ConnectDB().Table("following").Create(&follow)
		}

		respond(w, message(200, "ACCOUNT CREATION SUCCESS"))

		return
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		session, _ := store.Get(r, "session-id")
		if session == nil {
			http.Error(w, "Already logged out", http.StatusBadRequest)
			return
		}
		delete(session.Values, "email")
		err := session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{ //empties the cookie
			Name:     "session-id",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		})
		respond(w, message(200, "Successfully logged out"))
	case "POST":
		respond(w, message(405, "Method not allowed"))
	}
}

func submitContent(w http.ResponseWriter, r *http.Request) {
	log.Println("submitContent")
	var data models.Submit
	creator_id := getUserBySession(r)
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data.Creator_id = creator_id
	switch data.Type {
	case "post":
		middleware.CreatePost(w, r)
	case "groupPost":
		middleware.CreateGroupPost(w, r)
	case "comment":
		middleware.CreateComment(w, r)
	case "groupComment":
		middleware.CreateGroupComment(w, r)
	case "followRequest":
		if !middleware.SendFollowRequest(data) {
			respond(w, message(200, "This should not happen on successful request,change to error later"))
			log.Println("This should not happen on successful request,change to error later")
			//http.Error(w, "Error trying to follow, try again later", http.StatusInternalServerError)
		} else {
			respond(w, message(200, "Follow request successful"))
		}
	case "followRequestUpdate":
		middleware.FollowRequestUpdate(data)
	case "groupJoinRequest":
		middleware.SendGroupJoinRequest(data)
	case "joinUpdate":
		middleware.JoinUpdate(data)
	case "groupInvite":
		middleware.SendGroupInvite(data)
	case "inviteUpdate":
		middleware.InviteUpdate(data)
	}
}

func getGroup(w http.ResponseWriter, r *http.Request) {
	//groupName := strings.Split(r.URL.String()[len(`/groups/`):], `/`)
	//group := GetGroup(groupName[0])
	//fmt.Println(group)
}

func profile(w http.ResponseWriter, r *http.Request) {
	//id := strings.Split(r.URL.String()[len(`/id/`):], `/`)
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "No id provided", http.StatusBadRequest)
		return
	}
	profile := middleware.GetProfile(id)
	if getUserBySession(r) == 0 {
		respond(w, message(200, "NOT LOGGED IN"))
		return
	}

	if !middleware.CheckIfUserExists(id) {
		respond(w, message(200, "USER DOES NOT EXIST"))
		return
	}
	requestUser := getUserBySession(r)
	secondUser := id
	data := make(map[string]interface{})

	if requestUser != secondUser && profile.Privacy == "0" { //if profile is private
		if !middleware.CheckIfFollows(requestUser, secondUser) { //and if not following said user
			//respond(w, message(200, "NOT_FOLLOWING"))
			//log.Println(getUserBySession(r), "NOT following", id[0], secondUser)
			limitedProfile := middleware.GetLimitedProfile(secondUser)
			data["profile"] = limitedProfile
			log.Println(data)
			respond(w, data)
			return
		}
	}
	data["profile"] = profile
	data["followers"] = middleware.GetFollowers(secondUser)
	data["following"] = middleware.GetRelationships(secondUser)
	respond(w, data)
}

func messageRequest(w http.ResponseWriter, r *http.Request) {
	var query models.MessageRequest
	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data := make(map[string]interface{})
	messages := websoct.GetMessageFromDB(query.MessageID, query.Amount, query.Chatroom_id)
	data["messages"] = messages
	respond(w, data)
}

func testFunction(w http.ResponseWriter, r *http.Request) {
	websoct.CreateChatroom(1, 3)
	if !validateSession(w, r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

}

func validateSessionHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-id")
	if err != nil {
		http.Error(w, "Unable to validate session", http.StatusInternalServerError)
		return
	}
	if session.Values["email"] != nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Error(w, "Unauthorized: Session is not valid", http.StatusForbidden)
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("imageHandler")
	imageName := r.URL.Path[len("api/images/"):]
	println(imageName)
	imagePath := fmt.Sprintf(imageName)
	println(imagePath)
	contentType := "image/jpeg"

	if strings.HasSuffix(imageName, ".png") {
		contentType = "image/png"
	}
	w.Header().Set("Content-Type", contentType)

	http.ServeFile(w, r, imagePath)
}
