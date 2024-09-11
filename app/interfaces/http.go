package interfaces

import (
	"fmt"
	"net/http"
	"tasks/app/domain"
	"github.com/labstack/echo/v4"
)

type HttpServer struct{}

func (server HttpServer) HandleHttpRequest(domain *domain.Domain) {
	fmt.Println("HTTP server have started.")
	e := echo.New()
	e.POST("/users", func(c echo.Context) error {
		email := c.QueryParam("email")
		password := c.QueryParam("password")
		user := domain.CreateUser(email, password)
		return c.String(http.StatusOK, fmt.Sprintf("%+v", user))
	})
	e.Logger.Fatal(e.Start(":1323"))
}