package interfaces

import (
	"fmt"
	"strconv"
	"net/http"
	"tasks/src/app/controller"
	"github.com/labstack/echo/v4"
	"tasks/src/app/domain/models"
)

type HttpServer struct{}

type RequestBody struct {
	Name string `json:"name"`
}

type Options struct {
    Message  string
    Data   map[string]interface{}
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
	
	//	Create user

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
			Data:    map[string]interface{}{"id": id},
		})
	})

	// Login user

	e.POST("/login", func(c echo.Context) (err error) {
		user := new(models.User)
		if err := c.Bind(user); err != nil {
			return server.Response(c, Options{
				Message: "data reading error",
			})
		}

		userDB, err := controller.Repo.GetUser(user.Email)
		if err != nil {
			return server.Response(c, Options{
				Message: "user is not exist",
			})
		}
		
		if user.Password != userDB.Password {
			return server.Response(c, Options{
				Message: "password incorrect",
			})
		}

		token := controller.CreateSession(user.Email)
		return server.Response(c, Options{
			Data:    map[string]interface{}{"token": token},
		})
	})

	// task/add

	e.POST("/task/add", func(c echo.Context) (err error) {
		email := c.Request().Header.Get("Email")

		var body RequestBody
		if err := c.Bind(&body); err != nil {
			return server.Response(c, Options{
				Message: "data reading error",
			})
		}

		user, err := controller.Repo.GetUser(email)
		if err != nil {
			return server.Response(c, Options{
				Message: "user is not exist",
			})
		}
		
		id, err := controller.CreateTask(user, body.Name)
		if err != nil {
			return server.Response(c, Options{
				Message: "data recording error",
			})
		}

		return server.Response(c, Options{
			Data:    map[string]interface{}{"id": id},
		})
	})

	// task/del

	e.DELETE("/task/:id", func(c echo.Context) (err error) {
		idParam := c.Param("id")
		id,_ := strconv.Atoi(idParam)
		err = controller.DeleteTask(id)
		if err != nil {
			return server.Response(c, Options{
				Message: "task is not exist",
			})
		}

		return server.Response(c, Options{
			Message: "task was deleted",
		})
	})

	// task/list

	e.GET("/task/list", func(c echo.Context) (err error) {
		tasks, err := controller.Repo.GetTasks()
		if err != nil {
			return server.Response(c, Options{
				Message: "task is not exist",
			})
		}
		return server.Response(c, Options{
			Data:    map[string]interface{}{"list of tasks": tasks},
		})
	})

	e.Logger.Fatal(e.Start(":1323"))
}