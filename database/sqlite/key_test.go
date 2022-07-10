package sqlite

import (
	"chronokeep/results/types"
	"testing"
	"time"
)

func setupKeyTests() {
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
		}
	}
}

func TestAddKey(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupKeyTests()
	account1, _ := db.AddAccount(accounts[0])
	account2, _ := db.AddAccount(accounts[1])
	times := []time.Time{
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
		time.Now().Add(time.Hour * 20).Truncate(time.Second),
		time.Date(2016, 4, 1, 4, 11, 5, 0, time.Local),
	}
	keys := []types.Key{
		{
			AccountIdentifier: account1.Identifier,
			Name:              "test1",
			Value:             "030001-1ACSDD-K2389A-00123B",
			Type:              "default",
			AllowedHosts:      "",
			ValidUntil:        &times[0],
		},
		{
			AccountIdentifier: account1.Identifier,
			Value:             "030001-1ACSDD-K2389A-22123B",
			Type:              "write",
			AllowedHosts:      "",
			ValidUntil:        &times[1],
		},
		{
			AccountIdentifier: account2.Identifier,
			Name:              "test2",
			Value:             "030001-1ACSDD-KH789A-00123B",
			Type:              "delete",
			AllowedHosts:      "https://test.com/",
			ValidUntil:        &times[2],
		},
		{
			AccountIdentifier: account2.Identifier,
			Value:             "030001-1ACSCT-K2389A-22123B",
			Type:              "write",
			AllowedHosts:      "",
			ValidUntil:        nil,
		},
	}
	key, err := db.AddKey(keys[0])
	if err != nil {
		t.Fatalf("Error adding key: %v", err)
	}
	if !key.Equal(&keys[0]) {
		t.Errorf("Expected key %+v, found %+v", keys[0], *key)
	}
	if key.Name != "test1" {
		t.Errorf("Expected key to be named %s, found %s.", "test1", key.Name)
	}
	key, err = db.AddKey(keys[1])
	if err != nil {
		t.Fatalf("Error adding key: %v", err)
	}
	if !key.Equal(&keys[1]) {
		t.Errorf("Expected key %+v, found %+v", keys[1], *key)
	}
	key, err = db.AddKey(keys[2])
	if err != nil {
		t.Fatalf("Error adding key: %v", err)
	}
	if !key.Equal(&keys[2]) {
		t.Errorf("Expected key %+v, found %+v", keys[2], *key)
	}
	if key.Name != "test2" {
		t.Errorf("Expected key to be named %s, found %s.", "test2", key.Name)
	}
	key, err = db.AddKey(keys[3])
	if err != nil {
		t.Fatalf("Error adding key: %v", err)
	}
	if !key.Equal(&keys[3]) {
		t.Errorf("Expected key %+v, found %+v", keys[3], *key)
	}
	key, err = db.AddKey(keys[3])
	if err == nil {
		t.Errorf("Expected error adding key that exists, found key %+v", key)
	}
}

func TestGetAccountKeys(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupKeyTests()
	account1, _ := db.AddAccount(accounts[0])
	account2, _ := db.AddAccount(accounts[1])
	times := []time.Time{
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
		time.Now().Add(time.Hour * 20).Truncate(time.Second),
		time.Date(2016, 4, 1, 4, 11, 5, 0, time.Local),
	}
	keys := []types.Key{
		{
			AccountIdentifier: account1.Identifier,
			Name:              "test1",
			Value:             "030001-1ACSDD-K2389A-00123B",
			Type:              "default",
			AllowedHosts:      "",
			ValidUntil:        &times[0],
		},
		{
			AccountIdentifier: account1.Identifier,
			Value:             "030001-1ACSDD-K2389A-22123B",
			Type:              "write",
			AllowedHosts:      "",
			ValidUntil:        &times[1],
		},
		{
			AccountIdentifier: account2.Identifier,
			Name:              "test2",
			Value:             "030001-1ACSDD-KH789A-00123B",
			Type:              "delete",
			AllowedHosts:      "https://test.com/",
			ValidUntil:        &times[2],
		},
		{
			AccountIdentifier: account2.Identifier,
			Value:             "030001-1ACSCT-K2389A-22123B",
			Type:              "write",
			AllowedHosts:      "",
			ValidUntil:        nil,
		},
	}
	k, err := db.GetAccountKeys(account1.Email)
	if err != nil {
		t.Fatalf("Error getting account keys: %v", err)
	}
	if len(k) != 0 {
		t.Errorf("Expected no keys found for account but found %v keys.", len(k))
	}
	db.AddKey(keys[0])
	db.AddKey(keys[2])
	k, err = db.GetAccountKeys(account1.Email)
	if err != nil {
		t.Fatalf("Error getting account keys: %v", err)
	}
	if len(k) != 1 {
		t.Errorf("Expected %v keys found for account but found %v keys.", 1, len(k))
	}
	k, err = db.GetAccountKeys(account2.Email)
	if err != nil {
		t.Fatalf("Error getting account keys: %v", err)
	}
	if len(k) != 1 {
		t.Errorf("Expected %v keys found for account but found %v keys.", 1, len(k))
	}
	db.AddKey(keys[1])
	db.AddKey(keys[3])
	k, err = db.GetAccountKeys(account1.Email)
	if err != nil {
		t.Fatalf("Error getting account keys: %v", err)
	}
	if len(k) != 2 {
		t.Errorf("Expected %v keys found for account but found %v keys.", 2, len(k))
	}
	k, err = db.GetAccountKeys(account2.Email)
	if err != nil {
		t.Fatalf("Error getting account keys: %v", err)
	}
	if len(k) != 2 {
		t.Errorf("Expected %v keys found for account but found %v keys.", 2, len(k))
	}
}

func TestGetKey(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupKeyTests()
	account1, _ := db.AddAccount(accounts[0])
	account2, _ := db.AddAccount(accounts[1])
	times := []time.Time{
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
		time.Now().Add(time.Hour * 20).Truncate(time.Second),
		time.Date(2016, 4, 1, 4, 11, 5, 0, time.Local),
	}
	keys := []types.Key{
		{
			AccountIdentifier: account1.Identifier,
			Name:              "test1",
			Value:             "030001-1ACSDD-K2389A-00123B",
			Type:              "default",
			AllowedHosts:      "",
			ValidUntil:        &times[0],
		},
		{
			AccountIdentifier: account1.Identifier,
			Value:             "030001-1ACSDD-K2389A-22123B",
			Type:              "write",
			AllowedHosts:      "",
			ValidUntil:        &times[1],
		},
		{
			AccountIdentifier: account2.Identifier,
			Name:              "test2",
			Value:             "030001-1ACSDD-KH789A-00123B",
			Type:              "delete",
			AllowedHosts:      "https://test.com/",
			ValidUntil:        &times[2],
		},
		{
			AccountIdentifier: account2.Identifier,
			Value:             "030001-1ACSCT-K2389A-22123B",
			Type:              "write",
			AllowedHosts:      "",
			ValidUntil:        nil,
		},
	}
	db.AddKey(keys[0])
	db.AddKey(keys[1])
	db.AddKey(keys[2])
	db.AddKey(keys[3])
	key, err := db.GetKey(keys[0].Value)
	if err != nil {
		t.Fatalf("Error getting key: %v", err)
	}
	if !key.Equal(&keys[0]) {
		t.Errorf("Expected key %+v, found %+v.", keys[0], *key)
	}
	if key.Name != "test1" {
		t.Errorf("Expected key to be named %s, found %s.", "test1", key.Name)
	}
	key, err = db.GetKey(keys[1].Value)
	if err != nil {
		t.Fatalf("Error getting key: %v", err)
	}
	if !key.Equal(&keys[1]) {
		t.Errorf("Expected key %+v, found %+v.", keys[1], *key)
	}
	key, err = db.GetKey(keys[2].Value)
	if err != nil {
		t.Fatalf("Error getting key: %v", err)
	}
	if !key.Equal(&keys[2]) {
		t.Errorf("Expected key %+v, found %+v.", keys[2], *key)
	}
	if key.Name != "test2" {
		t.Errorf("Expected key to be named %s, found %s.", "test2", key.Name)
	}
	key, err = db.GetKey(keys[3].Value)
	if err != nil {
		t.Fatalf("Error getting key: %v", err)
	}
	if !key.Equal(&keys[3]) {
		t.Errorf("Expected key %+v, found %+v.", keys[3], *key)
	}
	key, err = db.GetKey("test-value")
	if err != nil {
		t.Fatalf("Error getting non-existant key: %v", err)
	}
	if key != nil {
		t.Errorf("Expected no key but found %+v.", *key)
	}
}

func TestDeleteKey(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupKeyTests()
	account1, _ := db.AddAccount(accounts[0])
	account2, _ := db.AddAccount(accounts[1])
	times := []time.Time{
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
		time.Now().Add(time.Hour * 20).Truncate(time.Second),
		time.Date(2016, 4, 1, 4, 11, 5, 0, time.Local),
	}
	keys := []types.Key{
		{
			AccountIdentifier: account1.Identifier,
			Name:              "test1",
			Value:             "030001-1ACSDD-K2389A-00123B",
			Type:              "default",
			AllowedHosts:      "",
			ValidUntil:        &times[0],
		},
		{
			AccountIdentifier: account1.Identifier,
			Value:             "030001-1ACSDD-K2389A-22123B",
			Type:              "write",
			AllowedHosts:      "",
			ValidUntil:        &times[1],
		},
		{
			AccountIdentifier: account2.Identifier,
			Name:              "test2",
			Value:             "030001-1ACSDD-KH789A-00123B",
			Type:              "delete",
			AllowedHosts:      "https://test.com/",
			ValidUntil:        &times[2],
		},
		{
			AccountIdentifier: account2.Identifier,
			Value:             "030001-1ACSCT-K2389A-22123B",
			Type:              "write",
			AllowedHosts:      "",
			ValidUntil:        nil,
		},
	}
	db.AddKey(keys[0])
	db.AddKey(keys[1])
	db.AddKey(keys[2])
	db.AddKey(keys[3])
	err = db.DeleteKey(keys[0])
	if err != nil {
		t.Fatalf("Error deleting key: %v", err)
	}
	k, _ := db.GetKey(keys[0].Value)
	if k != nil {
		t.Errorf("Found deleted key: %+v", k)
	}
	err = db.DeleteKey(keys[1])
	if err != nil {
		t.Fatalf("Error deleting key: %v", err)
	}
	k, _ = db.GetKey(keys[1].Value)
	if k != nil {
		t.Errorf("Found deleted key: %+v", k)
	}
	err = db.DeleteKey(keys[2])
	if err != nil {
		t.Fatalf("Error deleting key: %v", err)
	}
	k, _ = db.GetKey(keys[2].Value)
	if k != nil {
		t.Errorf("Found deleted key: %+v", k)
	}
	err = db.DeleteKey(keys[3])
	if err != nil {
		t.Fatalf("Error deleting key: %v", err)
	}
	k, _ = db.GetKey(keys[3].Value)
	if k != nil {
		t.Errorf("Found deleted key: %+v", k)
	}
	/*// It appears that RowsAffected still throws a 1 if nothing was updated with the PGX library.
	err = db.DeleteKey(keys[3])
	if err == nil {
		t.Error("Expected error from deletion of already deleted key.")
	}//*/
}

func TestUpdateKey(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	setupKeyTests()
	account1, _ := db.AddAccount(accounts[0])
	times := []time.Time{
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
		time.Now().Add(time.Hour * 20).Truncate(time.Second),
		time.Date(2016, 4, 1, 4, 11, 5, 0, time.Local),
	}
	keys := []types.Key{
		{
			AccountIdentifier: account1.Identifier,
			Name:              "test1",
			Value:             "030001-1ACSDD-K2389A-00123B",
			Type:              "default",
			AllowedHosts:      "",
			ValidUntil:        &times[0],
		},
		{
			AccountIdentifier: account1.Identifier,
			Name:              "test2",
			Value:             "030001-1ACSDD-K2389A-22123B",
			Type:              "write",
			AllowedHosts:      "",
			ValidUntil:        &times[1],
		},
	}
	db.AddKey(keys[0])
	db.AddKey(keys[1])
	keys[0].Type = "write"
	keys[1].Name = "newtest1"
	keys[0].AllowedHosts = "test.lan,test.com,test.org"
	validTime := time.Now().Add(time.Minute * 30).Truncate(time.Second)
	keys[0].ValidUntil = &validTime
	err = db.UpdateKey(keys[0])
	if err != nil {
		t.Fatalf("Error updating key: %v", err)
	}
	key, _ := db.GetKey(keys[0].Value)
	if !key.Equal(&keys[0]) {
		t.Errorf("Expected key %+v, found %+v.", keys[0], *key)
	}
	if key.Name != keys[0].Name {
		t.Errorf("Expected key name to be %s, found %s.", keys[0].Name, key.Name)
	}
	keys[1].AccountIdentifier = account1.Identifier + 200
	keys[1].Value = "update-value-test"
	err = db.UpdateKey(keys[1])
	if err == nil {
		t.Error("Expected error from update with no changed values.")
	}
	key, _ = db.GetKey(keys[1].Value)
	if key != nil {
		t.Errorf("Found key with modified key value: %+v", key)
	}
}

func TestBadDatabaseKey(t *testing.T) {
	db := badTestSetup(t)
	_, err := db.GetAccountKeys("")
	if err == nil {
		t.Fatal("Expected error getting account keys.")
	}
	_, err = db.GetKey("")
	if err == nil {
		t.Fatal("Expected error getting key.")
	}
	_, err = db.AddKey(types.Key{})
	if err == nil {
		t.Fatal("Expected error adding key.")
	}
	err = db.DeleteKey(types.Key{})
	if err == nil {
		t.Fatal("Expected error deleting key.")
	}
	err = db.UpdateKey(types.Key{})
	if err == nil {
		t.Fatal("Expected error updating key.")
	}
}

func TestNoDatabaseKey(t *testing.T) {
	db := SQLite{}
	_, err := db.GetAccountKeys("")
	if err == nil {
		t.Fatal("Expected error getting account keys.")
	}
	_, err = db.GetKey("")
	if err == nil {
		t.Fatal("Expected error getting key.")
	}
	_, err = db.AddKey(types.Key{})
	if err == nil {
		t.Fatal("Expected error adding key.")
	}
	err = db.DeleteKey(types.Key{})
	if err == nil {
		t.Fatal("Expected error deleting key.")
	}
	err = db.UpdateKey(types.Key{})
	if err == nil {
		t.Fatal("Expected error updating key.")
	}
}
