package handlers

import (
	"chronokeep/results/types"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetSmsSubscriptions(t *testing.T) {
	// POST, /sms
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	// Test no key
	t.Log("Testing no key given.")
	body, err := json.Marshal(types.GetSmsSubscriptionsRequest{
		Slug: variables.events["event2"].Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/sms", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.GetSmsSubscriptions(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodPost, "/sms", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetSmsSubscriptions(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodPost, "/sms", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetSmsSubscriptions(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodPost, "/sms", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetSmsSubscriptions(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodPost, "/sms", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetSmsSubscriptions(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/sms", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetSmsSubscriptions(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test no slug
	t.Log("Testing no slug.")
	request = httptest.NewRequest(http.MethodPost, "/sms", strings.NewReader("{}"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetSmsSubscriptions(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/sms", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetSmsSubscriptions(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test a valid request with just a slug.
	t.Log("Testing slug only.")
	request = httptest.NewRequest(http.MethodPost, "/sms", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetSmsSubscriptions(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetSmsSubscriptionsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 2, len(resp.Subscriptions))
			// event2
			for _, outer := range resp.Subscriptions {
				assert.True(t, outer.Equals(&variables.sms[0]) || outer.Equals(&variables.sms[2]))
			}
		}
	}
	// Test with slug and year
	t.Log("Testing with slug and year.")
	year := variables.eventYears["event1"]["2020"].Year
	body, err = json.Marshal(types.GetSmsSubscriptionsRequest{
		Slug: variables.events["event1"].Slug,
		Year: &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/sms", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetSmsSubscriptions(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		var resp types.GetSmsSubscriptionsResponse
		if assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp)) {
			assert.Equal(t, 1, len(resp.Subscriptions))
			// event1 2020
			assert.Equal(t, variables.sms[1].Bib, resp.Subscriptions[0].Bib)
			assert.Equal(t, variables.sms[1].First, resp.Subscriptions[0].First)
			assert.Equal(t, variables.sms[1].Last, resp.Subscriptions[0].Last)
			assert.Equal(t, variables.sms[1].Phone, resp.Subscriptions[0].Phone)
		}
	}
	// Test invalid event
	t.Log("Testing event not found.")
	body, err = json.Marshal(types.GetSmsSubscriptionsRequest{
		Slug: "invalid event",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/sms", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetSmsSubscriptions(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid event invalid year
	year = "2000"
	body, err = json.Marshal(types.GetSmsSubscriptionsRequest{
		Slug: "event1",
		Year: &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/sms", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetSmsSubscriptions(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid key with restricted event
	t.Log("Testing restricted event but unauthorized key.")
	body, err = json.Marshal(types.GetSmsSubscriptionsRequest{
		Slug: variables.events["event2"].Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/sms", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.GetSmsSubscriptions(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}

func TestAddSmsSubscription(t *testing.T) {
	// POST, /sms/add
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	thisYear := fmt.Sprintf("%d", time.Now().Year())
	// Test no key
	t.Log("Testing no key given.")
	bib := "500"
	body, err := json.Marshal(types.AddSmsSubscriptionRequest{
		Slug:  variables.events["event3"].Slug,
		Bib:   &bib,
		Phone: "5551234567",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/sms/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.AddSmsSubscription(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test expired key
	t.Log("Testing expired key.")
	request = httptest.NewRequest(http.MethodPost, "/sms/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["expired2"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSmsSubscription(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid key
	t.Log("Testing invalid key.")
	request = httptest.NewRequest(http.MethodPost, "/sms/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer not-a-valid-key")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSmsSubscription(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid host
	t.Log("Testing invalid host.")
	request = httptest.NewRequest(http.MethodPost, "/sms/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["delete"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSmsSubscription(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test invalid authorization header
	t.Log("Testing invalid authorization header.")
	request = httptest.NewRequest(http.MethodPost, "/sms/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "not-a-valid-auth-header")
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSmsSubscription(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request = httptest.NewRequest(http.MethodPost, "/sms/add", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSmsSubscription(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test no slug
	t.Log("Testing no slug.")
	request = httptest.NewRequest(http.MethodPost, "/sms/add", strings.NewReader("{}"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSmsSubscription(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/sms/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSmsSubscription(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test a valid request with just a slug.
	t.Log("Testing slug only.")
	request = httptest.NewRequest(http.MethodPost, "/sms/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSmsSubscription(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		subs, err := database.GetSubscribedPhones(variables.eventYears["event3"][thisYear].Identifier)
		if assert.NoError(t, err) {
			assert.Equal(t, 1, len(subs))
			// event2
			found := false
			for _, outer := range subs {
				if outer.Bib == bib && outer.Phone == "5551234567" {
					found = true
				}
			}
			assert.True(t, found)
		}
	}
	t.Log("Testing repeat.")
	request = httptest.NewRequest(http.MethodPost, "/sms/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSmsSubscription(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		subs, err := database.GetSubscribedPhones(variables.eventYears["event3"][thisYear].Identifier)
		if assert.NoError(t, err) {
			assert.Equal(t, 1, len(subs))
			// event2
			found := false
			for _, outer := range subs {
				if outer.Bib == bib && outer.Phone == "5551234567" {
					found = true
				}
			}
			assert.True(t, found)
		}
	}
	// Test with slug and year
	t.Log("Testing with slug and year.")
	first := "John"
	last := "Smith"
	body, err = json.Marshal(types.AddSmsSubscriptionRequest{
		Slug:  variables.events["event3"].Slug,
		Year:  &thisYear,
		First: &first,
		Last:  &last,
		Phone: "5558765432",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/sms/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSmsSubscription(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		subs, err := database.GetSubscribedPhones(variables.eventYears["event3"][thisYear].Identifier)
		if assert.NoError(t, err) {
			assert.Equal(t, 2, len(subs))
			// event1 2020
			found := false
			for _, outer := range subs {
				if outer.First == first && outer.Last == last && outer.Phone == "5558765432" {
					found = true
				}
			}
			assert.True(t, found)
		}
	}
	// Test body with only a single name and no bib
	t.Log("Testing body with only a single name and no bib.")
	body, err = json.Marshal(types.AddSmsSubscriptionRequest{
		Slug:  variables.events["event3"].Slug,
		Year:  &thisYear,
		Last:  &last,
		Phone: "5558765432",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/sms/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSmsSubscription(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
		subs, err := database.GetSubscribedPhones(variables.eventYears["event3"][thisYear].Identifier)
		if assert.NoError(t, err) {
			assert.Equal(t, 2, len(subs))
		}
	}
	// Test body without phone
	t.Log("Testing body without phone.")
	body, err = json.Marshal(types.AddSmsSubscriptionRequest{
		Slug:  variables.events["event3"].Slug,
		Year:  &thisYear,
		First: &first,
		Last:  &last,
		Phone: "",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/sms/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSmsSubscription(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
		subs, err := database.GetSubscribedPhones(variables.eventYears["event3"][thisYear].Identifier)
		if assert.NoError(t, err) {
			assert.Equal(t, 2, len(subs))
		}
	}
	// Test invalid event
	t.Log("Testing event not found.")
	body, err = json.Marshal(types.AddSmsSubscriptionRequest{
		Slug: "invalid event",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/sms/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSmsSubscription(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid event invalid year
	t.Log("Testing valid event invalid year.")
	year := "2000"
	body, err = json.Marshal(types.AddSmsSubscriptionRequest{
		Slug: "event1",
		Year: &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/sms/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSmsSubscription(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test invalid bib/first/last
	t.Log("Testing invalid bib/first/last.")
	year = "2000"
	bib = ""
	first = ""
	last = ""
	body, err = json.Marshal(types.AddSmsSubscriptionRequest{
		Slug:  variables.events["event3"].Slug,
		Year:  &thisYear,
		Bib:   &bib,
		First: &first,
		Last:  &last,
		Phone: "1122334567",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/sms/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["read"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSmsSubscription(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test valid key with restricted event
	t.Log("Testing restricted event but unauthorized key.")
	body, err = json.Marshal(types.AddSmsSubscriptionRequest{
		Slug: variables.events["event3"].Slug,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/sms/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write"])
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.AddSmsSubscription(c)) {
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}

func TestRemoveSmsSubscription(t *testing.T) {
	// POST, /sms/remove
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	body, err := json.Marshal(types.RemoveSmsSubscriptionRequest{
		Slug:  variables.events["event2"].Slug,
		Phone: "1235557890",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	// Test bad request
	t.Log("Testing bad request.")
	request := httptest.NewRequest(http.MethodPost, "/sms/remove", strings.NewReader("////"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	c := e.NewContext(request, response)
	if assert.NoError(t, h.RemoveSmsSubscription(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test no slug
	t.Log("Testing no slug.")
	request = httptest.NewRequest(http.MethodPost, "/sms/remove", strings.NewReader("{}"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RemoveSmsSubscription(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test wrong content type
	t.Log("Testing wrong content type.")
	request = httptest.NewRequest(http.MethodPost, "/sms/remove", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RemoveSmsSubscription(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Test a valid request with just a slug.
	t.Log("Testing slug only.")
	subs, err := database.GetSubscribedPhones(variables.eventYears["event2"]["2021"].Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, 2, len(subs))
	}
	request = httptest.NewRequest(http.MethodPost, "/sms/remove", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RemoveSmsSubscription(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		subs, err := database.GetSubscribedPhones(variables.eventYears["event2"]["2021"].Identifier)
		if assert.NoError(t, err) {
			assert.Equal(t, 0, len(subs))
		}
	}
	// Test a valid request with just a slug.
	t.Log("Testing repeat valid request.")
	subs, err = database.GetSubscribedPhones(variables.eventYears["event2"]["2021"].Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(subs))
	}
	request = httptest.NewRequest(http.MethodPost, "/sms/remove", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RemoveSmsSubscription(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		subs, err := database.GetSubscribedPhones(variables.eventYears["event2"]["2021"].Identifier)
		if assert.NoError(t, err) {
			assert.Equal(t, 0, len(subs))
		}
	}
	// Test with slug and year
	t.Log("Testing with slug and year.")
	subs, err = database.GetSubscribedPhones(variables.eventYears["event2"]["2020"].Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, 3, len(subs))
	}
	year := variables.eventYears["event2"]["2020"].Year
	body, err = json.Marshal(types.RemoveSmsSubscriptionRequest{
		Slug:  variables.events["event2"].Slug,
		Year:  &year,
		Phone: "1235557890",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/sms/remove", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RemoveSmsSubscription(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		subs, err := database.GetSubscribedPhones(variables.eventYears["event2"]["2020"].Identifier)
		if assert.NoError(t, err) {
			assert.Equal(t, 1, len(subs))
		}
	}
	// Test body without phone
	t.Log("Testing body without phone.")
	year = variables.eventYears["event1"]["2020"].Year
	body, err = json.Marshal(types.RemoveSmsSubscriptionRequest{
		Slug:  variables.events["event1"].Slug,
		Year:  &year,
		Phone: "",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/sms/remove", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RemoveSmsSubscription(c)) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
		subs, err := database.GetSubscribedPhones(variables.eventYears["event1"]["2020"].Identifier)
		if assert.NoError(t, err) {
			assert.Equal(t, 1, len(subs))
		}
	}
	// Test invalid event
	t.Log("Testing event not found.")
	body, err = json.Marshal(types.RemoveSmsSubscriptionRequest{
		Slug: "invalid event",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/sms/remove", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RemoveSmsSubscription(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid event invalid year
	year = "2000"
	body, err = json.Marshal(types.RemoveSmsSubscriptionRequest{
		Slug: "event1",
		Year: &year,
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/sms/remove", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RemoveSmsSubscription(c)) {
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
	// Test valid key with restricted event
	t.Log("Testing restricted event but unauthorized key.")
	year = variables.eventYears["event2"]["2020"].Year
	body, err = json.Marshal(types.RemoveSmsSubscriptionRequest{
		Slug:  variables.events["event2"].Slug,
		Year:  &year,
		Phone: "1325557890",
	})
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request = httptest.NewRequest(http.MethodPost, "/sms/remove", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response = httptest.NewRecorder()
	c = e.NewContext(request, response)
	if assert.NoError(t, h.RemoveSmsSubscription(c)) {
		assert.Equal(t, http.StatusOK, response.Code)
		subs, err := database.GetSubscribedPhones(variables.eventYears["event2"]["2020"].Identifier)
		if assert.NoError(t, err) {
			assert.Equal(t, 0, len(subs))
		}
	}
}
