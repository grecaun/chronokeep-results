package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) GetAccount(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}

func (h Handler) GetAccounts(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}

func (h Handler) AddAccount(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}

func (h Handler) UpdateAccount(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}

func (h Handler) DeleteAccount(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}
