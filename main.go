package main

import (
	"log"
	"os"

	admin_controller "github.com/asishshaji/dotedon-api/controller/admin"
	student_controller "github.com/asishshaji/dotedon-api/controller/student"
	admin_repository "github.com/asishshaji/dotedon-api/repositories/admin"
	student_repository "github.com/asishshaji/dotedon-api/repositories/student"
	admin_service "github.com/asishshaji/dotedon-api/services/admin"
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

	adminRepo := admin_repository.NewAdminRepository(logger, db)
	adminService := admin_service.NewAdminService(logger, adminRepo)
	adminController := admin_controller.NewAdminController(logger, adminService)

	controller := Controllers{
		StudentController: studentController,
		AdminController:   adminController,
	}

	app := NewApp(env.ServerPort, controller)
	app.RunServer()
}
