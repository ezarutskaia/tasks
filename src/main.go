package main

import (
	"fmt"
	"tasks/src/app"
	"tasks/src/app/domain"
	//"tasks/src/app/infra"
	"tasks/src/app/interfaces"
)

func main() {
	fmt.Println("Initialize app.")

	app := &app.App{
		Domain: &domain.Domain{},
		Interfaces: &interfaces.Interfaces{},
	}

	app.Interfaces.HttpServer.HandleHttpRequest(app.Domain)
}