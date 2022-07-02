package mysql

import (
	"chronokeep/results/types"
	"testing"
	"time"
)

func TestGetPerson(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up tests. %v", err)
	}
	defer finalize(t)
	setupResultTests()
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
	person, err := db.GetPerson(event.Slug, eventYear.Year, results[0].Bib)
	if err != nil {
		t.Errorf("Error finding non existent person: %v", err)
	}
	if person != nil {
		t.Errorf("Found someone when no one should exist: %v", person)
	}
	db.AddResults(eventYear.Identifier, results)
	for _, res := range results {
		person, err = db.GetPerson(event.Slug, eventYear.Year, res.Bib)
		if err != nil {
			t.Errorf("Error finding non existent person: %v", err)
		}
		if person == nil {
			t.Errorf("Expected to find someone but did not.")
		} else if person.Bib != res.Bib {
			t.Errorf("Bib doesn't match: %v+ %v+", *person, res)
		} else if person.First != res.First {
			t.Errorf("First name doesn't match: %v+ %v+", *person, res)
		} else if person.Last != res.Last {
			t.Errorf("Last name doesn't match: %v+ %v+", *person, res)
		} else if person.Age != res.Age {
			t.Errorf("Age doesn't match: %v+ %v+", *person, res)
		} else if person.AgeGroup != res.AgeGroup {
			t.Errorf("Age Group doesn't match: %v+ %v+", *person, res)
		} else if person.Gender != res.Gender {
			t.Errorf("Gender doesn't match: %v+ %v+", *person, res)
		} else if person.Distance != res.Distance {
			t.Errorf("Distance doesn't match: %v+ %v+", *person, res)
		}
	}
}

func TestBadDatabasePerson(t *testing.T) {
	db := badTestSetup(t)
	_, err := db.GetPerson("", "", "")
	if err == nil {
		t.Fatalf("Expected error getting person")
	}
}

func TestNoDatabasePerson(t *testing.T) {
	db := MySQL{}
	_, err := db.GetPerson("", "", "")
	if err == nil {
		t.Fatalf("Expected error getting person")
	}
}
