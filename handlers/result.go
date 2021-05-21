package handlers

import (
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
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Key/Account", err)
	}
	if mkey == nil || mkey.Key == nil || mkey.Account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key/Account Not Found", nil)
	}
	// And Event for verification of whether or not we can allow access to this key
	mult, err := database.GetEventAndYear(request.Slug, request.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event/Year", err)
	}
	if mult == nil || mult.Event == nil || mult.EventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Event/Year Not FOund", nil)
	}
	if mult.Event.AccessRestricted && mkey.Account.Identifier != mult.Event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Restricted Event", nil)
	}
	results, err := database.GetResults(mult.EventYear.Identifier)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Results", err)
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
	var resToAdd []types.Result
	for _, res := range request.Results {
		// Validate all results, only add the results that pass validation.
		if err := res.Validate(h.validate); err == nil {
			resToAdd = append(resToAdd, res)
		}
	}
	// Get Key :: TODO :: Add verification of HOST value.
	mkey, err := database.GetKeyAndAccount(request.Key)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Key/Account", err)
	}
	if mkey == nil || mkey.Key == nil || mkey.Account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key/Account Not Found", nil)
	}
	if mkey.Key.Type == "read" {
		return getAPIError(c, http.StatusUnauthorized, "Key is ReadOnly", nil)
	}
	// And Event for verification of whether or not we can allow access to this key
	mult, err := database.GetEventAndYear(request.Slug, request.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event/Year", err)
	}
	if mult == nil || mult.Event == nil || mult.EventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Event/Year Not Found", nil)
	}
	// Check if they own this event.
	if mult.Event.AccountIdentifier != mkey.Account.Identifier {
		return getAPIError(c, http.StatusUnauthorized, "Ownership Error", nil)
	}
	results, err := database.AddResults(mult.EventYear.Identifier, resToAdd)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Adding Results", err)
	}
	return c.JSON(http.StatusOK, types.AddResultsResponse{
		Count: len(results),
	})
}

func (h Handler) DeleteResults(c echo.Context) error {
	var request types.GetResultsRequest
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
	if mkey.Key.Type != "delete" {
		return getAPIError(c, http.StatusUnauthorized, "Key is ReadOnly/Write", nil)
	}
	// And Event for verification of whether or not we can allow access to this key
	mult, err := database.GetEventAndYear(request.Slug, request.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event/Year", err)
	}
	if mult == nil || mult.Event == nil || mult.EventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Event/Year Not Found", nil)
	}
	// Check if the account is an admin or if they own this event.
	if mult.Event.AccountIdentifier != mkey.Account.Identifier {
		return getAPIError(c, http.StatusUnauthorized, "Ownership Error", nil)
	}
	count, err := database.DeleteEventResults(mult.EventYear.Identifier)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Deleting Results", err)
	}
	return c.JSON(http.StatusOK, types.AddResultsResponse{
		Count: int(count),
	})
}
