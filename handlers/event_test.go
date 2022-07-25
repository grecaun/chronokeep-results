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

func TestGetEvent(t *testing.T) {
	// POST, /event
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	// Test no key
	t.Log("Testing no key given.")
	body, err := json.Marshal(types.GetEventRequest{
		Slug: variables.events["event2"].Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/event", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.GetEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodPost, "/event", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodPost, "/event", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test restricted event, wrong account key
	t.Log("Testing restricted event, wrong account key.")
	request = httptest.NewRequest(http.MethodPost, "/event", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test unrestricted event, wrong account key
	body, err = json.Marshal(types.GetEventRequest{
		Slug: variables.events["event1"].Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/event", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEvent(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodPost, "/event", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodPost, "/event", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/event", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test empty request
	t.Log("Testing empty request.")
	request = httptest.NewRequest(http.MethodPost, "/event", strings.NewReader("{}"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEvent(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/event", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test valid request
	t.Log("Testing valid request.")
	body, err = json.Marshal(types.GetEventRequest{
		Slug: variables.events["event2"].Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/event", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEvent(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetEventResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, variables.events["event2"].Name, resp.Event.Name)
			assert.Equal(t, variables.events["event2"].Slug, resp.Event.Slug)
			assert.Equal(t, variables.events["event2"].Website, resp.Event.Website)
			assert.Equal(t, variables.events["event2"].Image, resp.Event.Image)
			assert.Equal(t, variables.events["event2"].ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, variables.events["event2"].AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, variables.events["event2"].Type, resp.Event.Type)
			assert.Equal(t, len(variables.eventYears["event2"]), len(resp.EventYears))
			if assert.NotNil(t, resp.Year) {
				assert.Equal(t, variables.eventYears["event2"]["2021"].Year, resp.Year.Year)
				assert.Equal(t, variables.eventYears["event2"]["2021"].Live, resp.Year.Live)
				assert.Equal(t, variables.eventYears["event2"]["2021"].DateTime.Local(), resp.Year.DateTime.Local())
			}
		}
	}
	// Test invalid event
	t.Log("Testing invalid event.")
	body, err = json.Marshal(types.GetEventRequest{
		Slug: "not-a-real-event",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/event", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEvent(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
}

func TestGetEvents(t *testing.T) {
	// GET, /event/all
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	// Test no key
	t.Log("Testing no key given.")
	request := httptest.NewRequest(http.MethodPost, "/event/all", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.GetEvents(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodPost, "/event/all", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEvents(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodPost, "/event/all", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEvents(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodPost, "/event/all", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEvents(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodPost, "/event/all", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEvents(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request (body doesn't matter)
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/event/all", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEvents(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	// Test wrong content type (body doesn't matter)
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/event/all", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEvents(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	// Test valid request without restricted events
	t.Log("Testing request with account without restricted events.")
	request = httptest.NewRequest(http.MethodPost, "/event/all", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEvents(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetEventsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 1, len(resp.Events))
		}
	}
	// Test valid request with account with restricted events
	t.Log("Testing request with account with restricted events.  This should not pull up restricted events.")
	request = httptest.NewRequest(http.MethodPost, "/event/all", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetEvents(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetEventsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 1, len(resp.Events))
		}
	}
}

func TestGetMyEvents(t *testing.T) {
	// GET, /event/my
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	// Test no key
	t.Log("Testing no key given.")
	request := httptest.NewRequest(http.MethodPost, "/event/my", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.GetMyEvents(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodPost, "/event/my", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetMyEvents(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodPost, "/event/my", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetMyEvents(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodPost, "/event/my", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetMyEvents(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodPost, "/event/my", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetMyEvents(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/event/my", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetMyEvents(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/event/my", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetMyEvents(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	// Test valid request
	t.Log("Testing request with account without restricted events.")
	request = httptest.NewRequest(http.MethodPost, "/event/my", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetMyEvents(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetEventsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 1, len(resp.Events))
		}
	}
	// Test with restricted events
	t.Log("Testing request with account with restricted events.  This should pull up restricted events.")
	request = httptest.NewRequest(http.MethodPost, "/event/my", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetMyEvents(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetEventsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 1, len(resp.Events))
		}
	}
}

func TestAddEvent(t *testing.T) {
	// POST, /event/add
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
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
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	// Test no key
	t.Log("Testing no key given.")
	request := httptest.NewRequest(http.MethodPost, "/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.AddEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodPost, "/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodPost, "/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test read key
	t.Log("Testing read host.")
	request = httptest.NewRequest(http.MethodPost, "/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodPost, "/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodPost, "/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/event/add", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test empty request
	t.Log("Testing empty request.")
	request = httptest.NewRequest(http.MethodPost, "/event/add", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test valid request
	t.Log("Testing valid request.")
	request = httptest.NewRequest(http.MethodPost, "/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEvent(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		nEv, err := database.GetEvent("event3")
		if assert.NoError(t, err) {
			assert.Equal(t, event.Name, nEv.Name)
			assert.Equal(t, event.Slug, nEv.Slug)
			assert.Equal(t, event.ContactEmail, nEv.ContactEmail)
			assert.Equal(t, event.AccessRestricted, nEv.AccessRestricted)
			assert.Equal(t, event.Type, nEv.Type)
		}
	}
	// Test slug collision
	body, err = json.Marshal(types.AddEventRequest{
		Event: types.Event{
			Name:             "Test Event 4",
			Slug:             "event2",
			ContactEmail:     "email@test.com",
			AccessRestricted: false,
			Type:             "distance",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing slug collision.")
	request = httptest.NewRequest(http.MethodPost, "/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEvent(c)) {
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	}
	// Test validation errors - invalid slug
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
	request = httptest.NewRequest(http.MethodPost, "/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test validation errors - no slug
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
	request = httptest.NewRequest(http.MethodPost, "/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test validation errors - invalid name
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
	request = httptest.NewRequest(http.MethodPost, "/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test validation errors - no name
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
	request = httptest.NewRequest(http.MethodPost, "/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test validation errors - invalid type
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
	request = httptest.NewRequest(http.MethodPost, "/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test validation errors - no type
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
	request = httptest.NewRequest(http.MethodPost, "/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test validation errors - invalid website
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
	request = httptest.NewRequest(http.MethodPost, "/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test validation errors - invalid image
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
	request = httptest.NewRequest(http.MethodPost, "/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test validation errors - invalid email
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
	request = httptest.NewRequest(http.MethodPost, "/event/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
}

func TestUpdateEvent(t *testing.T) {
	// PUT, /event/update
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	event := types.Event{
		Name:             "Test Event 3",
		Slug:             "event2",
		ContactEmail:     "email@test2.com",
		AccessRestricted: false,
		Type:             "time",
	}
	body, err := json.Marshal(types.UpdateEventRequest{
		Event: event,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	// Test no key
	t.Log("Testing no key given.")
	request := httptest.NewRequest(http.MethodPost, "/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodPost, "/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodPost, "/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test read key
	t.Log("Testing read host.")
	request = httptest.NewRequest(http.MethodPost, "/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test wrong account key
	t.Log("Testing wrong account key.")
	request = httptest.NewRequest(http.MethodPost, "/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodPost, "/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodPost, "/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/event/update", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test empty request
	t.Log("Testing empty request.")
	request = httptest.NewRequest(http.MethodPost, "/event/update", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test valid request
	t.Log("Testing valid request.")
	request = httptest.NewRequest(http.MethodPost, "/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEvent(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		nEv, err := database.GetEvent("event2")
		if assert.NoError(t, err) {
			assert.Equal(t, event.Name, nEv.Name)
			assert.Equal(t, event.Slug, nEv.Slug)
			assert.Equal(t, event.ContactEmail, nEv.ContactEmail)
			assert.Equal(t, event.AccessRestricted, nEv.AccessRestricted)
			assert.Equal(t, event.Type, nEv.Type)
		}
	}
	// Test unknown event
	body, err = json.Marshal(types.UpdateEventRequest{
		Event: types.Event{
			Name:             "Test Event 4",
			Slug:             "event12",
			ContactEmail:     "email@test.com",
			AccessRestricted: false,
			Type:             "distance",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing unknown event.")
	request = httptest.NewRequest(http.MethodPost, "/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEvent(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test validation errors - invalid slug
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
	request = httptest.NewRequest(http.MethodPost, "/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test validation errors - no slug
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
	request = httptest.NewRequest(http.MethodPost, "/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test validation errors - invalid name
	body, err = json.Marshal(types.UpdateEventRequest{
		Event: types.Event{
			Name:             "Test Event@#$ 4",
			Slug:             "event2",
			ContactEmail:     "email@test.com",
			AccessRestricted: false,
			Type:             "distance",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing invalid name.")
	request = httptest.NewRequest(http.MethodPost, "/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test validation errors - no name
	body, err = json.Marshal(types.UpdateEventRequest{
		Event: types.Event{
			Name:             "",
			Slug:             "event2",
			ContactEmail:     "email@test.com",
			AccessRestricted: false,
			Type:             "distance",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing no name.")
	request = httptest.NewRequest(http.MethodPost, "/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test validation errors - invalid type
	body, err = json.Marshal(types.UpdateEventRequest{
		Event: types.Event{
			Name:             "Test Event 4",
			Slug:             "event2",
			ContactEmail:     "email@test.com",
			AccessRestricted: false,
			Type:             "invalid",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing invalid type.")
	request = httptest.NewRequest(http.MethodPost, "/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test validation errors - no type
	body, err = json.Marshal(types.UpdateEventRequest{
		Event: types.Event{
			Name:             "Test Event 4",
			Slug:             "event2",
			ContactEmail:     "email@test.com",
			AccessRestricted: false,
			Type:             "",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing no type.")
	request = httptest.NewRequest(http.MethodPost, "/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test validation errors - invalid website
	body, err = json.Marshal(types.UpdateEventRequest{
		Event: types.Event{
			Name:             "Test Event 4",
			Slug:             "event2",
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
	request = httptest.NewRequest(http.MethodPost, "/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test validation errors - invalid image
	body, err = json.Marshal(types.UpdateEventRequest{
		Event: types.Event{
			Name:             "Test Event 4",
			Slug:             "event2",
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
	request = httptest.NewRequest(http.MethodPost, "/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test validation errors - invalid email
	body, err = json.Marshal(types.UpdateEventRequest{
		Event: types.Event{
			Name:             "Test Event 4",
			Slug:             "event2",
			ContactEmail:     "emailtest.com",
			AccessRestricted: false,
			Type:             "distance",
		},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing invalid email.")
	request = httptest.NewRequest(http.MethodPost, "/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.UpdateEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
}

func TestDeleteEvent(t *testing.T) {
	// DELETE, /event/delete
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	body, err := json.Marshal(types.DeleteEventRequest{
		Slug: "event3",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	// Test no key
	t.Log("Testing no key given.")
	request := httptest.NewRequest(http.MethodPost, "/event/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodPost, "/event/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodPost, "/event/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test read key
	t.Log("Testing read host.")
	request = httptest.NewRequest(http.MethodPost, "/event/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test write key
	t.Log("Testing read host.")
	request = httptest.NewRequest(http.MethodPost, "/event/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test wrong account key
	body, err = json.Marshal(types.DeleteEventRequest{
		Slug: "event1",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing wrong account key.")
	request = httptest.NewRequest(http.MethodPost, "/event/update", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid host
	body, err = json.Marshal(types.DeleteEventRequest{
		Slug: "event32",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodPost, "/event/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodPost, "/event/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEvent(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/event/delete", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test empty request
	t.Log("Testing empty request.")
	request = httptest.NewRequest(http.MethodPost, "/event/delete", strings.NewReader(string("")))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEvent(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/event/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEvent(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test valid request
	body, err = json.Marshal(types.DeleteEventRequest{
		Slug: "event2",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing valid request.")
	request = httptest.NewRequest(http.MethodPost, "/event/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEvent(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		nEv, err := database.GetEvent("event3")
		if assert.NoError(t, err) {
			assert.Nil(t, nEv)
		}
	}
	// Test unknown event
	body, err = json.Marshal(types.DeleteEventRequest{
		Slug: "event12",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing unknown event.")
	request = httptest.NewRequest(http.MethodPost, "/event/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteEvent(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
}
