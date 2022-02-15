package main

import (
	"log"
	"os"

	student_controller "github.com/asishshaji/dotedon-api/controller/student"
	student_repository "github.com/asishshaji/dotedon-api/repositories/student"
	student_service "github.com/asishshaji/dotedon-api/services/student"
	"github.com/asishshaji/dotedon-api/utils"
)

func main() {

	logger := log.New(os.Stdout, "dotedon-api", log.LstdFlags)

	env := utils.LoadEnv(logger)
	db := env.ConnectToDB()

	studentRepo := student_repository.NewStudentAuthRepo(logger, db)
	studentService := student_service.NewStudentService(logger, studentRepo)
	studentController := student_controller.NewStudentController(logger, studentService)

	controller := Controllers{
		StudentController: studentController,
	}

	app := NewApp(env.ServerPort, controller)
	app.RunServer()
}
