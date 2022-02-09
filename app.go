package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	user_controller "github.com/asishshaji/dotedon-api/controller/user"
	"github.com/labstack/echo/v4"
)

type App struct {
	app  *echo.Echo
	port string
}

func NewApp(port string, userController user_controller.IUserController) *App {
	e := echo.New()

	uG := e.Group("/user")
	uG.POST("", userController.RegisterUser)
	uG.POST("/login", userController.LoginUser)
	uG.GET("/mentors", userController.GetMentors)

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
