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
	"context"
	"net"
	"time"
)

// TCPUserTimeout is how long data sent to a client may go unacknowledged before its
// connection is dropped. Without it, a client that vanished without closing its connection
// (a sleeping laptop, a phone that lost signal) keeps a results stream open until the
// system gives up retransmitting, which takes around 15 minutes on Linux. It's only set
// where the system supports it.
const TCPUserTimeout = 60 * time.Second

// Listen opens the API's TCP listener at addr. Connections it accepts use TCPUserTimeout.
func Listen(addr string) (net.Listener, error) {
	config := net.ListenConfig{Control: setUserTimeout}
	return config.Listen(context.Background(), "tcp", addr)
}
