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

func TestRGetEvents(t *testing.T) {
	// POST, /r/event
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	// Test empty auth header
	t.Log("Testing empty auth header.")
	request := httptest.NewRequest(http.MethodPost, "/r/event", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.RGetEvents(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid auth header
	t.Log("Testing invalid auth header.")
	request = httptest.NewRequest(http.MethodPost, "/r/event", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEvents(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid token
	t.Log("Testing invalid token.")
	request = httptest.NewRequest(http.MethodPost, "/r/event", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEvents(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test unknown email
	t.Log("Testing unknown email in token.")
	token, _, err := createTokens("unknown@test.com")
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEvents(c)) {
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
	request = httptest.NewRequest(http.MethodPost, "/r/event", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEvents(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test not logged in
	t.Log("Testing not logged in.")
	token, _, err = createTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEvents(c)) {
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
	request = httptest.NewRequest(http.MethodPost, "/r/event", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEvents(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	account.Locked = true
	err = database.UnlockAccount(account)
	if err != nil {
		t.Fatalf("Error unlocking account during test: %v", err)
	}
	// test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/r/event", strings.NewReader(string("////")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEvents(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// test invalid content type
	t.Log("Testing invalid content type.")
	body, err := json.Marshal(types.GetREventsRequest{
		Email: &variables.accounts[0].Email,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEvents(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// test valid request - self
	t.Log("Testing valid request -- self add.")
	token, refresh, err = createTokens(variables.accounts[1].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	account = variables.accounts[1]
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEvents(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetEventsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			events, err := database.GetAccountEvents(variables.accounts[1].Email)
			if assert.NoError(t, err) {
				assert.Equal(t, len(events), len(resp.Events))
				for _, outer := range events {
					found := false
					for _, inner := range resp.Events {
						if outer.Slug == inner.Slug {
							found = true
							assert.Equal(t, outer.Name, inner.Name)
							assert.Equal(t, outer.AccessRestricted, inner.AccessRestricted)
							assert.Equal(t, outer.ContactEmail, inner.ContactEmail)
							assert.Equal(t, outer.Image, inner.Image)
							assert.Equal(t, outer.Type, inner.Type)
							assert.Equal(t, outer.Website, inner.Website)
						}
					}
					assert.True(t, found)
				}
			}
		}
	}
	// test valid request - admin for other
	t.Log("Testing valid request -- admin for other.")
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
	body, err = json.Marshal(types.GetREventsRequest{
		Email: &variables.accounts[1].Email,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEvents(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetEventsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			events, err := database.GetAccountEvents(variables.accounts[1].Email)
			if assert.NoError(t, err) {
				assert.Equal(t, len(events), len(resp.Events))
				for _, outer := range events {
					found := false
					for _, inner := range resp.Events {
						if outer.Slug == inner.Slug {
							found = true
							assert.Equal(t, outer.Name, inner.Name)
							assert.Equal(t, outer.AccessRestricted, inner.AccessRestricted)
							assert.Equal(t, outer.ContactEmail, inner.ContactEmail)
							assert.Equal(t, outer.Image, inner.Image)
							assert.Equal(t, outer.Type, inner.Type)
							assert.Equal(t, outer.Website, inner.Website)
						}
					}
					assert.True(t, found)
				}
			}
		}
	}
	// test invalid request - non-admin for other
	t.Log("Testing invalid request -- non-admin for other.")
	token, refresh, err = createTokens(variables.accounts[2].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	account = variables.accounts[2]
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	body, err = json.Marshal(types.GetREventsRequest{
		Email: &variables.accounts[1].Email,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEvents(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// test unknown email
	t.Log("Testing unknown email.")
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
	email := "unknown@test.com"
	body, err = json.Marshal(types.GetREventsRequest{
		Email: &email,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEvents(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetEventsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 0, len(resp.Events))
		}
	}
	// test token for wrong account //->//
	t.Log("Test token with embeded email not belonging to account it is attached to.")
	account = variables.accounts[0]
	account.Token = "totally-not-valid-token"
	account.RefreshToken = "not-a-valid-refresh-token"
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	token, refresh, err = createTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	account = variables.accounts[1]
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEvents(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}

func TestRAddEvent(t *testing.T) {
	// POST, /r/event/add
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	// Test empty auth header
	t.Log("Testing empty auth header.")
	request := httptest.NewRequest(http.MethodPost, "/r/event/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.RAddEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid auth header
	t.Log("Testing invalid auth header.")
	request = httptest.NewRequest(http.MethodPost, "/r/event/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid token
	t.Log("Testing invalid token.")
	request = httptest.NewRequest(http.MethodPost, "/r/event/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test unknown email
	t.Log("Testing unknown email in token.")
	token, _, err := createTokens("unknown@test.com")
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEvent(c)) {
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
	request = httptest.NewRequest(http.MethodPost, "/r/event/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test not logged in
	t.Log("Testing not logged in.")
	token, _, err = createTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEvent(c)) {
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
	request = httptest.NewRequest(http.MethodPost, "/r/event/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	account.Locked = true
	err = database.UnlockAccount(account)
	if err != nil {
		t.Fatalf("Error unlocking account during test: %v", err)
	}
	// test empty request
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
	t.Log("Testing empty request.")
	request = httptest.NewRequest(http.MethodPost, "/r/event/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/r/event/add", strings.NewReader(string("////")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// test invalid content type
	t.Log("Testing invalid content type.")
	event := types.Event{
		Name:             "Test Event 3",
		Slug:             "event3",
		ContactEmail:     "email@test.com",
		AccessRestricted: false,
		Type:             "distance",
	}
	body, err := json.Marshal(types.AddEventRequest{
		Event: event,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// test valid request - self
	t.Log("Testing valid request -- self add.")
	token, refresh, err = createTokens(variables.accounts[1].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	account = variables.accounts[1]
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	event = types.Event{
		Name:             "Test Event 5",
		Slug:             "event5",
		ContactEmail:     "email@test.com",
		AccessRestricted: false,
		Type:             "distance",
	}
	body, err = json.Marshal(types.AddEventRequest{
		Event: event,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEvent(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		nEv, err := database.GetEvent(event.Slug)
		if assert.NoError(t, err) && assert.NotNil(t, nEv) {
			assert.Equal(t, event.Name, nEv.Name)
			assert.Equal(t, event.Slug, nEv.Slug)
			assert.Equal(t, event.ContactEmail, nEv.ContactEmail)
			assert.Equal(t, event.AccessRestricted, nEv.AccessRestricted)
			assert.Equal(t, event.Type, nEv.Type)
			assert.Equal(t, event.Website, nEv.Website)
			assert.Equal(t, event.Image, nEv.Image)
		}
		var resp types.ModifyEventResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, event.Name, resp.Event.Name)
			assert.Equal(t, event.Slug, resp.Event.Slug)
			assert.Equal(t, event.ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, event.AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, event.Type, resp.Event.Type)
			assert.Equal(t, event.Website, resp.Event.Website)
			assert.Equal(t, event.Image, resp.Event.Image)
		}
	}
	// test valid request - admin for other
	t.Log("Testing valid request -- admin for other.")
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
	event = types.Event{
		Name:             "Test Event 6",
		Slug:             "event6",
		ContactEmail:     "email@test.com",
		AccessRestricted: false,
		Type:             "distance",
	}
	body, err = json.Marshal(types.AddEventRequest{
		Email: &variables.accounts[2].Email,
		Event: event,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEvent(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		nEv, err := database.GetEvent(event.Slug)
		if assert.NoError(t, err) {
			assert.Equal(t, event.Name, nEv.Name)
			assert.Equal(t, event.Slug, nEv.Slug)
			assert.Equal(t, event.ContactEmail, nEv.ContactEmail)
			assert.Equal(t, event.AccessRestricted, nEv.AccessRestricted)
			assert.Equal(t, event.Type, nEv.Type)
			assert.Equal(t, event.Website, nEv.Website)
			assert.Equal(t, event.Image, nEv.Image)
		}
		var resp types.ModifyEventResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, event.Name, resp.Event.Name)
			assert.Equal(t, event.Slug, resp.Event.Slug)
			assert.Equal(t, event.ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, event.AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, event.Type, resp.Event.Type)
			assert.Equal(t, event.Website, resp.Event.Website)
			assert.Equal(t, event.Image, resp.Event.Image)
		}
	}
	// test duplicate
	t.Log("Testing duplicate entry.")
	request = httptest.NewRequest(http.MethodPost, "/r/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEvent(c)) {
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	}
	// test invalid request - non-admin for other
	t.Log("Testing invalid request -- non-admin for other.")
	token, refresh, err = createTokens(variables.accounts[1].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	account = variables.accounts[1]
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	event = types.Event{
		Name:             "Test Event 9",
		Slug:             "event9",
		ContactEmail:     "email@test.com",
		AccessRestricted: false,
		Type:             "distance",
	}
	body, err = json.Marshal(types.AddEventRequest{
		Email: &variables.accounts[2].Email,
		Event: event,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// test unknown email
	t.Log("Testing unknown email.")
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
	email := "unknown@test.com"
	event = types.Event{
		Name:             "Test Event 4",
		Slug:             "event4",
		ContactEmail:     "email@test.com",
		AccessRestricted: false,
		Type:             "distance",
	}
	body, err = json.Marshal(types.AddEventRequest{
		Email: &email,
		Event: event,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEvent(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// validation -- invalid slug
	body, err = json.Marshal(types.AddEventRequest{
		Event: types.Event{
			Name:             "Test Event 4",
			Slug:             "event@#$7",
			ContactEmail:     "email@test.com",
			AccessRestricted: false,
			Type:             "distance",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing invalid slug.")
	request = httptest.NewRequest(http.MethodPost, "/r/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// validation -- no slug
	body, err = json.Marshal(types.AddEventRequest{
		Event: types.Event{
			Name:             "Test Event 4",
			Slug:             "",
			ContactEmail:     "email@test.com",
			AccessRestricted: false,
			Type:             "distance",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing no slug.")
	request = httptest.NewRequest(http.MethodPost, "/r/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// validation -- invalid name
	body, err = json.Marshal(types.AddEventRequest{
		Event: types.Event{
			Name:             "Test Event@#$ 4",
			Slug:             "event7",
			ContactEmail:     "email@test.com",
			AccessRestricted: false,
			Type:             "distance",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing invalid name.")
	request = httptest.NewRequest(http.MethodPost, "/r/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// validation -- no name
	body, err = json.Marshal(types.AddEventRequest{
		Event: types.Event{
			Name:             "",
			Slug:             "event7",
			ContactEmail:     "email@test.com",
			AccessRestricted: false,
			Type:             "distance",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing no name.")
	request = httptest.NewRequest(http.MethodPost, "/r/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// validation -- invalid type
	body, err = json.Marshal(types.AddEventRequest{
		Event: types.Event{
			Name:             "Test Event 4",
			Slug:             "event7",
			ContactEmail:     "email@test.com",
			AccessRestricted: false,
			Type:             "invalid",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing invalid type.")
	request = httptest.NewRequest(http.MethodPost, "/r/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// validation -- no type
	body, err = json.Marshal(types.AddEventRequest{
		Event: types.Event{
			Name:             "Test Event 4",
			Slug:             "event7",
			ContactEmail:     "email@test.com",
			AccessRestricted: false,
			Type:             "",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing no type.")
	request = httptest.NewRequest(http.MethodPost, "/r/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// validation -- invalid website
	body, err = json.Marshal(types.AddEventRequest{
		Event: types.Event{
			Name:             "Test Event 4",
			Slug:             "event7",
			ContactEmail:     "email@test.com",
			AccessRestricted: false,
			Type:             "distance",
			Website:          "//231a/a/d/c",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing invalid website.")
	request = httptest.NewRequest(http.MethodPost, "/r/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// validation -- invalid image
	body, err = json.Marshal(types.AddEventRequest{
		Event: types.Event{
			Name:             "Test Event 4",
			Slug:             "event7",
			ContactEmail:     "email@test.com",
			AccessRestricted: false,
			Type:             "distance",
			Image:            "/asd//a//czza//dda",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing invalid image.")
	request = httptest.NewRequest(http.MethodPost, "/r/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// validation -- invalid email
	body, err = json.Marshal(types.AddEventRequest{
		Event: types.Event{
			Name:             "Test Event 4",
			Slug:             "event7",
			ContactEmail:     "emailtest.com",
			AccessRestricted: false,
			Type:             "distance",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing invalid email.")
	request = httptest.NewRequest(http.MethodPost, "/r/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// test token for wrong account //->//
	t.Log("Test token with embeded email not belonging to account it is attached to.")
	account = variables.accounts[0]
	account.Token = "totally-not-valid-token"
	account.RefreshToken = "not-a-valid-refresh-token"
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	token, refresh, err = createTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	account = variables.accounts[1]
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}

func TestRUpdateEvent(t *testing.T) {
	// PUT, /r/event/update
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	// Test empty auth header
	t.Log("Testing empty auth header.")
	request := httptest.NewRequest(http.MethodPut, "/r/event/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid auth header
	t.Log("Testing invalid auth header.")
	request = httptest.NewRequest(http.MethodPut, "/r/event/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid token
	t.Log("Testing invalid token.")
	request = httptest.NewRequest(http.MethodPut, "/r/event/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test unknown email
	t.Log("Testing unknown email in token.")
	token, _, err := createTokens("unknown@test.com")
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPut, "/r/event/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEvent(c)) {
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
	request = httptest.NewRequest(http.MethodPut, "/r/event/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test not logged in
	t.Log("Testing not logged in.")
	token, _, err = createTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPut, "/r/event/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEvent(c)) {
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
	request = httptest.NewRequest(http.MethodPut, "/r/event/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	account.Locked = true
	err = database.UnlockAccount(account)
	if err != nil {
		t.Fatalf("Error unlocking account during test: %v", err)
	}
	// test empty request
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
	t.Log("Testing empty request.")
	request = httptest.NewRequest(http.MethodPut, "/r/event/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPut, "/r/event/update", strings.NewReader(string("////")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// test invalid content type
	t.Log("Testing invalid content type.")
	event := types.Event{
		Name:             "Test Event 3",
		Slug:             "event3",
		ContactEmail:     "email@test.com",
		AccessRestricted: false,
		Type:             "distance",
	}
	body, err := json.Marshal(types.UpdateEventRequest{
		Event: event,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPut, "/r/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// test valid request - self
	t.Log("Testing valid request -- self.")
	token, refresh, err = createTokens(variables.accounts[1].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	account = variables.accounts[1]
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	event = types.Event{
		Name:             "Updated Event Title 1",
		Slug:             variables.events["event2"].Slug,
		ContactEmail:     "email@test.com",
		AccessRestricted: false,
		Type:             "time",
		Image:            "http://google.com",
		Website:          "http://google.com",
	}
	body, err = json.Marshal(types.UpdateEventRequest{
		Event: event,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPut, "/r/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEvent(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		nEv, err := database.GetEvent(event.Slug)
		if assert.NoError(t, err) && assert.NotNil(t, nEv) {
			assert.Equal(t, event.Name, nEv.Name)
			assert.Equal(t, event.Slug, nEv.Slug)
			assert.Equal(t, event.ContactEmail, nEv.ContactEmail)
			assert.Equal(t, event.AccessRestricted, nEv.AccessRestricted)
			assert.Equal(t, event.Type, nEv.Type)
			assert.Equal(t, event.Website, nEv.Website)
			assert.Equal(t, event.Image, nEv.Image)
		}
		var resp types.ModifyEventResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, event.Name, resp.Event.Name)
			assert.Equal(t, event.Slug, resp.Event.Slug)
			assert.Equal(t, event.ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, event.AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, event.Type, resp.Event.Type)
			assert.Equal(t, event.Image, resp.Event.Image)
			assert.Equal(t, event.Website, resp.Event.Website)
		}
	}
	// test valid request - admin for other
	t.Log("Testing valid request -- admin for other.")
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
	event = types.Event{
		Name:             "Other Updated Event Title 1",
		Slug:             variables.events["event2"].Slug,
		ContactEmail:     "email2@test.com",
		AccessRestricted: true,
		Type:             "backyardultra",
	}
	body, err = json.Marshal(types.UpdateEventRequest{
		Event: event,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPut, "/r/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEvent(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		nEv, err := database.GetEvent(event.Slug)
		if assert.NoError(t, err) && assert.NotNil(t, nEv) {
			assert.Equal(t, event.Name, nEv.Name)
			assert.Equal(t, event.Slug, nEv.Slug)
			assert.Equal(t, event.ContactEmail, nEv.ContactEmail)
			assert.Equal(t, event.AccessRestricted, nEv.AccessRestricted)
			assert.Equal(t, event.Type, nEv.Type)
			assert.Equal(t, event.Website, nEv.Website)
			assert.Equal(t, event.Image, nEv.Image)
		}
		var resp types.ModifyEventResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, event.Name, resp.Event.Name)
			assert.Equal(t, event.Slug, resp.Event.Slug)
			assert.Equal(t, event.ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, event.AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, event.Type, resp.Event.Type)
			assert.Equal(t, event.Website, resp.Event.Website)
			assert.Equal(t, event.Image, resp.Event.Image)
		}
	}
	// test invalid request - non-admin for other
	t.Log("Testing invalid request -- non-admin for other.")
	token, refresh, err = createTokens(variables.accounts[2].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	account = variables.accounts[2]
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	event = types.Event{
		Name:             "iOther Updated Event Title 1",
		Slug:             variables.events["event2"].Slug,
		ContactEmail:     "iemail2@test.com",
		AccessRestricted: false,
		Type:             "distance",
	}
	body, err = json.Marshal(types.UpdateEventRequest{
		Event: event,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPut, "/r/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// test unknown event
	t.Log("Testing unknown event.")
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
	event = types.Event{
		Name:             "Test Event 4",
		Slug:             "unknownevent",
		ContactEmail:     "email@test.com",
		AccessRestricted: false,
		Type:             "distance",
	}
	body, err = json.Marshal(types.UpdateEventRequest{
		Event: event,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPut, "/r/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEvent(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// validation -- invalid slug
	body, err = json.Marshal(types.UpdateEventRequest{
		Event: types.Event{
			Name:             "Test Event 4",
			Slug:             "event@#$7",
			ContactEmail:     "email@test.com",
			AccessRestricted: false,
			Type:             "distance",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing invalid slug.")
	request = httptest.NewRequest(http.MethodPut, "/r/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// validation -- no slug
	body, err = json.Marshal(types.UpdateEventRequest{
		Event: types.Event{
			Name:             "Test Event 4",
			Slug:             "",
			ContactEmail:     "email@test.com",
			AccessRestricted: false,
			Type:             "distance",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing no slug.")
	request = httptest.NewRequest(http.MethodPut, "/r/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// validation -- invalid name
	body, err = json.Marshal(types.UpdateEventRequest{
		Event: types.Event{
			Name:             "Test Event@#$ 4",
			Slug:             "event7",
			ContactEmail:     "email@test.com",
			AccessRestricted: false,
			Type:             "distance",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing invalid name.")
	request = httptest.NewRequest(http.MethodPut, "/r/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// validation -- no name
	body, err = json.Marshal(types.UpdateEventRequest{
		Event: types.Event{
			Name:             "",
			Slug:             "event7",
			ContactEmail:     "email@test.com",
			AccessRestricted: false,
			Type:             "distance",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing no name.")
	request = httptest.NewRequest(http.MethodPut, "/r/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// validation -- invalid type
	body, err = json.Marshal(types.UpdateEventRequest{
		Event: types.Event{
			Name:             "Test Event 4",
			Slug:             "event7",
			ContactEmail:     "email@test.com",
			AccessRestricted: false,
			Type:             "invalid",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing invalid type.")
	request = httptest.NewRequest(http.MethodPut, "/r/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// validation -- no type
	body, err = json.Marshal(types.UpdateEventRequest{
		Event: types.Event{
			Name:             "Test Event 4",
			Slug:             "event7",
			ContactEmail:     "email@test.com",
			AccessRestricted: false,
			Type:             "",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing no type.")
	request = httptest.NewRequest(http.MethodPut, "/r/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// validation -- invalid website
	body, err = json.Marshal(types.UpdateEventRequest{
		Event: types.Event{
			Name:             "Test Event 4",
			Slug:             "event7",
			ContactEmail:     "email@test.com",
			AccessRestricted: false,
			Type:             "distance",
			Website:          "//231a/a/d/c",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing invalid website.")
	request = httptest.NewRequest(http.MethodPut, "/r/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// validation -- invalid image
	body, err = json.Marshal(types.UpdateEventRequest{
		Event: types.Event{
			Name:             "Test Event 4",
			Slug:             "event7",
			ContactEmail:     "email@test.com",
			AccessRestricted: false,
			Type:             "distance",
			Image:            "/asd//a//czza//dda",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing invalid image.")
	request = httptest.NewRequest(http.MethodPut, "/r/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// validation -- invalid email
	body, err = json.Marshal(types.UpdateEventRequest{
		Event: types.Event{
			Name:             "Test Event 4",
			Slug:             "event7",
			ContactEmail:     "emailtest.com",
			AccessRestricted: false,
			Type:             "distance",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing invalid email.")
	request = httptest.NewRequest(http.MethodPut, "/r/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// test token for wrong account //->//
	t.Log("Test token with embeded email not belonging to account it is attached to.")
	account = variables.accounts[0]
	account.Token = "totally-not-valid-token"
	account.RefreshToken = "not-a-valid-refresh-token"
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	token, refresh, err = createTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	account = variables.accounts[1]
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPut, "/r/event/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}

func TestRDeleteEvent(t *testing.T) {
	// DELETE, /r/event/delete
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	// Test empty auth header
	t.Log("Testing empty auth header.")
	request := httptest.NewRequest(http.MethodDelete, "/r/event/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid auth header
	t.Log("Testing invalid auth header.")
	request = httptest.NewRequest(http.MethodDelete, "/r/event/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid token
	t.Log("Testing invalid token.")
	request = httptest.NewRequest(http.MethodDelete, "/r/event/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test unknown email
	t.Log("Testing unknown email in token.")
	token, _, err := createTokens("unknown@test.com")
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/r/event/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEvent(c)) {
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
	request = httptest.NewRequest(http.MethodDelete, "/r/event/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test not logged in
	t.Log("Testing not logged in.")
	token, _, err = createTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/r/event/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEvent(c)) {
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
	request = httptest.NewRequest(http.MethodDelete, "/r/event/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	account.Locked = true
	err = database.UnlockAccount(account)
	if err != nil {
		t.Fatalf("Error unlocking account during test: %v", err)
	}
	// test empty request
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
	t.Log("Testing empty request.")
	request = httptest.NewRequest(http.MethodDelete, "/r/event/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodDelete, "/r/event/delete", strings.NewReader(string("////")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// test invalid content type
	t.Log("Testing invalid content type.")
	event := types.Event{
		AccountIdentifier: variables.accounts[1].Identifier,
		Name:              "Test Event 3",
		Slug:              "event3",
		ContactEmail:      "email@test.com",
		AccessRestricted:  false,
		Type:              "distance",
	}
	_, err = database.AddEvent(event)
	if err != nil {
		t.Fatalf("Error adding test event: %v", err)
	}
	body, err := json.Marshal(types.DeleteEventRequest{
		Slug: event.Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/r/event/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// test valid request - self
	t.Log("Testing valid request -- self.")
	token, refresh, err = createTokens(variables.accounts[1].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	account = variables.accounts[1]
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	event = types.Event{
		AccountIdentifier: variables.accounts[1].Identifier,
		Name:              "Test Event 5",
		Slug:              "event5",
		ContactEmail:      "email@test.com",
		AccessRestricted:  false,
		Type:              "distance",
	}
	_, err = database.AddEvent(event)
	if err != nil {
		t.Fatalf("Error adding test event: %v", err)
	}
	body, err = json.Marshal(types.DeleteEventRequest{
		Slug: event.Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/r/event/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEvent(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		nEv, err := database.GetEvent(event.Slug)
		if assert.NoError(t, err) {
			assert.Nil(t, nEv)
		}
	}
	// test valid request - admin for other
	t.Log("Testing valid request -- admin for other.")
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
	event = types.Event{
		AccountIdentifier: variables.accounts[2].Identifier,
		Name:              "Test Event 7",
		Slug:              "event7",
		ContactEmail:      "email@test.com",
		AccessRestricted:  false,
		Type:              "distance",
	}
	_, err = database.AddEvent(event)
	if err != nil {
		t.Fatalf("Error adding test event: %v", err)
	}
	body, err = json.Marshal(types.DeleteEventRequest{
		Slug: event.Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/r/event/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEvent(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		nEv, err := database.GetEvent(event.Slug)
		if assert.NoError(t, err) {
			assert.Nil(t, nEv)
		}
	}
	// test invalid request - non-admin for other
	t.Log("Testing invalid request -- non-admin for other.")
	token, refresh, err = createTokens(variables.accounts[1].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	account = variables.accounts[1]
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	event = types.Event{
		AccountIdentifier: variables.accounts[2].Identifier,
		Name:              "Test Event 9",
		Slug:              "event9",
		ContactEmail:      "email@test.com",
		AccessRestricted:  false,
		Type:              "distance",
	}
	_, err = database.AddEvent(event)
	if err != nil {
		t.Fatalf("Error adding test event: %v", err)
	}
	body, err = json.Marshal(types.DeleteEventRequest{
		Slug: event.Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/r/event/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// test unknown event
	t.Log("Testing unknown event.")
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
	body, err = json.Marshal(types.DeleteEventRequest{
		Slug: "unknownevent",
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/r/event/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEvent(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// test token for wrong account //->//
	t.Log("Test token with embeded email not belonging to account it is attached to.")
	account = variables.accounts[0]
	account.Token = "totally-not-valid-token"
	account.RefreshToken = "not-a-valid-refresh-token"
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	token, refresh, err = createTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	account = variables.accounts[1]
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/r/event/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}
