package infrastructure

import (
	"idp/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func DbInit() {
	db, err := gorm.Open(sqlite.Open("idp.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Client{})
	db.AutoMigrate(&models.AccessToken{})
}

func Connect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("idp.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}
