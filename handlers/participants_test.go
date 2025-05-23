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

func TestGetParticipants(t *testing.T) {
	// GET, /participants
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	// Test no key
	t.Log("Testing no key given.")
	year := variables.eventYears["event2"]["2020"].Year
	body, err := json.Marshal(types.GetParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request := httptest.NewRequest(http.MethodGet, "/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodGet, "/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodGet, "/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test read key
	t.Log("Testing read host.")
	request = httptest.NewRequest(http.MethodPost, "/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test restricted event, wrong account key
	t.Log("Testing restricted event, wrong account key.")
	request = httptest.NewRequest(http.MethodGet, "/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test unrestricted event, wrong account key
	t.Log("Testing unrestricted event, wrong account key.")
	year = variables.eventYears["event1"]["2021"].Year
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug: variables.events["event1"].Slug,
		Year: &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodGet, "/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test unrestricted event, admin key
	t.Log("Testing unrestricted event, wrong account key.")
	year = variables.eventYears["event2"]["2021"].Year
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodGet, "/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodGet, "/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodGet, "/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodGet, "/participants", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test empty request
	t.Log("Testing empty request.")
	request = httptest.NewRequest(http.MethodGet, "/participants", strings.NewReader("{}"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodGet, "/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test valid request
	t.Log("Testing valid request.")
	year = variables.eventYears["event2"]["2020"].Year
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodGet, "/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
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
			assert.Equal(t, variables.eventYears["event2"]["2020"].Year, resp.Year.Year)
			assert.True(t, variables.eventYears["event2"]["2020"].DateTime.Equal(resp.Year.DateTime))
			assert.Equal(t, variables.eventYears["event2"]["2020"].Live, resp.Year.Live)
			assert.Equal(t, 4, len(resp.Participants))
		}
	}
	// Test valid request -- UpdatedAfter
	t.Log("Testing valid request -- UpdatedAfter.")
	year = variables.eventYears["event2"]["2020"].Year
	updatedAfter := int64(100)
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug:         variables.events["event2"].Slug,
		Year:         &year,
		UpdatedAfter: &updatedAfter,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodGet, "/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
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
			assert.Equal(t, variables.eventYears["event2"]["2020"].Year, resp.Year.Year)
			assert.True(t, variables.eventYears["event2"]["2020"].DateTime.Equal(resp.Year.DateTime))
			assert.Equal(t, variables.eventYears["event2"]["2020"].Live, resp.Year.Live)
			assert.Equal(t, 2, len(resp.Participants))
		}
	}
	// Test valid request - Limit and Page defined
	t.Log("Testing valid request. Limit and Page defined")
	year = variables.eventYears["event1"]["2021"].Year
	limit := 50
	page := 1
	fetchedParticipants := make([]types.Participant, 0)
	for {
		body, err = json.Marshal(types.GetParticipantsRequest{
			Slug:  variables.events["event1"].Slug,
			Year:  &year,
			Limit: &limit,
			Page:  &page,
		})
		if err != nil {
			t.Fatalf("Error encoding request body into json object: %v", err)
		}
		request = httptest.NewRequest(http.MethodGet, "/participants", strings.NewReader(string(body)))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete3"])
		response = httptest.NewRecorder()
		c = e.NewContext(request, response)
		if assert.NoError(t, h.GetParticipants(c)) {
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
				assert.True(t, len(resp.Participants) == 50 || len(resp.Participants) == 0)
				for _, outer := range resp.Participants {
					found := false
					for _, inner := range fetchedParticipants {
						if outer.AlternateId == inner.AlternateId {
							found = true
						}
					}
					assert.False(t, found)
				}
				fetchedParticipants = append(fetchedParticipants, resp.Participants...)
				if len(resp.Participants) != 50 {
					break
				}
			} else {
				break
			}
		} else {
			break
		}
		page += 1
	}
	assert.Equal(t, 500, len(fetchedParticipants))
	// Test valid request - no year
	t.Log("Testing valid request.")
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug: variables.events["event2"].Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodGet, "/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
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
	// Test invalid event
	t.Log("Testing invalid event.")
	year = variables.eventYears["event2"]["2020"].Year
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug: "not-a-real-event",
		Year: &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodGet, "/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid event invalid year
	t.Log("Testing valid event, invalid year.")
	year = "invalid-year"
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodGet, "/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
}

func TestAddParticipants(t *testing.T) {
	// POST, /participants/add
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	year, err := database.AddEventYear(types.EventYear{
		EventIdentifier: variables.events["event1"].Identifier,
		Year:            "2023",
		DateTime:        time.Date(2023, 04, 05, 9, 0, 0, 0, time.Local),
		Live:            false,
		DaysAllowed:     1,
		RankingType:     "chip",
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
			Gender:      "Man",
			AgeGroup:    "20-29",
			Distance:    "1 Mile",
			Anonymous:   false,
			SMSEnabled:  false,
			Mobile:      "",
			Apparel:     "",
			UpdatedAt:   50,
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
			UpdatedAt:   150,
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
			UpdatedAt:   0,
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
			UpdatedAt:   0,
		},
	}
	body, err := json.Marshal(types.AddParticipantsRequest{
		Slug:         variables.events["event1"].Slug,
		Year:         year.Year,
		Participants: parts,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	// Test no key
	t.Log("Testing no key given.")
	request := httptest.NewRequest(http.MethodPost, "/participants/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.AddParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodPost, "/participants/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodPost, "/participants/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test read key
	t.Log("Testing read host.")
	request = httptest.NewRequest(http.MethodPost, "/participants/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodPost, "/participants/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodPost, "/participants/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/participants/add", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete3"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test empty request
	t.Log("Testing empty request.")
	request = httptest.NewRequest(http.MethodPost, "/participants/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete3"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/participants/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete3"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test unknown event
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug:         "unknown-event",
		Year:         year.Year,
		Participants: parts,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing unknown event.")
	request = httptest.NewRequest(http.MethodPost, "/participants/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete3"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddParticipants(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test known event unknown year
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug:         variables.events["event1"].Slug,
		Year:         "invalid",
		Participants: parts,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing known event unknown year.")
	request = httptest.NewRequest(http.MethodPost, "/participants/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete3"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddParticipants(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test wrong account key
	year.EventIdentifier = variables.events["event2"].Identifier
	year, err = database.AddEventYear(*year)
	if err != nil {
		t.Fatalf("Error adding event year to database: %v", err)
	}
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug:         variables.events["event2"].Slug,
		Year:         year.Year,
		Participants: parts,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing wrong account key.")
	request = httptest.NewRequest(http.MethodPost, "/participants/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test valid request
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug:         variables.events["event2"].Slug,
		Year:         year.Year,
		Participants: parts[0:2],
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing valid request.")
	request = httptest.NewRequest(http.MethodPost, "/participants/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddParticipants(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		part, err := database.GetParticipants(year.Identifier, 0, 0, nil)
		if assert.NoError(t, err) {
			assert.Equal(t, len(parts[0:2]), len(part))
			for _, outer := range parts[0:2] {
				found := false
				for _, inner := range part {
					if outer.AlternateId == inner.AlternateId {
						assert.True(t, outer.TestEquals(&inner))
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
		var resp types.AddParticipantsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(parts[0:2]), resp.Count)
			assert.Equal(t, 0, len(resp.Updated))
		}
	}
	// Test valid request - UpdatedAfter
	t.Log("Testing valid request -- UpdatedAfter.")
	database.DeleteParticipants(year.Identifier, []string{})
	updatedAfter := time.Now().UTC().Unix()
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug:         variables.events["event2"].Slug,
		Year:         year.Year,
		Participants: parts[2:4],
		UpdatedAfter: &updatedAfter,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/participants/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddParticipants(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		part, err := database.GetParticipants(year.Identifier, 0, 0, nil)
		if assert.NoError(t, err) {
			assert.Equal(t, len(parts[2:4]), len(part))
			for _, outer := range parts[2:4] {
				found := false
				for _, inner := range part {
					if outer.AlternateId == inner.AlternateId {
						assert.True(t, outer.TestEquals(&inner))
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
		var resp types.AddParticipantsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(parts[2:4]), resp.Count)
			assert.Equal(t, len(parts[2:4]), len(resp.Updated))
		}
	}
	// validation -- age
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: year.Year,
		Participants: []types.Participant{
			{
				AlternateId: "1011111",
				Bib:         "10",
				Birthdate:   "1/1/3050",
				AgeGroup:    "10-20",
				First:       "John",
				Last:        "Jacob",
				Distance:    "1 Mile",
				Gender:      "Man",
				Anonymous:   false,
				SMSEnabled:  false,
				Mobile:      "",
				Apparel:     "",
			},
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing validation -- age.")
	request = httptest.NewRequest(http.MethodPost, "/participants/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// validation -- gender
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: year.Year,
		Participants: []types.Participant{
			{
				AlternateId: "1011111",
				Bib:         "10",
				Birthdate:   "1/1/2004",
				AgeGroup:    "10-20",
				First:       "John",
				Last:        "Jacob",
				Distance:    "1 Mile",
				Gender:      "G",
				Anonymous:   false,
				SMSEnabled:  false,
				Mobile:      "",
				Apparel:     "",
			},
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing validation -- gender.") // Previously gender was checked, it no longer is.
	request = httptest.NewRequest(http.MethodPost, "/participants/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddParticipants(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	// validation -- distance
	body, err = json.Marshal(types.AddParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: year.Year,
		Participants: []types.Participant{
			{
				AlternateId: "1011111",
				Bib:         "10",
				Birthdate:   "1/1/2004",
				AgeGroup:    "10-20",
				First:       "John",
				Last:        "Jacob",
				Distance:    "",
				Gender:      "Man",
				Anonymous:   false,
				SMSEnabled:  false,
				Mobile:      "",
				Apparel:     "",
			},
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing validation -- distance.")
	request = httptest.NewRequest(http.MethodPost, "/participants/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
}

func TestDeleteParticipants(t *testing.T) {
	// DELETE, /participants/delete
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	year, err := database.AddEventYear(types.EventYear{
		EventIdentifier: variables.events["event1"].Identifier,
		Year:            "2023",
		DateTime:        time.Date(2023, 04, 05, 9, 0, 0, 0, time.Local),
		Live:            false,
		DaysAllowed:     1,
		RankingType:     "chip",
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
			Gender:      "Man",
			AgeGroup:    "20-29",
			Distance:    "1 Mile",
			Anonymous:   false,
			SMSEnabled:  false,
			Mobile:      "",
			Apparel:     "",
			UpdatedAt:   50,
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
			UpdatedAt:   150,
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
			UpdatedAt:   0,
		},
		{
			AlternateId: "1364",
			Bib:         "1364",
			First:       "Lilly",
			Last:        "Smith",
			Birthdate:   "1/1/2014",
			Gender:      "Non-Binary",
			AgeGroup:    "0-19",
			Distance:    "1 Mile",
			Anonymous:   false,
			SMSEnabled:  false,
			Mobile:      "",
			Apparel:     "",
			UpdatedAt:   0,
		},
	}
	p, err := database.AddParticipants(year.Identifier, parts)
	if err != nil {
		t.Fatalf("Error adding participants to database for test: %v", err)
	}
	assert.Equal(t, len(parts), len(p))
	body, err := json.Marshal(types.DeleteParticipantsRequest{
		Slug: variables.events["event1"].Slug,
		Year: year.Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	// Test no key
	t.Log("Testing no key given.")
	request := httptest.NewRequest(http.MethodDelete, "/participants/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.DeleteParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodDelete, "/participants/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodDelete, "/participants/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// test read key
	t.Log("Testing read key.")
	request = httptest.NewRequest(http.MethodDelete, "/participants/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// test write key
	t.Log("Testing write key.")
	request = httptest.NewRequest(http.MethodDelete, "/participants/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteParticipants(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	// Test wrong account key
	t.Log("Testing wrong account key.")
	request = httptest.NewRequest(http.MethodDelete, "/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodDelete, "/participants/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodDelete, "/participants/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodDelete, "/participants/delete", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete3"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test empty request
	t.Log("Testing empty request.")
	request = httptest.NewRequest(http.MethodDelete, "/participants/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete3"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodDelete, "/participants/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete3"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test valid request
	p, err = database.AddParticipants(year.Identifier, parts)
	if err != nil {
		t.Fatalf("Error adding participants to database for test: %v", err)
	}
	assert.Equal(t, len(parts), len(p))
	t.Log("Testing valid request.")
	request = httptest.NewRequest(http.MethodDelete, "/participants/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete3"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteParticipants(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		newP, err := database.GetParticipants(year.Identifier, 0, 0, nil)
		if assert.NoError(t, err) {
			assert.Equal(t, 0, len(newP))
		}
	}
	// Test valid request with valid bibs
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
		Slug:        variables.events["event1"].Slug,
		Year:        year.Year,
		Identifiers: idents[2:],
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing valid request - bibs specified.")
	request = httptest.NewRequest(http.MethodDelete, "/participants/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete3"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteParticipants(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		newP, err := database.GetParticipants(year.Identifier, 0, 0, nil)
		if assert.NoError(t, err) {
			assert.Equal(t, len(parts)-2, len(newP))
		}
	}
	// Test valid request with unknown bibs
	p, err = database.AddParticipants(year.Identifier, parts)
	if err != nil {
		t.Fatalf("Error adding participants to database for test: %v", err)
	}
	assert.Equal(t, len(parts), len(p))
	idents = make([]string, 0)
	for _, person := range parts {
		idents = append(idents, person.Bib)
	}
	idents = append(idents, "not-valid")
	body, err = json.Marshal(types.DeleteParticipantsRequest{
		Slug:        variables.events["event1"].Slug,
		Year:        year.Year,
		Identifiers: idents[2:],
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing valid request - bibs specified.")
	request = httptest.NewRequest(http.MethodDelete, "/participants/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete3"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteParticipants(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		newP, err := database.GetParticipants(year.Identifier, 0, 0, nil)
		if assert.NoError(t, err) {
			assert.Equal(t, len(parts)-2, len(newP))
		}
	}
	// Test admin delete other
	p, err = database.GetParticipants(variables.eventYears["event2"]["2020"].Identifier, 0, 0, nil)
	if err != nil {
		t.Fatalf("Error getting participants from database: %v", err)
	}
	assert.NotEqual(t, 0, len(p))
	body, err = json.Marshal(types.DeleteParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: variables.eventYears["event2"]["2020"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing admin delete for other.")
	request = httptest.NewRequest(http.MethodDelete, "/participants/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete3"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test unknown event
	body, err = json.Marshal(types.DeleteParticipantsRequest{
		Slug: "event12",
		Year: "2021",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing unknown event.")
	request = httptest.NewRequest(http.MethodDelete, "/participants/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete3"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteParticipants(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test known event, unknown year
	body, err = json.Marshal(types.DeleteParticipantsRequest{
		Slug: variables.events["event1"].Slug,
		Year: "2025",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing unknown event.")
	request = httptest.NewRequest(http.MethodDelete, "/participants/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete3"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteParticipants(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
}
