package interfaces

import "github.com/labstack/echo/v4"

type IController interface {
	LoadRoutes(*echo.Group)
}
