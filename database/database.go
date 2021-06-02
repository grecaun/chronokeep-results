package database

import (
	"chronokeep/results/types"
	"chronokeep/results/util"
)

type Database interface {
	// Database Base Functions
	Setup(config *util.Config) error
	SetSetting(name, value string) error
	// Account Functions
	GetAccount(email string) (*types.Account, error)
	GetAccountByKey(key string) (*types.Account, error)
	GetAccountByID(id int64) (*types.Account, error)
	GetAccounts() ([]types.Account, error)
	AddAccount(account types.Account) (*types.Account, error)
	DeleteAccount(id int64) error
	ResurrectAccount(emnail string) error
	GetDeletedAccount(email string) (*types.Account, error)
	UpdateAccount(account types.Account) error
	ChangePassword(email, newPassword string, logout ...bool) error
	ChangeEmail(oldEmail, newEmail string) error
	InvalidPassword(account types.Account) error
	ValidPassword(account types.Account) error
	UnlockAccount(account types.Account) error
	UpdateTokens(account types.Account) error
	// Call Record Functions
	GetAccountCallRecords(email string) ([]types.CallRecord, error)
	GetCallRecord(email string, inTime int64) (*types.CallRecord, error)
	AddCallRecord(record types.CallRecord) error
	AddCallRecords(records []types.CallRecord) error
	// EventYear Functions
	GetEventYear(event_slug, year string) (*types.EventYear, error)
	GetEventYears(event_slug string) ([]types.EventYear, error)
	AddEventYear(year types.EventYear) (*types.EventYear, error)
	DeleteEventYear(year types.EventYear) error
	UpdateEventYear(year types.EventYear) error
	// Event Functions
	GetEvent(slug string) (*types.Event, error)
	GetEvents() ([]types.Event, error)
	GetAccountEvents(email string) ([]types.Event, error)
	AddEvent(event types.Event) (*types.Event, error)
	DeleteEvent(event types.Event) error
	UpdateEvent(event types.Event) error
	// Key Functions
	GetAccountKeys(email string) ([]types.Key, error)
	GetKey(key string) (*types.Key, error)
	AddKey(key types.Key) (*types.Key, error)
	DeleteKey(key types.Key) error
	UpdateKey(key types.Key) error
	// Result Functions
	GetResults(eventYearID int64) ([]types.Result, error)
	GetLastResults(eventYearID int64) ([]types.Result, error)
	GetBibResults(eventYearID int64, bib string) ([]types.Result, error)
	DeleteResults(eventYearID int64, results []types.Result) error
	DeleteEventResults(eventYearID int64) (int64, error)
	AddResults(eventYearID int64, results []types.Result) ([]types.Result, error)
	// Multi-Get Functions
	GetAccountAndEvent(slug string) (*types.MultiGet, error)
	GetAccountEventAndYear(slug, year string) (*types.MultiGet, error)
	GetEventAndYear(slug, year string) (*types.MultiGet, error)
	GetKeyAndAccount(key string) (*types.MultiKey, error)
	// Close the database.
	Close()
}
