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

func TestGetMultiResults(t *testing.T) {
	// POST, /results/multi
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	// Test no key
	t.Log("Testing no key given.")
	body, err := json.Marshal(types.GetMultiResultsRequest{
		Slug:  variables.events["event1"].Slug,
		Years: []string{},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/results/multi", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.GetMultiResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodPost, "/results/multi", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetMultiResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodPost, "/results/multi", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetMultiResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodPost, "/results/multi", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetMultiResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/results/multi", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetMultiResults(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test no slug
	t.Log("Testing no slug.")
	request = httptest.NewRequest(http.MethodPost, "/results/multi", strings.NewReader("{}"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetMultiResults(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/results/multi", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetMultiResults(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodPost, "/results/multi", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetMultiResults(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid request -- no years specified
	t.Log("Testing invalid request - no years specified.")
	request = httptest.NewRequest(http.MethodPost, "/results/multi", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetMultiResults(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test a valid request.
	t.Log("Testing valid request.")
	body, err = json.Marshal(types.GetMultiResultsRequest{
		Slug:  variables.events["event1"].Slug,
		Years: []string{"2021", "2020"},
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/results/multi", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetMultiResults(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetMultiResultsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			// number of distances
			assert.Equal(t, 2, len(resp.Results))
			// verify the actual results returned
			for _, res := range resp.Results["2021"] {
				for _, outer := range res {
					found := false
					for _, inner := range variables.results["event1"]["2021"] {
						if inner.Equals(&outer) {
							found = true
							break
						}
					}
					assert.True(t, found)
				}
			}
			for _, res := range resp.Results["2020"] {
				for _, outer := range res {
					found := false
					for _, inner := range variables.results["event1"]["2020"] {
						if inner.Equals(&outer) {
							found = true
							break
						}
					}
					assert.True(t, found)
				}
			}
			assert.Equal(t, variables.events["event1"].Slug, resp.Event.Slug)
			assert.Equal(t, variables.events["event1"].Website, resp.Event.Website)
			assert.Equal(t, variables.events["event1"].Image, resp.Event.Image)
			assert.Equal(t, variables.events["event1"].ContactEmail, resp.Event.ContactEmail)
			assert.Equal(t, variables.events["event1"].AccessRestricted, resp.Event.AccessRestricted)
			assert.Equal(t, variables.events["event1"].Type, resp.Event.Type)
			event := variables.events["event1"]
			assert.True(t, event.Equals(&resp.Event))
			assert.Equal(t, 2, len(resp.Results))
		}
	}
}
