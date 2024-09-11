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
	year := variables.eventYears["event1"]["2021"].Year
	body, err := json.Marshal(types.GetParticipantsRequest{
		Slug: variables.events["event1"].Slug,
		Year: &year,
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
	year = variables.eventYears["event2"]["2021"].Year
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: &year,
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
	// Test valid request - self & no year
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
	year = variables.eventYears["event2"]["2021"].Year
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: &year,
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
	year = variables.eventYears["event1"]["2020"].Year
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug: variables.events["event1"].Slug,
		Year: &year,
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
	year = "2020"
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug: "unknown",
		Year: &year,
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
	year = "2000"
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug: variables.events["event1"].Slug,
		Year: &year,
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
	// Test linked account
	t.Log("Testing valid request -- linked account.")
	err = database.LinkAccounts(variables.accounts[0], variables.accounts[4])
	assert.NoError(t, err)
	token, refresh, err = createTokens(variables.accounts[4].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	account = variables.accounts[4]
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug: variables.events["event1"].Slug,
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
			assert.Equal(t, variables.events["event1"].Name, resp.Event.Name)
			assert.Equal(t, variables.events["event1"].Slug, resp.Event.Slug)
			assert.Equal(t, variables.events["event1"].Website, resp.Event.Website)
			assert.Equal(t, variables.events["event1"].Image, resp.Event.Image)
			assert.Equal(t, variables.events["event1"].ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, variables.events["event1"].AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, variables.events["event1"].Type, resp.Event.Type)
			assert.Equal(t, variables.events["event1"].RecentTime, resp.Event.RecentTime)
			assert.Equal(t, variables.eventYears["event1"]["2021"].Year, resp.Year.Year)
			assert.True(t, variables.eventYears["event1"]["2021"].DateTime.Equal(resp.Year.DateTime))
			assert.Equal(t, variables.eventYears["event1"]["2021"].Live, resp.Year.Live)
			assert.Equal(t, 500, len(resp.Participants))
		}
	}
}

func TestRAddParticipant(t *testing.T) {
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
	year2, err := database.AddEventYear(types.EventYear{
		EventIdentifier: variables.events["event1"].Identifier,
		Year:            "2024",
		DateTime:        time.Date(2024, 04, 05, 9, 0, 0, 0, time.Local),
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
	if assert.NoError(t, h.RAddParticipant(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid auth header
	t.Log("Testing invalid auth header.")
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddParticipant(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid token
	t.Log("Testing invalid token.")
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddParticipant(c)) {
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
	if assert.NoError(t, h.RAddParticipant(c)) {
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
	if assert.NoError(t, h.RAddParticipant(c)) {
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
	if assert.NoError(t, h.RAddParticipant(c)) {
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
	if assert.NoError(t, h.RAddParticipant(c)) {
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
	if assert.NoError(t, h.RAddParticipant(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add", strings.NewReader(string("////")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddParticipant(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test invalid content type
	t.Log("Testing invalid content type.")
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddParticipant(c)) {
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
	if assert.NoError(t, h.RAddParticipant(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		if assert.NoError(t, err) {
			part, err := database.GetParticipants(year.Identifier, 0, 0)
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
	if assert.NoError(t, h.RAddParticipant(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		part, err := database.GetParticipants(year.Identifier, 0, 0)
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
	if assert.NoError(t, h.RAddParticipant(c)) {
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
	if assert.NoError(t, h.RAddParticipant(c)) {
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
	if assert.NoError(t, h.RAddParticipant(c)) {
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
	if assert.NoError(t, h.RAddParticipant(c)) {
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
	if assert.NoError(t, h.RAddParticipant(c)) {
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
	if assert.NoError(t, h.RAddParticipant(c)) {
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
	if assert.NoError(t, h.RAddParticipant(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test linked account
	t.Log("Testing valid request -- linked account.")
	err = database.LinkAccounts(variables.accounts[0], variables.accounts[4])
	assert.NoError(t, err)
	token, refresh, err = createTokens(variables.accounts[4].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	account = variables.accounts[4]
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	body, err = json.Marshal(types.AddParticipantRequest{
		Slug:        variables.events["event1"].Slug,
		Year:        year2.Year,
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
	if assert.NoError(t, h.RAddParticipant(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		if assert.NoError(t, err) {
			part, err := database.GetParticipants(year2.Identifier, 0, 0)
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
		newP, err := database.GetParticipants(year.Identifier, 0, 0)
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
		newP, err := database.GetParticipants(year.Identifier, 0, 0)
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
		newP, err := database.GetParticipants(year.Identifier, 0, 0)
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
		newP, err := database.GetParticipants(year.Identifier, 0, 0)
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

func TestRUpdateParticipant(t *testing.T) {
	// POST, /r/participants/update
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	parts := []types.Participant{
		{
			AlternateId: "1024",
			Bib:         "1024",
			First:       "John",
			Last:        "Smnit",
			Birthdate:   "1/1/2004",
			Gender:      "Man",
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
			Gender:      "Man",
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
			Gender:      "Woman",
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
			Gender:      "Woman",
			AgeGroup:    "0-19",
			Distance:    "1 Mile",
			Anonymous:   false,
			SMSEnabled:  false,
			Mobile:      "",
			Apparel:     "",
		},
	}
	_, err := database.AddParticipants(variables.eventYears["event2"]["2020"].Identifier, parts)
	if err != nil {
		t.Fatalf("Error adding participants to database: %v", err)
	}
	_, err = database.AddParticipants(variables.eventYears["event2"]["2021"].Identifier, parts)
	if err != nil {
		t.Fatalf("Error adding participants to database: %v", err)
	}
	_, err = database.AddParticipants(variables.eventYears["event1"]["2021"].Identifier, parts)
	if err != nil {
		t.Fatalf("Error adding participants to database: %v", err)
	}
	// Test empty auth header
	t.Log("Testing empty auth header.")
	request := httptest.NewRequest(http.MethodPost, "/r/participants/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateParticipant(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid auth header
	t.Log("Testing invalid auth header.")
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateParticipant(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid token
	t.Log("Testing invalid token.")
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateParticipant(c)) {
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
	if assert.NoError(t, h.RUpdateParticipant(c)) {
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
	if assert.NoError(t, h.RUpdateParticipant(c)) {
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
	if assert.NoError(t, h.RUpdateParticipant(c)) {
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
	if assert.NoError(t, h.RUpdateParticipant(c)) {
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
	if assert.NoError(t, h.RUpdateParticipant(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update", strings.NewReader(string("////")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateParticipant(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test invalid content type
	t.Log("Testing invalid content type.")
	body, err := json.Marshal(types.UpdateParticipantRequest{
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
	if assert.NoError(t, h.RUpdateParticipant(c)) {
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
	updated := parts[0]
	updated.Bib = "newbib"
	updated.AgeGroup = "testagegroup"
	updated.Anonymous = !updated.Anonymous
	updated.SMSEnabled = !updated.SMSEnabled
	updated.Apparel = "newapparel"
	updated.Birthdate = "new bd"
	updated.Distance = "testdist"
	updated.First = "uTom"
	updated.Last = "uSmith"
	updated.Gender = "Unkn"
	updated.Mobile = "notanum"
	body, err = json.Marshal(types.UpdateParticipantRequest{
		Slug:        variables.events["event2"].Slug,
		Year:        variables.eventYears["event2"]["2020"].Year,
		Participant: updated,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateParticipant(c)) {
		p, err := database.GetParticipants(variables.eventYears["event2"]["2020"].Identifier, 0, 0)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.UpdateParticipantResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, updated.AgeGroup, resp.Participant.AgeGroup)
			assert.Equal(t, updated.AlternateId, resp.Participant.AlternateId)
			assert.Equal(t, updated.Anonymous, resp.Participant.Anonymous)
			assert.Equal(t, updated.Apparel, resp.Participant.Apparel)
			assert.Equal(t, updated.Bib, resp.Participant.Bib)
			assert.Equal(t, updated.Birthdate, resp.Participant.Birthdate)
			assert.Equal(t, updated.Distance, resp.Participant.Distance)
			assert.Equal(t, updated.First, resp.Participant.First)
			assert.Equal(t, updated.Gender, resp.Participant.Gender)
			assert.Equal(t, updated.Last, resp.Participant.Last)
			assert.Equal(t, updated.Mobile, resp.Participant.Mobile)
			assert.Equal(t, updated.SMSEnabled, resp.Participant.SMSEnabled)
		}
		found := false
		for _, outer := range p {
			if updated.AlternateId == outer.AlternateId {
				assert.Equal(t, updated.AgeGroup, outer.AgeGroup)
				assert.Equal(t, updated.AlternateId, outer.AlternateId)
				assert.Equal(t, updated.Anonymous, outer.Anonymous)
				assert.Equal(t, updated.Apparel, outer.Apparel)
				assert.Equal(t, updated.Bib, outer.Bib)
				assert.Equal(t, updated.Birthdate, outer.Birthdate)
				assert.Equal(t, updated.Distance, outer.Distance)
				assert.Equal(t, updated.First, outer.First)
				assert.Equal(t, updated.Gender, outer.Gender)
				assert.Equal(t, updated.Last, outer.Last)
				assert.Equal(t, updated.Mobile, outer.Mobile)
				assert.Equal(t, updated.SMSEnabled, outer.SMSEnabled)
				found = true
				break
			}
		}
		assert.True(t, found)
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
	body, err = json.Marshal(types.UpdateParticipantRequest{
		Slug:        variables.events["event2"].Slug,
		Year:        variables.eventYears["event2"]["2021"].Year,
		Participant: updated,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateParticipant(c)) {
		p, err := database.GetParticipants(variables.eventYears["event2"]["2021"].Identifier, 0, 0)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.UpdateParticipantResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, updated.AgeGroup, resp.Participant.AgeGroup)
			assert.Equal(t, updated.AlternateId, resp.Participant.AlternateId)
			assert.Equal(t, updated.Anonymous, resp.Participant.Anonymous)
			assert.Equal(t, updated.Apparel, resp.Participant.Apparel)
			assert.Equal(t, updated.Bib, resp.Participant.Bib)
			assert.Equal(t, updated.Birthdate, resp.Participant.Birthdate)
			assert.Equal(t, updated.Distance, resp.Participant.Distance)
			assert.Equal(t, updated.First, resp.Participant.First)
			assert.Equal(t, updated.Gender, resp.Participant.Gender)
			assert.Equal(t, updated.Last, resp.Participant.Last)
			assert.Equal(t, updated.Mobile, resp.Participant.Mobile)
			assert.Equal(t, updated.SMSEnabled, resp.Participant.SMSEnabled)
		}
		found := false
		for _, outer := range p {
			if updated.AlternateId == outer.AlternateId {
				assert.Equal(t, updated.AgeGroup, outer.AgeGroup)
				assert.Equal(t, updated.AlternateId, outer.AlternateId)
				assert.Equal(t, updated.Anonymous, outer.Anonymous)
				assert.Equal(t, updated.Apparel, outer.Apparel)
				assert.Equal(t, updated.Bib, outer.Bib)
				assert.Equal(t, updated.Birthdate, outer.Birthdate)
				assert.Equal(t, updated.Distance, outer.Distance)
				assert.Equal(t, updated.First, outer.First)
				assert.Equal(t, updated.Gender, outer.Gender)
				assert.Equal(t, updated.Last, outer.Last)
				assert.Equal(t, updated.Mobile, outer.Mobile)
				assert.Equal(t, updated.SMSEnabled, outer.SMSEnabled)
				found = true
				break
			}
		}
		assert.True(t, found)
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
	body, err = json.Marshal(types.UpdateParticipantRequest{
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
	if assert.NoError(t, h.RUpdateParticipant(c)) {
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
	body, err = json.Marshal(types.UpdateParticipantRequest{
		Slug: "unknown",
		Year: "unknown",
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateParticipant(c)) {
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
	body, err = json.Marshal(types.UpdateParticipantRequest{
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
	if assert.NoError(t, h.RUpdateParticipant(c)) {
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
	if assert.NoError(t, h.RUpdateParticipant(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test valid request - self
	t.Log("Testing valid request -- self.")
	err = database.LinkAccounts(variables.accounts[0], variables.accounts[4])
	assert.NoError(t, err)
	token, refresh, err = createTokens(variables.accounts[4].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	account = variables.accounts[4]
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	updated = parts[0]
	updated.Bib = "newbib"
	updated.AgeGroup = "testagegroup"
	updated.Anonymous = !updated.Anonymous
	updated.SMSEnabled = !updated.SMSEnabled
	updated.Apparel = "newapparel"
	updated.Birthdate = "new bd"
	updated.Distance = "testdist"
	updated.First = "uTom"
	updated.Last = "uSmith"
	updated.Gender = "Unkn"
	updated.Mobile = "notanum"
	body, err = json.Marshal(types.UpdateParticipantRequest{
		Slug:        variables.events["event1"].Slug,
		Year:        variables.eventYears["event1"]["2021"].Year,
		Participant: updated,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateParticipant(c)) {
		p, err := database.GetParticipants(variables.eventYears["event1"]["2021"].Identifier, 0, 0)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.UpdateParticipantResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, updated.AgeGroup, resp.Participant.AgeGroup)
			assert.Equal(t, updated.AlternateId, resp.Participant.AlternateId)
			assert.Equal(t, updated.Anonymous, resp.Participant.Anonymous)
			assert.Equal(t, updated.Apparel, resp.Participant.Apparel)
			assert.Equal(t, updated.Bib, resp.Participant.Bib)
			assert.Equal(t, updated.Birthdate, resp.Participant.Birthdate)
			assert.Equal(t, updated.Distance, resp.Participant.Distance)
			assert.Equal(t, updated.First, resp.Participant.First)
			assert.Equal(t, updated.Gender, resp.Participant.Gender)
			assert.Equal(t, updated.Last, resp.Participant.Last)
			assert.Equal(t, updated.Mobile, resp.Participant.Mobile)
			assert.Equal(t, updated.SMSEnabled, resp.Participant.SMSEnabled)
		}
		found := false
		for _, outer := range p {
			if updated.AlternateId == outer.AlternateId {
				assert.Equal(t, updated.AgeGroup, outer.AgeGroup)
				assert.Equal(t, updated.AlternateId, outer.AlternateId)
				assert.Equal(t, updated.Anonymous, outer.Anonymous)
				assert.Equal(t, updated.Apparel, outer.Apparel)
				assert.Equal(t, updated.Bib, outer.Bib)
				assert.Equal(t, updated.Birthdate, outer.Birthdate)
				assert.Equal(t, updated.Distance, outer.Distance)
				assert.Equal(t, updated.First, outer.First)
				assert.Equal(t, updated.Gender, outer.Gender)
				assert.Equal(t, updated.Last, outer.Last)
				assert.Equal(t, updated.Mobile, outer.Mobile)
				assert.Equal(t, updated.SMSEnabled, outer.SMSEnabled)
				found = true
				break
			}
		}
		assert.True(t, found)
	}
}

func TestRUpdateManyParticipants(t *testing.T) {
	// POST, /r/participants/update-many
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	parts := []types.Participant{
		{
			AlternateId: "1024",
			Bib:         "1024",
			First:       "John",
			Last:        "Smnit",
			Birthdate:   "1/1/2004",
			Gender:      "Man",
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
			Gender:      "Man",
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
			Gender:      "Woman",
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
			Gender:      "Woman",
			AgeGroup:    "0-19",
			Distance:    "1 Mile",
			Anonymous:   false,
			SMSEnabled:  false,
			Mobile:      "",
			Apparel:     "",
		},
	}
	_, err := database.AddParticipants(variables.eventYears["event2"]["2020"].Identifier, parts)
	if err != nil {
		t.Fatalf("Error adding participants to database: %v", err)
	}
	_, err = database.AddParticipants(variables.eventYears["event2"]["2021"].Identifier, parts)
	if err != nil {
		t.Fatalf("Error adding participants to database: %v", err)
	}
	_, err = database.AddParticipants(variables.eventYears["event1"]["2021"].Identifier, parts)
	if err != nil {
		t.Fatalf("Error adding participants to database: %v", err)
	}
	// Test empty auth header
	t.Log("Testing empty auth header.")
	request := httptest.NewRequest(http.MethodPost, "/r/participants/update-many", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateManyParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid auth header
	t.Log("Testing invalid auth header.")
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update-many", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateManyParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid token
	t.Log("Testing invalid token.")
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update-many", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateManyParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test unknown email
	t.Log("Testing unknown email in token.")
	token, _, err := createTokens("unknown@test.com")
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update-many", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateManyParticipants(c)) {
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
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update-many", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateManyParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test not logged in
	t.Log("Testing not logged in.")
	token, _, err = createTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update-many", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateManyParticipants(c)) {
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
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update-many", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateManyParticipants(c)) {
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
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update-many", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateManyParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update-many", strings.NewReader(string("////")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateManyParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test invalid content type
	t.Log("Testing invalid content type.")
	body, err := json.Marshal(types.AddParticipantsRequest{
		Slug:         variables.events["event1"].Slug,
		Year:         variables.eventYears["event1"]["2021"].Year,
		Participants: make([]types.Participant, 0),
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update-many", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateManyParticipants(c)) {
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
	updated := parts[0]
	updated.Bib = "newbib"
	updated.AgeGroup = "testagegroup"
	updated.Anonymous = !updated.Anonymous
	updated.SMSEnabled = !updated.SMSEnabled
	updated.Apparel = "newapparel"
	updated.Birthdate = "1/1/1999"
	updated.Distance = "testdist"
	updated.First = "uTom"
	updated.Last = "uSmith"
	updated.Gender = "Unkn"
	updated.Mobile = "notanum"
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug:         variables.events["event2"].Slug,
		Year:         variables.eventYears["event2"]["2020"].Year,
		Participants: []types.Participant{updated},
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update-many", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateManyParticipants(c)) {
		p, err := database.GetParticipants(variables.eventYears["event2"]["2020"].Identifier, 0, 0)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.UpdateParticipantsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, updated.AgeGroup, resp.Participants[0].AgeGroup)
			assert.Equal(t, updated.AlternateId, resp.Participants[0].AlternateId)
			assert.Equal(t, updated.Anonymous, resp.Participants[0].Anonymous)
			assert.Equal(t, updated.Apparel, resp.Participants[0].Apparel)
			assert.Equal(t, updated.Bib, resp.Participants[0].Bib)
			assert.Equal(t, updated.Birthdate, resp.Participants[0].Birthdate)
			assert.Equal(t, updated.Distance, resp.Participants[0].Distance)
			assert.Equal(t, updated.First, resp.Participants[0].First)
			assert.Equal(t, updated.Gender, resp.Participants[0].Gender)
			assert.Equal(t, updated.Last, resp.Participants[0].Last)
			assert.Equal(t, updated.Mobile, resp.Participants[0].Mobile)
			assert.Equal(t, updated.SMSEnabled, resp.Participants[0].SMSEnabled)
		}
		found := false
		for _, outer := range p {
			if updated.AlternateId == outer.AlternateId {
				assert.Equal(t, updated.AgeGroup, outer.AgeGroup)
				assert.Equal(t, updated.AlternateId, outer.AlternateId)
				assert.Equal(t, updated.Anonymous, outer.Anonymous)
				assert.Equal(t, updated.Apparel, outer.Apparel)
				assert.Equal(t, updated.Bib, outer.Bib)
				assert.Equal(t, updated.Birthdate, outer.Birthdate)
				assert.Equal(t, updated.Distance, outer.Distance)
				assert.Equal(t, updated.First, outer.First)
				assert.Equal(t, updated.Gender, outer.Gender)
				assert.Equal(t, updated.Last, outer.Last)
				assert.Equal(t, updated.Mobile, outer.Mobile)
				assert.Equal(t, updated.SMSEnabled, outer.SMSEnabled)
				found = true
				break
			}
		}
		assert.True(t, found)
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
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug:         variables.events["event2"].Slug,
		Year:         variables.eventYears["event2"]["2021"].Year,
		Participants: []types.Participant{updated},
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update-many", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateManyParticipants(c)) {
		p, err := database.GetParticipants(variables.eventYears["event2"]["2021"].Identifier, 0, 0)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.UpdateParticipantsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, updated.AgeGroup, resp.Participants[0].AgeGroup)
			assert.Equal(t, updated.AlternateId, resp.Participants[0].AlternateId)
			assert.Equal(t, updated.Anonymous, resp.Participants[0].Anonymous)
			assert.Equal(t, updated.Apparel, resp.Participants[0].Apparel)
			assert.Equal(t, updated.Bib, resp.Participants[0].Bib)
			assert.Equal(t, updated.Birthdate, resp.Participants[0].Birthdate)
			assert.Equal(t, updated.Distance, resp.Participants[0].Distance)
			assert.Equal(t, updated.First, resp.Participants[0].First)
			assert.Equal(t, updated.Gender, resp.Participants[0].Gender)
			assert.Equal(t, updated.Last, resp.Participants[0].Last)
			assert.Equal(t, updated.Mobile, resp.Participants[0].Mobile)
			assert.Equal(t, updated.SMSEnabled, resp.Participants[0].SMSEnabled)
		}
		found := false
		for _, outer := range p {
			if updated.AlternateId == outer.AlternateId {
				assert.Equal(t, updated.AgeGroup, outer.AgeGroup)
				assert.Equal(t, updated.AlternateId, outer.AlternateId)
				assert.Equal(t, updated.Anonymous, outer.Anonymous)
				assert.Equal(t, updated.Apparel, outer.Apparel)
				assert.Equal(t, updated.Bib, outer.Bib)
				assert.Equal(t, updated.Birthdate, outer.Birthdate)
				assert.Equal(t, updated.Distance, outer.Distance)
				assert.Equal(t, updated.First, outer.First)
				assert.Equal(t, updated.Gender, outer.Gender)
				assert.Equal(t, updated.Last, outer.Last)
				assert.Equal(t, updated.Mobile, outer.Mobile)
				assert.Equal(t, updated.SMSEnabled, outer.SMSEnabled)
				found = true
				break
			}
		}
		assert.True(t, found)
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
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug: variables.events["event1"].Slug,
		Year: variables.eventYears["event1"]["2020"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update-many", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateManyParticipants(c)) {
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
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug: "unknown",
		Year: "unknown",
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update-many", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateManyParticipants(c)) {
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
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug: variables.events["event1"].Slug,
		Year: "2000",
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update-many", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateManyParticipants(c)) {
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
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update-many", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateManyParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test valid request - self
	t.Log("Testing valid request -- self.")
	err = database.LinkAccounts(variables.accounts[0], variables.accounts[4])
	assert.NoError(t, err)
	token, refresh, err = createTokens(variables.accounts[4].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	account = variables.accounts[4]
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	updated = parts[0]
	updated.Bib = "newbib"
	updated.AgeGroup = "testagegroup"
	updated.Anonymous = !updated.Anonymous
	updated.SMSEnabled = !updated.SMSEnabled
	updated.Apparel = "newapparel"
	updated.Birthdate = "12/12/1985"
	updated.Distance = "testdist"
	updated.First = "uTom"
	updated.Last = "uSmith"
	updated.Gender = "Unkn"
	updated.Mobile = "notanum"
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug:         variables.events["event1"].Slug,
		Year:         variables.eventYears["event1"]["2021"].Year,
		Participants: []types.Participant{updated},
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/update-many", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RUpdateManyParticipants(c)) {
		p, err := database.GetParticipants(variables.eventYears["event1"]["2021"].Identifier, 0, 0)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.UpdateParticipantsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, updated.AgeGroup, resp.Participants[0].AgeGroup)
			assert.Equal(t, updated.AlternateId, resp.Participants[0].AlternateId)
			assert.Equal(t, updated.Anonymous, resp.Participants[0].Anonymous)
			assert.Equal(t, updated.Apparel, resp.Participants[0].Apparel)
			assert.Equal(t, updated.Bib, resp.Participants[0].Bib)
			assert.Equal(t, updated.Birthdate, resp.Participants[0].Birthdate)
			assert.Equal(t, updated.Distance, resp.Participants[0].Distance)
			assert.Equal(t, updated.First, resp.Participants[0].First)
			assert.Equal(t, updated.Gender, resp.Participants[0].Gender)
			assert.Equal(t, updated.Last, resp.Participants[0].Last)
			assert.Equal(t, updated.Mobile, resp.Participants[0].Mobile)
			assert.Equal(t, updated.SMSEnabled, resp.Participants[0].SMSEnabled)
		}
		found := false
		for _, outer := range p {
			if updated.AlternateId == outer.AlternateId {
				assert.Equal(t, updated.AgeGroup, outer.AgeGroup)
				assert.Equal(t, updated.AlternateId, outer.AlternateId)
				assert.Equal(t, updated.Anonymous, outer.Anonymous)
				assert.Equal(t, updated.Apparel, outer.Apparel)
				assert.Equal(t, updated.Bib, outer.Bib)
				assert.Equal(t, updated.Birthdate, outer.Birthdate)
				assert.Equal(t, updated.Distance, outer.Distance)
				assert.Equal(t, updated.First, outer.First)
				assert.Equal(t, updated.Gender, outer.Gender)
				assert.Equal(t, updated.Last, outer.Last)
				assert.Equal(t, updated.Mobile, outer.Mobile)
				assert.Equal(t, updated.SMSEnabled, outer.SMSEnabled)
				found = true
				break
			}
		}
		assert.True(t, found)
	}
}

func TestRAddManyParticipantss(t *testing.T) {
	// POST, /r/participants/add-many
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
	year2, err := database.AddEventYear(types.EventYear{
		EventIdentifier: variables.events["event1"].Identifier,
		Year:            "2024",
		DateTime:        time.Date(2024, 04, 05, 9, 0, 0, 0, time.Local),
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
	body, err := json.Marshal(types.AddParticipantsRequest{
		Slug:         variables.events["event2"].Slug,
		Year:         year.Year,
		Participants: parts[0:1],
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	// Test empty auth header
	t.Log("Testing empty auth header.")
	request := httptest.NewRequest(http.MethodPost, "/r/participants/add-many", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.RAddManyParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid auth header
	t.Log("Testing invalid auth header.")
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add-many", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddManyParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid token
	t.Log("Testing invalid token.")
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add-many", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer invalid-bearer-token")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddManyParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test unknown email
	t.Log("Testing unknown email in token.")
	token, _, err := createTokens("unknown@test.com")
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add-many", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddManyParticipants(c)) {
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
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add-many", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddManyParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test not logged in
	t.Log("Testing not logged in.")
	token, _, err = createTokens(variables.accounts[0].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add-many", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddManyParticipants(c)) {
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
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add-many", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddManyParticipants(c)) {
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
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add-many", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddManyParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add-many", strings.NewReader(string("////")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddManyParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test invalid content type
	t.Log("Testing invalid content type.")
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add-many", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddManyParticipants(c)) {
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
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug:         variables.events["event2"].Slug,
		Year:         year.Year,
		Participants: parts[0:1],
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add-many", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddManyParticipants(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		if assert.NoError(t, err) {
			part, err := database.GetParticipants(year.Identifier, 0, 0)
			if assert.NoError(t, err) {
				outer := parts[0]
				found := false
				var resp types.UpdateParticipantsResponse
				if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
					assert.True(t, outer.Equals(&resp.Participants[0]))
					assert.Equal(t, outer.AlternateId, resp.Participants[0].AlternateId)
					assert.Equal(t, outer.Bib, resp.Participants[0].Bib)
					assert.Equal(t, outer.First, resp.Participants[0].First)
					assert.Equal(t, outer.Last, resp.Participants[0].Last)
					assert.Equal(t, outer.Birthdate, resp.Participants[0].Birthdate)
					assert.Equal(t, outer.Gender, resp.Participants[0].Gender)
					assert.Equal(t, outer.AgeGroup, resp.Participants[0].AgeGroup)
					assert.Equal(t, outer.Distance, resp.Participants[0].Distance)
					assert.Equal(t, outer.Anonymous, resp.Participants[0].Anonymous)
					assert.Equal(t, outer.SMSEnabled, resp.Participants[0].SMSEnabled)
					assert.Equal(t, outer.Mobile, resp.Participants[0].Mobile)
					assert.Equal(t, outer.Apparel, resp.Participants[0].Apparel)
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
	body, _ = json.Marshal(types.AddParticipantsRequest{
		Slug:         variables.events["event2"].Slug,
		Year:         year.Year,
		Participants: parts[1:3],
	})
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add-many", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddManyParticipants(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		part, err := database.GetParticipants(year.Identifier, 0, 0)
		if assert.NoError(t, err) {
			outer := parts[1]
			second := parts[2]
			found := false
			secFound := false
			var resp types.UpdateParticipantsResponse
			if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
				outerFound := false
				secondFound := false
				for _, part := range resp.Participants {
					if outer.Equals(&part) {
						outerFound = true
						assert.Equal(t, outer.AlternateId, part.AlternateId)
						assert.Equal(t, outer.Bib, part.Bib)
						assert.Equal(t, outer.First, part.First)
						assert.Equal(t, outer.Last, part.Last)
						assert.Equal(t, outer.Birthdate, part.Birthdate)
						assert.Equal(t, outer.Gender, part.Gender)
						assert.Equal(t, outer.AgeGroup, part.AgeGroup)
						assert.Equal(t, outer.Distance, part.Distance)
						assert.Equal(t, outer.Anonymous, part.Anonymous)
						assert.Equal(t, outer.SMSEnabled, part.SMSEnabled)
						assert.Equal(t, outer.Mobile, part.Mobile)
						assert.Equal(t, outer.Apparel, part.Apparel)
					} else if second.Equals(&part) {
						secondFound = true
						assert.Equal(t, second.AlternateId, part.AlternateId)
						assert.Equal(t, second.Bib, part.Bib)
						assert.Equal(t, second.First, part.First)
						assert.Equal(t, second.Last, part.Last)
						assert.Equal(t, second.Birthdate, part.Birthdate)
						assert.Equal(t, second.Gender, part.Gender)
						assert.Equal(t, second.AgeGroup, part.AgeGroup)
						assert.Equal(t, second.Distance, part.Distance)
						assert.Equal(t, second.Anonymous, part.Anonymous)
						assert.Equal(t, second.SMSEnabled, part.SMSEnabled)
						assert.Equal(t, second.Mobile, part.Mobile)
						assert.Equal(t, second.Apparel, part.Apparel)
					}
				}
				assert.True(t, outerFound)
				assert.True(t, secondFound)
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
				if second.Bib == inner.Bib {
					assert.True(t, second.Equals(&inner))
					assert.Equal(t, second.AlternateId, inner.AlternateId)
					assert.Equal(t, second.Bib, inner.Bib)
					assert.Equal(t, second.First, inner.First)
					assert.Equal(t, second.Last, inner.Last)
					assert.Equal(t, second.Birthdate, inner.Birthdate)
					assert.Equal(t, second.Gender, inner.Gender)
					assert.Equal(t, second.AgeGroup, inner.AgeGroup)
					assert.Equal(t, second.Distance, inner.Distance)
					assert.Equal(t, second.Anonymous, inner.Anonymous)
					assert.Equal(t, second.SMSEnabled, inner.SMSEnabled)
					assert.Equal(t, second.Mobile, inner.Mobile)
					assert.Equal(t, second.Apparel, inner.Apparel)
					secFound = true
				}
			}
			assert.True(t, found)
			assert.True(t, secFound)
		}
	}
	// Test invalid request - non-admin for other
	t.Log("Testing invalid request -- non-admin for other.")
	year.EventIdentifier = variables.events["event1"].Identifier
	year, err = database.AddEventYear(*year)
	if err != nil {
		t.Fatalf("Error adding event year to database: %v", err)
	}
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug:         variables.events["event1"].Slug,
		Year:         year.Year,
		Participants: parts[0:1],
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
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add-many", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddManyParticipants(c)) {
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
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug:         "unknown",
		Year:         "2020",
		Participants: parts[0:1],
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add-many", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddManyParticipants(c)) {
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
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug:         variables.events["event1"].Slug,
		Year:         "2025",
		Participants: parts[0:1],
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add-many", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddManyParticipants(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test validation -- age
	t.Log("Testing validation -- age")
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: year.Year,
		Participants: []types.Participant{{
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
		}},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add-many", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddManyParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test validation -- gender
	t.Log("Testing validation -- gender") // Gender is no longer validated to allow for unknown genders
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: year.Year,
		Participants: []types.Participant{{
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
		}},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add-many", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddManyParticipants(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	// Test validation -- distance
	t.Log("Testing validation -- distance")
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: year.Year,
		Participants: []types.Participant{{
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
		}},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add-many", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddManyParticipants(c)) {
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
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add-many", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddManyParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test linked account
	t.Log("Testing valid request -- linked account.")
	err = database.LinkAccounts(variables.accounts[0], variables.accounts[4])
	assert.NoError(t, err)
	token, refresh, err = createTokens(variables.accounts[4].Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	account = variables.accounts[4]
	account.Token = *token
	account.RefreshToken = *refresh
	err = database.UpdateTokens(account)
	if err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug:         variables.events["event1"].Slug,
		Year:         year2.Year,
		Participants: parts[0:1],
	})
	if err != nil {
		t.Fatalf("Error encoding request into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add-many", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddManyParticipants(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		if assert.NoError(t, err) {
			part, err := database.GetParticipants(year2.Identifier, 0, 0)
			if assert.NoError(t, err) {
				outer := parts[0]
				found := false
				var resp types.UpdateParticipantsResponse
				if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
					assert.True(t, outer.Equals(&resp.Participants[0]))
					assert.Equal(t, outer.AlternateId, resp.Participants[0].AlternateId)
					assert.Equal(t, outer.Bib, resp.Participants[0].Bib)
					assert.Equal(t, outer.First, resp.Participants[0].First)
					assert.Equal(t, outer.Last, resp.Participants[0].Last)
					assert.Equal(t, outer.Birthdate, resp.Participants[0].Birthdate)
					assert.Equal(t, outer.Gender, resp.Participants[0].Gender)
					assert.Equal(t, outer.AgeGroup, resp.Participants[0].AgeGroup)
					assert.Equal(t, outer.Distance, resp.Participants[0].Distance)
					assert.Equal(t, outer.Anonymous, resp.Participants[0].Anonymous)
					assert.Equal(t, outer.SMSEnabled, resp.Participants[0].SMSEnabled)
					assert.Equal(t, outer.Mobile, resp.Participants[0].Mobile)
					assert.Equal(t, outer.Apparel, resp.Participants[0].Apparel)
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
}
