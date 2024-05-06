package gateway

import (
	"github.com/k-narusawa/go-idp/authorization/domain/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func DbInit() {
	db, err := gorm.Open(sqlite.Open("idp.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Client{})
	db.AutoMigrate(&models.OidcSession{})
	db.AutoMigrate(&models.AuthorizationCode{})
	db.AutoMigrate(&models.AccessToken{})
	db.AutoMigrate(&models.RefreshToken{})
	db.AutoMigrate(&models.PKCE{})
	db.AutoMigrate(&models.LoginSkipSession{})

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.WebauthnCredential{})
	db.AutoMigrate(&models.WebauthnSessionData{})

	testUser := models.NewUser("test@example.com", "!Password0")
	db.Save(&testUser)
}

func Connect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("idp.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}
