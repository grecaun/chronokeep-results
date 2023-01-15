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

func TestGetResults(t *testing.T) {
	// POST, /results
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	// Test no key
	t.Log("Testing no key given.")
	body, err := json.Marshal(types.GetResultsRequest{
		Slug: variables.events["event1"].Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/results", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.GetResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodPost, "/results", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodPost, "/results", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodPost, "/results", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/results", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetResults(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test no slug
	t.Log("Testing no slug.")
	request = httptest.NewRequest(http.MethodPost, "/results", strings.NewReader("{}"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetResults(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/results", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetResults(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodPost, "/results", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test a valid request with just a slug.
	t.Log("Testing slug only.")
	request = httptest.NewRequest(http.MethodPost, "/results", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 500, resp.Count)
			// number of distances
			assert.Equal(t, 2, len(resp.Results))
			assert.Equal(t, variables.events["event1"].Slug, resp.Event.Slug)
			assert.Equal(t, variables.events["event1"].Website, resp.Event.Website)
			assert.Equal(t, variables.events["event1"].Image, resp.Event.Image)
			assert.Equal(t, variables.events["event1"].ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, variables.events["event1"].AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, variables.events["event1"].Type, resp.Event.Type)
			assert.Equal(t, variables.events["event1"].RecentTime, resp.Event.RecentTime)
			assert.Equal(t, variables.eventYears["event1"]["2021"].Year, resp.EventYear.Year)
			assert.Equal(t, variables.eventYears["event1"]["2021"].DateTime.Local(), resp.EventYear.DateTime)
			assert.Equal(t, variables.eventYears["event1"]["2021"].Live, resp.EventYear.Live)
			assert.Equal(t, 2, len(resp.Years))
		}
	}
	// Test with slug, limits
	t.Log("Testing slug with limits/page.")
	var limit = 51
	var page = 1
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:  variables.events["event1"].Slug,
		Limit: &limit,
		Page:  &page,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	var testCorrect types.Result
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 51, resp.Count)
			testCorrect = resp.Results["1 Mile"][50]
		}
	}
	// verify pages work properly
	t.Log("Testing slug with limits/page - verifying second page returns next values.")
	limit = 50
	page = 2
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:  variables.events["event1"].Slug,
		Limit: &limit,
		Page:  &page,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 50, resp.Count)
			assert.Equal(t, testCorrect, resp.Results["1 Mile"][0])
		}
	}
	// Test with slug, distance
	t.Log("Testing with slug and distance defined.")
	distance := "1 Mile"
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:     variables.events["event1"].Slug,
		Distance: &distance,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 300, resp.Count)
			// number of distances
			assert.Equal(t, 1, len(resp.Results))
			assert.Equal(t, variables.events["event1"].Slug, resp.Event.Slug)
			assert.Equal(t, variables.events["event1"].Website, resp.Event.Website)
			assert.Equal(t, variables.events["event1"].Image, resp.Event.Image)
			assert.Equal(t, variables.events["event1"].ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, variables.events["event1"].AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, variables.events["event1"].Type, resp.Event.Type)
			assert.Equal(t, variables.events["event1"].RecentTime, resp.Event.RecentTime)
			assert.Equal(t, variables.eventYears["event1"]["2021"].Year, resp.EventYear.Year)
			assert.Equal(t, variables.eventYears["event1"]["2021"].DateTime.Local(), resp.EventYear.DateTime)
			assert.Equal(t, variables.eventYears["event1"]["2021"].Live, resp.EventYear.Live)
			assert.Equal(t, 2, len(resp.Years))
		}
	}
	// Test with slug, distance, limits
	t.Log("Testing with slug and distance defined as well as limit/page.")
	limit = 51
	page = 1
	distance = "2 Mile"
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:     variables.events["event1"].Slug,
		Limit:    &limit,
		Page:     &page,
		Distance: &distance,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 51, resp.Count)
			testCorrect = resp.Results["2 Mile"][50]
		}
	}
	// verify pages work properly
	t.Log("Testing with slug and distance defined as well as limit/page - verifying page correctness.")
	limit = 50
	page = 2
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:     variables.events["event1"].Slug,
		Limit:    &limit,
		Page:     &page,
		Distance: &distance,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 50, resp.Count)
			assert.Equal(t, testCorrect, resp.Results["2 Mile"][0])
		}
	}
	// Test with slug and year
	t.Log("Testing with slug and year.")
	year := variables.eventYears["event1"]["2020"].Year
	body, err = json.Marshal(types.GetResultsRequest{
		Slug: variables.events["event1"].Slug,
		Year: &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 4, resp.Count)
			// number of distances
			assert.Equal(t, 2, len(resp.Results))
			assert.Equal(t, variables.events["event1"].Slug, resp.Event.Slug)
			assert.Equal(t, variables.events["event1"].Website, resp.Event.Website)
			assert.Equal(t, variables.events["event1"].Image, resp.Event.Image)
			assert.Equal(t, variables.events["event1"].ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, variables.events["event1"].AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, variables.events["event1"].Type, resp.Event.Type)
			assert.Equal(t, variables.events["event1"].RecentTime, resp.Event.RecentTime)
			assert.Equal(t, variables.eventYears["event1"]["2020"].Year, resp.EventYear.Year)
			assert.Equal(t, variables.eventYears["event1"]["2020"].DateTime.Local(), resp.EventYear.DateTime)
			assert.Equal(t, variables.eventYears["event1"]["2020"].Live, resp.EventYear.Live)
			assert.Equal(t, 2, len(resp.Years))
		}
	}
	// Test with slug, year and limits
	t.Log("Testing with slug and year as well as limit/page.")
	limit = 51
	page = 1
	year = "2021"
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:  variables.events["event1"].Slug,
		Limit: &limit,
		Page:  &page,
		Year:  &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 51, resp.Count)
			testCorrect = resp.Results["1 Mile"][50]
		}
	}
	// verify pages work properly
	t.Log("Testing with slug and year as well as limit/page - verifying page correctness.")
	limit = 50
	page = 2
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:  variables.events["event1"].Slug,
		Limit: &limit,
		Page:  &page,
		Year:  &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 50, resp.Count)
			assert.Equal(t, testCorrect, resp.Results["1 Mile"][0])
		}
	}
	// Test with slug, year, distance
	t.Log("Testing with slug, year, and distance.")
	distance = "2 Mile"
	year = "2020"
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:     variables.events["event1"].Slug,
		Year:     &year,
		Distance: &distance,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 1, resp.Count)
			// number of distances
			assert.Equal(t, 1, len(resp.Results))
			assert.Equal(t, variables.events["event1"].Slug, resp.Event.Slug)
			assert.Equal(t, variables.events["event1"].Website, resp.Event.Website)
			assert.Equal(t, variables.events["event1"].Image, resp.Event.Image)
			assert.Equal(t, variables.events["event1"].ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, variables.events["event1"].AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, variables.events["event1"].Type, resp.Event.Type)
			assert.Equal(t, variables.events["event1"].RecentTime, resp.Event.RecentTime)
			assert.Equal(t, variables.eventYears["event1"]["2020"].Year, resp.EventYear.Year)
			assert.Equal(t, variables.eventYears["event1"]["2020"].DateTime.Local(), resp.EventYear.DateTime)
			assert.Equal(t, variables.eventYears["event1"]["2020"].Live, resp.EventYear.Live)
			assert.Equal(t, 2, len(resp.Years))
		}
	}
	// Test with slug, year, distance, limits
	t.Log("Testing with slug, year, and distance as well as limit/page.")
	limit = 51
	page = 1
	distance = "2 Mile"
	year = "2021"
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:     variables.events["event1"].Slug,
		Limit:    &limit,
		Page:     &page,
		Year:     &year,
		Distance: &distance,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 51, resp.Count)
			testCorrect = resp.Results["2 Mile"][50]
		}
	}
	// verify pages work properly
	t.Log("Testing with slug, year, and distance as well as limit/page - verifying page correctness.")
	limit = 50
	page = 2
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:     variables.events["event1"].Slug,
		Limit:    &limit,
		Page:     &page,
		Year:     &year,
		Distance: &distance,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 50, resp.Count)
			assert.Equal(t, testCorrect, resp.Results["2 Mile"][0])
		}
	}
	// Test invalid event
	t.Log("Testing event not found.")
	body, err = json.Marshal(types.GetResultsRequest{
		Slug: "invalid event",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetResults(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid event invalid year
	year = "2000"
	body, err = json.Marshal(types.GetResultsRequest{
		Slug: "invalid event",
		Year: &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetResults(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid event, valid year, invalid distance
	year = "2021"
	distance = "invalid-distance"
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:     "invalid event",
		Year:     &year,
		Distance: &distance,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetResults(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid key with restricted event
	t.Log("Testing restricted event but unauthorized key.")
	body, err = json.Marshal(types.GetResultsRequest{
		Slug: variables.events["event2"].Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}

func TestGetAllResults(t *testing.T) {
	// POST, /results/all
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	// Test no key
	t.Log("Testing no key given.")
	body, err := json.Marshal(types.GetResultsRequest{
		Slug: variables.events["event1"].Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/results/all", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.GetAllResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodPost, "/results/all", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAllResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodPost, "/results/all", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAllResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodPost, "/results/all", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAllResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/results/all", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAllResults(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test no slug
	t.Log("Testing no slug.")
	request = httptest.NewRequest(http.MethodPost, "/results/all", strings.NewReader("{}"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAllResults(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/results/all", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAllResults(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodPost, "/results/all", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAllResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test a valid request with just a slug.
	t.Log("Testing slug only.")
	request = httptest.NewRequest(http.MethodPost, "/results/all", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAllResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 800, resp.Count)
			// number of distances
			assert.Equal(t, 2, len(resp.Results))
			assert.Equal(t, variables.events["event1"].Slug, resp.Event.Slug)
			assert.Equal(t, variables.events["event1"].Website, resp.Event.Website)
			assert.Equal(t, variables.events["event1"].Image, resp.Event.Image)
			assert.Equal(t, variables.events["event1"].ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, variables.events["event1"].AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, variables.events["event1"].Type, resp.Event.Type)
			assert.Equal(t, variables.events["event1"].RecentTime, resp.Event.RecentTime)
			assert.Equal(t, variables.eventYears["event1"]["2021"].Year, resp.EventYear.Year)
			assert.Equal(t, variables.eventYears["event1"]["2021"].DateTime.Local(), resp.EventYear.DateTime)
			assert.Equal(t, variables.eventYears["event1"]["2021"].Live, resp.EventYear.Live)
			assert.Equal(t, 2, len(resp.Years))
		}
	}
	// Test with slug, limits
	t.Log("Testing slug with limits/page.")
	var limit = 51
	var page = 1
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:  variables.events["event1"].Slug,
		Limit: &limit,
		Page:  &page,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/all", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	var testCorrect types.Result
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAllResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 51, resp.Count)
			testCorrect = resp.Results["1 Mile"][50]
		}
	}
	// verify pages work properly
	t.Log("Testing slug with limits/page - verifying second page returns next values.")
	limit = 50
	page = 2
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:  variables.events["event1"].Slug,
		Limit: &limit,
		Page:  &page,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/all", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAllResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 50, resp.Count)
			assert.Equal(t, testCorrect, resp.Results["1 Mile"][0])
		}
	}
	// Test with slug, distance
	t.Log("Testing with slug and distance defined.")
	distance := "1 Mile"
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:     variables.events["event1"].Slug,
		Distance: &distance,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/all", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAllResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 600, resp.Count)
			// number of distances
			assert.Equal(t, 1, len(resp.Results))
			assert.Equal(t, variables.events["event1"].Slug, resp.Event.Slug)
			assert.Equal(t, variables.events["event1"].Website, resp.Event.Website)
			assert.Equal(t, variables.events["event1"].Image, resp.Event.Image)
			assert.Equal(t, variables.events["event1"].ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, variables.events["event1"].AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, variables.events["event1"].Type, resp.Event.Type)
			assert.Equal(t, variables.events["event1"].RecentTime, resp.Event.RecentTime)
			assert.Equal(t, variables.eventYears["event1"]["2021"].Year, resp.EventYear.Year)
			assert.Equal(t, variables.eventYears["event1"]["2021"].DateTime.Local(), resp.EventYear.DateTime)
			assert.Equal(t, variables.eventYears["event1"]["2021"].Live, resp.EventYear.Live)
			assert.Equal(t, 2, len(resp.Years))
		}
	}
	// Test with slug, distance, limits
	t.Log("Testing with slug and distance defined as well as limit/page.")
	limit = 51
	page = 1
	distance = "2 Mile"
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:     variables.events["event1"].Slug,
		Limit:    &limit,
		Page:     &page,
		Distance: &distance,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/all", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAllResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 51, resp.Count)
			testCorrect = resp.Results["2 Mile"][50]
		}
	}
	// verify pages work properly
	t.Log("Testing with slug and distance defined as well as limit/page - verifying page correctness.")
	limit = 50
	page = 2
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:     variables.events["event1"].Slug,
		Limit:    &limit,
		Page:     &page,
		Distance: &distance,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/all", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAllResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 50, resp.Count)
			assert.Equal(t, testCorrect, resp.Results["2 Mile"][0])
		}
	}
	// Test with slug and year
	t.Log("Testing with slug and year.")
	year := variables.eventYears["event1"]["2020"].Year
	body, err = json.Marshal(types.GetResultsRequest{
		Slug: variables.events["event1"].Slug,
		Year: &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/all", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAllResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 5, resp.Count)
			// number of distances
			assert.Equal(t, 2, len(resp.Results))
			assert.Equal(t, variables.events["event1"].Slug, resp.Event.Slug)
			assert.Equal(t, variables.events["event1"].Website, resp.Event.Website)
			assert.Equal(t, variables.events["event1"].Image, resp.Event.Image)
			assert.Equal(t, variables.events["event1"].ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, variables.events["event1"].AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, variables.events["event1"].Type, resp.Event.Type)
			assert.Equal(t, variables.events["event1"].RecentTime, resp.Event.RecentTime)
			assert.Equal(t, variables.eventYears["event1"]["2020"].Year, resp.EventYear.Year)
			assert.Equal(t, variables.eventYears["event1"]["2020"].DateTime.Local(), resp.EventYear.DateTime)
			assert.Equal(t, variables.eventYears["event1"]["2020"].Live, resp.EventYear.Live)
			assert.Equal(t, 2, len(resp.Years))
		}
	}
	// Test with slug, year and limits
	t.Log("Testing with slug and year as well as limit/page.")
	limit = 51
	page = 1
	year = "2021"
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:  variables.events["event1"].Slug,
		Limit: &limit,
		Page:  &page,
		Year:  &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/all", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAllResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 51, resp.Count)
			testCorrect = resp.Results["1 Mile"][50]
		}
	}
	// verify pages work properly
	t.Log("Testing with slug and year as well as limit/page - verifying page correctness.")
	limit = 50
	page = 2
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:  variables.events["event1"].Slug,
		Limit: &limit,
		Page:  &page,
		Year:  &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/all", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAllResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 50, resp.Count)
			assert.Equal(t, testCorrect, resp.Results["1 Mile"][0])
		}
	}
	// Test with slug, year, distance
	t.Log("Testing with slug, year, and distance.")
	distance = "2 Mile"
	year = "2020"
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:     variables.events["event1"].Slug,
		Year:     &year,
		Distance: &distance,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/all", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAllResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 2, resp.Count)
			// number of distances
			assert.Equal(t, 1, len(resp.Results))
			assert.Equal(t, variables.events["event1"].Slug, resp.Event.Slug)
			assert.Equal(t, variables.events["event1"].Website, resp.Event.Website)
			assert.Equal(t, variables.events["event1"].Image, resp.Event.Image)
			assert.Equal(t, variables.events["event1"].ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, variables.events["event1"].AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, variables.events["event1"].Type, resp.Event.Type)
			assert.Equal(t, variables.events["event1"].RecentTime, resp.Event.RecentTime)
			assert.Equal(t, variables.eventYears["event1"]["2020"].Year, resp.EventYear.Year)
			assert.Equal(t, variables.eventYears["event1"]["2020"].DateTime.Local(), resp.EventYear.DateTime)
			assert.Equal(t, variables.eventYears["event1"]["2020"].Live, resp.EventYear.Live)
			assert.Equal(t, 2, len(resp.Years))
		}
	}
	// Test with slug, year, distance, limits
	t.Log("Testing with slug, year, and distance as well as limit/page.")
	limit = 51
	page = 1
	distance = "2 Mile"
	year = "2021"
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:     variables.events["event1"].Slug,
		Limit:    &limit,
		Page:     &page,
		Year:     &year,
		Distance: &distance,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/all", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAllResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 51, resp.Count)
			testCorrect = resp.Results["2 Mile"][50]
		}
	}
	// verify pages work properly
	t.Log("Testing with slug, year, and distance as well as limit/page - verifying page correctness.")
	limit = 50
	page = 2
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:     variables.events["event1"].Slug,
		Limit:    &limit,
		Page:     &page,
		Year:     &year,
		Distance: &distance,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/all", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAllResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 50, resp.Count)
			assert.Equal(t, testCorrect, resp.Results["2 Mile"][0])
		}
	}
	// Test invalid event
	t.Log("Testing event not found.")
	body, err = json.Marshal(types.GetResultsRequest{
		Slug: "invalid event",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/all", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAllResults(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid event invalid year
	year = "2000"
	body, err = json.Marshal(types.GetResultsRequest{
		Slug: "invalid event",
		Year: &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/all", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAllResults(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid event, valid year, invalid distance
	year = "2021"
	distance = "invalid-distance"
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:     "invalid event",
		Year:     &year,
		Distance: &distance,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/all", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAllResults(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid key with restricted event
	t.Log("Testing restricted event but unauthorized key.")
	body, err = json.Marshal(types.GetResultsRequest{
		Slug: variables.events["event2"].Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/all", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetAllResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}

func TestGetFinishResults(t *testing.T) {
	// POST, /results/finish
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	// Test no key
	t.Log("Testing no key given.")
	body, err := json.Marshal(types.GetResultsRequest{
		Slug: variables.events["event1"].Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/results/finish", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.GetFinishResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodPost, "/results/finish", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetFinishResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodPost, "/results/finish", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetFinishResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodPost, "/results/finish", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetFinishResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/results/finish", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetFinishResults(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test no slug
	t.Log("Testing no slug.")
	request = httptest.NewRequest(http.MethodPost, "/results/finish", strings.NewReader("{}"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetFinishResults(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/results/finish", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetFinishResults(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodPost, "/results/finish", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetFinishResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test a valid request with just a slug.
	t.Log("Testing slug only.")
	request = httptest.NewRequest(http.MethodPost, "/results/finish", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetFinishResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 500, resp.Count)
			// number of distances
			assert.Equal(t, 2, len(resp.Results))
			assert.Equal(t, variables.events["event1"].Slug, resp.Event.Slug)
			assert.Equal(t, variables.events["event1"].Website, resp.Event.Website)
			assert.Equal(t, variables.events["event1"].Image, resp.Event.Image)
			assert.Equal(t, variables.events["event1"].ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, variables.events["event1"].AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, variables.events["event1"].Type, resp.Event.Type)
			assert.Equal(t, variables.events["event1"].RecentTime, resp.Event.RecentTime)
			assert.Equal(t, variables.eventYears["event1"]["2021"].Year, resp.EventYear.Year)
			assert.Equal(t, variables.eventYears["event1"]["2021"].DateTime.Local(), resp.EventYear.DateTime)
			assert.Equal(t, variables.eventYears["event1"]["2021"].Live, resp.EventYear.Live)
			assert.Equal(t, 2, len(resp.Years))
		}
	}
	// Test with slug, limits
	t.Log("Testing slug with limits/page.")
	var limit = 51
	var page = 1
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:  variables.events["event1"].Slug,
		Limit: &limit,
		Page:  &page,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/finish", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	var testCorrect types.Result
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetFinishResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 51, resp.Count)
			testCorrect = resp.Results["1 Mile"][50]
		}
	}
	// verify pages work properly
	t.Log("Testing slug with limits/page - verifying second page returns next values.")
	limit = 50
	page = 2
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:  variables.events["event1"].Slug,
		Limit: &limit,
		Page:  &page,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/finish", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetFinishResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 50, resp.Count)
			assert.Equal(t, testCorrect, resp.Results["1 Mile"][0])
		}
	}
	// Test with slug, distance
	t.Log("Testing with slug and distance defined.")
	distance := "1 Mile"
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:     variables.events["event1"].Slug,
		Distance: &distance,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/finish", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetFinishResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 300, resp.Count)
			// number of distances
			assert.Equal(t, 1, len(resp.Results))
			assert.Equal(t, variables.events["event1"].Slug, resp.Event.Slug)
			assert.Equal(t, variables.events["event1"].Website, resp.Event.Website)
			assert.Equal(t, variables.events["event1"].Image, resp.Event.Image)
			assert.Equal(t, variables.events["event1"].ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, variables.events["event1"].AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, variables.events["event1"].Type, resp.Event.Type)
			assert.Equal(t, variables.events["event1"].RecentTime, resp.Event.RecentTime)
			assert.Equal(t, variables.eventYears["event1"]["2021"].Year, resp.EventYear.Year)
			assert.Equal(t, variables.eventYears["event1"]["2021"].DateTime.Local(), resp.EventYear.DateTime)
			assert.Equal(t, variables.eventYears["event1"]["2021"].Live, resp.EventYear.Live)
			assert.Equal(t, 2, len(resp.Years))
		}
	}
	// Test with slug, distance, limits
	t.Log("Testing with slug and distance defined as well as limit/page.")
	limit = 51
	page = 1
	distance = "2 Mile"
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:     variables.events["event1"].Slug,
		Limit:    &limit,
		Page:     &page,
		Distance: &distance,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/finish", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetFinishResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 51, resp.Count)
			testCorrect = resp.Results["2 Mile"][50]
		}
	}
	// verify pages work properly
	t.Log("Testing with slug and distance defined as well as limit/page - verifying page correctness.")
	limit = 50
	page = 2
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:     variables.events["event1"].Slug,
		Limit:    &limit,
		Page:     &page,
		Distance: &distance,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/finish", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetFinishResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 50, resp.Count)
			assert.Equal(t, testCorrect, resp.Results["2 Mile"][0])
		}
	}
	// Test with slug and year
	t.Log("Testing with slug and year.")
	year := variables.eventYears["event1"]["2020"].Year
	body, err = json.Marshal(types.GetResultsRequest{
		Slug: variables.events["event1"].Slug,
		Year: &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/finish", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetFinishResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 4, resp.Count)
			// number of distances
			assert.Equal(t, 2, len(resp.Results))
			assert.Equal(t, variables.events["event1"].Slug, resp.Event.Slug)
			assert.Equal(t, variables.events["event1"].Website, resp.Event.Website)
			assert.Equal(t, variables.events["event1"].Image, resp.Event.Image)
			assert.Equal(t, variables.events["event1"].ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, variables.events["event1"].AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, variables.events["event1"].Type, resp.Event.Type)
			assert.Equal(t, variables.events["event1"].RecentTime, resp.Event.RecentTime)
			assert.Equal(t, variables.eventYears["event1"]["2020"].Year, resp.EventYear.Year)
			assert.Equal(t, variables.eventYears["event1"]["2020"].DateTime.Local(), resp.EventYear.DateTime)
			assert.Equal(t, variables.eventYears["event1"]["2020"].Live, resp.EventYear.Live)
			assert.Equal(t, 2, len(resp.Years))
		}
	}
	// Test with slug, year and limits
	t.Log("Testing with slug and year as well as limit/page.")
	limit = 51
	page = 1
	year = "2021"
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:  variables.events["event1"].Slug,
		Limit: &limit,
		Page:  &page,
		Year:  &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/finish", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetFinishResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 51, resp.Count)
			testCorrect = resp.Results["1 Mile"][50]
		}
	}
	// verify pages work properly
	t.Log("Testing with slug and year as well as limit/page - verifying page correctness.")
	limit = 50
	page = 2
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:  variables.events["event1"].Slug,
		Limit: &limit,
		Page:  &page,
		Year:  &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/finish", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetFinishResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 50, resp.Count)
			assert.Equal(t, testCorrect, resp.Results["1 Mile"][0])
		}
	}
	// Test with slug, year, distance
	t.Log("Testing with slug, year, and distance.")
	distance = "2 Mile"
	year = "2020"
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:     variables.events["event1"].Slug,
		Year:     &year,
		Distance: &distance,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/finish", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetFinishResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 1, resp.Count)
			// number of distances
			assert.Equal(t, 1, len(resp.Results))
			assert.Equal(t, variables.events["event1"].Slug, resp.Event.Slug)
			assert.Equal(t, variables.events["event1"].Website, resp.Event.Website)
			assert.Equal(t, variables.events["event1"].Image, resp.Event.Image)
			assert.Equal(t, variables.events["event1"].ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, variables.events["event1"].AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, variables.events["event1"].Type, resp.Event.Type)
			assert.Equal(t, variables.events["event1"].RecentTime, resp.Event.RecentTime)
			assert.Equal(t, variables.eventYears["event1"]["2020"].Year, resp.EventYear.Year)
			assert.Equal(t, variables.eventYears["event1"]["2020"].DateTime.Local(), resp.EventYear.DateTime)
			assert.Equal(t, variables.eventYears["event1"]["2020"].Live, resp.EventYear.Live)
			assert.Equal(t, 2, len(resp.Years))
		}
	}
	// Test with slug, year, distance, limits
	t.Log("Testing with slug, year, and distance as well as limit/page.")
	limit = 51
	page = 1
	distance = "2 Mile"
	year = "2021"
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:     variables.events["event1"].Slug,
		Limit:    &limit,
		Page:     &page,
		Year:     &year,
		Distance: &distance,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/finish", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetFinishResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 51, resp.Count)
			testCorrect = resp.Results["2 Mile"][50]
		}
	}
	// verify pages work properly
	t.Log("Testing with slug, year, and distance as well as limit/page - verifying page correctness.")
	limit = 50
	page = 2
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:     variables.events["event1"].Slug,
		Limit:    &limit,
		Page:     &page,
		Year:     &year,
		Distance: &distance,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/finish", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetFinishResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 50, resp.Count)
			assert.Equal(t, testCorrect, resp.Results["2 Mile"][0])
		}
	}
	// Test invalid event
	t.Log("Testing event not found.")
	body, err = json.Marshal(types.GetResultsRequest{
		Slug: "invalid event",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/finish", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetFinishResults(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid event invalid year
	year = "2000"
	body, err = json.Marshal(types.GetResultsRequest{
		Slug: "invalid event",
		Year: &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/finish", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetFinishResults(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid event, valid year, invalid distance
	year = "2021"
	distance = "invalid-distance"
	body, err = json.Marshal(types.GetResultsRequest{
		Slug:     "invalid event",
		Year:     &year,
		Distance: &distance,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/finish", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetFinishResults(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid key with restricted event
	t.Log("Testing restricted event but unauthorized key.")
	body, err = json.Marshal(types.GetResultsRequest{
		Slug: variables.events["event2"].Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/finish", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetFinishResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}

func TestGetBibResults(t *testing.T) {
	// POST, /results/bib
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	// Test no key
	t.Log("Testing no key given.")
	body, err := json.Marshal(types.GetBibResultsRequest{
		Slug: variables.events["event1"].Slug,
		Year: variables.eventYears["event1"]["2021"].Year,
		Bib:  "0",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/results/bib", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.GetBibResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodPost, "/results/bib", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetBibResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodPost, "/results/bib", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetBibResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodPost, "/results/bib", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetBibResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/results/bib", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetBibResults(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test no slug
	t.Log("Testing no slug.")
	request = httptest.NewRequest(http.MethodPost, "/results/bib", strings.NewReader("{}"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetBibResults(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/results/bib", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetBibResults(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodPost, "/results/bib", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetBibResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test a valid request
	t.Log("Testing slug only.")
	request = httptest.NewRequest(http.MethodPost, "/results/bib", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetBibResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetBibResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 2, len(resp.Results))
			assert.Equal(t, variables.events["event1"].Slug, resp.Event.Slug)
			assert.Equal(t, variables.events["event1"].Website, resp.Event.Website)
			assert.Equal(t, variables.events["event1"].Image, resp.Event.Image)
			assert.Equal(t, variables.events["event1"].ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, variables.events["event1"].AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, variables.events["event1"].Type, resp.Event.Type)
			assert.Equal(t, variables.events["event1"].RecentTime, resp.Event.RecentTime)
			assert.Equal(t, variables.eventYears["event1"]["2021"].Year, resp.EventYear.Year)
			assert.Equal(t, variables.eventYears["event1"]["2021"].DateTime.Local(), resp.EventYear.DateTime)
			assert.Equal(t, variables.eventYears["event1"]["2021"].Live, resp.EventYear.Live)
			assert.NotNil(t, resp.Person)
			// first result
			assert.Equal(t, variables.results["event1"]["2021"][0], resp.Results[0])
			// second result
			assert.Equal(t, variables.results["event1"]["2021"][300], resp.Results[1])
		}
	}
	// Test invalid event
	t.Log("Testing invalid event, valid year/bib.")
	body, err = json.Marshal(types.GetBibResultsRequest{
		Slug: "invalid-event",
		Year: "2021",
		Bib:  "0",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/bib", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetBibResults(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid event, invalid year
	t.Log("Testing valid event/bib, invalid year.")
	body, err = json.Marshal(types.GetBibResultsRequest{
		Slug: variables.events["event1"].Slug,
		Year: "invalid-year",
		Bib:  "0",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/bib", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetBibResults(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid event invalid bib
	t.Log("Testing valid event/year, invalid bib.")
	body, err = json.Marshal(types.GetBibResultsRequest{
		Slug: variables.events["event1"].Slug,
		Year: variables.eventYears["event1"]["2021"].Year,
		Bib:  "invalid-bib",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/bib", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetBibResults(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid key with restricted event
	t.Log("Testing valid key with restricted event.")
	body, err = json.Marshal(types.GetBibResultsRequest{
		Slug: variables.events["event2"].Slug,
		Year: variables.eventYears["event1"]["2021"].Year,
		Bib:  "100",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/bib", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetBibResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}

func TestAddResults(t *testing.T) {
	// POST, /results/add
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	// add eventYear to upload to
	eventYear, err := database.AddEventYear(types.EventYear{
		EventIdentifier: variables.events["event2"].Identifier,
		Year:            "2023",
		DateTime:        time.Date(2023, 04, 05, 11, 0, 0, 0, time.Local),
		Live:            false,
	})
	assert.NoError(t, err)
	// Test no key
	results := variables.results["event2"]["2021"]
	t.Log("Testing no key given.")
	body, err := json.Marshal(types.AddResultsRequest{
		Slug:    variables.events["event2"].Slug,
		Year:    "2023",
		Results: results,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/results/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.AddResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodPost, "/results/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodPost, "/results/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test read key (vs delete/write)
	t.Log("Testing read key.")
	request = httptest.NewRequest(http.MethodPost, "/results/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test different account's key
	t.Log("Testing different account's key.")
	request = httptest.NewRequest(http.MethodPost, "/results/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodPost, "/results/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/results/add", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddResults(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test without slug/other info
	t.Log("Testing no slug.")
	request = httptest.NewRequest(http.MethodPost, "/results/add", strings.NewReader("{}"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddResults(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/results/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddResults(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test a valid request with a write key
	t.Log("Testing valid write key.")
	request = httptest.NewRequest(http.MethodPost, "/results/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.AddResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(results), resp.Count)
		}
	}
	// verify add results adds information
	t.Log("Verifying information was added.")
	uploaded, err := database.GetResults(eventYear.Identifier, 0, 0)
	if assert.NoError(t, err) {
		found := 0
		for _, res := range uploaded {
			for _, inner := range results {
				if res == inner {
					found++
				}
			}
		}
		assert.Equal(t, found, len(results))
	}
	for ix := range results {
		results[ix].Seconds = results[ix].Seconds + 2
	}
	// Test a valid request with a delete key
	t.Log("Testing valid write key.")
	body, err = json.Marshal(types.AddResultsRequest{
		Slug:    variables.events["event2"].Slug,
		Year:    "2023",
		Results: results,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.AddResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(results), resp.Count)
		}
	}
	t.Log("Verifying information was updated.")
	uploaded, err = database.GetResults(eventYear.Identifier, 0, 0)
	if assert.NoError(t, err) {
		found := 0
		for _, res := range uploaded {
			for _, inner := range results {
				if res == inner {
					found++
				}
			}
		}
		assert.Equal(t, found, len(results))
	}
	// Test invalid event
	t.Log("Testing invalid event.")
	body, err = json.Marshal(types.AddResultsRequest{
		Slug:    "invalid event",
		Year:    "2023",
		Results: results,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddResults(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid event invalid year
	t.Log("Testing invalid event.")
	body, err = json.Marshal(types.AddResultsRequest{
		Slug:    variables.events["event2"].Slug,
		Year:    "invalid-year",
		Results: results,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddResults(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodPost, "/results/add", strings.NewReader(""))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test validation - Bib
	t.Log("Test validation check - Bib.")
	results[0].Bib = ""
	body, err = json.Marshal(types.AddResultsRequest{
		Slug:    variables.events["event2"].Slug,
		Year:    "2023",
		Results: results,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.AddResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(results)-1, resp.Count)
		}
	}
	// Test validation - First
	t.Log("Test validation check - First.")
	results[0].Bib = "55235"
	results[0].First = ""
	body, err = json.Marshal(types.AddResultsRequest{
		Slug:    variables.events["event2"].Slug,
		Year:    "2023",
		Results: results,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.AddResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(results), resp.Count)
		}
	}
	// Test validation - Last
	t.Log("Test validation check - Last.")
	results[0].First = "John"
	results[0].Last = ""
	body, err = json.Marshal(types.AddResultsRequest{
		Slug:    variables.events["event2"].Slug,
		Year:    "2023",
		Results: results,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.AddResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(results), resp.Count)
		}
	}
	// Test validation - Age1
	t.Log("Test validation check - Age1.")
	results[0].Last = "Smith-Johnson"
	results[0].Age = -1
	body, err = json.Marshal(types.AddResultsRequest{
		Slug:    variables.events["event2"].Slug,
		Year:    "2023",
		Results: results,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.AddResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(results)-1, resp.Count)
		}
	}
	// Test validation - Age2
	t.Log("Test validation check - Age2.")
	results[0].Age = 135
	body, err = json.Marshal(types.AddResultsRequest{
		Slug:    variables.events["event2"].Slug,
		Year:    "2023",
		Results: results,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.AddResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(results)-1, resp.Count)
		}
	}
	// Test validation - Seconds
	t.Log("Test validation check - Seconds.")
	results[0].Age = 13
	results[0].Seconds = -100
	body, err = json.Marshal(types.AddResultsRequest{
		Slug:    variables.events["event2"].Slug,
		Year:    "2023",
		Results: results,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.AddResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(results)-1, resp.Count)
		}
	}
	// Test validation - ChipSeconds
	t.Log("Test validation check - ChipSeconds.")
	results[0].Seconds = 300
	results[0].ChipSeconds = -51
	body, err = json.Marshal(types.AddResultsRequest{
		Slug:    variables.events["event2"].Slug,
		Year:    "2023",
		Results: results,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.AddResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(results)-1, resp.Count)
		}
	}
	// Test validation - Location
	t.Log("Test validation check - Location.")
	results[0].ChipSeconds = -51
	results[0].Location = ""
	body, err = json.Marshal(types.AddResultsRequest{
		Slug:    variables.events["event2"].Slug,
		Year:    "2023",
		Results: results,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.AddResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(results)-1, resp.Count)
		}
	}
	// Test validation - Occurence
	t.Log("Test validation check - Occurence.")
	results[0].Location = "Start/Finish"
	results[0].Occurence = -1
	body, err = json.Marshal(types.AddResultsRequest{
		Slug:    variables.events["event2"].Slug,
		Year:    "2023",
		Results: results,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.AddResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(results)-1, resp.Count)
		}
	}
	// Test validation - Gender
	t.Log("Test validation check - Gender.")
	results[0].Occurence = 1
	results[0].Gender = ""
	body, err = json.Marshal(types.AddResultsRequest{
		Slug:    variables.events["event2"].Slug,
		Year:    "2023",
		Results: results,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.AddResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(results)-1, resp.Count)
		}
	}
}
func TestDeleteResults(t *testing.T) {
	// DELETE, /results/delete
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	// Test no key
	t.Log("Testing no key given.")
	year := variables.eventYears["event2"]["2021"].Year
	body, err := json.Marshal(types.GetResultsRequest{
		Slug: variables.events["event2"].Slug,
		Year: &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request := httptest.NewRequest(http.MethodDelete, "/results/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.DeleteResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodDelete, "/results/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodDelete, "/results/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test read key
	t.Log("Testing read key.")
	request = httptest.NewRequest(http.MethodDelete, "/results/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test different account's key
	t.Log("Testing different account's key.")
	request = httptest.NewRequest(http.MethodDelete, "/results/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodDelete, "/results/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodDelete, "/results/delete", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteResults(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test without slug/other info
	t.Log("Testing no slug.")
	request = httptest.NewRequest(http.MethodDelete, "/results/delete", strings.NewReader("{}"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteResults(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodDelete, "/results/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteResults(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// verify information is there to delete
	t.Log("Verifying information before deletion.")
	deleted, err := database.GetResults(variables.eventYears["event2"]["2021"].Identifier, 0, 0)
	if assert.NoError(t, err) {
		assert.Equal(t, len(variables.results["event2"]["2021"]), len(deleted))
	}
	// Test write key
	t.Log("Testing write key.")
	request = httptest.NewRequest(http.MethodDelete, "/results/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.AddResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(deleted), resp.Count)
		}
	}
	// verify delete remotes information
	t.Log("Verifying information was deleted.")
	deleted, err = database.GetResults(variables.eventYears["event2"]["2021"].Identifier, 0, 0)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(deleted))
	}
	// verify information is there to delete
	t.Log("Verifying information before deletion.")
	deleted, err = database.GetResults(variables.eventYears["event2"]["2020"].Identifier, 0, 0)
	if assert.NoError(t, err) {
		assert.Equal(t, len(variables.results["event2"]["2020"]), len(deleted))
	}
	// Test a delete request
	year = variables.eventYears["event2"]["2020"].Year
	body, err = json.Marshal(types.GetResultsRequest{
		Slug: variables.events["event2"].Slug,
		Year: &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing delete key.")
	request = httptest.NewRequest(http.MethodDelete, "/results/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.AddResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(deleted), resp.Count)
		}
	}
	// verify delete remotes information
	t.Log("Verifying information was deleted.")
	deleted, err = database.GetResults(variables.eventYears["event2"]["2020"].Identifier, 0, 0)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(deleted))
	}
	// Test invalid event
	t.Log("Testing invalid event.")
	body, err = json.Marshal(types.GetResultsRequest{
		Slug: "invalid event",
		Year: &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/results/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteResults(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid event invalid year
	t.Log("Testing invalid event.")
	year = "invalid-year"
	body, err = json.Marshal(types.GetResultsRequest{
		Slug: variables.events["event2"].Slug,
		Year: &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/results/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteResults(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodDelete, "/results/delete", strings.NewReader(""))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}
