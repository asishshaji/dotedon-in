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

	task := models.Task{}

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

	taskId := c.FormValue("task_id")
	taskObjId, _ := primitive.ObjectIDFromHex(taskId)

	task := models.TaskUpdate{}

	if err := json.NewDecoder(c.Request().Body).Decode(&task); err != nil {
		aC.l.Println(err)
		return echo.ErrBadRequest

	}

	task.Id = taskObjId

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
