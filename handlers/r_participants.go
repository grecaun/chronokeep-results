package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) RGetParticipants(c echo.Context) error {
	return getAPIError(c, http.StatusNotImplemented, "Not Implemented", errors.New("endpoint not implemented yet"))
}

func (h Handler) RAddParticipants(c echo.Context) error {
	return getAPIError(c, http.StatusNotImplemented, "Not Implemented", errors.New("endpoint not implemented yet"))
}

func (h Handler) RDeleteParticipants(c echo.Context) error {
	return getAPIError(c, http.StatusNotImplemented, "Not Implemented", errors.New("endpoint not implemented yet"))
}
