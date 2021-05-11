package handlers

import (
	"chronokeep/results/database"
	"chronokeep/results/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) GetEventYear(c echo.Context) error {
	var request types.GetEventYearRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	// Get Key
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
	event, err := database.GetEvent(request.Slug)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	eventYear, err := database.GetEventYear(request.Slug, request.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if event == nil || eventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Event Year Not Found", nil)
	}
	if account.Type != "admin" && account.Identifier != event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	return c.JSON(http.StatusOK, types.EventYearResponse{
		Event:     *event,
		EventYear: *eventYear,
	})
}

func (h Handler) AddEventYear(c echo.Context) error {
	var request types.ModifyEventYearRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	// Get Key
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
	event, err := database.GetEvent(request.Slug)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if event == nil {
		return getAPIError(c, http.StatusNotFound, "Event Not Found", nil)
	}
	// Verify they're allowed to add this event.
	if account.Type != "admin" && account.Identifier != event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	eventYear, err := database.AddEventYear(types.EventYear{
		EventIdentifier: event.Identifier,
		Year:            request.EventYear.Year,
		DateTime:        request.EventYear.DateTime,
		Live:            request.EventYear.Live,
	})
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	return c.JSON(http.StatusOK, types.EventYearResponse{
		Event:     *event,
		EventYear: *eventYear,
	})
}

func (h Handler) UpdateEventYear(c echo.Context) error {
	var request types.ModifyEventYearRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	// Get Key
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
	event, err := database.GetEvent(request.Slug)
	if err != nil || event == nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	// Verify they're allowed to modify this event year.
	if account.Type != "admin" && account.Identifier != event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	eventYear, err := database.GetEventYear(event.Slug, request.EventYear.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if eventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Event Year Not Found", nil)
	}
	err = database.UpdateEventYear(types.EventYear{
		EventIdentifier: eventYear.Identifier,
		Identifier:      eventYear.Identifier,
		Year:            eventYear.Year,
		DateTime:        request.EventYear.DateTime,
		Live:            request.EventYear.Live,
	})
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	eventYear, err = database.GetEventYear(event.Slug, request.EventYear.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	return c.JSON(http.StatusOK, types.EventYearResponse{
		Event:     *event,
		EventYear: *eventYear,
	})
}

func (h Handler) DeleteEventYear(c echo.Context) error {
	var request types.DeleteEventYearRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	// Get Key
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
	event, err := database.GetEvent(request.Slug)
	if err != nil || event == nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	// Verify they're allowed to modify this event year.
	if account.Type != "admin" && account.Identifier != event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	eventYear, err := database.GetEventYear(event.Slug, request.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if eventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Event Year Not Found", nil)
	}
	err = database.DeleteEventYear(*eventYear)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	return c.NoContent(http.StatusOK)
}
