package postgres

import (
	"chronokeep/results/auth"
	"chronokeep/results/database"
	"chronokeep/results/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	accounts      []types.Account
	subAccounts   []types.Account
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
	if len(subAccounts) < 1 {
		subAccounts = []types.Account{
			{
				Name:     "Registration1",
				Email:    "registration1@test.com",
				Type:     "registration",
				Password: testHashPassword(testPassword1),
			},
			{
				Name:     "Registration2",
				Email:    "registration2@test.com",
				Type:     "registration",
				Password: testHashPassword(testPassword1),
			},
			{
				Name:     "Registration3",
				Email:    "registration3@test.com",
				Type:     "registration",
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
	times := []time.Time{
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
		time.Now().Add(time.Hour * 20).Truncate(time.Second),
	}
	keys := []types.Key{
		{
			AccountIdentifier: nAccount.Identifier,
			Value:             "030001-1ACSDD-K2389A-00123B",
			Type:              "default",
			AllowedHosts:      "",
			ValidUntil:        &times[0],
		},
		{
			AccountIdentifier: nAccount.Identifier,
			Value:             "030001-1ACSDD-K2389A-22123B",
			Type:              "write",
			AllowedHosts:      "",
			ValidUntil:        &times[1],
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
	dAccount, err := db.GetDeletedAccount(nAccount.Email)
	if err != nil {
		t.Fatalf("Error getting deleted account %v", err)
	}
	if dAccount != nil {
		t.Errorf("Deleted account found: %v", nAccount.Email)
	}
	db.DeleteAccount(nAccount.Identifier)
	dAccount, err = db.GetDeletedAccount(nAccount.Email)
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
	if nAccount == nil {
		t.Fatal("get account failure")
	}
	err = auth.VerifyPassword(nAccount.Password, testPassword2)
	if err != nil {
		t.Errorf("password doesn't match: %v", err)
	}
	nAccount.Token = "testToken1"
	nAccount.RefreshToken = "testToken2"
	_ = db.UpdateTokens(*nAccount)
	nAccount, _ = db.GetAccount(nAccount.Email)
	if nAccount.Token == "" || nAccount.RefreshToken == "" {
		t.Error("Expected tokens to be set.")
	}
	err = db.ChangePassword(nAccount.Email, hashPass, true)
	if err != nil {
		t.Fatalf("error changing password: %v", err)
	}
	nAccount, _ = db.GetAccount(nAccount.Email)
	if nAccount.Token != "" || nAccount.RefreshToken != "" {
		t.Errorf("Expected tokens not to be set. Found Token %v and Refresh Token %v.", nAccount.Token, nAccount.RefreshToken)
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
	nAccount.Token = "testToken1"
	nAccount.RefreshToken = "testToken2"
	_ = db.UpdateTokens(*nAccount)
	nAccount, _ = db.GetAccount(nAccount.Email)
	if nAccount.Token == "" || nAccount.RefreshToken == "" {
		t.Error("Expected tokens to be set.")
	}
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
	} else if nAccount.Token != "" || nAccount.RefreshToken != "" {
		t.Errorf("Expected tokens not to be set. Found Token %v and Refresh Token %v.", nAccount.Token, nAccount.RefreshToken)
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
	nAccount.Token = "testToken1"
	nAccount.RefreshToken = "testToken2"
	_ = db.UpdateTokens(*nAccount)
	nAccount, _ = db.GetAccount(nAccount.Email)
	var dAccount *types.Account
	if nAccount.Token == "" || nAccount.RefreshToken == "" {
		t.Error("Expected tokens to be set.")
	}
	for i := 1; i <= database.MaxLoginAttempts+3; i++ {
		err = db.InvalidPassword(*nAccount)
		if err != nil {
			t.Fatalf("(%v) error telling the database about an invalid password: %v", i, err)
		}
		dAccount, _ = db.GetAccount(nAccount.Email)
		if dAccount.WrongPassAttempts > database.MaxLoginAttempts && dAccount.Locked == false {
			t.Errorf("account is not locked after (%v) invalid password attempts; should be after (%v)", i, database.MaxLoginAttempts+1)
			if dAccount.Token != "" || dAccount.RefreshToken != "" {
				t.Errorf("Expected tokens not to be set. Found Token %v and Refresh Token %v.", dAccount.Token, dAccount.RefreshToken)
			}
		} else if dAccount.WrongPassAttempts <= database.MaxLoginAttempts && dAccount.Locked == true {
			t.Errorf("account is locked after (%v) invalid password attempts; should be (%v)", i, database.MaxLoginAttempts+1)
			if dAccount.Token == "" || dAccount.RefreshToken == "" {
				t.Error("Expected tokens to be set.")
			}
		}
		if dAccount.WrongPassAttempts != i {
			t.Errorf("wrong password attempts set to %v, should be %v", dAccount.WrongPassAttempts, i)
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
	times := []time.Time{
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
		time.Now().Add(time.Hour * 20).Truncate(time.Second),
	}
	keys := []types.Key{
		{
			AccountIdentifier: nAccount1.Identifier,
			Value:             "030001-1ACSDD-K2389A-00123B",
			Type:              "default",
			AllowedHosts:      "",
			ValidUntil:        &times[0],
		},
		{
			AccountIdentifier: nAccount2.Identifier,
			Value:             "030001-1ACSDD-K2389A-22123B",
			Type:              "write",
			AllowedHosts:      "",
			ValidUntil:        &times[1],
		},
	}
	db.AddKey(keys[0])
	db.AddKey(keys[1])
	dAccount, err := db.GetAccountByKey(keys[0].Value)
	if err != nil {
		t.Fatalf("Error getting account: %v", err)
	}
	if dAccount == nil {
		t.Fatalf("Account not found. (1)")
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
	if dAccount == nil {
		t.Fatalf("Account not found. (2)")
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

func TestUpdateTokens(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupAccountTests()
	oAccount := accounts[0]
	nAccount, _ := db.AddAccount(oAccount)
	nAccount.Token = "testtoken1"
	nAccount.RefreshToken = "refreshtoken1"
	err = db.UpdateTokens(*nAccount)
	if err != nil {
		t.Fatalf("Error updating tokens: %v", err)
	}
	dAccount, _ := db.GetAccount(nAccount.Email)
	if dAccount.Token != nAccount.Token {
		t.Errorf("Expected token %v, found %v.", nAccount.Token, dAccount.Token)
	}
	if dAccount.RefreshToken != nAccount.RefreshToken {
		t.Errorf("Expected refresh token %v, found %v.", nAccount.RefreshToken, dAccount.RefreshToken)
	}
}

func TestValidPassword(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupAccountTests()
	nAccount, _ := db.AddAccount(accounts[0])
	for i := 1; i <= database.MaxLoginAttempts-2; i++ {
		err = db.InvalidPassword(*nAccount)
		if err != nil {
			t.Fatalf("(%v) error telling the database about an invalid password: %v", i, err)
		}
	}
	nAccount, _ = db.GetAccount(nAccount.Email)
	if nAccount.WrongPassAttempts < 1 {
		t.Errorf("Expected more than 1 wrong pass attempts; found %v.", nAccount.WrongPassAttempts)
	}
	err = db.ValidPassword(*nAccount)
	if err != nil {
		t.Fatalf("Valid password threw an error: %v", err)
	}
	nAccount, _ = db.GetAccount(nAccount.Email)
	if nAccount.WrongPassAttempts != 0 {
		t.Errorf("Expected zero wrong pass attempts; found %v.", nAccount.WrongPassAttempts)
	}
	// Test to make sure we don't unlock if locked.
	for i := 1; i <= database.MaxLoginAttempts+3; i++ {
		err = db.InvalidPassword(*nAccount)
		if err != nil {
			t.Fatalf("(%v) error telling the database about an invalid password: %v", i, err)
		}
	}
	err = db.ValidPassword(*nAccount)
	if err == nil {
		t.Fatal("Expected an error on valid password attempt for locked account.")
	}
	nAccount, _ = db.GetAccount(nAccount.Email)
	if nAccount.WrongPassAttempts == 0 {
		t.Errorf("Expected wrong password attempts; found %v.", nAccount.WrongPassAttempts)
	}
}

func TestUnlockAccount(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupAccountTests()
	nAccount, _ := db.AddAccount(accounts[0])
	// Should throw error if account isn't locked
	err = db.UnlockAccount(*nAccount)
	if err == nil {
		t.Fatal("no error thrown on unlock of unlocked account")
	}
	for i := 1; i <= database.MaxLoginAttempts+3; i++ {
		err = db.InvalidPassword(*nAccount)
		if err != nil {
			t.Fatalf("(%v) error telling the database about an invalid password: %v", i, err)
		}
	}
	nAccount, _ = db.GetAccount(nAccount.Email)
	err = db.UnlockAccount(*nAccount)
	if err != nil {
		t.Fatalf("Unexpected error on unlock account: %v", err)
	}
	nAccount, _ = db.GetAccount(nAccount.Email)
	if nAccount.WrongPassAttempts != 0 {
		t.Errorf("Expected wrong pass attempts to be reset to 0; found %v.", nAccount.WrongPassAttempts)
	}
}

func TestInternalAccount(t *testing.T) {
	// test nil, nil, nil for get internal account
	db, finalize, _ := setupTests(t)
	defer finalize(t)
	_, err := db.getAccountInternal(nil, nil, nil)
	if err == nil {
		t.Fatalf("Expected error getting account internal with no values given.")
	}
}

func TestBadDatabaseAccount(t *testing.T) {
	// test bad database connection
	db := badTestSetup(t)
	_, err := db.GetAccount("")
	if err == nil {
		t.Fatalf("Expected error getting account by email.")
	}
	_, err = db.GetAccountByKey("")
	if err == nil {
		t.Fatalf("Expected error getting account by key.")
	}
	_, err = db.GetAccountByID(0)
	if err == nil {
		t.Fatalf("Expected error getting account by id.")
	}
	_, err = db.GetAccounts()
	if err == nil {
		t.Fatalf("Expected error getting accounts.")
	}
	_, err = db.AddAccount(types.Account{})
	if err == nil {
		t.Fatalf("Expected error adding empty account.")
	}
	err = db.DeleteAccount(0)
	if err == nil {
		t.Fatalf("Expected error deleting account.")
	}
	err = db.ResurrectAccount("")
	if err == nil {
		t.Fatalf("Expected error resurrecting account.")
	}
	_, err = db.GetDeletedAccount("")
	if err == nil {
		t.Fatalf("Expected error getting deleted account.")
	}
	err = db.UpdateAccount(types.Account{})
	if err == nil {
		t.Fatalf("Expected error updating account.")
	}
	err = db.ChangePassword("", "", true)
	if err == nil {
		t.Fatalf("Expected error changing password.")
	}
	err = db.UpdateTokens(types.Account{})
	if err == nil {
		t.Fatalf("Expected error updating tokens.")
	}
	err = db.ChangeEmail("", "")
	if err == nil {
		t.Fatalf("Expected error changing email.")
	}
	err = db.InvalidPassword(types.Account{})
	if err == nil {
		t.Fatalf("Expected error setting invalid password.")
	}
	err = db.ValidPassword(types.Account{})
	if err == nil {
		t.Fatalf("Expected error setting valid password.")
	}
	err = db.UnlockAccount(types.Account{})
	if err == nil {
		t.Fatalf("Expected error unlocking account.")
	}
	_, err = db.getAccountInternal(nil, nil, nil)
	if err == nil {
		t.Fatalf("Expected error getting account internal with no values given.")
	}
}

func TestNoDatabaseAccount(t *testing.T) {
	// test whether or not we've connected to a database
	db := Postgres{}
	_, err := db.GetAccount("")
	if err == nil {
		t.Fatalf("Expected error getting account by email.")
	}
	_, err = db.GetAccountByKey("")
	if err == nil {
		t.Fatalf("Expected error getting account by key.")
	}
	_, err = db.GetAccountByID(0)
	if err == nil {
		t.Fatalf("Expected error getting account by id.")
	}
	_, err = db.GetAccounts()
	if err == nil {
		t.Fatalf("Expected error getting accounts.")
	}
	_, err = db.AddAccount(types.Account{})
	if err == nil {
		t.Fatalf("Expected error adding empty account.")
	}
	err = db.DeleteAccount(0)
	if err == nil {
		t.Fatalf("Expected error deleting account.")
	}
	err = db.ResurrectAccount("")
	if err == nil {
		t.Fatalf("Expected error resurrecting account.")
	}
	_, err = db.GetDeletedAccount("")
	if err == nil {
		t.Fatalf("Expected error getting deleted account.")
	}
	err = db.UpdateAccount(types.Account{})
	if err == nil {
		t.Fatalf("Expected error updating account.")
	}
	err = db.ChangePassword("", "", true)
	if err == nil {
		t.Fatalf("Expected error changing password.")
	}
	err = db.UpdateTokens(types.Account{})
	if err == nil {
		t.Fatalf("Expected error updating tokens.")
	}
	err = db.ChangeEmail("", "")
	if err == nil {
		t.Fatalf("Expected error changing email.")
	}
	err = db.InvalidPassword(types.Account{})
	if err == nil {
		t.Fatalf("Expected error setting invalid password.")
	}
	err = db.ValidPassword(types.Account{})
	if err == nil {
		t.Fatalf("Expected error setting valid password.")
	}
	err = db.UnlockAccount(types.Account{})
	if err == nil {
		t.Fatalf("Expected error unlocking account.")
	}
	_, err = db.getAccountInternal(nil, nil, nil)
	if err == nil {
		t.Fatalf("Expected error getting account internal with no values given.")
	}
}

func TestGetLinkedAccounts(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	t.Log("Adding accounts")
	setupAccountTests()
	acc1, err := db.AddAccount(accounts[0])
	assert.NoError(t, err)
	acc2, err := db.AddAccount(accounts[1])
	assert.NoError(t, err)
	subAcc1, err := db.AddAccount(subAccounts[0])
	assert.NoError(t, err)
	subAcc2, err := db.AddAccount(subAccounts[1])
	assert.NoError(t, err)
	subAcc3, err := db.AddAccount(subAccounts[2])
	assert.NoError(t, err)
	// verify empty
	a, err := db.GetLinkedAccounts(acc1.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(a))
	}
	// test add valid
	err = db.LinkAccounts(*acc1, *subAcc1)
	assert.NoError(t, err)
	a, err = db.GetLinkedAccounts(acc1.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(a))
	}
	// test add more
	err = db.LinkAccounts(*acc1, *subAcc2)
	assert.NoError(t, err)
	err = db.LinkAccounts(*acc1, *subAcc3)
	assert.NoError(t, err)
	a, err = db.GetLinkedAccounts(acc1.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 3, len(a))
	}
	// verify empty
	a, err = db.GetLinkedAccounts(acc2.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(a))
	}
	// test add to second
	err = db.LinkAccounts(*acc2, *subAcc1)
	assert.NoError(t, err)
	a, err = db.GetLinkedAccounts(acc2.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(a))
	}
}

func TestLinkAccounts(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	t.Log("Adding accounts")
	setupAccountTests()
	acc1, err := db.AddAccount(accounts[0])
	assert.NoError(t, err)
	acc2, err := db.AddAccount(accounts[1])
	assert.NoError(t, err)
	subAcc1, err := db.AddAccount(subAccounts[0])
	assert.NoError(t, err)
	subAcc2, err := db.AddAccount(subAccounts[1])
	assert.NoError(t, err)
	subAcc3, err := db.AddAccount(subAccounts[2])
	assert.NoError(t, err)
	// verify empty
	a, err := db.GetLinkedAccounts(acc1.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(a))
	}
	// test add valid
	err = db.LinkAccounts(*acc1, *subAcc1)
	assert.NoError(t, err)
	a, err = db.GetLinkedAccounts(acc1.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(a))
	}
	// test re-add
	err = db.LinkAccounts(*acc1, *subAcc1)
	assert.NoError(t, err)
	a, err = db.GetLinkedAccounts(acc1.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(a))
	}
	// test add more
	err = db.LinkAccounts(*acc1, *subAcc2)
	assert.NoError(t, err)
	err = db.LinkAccounts(*acc1, *subAcc3)
	assert.NoError(t, err)
	a, err = db.GetLinkedAccounts(acc1.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 3, len(a))
	}
	// test invalid add - invalid main account
	err = db.LinkAccounts(*subAcc1, *subAcc2)
	assert.Error(t, err)
	// test invalid add - invalid sub account
	err = db.LinkAccounts(*acc2, *acc1)
	assert.Error(t, err)
	// Verify we can't link the same account to itself
	err = db.LinkAccounts(*acc2, *acc2)
	assert.Error(t, err)
	err = db.LinkAccounts(*subAcc1, *subAcc1)
	assert.Error(t, err)
	// verify empty
	a, err = db.GetLinkedAccounts(acc2.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(a))
	}
	// test add to second
	err = db.LinkAccounts(*acc2, *subAcc1)
	assert.NoError(t, err)
	a, err = db.GetLinkedAccounts(acc2.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(a))
	}
}

func TestUnlinkAccounts(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	t.Log("Adding accounts")
	setupAccountTests()
	acc1, err := db.AddAccount(accounts[0])
	assert.NoError(t, err)
	acc2, err := db.AddAccount(accounts[1])
	assert.NoError(t, err)
	subAcc1, err := db.AddAccount(subAccounts[0])
	assert.NoError(t, err)
	subAcc2, err := db.AddAccount(subAccounts[1])
	assert.NoError(t, err)
	subAcc3, err := db.AddAccount(subAccounts[2])
	assert.NoError(t, err)
	// verify empty
	a, err := db.GetLinkedAccounts(acc1.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(a))
	}
	// test add valid
	err = db.LinkAccounts(*acc1, *subAcc1)
	assert.NoError(t, err)
	err = db.LinkAccounts(*acc1, *subAcc2)
	assert.NoError(t, err)
	err = db.LinkAccounts(*acc1, *subAcc3)
	assert.NoError(t, err)
	a, err = db.GetLinkedAccounts(acc1.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 3, len(a))
	}
	// Test unlink
	err = db.UnlinkAccounts(*acc1, *subAcc1)
	assert.NoError(t, err)
	a, err = db.GetLinkedAccounts(acc1.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 2, len(a))
	}
	// Test re-do of unlink
	err = db.UnlinkAccounts(*acc1, *subAcc1)
	assert.NoError(t, err)
	a, err = db.GetLinkedAccounts(acc1.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 2, len(a))
	}
	// Test re-do of unlink
	err = db.UnlinkAccounts(*acc1, *subAcc2)
	assert.NoError(t, err)
	err = db.UnlinkAccounts(*acc1, *subAcc3)
	assert.NoError(t, err)
	a, err = db.GetLinkedAccounts(acc1.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(a))
	}
	// verify empty
	a, err = db.GetLinkedAccounts(acc2.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(a))
	}
	// test incorrect main
	err = db.UnlinkAccounts(*subAcc2, *subAcc2)
	assert.NoError(t, err)
	// test add to second
	err = db.LinkAccounts(*acc2, *subAcc1)
	assert.NoError(t, err)
	a, err = db.GetLinkedAccounts(acc2.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(a))
	}
	// test unlink second
	err = db.UnlinkAccounts(*acc2, *subAcc1)
	assert.NoError(t, err)
	a, err = db.GetLinkedAccounts(acc2.Email)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(a))
	}
}
