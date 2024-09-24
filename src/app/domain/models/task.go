package models

import (
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ID int `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	UserID   int `json:"user_id"`
}

func (task *Task) GetID() int {
	return task.ID
}