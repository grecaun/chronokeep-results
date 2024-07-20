package sqlite

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
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupMultiTests()
	account1, _ := db.AddAccount(accounts[0])
	account2, _ := db.AddAccount(accounts[1])
	event1 := types.Event{
		AccountIdentifier: account1.Identifier,
		Name:              "Event 1",
		Slug:              "event1",
		ContactEmail:      "event1@test.com",
		AccessRestricted:  false,
		Type:              "time",
	}
	event2 := types.Event{
		AccountIdentifier: account2.Identifier,
		Name:              "Event 2",
		Slug:              "event2",
		ContactEmail:      "event2@test.com",
		AccessRestricted:  true,
		Type:              "distance",
	}
	db.AddEvent(event1)
	db.AddEvent(event2)
	mult, err := db.GetAccountAndEvent(event1.Slug)
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
	mult, err = db.GetAccountAndEvent(event2.Slug)
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
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupMultiTests()
	account1, _ := db.AddAccount(accounts[0])
	account2, _ := db.AddAccount(accounts[1])
	event1 := &types.Event{
		AccountIdentifier: account1.Identifier,
		Name:              "Event 1",
		Slug:              "event1",
		ContactEmail:      "event1@test.com",
		AccessRestricted:  false,
		Type:              "time",
	}
	event2 := &types.Event{
		AccountIdentifier: account2.Identifier,
		Name:              "Event 2",
		Slug:              "event2",
		ContactEmail:      "event2@test.com",
		AccessRestricted:  true,
		Type:              "distance",
	}
	event1, _ = db.AddEvent(*event1)
	event2, _ = db.AddEvent(*event2)
	eventYear1 := types.EventYear{
		EventIdentifier: event1.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 10, 06, 9, 0, 0, 0, time.Local),
		Live:            false,
		DaysAllowed:     1,
	}
	eventYear2 := types.EventYear{
		EventIdentifier: event2.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 04, 05, 11, 0, 0, 0, time.Local),
		Live:            false,
		DaysAllowed:     2,
	}
	db.AddEventYear(eventYear1)
	db.AddEventYear(eventYear2)
	mult, err := db.GetAccountEventAndYear(event1.Slug, eventYear1.Year)
	if err != nil {
		t.Fatalf("Error getting account, event, and year: %v", err)
	}
	if mult == nil || mult.Account == nil || mult.Event == nil || mult.EventYear == nil {
		t.Fatal("Account, Event, or EventYear was nil.")
	}
	if !mult.Account.Equals(&accounts[0]) || !mult.Event.Equals(event1) || !mult.EventYear.Equals(&eventYear1) {
		t.Errorf("Account expected: %+v; Found: %+v;\nEvent expected: %+v; Found: %+v;\nEventYear expected: %+v; Found %+v;", accounts[0], *mult.Account, *event1, *mult.Event, eventYear1, *mult.EventYear)
	}
	mult, err = db.GetAccountEventAndYear(event1.Slug, "")
	if err != nil {
		t.Fatalf("Error getting account, event, and year: %v", err)
	}
	if mult == nil || mult.Account == nil || mult.Event == nil || mult.EventYear == nil {
		t.Fatal("Account, Event, or EventYear was nil.")
	}
	if !mult.Account.Equals(&accounts[0]) || !mult.Event.Equals(event1) || !mult.EventYear.Equals(&eventYear1) {
		t.Errorf("Account expected: %+v; Found: %+v;\nEvent expected: %+v; Found: %+v;\nEventYear expected: %+v; Found %+v;", accounts[0], *mult.Account, *event1, *mult.Event, eventYear1, *mult.EventYear)
	}
	mult, err = db.GetAccountEventAndYear(event2.Slug, eventYear2.Year)
	if err != nil {
		t.Fatalf("Error getting account, event, and year: %v", err)
	}
	if mult == nil || mult.Account == nil || mult.Event == nil || mult.EventYear == nil {
		t.Fatal("Account, Event, or EventYear was nil.")
	}
	if !mult.Account.Equals(&accounts[1]) || !mult.Event.Equals(event2) || !mult.EventYear.Equals(&eventYear2) {
		t.Errorf("Account expected: %+v; Found: %+v;\nEvent expected: %+v; Found: %+v;\nEventYear expected: %+v; Found %+v;", accounts[1], *mult.Account, *event2, *mult.Event, eventYear2, *mult.EventYear)
	}
}

func TestGetEventAndYear(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupMultiTests()
	account1, _ := db.AddAccount(accounts[0])
	account2, _ := db.AddAccount(accounts[1])
	event1 := &types.Event{
		AccountIdentifier: account1.Identifier,
		Name:              "Event 1",
		Slug:              "event1",
		ContactEmail:      "event1@test.com",
		AccessRestricted:  false,
		Type:              "time",
	}
	event2 := &types.Event{
		AccountIdentifier: account2.Identifier,
		Name:              "Event 2",
		Slug:              "event2",
		ContactEmail:      "event2@test.com",
		AccessRestricted:  true,
		Type:              "distance",
	}
	event1, _ = db.AddEvent(*event1)
	event2, _ = db.AddEvent(*event2)
	eventYear1 := types.EventYear{
		EventIdentifier: event1.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 10, 06, 9, 0, 0, 0, time.Local),
		Live:            false,
		DaysAllowed:     1,
	}
	eventYear2 := types.EventYear{
		EventIdentifier: event2.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 04, 05, 11, 0, 0, 0, time.Local),
		Live:            false,
		DaysAllowed:     2,
	}
	db.AddEventYear(eventYear1)
	db.AddEventYear(eventYear2)
	mult, err := db.GetEventAndYear(event1.Slug, eventYear1.Year)
	if err != nil {
		t.Fatalf("Error getting event and year: %v", err)
	}
	if mult == nil || mult.Event == nil || mult.EventYear == nil {
		t.Fatal("Event or EventYear was nil.")
	}
	if !mult.Event.Equals(event1) || !mult.EventYear.Equals(&eventYear1) {
		t.Errorf("Event expected: %+v; Found: %+v;\nEventYear expected: %+v; Found %+v;", *event1, *mult.Event, eventYear1, *mult.EventYear)
	}
	mult, err = db.GetEventAndYear(event1.Slug, "")
	if err != nil {
		t.Fatalf("Error getting event and year: %v", err)
	}
	if mult == nil || mult.Event == nil || mult.EventYear == nil {
		t.Fatal("Event or EventYear was nil.")
	}
	if !mult.Event.Equals(event1) || !mult.EventYear.Equals(&eventYear1) {
		t.Errorf("Event expected: %+v; Found: %+v;\nEventYear expected: %+v; Found %+v;", *event1, *mult.Event, eventYear1, *mult.EventYear)
	}
	mult, err = db.GetEventAndYear(event2.Slug, eventYear2.Year)
	if err != nil {
		t.Fatalf("Error getting event and year: %v", err)
	}
	if mult == nil || mult.Event == nil || mult.EventYear == nil {
		t.Fatal("Event or EventYear was nil.")
	}
	if !mult.Event.Equals(event2) || !mult.EventYear.Equals(&eventYear2) {
		t.Errorf("Event expected: %+v; Found: %+v;\nEventYear expected: %+v; Found %+v;", *event2, *mult.Event, eventYear2, *mult.EventYear)
	}
}

func TestGetKeyAndAccount(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupMultiTests()
	account1, _ := db.AddAccount(accounts[0])
	account2, _ := db.AddAccount(accounts[1])
	times := []time.Time{
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
		time.Date(2016, 4, 1, 4, 11, 5, 0, time.Local),
	}
	keys := []types.Key{
		{
			AccountIdentifier: account1.Identifier,
			Value:             "030001-1ACSDD-K2389A-00123B",
			Type:              "default",
			AllowedHosts:      "",
			ValidUntil:        &times[0],
		},
		{
			AccountIdentifier: account2.Identifier,
			Value:             "030001-1ACSDD-KH789A-00123B",
			Type:              "delete",
			AllowedHosts:      "https://test.com/",
			ValidUntil:        &times[1],
		},
	}
	db.AddKey(keys[0])
	db.AddKey(keys[1])
	mult, err := db.GetKeyAndAccount(keys[0].Value)
	if err != nil {
		t.Fatalf("Error getting key and account: %v", err)
	}
	if mult == nil || mult.Key == nil || mult.Account == nil {
		t.Fatal("Key or Account was nil.")
	}
	if !mult.Account.Equals(account1) || !mult.Key.Equal(&keys[0]) {
		t.Errorf("Account expected: %+v; Found %+v;\nKey expected: %+v; Found %+v;", *account1, *mult.Account, keys[0], *mult.Key)
	}
	mult, err = db.GetKeyAndAccount(keys[1].Value)
	if err != nil {
		t.Fatalf("Error getting key and account: %v", err)
	}
	if mult == nil || mult.Key == nil || mult.Account == nil {
		t.Fatal("Key or Account was nil.")
	}
	if !mult.Account.Equals(account2) || !mult.Key.Equal(&keys[1]) {
		t.Errorf("Account expected: %+v; Found %+v;\nKey expected: %+v; Found %+v;", *account2, *mult.Account, keys[1], *mult.Key)
	}
}

func TestBadDatabaseMultiGet(t *testing.T) {
	db := badTestSetup(t)
	_, err := db.GetAccountAndEvent("")
	if err == nil {
		t.Fatal("Expected error on get account and event.")
	}
	_, err = db.GetAccountEventAndYear("", "")
	if err == nil {
		t.Fatal("Expected error on get account, event, and year.")
	}
	_, err = db.GetEventAndYear("", "")
	if err == nil {
		t.Fatal("Expected error on get event and year.")
	}
	_, err = db.GetKeyAndAccount("")
	if err == nil {
		t.Fatal("Expected error on get account and key.")
	}
}

func TestNoDatabaseMultiGet(t *testing.T) {
	db := SQLite{}
	_, err := db.GetAccountAndEvent("")
	if err == nil {
		t.Fatal("Expected error on get account and event.")
	}
	_, err = db.GetAccountEventAndYear("", "")
	if err == nil {
		t.Fatal("Expected error on get account, event, and year.")
	}
	_, err = db.GetEventAndYear("", "")
	if err == nil {
		t.Fatal("Expected error on get event and year.")
	}
	_, err = db.GetKeyAndAccount("")
	if err == nil {
		t.Fatal("Expected error on get account and key.")
	}
}
