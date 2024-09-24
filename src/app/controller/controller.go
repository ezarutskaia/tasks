package controller

import (
	"bytes"
    "encoding/json"
    "net/http"
	"tasks/src/app/domain"
	"tasks/src/app/infra/pdf"
	"tasks/src/app/domain/models"
	"tasks/src/app/infra/repository"
)

type Controller struct {
	Repo   *repository.Repository
	Domain *domain.Domain
}

func (controller *Controller) CreateUser(email string, password string) (id int, err error) {
	user := controller.Domain.CreateUser(email, password)
	id, err = controller.Repo.SaveModel(user)
	return id, err
}

func (controller *Controller) CreateSession(email string) string {
	_, err := controller.Repo.GetUser(email)
	if err == nil {
		session := controller.Domain.CreateSession(email)
		controller.Repo.SaveModel(session)
		token := controller.Domain.CreateToken(email)
		return token.Value
	}
	return ""
}

func (controller *Controller) ValidationSession(token *models.Token, email string) (user *models.User, err error) {
	_, err = controller.Repo.GetSession(email)
		if err == nil {
			if token.ValidToken(email) {
				user, err := controller.Repo.GetUser(email)
				return user, err
			}
		}
		return &models.User{}, err
}

func (controller *Controller) CreateTask(user *models.User, name string) (id int, err error) {
	task := controller.Domain.CreateTask(name, user)
	id, err = controller.Repo.SaveModel(task)
	return id, err
}

func (controller *Controller) DeleteTask(id int) (err error) {
	task, err := controller.Repo.GetTask(id)
	if err == nil {
		result := controller.Repo.DeleteNote(task, id)
		return result
		}
	return err
}

func (controller *Controller) PrintTasks(tasks []*pdf.TaskDTO) (err error) {
	url := "http://127.0.0.1:8050/pdf"

	taskJson, err := json.Marshal(tasks)
    if err != nil {
        return err
    }

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(taskJson))
    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

	return nil
}