package models

import (
	"gorm.io/gorm"
	"time"
)

type Session struct {
	gorm.Model
	ID    int `gorm:"primaryKey"`
	Email string 
	Endsession time.Time
}

func (session *Session) GetID() int {
	return session.ID
}