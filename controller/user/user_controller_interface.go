package user_controller

import "github.com/labstack/echo/v4"

type IUserController interface {
	RegisterUser(c echo.Context) error
	LoginUser(c echo.Context) error
	GetMentors(c echo.Context) error
	AddMentorToUser(c echo.Context) error
}
