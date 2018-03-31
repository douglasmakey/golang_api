package routes

import "github.com/labstack/echo"

func Init(e *echo.Echo) {
	RegisterUser(e)

	// Your routes
}
