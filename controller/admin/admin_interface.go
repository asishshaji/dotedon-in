package admin_controller

import "github.com/labstack/echo/v4"

type IAdminController interface {
	Login(c echo.Context) error
	AddTask(c echo.Context) error
	UpdateTask(c echo.Context) error
	GetTasks(c echo.Context) error
	DeleteTask(c echo.Context) error
	GetUsers(c echo.Context) error
	CreateType(c echo.Context) error
	GetTaskSubmissions(c echo.Context) error
}
