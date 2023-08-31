package v1

import (
	_ "avito-rest-api/docs"
	"avito-rest-api/internal/service"
	customError "avito-rest-api/package/error"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

func NewRouter(handler *echo.Echo, services *service.Services) {
	handler.GET("/health", func(c echo.Context) error { return c.NoContent(200) })
	handler.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := handler.Group("/api/v1")

	newUserRoutes(v1.Group("/users"), services.User)
	newSegmentRoutes(v1.Group("/segments"), services.Segment)
}

func errorHandler(c echo.Context, err error) error {
	switch t := err.(type) {
	// Ошибки пользователя
	case customError.ErrUserValidationError:
		t.Title = "ErrUserValidationError"
		return c.JSON(http.StatusBadRequest, t)
	case customError.ErrUserNotFound:
		t.Title = "ErrUserNotFound"
		return c.JSON(http.StatusNotFound, t)
	case customError.ErrUserDeleted:
		t.Title = "ErrUserDeleted"
		return c.JSON(http.StatusNotFound, t)

	// Ошибки сегмента
	case customError.ErrSegmentValidationError:
		t.Title = "ErrSegmentValidationError"
		return c.JSON(http.StatusBadRequest, t)
	case customError.ErrSegmentNotFound:
		t.Title = "ErrSegmentNotFound"
		return c.JSON(http.StatusNotFound, t)
	case customError.ErrSegmentDeleted:
		t.Title = "ErrSegmentDeleted"
		return c.JSON(http.StatusNotFound, t)
	case customError.ErrSegmentAlreadyExists:
		t.Title = "ErrSegmentAlreadyExists"
		return c.JSON(http.StatusBadRequest, t)

	// Внутренняя ошибка сервера
	case customError.ErrInternalServerError:
		t.Title = "ErrInternalServerError"
		return c.JSON(http.StatusInternalServerError, t)

	default:
		return c.JSON(http.StatusInternalServerError, err)
	}
}
