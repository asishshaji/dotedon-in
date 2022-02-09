package authentication_controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/asishshaji/dotedon-api/models"
	authentication_service "github.com/asishshaji/dotedon-api/services/authentication"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// RegisterUser godoc
// @Summary  Registers user
// @Accept   json
// @Product  json
// @Param    user  body      models.User  true  "User body"
// @Success  200   {object}  map[string]string
// @Router   /user [post]
func (uC UserController) RegisterUser(c echo.Context) error {

	userModel := models.User{}

	if err := json.NewDecoder(c.Request().Body).Decode(&userModel); err != nil {
		uC.l.Println(err)
		return echo.ErrBadRequest

	}

	userModel.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	userModel.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	err := uC.userService.RegisterUser(c.Request().Context(), &userModel)

	if err != nil {
		uC.l.Println(err)
		return c.JSON(http.StatusNotAcceptable, models.Response{
			Message: err.Error(),
		})

	}

	return c.JSON(http.StatusCreated, models.Response{
		Message: "created user",
	})
}

func (uC UserController) LoginUser(c echo.Context) error {

	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)
	if err != nil {
		uC.l.Println(err)
		return echo.ErrInternalServerError
	}

	username := fmt.Sprintf("%v", json_map["username"])
	password := fmt.Sprintf("%v", json_map["password"])

	if len(username) == 0 || len(password) == 0 {
		return echo.ErrBadRequest
	}

	token, err := uC.userService.LoginUser(c.Request().Context(), username, password)

	if err != nil {
		uC.l.Println(err)
		return c.JSON(http.StatusForbidden, models.Response{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: token,
	})

}
