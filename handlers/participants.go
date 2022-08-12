package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) GetParticipants(c echo.Context) error {
	return getAPIError(c, http.StatusNotImplemented, "Not Implemented", errors.New("endpoint not implemented yet"))
}

func (h Handler) AddParticipants(c echo.Context) error {
	return getAPIError(c, http.StatusNotImplemented, "Not Implemented", errors.New("endpoint not implemented yet"))
}

func (h Handler) DeleteParticipants(c echo.Context) error {
	return getAPIError(c, http.StatusNotImplemented, "Not Implemented", errors.New("endpoint not implemented yet"))
}
