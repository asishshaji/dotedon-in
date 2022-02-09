package main

import (
	"log"
	"os"

	user_controller "github.com/asishshaji/dotedon-api/controller/user"
	user_repository "github.com/asishshaji/dotedon-api/repositories/user"
	user_service "github.com/asishshaji/dotedon-api/services/user"
	"github.com/asishshaji/dotedon-api/utils"
)

func main() {

	logger := log.New(os.Stdout, "dotedon-api", log.LstdFlags)

	env := utils.LoadEnv(logger)
	db := env.ConnectToDB()

	userRepo := user_repository.NewUserAuthRepo(logger, db)
	userService := user_service.NewUserService(logger, userRepo)
	userController := user_controller.NewUserController(logger, userService)

	app := NewApp(env.ServerPort, userController)
	app.RunServer()
}
