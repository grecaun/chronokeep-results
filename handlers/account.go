package handlers

import (
	"chronokeep/results/auth"
	"chronokeep/results/types"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (h Handler) GetAccount(c echo.Context) error {
	var request types.GetAccountRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	// Get session.
	sess, err := session.Get("session", c)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Session Retrieval/Creation Error", err)
	}
	val, ok := sess.Values["email"]
	if !ok {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	email, ok := val.(string)
	if !ok {
		return getAPIError(c, http.StatusInternalServerError, "Invalid Session Value Type", nil)
	}
	account, err := database.GetAccount(email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	// Verify the account was found and they're either an admin or they own the account.
	if account == nil || (account.Type != "admin" && account.Email != request.Email) {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	keys, err := database.GetAccountKeys(account.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	events, err := database.GetAccountEvents(request.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	return c.JSON(http.StatusOK, types.GetAccountResponse{
		Account: *account,
		Keys:    keys,
		Events:  events,
	})
}

func (h Handler) GetAccounts(c echo.Context) error {
	// Get session.
	sess, err := session.Get("session", c)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Session Retrieval/Creation Error", err)
	}
	val, ok := sess.Values["email"]
	if !ok {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	email, ok := val.(string)
	if !ok {
		return getAPIError(c, http.StatusInternalServerError, "Invalid Session Value Type", nil)
	}
	account, err := database.GetAccount(email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	// Verify the account was found and they're an admin.
	if account == nil || account.Type != "admin" {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	accounts, err := database.GetAccounts()
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	return c.JSON(http.StatusOK, types.GetAllAccountsResponse{
		Accounts: accounts,
	})
}

func (h Handler) AddAccount(c echo.Context) error {
	// Get session.
	sess, err := session.Get("session", c)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Session Retrieval/Creation Error", err)
	}
	val, ok := sess.Values["email"]
	if !ok {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	email, ok := val.(string)
	if !ok {
		return getAPIError(c, http.StatusInternalServerError, "Invalid Session Value Type", nil)
	}
	account, err := database.GetAccount(email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	// Verify the account was found and they're an admin.
	if account == nil || account.Type != "admin" {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	var request types.AddAccountRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	password, err := auth.HashPassword(request.Password)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Server Error", err)
	}
	request.Account.Password = password
	account, err = database.AddAccount(request.Account)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Server Error", err)
	}
	return c.JSON(http.StatusOK, types.ModifyAccountResponse{
		Account: *account,
	})
}

func (h Handler) UpdateAccount(c echo.Context) error {
	var request types.UpdateAccountRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	// Get session.
	sess, err := session.Get("session", c)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Session Retrieval/Creation Error", err)
	}
	val, ok := sess.Values["email"]
	if !ok {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	email, ok := val.(string)
	if !ok {
		return getAPIError(c, http.StatusInternalServerError, "Invalid Session Value Type", nil)
	}
	account, err := database.GetAccount(email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	// Verify the account was found and they're either an admin or they own the account.
	if account == nil || (account.Type != "admin" && account.Email != request.Account.Email) {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	err = database.UpdateAccount(request.Account)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Server Error", err)
	}
	account, err = database.GetAccount(request.Account.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Server Error", err)
	}
	return c.JSON(http.StatusOK, types.ModifyAccountResponse{
		Account: *account,
	})
}

func (h Handler) DeleteAccount(c echo.Context) error {
	var request types.DeleteAccountRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	// Get session.
	sess, err := session.Get("session", c)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Session Retrieval/Creation Error", err)
	}
	val, ok := sess.Values["email"]
	if !ok {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	email, ok := val.(string)
	if !ok {
		return getAPIError(c, http.StatusInternalServerError, "Invalid Session Value Type", nil)
	}
	account, err := database.GetAccount(email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	// Verify the account was found and they're either an admin or they own the account.
	if account == nil || (account.Type != "admin" && account.Email != request.Email) {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	err = database.DeleteAccount(account.Identifier)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Server Error", err)
	}
	return c.NoContent(http.StatusOK)
}

func (h Handler) ChangePassword(c echo.Context) error {
	var request types.ChangePasswordRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	hashedPassword, err := auth.HashPassword(request.NewPassword)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Server Error", err)
	}
	// Get session.
	sess, err := session.Get("session", c)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Session Retrieval/Creation Error", err)
	}
	val, ok := sess.Values["email"]
	if !ok {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	email, ok := val.(string)
	if !ok {
		return getAPIError(c, http.StatusInternalServerError, "Invalid Session Value Type", nil)
	}
	account, err := database.GetAccount(email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	// Let admins change passwords.
	if account.Type == "admin" {
		err = database.ChangePassword(request.Email, hashedPassword)
		if err != nil {
			return getAPIError(c, http.StatusInternalServerError, "Server Error", err)
		}
		return c.NoContent(http.StatusOK)
	}
	// Verify we're changing a specific user's password.
	if account.Email != request.Email {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	// Verify they knew their old password.
	err = auth.VerifyPassword(account.Password, request.OldPassword)
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Invalid Credentials", err)
	}
	err = database.ChangePassword(request.Email, hashedPassword)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Server Error", err)
	}
	return c.NoContent(http.StatusOK)
}

func (h Handler) ChangeEmail(c echo.Context) error {
	var request types.ChangeEmailRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	// Get session.
	sess, err := session.Get("session", c)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Session Retrieval/Creation Error", err)
	}
	val, ok := sess.Values["email"]
	if !ok {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	email, ok := val.(string)
	if !ok {
		return getAPIError(c, http.StatusInternalServerError, "Invalid Session Value Type", nil)
	}
	account, err := database.GetAccount(email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	// Only let admins change emails.
	if account.Type != "admin" {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	err = database.ChangeEmail(request.OldEmail, request.NewEmail)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Server Error", err)
	}
	return c.NoContent(http.StatusOK)
}

func (h Handler) Login(c echo.Context) error {
	var request types.LoginRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	// Get User
	account, err := database.GetAccount(request.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Invalid Credentials", errors.New("user not found"))
	}
	// Check if account locked. Do this before verifying password.
	// If done after a bad actor could potentially figure out if they had a correct password by trying
	// even after it was locked until they received the locked message.
	if account.Locked {
		return getAPIError(c, http.StatusUnauthorized, "Account Locked", fmt.Errorf("account locked: %+v", account))
	}
	err = auth.VerifyPassword(account.Password, request.Password)
	if err != nil {
		database.InvalidPassword(*account)
		return getAPIError(c, http.StatusUnauthorized, "Invalid Credentials", err)
	}
	// Create/get session
	sess, err := session.Get("session", c)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Session Retrieval/Creation Error", err)
	}
	sess.Values["email"] = account.Email
	sess.Save(c.Request(), c.Response())
	return c.NoContent(http.StatusOK)
}
