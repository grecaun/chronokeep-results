package handlers

import (
	"chronokeep/results/auth"
	"chronokeep/results/database"
	"chronokeep/results/types"
	"chronokeep/results/util"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

const (
	expirationWindow = time.Hour * 3
)

func (h Handler) GetAccount(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	var request types.GetAccountRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	if claims["type"].(string) != "admin" && claims["email"].(string) != request.Email {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	account, err := database.GetAccount(request.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if account == nil {
		return getAPIError(c, http.StatusNotFound, "Account Not Found", nil)
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
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	if claims["type"].(string) != "admin" {
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
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	if claims["type"].(string) != "admin" {
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
	account, err := database.AddAccount(request.Account)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Server Error", err)
	}
	return c.JSON(http.StatusOK, types.ModifyAccountResponse{
		Account: *account,
	})
}

func (h Handler) UpdateAccount(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	var request types.UpdateAccountRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	if claims["type"].(string) != "admin" && claims["email"].(string) != request.Account.Email {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	err := database.UpdateAccount(request.Account)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Server Error", err)
	}
	account, err := database.GetAccount(request.Account.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Server Error", err)
	}
	return c.JSON(http.StatusOK, types.ModifyAccountResponse{
		Account: *account,
	})
}

func (h Handler) DeleteAccount(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	var request types.DeleteAccountRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	if claims["type"].(string) != "admin" && claims["email"].(string) != request.Email {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	err := database.DeleteAccount(request.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Server Error", err)
	}
	return c.NoContent(http.StatusOK)
}

func (h Handler) ChangePassword(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	var request types.ChangePasswordRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	hashedPassword, err := auth.HashPassword(request.NewPassword)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Server Error", err)
	}
	// Let admins change passwords.
	if claims["type"].(string) == "admin" {
		err = database.ChangePassword(request.Email, hashedPassword)
		if err != nil {
			return getAPIError(c, http.StatusInternalServerError, "Server Error", err)
		}
		return c.NoContent(http.StatusOK)
	}
	// Verify we're changing a specific user's password.
	if claims["email"].(string) != request.Email {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	// Verify they knew their old password.
	account, err := database.GetAccount(request.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Server Error", err)
	}
	if account == nil {
		return getAPIError(c, http.StatusInternalServerError, "Account Not Found", nil)
	}
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
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	var request types.ChangeEmailRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	// Only let admins change passwords.
	if claims["type"].(string) != "admin" {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	err := database.ChangeEmail(request.OldEmail, request.NewEmail)
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
	config, err := util.GetConfig()
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Config Error", err)
	}
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = account.Name
	claims["email"] = account.Email
	claims["type"] = account.Type
	claims["exp"] = time.Now().Add(expirationWindow)

	t, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Token Generation Error", err)
	}
	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}
