package interfaces

import (
	"fmt"
	"net/http"
	"tasks/src/app/controller"
	"github.com/labstack/echo/v4"
	"tasks/src/app/domain/models"
)

type HttpServer struct{}

type Options struct {
    Message  string
    Data   map[string]int
}

func (server HttpServer) Response (c echo.Context, options Options) (error) {
	return c.JSON(http.StatusOK, map[string]interface{}{
        "message": options.Message,
        "data":    options.Data,
    })
}

func (server HttpServer) HandleHttpRequest(controller *controller.Controller) {
	fmt.Println("HTTP server have started.")
	e := echo.New()
	
	/*
	*	Create user
	*/
	
	e.POST("/user/add", func(c echo.Context) (err error) {
		user := new(models.User)
		if err := c.Bind(user); err != nil {
			return server.Response(c, Options{
				Message: "data reading error",
			})
		  }
		id, err := controller.CreateUser(user.Email, user.Password)
		if err != nil {
			return server.Response(c, Options{
				Message: "data recording error",
			})
		  }
		return server.Response(c, Options{
			Data:    map[string]int{"id": id},
		})
	})

	e.Logger.Fatal(e.Start(":1323"))
}