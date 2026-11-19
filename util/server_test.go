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
	"bufio"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

var testTimeouts = ServerTimeouts{
	ReadHeader: 200 * time.Millisecond,
	Read:       500 * time.Millisecond,
	Write:      500 * time.Millisecond,
	Idle:       500 * time.Millisecond,
}

// startServer runs NewServer on a local port and returns its address.
func startServer(t *testing.T, handler http.Handler) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Error listening: %v", err)
	}
	server := NewServer(listener.Addr().String(), handler, testTimeouts)
	go server.Serve(listener)
	t.Cleanup(func() { server.Close() })
	return listener.Addr().String()
}

// closedWithin reports whether the server closes conn within limit.
func closedWithin(t *testing.T, conn net.Conn, limit time.Duration) bool {
	t.Helper()
	conn.SetReadDeadline(time.Now().Add(limit))
	_, err := io.Copy(io.Discard, conn)
	if ne, ok := err.(net.Error); ok && ne.Timeout() {
		return false
	}
	return true
}

func TestNewServerClosesSilentConnections(t *testing.T) {
	addr := startServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Error connecting: %v", err)
	}
	defer conn.Close()
	// send part of a request and then nothing
	conn.Write([]byte("GET / HTTP/1.1\r\nHost: x\r\n"))
	if !closedWithin(t, conn, 3*time.Second) {
		t.Error("Expected the server to close a connection that never finished its headers.")
	}
}

func TestNewServerStopsSlowResponses(t *testing.T) {
	addr := startServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.Write([]byte("too late"))
	}))
	client := &http.Client{Timeout: 5 * time.Second}
	start := time.Now()
	response, err := client.Get("http://" + addr + "/")
	if err == nil {
		body, _ := io.ReadAll(response.Body)
		response.Body.Close()
		if string(body) == "too late" {
			t.Error("Expected the response to be cut off by the write timeout.")
		}
	}
	if time.Since(start) > 4*time.Second {
		t.Errorf("Expected the slow response to end by the write timeout, took %v.", time.Since(start))
	}
}

func TestNewServerClosesIdleConnections(t *testing.T) {
	addr := startServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Error connecting: %v", err)
	}
	defer conn.Close()
	conn.Write([]byte("GET / HTTP/1.1\r\nHost: x\r\n\r\n"))
	reader := bufio.NewReader(conn)
	response, err := http.ReadResponse(reader, nil)
	if err != nil {
		t.Fatalf("Error reading response: %v", err)
	}
	body, _ := io.ReadAll(response.Body)
	response.Body.Close()
	if strings.TrimSpace(string(body)) != "ok" {
		t.Fatalf("Unexpected response %q.", body)
	}
	// the connection is kept alive, then closed once it's been idle too long
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	_, err = reader.ReadByte()
	if ne, ok := err.(net.Error); ok && ne.Timeout() {
		t.Error("Expected the server to close an idle connection.")
	}
}
