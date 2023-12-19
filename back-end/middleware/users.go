package middleware

import (
	"SocialNetwork/DB"
	"SocialNetwork/models"
	"SocialNetwork/websoct"
	"fmt"
	"log"
)

func SendFollowRequest(data models.Submit) bool {
	from := data.Creator_id
	to := data.Target_user_id
	println(from, " wants to follow ", to)

	//send follow request notification to said user
	follower := models.Followers{First_user: from, Second_user: to, Status: "PENDING"}
	//SQL
	if result := DB.ConnectDB().Table("following").Model(&follower).Where("first_user = ? AND second_user = ?", from, to).Update("Status", "pending"); result != nil {
		fmt.Println("Bad following")
		log.Println(result, follower)
		return false
	}
	websoct.Notify("FollowRequest", from, "", []string{GetUsernameByUserID(to)})
	return true
	//change in db first_user(data.creator_id), second_user(data.target_user_id), status PENDING
}

func FollowRequestUpdate(data models.Submit) {
	to := data.Creator_id
	from := data.Target_user_id
	update := data.Content
	var following models.Followers
	log.Println(to, " has ", update, "ed follow request from ", from)
	//SQL
	//change pending to false or true
	switch update {
	case "true":
		if result := DB.ConnectDB().Table("following").Model(&following).Where("first_user = ? AND second_user = ?", to, from).Update("Status", "true"); result != nil {
			fmt.Println("Bad following")
			return
		}
		websoct.Notify("FollowAccept", from, "", []string{GetUsernameByUserID(to)})
		websoct.CreateChatroom(to, from)
	case "false":
		if result := DB.ConnectDB().Table("following").Model(&following).Where("first_user = ? AND second_user = ?", to, from).Update("Status", "false"); result != nil {
			fmt.Println("Bad following")
			return
		}
		websoct.Notify("FollowDecline", from, "", []string{GetUsernameByUserID(to)})
	}

	//if following is created either way, add a chatroom to db
}

func GetProfile(id int) models.User {
	var profile models.User
	Database.Raw(` 
	SELECT
	*

	FROM users 
	WHERE users.id = ?`, id).Scan(&profile)
	profile.Password = ""
	//profile.Followers = GetFollowers(GetUserIDByUsername(username))
	//profile.Following = GetRelationships(GetUserIDByUsername(username))
	return profile
}

func GetUserIDByUsername(username string) int {
	var userid int
	Database.Raw(` 
	SELECT
	users.id

	FROM users 
	WHERE users.username = ?`, username).Scan(&userid)
	return userid
}

func GetUserIDByEmail(email string) int {
	var userid int
	Database.Raw(` 
	SELECT
	users.id

	FROM users 
	WHERE users.email = ?`, email).Scan(&userid)
	return userid
}

func GetUsernameByUserID(userID int) string {
	var username string
	Database.Raw(` 
	SELECT
	users.email

	FROM users 
	WHERE users.id = ?`, userID).Scan(&username)
	return username
}

func GetRelationships(user int) []models.FollowersInfo {
	var friendships []models.FollowersInfo
	Database.Raw(` 
	SELECT f.id, f.first_user, f.second_user, f.status, u.first_name AS FirstName, u.last_name AS LastName
	FROM following f
	INNER JOIN users u ON f.second_user = u.id
	WHERE f.first_user = ?`, user).Scan(&friendships)
	return friendships
}

func CheckIfFollows(user1, user2 int) bool {
	following := GetRelationships(user1)
	follows := false
	for i := range following {
		if following[i].Second_user == user2 && following[i].Status == "TRUE" {
			follows = true
			break
		}
	}
	return follows
}

func GetUsernameByEmail(email string) string {
	var username string
	Database.Raw(` 
	SELECT
	username

	FROM users 
	WHERE users.email = ?`, email).Scan(&username)
	return username
}

func GetFollowers(user int) []models.FollowersInfo {
	var followers []models.FollowersInfo
	Database.Raw(` 
	SELECT f.id, f.first_user, f.second_user, f.status, u.first_name AS FirstName, u.last_name AS LastName
	FROM following f
	INNER JOIN users u ON f.first_user = u.id
	WHERE f.second_user = ?`, user).Scan(&followers)
	return followers
}

func GetLimitedProfile(user int) models.User {
	var profile models.User
	Database.Raw(` 
	SELECT
	users.id,
	users.first_name,
	users.last_name,
	users.username,
	users.email,
	users.avatar,
	users.privacy

	FROM users 
	WHERE users.id = ?`, user).Scan(&profile)
	return profile
}
