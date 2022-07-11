package handlers

import (
	"chronokeep/results/auth"
	"chronokeep/results/database/sqlite"
	"chronokeep/results/types"
	"chronokeep/results/util"
	"os"
	"strconv"
	"testing"
	"time"
)

func setupTests(t *testing.T) (SetupVariables, func(t *testing.T)) {
	t.Log("Setting up sqlite database.")
	database = &sqlite.SQLite{}
	config = &util.Config{
		DBName:     "./results_test.sqlite",
		DBHost:     "",
		DBUser:     "",
		DBPassword: "",
		DBPort:     0,
		DBDriver:   "sqlite3",
	}
	database.Setup(config)
	t.Log("Setting up config variables to export.")
	output := SetupVariables{
		testPassword1: "amazingpassword",
		testPassword2: "othergoodpassword",
	}
	// add accounts
	t.Log("Adding accounts.")
	for _, account := range []types.Account{
		{
			Name:     "John Smith",
			Email:    "j@test.com",
			Type:     "admin",
			Password: testHashPassword(output.testPassword1),
		},
		{
			Name:     "Jerry Garcia",
			Email:    "jgarcia@test.com",
			Type:     "free",
			Password: testHashPassword(output.testPassword1),
		},
		{
			Name:     "Rose MacDonald",
			Email:    "rose2004@test.com",
			Type:     "paid",
			Password: testHashPassword(output.testPassword2),
		},
	} {
		database.AddAccount(account)
	}
	var err error
	output.accounts, err = database.GetAccounts()
	if err != nil {
		t.Fatalf("Unexpected error adding accounts: %v", err)
	}
	t.Log("Adding Keys.")
	// add keys, one expired, one with a timer, one write, one read, one delete, two different accounts
	times := []time.Time{
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
		time.Now().Add(time.Hour * 20).Truncate(time.Second),
	}
	for _, key := range []types.Key{
		{
			AccountIdentifier: output.accounts[0].Identifier,
			Name:              "expired1",
			Value:             "030001-1ACSDD-K2389A-00123B",
			Type:              "read",
			AllowedHosts:      "",
			ValidUntil:        &times[0],
		},
		{
			AccountIdentifier: output.accounts[0].Identifier,
			Value:             "030001-1ACSDD-K2389A-22123B",
			Type:              "write",
			AllowedHosts:      "",
			ValidUntil:        &times[1],
		},
		{
			AccountIdentifier: output.accounts[1].Identifier,
			Value:             "030001-1ACSCT-K2389A-22023B",
			Type:              "delete",
			AllowedHosts:      "",
			ValidUntil:        nil,
		},
		{
			AccountIdentifier: output.accounts[1].Identifier,
			Value:             "030001-1ACSCT-K2389A-22423B",
			Type:              "read",
			AllowedHosts:      "chronokeep.com",
			ValidUntil:        nil,
		},
	} {
		database.AddKey(key)
	}
	output.keys = make([]types.Key, 0)
	acckeys, err := database.GetAccountKeys(output.accounts[0].Email)
	if err != nil {
		t.Fatalf("Unexptected error getting keys: %v", err)
	}
	output.keys = append(output.keys, acckeys...)
	acckeys, err = database.GetAccountKeys(output.accounts[1].Email)
	if err != nil {
		t.Fatalf("Unexptected error getting keys: %v", err)
	}
	output.keys = append(output.keys, acckeys...)
	t.Log("Adding events.")
	// add two events, one for each account
	for _, event := range []types.Event{
		{
			AccountIdentifier: output.accounts[0].Identifier,
			Name:              "Event 1",
			Slug:              "event1",
			ContactEmail:      "event1@test.com",
			AccessRestricted:  false,
		},
		{
			AccountIdentifier: output.accounts[1].Identifier,
			Name:              "Event 2",
			Slug:              "event2",
			ContactEmail:      "event2@test.com",
			AccessRestricted:  true,
		},
	} {
		database.AddEvent(event)
	}
	for _, account := range output.accounts {
		tmp, err := database.GetAccountEvents(account.Email)
		if err != nil {
			t.Fatalf("Unexpected error getting events: %v", err)
		}
		output.events = append(output.events, tmp...)
	}
	// add event years, two per event
	t.Log("Adding event years.")
	for _, eventYear := range []types.EventYear{
		{
			EventIdentifier: output.events[0].Identifier,
			Year:            "2021",
			DateTime:        time.Date(2021, 10, 06, 9, 0, 0, 0, time.Local),
			Live:            false,
		},
		{
			EventIdentifier: output.events[0].Identifier,
			Year:            "2020",
			DateTime:        time.Date(2020, 10, 05, 9, 0, 0, 0, time.Local),
			Live:            false,
		},
		{
			EventIdentifier: output.events[1].Identifier,
			Year:            "2021",
			DateTime:        time.Date(2021, 04, 05, 11, 0, 0, 0, time.Local),
			Live:            false,
		},
		{
			EventIdentifier: output.events[1].Identifier,
			Year:            "2020",
			DateTime:        time.Date(2021, 04, 05, 11, 0, 0, 0, time.Local),
			Live:            false,
		},
	} {
		database.AddEventYear(eventYear)
	}
	evYear, err := database.GetEventYears(output.events[0].Slug)
	if err != nil {
		t.Fatalf("Unexpected error getting event years: %v", err)
	}
	output.eventYears = append(output.eventYears, evYear...)
	evYear, err = database.GetEventYears(output.events[1].Slug)
	if err != nil {
		t.Fatalf("Unexpected error getting event years: %v", err)
	}
	output.eventYears = append(output.eventYears, evYear...)
	// add results
	t.Log("Adding results.")
	res := make([]types.Result, 0)
	for i := 0; i < 500; i++ {
		tmpStr := strconv.Itoa(i)
		res = append(res, types.Result{
			Bib:           tmpStr,
			First:         "John" + tmpStr,
			Last:          "Smith",
			Age:           24,
			Gender:        "M",
			AgeGroup:      "20-29",
			Distance:      "1 Mile",
			Seconds:       377 + i*5,
			Milliseconds:  0,
			Segment:       "",
			Location:      "Start/Finish",
			Occurence:     1,
			Ranking:       i + 1,
			AgeRanking:    i + 1,
			GenderRanking: i + 1,
			Finish:        true,
		})
	}
	_, err = database.AddResults(output.eventYears[0].Identifier, res)
	if err != nil {
		t.Fatalf("Unexpected error adding 500 results to database: %v", err)
	}
	for _, eventYear := range output.eventYears[1:] {
		_, err = database.AddResults(eventYear.Identifier, []types.Result{
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
				Finish:        false,
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
				Distance:      "2 Mile",
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
				Distance:      "2 Mile",
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
		})
		if err != nil {
			t.Fatalf("Unexpected error adding small number of results to database: %v", err)
		}
	}
	for _, eventYear := range output.eventYears {
		res, err = database.GetResults(eventYear.Identifier, 0, 0)
		if err != nil {
			t.Fatalf("Unexpected error getting results from database: %v", err)
		}
		output.results = append(output.results, res...)
	}
	return output, func(t *testing.T) {
		t.Log("Deleting old database.")
		database.Close()
		err := os.Remove(config.DBName)
		if err != nil {
			t.Fatalf("Error deleting database: %v", err)
		}
		t.Log("Cleanup successful.")
	}
}

func testHashPassword(pass string) string {
	hash, _ := auth.HashPassword(pass)
	return hash
}

type SetupVariables struct {
	accounts      []types.Account
	testPassword1 string
	testPassword2 string
	keys          []types.Key
	events        []types.Event
	eventYears    []types.EventYear
	results       []types.Result
}
