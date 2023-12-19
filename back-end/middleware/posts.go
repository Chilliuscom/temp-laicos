package middleware

import (
	"SocialNetwork/DB"
	"SocialNetwork/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const imgDir = "./pkg/img"

func CreateComment(w http.ResponseWriter, r *http.Request) {
	//send notification to post creator and other commenters
	r.ParseMultipartForm(10 << 20)
	master_post, errs := strconv.Atoi(r.PostFormValue("masterPost"))
	content := r.PostFormValue("content")
	creator_id, errs := strconv.Atoi(r.PostFormValue("creator_id"))
	created_at := time.Now().Format("02.01.2006 15:04")
	if errs != nil {
		fmt.Println(errs)
	}

	_, _, err := r.FormFile("avatar")
	var img string
	if err != nil {
		img = ""
	}
	if err == nil {
		img = saveIMG(w, r)
	}
	log.Println(master_post, content, creator_id, created_at, img)
	comment := models.Comment{Text: content, User_id: creator_id, Post_id: master_post, Created_at: created_at, Image: img}
	if result := DB.ConnectDB().Table("comment").Create(&comment); result != nil {
		respond(w, message(500, "Some error. Try again later or contact the administration"))
	}
	respond(w, message(200, "COMMENT CREATION SUCCESS"))
}

func CreateGroupPost(w http.ResponseWriter, r *http.Request) {
	//send notification to group users
	r.ParseMultipartForm(10 << 20)
	group_id, errs := strconv.Atoi(r.PostFormValue("group_id"))
	header := r.PostFormValue("header")
	content := r.PostFormValue("content")
	creator_id, errs := strconv.Atoi(r.PostFormValue("creator_id"))
	created_at := time.Now().Format("02.01.2006 15:04")
	if errs != nil {
		fmt.Println(errs)
	}

	_, _, err := r.FormFile("avatar")
	var img string
	if err != nil {
		img = ""
	}
	if err == nil {
		img = saveIMG(w, r)
	}
	log.Println(group_id, header, content, creator_id, created_at, img)
	//SQL
	groupPost := models.Grouppost{Header: header, Text: content, User_id: creator_id, Group_id: group_id, Created_at: created_at, Image: img}
	if result := DB.ConnectDB().Table("group_posts").Create(&groupPost); result != nil {
		respond(w, message(500, "Some error. Try again later or contact the administration"))
	}
	respond(w, message(200, "GROUP POST CREATION SUCCESS"))
	//add to db header(header), text(content), user_id(creator_id), group_id(group_id), current time(created_at), img

}

func CreateGroupComment(w http.ResponseWriter, r *http.Request) {
	//send notification to post creator and other commenters
	r.ParseMultipartForm(10 << 20)
	master_post, errs := strconv.Atoi(r.PostFormValue("masterPost"))
	content := r.PostFormValue("content")
	creator_id, errs := strconv.Atoi(r.PostFormValue("creator_id"))
	created_at := time.Now().Format("02.01.2006 15:04")
	if errs != nil {
		fmt.Println(errs)
	}
	_, _, err := r.FormFile("avatar")
	var img string
	if err != nil {
		img = ""
	}
	if err == nil {
		img = saveIMG(w, r)
	}
	log.Println(master_post, content, creator_id, created_at, img)
	//SQL
	groupComment := models.Comment{Text: content, User_id: creator_id, Post_id: master_post, Created_at: created_at, Image: img}
	if result := DB.ConnectDB().Table("comment").Create(&groupComment); result != nil {
		respond(w, message(500, "Some error. Try again later or contact the administration"))
	}
	respond(w, message(200, "GROUP COMMENT CREATION SUCCESS"))
	//add to db text(content), user_id(creator_id), post_id(master_post), current time(created_at), IMG

}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	//Send notification to followers
	r.ParseMultipartForm(10 << 20)
	header := r.PostFormValue("header")
	content := r.PostFormValue("content")
	creator_id, errs := strconv.Atoi(r.PostFormValue("creator_id"))
	created_at := time.Now().Format("02.01.2006 15:04")
	if errs != nil {
		fmt.Println(errs)
	}

	_, _, err := r.FormFile("avatar")
	var img string
	if err != nil {
		img = ""
	}
	if err == nil {
		img = saveIMG(w, r)
	}
	log.Println(header, content, creator_id, created_at, img)
	//SQL

	post := models.Post{Header: header, Text: content, User_id: creator_id, Created_at: created_at, Image: img}
	if result := DB.ConnectDB().Table("post").Create(&post); result != nil {
		respond(w, message(500, "Some error. Try again later or contact the administration"))
	}
	respond(w, message(200, "POST CREATION SUCCESS"))
	//add to db header(header), text(content), user_id(creator_id), current time(created_at), img

}

func saveIMG(w http.ResponseWriter, r *http.Request) string {
	file, header, _ := r.FormFile("avatar")
	defer file.Close()
	contentType := header.Header.Get("Content-Type")
	fileExtension := ".jpg" // Default extension if not recognized

	if contentType == "image/jpeg" {
		fileExtension = ".jpg"
	} else if contentType == "image/png" {
		fileExtension = ".png"
	} else if contentType == "image/gif" {
		fileExtension = ".gif"
	}
	uniqueFilename := uuid.New().String() + fileExtension
	filePath := filepath.Join(imgDir, uniqueFilename)
	outFile, err := os.Create(filePath)
	if err != nil {
		respond(w, message(500, "Failed to create file"))
		return ""
	}
	defer outFile.Close()
	_, err = io.Copy(outFile, file)
	if err != nil {
		respond(w, message(500, "Failed to save file"))
		return ""
	}
	return filePath
}

func message(status int, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
