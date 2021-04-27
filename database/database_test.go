package database

import (
	"chronokeep/results/util"
	"fmt"
	"testing"
)

const (
	DBName     = "results_test"
	DBHost     = "localhost"
	DBUser     = "user"
	DBPassword = "password"
)

func setupTests(t *testing.T) (func(t *testing.T), error) {
	t.Log("Setting up testing database variables.")
	config, err := util.GetConfig()
	if err != nil {
		t.Log(fmt.Sprintf("Error setting up config values. %v", err))
		return nil, err
	}
	config.DBHost = DBHost
	config.DBName = DBName
	config.DBUser = DBUser
	config.DBPassword = DBPassword
	t.Log("Initializing database.")
	err = Setup(config)
	if err != nil {
		t.Log(fmt.Sprintf("Error initializing database. %v", err))
		return nil, err
	}
	t.Log("Database initialized.")
	return func(t *testing.T) {
		t.Log("Deleting old database.")
		err = deleteDatabase()
		if err != nil {
			t.Log(fmt.Sprintf("Error deleting database. %v", err))
			return
		}
		t.Log("Database successfully deleted.")
	}, nil
}

func TestSetup(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Log("setup error")
	}
	defer finalize(t)
}

func TestGetDatabase(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Log("setup error")
	}
	defer finalize(t)
}

func TestGetDB(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Log("setup error")
	}
	defer finalize(t)
}
