package database

import (
	"chronokeep/results/types"
	"testing"
	"time"
)

var (
	results []types.Result
)

func setupResultTests() {
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
	if len(results) < 1 {
		results = []types.Result{
			{
				Bib:           "100",
				First:         "John",
				Last:          "Smith",
				Age:           24,
				Gender:        "M",
				AgeGroup:      "20-29",
				Distance:      "1 Mile",
				Seconds:       377,
				Milliseconds:  0,
				Segment:       "",
				Location:      "Start/Finish",
				Occurence:     1,
				Ranking:       1,
				AgeRanking:    1,
				GenderRanking: 1,
				Finish:        true,
			},
			{
				Bib:           "106",
				First:         "Rose",
				Last:          "Johnson",
				Age:           55,
				Gender:        "F",
				AgeGroup:      "50-59",
				Distance:      "1 Mile",
				Seconds:       577,
				Milliseconds:  100,
				Segment:       "",
				Location:      "Start/Finish",
				Occurence:     1,
				Ranking:       3,
				AgeRanking:    1,
				GenderRanking: 1,
				Finish:        true,
			},
			{
				Bib:           "209",
				First:         "Tony",
				Last:          "Starke",
				Age:           45,
				Gender:        "M",
				AgeGroup:      "40-49",
				Distance:      "1 Mile",
				Seconds:       405,
				Milliseconds:  20,
				Segment:       "",
				Location:      "Start/Finish",
				Occurence:     1,
				Ranking:       2,
				AgeRanking:    1,
				GenderRanking: 2,
				Finish:        true,
			},
			{
				Bib:           "287",
				First:         "Jamie",
				Last:          "Fischer",
				Age:           35,
				Gender:        "F",
				AgeGroup:      "30-39",
				Distance:      "1 Mile",
				Seconds:       653,
				Milliseconds:  0,
				Segment:       "",
				Location:      "Start/Finish",
				Occurence:     1,
				Ranking:       4,
				AgeRanking:    1,
				GenderRanking: 2,
				Finish:        false,
			},
			{
				Bib:           "287",
				First:         "Jamie",
				Last:          "Fischer",
				Age:           35,
				Gender:        "F",
				AgeGroup:      "30-39",
				Distance:      "1 Mile",
				Seconds:       1003,
				Milliseconds:  0,
				Segment:       "",
				Location:      "Start/Finish",
				Occurence:     2,
				Ranking:       4,
				AgeRanking:    1,
				GenderRanking: 2,
				Finish:        true,
			},
		}
	}
}

func TestAddResults(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupResultTests()
	account, _ := AddAccount(accounts[0])
	event := &types.Event{
		AccountIdentifier: account.Identifier,
		Name:              "Event 1",
		Slug:              "event1",
	}
	event, _ = AddEvent(*event)
	eventYear := &types.EventYear{
		EventIdentifier: event.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 04, 20, 9, 0, 0, 0, time.Local),
	}
	eventYear, _ = AddEventYear(*eventYear)
	res, err := AddResults(eventYear.Identifier, results)
	if err != nil {
		t.Fatalf("Error adding results: %v", err)
	}
	if len(res) != len(results) {
		t.Errorf("Expected %v results to be added, %v added.", len(results), len(res))
	}
	// Test the update feature of AddResults.
	results[0].Seconds = 30
	results[0].First = "Johnny"
	results[0].Last = "Cash"
	results[0].Age = 100
	results[0].Gender = "U"
	results[0].AgeGroup = "Old"
	results[0].Distance = "Too Far"
	results[0].Milliseconds = 300000000
	results[0].Segment = "Something"
	results[0].Ranking = 12
	results[0].AgeRanking = 23
	results[0].GenderRanking = 45
	results[0].Finish = false
	res, err = AddResults(eventYear.Identifier, results[0:1])
	if err != nil {
		t.Fatalf("Error adding results: %v", err)
	}
	if len(res) != 1 {
		t.Errorf("Expected %v results to be added, %v added.", 1, len(res))
	}
	res, _ = GetResults(eventYear.Identifier)
	if len(res) != len(results) {
		t.Errorf("Expected %v results to be added, %v added.", len(results), len(res))
	}
	// Test the update feature of AddResults.
	results[1].Seconds = 30
	results[1].First = "Rebecca"
	results[1].Last = "Small"
	results[1].Age = 10
	results[1].Gender = "U"
	results[1].AgeGroup = "Young"
	results[1].Distance = "Not Far Enough"
	results[1].Milliseconds = 300
	results[1].Segment = "Something"
	results[1].Occurence = 12
	results[1].Ranking = 121
	results[1].AgeRanking = 231
	results[1].GenderRanking = 451
	res, err = AddResults(eventYear.Identifier, results[1:2])
	if err != nil {
		t.Fatalf("Error adding results: %v", err)
	}
	if len(res) != 1 {
		t.Errorf("Expected %v results to be added, %v added.", 1, len(res))
	}
	res, _ = GetResults(eventYear.Identifier)
	if len(res) != (len(results) + 1) {
		t.Errorf("Expected %v results to be added, %v added.", (len(results) + 1), len(res))
	}
}

func TestGetResults(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupResultTests()
	account, _ := AddAccount(accounts[0])
	event := &types.Event{
		AccountIdentifier: account.Identifier,
		Name:              "Event 1",
		Slug:              "event1",
	}
	event, _ = AddEvent(*event)
	eventYear := &types.EventYear{
		EventIdentifier: event.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 04, 20, 9, 0, 0, 0, time.Local),
	}
	eventYear, _ = AddEventYear(*eventYear)
	res, err := GetResults(eventYear.Identifier)
	if err != nil {
		t.Fatalf("Error getting results: %v", err)
	}
	if len(res) != 0 {
		t.Errorf("Results not added but we've found %v results.", len(res))
	}
	AddResults(eventYear.Identifier, results[0:1])
	res, err = GetResults(eventYear.Identifier)
	if err != nil {
		t.Fatalf("Error getting results: %v", err)
	}
	if len(res) != 1 {
		t.Errorf("Expected %v results to be added, %v added.", 1, len(res))
	}
	if res[0] != results[0] {
		t.Errorf("Expected results %+v, found %+v.", results[0], res[0])
	}
	AddResults(eventYear.Identifier, results)
	res, err = GetResults(eventYear.Identifier)
	if err != nil {
		t.Fatalf("Error getting results: %v", err)
	}
	if len(res) != len(results) {
		t.Errorf("Expected %v results to be added, %v added.", len(results), len(res))
	}
	// Verify that we've got correct information for our results.
	found := false
	for _, result := range res {
		if result == results[0] {
			found = true
		}
	}
	if !found {
		t.Errorf("Unable to find our first result in the database. %+v", res)
	}
}

func TestDeleteResults(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupResultTests()
	account, _ := AddAccount(accounts[0])
	event := &types.Event{
		AccountIdentifier: account.Identifier,
		Name:              "Event 1",
		Slug:              "event1",
	}
	event, _ = AddEvent(*event)
	eventYear := &types.EventYear{
		EventIdentifier: event.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 04, 20, 9, 0, 0, 0, time.Local),
	}
	eventYear, _ = AddEventYear(*eventYear)
	AddResults(eventYear.Identifier, results)
	err = DeleteResults(eventYear.Identifier, results[1:2])
	if err != nil {
		t.Fatalf("Error deleting specific results: %v", err)
	}
	res, _ := GetResults(eventYear.Identifier)
	if len(res) != (len(results) - 1) {
		t.Errorf("Expected %v results after delete but found %v.", (len(results) - 1), len(res))
	}
	found := false
	for _, result := range res {
		if result == results[1] {
			found = true
		}
	}
	if found {
		t.Errorf("Found deleted result.")
	}
}

func TestDeleteEventResults(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupResultTests()
	account, _ := AddAccount(accounts[0])
	event := &types.Event{
		AccountIdentifier: account.Identifier,
		Name:              "Event 1",
		Slug:              "event1",
	}
	event, _ = AddEvent(*event)
	eventYear := &types.EventYear{
		EventIdentifier: event.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 04, 20, 9, 0, 0, 0, time.Local),
	}
	eventYear, _ = AddEventYear(*eventYear)
	AddResults(eventYear.Identifier, results)
	count, err := DeleteEventResults(eventYear.Identifier)
	if err != nil {
		t.Fatalf("Error deleting specific results: %v", err)
	}
	if (count) != int64(len(results)) {
		t.Errorf("Expected to find out %v results were deleted, %v were deleted.", len(results), count)
	}
	res, _ := GetResults(eventYear.Identifier)
	if len(res) != 0 {
		t.Errorf("Expected %v results after delete but found %v.", 0, len(res))
	}
}

func TestGetBibResults(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupResultTests()
	account, _ := AddAccount(accounts[0])
	event := &types.Event{
		AccountIdentifier: account.Identifier,
		Name:              "Event 1",
		Slug:              "event1",
	}
	event, _ = AddEvent(*event)
	eventYear := &types.EventYear{
		EventIdentifier: event.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 04, 20, 9, 0, 0, 0, time.Local),
	}
	eventYear, _ = AddEventYear(*eventYear)
	lastIX := len(results) - 1
	res, err := GetBibResults(eventYear.Identifier, results[lastIX].Bib)
	if err != nil {
		t.Fatalf("Error getting bib results: %v", err)
	}
	if len(res) != 0 {
		t.Errorf("Expected %v results to be added, %v added.", 0, len(res))
	}
	AddResults(eventYear.Identifier, results[0:lastIX])
	res, err = GetBibResults(eventYear.Identifier, results[lastIX].Bib)
	if err != nil {
		t.Fatalf("Error getting bib results: %v", err)
	}
	if len(res) != 1 {
		t.Errorf("Expected %v results to be added, %v added.", 1, len(res))
	}
	AddResults(eventYear.Identifier, results)
	res, err = GetBibResults(eventYear.Identifier, results[lastIX].Bib)
	if err != nil {
		t.Fatalf("Error getting bib results: %v", err)
	}
	if len(res) != 2 {
		t.Errorf("Expected %v results to be added, %v added.", 2, len(res))
	}
	// Verify that we've got correct information for our results.
	found := false
	for _, result := range res {
		if result == results[lastIX] {
			found = true
		}
	}
	if !found {
		t.Errorf("Unable to find our first result in the database. %+v", res)
	}
}
