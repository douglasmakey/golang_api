package routes

import (
	"github.com/labstack/echo"

	"github.com/douglasmakey/backend_base/resources"

)

func RegisterUser(e *echo.Echo) {
	userController := new(resources.UserController)

	e.POST("auth/register", userController.Register)
	e.POST("auth/login", userController.Login)

}
