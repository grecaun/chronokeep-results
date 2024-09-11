package handlers

import (
	"chronokeep/results/types"
	"errors"
	"fmt"
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
	year := ""
	if request.Year != nil {
		year = *request.Year
	}
	if len(request.Slug) < 1 {
		return getAPIError(c, http.StatusBadRequest, "Bad Request", errors.New("no slug/year specified"))
	}
	multi, err := database.GetAccountEventAndYear(request.Slug, year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event/Year", nil)
	}
	if multi == nil || multi.Event == nil || multi.EventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Event/Year Not Found", nil)
	}
	linked, err := database.GetLinkedAccounts(multi.Account.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Linked Accounts", nil)
	}
	is_linked := false
	for _, acc := range linked {
		if acc.Identifier == account.Identifier {
			is_linked = true
			break
		}
	}
	// Verify they're allowed to pull these identifiers
	if account.Type != "admin" && account.Identifier != multi.Event.AccountIdentifier && !is_linked {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", errors.New("ownership error"))
	}
	limit := 0
	page := 0
	if request.Limit != nil {
		limit = *request.Limit
	}
	if request.Page != nil {
		page = *request.Page
		if page > 0 {
			page--
		}
	}
	participants, err := database.GetParticipants(multi.EventYear.Identifier, limit, page)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Participants", err)
	}
	return c.JSON(http.StatusOK, types.GetParticipantsResponse{
		Event:        *multi.Event,
		Year:         *multi.EventYear,
		Participants: participants,
	})
}

func (h Handler) RAddParticipant(c echo.Context) error {
	var request types.AddParticipantRequest
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
	linked, err := database.GetLinkedAccounts(multi.Account.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Linked Accounts", nil)
	}
	is_linked := false
	for _, acc := range linked {
		if acc.Identifier == account.Identifier {
			is_linked = true
			break
		}
	}
	// Verify they're allowed to pull these identifiers
	if account.Type != "admin" && account.Identifier != multi.Event.AccountIdentifier && !is_linked {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", errors.New("ownership error"))
	}
	// validate participants
	var partToAdd []types.Participant
	// Validate, only add if it passes validation.
	if err := request.Participant.Validate(h.validate); err == nil {
		partToAdd = append(partToAdd, request.Participant)
	}
	if len(partToAdd) < 1 {
		return getAPIError(c, http.StatusBadRequest, "Invalid", nil)
	}
	// create random alternate id if none is set
	if len(partToAdd[0].AlternateId) < 1 {
		partToAdd[0].AlternateId = fmt.Sprintf("new%s%s", partToAdd[0].First, partToAdd[0].Last)
	}
	participants, err := database.AddParticipants(multi.EventYear.Identifier, partToAdd)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Adding Participants", err)
	}
	if len(participants) > 1 {
		return getAPIError(c, http.StatusInternalServerError, "Multiple Participants Added", nil)
	}
	return c.JSON(http.StatusOK, types.UpdateParticipantResponse{
		Participant: participants[0],
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
	count, err := database.DeleteParticipants(multi.EventYear.Identifier, request.Identifiers)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Deleting Participants", err)
	}
	return c.JSON(http.StatusOK, types.AddResultsResponse{
		Count: int(count),
	})
}

func (h Handler) RUpdateParticipant(c echo.Context) error {
	var request types.UpdateParticipantRequest
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
	linked, err := database.GetLinkedAccounts(multi.Account.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Linked Accounts", nil)
	}
	is_linked := false
	for _, acc := range linked {
		if acc.Identifier == account.Identifier {
			is_linked = true
			break
		}
	}
	// Verify they're allowed to pull these identifiers
	if account.Type != "admin" && account.Identifier != multi.Event.AccountIdentifier && !is_linked {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", errors.New("ownership error"))
	}
	part, err := database.UpdateParticipant(multi.EventYear.Identifier, request.Participant)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Updating Participant", err)
	}
	return c.JSON(http.StatusOK, types.UpdateParticipantRequest{
		Participant: *part,
	})
}

func (h Handler) RUpdateManyParticipants(c echo.Context) error {
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
	linked, err := database.GetLinkedAccounts(multi.Account.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Linked Accounts", nil)
	}
	is_linked := false
	for _, acc := range linked {
		if acc.Identifier == account.Identifier {
			is_linked = true
			break
		}
	}
	// Verify they're allowed to pull these identifiers
	if account.Type != "admin" && account.Identifier != multi.Event.AccountIdentifier && !is_linked {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", errors.New("ownership error"))
	}
	// validate participants
	var partsToAdd []types.Participant
	// Validate, only add if it passes validation.
	for _, part := range request.Participants {
		if err := part.Validate(h.validate); err == nil {
			partsToAdd = append(partsToAdd, part)
		}
	}
	if len(partsToAdd) < 1 {
		return getAPIError(c, http.StatusBadRequest, "No Valid Participants", nil)
	}
	participants, err := database.UpdateParticipants(multi.EventYear.Identifier, partsToAdd)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Adding Participants", err)
	}
	return c.JSON(http.StatusOK, types.UpdateParticipantsResponse{
		Participants: participants,
	})
}

func (h Handler) RAddManyParticipants(c echo.Context) error {
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
	linked, err := database.GetLinkedAccounts(multi.Account.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Linked Accounts", nil)
	}
	is_linked := false
	for _, acc := range linked {
		if acc.Identifier == account.Identifier {
			is_linked = true
			break
		}
	}
	// Verify they're allowed to pull these identifiers
	if account.Type != "admin" && account.Identifier != multi.Event.AccountIdentifier && !is_linked {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", errors.New("ownership error"))
	}
	// validate participants
	var partsToAdd []types.Participant
	// Validate, only add if it passes validation.
	for _, part := range request.Participants {
		if err := part.Validate(h.validate); err == nil {
			// create random alternate id if none is set
			if len(part.AlternateId) < 1 || part.AlternateId == "-1" || part.AlternateId == "0" {
				part.AlternateId = fmt.Sprintf("new%s%s", part.First, part.Last)
			}
			partsToAdd = append(partsToAdd, part)
		}
	}
	if len(partsToAdd) < 1 {
		return getAPIError(c, http.StatusBadRequest, "Invalid", nil)
	}
	participants, err := database.AddParticipants(multi.EventYear.Identifier, partsToAdd)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Adding Participants", err)
	}
	return c.JSON(http.StatusOK, types.UpdateParticipantsResponse{
		Participants: participants,
	})
}
