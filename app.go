package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	student_controller "github.com/asishshaji/dotedon-api/controller"
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
}

func NewApp(port string, controller Controllers) *App {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(middleware.Secure())

	e.POST("", controller.StudentController.RegisterStudent)
	e.POST("/login", controller.StudentController.LoginStudent)

	studentJwtConfig := middleware.JWTConfig{
		Claims:     &models.StudentJWTClaims{},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}

	r := e.Group("/restricted")

	r.Use(middleware.JWTWithConfig(studentJwtConfig))
	r.Use(utils.StudentAuthenticationMiddleware)

	r.GET("/mentors", controller.StudentController.GetMentors)
	r.POST("/mentors", controller.StudentController.FollowMentor)
	r.POST("/task/submit", controller.StudentController.CreateTaskSubmisson)
	r.PUT("/task/submit", controller.StudentController.UpdateTaskSubmission)
	r.GET("/task", controller.StudentController.GetTasks)

	return &App{
		app:  e,
		port: port,
	}
}

func (a *App) RunServer() {

	go func() {
		a.app.Logger.Fatal(a.app.Start(a.port))
	}()

	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sigData := <-sigChan

	log.Printf("Signal received : %v\n", sigData)
	tc, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	a.app.Shutdown(tc)
}
