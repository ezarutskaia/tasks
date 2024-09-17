package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID    int  `gorm:"primaryKey"`
	Email string `gorm:"unique"`
	Password string
	Tasks []Task
}

func (user *User) GetID() int {
	return user.ID
}