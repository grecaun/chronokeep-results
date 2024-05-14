package mysql

import (
	"chronokeep/results/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	segments []types.Segment
)

func setupSegmentTests() {
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
	segments = []types.Segment{
		{
			Location:      "Finish",
			DistanceName:  "10k",
			Name:          "5k",
			DistanceValue: 5.0,
			DistanceUnit:  "kilometer",
			GPS:           "rando",
			MapLink:       "some.link",
		},
		{
			Location:      "Aid1",
			DistanceName:  "Half Marathon",
			Name:          "Mile 6.2",
			DistanceValue: 6.2,
			DistanceUnit:  "mile",
			GPS:           "",
			MapLink:       "",
		},
		{
			Location:      "Aid2",
			DistanceName:  "Half Marathon",
			Name:          "Mile 8",
			DistanceValue: 8.0,
			DistanceUnit:  "mile",
			GPS:           "",
			MapLink:       "another.link",
		},
		{
			Location:      "Aid2",
			DistanceName:  "Marathon",
			Name:          "Mile 8",
			DistanceValue: 8.0,
			DistanceUnit:  "mile",
			GPS:           "",
			MapLink:       "another2.link",
		},
		{
			Location:      "Aid1",
			DistanceName:  "Marathon",
			Name:          "Mile 6.2",
			DistanceValue: 6.2,
			DistanceUnit:  "mile",
			GPS:           "",
			MapLink:       "",
		},
	}
}

func TestAddSegments(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up tests. %v", err)
	}
	defer finalize(t)
	setupSegmentTests()
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
	s, err := db.GetSegments(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(s))
	}
	s, err = db.AddSegments(eventYear.Identifier, segments)
	if assert.NoError(t, err) {
		assert.Equal(t, len(segments), len(s))
		for _, outer := range segments {
			found := false
			for _, inner := range s {
				if outer.Name == inner.Name && outer.DistanceName == inner.DistanceName {
					assert.Equal(t, outer.Location, inner.Location)
					assert.Equal(t, outer.DistanceName, inner.DistanceName)
					assert.Equal(t, outer.Name, inner.Name)
					assert.Equal(t, outer.DistanceValue, inner.DistanceValue)
					assert.Equal(t, outer.DistanceUnit, inner.DistanceUnit)
					assert.Equal(t, outer.GPS, inner.GPS)
					assert.Equal(t, outer.MapLink, inner.MapLink)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
	s, err = db.GetSegments(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, len(segments), len(s))
		for _, outer := range segments {
			found := false
			for _, inner := range s {
				if outer.Name == inner.Name && outer.DistanceName == inner.DistanceName {
					assert.Equal(t, outer.Location, inner.Location)
					assert.Equal(t, outer.DistanceName, inner.DistanceName)
					assert.Equal(t, outer.Name, inner.Name)
					assert.Equal(t, outer.DistanceValue, inner.DistanceValue)
					assert.Equal(t, outer.DistanceUnit, inner.DistanceUnit)
					assert.Equal(t, outer.GPS, inner.GPS)
					assert.Equal(t, outer.MapLink, inner.MapLink)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
	// test update
	upd := make([]types.Segment, 0)
	for _, temp := range segments {
		upd = append(upd, types.Segment{
			DistanceName:  temp.DistanceName,
			Name:          temp.Name,
			Location:      "newLoc",
			DistanceValue: 10.0,
			DistanceUnit:  "newUnit",
			GPS:           "someGPS",
			MapLink:       "new.link",
		})
	}
	s, err = db.AddSegments(eventYear.Identifier, upd)
	if assert.NoError(t, err) {
		assert.Equal(t, len(upd), len(s))
		for _, outer := range upd {
			found := false
			for _, inner := range s {
				if outer.Name == inner.Name && outer.DistanceName == inner.DistanceName {
					assert.Equal(t, outer.Location, inner.Location)
					assert.Equal(t, outer.DistanceName, inner.DistanceName)
					assert.Equal(t, outer.Name, inner.Name)
					assert.Equal(t, outer.DistanceValue, inner.DistanceValue)
					assert.Equal(t, outer.DistanceUnit, inner.DistanceUnit)
					assert.Equal(t, outer.GPS, inner.GPS)
					assert.Equal(t, outer.MapLink, inner.MapLink)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
	s, err = db.GetSegments(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, len(upd), len(s))
		for _, outer := range upd {
			found := false
			for _, inner := range s {
				if outer.Name == inner.Name && outer.DistanceName == inner.DistanceName {
					assert.Equal(t, outer.Location, inner.Location)
					assert.Equal(t, outer.DistanceName, inner.DistanceName)
					assert.Equal(t, outer.Name, inner.Name)
					assert.Equal(t, outer.DistanceValue, inner.DistanceValue)
					assert.Equal(t, outer.DistanceUnit, inner.DistanceUnit)
					assert.Equal(t, outer.GPS, inner.GPS)
					assert.Equal(t, outer.MapLink, inner.MapLink)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
}

func TestGetSegments(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up tests. %v", err)
	}
	defer finalize(t)
	setupSegmentTests()
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
	s, err := db.GetSegments(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(s))
	}
	s, err = db.AddSegments(eventYear.Identifier, segments)
	if assert.NoError(t, err) {
		assert.Equal(t, len(segments), len(s))
	}
	s, err = db.GetSegments(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, len(segments), len(s))
		for _, outer := range segments {
			found := false
			for _, inner := range s {
				if outer.Name == inner.Name && outer.DistanceName == inner.DistanceName {
					assert.Equal(t, outer.Location, inner.Location)
					assert.Equal(t, outer.DistanceName, inner.DistanceName)
					assert.Equal(t, outer.Name, inner.Name)
					assert.Equal(t, outer.DistanceValue, inner.DistanceValue)
					assert.Equal(t, outer.DistanceUnit, inner.DistanceUnit)
					assert.Equal(t, outer.GPS, inner.GPS)
					assert.Equal(t, outer.MapLink, inner.MapLink)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
}

func TestDeleteSegments(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up tests. %v", err)
	}
	defer finalize(t)
	setupSegmentTests()
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
	s, err := db.GetSegments(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(s))
	}
	s, err = db.AddSegments(eventYear.Identifier, segments)
	if assert.NoError(t, err) {
		assert.Equal(t, len(segments), len(s))
	}
	s, err = db.GetSegments(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, len(segments), len(s))
	}
	count, err := db.DeleteSegments(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, int64(len(segments)), count)
	}
	s, err = db.GetSegments(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(s))
	}
}
