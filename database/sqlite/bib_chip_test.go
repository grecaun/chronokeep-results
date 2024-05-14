package sqlite

import (
	"chronokeep/results/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	bibChips []types.BibChip
)

func setupBibChipsTests() {
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
	bibChips = []types.BibChip{
		{
			Bib:  "100",
			Chip: "202",
		},
		{
			Bib:  "101",
			Chip: "203",
		},
		{
			Bib:  "102",
			Chip: "204",
		},
		{
			Bib:  "103",
			Chip: "2222",
		},
		{
			Bib:  "104",
			Chip: "2223",
		},
		{
			Bib:  "105",
			Chip: "2224",
		},
		{
			Bib:  "106",
			Chip: "2225",
		},
		{
			Bib:  "107",
			Chip: "abcf",
		},
		{
			Bib:  "108",
			Chip: "a255c",
		},
	}
}

func TestAddBibChips(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up tests. %v", err)
	}
	defer finalize(t)
	setupBibChipsTests()
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
	bc, err := db.GetBibChips(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(bc))
	}
	bc, err = db.AddBibChips(eventYear.Identifier, bibChips)
	if assert.NoError(t, err) {
		assert.Equal(t, len(bibChips), len(bc))
		for _, outer := range bibChips {
			found := false
			for _, inner := range bc {
				if outer.Chip == inner.Chip {
					assert.Equal(t, outer.Bib, inner.Bib)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
	bc, err = db.GetBibChips(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, len(bibChips), len(bc))
		for _, outer := range bibChips {
			found := false
			for _, inner := range bc {
				if outer.Chip == inner.Chip {
					assert.Equal(t, outer.Bib, inner.Bib)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
	// test update
	upd := make([]types.BibChip, 0)
	for _, temp := range bibChips {
		upd = append(upd, types.BibChip{
			Bib:  temp.Bib + "new",
			Chip: temp.Chip,
		})
	}
	bc, err = db.AddBibChips(eventYear.Identifier, upd)
	if assert.NoError(t, err) {
		for _, outer := range upd {
			found := false
			for _, inner := range bc {
				if outer.Chip == inner.Chip {
					assert.Equal(t, outer.Bib, inner.Bib)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
	bc, err = db.GetBibChips(eventYear.Identifier)
	if assert.NoError(t, err) {
		for _, outer := range upd {
			found := false
			for _, inner := range bc {
				if outer.Chip == inner.Chip {
					assert.Equal(t, outer.Bib, inner.Bib)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
}

func TestGetBibChips(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up tests. %v", err)
	}
	defer finalize(t)
	setupBibChipsTests()
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
	bc, err := db.GetBibChips(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(bc))
	}
	db.AddBibChips(eventYear.Identifier, bibChips)
	bc, err = db.GetBibChips(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, len(bibChips), len(bc))
		for _, outer := range bibChips {
			found := false
			for _, inner := range bc {
				if outer.Chip == inner.Chip {
					assert.Equal(t, outer.Bib, inner.Bib)
					found = true
				}
			}
			assert.True(t, found)
		}
	}
}

func TestDeleteBibChips(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up tests. %v", err)
	}
	defer finalize(t)
	setupBibChipsTests()
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
	bc, err := db.GetBibChips(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(bc))
	}
	db.AddBibChips(eventYear.Identifier, bibChips)
	bc, err = db.GetBibChips(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, len(bibChips), len(bc))
	}
	count, err := db.DeleteBibChips(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, int64(len(bibChips)), count)
	}
	bc, err = db.GetBibChips(eventYear.Identifier)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(bc))
	}
}
