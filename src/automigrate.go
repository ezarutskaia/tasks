package main

import (
	"tasks/src/app/domain/models"
	"tasks/src/app/infra/repository"
)

func main() {
	db := repository.SqlConnection()
	repo := &repository.Repository{DB: db,}
	repo.Automigrate(&models.User{})
}