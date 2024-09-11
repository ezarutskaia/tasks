package models

import (
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ID int `gorm:"primaryKey"`
	Name     string
	UserID   int
}

func (task *Task) GetID() int {
	return task.ID
}