package handlers

import (
	"chronokeep/results/types"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) GetParticipants(c echo.Context) error {
	// Get Key from Auth Header
	k, err := retrieveKey(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Error Getting Key From Authorization Header", err)
	}
	if k == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key Not Provided in Authorization Header", nil)
	}
	var request types.GetParticipantsRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	if len(request.Slug) < 1 || len(request.Year) < 1 {
		return getAPIError(c, http.StatusBadRequest, "Bad Request", errors.New("no slug/year specified"))
	}
	// Get Key
	mkey, err := database.GetKeyAndAccount(*k)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Key/Account", err)
	}
	if mkey == nil || mkey.Key == nil || mkey.Account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key/Account Not Found", nil)
	}
	// Check for expired key
	if mkey.Key.Expired() {
		return getAPIError(c, http.StatusUnauthorized, "Expired Key", nil)
	}
	// Check for host being allowed.
	if !mkey.Key.IsAllowed(c.Request().Referer()) {
		return getAPIError(c, http.StatusUnauthorized, "Host Not Allowed", nil)
	}
	mult, err := database.GetEventAndYear(request.Slug, request.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event Year", err)
	}
	if mult == nil || mult.Event == nil || mult.EventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Event/Year Not Found", nil)
	}
	// Only the account owner can access restricted events.
	if mult.Event.AccessRestricted && mkey.Account.Identifier != mult.Event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Restricted Event", nil)
	}
	participants, err := database.GetPeople(request.Slug, request.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Participants", err)
	}
	return c.JSON(http.StatusOK, types.GetParticipantsResponse{
		Participants: participants,
	})
}

func (h Handler) AddParticipants(c echo.Context) error {
	// Get Key from Auth Header
	k, err := retrieveKey(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Error Getting Key From Authorization Header", err)
	}
	if k == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key Not Provided in Authorization Header", nil)
	}
	var request types.AddParticipantsRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	if len(request.Slug) < 1 || len(request.Year) < 1 {
		return getAPIError(c, http.StatusBadRequest, "Bad Request", errors.New("no slug/year specified"))
	}
	// Get Key
	mkey, err := database.GetKeyAndAccount(*k)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Key/Account", err)
	}
	if mkey == nil || mkey.Key == nil || mkey.Account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key/Account Not Found", nil)
	}
	// Check for expired key
	if mkey.Key.Expired() {
		return getAPIError(c, http.StatusUnauthorized, "Expired Key", nil)
	}
	// Check for host being allowed.
	if !mkey.Key.IsAllowed(c.Request().Referer()) {
		return getAPIError(c, http.StatusUnauthorized, "Host Not Allowed", nil)
	}
	// Check to ensure key isn't read only
	if mkey.Key.Type != "write" && mkey.Key.Type != "delete" {
		return getAPIError(c, http.StatusUnauthorized, "Key Not Allowed", errors.New("read key not allowed to write"))
	}
	mult, err := database.GetEventAndYear(request.Slug, request.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event Year", err)
	}
	if mult == nil || mult.Event == nil || mult.EventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Event/Year Not Found", nil)
	}
	// Only the account owner or admins can add.
	if mkey.Account.Identifier != mult.Event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Restricted Event", nil)
	}
	// validate participants
	var partToAdd []types.Person
	for _, part := range request.Participants {
		// Validate all results, only add the results that pass validation.
		if err := part.Validate(h.validate); err == nil {
			partToAdd = append(partToAdd, part)
		}
	}
	if len(partToAdd) < 1 {
		return getAPIError(c, http.StatusBadRequest, "No Valid Participants", nil)
	}
	participants, err := database.AddPeople(mult.EventYear.Identifier, partToAdd)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Adding Participants", err)
	}
	return c.JSON(http.StatusOK, types.AddResultsResponse{
		Count: len(participants),
	})
}

func (h Handler) DeleteParticipants(c echo.Context) error {
	// Get Key from Auth Header
	k, err := retrieveKey(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Error Getting Key From Authorization Header", err)
	}
	if k == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key Not Provided in Authorization Header", nil)
	}
	var request types.DeleteParticipantsRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	if len(request.Slug) < 1 || len(request.Year) < 1 {
		return getAPIError(c, http.StatusBadRequest, "Bad Request", errors.New("no slug/year specified"))
	}
	// Get Key
	mkey, err := database.GetKeyAndAccount(*k)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Key/Account", err)
	}
	if mkey == nil || mkey.Key == nil || mkey.Account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key/Account Not Found", nil)
	}
	// Check for expired key
	if mkey.Key.Expired() {
		return getAPIError(c, http.StatusUnauthorized, "Expired Key", nil)
	}
	// Check to ensure not a read key
	if mkey.Key.Type != "write" && mkey.Key.Type != "delete" {
		return getAPIError(c, http.StatusUnauthorized, "Invalid Key Type", nil)
	}
	// Check for host being allowed.
	if !mkey.Key.IsAllowed(c.Request().Referer()) {
		return getAPIError(c, http.StatusUnauthorized, "Host Not Allowed", nil)
	}
	mult, err := database.GetEventAndYear(request.Slug, request.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event Year", err)
	}
	if mult == nil || mult.Event == nil || mult.EventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Event/Year Not Found", nil)
	}
	// Only the account owner or admins can delete.
	if mkey.Account.Identifier != mult.Event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Restricted Event", nil)
	}
	count, err := database.DeletePeople(mult.EventYear.Identifier, request.Bibs)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Deleting Participants", err)
	}
	return c.JSON(http.StatusOK, types.AddResultsResponse{
		Count: int(count),
	})
}
