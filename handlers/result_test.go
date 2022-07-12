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
	// Test no key
	// Test expired key
	// Test invalid key
	// Test invalid authorization header
	// Test bad request
	// Test wrong content type
	// Test invalid host
	// Test a valid request with just a slug.
	// Test with slug, limits
	// Test with slug, distance
	// Test with slug, distance, limits
	// Test with slug and year
	// Test with slug, year and limits
	// Test with slug, year, distance
	// Test with slug, year, distance, limits
	// Test invalid event
	// Test valid event invalid year
	// Test valid event, valid year, invalid distance
	// Test valid key with restricted event
}

func TestGetBibResults(t *testing.T) {
	// POST, /results/bib
	// Test no key
	// Test expired key
	// Test invalid key
	// Test invalid authorization header
	// Test bad request
	// Test wrong content type
	// Test invalid host
	// Test a valid request with just a slug.
	// Test with slug, limits
	// Test with slug, distance
	// Test with slug, distance, limits
	// Test with slug and year
	// Test with slug, year and limits
	// Test with slug, year, distance
	// Test with slug, year, distance, limits
	// Test invalid event
	// Test valid event invalid year
	// Test valid event, valid year, invalid distance
	// Test valid key with restricted event
}

func TestAddResults(t *testing.T) {
	// POST, /results/add
	// Test no key
	// Test expired key
	// Test invalid key
	// Test invalid authorization header
	// Test bad request
	// Test wrong content type
	// Test invalid host
	// Test a valid request with just a slug.
	// Test with slug, limits
	// Test with slug, distance
	// Test with slug, distance, limits
	// Test with slug and year
	// Test with slug, year and limits
	// Test with slug, year, distance
	// Test with slug, year, distance, limits
	// Test invalid event
	// Test valid event invalid year
	// Test valid event, valid year, invalid distance
	// Test valid key with restricted event
}

func TestDeleteResults(t *testing.T) {
	// DELETE, /results/delete
	// Test no key
	// Test expired key
	// Test invalid key
	// Test invalid authorization header
	// Test bad request
	// Test wrong content type
	// Test invalid host
	// Test a valid request with just a slug.
	// Test with slug, limits
	// Test with slug, distance
	// Test with slug, distance, limits
	// Test with slug and year
	// Test with slug, year and limits
	// Test with slug, year, distance
	// Test with slug, year, distance, limits
	// Test invalid event
	// Test valid event invalid year
	// Test valid event, valid year, invalid distance
	// Test valid key with restricted event
}
