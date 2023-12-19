package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

func message(status int, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func saveAvatar(w http.ResponseWriter, r *http.Request) string {
	file, header, _ := r.FormFile("avatar")
	defer file.Close()
	contentType := header.Header.Get("Content-Type")
	fileExtension := ".jpg" // Default extension if not recognized

	if contentType == "image/jpeg" {
		fileExtension = ".jpg"
	} else if contentType == "image/png" {
		fileExtension = ".png"
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
