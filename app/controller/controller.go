package controller

import (
	"errors"
	"tasks/app/domain"
	"tasks/app/domain/models"
	"tasks/app/infra/repository"
)

type Controller struct {
	Repo   *repository.Repository
	Domain *domain.Domain
}

func (controller *Controller) CreateUser(email string, password string) (id int, err error) {
	user := controller.Domain.CreateUser(email, password)
	id, err := controller.Repo.SaveModel(user)
	return id, err
}

func (controller *Controller) CreateSession(email string) (id int, err error) {
	session := controller.Domain.CreateSession(email)
	id, err := controller.Repo.SaveModel(session)
	return id, err
}

func (controller *Controller) CreateTask(name string, email string) (id int, err error) {
	user, err := controller.Repo.GetUser(email)
	if err == nil {
		session, err := controller.Repo.GetSession(email)
		if err == nil {
			task := controller.Domain.CreateTask(name, user)
			id, err := controller.Repo.SaveModel(task)
			return id, err
		}
	}
	return 0, err
}

func (controller *Controller) DeleteTask(id int, email string) (err error) {
	task, err := controller.Repo.GetTask(id)
	if err == nil {
		session, err := controller.Repo.GetSession(email)
		if err == nil {
			result := controller.Repo.DeleteNote(task, id)
			return nil
		}
	}
	return result.Error
}