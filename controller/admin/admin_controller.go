package admin_controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/asishshaji/dotedon-api/models"
	admin_service "github.com/asishshaji/dotedon-api/services/admin"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminController struct {
	l            *log.Logger
	adminService admin_service.IAdminService
}

func NewAdminController(l *log.Logger, adminService admin_service.IAdminService) IAdminController {
	return AdminController{
		l:            l,
		adminService: adminService,
	}
}

func (aC AdminController) Login(c echo.Context) error {
	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)
	if err != nil {
		aC.l.Println(err)
		return echo.ErrInternalServerError
	}
	username := fmt.Sprintf("%v", json_map["username"])
	password := fmt.Sprintf("%v", json_map["password"])

	if len(username) == 0 || len(password) == 0 {
		return echo.ErrBadRequest
	}

	token, err := aC.adminService.Login(c.Request().Context(), username, password)
	if err != nil {
		aC.l.Println(err)
		return c.JSON(http.StatusForbidden, models.Response{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: token,
	})
}

func (aC AdminController) AddTask(c echo.Context) error {

	adminId := c.Get("admin_id").(primitive.ObjectID)

	task := models.TaskDTO{}

	if err := json.NewDecoder(c.Request().Body).Decode(&task); err != nil {
		aC.l.Println(err)
		return echo.ErrBadRequest

	}

	err := aC.adminService.AddTask(c.Request().Context(), task, adminId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Error creating task",
		})
	}

	return c.JSON(http.StatusCreated, models.Response{
		Message: "Created task",
	})
}

func (aC AdminController) UpdateTask(c echo.Context) error {

	task := models.TaskDTO{}

	if err := json.NewDecoder(c.Request().Body).Decode(&task); err != nil {
		aC.l.Println(err)
		return echo.ErrBadRequest

	}

	err := task.Validate()
	if err != nil {
		aC.l.Println(err)
		return echo.ErrInternalServerError
	}

	err = aC.adminService.UpdateTask(c.Request().Context(), task)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Error updating task",
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: "Updated task",
	})
}

func (aC AdminController) GetTasks(c echo.Context) error {
	tasks, err := aC.adminService.GetTasks(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Error getting tasks",
		})
	}

	return c.JSON(http.StatusOK, tasks)
}

func (aC AdminController) GetUsers(c echo.Context) error {
	students, err := aC.adminService.GetUsers(c.Request().Context())

	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, students)
}

func (aC AdminController) CreateType(c echo.Context) error {
	return nil
}

func (aC AdminController) DeleteTask(c echo.Context) error {
	id := c.FormValue("task_id")
	taskId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		aC.l.Println("Error parsing task id")
		return echo.ErrInternalServerError
	}

	err = aC.adminService.DeleteTask(c.Request().Context(), taskId)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusAccepted, models.Response{
		Message: "deleted task",
	})
}

func (aC AdminController) GetTaskSubmissions(c echo.Context) error {
	res, _ := aC.adminService.GetTaskSubmissions(c.Request().Context())
	return c.JSON(http.StatusOK, res)
}

func (aC AdminController) EditTaskSubmissionStatus(c echo.Context) error {
	// get task id, check if task submission exists
	// get task status, update task submission
	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)
	if err != nil {
		aC.l.Println(err)
		return echo.ErrInternalServerError
	}
	statusString := fmt.Sprintf("%v", json_map["status"])

	taskId := fmt.Sprintf("%v", json_map["task_id"])

	if statusString == "" {
		aC.l.Println("Error parsing status")
		return echo.ErrInternalServerError
	}

	aC.l.Println(statusString, taskId)

	// status, err = primitive.ObjectIDFromHex(models.Status(status).String())
	taskIdObj, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		aC.l.Println(err)
		return echo.ErrInternalServerError
	}

	// if !primitive.IsValidObjectID(taskIdObj.String()) {
	// 	aC.l.Println("Invalid task id")
	// 	return c.JSON(http.StatusInternalServerError, models.Response{
	// 		Message: "invalid task id",
	// 	})
	// }

	aC.l.Println(taskIdObj)

	err = aC.adminService.EditTaskSubmission(c.Request().Context(), taskIdObj, models.Status(statusString))
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusAccepted, models.Response{
		Message: "edited task submission",
	})
}

func (aC AdminController) GetTaskSubmissionForUser(c echo.Context) error {
	userId := c.Param("id")
	userIdObj, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		aC.l.Println("Error parsing user id")
		return echo.ErrInternalServerError
	}
	tasks, err := aC.adminService.GetTaskSubmissionsForUser(c.Request().Context(), userIdObj)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, tasks)
}

func (aC AdminController) CreateMentor(c echo.Context) error {
	mentor := new(models.MentorDTO)

	if err := json.NewDecoder(c.Request().Body).Decode(mentor); err != nil {
		aC.l.Println("Error parsing mentor body")
		return echo.ErrInternalServerError
	}

	if err := mentor.Validate(); err != nil {
		aC.l.Println(err)
		return echo.ErrInternalServerError
	}

	err := aC.adminService.CreateMentor(c.Request().Context(), *mentor)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, models.Response{
		Message: "created mentor",
	})
}

func (aC AdminController) UpdateMentor(c echo.Context) error {

	mentor := models.MentorDTO{}

	if err := json.NewDecoder(c.Request().Body).Decode(&mentor); err != nil {
		aC.l.Println("Error parsing mentor")
		return echo.ErrInternalServerError
	}

	if err := mentor.Validate(); err != nil {
		aC.l.Println(err)
		return echo.ErrInternalServerError
	}

	err := aC.adminService.UpdateMentor(c.Request().Context(), mentor)
	if err != nil {
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusAccepted, models.Response{
		Message: "updated mentor",
	})
}

func (aC AdminController) GetMentors(c echo.Context) error {
	mentors, err := aC.adminService.GetMentors(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, mentors)
}
