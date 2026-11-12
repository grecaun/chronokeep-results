/*
Chronokeep Results - Race Results API
Copyright (C) 2026 James Sentinella

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package handlers

import (
	"bufio"
	"chronokeep/results/types"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

func TestResultsBrokerFoldsBurstsIntoOneNotification(t *testing.T) {
	old := resultsNotifyDelay
	resultsNotifyDelay = 50 * time.Millisecond
	defer func() { resultsNotifyDelay = old }()
	b := newResultsBroker()
	first, stopFirst := b.subscribe(1)
	second, stopSecond := b.subscribe(1)
	other, stopOther := b.subscribe(2)
	defer stopOther()
	assert.Equal(t, 2, b.watchers(1))

	for i := 0; i < 5; i++ {
		b.notify(1)
	}
	for name, ch := range map[string]<-chan struct{}{"first": first, "second": second} {
		select {
		case <-ch:
		case <-time.After(time.Second):
			t.Fatalf("%s watcher wasn't notified", name)
		}
		select {
		case <-ch:
			t.Fatalf("%s watcher got more than one notification for one burst", name)
		case <-time.After(150 * time.Millisecond):
		}
	}
	select {
	case <-other:
		t.Fatal("a watcher of another event year was notified")
	default:
	}

	stopFirst()
	stopSecond()
	assert.Equal(t, 0, b.watchers(1))
	// with nobody watching, notify does nothing and doesn't leave anything pending
	b.notify(1)
	b.mu.Lock()
	assert.False(t, b.pending[1])
	b.mu.Unlock()
}

// streamServer runs the API on a real test server so the stream is flushed as it would be in production.
func streamServer(t *testing.T) *httptest.Server {
	e := echo.New()
	h := Handler{}
	h.Setup()
	h.Bind(e.Group(""))
	server := httptest.NewServer(e)
	t.Cleanup(server.Close)
	return server
}

func streamURL(server *httptest.Server, values map[string]string) string {
	q := url.Values{}
	for k, v := range values {
		q.Set(k, v)
	}
	return server.URL + "/results/stream?" + q.Encode()
}

// openStream connects to the stream and returns a function that reads the next event's name and data.
func openStream(t *testing.T, ctx context.Context, address string) (*http.Response, func() (string, string)) {
	t.Helper()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, address, nil)
	assert.NoError(t, err)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("Error opening stream: %v", err)
	}
	reader := bufio.NewReader(response.Body)
	next := func() (string, string) {
		name, data := "", ""
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return name, data
			}
			line = strings.TrimRight(line, "\r\n")
			switch {
			case line == "" && name != "":
				return name, data
			case strings.HasPrefix(line, "event: "):
				name = strings.TrimPrefix(line, "event: ")
			case strings.HasPrefix(line, "data: "):
				data = strings.TrimPrefix(line, "data: ")
			}
		}
	}
	return response, next
}

func TestStreamResultsChecksAccess(t *testing.T) {
	variables, finalize := setupTests(t)
	defer finalize(t)
	server := streamServer(t)
	for _, tc := range []struct {
		name   string
		values map[string]string
		status int
	}{
		{"no key", map[string]string{"slug": "event1"}, http.StatusUnauthorized},
		{"unknown key", map[string]string{"slug": "event1", "key": "not-a-key"}, http.StatusUnauthorized},
		{"expired key", map[string]string{"slug": "event1", "key": variables.knownValues["expired"]}, http.StatusUnauthorized},
		{"no slug", map[string]string{"key": variables.knownValues["read"]}, http.StatusBadRequest},
		{"unknown event", map[string]string{"slug": "nope", "key": variables.knownValues["read"]}, http.StatusNotFound},
		{"unknown year", map[string]string{"slug": "event1", "year": "1999", "key": variables.knownValues["read"]}, http.StatusNotFound},
		// the write key belongs to a different account than event2's owner
		{"someone else's restricted event", map[string]string{"slug": "event2", "key": variables.knownValues["write"]}, http.StatusUnauthorized},
	} {
		t.Run(tc.name, func(t *testing.T) {
			response, err := http.Get(streamURL(server, tc.values))
			if err != nil {
				t.Fatalf("Error calling stream: %v", err)
			}
			defer response.Body.Close()
			assert.Equal(t, tc.status, response.StatusCode)
		})
	}
}

func TestStreamResultsSendsReadyThenResultChanges(t *testing.T) {
	variables, finalize := setupTests(t)
	defer finalize(t)
	old := resultsNotifyDelay
	resultsNotifyDelay = 50 * time.Millisecond
	defer func() { resultsNotifyDelay = old }()
	server := streamServer(t)
	eventYear := variables.eventYears["event2"]["2021"]

	ctx, cancel := context.WithCancel(context.Background())
	// the owner can watch their own restricted event, and the key works as a header too
	request, _ := http.NewRequestWithContext(ctx, http.MethodGet, streamURL(server, map[string]string{"slug": "event2", "year": "2021"}), nil)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write2"])
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("Error opening stream: %v", err)
	}
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "text/event-stream", response.Header.Get("Content-Type"))
	response.Body.Close()
	cancel()

	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	response, next := openStream(t, ctx, streamURL(server, map[string]string{"slug": "event2", "year": "2021", "key": variables.knownValues["write2"]}))
	defer response.Body.Close()
	assert.Equal(t, http.StatusOK, response.StatusCode)
	name, data := next()
	assert.Equal(t, "ready", name)
	var ready struct {
		Year string `json:"year"`
		Live bool   `json:"live"`
	}
	assert.NoError(t, json.Unmarshal([]byte(data), &ready))
	assert.Equal(t, "2021", ready.Year)
	assert.False(t, ready.Live)

	// adding results notifies the stream and makes the event year live
	body, _ := json.Marshal(types.AddResultsRequest{
		Slug:    "event2",
		Year:    "2021",
		Results: variables.results["event2"]["2021"][:1],
	})
	addRequest := httptest.NewRequest(http.MethodPost, "/results/add", strings.NewReader(string(body)))
	addRequest.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	addRequest.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write2"])
	addResponse := httptest.NewRecorder()
	h := Handler{}
	h.Setup()
	if assert.NoError(t, h.AddResults(echo.New().NewContext(addRequest, addResponse))) {
		assert.Equal(t, http.StatusOK, addResponse.Code)
	}
	events := make(chan string, 1)
	go func() {
		name, _ := next()
		events <- name
	}()
	select {
	case name := <-events:
		assert.Equal(t, "results", name)
	case <-time.After(5 * time.Second):
		t.Fatal("stream didn't send a results event after results were added")
	}
	year, err := database.GetEventYear("event2", "2021")
	if assert.NoError(t, err) && assert.NotNil(t, year) {
		assert.Equal(t, eventYear.Identifier, year.Identifier)
		assert.True(t, year.Live)
	}

	// closing the stream stops watching
	cancel()
	deadline := time.Now().Add(5 * time.Second)
	for broker.watchers(eventYear.Identifier) > 0 && time.Now().Before(deadline) {
		time.Sleep(20 * time.Millisecond)
	}
	assert.Equal(t, 0, broker.watchers(eventYear.Identifier))
}

func TestStreamResultsWithoutYearWatchesTheLatestYear(t *testing.T) {
	variables, finalize := setupTests(t)
	defer finalize(t)
	server := streamServer(t)
	latest, err := database.GetEventAndYear("event1", "")
	if err != nil || latest == nil || latest.EventYear == nil {
		t.Fatalf("Error getting latest year: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	response, next := openStream(t, ctx, streamURL(server, map[string]string{"slug": "event1", "key": variables.knownValues["read"]}))
	defer response.Body.Close()
	name, data := next()
	assert.Equal(t, "ready", name)
	assert.Contains(t, data, `"year":"`+latest.EventYear.Year+`"`)
}

func TestAddResultsWithNoResultsDoesNotMakeEventLive(t *testing.T) {
	variables, finalize := setupTests(t)
	defer finalize(t)
	body, _ := json.Marshal(types.AddResultsRequest{Slug: "event2", Year: "2021", Results: []types.Result{}})
	request := httptest.NewRequest(http.MethodPost, "/results/add", strings.NewReader(string(body)))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+variables.knownValues["write2"])
	response := httptest.NewRecorder()
	h := Handler{}
	h.Setup()
	if assert.NoError(t, h.AddResults(echo.New().NewContext(request, response))) {
		assert.Equal(t, http.StatusOK, response.Code)
	}
	year, err := database.GetEventYear("event2", "2021")
	if assert.NoError(t, err) && assert.NotNil(t, year) {
		assert.False(t, year.Live)
	}
}
