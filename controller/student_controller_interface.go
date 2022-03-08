package student_controller

import "github.com/labstack/echo/v4"

type IStudentController interface {
	RegisterStudent(c echo.Context) error
	LoginStudent(c echo.Context) error
	GetMentors(c echo.Context) error
	FollowMentor(c echo.Context) error

	UpdateTaskSubmission(c echo.Context) error //separate for create and upate
	CreateTaskSubmisson(c echo.Context) error
	GetTasks(c echo.Context) error
}
