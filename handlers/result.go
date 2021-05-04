package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) GetResults(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}

func (h Handler) AddResults(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}

func (h Handler) DeleteResults(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}
