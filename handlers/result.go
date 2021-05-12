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
	mkey, err := database.GetKeyAndAccount(request.Key)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if mkey == nil || mkey.Key == nil || mkey.Account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	// And Event for verification of whether or not we can allow access to this key
	mult, err := database.GetEventAndYear(request.Slug, request.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if mult == nil || mult.Event == nil || mult.EventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Not Found", nil)
	}
	if mkey.Account.Type != "admin" && mult.Event.AccessRestricted && mkey.Account.Identifier != mult.Event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	results, err := database.GetResults(mult.EventYear.Identifier)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	return c.JSON(http.StatusOK, types.GetResultsResponse{
		Event:     *mult.Event,
		EventYear: *mult.EventYear,
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
	mkey, err := database.GetKeyAndAccount(request.Key)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if mkey == nil || mkey.Key == nil || mkey.Account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	// And Event for verification of whether or not we can allow access to this key
	mult, err := database.GetEventAndYear(request.Slug, request.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if mult == nil || mult.Event == nil || mult.EventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Not Found", nil)
	}
	// Check if the account is an admin or if they own this event.
	if mkey.Account.Type != "admin" && mult.Event.AccountIdentifier != mkey.Account.Identifier {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	results, err := database.AddResults(mult.EventYear.Identifier, request.Results)
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
	mkey, err := database.GetKeyAndAccount(request.Key)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if mkey == nil || mkey.Key == nil || mkey.Account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	// And Event for verification of whether or not we can allow access to this key
	mult, err := database.GetEventAndYear(request.Slug, request.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if mult == nil || mult.Event == nil || mult.EventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Not Found", nil)
	}
	// Check if the account is an admin or if they own this event.
	if mkey.Account.Type != "admin" && mult.Event.AccountIdentifier != mkey.Account.Identifier {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	count, err := database.DeleteEventResults(mult.EventYear.Identifier)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	return c.JSON(http.StatusOK, types.AddResultsResponse{
		Count: int(count),
	})
}
