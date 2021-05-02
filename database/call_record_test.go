package database

import (
	"chronokeep/results/types"
	"testing"
	"time"
)

func TestAddCallRecord(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	account := types.Account{
		Name:  "John Smith",
		Email: "j@test.com",
		Type:  "admin",
	}
	account2 := types.Account{
		Name:  "Rose MacDonald",
		Email: "rose2004@test.com",
		Type:  "paid",
	}
	acc1, _ := AddAccount(account)
	acc2, _ := AddAccount(account2)
	records := []types.CallRecord{
		{
			AccountIdentifier: acc1.Identifier,
			DateTime:          time.Date(2020, 10, 5, 10, 5, 0, 0, nil).Unix(),
			Count:             15,
		},
		{
			AccountIdentifier: acc1.Identifier,
			DateTime:          time.Date(2020, 10, 5, 10, 10, 0, 0, nil).Unix(),
			Count:             500,
		},
		{
			AccountIdentifier: acc1.Identifier,
			DateTime:          time.Date(2021, 10, 5, 10, 5, 0, 0, nil).Unix(),
			Count:             1,
		},
		{
			AccountIdentifier: acc2.Identifier,
			DateTime:          time.Date(2021, 4, 5, 10, 5, 0, 0, nil).Unix(),
			Count:             150,
		},
		{
			AccountIdentifier: acc2.Identifier,
			DateTime:          time.Date(2020, 2, 5, 18, 5, 0, 0, nil).Unix(),
			Count:             70,
		},
		{
			AccountIdentifier: acc2.Identifier,
			DateTime:          time.Date(2019, 10, 5, 19, 5, 0, 0, nil).Unix(),
			Count:             35,
		},
	}
	err = AddCallRecord(records[0])
	if err != nil {
		t.Errorf("Error adding multiple call records: %v", err)
	}
	err = AddCallRecord(records[1])
	if err != nil {
		t.Errorf("Error adding multiple call records: %v", err)
	}
	err = AddCallRecord(records[2])
	if err != nil {
		t.Errorf("Error adding multiple call records: %v", err)
	}
	err = AddCallRecord(records[3])
	if err != nil {
		t.Errorf("Error adding multiple call records: %v", err)
	}
	err = AddCallRecord(records[4])
	if err != nil {
		t.Errorf("Error adding multiple call records: %v", err)
	}
	err = AddCallRecord(records[5])
	if err != nil {
		t.Errorf("Error adding multiple call records: %v", err)
	}
	records[1].Count = 700
	err = AddCallRecord(records[1])
	if err != nil {
		t.Errorf("Error adding multiple call records: %v", err)
	}
	rec, err := GetCallRecord(account.Email, records[1].DateTime)
	if err != nil {
		t.Errorf("Error getting call record: %v", err)
	}
	if rec.Count != 700 {
		t.Errorf("Expected count of %v, found %v.", 700, rec.Count)
	}
}

func TestAddCallRecords(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	account := types.Account{
		Name:  "John Smith",
		Email: "j@test.com",
		Type:  "admin",
	}
	account2 := types.Account{
		Name:  "Rose MacDonald",
		Email: "rose2004@test.com",
		Type:  "paid",
	}
	acc1, _ := AddAccount(account)
	acc2, _ := AddAccount(account2)
	records := []types.CallRecord{
		{
			AccountIdentifier: acc1.Identifier,
			DateTime:          time.Date(2020, 10, 5, 10, 5, 0, 0, nil).Unix(),
			Count:             15,
		},
		{
			AccountIdentifier: acc1.Identifier,
			DateTime:          time.Date(2020, 10, 5, 10, 10, 0, 0, nil).Unix(),
			Count:             500,
		},
		{
			AccountIdentifier: acc1.Identifier,
			DateTime:          time.Date(2021, 10, 5, 10, 5, 0, 0, nil).Unix(),
			Count:             1,
		},
		{
			AccountIdentifier: acc2.Identifier,
			DateTime:          time.Date(2021, 4, 5, 10, 5, 0, 0, nil).Unix(),
			Count:             150,
		},
		{
			AccountIdentifier: acc2.Identifier,
			DateTime:          time.Date(2020, 2, 5, 18, 5, 0, 0, nil).Unix(),
			Count:             70,
		},
		{
			AccountIdentifier: acc2.Identifier,
			DateTime:          time.Date(2019, 10, 5, 19, 5, 0, 0, nil).Unix(),
			Count:             35,
		},
	}
	err = AddCallRecords(records[0:3])
	if err != nil {
		t.Errorf("Error adding multiple call records: %v", err)
	}
	recs, _ := GetAccountCallRecords(account.Email)
	if len(recs) != 3 {
		t.Errorf("Expected %v call records, found %v.", 3, len(recs))
	}
	err = AddCallRecords(records[2:6])
	if err != nil {
		t.Errorf("Error adding multiple call records: %v", err)
	}
	recs, _ = GetAccountCallRecords(account.Email)
	if len(recs) != 3 {
		t.Errorf("Expected %v call records, found %v.", 3, len(recs))
	}
	recs, _ = GetAccountCallRecords(account2.Email)
	if len(recs) != 3 {
		t.Errorf("Expected %v call records, found %v.", 3, len(recs))
	}
}

func TestGetCallRecord(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	account := types.Account{
		Name:  "John Smith",
		Email: "j@test.com",
		Type:  "admin",
	}
	account2 := types.Account{
		Name:  "Rose MacDonald",
		Email: "rose2004@test.com",
		Type:  "paid",
	}
	acc1, _ := AddAccount(account)
	acc2, _ := AddAccount(account2)
	records := []types.CallRecord{
		{
			AccountIdentifier: acc1.Identifier,
			DateTime:          time.Date(2020, 10, 5, 10, 5, 0, 0, nil).Unix(),
			Count:             15,
		},
		{
			AccountIdentifier: acc1.Identifier,
			DateTime:          time.Date(2020, 10, 5, 10, 10, 0, 0, nil).Unix(),
			Count:             500,
		},
		{
			AccountIdentifier: acc1.Identifier,
			DateTime:          time.Date(2021, 10, 5, 10, 5, 0, 0, nil).Unix(),
			Count:             1,
		},
		{
			AccountIdentifier: acc2.Identifier,
			DateTime:          time.Date(2021, 4, 5, 10, 5, 0, 0, nil).Unix(),
			Count:             150,
		},
		{
			AccountIdentifier: acc2.Identifier,
			DateTime:          time.Date(2020, 2, 5, 18, 5, 0, 0, nil).Unix(),
			Count:             70,
		},
		{
			AccountIdentifier: acc2.Identifier,
			DateTime:          time.Date(2019, 10, 5, 19, 5, 0, 0, nil).Unix(),
			Count:             35,
		},
	}
	oldRec := records[0]
	AddCallRecord(oldRec)
	newRec, err := GetCallRecord(account.Email, time.Date(2020, 10, 5, 10, 5, 0, 0, nil).Unix())
	if err != nil {
		t.Errorf("Error getting call record: %v", err)
	}
	if *newRec != oldRec {
		t.Errorf("Expected record %v+, found %v+.", oldRec, newRec)
	}
	oldRec = records[2]
	AddCallRecord(oldRec)
	newRec, err = GetCallRecord(account.Email, time.Date(2020, 10, 5, 10, 5, 0, 0, nil).Unix())
	if err != nil {
		t.Errorf("Error getting call record: %v", err)
	}
	if *newRec != oldRec {
		t.Errorf("Expected record %v+, found %v+.", oldRec, newRec)
	}
	oldRec = records[4]
	AddCallRecord(oldRec)
	newRec, err = GetCallRecord(account2.Email, time.Date(2020, 10, 5, 10, 5, 0, 0, nil).Unix())
	if err != nil {
		t.Errorf("Error getting call record: %v", err)
	}
	if *newRec != oldRec {
		t.Errorf("Expected record %v+, found %v+.", oldRec, newRec)
	}
	oldRec = records[5]
	AddCallRecord(oldRec)
	newRec, err = GetCallRecord(account2.Email, time.Date(2020, 10, 5, 10, 5, 0, 0, nil).Unix())
	if err != nil {
		t.Errorf("Error getting call record: %v", err)
	}
	if *newRec != oldRec {
		t.Errorf("Expected record %v+, found %v+.", oldRec, newRec)
	}
}

func TestGetAccountCallRecords(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	account := types.Account{
		Name:  "John Smith",
		Email: "j@test.com",
		Type:  "admin",
	}
	account2 := types.Account{
		Name:  "Rose MacDonald",
		Email: "rose2004@test.com",
		Type:  "paid",
	}
	acc1, _ := AddAccount(account)
	acc2, _ := AddAccount(account2)
	records := []types.CallRecord{
		{
			AccountIdentifier: acc1.Identifier,
			DateTime:          time.Date(2020, 10, 5, 10, 5, 0, 0, nil).Unix(),
			Count:             15,
		},
		{
			AccountIdentifier: acc1.Identifier,
			DateTime:          time.Date(2020, 10, 5, 10, 10, 0, 0, nil).Unix(),
			Count:             500,
		},
		{
			AccountIdentifier: acc1.Identifier,
			DateTime:          time.Date(2021, 10, 5, 10, 5, 0, 0, nil).Unix(),
			Count:             1,
		},
		{
			AccountIdentifier: acc2.Identifier,
			DateTime:          time.Date(2021, 4, 5, 10, 5, 0, 0, nil).Unix(),
			Count:             150,
		},
		{
			AccountIdentifier: acc2.Identifier,
			DateTime:          time.Date(2020, 2, 5, 18, 5, 0, 0, nil).Unix(),
			Count:             70,
		},
		{
			AccountIdentifier: acc2.Identifier,
			DateTime:          time.Date(2019, 10, 5, 19, 5, 0, 0, nil).Unix(),
			Count:             35,
		},
	}
	AddCallRecords(records)
	recs, err := GetAccountCallRecords(account.Email)
	if err != nil {
		t.Errorf("Error retrieving account call records: %v", err)
	}
	if len(recs) != 3 {
		t.Errorf("Expected %v call records, found %v", 3, len(recs))
	}
	recs, err = GetAccountCallRecords(account2.Email)
	if err != nil {
		t.Errorf("Error retrieving account call records: %v", err)
	}
	if len(recs) != 3 {
		t.Errorf("Expected %v call records, found %v", 3, len(recs))
	}
}
