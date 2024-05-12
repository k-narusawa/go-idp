package gateway

import (
	"github.com/k-narusawa/go-idp/authorization/domain/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dsn string

func DbInit(mode, inputDsn string) {
	if mode == "sqlite3" {
		dsn = "idp.db"
	} else if mode == "postgres" {
		dsn = inputDsn
	}

	dsn = inputDsn

	var db *gorm.DB
	var err error
	switch mode {
	case "sqlite3":
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})

	case "postgres":
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		panic("failed to connect database")
	}

	if mode == "sqlite3" {
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
	}

	testUser := models.NewUser("test@example.com", "!Password0")
	db.Save(&testUser)
}

func Connect() *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}
