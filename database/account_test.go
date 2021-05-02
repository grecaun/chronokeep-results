package database

import (
	"chronokeep/results/types"
	"testing"
)

func TestAddAccount(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	accounts := []types.Account{
		{
			Name:  "John Smith",
			Email: "j@test.com",
			Type:  "admin",
		},
		{
			Name:  "Jerry Garcia",
			Email: "jgarcia@test.com",
			Type:  "free",
		},
		{
			Name:  "Rose MacDonald",
			Email: "rose2004@test.com",
			Type:  "paid",
		},
		{
			Name:  "Tia Johnson",
			Email: "tiatheway@test.com",
			Type:  "free",
		},
		{
			Name:  "Thomas Donaldson",
			Email: "tdon@test.com",
			Type:  "admin",
		},
		{
			Name:  "Ester White",
			Email: "white@test.com",
			Type:  "test",
		},
		{
			Name:  "Ricky Reagan",
			Email: "rreagan@test.com",
			Type:  "free",
		},
	}
	// Ensure adding accounts works properly.
	oAccount := accounts[0]
	nAccount, err := AddAccount(oAccount)
	if err != nil {
		t.Errorf("Error adding account: %v", err)
	}
	t.Logf("New account ID: %v", nAccount.Identifier)
	if !oAccount.Equals(nAccount) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", oAccount, nAccount)
	}
	oAccount = accounts[1]
	nAccount, err = AddAccount(oAccount)
	if err != nil {
		t.Errorf("Error adding account: %v", err)
	}
	t.Logf("New account ID: %v", nAccount.Identifier)
	if !oAccount.Equals(nAccount) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", oAccount, nAccount)
	}
	oAccount = accounts[2]
	nAccount, err = AddAccount(oAccount)
	if err != nil {
		t.Errorf("Error adding account: %v", err)
	}
	t.Logf("New account ID: %v", nAccount.Identifier)
	if !oAccount.Equals(nAccount) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", oAccount, nAccount)
	}
	oAccount = accounts[3]
	nAccount, err = AddAccount(oAccount)
	if err != nil {
		t.Errorf("Error adding account: %v", err)
	}
	t.Logf("New account ID: %v", nAccount.Identifier)
	if !oAccount.Equals(nAccount) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", oAccount, nAccount)
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
	accounts := []types.Account{
		{
			Name:  "John Smith",
			Email: "j@test.com",
			Type:  "admin",
		},
		{
			Name:  "Jerry Garcia",
			Email: "jgarcia@test.com",
			Type:  "free",
		},
		{
			Name:  "Rose MacDonald",
			Email: "rose2004@test.com",
			Type:  "paid",
		},
		{
			Name:  "Tia Johnson",
			Email: "tiatheway@test.com",
			Type:  "free",
		},
		{
			Name:  "Thomas Donaldson",
			Email: "tdon@test.com",
			Type:  "admin",
		},
		{
			Name:  "Ester White",
			Email: "white@test.com",
			Type:  "test",
		},
		{
			Name:  "Ricky Reagan",
			Email: "rreagan@test.com",
			Type:  "free",
		},
	}
	// Test getting known accounts.
	oAccount := accounts[0]
	nAccount, err := AddAccount(oAccount)
	if err != nil {
		t.Errorf("Error adding account: %v", err)
	}
	t.Logf("New account ID: %v", nAccount.Identifier)
	dAccount, err := GetAccount(oAccount.Email)
	if !dAccount.Equals(nAccount) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", oAccount, nAccount)
	}
	if dAccount.Identifier != nAccount.Identifier {
		t.Errorf("Account id expected to be %v but found %v.", dAccount.Identifier, nAccount.Identifier)
	}
	oAccount = accounts[1]
	nAccount, err = AddAccount(oAccount)
	if err != nil {
		t.Errorf("Error adding account: %v", err)
	}
	t.Logf("New account ID: %v", nAccount.Identifier)
	dAccount, err = GetAccount(oAccount.Email)
	if !dAccount.Equals(nAccount) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", oAccount, nAccount)
	}
	if dAccount.Identifier != nAccount.Identifier {
		t.Errorf("Account id expected to be %v but found %v.", dAccount.Identifier, nAccount.Identifier)
	}
	oAccount = accounts[2]
	nAccount, err = AddAccount(oAccount)
	if err != nil {
		t.Errorf("Error adding account: %v", err)
	}
	t.Logf("New account ID: %v", nAccount.Identifier)
	dAccount, err = GetAccount(oAccount.Email)
	if !dAccount.Equals(nAccount) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", oAccount, nAccount)
	}
	if dAccount.Identifier != nAccount.Identifier {
		t.Errorf("Account id expected to be %v but found %v.", dAccount.Identifier, nAccount.Identifier)
	}
	oAccount = accounts[3]
	nAccount, err = AddAccount(oAccount)
	if err != nil {
		t.Errorf("Error adding account: %v", err)
	}
	t.Logf("New account ID: %v", nAccount.Identifier)
	dAccount, err = GetAccount(oAccount.Email)
	if !dAccount.Equals(nAccount) {
		t.Errorf("Account expected to be equal. %+v was expected, found %+v", oAccount, nAccount)
	}
	if dAccount.Identifier != nAccount.Identifier {
		t.Errorf("Account id expected to be %v but found %v.", dAccount.Identifier, nAccount.Identifier)
	}
	// Test getting unknown accounts.
	dAccount, err = GetAccount("random@test.com")
	if err != nil {
		t.Errorf("Error finding account not in existence: %v", err)
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
	accounts := []types.Account{
		{
			Name:  "John Smith",
			Email: "j@test.com",
			Type:  "admin",
		},
		{
			Name:  "Jerry Garcia",
			Email: "jgarcia@test.com",
			Type:  "free",
		},
		{
			Name:  "Rose MacDonald",
			Email: "rose2004@test.com",
			Type:  "paid",
		},
		{
			Name:  "Tia Johnson",
			Email: "tiatheway@test.com",
			Type:  "free",
		},
		{
			Name:  "Thomas Donaldson",
			Email: "tdon@test.com",
			Type:  "admin",
		},
		{
			Name:  "Ester White",
			Email: "white@test.com",
			Type:  "test",
		},
		{
			Name:  "Ricky Reagan",
			Email: "rreagan@test.com",
			Type:  "free",
		},
	}
	AddAccount(accounts[0])
	AddAccount(accounts[1])
	AddAccount(accounts[2])
	retAccounts, err := GetAccounts()
	if err != nil {
		t.Errorf("Error getting accounts: %v", err)
	}
	if len(retAccounts) != 3 {
		t.Errorf("Expected number of accounts is %v but %v were found.", 3, len(accounts))
	}
	AddAccount(accounts[3])
	AddAccount(accounts[4])
	AddAccount(accounts[5])
	AddAccount(accounts[6])
	accounts, err = GetAccounts()
	if err != nil {
		t.Errorf("Error getting accounts: %v", err)
	}
	if len(retAccounts) != 7 {
		t.Errorf("Expected number of accounts is %v but %v were found.", 7, len(accounts))
	}
}

func TestUpdateAccount(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	accounts := []types.Account{
		{
			Name:  "John Smith",
			Email: "j@test.com",
			Type:  "admin",
		},
		{
			Name:  "Jerry Garcia",
			Email: "jgarcia@test.com",
			Type:  "free",
		},
		{
			Name:  "Rose MacDonald",
			Email: "rose2004@test.com",
			Type:  "paid",
		},
		{
			Name:  "Tia Johnson",
			Email: "tiatheway@test.com",
			Type:  "free",
		},
		{
			Name:  "Thomas Donaldson",
			Email: "tdon@test.com",
			Type:  "admin",
		},
		{
			Name:  "Ester White",
			Email: "white@test.com",
			Type:  "test",
		},
		{
			Name:  "Ricky Reagan",
			Email: "rreagan@test.com",
			Type:  "free",
		},
	}
	// Ensure adding accounts works properly.
	nAccount, err := AddAccount(accounts[0])
	nAccount.Name = "New Name 1"
	err = UpdateAccount(*nAccount)
	if err != nil {
		t.Errorf("Error updating account: %v", err)
	}
	dAccount, _ := GetAccount(nAccount.Email)
	if nAccount.Identifier != dAccount.Identifier {
		t.Errorf("Account ID expected to be %v but found %v instead.", nAccount, dAccount)
	}
	if dAccount.Name != "New Name 1" {
		t.Errorf("Account name expected to be %v but found %v instead.", "New Name 1", dAccount.Name)
	}
	nAccount, err = AddAccount(accounts[1])
	nAccount.Type = "New Type 1"
	err = UpdateAccount(*nAccount)
	if err != nil {
		t.Errorf("Error updating account: %v", err)
	}
	dAccount, _ = GetAccount(nAccount.Email)
	if nAccount.Identifier != dAccount.Identifier {
		t.Errorf("Account ID expected to be %v but found %v instead.", nAccount, dAccount)
	}
	if dAccount.Type != "New Type 1" {
		t.Errorf("Account name expected to be %v but found %v instead.", "New Type 1", dAccount.Type)
	}
	nAccount, err = AddAccount(accounts[2])
	nAccount.Name = "New Name 2"
	err = UpdateAccount(*nAccount)
	dAccount, _ = GetAccount(nAccount.Email)
	if err != nil {
		t.Errorf("Error updating account: %v", err)
	}
	if nAccount.Identifier != dAccount.Identifier {
		t.Errorf("Account ID expected to be %v but found %v instead.", nAccount, dAccount)
	}
	if dAccount.Name != "New Name 2" {
		t.Errorf("Account name expected to be %v but found %v instead.", "New Name 2", dAccount.Name)
	}
	nAccount, err = AddAccount(accounts[3])
	nAccount.Type = "New Type 2"
	err = UpdateAccount(*nAccount)
	if err != nil {
		t.Errorf("Error updating account: %v", err)
	}
	dAccount, _ = GetAccount(nAccount.Email)
	if nAccount.Identifier != dAccount.Identifier {
		t.Errorf("Account ID expected to be %v but found %v instead.", nAccount, dAccount)
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
	accounts := []types.Account{
		{
			Name:  "John Smith",
			Email: "j@test.com",
			Type:  "admin",
		},
		{
			Name:  "Jerry Garcia",
			Email: "jgarcia@test.com",
			Type:  "free",
		},
		{
			Name:  "Rose MacDonald",
			Email: "rose2004@test.com",
			Type:  "paid",
		},
		{
			Name:  "Tia Johnson",
			Email: "tiatheway@test.com",
			Type:  "free",
		},
		{
			Name:  "Thomas Donaldson",
			Email: "tdon@test.com",
			Type:  "admin",
		},
		{
			Name:  "Ester White",
			Email: "white@test.com",
			Type:  "test",
		},
		{
			Name:  "Ricky Reagan",
			Email: "rreagan@test.com",
			Type:  "free",
		},
	}
	nAccount, _ := AddAccount(accounts[0])
	err = DeleteAccount(*nAccount)
	if err != nil {
		t.Errorf("Error deleting account: %v", err)
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
	err = DeleteAccount(*nAccount)
	if err != nil {
		t.Errorf("Error deleting account: %v", err)
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
	err = DeleteAccount(*nAccount)
	if err != nil {
		t.Errorf("Error deleting account: %v", err)
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
	accounts := []types.Account{
		{
			Name:  "John Smith",
			Email: "j@test.com",
			Type:  "admin",
		},
		{
			Name:  "Jerry Garcia",
			Email: "jgarcia@test.com",
			Type:  "free",
		},
		{
			Name:  "Rose MacDonald",
			Email: "rose2004@test.com",
			Type:  "paid",
		},
		{
			Name:  "Tia Johnson",
			Email: "tiatheway@test.com",
			Type:  "free",
		},
		{
			Name:  "Thomas Donaldson",
			Email: "tdon@test.com",
			Type:  "admin",
		},
		{
			Name:  "Ester White",
			Email: "white@test.com",
			Type:  "test",
		},
		{
			Name:  "Ricky Reagan",
			Email: "rreagan@test.com",
			Type:  "free",
		},
	}
	nAccount, _ := AddAccount(accounts[0])
	DeleteAccount(*nAccount)
	err = ResurrectAccount(nAccount.Email)
	if err != nil {
		t.Errorf("Error resurrecting account: %v", err)
	}
	dAccount, _ := GetAccount(nAccount.Email)
	if dAccount == nil {
		t.Error("Account was not resurrected.")
	}
	nAccount, _ = AddAccount(accounts[1])
	DeleteAccount(*nAccount)
	err = ResurrectAccount(nAccount.Email)
	if err != nil {
		t.Errorf("Error resurrecting account: %v", err)
	}
	dAccount, _ = GetAccount(nAccount.Email)
	if dAccount == nil {
		t.Error("Account was not resurrected.")
	}
	nAccount, _ = AddAccount(accounts[4])
	DeleteAccount(*nAccount)
	err = ResurrectAccount(nAccount.Email)
	if err != nil {
		t.Errorf("Error resurrecting account: %v", err)
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
	accounts := []types.Account{
		{
			Name:  "John Smith",
			Email: "j@test.com",
			Type:  "admin",
		},
		{
			Name:  "Jerry Garcia",
			Email: "jgarcia@test.com",
			Type:  "free",
		},
		{
			Name:  "Rose MacDonald",
			Email: "rose2004@test.com",
			Type:  "paid",
		},
		{
			Name:  "Tia Johnson",
			Email: "tiatheway@test.com",
			Type:  "free",
		},
		{
			Name:  "Thomas Donaldson",
			Email: "tdon@test.com",
			Type:  "admin",
		},
		{
			Name:  "Ester White",
			Email: "white@test.com",
			Type:  "test",
		},
		{
			Name:  "Ricky Reagan",
			Email: "rreagan@test.com",
			Type:  "free",
		},
	}
	nAccount, _ := AddAccount(accounts[0])
	DeleteAccount(*nAccount)
	dAccount, err := GetDeletedAccount(nAccount.Email)
	if err != nil {
		t.Errorf("Error getting deleted account %v", err)
	}
	if dAccount == nil {
		t.Error("Deleted account not found.")
	}
	nAccount, _ = AddAccount(accounts[3])
	DeleteAccount(*nAccount)
	dAccount, err = GetDeletedAccount(nAccount.Email)
	if err != nil {
		t.Errorf("Error getting deleted account %v", err)
	}
	if dAccount == nil {
		t.Error("Deleted account not found.")
	}
	nAccount, _ = AddAccount(accounts[5])
	DeleteAccount(*nAccount)
	dAccount, err = GetDeletedAccount(nAccount.Email)
	if err != nil {
		t.Errorf("Error getting deleted account %v", err)
	}
	if dAccount == nil {
		t.Error("Deleted account not found.")
	}
	nAccount, _ = AddAccount(accounts[6])
	DeleteAccount(*nAccount)
	dAccount, err = GetDeletedAccount(nAccount.Email)
	if err != nil {
		t.Errorf("Error getting deleted account %v", err)
	}
	if dAccount == nil {
		t.Error("Deleted account not found.")
	}

}
