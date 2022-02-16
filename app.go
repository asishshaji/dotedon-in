package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	admin_controller "github.com/asishshaji/dotedon-api/controller/admin"
	student_controller "github.com/asishshaji/dotedon-api/controller/student"
	"github.com/asishshaji/dotedon-api/models"
	"github.com/asishshaji/dotedon-api/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	app  *echo.Echo
	port string
}

type Controllers struct {
	StudentController student_controller.IStudentController
	AdminController   admin_controller.IAdminController
}

func NewApp(port string, controller Controllers) *App {
	e := echo.New()

	uG := e.Group("/student")
	uG.POST("", controller.StudentController.RegisterStudent)
	uG.POST("/login", controller.StudentController.LoginStudent)

	studentJwtConfig := middleware.JWTConfig{
		Claims:     &models.StudentJWTClaims{},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}

	r := e.Group("/restricted")

	r.Use(middleware.JWTWithConfig(studentJwtConfig))
	r.Use(utils.StudentAuthenticationMiddleware)

	r.GET("/mentors", controller.StudentController.GetMentors)
	r.POST("/mentors", controller.StudentController.AddMentorToStudent)
	r.POST("/task/submit", controller.StudentController.TaskSubmission)
	r.GET("/task", controller.StudentController.GetTasks)

	adminGroup := e.Group("/admin")
	adminGroup.POST("/login", controller.AdminController.Login)

	adminGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &models.AdminJWTClaims{},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))

	adminGroup.Use(utils.AdminAuthenticationMiddleware)

	adminGroup.POST("/task", controller.AdminController.AddTask)
	adminGroup.PUT("/task", controller.AdminController.UpdateTask)
	adminGroup.GET("/task", controller.AdminController.GetTasks)

	return &App{
		app:  e,
		port: port,
	}
}

func (a *App) RunServer() {

	go func() {
		a.app.Logger.Fatal(a.app.Start(a.port))
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sigData := <-sigChan

	log.Printf("Signal received : %v\n", sigData)
	tc, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	a.app.Shutdown(tc)
}
