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
	// POST, /participants
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	// Test no key
	t.Log("Testing no key given.")
	body, err := json.Marshal(types.GetParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: variables.eventYears["event2"]["2020"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodPost, "/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodPost, "/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test restricted event, wrong account key
	t.Log("Testing restricted event, wrong account key.")
	request = httptest.NewRequest(http.MethodPost, "/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test unrestricted event, wrong account key
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug: variables.events["event1"].Slug,
		Year: variables.eventYears["event1"]["2021"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetParticipantsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 500, len(resp.Participants))
		}
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodPost, "/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodPost, "/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/participants", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test empty request
	t.Log("Testing empty request.")
	request = httptest.NewRequest(http.MethodPost, "/participants", strings.NewReader("{}"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test valid request
	t.Log("Testing valid request.")
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: variables.eventYears["event2"]["2020"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetParticipantsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 4, len(resp.Participants))
		}
	}
	// Test invalid event
	t.Log("Testing invalid event.")
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug: "not-a-real-event",
		Year: variables.eventYears["event2"]["2020"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/participants", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetParticipants(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid event invalid year
	t.Log("Testing valid event, invalid year.")
	body, err = json.Marshal(types.GetParticipantsRequest{
		Slug: variables.events["event2"].Slug,
		Year: "invalid-year",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/participants", strings.NewReader(string(body)))
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
		Slug:         variables.events["event1"].Slug,
		Year:         year.Year,
		Participants: people,
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
		Participants: people,
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
		Participants: people,
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
		Participants: people,
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
		Participants: people,
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
	// validation -- age
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
	p, err = database.AddPeople(year.Identifier, people)
	if err != nil {
		t.Fatalf("Error adding people to database for test: %v", err)
	}
	assert.Equal(t, len(people), len(p))
	t.Log("Testing valid request.")
	request = httptest.NewRequest(http.MethodDelete, "/participants/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete3"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteParticipants(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		newP, err := database.GetPeople(variables.events["event1"].Slug, year.Year)
		if assert.NoError(t, err) {
			assert.Equal(t, 0, len(newP))
		}
	}
	// Test valid request with valid bibs
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
		Slug: variables.events["event1"].Slug,
		Year: year.Year,
		Bibs: bibs[2:],
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
		newP, err := database.GetPeople(variables.events["event1"].Slug, year.Year)
		if assert.NoError(t, err) {
			assert.Equal(t, len(people)-2, len(newP))
		}
	}
	// Test valid request with unknown bibs
	p, err = database.AddPeople(year.Identifier, people)
	if err != nil {
		t.Fatalf("Error adding people to database for test: %v", err)
	}
	assert.Equal(t, len(people), len(p))
	bibs = make([]string, 0)
	for _, person := range people {
		bibs = append(bibs, person.Bib)
	}
	bibs = append(bibs, "not-valid")
	body, err = json.Marshal(types.DeleteParticipantsRequest{
		Slug: variables.events["event1"].Slug,
		Year: year.Year,
		Bibs: bibs[2:],
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
		newP, err := database.GetPeople(variables.events["event1"].Slug, year.Year)
		if assert.NoError(t, err) {
			assert.Equal(t, len(people)-2, len(newP))
		}
	}
	// Test admin delete other
	p, err = database.GetPeople(variables.events["event2"].Slug, variables.eventYears["event2"]["2020"].Year)
	if err != nil {
		t.Fatalf("Error getting people from database: %v", err)
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
