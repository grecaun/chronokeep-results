package database

import (
	"chronokeep/results/auth"
	"chronokeep/results/types"
	"testing"
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
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	// Ensure adding accounts works properly.
	t.Log("Adding accounts")
	setupAccountTests()
	oAccount := accounts[0]
	nAccount, err := AddAccount(oAccount)
	if err != nil {
		t.Fatalf("Error adding account: %v", err)
	}
	t.Logf("New account ID: %v", nAccount.Identifier)
	if !oAccount.Equals(nAccount) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", oAccount, *nAccount)
	}
	oAccount = accounts[1]
	nAccount, err = AddAccount(oAccount)
	if err != nil {
		t.Fatalf("Error adding account: %v", err)
	}
	t.Logf("New account ID: %v", nAccount.Identifier)
	if !oAccount.Equals(nAccount) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", oAccount, *nAccount)
	}
	oAccount = accounts[2]
	nAccount, err = AddAccount(oAccount)
	if err != nil {
		t.Fatalf("Error adding account: %v", err)
	}
	t.Logf("New account ID: %v", nAccount.Identifier)
	if !oAccount.Equals(nAccount) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", oAccount, *nAccount)
	}
	oAccount = accounts[3]
	nAccount, err = AddAccount(oAccount)
	if err != nil {
		t.Fatalf("Error adding account: %v", err)
	}
	t.Logf("New account ID: %v", nAccount.Identifier)
	if !oAccount.Equals(nAccount) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", oAccount, *nAccount)
	}
	// Test for collisions.
	_, err = AddAccount(accounts[2])
	if err == nil {
		t.Error("Expected error adding account with duplicate email.")
	}
}

func TestGetAccount(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupAccountTests()
	// Test getting known accounts.
	oAccount := accounts[0]
	nAccount, err := AddAccount(oAccount)
	if err != nil {
		t.Fatalf("Error adding account: %v", err)
	}
	t.Logf("New account ID: %v", nAccount.Identifier)
	dAccount, err := GetAccount(oAccount.Email)
	if err != nil {
		t.Fatalf("Error getting account: %v", err)
	}
	if !dAccount.Equals(nAccount) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", oAccount, *nAccount)
	}
	if dAccount.Identifier != nAccount.Identifier {
		t.Errorf("Account id expected to be %v but found %v.", dAccount.Identifier, nAccount.Identifier)
	}
	oAccount = accounts[1]
	nAccount, err = AddAccount(oAccount)
	if err != nil {
		t.Fatalf("Error adding account: %v", err)
	}
	t.Logf("New account ID: %v", nAccount.Identifier)
	dAccount, err = GetAccount(oAccount.Email)
	if err != nil {
		t.Fatalf("Error getting account: %v", err)
	}
	if !dAccount.Equals(nAccount) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", oAccount, *nAccount)
	}
	if dAccount.Identifier != nAccount.Identifier {
		t.Errorf("Account id expected to be %v but found %v.", dAccount.Identifier, nAccount.Identifier)
	}
	oAccount = accounts[2]
	nAccount, err = AddAccount(oAccount)
	if err != nil {
		t.Fatalf("Error adding account: %v", err)
	}
	t.Logf("New account ID: %v", nAccount.Identifier)
	dAccount, err = GetAccount(oAccount.Email)
	if err != nil {
		t.Fatalf("Error getting account: %v", err)
	}
	if !dAccount.Equals(nAccount) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", oAccount, *nAccount)
	}
	if dAccount.Identifier != nAccount.Identifier {
		t.Errorf("Account id expected to be %v but found %v.", dAccount.Identifier, nAccount.Identifier)
	}
	oAccount = accounts[3]
	nAccount, err = AddAccount(oAccount)
	if err != nil {
		t.Fatalf("Error adding account: %v", err)
	}
	t.Logf("New account ID: %v", nAccount.Identifier)
	dAccount, err = GetAccount(oAccount.Email)
	if err != nil {
		t.Fatalf("Error getting account: %v", err)
	}
	if !dAccount.Equals(nAccount) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", oAccount, *nAccount)
	}
	if dAccount.Identifier != nAccount.Identifier {
		t.Errorf("Account id expected to be %v but found %v.", dAccount.Identifier, nAccount.Identifier)
	}
	// Test getting unknown accounts.
	dAccount, err = GetAccount("random@test.com")
	if err != nil {
		t.Fatalf("Error finding account not in existence: %v", err)
	}
	if dAccount != nil {
		t.Error("Expected not to find an account but one was found.")
	}
}

func TestGetAccounts(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupAccountTests()
	retAccounts, err := GetAccounts()
	if err != nil {
		t.Fatalf("Error getting accounts: %v", err)
	}
	if len(retAccounts) != 0 {
		t.Errorf("Expected number of accounts is %v but %v were found.", 0, len(retAccounts))
	}
	AddAccount(accounts[0])
	AddAccount(accounts[1])
	AddAccount(accounts[2])
	retAccounts, err = GetAccounts()
	if err != nil {
		t.Fatalf("Error getting accounts: %v", err)
	}
	if len(retAccounts) != 3 {
		t.Errorf("Expected number of accounts is %v but %v were found.", 3, len(retAccounts))
	}
	AddAccount(accounts[3])
	AddAccount(accounts[4])
	AddAccount(accounts[5])
	AddAccount(accounts[6])
	retAccounts, err = GetAccounts()
	if err != nil {
		t.Fatalf("Error getting accounts: %v", err)
	}
	if len(retAccounts) != 7 {
		t.Errorf("Expected number of accounts is %v but %v were found.", 7, len(retAccounts))
	}
}

func TestUpdateAccount(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupAccountTests()
	// Ensure adding accounts works properly.
	nAccount, err := AddAccount(accounts[0])
	if err != nil {
		t.Fatalf("Error adding account: %v", err)
	}
	nAccount.Name = "New Name 1"
	err = UpdateAccount(*nAccount)
	if err != nil {
		t.Fatalf("Error updating account: %v", err)
	}
	dAccount, _ := GetAccount(nAccount.Email)
	if nAccount.Identifier != dAccount.Identifier {
		t.Errorf("Account ID expected to be %v but found %v instead.", nAccount.Identifier, dAccount.Identifier)
	}
	if dAccount.Name != "New Name 1" {
		t.Errorf("Account name expected to be %v but found %v instead.", "New Name 1", dAccount.Name)
	}
	nAccount, err = AddAccount(accounts[1])
	if err != nil {
		t.Fatalf("Error adding account: %v", err)
	}
	nAccount.Type = "New Type 1"
	err = UpdateAccount(*nAccount)
	if err != nil {
		t.Fatalf("Error updating account: %v", err)
	}
	dAccount, _ = GetAccount(nAccount.Email)
	if nAccount.Identifier != dAccount.Identifier {
		t.Errorf("Account ID expected to be %v but found %v instead.", nAccount.Identifier, dAccount.Identifier)
	}
	if dAccount.Type != "New Type 1" {
		t.Errorf("Account name expected to be %v but found %v instead.", "New Type 1", dAccount.Type)
	}
	nAccount, err = AddAccount(accounts[2])
	if err != nil {
		t.Fatalf("Error adding account: %v", err)
	}
	nAccount.Name = "New Name 2"
	err = UpdateAccount(*nAccount)
	dAccount, _ = GetAccount(nAccount.Email)
	if err != nil {
		t.Fatalf("Error updating account: %v", err)
	}
	if nAccount.Identifier != dAccount.Identifier {
		t.Errorf("Account ID expected to be %v but found %v instead.", nAccount.Identifier, dAccount.Identifier)
	}
	if dAccount.Name != "New Name 2" {
		t.Errorf("Account name expected to be %v but found %v instead.", "New Name 2", dAccount.Name)
	}
	nAccount, err = AddAccount(accounts[3])
	if err != nil {
		t.Fatalf("Error adding account: %v", err)
	}
	nAccount.Type = "New Type 2"
	err = UpdateAccount(*nAccount)
	if err != nil {
		t.Fatalf("Error updating account: %v", err)
	}
	dAccount, _ = GetAccount(nAccount.Email)
	if nAccount.Identifier != dAccount.Identifier {
		t.Errorf("Account ID expected to be %v but found %v instead.", nAccount.Identifier, dAccount.Identifier)
	}
	if dAccount.Type != "New Type 2" {
		t.Errorf("Account name expected to be %v but found %v instead.", "New Type 2", dAccount.Type)
	}
	// Test for collisions.
	_, err = AddAccount(accounts[2])
	if err == nil {
		t.Error("Expected error adding account with duplicate email.")
	}
}

func TestDeleteAccount(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupAccountTests()
	nAccount, _ := AddAccount(accounts[0])
	err = DeleteAccount(nAccount.Email)
	if err != nil {
		t.Fatalf("Error deleting account: %v", err)
	}
	dAccount, _ := GetAccount(nAccount.Email)
	if dAccount != nil {
		t.Error("Unexpectedly found a deleted account.")
	}
	_, err = AddAccount(accounts[0])
	if err == nil {
		t.Error("No error found when trying to add a deleted account.")
	}
	nAccount, _ = AddAccount(accounts[1])
	err = DeleteAccount(nAccount.Email)
	if err != nil {
		t.Fatalf("Error deleting account: %v", err)
	}
	dAccount, _ = GetAccount(nAccount.Email)
	if dAccount != nil {
		t.Error("Unexpectedly found a deleted account.")
	}
	_, err = AddAccount(accounts[1])
	if err == nil {
		t.Error("No error found when trying to add a deleted account.")
	}
	nAccount, _ = AddAccount(accounts[2])
	err = DeleteAccount(nAccount.Email)
	if err != nil {
		t.Fatalf("Error deleting account: %v", err)
	}
	dAccount, _ = GetAccount(nAccount.Email)
	if dAccount != nil {
		t.Error("Unexpectedly found a deleted account.")
	}
	_, err = AddAccount(accounts[2])
	if err == nil {
		t.Error("No error found when trying to add a deleted account.")
	}
}

func TestResurrectAccount(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupAccountTests()
	nAccount, _ := AddAccount(accounts[0])
	DeleteAccount(nAccount.Email)
	err = ResurrectAccount(nAccount.Email)
	if err != nil {
		t.Fatalf("Error resurrecting account: %v", err)
	}
	dAccount, _ := GetAccount(nAccount.Email)
	if dAccount == nil {
		t.Error("Account was not resurrected.")
	}
	nAccount, _ = AddAccount(accounts[1])
	DeleteAccount(nAccount.Email)
	err = ResurrectAccount(nAccount.Email)
	if err != nil {
		t.Fatalf("Error resurrecting account: %v", err)
	}
	dAccount, _ = GetAccount(nAccount.Email)
	if dAccount == nil {
		t.Error("Account was not resurrected.")
	}
	nAccount, _ = AddAccount(accounts[4])
	DeleteAccount(nAccount.Email)
	err = ResurrectAccount(nAccount.Email)
	if err != nil {
		t.Fatalf("Error resurrecting account: %v", err)
	}
	dAccount, _ = GetAccount(nAccount.Email)
	if dAccount == nil {
		t.Error("Account was not resurrected.")
	}
}

func TestGetDeletedAccount(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupAccountTests()
	nAccount, _ := AddAccount(accounts[0])
	DeleteAccount(nAccount.Email)
	dAccount, err := GetDeletedAccount(nAccount.Email)
	if err != nil {
		t.Fatalf("Error getting deleted account %v", err)
	}
	if dAccount == nil {
		t.Error("Deleted account not found.")
	}
	nAccount, _ = AddAccount(accounts[3])
	DeleteAccount(nAccount.Email)
	dAccount, err = GetDeletedAccount(nAccount.Email)
	if err != nil {
		t.Fatalf("Error getting deleted account %v", err)
	}
	if dAccount == nil {
		t.Error("Deleted account not found.")
	}
	nAccount, _ = AddAccount(accounts[5])
	DeleteAccount(nAccount.Email)
	dAccount, err = GetDeletedAccount(nAccount.Email)
	if err != nil {
		t.Fatalf("Error getting deleted account %v", err)
	}
	if dAccount == nil {
		t.Error("Deleted account not found.")
	}
	nAccount, _ = AddAccount(accounts[6])
	DeleteAccount(nAccount.Email)
	dAccount, err = GetDeletedAccount(nAccount.Email)
	if err != nil {
		t.Fatalf("Error getting deleted account %v", err)
	}
	if dAccount == nil {
		t.Error("Deleted account not found.")
	}
}

func TestChangePassword(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupAccountTests()
	nAccount, _ := AddAccount(accounts[0])
	hashPass, _ := auth.HashPassword(testPassword2)
	err = ChangePassword(nAccount.Email, hashPass)
	if err != nil {
		t.Fatalf("error changing password: %v", err)
	}
	nAccount, _ = GetAccount(nAccount.Email)
	err = auth.VerifyPassword(nAccount.Password, testPassword2)
	if err != nil {
		t.Errorf("password doesn't match: %v", err)
	}
}

func TestChangeEmail(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupAccountTests()
	nAccount, _ := AddAccount(accounts[0])
	newEmail := "new_email2020@test.com"
	err = ChangeEmail(nAccount.Email, newEmail)
	if err != nil {
		t.Fatalf("error changing email: %v", err)
	}
	nAccount, _ = GetAccount(nAccount.Email)
	if nAccount != nil {
		t.Errorf("account retrieved when the email should have changed: %v", nAccount)
	}
	nAccount, _ = GetAccount(newEmail)
	if nAccount == nil {
		t.Error("account with new email not found")
	}
}

func TestInvalidPassword(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupAccountTests()
	nAccount, _ := AddAccount(accounts[0])
	for i := 1; i <= MaxLoginAttempts+3; i++ {
		err = InvalidPassword(*nAccount)
		if err != nil {
			t.Fatalf("(%v) error telling the database about an invalid password: %v", i, err)
		}
		nAccount, _ = GetAccount(nAccount.Email)
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
