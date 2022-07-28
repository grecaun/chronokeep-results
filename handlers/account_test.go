package handlers

import (
	db "chronokeep/results/database"
	"chronokeep/results/types"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func createExpiredTokens(email string) (*string, *string, error) {
	// Create token
	claims := jwt.MapClaims{}
	claims["email"] = email
	claims["authorized"] = true
	claims["exp"] = time.Now().Add(-1 * expirationWindow).Unix()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString([]byte(config.SecretKey))
	if err != nil {
		return nil, nil, err
	}
	// Create refresh token
	claims = jwt.MapClaims{}
	claims["email"] = email
	claims["exp"] = time.Now().Add(-1 * refreshWindow).Unix()
	r := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refresh, err := r.SignedString([]byte(config.RefreshKey))
	if err != nil {
		return nil, nil, err
	}
	return &token, &refresh, nil
}

func lockAccount(t *testing.T, email string, e *echo.Echo, h Handler) {
	for i := 0; i <= db.MaxLoginAttempts; i++ {
		t.Logf("Wrong password attempt %d", i+2)
		body, err := json.Marshal(types.LoginRequest{
			Email:    email,
			Password: "totally-not-a-password",
		})
		if err != nil {
			t.Fatalf("Error encoding request body into json object: %v", err)
		}
		request := httptest.NewRequest(http.MethodPost, "/r/account/login", strings.NewReader(string(body)))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)
		if assert.NoError(t, h.Login(c)) {
			assert.Equal(t, http.StatusUnauthorized, response.Code)
		}
	}
}

func TestLogin(t *testing.T) {
	// POST, /r/account/login (no auth header required)
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	// User not found
	t.Log("Testing user not found.")
	body, err := json.Marshal(types.LoginRequest{
		Email:    "wrong-email",
		Password: "totally-not-a-password",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/r/account/login", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.Login(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Incorrect password
	t.Log("Testing incorrect password.")
	body, err = json.Marshal(types.LoginRequest{
		Email:    variables.accounts[0].Email,
		Password: "totally-not-a-password",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/account/login", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.Login(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Invalid content type
	t.Log("Testing invalid content type.")
	body, err = json.Marshal(types.LoginRequest{
		Email:    variables.accounts[0].Email,
		Password: variables.testPassword1,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/account/login", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.Login(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Empty body
	t.Log("Testing empty request body.")
	request = httptest.NewRequest(http.MethodPost, "/r/account/login", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.Login(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Invalid body.
	t.Log("Testing invalid request body.")
	request = httptest.NewRequest(http.MethodPost, "/r/account/login", strings.NewReader(string("/////")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.Login(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Account locked
	lockAccount(t, variables.accounts[0].Email, e, h)
	t.Log("Testing account locked.")
	account, err := database.GetAccount(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error getting account from database: %v", err)
	}
	if assert.NotNil(t, account) {
		assert.Equal(t, true, account.Locked)
	}
	body, err = json.Marshal(types.LoginRequest{
		Email:    variables.accounts[0].Email,
		Password: variables.testPassword1,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/account/login", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.Login(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// valid request
	t.Log("Testing valid login.")
	body, err = json.Marshal(types.LoginRequest{
		Email:    variables.accounts[1].Email,
		Password: variables.testPassword1,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/account/login", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.Login(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp map[string]string
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			account, err = database.GetAccount(variables.accounts[1].Email)
			if assert.NoError(t, err) {
				assert.Equal(t, account.Token, resp["access_token"])
				assert.Equal(t, account.RefreshToken, resp["refresh_token"])
			}
		}
	}
}

func TestRefresh(t *testing.T) {
	// POST, /r/account/refresh (no auth header required)
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	// Test invalid body
	t.Log("Testing invalid body.")
	request := httptest.NewRequest(http.MethodPost, "/r/account/refresh", strings.NewReader(string("////")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.Refresh(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test empty body
	t.Log("Testing empty body.")
	request = httptest.NewRequest(http.MethodPost, "/r/account/refresh", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.Refresh(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid token
	t.Log("Testing invalid token.")
	body, err := json.Marshal(types.RefreshTokenRequest{
		RefreshToken: "invalid-refresh-token",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/account/refresh", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.Refresh(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test token with unknown email
	t.Log("Testing token with unknown email.")
	_, refresh, err := createTokens("invalid@test.com")
	if err != nil {
		t.Fatalf("Unable to create test tokens: %v", err)
	}
	body, err = json.Marshal(types.RefreshTokenRequest{
		RefreshToken: *refresh,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/account/refresh", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.Refresh(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired token
	t.Log("Testing expired token.")
	_, refresh, err = createExpiredTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Unable to create test tokens: %v", err)
	}
	body, err = json.Marshal(types.RefreshTokenRequest{
		RefreshToken: *refresh,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/account/refresh", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.Refresh(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test valid token but account logged out
	t.Log("Testing valid token but not set.")
	_, refresh, err = createTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Unable to create test tokens: %v", err)
	}
	body, err = json.Marshal(types.RefreshTokenRequest{
		RefreshToken: *refresh,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/account/refresh", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.Refresh(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test valid token
	t.Log("Testing valid token.")
	token, refresh, err := createTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Unable to create test tokens: %v", err)
	}
	account := variables.accounts[0]
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Unable to add test tokens to account: %v", err)
	}
	body, err = json.Marshal(types.RefreshTokenRequest{
		RefreshToken: *refresh,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/account/refresh", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.Refresh(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp map[string]string
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			account, err := database.GetAccount(variables.accounts[0].Email)
			if assert.NoError(t, err) {
				assert.Equal(t, account.Token, resp["access_token"])
				assert.Equal(t, account.RefreshToken, resp["refresh_token"])
			}
		}
	}
	// Test old refresh token
	t.Log("Testing old token.")
	body, err = json.Marshal(types.RefreshTokenRequest{
		RefreshToken: *refresh,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/account/refresh", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.Refresh(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test locked account
	lockAccount(t, variables.accounts[1].Email, e, h)
	t.Log("Testing account locked.")
	acc, err := database.GetAccount(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error getting account from database: %v", err)
	}
	if assert.NotNil(t, acc) {
		assert.Equal(t, true, acc.Locked)
	}
	token, refresh, err = createTokens(variables.accounts[1].Email)
	if err != nil {
		t.Fatalf("Unable to create test tokens: %v", err)
	}
	account = variables.accounts[1]
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Unable to add test tokens to account: %v", err)
	}
	body, err = json.Marshal(types.RefreshTokenRequest{
		RefreshToken: *refresh,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/account/refresh", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.Refresh(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}

func TestLogout(t *testing.T) {
	// POST, /r/account/logout
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	// Test empty auth header
	t.Log("Testing empty auth header.")
	request := httptest.NewRequest(http.MethodPost, "/r/account/logout", strings.NewReader(string("")))
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.Logout(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid auth header
	t.Log("Testing invalid auth header.")
	request = httptest.NewRequest(http.MethodPost, "/r/account/logout", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderAuthorization, "invalid-authorization-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.Logout(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid token
	t.Log("Testing invalid token.")
	request = httptest.NewRequest(http.MethodPost, "/r/account/logout", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderAuthorization, "Bearer invalid-authorization-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.Logout(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test unknown email
	t.Log("Testing unknown account.")
	token, _, err := createTokens("unknown@test.com")
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/account/logout", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.Logout(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired token
	t.Log("Testing expired token.")
	token, _, err = createExpiredTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/account/logout", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.Logout(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test valid
	t.Log("Testing valid token.")
	token, refresh, err := createTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	account, err := database.GetAccount(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error getting database for test: %v", err)
	}
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(*account)
	if err != nil {
		t.Fatalf("Error updating tokens on account for test: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/account/logout", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.Logout(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		account, err := database.GetAccount(variables.accounts[0].Email)
		if assert.NoError(t, err) {
			assert.Equal(t, "", account.Token)
			assert.Equal(t, "", account.RefreshToken)
		}
	}
	// Verify token no longer registered
	t.Log("Verifying token de-registered from account after logout.")
	request = httptest.NewRequest(http.MethodPost, "/r/account/logout", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.Logout(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}

func TestGetAccount(t *testing.T) {
	// POST, /r/account
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	// Test empty auth header
	t.Log("Testing empty auth header.")
	request := httptest.NewRequest(http.MethodPost, "/r/account", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.GetAccount(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid auth header
	t.Log("Testing invalid auth header.")
	request = httptest.NewRequest(http.MethodPost, "/r/account", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAccount(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid token
	t.Log("Testing invalid token.")
	request = httptest.NewRequest(http.MethodPost, "/r/account", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAccount(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test unknown email
	t.Log("Testing unknown email in token.")
	token, _, err := createTokens("unknown@test.com")
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/account", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAccount(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test not logged in
	t.Log("Testing not logged in.")
	token, _, err = createTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/account", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAccount(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired token
	t.Log("Testing expired token.")
	token, refresh, err := createExpiredTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	account := variables.accounts[0]
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/account", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAccount(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test locked account
	lockAccount(t, variables.accounts[0].Email, e, h)
	t.Log("Testing locked account.")
	token, refresh, err = createTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	account = variables.accounts[0]
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/account", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAccount(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	account.Locked = true
	err = database.UnlockAccount(account)
	if err != nil {
		t.Fatalf("Error unlocking account during test: %v", err)
	}
	// Test invalid body
	t.Log("Testing invalid body.")
	request = httptest.NewRequest(http.MethodPost, "/r/account", strings.NewReader(string("////")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAccount(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test invalid content type
	t.Log("Testing invalid content type.")
	request = httptest.NewRequest(http.MethodPost, "/r/account", strings.NewReader(string("////")))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAccount(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test empty body - gets token's account
	t.Log("Testing valid request - no body.")
	request = httptest.NewRequest(http.MethodPost, "/r/account", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAccount(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetAccountResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			acc, err := database.GetAccount(account.Email)
			if assert.NoError(t, err) {
				assert.Equal(t, acc.Email, resp.Account.Email)
				assert.Equal(t, acc.Locked, resp.Account.Locked)
				assert.Equal(t, acc.Name, resp.Account.Name)
				assert.Equal(t, acc.Type, resp.Account.Type)
			}
		}
	}
	// Test email - admin, authorized
	t.Log("Testing email in body, admin.")
	body, err := json.Marshal(types.GetAccountRequest{
		Email: &variables.accounts[1].Email,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/account", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAccount(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetAccountResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, variables.accounts[1].Email, resp.Account.Email)
			assert.Equal(t, variables.accounts[1].Locked, resp.Account.Locked)
			assert.Equal(t, variables.accounts[1].Name, resp.Account.Name)
			assert.Equal(t, variables.accounts[1].Type, resp.Account.Type)
		}
	}
	// Test email - admin, authorized - unknown email
	t.Log("Testing email in body, admin, unknown email.")
	email := "unknown@test.com"
	body, err = json.Marshal(types.GetAccountRequest{
		Email: &email,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/account", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAccount(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test email - not admin, unauthorized
	t.Log("Testing email in body, not admin.")
	body, err = json.Marshal(types.GetAccountRequest{
		Email: &variables.accounts[1].Email,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	account = variables.accounts[2]
	token, refresh, err = createTokens(account.Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating tokens for test: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/account", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAccount(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}

func TestGetAccounts(t *testing.T) {
	// GET, /r/account/all
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	// Test empty auth header
	t.Log("Testing empty auth header.")
	request := httptest.NewRequest(http.MethodGet, "/r/account/all", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.GetAccounts(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid auth header
	t.Log("Testing invalid auth header.")
	request = httptest.NewRequest(http.MethodGet, "/r/account/all", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAccounts(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid token
	t.Log("Testing invalid token.")
	request = httptest.NewRequest(http.MethodGet, "/r/account/all", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAccounts(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test unknown email
	t.Log("Testing unknown email in token.")
	token, _, err := createTokens("unknown@test.com")
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodGet, "/r/account/all", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAccounts(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test not logged in
	t.Log("Testing not logged in.")
	token, _, err = createTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodGet, "/r/account/all", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAccounts(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired token
	t.Log("Testing expired token.")
	token, refresh, err := createExpiredTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	account := variables.accounts[0]
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodGet, "/r/account/all", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAccounts(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test locked account
	lockAccount(t, variables.accounts[0].Email, e, h)
	t.Log("Testing locked account.")
	token, refresh, err = createTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	account = variables.accounts[0]
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodGet, "/r/account/all", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAccounts(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	account.Locked = true
	err = database.UnlockAccount(account)
	if err != nil {
		t.Fatalf("Error unlocking account during test: %v", err)
	}
	// Test not admin
	t.Log("Testing  not admin.")
	account = variables.accounts[2]
	token, refresh, err = createTokens(account.Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating tokens for test: %v", err)
	}
	request = httptest.NewRequest(http.MethodGet, "/r/account/all", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAccounts(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test admin
	t.Log("Testing admin.")
	token, refresh, err = createTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	account = variables.accounts[0]
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodGet, "/r/account/all", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAccounts(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetAllAccountsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(variables.accounts), len(resp.Accounts))
			for _, acc := range resp.Accounts {
				found := false
				for _, expected := range variables.accounts {
					if acc.Email == expected.Email {
						found = true
						assert.Equal(t, expected.Email, acc.Email)
						assert.Equal(t, expected.Locked, acc.Locked)
						assert.Equal(t, expected.Name, acc.Name)
						assert.Equal(t, expected.Type, acc.Type)
					}
				}
				assert.Equal(t, true, found)
			}
		}
	}
}

func TestAddAccount(t *testing.T) {
	// POST, /r/account/add
	// Test empty auth header
	// Test invalid auth header
	// Test invalid token
	// Test unknown email
	// Test expired token
	// Test not logged in
	// Test locked account

	// test empty request
	// test bad request
	// test admin
	// test email collision
	// test locked admin
	// test non-admin
	// validation checks
}

func TestUpdateAccount(t *testing.T) {
	// PUT, /r/account/update
	// Test empty auth header
	// Test invalid auth header
	// Test invalid token
	// Test unknown email
	// Test expired token
	// Test not logged in
	// Test locked account

	// test empty request
	// test bad request
	// test default
	// test admin, other account
	// test locked account
	// test non-admin other account
	// validation checks
}

func TestChangePassword(t *testing.T) {
	// PUT, /r/account/password
	// Test empty auth header
	// Test invalid auth header
	// Test invalid token
	// Test unknown email
	// Test expired token
	// Test not logged in
	// Test locked account

	// test empty request
	// test bad request
	// test self change - wrong password
	// test self change - correct password
	// test admin change
	// test admin, unknown email
	// test locked account
	// validation checks
}

func TestChangeEmail(t *testing.T) {
	// PUT, /r/account/email
	// Test empty auth header
	// Test invalid auth header
	// Test invalid token
	// Test unknown email
	// Test expired token
	// Test not logged in
	// Test locked account

	// test empty request
	// test bad request
	// test non-admin
	// test admin
	// test admin, unknown email
	// test locked
	// validation checks
}

func TestUnlock(t *testing.T) {
	// POST, /r/account/unlock
	// Test empty auth header
	// Test invalid auth header
	// Test invalid token
	// Test unknown email
	// Test expired token
	// Test not logged in
	// Test locked account

	// test empty request
	// test bad request
	// test non-admin
	// test admin
	// test locked
	// test unknown account
}

func TestDelete(t *testing.T) {
	// DELETE, /r/account/delete
	// Test empty auth header
	// Test invalid auth header
	// Test invalid token
	// Test unknown email
	// Test expired token
	// Test not logged in
	// Test locked account

	// test empty request
	// test bad request
	// test non-admin
	// test admin
	// test locked
	// test unknown account
}
