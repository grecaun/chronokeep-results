package mysql

import (
	"chronokeep/results/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupSmsTests() []types.SmsSubscription {
	if len(accounts) < 1 {
		accounts = []types.Account{
			{
				Name:     "John Smith",
				Email:    "j@test.com",
				Type:     "admin",
				Password: testHashPassword("password"),
			},
		}
	}
	return []types.SmsSubscription{
		{
			Bib:   "1001",
			First: "",
			Last:  "",
			Phone: "1235557890",
		},
		{
			Bib:   "",
			First: "John",
			Last:  "Smith",
			Phone: "1325557890",
		},
		{
			Bib:   "100",
			First: "",
			Last:  "",
			Phone: "1235557890",
		},
	}
}

func TestAddSubscribedPhone(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	subs := setupSmsTests()
	account, _ := db.AddAccount(accounts[0])
	event := &types.Event{
		AccountIdentifier: account.Identifier,
		Name:              "Event 1",
		Slug:              "event1",
	}
	event, _ = db.AddEvent(*event)
	eventYear := &types.EventYear{
		EventIdentifier: event.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 04, 20, 9, 0, 0, 0, time.Local),
	}
	eventYear, _ = db.AddEventYear(*eventYear)
	defer finalize(t)
	err = db.AddSubscribedPhone(eventYear.Identifier, subs[0])
	if assert.NoError(t, err) {
		added, _ := db.GetSubscribedPhones(eventYear.Identifier)
		assert.Equal(t, 1, len(added))
		assert.Equal(t, subs[0].Bib, added[0].Bib)
		assert.Equal(t, subs[0].First, added[0].First)
		assert.Equal(t, subs[0].Last, added[0].Last)
		assert.Equal(t, subs[0].Phone, added[0].Phone)
	}
	err = db.AddSubscribedPhone(eventYear.Identifier, subs[0])
	if assert.NoError(t, err) {
		added, _ := db.GetSubscribedPhones(eventYear.Identifier)
		assert.Equal(t, 1, len(added))
	}
	err = db.AddSubscribedPhone(eventYear.Identifier, subs[1])
	assert.NoError(t, err)
	err = db.AddSubscribedPhone(eventYear.Identifier, subs[2])
	if assert.NoError(t, err) {
		added, _ := db.GetSubscribedPhones(eventYear.Identifier)
		assert.Equal(t, 3, len(added))
		for _, outer := range subs {
			found := false
			for _, inner := range added {
				if outer.Equals(&inner) {
					found = true
				}
			}
			assert.True(t, found)
		}
	}
}

func TestRemoveSubscribedPhone(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	subs := setupSmsTests()
	account, _ := db.AddAccount(accounts[0])
	event := &types.Event{
		AccountIdentifier: account.Identifier,
		Name:              "Event 1",
		Slug:              "event1",
	}
	event, _ = db.AddEvent(*event)
	eventYear := &types.EventYear{
		EventIdentifier: event.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 04, 20, 9, 0, 0, 0, time.Local),
	}
	eventYear, _ = db.AddEventYear(*eventYear)
	defer finalize(t)
	db.AddSubscribedPhone(eventYear.Identifier, subs[0])
	db.AddSubscribedPhone(eventYear.Identifier, subs[1])
	db.AddSubscribedPhone(eventYear.Identifier, subs[2])
	added, _ := db.GetSubscribedPhones(eventYear.Identifier)
	assert.Equal(t, 3, len(added))
	err = db.RemoveSubscribedPhone(eventYear.Identifier, subs[0].Phone)
	if assert.NoError(t, err) {
		added, _ = db.GetSubscribedPhones(eventYear.Identifier)
		assert.Equal(t, 1, len(added))
	}
	err = db.RemoveSubscribedPhone(eventYear.Identifier, subs[0].Phone)
	if assert.NoError(t, err) {
		added, _ = db.GetSubscribedPhones(eventYear.Identifier)
		assert.Equal(t, 1, len(added))
	}
	err = db.RemoveSubscribedPhone(eventYear.Identifier, subs[1].Phone)
	if assert.NoError(t, err) {
		added, _ = db.GetSubscribedPhones(eventYear.Identifier)
		assert.Equal(t, 0, len(added))
	}
}

func TestGetSubscribedPhones(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	subs := setupSmsTests()
	account, _ := db.AddAccount(accounts[0])
	event := &types.Event{
		AccountIdentifier: account.Identifier,
		Name:              "Event 1",
		Slug:              "event1",
	}
	event, _ = db.AddEvent(*event)
	eventYear := &types.EventYear{
		EventIdentifier: event.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 04, 20, 9, 0, 0, 0, time.Local),
	}
	eventYear, _ = db.AddEventYear(*eventYear)
	defer finalize(t)
	db.AddSubscribedPhone(eventYear.Identifier, subs[0])
	added, err := db.GetSubscribedPhones(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(added))
		assert.Equal(t, subs[0].Bib, added[0].Bib)
		assert.Equal(t, subs[0].First, added[0].First)
		assert.Equal(t, subs[0].Last, added[0].Last)
		assert.Equal(t, subs[0].Phone, added[0].Phone)
	}
	db.AddSubscribedPhone(eventYear.Identifier, subs[1])
	added, err = db.GetSubscribedPhones(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, 2, len(added))
	}
	db.AddSubscribedPhone(eventYear.Identifier, subs[2])
	added, err = db.GetSubscribedPhones(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, 3, len(added))
		for _, outer := range subs {
			found := false
			for _, inner := range added {
				if outer.Equals(&inner) {
					found = true
				}
			}
			assert.True(t, found)
		}
	}
}

func TestBadDatabaseSubscription(t *testing.T) {
	db := badTestSetup(t)
	subs := setupSmsTests()
	_, err := db.GetSubscribedPhones(0)
	assert.Error(t, err)
	err = db.AddSubscribedPhone(0, subs[0])
	assert.Error(t, err)
	err = db.RemoveSubscribedPhone(0, "")
	assert.Error(t, err)
}
