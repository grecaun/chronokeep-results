package handlers

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func (h Handler) GetCertificate(c echo.Context) error {
	name := c.Param("name")
	event := c.Param("event")
	time := c.Param("time")
	date := c.Param("date")
	log.WithFields(log.Fields{
		"name":  name,
		"event": event,
		"time":  time,
		"date":  date,
	}).Debug("Creating certificate.")
	name, err := url.QueryUnescape(name)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	event, err = url.QueryUnescape(event)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	time, err = url.QueryUnescape(time)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	date, err = url.QueryUnescape(date)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	img, err := CreateCertificate(name, event, time, date)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.Blob(http.StatusOK, "image/png", img)
}
