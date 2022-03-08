package main

import (
	"context"
	"log"
	"os"

	student_controller "github.com/asishshaji/dotedon-api/controller"
	student_repository "github.com/asishshaji/dotedon-api/repositories"
	image_service "github.com/asishshaji/dotedon-api/services/image"
	student_service "github.com/asishshaji/dotedon-api/services/student"
	"github.com/asishshaji/dotedon-api/utils"
	"github.com/go-redis/redis/v8"
)

func main() {

	logger := log.New(os.Stdout, "dotedon-api", log.LstdFlags)

	env := utils.LoadEnv(logger)
	db := env.ConnectToDB()

	// logger.Println("Creating indices for mongodb collections")
	// utils.CreateIndex(db, "mentor", "name", true)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := redisClient.Ping(context.TODO()).Result()
	if err != nil {
		// logger.Fatalln(err)
		logger.Println(err)
	} else {
		logger.Println("Connected to redis")
	}

	imageService := image_service.NewImageService(logger)

	studentRepo := student_repository.NewStudentAuthRepo(logger, db)
	studentService := student_service.NewStudentService(logger, studentRepo, redisClient, imageService)
	studentController := student_controller.NewStudentController(logger, studentService)

	controller := Controllers{
		StudentController: studentController,
	}

	app := NewApp(env.ServerPort, controller)
	app.RunServer()
}
