package database

import (
	"chronokeep/results/types"
	"chronokeep/results/util"
	"time"
)

const (
	MaxOpenConnections    = 20
	MaxIdleConnections    = 20
	MaxConnectionLifetime = time.Minute * 5
	CurrentVersion        = 6
	MaxLoginAttempts      = 4
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
	ResurrectAccount(email string) error
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
	GetResults(eventYearID int64, limit, page int) ([]types.Result, error)
	GetLastResults(eventYearID int64, limit, page int) ([]types.Result, error)
	GetDistanceResults(eventYearID int64, distance string, limit, page int) ([]types.Result, error)
	GetAllDistanceResults(eventYearID int64, distance string, limit, page int) ([]types.Result, error)
	GetFinishResults(eventYearID int64, distance string, limit, page int) ([]types.Result, error)
	GetBibResults(eventYearID int64, bib string) ([]types.Result, error)
	DeleteResults(eventYearID int64, results []types.Result) error
	DeleteEventResults(eventYearID int64) (int64, error)
	AddResults(eventYearID int64, results []types.Result) ([]types.Result, error)
	// Multi-Get Functions
	GetAccountAndEvent(slug string) (*types.MultiGet, error)
	GetAccountEventAndYear(slug, year string) (*types.MultiGet, error)
	GetEventAndYear(slug string, year string) (*types.MultiGet, error)
	GetKeyAndAccount(key string) (*types.MultiKey, error)
	// Person Functions
	GetPerson(slug, year, bib string) (*types.Person, error)
	GetPeople(slug, year string) ([]types.Person, error)
	AddPerson(eventYearID int64, person types.Person) (*types.Person, error)
	AddPeople(eventYearID int64, people []types.Person) ([]types.Person, error)
	DeletePeople(eventYearID int64, bibs []string) error
	// Close the database.
	Close()
}
