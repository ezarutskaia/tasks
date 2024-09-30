package interfaces

import (
	"fmt"
	"sync"
	"math/rand"
	"strconv"
	"net/http"
	"context"
	"errors"
	"encoding/json"
	"tasks/src/app/infra/pdf"
	"tasks/src/app/controller"
	"github.com/labstack/echo/v4"
	"tasks/src/app/domain/models"
)

type HttpServer struct{}

type RequestBody struct {
	Name string `json:"name"`
}

type ServiceResponse struct {
    Data    string `json:"data"`
    Message string `json:"message"`
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
		fmt.Println("Request on pdf")
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

		var servicesResp []ServiceResponse
		var printErrors []string
		var wg sync.WaitGroup
		var mu sync.Mutex

		wg.Add(len(tasksDTO))

		for _,taskDTO := range(tasksDTO) {
			task := taskDTO
			go func(taskDTO *pdf.TaskDTO) {
				defer wg.Done()

				body, err := controller.PrintTask(taskDTO)
				if err != nil {
					mu.Lock()
					if errors.Is(err, context.DeadlineExceeded) {
						printErrors = append(printErrors, fmt.Sprintf("Print service timed out for task %d", taskDTO.Id))
					} else {
						printErrors = append(printErrors, fmt.Sprintf("Print service is not available for task %d", taskDTO.Id))
					}
					mu.Unlock()
					return
				}

				var serviceResp ServiceResponse
				err = json.Unmarshal(body, &serviceResp)
				if err != nil {
					fmt.Printf("4 -> %s\n", serviceResp)
					mu.Lock()
					printErrors = append(printErrors, fmt.Sprintf("Reading response error for task %d", taskDTO.Id))
					mu.Unlock()
					return
				}
				fmt.Printf("5 -> %s\n", serviceResp)
				mu.Lock()
				servicesResp = append(servicesResp, serviceResp)
				mu.Unlock()

			} (task)
		}

		wg.Wait()

		if len(printErrors) > 0 {
			return server.Response(c, Options{
				Message: "Some tasks failed",
				Data:    map[string]interface{}{"count": len(printErrors)},
			})
		}
	
		return server.Response(c, Options{
			Message: "All tasks were printed successfully",
			Data:    map[string]interface{}{"count": len(servicesResp)},
		})
	})

	e.Logger.Fatal(e.Start(":1323"))
}