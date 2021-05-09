package database

import (
	"chronokeep/results/types"
	"testing"
	"time"
)

func setupEventYearTests() {
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
		}
	}
}

func TestAddEventYear(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupEventYearTests()
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
	eventYear1 := &types.EventYear{
		EventIdentifier: event1.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 10, 06, 9, 0, 0, 0, time.Local),
		Live:            false,
	}
	eventYear2 := &types.EventYear{
		EventIdentifier: event1.Identifier,
		Year:            "2020",
		DateTime:        time.Date(2020, 10, 05, 9, 0, 0, 0, time.Local),
		Live:            false,
	}
	eventYear3 := &types.EventYear{
		EventIdentifier: event2.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 04, 05, 11, 0, 0, 0, time.Local),
		Live:            false,
	}
	eventYear4 := &types.EventYear{
		EventIdentifier: event2.Identifier,
		Year:            "2020",
		DateTime:        time.Date(2021, 04, 05, 11, 0, 0, 0, time.Local),
		Live:            false,
	}
	out, err := AddEventYear(*eventYear1)
	if err != nil {
		t.Fatalf("Error adding event year: %v", err)
	}
	if !out.Equals(eventYear1) {
		t.Errorf("Expected event year %+v, found %+v.", *eventYear1, *out)
	}
	out, err = AddEventYear(*eventYear2)
	if err != nil {
		t.Fatalf("Error adding event year: %v", err)
	}
	if !out.Equals(eventYear2) {
		t.Errorf("Expected event year %+v, found %+v.", *eventYear2, *out)
	}
	out, err = AddEventYear(*eventYear3)
	if err != nil {
		t.Fatalf("Error adding event year: %v", err)
	}
	if !out.Equals(eventYear3) {
		t.Errorf("Expected event year %+v, found %+v.", *eventYear3, *out)
	}
	out, err = AddEventYear(*eventYear4)
	if err != nil {
		t.Fatalf("Error adding event year: %v", err)
	}
	if !out.Equals(eventYear4) {
		t.Errorf("Expected event year %+v, found %+v.", *eventYear4, *out)
	}
}

func TestGetEventYear(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupEventYearTests()
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
	eventYear1 := &types.EventYear{
		EventIdentifier: event1.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 10, 06, 9, 0, 0, 0, time.Local),
		Live:            false,
	}
	eventYear2 := &types.EventYear{
		EventIdentifier: event1.Identifier,
		Year:            "2020",
		DateTime:        time.Date(2020, 10, 05, 9, 0, 0, 0, time.Local),
		Live:            false,
	}
	eventYear3 := &types.EventYear{
		EventIdentifier: event2.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 04, 05, 11, 0, 0, 0, time.Local),
		Live:            false,
	}
	eventYear4 := &types.EventYear{
		EventIdentifier: event2.Identifier,
		Year:            "2020",
		DateTime:        time.Date(2021, 04, 05, 11, 0, 0, 0, time.Local),
		Live:            false,
	}
	AddEventYear(*eventYear1)
	AddEventYear(*eventYear2)
	AddEventYear(*eventYear3)
	AddEventYear(*eventYear4)
	eyear, err := GetEventYear(event1.Slug, eventYear1.Year)
	if err != nil {
		t.Fatalf("Error getting event year: %v", err)
	}
	if !eyear.Equals(eventYear1) {
		t.Errorf("Expected event year %+v, found %+v.", *eventYear1, *eyear)
	}
	eyear, err = GetEventYear(event1.Slug, eventYear2.Year)
	if err != nil {
		t.Fatalf("Error getting event year: %v", err)
	}
	if !eyear.Equals(eventYear2) {
		t.Errorf("Expected event year %+v, found %+v.", *eventYear2, *eyear)
	}
	eyear, err = GetEventYear(event2.Slug, eventYear3.Year)
	if err != nil {
		t.Fatalf("Error getting event year: %v", err)
	}
	if !eyear.Equals(eventYear3) {
		t.Errorf("Expected event year %+v, found %+v.", *eventYear3, *eyear)
	}
	eyear, err = GetEventYear(event2.Slug, eventYear4.Year)
	if err != nil {
		t.Fatalf("Error getting event year: %v", err)
	}
	if !eyear.Equals(eventYear4) {
		t.Errorf("Expected event year %+v, found %+v.", *eventYear4, *eyear)
	}
	eyear, err = GetEventYear(event1.Slug, "2000")
	if err != nil {
		t.Fatalf("Error getting event year: %v", err)
	}
	if eyear != nil {
		t.Errorf("Expected a nil event year but found %+v.", *eyear)
	}
	eyear, err = GetEventYear("testevent", "2000")
	if err != nil {
		t.Fatalf("Error getting event year: %v", err)
	}
	if eyear != nil {
		t.Errorf("Expected a nil event year but found %+v.", *eyear)
	}
}

func TestGetEventYears(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupEventYearTests()
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
	event3 := &types.Event{
		AccountIdentifier: account2.Identifier,
		Name:              "Event 3",
		Slug:              "event3",
		ContactEmail:      "event3@test.com",
		AccessRestricted:  true,
	}
	event1, _ = AddEvent(*event1)
	event2, _ = AddEvent(*event2)
	event3, _ = AddEvent(*event3)
	eventYear1 := &types.EventYear{
		EventIdentifier: event1.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 10, 06, 9, 0, 0, 0, time.Local),
		Live:            false,
	}
	eventYear2 := &types.EventYear{
		EventIdentifier: event1.Identifier,
		Year:            "2020",
		DateTime:        time.Date(2020, 10, 05, 9, 0, 0, 0, time.Local),
		Live:            false,
	}
	eventYear3 := &types.EventYear{
		EventIdentifier: event2.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 04, 05, 11, 0, 0, 0, time.Local),
		Live:            false,
	}
	eventYear4 := &types.EventYear{
		EventIdentifier: event2.Identifier,
		Year:            "2020",
		DateTime:        time.Date(2021, 04, 05, 11, 0, 0, 0, time.Local),
		Live:            false,
	}
	AddEventYear(*eventYear1)
	AddEventYear(*eventYear2)
	AddEventYear(*eventYear3)
	AddEventYear(*eventYear4)
	years, err := GetEventYears(event1.Slug)
	if err != nil {
		t.Fatalf("Error getting event years: %v", err)
	}
	if len(years) != 2 {
		t.Errorf("Expected %v years but found %v.", 2, len(years))
	}
	if (!years[0].Equals(eventYear1) && !years[0].Equals(eventYear2)) ||
		(!years[1].Equals(eventYear1) && !years[1].Equals(eventYear2)) {
		t.Errorf("Event years expected %+v %+v; Found %+v %+v;", *eventYear1, *eventYear2, years[0], years[1])
	}
	years, err = GetEventYears(event2.Slug)
	if err != nil {
		t.Fatalf("Error getting event years: %v", err)
	}
	if len(years) != 2 {
		t.Errorf("Expected %v years but found %v.", 2, len(years))
	}
	if (!years[0].Equals(eventYear3) && !years[0].Equals(eventYear4)) ||
		(!years[1].Equals(eventYear3) && !years[1].Equals(eventYear4)) {
		t.Errorf("Event years expected %+v %+v; Found %+v %+v;", *eventYear3, *eventYear4, years[0], years[1])
	}
	years, err = GetEventYears(event3.Slug)
	if err != nil {
		t.Fatalf("Error getting event years: %v", err)
	}
	if len(years) != 0 {
		t.Errorf("Expected %v years but found %v.", 0, len(years))
	}
}

func TestDeleteEventYear(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupEventYearTests()
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
	eventYear1 := &types.EventYear{
		EventIdentifier: event1.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 10, 06, 9, 0, 0, 0, time.Local),
		Live:            false,
	}
	eventYear2 := &types.EventYear{
		EventIdentifier: event1.Identifier,
		Year:            "2020",
		DateTime:        time.Date(2020, 10, 05, 9, 0, 0, 0, time.Local),
		Live:            false,
	}
	eventYear3 := &types.EventYear{
		EventIdentifier: event2.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 04, 05, 11, 0, 0, 0, time.Local),
		Live:            false,
	}
	eventYear4 := &types.EventYear{
		EventIdentifier: event2.Identifier,
		Year:            "2020",
		DateTime:        time.Date(2021, 04, 05, 11, 0, 0, 0, time.Local),
		Live:            false,
	}
	AddEventYear(*eventYear1)
	eventYear2, _ = AddEventYear(*eventYear2)
	eventYear3, _ = AddEventYear(*eventYear3)
	AddEventYear(*eventYear4)
	err = DeleteEventYear(*eventYear2)
	if err != nil {
		t.Fatalf("Error deleting event year: %v", err)
	}
	year, _ := GetEventYear(event1.Slug, eventYear2.Year)
	if year != nil {
		t.Errorf("Unexpectedly found deleted event year %+v.", year)
	}
	err = DeleteEventYear(*eventYear3)
	if err != nil {
		t.Fatalf("Error deleting event year: %v", err)
	}
	year, _ = GetEventYear(event2.Slug, eventYear3.Year)
	if year != nil {
		t.Errorf("Unexpectedly found deleted event year %+v.", *year)
	}
	years, _ := GetEventYears(event1.Slug)
	if len(years) != 1 {
		t.Errorf("Expected to find %v event years, found %v.", 1, len(years))
	}
}

func TestUpdateEventYear(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupEventYearTests()
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
	eventYear1 := &types.EventYear{
		EventIdentifier: event1.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 10, 06, 9, 0, 0, 0, time.Local),
		Live:            false,
	}
	eventYear2 := &types.EventYear{
		EventIdentifier: event1.Identifier,
		Year:            "2020",
		DateTime:        time.Date(2020, 10, 05, 9, 0, 0, 0, time.Local),
		Live:            false,
	}
	eventYear1, _ = AddEventYear(*eventYear1)
	eventYear2, _ = AddEventYear(*eventYear2)
	eventYear1.DateTime = time.Date(2000, 10, 05, 10, 10, 0, 0, time.Local)
	eventYear1.Live = true
	err = UpdateEventYear(*eventYear1)
	if err != nil {
		t.Fatalf("Error updating event year: %v", err)
	}
	event, _ := GetEventYear(event1.Slug, eventYear1.Year)
	if !event.Equals(eventYear1) {
		t.Errorf("Expected event year %+v, found %+v.", *eventYear1, *event)
	}
	eventYear2.EventIdentifier = event2.Identifier
	eventYear2.Identifier = eventYear2.Identifier + 200
	eventYear2.Year = "1999"
	err = UpdateEventYear(*eventYear2)
	if err != nil {
		t.Fatalf("Error updating event year: %v", err)
	}
	event, _ = GetEventYear(event1.Slug, "2020")
	if event.EventIdentifier == eventYear2.EventIdentifier {
		t.Errorf("Event identifier changed, found %v", event.EventIdentifier)
	}
	if event.Identifier == eventYear2.Identifier {
		t.Errorf("Event Year Identifier changed, found %v", event.Identifier)
	}
	if event.Year == eventYear2.Year {
		t.Errorf("Event Year year value changed, found %v", event.Year)
	}
}
