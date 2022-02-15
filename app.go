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

	uG := e.Group("/user")
	uG.POST("", controller.StudentController.RegisterStudent)
	uG.POST("/login", controller.StudentController.LoginStudent)

	config := middleware.JWTConfig{
		Claims:     &models.Claims{},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}

	r := e.Group("/restricted")

	r.Use(middleware.JWTWithConfig(config))

	r.GET("/mentors", controller.StudentController.GetMentors)
	r.POST("/mentors", controller.StudentController.AddMentorToStudent)

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
