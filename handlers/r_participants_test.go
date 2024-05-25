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

func TestRGetParticipants(t *testing.T) {
	// POST, /r/participants
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
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
	// Test empty request
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
	request = httptest.NewRequest(http.MethodPost, "/r/participants", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/r/participants", strings.NewReader(string("////")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test invalid content type
	t.Log("Testing invalid content type.")
	body, err := json.Marshal(types.GetParticipantsRequest{
		Slug: variables.events["event1"].Slug,
		Year: variables.eventYears["event1"]["2021"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test valid request - self
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
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: variables.eventYears["event2"]["2021"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetParticipants(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetParticipantsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, variables.events["event2"].Name, resp.Event.Name)
			assert.Equal(t, variables.events["event2"].Slug, resp.Event.Slug)
			assert.Equal(t, variables.events["event2"].Website, resp.Event.Website)
			assert.Equal(t, variables.events["event2"].Image, resp.Event.Image)
			assert.Equal(t, variables.events["event2"].ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, variables.events["event2"].AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, variables.events["event2"].Type, resp.Event.Type)
			assert.Equal(t, variables.events["event2"].RecentTime, resp.Event.RecentTime)
			assert.Equal(t, variables.eventYears["event2"]["2021"].Year, resp.Year.Year)
			assert.True(t, variables.eventYears["event2"]["2021"].DateTime.Equal(resp.Year.DateTime))
			assert.Equal(t, variables.eventYears["event2"]["2021"].Live, resp.Year.Live)
			assert.Equal(t, 4, len(resp.Participants))
		}
	}
	// Test valid request - admin for other
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
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: variables.eventYears["event2"]["2021"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetParticipants(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetParticipantsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, variables.events["event2"].Name, resp.Event.Name)
			assert.Equal(t, variables.events["event2"].Slug, resp.Event.Slug)
			assert.Equal(t, variables.events["event2"].Website, resp.Event.Website)
			assert.Equal(t, variables.events["event2"].Image, resp.Event.Image)
			assert.Equal(t, variables.events["event2"].ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, variables.events["event2"].AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, variables.events["event2"].Type, resp.Event.Type)
			assert.Equal(t, variables.events["event2"].RecentTime, resp.Event.RecentTime)
			assert.Equal(t, variables.eventYears["event2"]["2021"].Year, resp.Year.Year)
			assert.True(t, variables.eventYears["event2"]["2021"].DateTime.Equal(resp.Year.DateTime))
			assert.Equal(t, variables.eventYears["event2"]["2021"].Live, resp.Year.Live)
			assert.Equal(t, 4, len(resp.Participants))
		}
	}
	// Test invalid request - non-admin for other
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
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug: variables.events["event1"].Slug,
		Year: variables.eventYears["event1"]["2020"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test unknown event name
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
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug: "unknown",
		Year: "2020",
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetParticipants(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test known slug, unknown year
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
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug: variables.events["event1"].Slug,
		Year: "2000",
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetParticipants(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test token for wrong account //->//
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
	request = httptest.NewRequest(http.MethodPost, "/r/participants", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RGetParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}

func TestRAddParticipants(t *testing.T) {
	// POST, /r/participants/add
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	year, err := database.AddEventYear(types.EventYear{
		EventIdentifier: variables.events["event2"].Identifier,
		Year:            "2023",
		DateTime:        time.Date(2023, 04, 05, 9, 0, 0, 0, time.Local),
		Live:            false,
	})
	if err != nil {
		t.Fatalf("Error adding test event year to database: %v", err)
	}
	parts := []types.Participant{
		{
			AlternateId: "1024",
			Bib:         "1024",
			First:       "John",
			Last:        "Smnit",
			Birthdate:   "1/1/2004",
			Gender:      "M",
			AgeGroup:    "20-29",
			Distance:    "1 Mile",
			Anonymous:   false,
			SMSEnabled:  false,
			Mobile:      "",
			Apparel:     "",
		},
		{
			AlternateId: "2034",
			Bib:         "2034",
			First:       "Jason",
			Last:        "Jonson",
			Birthdate:   "1/1/1990",
			Gender:      "M",
			AgeGroup:    "30-39",
			Distance:    "5 Mile",
			Anonymous:   false,
			SMSEnabled:  false,
			Mobile:      "",
			Apparel:     "",
		},
		{
			AlternateId: "3521",
			Bib:         "3521",
			First:       "Rose",
			Last:        "McGowna",
			Birthdate:   "1/1/2008",
			Gender:      "F",
			AgeGroup:    "0-19",
			Distance:    "1 Mile",
			Anonymous:   false,
			SMSEnabled:  false,
			Mobile:      "",
			Apparel:     "",
		},
		{
			AlternateId: "1364",
			Bib:         "1364",
			First:       "Lilly",
			Last:        "Smith",
			Birthdate:   "1/1/2014",
			Gender:      "F",
			AgeGroup:    "0-19",
			Distance:    "1 Mile",
			Anonymous:   false,
			SMSEnabled:  false,
			Mobile:      "",
			Apparel:     "",
		},
	}
	body, err := json.Marshal(types.AddParticipantRequest{
		Slug:        variables.events["event2"].Slug,
		Year:        year.Year,
		Participant: parts[0],
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	// Test empty auth header
	t.Log("Testing empty auth header.")
	request := httptest.NewRequest(http.MethodPost, "/r/participants/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.RAddParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid auth header
	t.Log("Testing invalid auth header.")
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid token
	t.Log("Testing invalid token.")
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test unknown email
	t.Log("Testing unknown email in token.")
	token, _, err := createTokens("unknown@test.com")
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddParticipants(c)) {
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
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test not logged in
	t.Log("Testing not logged in.")
	token, _, err = createTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddParticipants(c)) {
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
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	account.Locked = true
	err = database.UnlockAccount(account)
	if err != nil {
		t.Fatalf("Error unlocking account during test: %v", err)
	}
	// Test empty request
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
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add", strings.NewReader(string("////")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test invalid content type
	t.Log("Testing invalid content type.")
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test valid request - self
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
	body, err = json.Marshal(types.AddParticipantRequest{
		Slug:        variables.events["event2"].Slug,
		Year:        year.Year,
		Participant: parts[0],
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddParticipants(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		if assert.NoError(t, err) {
			part, err := database.GetParticipants(year.Identifier)
			if assert.NoError(t, err) {
				outer := parts[0]
				found := false
				var resp types.UpdateParticipantResponse
				if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
					assert.True(t, outer.Equals(&resp.Participant))
					assert.Equal(t, outer.AlternateId, resp.Participant.AlternateId)
					assert.Equal(t, outer.Bib, resp.Participant.Bib)
					assert.Equal(t, outer.First, resp.Participant.First)
					assert.Equal(t, outer.Last, resp.Participant.Last)
					assert.Equal(t, outer.Birthdate, resp.Participant.Birthdate)
					assert.Equal(t, outer.Gender, resp.Participant.Gender)
					assert.Equal(t, outer.AgeGroup, resp.Participant.AgeGroup)
					assert.Equal(t, outer.Distance, resp.Participant.Distance)
					assert.Equal(t, outer.Anonymous, resp.Participant.Anonymous)
					assert.Equal(t, outer.SMSEnabled, resp.Participant.SMSEnabled)
					assert.Equal(t, outer.Mobile, resp.Participant.Mobile)
					assert.Equal(t, outer.Apparel, resp.Participant.Apparel)
				}
				for _, inner := range part {
					if outer.Bib == inner.Bib {
						assert.True(t, outer.Equals(&inner))
						assert.Equal(t, outer.AlternateId, inner.AlternateId)
						assert.Equal(t, outer.Bib, inner.Bib)
						assert.Equal(t, outer.First, inner.First)
						assert.Equal(t, outer.Last, inner.Last)
						assert.Equal(t, outer.Birthdate, inner.Birthdate)
						assert.Equal(t, outer.Gender, inner.Gender)
						assert.Equal(t, outer.AgeGroup, inner.AgeGroup)
						assert.Equal(t, outer.Distance, inner.Distance)
						assert.Equal(t, outer.Anonymous, inner.Anonymous)
						assert.Equal(t, outer.SMSEnabled, inner.SMSEnabled)
						assert.Equal(t, outer.Mobile, inner.Mobile)
						assert.Equal(t, outer.Apparel, inner.Apparel)
						found = true
					}
				}
				assert.True(t, found)
			}
		}
	}
	// Test valid request - admin for other
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
	body, _ = json.Marshal(types.AddParticipantRequest{
		Slug:        variables.events["event2"].Slug,
		Year:        year.Year,
		Participant: parts[1],
	})
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddParticipants(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		part, err := database.GetParticipants(year.Identifier)
		if assert.NoError(t, err) {
			outer := parts[1]
			found := false
			var resp types.UpdateParticipantResponse
			if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
				assert.True(t, outer.Equals(&resp.Participant))
				assert.Equal(t, outer.AlternateId, resp.Participant.AlternateId)
				assert.Equal(t, outer.Bib, resp.Participant.Bib)
				assert.Equal(t, outer.First, resp.Participant.First)
				assert.Equal(t, outer.Last, resp.Participant.Last)
				assert.Equal(t, outer.Birthdate, resp.Participant.Birthdate)
				assert.Equal(t, outer.Gender, resp.Participant.Gender)
				assert.Equal(t, outer.AgeGroup, resp.Participant.AgeGroup)
				assert.Equal(t, outer.Distance, resp.Participant.Distance)
				assert.Equal(t, outer.Anonymous, resp.Participant.Anonymous)
				assert.Equal(t, outer.SMSEnabled, resp.Participant.SMSEnabled)
				assert.Equal(t, outer.Mobile, resp.Participant.Mobile)
				assert.Equal(t, outer.Apparel, resp.Participant.Apparel)
			}
			for _, inner := range part {
				if outer.Bib == inner.Bib {
					assert.True(t, outer.Equals(&inner))
					assert.Equal(t, outer.AlternateId, inner.AlternateId)
					assert.Equal(t, outer.Bib, inner.Bib)
					assert.Equal(t, outer.First, inner.First)
					assert.Equal(t, outer.Last, inner.Last)
					assert.Equal(t, outer.Birthdate, inner.Birthdate)
					assert.Equal(t, outer.Gender, inner.Gender)
					assert.Equal(t, outer.AgeGroup, inner.AgeGroup)
					assert.Equal(t, outer.Distance, inner.Distance)
					assert.Equal(t, outer.Anonymous, inner.Anonymous)
					assert.Equal(t, outer.SMSEnabled, inner.SMSEnabled)
					assert.Equal(t, outer.Mobile, inner.Mobile)
					assert.Equal(t, outer.Apparel, inner.Apparel)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
	// Test invalid request - non-admin for other
	t.Log("Testing invalid request -- non-admin for other.")
	year.EventIdentifier = variables.events["event1"].Identifier
	year, err = database.AddEventYear(*year)
	if err != nil {
		t.Fatalf("Error adding event year to database: %v", err)
	}
	body, err = json.Marshal(types.AddParticipantRequest{
		Slug:        variables.events["event1"].Slug,
		Year:        year.Year,
		Participant: parts[0],
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
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
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test unknown event slug
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
	body, err = json.Marshal(types.AddParticipantRequest{
		Slug:        "unknown",
		Year:        "2020",
		Participant: parts[0],
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddParticipants(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test known slug, unknown year
	t.Log("Testing known slug, unknown year.")
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
	body, err = json.Marshal(types.AddParticipantRequest{
		Slug:        variables.events["event1"].Slug,
		Year:        "2025",
		Participant: parts[0],
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddParticipants(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test validation -- age
	t.Log("Testing validation -- age")
	body, err = json.Marshal(types.AddParticipantRequest{
		Slug: variables.events["event2"].Slug,
		Year: year.Year,
		Participant: types.Participant{
			Bib:        "10",
			First:      "John",
			Last:       "Jacob",
			Birthdate:  "1/1/3040",
			Gender:     "m",
			AgeGroup:   "10-20",
			Distance:   "1 Mile",
			Anonymous:  false,
			SMSEnabled: false,
			Mobile:     "",
			Apparel:    "",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test validation -- gender
	t.Log("Testing validation -- gender") // Gender is no longer validated to allow for unknown genders
	body, err = json.Marshal(types.AddParticipantRequest{
		Slug: variables.events["event2"].Slug,
		Year: year.Year,
		Participant: types.Participant{
			Bib:        "10",
			Birthdate:  "1/1/2004",
			AgeGroup:   "10-20",
			First:      "John",
			Last:       "Jacob",
			Distance:   "1 Mile",
			Gender:     "G",
			Anonymous:  false,
			SMSEnabled: false,
			Mobile:     "",
			Apparel:    "",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddParticipants(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	// Test validation -- distance
	t.Log("Testing validation -- distance")
	body, err = json.Marshal(types.AddParticipantRequest{
		Slug: variables.events["event2"].Slug,
		Year: year.Year,
		Participant: types.Participant{
			Bib:        "10",
			Birthdate:  "1/1/2004",
			AgeGroup:   "10-20",
			First:      "John",
			Last:       "Jacob",
			Distance:   "",
			Gender:     "M",
			Anonymous:  false,
			SMSEnabled: false,
			Mobile:     "",
			Apparel:    "",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test token for wrong account //->//
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
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}

func TestRDeleteParticipants(t *testing.T) {
	// POST, /r/participants/delete
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	year, err := database.AddEventYear(types.EventYear{
		EventIdentifier: variables.events["event2"].Identifier,
		Year:            "2023",
		DateTime:        time.Date(2023, 04, 05, 9, 0, 0, 0, time.Local),
		Live:            false,
	})
	if err != nil {
		t.Fatalf("Error adding test event year to database: %v", err)
	}
	parts := []types.Participant{
		{
			AlternateId: "1024",
			Bib:         "1024",
			First:       "John",
			Last:        "Smnit",
			Birthdate:   "1/1/2004",
			Gender:      "M",
			AgeGroup:    "20-29",
			Distance:    "1 Mile",
			Anonymous:   false,
			SMSEnabled:  false,
			Mobile:      "",
			Apparel:     "",
		},
		{
			AlternateId: "2034",
			Bib:         "2034",
			First:       "Jason",
			Last:        "Jonson",
			Birthdate:   "1/1/1990",
			Gender:      "M",
			AgeGroup:    "30-39",
			Distance:    "5 Mile",
			Anonymous:   false,
			SMSEnabled:  false,
			Mobile:      "",
			Apparel:     "",
		},
		{
			AlternateId: "3521",
			Bib:         "3521",
			First:       "Rose",
			Last:        "McGowna",
			Birthdate:   "1/1/2008",
			Gender:      "F",
			AgeGroup:    "0-19",
			Distance:    "1 Mile",
			Anonymous:   false,
			SMSEnabled:  false,
			Mobile:      "",
			Apparel:     "",
		},
		{
			AlternateId: "1364",
			Bib:         "1364",
			First:       "Lilly",
			Last:        "Smith",
			Birthdate:   "1/1/2014",
			Gender:      "F",
			AgeGroup:    "0-19",
			Distance:    "1 Mile",
			Anonymous:   false,
			SMSEnabled:  false,
			Mobile:      "",
			Apparel:     "",
		},
	}
	p, err := database.AddParticipants(year.Identifier, parts)
	if err != nil {
		t.Fatalf("Error adding participants to database for test: %v", err)
	}
	assert.Equal(t, len(parts), len(p))
	body, err := json.Marshal(types.DeleteParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: year.Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	// Test empty auth header
	t.Log("Testing empty auth header.")
	request := httptest.NewRequest(http.MethodDelete, "/r/participants/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid auth header
	t.Log("Testing invalid auth header.")
	request = httptest.NewRequest(http.MethodDelete, "/r/participants/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid token
	t.Log("Testing invalid token.")
	request = httptest.NewRequest(http.MethodDelete, "/r/participants/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test unknown email
	t.Log("Testing unknown email in token.")
	token, _, err := createTokens("unknown@test.com")
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/r/participants/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
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
	request = httptest.NewRequest(http.MethodDelete, "/r/participants/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test not logged in
	t.Log("Testing not logged in.")
	token, _, err = createTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/r/participants/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
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
	request = httptest.NewRequest(http.MethodDelete, "/r/participants/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	account.Locked = true
	err = database.UnlockAccount(account)
	if err != nil {
		t.Fatalf("Error unlocking account during test: %v", err)
	}
	// Test empty request
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
	request = httptest.NewRequest(http.MethodDelete, "/r/participants/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodDelete, "/r/participants/delete", strings.NewReader(string("////")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test invalid content type
	t.Log("Testing invalid content type.")
	request = httptest.NewRequest(http.MethodDelete, "/r/participants/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test valid request - self
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
	request = httptest.NewRequest(http.MethodDelete, "/r/participants/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		newP, err := database.GetParticipants(year.Identifier)
		if assert.NoError(t, err) {
			assert.Equal(t, 0, len(newP))
		}
	}
	// Test valid request - admin for other
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
	p, err = database.AddParticipants(year.Identifier, parts)
	if err != nil {
		t.Fatalf("Error adding participants to database for test: %v", err)
	}
	assert.Equal(t, len(parts), len(p))
	body, err = json.Marshal(types.DeleteParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: year.Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/r/participants/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		newP, err := database.GetParticipants(year.Identifier)
		if assert.NoError(t, err) {
			assert.Equal(t, 0, len(newP))
		}
	}
	// Test invalid request - non-admin for other
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
	body, err = json.Marshal(types.DeleteParticipantsRequest{
		Slug: variables.events["event1"].Slug,
		Year: variables.eventYears["event1"]["2021"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/r/participants/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test valid bibs
	t.Log("Testing valid bib based request.")
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
	p, err = database.AddParticipants(year.Identifier, parts)
	if err != nil {
		t.Fatalf("Error adding participants to database for test: %v", err)
	}
	assert.Equal(t, len(parts), len(p))
	idents := make([]string, 0)
	for _, person := range parts {
		idents = append(idents, person.AlternateId)
	}
	body, err = json.Marshal(types.DeleteParticipantsRequest{
		Slug:        variables.events["event2"].Slug,
		Year:        year.Year,
		Identifiers: idents[2:],
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/r/participants/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		newP, err := database.GetParticipants(year.Identifier)
		if assert.NoError(t, err) {
			assert.Equal(t, len(parts)-2, len(newP))
		}
	}
	// Test invalid bibs
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
	p, err = database.AddParticipants(year.Identifier, parts)
	if err != nil {
		t.Fatalf("Error adding participants to database for test: %v", err)
	}
	assert.Equal(t, len(parts), len(p))
	idents = []string{
		"invalid1",
		"invalid2",
	}
	body, err = json.Marshal(types.DeleteParticipantsRequest{
		Slug:        variables.events["event2"].Slug,
		Year:        year.Year,
		Identifiers: idents,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/r/participants/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		newP, err := database.GetParticipants(year.Identifier)
		if assert.NoError(t, err) {
			assert.Equal(t, len(parts), len(newP))
		}
	}
	// Test unknown event slug
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
	body, err = json.Marshal(types.DeleteParticipantsRequest{
		Slug: "unknown",
		Year: "2020",
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/r/participants/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test known slug, unknown year
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
	body, err = json.Marshal(types.DeleteParticipantsRequest{
		Slug: variables.events["event1"].Slug,
		Year: "2000",
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/r/participants/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test token for wrong account //->//
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
	request = httptest.NewRequest(http.MethodDelete, "/r/participants/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}

func TestRUpdateParticipants(t *testing.T) {
	// POST, /r/participants/update
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	// Test empty auth header
	t.Log("Testing empty auth header.")
	request := httptest.NewRequest(http.MethodPost, "/r/participants/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid auth header
	t.Log("Testing invalid auth header.")
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid token
	t.Log("Testing invalid token.")
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test unknown email
	t.Log("Testing unknown email in token.")
	token, _, err := createTokens("unknown@test.com")
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
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
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test not logged in
	t.Log("Testing not logged in.")
	token, _, err = createTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
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
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	account.Locked = true
	err = database.UnlockAccount(account)
	if err != nil {
		t.Fatalf("Error unlocking account during test: %v", err)
	}
	// Test empty request
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
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update", strings.NewReader(string("////")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test invalid content type
	t.Log("Testing invalid content type.")
	body, err := json.Marshal(types.GetParticipantsRequest{
		Slug: variables.events["event1"].Slug,
		Year: variables.eventYears["event1"]["2021"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test valid request - self
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
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: variables.eventYears["event2"]["2021"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	p, err := database.GetParticipants(variables.eventYears["event2"]["2021"].Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, 4, len(p))
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetParticipantsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 0, len(resp.Participants))
		}
	}
	// Test valid request - admin for other
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
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: variables.eventYears["event2"]["2020"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	p, err = database.GetParticipants(variables.eventYears["event2"]["2020"].Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, 4, len(p))
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetParticipantsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 0, len(resp.Participants))
		}
	}
	// Test invalid request - non-admin for other
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
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug: variables.events["event1"].Slug,
		Year: variables.eventYears["event1"]["2020"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test unknown event name
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
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug: "unknown",
		Year: "2020",
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test known slug, unknown year
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
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug: variables.events["event1"].Slug,
		Year: "2000",
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test token for wrong account //->//
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
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RDeleteParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}
