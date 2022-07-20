package handlers

import "testing"

func TestGetEvent(t *testing.T) {
	// POST, /event
	// Test no key
	// Test expired key
	// Test invalid key
	// Test wrong account key
	// Test invalid host
	// Test invalid authorization header
	// Test bad request
	// Test empty request
	// Test wrong content type
	// Test valid request
	// Test invalid event
}

func TestGetEvents(t *testing.T) {
	// GET, /event/all
	// Test no key
	// Test expired key
	// Test invalid key
	// Test invalid host
	// Test invalid authorization header
	// Test bad request
	// Test empty request
	// Test wrong content type
	// Test valid request
	// Verify that restricted events only show up for correct keys
}

func TestGetMyEvents(t *testing.T) {
	// GET, /event/my
	// Test no key
	// Test expired key
	// Test invalid key
	// Test invalid host
	// Test invalid authorization header
	// Test bad request
	// Test empty request
	// Test wrong content type
	// Test valid request
}

func TestAddEvent(t *testing.T) {
	// POST, /event/add
	// Test no key
	// Test expired key
	// Test invalid key
	// Test read key
	// Test invalid host
	// Test invalid authorization header
	// Test bad request
	// Test empty request
	// Test wrong content type
	// Test valid request
	// Test validation errors
}

func TestUpdateEvent(t *testing.T) {
	// PUT, /event/update
	// Test no key
	// Test expired key
	// Test invalid key
	// Test read key
	// Test wrong account key
	// Test invalid host
	// Test invalid authorization header
	// Test ownership
	// Test bad request
	// Test empty request
	// Test wrong content type
	// Test valid request
	// Test validation errors
	// Test unknown event
}

func TestDeleteEvent(t *testing.T) {
	// DELETE, /event/delete
	// Test no key
	// Test expired key
	// Test invalid key
	// Test read key
	// Test write key
	// Test wrong account key
	// Test invalid host
	// Test invalid authorization header
	// Test ownership
	// Test bad request
	// Test empty request
	// Test wrong content type
	// Test valid request
	// Test unknown event
}
