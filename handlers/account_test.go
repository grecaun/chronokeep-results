package handlers

import "testing"

func TestLogin(t *testing.T) {
	// POST, /r/account/login (no auth header required)
}

func TestRefresh(t *testing.T) {
	// POST, /r/account/refresh (no auth header required)
}

func TestGetAccount(t *testing.T) {
	// POST, /r/account
}

func TestGetAccounts(t *testing.T) {
	// GET, /r/account/all
}

func TestLogout(t *testing.T) {
	// POST, /r/account/logout
}

func TestAddAccount(t *testing.T) {
	// POST, /r/account/add
}

func TestUpdateAccount(t *testing.T) {
	// PUT, /r/account/update
}

func TestChangePassword(t *testing.T) {
	// PUT, /r/account/password
}

func TestChangeEmail(t *testing.T) {
	// PUT, /r/account/email
}

func TestUnlock(t *testing.T) {
	// POST, /r/account/unlock
}

func TestDelete(t *testing.T) {
	// DELETE, /r/account/delete
}
