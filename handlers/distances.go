package handlers

import (
	"chronokeep/results/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) GetDistances(c echo.Context) error {
	// Get Key from Authorization Header
	k, err := retrieveKey(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Error Getting Key From Authorization Header", err)
	}
	if k == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key Not Provided in Authorization Header", nil)
	}
	var request types.GetDistancesRequest
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
	// And Event for verification of whether or not we can allow access to this key
	year := ""
	if request.Year != nil {
		year = *request.Year
	}
	mult, err := database.GetEventAndYear(request.Slug, year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event/Year", err)
	}
	if mult == nil || mult.Event == nil || mult.EventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Event/Year Not Found", nil)
	}
	if mult.Event.AccessRestricted && mkey.Account.Identifier != mult.Event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Restricted Event", nil)
	}
	if request.Distance != nil {
		dist, err := database.GetDistance(mult.EventYear.Identifier, *request.Distance)
		if err != nil {
			return getAPIError(c, http.StatusInternalServerError, "Error Fetching Distance", err)
		}
		return c.JSON(http.StatusOK, types.GetDistanceResponse{
			Distance: dist,
		})
	}
	distances, err := database.GetDistances(mult.EventYear.Identifier)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Fetching Distances", err)
	}
	return c.JSON(http.StatusOK, types.GetDistancesResponse{
		Distances: distances,
	})
}

func (h Handler) AddDistances(c echo.Context) error {
	// Get Key from Authorization Header
	k, err := retrieveKey(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Error Getting Key From Authorization Header", err)
	}
	if k == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key Not Provided in Authorization Header", nil)
	}
	var request types.AddDistancesRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	var distToAdd []types.Distance
	for _, dist := range request.Distances {
		if err := dist.Validate(h.validate); err == nil {
			distToAdd = append(distToAdd, dist)
		}
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
	distances, err := database.AddDistances(mult.EventYear.Identifier, distToAdd)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Adding Distances", err)
	}
	return c.JSON(http.StatusOK, types.GetDistancesResponse{
		Distances: distances,
	})
}

func (h Handler) DeleteDistances(c echo.Context) error {
	// Get Key from Authorization Header
	k, err := retrieveKey(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Error Getting Key From Authorization Header", err)
	}
	if k == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key Not Provided in Authorization Header", nil)
	}
	var request types.DeleteDistancesRequest
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
	// For results, let a write key delete.
	if mkey.Key.Type == "read" {
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
	count, err := database.DeleteDistances(mult.EventYear.Identifier)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Deleting Distances", nil)
	}
	return c.JSON(http.StatusOK, types.DeleteDistancesResponse{
		Count: count,
	})
}
