package student_controller

import "github.com/labstack/echo/v4"

type IStudentController interface {
	RegisterStudent(c echo.Context) error
	LoginStudent(c echo.Context) error
	GetMentors(c echo.Context) error
	AddMentorToStudent(c echo.Context) error
	TaskSubmission(c echo.Context) error
}
