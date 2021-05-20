package handlers

import (
	"chronokeep/results/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

// RGetEvents Used with a JWT for the website to get an account events.
func (h Handler) RGetEvents(c echo.Context) error {
	var request types.GetEventsRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	account, err := verifyToken(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized Token", err)
	}
	if account == nil {
		return getAPIError(c, http.StatusNotFound, "Account Not Found", nil)
	}
	// If the user is not an admin and they're requesting events for another account deny them.
	if account.Type != "admin" && request.Email != nil && account.Email != *request.Email {
		return getAPIError(c, http.StatusUnauthorized, "Ownership Error", nil)
	}
	// We're either pulling account events for the calling account, or the requesting email
	email := account.Email
	if request.Email != nil {
		email = *request.Email
	}
	events, err := database.GetAccountEvents(email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Events", nil)
	}
	return c.JSON(http.StatusOK, types.GetEventsResponse{
		Events: events,
	})
}

// RAddEvent Used with a JWT for the website to add an event.
func (h Handler) RAddEvent(c echo.Context) error {
	var request types.AddEventRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	account, err := verifyToken(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized Token", err)
	}
	if account == nil {
		return getAPIError(c, http.StatusNotFound, "Account Not Found", nil)
	}
	// Verify that the user has access to add an event if the email is set.
	if request.Email != nil && *request.Email != account.Email && account.Type != "admin" {
		return getAPIError(c, http.StatusUnauthorized, "Ownership Error", nil)
	}
	id := account.Identifier
	if request.Email != nil {
		a, err := database.GetAccount(*request.Email)
		if err != nil {
			return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Account", err)
		}
		id = a.Identifier
	}
	event, err := database.AddEvent(types.Event{
		AccountIdentifier: id,
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

// RUpdateEvent Used with a JWT for the website to update an event.
func (h Handler) RUpdateEvent(c echo.Context) error {
	var request types.UpdateEventRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	account, err := verifyToken(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized Token", err)
	}
	if account == nil {
		return getAPIError(c, http.StatusNotFound, "Account Not Found", nil)
	}
	event, err := database.GetEvent(request.Event.Slug)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event", err)
	}
	if event == nil {
		return getAPIError(c, http.StatusNotFound, "Event Not Found", nil)
	}
	// Verify that the user has access to update the event.
	if event.AccountIdentifier != account.Identifier && account.Type != "admin" {
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

// RDeleteEvent Used with a JWT for the website to delete an event.
func (h Handler) RDeleteEvent(c echo.Context) error {
	var request types.DeleteEventRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	account, err := verifyToken(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized Token", err)
	}
	if account == nil {
		return getAPIError(c, http.StatusNotFound, "Account Not Found", nil)
	}
	event, err := database.GetEvent(request.Slug)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event", err)
	}
	if event == nil {
		return getAPIError(c, http.StatusNotFound, "Event Not Found", nil)
	}
	// Verify that the user has access to delete the event.
	if event.AccountIdentifier != account.Identifier && account.Type != "admin" {
		return getAPIError(c, http.StatusUnauthorized, "Ownership Error", nil)
	}
	err = database.DeleteEvent(*event)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Deleting Event", err)
	}
	return c.NoContent(http.StatusOK)
}
