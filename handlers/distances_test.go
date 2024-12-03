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

func TestGetDistances(t *testing.T) {
	// POST, /distances
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	// Test no key
	t.Log("Testing no key given.")
	body, err := json.Marshal(types.GetDistancesRequest{
		Slug: variables.events["event1"].Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/distances", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.GetDistances(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodPost, "/distances", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetDistances(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodPost, "/distances", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetDistances(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodPost, "/distances", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetDistances(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/distances", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetDistances(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test no slug
	t.Log("Testing no slug.")
	request = httptest.NewRequest(http.MethodPost, "/distances", strings.NewReader("{}"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetDistances(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/distances", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetDistances(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodPost, "/distances", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetDistances(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test a valid request with just a slug, multi distance
	t.Log("Testing slug only - multi.")
	request = httptest.NewRequest(http.MethodPost, "/distances", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetDistances(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetDistancesResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(variables.distances["event1"]["2021"]), len(resp.Distances))
			for _, outer := range variables.distances["event1"]["2021"] {
				found := false
				for _, inner := range resp.Distances {
					if outer.Equals(inner) {
						found = true
					}
				}
				assert.True(t, found)
			}
		}
	}
	// Test a valid request with just a slug, single distance
	t.Log("Testing slug only - single.")
	body, err = json.Marshal(types.GetDistancesRequest{
		Slug:     variables.events["event1"].Slug,
		Distance: &variables.distances["event1"]["2021"][0].Name,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/distances", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetDistances(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetDistanceResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, variables.distances["event1"]["2021"][0].Name, resp.Distance.Name)
			assert.Equal(t, variables.distances["event1"]["2021"][0].Certification, resp.Distance.Certification)
		}
	}
	// Test with slug and year - multi
	t.Log("Testing with slug and year - multi.")
	year := variables.eventYears["event1"]["2020"].Year
	body, err = json.Marshal(types.GetDistancesRequest{
		Slug: variables.events["event1"].Slug,
		Year: &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/distances", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetDistances(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetDistancesResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(variables.distances["event1"]["2020"]), len(resp.Distances))
			for _, outer := range variables.distances["event1"]["2020"] {
				found := false
				for _, inner := range resp.Distances {
					if outer.Equals(inner) {
						found = true
					}
				}
				assert.True(t, found)
			}
		}
	}
	// Test with slug and year - single
	t.Log("Testing with slug and year - single.")
	body, err = json.Marshal(types.GetDistancesRequest{
		Slug:     variables.events["event1"].Slug,
		Year:     &year,
		Distance: &variables.distances["event1"]["2021"][1].Name,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/distances", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetDistances(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetDistanceResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, variables.distances["event1"]["2021"][1].Name, resp.Distance.Name)
			assert.Equal(t, variables.distances["event1"]["2021"][1].Certification, resp.Distance.Certification)
		}
	}
	// Test valid event invalid year
	year = "2000"
	body, err = json.Marshal(types.GetDistancesRequest{
		Slug: "event2",
		Year: &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/distances", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetDistances(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid key with restricted event
	t.Log("Testing restricted event but unauthorized key.")
	body, err = json.Marshal(types.GetDistancesRequest{
		Slug: variables.events["event2"].Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/distances", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetDistances(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}

func TestAddDistances(t *testing.T) {
	// POST, /distances/add
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
	distances := variables.distances["event1"]["2020"]
	t.Log("Testing no key given.")
	body, err := json.Marshal(types.AddDistancesRequest{
		Slug:      variables.events["event2"].Slug,
		Year:      "2023",
		Distances: distances,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/distances/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.AddDistances(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodPost, "/distances/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddDistances(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodPost, "/distances/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddDistances(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test read key (vs delete/write)
	t.Log("Testing read key.")
	request = httptest.NewRequest(http.MethodPost, "/distances/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddDistances(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test different account's key
	t.Log("Testing different account's key.")
	request = httptest.NewRequest(http.MethodPost, "/distances/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddDistances(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodPost, "/distances/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddDistances(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/distances/add", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddDistances(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test without slug/other info
	t.Log("Testing no slug.")
	request = httptest.NewRequest(http.MethodPost, "/distances/add", strings.NewReader("{}"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddDistances(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/distances/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddDistances(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test a valid request with a write key
	t.Log("Testing valid write key.")
	request = httptest.NewRequest(http.MethodPost, "/distances/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddDistances(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetDistancesResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(distances), len(resp.Distances))
			for _, outer := range distances {
				found := false
				for _, inner := range resp.Distances {
					if outer.Equals(inner) {
						found = true
					}
				}
				assert.True(t, found)
			}
		}
	}
	database.DeleteDistances(eventYear.Identifier)
	// Test a valid request with a delete key
	t.Log("Testing valid write key.")
	body, err = json.Marshal(types.AddDistancesRequest{
		Slug:      variables.events["event2"].Slug,
		Year:      "2023",
		Distances: distances,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/distances/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddDistances(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetDistancesResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(distances), len(resp.Distances))
			for _, outer := range distances {
				found := false
				for _, inner := range resp.Distances {
					if outer.Equals(inner) {
						found = true
					}
				}
				assert.True(t, found)
			}
		}
	}
	t.Log("Verifying information was updated.")
	uploaded, err := database.GetDistances(eventYear.Identifier)
	if assert.NoError(t, err) {
		found := 0
		for _, seg := range uploaded {
			for _, inner := range distances {
				if seg.Equals(inner) {
					found++
				}
			}
		}
		assert.Equal(t, len(distances), found)
	}
	// Test invalid event
	t.Log("Testing invalid event.")
	body, err = json.Marshal(types.AddDistancesRequest{
		Slug:      "invalid event",
		Year:      "2023",
		Distances: distances,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/distances/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddDistances(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid event invalid year
	t.Log("Testing valid event invalid year.")
	body, err = json.Marshal(types.AddDistancesRequest{
		Slug:      variables.events["event2"].Slug,
		Year:      "invalid-year",
		Distances: distances,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/distances/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddDistances(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodPost, "/distances/add", strings.NewReader(""))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddDistances(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test validation - Certification
	t.Log("Test validation check - Certification.")
	database.DeleteDistances(eventYear.Identifier)
	distances[0].Certification = ""
	body, err = json.Marshal(types.AddDistancesRequest{
		Slug:      variables.events["event2"].Slug,
		Year:      "2023",
		Distances: distances,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/distances/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddDistances(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetDistancesResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(distances)-1, len(resp.Distances))
		}
	}
	// Test validation - Name
	t.Log("Test validation check - Name.")
	database.DeleteDistances(eventYear.Identifier)
	distances[0].Certification = "USATF Certification #WA512NewCert"
	distances[0].Name = ""
	body, err = json.Marshal(types.AddDistancesRequest{
		Slug:      variables.events["event2"].Slug,
		Year:      "2023",
		Distances: distances,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/distances/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddDistances(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetDistancesResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(distances)-1, len(resp.Distances))
		}
	}
}

func TestDeleteDistances(t *testing.T) {
	// DELETE, /distances/delete
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	// Test no key
	t.Log("Testing no key given.")
	body, err := json.Marshal(types.DeleteDistancesRequest{
		Slug: variables.events["event2"].Slug,
		Year: variables.eventYears["event2"]["2020"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request := httptest.NewRequest(http.MethodDelete, "/distances/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.DeleteDistances(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodDelete, "/distances/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteDistances(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodDelete, "/distances/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteDistances(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test read key
	t.Log("Testing read key.")
	request = httptest.NewRequest(http.MethodDelete, "/distances/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteDistances(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test different account's key
	t.Log("Testing different account's key.")
	request = httptest.NewRequest(http.MethodDelete, "/distances/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteDistances(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodDelete, "/distances/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteDistances(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodDelete, "/distances/delete", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteDistances(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test without slug/other info
	t.Log("Testing no slug.")
	request = httptest.NewRequest(http.MethodDelete, "/distances/delete", strings.NewReader("{}"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteDistances(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodDelete, "/distances/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteDistances(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test write key
	// verify information is there to delete
	t.Log("Verifying information before deletion.")
	deleted, err := database.GetDistances(variables.eventYears["event2"]["2020"].Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, len(variables.distances["event2"]["2020"]), len(deleted))
	}
	t.Log("Testing write key.")
	request = httptest.NewRequest(http.MethodDelete, "/distances/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteDistances(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.DeleteDistancesResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(deleted), int(resp.Count))
		}
	}
	// verify delete remotes information
	t.Log("Verifying information was deleted.")
	deleted, err = database.GetDistances(variables.eventYears["event2"]["2020"].Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(deleted))
	}
	// Test a delete request
	// verify information is there to delete
	t.Log("Verifying information before deletion.")
	database.AddDistances(variables.eventYears["event2"]["2020"].Identifier, variables.distances["event2"]["2020"])
	deleted, err = database.GetDistances(variables.eventYears["event2"]["2020"].Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, len(variables.distances["event2"]["2020"]), len(deleted))
	}
	body, err = json.Marshal(types.DeleteDistancesRequest{
		Slug: variables.events["event2"].Slug,
		Year: variables.eventYears["event2"]["2020"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing delete key.")
	request = httptest.NewRequest(http.MethodDelete, "/distances/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteDistances(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.DeleteDistancesResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(deleted), int(resp.Count))
		}
	}
	// verify delete remotes information
	t.Log("Verifying information was deleted.")
	deleted, err = database.GetDistances(variables.eventYears["event2"]["2020"].Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(deleted))
	}
	// Test invalid event
	t.Log("Testing invalid event.")
	body, err = json.Marshal(types.DeleteDistancesRequest{
		Slug: "invalid-event",
		Year: variables.eventYears["event2"]["2019"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/distances/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteDistances(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid event invalid year
	t.Log("Testing valid event invalid year.")
	body, err = json.Marshal(types.DeleteDistancesRequest{
		Slug: variables.events["event2"].Slug,
		Year: "invalid-year",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/distances/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteDistances(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodDelete, "/distances/delete", strings.NewReader(""))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteDistances(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}
