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
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v5"
)

var (
	// resultsNotifyDelay groups bursts of result changes into one notification per event year.
	resultsNotifyDelay = 2 * time.Second
	// streamPingInterval keeps idle streams from being closed by proxies and browsers.
	streamPingInterval = 30 * time.Second
	// broker tells open results streams when an event year's results change.
	broker = newResultsBroker()
)

// resultsBroker keeps the open results streams for each event year and tells them when
// that event year's results change. Everything is in memory, so it only reaches streams
// on this server.
type resultsBroker struct {
	mu          sync.Mutex
	subscribers map[int64]map[chan struct{}]struct{}
	pending     map[int64]bool
}

func newResultsBroker() *resultsBroker {
	return &resultsBroker{
		subscribers: make(map[int64]map[chan struct{}]struct{}),
		pending:     make(map[int64]bool),
	}
}

// subscribe opens a channel that receives a value when the event year's results change.
// Call the returned function to stop listening.
func (b *resultsBroker) subscribe(eventYearID int64) (<-chan struct{}, func()) {
	ch := make(chan struct{}, 1)
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.subscribers[eventYearID] == nil {
		b.subscribers[eventYearID] = make(map[chan struct{}]struct{})
	}
	b.subscribers[eventYearID][ch] = struct{}{}
	return ch, func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		delete(b.subscribers[eventYearID], ch)
		if len(b.subscribers[eventYearID]) == 0 {
			delete(b.subscribers, eventYearID)
		}
	}
}

// notify tells everyone watching an event year that its results changed. Calls that come
// in while a notification is waiting to go out are folded into it.
func (b *resultsBroker) notify(eventYearID int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.subscribers[eventYearID]) == 0 || b.pending[eventYearID] {
		return
	}
	b.pending[eventYearID] = true
	time.AfterFunc(resultsNotifyDelay, func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		delete(b.pending, eventYearID)
		for ch := range b.subscribers[eventYearID] {
			// a full channel already has a notification waiting, which covers this one
			select {
			case ch <- struct{}{}:
			default:
			}
		}
	})
}

// watchers returns how many streams are open for an event year.
func (b *resultsBroker) watchers(eventYearID int64) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.subscribers[eventYearID])
}

// StreamResults keeps a server-sent events stream open for an event year and sends a
// "results" event whenever its results change, so pages know to fetch them again.
// The browser's EventSource can't set headers, so the key may also be passed as the
// "key" query parameter.
func (h Handler) StreamResults(c *echo.Context) error {
	k, err := retrieveKey(c.Request())
	if err != nil || k == nil {
		key := c.QueryParam("key")
		if key == "" {
			return getAPIError(c, http.StatusUnauthorized, "Key Not Provided", err)
		}
		k = &key
	}
	slug := c.QueryParam("slug")
	if slug == "" {
		return getAPIError(c, http.StatusBadRequest, "Slug Not Provided", nil)
	}
	mkey, err := database.GetKeyAndAccount(*k)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Key/Account", err)
	}
	if mkey == nil || mkey.Key == nil || mkey.Account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key/Account Not Found", nil)
	}
	if mkey.Key.Expired() {
		return getAPIError(c, http.StatusUnauthorized, "Expired Key", nil)
	}
	if !mkey.Key.IsAllowed(c.Request().Referer()) {
		return getAPIError(c, http.StatusUnauthorized, "Host Not Allowed", nil)
	}
	mult, err := database.GetEventAndYear(slug, c.QueryParam("year"))
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event/Year", err)
	}
	if mult == nil || mult.Event == nil || mult.EventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Event/Year Not Found", nil)
	}
	if mult.Event.AccessRestricted && mkey.Account.Identifier != mult.Event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Restricted Event", nil)
	}

	updates, stop := broker.subscribe(mult.EventYear.Identifier)
	defer stop()

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)
	rc := http.NewResponseController(w)
	send := func(msg string) bool {
		if _, err := fmt.Fprint(w, msg); err != nil {
			return false
		}
		return rc.Flush() == nil
	}
	// the ready event says which year is being watched, since the year can be left out
	ready, _ := json.Marshal(map[string]any{"year": mult.EventYear.Year, "live": mult.EventYear.Live})
	if !send(fmt.Sprintf("retry: 5000\nevent: ready\ndata: %s\n\n", ready)) {
		return nil
	}
	ping := time.NewTicker(streamPingInterval)
	defer ping.Stop()
	done := c.Request().Context().Done()
	for {
		var msg string
		select {
		case <-done:
			return nil
		case <-updates:
			msg = "event: results\ndata: {}\n\n"
		case <-ping.C:
			msg = ": ping\n\n"
		}
		if !send(msg) {
			return nil
		}
	}
}
