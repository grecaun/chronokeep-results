package handlers

import (
	"chronokeep/results/types"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func (h Handler) GetEvents(c echo.Context) error {
	log.Info(fmt.Sprintf("Host: %v", c.Request().Host))
	log.Info(fmt.Sprintf("Referer: %v", c.Request().Referer()))
	log.Info(fmt.Sprintf("RemoteAddr: %v", c.Request().RemoteAddr))
	var request types.GetEventsRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	// Get Key :: TODO :: Add verification of HOST value.
	mkey, err := database.GetKeyAndAccount(request.Key)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Key/Account", err)
	}
	if mkey == nil || mkey.Key == nil || mkey.Account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key/Account Not Found", nil)
	}
	if mkey.Key.Expired() {
		return getAPIError(c, http.StatusUnauthorized, "Expired Key", nil)
	}
	var events []types.Event
	if request.Email != nil {
		// Only account owners can pull their account events with a key.
		if mkey.Account.Email != *request.Email {
			return getAPIError(c, http.StatusUnauthorized, "Ownership Error", nil)
		}
		events, err = database.GetAccountEvents(*request.Email)
	} else {
		events, err = database.GetEvents()
	}
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Events", err)
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
	mkey, err := database.GetKeyAndAccount(request.Key)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Key/Account", err)
	}
	if mkey == nil || mkey.Key == nil || mkey.Account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key/Account Not Found", nil)
	}
	if mkey.Key.Expired() {
		return getAPIError(c, http.StatusUnauthorized, "Expired Key", nil)
	}
	event, err := database.GetEvent(request.Slug)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event", err)
	}
	if event == nil {
		return getAPIError(c, http.StatusNotFound, "Event Not Found", nil)
	}
	// Only account owner can access restricted events
	if event.AccessRestricted && mkey.Account.Identifier != event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Restricted Event", nil)
	}
	eventYears, err := database.GetEventYears(event.Slug)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event Years", err)
	}
	var recent *types.EventYear
	now := time.Now()
	if len(eventYears) > 0 {
		recent = &eventYears[0]
		for _, y := range eventYears[1:] {
			if recent.DateTime.Before(y.DateTime) && y.DateTime.Before(now) {
				recent = &y
			}
		}
	}
	var res []types.Result
	if recent != nil {
		res, err = database.GetResults(recent.Identifier)
		if err != nil {
			return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Results", err)
		}
	}
	return c.JSON(http.StatusOK, types.GetEventResponse{
		Event:      *event,
		EventYears: eventYears,
		Year:       recent,
		Results:    res,
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
	mkey, err := database.GetKeyAndAccount(request.Key)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Key/Account", err)
	}
	if mkey == nil || mkey.Key == nil || mkey.Account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key/Account Not Found", nil)
	}
	if mkey.Key.Expired() {
		return getAPIError(c, http.StatusUnauthorized, "Expired Key", nil)
	}
	// Verify key access level.  Readonly cannot write or modify values.
	if mkey.Key.Type == "read" {
		return getAPIError(c, http.StatusUnauthorized, "Key is ReadOnly", nil)
	}
	event, err := database.AddEvent(types.Event{
		AccountIdentifier: mkey.Account.Identifier,
		Name:              request.Event.Name,
		Slug:              request.Event.Slug,
		Website:           request.Event.Website,
		Image:             request.Event.Image,
		ContactEmail:      request.Event.ContactEmail,
		AccessRestricted:  request.Event.AccessRestricted,
	})
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Adding Event (Duplicate Slug/Name Likely)", err)
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
	mkey, err := database.GetKeyAndAccount(request.Key)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Key/Account", err)
	}
	if mkey == nil || mkey.Key == nil || mkey.Account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key/Account Not Found", nil)
	}
	if mkey.Key.Expired() {
		return getAPIError(c, http.StatusUnauthorized, "Expired Key", nil)
	}
	// Verify key access level.  Readonly cannot write or modify values.
	if mkey.Key.Type == "read" {
		return getAPIError(c, http.StatusUnauthorized, "Key is ReadOnly", nil)
	}
	event, err := database.GetEvent(request.Event.Slug)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event", err)
	}
	if event == nil {
		return getAPIError(c, http.StatusNotFound, "Event Not Found", nil)
	}
	if event.AccountIdentifier != mkey.Account.Identifier {
		return getAPIError(c, http.StatusUnauthorized, "Ownership Error", nil)
	}
	err = database.UpdateEvent(types.Event{
		Identifier:       event.Identifier,
		Name:             request.Event.Name,
		Slug:             request.Event.Slug,
		ContactEmail:     request.Event.ContactEmail,
		Website:          request.Event.Website,
		Image:            request.Event.Image,
		AccessRestricted: request.Event.AccessRestricted,
	})
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Updating Event (Nothing to Update / Name Conflict)", err)
	}
	event, err = database.GetEvent(request.Event.Slug)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event", err)
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
	mkey, err := database.GetKeyAndAccount(request.Key)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Key/Account", err)
	}
	if mkey == nil || mkey.Key == nil || mkey.Account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key/Account Not Found", nil)
	}
	if mkey.Key.Expired() {
		return getAPIError(c, http.StatusUnauthorized, "Expired Key", nil)
	}
	// Verify access level. Delete is the only level that can delete values.
	if mkey.Key.Type != "delete" {
		return getAPIError(c, http.StatusUnauthorized, "Key is ReadOnly/Write", nil)
	}
	event, err := database.GetEvent(request.Slug)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event", err)
	}
	if event == nil {
		return getAPIError(c, http.StatusNotFound, "Event Not Found", nil)
	}
	if event.AccountIdentifier != mkey.Account.Identifier {
		return getAPIError(c, http.StatusUnauthorized, "Ownership Error", nil)
	}
	err = database.DeleteEvent(*event)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Deleting Event", err)
	}
	return c.NoContent(http.StatusOK)
}
