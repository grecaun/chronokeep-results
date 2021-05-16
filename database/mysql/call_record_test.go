package mysql

import (
	"chronokeep/results/types"
	"testing"
	"time"
)

func setupCallRecordTests() {
	if len(accounts) < 1 {
		accounts = []types.Account{
			{
				Name:     "John Smith",
				Email:    "j@test.com",
				Type:     "admin",
				Password: testHashPassword("password"),
			},
			{
				Name:     "Rose MacDonald",
				Email:    "rose2004@test.com",
				Type:     "paid",
				Password: testHashPassword("password"),
			},
			{
				Name:     "Tia Johnson",
				Email:    "tiatheway@test.com",
				Type:     "free",
				Password: testHashPassword("password"),
			},
		}
	}
}

func TestAddCallRecord(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupCallRecordTests()
	acc1, _ := db.AddAccount(accounts[0])
	acc2, _ := db.AddAccount(accounts[1])
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
	err = db.AddCallRecord(records[0])
	if err != nil {
		t.Fatalf("Error adding call record: %v", err)
	}
	err = db.AddCallRecord(records[1])
	if err != nil {
		t.Fatalf("Error adding call record: %v", err)
	}
	err = db.AddCallRecord(records[2])
	if err != nil {
		t.Fatalf("Error adding call record: %v", err)
	}
	err = db.AddCallRecord(records[3])
	if err != nil {
		t.Fatalf("Error adding call record: %v", err)
	}
	err = db.AddCallRecord(records[4])
	if err != nil {
		t.Fatalf("Error adding call record: %v", err)
	}
	err = db.AddCallRecord(records[5])
	if err != nil {
		t.Fatalf("Error adding call record: %v", err)
	}
	records[1].Count = 700
	err = db.AddCallRecord(records[1])
	if err != nil {
		t.Fatalf("Error adding call record: %v", err)
	}
	rec, err := db.GetCallRecord(accounts[0].Email, records[1].DateTime)
	if err != nil {
		t.Fatalf("Error getting call record: %v", err)
	}
	if rec.Count != 700 {
		t.Errorf("Expected count of %v, found %v.", 700, rec.Count)
	}
}

func TestAddCallRecords(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupCallRecordTests()
	acc1, _ := db.AddAccount(accounts[0])
	acc2, _ := db.AddAccount(accounts[1])
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
	err = db.AddCallRecords(records[0:3])
	if err != nil {
		t.Fatalf("Error adding multiple call records: %v", err)
	}
	recs, _ := db.GetAccountCallRecords(accounts[0].Email)
	if len(recs) != 3 {
		t.Errorf("Expected %v call records, found %v.", 3, len(recs))
	}
	err = db.AddCallRecords(records[2:6])
	if err != nil {
		t.Fatalf("Error adding multiple call records: %v", err)
	}
	recs, _ = db.GetAccountCallRecords(accounts[0].Email)
	if len(recs) != 3 {
		t.Errorf("Expected %v call records, found %v.", 3, len(recs))
	}
	recs, _ = db.GetAccountCallRecords(accounts[0].Email)
	if len(recs) != 3 {
		t.Errorf("Expected %v call records, found %v.", 3, len(recs))
	}
}

func TestGetCallRecord(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupCallRecordTests()
	acc1, _ := db.AddAccount(accounts[0])
	acc2, _ := db.AddAccount(accounts[1])
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
	db.AddCallRecord(oldRec)
	newRec, err := db.GetCallRecord(accounts[0].Email, oldRec.DateTime)
	if err != nil {
		t.Fatalf("Error getting call record: %v", err)
	}
	if *newRec != oldRec {
		t.Errorf("Expected record %+v, found %+v.", oldRec, *newRec)
	}
	t.Log("Record 2")
	oldRec = records[2]
	db.AddCallRecord(oldRec)
	newRec, err = db.GetCallRecord(accounts[0].Email, oldRec.DateTime)
	if err != nil {
		t.Fatalf("Error getting call record: %v", err)
	}
	if *newRec != oldRec {
		t.Errorf("Expected record %+v, found %+v.", oldRec, *newRec)
	}
	t.Log("Record 3")
	oldRec = records[4]
	db.AddCallRecord(oldRec)
	newRec, err = db.GetCallRecord(accounts[1].Email, oldRec.DateTime)
	if err != nil {
		t.Fatalf("Error getting call record: %v", err)
	}
	if *newRec != oldRec {
		t.Errorf("Expected record %+v, found %+v.", oldRec, *newRec)
	}
	t.Log("Record 4")
	oldRec = records[5]
	db.AddCallRecord(oldRec)
	newRec, err = db.GetCallRecord(accounts[1].Email, oldRec.DateTime)
	if err != nil {
		t.Fatalf("Error getting call record: %v", err)
	}
	if *newRec != oldRec {
		t.Errorf("Expected record %+v, found %+v.", oldRec, *newRec)
	}
	t.Log("Record 5 (empty)")
	newRec, err = db.GetCallRecord(accounts[1].Email, time.Date(2020, 10, 5, 10, 35, 0, 0, time.Local).Unix())
	if err != nil {
		t.Fatalf("Error getting call record: %v", err)
	}
	if newRec != nil {
		t.Errorf("Found an unexpected call record: %+v", *newRec)
	}
}

func TestGetAccountCallRecords(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupCallRecordTests()
	acc1, _ := db.AddAccount(accounts[0])
	acc2, _ := db.AddAccount(accounts[1])
	db.AddAccount(accounts[2])
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
	db.AddCallRecords(records)
	recs, err := db.GetAccountCallRecords(accounts[0].Email)
	if err != nil {
		t.Fatalf("Error retrieving account1 call records: %v", err)
	}
	if len(recs) != 3 {
		t.Errorf("Expected %v call records, found %v", 3, len(recs))
	}
	recs, err = db.GetAccountCallRecords(accounts[1].Email)
	if err != nil {
		t.Fatalf("Error retrieving account1 call records: %v", err)
	}
	if len(recs) != 3 {
		t.Errorf("Expected %v call records, found %v", 3, len(recs))
	}
	recs, err = db.GetAccountCallRecords(accounts[2].Email)
	if err != nil {
		t.Fatalf("Error retrieving account1 call records: %v", err)
	}
	if len(recs) != 0 {
		t.Errorf("Expected %v call records, found %v", 0, len(recs))
	}
}
