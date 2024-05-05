package models

import (
	"math/rand"
	"time"
)

type LoginSkipSession struct {
	SessionID uint64 `gorm:"primaryKey;autoIncrement:true;"`
	Token     string `gorm:"type:varchar(255);not null;unique;index"`
	UserID    string
	ExpiresAt time.Time
}

func NewLoginSkipSession(userId string) *LoginSkipSession {
	return &LoginSkipSession{
		Token:     generateToken(),
		UserID:    userId,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
}

func generateToken() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
