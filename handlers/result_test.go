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
	// Test a basic request with just a slug.
	body, err := json.Marshal(types.GetResultsRequest{
		Slug: variables.events[0].Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/results", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.keys[1].Value)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	h := Handler{}
	err = h.GetResults(c)
	t.Logf("GetResults response: %v", err)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, response.Code)
		t.Logf("Response body: %v", response.Body.String())
	}
	// Test with slug, limits
	// Test with slug, distance
	// Test with slug, distance, limits
	// Test with slug and year
	// Test with slug, year and limits
	// Test with slug, year, distance
	// Test with slug, year, distance, limits
	// Test expired key
	// Test valid key with restricted event
	// Test no key
	// Test invalid key
	// Test invalid authorization header
	// Test bad request
	// Test wrong content type
	// Test invalid host
	// Test invalid event
	// Test valid event invalid year
	// Test valid event, valid year, invalid distance
	// Test invalid method (GET / DELETE)
}

func TestGetAllResutls(t *testing.T) {
	// POST, /results/all
}

func TestGetFinishResults(t *testing.T) {
	// POST, /results/finish
}

func TestGetBibResults(t *testing.T) {
	// POST, /results/bib
}

func TestAddResults(t *testing.T) {
	// POST, /results/add
}

func TestDeleteResults(t *testing.T) {
	// DELETE, /results/delete
}
