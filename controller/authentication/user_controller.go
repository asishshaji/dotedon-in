package authentication_controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/asishshaji/dotedon-api/models"
	authentication_service "github.com/asishshaji/dotedon-api/services/authentication"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	userService authentication_service.IAuthenticationService
	l           *log.Logger
}

func NewUserController(l *log.Logger, uS authentication_service.IAuthenticationService) IAuthenticationControllerInterface {
	return UserController{
		userService: uS,
		l:           l,
	}
}

func (uC UserController) RegisterUser(c echo.Context) error {

	userModel := models.User{}

	if err := json.NewDecoder(c.Request().Body).Decode(&userModel); err != nil {
		uC.l.Println(err)
		return echo.ErrBadRequest

	}
	err := uC.userService.RegisterUser(c.Request().Context(), &userModel)

	if err != nil {
		uC.l.Println(err)
		return echo.ErrBadRequest

	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "created user",
	})
}
