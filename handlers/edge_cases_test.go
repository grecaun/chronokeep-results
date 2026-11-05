package handlers

import (
	"chronokeep/results/types"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/twilio/twilio-go/client"
)

func intPtr(i int) *int       { return &i }
func strPtr(s string) *string { return &s }
func int64Ptr(i int64) *int64 { return &i }

func TestGetLimitAndPage(t *testing.T) {
	assert.Equal(t, 0, getLimit(nil))
	assert.Equal(t, 0, getLimit(intPtr(-5)))
	assert.Equal(t, 0, getLimit(intPtr(0)))
	assert.Equal(t, 25, getLimit(intPtr(25)))
	assert.Equal(t, 0, getPage(nil))
	assert.Equal(t, 0, getPage(intPtr(-3)), "negative pages must not produce a negative OFFSET")
	assert.Equal(t, 0, getPage(intPtr(0)))
	assert.Equal(t, 0, getPage(intPtr(1)))
	assert.Equal(t, 4, getPage(intPtr(5)))
}

func TestMostRecentYear(t *testing.T) {
	now := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	y := func(year string, dt time.Time) types.EventYear { return types.EventYear{Year: year, DateTime: dt} }
	assert.Nil(t, mostRecentYear(nil, now))
	assert.Nil(t, mostRecentYear([]types.EventYear{}, now))
	// A future year listed first used to be returned even though past years exist.
	years := []types.EventYear{
		y("2025", now.AddDate(1, 0, 0)),
		y("2022", now.AddDate(-2, 0, 0)),
		y("2023", now.AddDate(-1, 0, 0)),
	}
	assert.Equal(t, "2023", mostRecentYear(years, now).Year)
	years = []types.EventYear{
		y("2023", now.AddDate(-1, 0, 0)),
		y("2022", now.AddDate(-2, 0, 0)),
	}
	assert.Equal(t, "2023", mostRecentYear(years, now).Year)
	// Only future years: the next upcoming one is returned.
	years = []types.EventYear{
		y("2027", now.AddDate(3, 0, 0)),
		y("2025", now.AddDate(1, 0, 0)),
	}
	assert.Equal(t, "2025", mostRecentYear(years, now).Year)
	// Returned pointer refers to the slice element, not a loop copy.
	years = []types.EventYear{y("2023", now.AddDate(-1, 0, 0))}
	assert.Same(t, &years[0], mostRecentYear(years, now))
}

func TestNormalizePhone(t *testing.T) {
	valid := map[string]string{
		"5551234567":        "5551234567",
		"555-123-4567":      "5551234567",
		"(555) 123-4567":    "5551234567",
		"555.123.4567":      "5551234567",
		"+1 (555) 123-4567": "5551234567",
		"15551234567":       "5551234567",
	}
	for in, want := range valid {
		got, ok := normalizePhone(in)
		assert.True(t, ok, in)
		assert.Equal(t, want, got, in)
	}
	for _, in := range []string{"", "555123456", "555123456789", "25551234567", "abcdefghij", "555-123-456a", "555_123_4567"} {
		_, ok := normalizePhone(in)
		assert.False(t, ok, in)
	}
}

func TestMessageWords(t *testing.T) {
	words := messageWords("please STOP, now!")
	assert.True(t, words["please"])
	assert.True(t, words["STOP"])
	assert.True(t, words["now"])
	assert.False(t, words["stop"], "callers are expected to lower case first")
	words = messageWords("send me my results")
	assert.False(t, words["end"])
	assert.Empty(t, messageWords(""))
}

func jsonRequest(t *testing.T, method, path string, body any, auth string) (*http.Request, *httptest.ResponseRecorder) {
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Error encoding request body into json object: %v", err)
	}
	request := httptest.NewRequest(method, path, strings.NewReader(string(b)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if auth != "" {
		request.Header.Set(echo.HeaderAuthorization, "Bearer "+auth)
	}
	return request, httptest.NewRecorder()
}

func loginToken(t *testing.T, account types.Account) string {
	token, refresh, err := createTokens(account.Email)
	if err != nil {
		t.Fatalf("Error creating test tokens: %v", err)
	}
	account.Token = *token
	account.RefreshToken = *refresh
	if err := database.UpdateTokens(account); err != nil {
		t.Fatalf("Error updating test tokens: %v", err)
	}
	return *token
}

func twilioSignature(token, webhook string, params url.Values) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sb strings.Builder
	sb.WriteString(webhook)
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteString(params.Get(k))
	}
	mac := hmac.New(sha1.New, []byte(token))
	mac.Write([]byte(sb.String()))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func TestTwilioEdgeCases(t *testing.T) {
	_, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	const authToken = "test-twilio-token"
	const webhook = "https://example.com/twilio"
	config.TwilioAuthToken = authToken
	config.TwilioResponseWebhookURL = webhook
	twilioRequestValidator = client.NewRequestValidator(authToken)

	send := func(params url.Values, signature string) *httptest.ResponseRecorder {
		request := httptest.NewRequest(http.MethodPost, "/twilio", strings.NewReader(params.Encode()))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		if signature != "" {
			request.Header.Set("X-Twilio-Signature", signature)
		}
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)
		assert.NotPanics(t, func() {
			assert.NoError(t, h.Twilio(c))
		})
		return response
	}
	signed := func(from, body string) *httptest.ResponseRecorder {
		params := url.Values{"From": {from}, "Body": {body}}
		return send(params, twilioSignature(authToken, webhook, params))
	}
	isBlocked := func(phone string) bool {
		phones, err := database.GetBlockedPhones()
		assert.NoError(t, err)
		for _, p := range phones {
			if p == phone {
				return true
			}
		}
		return false
	}

	// Missing fields used to panic with an index out of range.
	t.Log("Testing malformed twilio requests.")
	assert.Equal(t, http.StatusBadRequest, send(url.Values{}, "").Code)
	assert.Equal(t, http.StatusBadRequest, send(url.Values{"From": {"+15550001111"}, "Body": {"hi"}}, "").Code)
	params := url.Values{"Body": {"hi"}}
	assert.Equal(t, http.StatusBadRequest, send(params, twilioSignature(authToken, webhook, params)).Code)
	// Bad signature.
	assert.Equal(t, http.StatusUnauthorized, send(url.Values{"From": {"+15550001111"}, "Body": {"stop"}}, "bad").Code)
	assert.False(t, isBlocked("+15550001111"))

	// A message without a body (e.g. MMS) is fine.
	assert.Equal(t, http.StatusOK, signed("+15550001111", "").Code)

	// Words that merely contain a keyword must not unsubscribe ("send" contains "end").
	t.Log("Testing keyword substrings.")
	for _, msg := range []string{"send me my results", "weekend race?", "is the race starting?", "Quitting time"} {
		resp := signed("+15550001111", msg)
		assert.Equal(t, http.StatusOK, resp.Code, msg)
		assert.Empty(t, resp.Body.String(), msg)
		assert.False(t, isBlocked("+15550001111"), msg)
	}

	t.Log("Testing exact stop keyword.")
	resp := signed("+15550001111", " STOP ")
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Empty(t, resp.Body.String(), "twilio handles the reply for exact keywords")
	assert.True(t, isBlocked("+15550001111"))
	// Blocked numbers get no help response.
	resp = signed("+15550001111", "help")
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Empty(t, resp.Body.String())

	t.Log("Testing start keyword in a sentence.")
	resp = signed("+15550001111", "Please start again")
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "re-subscribed")
	assert.False(t, isBlocked("+15550001111"))

	t.Log("Testing stop keyword in a sentence.")
	resp = signed("+15550002222", "please stop texting me")
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "unsubscribed")
	assert.True(t, isBlocked("+15550002222"))

	t.Log("Testing help.")
	resp = signed("+15550003333", "HELP!")
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Reply STOP")
	resp = signed("+15550003333", "helpful info")
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Empty(t, resp.Body.String())
}

func TestSmsSubscriptionPhoneFormats(t *testing.T) {
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	slug := variables.events["event3"].Slug
	eventYear := variables.eventYears["event3"][time.Now().Format("2006")]
	year := eventYear.Year
	auth := variables.knownValues["delete2"]
	subscribed := func(phone string) bool {
		subs, err := database.GetSubscribedPhones(eventYear.Identifier)
		assert.NoError(t, err)
		for _, s := range subs {
			if s.Phone == phone {
				return true
			}
		}
		return false
	}

	// Formatted numbers are stored normalized so they can later be removed.
	for _, phone := range []string{"(555) 123-9876", "+1 555.123.9876", "555 123 9876"} {
		request, response := jsonRequest(t, http.MethodPost, "/sms/add", types.AddSmsSubscriptionRequest{
			Slug: slug, Year: &year, Bib: strPtr("1"), Phone: phone,
		}, auth)
		if assert.NoError(t, h.AddSmsSubscription(e.NewContext(request, response))) {
			assert.Equal(t, http.StatusOK, response.Code, phone)
		}
		assert.True(t, subscribed("5551239876"), phone)
		assert.False(t, subscribed(phone), phone)
		request, response = jsonRequest(t, http.MethodPost, "/sms/remove", types.RemoveSmsSubscriptionRequest{
			Slug: slug, Year: &year, Phone: phone,
		}, auth)
		if assert.NoError(t, h.RemoveSmsSubscription(e.NewContext(request, response))) {
			assert.Equal(t, http.StatusOK, response.Code, phone)
		}
		assert.False(t, subscribed("5551239876"), phone)
	}

	for _, phone := range []string{"abcdefghij", "555-123-98a6", "555123", "25551239876"} {
		request, response := jsonRequest(t, http.MethodPost, "/sms/add", types.AddSmsSubscriptionRequest{
			Slug: slug, Year: &year, Bib: strPtr("1"), Phone: phone,
		}, auth)
		if assert.NoError(t, h.AddSmsSubscription(e.NewContext(request, response))) {
			assert.Equal(t, http.StatusBadRequest, response.Code, phone)
		}
		request, response = jsonRequest(t, http.MethodPost, "/sms/remove", types.RemoveSmsSubscriptionRequest{
			Slug: slug, Year: &year, Phone: phone,
		}, auth)
		if assert.NoError(t, h.RemoveSmsSubscription(e.NewContext(request, response))) {
			assert.Equal(t, http.StatusBadRequest, response.Code, phone)
		}
	}

	// Bib set but empty, first/last missing.
	request, response := jsonRequest(t, http.MethodPost, "/sms/add", types.AddSmsSubscriptionRequest{
		Slug: slug, Year: &year, Bib: strPtr(""), Phone: "5551239876",
	}, auth)
	if assert.NoError(t, h.AddSmsSubscription(e.NewContext(request, response))) {
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
}

func TestNegativePagination(t *testing.T) {
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	slug := variables.events["event1"].Slug
	year := "2021"
	// Negative pages are treated as the first page instead of being passed to the
	// database as a negative OFFSET (an error in postgres and mysql).
	for _, page := range []int{-10, -1, 0, 1} {
		req := types.GetResultsRequest{Slug: slug, Year: &year, Limit: intPtr(10), Page: intPtr(page)}
		for name, fn := range map[string]func(*echo.Context) error{
			"results":     h.GetResults,
			"all results": h.GetAllResults,
		} {
			request, response := jsonRequest(t, http.MethodPost, "/results", req, variables.knownValues["write"])
			if assert.NoError(t, fn(e.NewContext(request, response))) && assert.Equal(t, http.StatusOK, response.Code, name) {
				var resp types.GetResultsResponse
				assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp))
				assert.Equal(t, 10, resp.Count, "%s page %d", name, page)
			}
		}
		request, response := jsonRequest(t, http.MethodPost, "/participants", types.GetParticipantsRequest{
			Slug: slug, Year: &year, Limit: intPtr(10), Page: intPtr(page),
		}, variables.knownValues["write"])
		if assert.NoError(t, h.GetParticipants(e.NewContext(request, response))) && assert.Equal(t, http.StatusOK, response.Code) {
			var resp types.GetParticipantsResponse
			assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp))
			assert.Equal(t, 10, len(resp.Participants), "page %d", page)
		}
	}
	// Negative limit means no limit.
	req := types.GetResultsRequest{Slug: slug, Year: &year, Limit: intPtr(-1)}
	request, response := jsonRequest(t, http.MethodPost, "/results/all", req, variables.knownValues["write"])
	if assert.NoError(t, h.GetAllResults(e.NewContext(request, response))) && assert.Equal(t, http.StatusOK, response.Code) {
		var resp types.GetResultsResponse
		assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp))
		all, err := database.GetAllDistanceResults(variables.eventYears["event1"]["2021"].Identifier, "", 0, 0)
		assert.NoError(t, err)
		assert.Greater(t, len(all), 10)
		assert.Equal(t, len(all), resp.Count)
	}
	// Restricted endpoint.
	token := loginToken(t, variables.accounts[0])
	request, response = jsonRequest(t, http.MethodPost, "/r/participants", types.GetParticipantsRequest{
		Slug: slug, Year: &year, Limit: intPtr(5), Page: intPtr(-2),
	}, token)
	if assert.NoError(t, h.RGetParticipants(e.NewContext(request, response))) && assert.Equal(t, http.StatusOK, response.Code) {
		var resp types.GetParticipantsResponse
		assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp))
		assert.Equal(t, 5, len(resp.Participants))
	}
}

func TestAddKeyDefaultTypeBecomesRead(t *testing.T) {
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	token := loginToken(t, variables.accounts[0])
	request, response := jsonRequest(t, http.MethodPost, "/key/add", types.AddKeyRequest{
		Key: types.RequestKey{Name: "default key", Type: "default", AllowedHosts: "a.com, b.com"},
	}, token)
	if assert.NoError(t, h.AddKey(e.NewContext(request, response))) && assert.Equal(t, http.StatusOK, response.Code) {
		var resp types.ModifyKeyResponse
		assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp))
		// Previously stored as "default", which is not treated as read only by write endpoints.
		assert.Equal(t, "read", resp.Key.Type)
		stored, err := database.GetKey(resp.Key.Value)
		if assert.NoError(t, err) && assert.NotNil(t, stored) {
			assert.Equal(t, "read", stored.Type)
		}
		// The key must not be able to write.
		request, response = jsonRequest(t, http.MethodPost, "/event/add", types.AddEventRequest{
			Event: types.Event{Name: "Default Key Event", Slug: "default-key-event", ContactEmail: "a@b.com", Type: "distance"},
		}, resp.Key.Value)
		request.Header.Set("Referer", "https://b.com/page")
		if assert.NoError(t, h.AddEvent(e.NewContext(request, response))) {
			assert.Equal(t, http.StatusUnauthorized, response.Code)
		}
	}
}

func TestGetEventPicksMostRecentPastYear(t *testing.T) {
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	_, err := database.AddEventYear(types.EventYear{
		EventIdentifier: variables.events["event1"].Identifier,
		Year:            "2099",
		DateTime:        time.Date(2099, 1, 1, 9, 0, 0, 0, time.Local),
		DaysAllowed:     1,
		RankingType:     "chip",
	})
	if err != nil {
		t.Fatalf("Error adding future year: %v", err)
	}
	request, response := jsonRequest(t, http.MethodPost, "/event", types.GetEventRequest{Slug: "event1"}, variables.knownValues["write"])
	if assert.NoError(t, h.GetEvent(e.NewContext(request, response))) && assert.Equal(t, http.StatusOK, response.Code) {
		var resp types.GetEventResponse
		assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp))
		assert.Equal(t, 3, len(resp.EventYears))
		if assert.NotNil(t, resp.Year) {
			assert.Equal(t, "2021", resp.Year.Year)
		}
		assert.NotEmpty(t, resp.Results)
	}
}

func TestRUpdateParticipantValidation(t *testing.T) {
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	token := loginToken(t, variables.accounts[0])
	valid := types.Participant{
		AlternateId: "5", Bib: "5", First: "Jane", Last: "Doe", Birthdate: "2/3/1990",
		Gender: "Woman", AgeGroup: "30-39", Distance: "Marathon",
	}
	for name, mod := range map[string]func(*types.Participant){
		"bad birthdate":    func(p *types.Participant) { p.Birthdate = "not a date" },
		"future birthdate": func(p *types.Participant) { p.Birthdate = time.Now().AddDate(1, 0, 0).Format("1/2/2006") },
		"no distance":      func(p *types.Participant) { p.Distance = "" },
	} {
		p := valid
		mod(&p)
		request, response := jsonRequest(t, http.MethodPost, "/r/participants/update", types.UpdateParticipantRequest{
			Slug: "event1", Year: "2021", Participant: p,
		}, token)
		if assert.NoError(t, h.RUpdateParticipant(e.NewContext(request, response)), name) {
			assert.Equal(t, http.StatusBadRequest, response.Code, name)
		}
	}
	request, response := jsonRequest(t, http.MethodPost, "/r/participants/update", types.UpdateParticipantRequest{
		Slug: "event1", Year: "2021", Participant: valid, UpdatedAfter: int64Ptr(0),
	}, token)
	if assert.NoError(t, h.RUpdateParticipant(e.NewContext(request, response))) && assert.Equal(t, http.StatusOK, response.Code) {
		var resp types.UpdateParticipantResponse
		assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp))
		assert.True(t, valid.TestEquals(&resp.Participant))
	}
	// Unknown participant.
	unknown := valid
	unknown.AlternateId = "does-not-exist"
	request, response = jsonRequest(t, http.MethodPost, "/r/participants/update", types.UpdateParticipantRequest{
		Slug: "event1", Year: "2021", Participant: unknown,
	}, token)
	if assert.NoError(t, h.RUpdateParticipant(e.NewContext(request, response))) {
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	}
}

func TestDeleteParticipantsCountsDeletedRows(t *testing.T) {
	variables, finalize := setupTests(t)
	defer finalize(t)
	e := echo.New()
	h := Handler{}
	h.Setup()
	// Only two of the four identifiers exist; the count used to be the number requested.
	request, response := jsonRequest(t, http.MethodDelete, "/participants/delete", types.DeleteParticipantsRequest{
		Slug: "event1", Year: "2021", Identifiers: []string{"0", "1", "missing", "also-missing"},
	}, variables.knownValues["write"])
	if assert.NoError(t, h.DeleteParticipants(e.NewContext(request, response))) && assert.Equal(t, http.StatusOK, response.Code) {
		var resp types.AddResultsResponse
		assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp))
		assert.Equal(t, 2, resp.Count)
	}
	// Deleting the same identifiers again removes nothing.
	request, response = jsonRequest(t, http.MethodDelete, "/participants/delete", types.DeleteParticipantsRequest{
		Slug: "event1", Year: "2021", Identifiers: []string{"0", "1"},
	}, variables.knownValues["write"])
	if assert.NoError(t, h.DeleteParticipants(e.NewContext(request, response))) && assert.Equal(t, http.StatusOK, response.Code) {
		var resp types.AddResultsResponse
		assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &resp))
		assert.Equal(t, 0, resp.Count)
	}
}

func TestBindRegistersAllRoutes(t *testing.T) {
	e := echo.New()
	h := Handler{}
	h.Bind(e.Group(""))
	h.BindRestricted(e.Group(""))
	paths := make(map[string]bool)
	for _, r := range e.Router().Routes() {
		paths[r.Method+" "+r.Path] = true
	}
	for _, want := range []string{
		"POST /r/participants/update",
		"POST /r/participants/update-many",
		"POST /twilio",
		"POST /sms/add",
		"DELETE /participants/delete",
	} {
		assert.True(t, paths[want], want)
	}
}
