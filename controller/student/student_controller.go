package student_controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/asishshaji/dotedon-api/models"
	student_service "github.com/asishshaji/dotedon-api/services/student"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentController struct {
	studentService student_service.IStudentService
	l              *log.Logger
}

func NewStudentController(l *log.Logger, uS student_service.IStudentService) IStudentController {
	return StudentController{
		studentService: uS,
		l:              l,
	}
}

func (uC StudentController) RegisterStudent(c echo.Context) error {

	userModel := models.Student{}

	if err := json.NewDecoder(c.Request().Body).Decode(&userModel); err != nil {
		uC.l.Println(err)
		return echo.ErrBadRequest

	}

	userModel.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	userModel.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	err := uC.studentService.RegisterStudent(c.Request().Context(), &userModel)

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

func (uC StudentController) LoginStudent(c echo.Context) error {

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

	token, err := uC.studentService.LoginStudent(c.Request().Context(), username, password)

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

func (uC StudentController) GetMentors(c echo.Context) error {
	responseDtos, err := uC.studentService.GetMentors(c.Request().Context())
	if err != nil {
		uC.l.Println(err)
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, responseDtos)
}

func (uC StudentController) AddMentorToStudent(c echo.Context) error {
	// change updated time of user

	mentorId := c.FormValue("mentor_id")

	mentorObjId, err := primitive.ObjectIDFromHex(mentorId)

	if err != nil {
		uC.l.Println(err)
		return err
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.Claims)
	uC.l.Println(claims.StudentId)
	uC.studentService.AddMentorToStudent(c.Request().Context(), claims.StudentId, mentorObjId)
	return nil
}
