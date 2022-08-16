package handlers

import (
	"chronokeep/results/types"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) RGetParticipants(c echo.Context) error {
	var request types.GetParticipantsRequest
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
	if account.Locked {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", errors.New("account locked"))
	}
	if len(request.Slug) < 1 || len(request.Year) < 1 {
		return getAPIError(c, http.StatusBadRequest, "Bad Request", errors.New("no slug/year specified"))
	}
	multi, err := database.GetAccountEventAndYear(request.Slug, request.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event/Year", nil)
	}
	if multi == nil || multi.Event == nil || multi.EventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Event/Year Not Found", nil)
	}
	// Verify they're allowed to pull these identifiers
	if account.Type != "admin" && account.Identifier != multi.Event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", errors.New("ownership error"))
	}
	participants, err := database.GetPeople(request.Slug, request.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Participants", err)
	}
	return c.JSON(http.StatusOK, types.GetParticipantsResponse{
		Participants: participants,
	})
}

func (h Handler) RAddParticipants(c echo.Context) error {
	var request types.AddParticipantsRequest
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
	if account.Locked {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", errors.New("account locked"))
	}
	if len(request.Slug) < 1 || len(request.Year) < 1 {
		return getAPIError(c, http.StatusBadRequest, "Bad Request", errors.New("no slug/year specified"))
	}
	multi, err := database.GetAccountEventAndYear(request.Slug, request.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event/Year", nil)
	}
	if multi == nil || multi.Event == nil || multi.EventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Event/Year Not Found", nil)
	}
	// Verify they're allowed to pull these identifiers
	if account.Type != "admin" && account.Identifier != multi.Event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", errors.New("ownership error"))
	}
	participants, err := database.AddPeople(multi.EventYear.Identifier, request.Participants)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Adding Participants", err)
	}
	return c.JSON(http.StatusOK, types.GetParticipantsResponse{
		Participants: participants,
	})
}

func (h Handler) RDeleteParticipants(c echo.Context) error {
	var request types.DeleteParticipantsRequest
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
	if account.Locked {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", errors.New("account locked"))
	}
	if len(request.Slug) < 1 || len(request.Year) < 1 {
		return getAPIError(c, http.StatusBadRequest, "Bad Request", errors.New("no slug/year specified"))
	}
	multi, err := database.GetAccountEventAndYear(request.Slug, request.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event/Year", nil)
	}
	if multi == nil || multi.Event == nil || multi.EventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Event/Year Not Found", nil)
	}
	// Verify they're allowed to pull these identifiers
	if account.Type != "admin" && account.Identifier != multi.Event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", errors.New("ownership error"))
	}
	err = database.DeletePeople(multi.EventYear.Identifier, request.Bibs)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Deleting Participants", err)
	}
	return c.NoContent(http.StatusOK)
}
