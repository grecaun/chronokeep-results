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
	request = httptest.NewRequest(http.MethodPost, "/r/event", strings.NewReader(string("")))
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
	people := []types.Person{
		{
			Bib:      "1024",
			First:    "John",
			Last:     "Smnit",
			Age:      20,
			Gender:   "M",
			AgeGroup: "20-29",
			Distance: "1 Mile",
		},
		{
			Bib:      "2034",
			First:    "Jason",
			Last:     "Jonson",
			Age:      34,
			Gender:   "M",
			AgeGroup: "30-39",
			Distance: "5 Mile",
		},
		{
			Bib:      "3521",
			First:    "Rose",
			Last:     "McGowna",
			Age:      16,
			Gender:   "F",
			AgeGroup: "0-19",
			Distance: "1 Mile",
		},
		{
			Bib:      "1364",
			First:    "Lilly",
			Last:     "Smith",
			Age:      10,
			Gender:   "F",
			AgeGroup: "0-19",
			Distance: "1 Mile",
		},
	}
	body, err := json.Marshal(types.AddParticipantsRequest{
		Slug:         variables.events["event2"].Slug,
		Year:         year.Year,
		Participants: people,
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
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug:         variables.events["event2"].Slug,
		Year:         year.Year,
		Participants: people,
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
		part, err := database.GetPeople(variables.events["event2"].Slug, year.Year)
		if assert.NoError(t, err) {
			assert.Equal(t, len(people), len(part))
			for _, outer := range people {
				found := false
				for _, inner := range part {
					if outer.Bib == inner.Bib {
						assert.Equal(t, outer.Age, inner.Age)
						assert.Equal(t, outer.AgeGroup, inner.AgeGroup)
						assert.Equal(t, outer.Bib, inner.Bib)
						assert.Equal(t, outer.Distance, inner.Distance)
						assert.Equal(t, outer.First, inner.First)
						assert.Equal(t, outer.Gender, inner.Gender)
						assert.Equal(t, outer.Last, inner.Last)
						found = true
					}
				}
				assert.True(t, found)
			}
		}
		var resp types.AddResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(people), resp.Count)
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
	request = httptest.NewRequest(http.MethodPost, "/r/participants/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+*token)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RAddParticipants(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		part, err := database.GetPeople(variables.events["event2"].Slug, year.Year)
		if assert.NoError(t, err) {
			assert.Equal(t, len(people), len(part))
			for _, outer := range people {
				found := false
				for _, inner := range part {
					if outer.Bib == inner.Bib {
						assert.Equal(t, outer.Age, inner.Age)
						assert.Equal(t, outer.AgeGroup, inner.AgeGroup)
						assert.Equal(t, outer.Bib, inner.Bib)
						assert.Equal(t, outer.Distance, inner.Distance)
						assert.Equal(t, outer.First, inner.First)
						assert.Equal(t, outer.Gender, inner.Gender)
						assert.Equal(t, outer.Last, inner.Last)
						found = true
					}
				}
				assert.True(t, found)
			}
		}
		var resp types.AddResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(people), resp.Count)
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
		Participants: people,
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
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug:         "unknown",
		Year:         "2020",
		Participants: people,
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
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug:         variables.events["event1"].Slug,
		Year:         "2025",
		Participants: people,
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
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: year.Year,
		Participants: []types.Person{
			{
				Bib:      "10",
				Age:      -20,
				AgeGroup: "10-20",
				First:    "John",
				Last:     "Jacob",
				Distance: "1 Mile",
				Gender:   "m",
			},
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
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: year.Year,
		Participants: []types.Person{
			{
				Bib:      "10",
				Age:      20,
				AgeGroup: "10-20",
				First:    "John",
				Last:     "Jacob",
				Distance: "1 Mile",
				Gender:   "G",
			},
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
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: year.Year,
		Participants: []types.Person{
			{
				Bib:      "10",
				Age:      20,
				AgeGroup: "10-20",
				First:    "John",
				Last:     "Jacob",
				Distance: "",
				Gender:   "M",
			},
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
	people := []types.Person{
		{
			Bib:      "1024",
			First:    "John",
			Last:     "Smnit",
			Age:      20,
			Gender:   "M",
			AgeGroup: "20-29",
			Distance: "1 Mile",
		},
		{
			Bib:      "2034",
			First:    "Jason",
			Last:     "Jonson",
			Age:      34,
			Gender:   "M",
			AgeGroup: "30-39",
			Distance: "5 Mile",
		},
		{
			Bib:      "3521",
			First:    "Rose",
			Last:     "McGowna",
			Age:      16,
			Gender:   "F",
			AgeGroup: "0-19",
			Distance: "1 Mile",
		},
		{
			Bib:      "1364",
			First:    "Lilly",
			Last:     "Smith",
			Age:      10,
			Gender:   "F",
			AgeGroup: "0-19",
			Distance: "1 Mile",
		},
	}
	p, err := database.AddPeople(year.Identifier, people)
	if err != nil {
		t.Fatalf("Error adding people to database for test: %v", err)
	}
	assert.Equal(t, len(people), len(p))
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
		newP, err := database.GetPeople(variables.events["event2"].Slug, year.Year)
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
	p, err = database.AddPeople(year.Identifier, people)
	if err != nil {
		t.Fatalf("Error adding people to database for test: %v", err)
	}
	assert.Equal(t, len(people), len(p))
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
		newP, err := database.GetPeople(variables.events["event2"].Slug, year.Year)
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
	p, err = database.AddPeople(year.Identifier, people)
	if err != nil {
		t.Fatalf("Error adding people to database for test: %v", err)
	}
	assert.Equal(t, len(people), len(p))
	bibs := make([]string, 0)
	for _, person := range people {
		bibs = append(bibs, person.Bib)
	}
	body, err = json.Marshal(types.DeleteParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: year.Year,
		Bibs: bibs[2:],
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
		newP, err := database.GetPeople(variables.events["event2"].Slug, year.Year)
		if assert.NoError(t, err) {
			assert.Equal(t, len(people)-2, len(newP))
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
	p, err = database.AddPeople(year.Identifier, people)
	if err != nil {
		t.Fatalf("Error adding people to database for test: %v", err)
	}
	assert.Equal(t, len(people), len(p))
	bibs = []string{
		"invalid1",
		"invalid2",
	}
	body, err = json.Marshal(types.DeleteParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: year.Year,
		Bibs: bibs,
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
		newP, err := database.GetPeople(variables.events["event2"].Slug, year.Year)
		if assert.NoError(t, err) {
			assert.Equal(t, len(people), len(newP))
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

}
