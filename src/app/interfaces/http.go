package interfaces

import (
	"fmt"
	"math/rand"
	"strconv"
	"net/http"
	"tasks/src/app/infra/pdf"
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

	taskGroup := e.Group("/task")

	taskGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			email := c.Request().Header.Get("Email")
			tokenString := c.Request().Header.Get("Token")
			if email == "" || tokenString == "" {
				return server.Response(c, Options{
					Message: "no token or email",
				})
			}

			token := &models.Token{Value: tokenString}

			_, err := controller.ValidationSession(token, email)
			if err != nil {
				return server.Response(c, Options{
					Message: "invalid token or email",
				})
			}

			return next(c)
		}
	})
	
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

	taskGroup.POST("/add", func(c echo.Context) (err error) {
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

	taskGroup.DELETE("/delete/:id", func(c echo.Context) (err error) {
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

	taskGroup.GET("/list", func(c echo.Context) (err error) {
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

	// task/pdf

	taskGroup.POST("/pdf", func(c echo.Context) (err error) {
		var ids []int
		if err := c.Bind(&ids); err != nil {
			return server.Response(c, Options{
				Message: "data reading error",
			})
		}
		
		tasks, err := controller.Repo.GetTasksByIds(ids)
		if err != nil {
			return server.Response(c, Options{
				Message: "task is not exist",
			})
		}

		var tasksDTO []*pdf.TaskDTO
		var weight int

		for _,task := range(tasks) {
			user, err := controller.Repo.GetUserById(task.UserID)
			if err != nil {
				return server.Response(c, Options{
					Message: "user is not exist",
				})
			}
			weight = rand.Intn(10000)
			tasksDTO = append(tasksDTO, &pdf.TaskDTO{
				Id: task.ID, 
				Title: task.Name, 
				User: user.Email, 
				Weight: weight})
		}

		err = controller.PrintTasks(tasksDTO)
		if err != nil {
			return server.Response(c, Options{
				Message: "error of print of tasks",
			})
		}

		return server.Response(c, Options{
			Data:    map[string]interface{}{"list of tasks": tasks},
		})

	})

	e.Logger.Fatal(e.Start(":1323"))
}