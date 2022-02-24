package handlers

import (
	"bytes"
	"chronokeep/results/types"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	log "github.com/sirupsen/logrus"
)

const (
	expirationWindow = time.Minute * 15
	refreshWindow    = time.Hour * 24 * 7
)

func (h Handler) GetAccount(c echo.Context) error {
	var request types.GetAccountRequest
	_ = c.Bind(&request)
	token, err := retrieveKey(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized Token", err)
	}
	requestBody, err := json.Marshal(types.AccountsGetAccountRequest{
		Ident: request.Ident,
		Token: *token,
	})
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Unable To Create JSON Request for Account Service", err)
	}
	resp, err := http.Post(config.AccountsURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Account Service Error", err)
	}
	defer resp.Body.Close()
	code := resp.StatusCode
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Reading Account Service Response", err)
	}
	var accountResponse types.AccountsSingleResponse
	json.Unmarshal(body, &accountResponse)
	if accountResponse.Message != nil {
		return getAPIError(c, code, *accountResponse.Message, nil)
	}
	if accountResponse.Account == nil && accountResponse.UserAccount == nil {
		return getAPIError(c, http.StatusInternalServerError, "Account Service Failed To Return Account", nil)
	}
	if code != 200 {
		return getAPIError(c, code, "Unusual Code Returned Without Error Message", nil)
	}
	account := accountResponse.Account
	userAccount := accountResponse.UserAccount
	// only allowe admins to view accounts not their own
	if userAccount != nil && account != nil && userAccount.Type != "admin" {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	unique := userAccount.Unique
	if account != nil {
		unique = account.Unique
	}
	// pull up the account we're giving information about
	outAccount, err := database.GetAccount(unique)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if outAccount == nil {
		return getAPIError(c, http.StatusNotFound, "Account Not Found", nil)
	}
	keys, err := database.GetAccountKeys(unique)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	events, err := database.GetAccountEvents(unique)
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
	var request types.GetAccountRequest
	_ = c.Bind(&request)
	token, err := retrieveKey(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized Token", err)
	}
	requestBody, err := json.Marshal(types.AccountsGetAccountRequest{
		Email: request.Email,
		Token: *token,
	})
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Unable To Create JSON Request for Account Service", err)
	}
	resp, err := http.Post(config.AccountsURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Account Service Error", err)
	}
	defer resp.Body.Close()
	code := resp.StatusCode
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Reading Account Service Response", err)
	}
	var accountResponse types.AccountsSingleResponse
	json.Unmarshal(body, &accountResponse)
	if accountResponse.Message != nil {
		return getAPIError(c, code, *accountResponse.Message, nil)
	}
	if accountResponse.Account == nil {
		return getAPIError(c, http.StatusInternalServerError, "Account Service Failed To Return Account", nil)
	}
	account := accountResponse.Account
	if !account.ResultsAPI {
		return getAPIError(c, http.StatusUnauthorized, "Account Not Authorized for Results API", nil)
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
	token, err := retrieveKey(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized Token", err)
	}
	requestBody, err := json.Marshal(types.AccountsAddAccountRequest{
		Token:    *token,
		Account:  request.Account,
		Password: request.Password,
	})
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Unable To Create JSON Request for Account Service", err)
	}
	resp, err := http.Post(config.AccountsURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Account Service Error", err)
	}
	defer resp.Body.Close()
	code := resp.StatusCode
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Reading Account Service Response", err)
	}
	var accountResponse types.AccountsSingleResponse
	json.Unmarshal(body, &accountResponse)
	if accountResponse.Message != nil {
		return getAPIError(c, code, *accountResponse.Message, nil)
	}
	if accountResponse.Account == nil {
		return getAPIError(c, http.StatusInternalServerError, "Account Service Failed To Return Account", nil)
	}
	if code != 200 {
		return getAPIError(c, code, "Incorrect Code Return Without Error Message", nil)
	}
	account, err := database.AddAccount(request.Account)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Unable To Add Account", err)
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
	token, err := retrieveKey(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized Token", err)
	}
	requestBody, err := json.Marshal(types.AccountsGetAccountRequest{
		Token: *token,
	})
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Unable To Create JSON Request for Account Service", err)
	}
	resp, err := http.Post(config.AccountsURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Account Service Error", err)
	}
	defer resp.Body.Close()
	code := resp.StatusCode
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Reading Account Service Response", err)
	}
	var accountResponse types.AccountsSingleResponse
	json.Unmarshal(body, &accountResponse)
	if accountResponse.Message != nil {
		return getAPIError(c, code, *accountResponse.Message, nil)
	}
	if accountResponse.Account == nil {
		return getAPIError(c, http.StatusInternalServerError, "Account Service Failed To Return Account", nil)
	}
	err = database.UpdateAccount(request.Account)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Unable To Update Account", err)
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
	var request types.DeleteAccountRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	token, err := retrieveKey(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized Token", err)
	}
	requestBody, err := json.Marshal(types.AccountsGetAccountRequest{
		Token: *token,
	})
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Unable To Create JSON Request for Account Service", err)
	}
	resp, err := http.Post(config.AccountsURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Account Service Error", err)
	}
	defer resp.Body.Close()
	code := resp.StatusCode
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Reading Account Service Response", err)
	}
	var accountResponse types.AccountsSingleResponse
	json.Unmarshal(body, &accountResponse)
	if accountResponse.Message != nil {
		return getAPIError(c, code, *accountResponse.Message, nil)
	}
	if accountResponse.Account == nil {
		return getAPIError(c, http.StatusInternalServerError, "Account Service Failed To Return Account", nil)
	}
	account := accountResponse.Account
	// Only admins and the owner of an account can delete it.
	if account.Type != "admin" && account.Identifier != request.Identifier {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	account, err = database.GetAccount(request.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if account == nil {
		return getAPIError(c, http.StatusNotFound, "Account Not Found", nil)
	}
	err = database.DeleteAccount(account.Identifier)
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
