package handlers

import (
	"chronokeep/results/auth"
	"chronokeep/results/database/sqlite"
	"chronokeep/results/types"
	"chronokeep/results/util"
	"fmt"
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
		knownValues:   make(map[string]string),
		keys:          make(map[string][]types.Key),
		events:        make(map[string]types.Event),
		eventYears:    make(map[string]map[string]types.EventYear),
		results:       make(map[string]map[string][]types.Result),
		segments:      make(map[string]map[string][]types.Segment),
	}
	// add accounts
	t.Log("Adding accounts.")
	output.knownValues["admin"] = "j@test.com"
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
		{
			Name:     "Registration",
			Email:    "registration@test.com",
			Type:     "registration",
			Password: testHashPassword(output.testPassword2),
		},
		{
			Name:     "Registration2",
			Email:    "registration2@test.com",
			Type:     "registration",
			Password: testHashPassword(output.testPassword2),
		},
		{
			Name:     "Registration3",
			Email:    "registration3@test.com",
			Type:     "registration",
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
	output.knownValues["expired"] = "030001-1ACSDD-K2389A-00123B"
	output.knownValues["expired2"] = "030001-1ACSDD-K2389A-001230B"
	output.knownValues["delete"] = "030001-1ACSCT-K2389A-22023B"
	output.knownValues["delete2"] = "030001-1ACSCT-K2389A-22023BAA"
	output.knownValues["delete3"] = "0030001-1ACSCT-K2389A-22023BAA"
	output.knownValues["read"] = "030001-1ACSCT-K2389A-22423B"
	output.knownValues["write"] = "030001-1ACSDD-K2389A-22123B"
	output.knownValues["write2"] = "030001-1ACSCT-K2389A-22423BAA"
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
			Name:              "write",
			Value:             "030001-1ACSDD-K2389A-22123B",
			Type:              "write",
			AllowedHosts:      "",
			ValidUntil:        &times[1],
		},
		{
			AccountIdentifier: output.accounts[0].Identifier,
			Value:             "0030001-1ACSCT-K2389A-22023BAA",
			Name:              "delete3",
			Type:              "delete",
			AllowedHosts:      "",
			ValidUntil:        nil,
		},
		{
			AccountIdentifier: output.accounts[1].Identifier,
			Value:             "030001-1ACSCT-K2389A-22023B",
			Name:              "delete",
			Type:              "delete",
			AllowedHosts:      "chronokeep.com",
			ValidUntil:        nil,
		},
		{
			AccountIdentifier: output.accounts[1].Identifier,
			Value:             "030001-1ACSCT-K2389A-22423B",
			Name:              "read",
			Type:              "read",
			AllowedHosts:      "",
			ValidUntil:        nil,
		},
		{
			AccountIdentifier: output.accounts[1].Identifier,
			Value:             "030001-1ACSCT-K2389A-22423BAA",
			Name:              "write2",
			Type:              "write",
			AllowedHosts:      "",
			ValidUntil:        nil,
		},
		{
			AccountIdentifier: output.accounts[1].Identifier,
			Value:             "030001-1ACSCT-K2389A-22023BAA",
			Name:              "delete2",
			Type:              "delete",
			AllowedHosts:      "",
			ValidUntil:        nil,
		},
		{
			AccountIdentifier: output.accounts[1].Identifier,
			Name:              "expired2",
			Value:             "030001-1ACSDD-K2389A-001230B",
			Type:              "delete",
			AllowedHosts:      "",
			ValidUntil:        &times[0],
		},
	} {
		database.AddKey(key)
	}
	output.keys[output.accounts[0].Email], err = database.GetAccountKeys(output.accounts[0].Email)
	if err != nil {
		t.Fatalf("Unexptected error getting keys: %v", err)
	}
	output.keys[output.accounts[1].Email], err = database.GetAccountKeys(output.accounts[1].Email)
	if err != nil {
		t.Fatalf("Unexptected error getting keys: %v", err)
	}
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
		{
			AccountIdentifier: output.accounts[1].Identifier,
			Name:              "Event 3",
			Slug:              "event3",
			ContactEmail:      "event3@test.com",
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
		for _, event := range tmp {
			output.events[event.Slug] = event
			output.eventYears[event.Slug] = make(map[string]types.EventYear)
			output.results[event.Slug] = make(map[string][]types.Result)
			output.segments[event.Slug] = make(map[string][]types.Segment)
		}
	}
	// add event years, two per event
	t.Log("Adding event years.")
	for _, eventYear := range []types.EventYear{
		{
			EventIdentifier: output.events["event1"].Identifier,
			Year:            "2021",
			DateTime:        time.Date(2021, 10, 06, 9, 0, 0, 0, time.Local),
			Live:            false,
		},
		{
			EventIdentifier: output.events["event1"].Identifier,
			Year:            "2020",
			DateTime:        time.Date(2020, 10, 05, 9, 0, 0, 0, time.Local),
			Live:            false,
		},
		{
			EventIdentifier: output.events["event2"].Identifier,
			Year:            "2021",
			DateTime:        time.Date(2021, 04, 05, 11, 0, 0, 0, time.Local),
			Live:            false,
		},
		{
			EventIdentifier: output.events["event2"].Identifier,
			Year:            "2020",
			DateTime:        time.Date(2020, 04, 05, 11, 0, 0, 0, time.Local),
			Live:            false,
		},
		{
			EventIdentifier: output.events["event2"].Identifier,
			Year:            "2019",
			DateTime:        time.Date(2019, 04, 05, 11, 0, 0, 0, time.Local),
			Live:            false,
		},
		{
			EventIdentifier: output.events["event3"].Identifier,
			Year:            fmt.Sprintf("%v", time.Now().Year()),
			DateTime:        time.Now(),
			Live:            false,
			DaysAllowed:     2,
		},
	} {
		database.AddEventYear(eventYear)
	}
	evYear, err := database.GetEventYears(output.events["event1"].Slug)
	if err != nil {
		t.Fatalf("Unexpected error getting event years: %v", err)
	}
	for _, eventYear := range evYear {
		output.eventYears[output.events["event1"].Slug][eventYear.Year] = eventYear
	}
	evYear, err = database.GetEventYears(output.events["event2"].Slug)
	if err != nil {
		t.Fatalf("Unexpected error getting event years: %v", err)
	}
	for _, eventYear := range evYear {
		output.eventYears[output.events["event2"].Slug][eventYear.Year] = eventYear
	}
	evYear, err = database.GetEventYears(output.events["event3"].Slug)
	if err != nil {
		t.Fatalf("Unexpected error getting event years: %v", err)
	}
	for _, eventYear := range evYear {
		fmt.Printf("%v\n", eventYear)
		output.eventYears[output.events["event3"].Slug][eventYear.Year] = eventYear
	}
	// add results
	t.Log("Adding results.")
	res := make([]types.Result, 0)
	parts := make([]types.Participant, 0)
	bibChips := make([]types.BibChip, 0)
	for i := 0; i < 300; i++ {
		tmpStr := strconv.Itoa(i)
		res = append(res, types.Result{
			PersonId:      tmpStr,
			Bib:           tmpStr,
			First:         "John" + tmpStr,
			Last:          "Smith",
			Age:           24,
			Gender:        "Man",
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
		parts = append(parts, types.Participant{
			AlternateId: tmpStr,
			Bib:         tmpStr,
			First:       "John" + tmpStr,
			Last:        "Smith",
			Birthdate:   "1/1/2000",
			Gender:      "Man",
			AgeGroup:    "20-29",
			Distance:    "1 Mile",
			Anonymous:   false,
			SMSEnabled:  false,
			Mobile:      "",
			Apparel:     "",
		})
		bibChips = append(bibChips, types.BibChip{
			Bib:  tmpStr,
			Chip: "chip" + tmpStr,
		})
	}
	for i := 300; i < 500; i++ {
		tmpStr := strconv.Itoa(i)
		res = append(res, types.Result{
			PersonId:      tmpStr,
			Bib:           tmpStr,
			First:         "Jane" + tmpStr,
			Last:          "Smithson",
			Age:           27,
			Gender:        "Woman",
			AgeGroup:      "20-29",
			Distance:      "2 Mile",
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
		parts = append(parts, types.Participant{
			AlternateId: tmpStr,
			Bib:         tmpStr,
			First:       "Jane" + tmpStr,
			Last:        "Smithson",
			Birthdate:   "1/1/1997",
			Gender:      "Woman",
			AgeGroup:    "20-29",
			Distance:    "2 Mile",
			Anonymous:   false,
			SMSEnabled:  false,
			Mobile:      "",
			Apparel:     "",
		})
		bibChips = append(bibChips, types.BibChip{
			Bib:  tmpStr,
			Chip: "chip" + tmpStr,
		})
	}
	for i := 0; i < 300; i++ {
		tmpStr := strconv.Itoa(i)
		res = append(res, types.Result{
			PersonId:      tmpStr,
			Bib:           tmpStr,
			First:         "John" + tmpStr,
			Last:          "Smith",
			Age:           24,
			Gender:        "Man",
			AgeGroup:      "20-29",
			Distance:      "1 Mile",
			Seconds:       i,
			Milliseconds:  0,
			Segment:       "",
			Location:      "Start/Finish",
			Occurence:     0,
			Ranking:       i + 1,
			AgeRanking:    i + 1,
			GenderRanking: i + 1,
			Finish:        false,
		})
		parts = append(parts, types.Participant{
			AlternateId: tmpStr,
			Bib:         tmpStr,
			First:       "John" + tmpStr,
			Last:        "Smith",
			Birthdate:   "1/1/2000",
			Gender:      "Man",
			AgeGroup:    "20-29",
			Distance:    "1 Mile",
			Anonymous:   false,
			SMSEnabled:  false,
			Mobile:      "",
			Apparel:     "",
		})
		bibChips = append(bibChips, types.BibChip{
			Bib:  tmpStr,
			Chip: "chip" + tmpStr,
		})
	}
	_, err = database.AddResults(output.eventYears["event1"]["2021"].Identifier, res)
	if err != nil {
		t.Fatalf("Unexpected error adding 500 results to database: %v", err)
	}
	_, err = database.AddParticipants(output.eventYears["event1"]["2021"].Identifier, parts)
	if err != nil {
		t.Fatalf("Unexpected error adding 500 participants to database: %v", err)
	}
	_, err = database.AddBibChips(output.eventYears["event1"]["2021"].Identifier, bibChips)
	if err != nil {
		t.Fatalf("Unexpected error adding 500 bibchips to database: %v", err)
	}
	_, err = database.AddResults(output.eventYears["event1"]["2020"].Identifier, []types.Result{
		{
			PersonId:      "100",
			Bib:           "100",
			First:         "John",
			Last:          "Smith",
			Age:           24,
			Gender:        "Non-Binary",
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
			PersonId:      "106",
			Bib:           "106",
			First:         "Rose",
			Last:          "Johnson",
			Age:           55,
			Gender:        "Woman",
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
			PersonId:      "209",
			Bib:           "209",
			First:         "Tony",
			Last:          "Starke",
			Age:           45,
			Gender:        "Man",
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
			PersonId:      "287",
			Bib:           "287",
			First:         "Jamie",
			Last:          "Fischer",
			Age:           35,
			Gender:        "Woman",
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
			PersonId:      "287",
			Bib:           "287",
			First:         "Jamie",
			Last:          "Fischer",
			Age:           35,
			Gender:        "Woman",
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
	_, err = database.AddParticipants(output.eventYears["event1"]["2020"].Identifier, []types.Participant{
		{
			AlternateId: "100",
			Bib:         "100",
			First:       "John",
			Last:        "Smith",
			Birthdate:   "1/1/2000",
			Gender:      "Non-Binary",
			AgeGroup:    "20-29",
			Distance:    "1 Mile",
			Anonymous:   false,
			SMSEnabled:  false,
			Mobile:      "",
			Apparel:     "",
		},
		{
			AlternateId: "106",
			Bib:         "106",
			First:       "Rose",
			Last:        "Johnson",
			Birthdate:   "1/1/1969",
			Gender:      "Woman",
			AgeGroup:    "50-59",
			Distance:    "1 Mile",
			Anonymous:   false,
			SMSEnabled:  false,
			Mobile:      "",
			Apparel:     "",
		},
		{
			AlternateId: "209",
			Bib:         "209",
			First:       "Tony",
			Last:        "Starke",
			Birthdate:   "1/1/1979",
			Gender:      "Man",
			AgeGroup:    "40-49",
			Distance:    "1 Mile",
			Anonymous:   false,
			SMSEnabled:  false,
			Mobile:      "",
			Apparel:     "",
		},
		{
			AlternateId: "287",
			Bib:         "287",
			First:       "Jamie",
			Last:        "Fischer",
			Birthdate:   "1/1/1989",
			Gender:      "Woman",
			AgeGroup:    "30-39",
			Distance:    "2 Mile",
			Anonymous:   false,
			SMSEnabled:  false,
			Mobile:      "",
			Apparel:     "",
		},
	})
	if err != nil {
		t.Fatalf("Unexpected error adding small number of participants to database: %v", err)
	}
	_, err = database.AddBibChips(output.eventYears["event1"]["2020"].Identifier, []types.BibChip{
		{
			Chip: "chip100",
			Bib:  "100",
		},
		{
			Chip: "chip106",
			Bib:  "106",
		},
		{
			Chip: "chip209",
			Bib:  "209",
		},
		{
			Chip: "chip287",
			Bib:  "287",
		},
	})
	if err != nil {
		t.Fatalf("Unexpected error adding small number of bibchips to database: %v", err)
	}
	for _, eventYear := range output.eventYears["event2"] {
		_, err = database.AddResults(eventYear.Identifier, []types.Result{
			{
				PersonId:      "1001",
				Bib:           "100",
				First:         "John",
				Last:          "Smith",
				Age:           24,
				Gender:        "Man",
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
				PersonId:      "1006",
				Bib:           "106",
				First:         "Rose",
				Last:          "Johnson",
				Age:           55,
				Gender:        "Woman",
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
				PersonId:      "2009",
				Bib:           "209",
				First:         "Tony",
				Last:          "Starke",
				Age:           45,
				Gender:        "Man",
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
				PersonId:      "287",
				Bib:           "287",
				First:         "Jamie",
				Last:          "Fischer",
				Age:           35,
				Gender:        "Woman",
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
				PersonId:      "287",
				Bib:           "287",
				First:         "Jamie",
				Last:          "Fischer",
				Age:           35,
				Gender:        "Woman",
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
		_, err = database.AddParticipants(eventYear.Identifier, []types.Participant{
			{
				AlternateId: "100",
				Bib:         "100",
				First:       "John",
				Last:        "Smith",
				Birthdate:   "1/1/2024",
				Gender:      "Man",
				AgeGroup:    "20-29",
				Distance:    "1 Mile",
				Anonymous:   false,
				SMSEnabled:  false,
				Mobile:      "",
				Apparel:     "",
			},
			{
				AlternateId: "106",
				Bib:         "106",
				First:       "Rose",
				Last:        "Johnson",
				Birthdate:   "1/1/1969",
				Gender:      "Woman",
				AgeGroup:    "50-59",
				Distance:    "1 Mile",
				Anonymous:   false,
				SMSEnabled:  false,
				Mobile:      "",
				Apparel:     "",
			},
			{
				AlternateId: "209",
				Bib:         "209",
				First:       "Tony",
				Last:        "Starke",
				Birthdate:   "1/1/1979",
				Gender:      "Man",
				AgeGroup:    "40-49",
				Distance:    "1 Mile",
				Anonymous:   false,
				SMSEnabled:  false,
				Mobile:      "",
				Apparel:     "",
			},
			{
				AlternateId: "287",
				Bib:         "287",
				First:       "Jamie",
				Last:        "Fischer",
				Birthdate:   "1/1/1989",
				Gender:      "Woman",
				AgeGroup:    "30-39",
				Distance:    "2 Mile",
				Anonymous:   false,
				SMSEnabled:  false,
				Mobile:      "",
				Apparel:     "",
			},
		})
		if err != nil {
			t.Fatalf("Unexpected error adding small number of participants to database: %v", err)
		}
		_, err = database.AddBibChips(eventYear.Identifier, []types.BibChip{
			{
				Chip: "chip100",
				Bib:  "100",
			},
			{
				Chip: "chip106",
				Bib:  "106",
			},
			{
				Chip: "chip209",
				Bib:  "209",
			},
			{
				Chip: "chip287",
				Bib:  "287",
			},
		})
		if err != nil {
			t.Fatalf("Unexpected error adding small number of bibchips to database: %v", err)
		}
	}
	for eventKey, eventYears := range output.eventYears {
		for yearKey, eventYear := range eventYears {
			output.results[eventKey][yearKey], err = database.GetResults(eventYear.Identifier, 0, 0)
			if err != nil {
				t.Fatalf("Unexpected error getting results from database: %v", err)
			}
		}
	}
	output.sms = []types.SmsSubscription{
		{
			Bib:   "1001",
			First: "",
			Last:  "",
			Phone: "1235557890",
		},
		{
			Bib:   "",
			First: "John",
			Last:  "Smith",
			Phone: "1325557890",
		},
		{
			Bib:   "100",
			First: "",
			Last:  "",
			Phone: "1235557890",
		},
	}
	database.AddSubscribedPhone(output.eventYears["event1"]["2020"].Identifier, output.sms[1])
	database.AddSubscribedPhone(output.eventYears["event1"]["2021"].Identifier, output.sms[1])
	database.AddSubscribedPhone(output.eventYears["event2"]["2020"].Identifier, output.sms[0])
	database.AddSubscribedPhone(output.eventYears["event2"]["2020"].Identifier, output.sms[1])
	database.AddSubscribedPhone(output.eventYears["event2"]["2020"].Identifier, output.sms[2])
	database.AddSubscribedPhone(output.eventYears["event2"]["2021"].Identifier, output.sms[0])
	database.AddSubscribedPhone(output.eventYears["event2"]["2021"].Identifier, output.sms[2])
	segs := []types.Segment{
		{
			Location:      "Half Marathon",
			DistanceName:  "Marathon",
			Name:          "13.1 Miles",
			DistanceValue: 13.1,
			DistanceUnit:  "Mi",
			GPS:           "https://maps.google.com",
			MapLink:       "https://maps.google.com",
		},
		{
			Location:      "7 Mile",
			DistanceName:  "Marathon",
			Name:          "7 Miles",
			DistanceValue: 7.0,
			DistanceUnit:  "Mi",
			GPS:           "https://maps.google.com",
			MapLink:       "https://maps.google.com",
		},
		{
			Location:      "21 Mile",
			DistanceName:  "Marathon",
			Name:          "21 Miles",
			DistanceValue: 21.0,
			DistanceUnit:  "Mi",
			GPS:           "https://maps.google.com",
			MapLink:       "https://maps.google.com",
		},
		{
			Location:      "21 Mile",
			DistanceName:  "Half Marathon",
			Name:          "8 Miles",
			DistanceValue: 8.0,
			DistanceUnit:  "Mi",
			GPS:           "https://maps.google.com",
			MapLink:       "https://maps.google.com",
		},
	}
	output.segments["event1"]["2020"], _ = database.AddSegments(output.eventYears["event1"]["2020"].Identifier, segs)
	output.segments["event1"]["2021"], _ = database.AddSegments(output.eventYears["event1"]["2021"].Identifier, segs)
	output.segments["event2"]["2020"], _ = database.AddSegments(output.eventYears["event2"]["2020"].Identifier, segs)
	return output, func(t *testing.T) {
		t.Log("Deleting old database.")
		database.Close()
		// sleep here for a small period of time so the program has time to close the database and we don't default to
		// an error on the os.Remove call
		time.Sleep(200 * time.Millisecond)
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
	keys          map[string][]types.Key
	events        map[string]types.Event
	eventYears    map[string]map[string]types.EventYear
	results       map[string]map[string][]types.Result
	knownValues   map[string]string
	sms           []types.SmsSubscription
	segments      map[string]map[string][]types.Segment
}
