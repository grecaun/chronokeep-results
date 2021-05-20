package handlers

import (
	"chronokeep/results/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) RGetEventYears(c echo.Context) error {
	var request types.GetEventYearRequest
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
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event", nil)
	}
	if event == nil {
		return getAPIError(c, http.StatusNotFound, "Event Not Found", nil)
	}
	// Verify they're allowed to pull these identifiers
	if account.Type != "admin" && account.Identifier != event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Ownership Error", nil)
	}
	eventYears, err := database.GetEventYears(request.Slug)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event Years", nil)
	}
	return c.JSON(http.StatusOK, types.EventYearsResponse{
		EventYears: eventYears,
	})
}

func (h Handler) RAddEventYear(c echo.Context) error {
	var request types.ModifyEventYearRequest
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
	// Verify they're allowed to add this event.
	if account.Identifier != event.AccountIdentifier && account.Type != "admin" {
		return getAPIError(c, http.StatusUnauthorized, "Ownership Error", nil)
	}
	eventYear, err := database.AddEventYear(types.EventYear{
		EventIdentifier: event.Identifier,
		Year:            request.EventYear.Year,
		DateTime:        request.EventYear.GetDateTime(),
		Live:            request.EventYear.Live,
	})
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Adding Event Year", err)
	}
	return c.JSON(http.StatusOK, types.EventYearResponse{
		Event:     *event,
		EventYear: *eventYear,
	})
}

func (h Handler) RUpdateEventYear(c echo.Context) error {
	var request types.ModifyEventYearRequest
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
	mult, err := database.GetEventAndYear(request.Slug, request.EventYear.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event/Year", err)
	}
	if mult == nil || mult.Event == nil || mult.EventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Event/Year Not Found", nil)
	}
	// Verify they're allowed to modify this event year.
	if account.Identifier != mult.Event.AccountIdentifier && account.Type != "admin" {
		return getAPIError(c, http.StatusUnauthorized, "Ownership Error", nil)
	}
	err = database.UpdateEventYear(types.EventYear{
		EventIdentifier: mult.EventYear.EventIdentifier,
		Identifier:      mult.EventYear.Identifier,
		Year:            mult.EventYear.Year,
		DateTime:        request.EventYear.GetDateTime(),
		Live:            request.EventYear.Live,
	})
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Updating Event", err)
	}
	eventYear, err := database.GetEventYear(request.Slug, request.EventYear.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event Year", err)
	}
	return c.JSON(http.StatusOK, types.EventYearResponse{
		Event:     *mult.Event,
		EventYear: *eventYear,
	})
}

func (h Handler) RDeleteEventYear(c echo.Context) error {
	var request types.DeleteEventYearRequest
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
	mult, err := database.GetEventAndYear(request.Slug, request.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event/Year", err)
	}
	if mult == nil || mult.Event == nil || mult.EventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Event/Year Not Found", nil)
	}
	// Verify they're allowed to modify this event year.
	if account.Identifier != mult.Event.AccountIdentifier && account.Type != "admin" {
		return getAPIError(c, http.StatusUnauthorized, "Ownership Error", nil)
	}
	err = database.DeleteEventYear(*mult.EventYear)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Deleting Event Year", err)
	}
	return c.NoContent(http.StatusOK)
}
