package handlers

import (
	"chronokeep/results/types"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (h Handler) GetEvents(c echo.Context) error {
	// Get Key from Authorization Header
	k, err := retrieveKey(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Error Getting Key From Authorization Header", err)
	}
	if k == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key Not Provided in Authorization Header", nil)
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
	events, err := database.GetEvents()
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Events", err)
	}
	return c.JSON(http.StatusOK, types.GetEventsResponse{
		Events: events,
	})
}

func (h Handler) GetMyEvents(c echo.Context) error {
	// Get Key from Authorization Header
	k, err := retrieveKey(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Error Getting Key From Authorization Header", err)
	}
	if k == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key Not Provided in Authorization Header", nil)
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
	events, err := database.GetAccountEvents(mkey.Account.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Events", err)
	}
	return c.JSON(http.StatusOK, types.GetEventsResponse{
		Events: events,
	})
}

func (h Handler) GetEvent(c echo.Context) error {
	// Get Key from Authorization Header
	k, err := retrieveKey(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Error Getting Key From Authorization Header", err)
	}
	if k == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key Not Provided in Authorization Header", nil)
	}
	var request types.GetEventRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
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
		res, err = database.GetResults(recent.Identifier, 0, 0)
		if err != nil {
			return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Results", err)
		}
	}
	outRes := make(map[string][]types.Result)
	for _, result := range res {
		if _, ok := outRes[result.Distance]; !ok {
			outRes[result.Distance] = make([]types.Result, 0, 1)
		}
		outRes[result.Distance] = append(outRes[result.Distance], result)
	}
	return c.JSON(http.StatusOK, types.GetEventResponse{
		Event:      *event,
		EventYears: eventYears,
		Year:       recent,
		Results:    outRes,
	})
}

func (h Handler) AddEvent(c echo.Context) error {
	// Get Key from Authorization Header
	k, err := retrieveKey(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Error Getting Key From Authorization Header", err)
	}
	if k == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key Not Provided in Authorization Header", nil)
	}
	var request types.AddEventRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	// Validate the Event
	if err := request.Event.Validate(h.validate); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Validation Error", err)
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
		Type:              request.Event.Type,
	})
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Adding Event (Duplicate Slug/Name Likely)", err)
	}
	return c.JSON(http.StatusOK, types.ModifyEventResponse{
		Event: *event,
	})
}

func (h Handler) UpdateEvent(c echo.Context) error {
	// Get Key from Authorization Header
	k, err := retrieveKey(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Error Getting Key From Authorization Header", err)
	}
	if k == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key Not Provided in Authorization Header", nil)
	}
	var request types.UpdateEventRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	// Validate the Event Year
	if err := request.Event.Validate(h.validate); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Validation Error", err)
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
		Type:             request.Event.Type,
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
	// Get Key from Authorization Header
	k, err := retrieveKey(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Error Getting Key From Authorization Header", err)
	}
	if k == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key Not Provided in Authorization Header", nil)
	}
	var request types.DeleteEventRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
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
