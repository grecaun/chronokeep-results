package handlers

import (
	"chronokeep/results/auth"
	"chronokeep/results/types"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	log "github.com/sirupsen/logrus"
)

const (
	expirationWindow = time.Minute * 15
	refreshWindow    = time.Hour * 24 * 7
)

func (h Handler) GetAccount(c echo.Context) error {
	var request types.GetAccountRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	account, err := verifyToken(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized Token", err)
	}
	if account == nil || (account.Type != "admin" && account.Email != request.Email) {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	if account.Locked {
		return getAPIError(c, http.StatusUnauthorized, "Account Locked", nil)
	}
	account, err = database.GetAccount(request.Email)
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
	account, err := verifyToken(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized Token", err)
	}
	if account == nil || account.Type != "admin" {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	if account.Locked {
		return getAPIError(c, http.StatusUnauthorized, "Account Locked", nil)
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
	var request types.AddAccountRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	account, err := verifyToken(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized Token", err)
	}
	if account == nil || account.Type != "admin" {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	if account.Locked {
		return getAPIError(c, http.StatusUnauthorized, "Account Locked", nil)
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
	account, err := verifyToken(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized Token", err)
	}
	// Only admins and the owner of an account can update it.
	if account.Type != "admin" && account.Email != request.Account.Email {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	if account.Type != "admin" && account.Locked {
		account.RefreshToken = ""
		account.Token = ""
		err = database.UpdateTokens(*account)
		if err != nil {
			return getAPIError(c, http.StatusInternalServerError, "Database Error", nil)
		}
		return getAPIError(c, http.StatusUnauthorized, "Account Locked", nil)
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
	account, err := verifyToken(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized Token", err)
	}
	// Only admins and the owner of an account can delete it.
	if account.Type != "admin" && account.Email != request.Email {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	account, err = database.GetAccount(request.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
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
	account, err := verifyToken(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized Token", err)
	}
	if account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Account Not Found", nil)
	}
	if account.Locked {
		return getAPIError(c, http.StatusUnauthorized, "Account Locked", nil)
	}
	hashedPassword, err := auth.HashPassword(request.NewPassword)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Server Error", err)
	}
	// If the user is changing their own password they need to know their old password.
	if account.Email == request.Email {
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
		// Otherwise if an admin is changing a password for a user let them.
	} else if account.Type == "admin" {
		// Admin should log the person out when changing their password
		err = database.ChangePassword(request.Email, hashedPassword, true)
		if err != nil {
			return getAPIError(c, http.StatusInternalServerError, "Server Error", err)
		}
		return c.NoContent(http.StatusOK)
	}
	// Not their own account and not an admin, unauthorized.
	return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
}

func (h Handler) ChangeEmail(c echo.Context) error {
	var request types.ChangeEmailRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	account, err := verifyToken(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized Token", err)
	}
	// Only let admins change emails.
	if account.Locked || account.Type != "admin" {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	err = database.ChangeEmail(request.OldEmail, request.NewEmail)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Server Error", err)
	}
	return c.NoContent(http.StatusOK)
}

func (h Handler) Login(c echo.Context) error {
	log.Info("Logging in.")
	var request types.LoginRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	log.Info("Bind success, getting account.")
	// Get User
	account, err := database.GetAccount(request.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Invalid Credentials", errors.New("user not found"))
	}
	log.Info("User found.")
	// Check if account locked. Do this before verifying password.
	// If done after a bad actor could potentially figure out if they had a correct password by trying
	// even after it was locked until they received the locked message.
	if account.Locked {
		return getAPIError(c, http.StatusUnauthorized, "Account Locked", fmt.Errorf("account locked: %+v", account))
	}
	log.Info("Verifying password.")
	err = auth.VerifyPassword(account.Password, request.Password)
	if err != nil {
		database.InvalidPassword(*account)
		return getAPIError(c, http.StatusUnauthorized, "Invalid Credentials", err)
	}
	err = database.ValidPassword(*account)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	log.Info("Generating tokens.")
	token, refresh, err := createTokens(account.Email)
	if err != nil || token == nil || refresh == nil {
		return getAPIError(c, http.StatusInternalServerError, "Token Generation Error", err)
	}
	log.Info("Updating tokens on account.")
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(*account)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	return c.JSON(http.StatusOK, map[string]string{
		"access_token":  *token,
		"refresh_token": *refresh,
	})
}

func (h Handler) Logout(c echo.Context) error {
	account, err := verifyToken(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized Token", err)
	}
	account.Token = ""
	account.RefreshToken = ""
	err = database.UpdateTokens(*account)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	return c.NoContent(http.StatusOK)
}

func (h Handler) Refresh(c echo.Context) error {
	request := types.RefreshTokenRequest{}
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	rtoken, err := jwt.Parse(request.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.RefreshKey), nil
	})
	// Probably expired or doesn't exist.
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", err)
	}
	// Check if valid or the claims are set
	claims, ok := rtoken.Claims.(jwt.MapClaims)
	if !ok || !rtoken.Valid {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", errors.New("token not valid or claims issue"))
	}
	// Valid not expired token.
	email, ok := claims["email"].(string)
	if !ok {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", errors.New("email not set in token"))
	}
	account, err := database.GetAccount(email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", errors.New("account not found"))
	}
	if account.Locked {
		account.RefreshToken = ""
		account.Token = ""
		err = database.UpdateTokens(*account)
		if err != nil {
			return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
		}
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", errors.New("account locked"))
	}
	// Verify the token matches, throw in an empty check for if the user logged out as well.
	if account.RefreshToken != request.RefreshToken || account.RefreshToken == "" {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", errors.New("refresh token does not match account token"))
	}
	token, refresh, err := createTokens(account.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Token Generation Error", err)
	}
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(*account)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	return c.JSON(http.StatusOK, map[string]string{
		"access_token":  *token,
		"refresh_token": *refresh,
	})
}

func (h Handler) Unlock(c echo.Context) error {
	var request types.DeleteAccountRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	account, err := verifyToken(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized Token", err)
	}
	// Only let admins unlock accounts.
	if account.Locked || account.Type != "admin" {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	toUnlock, err := database.GetAccount(request.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if toUnlock == nil {
		return getAPIError(c, http.StatusNotFound, "Account Not Found", nil)
	}
	err = database.UnlockAccount(*toUnlock)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Server Error", err)
	}
	return c.NoContent(http.StatusOK)
}

func verifyToken(r *http.Request) (*types.Account, error) {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) != 2 {
		return nil, errors.New("unknown authorization header")
	}
	token, err := jwt.Parse(strArr[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("claims not set or token is not valid")
	}
	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("email not found in token claims")
	}
	account, err := database.GetAccount(email)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, errors.New("account not found")
	}
	if account.Token != strArr[1] || account.Token == "" {
		return nil, errors.New("token no longer valid")
	}
	return account, nil
}

func createTokens(email string) (*string, *string, error) {
	// Create token
	claims := jwt.MapClaims{}
	claims["email"] = email
	claims["authorized"] = true
	claims["exp"] = time.Now().Add(expirationWindow).Unix()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString([]byte(config.SecretKey))
	if err != nil {
		return nil, nil, err
	}
	// Create refresh token
	claims = jwt.MapClaims{}
	claims["email"] = email
	claims["exp"] = time.Now().Add(refreshWindow).Unix()
	r := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refresh, err := r.SignedString([]byte(config.RefreshKey))
	if err != nil {
		return nil, nil, err
	}
	return &token, &refresh, nil
}
