package domain

import (
	"tasks/app/domain/models"
	"time"
)

type Domain struct {}

func (domain *Domain) CreateUser(email string, password string) (*models.User){
	return &models.User{Email: email, Password: password}
}

func (domain *Domain) CreateSession(email string) (*models.Session) {
	now := time.Now()
	return &models.Session{Email: email, Endsession: now.Add(time.Hour)}
}

func (domain *Domain) CreateTask(name string, user *models.User) (*models.Task) {
	return &models.Task{Name: name, UserID: user.ID}
}

func (domain *Domain) CreateToken(email string) (*models.Token) {
	token := new(models.Token)
	token.GenerateToken(email)
	return token
}