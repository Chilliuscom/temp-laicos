package middleware

import (
	"SocialNetwork/DB"
	"SocialNetwork/models"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var Database = DB.ConnectDB()

func FindUserAndComparePassword(db *gorm.DB, email, password string) string {
	//if user exists and pw correct return "OK"
	//if user exists but pw incorrect return "INVALIDPW"
	//if user doesn't exist return "NOUSER"
	var user models.User
	res := db.First(&user, "email = ?", email).Error
	if errors.Is(res, gorm.ErrRecordNotFound) {
		return "NOUSER"
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println(err)
		return "INVALIDPW"
	}
	return "OK"
}

func CheckIfUserExists(id int) bool {
	//If user with this email exists return true
	var user models.User
	res := DB.ConnectDB().First(&user, "id = ?", id).Error
	if errors.Is(res, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
