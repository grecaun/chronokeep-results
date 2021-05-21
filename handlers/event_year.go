package handlers

import (
	"chronokeep/results/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) GetEventYear(c echo.Context) error {
	var request types.GetEventYearRequest
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
	mult, err := database.GetEventAndYear(request.Slug, request.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event/Year", err)
	}
	if mult == nil || mult.Event == nil || mult.EventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Event/Year Not Found", nil)
	}
	// Only the account owner can access restricted events.
	if mult.Event.AccessRestricted && mkey.Account.Identifier != mult.Event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Restricted Event", nil)
	}
	return c.JSON(http.StatusOK, types.EventYearResponse{
		Event:     *mult.Event,
		EventYear: *mult.EventYear,
	})
}

func (h Handler) AddEventYear(c echo.Context) error {
	var request types.ModifyEventYearRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	// Validate the Event Year
	if err := request.EventYear.Validate(h.validate); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Validation Error", err)
	}
	// Get Key :: TODO :: Add verification of HOST value.
	mkey, err := database.GetKeyAndAccount(request.Key)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Key/Account", err)
	}
	if mkey == nil || mkey.Key == nil || mkey.Account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key/Account Not Found", nil)
	}
	// Verify key access level.  Readonly cannot write or modify values.
	if mkey.Key.Type == "read" {
		return getAPIError(c, http.StatusUnauthorized, "Key is ReadOnly", nil)
	}
	event, err := database.GetEvent(request.Slug)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event", err)
	}
	if event == nil {
		return getAPIError(c, http.StatusNotFound, "Event Not Found", nil)
	}
	// Verify they're allowed to add this event.
	if mkey.Account.Identifier != event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Ownership Error", nil)
	}
	eventYear, err := database.AddEventYear(types.EventYear{
		EventIdentifier: event.Identifier,
		Year:            request.EventYear.Year,
		DateTime:        request.EventYear.GetDateTime(),
		Live:            request.EventYear.Live,
	})
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Adding Event Year (Duplicate Year Likely)", err)
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
	// Validate the Event Year
	if err := request.EventYear.Validate(h.validate); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Validation Error", err)
	}
	// Get Key :: TODO :: Add verification of HOST value.
	mkey, err := database.GetKeyAndAccount(request.Key)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Key/Account", err)
	}
	if mkey == nil || mkey.Key == nil || mkey.Account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key/Account Not Found", nil)
	}
	// Verify key access level.  Readonly cannot write or modify values.
	if mkey.Key.Type == "read" {
		return getAPIError(c, http.StatusUnauthorized, "Key is ReadOnly", nil)
	}
	mult, err := database.GetEventAndYear(request.Slug, request.EventYear.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event/Year", err)
	}
	if mult == nil || mult.Event == nil || mult.EventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Event/Year Not Found", nil)
	}
	// Verify they're allowed to modify this event year.
	if mkey.Account.Identifier != mult.Event.AccountIdentifier {
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

func (h Handler) DeleteEventYear(c echo.Context) error {
	var request types.DeleteEventYearRequest
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
	// Verify access level. Delete is the only level that can delete values.
	if mkey.Key.Type != "delete" {
		return getAPIError(c, http.StatusUnauthorized, "Key is ReadOnly/Write", nil)
	}
	mult, err := database.GetEventAndYear(request.Slug, request.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event/Year", err)
	}
	if mult == nil || mult.Event == nil || mult.EventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Event/Year Not Found", nil)
	}
	// Verify they're allowed to modify this event year.
	if mkey.Account.Identifier != mult.Event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Ownership Error", nil)
	}
	err = database.DeleteEventYear(*mult.EventYear)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Deleting Event Year", err)
	}
	return c.NoContent(http.StatusOK)
}
