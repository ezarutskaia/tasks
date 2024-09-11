package main

import (
	"fmt"
	"tasks/app"
	"tasks/app/domain"
	//"tasks/app/infra"
	"tasks/app/interfaces"
)

func main() {
	fmt.Println("Initialize app.")

	app := &app.App{
		Domain: &domain.Domain{},
		Interfaces: &interfaces.Interfaces{},
	}

	app.Interfaces.HttpServer.HandleHttpRequest(app.Domain)
}