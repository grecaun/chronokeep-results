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

func setupPeopleTests() {
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
			AlternateId: "1001",
			Bib:         "100",
			First:       "John",
			Last:        "Smith",
			Age:         24,
			Gender:      "M",
			AgeGroup:    "20-29",
			Distance:    "1 Mile",
			Anonymous:   true,
		},
		{
			AlternateId: "1061",
			Bib:         "106",
			First:       "Rose",
			Last:        "Johnson",
			Age:         55,
			Gender:      "F",
			AgeGroup:    "50-59",
			Distance:    "1 Mile",
			Anonymous:   false,
		},
		{
			AlternateId: "10",
			Bib:         "209",
			First:       "Tony",
			Last:        "Starke",
			Age:         45,
			Gender:      "M",
			AgeGroup:    "40-49",
			Distance:    "1 Mile",
		},
		{
			AlternateId: "285",
			Bib:         "287",
			First:       "Jamie",
			Last:        "Fischer",
			Age:         35,
			Gender:      "NB",
			AgeGroup:    "30-39",
			Distance:    "5 Mile",
		},
		{
			AlternateId: "132",
			Bib:         "2871",
			First:       "Jamie",
			Last:        "Fischer",
			Age:         35,
			Gender:      "F",
			AgeGroup:    "30-39",
			Distance:    "5 Mile",
		},
	}
}

func TestGetPerson(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up tests. %v", err)
	}
	defer finalize(t)
	setupPeopleTests()
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
		Live:            false,
		DaysAllowed:     1,
		RankingType:     "chip",
	}
	eventYear, _ = db.AddEventYear(*eventYear)
	person, err := db.GetPerson(event.Slug, eventYear.Year, people[0].Bib)
	if err != nil {
		t.Errorf("Error finding non existent person: %v", err)
	}
	if person != nil {
		t.Errorf("Found someone when no one should exist: %v", person)
	}
	db.AddPeople(eventYear.Identifier, people)
	for _, p := range people {
		person, err = db.GetPerson(event.Slug, eventYear.Year, p.Bib)
		if err != nil {
			t.Errorf("Error finding non existent person: %v", err)
		}
		assert.True(t, person.Equals(&p))
		if person == nil {
			t.Errorf("Expected to find someone but did not.")
		} else if person.Bib != p.Bib {
			t.Errorf("Bib doesn't match: %v+ %v+", *person, p)
		} else if person.First != p.First {
			t.Errorf("First name doesn't match: %v+ %v+", *person, p)
		} else if person.Last != p.Last {
			t.Errorf("Last name doesn't match: %v+ %v+", *person, p)
		} else if person.Age != p.Age {
			t.Errorf("Age doesn't match: %v+ %v+", *person, p)
		} else if person.AgeGroup != p.AgeGroup {
			t.Errorf("Age Group doesn't match: %v+ %v+", *person, p)
		} else if person.Gender != p.Gender {
			t.Errorf("Gender doesn't match: %v+ %v+", *person, p)
		} else if person.Distance != p.Distance {
			t.Errorf("Distance doesn't match: %v+ %v+", *person, p)
		} else if person.Anonymous != p.Anonymous {
			t.Errorf("Anonymous doesn't match: %v+ %v+", *person, p)
		} else if person.AlternateId != p.AlternateId {
			t.Errorf("AlternateId doesn't match: %v+ %v+", *person, p)
		}
	}
}

func TestGetPeople(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up tests. %v", err)
	}
	defer finalize(t)
	setupPeopleTests()
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
		Live:            false,
		DaysAllowed:     1,
		RankingType:     "chip",
	}
	eventYear, _ = db.AddEventYear(*eventYear)
	iPeople, err := db.GetPeople(event.Slug, eventYear.Year)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(iPeople))
	}
	db.AddPeople(eventYear.Identifier, people)
	iPeople, err = db.GetPeople(event.Slug, eventYear.Year)
	if assert.NoError(t, err) {
		assert.Equal(t, len(people), len(iPeople))
		for _, outer := range people {
			found := false
			for _, inner := range iPeople {
				if outer.AlternateId == inner.AlternateId {
					assert.True(t, outer.Equals(&inner))
					assert.Equal(t, outer.Age, inner.Age)
					assert.Equal(t, outer.AgeGroup, inner.AgeGroup)
					assert.Equal(t, outer.Bib, inner.Bib)
					assert.Equal(t, outer.Distance, inner.Distance)
					assert.Equal(t, outer.First, inner.First)
					assert.Equal(t, outer.Gender, inner.Gender)
					assert.Equal(t, outer.Last, inner.Last)
					assert.Equal(t, outer.Anonymous, inner.Anonymous)
					assert.Equal(t, outer.AlternateId, inner.AlternateId)
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
	setupPeopleTests()
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
		Live:            false,
		DaysAllowed:     1,
		RankingType:     "chip",
	}
	eventYear, _ = db.AddEventYear(*eventYear)
	person, err := db.AddPerson(eventYear.Identifier, people[0])
	if assert.NoError(t, err) {
		assert.True(t, people[0].Equals(person))
		assert.Equal(t, people[0].Age, person.Age)
		assert.Equal(t, people[0].AgeGroup, person.AgeGroup)
		assert.Equal(t, people[0].Bib, person.Bib)
		assert.Equal(t, people[0].Distance, person.Distance)
		assert.Equal(t, people[0].First, person.First)
		assert.Equal(t, people[0].Gender, person.Gender)
		assert.Equal(t, people[0].Last, person.Last)
		assert.Equal(t, people[0].Anonymous, person.Anonymous)
		assert.Equal(t, people[0].AlternateId, person.AlternateId)
	}
	id := person.Identifier
	person, err = db.GetPerson(event.Slug, eventYear.Year, people[0].Bib)
	if assert.NoError(t, err) {
		assert.True(t, people[0].Equals(person))
		assert.Equal(t, people[0].Age, person.Age)
		assert.Equal(t, people[0].AgeGroup, person.AgeGroup)
		assert.Equal(t, people[0].Bib, person.Bib)
		assert.Equal(t, people[0].Distance, person.Distance)
		assert.Equal(t, people[0].First, person.First)
		assert.Equal(t, people[0].Gender, person.Gender)
		assert.Equal(t, people[0].Last, person.Last)
		assert.Equal(t, id, person.Identifier)
		assert.Equal(t, people[0].Anonymous, person.Anonymous)
		assert.Equal(t, people[0].AlternateId, person.AlternateId)
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
		assert.True(t, temp.Equals(person))
		assert.Equal(t, temp.Age, person.Age)
		assert.Equal(t, temp.AgeGroup, person.AgeGroup)
		assert.Equal(t, temp.Bib, person.Bib)
		assert.Equal(t, temp.Distance, person.Distance)
		assert.Equal(t, temp.First, person.First)
		assert.Equal(t, temp.Gender, person.Gender)
		assert.Equal(t, temp.Last, person.Last)
		assert.Equal(t, id, person.Identifier)
		assert.Equal(t, temp.Anonymous, person.Anonymous)
		assert.Equal(t, temp.AlternateId, person.AlternateId)
	}
	temp = people[1]
	temp.Age = 4
	temp.AgeGroup = "Youngling"
	temp.First = "Update!"
	temp.Last = "Test"
	temp.Distance = "12 Mile Fun"
	temp.Gender = "NB"
	temp.Anonymous = true
	person, err = db.AddPerson(eventYear.Identifier, temp)
	if assert.NoError(t, err) && assert.NotNil(t, person) {
		assert.True(t, temp.Equals(person))
		assert.Equal(t, temp.Age, person.Age)
		assert.Equal(t, temp.AgeGroup, person.AgeGroup)
		assert.Equal(t, temp.Bib, person.Bib)
		assert.Equal(t, temp.Distance, person.Distance)
		assert.Equal(t, temp.First, person.First)
		assert.Equal(t, temp.Gender, person.Gender)
		assert.Equal(t, temp.Last, person.Last)
		assert.Equal(t, temp.Anonymous, person.Anonymous)
		assert.Equal(t, temp.AlternateId, person.AlternateId)
	}
}

func TestAddPeople(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up tests. %v", err)
	}
	defer finalize(t)
	setupPeopleTests()
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
		Live:            false,
		DaysAllowed:     1,
		RankingType:     "chip",
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
				if outer.AlternateId == inner.AlternateId {
					assert.True(t, outer.Equals(&inner))
					assert.Equal(t, outer.Age, inner.Age)
					assert.Equal(t, outer.AgeGroup, inner.AgeGroup)
					assert.Equal(t, outer.Bib, inner.Bib)
					assert.Equal(t, outer.Distance, inner.Distance)
					assert.Equal(t, outer.First, inner.First)
					assert.Equal(t, outer.Gender, inner.Gender)
					assert.Equal(t, outer.Last, inner.Last)
					assert.Equal(t, outer.Anonymous, inner.Anonymous)
					assert.Equal(t, outer.AlternateId, inner.AlternateId)
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
			AlternateId: temp.AlternateId,
			Bib:         temp.Bib,
			Age:         4,
			AgeGroup:    "Youngling",
			First:       "Update!",
			Last:        "Test",
			Distance:    "12 Mile Fun",
			Gender:      "U",
			Anonymous:   true,
		})
	}
	p, err = db.AddPeople(eventYear.Identifier, upd)
	if assert.NoError(t, err) {
		assert.Equal(t, len(people), len(p))
		for _, outer := range people {
			found := false
			for _, inner := range p {
				if outer.AlternateId == inner.AlternateId {
					assert.Equal(t, 4, inner.Age)
					assert.Equal(t, "Youngling", inner.AgeGroup)
					assert.Equal(t, outer.Bib, inner.Bib)
					assert.Equal(t, "12 Mile Fun", inner.Distance)
					assert.Equal(t, "Update!", inner.First)
					assert.Equal(t, "U", inner.Gender)
					assert.Equal(t, "Test", inner.Last)
					assert.Equal(t, true, inner.Anonymous)
					assert.Equal(t, outer.AlternateId, inner.AlternateId)
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
	setupPeopleTests()
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
		Live:            false,
		DaysAllowed:     1,
		RankingType:     "chip",
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
	count, err := db.DeletePeople(eventYear.Identifier, nil)
	assert.NoError(t, err)
	p, err = db.GetPeople(event.Slug, eventYear.Year)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(p))
	}
	assert.Equal(t, count, int64(len(people)))
	_, err = db.AddPeople(eventYear.Identifier, people)
	assert.NoError(t, err)
	p, err = db.GetPeople(event.Slug, eventYear.Year)
	if assert.NoError(t, err) {
		assert.Equal(t, len(people), len(p))
	}
	toDelete := []string{
		people[0].AlternateId,
		people[1].AlternateId,
	}
	count, err = db.DeletePeople(eventYear.Identifier, toDelete)
	assert.NoError(t, err)
	assert.Equal(t, count, int64(len(toDelete)))
	p, err = db.GetPeople(event.Slug, eventYear.Year)
	if assert.NoError(t, err) {
		assert.Equal(t, len(people)-2, len(p))
		for _, outer := range people[2:] {
			found := false
			for _, inner := range p {
				if outer.AlternateId == inner.AlternateId {
					assert.Equal(t, outer.Age, inner.Age)
					assert.Equal(t, outer.AgeGroup, inner.AgeGroup)
					assert.Equal(t, outer.Bib, inner.Bib)
					assert.Equal(t, outer.Distance, inner.Distance)
					assert.Equal(t, outer.First, inner.First)
					assert.Equal(t, outer.Gender, inner.Gender)
					assert.Equal(t, outer.Last, inner.Last)
					assert.Equal(t, outer.Anonymous, inner.Anonymous)
					assert.Equal(t, outer.AlternateId, inner.AlternateId)
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
	_, err = db.DeletePeople(0, nil)
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
	_, err = db.DeletePeople(0, nil)
	if err == nil {
		t.Fatalf("Expected error deleting people")
	}
}

func TestUpdatePerson(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up tests. %v", err)
	}
	defer finalize(t)
	setupPeopleTests()
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
		Live:            false,
		DaysAllowed:     1,
		RankingType:     "chip",
	}
	eventYear, _ = db.AddEventYear(*eventYear)
	person, _ := db.AddPerson(eventYear.Identifier, people[0])
	assert.NotNil(t, person)
	id := person.Identifier
	// test update
	temp := people[0]
	temp.Bib = "-200"
	temp.Age = 4
	temp.AgeGroup = "Youngling"
	temp.First = "Update!"
	temp.Last = "Test"
	temp.Distance = "12 Mile Fun"
	temp.Gender = "U"
	temp.Anonymous = !people[0].Anonymous
	person, err = db.UpdatePerson(eventYear.Identifier, temp)
	if assert.NoError(t, err) {
		assert.True(t, temp.Equals(person))
		assert.Equal(t, temp.Age, person.Age)
		assert.Equal(t, temp.AgeGroup, person.AgeGroup)
		assert.Equal(t, temp.Bib, person.Bib)
		assert.Equal(t, temp.Distance, person.Distance)
		assert.Equal(t, temp.First, person.First)
		assert.Equal(t, temp.Gender, person.Gender)
		assert.Equal(t, temp.Last, person.Last)
		assert.Equal(t, id, person.Identifier)
		assert.Equal(t, temp.Anonymous, person.Anonymous)
		assert.Equal(t, temp.AlternateId, person.AlternateId)
	}
	// test invalid update
	temp = people[1]
	temp.Bib = "newbib"
	person, err = db.UpdatePerson(eventYear.Identifier, temp)
	assert.Error(t, err)
	assert.Nil(t, person)
}
