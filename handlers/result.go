package handlers

import (
	"chronokeep/results/database"
	"chronokeep/results/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) GetResults(c echo.Context) error {
	var request types.GetResultsRequest
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
	// Get Account
	account, err := database.GetAccountByID(key.AccountIdentifier)
	if err != nil || account == nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	// And Event for verification of whether or not we can allow access to this key
	event, err := database.GetEvent(request.Slug)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if event == nil {
		return getAPIError(c, http.StatusNotFound, "Event Not Found", nil)
	}
	if account.Type != "admin" && event.AccessRestricted && account.Identifier != event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	eventYear, err := database.GetEventYear(request.Slug, request.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if eventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Year Not Found", nil)
	}
	results, err := database.GetResults(eventYear.Identifier)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	return c.JSON(http.StatusOK, types.GetResultsResponse{
		Event:     *event,
		EventYear: *eventYear,
		Results:   results,
		Count:     len(results),
	})
}

func (h Handler) AddResults(c echo.Context) error {
	var request types.AddResultsRequest
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
	// Get Account
	account, err := database.GetAccountByID(key.AccountIdentifier)
	if err != nil || account == nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	// And Event for verification of whether or not we can allow access to this key
	event, err := database.GetEvent(request.Slug)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if event == nil {
		return getAPIError(c, http.StatusNotFound, "Event Not Found", nil)
	}
	// Check if the account is an admin or if they own this event.
	if account.Type != "admin" && event.AccountIdentifier != account.Identifier {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	eventYear, err := database.GetEventYear(request.Slug, request.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if eventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Year Not Found", nil)
	}
	results, err := database.AddResults(eventYear.Identifier, request.Results)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	return c.JSON(http.StatusOK, types.AddResultsResponse{
		Count: len(results),
	})
}

func (h Handler) DeleteResults(c echo.Context) error {
	var request types.DeleteResultsRequest
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
	// Get Account
	account, err := database.GetAccountByID(key.AccountIdentifier)
	if err != nil || account == nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	// And Event for verification of whether or not we can allow access to this key
	event, err := database.GetEvent(request.Slug)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if event == nil {
		return getAPIError(c, http.StatusNotFound, "Event Not Found", nil)
	}
	// Check if the account is an admin or if they own this event.
	if account.Type != "admin" && event.AccountIdentifier != account.Identifier {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	eventYear, err := database.GetEventYear(request.Slug, request.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if eventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Year Not Found", nil)
	}
	count, err := database.DeleteEventResults(eventYear.Identifier)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	return c.JSON(http.StatusOK, types.AddResultsResponse{
		Count: int(count),
	})
}
