package handlers

import (
	"chronokeep/results/types"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/twilio/twilio-go/client"
)

func TestCanonicalPhone(t *testing.T) {
	valid := map[string]string{
		"5551234567":        "+15551234567",
		"(555) 123-4567":    "+15551234567",
		"555.123.4567":      "+15551234567",
		"15551234567":       "+15551234567",
		"+15551234567":      "+15551234567",
		"+1 (555) 123-4567": "+15551234567",
		"+447911123456":     "+447911123456",
		" +44 7911 123456 ": "+447911123456",
	}
	for in, want := range valid {
		got, ok := canonicalPhone(in)
		assert.True(t, ok, in)
		assert.Equal(t, want, got, in)
	}
	for _, in := range []string{"", "+", "1234", "555123456", "25551234567", "447911123456", "abcdefghij", "555-123-456a", "+1234567890123456"} {
		_, ok := canonicalPhone(in)
		assert.False(t, ok, in)
	}
	assert.ElementsMatch(t, []string{"+15551234567", "15551234567", "5551234567"}, phoneVariants("+15551234567"))
	assert.ElementsMatch(t, []string{"+447911123456", "447911123456"}, phoneVariants("+447911123456"))
	assert.True(t, samePhone("5551234567", "+15551234567"))
	assert.True(t, samePhone("(555) 123-4567", "15551234567"))
	assert.False(t, samePhone("5551234567", "5551234568"))
	assert.False(t, samePhone("", ""))
	// Substrings no longer count as the same number.
	assert.False(t, samePhone("1234567", "+15551234567"))
}

func TestTextUnsubscribeAcrossFormats(t *testing.T) {
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	const authToken = "test-twilio-token"
	const webhook = "https://example.com/twilio"
	config.TwilioAuthToken = authToken
	config.TwilioResponseWebhookURL = webhook
	twilioRequestValidator = client.NewRequestValidator(authToken)
	text := func(from, body string) *httptest.ResponseRecorder {
		params := url.Values{"From": {from}, "Body": {body}}
		request := httptest.NewRequest(http.MethodPost, "/twilio", strings.NewReader(params.Encode()))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		request.Header.Set("X-Twilio-Signature", twilioSignature(authToken, webhook, params))
		response := httptest.NewRecorder()
		assert.NoError(t, h.Twilio(e.NewContext(request, response)))
		return response
	}
	blocked := func() []string {
		phones, err := database.GetBlockedPhones()
		assert.NoError(t, err)
		return phones
	}

	// Texting STOP stores the number in E.164 format, which the desktop app matches against.
	t.Log("Testing STOP stores canonical number.")
	assert.Equal(t, http.StatusOK, text("+15551110000", "STOP").Code)
	assert.Equal(t, []string{"+15551110000"}, blocked())
	assert.Equal(t, http.StatusOK, text("+15551110000", "stop").Code)
	assert.Equal(t, []string{"+15551110000"}, blocked(), "repeat STOP must not duplicate")
	t.Log("Testing START removes it.")
	assert.Equal(t, http.StatusOK, text("+15551110000", "START").Code)
	assert.Empty(t, blocked())

	// Numbers blocked in older formats (10 or 11 digits without +) can be unblocked by text.
	t.Log("Testing START removes legacy formats.")
	assert.NoError(t, database.AddBlockedPhones([]string{"5552220000", "15552220000", "+15552220000"}))
	assert.Equal(t, http.StatusOK, text("+15552220000", "start").Code)
	assert.Empty(t, blocked())

	// A number blocked through the API in any format is recognized when it texts in.
	t.Log("Testing API blocked number recognized by twilio.")
	request, response := jsonRequest(t, http.MethodPost, "/blocked/phones/add", types.ModifyBannedPhoneRequest{Phone: "(555) 333-0000"}, "")
	if assert.NoError(t, h.AddBannedPhone(e.NewContext(request, response))) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	assert.Equal(t, []string{"+15553330000"}, blocked())
	resp := text("+15553330000", "help")
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Empty(t, resp.Body.String(), "blocked numbers get no reply")
	// Legacy 10 digit entries are recognized too.
	assert.NoError(t, database.AddBlockedPhone("5554440000"))
	resp = text("+15554440000", "help")
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Empty(t, resp.Body.String())
	// A number that merely contains a blocked entry as a substring is not blocked.
	assert.NoError(t, database.AddBlockedPhone("4440000"))
	resp = text("+15559994444", "help")
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Reply STOP")

	// Unblocking through the API removes all formats.
	t.Log("Testing API unblock removes all formats.")
	assert.NoError(t, database.AddBlockedPhones([]string{"5553330000", "15553330000"}))
	token := loginToken(t, variables.accounts[0])
	request, response = jsonRequest(t, http.MethodPost, "/blocked/phones/unblock", types.ModifyBannedPhoneRequest{Phone: "555-333-0000"}, token)
	if assert.NoError(t, h.RemoveBannedPhone(e.NewContext(request, response))) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	for _, p := range blocked() {
		assert.False(t, samePhone(p, "5553330000"), p)
	}

	// Invalid numbers are rejected.
	for _, phone := range []string{"invalid phone", "1234", "+1"} {
		request, response = jsonRequest(t, http.MethodPost, "/blocked/phones/add", types.ModifyBannedPhoneRequest{Phone: phone}, "")
		if assert.NoError(t, h.AddBannedPhone(e.NewContext(request, response))) {
			assert.Equal(t, http.StatusBadRequest, response.Code, phone)
		}
	}
}

func TestAddEventWithoutContactEmail(t *testing.T) {
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	request, response := jsonRequest(t, http.MethodPost, "/event/add", types.AddEventRequest{
		Event: types.Event{Name: "No Contact Event", Slug: "no-contact", Type: "distance"},
	}, variables.knownValues["write"])
	if assert.NoError(t, h.AddEvent(e.NewContext(request, response))) && assert.Equal(t, http.StatusOK, response.Code) {
		var resp types.ModifyEventResponse
		assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp))
		assert.Equal(t, "", resp.Event.ContactEmail)
		assert.Equal(t, "no-contact", resp.Event.Slug)
	}
	// An invalid email is still rejected.
	request, response = jsonRequest(t, http.MethodPost, "/event/add", types.AddEventRequest{
		Event: types.Event{Name: "Bad Contact Event", Slug: "bad-contact", Type: "distance", ContactEmail: "nope"},
	}, variables.knownValues["write"])
	if assert.NoError(t, h.AddEvent(e.NewContext(request, response))) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
	// Removing the contact email from an existing event works.
	token := loginToken(t, variables.accounts[0])
	request, response = jsonRequest(t, http.MethodPut, "/r/event/update", types.UpdateEventRequest{
		Event: types.Event{Name: "Event 1", Slug: "event1", Type: "distance"},
	}, token)
	if assert.NoError(t, h.RUpdateEvent(e.NewContext(request, response))) && assert.Equal(t, http.StatusOK, response.Code) {
		var resp types.ModifyEventResponse
		assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp))
		assert.Equal(t, "", resp.Event.ContactEmail)
	}
}
