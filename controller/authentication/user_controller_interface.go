package authentication_controller

import "github.com/labstack/echo/v4"

type IAuthenticationControllerInterface interface {
	RegisterUser(c echo.Context) error
}
