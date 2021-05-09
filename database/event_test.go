package database

import (
	"chronokeep/results/types"
	"testing"
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
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupEventTests()
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
	event, err := AddEvent(event1)
	if err != nil {
		t.Fatalf("Error adding event: %v", err)
	}
	t.Logf("New event ID: %v", event.Identifier)
	if !event.Equals(&event1) {
		t.Errorf("Event expected: %+v; Event found: %+v;", event1, *event)
	}
	event, err = AddEvent(event2)
	if err != nil {
		t.Fatalf("Error adding event: %v", err)
	}
	t.Logf("New event ID: %v", event.Identifier)
	if !event.Equals(&event2) {
		t.Errorf("Event expected: %+v; Event found: %+v;", event2, *event)
	}
	_, err = AddEvent(event2)
	if err == nil {
		t.Error("Expected an error when adding duplicate event.")
	}
}

func TestGetEvent(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupEventTests()
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
	testEvent, err := GetEvent(event1.Slug)
	if err != nil {
		t.Fatalf("Error getting event: %v", err)
	}
	if !testEvent.Equals(&event1) {
		t.Errorf("Event expected: %+v; Event found: %+v;", event1, *testEvent)
	}
	testEvent, err = GetEvent(event2.Slug)
	if err != nil {
		t.Fatalf("Error getting event: %v", err)
	}
	if !testEvent.Equals(&event2) {
		t.Errorf("Event expected: %+v; Event found: %+v;", event2, *testEvent)
	}
	testEvent, err = GetEvent("test")
	if err != nil {
		t.Fatalf("Error getting event: %v", err)
	}
	if testEvent != nil {
		t.Errorf("Unexpected event found: %+v;", *testEvent)
	}
}

func TestGetEvents(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupEventTests()
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
	event3 := types.Event{
		AccountIdentifier: account1.Identifier,
		Name:              "Event 3",
		Slug:              "event3",
		ContactEmail:      "event3@test.com",
		AccessRestricted:  false,
	}
	event4 := types.Event{
		AccountIdentifier: account2.Identifier,
		Name:              "Event 4",
		Slug:              "event4",
		ContactEmail:      "event4@test.com",
		AccessRestricted:  false,
	}
	events, err := GetEvents()
	if err != nil {
		t.Fatalf("Error attempting to get events: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("Expected %v events but found %v.", 0, len(events))
	}
	AddEvent(event1)
	AddEvent(event2)
	events, err = GetEvents()
	if err != nil {
		t.Fatalf("Error attempting to get events: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("Expected %v events but found %v.", 2, len(events))
	}
	AddEvent(event3)
	AddEvent(event4)
}

func TestGetAccountEvents(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupEventTests()
	account1, _ := AddAccount(accounts[0])
	account2, _ := AddAccount(accounts[1])
	account3, _ := AddAccount(accounts[2])
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
	event3 := types.Event{
		AccountIdentifier: account1.Identifier,
		Name:              "Event 3",
		Slug:              "event3",
		ContactEmail:      "event3@test.com",
		AccessRestricted:  false,
	}
	event4 := types.Event{
		AccountIdentifier: account2.Identifier,
		Name:              "Event 4",
		Slug:              "event4",
		ContactEmail:      "event4@test.com",
		AccessRestricted:  false,
	}
	t.Log("Adding one event for each account.")
	AddEvent(event1)
	AddEvent(event2)
	t.Logf("Account email: %v", account1.Email)
	events, err := GetAccountEvents(account1.Email)
	if err != nil {
		t.Fatalf("Error attempting to get events: %v", err)
	}
	if len(events) != 1 {
		t.Errorf("Expected %v events but found %v.", 1, len(events))
	}
	t.Logf("Account email: %v", account2.Email)
	events, err = GetAccountEvents(account2.Email)
	if err != nil {
		t.Fatalf("Error attempting to get events: %v", err)
	}
	if len(events) != 1 {
		t.Errorf("Expected %v events but found %v.", 1, len(events))
	}
	t.Log("Adding the final two events.")
	AddEvent(event3)
	AddEvent(event4)
	events, err = GetAccountEvents(account1.Email)
	if err != nil {
		t.Fatalf("Error attempting to get events: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("Expected %v events but found %v.", 2, len(events))
	}
	events, err = GetAccountEvents(account2.Email)
	if err != nil {
		t.Fatalf("Error attempting to get events: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("Expected %v events but found %v.", 2, len(events))
	}
	t.Log("Testing an account with no events.")
	events, err = GetAccountEvents(account3.Email)
	if err != nil {
		t.Fatalf("Error attempting to get events: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("Expected %v events but found %v.", 0, len(events))
	}
}

func TestDeleteEvent(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupEventTests()
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
	AddEvent(*event2)
	err = DeleteEvent(*event1)
	if err != nil {
		t.Fatalf("Error deleting event: %v", err)
	}
	event, _ := GetEvent(event1.Slug)
	if event != nil {
		t.Errorf("Found deleted event: %+v", *event)
	}
	events, _ := GetEvents()
	if len(events) != 1 {
		t.Errorf("Expected %v events but found %v.", 1, len(events))
	}
}

func TestUpdateEvent(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupEventTests()
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
	event1.AccountIdentifier = account2.Identifier
	event1.Name = "Updated Event Name"
	event1.Slug = "update"
	event1.ContactEmail = "event1_test@test.com"
	event1.Website = "https://test.com/"
	err = UpdateEvent(*event1)
	if err != nil {
		t.Fatalf("Error updating event: %v", err)
	}
	event, _ := GetEvent("event1")
	if event.AccountIdentifier != account1.Identifier {
		t.Errorf("Event account id changed from %v to %v.", account1.Identifier, event.AccountIdentifier)
	}
	if event.Name != "Event 1" {
		t.Errorf("Event name changed from %v to %v.", "Event 1", event.Name)
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
	err = UpdateEvent(*event2)
	if err != nil {
		t.Fatalf("Error updating event: %v", err)
	}
	event, _ = GetEvent("event2")
	if event.AccessRestricted != event2.AccessRestricted {
		t.Errorf("Expected access restricted value %v, found %v.", event2.AccessRestricted, event.AccessRestricted)
	}
	if event.Image != event2.Image {
		t.Errorf("Expected image %v, found %v.", event2.Image, event.Image)
	}
}
