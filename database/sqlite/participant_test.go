package sqlite

import (
	"chronokeep/results/types"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	participants     []types.Participant
	manyParticipants []types.Participant
)

func setupParticipantTests() {
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
	participants = []types.Participant{
		{
			AlternateId: "100",
			Bib:         "100",
			First:       "John",
			Last:        "Smith",
			Birthdate:   "1/1/2000",
			Gender:      "M",
			AgeGroup:    "20-29",
			Distance:    "1 Mile",
			Anonymous:   true,
			SMSEnabled:  false,
			Mobile:      "3555521234",
			Apparel:     "Medium",
			UpdatedAt:   0,
		},
		{
			AlternateId: "106",
			Bib:         "106",
			First:       "Rose",
			Last:        "Johnson",
			Birthdate:   "1/1/1975",
			Gender:      "F",
			AgeGroup:    "50-59",
			Distance:    "1 Mile",
			Anonymous:   false,
			SMSEnabled:  true,
			Mobile:      "3543421234",
			Apparel:     "Small",
			UpdatedAt:   0,
		},
		{
			AlternateId: "10",
			Bib:         "209",
			First:       "Tony",
			Last:        "Starke",
			Birthdate:   "1/1/1983",
			Gender:      "M",
			AgeGroup:    "40-49",
			Distance:    "1 Mile",
			Anonymous:   false,
			SMSEnabled:  false,
			Mobile:      "3525521234",
			Apparel:     "",
			UpdatedAt:   0,
		},
		{
			AlternateId: "285",
			Bib:         "287",
			First:       "Jamie",
			Last:        "Fischer",
			Birthdate:   "1/1/1993",
			Gender:      "NB",
			AgeGroup:    "30-39",
			Distance:    "5 Mile",
			Anonymous:   false,
			SMSEnabled:  false,
			Mobile:      "",
			Apparel:     "X-Large",
			UpdatedAt:   0,
		},
		{
			AlternateId: "132",
			Bib:         "2871",
			First:       "Jamie",
			Last:        "Fischer",
			Birthdate:   "1/1/1993",
			Gender:      "F",
			AgeGroup:    "30-39",
			Distance:    "5 Mile",
			Anonymous:   false,
			SMSEnabled:  true,
			Mobile:      "3215521234",
			Apparel:     "X-Small",
			UpdatedAt:   100,
		},
	}
	manyParticipants = make([]types.Participant, 0)
	for i := 0; i < 220; i++ {
		manyParticipants = append(manyParticipants, types.Participant{
			AlternateId: fmt.Sprintf("%d", i),
			Bib:         fmt.Sprintf("%d", i),
			First:       fmt.Sprintf("%s-%d", "John", i),
			Last:        fmt.Sprintf("%s-%d", "Smith", i),
			Birthdate:   "1/1/2000",
			Gender:      "M",
			AgeGroup:    "20-29",
			Distance:    "1 Mile",
			Anonymous:   true,
			SMSEnabled:  false,
			Mobile:      "3555521234",
			Apparel:     "Medium",
			UpdatedAt:   int64(i * 10),
		})
	}
}

func TestAddParticipants(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up tests. %v", err)
	}
	defer finalize(t)
	setupParticipantTests()
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
	p, err := db.GetParticipants(eventYear.Identifier, 0, 0, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(p))
	}
	p, err = db.AddParticipants(eventYear.Identifier, participants)
	if assert.NoError(t, err) {
		assert.Equal(t, len(participants), len(p))
		for _, outer := range participants {
			found := false
			for _, inner := range p {
				if outer.AlternateId == inner.AlternateId {
					assert.True(t, outer.Equals(&inner))
					assert.Equal(t, outer.Birthdate, inner.Birthdate)
					assert.Equal(t, outer.AgeGroup, inner.AgeGroup)
					assert.Equal(t, outer.Bib, inner.Bib)
					assert.Equal(t, outer.Distance, inner.Distance)
					assert.Equal(t, outer.First, inner.First)
					assert.Equal(t, outer.Gender, inner.Gender)
					assert.Equal(t, outer.Last, inner.Last)
					assert.Equal(t, outer.Anonymous, inner.Anonymous)
					assert.Equal(t, outer.AlternateId, inner.AlternateId)
					assert.Equal(t, outer.SMSEnabled, inner.SMSEnabled)
					assert.Equal(t, outer.Mobile, inner.Mobile)
					assert.Equal(t, outer.Apparel, inner.Apparel)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
	p, err = db.GetParticipants(eventYear.Identifier, 0, 0, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, len(participants), len(p))
		for _, outer := range participants {
			found := false
			for _, inner := range p {
				if outer.AlternateId == inner.AlternateId {
					assert.True(t, outer.Equals(&inner))
					assert.Equal(t, outer.Birthdate, inner.Birthdate)
					assert.Equal(t, outer.AgeGroup, inner.AgeGroup)
					assert.Equal(t, outer.Bib, inner.Bib)
					assert.Equal(t, outer.Distance, inner.Distance)
					assert.Equal(t, outer.First, inner.First)
					assert.Equal(t, outer.Gender, inner.Gender)
					assert.Equal(t, outer.Last, inner.Last)
					assert.Equal(t, outer.Anonymous, inner.Anonymous)
					assert.Equal(t, outer.AlternateId, inner.AlternateId)
					assert.Equal(t, outer.SMSEnabled, inner.SMSEnabled)
					assert.Equal(t, outer.Mobile, inner.Mobile)
					assert.Equal(t, outer.Apparel, inner.Apparel)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
	// test update
	upd := make([]types.Participant, 0)
	for _, temp := range participants {
		upd = append(upd, types.Participant{
			AlternateId: temp.AlternateId,
			Bib:         temp.Bib,
			Birthdate:   "2/2/1975",
			AgeGroup:    "newgroup",
			First:       "Update!",
			Last:        "test",
			Distance:    "12 Mile Fun",
			Gender:      "U",
			Anonymous:   true,
			SMSEnabled:  true,
			Mobile:      "empty",
			Apparel:     "moreempty",
		})
	}
	p, err = db.AddParticipants(eventYear.Identifier, upd)
	if assert.NoError(t, err) {
		assert.Equal(t, len(upd), len(p))
		for _, outer := range upd {
			found := false
			for _, inner := range p {
				if outer.AlternateId == inner.AlternateId {
					assert.True(t, outer.Equals(&inner))
					assert.Equal(t, outer.Birthdate, inner.Birthdate)
					assert.Equal(t, outer.AgeGroup, inner.AgeGroup)
					assert.Equal(t, outer.Bib, inner.Bib)
					assert.Equal(t, outer.Distance, inner.Distance)
					assert.Equal(t, outer.First, inner.First)
					assert.Equal(t, outer.Gender, inner.Gender)
					assert.Equal(t, outer.Last, inner.Last)
					assert.Equal(t, outer.Anonymous, inner.Anonymous)
					assert.Equal(t, outer.AlternateId, inner.AlternateId)
					assert.Equal(t, outer.SMSEnabled, inner.SMSEnabled)
					assert.Equal(t, outer.Mobile, inner.Mobile)
					assert.Equal(t, outer.Apparel, inner.Apparel)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
	p, err = db.GetParticipants(eventYear.Identifier, 0, 0, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, len(upd), len(p))
		for _, outer := range upd {
			found := false
			for _, inner := range p {
				if outer.AlternateId == inner.AlternateId {
					assert.True(t, outer.Equals(&inner))
					assert.Equal(t, outer.Birthdate, inner.Birthdate)
					assert.Equal(t, outer.AgeGroup, inner.AgeGroup)
					assert.Equal(t, outer.Bib, inner.Bib)
					assert.Equal(t, outer.Distance, inner.Distance)
					assert.Equal(t, outer.First, inner.First)
					assert.Equal(t, outer.Gender, inner.Gender)
					assert.Equal(t, outer.Last, inner.Last)
					assert.Equal(t, outer.Anonymous, inner.Anonymous)
					assert.Equal(t, outer.AlternateId, inner.AlternateId)
					assert.Equal(t, outer.SMSEnabled, inner.SMSEnabled)
					assert.Equal(t, outer.Mobile, inner.Mobile)
					assert.Equal(t, outer.Apparel, inner.Apparel)
					found = true
				}
			}
			assert.True(t, found)
		}
	} // test update (bib should not update)
	upd = make([]types.Participant, 0)
	for _, temp := range participants {
		upd = append(upd, types.Participant{
			AlternateId: temp.AlternateId,
			Bib:         temp.Bib + "new",
			Birthdate:   "2/2/1975",
			AgeGroup:    "newgroup",
			First:       "Update!",
			Last:        "test",
			Distance:    "12 Mile Fun",
			Gender:      "U",
			Anonymous:   true,
			SMSEnabled:  true,
			Mobile:      "empty",
			Apparel:     "moreempty",
		})
	}
	p, err = db.AddParticipants(eventYear.Identifier, upd)
	if assert.NoError(t, err) {
		assert.Equal(t, len(upd), len(p))
		for _, outer := range upd {
			found := false
			for _, inner := range p {
				if outer.AlternateId == inner.AlternateId {
					assert.True(t, outer.Equals(&inner))
					assert.Equal(t, outer.Birthdate, inner.Birthdate)
					assert.Equal(t, outer.AgeGroup, inner.AgeGroup)
					assert.Equal(t, outer.Bib, inner.Bib)
					assert.Equal(t, outer.Distance, inner.Distance)
					assert.Equal(t, outer.First, inner.First)
					assert.Equal(t, outer.Gender, inner.Gender)
					assert.Equal(t, outer.Last, inner.Last)
					assert.Equal(t, outer.Anonymous, inner.Anonymous)
					assert.Equal(t, outer.AlternateId, inner.AlternateId)
					assert.Equal(t, outer.SMSEnabled, inner.SMSEnabled)
					assert.Equal(t, outer.Mobile, inner.Mobile)
					assert.Equal(t, outer.Apparel, inner.Apparel)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
	p, err = db.GetParticipants(eventYear.Identifier, 0, 0, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, len(upd), len(p))
		for _, outer := range upd {
			found := false
			for _, inner := range p {
				if outer.AlternateId == inner.AlternateId {
					assert.False(t, outer.Equals(&inner))
					assert.Equal(t, outer.Birthdate, inner.Birthdate)
					assert.Equal(t, outer.AgeGroup, inner.AgeGroup)
					assert.NotEqual(t, outer.Bib, inner.Bib)
					assert.Equal(t, outer.Distance, inner.Distance)
					assert.Equal(t, outer.First, inner.First)
					assert.Equal(t, outer.Gender, inner.Gender)
					assert.Equal(t, outer.Last, inner.Last)
					assert.Equal(t, outer.Anonymous, inner.Anonymous)
					assert.Equal(t, outer.AlternateId, inner.AlternateId)
					assert.Equal(t, outer.SMSEnabled, inner.SMSEnabled)
					assert.Equal(t, outer.Mobile, inner.Mobile)
					assert.Equal(t, outer.Apparel, inner.Apparel)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
}

func TestGetParticipants(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up tests. %v", err)
	}
	defer finalize(t)
	setupParticipantTests()
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
	iParts, err := db.GetParticipants(eventYear.Identifier, 50, 0, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(iParts))
	}
	db.AddParticipants(eventYear.Identifier, participants)
	iParts, err = db.GetParticipants(eventYear.Identifier, 50, 0, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, len(participants), len(iParts))
		for _, outer := range participants {
			found := false
			for _, inner := range iParts {
				if outer.AlternateId == inner.AlternateId {
					assert.True(t, outer.Equals(&inner))
					assert.Equal(t, outer.Birthdate, inner.Birthdate)
					assert.Equal(t, outer.AgeGroup, inner.AgeGroup)
					assert.Equal(t, outer.Bib, inner.Bib)
					assert.Equal(t, outer.Distance, inner.Distance)
					assert.Equal(t, outer.First, inner.First)
					assert.Equal(t, outer.Gender, inner.Gender)
					assert.Equal(t, outer.Last, inner.Last)
					assert.Equal(t, outer.Anonymous, inner.Anonymous)
					assert.Equal(t, outer.AlternateId, inner.AlternateId)
					assert.Equal(t, outer.SMSEnabled, inner.SMSEnabled)
					assert.Equal(t, outer.Mobile, inner.Mobile)
					assert.Equal(t, outer.Apparel, inner.Apparel)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
	updatedAt := int64(100)
	iParts, err = db.GetParticipants(eventYear.Identifier, 50, 0, &updatedAt)
	if assert.NoError(t, err) {
		if assert.Equal(t, 1, len(iParts)) {
			assert.True(t, participants[4].Equals(&iParts[0]))
			assert.True(t, iParts[0].UpdatedAt >= updatedAt)
		}
	}
	db.DeleteParticipants(eventYear.Identifier, nil)
	db.AddParticipants(eventYear.Identifier, manyParticipants)
	pulledParticipants := make([]types.Participant, 0)
	for i := 0; i < int(math.Ceil(float64(len(manyParticipants))/50)); i++ {
		iParts, err = db.GetParticipants(eventYear.Identifier, 50, i, nil)
		if assert.NoError(t, err) {
			assert.True(t, len(iParts) > 0 && len(iParts) <= 50)
			for _, outer := range iParts {
				found := false
				for _, inner := range manyParticipants {
					if outer.AlternateId == inner.AlternateId {
						assert.True(t, outer.Equals(&inner))
						assert.Equal(t, outer.Birthdate, inner.Birthdate)
						assert.Equal(t, outer.AgeGroup, inner.AgeGroup)
						assert.Equal(t, outer.Bib, inner.Bib)
						assert.Equal(t, outer.Distance, inner.Distance)
						assert.Equal(t, outer.First, inner.First)
						assert.Equal(t, outer.Gender, inner.Gender)
						assert.Equal(t, outer.Last, inner.Last)
						assert.Equal(t, outer.Anonymous, inner.Anonymous)
						assert.Equal(t, outer.AlternateId, inner.AlternateId)
						assert.Equal(t, outer.SMSEnabled, inner.SMSEnabled)
						assert.Equal(t, outer.Mobile, inner.Mobile)
						assert.Equal(t, outer.Apparel, inner.Apparel)
						found = true
					}
				}
				assert.True(t, found)
				found = false
				for _, inner := range pulledParticipants {
					if outer.AlternateId == inner.AlternateId {
						found = true
					}
				}
				assert.False(t, found)
			}
			pulledParticipants = append(pulledParticipants, iParts...)
		}
	}
	for _, outer := range manyParticipants {
		found := false
		for _, inner := range pulledParticipants {
			if outer.AlternateId == inner.AlternateId {
				assert.True(t, outer.Equals(&inner))
				assert.Equal(t, outer.Birthdate, inner.Birthdate)
				assert.Equal(t, outer.AgeGroup, inner.AgeGroup)
				assert.Equal(t, outer.Bib, inner.Bib)
				assert.Equal(t, outer.Distance, inner.Distance)
				assert.Equal(t, outer.First, inner.First)
				assert.Equal(t, outer.Gender, inner.Gender)
				assert.Equal(t, outer.Last, inner.Last)
				assert.Equal(t, outer.Anonymous, inner.Anonymous)
				assert.Equal(t, outer.AlternateId, inner.AlternateId)
				assert.Equal(t, outer.SMSEnabled, inner.SMSEnabled)
				assert.Equal(t, outer.Mobile, inner.Mobile)
				assert.Equal(t, outer.Apparel, inner.Apparel)
				found = true
			}
		}
		assert.True(t, found)
	}
	pulledParticipants = make([]types.Participant, 0)
	for i := range int(math.Ceil(float64(len(manyParticipants)) / 50)) {
		iParts, err = db.GetParticipants(eventYear.Identifier, 50, i, &updatedAt)
		if assert.NoError(t, err) {
			assert.True(t, len(iParts) > 0 && len(iParts) <= 50)
			for _, outer := range iParts {
				assert.GreaterOrEqual(t, outer.UpdatedAt, updatedAt)
				found := false
				for _, inner := range manyParticipants {
					if outer.AlternateId == inner.AlternateId {
						assert.True(t, outer.Equals(&inner))
						assert.Equal(t, outer.Birthdate, inner.Birthdate)
						assert.Equal(t, outer.AgeGroup, inner.AgeGroup)
						assert.Equal(t, outer.Bib, inner.Bib)
						assert.Equal(t, outer.Distance, inner.Distance)
						assert.Equal(t, outer.First, inner.First)
						assert.Equal(t, outer.Gender, inner.Gender)
						assert.Equal(t, outer.Last, inner.Last)
						assert.Equal(t, outer.Anonymous, inner.Anonymous)
						assert.Equal(t, outer.AlternateId, inner.AlternateId)
						assert.Equal(t, outer.SMSEnabled, inner.SMSEnabled)
						assert.Equal(t, outer.Mobile, inner.Mobile)
						assert.Equal(t, outer.Apparel, inner.Apparel)
						found = true
					}
				}
				assert.True(t, found)
				found = false
				for _, inner := range pulledParticipants {
					if outer.AlternateId == inner.AlternateId {
						found = true
					}
				}
				assert.False(t, found)
			}
			pulledParticipants = append(pulledParticipants, iParts...)
		}
	}
	assert.Equal(t, len(manyParticipants)-10, len(pulledParticipants))
}

func TestGetAllParticipants(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up tests. %v", err)
	}
	defer finalize(t)
	setupParticipantTests()
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
	iParts, err := db.GetParticipants(eventYear.Identifier, 0, 0, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(iParts))
	}
	db.AddParticipants(eventYear.Identifier, participants)
	iParts, err = db.GetParticipants(eventYear.Identifier, 0, 0, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, len(participants), len(iParts))
		for _, outer := range participants {
			found := false
			for _, inner := range iParts {
				if outer.AlternateId == inner.AlternateId {
					assert.True(t, outer.Equals(&inner))
					assert.Equal(t, outer.Birthdate, inner.Birthdate)
					assert.Equal(t, outer.AgeGroup, inner.AgeGroup)
					assert.Equal(t, outer.Bib, inner.Bib)
					assert.Equal(t, outer.Distance, inner.Distance)
					assert.Equal(t, outer.First, inner.First)
					assert.Equal(t, outer.Gender, inner.Gender)
					assert.Equal(t, outer.Last, inner.Last)
					assert.Equal(t, outer.Anonymous, inner.Anonymous)
					assert.Equal(t, outer.AlternateId, inner.AlternateId)
					assert.Equal(t, outer.SMSEnabled, inner.SMSEnabled)
					assert.Equal(t, outer.Mobile, inner.Mobile)
					assert.Equal(t, outer.Apparel, inner.Apparel)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
	db.DeleteParticipants(eventYear.Identifier, nil)
	db.AddParticipants(eventYear.Identifier, manyParticipants)
	updated_at := int64(50)
	iParts, err = db.GetParticipants(eventYear.Identifier, 0, 0, &updated_at)
	if assert.NoError(t, err) {
		assert.Equal(t, len(manyParticipants)-5, len(iParts))
		for _, outer := range iParts {
			assert.GreaterOrEqual(t, outer.UpdatedAt, updated_at)
			found := false
			for _, inner := range manyParticipants {
				if outer.AlternateId == inner.AlternateId {
					assert.True(t, outer.Equals(&inner))
					assert.Equal(t, outer.Birthdate, inner.Birthdate)
					assert.Equal(t, outer.AgeGroup, inner.AgeGroup)
					assert.Equal(t, outer.Bib, inner.Bib)
					assert.Equal(t, outer.Distance, inner.Distance)
					assert.Equal(t, outer.First, inner.First)
					assert.Equal(t, outer.Gender, inner.Gender)
					assert.Equal(t, outer.Last, inner.Last)
					assert.Equal(t, outer.Anonymous, inner.Anonymous)
					assert.Equal(t, outer.AlternateId, inner.AlternateId)
					assert.Equal(t, outer.SMSEnabled, inner.SMSEnabled)
					assert.Equal(t, outer.Mobile, inner.Mobile)
					assert.Equal(t, outer.Apparel, inner.Apparel)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
}

func TestDeleteParticipants(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up tests. %v", err)
	}
	defer finalize(t)
	setupParticipantTests()
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
	p, err := db.GetParticipants(eventYear.Identifier, 0, 0, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(p))
	}
	_, err = db.AddParticipants(eventYear.Identifier, participants)
	assert.NoError(t, err)
	p, err = db.GetParticipants(eventYear.Identifier, 0, 0, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, len(participants), len(p))
	}
	count, err := db.DeleteParticipants(eventYear.Identifier, nil)
	assert.NoError(t, err)
	p, err = db.GetParticipants(eventYear.Identifier, 0, 0, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(p))
	}
	assert.Equal(t, count, int64(len(participants)))
	_, err = db.AddParticipants(eventYear.Identifier, participants)
	assert.NoError(t, err)
	p, err = db.GetParticipants(eventYear.Identifier, 0, 0, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, len(participants), len(p))
	}
	toDelete := []string{
		participants[0].AlternateId,
		participants[1].AlternateId,
	}
	count, err = db.DeleteParticipants(eventYear.Identifier, toDelete)
	assert.NoError(t, err)
	assert.Equal(t, count, int64(len(toDelete)))
	p, err = db.GetParticipants(eventYear.Identifier, 0, 0, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, len(participants)-2, len(p))
		for _, outer := range participants[2:] {
			found := false
			for _, inner := range p {
				if outer.AlternateId == inner.AlternateId {
					assert.True(t, outer.Equals(&inner))
					assert.Equal(t, outer.Birthdate, inner.Birthdate)
					assert.Equal(t, outer.AgeGroup, inner.AgeGroup)
					assert.Equal(t, outer.Bib, inner.Bib)
					assert.Equal(t, outer.Distance, inner.Distance)
					assert.Equal(t, outer.First, inner.First)
					assert.Equal(t, outer.Gender, inner.Gender)
					assert.Equal(t, outer.Last, inner.Last)
					assert.Equal(t, outer.Anonymous, inner.Anonymous)
					assert.Equal(t, outer.AlternateId, inner.AlternateId)
					assert.Equal(t, outer.SMSEnabled, inner.SMSEnabled)
					assert.Equal(t, outer.Mobile, inner.Mobile)
					assert.Equal(t, outer.Apparel, inner.Apparel)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
}

func TestUpdateParticipant(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up tests. %v", err)
	}
	defer finalize(t)
	setupParticipantTests()
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
	parts, _ := db.AddParticipants(eventYear.Identifier, participants)
	assert.NotNil(t, parts)
	id := parts[0].Identifier
	// test update
	temp := participants[0]
	temp.Bib = "-200"
	temp.Birthdate = "4/12/2020"
	temp.AgeGroup = "Youngling"
	temp.First = "Update!"
	temp.Last = "Test"
	temp.Distance = "12 Mile Fun"
	temp.Gender = "U"
	temp.Anonymous = !participants[0].Anonymous
	part, err := db.UpdateParticipant(eventYear.Identifier, temp)
	if assert.NoError(t, err) {
		assert.True(t, temp.Equals(part))
		assert.Equal(t, temp.Birthdate, part.Birthdate)
		assert.Equal(t, temp.AgeGroup, part.AgeGroup)
		assert.Equal(t, temp.Bib, part.Bib)
		assert.Equal(t, temp.Distance, part.Distance)
		assert.Equal(t, temp.First, part.First)
		assert.Equal(t, temp.Gender, part.Gender)
		assert.Equal(t, temp.Last, part.Last)
		assert.Equal(t, id, part.Identifier)
		assert.Equal(t, temp.Anonymous, part.Anonymous)
		assert.Equal(t, temp.AlternateId, part.AlternateId)
		assert.Equal(t, temp.SMSEnabled, part.SMSEnabled)
		assert.Equal(t, temp.Mobile, part.Mobile)
		assert.Equal(t, temp.Apparel, part.Apparel)
	}
	// test invalid update
	temp = participants[1]
	temp.AlternateId = "newbib"
	part, err = db.UpdateParticipant(eventYear.Identifier, temp)
	assert.Error(t, err)
	assert.Nil(t, part)
}

func TestUpdateParticipants(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up tests. %v", err)
	}
	defer finalize(t)
	setupParticipantTests()
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
	parts, _ := db.AddParticipants(eventYear.Identifier, participants)
	assert.NotNil(t, parts)
	id := parts[0].Identifier
	// test update
	temp := participants[0]
	temp.Bib = "-200"
	temp.Birthdate = "4/12/2020"
	temp.AgeGroup = "Youngling"
	temp.First = "Update!"
	temp.Last = "Test"
	temp.Distance = "12 Mile Fun"
	temp.Gender = "U"
	temp.Anonymous = !participants[0].Anonymous
	part, err := db.UpdateParticipants(eventYear.Identifier, []types.Participant{temp})
	if assert.NoError(t, err) {
		assert.True(t, temp.Equals(&part[0]))
		assert.Equal(t, 1, len(part))
		assert.Equal(t, temp.Birthdate, part[0].Birthdate)
		assert.Equal(t, temp.AgeGroup, part[0].AgeGroup)
		assert.Equal(t, temp.Bib, part[0].Bib)
		assert.Equal(t, temp.Distance, part[0].Distance)
		assert.Equal(t, temp.First, part[0].First)
		assert.Equal(t, temp.Gender, part[0].Gender)
		assert.Equal(t, temp.Last, part[0].Last)
		assert.Equal(t, id, part[0].Identifier)
		assert.Equal(t, temp.Anonymous, part[0].Anonymous)
		assert.Equal(t, temp.AlternateId, part[0].AlternateId)
		assert.Equal(t, temp.SMSEnabled, part[0].SMSEnabled)
		assert.Equal(t, temp.Mobile, part[0].Mobile)
		assert.Equal(t, temp.Apparel, part[0].Apparel)
	}
	// test invalid update
	temp = participants[1]
	temp.AlternateId = "newbib"
	part, err = db.UpdateParticipants(eventYear.Identifier, []types.Participant{temp})
	assert.Error(t, err)
	assert.Nil(t, part)
}
