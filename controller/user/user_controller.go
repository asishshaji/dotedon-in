package user_controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/asishshaji/dotedon-api/models"
	user_service "github.com/asishshaji/dotedon-api/services/user"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	userService user_service.IUserService
	l           *log.Logger
}

func NewUserController(l *log.Logger, uS user_service.IUserService) IUserController {
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

func (uC UserController) GetMentors(c echo.Context) error {
	responseDtos, err := uC.userService.GetMentors(c.Request().Context())
	if err != nil {
		uC.l.Println(err)
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, responseDtos)
}

func (uC UserController) AddMentorToUser(c echo.Context) error {
	// change updated time of user

	mentorId := c.FormValue("mentor_id")

	mentorObjId, err := primitive.ObjectIDFromHex(mentorId)

	if err != nil {
		uC.l.Println(err)
		return err
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.Claims)
	uC.l.Println(claims.UserId)
	uC.userService.AddMentorToUser(c.Request().Context(), claims.UserId, mentorObjId)
	return nil
}
