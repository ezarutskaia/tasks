package repository

import (
	"gorm.io/gorm"
	"tasks/app/domain/models"
)

type Repository struct {
	DB *gorm.DB
}

func (repo *Repository) Automigrate (model interface{}) {
	repo.DB.AutoMigrate(model)
}

func (repo *Repository) SaveModel (model models.HasID) (id int, err error) {
	result := (*repo.DB).Create(model)
	return model.GetID(), result.Error
}

func (repo *Repository) GetUser (email string) (user *models.User, err error) {
	result := (*repo.DB).Where("email = ?", email).First(&user)
	return user, result.Error
}

func (repo *Repository) GetSession (email string) (session *models.Session, err error) {
	result := (*repo.DB).Where("email = ? AND endsession > NOW()", email).Last(&session)
	return session, result.Error
}

func (repo *Repository) GetTask (id int) (task *models.Task, err error) {
	result := (*repo.DB).Where("id = ?", id).First(&task)
	return task, result.Error
}

func (repo *Repository) DeleteNote (model interface{}, id string) error {
	result := (*repo.DB).Delete(model, id)
	return result.Error
}