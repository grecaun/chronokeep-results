package handlers

import (
	"chronokeep/results/types"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAddBannedPhone(t *testing.T) {
	// POST
	_, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	// Test bad request
	body, err := json.Marshal(types.ModifyBannedPhoneRequest{
		Phone: "3525551478",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing bad request.")
	request := httptest.NewRequest(http.MethodPost, "/blocked/phones/add", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.AddBannedPhone(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test wrong method
	t.Log("Testing wrong method.")
	request = httptest.NewRequest(http.MethodGet, "/blocked/phones/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddBannedPhone(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test empty request
	t.Log("Testing empty request.")
	request = httptest.NewRequest(http.MethodPost, "/blocked/phones/add", strings.NewReader("{}"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddBannedPhone(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/blocked/phones/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddBannedPhone(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test non phone number
	t.Log("Testing non phone number.")
	body, err = json.Marshal(types.ModifyBannedPhoneRequest{
		Phone: "invalid phone",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/blocked/phones/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddBannedPhone(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test too short phone number
	t.Log("Testing too short phone number.")
	body, err = json.Marshal(types.ModifyBannedPhoneRequest{
		Phone: "1234",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/blocked/phones/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddBannedPhone(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test valid request
	t.Log("Testing valid request.")
	body, err = json.Marshal(types.ModifyBannedPhoneRequest{
		Phone: "3455551234",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/blocked/phones/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddBannedPhone(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	phones, err := database.GetBlockedPhones()
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(phones))
	}
	// Test second valid request
	t.Log("Testing second valid request.")
	body, err = json.Marshal(types.ModifyBannedPhoneRequest{
		Phone: "3455554534",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/blocked/phones/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddBannedPhone(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	phones, err = database.GetBlockedPhones()
	if assert.NoError(t, err) {
		assert.Equal(t, 2, len(phones))
	}
	// Test repeat valid request
	t.Log("Testing repeat valid request.")
	body, err = json.Marshal(types.ModifyBannedPhoneRequest{
		Phone: "3455551234",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/blocked/phones/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddBannedPhone(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	phones, err = database.GetBlockedPhones()
	if assert.NoError(t, err) {
		assert.Equal(t, 2, len(phones))
	}
}

func TestGetBannedPhones(t *testing.T) {
	// GET
	_, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	// Test wrong method
	t.Log("Testing wrong method.")
	request := httptest.NewRequest(http.MethodPost, "/blocked/phones/get", strings.NewReader(""))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.GetBannedPhones(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test valid request
	phones := []string{
		"3455554534",
		"3455551234",
	}
	err := database.AddBlockedPhone(phones[0])
	assert.NoError(t, err)
	t.Log("Testing valid request.")
	request = httptest.NewRequest(http.MethodGet, "/blocked/phones/get", strings.NewReader(""))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetBannedPhones(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetBannedPhonesResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 1, len(resp.Phones))
			assert.Equal(t, phones[0], resp.Phones[0])
		}
	}
	// Test second valid request
	err = database.AddBlockedPhones(phones)
	assert.NoError(t, err)
	request = httptest.NewRequest(http.MethodGet, "/blocked/phones/get", strings.NewReader(""))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetBannedPhones(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetBannedPhonesResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(phones), len(resp.Phones))
			found := 0
			for _, outer := range phones {
				for _, inner := range resp.Phones {
					if outer == inner {
						found += 1
					}
				}
			}
			assert.Equal(t, len(phones), found)
		}
	}
	// Test repeat valid request
	t.Log("Testing repeat valid request.")
	request = httptest.NewRequest(http.MethodGet, "/blocked/phones/get", strings.NewReader(""))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetBannedPhones(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	phones, err = database.GetBlockedPhones()
	if assert.NoError(t, err) {
		assert.Equal(t, 2, len(phones))
		var resp types.GetBannedPhonesResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(phones), len(resp.Phones))
			found := 0
			for _, outer := range phones {
				for _, inner := range resp.Phones {
					if outer == inner {
						found += 1
					}
				}
			}
			assert.Equal(t, len(phones), found)
		}
	}
}

func TestRemoveBannedPhone(t *testing.T) {
	// POST
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	// Test empty auth header
	t.Log("Testing empty auth header.")
	request := httptest.NewRequest(http.MethodPost, "/r/participants", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.RGetParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid auth header
	t.Log("Testing invalid auth header.")
	request = httptest.NewRequest(http.MethodPost, "/r/participants", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid token
	t.Log("Testing invalid token.")
	request = httptest.NewRequest(http.MethodPost, "/r/participants", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test unknown email
	t.Log("Testing unknown email in token.")
	token, _, err := createTokens("unknown@test.com")
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetParticipants(c)) {
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
	request = httptest.NewRequest(http.MethodPost, "/r/participants", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test not logged in
	t.Log("Testing not logged in.")
	token, _, err = createTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetParticipants(c)) {
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
	request = httptest.NewRequest(http.MethodPost, "/r/participants", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	account.Locked = true
	err = database.UnlockAccount(account)
	if err != nil {
		t.Fatalf("Error unlocking account during test: %v", err)
	}
	// Test bad request
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
	body, err := json.Marshal(types.ModifyBannedPhoneRequest{
		Phone: "3525551478",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/blocked/phones/unblock", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RemoveBannedPhone(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test wrong method
	t.Log("Testing wrong method.")
	request = httptest.NewRequest(http.MethodGet, "/blocked/phones/unblock", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RemoveBannedPhone(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test empty request
	t.Log("Testing empty request.")
	request = httptest.NewRequest(http.MethodPost, "/blocked/phones/unblock", strings.NewReader("{}"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RemoveBannedPhone(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/blocked/phones/unblock", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RemoveBannedPhone(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test non phone number
	t.Log("Testing non phone number.")
	body, err = json.Marshal(types.ModifyBannedPhoneRequest{
		Phone: "invalid phone",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/blocked/phones/unblock", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RemoveBannedPhone(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test too short phone number
	t.Log("Testing too short phone number.")
	body, err = json.Marshal(types.ModifyBannedPhoneRequest{
		Phone: "1234",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/blocked/phones/unblock", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RemoveBannedPhone(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test valid request
	err = database.AddBlockedPhones([]string{
		"3455554534",
		"3455551234",
	})
	assert.NoError(t, err)
	t.Log("Testing valid request.")
	body, err = json.Marshal(types.ModifyBannedPhoneRequest{
		Phone: "3455551234",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/blocked/phones/unblock", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RemoveBannedPhone(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	phones, err := database.GetBlockedPhones()
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(phones))
	}
	// Test second valid request
	t.Log("Testing second valid request.")
	body, err = json.Marshal(types.ModifyBannedPhoneRequest{
		Phone: "3455554534",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/blocked/phones/unblock", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RemoveBannedPhone(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	phones, err = database.GetBlockedPhones()
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(phones))
	}
	// Test repeat valid request
	t.Log("Testing repeat valid request.")
	body, err = json.Marshal(types.ModifyBannedPhoneRequest{
		Phone: "3455551234",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/blocked/phones/unblock", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RemoveBannedPhone(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	phones, err = database.GetBlockedPhones()
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(phones))
	}
}

func TestAddBannedEmail(t *testing.T) {
	// POST
	_, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	// Test bad request
	body, err := json.Marshal(types.ModifyBannedEmailRequest{
		Email: "admin@test.com",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into jsn object: %v", err)
	}
	t.Log("Testing bad request.")
	request := httptest.NewRequest(http.MethodPost, "/blocked/emails/add", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.AddBannedEmail(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test wrong method
	t.Log("Testing wrong method.")
	request = httptest.NewRequest(http.MethodGet, "/blocked/emails/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddBannedEmail(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test empty request
	t.Log("Testing empty request.")
	request = httptest.NewRequest(http.MethodPost, "/blocked/emails/add", strings.NewReader("{}"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddBannedEmail(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/blocked/emails/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddBannedEmail(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test invalid email
	t.Log("Testing invalid email.")
	body, err = json.Marshal(types.ModifyBannedEmailRequest{
		Email: "invalid email",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/blocked/emails/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddBannedEmail(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test valid request
	t.Log("Testing valid request.")
	body, err = json.Marshal(types.ModifyBannedEmailRequest{
		Email: "test@test.com",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/blocked/emails/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddBannedEmail(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	emails, err := database.GetBlockedEmails()
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(emails))
	}
	// Test second valid request
	t.Log("Testing second valid request.")
	body, err = json.Marshal(types.ModifyBannedEmailRequest{
		Email: "no2@test.com",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/blocked/emails/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddBannedEmail(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	emails, err = database.GetBlockedEmails()
	if assert.NoError(t, err) {
		assert.Equal(t, 2, len(emails))
	}
	// Test repeat valid request
	t.Log("Testing repeat valid request.")
	body, err = json.Marshal(types.ModifyBannedEmailRequest{
		Email: "test@test.com",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/blocked/emails/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddBannedEmail(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	emails, err = database.GetBlockedEmails()
	if assert.NoError(t, err) {
		assert.Equal(t, 2, len(emails))
	}
}

func TestGetBannedEmails(t *testing.T) {
	// GET
	_, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	// Test wrong method
	t.Log("Testing wrong method.")
	request := httptest.NewRequest(http.MethodPost, "/blocked/emails/get", strings.NewReader(""))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.GetBannedEmails(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test valid request
	emails := []string{
		"test@test.com",
		"no2@test.com",
	}
	err := database.AddBlockedEmail(emails[0])
	assert.NoError(t, err)
	t.Log("Testing valid request.")
	request = httptest.NewRequest(http.MethodGet, "/blocked/emails/get", strings.NewReader(""))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetBannedEmails(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetBannedEmailsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 1, len(resp.Emails))
			assert.Equal(t, emails[0], resp.Emails[0])
		}
	}
	// Test second valid request
	err = database.AddBlockedEmails(emails)
	assert.NoError(t, err)
	request = httptest.NewRequest(http.MethodGet, "/blocked/emails/get", strings.NewReader(""))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetBannedEmails(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetBannedEmailsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(emails), len(resp.Emails))
			found := 0
			for _, outer := range emails {
				for _, inner := range resp.Emails {
					if outer == inner {
						found += 1
					}
				}
			}
			assert.Equal(t, len(emails), found)
		}
	}
	// Test repeat valid request
	t.Log("Testing repeat valid request.")
	request = httptest.NewRequest(http.MethodGet, "/blocked/emails/get", strings.NewReader(""))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetBannedEmails(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	emails, err = database.GetBlockedEmails()
	if assert.NoError(t, err) {
		assert.Equal(t, 2, len(emails))
		var resp types.GetBannedEmailsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(emails), len(resp.Emails))
			found := 0
			for _, outer := range emails {
				for _, inner := range resp.Emails {
					if outer == inner {
						found += 1
					}
				}
			}
			assert.Equal(t, len(emails), found)
		}
	}
}

func TestRemoveBannedEmail(t *testing.T) {
	// POST
	_, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	// Test bad request
	body, err := json.Marshal(types.ModifyBannedEmailRequest{
		Email: "test@test.com",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing bad request.")
	request := httptest.NewRequest(http.MethodPost, "/blocked/emails/unblock", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.RemoveBannedEmail(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test wrong method
	t.Log("Testing wrong method.")
	request = httptest.NewRequest(http.MethodGet, "/blocked/emails/unblock", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RemoveBannedEmail(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test empty request
	t.Log("Testing empty request.")
	request = httptest.NewRequest(http.MethodPost, "/blocked/emails/unblock", strings.NewReader("{}"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RemoveBannedEmail(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/blocked/emails/unblock", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RemoveBannedEmail(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test invalid email
	t.Log("Testing invalid email.")
	body, err = json.Marshal(types.ModifyBannedEmailRequest{
		Email: "invalid email",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/blocked/emails/unblock", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RemoveBannedEmail(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test valid request
	err = database.AddBlockedEmails([]string{
		"test@test.com",
		"no2@test.com",
	})
	assert.NoError(t, err)
	t.Log("Testing valid request.")
	body, err = json.Marshal(types.ModifyBannedEmailRequest{
		Email: "test@test.com",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/blocked/emails/unblock", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RemoveBannedEmail(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	emails, err := database.GetBlockedEmails()
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(emails))
	}
	// Test second valid request
	t.Log("Testing second valid request.")
	body, err = json.Marshal(types.ModifyBannedEmailRequest{
		Email: "no2@test.com",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/blocked/emails/unblock", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RemoveBannedEmail(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	emails, err = database.GetBlockedEmails()
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(emails))
	}
	// Test repeat valid request
	t.Log("Testing repeat valid request.")
	body, err = json.Marshal(types.ModifyBannedEmailRequest{
		Email: "test@test.com",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/blocked/emails/unblock", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RemoveBannedEmail(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	emails, err = database.GetBlockedEmails()
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(emails))
	}
}
