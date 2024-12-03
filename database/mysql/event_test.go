package mysql

import (
	"chronokeep/results/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupEventTests() {
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

func TestAddEvent(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupEventTests()
	account1, _ := db.AddAccount(accounts[0])
	account2, _ := db.AddAccount(accounts[1])
	event1 := types.Event{
		AccountIdentifier: account1.Identifier,
		Name:              "Event 1",
		CertificateName:   "An Event",
		Slug:              "cascade-express-half-marathon",
		ContactEmail:      "event1@test.com",
		AccessRestricted:  false,
	}
	event2 := types.Event{
		AccountIdentifier: account2.Identifier,
		Name:              "Event 2",
		CertificateName:   "Another Event",
		Slug:              "event2",
		ContactEmail:      "event2@test.com",
		AccessRestricted:  true,
	}
	event, err := db.AddEvent(event1)
	if err != nil {
		t.Fatalf("Error adding event: %v", err)
	}
	t.Logf("New event ID: %v", event.Identifier)
	if !event.Equals(&event1) {
		t.Errorf("Event expected: %+v; Event found: %+v;", event1, *event)
	}
	event, err = db.AddEvent(event2)
	if err != nil {
		t.Fatalf("Error adding event: %v", err)
	}
	t.Logf("New event ID: %v", event.Identifier)
	if !event.Equals(&event2) {
		t.Errorf("Event expected: %+v; Event found: %+v;", event2, *event)
	}
	_, err = db.AddEvent(event2)
	if err == nil {
		t.Error("Expected an error when adding duplicate event.")
	}
}

func TestGetEvent(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupEventTests()
	account1, _ := db.AddAccount(accounts[0])
	account2, _ := db.AddAccount(accounts[1])
	event1 := types.Event{
		AccountIdentifier: account1.Identifier,
		Name:              "Event 1",
		CertificateName:   "An Event",
		Slug:              "event1",
		ContactEmail:      "event1@test.com",
		AccessRestricted:  false,
	}
	event2 := types.Event{
		AccountIdentifier: account2.Identifier,
		Name:              "Event 2",
		CertificateName:   "Another Event",
		Slug:              "event2",
		ContactEmail:      "event2@test.com",
		AccessRestricted:  true,
	}
	nEvent1, _ := db.AddEvent(event1)
	db.AddEvent(event2)
	eventYear1 := &types.EventYear{
		EventIdentifier: nEvent1.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 10, 06, 9, 1, 15, 0, time.Local),
		Live:            false,
	}
	db.AddEventYear(*eventYear1)
	testEvent, err := db.GetEvent(event1.Slug)
	if err != nil {
		t.Fatalf("Error getting event: %v", err)
	}
	if !testEvent.Equals(&event1) {
		t.Errorf("Event expected: %+v; Event found: %+v;", event1, *testEvent)
	}
	if testEvent.RecentTime == nil {
		t.Errorf("Expected time for an event but didn't find anything.")
	} else if !testEvent.RecentTime.Equal(eventYear1.DateTime) {
		t.Errorf("Expected time to be %v, found %v", eventYear1.DateTime, testEvent.RecentTime)
	}
	testEvent, err = db.GetEvent(event2.Slug)
	if err != nil {
		t.Fatalf("Error getting event: %v", err)
	}
	if !testEvent.Equals(&event2) {
		t.Errorf("Event expected: %+v; Event found: %+v;", event2, *testEvent)
	}
	if testEvent.RecentTime != nil {
		t.Errorf("Expected nil time for an event that has no event years. (2)")
	}
	testEvent, err = db.GetEvent("test")
	if err != nil {
		t.Fatalf("Error getting event: %v", err)
	}
	if testEvent != nil {
		t.Errorf("Unexpected event found: %+v;", *testEvent)
	}
}

func TestGetEvents(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupEventTests()
	account1, _ := db.AddAccount(accounts[0])
	account2, _ := db.AddAccount(accounts[1])
	event1 := types.Event{
		AccountIdentifier: account1.Identifier,
		Name:              "Event 1",
		CertificateName:   "An Event",
		Slug:              "event1",
		ContactEmail:      "event1@test.com",
		AccessRestricted:  false,
	}
	event2 := types.Event{
		AccountIdentifier: account2.Identifier,
		Name:              "Event 2",
		CertificateName:   "Another Event",
		Slug:              "event2",
		ContactEmail:      "event2@test.com",
		AccessRestricted:  true,
	}
	event3 := types.Event{
		AccountIdentifier: account1.Identifier,
		Name:              "Event 3",
		CertificateName:   "Oops Event",
		Slug:              "event3",
		ContactEmail:      "event3@test.com",
		AccessRestricted:  false,
	}
	event4 := types.Event{
		AccountIdentifier: account2.Identifier,
		Name:              "Event 4",
		CertificateName:   "Fun Event",
		Slug:              "event4",
		ContactEmail:      "event4@test.com",
		AccessRestricted:  false,
	}
	events, err := db.GetEvents()
	if err != nil {
		t.Fatalf("Error attempting to get events: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("Expected %v events but found %v.", 0, len(events))
	}
	nEvent1, _ := db.AddEvent(event1)
	db.AddEvent(event2)
	eventYear1 := &types.EventYear{
		EventIdentifier: nEvent1.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 10, 06, 9, 1, 15, 0, time.Local),
		Live:            false,
	}
	db.AddEventYear(*eventYear1)
	events, err = db.GetEvents()
	if err != nil {
		t.Fatalf("Error attempting to get events: %v", err)
	}
	// db.GetEvents does not get restricted events.  The second event added is restricted.
	if len(events) != 1 {
		t.Errorf("Expected %v events but found %v.", 1, len(events))
	}
	db.AddEvent(event3)
	db.AddEvent(event4)
	events, err = db.GetEvents()
	if err != nil {
		t.Fatalf("Error attempting to get events: %v", err)
	}
	// db.GetEvents does not get restricted events.  The second event added is restricted.
	if len(events) != 3 {
		t.Errorf("Expected %v events but found %v.", 3, len(events))
	}
	found := 0
	for _, ev := range events {
		if ev.RecentTime != nil {
			found++
			if !ev.RecentTime.Equal(eventYear1.DateTime) {
				t.Errorf("Expected time %v, found %v", eventYear1.DateTime, ev.RecentTime)
			}
		}
	}
	if found != 1 {
		t.Errorf("Expected to find %v events with times, found %v.", 1, found)
	}
}

func TestGetAllEvents(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupEventTests()
	account1, _ := db.AddAccount(accounts[0])
	account2, _ := db.AddAccount(accounts[1])
	event1 := types.Event{
		AccountIdentifier: account1.Identifier,
		Name:              "Event 1",
		CertificateName:   "An Event",
		Slug:              "event1",
		ContactEmail:      "event1@test.com",
		AccessRestricted:  false,
	}
	event2 := types.Event{
		AccountIdentifier: account2.Identifier,
		Name:              "Event 2",
		CertificateName:   "Another Event",
		Slug:              "event2",
		ContactEmail:      "event2@test.com",
		AccessRestricted:  true,
	}
	event3 := types.Event{
		AccountIdentifier: account1.Identifier,
		Name:              "Event 3",
		CertificateName:   "Oops Event",
		Slug:              "event3",
		ContactEmail:      "event3@test.com",
		AccessRestricted:  false,
	}
	event4 := types.Event{
		AccountIdentifier: account2.Identifier,
		Name:              "Event 4",
		CertificateName:   "Fun Event",
		Slug:              "event4",
		ContactEmail:      "event4@test.com",
		AccessRestricted:  false,
	}
	events, err := db.GetAllEvents()
	if err != nil {
		t.Fatalf("Error attempting to get events: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("Expected %v events but found %v.", 0, len(events))
	}
	nEvent1, _ := db.AddEvent(event1)
	db.AddEvent(event2)
	eventYear1 := &types.EventYear{
		EventIdentifier: nEvent1.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 10, 06, 9, 1, 15, 0, time.Local),
		Live:            false,
	}
	db.AddEventYear(*eventYear1)
	events, err = db.GetAllEvents()
	if err != nil {
		t.Fatalf("Error attempting to get events: %v", err)
	}
	// db.GetAllEvents does get restricted events.  The second event added is restricted.
	if len(events) != 2 {
		t.Errorf("Expected %v events but found %v.", 2, len(events))
	}
	db.AddEvent(event3)
	db.AddEvent(event4)
	events, err = db.GetAllEvents()
	if err != nil {
		t.Fatalf("Error attempting to get events: %v", err)
	}
	// db.GetAllEvents does get restricted events.  The second event added is restricted.
	if len(events) != 4 {
		t.Errorf("Expected %v events but found %v.", 4, len(events))
	}
	found := 0
	for _, ev := range events {
		if ev.RecentTime != nil {
			found++
			if !ev.RecentTime.Equal(eventYear1.DateTime) {
				t.Errorf("Expected time %v, found %v", eventYear1.DateTime, ev.RecentTime)
			}
		}
	}
	if found != 1 {
		t.Errorf("Expected to find %v events with times, found %v.", 1, found)
	}
}

func TestGetAccountEvents(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupEventTests()
	account1, _ := db.AddAccount(accounts[0])
	account2, _ := db.AddAccount(accounts[1])
	account3, _ := db.AddAccount(accounts[2])
	registration, _ := db.AddAccount(
		types.Account{
			Name:     "Registration",
			Email:    "registration@test.com",
			Type:     "registration",
			Password: testHashPassword("password"),
		})
	_ = db.LinkAccounts(*account1, *registration)
	assert.NoError(t, err)
	event1 := types.Event{
		AccountIdentifier: account1.Identifier,
		Name:              "Event 1",
		CertificateName:   "An Event",
		Slug:              "event1",
		ContactEmail:      "event1@test.com",
		AccessRestricted:  false,
	}
	event2 := types.Event{
		AccountIdentifier: account2.Identifier,
		Name:              "Event 2",
		CertificateName:   "Another Event",
		Slug:              "event2",
		ContactEmail:      "event2@test.com",
		AccessRestricted:  true,
	}
	event3 := types.Event{
		AccountIdentifier: account1.Identifier,
		Name:              "Event 3",
		CertificateName:   "Oops Event",
		Slug:              "event3",
		ContactEmail:      "event3@test.com",
		AccessRestricted:  false,
	}
	event4 := types.Event{
		AccountIdentifier: account2.Identifier,
		Name:              "Event 4",
		CertificateName:   "Fun Event",
		Slug:              "event4",
		ContactEmail:      "event4@test.com",
		AccessRestricted:  false,
	}
	t.Log("Verifying no events added for accounts.")
	t.Logf("Account email: %v", account1.Email)
	events, err := db.GetAccountEvents(account1.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(events))
	}
	t.Logf("Account email: %v", account2.Email)
	events, err = db.GetAccountEvents(account2.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(events))
	}
	t.Log("Adding one event for each account.")
	db.AddEvent(event1)
	db.AddEvent(event2)
	t.Logf("Account email: %v", account1.Email)
	events, err = db.GetAccountEvents(account1.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(events))
	}
	t.Logf("Account email: %v", account2.Email)
	events, err = db.GetAccountEvents(account2.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(events))
	}
	t.Log("Adding the final two events.")
	db.AddEvent(event3)
	db.AddEvent(event4)
	events, err = db.GetAccountEvents(account1.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 2, len(events))
	}
	events, err = db.GetAccountEvents(account2.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 2, len(events))
	}
	t.Log("Testing an account with no events.")
	events, err = db.GetAccountEvents(account3.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(events))
	}
	t.Log("Testing linked account fetch.")
	events, err = db.GetAccountEvents(registration.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 2, len(events))
	}
}

func TestDeleteEvent(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupEventTests()
	account1, _ := db.AddAccount(accounts[0])
	account2, _ := db.AddAccount(accounts[1])
	event1 := &types.Event{
		AccountIdentifier: account1.Identifier,
		Name:              "Event 1",
		CertificateName:   "An Event",
		Slug:              "event1",
		ContactEmail:      "event1@test.com",
		AccessRestricted:  false,
	}
	event2 := &types.Event{
		AccountIdentifier: account2.Identifier,
		Name:              "Event 2",
		CertificateName:   "Another Event",
		Slug:              "event2",
		ContactEmail:      "event2@test.com",
		AccessRestricted:  false,
	}
	event1, _ = db.AddEvent(*event1)
	eventYear1 := &types.EventYear{
		EventIdentifier: event1.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 10, 06, 9, 1, 15, 0, time.Local),
		Live:            false,
	}
	eventYear2 := &types.EventYear{
		EventIdentifier: event1.Identifier,
		Year:            "2020",
		DateTime:        time.Date(2020, 10, 05, 9, 0, 0, 0, time.Local),
		Live:            false,
	}
	db.AddEventYear(*eventYear1)
	db.AddEventYear(*eventYear2)
	db.AddEvent(*event2)
	err = db.DeleteEvent(*event1)
	if err != nil {
		t.Fatalf("Error deleting event: %v", err)
	}
	event, _ := db.GetEvent(event1.Slug)
	if event != nil {
		t.Errorf("Found deleted event: %+v", *event)
	}
	events, _ := db.GetEvents()
	if len(events) != 1 {
		t.Errorf("Expected %v events but found %v.", 1, len(events))
	}
	eventYears, _ := db.GetEventYears(event1.Slug)
	if len(eventYears) != 0 {
		t.Errorf("Expected to find %v event years after deletion of event but found %v.", 0, len(eventYears))
	}
}

func TestUpdateEvent(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupEventTests()
	account1, _ := db.AddAccount(accounts[0])
	account2, _ := db.AddAccount(accounts[1])
	event1 := &types.Event{
		AccountIdentifier: account1.Identifier,
		Name:              "Event 1",
		CertificateName:   "An Event",
		Slug:              "event1",
		ContactEmail:      "event1@test.com",
		AccessRestricted:  false,
	}
	event2 := &types.Event{
		AccountIdentifier: account2.Identifier,
		Name:              "Event 2",
		CertificateName:   "Another Event",
		Slug:              "event2",
		ContactEmail:      "event2@test.com",
		AccessRestricted:  true,
	}
	event1, _ = db.AddEvent(*event1)
	event2, _ = db.AddEvent(*event2)
	event1.AccountIdentifier = account2.Identifier
	event1.Name = "Updated Event Name"
	event1.CertificateName = "An Updated Event"
	event1.Slug = "update"
	event1.ContactEmail = "event1_test@test.com"
	event1.Website = "https://test.com/"
	err = db.UpdateEvent(*event1)
	if err != nil {
		t.Fatalf("Error updating event: %v", err)
	}
	event, _ := db.GetEvent("event1")
	if event.AccountIdentifier != account1.Identifier {
		t.Errorf("Event account id changed from %v to %v.", account1.Identifier, event.AccountIdentifier)
	}
	if event.Name == "Event 1" {
		t.Errorf("Event name not changed: %v to %v.", "Event 1", event.Name)
	}
	if event.CertificateName == "An Event" {
		t.Errorf("Certificate name not changed: %v to %v.", "An Event", event.CertificateName)
	}
	if event.Slug != "event1" {
		t.Errorf("Event name changed from %v to %v.", "event1", event.Slug)
	}
	if event.ContactEmail != event1.ContactEmail {
		t.Errorf("Expected contact email %v, found %v.", event1.ContactEmail, event.ContactEmail)
	}
	if event.Website != event1.Website {
		t.Errorf("Expected website %v, found %v.", event1.Website, event.Website)
	}
	event2.AccessRestricted = false
	event2.Image = "https://test.com/"
	err = db.UpdateEvent(*event2)
	if err != nil {
		t.Fatalf("Error updating event: %v", err)
	}
	event, _ = db.GetEvent("event2")
	if event.AccessRestricted != event2.AccessRestricted {
		t.Errorf("Expected access restricted value %v, found %v.", event2.AccessRestricted, event.AccessRestricted)
	}
	if event.Image != event2.Image {
		t.Errorf("Expected image %v, found %v.", event2.Image, event.Image)
	}
}

func TestBadDatabaseEvent(t *testing.T) {
	db := badTestSetup(t)
	_, err := db.GetEvent("")
	if err == nil {
		t.Fatal("Expected error getting event.")
	}
	_, err = db.GetEvents()
	if err == nil {
		t.Fatal("Expected error getting events.")
	}
	_, err = db.GetAccountEvents("")
	if err == nil {
		t.Fatal("Expected error getting account events.")
	}
	_, err = db.AddEvent(types.Event{})
	if err == nil {
		t.Fatal("Expected error adding event.")
	}
	err = db.DeleteEvent(types.Event{})
	if err == nil {
		t.Fatal("Expected error deleting event.")
	}
	err = db.RealDeleteEvent(types.Event{})
	if err == nil {
		t.Fatal("Expected error really deleting an event.")
	}
	err = db.UpdateEvent(types.Event{})
	if err == nil {
		t.Fatal("Expected error updating event.")
	}
}

func TestNoDatabaseEvent(t *testing.T) {
	db := MySQL{}
	_, err := db.GetEvent("")
	if err == nil {
		t.Fatal("Expected error getting event.")
	}
	_, err = db.GetEvents()
	if err == nil {
		t.Fatal("Expected error getting events.")
	}
	_, err = db.GetAccountEvents("")
	if err == nil {
		t.Fatal("Expected error getting account events.")
	}
	_, err = db.AddEvent(types.Event{})
	if err == nil {
		t.Fatal("Expected error adding event.")
	}
	err = db.DeleteEvent(types.Event{})
	if err == nil {
		t.Fatal("Expected error deleting event.")
	}
	err = db.RealDeleteEvent(types.Event{})
	if err == nil {
		t.Fatal("Expected error really deleting an event.")
	}
	err = db.UpdateEvent(types.Event{})
	if err == nil {
		t.Fatal("Expected error updating event.")
	}
}
