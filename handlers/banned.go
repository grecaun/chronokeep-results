package handlers

import (
	"chronokeep/results/types"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) AddBannedPhone(c echo.Context) error {
	if c.Request().Method != echo.POST {
		return getAPIError(c, http.StatusBadRequest, "Invalid Method", nil)
	}
	// No need for keys for any of these calls
	var request types.ModifyBannedPhoneRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	err := h.validate.Struct(request)
	if len(request.Phone) < 10 || err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Phone Field", nil)
	}
	err = database.AddBlockedPhone(request.Phone)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Adding Phone", err)
	}
	return c.NoContent(http.StatusOK)
}

func (h Handler) GetBannedPhones(c echo.Context) error {
	if c.Request().Method != echo.GET {
		return getAPIError(c, http.StatusBadRequest, "Invalid Method", nil)
	}
	phones, err := database.GetBlockedPhones()
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Fetching Blocked Phones", err)
	}
	return c.JSON(http.StatusOK, types.GetBannedPhonesResponse{
		Phones: phones,
	})
}

func (h Handler) RemoveBannedPhone(c echo.Context) error {
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
	if account.Type != "admin" {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", errors.New("not admin"))
	}
	if c.Request().Method != echo.POST {
		return getAPIError(c, http.StatusBadRequest, "Invalid Method", nil)
	}
	var request types.ModifyBannedPhoneRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	err = h.validate.Struct(request)
	if len(request.Phone) < 10 || err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Phone Field", nil)
	}
	err = database.UnblockPhone(request.Phone)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Unblocking Phone", err)
	}
	return c.NoContent(http.StatusOK)
}

func (h Handler) AddBannedEmail(c echo.Context) error {
	if c.Request().Method != echo.POST {
		return getAPIError(c, http.StatusBadRequest, "Invalid Method", nil)
	}
	// No need for keys for any of these calls
	var request types.ModifyBannedEmailRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	err := h.validate.Struct(request)
	if err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Email Field", nil)
	}
	err = database.AddBlockedEmail(request.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Adding Email", err)
	}
	return c.NoContent(http.StatusOK)
}

func (h Handler) GetBannedEmails(c echo.Context) error {
	if c.Request().Method != echo.GET {
		return getAPIError(c, http.StatusBadRequest, "Invalid Method", nil)
	}
	emails, err := database.GetBlockedEmails()
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Fetching Blocked Emails", err)
	}
	return c.JSON(http.StatusOK, types.GetBannedEmailsResponse{
		Emails: emails,
	})
}

func (h Handler) RemoveBannedEmail(c echo.Context) error {
	if c.Request().Method != echo.POST {
		return getAPIError(c, http.StatusBadRequest, "Invalid Method", nil)
	}
	// No need for keys for any of these calls
	var request types.ModifyBannedEmailRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	err := h.validate.Struct(request)
	if err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Email Field", nil)
	}
	err = database.UnblockEmail(request.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Unblocking Email", err)
	}
	return c.NoContent(http.StatusOK)
}
