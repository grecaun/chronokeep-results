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
	"net"
	"syscall"
	"testing"

	"golang.org/x/sys/unix"
)

func TestListenSetsUserTimeoutOnAcceptedConnections(t *testing.T) {
	listener, err := Listen("127.0.0.1:0")
	if err != nil {
		t.Fatalf("Error listening: %v", err)
	}
	defer listener.Close()
	client, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatalf("Error connecting: %v", err)
	}
	defer client.Close()
	conn, err := listener.Accept()
	if err != nil {
		t.Fatalf("Error accepting: %v", err)
	}
	defer conn.Close()
	raw, err := conn.(syscall.Conn).SyscallConn()
	if err != nil {
		t.Fatalf("Error getting raw connection: %v", err)
	}
	var timeout int
	var sockErr error
	raw.Control(func(fd uintptr) {
		timeout, sockErr = unix.GetsockoptInt(int(fd), unix.IPPROTO_TCP, unix.TCP_USER_TIMEOUT)
	})
	if sockErr != nil {
		t.Fatalf("Error reading TCP_USER_TIMEOUT: %v", sockErr)
	}
	if timeout != int(TCPUserTimeout.Milliseconds()) {
		t.Errorf("Expected accepted connections to have a %v user timeout, found %dms.", TCPUserTimeout, timeout)
	}
}
