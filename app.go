package main

import (
	authentication_controller "github.com/asishshaji/dotedon-api/controller/authentication"
	"github.com/labstack/echo/v4"
)

type App struct {
	app  *echo.Echo
	port string
}

func NewApp(port string, userController authentication_controller.IAuthenticationControllerInterface) *App {
	e := echo.New()

	uG := e.Group("/user")
	uG.POST("", userController.RegisterUser)

	return &App{
		app:  e,
		port: port,
	}
}

func (a *App) RunServer() {
	a.app.Logger.Fatal(a.app.Start(a.port))
}
