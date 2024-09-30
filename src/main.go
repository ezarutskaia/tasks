package main

import (
	"fmt"
	"tasks/src/app"
	"tasks/src/app/domain"
	"tasks/src/app/infra"
	"tasks/src/app/infra/pdf"
	"tasks/src/app/controller"
	"tasks/src/app/interfaces"
	"tasks/src/app/infra/repository"
)

func main() {
	fmt.Println("Initialize app.")
	db := repository.SqlConnection()
	url := "http://localhost:8050/pdf"

	app := &app.App{
		Domain: &domain.Domain{},
		Infra: &infra.Infra{
			Repository: &repository.Repository{
				DB: db,
			},
			Pdf: &pdf.Pdf{
				URL: url,
			},
		},
		Interfaces: &interfaces.Interfaces{},
	}

	app.Interfaces.HttpServer.HandleHttpRequest(&controller.Controller{
		Repo: app.Infra.Repository,
		Domain: app.Domain,
		Pdf: app.Infra.Pdf,
	})
}