package mysql

import (
	"chronokeep/results/auth"
	"chronokeep/results/types"
	"testing"
	"time"
)

var (
	accounts      []types.Account
	testPassword1 string = "password"
	testPassword2 string = "newpassword"
)

func setupAccountTests() {
	if len(accounts) < 1 {
		accounts = []types.Account{
			{
				Name:     "John Smith",
				Email:    "j@test.com",
				Type:     "admin",
				Password: testHashPassword(testPassword1),
			},
			{
				Name:     "Jerry Garcia",
				Email:    "jgarcia@test.com",
				Type:     "free",
				Password: testHashPassword(testPassword1),
			},
			{
				Name:     "Rose MacDonald",
				Email:    "rose2004@test.com",
				Type:     "paid",
				Password: testHashPassword(testPassword1),
			},
			{
				Name:     "Tia Johnson",
				Email:    "tiatheway@test.com",
				Type:     "free",
				Password: testHashPassword(testPassword1),
			},
			{
				Name:     "Thomas Donaldson",
				Email:    "tdon@test.com",
				Type:     "admin",
				Password: testHashPassword(testPassword1),
			},
			{
				Name:     "Ester White",
				Email:    "white@test.com",
				Type:     "test",
				Password: testHashPassword(testPassword1),
			},
			{
				Name:     "Ricky Reagan",
				Email:    "rreagan@test.com",
				Type:     "free",
				Password: testHashPassword(testPassword1),
			},
		}
	}
}

func TestAddAccount(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	// Ensure adding accounts works properly.
	t.Log("Adding accounts")
	setupAccountTests()
	oAccount := accounts[0]
	nAccount, err := db.AddAccount(oAccount)
	if err != nil {
		t.Fatalf("Error adding account: %v", err)
	}
	t.Logf("New account ID: %v", nAccount.Identifier)
	if !oAccount.Equals(nAccount) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", oAccount, *nAccount)
	}
	oAccount = accounts[1]
	nAccount, err = db.AddAccount(oAccount)
	if err != nil {
		t.Fatalf("Error adding account: %v", err)
	}
	t.Logf("New account ID: %v", nAccount.Identifier)
	if !oAccount.Equals(nAccount) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", oAccount, *nAccount)
	}
	oAccount = accounts[2]
	nAccount, err = db.AddAccount(oAccount)
	if err != nil {
		t.Fatalf("Error adding account: %v", err)
	}
	t.Logf("New account ID: %v", nAccount.Identifier)
	if !oAccount.Equals(nAccount) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", oAccount, *nAccount)
	}
	oAccount = accounts[3]
	nAccount, err = db.AddAccount(oAccount)
	if err != nil {
		t.Fatalf("Error adding account: %v", err)
	}
	t.Logf("New account ID: %v", nAccount.Identifier)
	if !oAccount.Equals(nAccount) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", oAccount, *nAccount)
	}
	// Test for collisions.
	_, err = db.AddAccount(accounts[2])
	if err == nil {
		t.Error("Expected error adding account with duplicate email.")
	}
}

func TestGetAccount(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupAccountTests()
	// Test getting known accounts.
	oAccount := accounts[0]
	nAccount, err := db.AddAccount(oAccount)
	if err != nil {
		t.Fatalf("Error adding account: %v", err)
	}
	t.Logf("New account ID: %v", nAccount.Identifier)
	dAccount, err := db.GetAccount(oAccount.Email)
	if err != nil {
		t.Fatalf("Error getting account: %v", err)
	}
	if !dAccount.Equals(nAccount) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", *nAccount, *dAccount)
	}
	if dAccount.Identifier != nAccount.Identifier {
		t.Errorf("Account id expected to be %v but found %v.", nAccount.Identifier, dAccount.Identifier)
	}
	oAccount = accounts[1]
	nAccount, err = db.AddAccount(oAccount)
	if err != nil {
		t.Fatalf("Error adding account: %v", err)
	}
	t.Logf("New account ID: %v", nAccount.Identifier)
	dAccount, err = db.GetAccount(oAccount.Email)
	if err != nil {
		t.Fatalf("Error getting account: %v", err)
	}
	if !dAccount.Equals(nAccount) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", *nAccount, *dAccount)
	}
	if dAccount.Identifier != nAccount.Identifier {
		t.Errorf("Account id expected to be %v but found %v.", nAccount.Identifier, dAccount.Identifier)
	}
	oAccount = accounts[2]
	nAccount, err = db.AddAccount(oAccount)
	if err != nil {
		t.Fatalf("Error adding account: %v", err)
	}
	t.Logf("New account ID: %v", nAccount.Identifier)
	dAccount, err = db.GetAccount(oAccount.Email)
	if err != nil {
		t.Fatalf("Error getting account: %v", err)
	}
	if !dAccount.Equals(nAccount) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", *nAccount, *dAccount)
	}
	if dAccount.Identifier != nAccount.Identifier {
		t.Errorf("Account id expected to be %v but found %v.", nAccount.Identifier, dAccount.Identifier)
	}
	oAccount = accounts[3]
	nAccount, err = db.AddAccount(oAccount)
	if err != nil {
		t.Fatalf("Error adding account: %v", err)
	}
	t.Logf("New account ID: %v", nAccount.Identifier)
	dAccount, err = db.GetAccount(oAccount.Email)
	if err != nil {
		t.Fatalf("Error getting account: %v", err)
	}
	if !dAccount.Equals(nAccount) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", *nAccount, *dAccount)
	}
	if dAccount.Identifier != nAccount.Identifier {
		t.Errorf("Account id expected to be %v but found %v.", nAccount.Identifier, dAccount.Identifier)
	}
	// Test getting unknown accounts.
	dAccount, err = db.GetAccount("random@test.com")
	if err != nil {
		t.Fatalf("Error finding account not in existence: %v", err)
	}
	if dAccount != nil {
		t.Error("Expected not to find an account but one was found.")
	}
}

func TestGetAccounts(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupAccountTests()
	retAccounts, err := db.GetAccounts()
	if err != nil {
		t.Fatalf("Error getting accounts: %v", err)
	}
	if len(retAccounts) != 0 {
		t.Errorf("Expected number of accounts is %v but %v were found.", 0, len(retAccounts))
	}
	db.AddAccount(accounts[0])
	db.AddAccount(accounts[1])
	db.AddAccount(accounts[2])
	retAccounts, err = db.GetAccounts()
	if err != nil {
		t.Fatalf("Error getting accounts: %v", err)
	}
	if len(retAccounts) != 3 {
		t.Errorf("Expected number of accounts is %v but %v were found.", 3, len(retAccounts))
	}
	db.AddAccount(accounts[3])
	db.AddAccount(accounts[4])
	db.AddAccount(accounts[5])
	db.AddAccount(accounts[6])
	retAccounts, err = db.GetAccounts()
	if err != nil {
		t.Fatalf("Error getting accounts: %v", err)
	}
	if len(retAccounts) != 7 {
		t.Errorf("Expected number of accounts is %v but %v were found.", 7, len(retAccounts))
	}
}

func TestUpdateAccount(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupAccountTests()
	// Ensure adding accounts works properly.
	nAccount, err := db.AddAccount(accounts[0])
	if err != nil {
		t.Fatalf("Error adding account: %v", err)
	}
	nAccount.Name = "New Name 1"
	err = db.UpdateAccount(*nAccount)
	if err != nil {
		t.Fatalf("Error updating account: %v", err)
	}
	dAccount, _ := db.GetAccount(nAccount.Email)
	if nAccount.Identifier != dAccount.Identifier {
		t.Errorf("Account ID expected to be %v but found %v instead.", nAccount.Identifier, dAccount.Identifier)
	}
	if dAccount.Name != "New Name 1" {
		t.Errorf("Account name expected to be %v but found %v instead.", "New Name 1", dAccount.Name)
	}
	nAccount, err = db.AddAccount(accounts[1])
	if err != nil {
		t.Fatalf("Error adding account: %v", err)
	}
	nAccount.Type = "New Type 1"
	err = db.UpdateAccount(*nAccount)
	if err != nil {
		t.Fatalf("Error updating account: %v", err)
	}
	dAccount, _ = db.GetAccount(nAccount.Email)
	if nAccount.Identifier != dAccount.Identifier {
		t.Errorf("Account ID expected to be %v but found %v instead.", nAccount.Identifier, dAccount.Identifier)
	}
	if dAccount.Type != "New Type 1" {
		t.Errorf("Account name expected to be %v but found %v instead.", "New Type 1", dAccount.Type)
	}
	nAccount, err = db.AddAccount(accounts[2])
	if err != nil {
		t.Fatalf("Error adding account: %v", err)
	}
	nAccount.Name = "New Name 2"
	err = db.UpdateAccount(*nAccount)
	dAccount, _ = db.GetAccount(nAccount.Email)
	if err != nil {
		t.Fatalf("Error updating account: %v", err)
	}
	if nAccount.Identifier != dAccount.Identifier {
		t.Errorf("Account ID expected to be %v but found %v instead.", nAccount.Identifier, dAccount.Identifier)
	}
	if dAccount.Name != "New Name 2" {
		t.Errorf("Account name expected to be %v but found %v instead.", "New Name 2", dAccount.Name)
	}
	nAccount, err = db.AddAccount(accounts[3])
	if err != nil {
		t.Fatalf("Error adding account: %v", err)
	}
	nAccount.Type = "New Type 2"
	err = db.UpdateAccount(*nAccount)
	if err != nil {
		t.Fatalf("Error updating account: %v", err)
	}
	dAccount, _ = db.GetAccount(nAccount.Email)
	if nAccount.Identifier != dAccount.Identifier {
		t.Errorf("Account ID expected to be %v but found %v instead.", nAccount.Identifier, dAccount.Identifier)
	}
	if dAccount.Type != "New Type 2" {
		t.Errorf("Account name expected to be %v but found %v instead.", "New Type 2", dAccount.Type)
	}
	// Test for collisions.
	_, err = db.AddAccount(accounts[2])
	if err == nil {
		t.Error("Expected error adding account with duplicate email.")
	}
}

func TestDeleteAccount(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupAccountTests()
	nAccount, _ := db.AddAccount(accounts[0])
	keys := []types.Key{
		{
			AccountIdentifier: nAccount.Identifier,
			Value:             "030001-1ACSDD-K2389A-00123B",
			Type:              "default",
			AllowedHosts:      "",
			ValidUntil:        time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
		},
		{
			AccountIdentifier: nAccount.Identifier,
			Value:             "030001-1ACSDD-K2389A-22123B",
			Type:              "write",
			AllowedHosts:      "",
			ValidUntil:        time.Now().Add(time.Hour * 20).Truncate(time.Second),
		},
	}
	db.AddKey(keys[0])
	db.AddKey(keys[1])
	err = db.DeleteAccount(nAccount.Identifier)
	if err != nil {
		t.Fatalf("Error deleting account: %v", err)
	}
	dAccount, _ := db.GetAccount(nAccount.Email)
	if dAccount != nil {
		t.Error("Unexpectedly found a deleted account.")
	}
	keys, _ = db.GetAccountKeys(nAccount.Email)
	if len(keys) != 0 {
		t.Errorf("expected to find %v keys after deleting account, found %v", 0, len(keys))
	}
	_, err = db.AddAccount(accounts[0])
	if err == nil {
		t.Error("No error found when trying to add a deleted account.")
	}
	nAccount, _ = db.AddAccount(accounts[1])
	err = db.DeleteAccount(nAccount.Identifier)
	if err != nil {
		t.Fatalf("Error deleting account: %v", err)
	}
	dAccount, _ = db.GetAccount(nAccount.Email)
	if dAccount != nil {
		t.Error("Unexpectedly found a deleted account.")
	}
	_, err = db.AddAccount(accounts[1])
	if err == nil {
		t.Error("No error found when trying to add a deleted account.")
	}
	nAccount, _ = db.AddAccount(accounts[2])
	err = db.DeleteAccount(nAccount.Identifier)
	if err != nil {
		t.Fatalf("Error deleting account: %v", err)
	}
	dAccount, _ = db.GetAccount(nAccount.Email)
	if dAccount != nil {
		t.Error("Unexpectedly found a deleted account.")
	}
	_, err = db.AddAccount(accounts[2])
	if err == nil {
		t.Error("No error found when trying to add a deleted account.")
	}
}

func TestResurrectAccount(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupAccountTests()
	nAccount, _ := db.AddAccount(accounts[0])
	db.DeleteAccount(nAccount.Identifier)
	err = db.ResurrectAccount(nAccount.Email)
	if err != nil {
		t.Fatalf("Error resurrecting account: %v", err)
	}
	dAccount, _ := db.GetAccount(nAccount.Email)
	if dAccount == nil {
		t.Error("Account was not resurrected.")
	}
	nAccount, _ = db.AddAccount(accounts[1])
	db.DeleteAccount(nAccount.Identifier)
	err = db.ResurrectAccount(nAccount.Email)
	if err != nil {
		t.Fatalf("Error resurrecting account: %v", err)
	}
	dAccount, _ = db.GetAccount(nAccount.Email)
	if dAccount == nil {
		t.Error("Account was not resurrected.")
	}
	nAccount, _ = db.AddAccount(accounts[4])
	db.DeleteAccount(nAccount.Identifier)
	err = db.ResurrectAccount(nAccount.Email)
	if err != nil {
		t.Fatalf("Error resurrecting account: %v", err)
	}
	dAccount, _ = db.GetAccount(nAccount.Email)
	if dAccount == nil {
		t.Error("Account was not resurrected.")
	}
}

func TestGetDeletedAccount(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupAccountTests()
	nAccount, _ := db.AddAccount(accounts[0])
	db.DeleteAccount(nAccount.Identifier)
	dAccount, err := db.GetDeletedAccount(nAccount.Email)
	if err != nil {
		t.Fatalf("Error getting deleted account %v", err)
	}
	if dAccount == nil {
		t.Error("Deleted account not found.")
	}
	nAccount, _ = db.AddAccount(accounts[3])
	db.DeleteAccount(nAccount.Identifier)
	dAccount, err = db.GetDeletedAccount(nAccount.Email)
	if err != nil {
		t.Fatalf("Error getting deleted account %v", err)
	}
	if dAccount == nil {
		t.Error("Deleted account not found.")
	}
	nAccount, _ = db.AddAccount(accounts[5])
	db.DeleteAccount(nAccount.Identifier)
	dAccount, err = db.GetDeletedAccount(nAccount.Email)
	if err != nil {
		t.Fatalf("Error getting deleted account %v", err)
	}
	if dAccount == nil {
		t.Error("Deleted account not found.")
	}
	nAccount, _ = db.AddAccount(accounts[6])
	db.DeleteAccount(nAccount.Identifier)
	dAccount, err = db.GetDeletedAccount(nAccount.Email)
	if err != nil {
		t.Fatalf("Error getting deleted account %v", err)
	}
	if dAccount == nil {
		t.Error("Deleted account not found.")
	}
}

func TestChangePassword(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupAccountTests()
	nAccount, _ := db.AddAccount(accounts[0])
	hashPass, _ := auth.HashPassword(testPassword2)
	err = db.ChangePassword(nAccount.Email, hashPass)
	if err != nil {
		t.Fatalf("error changing password: %v", err)
	}
	nAccount, _ = db.GetAccount(nAccount.Email)
	err = auth.VerifyPassword(nAccount.Password, testPassword2)
	if err != nil {
		t.Errorf("password doesn't match: %v", err)
	}
}

func TestChangeEmail(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupAccountTests()
	nAccount, _ := db.AddAccount(accounts[0])
	newEmail := "new_email2020@test.com"
	err = db.ChangeEmail(nAccount.Email, newEmail)
	if err != nil {
		t.Fatalf("error changing email: %v", err)
	}
	nAccount, _ = db.GetAccount(nAccount.Email)
	if nAccount != nil {
		t.Errorf("account retrieved when the email should have changed: %v", nAccount)
	}
	nAccount, _ = db.GetAccount(newEmail)
	if nAccount == nil {
		t.Error("account with new email not found")
	}
}

func TestInvalidPassword(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupAccountTests()
	nAccount, _ := db.AddAccount(accounts[0])
	for i := 1; i <= MaxLoginAttempts+3; i++ {
		err = db.InvalidPassword(*nAccount)
		if err != nil {
			t.Fatalf("(%v) error telling the database about an invalid password: %v", i, err)
		}
		nAccount, _ = db.GetAccount(nAccount.Email)
		if nAccount.WrongPassAttempts > MaxLoginAttempts && nAccount.Locked == false {
			t.Errorf("account is not locked after (%v) invalid password attempts; should be after (%v)", i, MaxLoginAttempts+1)
		} else if nAccount.WrongPassAttempts <= MaxLoginAttempts && nAccount.Locked == true {
			t.Errorf("account is locked after (%v) invalid password attempts; should be (%v)", i, MaxLoginAttempts+1)
		}
		if nAccount.WrongPassAttempts != i {
			t.Errorf("wrong password attempts set to %v, should be %v", nAccount.WrongPassAttempts, i)
		}
	}
}

func TestGetAccountByKey(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupAccountTests()
	// Test getting known accounts.
	nAccount1, _ := db.AddAccount(accounts[0])
	nAccount2, _ := db.AddAccount(accounts[1])
	keys := []types.Key{
		{
			AccountIdentifier: nAccount1.Identifier,
			Value:             "030001-1ACSDD-K2389A-00123B",
			Type:              "default",
			AllowedHosts:      "",
			ValidUntil:        time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
		},
		{
			AccountIdentifier: nAccount2.Identifier,
			Value:             "030001-1ACSDD-K2389A-22123B",
			Type:              "write",
			AllowedHosts:      "",
			ValidUntil:        time.Now().Add(time.Hour * 20).Truncate(time.Second),
		},
	}
	db.AddKey(keys[0])
	db.AddKey(keys[1])
	dAccount, err := db.GetAccountByKey(keys[0].Value)
	if err != nil {
		t.Fatalf("Error getting account: %v", err)
	}
	if !dAccount.Equals(nAccount1) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", *nAccount1, *dAccount)
	}
	if dAccount.Identifier != nAccount1.Identifier {
		t.Errorf("Account id expected to be %v but found %v.", nAccount1.Identifier, dAccount.Identifier)
	}
	dAccount, err = db.GetAccountByKey(keys[1].Value)
	if err != nil {
		t.Fatalf("Error getting account: %v", err)
	}
	if !dAccount.Equals(nAccount2) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", *nAccount2, *dAccount)
	}
	if dAccount.Identifier != nAccount2.Identifier {
		t.Errorf("Account id expected to be %v but found %v.", nAccount2.Identifier, dAccount.Identifier)
	}
}

func TestGetAccountByID(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupAccountTests()
	// Test getting known accounts.
	oAccount := accounts[0]
	nAccount, _ := db.AddAccount(oAccount)
	dAccount, err := db.GetAccountByID(nAccount.Identifier)
	if err != nil {
		t.Fatalf("Error getting account: %v", err)
	}
	if !dAccount.Equals(nAccount) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", *nAccount, *dAccount)
	}
	if dAccount.Identifier != nAccount.Identifier {
		t.Errorf("Account id expected to be %v but found %v.", nAccount.Identifier, dAccount.Identifier)
	}
	oAccount = accounts[1]
	nAccount, _ = db.AddAccount(oAccount)
	dAccount, err = db.GetAccountByID(nAccount.Identifier)
	if err != nil {
		t.Fatalf("Error getting account: %v", err)
	}
	if !dAccount.Equals(nAccount) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", *nAccount, *dAccount)
	}
	if dAccount.Identifier != nAccount.Identifier {
		t.Errorf("Account id expected to be %v but found %v.", nAccount.Identifier, dAccount.Identifier)
	}
}
