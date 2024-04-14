package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"idp/hash"
)

type User struct {
	gorm.Model
	UserID   string `gorm:"type:varchar(36);not null;primary_key"`
	Username string `gorm:"type:varchar(255);not null;unique"`
	Password []byte `gorm:"type:blob"`
}

func NewUser(username, password string) *User {
	uuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}

	hashedPassword, err := hash.HashPassword(password)
	if err != nil {
		panic(err)
	}

	return &User{
		UserID:   uuid.String(),
		Username: username,
		Password: hashedPassword,
	}
}

func (u *User) Authenticate(password string) error {
	return hash.Compare([]byte(password), u.Password)
}
