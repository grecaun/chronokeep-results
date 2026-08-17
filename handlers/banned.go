/*
Chronokeep Desktop - Race Scoring Software
Copyright (C) 2026 James Sentinella

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package handlers

import (
	"chronokeep/results/types"
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
)

func (h Handler) AddBannedPhone(c *echo.Context) error {
	if c.Request().Method != http.MethodPost {
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

func (h Handler) GetBannedPhones(c *echo.Context) error {
	if c.Request().Method != http.MethodGet {
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

func (h Handler) RemoveBannedPhone(c *echo.Context) error {
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
	if c.Request().Method != http.MethodPost {
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

func (h Handler) AddBannedEmail(c *echo.Context) error {
	if c.Request().Method != http.MethodPost {
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

func (h Handler) GetBannedEmails(c *echo.Context) error {
	if c.Request().Method != http.MethodGet {
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

func (h Handler) RemoveBannedEmail(c *echo.Context) error {
	if c.Request().Method != http.MethodPost {
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

