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

func TestGetEventYear(t *testing.T) {
	// POST, /event-year
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	// Test no key
	t.Log("Testing no key given.")
	body, err := json.Marshal(types.GetEventYearRequest{
		Slug: variables.events["event2"].Slug,
		Year: variables.eventYears["event2"]["2020"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/event-year", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.GetEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodPost, "/event-year", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodPost, "/event-year", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test restricted event, wrong account key
	t.Log("Testing restricted event, wrong account key.")
	request = httptest.NewRequest(http.MethodPost, "/event-year", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test unrestricted event, wrong account key
	body, err = json.Marshal(types.GetEventYearRequest{
		Slug: variables.events["event1"].Slug,
		Year: variables.eventYears["event1"]["2020"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/event-year", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEventYear(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.EventYearResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, variables.events["event1"].Name, resp.Event.Name)
			assert.Equal(t, variables.events["event1"].Slug, resp.Event.Slug)
			assert.Equal(t, variables.events["event1"].Website, resp.Event.Website)
			assert.Equal(t, variables.events["event1"].Image, resp.Event.Image)
			assert.Equal(t, variables.events["event1"].ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, variables.events["event1"].AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, variables.events["event1"].Type, resp.Event.Type)
			assert.Equal(t, variables.eventYears["event1"]["2020"].Year, resp.EventYear.Year)
			assert.Equal(t, variables.eventYears["event1"]["2020"].Live, resp.EventYear.Live)
			assert.Equal(t, variables.eventYears["event1"]["2020"].DaysAllowed, resp.EventYear.DaysAllowed)
			assert.True(t, variables.eventYears["event1"]["2020"].DateTime.Equal(resp.EventYear.DateTime))
		}
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodPost, "/event-year", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodPost, "/event-year", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/event-year", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEventYear(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test empty request
	t.Log("Testing empty request.")
	request = httptest.NewRequest(http.MethodPost, "/event-year", strings.NewReader("{}"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEventYear(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/event-year", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEventYear(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test valid request
	t.Log("Testing valid request.")
	body, err = json.Marshal(types.GetEventYearRequest{
		Slug: variables.events["event2"].Slug,
		Year: variables.eventYears["event2"]["2020"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/event-year", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEventYear(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.EventYearResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, variables.events["event2"].Name, resp.Event.Name)
			assert.Equal(t, variables.events["event2"].Slug, resp.Event.Slug)
			assert.Equal(t, variables.events["event2"].Website, resp.Event.Website)
			assert.Equal(t, variables.events["event2"].Image, resp.Event.Image)
			assert.Equal(t, variables.events["event2"].ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, variables.events["event2"].AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, variables.events["event2"].Type, resp.Event.Type)
			assert.Equal(t, variables.eventYears["event2"]["2020"].Year, resp.EventYear.Year)
			assert.Equal(t, variables.eventYears["event2"]["2020"].Live, resp.EventYear.Live)
			assert.Equal(t, variables.eventYears["event2"]["2020"].DaysAllowed, resp.EventYear.DaysAllowed)
			assert.Equal(t, variables.eventYears["event2"]["2020"].DateTime.Local(), resp.EventYear.DateTime.Local())
		}
	}
	// Test invalid event
	t.Log("Testing invalid event.")
	body, err = json.Marshal(types.GetEventYearRequest{
		Slug: "not-a-real-event",
		Year: variables.eventYears["event2"]["2020"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/event-year", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEventYear(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid event invalid year
	t.Log("Testing valid event, invalid year.")
	body, err = json.Marshal(types.GetEventYearRequest{
		Slug: variables.events["event2"].Slug,
		Year: "invalid-year",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/event-year", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEventYear(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
}

func TestGetEventYears(t *testing.T) {
	// POST, /event-year/event
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	body, err := json.Marshal(types.GetEventRequest{
		Slug: variables.events["event1"].Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	// Test no key
	t.Log("Testing no key given.")
	request := httptest.NewRequest(http.MethodPost, "/event-year/event", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.GetEventYears(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodPost, "/event-year/event", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEventYears(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodPost, "/event-year/event", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEventYears(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodPost, "/event-year/event", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEventYears(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodPost, "/event-year/event", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEventYears(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/event-year/event", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEventYears(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/event-year/event", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEventYears(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test empty request
	t.Log("Testing empty request.")
	request = httptest.NewRequest(http.MethodPost, "/event-year/event", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEventYears(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test restricted event, wrong account key
	body, err = json.Marshal(types.GetEventRequest{
		Slug: variables.events["event2"].Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing restricted event, wrong account key.")
	request = httptest.NewRequest(http.MethodPost, "/event-year/event", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEventYears(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test restricted event, correct account key
	body, err = json.Marshal(types.GetEventRequest{
		Slug: variables.events["event2"].Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing restricted event, correct account key.")
	request = httptest.NewRequest(http.MethodPost, "/event-year/event", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEventYears(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.EventYearsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(variables.eventYears["event2"]), len(resp.EventYears))
			for _, outer := range variables.eventYears["event2"] {
				found := false
				for _, inner := range resp.EventYears {
					if outer.Equals(&inner) {
						found = true
					}
				}
				assert.True(t, found)
			}
		}
	}
	// Test valid request
	body, err = json.Marshal(types.GetEventRequest{
		Slug: variables.events["event1"].Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing valid request.")
	request = httptest.NewRequest(http.MethodPost, "/event-year/event", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEventYears(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.EventYearsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(variables.eventYears["event1"]), len(resp.EventYears))
			for _, outer := range variables.eventYears["event1"] {
				found := false
				for _, inner := range resp.EventYears {
					if outer.Equals(&inner) {
						found = true
					}
				}
				assert.True(t, found)
			}
		}
	}
	// Test invalid event
	body, err = json.Marshal(types.GetEventRequest{
		Slug: "invalid-event",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing invalid event.")
	request = httptest.NewRequest(http.MethodPost, "/event-year/event", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEventYears(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
}

func TestAddEventYear(t *testing.T) {
	// POST, /event-year/add
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	body, err := json.Marshal(types.ModifyEventYearRequest{
		Slug: variables.events["event1"].Slug,
		EventYear: types.RequestYear{
			Year:     "2022",
			DateTime: "2022/04/05 9:00:00 -07:00",
			Live:     false,
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	// Test no key
	t.Log("Testing no key given.")
	request := httptest.NewRequest(http.MethodPost, "/event-year/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.AddEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodPost, "/event-year/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodPost, "/event-year/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test read key
	t.Log("Testing read host.")
	request = httptest.NewRequest(http.MethodPost, "/event-year/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodPost, "/event-year/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodPost, "/event-year/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/event-year/add", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEventYear(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test empty request
	t.Log("Testing empty request.")
	request = httptest.NewRequest(http.MethodPost, "/event-year/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEventYear(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/event-year/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEventYear(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test unknown event
	body, err = json.Marshal(types.ModifyEventYearRequest{
		Slug: "unknown-event",
		EventYear: types.RequestYear{
			Year:     "2022",
			DateTime: "2022/04/05 9:00:00 -07:00",
			Live:     false,
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing unknown event.")
	request = httptest.NewRequest(http.MethodPost, "/event-year/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEventYear(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test wrong account key
	eYear := types.RequestYear{
		Year:     "2022",
		DateTime: "2022/04/05 09:00:00 -07:00",
		Live:     false,
	}
	body, err = json.Marshal(types.ModifyEventYearRequest{
		Slug:      variables.events["event2"].Slug,
		EventYear: eYear,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing wrong account key.")
	request = httptest.NewRequest(http.MethodPost, "/event-year/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test valid request
	eYear = types.RequestYear{
		Year:        "2022",
		DateTime:    "2022/04/05 09:00:00 +00:00",
		Live:        false,
		DaysAllowed: 1,
	}
	body, err = json.Marshal(types.ModifyEventYearRequest{
		Slug:      variables.events["event2"].Slug,
		EventYear: eYear,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing valid request.")
	request = httptest.NewRequest(http.MethodPost, "/event-year/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEventYear(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		nEv, err := database.GetEventYear("event2", "2022")
		if assert.NoError(t, err) {
			assert.Equal(t, eYear.Year, nEv.Year)
			assert.Equal(t, eYear.Live, nEv.Live)
			assert.Equal(t, eYear.DaysAllowed, nEv.DaysAllowed)
			assert.Equal(t, eYear.DateTime, nEv.DateTime.Format("2006/01/02 15:04:05 -07:00"))
		}
		var resp types.EventYearResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, variables.events["event2"].Name, resp.Event.Name)
			assert.Equal(t, variables.events["event2"].Slug, resp.Event.Slug)
			assert.Equal(t, variables.events["event2"].Image, resp.Event.Image)
			assert.Equal(t, variables.events["event2"].Website, resp.Event.Website)
			assert.Equal(t, variables.events["event2"].Type, resp.Event.Type)
			assert.Equal(t, variables.events["event2"].AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, variables.events["event2"].ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, eYear.Year, resp.EventYear.Year)
			assert.Equal(t, eYear.Live, resp.EventYear.Live)
			assert.Equal(t, eYear.DaysAllowed, resp.EventYear.DaysAllowed)
			assert.Equal(t, eYear.DateTime, resp.EventYear.DateTime.Format("2006/01/02 15:04:05 -07:00"))
		}
	}
	// Test valid request 2
	eYear = types.RequestYear{
		Year:     "2022-2",
		DateTime: "2022/05/05 09:00:00 +00:00",
		Live:     false,
	}
	body, err = json.Marshal(types.ModifyEventYearRequest{
		Slug:      variables.events["event2"].Slug,
		EventYear: eYear,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing valid request 2.")
	request = httptest.NewRequest(http.MethodPost, "/event-year/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEventYear(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		nEv, err := database.GetEventYear("event2", "2022-2")
		if assert.NoError(t, err) {
			assert.Equal(t, eYear.Year, nEv.Year)
			assert.Equal(t, eYear.Live, nEv.Live)
			assert.Equal(t, eYear.DaysAllowed, nEv.DaysAllowed)
			assert.Equal(t, eYear.DateTime, nEv.DateTime.Format("2006/01/02 15:04:05 -07:00"))
		}
		var resp types.EventYearResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, variables.events["event2"].Name, resp.Event.Name)
			assert.Equal(t, variables.events["event2"].Slug, resp.Event.Slug)
			assert.Equal(t, variables.events["event2"].Image, resp.Event.Image)
			assert.Equal(t, variables.events["event2"].Website, resp.Event.Website)
			assert.Equal(t, variables.events["event2"].Type, resp.Event.Type)
			assert.Equal(t, variables.events["event2"].AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, variables.events["event2"].ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, eYear.Year, resp.EventYear.Year)
			assert.Equal(t, eYear.Live, resp.EventYear.Live)
			assert.Equal(t, eYear.DaysAllowed, resp.EventYear.DaysAllowed)
			assert.Equal(t, eYear.DateTime, resp.EventYear.DateTime.Format("2006/01/02 15:04:05 -07:00"))
		}
	}
	// Test year collision
	t.Log("Testing year collision.")
	request = httptest.NewRequest(http.MethodPost, "/event-year/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEventYear(c)) {
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	}
	// Test validation errors -- year (valid)
	body, err = json.Marshal(types.ModifyEventYearRequest{
		Slug: variables.events["event2"].Slug,
		EventYear: types.RequestYear{
			Year:     "invalid-year",
			DateTime: "2022/04/05 9:00:00 +00:00",
			Live:     false,
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing validation errors -- year (valid).")
	request = httptest.NewRequest(http.MethodPost, "/event-year/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEventYear(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	// Test validation errors -- year
	body, err = json.Marshal(types.ModifyEventYearRequest{
		Slug: variables.events["event2"].Slug,
		EventYear: types.RequestYear{
			Year:     "invalid-year$_",
			DateTime: "2022/04/05 9:00:00 +00:00",
			Live:     false,
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing validation errors -- year.")
	request = httptest.NewRequest(http.MethodPost, "/event-year/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEventYear(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test validation errors -- date
	body, err = json.Marshal(types.ModifyEventYearRequest{
		Slug: variables.events["event2"].Slug,
		EventYear: types.RequestYear{
			Year:     "2025",
			DateTime: "2025/04 9:00:00",
			Live:     false,
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing validation errors -- date.")
	request = httptest.NewRequest(http.MethodPost, "/event-year/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEventYear(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestUpdateEventYear(t *testing.T) {
	// PUT, /event-year/update
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	body, err := json.Marshal(types.ModifyEventYearRequest{
		Slug: variables.events["event2"].Slug,
		EventYear: types.RequestYear{
			Year:        variables.eventYears["event2"]["2021"].Year,
			DateTime:    "2022/04/05 9:00:00 -07:00",
			Live:        false,
			DaysAllowed: 9,
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	// Test no key
	t.Log("Testing no key given.")
	request := httptest.NewRequest(http.MethodPut, "/event-year/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodPut, "/event-year/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodPut, "/event-year/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test read key
	t.Log("Testing read host.")
	request = httptest.NewRequest(http.MethodPut, "/event-year/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test wrong account key
	t.Log("Testing wrong account key.")
	request = httptest.NewRequest(http.MethodPut, "/event-year/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodPut, "/event-year/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodPut, "/event-year/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPut, "/event-year/update", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEventYear(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test empty request
	t.Log("Testing empty request.")
	request = httptest.NewRequest(http.MethodPut, "/event-year/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEventYear(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPut, "/event-year/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEventYear(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test valid request
	eYear := types.RequestYear{
		Year:        variables.eventYears["event2"]["2021"].Year,
		DateTime:    variables.eventYears["event2"]["2021"].DateTime.Add(time.Hour * 48).Format("2006/01/02 15:04:05 -07:00"),
		Live:        !variables.eventYears["event2"]["2021"].Live,
		DaysAllowed: variables.eventYears["event2"]["2021"].DaysAllowed + 2,
	}
	body, err = json.Marshal(types.ModifyEventYearRequest{
		Slug:      variables.events["event2"].Slug,
		EventYear: eYear,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing valid request.")
	request = httptest.NewRequest(http.MethodPut, "/event-year/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEventYear(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		nEv, err := database.GetEventYear(variables.events["event2"].Slug, variables.eventYears["event2"]["2021"].Year)
		if assert.NoError(t, err) {
			assert.Equal(t, eYear.Year, nEv.Year)
			assert.Equal(t, eYear.Live, nEv.Live)
			assert.Equal(t, eYear.DaysAllowed, nEv.DaysAllowed)
			assert.Equal(t, eYear.DateTime, nEv.DateTime.Format("2006/01/02 15:04:05 -07:00"))
		}
		var resp types.EventYearResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, variables.events["event2"].Name, resp.Event.Name)
			assert.Equal(t, variables.events["event2"].Slug, resp.Event.Slug)
			assert.Equal(t, variables.events["event2"].Image, resp.Event.Image)
			assert.Equal(t, variables.events["event2"].Website, resp.Event.Website)
			assert.Equal(t, variables.events["event2"].Type, resp.Event.Type)
			assert.Equal(t, variables.events["event2"].AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, variables.events["event2"].ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, eYear.Year, resp.EventYear.Year)
			assert.Equal(t, eYear.Live, resp.EventYear.Live)
			assert.Equal(t, eYear.DaysAllowed, resp.EventYear.DaysAllowed)
			assert.Equal(t, eYear.DateTime, resp.EventYear.DateTime.Format("2006/01/02 15:04:05 -07:00"))
		}
	}
	// Test attempt to update year
	eYear = types.RequestYear{
		Year:        variables.eventYears["event2"]["2021"].Year + "-2",
		DateTime:    variables.eventYears["event2"]["2021"].DateTime.Add(time.Hour * 48).Format("2006/01/02 15:04:05 -07:00"),
		Live:        variables.eventYears["event2"]["2021"].Live,
		DaysAllowed: variables.eventYears["event2"]["2021"].DaysAllowed + 3,
	}
	body, err = json.Marshal(types.ModifyEventYearRequest{
		Slug:      variables.events["event2"].Slug,
		EventYear: eYear,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing attempt to update year - invalid.")
	request = httptest.NewRequest(http.MethodPut, "/event-year/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEventYear(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test unknown event
	body, err = json.Marshal(types.ModifyEventYearRequest{
		Slug:      "event12",
		EventYear: eYear,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing unknown event.")
	request = httptest.NewRequest(http.MethodPut, "/event-year/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEventYear(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test validation errors -- year
	body, err = json.Marshal(types.ModifyEventYearRequest{
		Slug: variables.events["event2"].Slug,
		EventYear: types.RequestYear{
			Year:     "invalid-year$",
			DateTime: "2022/04/05 9:00:00 -07:00",
			Live:     false,
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing validation errors -- year.")
	request = httptest.NewRequest(http.MethodPut, "/event-year/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEventYear(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test validation errors -- date
	body, err = json.Marshal(types.ModifyEventYearRequest{
		Slug: variables.events["event2"].Slug,
		EventYear: types.RequestYear{
			Year:     variables.eventYears["event2"]["2021"].Year,
			DateTime: "2025/04 9:00:00 -07:00",
			Live:     variables.eventYears["event2"]["2021"].Live,
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing validation errors -- date.")
	request = httptest.NewRequest(http.MethodPut, "/event-year/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEventYear(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestDeleteEventYear(t *testing.T) {
	// DELETE, /event-year/delete
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	body, err := json.Marshal(types.DeleteEventYearRequest{
		Slug: "event3",
		Year: "2025",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	// Test no key
	t.Log("Testing no key given.")
	request := httptest.NewRequest(http.MethodDelete, "/event-year/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodDelete, "/event-year/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodDelete, "/event-year/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test read key
	t.Log("Testing read key.")
	request = httptest.NewRequest(http.MethodDelete, "/event-year/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test write key
	t.Log("Testing read key.")
	request = httptest.NewRequest(http.MethodDelete, "/event-year/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test wrong account key
	body, err = json.Marshal(types.DeleteEventYearRequest{
		Slug: variables.events["event1"].Slug,
		Year: variables.eventYears["event1"]["2021"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing wrong account key.")
	request = httptest.NewRequest(http.MethodDelete, "/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid host
	body, err = json.Marshal(types.DeleteEventYearRequest{
		Slug: variables.events["event2"].Slug,
		Year: variables.eventYears["event2"]["2021"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodDelete, "/event-year/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodDelete, "/event-year/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEventYear(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodDelete, "/event-year/delete", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEventYear(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test empty request
	t.Log("Testing empty request.")
	request = httptest.NewRequest(http.MethodDelete, "/event-year/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEventYear(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodDelete, "/event-year/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEventYear(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test valid request
	t.Log("Testing valid request.")
	request = httptest.NewRequest(http.MethodDelete, "/event-year/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEventYear(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		nEv, err := database.GetEventYear(variables.events["event2"].Slug, variables.eventYears["event2"]["2021"].Year)
		if assert.NoError(t, err) {
			assert.Nil(t, nEv)
		}
	}
	// Test unknown event
	body, err = json.Marshal(types.DeleteEventYearRequest{
		Slug: "event12",
		Year: "2021",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing unknown event.")
	request = httptest.NewRequest(http.MethodDelete, "/event-year/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEventYear(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test known event, unknown year
	body, err = json.Marshal(types.DeleteEventYearRequest{
		Slug: variables.events["event1"].Slug,
		Year: "2025",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing unknown event.")
	request = httptest.NewRequest(http.MethodDelete, "/event-year/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEventYear(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
}
