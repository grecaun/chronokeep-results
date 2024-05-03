package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) GetCertificate(c echo.Context) error {
	name := c.Param("name")
	event := c.Param("event")
	time := c.Param("time")
	date := c.Param("date")
	img, err := CreateCertificate(name, event, time, date)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.Blob(http.StatusOK, "image/png", img)
}
