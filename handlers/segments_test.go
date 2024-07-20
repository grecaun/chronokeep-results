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

func TestGetSegments(t *testing.T) {
	// POST, /segments
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	// Test no key
	t.Log("Testing no key given.")
	body, err := json.Marshal(types.GetSegmentsRequest{
		Slug: variables.events["event1"].Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/segments", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.GetSegments(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodPost, "/segments", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetSegments(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodPost, "/segments", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetSegments(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodPost, "/segments", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetSegments(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/segments", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetSegments(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test no slug
	t.Log("Testing no slug.")
	request = httptest.NewRequest(http.MethodPost, "/segments", strings.NewReader("{}"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetSegments(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/segments", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetSegments(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodPost, "/segments", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetSegments(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test a valid request with just a slug.
	t.Log("Testing slug only.")
	request = httptest.NewRequest(http.MethodPost, "/segments", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetSegments(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetSegmentsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(variables.segments["event1"]["2021"]), len(resp.Segments))
			for _, outer := range variables.segments["event1"]["2021"] {
				found := false
				for _, inner := range resp.Segments {
					if outer.Equals(inner) {
						found = true
					}
				}
				assert.True(t, found)
			}
		}
	}
	// Test with slug and year
	t.Log("Testing with slug and year.")
	year := variables.eventYears["event1"]["2020"].Year
	body, err = json.Marshal(types.GetSegmentsRequest{
		Slug: variables.events["event1"].Slug,
		Year: &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/segments", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetSegments(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetSegmentsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(variables.segments["event1"]["2020"]), len(resp.Segments))
			for _, outer := range variables.segments["event1"]["2020"] {
				found := false
				for _, inner := range resp.Segments {
					if outer.Equals(inner) {
						found = true
					}
				}
				assert.True(t, found)
			}
		}
	}
	// Test valid event invalid year
	year = "2000"
	body, err = json.Marshal(types.GetSegmentsRequest{
		Slug: "event2",
		Year: &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/segments", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetSegments(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid key with restricted event
	t.Log("Testing restricted event but unauthorized key.")
	body, err = json.Marshal(types.GetSegmentsRequest{
		Slug: variables.events["event2"].Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/segments", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetSegments(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}

func TestAddSegments(t *testing.T) {
	// POST, /segments/add
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
	segments := variables.segments["event1"]["2020"]
	t.Log("Testing no key given.")
	body, err := json.Marshal(types.AddSegmentsRequest{
		Slug:     variables.events["event2"].Slug,
		Year:     "2023",
		Segments: segments,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/segments/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.AddSegments(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodPost, "/segments/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSegments(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodPost, "/segments/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSegments(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test read key (vs delete/write)
	t.Log("Testing read key.")
	request = httptest.NewRequest(http.MethodPost, "/segments/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSegments(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test different account's key
	t.Log("Testing different account's key.")
	request = httptest.NewRequest(http.MethodPost, "/segments/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSegments(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodPost, "/segments/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSegments(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/segments/add", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSegments(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test without slug/other info
	t.Log("Testing no slug.")
	request = httptest.NewRequest(http.MethodPost, "/segments/add", strings.NewReader("{}"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSegments(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/segments/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSegments(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test a valid request with a write key
	t.Log("Testing valid write key.")
	request = httptest.NewRequest(http.MethodPost, "/segments/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSegments(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.AddSegmentsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(segments), len(resp.Segments))
		}
	}
	database.DeleteSegments(eventYear.Identifier)
	// Test a valid request with a delete key
	t.Log("Testing valid write key.")
	body, err = json.Marshal(types.AddSegmentsRequest{
		Slug:     variables.events["event2"].Slug,
		Year:     "2023",
		Segments: segments,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/segments/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSegments(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.AddSegmentsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(segments), len(resp.Segments))
		}
	}
	t.Log("Verifying information was updated.")
	uploaded, err := database.GetSegments(eventYear.Identifier)
	if assert.NoError(t, err) {
		found := 0
		for _, seg := range uploaded {
			for _, inner := range segments {
				if seg.Equals(inner) {
					found++
				}
			}
		}
		assert.Equal(t, len(segments), found)
	}
	// Test invalid event
	t.Log("Testing invalid event.")
	body, err = json.Marshal(types.AddSegmentsRequest{
		Slug:     "invalid event",
		Year:     "2023",
		Segments: segments,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/segments/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSegments(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid event invalid year
	t.Log("Testing valid event invalid year.")
	body, err = json.Marshal(types.AddSegmentsRequest{
		Slug:     variables.events["event2"].Slug,
		Year:     "invalid-year",
		Segments: segments,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/segments/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSegments(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodPost, "/segments/add", strings.NewReader(""))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSegments(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test validation - Location
	t.Log("Test validation check - Location.")
	database.DeleteSegments(eventYear.Identifier)
	segments[0].Location = ""
	body, err = json.Marshal(types.AddSegmentsRequest{
		Slug:     variables.events["event2"].Slug,
		Year:     "2023",
		Segments: segments,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/segments/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSegments(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.AddSegmentsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(segments)-1, len(resp.Segments))
		}
	}
	// Test validation - DistanceName
	t.Log("Test validation check - DistanceName.")
	database.DeleteSegments(eventYear.Identifier)
	segments[0].Location = "Half Marathon"
	segments[0].DistanceName = ""
	body, err = json.Marshal(types.AddSegmentsRequest{
		Slug:     variables.events["event2"].Slug,
		Year:     "2023",
		Segments: segments,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/segments/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSegments(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.AddSegmentsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(segments)-1, len(resp.Segments))
		}
	}
	// Test validation - Name
	t.Log("Test validation check - Name.")
	database.DeleteSegments(eventYear.Identifier)
	segments[0].DistanceName = "Marathon"
	segments[0].Name = ""
	body, err = json.Marshal(types.AddSegmentsRequest{
		Slug:     variables.events["event2"].Slug,
		Year:     "2023",
		Segments: segments,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/segments/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSegments(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.AddSegmentsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(segments)-1, len(resp.Segments))
		}
	}
	// Test validation - DistanceValue
	t.Log("Test validation check - DistanceValue.")
	database.DeleteSegments(eventYear.Identifier)
	segments[0].Name = "13.1 Miles"
	segments[0].DistanceValue = 0.0
	body, err = json.Marshal(types.AddSegmentsRequest{
		Slug:     variables.events["event2"].Slug,
		Year:     "2023",
		Segments: segments,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/segments/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSegments(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.AddSegmentsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(segments)-1, len(resp.Segments))
		}
	}
	// Test validation - DistanceUnit
	t.Log("Test validation check - DistanceUnit.")
	database.DeleteSegments(eventYear.Identifier)
	segments[0].DistanceValue = 13.1
	segments[0].DistanceUnit = ""
	body, err = json.Marshal(types.AddSegmentsRequest{
		Slug:     variables.events["event2"].Slug,
		Year:     "2023",
		Segments: segments,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/segments/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSegments(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.AddSegmentsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(segments)-1, len(resp.Segments))
		}
	}
}

func TestDeleteSegments(t *testing.T) {
	// DELETE, /segments/delete
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	// Test no key
	t.Log("Testing no key given.")
	body, err := json.Marshal(types.DeleteSegmentsRequest{
		Slug: variables.events["event2"].Slug,
		Year: variables.eventYears["event2"]["2020"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request := httptest.NewRequest(http.MethodDelete, "/segments/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.DeleteSegments(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodDelete, "/segments/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteSegments(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodDelete, "/segments/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteSegments(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test read key
	t.Log("Testing read key.")
	request = httptest.NewRequest(http.MethodDelete, "/segments/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteSegments(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test different account's key
	t.Log("Testing different account's key.")
	request = httptest.NewRequest(http.MethodDelete, "/segments/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteSegments(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodDelete, "/segments/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteSegments(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodDelete, "/segments/delete", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteSegments(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test without slug/other info
	t.Log("Testing no slug.")
	request = httptest.NewRequest(http.MethodDelete, "/segments/delete", strings.NewReader("{}"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteSegments(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodDelete, "/segments/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteSegments(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test write key
	// verify information is there to delete
	t.Log("Verifying information before deletion.")
	deleted, err := database.GetSegments(variables.eventYears["event2"]["2020"].Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, len(variables.segments["event2"]["2020"]), len(deleted))
	}
	t.Log("Testing write key.")
	request = httptest.NewRequest(http.MethodDelete, "/segments/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteSegments(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.DeleteSegmentsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(deleted), int(resp.Count))
		}
	}
	// verify delete remotes information
	t.Log("Verifying information was deleted.")
	deleted, err = database.GetSegments(variables.eventYears["event2"]["2020"].Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(deleted))
	}
	// Test a delete request
	// verify information is there to delete
	t.Log("Verifying information before deletion.")
	database.AddSegments(variables.eventYears["event2"]["2020"].Identifier, variables.segments["event2"]["2020"])
	deleted, err = database.GetSegments(variables.eventYears["event2"]["2020"].Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, len(variables.segments["event2"]["2020"]), len(deleted))
	}
	body, err = json.Marshal(types.DeleteSegmentsRequest{
		Slug: variables.events["event2"].Slug,
		Year: variables.eventYears["event2"]["2020"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	t.Log("Testing delete key.")
	request = httptest.NewRequest(http.MethodDelete, "/segments/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteSegments(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.DeleteSegmentsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, len(deleted), int(resp.Count))
		}
	}
	// verify delete remotes information
	t.Log("Verifying information was deleted.")
	deleted, err = database.GetSegments(variables.eventYears["event2"]["2020"].Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(deleted))
	}
	// Test invalid event
	t.Log("Testing invalid event.")
	body, err = json.Marshal(types.DeleteSegmentsRequest{
		Slug: "invalid-event",
		Year: variables.eventYears["event2"]["2019"].Year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/segments/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteSegments(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid event invalid year
	t.Log("Testing valid event invalid year.")
	body, err = json.Marshal(types.DeleteSegmentsRequest{
		Slug: variables.events["event2"].Slug,
		Year: "invalid-year",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/segments/delete", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteSegments(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodDelete, "/segments/delete", strings.NewReader(""))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.DeleteSegments(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}
