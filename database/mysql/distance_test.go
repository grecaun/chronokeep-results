package mysql

import (
	"chronokeep/results/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	distances []types.Distance
)

func setupDistanceTests() {
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
	distances = []types.Distance{
		{
			Name:          "Marathon",
			Certification: "USATF Certification #WA555221",
		},
		{
			Name:          "Half Marathon",
			Certification: "USATF Certification #WA555222",
		},
		{
			Name:          "10k",
			Certification: "USATF Certification #WA555121",
		},
	}
}

func TestAddDistance(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up tests. %v", err)
	}
	defer finalize(t)
	setupDistanceTests()
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
	d, err := db.GetDistances(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(d))
	}
	d, err = db.AddDistances(eventYear.Identifier, distances)
	if assert.NoError(t, err) {
		assert.Equal(t, len(distances), len(d))
		for _, outer := range distances {
			found := false
			for _, inner := range d {
				if outer.Name == inner.Name {
					assert.Equal(t, outer.Certification, inner.Certification)
					assert.Equal(t, outer.Name, inner.Name)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
	d, err = db.GetDistances(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, len(distances), len(d))
		for _, outer := range distances {
			found := false
			for _, inner := range d {
				if outer.Name == inner.Name {
					assert.Equal(t, outer.Certification, inner.Certification)
					assert.Equal(t, outer.Name, inner.Name)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
	// test update feature
	upd := make([]types.Distance, 0)
	for _, temp := range distances {
		upd = append(upd, types.Distance{
			Name:          temp.Name,
			Certification: "New Certification #28731" + temp.Name,
		})
	}
	d, err = db.AddDistances(eventYear.Identifier, upd)
	if assert.NoError(t, err) {
		assert.Equal(t, len(upd), len(d))
		for _, outer := range upd {
			found := false
			for _, inner := range d {
				if outer.Name == inner.Name {
					assert.Equal(t, outer.Certification, inner.Certification)
					assert.Equal(t, outer.Name, inner.Name)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
	d, err = db.GetDistances(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, len(upd), len(d))
		for _, outer := range upd {
			found := false
			for _, inner := range d {
				if outer.Name == inner.Name {
					assert.Equal(t, outer.Certification, inner.Certification)
					assert.Equal(t, outer.Name, inner.Name)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
}

func TestGetDistance(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up tests. %v", err)
	}
	defer finalize(t)
	setupDistanceTests()
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
	d, err := db.GetDistance(eventYear.Identifier, distances[0].Name)
	if assert.NoError(t, err) {
		assert.Nil(t, d)
	}
	dists, err := db.AddDistances(eventYear.Identifier, distances)
	if assert.NoError(t, err) {
		assert.Equal(t, len(distances), len(dists))
	}
	d, err = db.GetDistance(eventYear.Identifier, distances[0].Name)
	if assert.NoError(t, err) {
		assert.Equal(t, distances[0].Name, d.Name)
		assert.Equal(t, distances[0].Certification, d.Certification)
	}
}

func TestGetDistances(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up tests. %v", err)
	}
	defer finalize(t)
	setupDistanceTests()
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
	d, err := db.GetDistances(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(d))
	}
	d, err = db.AddDistances(eventYear.Identifier, distances)
	if assert.NoError(t, err) {
		assert.Equal(t, len(distances), len(d))
	}
	d, err = db.GetDistances(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, len(distances), len(d))
		for _, outer := range distances {
			found := false
			for _, inner := range d {
				if outer.Name == inner.Name {
					assert.Equal(t, outer.Certification, inner.Certification)
					assert.Equal(t, outer.Name, inner.Name)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
}

func TestDeleteDistances(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up tests. %v", err)
	}
	defer finalize(t)
	setupDistanceTests()
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
	d, err := db.GetDistances(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(d))
	}
	d, err = db.AddDistances(eventYear.Identifier, distances)
	if assert.NoError(t, err) {
		assert.Equal(t, len(distances), len(d))
	}
	d, err = db.GetDistances(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, len(distances), len(d))
	}
	count, err := db.DeleteDistances(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, int64(len(distances)), count)
	}
	d, err = db.GetDistances(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(d))
	}
}
