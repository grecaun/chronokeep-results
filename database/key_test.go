package database

import (
	"chronokeep/results/types"
	"testing"
	"time"
)

func TestAddKey(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	account1 := &types.Account{
		Name:  "John Smith",
		Email: "j@test.com",
		Type:  "admin",
	}
	account2 := &types.Account{
		Name:  "Rose MacDonald",
		Email: "rose2004@test.com",
		Type:  "paid",
	}
	account1, _ = AddAccount(*account1)
	account2, _ = AddAccount(*account2)
	keys := []types.Key{
		{
			AccountIdentifier: account1.Identifier,
			Value:             "030001-1ACSDD-K2389A-00123B",
			Type:              "default",
			AllowedHosts:      "",
			ValidUntil:        time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
		},
		{
			AccountIdentifier: account1.Identifier,
			Value:             "030001-1ACSDD-K2389A-22123B",
			Type:              "write",
			AllowedHosts:      "",
			ValidUntil:        time.Now().Add(time.Hour * 20).Truncate(time.Second),
		},
		{
			AccountIdentifier: account2.Identifier,
			Value:             "030001-1ACSDD-KH789A-00123B",
			Type:              "delete",
			AllowedHosts:      "https://test.com/",
			ValidUntil:        time.Date(2016, 4, 1, 4, 11, 5, 0, time.Local),
		},
		{
			AccountIdentifier: account2.Identifier,
			Value:             "030001-1ACSCT-K2389A-22123B",
			Type:              "write",
			AllowedHosts:      "",
			ValidUntil:        time.Now().Add(time.Hour * 3).Truncate(time.Second),
		},
	}
	key, err := AddKey(keys[0])
	if err != nil {
		t.Fatalf("Error adding key: %v", err)
	}
	if !key.Equal(&keys[0]) {
		t.Errorf("Expected key %+v, found %+v", keys[0], *key)
	}
	key, err = AddKey(keys[1])
	if err != nil {
		t.Fatalf("Error adding key: %v", err)
	}
	if !key.Equal(&keys[1]) {
		t.Errorf("Expected key %+v, found %+v", keys[1], *key)
	}
	key, err = AddKey(keys[2])
	if err != nil {
		t.Fatalf("Error adding key: %v", err)
	}
	if !key.Equal(&keys[2]) {
		t.Errorf("Expected key %+v, found %+v", keys[2], *key)
	}
	key, err = AddKey(keys[3])
	if err != nil {
		t.Fatalf("Error adding key: %v", err)
	}
	if !key.Equal(&keys[3]) {
		t.Errorf("Expected key %+v, found %+v", keys[3], *key)
	}
	key, err = AddKey(keys[3])
	if err == nil {
		t.Errorf("Expected error adding key that exists, found key %+v", key)
	}
}

func TestGetAccountKeys(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	account1 := &types.Account{
		Name:  "John Smith",
		Email: "j@test.com",
		Type:  "admin",
	}
	account2 := &types.Account{
		Name:  "Rose MacDonald",
		Email: "rose2004@test.com",
		Type:  "paid",
	}
	account1, _ = AddAccount(*account1)
	account2, _ = AddAccount(*account2)
	keys := []types.Key{
		{
			AccountIdentifier: account1.Identifier,
			Value:             "030001-1ACSDD-K2389A-00123B",
			Type:              "default",
			AllowedHosts:      "",
			ValidUntil:        time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
		},
		{
			AccountIdentifier: account1.Identifier,
			Value:             "030001-1ACSDD-K2389A-22123B",
			Type:              "write",
			AllowedHosts:      "",
			ValidUntil:        time.Now().Add(time.Hour * 20).Truncate(time.Second),
		},
		{
			AccountIdentifier: account2.Identifier,
			Value:             "030001-1ACSDD-KH789A-00123B",
			Type:              "delete",
			AllowedHosts:      "https://test.com/",
			ValidUntil:        time.Date(2016, 4, 1, 4, 11, 5, 0, time.Local),
		},
		{
			AccountIdentifier: account2.Identifier,
			Value:             "030001-1ACSCT-K2389A-22123B",
			Type:              "write",
			AllowedHosts:      "",
			ValidUntil:        time.Now().Add(time.Hour * 3).Truncate(time.Second),
		},
	}
	k, err := GetAccountKeys(account1.Identifier)
	if err != nil {
		t.Fatalf("Error getting account keys: %v", err)
	}
	if len(k) != 0 {
		t.Errorf("Expected no keys found for account but found %v keys.", len(k))
	}
	AddKey(keys[0])
	AddKey(keys[2])
	k, err = GetAccountKeys(account1.Identifier)
	if err != nil {
		t.Fatalf("Error getting account keys: %v", err)
	}
	if len(k) != 1 {
		t.Errorf("Expected %v keys found for account but found %v keys.", 1, len(k))
	}
	k, err = GetAccountKeys(account2.Identifier)
	if err != nil {
		t.Fatalf("Error getting account keys: %v", err)
	}
	if len(k) != 1 {
		t.Errorf("Expected %v keys found for account but found %v keys.", 1, len(k))
	}
	AddKey(keys[1])
	AddKey(keys[3])
	k, err = GetAccountKeys(account1.Identifier)
	if err != nil {
		t.Fatalf("Error getting account keys: %v", err)
	}
	if len(k) != 2 {
		t.Errorf("Expected %v keys found for account but found %v keys.", 2, len(k))
	}
	k, err = GetAccountKeys(account2.Identifier)
	if err != nil {
		t.Fatalf("Error getting account keys: %v", err)
	}
	if len(k) != 2 {
		t.Errorf("Expected %v keys found for account but found %v keys.", 2, len(k))
	}
}

func TestGetKey(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	account1 := &types.Account{
		Name:  "John Smith",
		Email: "j@test.com",
		Type:  "admin",
	}
	account2 := &types.Account{
		Name:  "Rose MacDonald",
		Email: "rose2004@test.com",
		Type:  "paid",
	}
	account1, _ = AddAccount(*account1)
	account2, _ = AddAccount(*account2)
	keys := []types.Key{
		{
			AccountIdentifier: account1.Identifier,
			Value:             "030001-1ACSDD-K2389A-00123B",
			Type:              "default",
			AllowedHosts:      "",
			ValidUntil:        time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
		},
		{
			AccountIdentifier: account1.Identifier,
			Value:             "030001-1ACSDD-K2389A-22123B",
			Type:              "write",
			AllowedHosts:      "",
			ValidUntil:        time.Now().Add(time.Hour * 20).Truncate(time.Second),
		},
		{
			AccountIdentifier: account2.Identifier,
			Value:             "030001-1ACSDD-KH789A-00123B",
			Type:              "delete",
			AllowedHosts:      "https://test.com/",
			ValidUntil:        time.Date(2016, 4, 1, 4, 11, 5, 0, time.Local),
		},
		{
			AccountIdentifier: account2.Identifier,
			Value:             "030001-1ACSCT-K2389A-22123B",
			Type:              "write",
			AllowedHosts:      "",
			ValidUntil:        time.Now().Add(time.Hour * 3).Truncate(time.Second),
		},
	}
	AddKey(keys[0])
	AddKey(keys[1])
	AddKey(keys[2])
	AddKey(keys[3])
	key, err := GetKey(keys[0].Value)
	if err != nil {
		t.Fatalf("Error getting key: %v", err)
	}
	if !key.Equal(&keys[0]) {
		t.Errorf("Expected key %+v, found %+v.", keys[0], *key)
	}
	key, err = GetKey(keys[1].Value)
	if err != nil {
		t.Fatalf("Error getting key: %v", err)
	}
	if !key.Equal(&keys[1]) {
		t.Errorf("Expected key %+v, found %+v.", keys[1], *key)
	}
	key, err = GetKey(keys[2].Value)
	if err != nil {
		t.Fatalf("Error getting key: %v", err)
	}
	if !key.Equal(&keys[2]) {
		t.Errorf("Expected key %+v, found %+v.", keys[2], *key)
	}
	key, err = GetKey(keys[3].Value)
	if err != nil {
		t.Fatalf("Error getting key: %v", err)
	}
	if !key.Equal(&keys[3]) {
		t.Errorf("Expected key %+v, found %+v.", keys[3], *key)
	}
	key, err = GetKey("test-value")
	if err != nil {
		t.Fatalf("Error getting non-existant key: %v", err)
	}
	if key != nil {
		t.Errorf("Expected no key but found %+v.", *key)
	}
}

func TestDeleteKey(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	account1 := &types.Account{
		Name:  "John Smith",
		Email: "j@test.com",
		Type:  "admin",
	}
	account2 := &types.Account{
		Name:  "Rose MacDonald",
		Email: "rose2004@test.com",
		Type:  "paid",
	}
	account1, _ = AddAccount(*account1)
	account2, _ = AddAccount(*account2)
	keys := []types.Key{
		{
			AccountIdentifier: account1.Identifier,
			Value:             "030001-1ACSDD-K2389A-00123B",
			Type:              "default",
			AllowedHosts:      "",
			ValidUntil:        time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
		},
		{
			AccountIdentifier: account1.Identifier,
			Value:             "030001-1ACSDD-K2389A-22123B",
			Type:              "write",
			AllowedHosts:      "",
			ValidUntil:        time.Now().Add(time.Hour * 20).Truncate(time.Second),
		},
		{
			AccountIdentifier: account2.Identifier,
			Value:             "030001-1ACSDD-KH789A-00123B",
			Type:              "delete",
			AllowedHosts:      "https://test.com/",
			ValidUntil:        time.Date(2016, 4, 1, 4, 11, 5, 0, time.Local),
		},
		{
			AccountIdentifier: account2.Identifier,
			Value:             "030001-1ACSCT-K2389A-22123B",
			Type:              "write",
			AllowedHosts:      "",
			ValidUntil:        time.Now().Add(time.Hour * 3).Truncate(time.Second),
		},
	}
	AddKey(keys[0])
	AddKey(keys[1])
	AddKey(keys[2])
	AddKey(keys[3])
	err = DeleteKey(keys[0])
	if err != nil {
		t.Fatalf("Error deleting key: %v", err)
	}
	k, _ := GetKey(keys[0].Value)
	if k != nil {
		t.Errorf("Found deleted key: %+v", k)
	}
	err = DeleteKey(keys[1])
	if err != nil {
		t.Fatalf("Error deleting key: %v", err)
	}
	k, _ = GetKey(keys[1].Value)
	if k != nil {
		t.Errorf("Found deleted key: %+v", k)
	}
	err = DeleteKey(keys[2])
	if err != nil {
		t.Fatalf("Error deleting key: %v", err)
	}
	k, _ = GetKey(keys[2].Value)
	if k != nil {
		t.Errorf("Found deleted key: %+v", k)
	}
	err = DeleteKey(keys[3])
	if err != nil {
		t.Fatalf("Error deleting key: %v", err)
	}
	k, _ = GetKey(keys[3].Value)
	if k != nil {
		t.Errorf("Found deleted key: %+v", k)
	}
	err = DeleteKey(keys[3])
	if err == nil {
		t.Error("Expected error from deletion of already deleted key.")
	}
}

func TestUpdateKey(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	account1 := &types.Account{
		Name:  "John Smith",
		Email: "j@test.com",
		Type:  "admin",
	}
	account1, _ = AddAccount(*account1)
	keys := []types.Key{
		{
			AccountIdentifier: account1.Identifier,
			Value:             "030001-1ACSDD-K2389A-00123B",
			Type:              "default",
			AllowedHosts:      "",
			ValidUntil:        time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
		},
		{
			AccountIdentifier: account1.Identifier,
			Value:             "030001-1ACSDD-K2389A-22123B",
			Type:              "write",
			AllowedHosts:      "",
			ValidUntil:        time.Now().Add(time.Hour * 20).Truncate(time.Second),
		},
	}
	AddKey(keys[0])
	AddKey(keys[1])
	keys[0].Type = "write"
	keys[0].AllowedHosts = "test.lan,test.com,test.org"
	keys[0].ValidUntil = time.Now().Add(time.Minute * 30).Truncate(time.Second)
	err = UpdateKey(keys[0])
	if err != nil {
		t.Fatalf("Error updating key: %v", err)
	}
	key, _ := GetKey(keys[0].Value)
	if !key.Equal(&keys[0]) {
		t.Errorf("Expected key %+v, found %+v.", keys[0], *key)
	}
	keys[1].AccountIdentifier = account1.Identifier + 200
	keys[1].Value = "update-value-test"
	err = UpdateKey(keys[1])
	if err == nil {
		t.Error("Expected error from update with no changed values.")
	}
	key, _ = GetKey(keys[1].Value)
	if key != nil {
		t.Errorf("Found key with modified key value: %+v", key)
	}
}
