package student_controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/asishshaji/dotedon-api/models"
	student_service "github.com/asishshaji/dotedon-api/services/student"
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

	fmt.Println("Hello")

	user := models.StudentDTO{}

	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		uC.l.Println(err)
		return echo.ErrBadRequest

	}
	if err := user.Validate(); err != nil {
		uC.l.Println(err)
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: err.Error(),
		})
	}

	err := uC.studentService.RegisterStudent(c.Request().Context(), &user)

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

	if json_map["email"] == nil || json_map["password"] == nil {
		return echo.ErrBadRequest
	}

	username := fmt.Sprintf("%v", json_map["email"])
	password := fmt.Sprintf("%v", json_map["password"])

	if len(username) == 0 || len(password) == 0 {
		return echo.ErrBadRequest
	}

	res, err := uC.studentService.LoginStudent(c.Request().Context(), username, password)

	if err != nil {
		uC.l.Println(err)
		return c.JSON(http.StatusForbidden, models.Response{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)

}

func (sC StudentController) UpdateStudent(c echo.Context) error {
	sDto := models.StudentDTO{}
	if err := json.NewDecoder(c.Request().Body).Decode(&sDto); err != nil {
		sC.l.Println("Error parsing user details")
		return echo.ErrInternalServerError
	}
	err := sC.studentService.UpdateStudent(c.Request().Context(), sDto)
	if err != nil {
		sC.l.Println(err)
		return echo.ErrInternalServerError
	}
	return nil
}

func (uC StudentController) GetUser(c echo.Context) error {
	res, err := uC.studentService.GetStudent(c.Request().Context(), c.Get("student_id").(primitive.ObjectID))
	if err != nil {
		uC.l.Println(err.Error())
		return c.JSON(http.StatusNoContent, nil)
	}
	return c.JSON(http.StatusOK, res)
}
func (uC StudentController) GetMentors(c echo.Context) error {

	responseDtos, err := uC.studentService.GetMentors(c.Request().Context(), c.Get("student_id").(primitive.ObjectID))
	if err != nil {
		uC.l.Println(err)
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, responseDtos)
}

func (uC StudentController) FollowMentor(c echo.Context) error {
	// change updated time of user

	mentorId := c.FormValue("mentor_id")

	mentorObjId, err := primitive.ObjectIDFromHex(mentorId)

	if err != nil {
		uC.l.Println(err)
		return err
	}

	uC.studentService.AddMentorToStudent(c.Request().Context(), c.Get("student_id").(primitive.ObjectID), mentorObjId)
	return nil
}

func (sU StudentController) CreateTaskSubmisson(c echo.Context) error {

	tSDto := models.TaskSubmissionDTO{}
	if err := json.NewDecoder(c.Request().Body).Decode(&tSDto); err != nil {
		sU.l.Println("Error parsing submission details")
		return echo.ErrInternalServerError
	}

	err := sU.studentService.CreateTaskSubmission(c.Request().Context(), tSDto, c.Get("student_id").(primitive.ObjectID))
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}

	return c.JSON(http.StatusCreated, models.Response{
		Message: "Created submission",
	})
}

func (sU StudentController) UpdateTaskSubmission(c echo.Context) error {

	tSDto := models.TaskSubmissionDTO{}
	if err := json.NewDecoder(c.Request().Body).Decode(&tSDto); err != nil {
		sU.l.Println("Error parsing submission details")
		return echo.ErrInternalServerError
	}

	err := sU.studentService.UpdateTaskSubmission(c.Request().Context(), tSDto, c.Get("student_id").(primitive.ObjectID))

	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, models.Response{
		Message: "Submitted task",
	})
}

func (sC StudentController) GetTasks(c echo.Context) error {

	taskStudentResponseArr, err := sC.studentService.GetTasks(c.Request().Context(), c.Get("student_id").(primitive.ObjectID))
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, taskStudentResponseArr)

}

func (sC StudentController) GetData(c echo.Context) error {

	domains, err := sC.studentService.GetData(c.Request().Context())
	if err != nil {
		sC.l.Println(err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, domains)
}

func (sC StudentController) UploadFile(c echo.Context) error {
	image, _, _ := c.Request().FormFile("file")

	url, err := sC.studentService.UploadFile(c.Request().Context(), image)
	if err != nil {
		sC.l.Println(err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"url": url,
	})

}

func (sC StudentController) InsertToken(c echo.Context) error {
	tK := models.TokenDto{}

	if err := json.NewDecoder(c.Request().Body).Decode(&tK); err != nil {
		return echo.ErrBadRequest
	}

	err := sC.studentService.InsertToken(c.Request().Context(), tK, c.Get("student_id").(primitive.ObjectID))

	if err != nil {
		sC.l.Println(err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, nil)
}
