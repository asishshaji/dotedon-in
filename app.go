package main

import (
	authentication_controller "github.com/asishshaji/dotedon-api/controller/authentication"
	"github.com/asishshaji/dotedon-api/docs"

	_ "github.com/asishshaji/dotedon-api/docs"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type App struct {
	app  *echo.Echo
	port string
}

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func NewApp(port string, userController authentication_controller.IAuthenticationControllerInterface) *App {
	docs.SwaggerInfo.Title = "Dotedon API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:9091"

	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	uG := e.Group("/user")
	uG.POST("", userController.RegisterUser)
	uG.POST("/login", userController.LoginUser)

	return &App{
		app:  e,
		port: port,
	}
}

func (a *App) RunServer() {
	a.app.Logger.Fatal(a.app.Start(a.port))
}
