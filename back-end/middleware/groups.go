package middleware

import (
	"SocialNetwork/DB"
	"SocialNetwork/models"
	"SocialNetwork/websoct"

	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

func SendGroupJoinRequest(data models.Submit) {
	//send group join request notification to group owner
	creator := data.Creator_id
	group := data.Group_id
	log.Println(creator, " wants to join group ", group)
	groupData := GetGroupByID(group)
	request := models.GroupRequest{Group_id: group, User_id: creator, Pending_type: "REQUEST", Status: "PENDING"}
	if result := DB.ConnectDB().Table("members_of_group").Create(&request); result != nil {
		fmt.Println("Some errors in groupJoinRequest")
		return
	}
	websoct.Notify("GroupJoinRequest", creator, groupData.Name, []string{GetUsernameByUserID(groupData.Creator)})
}

func JoinUpdate(data models.Submit) {
	//notify all group users and requester
	to := data.Creator_id
	group_id := data.Group_id
	groupData := GetGroupByID(group_id)

	update := data.Content
	var request models.GroupRequest
	log.Println(to, " join request to group with id ", group_id, " has been ", update, "ed")
	//SQL
	//pendingtype "" status true or false
	switch update {
	case "true":
		if result := DB.ConnectDB().Table("members_of_group").Model(&request).Where("group_id = ? AND user_id = ?", group_id, to).Updates(map[string]interface{}{"status": "true", "pending_type": ""}); result != nil {
			fmt.Println("Bad update request true")
			return
		}
		websoct.Notify("GroupJoinRequestAccept", groupData.Creator, groupData.Name, []string{GetUsernameByUserID(to)})

	case "false":
		if result := DB.ConnectDB().Table("members_of_group").Model(&request).Where("group_id = ? AND user_id = ?", group_id, to).Updates(map[string]interface{}{"status": "false", "pending_type": ""}); result != nil {
			fmt.Println("Bad update request false")
			return
		}
		websoct.Notify("GroupJoinRequestDecline", groupData.Creator, groupData.Name, []string{GetUsernameByUserID(to)})

	}
}

func SendGroupInvite(data models.Submit) {
	//send group invite notification to the invited person
	target := data.Target_user_id
	from := data.Creator_id
	group := data.Group_id
	groupData := GetGroupByID(group)

	log.Println(target, " has been invited by ", from, " to join group ", group)
	//SQL
	//add to DB group members status PENDING pendingtype INVITE
	request := models.GroupRequest{Group_id: group, User_id: target, Pending_type: "INVITE", Status: "PENDING"}
	if result := DB.ConnectDB().Table("members_of_group").Create(&request); result != nil {
		fmt.Println("Some errors in groupJoinRequest")
		return
	}
	websoct.Notify("GroupInvite", from, groupData.Name, []string{GetUsernameByUserID(target)})

}

func InviteUpdate(data models.Submit) {
	//notify all group users and owner
	to := data.Target_user_id
	from := data.Creator_id
	group_id := data.Group_id
	groupData := GetGroupByID(group_id)

	update := data.Content
	log.Println(to, " has ", update, "ed ", from, "'s invite to group with id ", group_id)
	var request models.GroupRequest
	//SQL
	//pendingtype "" status true or false
	switch update {
	case "true":
		if result := DB.ConnectDB().Table("members_of_group").Model(&request).Where("group_id = ? AND user_id = ?", group_id, to).Updates(map[string]interface{}{"status": "true", "pending_type": ""}); result != nil {
			fmt.Println("Bad update request true")
			return
		}
		websoct.Notify("GroupInviteAccept", to, groupData.Name, []string{GetUsernameByUserID(from)})

	case "false":
		if result := DB.ConnectDB().Table("members_of_group").Model(&request).Where("group_id = ? AND user_id = ?", group_id, to).Updates(map[string]interface{}{"status": "false", "pending_type": ""}); result != nil {
			fmt.Println("Bad update request false")
			return
		}
		websoct.Notify("GroupInviteDecline", to, groupData.Name, []string{GetUsernameByUserID(from)})

	}
}

func GetGroup(name string) models.Group {
	var group models.Group
	//SQL
	Database.Raw(` 
	SELECT
	*

	FROM groups 
	WHERE groups.name = ?`, name).Scan(&group)
	return group
}

func GetGroupByID(id int) models.Group {
	var group models.Group
	//SQL
	Database.Raw(` 
	SELECT
	*

	FROM groups 
	WHERE groups.id = ?`, id).Scan(&group)
	return group
}

func CheckIfGroupExists(name string) bool {
	//If user with this email exists return true
	var group models.Group
	//SQL
	res := DB.ConnectDB().Table("group").First(&group, "name = ?", name).Error //no such table "groups"?
	if errors.Is(res, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func CheckIfUserIsMember(group_id, user_id int) bool {

	//SQL
	var members models.GroupRequest
	//group_id and user id, if status True then return true that mean user already is member
	res := DB.ConnectDB().Table("members_of_group").First(&members, "user_id = ? AND group_id = ?", user_id, group_id).Error //no such table "groups"?
	if errors.Is(res, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func GetGroupPosts(groupID, postID int) []models.Grouppost {
	var groupPosts []models.Grouppost
	DB.ConnectDB().Where("group_id = ? AND id = ?", groupID, postID).Table("post").Find(&groupPosts)
	return groupPosts
}
