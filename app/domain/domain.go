package domain

import {
	"tasks/app/domain/models"
	"time"
}

type Domain struct {}

func (domain *Domain) CreateUser(email string, password string) (*User){
	return &User{Email: email, Password: password}
}

func (domain *Domain) CreateSession(email string) (*Session) {
	now := time.Now()
	return &Session{Email: email, Endsession: now.Add(time.Hour)}
}

func (domain *Domain) CreateTask(name string, user *User) (*Task) {
	return &Task{Name: name, UserID: user.ID}
}