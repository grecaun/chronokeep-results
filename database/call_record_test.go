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
	account1 := types.Account{
		Name:  "John Smith",
		Email: "j@test.com",
		Type:  "admin",
	}
	account2 := types.Account{
		Name:  "Rose MacDonald",
		Email: "rose2004@test.com",
		Type:  "paid",
	}
	acc1, _ := AddAccount(account1)
	acc2, _ := AddAccount(account2)
	records := []types.CallRecord{
		{
			AccountIdentifier: acc1.Identifier,
			DateTime:          time.Date(2020, 10, 5, 10, 5, 0, 0, time.Local).Unix(),
			Count:             15,
		},
		{
			AccountIdentifier: acc1.Identifier,
			DateTime:          time.Date(2020, 10, 5, 10, 10, 0, 0, time.Local).Unix(),
			Count:             500,
		},
		{
			AccountIdentifier: acc1.Identifier,
			DateTime:          time.Date(2021, 10, 5, 10, 5, 0, 0, time.Local).Unix(),
			Count:             1,
		},
		{
			AccountIdentifier: acc2.Identifier,
			DateTime:          time.Date(2021, 4, 5, 10, 5, 0, 0, time.Local).Unix(),
			Count:             150,
		},
		{
			AccountIdentifier: acc2.Identifier,
			DateTime:          time.Date(2020, 2, 5, 18, 5, 0, 0, time.Local).Unix(),
			Count:             70,
		},
		{
			AccountIdentifier: acc2.Identifier,
			DateTime:          time.Date(2019, 10, 5, 19, 5, 0, 0, time.Local).Unix(),
			Count:             35,
		},
	}
	err = AddCallRecord(records[0])
	if err != nil {
		t.Fatalf("Error adding call record: %v", err)
	}
	err = AddCallRecord(records[1])
	if err != nil {
		t.Fatalf("Error adding call record: %v", err)
	}
	err = AddCallRecord(records[2])
	if err != nil {
		t.Fatalf("Error adding call record: %v", err)
	}
	err = AddCallRecord(records[3])
	if err != nil {
		t.Fatalf("Error adding call record: %v", err)
	}
	err = AddCallRecord(records[4])
	if err != nil {
		t.Fatalf("Error adding call record: %v", err)
	}
	err = AddCallRecord(records[5])
	if err != nil {
		t.Fatalf("Error adding call record: %v", err)
	}
	records[1].Count = 700
	err = AddCallRecord(records[1])
	if err != nil {
		t.Fatalf("Error adding call record: %v", err)
	}
	rec, err := GetCallRecord(account1.Email, records[1].DateTime)
	if err != nil {
		t.Fatalf("Error getting call record: %v", err)
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
	account1 := types.Account{
		Name:  "John Smith",
		Email: "j@test.com",
		Type:  "admin",
	}
	account2 := types.Account{
		Name:  "Rose MacDonald",
		Email: "rose2004@test.com",
		Type:  "paid",
	}
	acc1, _ := AddAccount(account1)
	acc2, _ := AddAccount(account2)
	records := []types.CallRecord{
		{
			AccountIdentifier: acc1.Identifier,
			DateTime:          time.Date(2020, 10, 5, 10, 5, 0, 0, time.Local).Unix(),
			Count:             15,
		},
		{
			AccountIdentifier: acc1.Identifier,
			DateTime:          time.Date(2020, 10, 5, 10, 10, 0, 0, time.Local).Unix(),
			Count:             500,
		},
		{
			AccountIdentifier: acc1.Identifier,
			DateTime:          time.Date(2021, 10, 5, 10, 5, 0, 0, time.Local).Unix(),
			Count:             1,
		},
		{
			AccountIdentifier: acc2.Identifier,
			DateTime:          time.Date(2021, 4, 5, 10, 5, 0, 0, time.Local).Unix(),
			Count:             150,
		},
		{
			AccountIdentifier: acc2.Identifier,
			DateTime:          time.Date(2020, 2, 5, 18, 5, 0, 0, time.Local).Unix(),
			Count:             70,
		},
		{
			AccountIdentifier: acc2.Identifier,
			DateTime:          time.Date(2019, 10, 5, 19, 5, 0, 0, time.Local).Unix(),
			Count:             35,
		},
	}
	err = AddCallRecords(records[0:3])
	if err != nil {
		t.Fatalf("Error adding multiple call records: %v", err)
	}
	recs, _ := GetAccountCallRecords(account1.Email)
	if len(recs) != 3 {
		t.Errorf("Expected %v call records, found %v.", 3, len(recs))
	}
	err = AddCallRecords(records[2:6])
	if err != nil {
		t.Fatalf("Error adding multiple call records: %v", err)
	}
	recs, _ = GetAccountCallRecords(account1.Email)
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
	account1 := types.Account{
		Name:  "John Smith",
		Email: "j@test.com",
		Type:  "admin",
	}
	account2 := types.Account{
		Name:  "Rose MacDonald",
		Email: "rose2004@test.com",
		Type:  "paid",
	}
	acc1, _ := AddAccount(account1)
	acc2, _ := AddAccount(account2)
	records := []types.CallRecord{
		{
			AccountIdentifier: acc1.Identifier,
			DateTime:          time.Date(2020, 10, 5, 10, 5, 0, 0, time.Local).Unix(),
			Count:             15,
		},
		{
			AccountIdentifier: acc1.Identifier,
			DateTime:          time.Date(2020, 10, 5, 10, 10, 0, 0, time.Local).Unix(),
			Count:             500,
		},
		{
			AccountIdentifier: acc1.Identifier,
			DateTime:          time.Date(2021, 10, 5, 10, 5, 0, 0, time.Local).Unix(),
			Count:             1,
		},
		{
			AccountIdentifier: acc2.Identifier,
			DateTime:          time.Date(2021, 4, 5, 10, 5, 0, 0, time.Local).Unix(),
			Count:             150,
		},
		{
			AccountIdentifier: acc2.Identifier,
			DateTime:          time.Date(2020, 2, 5, 18, 5, 0, 0, time.Local).Unix(),
			Count:             70,
		},
		{
			AccountIdentifier: acc2.Identifier,
			DateTime:          time.Date(2019, 10, 5, 19, 5, 0, 0, time.Local).Unix(),
			Count:             35,
		},
	}
	t.Log("Record 1")
	oldRec := records[0]
	AddCallRecord(oldRec)
	newRec, err := GetCallRecord(account1.Email, oldRec.DateTime)
	if err != nil {
		t.Fatalf("Error getting call record: %v", err)
	}
	if *newRec != oldRec {
		t.Errorf("Expected record %+v, found %+v.", oldRec, *newRec)
	}
	t.Log("Record 2")
	oldRec = records[2]
	AddCallRecord(oldRec)
	newRec, err = GetCallRecord(account1.Email, oldRec.DateTime)
	if err != nil {
		t.Fatalf("Error getting call record: %v", err)
	}
	if *newRec != oldRec {
		t.Errorf("Expected record %+v, found %+v.", oldRec, *newRec)
	}
	t.Log("Record 3")
	oldRec = records[4]
	AddCallRecord(oldRec)
	newRec, err = GetCallRecord(account2.Email, oldRec.DateTime)
	if err != nil {
		t.Fatalf("Error getting call record: %v", err)
	}
	if *newRec != oldRec {
		t.Errorf("Expected record %+v, found %+v.", oldRec, *newRec)
	}
	t.Log("Record 4")
	oldRec = records[5]
	AddCallRecord(oldRec)
	newRec, err = GetCallRecord(account2.Email, oldRec.DateTime)
	if err != nil {
		t.Fatalf("Error getting call record: %v", err)
	}
	if *newRec != oldRec {
		t.Errorf("Expected record %+v, found %+v.", oldRec, *newRec)
	}
	t.Log("Record 5 (empty)")
	newRec, err = GetCallRecord(account2.Email, time.Date(2020, 10, 5, 10, 35, 0, 0, time.Local).Unix())
	if err != nil {
		t.Fatalf("Error getting call record: %v", err)
	}
	if newRec != nil {
		t.Errorf("Found an unexpected call record: %+v", *newRec)
	}
}

func TestGetAccountCallRecords(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	account1 := types.Account{
		Name:  "John Smith",
		Email: "j@test.com",
		Type:  "admin",
	}
	account2 := types.Account{
		Name:  "Rose MacDonald",
		Email: "rose2004@test.com",
		Type:  "paid",
	}
	account3 := types.Account{
		Name:  "Tia Johnson",
		Email: "tiatheway@test.com",
		Type:  "free",
	}
	acc1, _ := AddAccount(account1)
	acc2, _ := AddAccount(account2)
	AddAccount(account3)
	records := []types.CallRecord{
		{
			AccountIdentifier: acc1.Identifier,
			DateTime:          time.Date(2020, 10, 5, 10, 5, 0, 0, time.Local).Unix(),
			Count:             15,
		},
		{
			AccountIdentifier: acc1.Identifier,
			DateTime:          time.Date(2020, 10, 5, 10, 10, 0, 0, time.Local).Unix(),
			Count:             500,
		},
		{
			AccountIdentifier: acc1.Identifier,
			DateTime:          time.Date(2021, 10, 5, 10, 5, 0, 0, time.Local).Unix(),
			Count:             1,
		},
		{
			AccountIdentifier: acc2.Identifier,
			DateTime:          time.Date(2021, 4, 5, 10, 5, 0, 0, time.Local).Unix(),
			Count:             150,
		},
		{
			AccountIdentifier: acc2.Identifier,
			DateTime:          time.Date(2020, 2, 5, 18, 5, 0, 0, time.Local).Unix(),
			Count:             70,
		},
		{
			AccountIdentifier: acc2.Identifier,
			DateTime:          time.Date(2019, 10, 5, 19, 5, 0, 0, time.Local).Unix(),
			Count:             35,
		},
	}
	AddCallRecords(records)
	recs, err := GetAccountCallRecords(account1.Email)
	if err != nil {
		t.Fatalf("Error retrieving account1 call records: %v", err)
	}
	if len(recs) != 3 {
		t.Errorf("Expected %v call records, found %v", 3, len(recs))
	}
	recs, err = GetAccountCallRecords(account2.Email)
	if err != nil {
		t.Fatalf("Error retrieving account1 call records: %v", err)
	}
	if len(recs) != 3 {
		t.Errorf("Expected %v call records, found %v", 3, len(recs))
	}
	recs, err = GetAccountCallRecords(account3.Email)
	if err != nil {
		t.Fatalf("Error retrieving account1 call records: %v", err)
	}
	if len(recs) != 0 {
		t.Errorf("Expected %v call records, found %v", 0, len(recs))
	}
}
