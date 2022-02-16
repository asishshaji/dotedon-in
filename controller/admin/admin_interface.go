package admin_controller

import "github.com/labstack/echo/v4"

type IAdminController interface {
	Login(c echo.Context) error
	AddTask(c echo.Context) error
	UpdateTask(c echo.Context) error
	GetTasks(c echo.Context) error
}
