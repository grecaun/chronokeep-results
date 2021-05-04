package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) GetEvents(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}

func (h Handler) GetEvent(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}

func (h Handler) AddEvent(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}

func (h Handler) UpdateEvent(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}

func (h Handler) DeleteEvent(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}
