package DB

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("social-network.db"), &gorm.Config{TranslateError: true})
	if err != nil {
		fmt.Println(err)
	}
	return db
}
