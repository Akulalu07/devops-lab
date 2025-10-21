package internal

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HelloHandler(c echo.Context) error {
	totalRequests.Inc()
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Hello from Echo microservice!",
	})
}
