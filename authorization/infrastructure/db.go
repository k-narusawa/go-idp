package infrastructure

import (
	"idp/authorization/models"

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
	db.AutoMigrate(&models.IDSession{})
	db.AutoMigrate(&models.AuthorizationCode{})
	db.AutoMigrate(&models.AccessToken{})
	db.AutoMigrate(&models.RefreshToken{})

	testUser := models.NewUser("test@example.com", "password")
	db.Create(&testUser)
	db.Create(&models.Client{
		ID:             "my-client",
		Secret:         []byte(`$2a$10$IxMdI6d.LIRZPpSfEwNoeu4rY3FhDREsxFJXikcgdRRAStxUlsuEO`), // = "foobar"
		RotatedSecrets: `$2y$10$X51gLxUQJ.hGw1epgHTE5u0bt64xM0COU7K9iAp.OFg8p2pUd.1zC `,
		RedirectURIs:   "http://localhost:3846/callback",
		ResponseTypes:  "id_token,code,token,id_token token,code id_token,code token,code id_token token",
		GrantTypes:     "refresh_token,authorization_code,client_credentials",
		Scopes:         "openid,offline",
	})
}

func Connect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("idp.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}
