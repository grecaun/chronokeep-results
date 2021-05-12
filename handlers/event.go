package handlers

import (
	"chronokeep/results/database"
	"chronokeep/results/types"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func (h Handler) GetEvents(c echo.Context) error {
	log.Info(fmt.Sprintf("Host: %v", c.Request().Host))
	log.Info(fmt.Sprintf("Referer: %v", c.Request().Referer()))
	log.Info(fmt.Sprintf("RemoteAddr: %v", c.Request().RemoteAddr))
	var request types.GenderalRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	// Get Key :: TODO :: Add verification of HOST value.
	key, err := database.GetKey(request.Key)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if key == nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	events, err := database.GetEvents()
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	return c.JSON(http.StatusOK, types.GetEventsResponse{
		Events: events,
	})
}

func (h Handler) GetEvent(c echo.Context) error {
	log.Info(fmt.Sprintf("Host: %v", c.Request().Host))
	log.Info(fmt.Sprintf("Referer: %v", c.Request().Referer()))
	log.Info(fmt.Sprintf("RemoteAddr: %v", c.Request().RemoteAddr))
	var request types.GetEventRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	// Get Key :: TODO :: Add verification of HOST value.
	key, err := database.GetKey(request.Key)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if key == nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	account, err := database.GetAccountByID(key.AccountIdentifier)
	if err != nil || account == nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	event, err := database.GetEvent(request.EventSlug)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if event == nil {
		return getAPIError(c, http.StatusNotFound, "Event Not Found", nil)
	}
	if account.Type != "admin" && event.AccessRestricted && account.Identifier != event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	eventYears, err := database.GetEventYears(event.Slug)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	return c.JSON(http.StatusOK, types.GetEventResponse{
		Event:      *event,
		EventYears: eventYears,
	})
}

func (h Handler) AddEvent(c echo.Context) error {
	log.Info(fmt.Sprintf("Host: %v", c.Request().Host))
	log.Info(fmt.Sprintf("Referer: %v", c.Request().Referer()))
	log.Info(fmt.Sprintf("RemoteAddr: %v", c.Request().RemoteAddr))
	var request types.AddEventRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	// Get Key :: TODO :: Add verification of HOST value.
	key, err := database.GetKey(request.Key)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if key == nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	// Verify key access level.  Readonly cannot write or modify values.
	if key.Type == "read" {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	account, err := database.GetAccountByID(key.AccountIdentifier)
	if err != nil || account == nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if account.Type != "admin" && request.AccountEmail != account.Email {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	event, err := database.AddEvent(types.Event{
		AccountIdentifier: account.Identifier,
		Name:              request.Event.Name,
		Slug:              request.Event.Slug,
		Website:           request.Event.Website,
		Image:             request.Event.Image,
		ContactEmail:      request.Event.ContactEmail,
		AccessRestricted:  request.Event.AccessRestricted,
	})
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	return c.JSON(http.StatusOK, types.ModifyEventResponse{
		Event: *event,
	})
}

func (h Handler) UpdateEvent(c echo.Context) error {
	log.Info(fmt.Sprintf("Host: %v", c.Request().Host))
	log.Info(fmt.Sprintf("Referer: %v", c.Request().Referer()))
	log.Info(fmt.Sprintf("RemoteAddr: %v", c.Request().RemoteAddr))
	var request types.UpdateEventRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	// Get Key :: TODO :: Add verification of HOST value.
	key, err := database.GetKey(request.Key)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if key == nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	// Verify key access level.  Readonly cannot write or modify values.
	if key.Type == "read" {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	account, err := database.GetAccountByID(key.AccountIdentifier)
	if err != nil || account == nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	event, err := database.GetEvent(request.Event.Slug)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if event == nil {
		return getAPIError(c, http.StatusNotFound, "Event Not Found", nil)
	}
	if account.Type != "admin" && event.AccountIdentifier != account.Identifier {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	err = database.UpdateEvent(request.Event)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	event, err = database.GetEvent(request.Event.Slug)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	return c.JSON(http.StatusOK, types.ModifyEventResponse{
		Event: *event,
	})
}

func (h Handler) DeleteEvent(c echo.Context) error {
	log.Info(fmt.Sprintf("Host: %v", c.Request().Host))
	log.Info(fmt.Sprintf("Referer: %v", c.Request().Referer()))
	log.Info(fmt.Sprintf("RemoteAddr: %v", c.Request().RemoteAddr))
	var request types.DeleteEventRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	// Get Key :: TODO :: Add verification of HOST value.
	key, err := database.GetKey(request.Key)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if key == nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	// Verify access level. Delete is the only level that can delete values.
	if key.Type != "delete" {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	account, err := database.GetAccountByID(key.AccountIdentifier)
	if err != nil || account == nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	event, err := database.GetEvent(request.Slug)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if event == nil {
		return getAPIError(c, http.StatusNotFound, "Event Not Found", nil)
	}
	if account.Type != "admin" && event.AccountIdentifier != account.Identifier {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	err = database.DeleteEvent(*event)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	return c.NoContent(http.StatusOK)
}
