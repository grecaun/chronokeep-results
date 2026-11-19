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

package util

import (
	"net/http"
	"time"
)

// ServerTimeouts limits how long the API waits on a connection, so slow or silent clients
// can't hold connections open forever.
type ServerTimeouts struct {
	// ReadHeader is how long a client has to send a request's headers.
	ReadHeader time.Duration
	// Read is how long a client has to send a whole request, including its body.
	Read time.Duration
	// Write is how long a response has to be written, counted from the end of the request's
	// headers. The results stream clears this for itself, since it stays open on purpose.
	Write time.Duration
	// Idle is how long a kept-alive connection may wait for its next request.
	Idle time.Duration
}

// DefaultServerTimeouts are the timeouts the API runs with. Read and Write leave room for
// large result uploads and responses on slow connections.
var DefaultServerTimeouts = ServerTimeouts{
	ReadHeader: 10 * time.Second,
	Read:       60 * time.Second,
	Write:      60 * time.Second,
	Idle:       120 * time.Second,
}

// NewServer makes the API's HTTP server with the given timeouts.
func NewServer(addr string, handler http.Handler, timeouts ServerTimeouts) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: timeouts.ReadHeader,
		ReadTimeout:       timeouts.Read,
		WriteTimeout:      timeouts.Write,
		IdleTimeout:       timeouts.Idle,
	}
}
