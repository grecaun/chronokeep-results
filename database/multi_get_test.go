package database

import (
	"chronokeep/results/types"
	"testing"
	"time"
)

func setupMultiTests() {
	if len(accounts) < 1 {
		accounts = []types.Account{
			{
				Name:     "John Smith",
				Email:    "j@test.com",
				Type:     "admin",
				Password: testHashPassword("password"),
			},
			{
				Name:     "Rose MacDonald",
				Email:    "rose2004@test.com",
				Type:     "paid",
				Password: testHashPassword("password"),
			},
			{
				Name:     "Tia Johnson",
				Email:    "tiatheway@test.com",
				Type:     "free",
				Password: testHashPassword("password"),
			},
		}
	}
}

func TestGetAccountAndEvent(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupMultiTests()
	account1, _ := AddAccount(accounts[0])
	account2, _ := AddAccount(accounts[1])
	event1 := types.Event{
		AccountIdentifier: account1.Identifier,
		Name:              "Event 1",
		Slug:              "event1",
		ContactEmail:      "event1@test.com",
		AccessRestricted:  false,
	}
	event2 := types.Event{
		AccountIdentifier: account2.Identifier,
		Name:              "Event 2",
		Slug:              "event2",
		ContactEmail:      "event2@test.com",
		AccessRestricted:  true,
	}
	AddEvent(event1)
	AddEvent(event2)
	mult, err := GetAccountAndEvent(event1.Slug)
	if err != nil {
		t.Fatalf("Error getting account and event: %v", err)
	}
	if mult == nil {
		t.Fatal("Mult was nil")
	}
	if mult.Account == nil || mult.Event == nil {
		t.Fatal("Account or Event was nil.")
	}
	if !mult.Account.Equals(&accounts[0]) || !mult.Event.Equals(&event1) {
		t.Errorf("Account expected: %+v; Found: %+v;\nEvent Expected: %+v; Found: %+v;", accounts[0], *mult.Account, event1, *mult.Event)
	}
	mult, err = GetAccountAndEvent(event2.Slug)
	if err != nil {
		t.Fatalf("Error getting account and event: %v", err)
	}
	if mult == nil {
		t.Fatal("Mult was nil")
	}
	if mult.Account == nil || mult.Event == nil {
		t.Fatal("Account or Event was nil.")
	}
	if !mult.Account.Equals(&accounts[1]) || !mult.Event.Equals(&event2) {
		t.Errorf("Account expected: %+v; Found: %+v;\nEvent Expected: %+v; Found: %+v;", accounts[1], *mult.Account, event2, *mult.Event)
	}
}

func TestGetAccountEventAndYear(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupMultiTests()
	account1, _ := AddAccount(accounts[0])
	account2, _ := AddAccount(accounts[1])
	event1 := &types.Event{
		AccountIdentifier: account1.Identifier,
		Name:              "Event 1",
		Slug:              "event1",
		ContactEmail:      "event1@test.com",
		AccessRestricted:  false,
	}
	event2 := &types.Event{
		AccountIdentifier: account2.Identifier,
		Name:              "Event 2",
		Slug:              "event2",
		ContactEmail:      "event2@test.com",
		AccessRestricted:  true,
	}
	event1, _ = AddEvent(*event1)
	event2, _ = AddEvent(*event2)
	eventYear1 := types.EventYear{
		EventIdentifier: event1.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 10, 06, 9, 0, 0, 0, time.Local),
		Live:            false,
	}
	eventYear2 := types.EventYear{
		EventIdentifier: event2.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 04, 05, 11, 0, 0, 0, time.Local),
		Live:            false,
	}
	AddEventYear(eventYear1)
	AddEventYear(eventYear2)
	mult, err := GetAccountEventAndYear(event1.Slug, eventYear1.Year)
	if err != nil {
		t.Fatalf("Error getting account and event: %v", err)
	}
	if mult == nil {
		t.Fatal("Mult was nil")
	}
	if mult.Account == nil || mult.Event == nil || mult.EventYear == nil {
		t.Fatal("Account or Event or EventYear was nil.")
	}
	if !mult.Account.Equals(&accounts[0]) || !mult.Event.Equals(event1) || !mult.EventYear.Equals(&eventYear1) {
		t.Errorf("Account expected: %+v; Found: %+v;\nEvent expected: %+v; Found: %+v;\nEventYear expected: %+v; Found %+v;", accounts[0], *mult.Account, *event1, *mult.Event, eventYear1, *mult.EventYear)
	}
	mult, err = GetAccountEventAndYear(event2.Slug, eventYear2.Year)
	if err != nil {
		t.Fatalf("Error getting account and event: %v", err)
	}
	if mult == nil {
		t.Fatal("Mult was nil")
	}
	if mult.Account == nil || mult.Event == nil || mult.EventYear == nil {
		t.Fatal("Account or Event or EventYear was nil.")
	}
	if !mult.Account.Equals(&accounts[1]) || !mult.Event.Equals(event2) || !mult.EventYear.Equals(&eventYear2) {
		t.Errorf("Account expected: %+v; Found: %+v;\nEvent expected: %+v; Found: %+v;\nEventYear expected: %+v; Found %+v;", accounts[1], *mult.Account, *event2, *mult.Event, eventYear2, *mult.EventYear)
	}
}
