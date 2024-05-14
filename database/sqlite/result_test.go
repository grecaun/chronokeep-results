package sqlite

import (
	"chronokeep/results/types"
	"strconv"
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
	results = []types.Result{
		{
			PersonId:      "100",
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
			Anonymous:     true,
		},
		{
			PersonId:      "106",
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
			Anonymous:     false,
		},
		{
			PersonId:      "209",
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
			PersonId:      "287",
			Bib:           "287",
			First:         "Jamie",
			Last:          "Fischer",
			Age:           35,
			Gender:        "F",
			AgeGroup:      "30-39",
			Distance:      "5 Mile",
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
			Gender:        "F",
			AgeGroup:      "30-39",
			Distance:      "5 Mile",
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

func setupPageResultTests() {
	setupResultTests()
	results = make([]types.Result, 0)
	tmpStr := ""
	for i := 0; i < 200; i++ {
		tmpStr = strconv.Itoa(i)
		results = append(results, types.Result{
			PersonId:      tmpStr,
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
	for i := 200; i < 300; i++ {
		tmpStr = strconv.Itoa(i)
		results = append(results, types.Result{
			PersonId:      tmpStr,
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
			Finish:        false,
		})
	}
	for i := 300; i < 400; i++ {
		tmpStr = strconv.Itoa(i)
		results = append(results, types.Result{
			PersonId:      tmpStr,
			Bib:           tmpStr,
			First:         "John" + tmpStr,
			Last:          "Smith",
			Age:           24,
			Gender:        "M",
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
	}
	for i := 0; i < 200; i++ {
		tmpStr = strconv.Itoa(i)
		results = append(results, types.Result{
			PersonId:      tmpStr,
			Bib:           tmpStr,
			First:         "John" + tmpStr,
			Last:          "Smith",
			Age:           24,
			Gender:        "M",
			AgeGroup:      "20-29",
			Distance:      "1 Mile",
			Seconds:       0 + i,
			Milliseconds:  0,
			Segment:       "",
			Location:      "Start/Finish",
			Occurence:     0,
			Ranking:       i + 1,
			AgeRanking:    i + 1,
			GenderRanking: i + 1,
			Finish:        false,
		})
	}
}

func TestAddResults(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
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
	res, err := db.AddResults(eventYear.Identifier, results)
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
	res, err = db.AddResults(eventYear.Identifier, results[0:1])
	if err != nil {
		t.Fatalf("Error adding results: %v", err)
	}
	if len(res) != 1 {
		t.Errorf("Expected %v results to be added, %v added.", 1, len(res))
	}
	if res[0].First != results[0].First {
		t.Errorf("Expected to find result with first name %v, found %v.", results[0].First, res[0].First)
	}
	if res[0].Last != results[0].Last {
		t.Errorf("Expected to find result with last name %v, found %v.", results[0].Last, res[0].Last)
	}
	if res[0].Age != results[0].Age {
		t.Errorf("Expected to find result with age %v, found %v.", results[0].Age, res[0].Age)
	}
	if res[0].Gender != results[0].Gender {
		t.Errorf("Expected to find result with gender %v, found %v.", results[0].Gender, res[0].Gender)
	}
	if res[0].AgeGroup != results[0].AgeGroup {
		t.Errorf("Expected to find result with age group %v, found %v.", results[0].AgeGroup, res[0].AgeGroup)
	}
	if res[0].Distance != results[0].Distance {
		t.Errorf("Expected to find result with distance %v, found %v.", results[0].Distance, res[0].Distance)
	}
	if res[0].Seconds != results[0].Seconds {
		t.Errorf("Expected to find result with seconds %v, found %v.", results[0].Seconds, res[0].Seconds)
	}
	if res[0].Milliseconds != results[0].Milliseconds {
		t.Errorf("Expected to find result with milliseconds %v, found %v.", results[0].Milliseconds, res[0].Milliseconds)
	}
	if res[0].Segment != results[0].Segment {
		t.Errorf("Expected to find result with segment %v, found %v.", results[0].Segment, res[0].Segment)
	}
	if res[0].Ranking != results[0].Ranking {
		t.Errorf("Expected to find result with ranking %v, found %v.", results[0].Ranking, res[0].Ranking)
	}
	if res[0].AgeRanking != results[0].AgeRanking {
		t.Errorf("Expected to find result with age ranking %v, found %v.", results[0].AgeRanking, res[0].AgeRanking)
	}
	if res[0].GenderRanking != results[0].GenderRanking {
		t.Errorf("Expected to find result with gender ranking %v, found %v.", results[0].GenderRanking, res[0].GenderRanking)
	}
	if res[0].Finish != results[0].Finish {
		t.Errorf("Expected to find result with finish %v, found %v.", results[0].Finish, res[0].Finish)
	}
	if res[0].Anonymous != results[0].Anonymous {
		t.Errorf("Expected to find result with anonymous %v, found %v.", results[0].Anonymous, res[0].Anonymous)
	}
	res, _ = db.GetResults(eventYear.Identifier, 0, 0)
	if len(res) != len(results) {
		t.Errorf("Expected %v results to be added, %v added.", len(results), len(res))
	}
	// Test the update feature of AddResults. (person update, but add another result)
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
	res, err = db.AddResults(eventYear.Identifier, results[1:2])
	if err != nil {
		t.Fatalf("Error adding results: %v", err)
	}
	if len(res) != 1 {
		t.Errorf("Expected %v results to be added, %v added.", 1, len(res))
	}
	if res[0].First != results[1].First {
		t.Errorf("Expected to find result with first name %v, found %v.", results[1].First, res[0].First)
	}
	if res[0].Last != results[1].Last {
		t.Errorf("Expected to find result with last name %v, found %v.", results[1].Last, res[0].Last)
	}
	if res[0].Age != results[1].Age {
		t.Errorf("Expected to find result with age %v, found %v.", results[1].Age, res[0].Age)
	}
	if res[0].Gender != results[1].Gender {
		t.Errorf("Expected to find result with gender %v, found %v.", results[1].Gender, res[0].Gender)
	}
	if res[0].AgeGroup != results[1].AgeGroup {
		t.Errorf("Expected to find result with age group %v, found %v.", results[1].AgeGroup, res[0].AgeGroup)
	}
	if res[0].Distance != results[1].Distance {
		t.Errorf("Expected to find result with distance %v, found %v.", results[1].Distance, res[0].Distance)
	}
	if res[0].Seconds != results[1].Seconds {
		t.Errorf("Expected to find result with seconds %v, found %v.", results[1].Seconds, res[0].Seconds)
	}
	if res[0].Milliseconds != results[1].Milliseconds {
		t.Errorf("Expected to find result with milliseconds %v, found %v.", results[1].Milliseconds, res[0].Milliseconds)
	}
	if res[0].Segment != results[1].Segment {
		t.Errorf("Expected to find result with segment %v, found %v.", results[1].Segment, res[0].Segment)
	}
	if res[0].Ranking != results[1].Ranking {
		t.Errorf("Expected to find result with ranking %v, found %v.", results[1].Ranking, res[0].Ranking)
	}
	if res[0].AgeRanking != results[1].AgeRanking {
		t.Errorf("Expected to find result with age ranking %v, found %v.", results[1].AgeRanking, res[0].AgeRanking)
	}
	if res[0].GenderRanking != results[1].GenderRanking {
		t.Errorf("Expected to find result with gender ranking %v, found %v.", results[1].GenderRanking, res[0].GenderRanking)
	}
	if res[0].Finish != results[1].Finish {
		t.Errorf("Expected to find result with finish %v, found %v.", results[1].Finish, res[0].Finish)
	}
	if res[0].Anonymous != results[1].Anonymous {
		t.Errorf("Expected to find result with anonymous %v, found %v.", results[1].Anonymous, res[0].Anonymous)
	}
	res, _ = db.GetResults(eventYear.Identifier, 0, 0)
	if len(res) != (len(results) + 1) {
		t.Errorf("Expected %v results to be added, %v added.", (len(results) + 1), len(res))
	}
	// check to make sure that we actually updated the information
	// because AddResults only returns the copy that it thinks it added
	var found1 = false
	var found2 = 0
	for _, r := range res {
		if r.Equals(&results[0]) {
			found1 = true
		} else if r.SamePerson(&results[1]) {
			found2++
		}
	}
	if !found1 {
		t.Errorf("Expected to find %v %v in the results but did not.", results[0].First, results[0].Last)
	}
	if found2 < 2 {
		t.Errorf("Expected to find %v %v in the results but did not.", results[1].First, results[1].Last)
	}
	setupPageResultTests()
	_, err = db.AddResults(eventYear.Identifier, results)
	if err != nil {
		t.Fatalf("Error adding large number of results at once: %v", err)
	}
}

func TestGetLastResults(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
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
	res, err := db.GetLastResults(eventYear.Identifier, 0, 0)
	if err != nil {
		t.Fatalf("Error getting last results: %v", err)
	}
	if len(res) != 0 {
		t.Errorf("Results not added but we've found %v results.", len(res))
	}
	db.AddResults(eventYear.Identifier, results[0:1])
	res, err = db.GetLastResults(eventYear.Identifier, 0, 0)
	if err != nil {
		t.Fatalf("Error getting last results: %v", err)
	}
	if len(res) != 1 {
		t.Errorf("Expected %v results to be added, %v added.", 1, len(res))
	}
	if res[0] != results[0] {
		t.Errorf("Expected results %+v, found %+v.", results[0], res[0])
	}
	db.AddResults(eventYear.Identifier, results)
	res, err = db.GetLastResults(eventYear.Identifier, 0, 0)
	if err != nil {
		t.Fatalf("Error getting last results: %v", err)
	}
	if len(res) != (len(results) - 1) {
		t.Errorf("Expected %v results to be added, %v added.", len(results)-1, len(res))
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

func TestGetLastResultsPage(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupPageResultTests()
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
	db.AddResults(eventYear.EventIdentifier, results)
	res, err := db.GetLastResults(eventYear.EventIdentifier, 50, 0)
	if err != nil {
		t.Fatalf("Error getting first page of results: %v", err)
	}
	if len(res) != 50 {
		t.Fatalf("Expected %v results, found %v.", 50, len(res))
	}
	// these are in a weird order because of the ORDER BY part of the sql statement
	if res[0] != results[0] {
		t.Fatalf("Found %v, expected %v", res[0], results[0])
	}
	res, err = db.GetLastResults(eventYear.EventIdentifier, 50, 1)
	if err != nil {
		t.Fatalf("Error getting second page of results: %v", err)
	}
	if len(res) != 50 {
		t.Fatalf("Expected %v results, found %v.", 50, len(res))
	}
	if res[0] != results[50] {
		t.Fatalf("Found %v, expected %v", res[0], results[50])
	}
	res, err = db.GetLastResults(eventYear.EventIdentifier, 50, 7)
	if err != nil {
		t.Fatalf("Error getting second last of results: %v", err)
	}
	if len(res) != 50 {
		t.Fatalf("Expected %v results, found %v.", 50, len(res))
	}
	if res[0] != results[350] {
		t.Fatalf("Found %v, expected %v", res[0], results[350])
	}
	res, err = db.GetLastResults(eventYear.EventIdentifier, 50, 8)
	if err != nil {
		t.Fatalf("Error getting second last of results: %v", err)
	}
	if len(res) != 0 {
		t.Fatalf("Expected %v results, found %v.", 0, len(res))
	}
}

func TestGetDistanceResults(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
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
	res, err := db.GetDistanceResults(eventYear.Identifier, "test", 0, 0)
	if err != nil {
		t.Fatalf("Error getting last results: %v", err)
	}
	if len(res) != 0 {
		t.Errorf("Results not added but we've found %v results.", len(res))
	}
	db.AddResults(eventYear.Identifier, results[0:1])
	res, err = db.GetDistanceResults(eventYear.Identifier, results[0].Distance, 0, 0)
	if err != nil {
		t.Fatalf("Error getting last results: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("Expected %v results to be added, %v added.", 1, len(res))
	}
	if res[0] != results[0] {
		t.Errorf("Expected results %+v, found %+v.", results[0], res[0])
	}
	db.AddResults(eventYear.Identifier, results)
	res, err = db.GetDistanceResults(eventYear.Identifier, results[3].Distance, 0, 0)
	if err != nil {
		t.Fatalf("Error getting last results: %v", err)
	}
	if len(res) != 1 {
		t.Errorf("Expected %v results to be added, %v added.", 1, len(res))
	}
	res, err = db.GetDistanceResults(eventYear.Identifier, results[0].Distance, 0, 0)
	if err != nil {
		t.Fatalf("Error getting last results: %v", err)
	}
	if len(res) != (len(results) - 2) {
		t.Errorf("Expected %v results to be added, %v added.", len(results)-2, len(res))
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

func TestGetDistanceResultsPage(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupPageResultTests()
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
	for i := 0; i < len(results); i += 10 {
		_, err = db.AddResults(eventYear.EventIdentifier, results[i:i+10])
		if err != nil {
			t.Fatalf("Something went wrong trying to add results: %v", err)
		}
	}
	// every distance
	res, err := db.GetDistanceResults(eventYear.EventIdentifier, "", 50, 0)
	if err != nil {
		t.Fatalf("Error getting first page of results: %v", err)
	}
	if len(res) != 50 {
		t.Fatalf("Expected %v results, found %v.", 50, len(res))
	}
	if res[0] != results[0] {
		t.Fatalf("Found %v, expected %v", res[0].First, results[0].First)
	}
	res, err = db.GetDistanceResults(eventYear.EventIdentifier, "", 50, 1)
	if err != nil {
		t.Fatalf("Error getting second page of results: %v", err)
	}
	if len(res) != 50 {
		t.Fatalf("Expected %v results, found %v.", 50, len(res))
	}
	if res[0] != results[50] {
		t.Fatalf("Found %v, expected %v", res[0], results[50])
	}
	// for this one, this should return the same as GetResults
	res, err = db.GetDistanceResults(eventYear.EventIdentifier, "", 50, 7)
	if err != nil {
		t.Fatalf("Error getting second last of results: %v", err)
	}
	if len(res) != 50 {
		t.Fatalf("Expected %v results, found %v.", 50, len(res))
	}
	if res[0] != results[350] {
		t.Fatalf("Found %v, expected %v", res[0], results[350])
	}
	res, err = db.GetDistanceResults(eventYear.EventIdentifier, "", 50, 8)
	if err != nil {
		t.Fatalf("Error getting second last of results: %v", err)
	}
	if len(res) != 0 {
		t.Fatalf("Expected %v results, found %v.", 0, len(res))
	}
	// just one distance
	res, err = db.GetDistanceResults(eventYear.EventIdentifier, "1 Mile", 50, 0)
	if err != nil {
		t.Fatalf("Error getting first page of results: %v", err)
	}
	if len(res) != 50 {
		t.Fatalf("Expected %v results, found %v.", 50, len(res))
	}
	if res[0] != results[0] {
		t.Fatalf("Found %v, expected %v", res[0].First, results[0].First)
	}
	res, err = db.GetDistanceResults(eventYear.EventIdentifier, "1 Mile", 50, 1)
	if err != nil {
		t.Fatalf("Error getting second page of results: %v", err)
	}
	if len(res) != 50 {
		t.Fatalf("Expected %v results, found %v.", 50, len(res))
	}
	if res[0] != results[50] {
		t.Fatalf("Found %v, expected %v", res[0], results[50])
	}
	// this one ignores the last 100 entries
	res, err = db.GetDistanceResults(eventYear.EventIdentifier, "1 Mile", 50, 5)
	if err != nil {
		t.Fatalf("Error getting second last of results: %v", err)
	}
	if len(res) != 50 {
		t.Fatalf("Expected %v results, found %v.", 50, len(res))
	}
	if res[0] != results[250] {
		t.Fatalf("Found %v, expected %v", res[0], results[250])
	}
	res, err = db.GetDistanceResults(eventYear.EventIdentifier, "1 Mile", 50, 6)
	if err != nil {
		t.Fatalf("Error getting second last of results: %v", err)
	}
	if len(res) != 0 {
		t.Fatalf("Expected %v results, found %v.", 0, len(res))
	}
}

func TestGetAllDistanceResults(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
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
	res, err := db.GetAllDistanceResults(eventYear.Identifier, "test", 0, 0)
	if err != nil {
		t.Fatalf("Error getting last results: %v", err)
	}
	if len(res) != 0 {
		t.Errorf("Results not added but we've found %v results.", len(res))
	}
	db.AddResults(eventYear.Identifier, results[0:1])
	res, err = db.GetAllDistanceResults(eventYear.Identifier, results[0].Distance, 0, 0)
	if err != nil {
		t.Fatalf("Error getting last results: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("Expected %v results to be added, %v added.", 1, len(res))
	}
	if res[0] != results[0] {
		t.Errorf("Expected results %+v, found %+v.", results[0], res[0])
	}
	db.AddResults(eventYear.Identifier, results)
	res, err = db.GetAllDistanceResults(eventYear.Identifier, results[3].Distance, 0, 0)
	if err != nil {
		t.Fatalf("Error getting last results: %v", err)
	}
	if len(res) != 2 {
		t.Errorf("Expected %v results to be added, %v added.", 2, len(res))
	}
	res, err = db.GetAllDistanceResults(eventYear.Identifier, results[0].Distance, 0, 0)
	if err != nil {
		t.Fatalf("Error getting last results: %v", err)
	}
	if len(res) != (len(results) - 2) {
		t.Errorf("Expected %v results to be added, %v added.", len(results)-2, len(res))
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

func TestGetAllDistanceResultsPage(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupPageResultTests()
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
	for i := 0; i < len(results); i += 10 {
		_, err = db.AddResults(eventYear.EventIdentifier, results[i:i+10])
		if err != nil {
			t.Fatalf("Something went wrong trying to add results: %v", err)
		}
	}
	// every distance
	res, err := db.GetAllDistanceResults(eventYear.EventIdentifier, "", 50, 0)
	if err != nil {
		t.Fatalf("Error getting first page of results: %v", err)
	}
	if len(res) != 50 {
		t.Fatalf("Expected %v results, found %v.", 50, len(res))
	}
	if res[0] != results[400] {
		t.Fatalf("Found %v, expected %v", res[0], results[400])
	}
	res, err = db.GetAllDistanceResults(eventYear.EventIdentifier, "", 50, 1)
	if err != nil {
		t.Fatalf("Error getting second page of results: %v", err)
	}
	if len(res) != 50 {
		t.Fatalf("Expected %v results, found %v.", 50, len(res))
	}
	if res[0] != results[450] {
		t.Fatalf("Found %v, expected %v", res[0], results[50])
	}
	// for this one, this should return the same as GetResults
	res, err = db.GetAllDistanceResults(eventYear.EventIdentifier, "", 50, 7)
	if err != nil {
		t.Fatalf("Error getting second last of results: %v", err)
	}
	if len(res) != 50 {
		t.Fatalf("Expected %v results, found %v.", 50, len(res))
	}
	if res[0] != results[150] {
		t.Fatalf("Found %v, expected %v", res[0], results[150])
	}
	res, err = db.GetAllDistanceResults(eventYear.EventIdentifier, "", 50, 12)
	if err != nil {
		t.Fatalf("Error getting second last of results: %v", err)
	}
	if len(res) != 0 {
		t.Fatalf("Expected %v results, found %v.", 0, len(res))
	}
	// just one distance
	res, err = db.GetAllDistanceResults(eventYear.EventIdentifier, "1 Mile", 50, 0)
	if err != nil {
		t.Fatalf("Error getting first page of results: %v", err)
	}
	if len(res) != 50 {
		t.Fatalf("Expected %v results, found %v.", 50, len(res))
	}
	if res[0] != results[400] {
		t.Fatalf("Found %v, expected %v", res[400], results[400])
	}
	res, err = db.GetAllDistanceResults(eventYear.EventIdentifier, "1 Mile", 50, 1)
	if err != nil {
		t.Fatalf("Error getting second page of results: %v", err)
	}
	if len(res) != 50 {
		t.Fatalf("Expected %v results, found %v.", 50, len(res))
	}
	if res[0] != results[450] {
		t.Fatalf("Found %v, expected %v", res[0], results[450])
	}
	// this one ignores the last 100 entries
	res, err = db.GetAllDistanceResults(eventYear.EventIdentifier, "1 Mile", 50, 5)
	if err != nil {
		t.Fatalf("Error getting second last of results: %v", err)
	}
	if len(res) != 50 {
		t.Fatalf("Expected %v results, found %v.", 50, len(res))
	}
	if res[0] != results[50] {
		t.Fatalf("Found %v, expected %v", res[0], results[50])
	}
	res, err = db.GetAllDistanceResults(eventYear.EventIdentifier, "1 Mile", 50, 10)
	if err != nil {
		t.Fatalf("Error getting second last of results: %v", err)
	}
	if len(res) != 0 {
		t.Fatalf("Expected %v results, found %v.", 0, len(res))
	}
}

func TestGetFinishResults(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
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
	res, err := db.GetFinishResults(eventYear.Identifier, "test", 0, 0)
	if err != nil {
		t.Fatalf("Error getting finish results (1): %v", err)
	}
	if len(res) != 0 {
		t.Errorf("Results not added but we've found %v results.", len(res))
	}
	db.AddResults(eventYear.Identifier, results[0:1])
	res, err = db.GetFinishResults(eventYear.Identifier, "", 0, 0)
	if err != nil {
		t.Fatalf("Error getting finish results (2): %v", err)
	}
	if len(res) != 0 {
		t.Fatalf("Expected %v results to be added, %v added.", 0, len(res))
	}
	db.AddResults(eventYear.Identifier, results[1:2])
	res, err = db.GetFinishResults(eventYear.Identifier, "", 0, 0)
	if err != nil {
		t.Fatalf("Error getting finish results (3): %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("Expected %v results to be added, %v added.", 1, len(res))
	}
	if res[0] != results[1] {
		t.Errorf("Expected results %+v, found %+v.", results[0], res[0])
	}
	db.AddResults(eventYear.Identifier, results)
	res, err = db.GetFinishResults(eventYear.Identifier, results[0].Distance, 0, 0)
	if err != nil {
		t.Fatalf("Error getting finish results (4): %v", err)
	}
	if len(res) != (len(results) - 3) {
		t.Errorf("Expected %v results to be added, %v added.", len(results)-3, len(res))
	}
	db.AddResults(eventYear.Identifier, results)
	res, err = db.GetFinishResults(eventYear.Identifier, "", 0, 0)
	if err != nil {
		t.Fatalf("Error getting finish results (5): %v", err)
	}
	if len(res) != (len(results) - 2) {
		t.Errorf("Expected %v results to be added, %v added.", len(results)-2, len(res))
	}
	// Verify that we've got correct information for our results.
	found := false
	for _, result := range res {
		if result == results[1] {
			found = true
		}
	}
	if !found {
		t.Errorf("Unable to find our first result in the database. %+v", res)
	}
}

func TestGetFinishResultsPage(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupPageResultTests()
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
	for i := 0; i < len(results); i += 10 {
		_, err = db.AddResults(eventYear.EventIdentifier, results[i:i+10])
		if err != nil {
			t.Fatalf("Something went wrong trying to add results: %v", err)
		}
	}
	// every distance
	res, err := db.GetFinishResults(eventYear.EventIdentifier, "", 50, 0)
	if err != nil {
		t.Fatalf("Error getting first page of results: %v", err)
	}
	if len(res) != 50 {
		t.Fatalf("Expected %v results, found %v.", 50, len(res))
	}
	if res[0] != results[0] {
		t.Fatalf("Found %v, expected %v.", res[0].First, results[0].First)
	}
	res, err = db.GetFinishResults(eventYear.EventIdentifier, "", 50, 1)
	if err != nil {
		t.Fatalf("Error getting second page of results: %v", err)
	}
	if len(res) != 50 {
		t.Fatalf("Expected %v results, found %v.", 50, len(res))
	}
	if res[0] != results[50] {
		t.Fatalf("Found %v, expected %v.", res[0], results[50])
	}
	// get last page
	res, err = db.GetFinishResults(eventYear.EventIdentifier, "", 50, 5)
	if err != nil {
		t.Fatalf("Error getting second last of results: %v", err)
	}
	if len(res) != 50 {
		t.Fatalf("Expected %v results, found %v.", 50, len(res))
	}
	// there are 2 pages between that we don't get because they're not finish times, thus *1
	if res[0] != results[350] {
		t.Fatalf("Found %v, expected %v.", res[0], results[350])
	}
	res, err = db.GetFinishResults(eventYear.EventIdentifier, "", 50, 6)
	if err != nil {
		t.Fatalf("Error getting second last of results: %v", err)
	}
	if len(res) != 0 {
		t.Fatalf("Expected %v results, found %v.", 0, len(res))
	}
	// just one distance
	res, err = db.GetFinishResults(eventYear.EventIdentifier, "1 Mile", 50, 0)
	if err != nil {
		t.Fatalf("Error getting first page of results: %v", err)
	}
	if len(res) != 50 {
		t.Fatalf("Expected %v results, found %v.", 50, len(res))
	}
	if res[0] != results[0] {
		t.Fatalf("Found %v, expected %v.", res[0].First, results[0].First)
	}
	res, err = db.GetFinishResults(eventYear.EventIdentifier, "1 Mile", 50, 1)
	if err != nil {
		t.Fatalf("Error getting second page of results: %v", err)
	}
	if len(res) != 50 {
		t.Fatalf("Expected %v results, found %v.", 50, len(res))
	}
	if res[0] != results[50] {
		t.Fatalf("Found %v, expected %v.", res[0], results[50])
	}
	res, err = db.GetFinishResults(eventYear.EventIdentifier, "1 Mile", 50, 3)
	if err != nil {
		t.Fatalf("Error getting second last of results: %v", err)
	}
	if len(res) != 50 {
		t.Fatalf("Expected %v results, found %v.", 50, len(res))
	}
	// unlike above, this ignores all of the last results which are 2 Mile not 1 Mile distances
	if res[0] != results[150] {
		t.Fatalf("Found %v, expected %v.", res[0], results[150])
	}
	res, err = db.GetFinishResults(eventYear.EventIdentifier, "1 Mile", 50, 4)
	if err != nil {
		t.Fatalf("Error getting second last of results: %v", err)
	}
	if len(res) != 0 {
		t.Fatalf("Expected %v results, found %v.", 0, len(res))
	}
}

func TestGetResults(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
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
	res, err := db.GetResults(eventYear.Identifier, 0, 0)
	if err != nil {
		t.Fatalf("Error getting results: %v", err)
	}
	if len(res) != 0 {
		t.Errorf("Results not added but we've found %v results.", len(res))
	}
	db.AddResults(eventYear.Identifier, results[0:1])
	res, err = db.GetResults(eventYear.Identifier, 0, 0)
	if err != nil {
		t.Fatalf("Error getting results: %v", err)
	}
	if len(res) != 1 {
		t.Errorf("Expected %v results to be added, %v added.", 1, len(res))
	}
	if res[0] != results[0] {
		t.Errorf("Expected results %+v, found %+v.", results[0], res[0])
	}
	db.AddResults(eventYear.Identifier, results)
	res, err = db.GetResults(eventYear.Identifier, 0, 0)
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

func TestGetResultsPage(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupPageResultTests()
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
	for i := 0; i < len(results); i += 10 {
		_, err = db.AddResults(eventYear.EventIdentifier, results[i:i+10])
		if err != nil {
			t.Fatalf("Something went wrong trying to add results: %v", err)
		}
	}
	res, err := db.GetResults(eventYear.EventIdentifier, 50, 0)
	if err != nil {
		t.Fatalf("Error getting first page of results: %v", err)
	}
	if len(res) != 50 {
		t.Fatalf("Expected %v results, found %v.", 50, len(res))
	}
	// these are in a weird order because of the ORDER BY part of the sql statement
	if res[0] != results[400] {
		t.Fatalf("Found %v, expected %v", res[0], results[400])
	}
	res, err = db.GetResults(eventYear.EventIdentifier, 50, 1)
	if err != nil {
		t.Fatalf("Error getting second page of results: %v", err)
	}
	if len(res) != 50 {
		t.Fatalf("Expected %v results, found %v.", 50, len(res))
	}
	if res[0] != results[450] {
		t.Fatalf("Found %v, expected %v", res[0], results[450])
	}
	res, err = db.GetResults(eventYear.EventIdentifier, 50, len(results)/50-1)
	if err != nil {
		t.Fatalf("Error getting second last of results: %v", err)
	}
	if len(res) != 50 {
		t.Fatalf("Expected %v results, found %v.", 50, len(res))
	}
	if res[0] != results[350] {
		t.Fatalf("Found %v, expected %v", res[0], results[350])
	}
	res, err = db.GetResults(eventYear.EventIdentifier, 50, len(results)/50)
	if err != nil {
		t.Fatalf("Error getting second last of results: %v", err)
	}
	if len(res) != 0 {
		t.Fatalf("Expected %v results, found %v.", 0, len(res))
	}
}

func TestPageNoLimit(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupPageResultTests()
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
	db.AddResults(eventYear.EventIdentifier, results)
	res, err := db.GetResults(eventYear.EventIdentifier, 0, 50)
	if err != nil {
		t.Fatalf("Error getting first page of results: %v", err)
	}
	if len(res) != len(results) {
		t.Fatalf("Expected %v results, found %v.", len(results), len(res))
	}
	res, err = db.GetDistanceResults(eventYear.EventIdentifier, "", 0, 50)
	if err != nil {
		t.Fatalf("Error getting first page of results: %v", err)
	}
	if len(res) != len(results)-200 {
		t.Fatalf("Expected %v results, found %v.", len(results)-200, len(res))
	}
	res, err = db.GetDistanceResults(eventYear.EventIdentifier, "1 Mile", 0, 50)
	if err != nil {
		t.Fatalf("Error getting first page of results: %v", err)
	}
	if len(res) != len(results)-300 {
		t.Fatalf("Expected %v results, found %v.", len(results)-300, len(res))
	}
	res, err = db.GetFinishResults(eventYear.EventIdentifier, "", 0, 50)
	if err != nil {
		t.Fatalf("Error getting first page of results: %v", err)
	}
	if len(res) != len(results)-300 {
		t.Fatalf("Expected %v results, found %v.", len(results)-300, len(res))
	}
	res, err = db.GetFinishResults(eventYear.EventIdentifier, "1 Mile", 0, 50)
	if err != nil {
		t.Fatalf("Error getting first page of results: %v", err)
	}
	if len(res) != len(results)-400 {
		t.Fatalf("Expected %v results, found %v.", len(results)-400, len(res))
	}
	res, err = db.GetLastResults(eventYear.EventIdentifier, 0, 50)
	if err != nil {
		t.Fatalf("Error getting first page of results: %v", err)
	}
	if len(res) != len(results)-200 {
		t.Fatalf("Expected %v results, found %v.", len(results)-200, len(res))
	}
}

func TestDeleteResults(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
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
	eventYear2 := &types.EventYear{
		EventIdentifier: event.Identifier,
		Year:            "2022",
		DateTime:        time.Date(2022, 04, 20, 9, 0, 0, 0, time.Local),
	}
	eventYear, _ = db.AddEventYear(*eventYear)
	eventYear2, _ = db.AddEventYear(*eventYear2)
	db.AddResults(eventYear.Identifier, results)
	db.AddResults(eventYear2.Identifier, results)
	err = db.DeleteResults(eventYear.Identifier, results[1:2])
	if err != nil {
		t.Fatalf("Error deleting specific results: %v", err)
	}
	res, _ := db.GetResults(eventYear.Identifier, 0, 0)
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
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupResultTests()
	account, _ := db.AddAccount(accounts[0])
	event := &types.Event{
		AccountIdentifier: account.Identifier,
		Name:              "Event 1",
		Slug:              "event1",
	}
	event2 := &types.Event{
		AccountIdentifier: account.Identifier,
		Name:              "Event 2",
		Slug:              "event2",
	}
	event, _ = db.AddEvent(*event)
	event2, _ = db.AddEvent(*event2)
	eventYear := &types.EventYear{
		EventIdentifier: event.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 04, 20, 9, 0, 0, 0, time.Local),
	}
	eventYear2 := &types.EventYear{
		EventIdentifier: event2.Identifier,
		Year:            "2022",
		DateTime:        time.Date(2022, 04, 20, 9, 0, 0, 0, time.Local),
	}
	eventYear, _ = db.AddEventYear(*eventYear)
	eventYear2, _ = db.AddEventYear(*eventYear2)
	db.AddResults(eventYear.Identifier, results)
	db.AddResults(eventYear2.Identifier, results)
	count, err := db.DeleteEventResults(eventYear.Identifier)
	if err != nil {
		t.Fatalf("Error deleting specific results: %v", err)
	}
	if (count) != int64(len(results)) {
		t.Errorf("Expected to find out %v results were deleted, %v were deleted.", len(results), count)
	}
	res, _ := db.GetResults(eventYear.Identifier, 0, 0)
	if len(res) != 0 {
		t.Errorf("Expected %v results after delete but found %v.", 0, len(res))
	}
	res, _ = db.GetResults(eventYear2.Identifier, 0, 0)
	if len(res) == 0 {
		t.Error("Expected to find results after delete but found none.")
	}
}

func TestGetBibResults(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
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
	lastIX := len(results) - 1
	res, err := db.GetBibResults(eventYear.Identifier, results[lastIX].Bib)
	if err != nil {
		t.Fatalf("Error getting bib results: %v", err)
	}
	if len(res) != 0 {
		t.Errorf("Expected %v results to be added, %v added.", 0, len(res))
	}
	db.AddResults(eventYear.Identifier, results[0:lastIX])
	res, err = db.GetBibResults(eventYear.Identifier, results[lastIX].Bib)
	if err != nil {
		t.Fatalf("Error getting bib results: %v", err)
	}
	if len(res) != 1 {
		t.Errorf("Expected %v results to be added, %v added.", 1, len(res))
	}
	db.AddResults(eventYear.Identifier, results)
	res, err = db.GetBibResults(eventYear.Identifier, results[lastIX].Bib)
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

func TestBadDatabaseResult(t *testing.T) {
	db := badTestSetup(t)
	_, err := db.GetResults(0, 0, 0)
	if err == nil {
		t.Fatalf("Expected error getting results by event year.")
	}
	_, err = db.GetLastResults(0, 0, 0)
	if err == nil {
		t.Fatalf("Expected error getting last results by event year.")
	}
	_, err = db.GetFinishResults(0, "", 0, 0)
	if err == nil {
		t.Fatalf("Expected error getting finish results by event year.")
	}
	_, err = db.GetDistanceResults(0, "", 0, 0)
	if err == nil {
		t.Fatalf("Expected error getting distance results by event year.")
	}
	_, err = db.GetBibResults(0, "fake bib")
	if err == nil {
		t.Fatalf("Expected error getting results by event year and bib.")
	}
	err = db.DeleteResults(0, make([]types.Result, 0))
	if err == nil {
		t.Fatalf("Expected error deleting results.")
	}
	_, err = db.DeleteEventResults(0)
	if err == nil {
		t.Fatalf("Expected error deleting event year results.")
	}
	_, err = db.AddResults(0, make([]types.Result, 0))
	if err == nil {
		t.Fatalf("Expected error adding results.")
	}
}

func TestNoDatabaseResult(t *testing.T) {
	db := SQLite{}
	_, err := db.GetResults(0, 0, 0)
	if err == nil {
		t.Fatalf("Expected error getting results by event year.")
	}
	_, err = db.GetLastResults(0, 0, 0)
	if err == nil {
		t.Fatalf("Expected error getting last results by event year.")
	}
	_, err = db.GetFinishResults(0, "", 0, 0)
	if err == nil {
		t.Fatalf("Expected error getting finish results by event year.")
	}
	_, err = db.GetDistanceResults(0, "", 0, 0)
	if err == nil {
		t.Fatalf("Expected error getting distance results by event year.")
	}
	_, err = db.GetBibResults(0, "fake bib")
	if err == nil {
		t.Fatalf("Expected error getting results by event year and bib.")
	}
	err = db.DeleteResults(0, make([]types.Result, 0))
	if err == nil {
		t.Fatalf("Expected error deleting results.")
	}
	_, err = db.DeleteEventResults(0)
	if err == nil {
		t.Fatalf("Expected error deleting event year results.")
	}
	_, err = db.AddResults(0, make([]types.Result, 0))
	if err == nil {
		t.Fatalf("Expected error adding results.")
	}
}
