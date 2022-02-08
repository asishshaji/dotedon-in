package main

import (
	"log"
	"os"

	authentication_controller "github.com/asishshaji/dotedon-api/controller/authentication"
	authentication_repository "github.com/asishshaji/dotedon-api/repositories/authentication"
	authentication_service "github.com/asishshaji/dotedon-api/services/authentication"
	"github.com/asishshaji/dotedon-api/utils"
)

func main() {

	logger := log.New(os.Stdout, "dotedon-api", log.LstdFlags)

	env := utils.LoadEnv(logger)
	db := env.ConnectToDB()

	userAuthRepo := authentication_repository.NewUserAuthRepo(logger, db)
	userAuthService := authentication_service.NewAuthenticationService(logger, userAuthRepo)
	userAuthController := authentication_controller.NewUserController(logger, userAuthService)

	app := NewApp(env.ServerPort, userAuthController)
	app.RunServer()
}
