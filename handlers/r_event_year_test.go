package handlers

import (
	"chronokeep/results/types"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRGetEventYears(t *testing.T) {
	// POST, /r/event-year
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	// Test empty auth header
	t.Log("Testing empty auth header.")
	request := httptest.NewRequest(http.MethodPost, "/r/event-year", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.RGetEventYears(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid auth header
	t.Log("Testing invalid auth header.")
	request = httptest.NewRequest(http.MethodPost, "/r/event-year", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEventYears(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid token
	t.Log("Testing invalid token.")
	request = httptest.NewRequest(http.MethodPost, "/r/event-year", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEventYears(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test unknown email
	t.Log("Testing unknown email in token.")
	token, _, err := createTokens("unknown@test.com")
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event-year", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEventYears(c)) {
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
	request = httptest.NewRequest(http.MethodPost, "/r/event-year", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEventYears(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test not logged in
	t.Log("Testing not logged in.")
	token, _, err = createTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event-year", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEventYears(c)) {
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
	request = httptest.NewRequest(http.MethodPost, "/r/event-year", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEventYears(c)) {
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
	request = httptest.NewRequest(http.MethodPost, "/r/event-year", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEventYears(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/r/event-year", strings.NewReader(string("////")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEventYears(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// test invalid content type
	t.Log("Testing invalid content type.")
	body, err := json.Marshal(types.GetEventYearRequest{
		Slug: "event1",
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event-year", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEventYears(c)) {
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
	body, err = json.Marshal(types.GetEventYearRequest{
		Slug: variables.events["event2"].Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event-year", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEventYears(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.EventYearsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(variables.eventYears["event2"]), len(resp.EventYears))
			for _, outer := range variables.eventYears["event2"] {
				found := false
				for _, inner := range resp.EventYears {
					if outer.Year == inner.Year {
						found = true
						assert.Equal(t, outer.DateTime.Local(), inner.DateTime.Local())
						assert.Equal(t, outer.Live, inner.Live)
					}
				}
				assert.True(t, found)
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
	body, err = json.Marshal(types.GetEventYearRequest{
		Slug: variables.events["event2"].Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event-year", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEventYears(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.EventYearsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(variables.eventYears["event2"]), len(resp.EventYears))
			for _, outer := range variables.eventYears["event2"] {
				found := false
				for _, inner := range resp.EventYears {
					if outer.Year == inner.Year {
						found = true
						assert.Equal(t, outer.DateTime.Local(), inner.DateTime.Local())
						assert.Equal(t, outer.Live, inner.Live)
					}
				}
				assert.True(t, found)
			}
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
	body, err = json.Marshal(types.GetEventYearRequest{
		Slug: variables.events["event1"].Slug,
		Year: variables.eventYears["event1"]["2020"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event-year", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEventYears(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// test unknown event name
	t.Log("Testing unknown event name.")
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
	body, err = json.Marshal(types.GetEventYearRequest{
		Slug: "unknown",
		Year: "2020",
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event-year", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEventYears(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// test known slug, unknown year
	t.Log("Testing unknown event year.")
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
	body, err = json.Marshal(types.GetEventYearRequest{
		Slug: variables.events["event1"].Slug,
		Year: "2000",
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event-year", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEventYears(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.EventYearsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(variables.eventYears["event1"]), len(resp.EventYears))
			for _, outer := range variables.eventYears["event1"] {
				found := false
				for _, inner := range resp.EventYears {
					if outer.Year == inner.Year {
						found = true
						assert.Equal(t, outer.DateTime.Local(), inner.DateTime.Local())
						assert.Equal(t, outer.Live, inner.Live)
					}
				}
				assert.True(t, found)
			}
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
	request = httptest.NewRequest(http.MethodPost, "/r/event-year", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetEventYears(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}

func TestRAddEventYear(t *testing.T) {
	// POST, /r/event-year/add
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	// Test empty auth header
	t.Log("Testing empty auth header.")
	request := httptest.NewRequest(http.MethodPost, "/r/event-year/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.RAddEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid auth header
	t.Log("Testing invalid auth header.")
	request = httptest.NewRequest(http.MethodPost, "/r/event-year/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid token
	t.Log("Testing invalid token.")
	request = httptest.NewRequest(http.MethodPost, "/r/event-year/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test unknown email
	t.Log("Testing unknown email in token.")
	token, _, err := createTokens("unknown@test.com")
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event-year/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEventYear(c)) {
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
	request = httptest.NewRequest(http.MethodPost, "/r/event-year/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test not logged in
	t.Log("Testing not logged in.")
	token, _, err = createTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event-year/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEventYear(c)) {
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
	request = httptest.NewRequest(http.MethodPost, "/r/event-year/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEventYear(c)) {
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
	request = httptest.NewRequest(http.MethodPost, "/r/event-year/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEventYear(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/r/event-year/add", strings.NewReader(string("////")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEventYear(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// test invalid content type
	t.Log("Testing invalid content type.")
	eYear := types.RequestYear{
		Year:     "2022",
		DateTime: "2022/04/06 09:00:00",
		Live:     false,
	}
	body, err := json.Marshal(types.ModifyEventYearRequest{
		Slug:      "event1",
		EventYear: eYear,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event-year/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEventYear(c)) {
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
	eYear = types.RequestYear{
		Year:     "2022",
		DateTime: "2022/04/06 9:01:15",
		Live:     false,
	}
	body, err = json.Marshal(types.ModifyEventYearRequest{
		Slug:      variables.events["event2"].Slug,
		EventYear: eYear,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event-year/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEventYear(c)) {
		if assert.Equal(t, http.StatusOK, response.Code) {
			var resp types.EventYearResponse
			if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
				assert.Equal(t, variables.events["event2"].Name, resp.Event.Name)
				assert.Equal(t, variables.events["event2"].AccessRestricted, resp.Event.AccessRestricted)
				assert.Equal(t, variables.events["event2"].ContactEmail, resp.Event.ContactEmail)
				assert.Equal(t, variables.events["event2"].Image, resp.Event.Image)
				assert.Equal(t, variables.events["event2"].Slug, resp.Event.Slug)
				assert.Equal(t, variables.events["event2"].Type, resp.Event.Type)
				assert.Equal(t, eYear.Year, resp.EventYear.Year)
				assert.Equal(t, eYear.Live, resp.EventYear.Live)
				assert.Equal(t, 2022, resp.EventYear.DateTime.Year())
				assert.Equal(t, time.Month(4), resp.EventYear.DateTime.Month())
				assert.Equal(t, 6, resp.EventYear.DateTime.Day())
				assert.Equal(t, 9, resp.EventYear.DateTime.Hour())
				assert.Equal(t, 1, resp.EventYear.DateTime.Minute())
				assert.Equal(t, 15, resp.EventYear.DateTime.Second())
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
	eYear = types.RequestYear{
		Year:     "2023",
		DateTime: "2023/07/12 6:03:02",
		Live:     true,
	}
	body, err = json.Marshal(types.ModifyEventYearRequest{
		Slug:      variables.events["event2"].Slug,
		EventYear: eYear,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event-year/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEventYear(c)) {
		if assert.Equal(t, http.StatusOK, response.Code) {
			var resp types.EventYearResponse
			if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
				assert.Equal(t, variables.events["event2"].Name, resp.Event.Name)
				assert.Equal(t, variables.events["event2"].AccessRestricted, resp.Event.AccessRestricted)
				assert.Equal(t, variables.events["event2"].ContactEmail, resp.Event.ContactEmail)
				assert.Equal(t, variables.events["event2"].Image, resp.Event.Image)
				assert.Equal(t, variables.events["event2"].Slug, resp.Event.Slug)
				assert.Equal(t, variables.events["event2"].Type, resp.Event.Type)
				assert.Equal(t, eYear.Year, resp.EventYear.Year)
				assert.Equal(t, eYear.Live, resp.EventYear.Live)
				assert.Equal(t, 2023, resp.EventYear.DateTime.Year())
				assert.Equal(t, time.Month(7), resp.EventYear.DateTime.Month())
				assert.Equal(t, 12, resp.EventYear.DateTime.Day())
				assert.Equal(t, 6, resp.EventYear.DateTime.Hour())
				assert.Equal(t, 3, resp.EventYear.DateTime.Minute())
				assert.Equal(t, 2, resp.EventYear.DateTime.Second())
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
	eYear = types.RequestYear{
		Year:     "2024",
		DateTime: "2024/07/12 06:03:02",
		Live:     false,
	}
	body, err = json.Marshal(types.ModifyEventYearRequest{
		Slug:      variables.events["event2"].Slug,
		EventYear: eYear,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event-year/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// duplicate
	t.Log("Testing duplicate.")
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
	eYear = types.RequestYear{
		Year:     "2022",
		DateTime: "2022/04/06 09:01:15",
		Live:     false,
	}
	body, err = json.Marshal(types.ModifyEventYearRequest{
		Slug:      variables.events["event2"].Slug,
		EventYear: eYear,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event-year/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEventYear(c)) {
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	}
	// test unknown event
	t.Log("Testing unknown event.")
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
	eYear = types.RequestYear{
		Year:     "2022",
		DateTime: "2022/04/06 09:01:15",
		Live:     false,
	}
	body, err = json.Marshal(types.ModifyEventYearRequest{
		Slug:      "unknown",
		EventYear: eYear,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event-year/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEventYear(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// validation -- year
	t.Log("Testing validation -- year.")
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
	eYear = types.RequestYear{
		Year:     "invalid",
		DateTime: "2022/04/06 09:01:15",
		Live:     false,
	}
	body, err = json.Marshal(types.ModifyEventYearRequest{
		Slug:      variables.events["event2"].Slug,
		EventYear: eYear,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event-year/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEventYear(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// validation -- date
	t.Log("Testing validation -- date.")
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
	eYear = types.RequestYear{
		Year:     "2052",
		DateTime: "2022/04 09:01:15",
		Live:     false,
	}
	body, err = json.Marshal(types.ModifyEventYearRequest{
		Slug:      variables.events["event2"].Slug,
		EventYear: eYear,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/event-year/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEventYear(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
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
	request = httptest.NewRequest(http.MethodPost, "/r/event-year/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}

func TestRUpdateEventYear(t *testing.T) {
	// PUT, /r/event-year/update
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	// Test empty auth header
	t.Log("Testing empty auth header.")
	request := httptest.NewRequest(http.MethodPut, "/r/event-year/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid auth header
	t.Log("Testing invalid auth header.")
	request = httptest.NewRequest(http.MethodPut, "/r/event-year/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid token
	t.Log("Testing invalid token.")
	request = httptest.NewRequest(http.MethodPut, "/r/event-year/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test unknown email
	t.Log("Testing unknown email in token.")
	token, _, err := createTokens("unknown@test.com")
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPut, "/r/event-year/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEventYear(c)) {
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
	request = httptest.NewRequest(http.MethodPut, "/r/event-year/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test not logged in
	t.Log("Testing not logged in.")
	token, _, err = createTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPut, "/r/event-year/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEventYear(c)) {
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
	request = httptest.NewRequest(http.MethodPut, "/r/event-year/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEventYear(c)) {
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
	request = httptest.NewRequest(http.MethodPut, "/r/event-year/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEventYear(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPut, "/r/event-year/update", strings.NewReader(string("////")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEventYear(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// test invalid content type
	t.Log("Testing invalid content type.")
	eYear := types.RequestYear{
		Year:     "2022",
		DateTime: "2022/04/06 09:00:00",
		Live:     false,
	}
	body, err := json.Marshal(types.ModifyEventYearRequest{
		Slug:      "event1",
		EventYear: eYear,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPut, "/r/event-year/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEventYear(c)) {
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
	eYear = types.RequestYear{
		Year:     variables.eventYears["event2"]["2020"].Year,
		DateTime: "2020/12/31 9:01:15",
		Live:     true,
	}
	body, err = json.Marshal(types.ModifyEventYearRequest{
		Slug:      variables.events["event2"].Slug,
		EventYear: eYear,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPut, "/r/event-year/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEventYear(c)) {
		if assert.Equal(t, http.StatusOK, response.Code) {
			var resp types.EventYearResponse
			if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
				assert.Equal(t, variables.events["event2"].Name, resp.Event.Name)
				assert.Equal(t, variables.events["event2"].AccessRestricted, resp.Event.AccessRestricted)
				assert.Equal(t, variables.events["event2"].ContactEmail, resp.Event.ContactEmail)
				assert.Equal(t, variables.events["event2"].Image, resp.Event.Image)
				assert.Equal(t, variables.events["event2"].Slug, resp.Event.Slug)
				assert.Equal(t, variables.events["event2"].Type, resp.Event.Type)
				assert.Equal(t, eYear.Year, resp.EventYear.Year)
				assert.Equal(t, eYear.Live, resp.EventYear.Live)
				assert.Equal(t, 2020, resp.EventYear.DateTime.Year())
				assert.Equal(t, time.Month(12), resp.EventYear.DateTime.Month())
				assert.Equal(t, 31, resp.EventYear.DateTime.Day())
				assert.Equal(t, 9, resp.EventYear.DateTime.Hour())
				assert.Equal(t, 1, resp.EventYear.DateTime.Minute())
				assert.Equal(t, 15, resp.EventYear.DateTime.Second())
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
	eYear = types.RequestYear{
		Year:     variables.eventYears["event2"]["2020"].Year,
		DateTime: "2023/07/12 6:03:02",
		Live:     false,
	}
	body, err = json.Marshal(types.ModifyEventYearRequest{
		Slug:      variables.events["event2"].Slug,
		EventYear: eYear,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPut, "/r/event-year/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEventYear(c)) {
		if assert.Equal(t, http.StatusOK, response.Code) {
			var resp types.EventYearResponse
			if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
				assert.Equal(t, variables.events["event2"].Name, resp.Event.Name)
				assert.Equal(t, variables.events["event2"].AccessRestricted, resp.Event.AccessRestricted)
				assert.Equal(t, variables.events["event2"].ContactEmail, resp.Event.ContactEmail)
				assert.Equal(t, variables.events["event2"].Image, resp.Event.Image)
				assert.Equal(t, variables.events["event2"].Slug, resp.Event.Slug)
				assert.Equal(t, variables.events["event2"].Type, resp.Event.Type)
				assert.Equal(t, eYear.Year, resp.EventYear.Year)
				assert.Equal(t, eYear.Live, resp.EventYear.Live)
				assert.Equal(t, 2023, resp.EventYear.DateTime.Year())
				assert.Equal(t, time.Month(7), resp.EventYear.DateTime.Month())
				assert.Equal(t, 12, resp.EventYear.DateTime.Day())
				assert.Equal(t, 6, resp.EventYear.DateTime.Hour())
				assert.Equal(t, 3, resp.EventYear.DateTime.Minute())
				assert.Equal(t, 2, resp.EventYear.DateTime.Second())
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
	eYear = types.RequestYear{
		Year:     variables.eventYears["event2"]["2020"].Year,
		DateTime: "2024/07/12 06:03:02",
		Live:     true,
	}
	body, err = json.Marshal(types.ModifyEventYearRequest{
		Slug:      variables.events["event2"].Slug,
		EventYear: eYear,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPut, "/r/event-year/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// test unknown event
	t.Log("Testing unknown event.")
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
	eYear = types.RequestYear{
		Year:     "2022",
		DateTime: "2022/04/06 09:01:15",
		Live:     false,
	}
	body, err = json.Marshal(types.ModifyEventYearRequest{
		Slug:      "unknown",
		EventYear: eYear,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPut, "/r/event-year/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEventYear(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// validation -- year (unknown)
	t.Log("Testing validation -- year.")
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
	eYear = types.RequestYear{
		Year:     "2024",
		DateTime: "2022/04/06 09:01:15",
		Live:     false,
	}
	body, err = json.Marshal(types.ModifyEventYearRequest{
		Slug:      variables.events["event2"].Slug,
		EventYear: eYear,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPut, "/r/event-year/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEventYear(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// validation -- date
	t.Log("Testing validation -- date.")
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
	eYear = types.RequestYear{
		Year:     variables.eventYears["event2"]["2020"].Year,
		DateTime: "2022/04 09:01:15",
		Live:     false,
	}
	body, err = json.Marshal(types.ModifyEventYearRequest{
		Slug:      variables.events["event2"].Slug,
		EventYear: eYear,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPut, "/r/event-year/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEventYear(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
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
	request = httptest.NewRequest(http.MethodPut, "/r/event-year/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}

func TestRDeleteEventYear(t *testing.T) {
	// DELETE, /r/event-year/delete
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	// Test empty auth header
	t.Log("Testing empty auth header.")
	request := httptest.NewRequest(http.MethodDelete, "/r/event-year/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid auth header
	t.Log("Testing invalid auth header.")
	request = httptest.NewRequest(http.MethodDelete, "/r/event-year/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid token
	t.Log("Testing invalid token.")
	request = httptest.NewRequest(http.MethodDelete, "/r/event-year/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test unknown email
	t.Log("Testing unknown email in token.")
	token, _, err := createTokens("unknown@test.com")
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/r/event-year/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEventYear(c)) {
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
	request = httptest.NewRequest(http.MethodDelete, "/r/event-year/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test not logged in
	t.Log("Testing not logged in.")
	token, _, err = createTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/r/event-year/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEventYear(c)) {
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
	request = httptest.NewRequest(http.MethodDelete, "/r/event-year/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEventYear(c)) {
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
	request = httptest.NewRequest(http.MethodDelete, "/r/event-year/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEventYear(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodDelete, "/r/event-year/delete", strings.NewReader(string("////")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEventYear(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// test invalid content type
	t.Log("Testing invalid content type.")
	body, err := json.Marshal(types.DeleteEventYearRequest{
		Slug: variables.events["event1"].Slug,
		Year: variables.eventYears["event1"]["2021"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/r/event-year/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEventYear(c)) {
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
	body, err = json.Marshal(types.DeleteEventYearRequest{
		Slug: variables.events["event2"].Slug,
		Year: variables.eventYears["event2"]["2021"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/r/event-year/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEventYear(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		nEv, err := database.GetEventYear(variables.events["event2"].Slug, variables.eventYears["event2"]["2021"].Year)
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
	eYear := types.EventYear{
		EventIdentifier: variables.events["event2"].Identifier,
		Year:            "2023",
		DateTime:        time.Date(2023, 12, 4, 6, 9, 1, 15, time.Local),
		Live:            true,
	}
	_, err = database.AddEventYear(eYear)
	if err != nil {
		t.Fatalf("Error adding new year to database: %v", err)
	}
	body, err = json.Marshal(types.DeleteEventYearRequest{
		Slug: variables.events["event2"].Slug,
		Year: "2023",
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/r/event-year/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEventYear(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		nEv, err := database.GetEventYear(variables.events["event2"].Slug, "2023")
		if assert.NoError(t, err) {
			assert.Nil(t, nEv)
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
	body, err = json.Marshal(types.DeleteEventYearRequest{
		Slug: variables.events["event2"].Slug,
		Year: "2020",
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/r/event-year/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// already deleted
	t.Log("Testing already deleted.")
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
	body, err = json.Marshal(types.DeleteEventYearRequest{
		Slug: variables.events["event2"].Slug,
		Year: variables.eventYears["event2"]["2021"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/r/event-year/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEventYear(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// test unknown event
	t.Log("Testing unknown event.")
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
	body, err = json.Marshal(types.DeleteEventYearRequest{
		Slug: "unknown",
		Year: variables.eventYears["event2"]["2021"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/r/event-year/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEventYear(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// test unknown year
	t.Log("Testing unknown year.")
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
	body, err = json.Marshal(types.DeleteEventYearRequest{
		Slug: variables.events["event2"].Slug,
		Year: "1999",
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/r/event-year/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEventYear(c)) {
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
	request = httptest.NewRequest(http.MethodDelete, "/r/event-year/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}
