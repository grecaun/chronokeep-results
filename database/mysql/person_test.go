package mysql

import (
	"chronokeep/results/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	people []types.Person
)

func setupPeople() {
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
	people = []types.Person{
		{
			Bib:      "100",
			First:    "John",
			Last:     "Smith",
			Age:      24,
			Gender:   "M",
			AgeGroup: "20-29",
			Distance: "1 Mile",
		},
		{
			Bib:      "106",
			First:    "Rose",
			Last:     "Johnson",
			Age:      55,
			Gender:   "F",
			AgeGroup: "50-59",
			Distance: "1 Mile",
		},
		{
			Bib:      "209",
			First:    "Tony",
			Last:     "Starke",
			Age:      45,
			Gender:   "M",
			AgeGroup: "40-49",
			Distance: "1 Mile",
		},
		{
			Bib:      "287",
			First:    "Jamie",
			Last:     "Fischer",
			Age:      35,
			Gender:   "NB",
			AgeGroup: "30-39",
			Distance: "5 Mile",
		},
		{
			Bib:      "2871",
			First:    "Jamie",
			Last:     "Fischer",
			Age:      35,
			Gender:   "F",
			AgeGroup: "30-39",
			Distance: "5 Mile",
		},
	}
}

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

func TestGetPeople(t *testing.T) {
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
	people, err := db.GetPeople(event.Slug, eventYear.Year)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(people))
	}
	db.AddResults(eventYear.Identifier, results)
	people, err = db.GetPeople(event.Slug, eventYear.Year)
	if assert.NoError(t, err) {
		assert.Equal(t, len(results)-1, len(people))
		for _, outer := range results {
			found := false
			for _, inner := range people {
				if outer.Bib == inner.Bib {
					assert.Equal(t, outer.Age, inner.Age)
					assert.Equal(t, outer.AgeGroup, inner.AgeGroup)
					assert.Equal(t, outer.Bib, inner.Bib)
					assert.Equal(t, outer.Distance, inner.Distance)
					assert.Equal(t, outer.First, inner.First)
					assert.Equal(t, outer.Gender, inner.Gender)
					assert.Equal(t, outer.Last, inner.Last)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
}

func TestAddPerson(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up tests. %v", err)
	}
	defer finalize(t)
	setupPeople()
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
	person, err := db.AddPerson(eventYear.Identifier, people[0])
	if assert.NoError(t, err) {
		assert.Equal(t, people[0].Age, person.Age)
		assert.Equal(t, people[0].AgeGroup, person.AgeGroup)
		assert.Equal(t, people[0].Bib, person.Bib)
		assert.Equal(t, people[0].Distance, person.Distance)
		assert.Equal(t, people[0].First, person.First)
		assert.Equal(t, people[0].Gender, person.Gender)
		assert.Equal(t, people[0].Last, person.Last)
	}
	id := person.Identifier
	person, err = db.GetPerson(event.Slug, eventYear.Year, people[0].Bib)
	if assert.NoError(t, err) {
		assert.Equal(t, people[0].Age, person.Age)
		assert.Equal(t, people[0].AgeGroup, person.AgeGroup)
		assert.Equal(t, people[0].Bib, person.Bib)
		assert.Equal(t, people[0].Distance, person.Distance)
		assert.Equal(t, people[0].First, person.First)
		assert.Equal(t, people[0].Gender, person.Gender)
		assert.Equal(t, people[0].Last, person.Last)
		assert.Equal(t, id, person.Identifier)
	}
	// test the update function
	temp := people[0]
	temp.Age = 4
	temp.AgeGroup = "Youngling"
	temp.First = "Update!"
	temp.Last = "Test"
	temp.Distance = "12 Mile Fun"
	temp.Gender = "U"
	person, err = db.AddPerson(eventYear.Identifier, temp)
	if assert.NoError(t, err) {
		assert.Equal(t, temp.Age, person.Age)
		assert.Equal(t, temp.AgeGroup, person.AgeGroup)
		assert.Equal(t, temp.Bib, person.Bib)
		assert.Equal(t, temp.Distance, person.Distance)
		assert.Equal(t, temp.First, person.First)
		assert.Equal(t, temp.Gender, person.Gender)
		assert.Equal(t, temp.Last, person.Last)
		assert.Equal(t, id, person.Identifier)
	}
	temp = people[1]
	temp.Age = 4
	temp.AgeGroup = "Youngling"
	temp.First = "Update!"
	temp.Last = "Test"
	temp.Distance = "12 Mile Fun"
	temp.Gender = "NB"
	person, err = db.AddPerson(eventYear.Identifier, temp)
	if assert.NoError(t, err) && assert.NotNil(t, person) {
		assert.Equal(t, temp.Age, person.Age)
		assert.Equal(t, temp.AgeGroup, person.AgeGroup)
		assert.Equal(t, temp.Bib, person.Bib)
		assert.Equal(t, temp.Distance, person.Distance)
		assert.Equal(t, temp.First, person.First)
		assert.Equal(t, temp.Gender, person.Gender)
		assert.Equal(t, temp.Last, person.Last)
	}
}

func TestAddPeople(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up tests. %v", err)
	}
	defer finalize(t)
	setupPeople()
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
	p, err := db.GetPeople(event.Slug, eventYear.Year)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(p))
	}
	p, err = db.AddPeople(eventYear.Identifier, people)
	if assert.NoError(t, err) {
		assert.Equal(t, len(people), len(p))
		for _, outer := range people {
			found := false
			for _, inner := range p {
				if outer.Bib == inner.Bib {
					assert.Equal(t, outer.Age, inner.Age)
					assert.Equal(t, outer.AgeGroup, inner.AgeGroup)
					assert.Equal(t, outer.Bib, inner.Bib)
					assert.Equal(t, outer.Distance, inner.Distance)
					assert.Equal(t, outer.First, inner.First)
					assert.Equal(t, outer.Gender, inner.Gender)
					assert.Equal(t, outer.Last, inner.Last)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
	p, err = db.GetPeople(event.Slug, eventYear.Year)
	if assert.NoError(t, err) {
		assert.Equal(t, len(people), len(p))
	}
	// test update
	upd := make([]types.Person, 0)
	for _, temp := range people {
		upd = append(upd, types.Person{
			Bib:      temp.Bib,
			Age:      4,
			AgeGroup: "Youngling",
			First:    "Update!",
			Last:     "Test",
			Distance: "12 Mile Fun",
			Gender:   "U",
		})
	}
	p, err = db.AddPeople(eventYear.Identifier, upd)
	if assert.NoError(t, err) {
		assert.Equal(t, len(people), len(p))
		for _, outer := range people {
			found := false
			for _, inner := range p {
				if outer.Bib == inner.Bib {
					assert.Equal(t, 4, inner.Age)
					assert.Equal(t, "Youngling", inner.AgeGroup)
					assert.Equal(t, outer.Bib, inner.Bib)
					assert.Equal(t, "12 Mile Fun", inner.Distance)
					assert.Equal(t, "Update!", inner.First)
					assert.Equal(t, "U", inner.Gender)
					assert.Equal(t, "Test", inner.Last)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
}

func TestDeletePeople(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up tests. %v", err)
	}
	defer finalize(t)
	setupPeople()
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
	p, err := db.GetPeople(event.Slug, eventYear.Year)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(p))
	}
	_, err = db.AddPeople(eventYear.Identifier, people)
	assert.NoError(t, err)
	p, err = db.GetPeople(event.Slug, eventYear.Year)
	if assert.NoError(t, err) {
		assert.Equal(t, len(people), len(p))
	}
	err = db.DeletePeople(eventYear.Identifier, nil)
	assert.NoError(t, err)
	p, err = db.GetPeople(event.Slug, eventYear.Year)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(p))
	}
	_, err = db.AddPeople(eventYear.Identifier, people)
	assert.NoError(t, err)
	p, err = db.GetPeople(event.Slug, eventYear.Year)
	if assert.NoError(t, err) {
		assert.Equal(t, len(people), len(p))
	}
	toDelete := []string{
		people[0].Bib,
		people[1].Bib,
	}
	err = db.DeletePeople(eventYear.Identifier, toDelete)
	assert.NoError(t, err)
	p, err = db.GetPeople(event.Slug, eventYear.Year)
	if assert.NoError(t, err) {
		assert.Equal(t, len(people)-2, len(p))
		for _, outer := range people[2:] {
			found := false
			for _, inner := range p {
				if outer.Bib == inner.Bib {
					assert.Equal(t, outer.Age, inner.Age)
					assert.Equal(t, outer.AgeGroup, inner.AgeGroup)
					assert.Equal(t, outer.Bib, inner.Bib)
					assert.Equal(t, outer.Distance, inner.Distance)
					assert.Equal(t, outer.First, inner.First)
					assert.Equal(t, outer.Gender, inner.Gender)
					assert.Equal(t, outer.Last, inner.Last)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
}

func TestBadDatabasePerson(t *testing.T) {
	db := badTestSetup(t)
	_, err := db.GetPerson("", "", "")
	if err == nil {
		t.Fatalf("Expected error getting person")
	}
	_, err = db.GetPeople("", "")
	if err == nil {
		t.Fatalf("Expected error getting people")
	}
	_, err = db.AddPerson(0, types.Person{})
	if err == nil {
		t.Fatalf("Expected error adding person")
	}
	_, err = db.AddPeople(0, nil)
	if err == nil {
		t.Fatalf("Expected error adding people")
	}
	err = db.DeletePeople(0, nil)
	if err == nil {
		t.Fatalf("Expected error deleting people")
	}
}

func TestNoDatabasePerson(t *testing.T) {
	db := MySQL{}
	_, err := db.GetPerson("", "", "")
	if err == nil {
		t.Fatalf("Expected error getting person")
	}
	_, err = db.GetPeople("", "")
	if err == nil {
		t.Fatalf("Expected error getting people")
	}
	_, err = db.AddPerson(0, types.Person{})
	if err == nil {
		t.Fatalf("Expected error adding person")
	}
	_, err = db.AddPeople(0, nil)
	if err == nil {
		t.Fatalf("Expected error adding people")
	}
	err = db.DeletePeople(0, nil)
	if err == nil {
		t.Fatalf("Expected error deleting people")
	}
}
