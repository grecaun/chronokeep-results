package postgres

import (
	"chronokeep/results/auth"
	"chronokeep/results/database"
	"chronokeep/results/util"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

const (
	dbName     = "results_test"
	dbHost     = "database.lan"
	dbUser     = "results_test"
	dbPassword = "results_test"
	dbPort     = 5432
	dbDriver   = "postgres"
)

func testHashPassword(pass string) string {
	hash, _ := auth.HashPassword(pass)
	return hash
}

func badTestSetup(t *testing.T) *Postgres {
	t.Log("Setting up bad test variables.")
	o := Postgres{}
	config := getTestConfig()
	config.DBName = "InvalidDatabaseName"
	o.GetDatabase(config)
	return &o
}

func setupTests(t *testing.T) (*Postgres, func(t *testing.T), error) {
	t.Log("Setting up testing database variables.")
	o := Postgres{}
	config := getTestConfig()
	t.Log("Initializing database.")
	// Connect to DB with database name.
	test, err := o.GetDatabase(config)
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to database: %v", err)
	}
	if test == nil {
		return nil, nil, errors.New("database returned was nil")
	}

	// Check our database version.
	dbVersion := o.checkVersion()

	// Error checking version, most likely means tables are not created.
	if dbVersion < 1 {
		err = o.createTables()
		if err != nil {
			return nil, nil, err
		}
		// Otherwise check if our database is out of date and update if necessary.
	} else if dbVersion < database.CurrentVersion {
		err = o.updateTables(dbVersion, database.CurrentVersion)
		if err != nil {
			return nil, nil, err
		}
	}
	t.Log("Database initialized.")
	return &o, func(t *testing.T) {
		t.Log("Deleting old database.")
		err = o.dropTables()
		if err != nil {
			t.Fatalf("Error deleting database. %v", err)
			return
		}
		t.Log("Database successfully deleted.")
	}, nil
}

func setupOld() (*Postgres, error) {
	o := Postgres{}
	config := getTestConfig()
	// Connect to DB with database name.
	db, err := o.GetDatabase(config)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}
	queries := []myQuery{
		// SETTINGS TABLE
		{
			name: "SettingsTable",
			query: "CREATE TABLE IF NOT EXISTS settings(" +
				"name VARCHAR(200) NOT NULL, " +
				"value VARCHAR(200) NOT NULL, " +
				"UNIQUE (name));",
		},
		// ACCOUNT TABLE
		{
			name: "AccountTable",
			query: "CREATE TABLE IF NOT EXISTS account(" +
				"account_id BIGSERIAL NOT NULL, " +
				"account_name VARCHAR(100) NOT NULL, " +
				"account_email VARCHAR(100) NOT NULL, " +
				"account_password VARCHAR(300) NOT NULL, " +
				"account_type VARCHAR(20) NOT NULL, " +
				"account_wrong_pass INT NOT NULL DEFAULT 0, " +
				"account_locked BOOL DEFAULT FALSE, " +
				"account_token VARCHAR(1000) NOT NULL DEFAULT '', " +
				"account_refresh_token VARCHAR(1000) NOT NULL DEFAULT '', " +
				"account_created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP, " +
				"account_updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP," +
				"account_deleted BOOL DEFAULT FALSE, " +
				"UNIQUE(account_email), " +
				"PRIMARY KEY (account_id)" +
				");",
		},
		// KEY TABLE
		{
			name: "KeyTable",
			query: "CREATE TABLE IF NOT EXISTS api_key(" +
				"account_id BIGINT NOT NULL, " +
				"key_value VARCHAR(100) NOT NULL, " +
				"key_type VARCHAR(20) NOT NULL, " +
				"allowed_hosts TEXT, " +
				"valid_until TIMESTAMPTZ DEFAULT NULL, " +
				"key_created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP, " +
				"key_updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP, " +
				"key_deleted BOOL DEFAULT FALSE, " +
				"UNIQUE(key_value), " +
				"FOREIGN KEY (account_id) REFERENCES account(account_id)" +
				");",
		},
		// EVENT TABLE
		{
			name: "EventTable",
			query: "CREATE TABLE IF NOT EXISTS event(" +
				"event_id BIGSERIAL NOT NULL, " +
				"account_id BIGINT NOT NULL, " +
				"event_name VARCHAR(100) NOT NULL, " +
				"slug VARCHAR(20) NOT NULL, " +
				"website VARCHAR(200), " +
				"image VARCHAR(200), " +
				"contact_email VARCHAR(100), " +
				"access_restricted BOOL DEFAULT FALSE, " +
				"event_created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP, " +
				"event_updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP," +
				"event_deleted BOOL DEFAULT FALSE, " +
				"UNIQUE(event_name), " +
				"UNIQUE(slug)," +
				"FOREIGN KEY (account_id) REFERENCES account(account_id)," +
				"PRIMARY KEY (event_id)" +
				");",
		},
		// EVENT YEAR TABLE
		{
			name: "EventYearTable",
			query: "CREATE TABLE IF NOT EXISTS event_year(" +
				"event_year_id BIGSERIAL NOT NULL, " +
				"event_id BIGINT NOT NULL, " +
				"year VARCHAR(20) NOT NULL, " +
				"date_time TIMESTAMPTZ NOT NULL, " +
				"live BOOL DEFAULT FALSE, " +
				"year_created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP, " +
				"year_updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP," +
				"year_deleted BOOL DEFAULT FALSE, " +
				"CONSTRAINT year_slug UNIQUE (event_id, year)," +
				"FOREIGN KEY (event_id) REFERENCES event(event_id)," +
				"PRIMARY KEY (event_year_id)" +
				");",
		},
		// RESULT TABLE
		{
			name: "ResultTable",
			query: "CREATE TABLE IF NOT EXISTS result(" +
				"event_year_id BIGINT NOT NULL, " +
				"bib VARCHAR(100) NOT NULL, " +
				"first VARCHAR(100) NOT NULL, " +
				"last VARCHAR(100) NOT NULL, " +
				"age INT NOT NULL, " +
				"gender CHAR(1) NOT NULL, " +
				"age_group VARCHAR(200), " +
				"distance VARCHAR(200) NOT NULL, " +
				"seconds INT DEFAULT 0, " +
				"milliseconds INT DEFAULT 0, " +
				"chip_seconds INT DEFAULT 0, " +
				"chip_milliseconds INT DEFAULT 0, " +
				"segment VARCHAR(500), " +
				"location VARCHAR(500), " +
				"occurence INT DEFAULT -1, " +
				"ranking INT DEFAULT -1, " +
				"age_ranking INT DEFAULT -1, " +
				"gender_ranking INT DEFAULT -1, " +
				"finish BOOL DEFAULT TRUE, " +
				"result_created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP, " +
				"result_updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP," +
				"CONSTRAINT one_occurrence UNIQUE (event_year_id, bib, location, occurence)," +
				"FOREIGN KEY (event_year_id) REFERENCES event_year(event_year_id)" +
				");",
		},
		// RECORD TABLE
		{
			name: "RecordTable",
			query: "CREATE TABLE IF NOT EXISTS call_record(" +
				"account_id BIGINT NOT NULL, " +
				"time BIGINT NOT NULL, " +
				"count INT DEFAULT 0, " +
				"CONSTRAINT account_time UNIQUE (account_id, time)," +
				"FOREIGN KEY (account_id) REFERENCES account(account_id)" +
				");",
		},
		// UPDATE ACCOUNT FUNC
		{
			name: "UpdateAccountFunc",
			query: "CREATE OR REPLACE FUNCTION account_timestamp_column() " +
				"RETURNS TRIGGER AS $$ " +
				"BEGIN " +
				"NEW.account_updated_at = now();" +
				"RETURN NEW;" +
				"END;" +
				"$$ language 'plpgsql';",
		},
		// UPDATE KEY FUNC
		{
			name: "UpdateKeyFunc",
			query: "CREATE OR REPLACE FUNCTION key_timestamp_column() " +
				"RETURNS TRIGGER AS $$ " +
				"BEGIN " +
				"NEW.key_updated_at = now();" +
				"RETURN NEW;" +
				"END;" +
				"$$ language 'plpgsql';",
		},
		// UPDATE EVENT FUNC
		{
			name: "UpdateEventFunc",
			query: "CREATE OR REPLACE FUNCTION event_timestamp_column() " +
				"RETURNS TRIGGER AS $$ " +
				"BEGIN " +
				"NEW.event_updated_at = now();" +
				"RETURN NEW;" +
				"END;" +
				"$$ language 'plpgsql';",
		},
		// UPDATE EVENT YEAR FUNC
		{
			name: "UpdateEventYearFunc",
			query: "CREATE OR REPLACE FUNCTION event_year_timestamp_column() " +
				"RETURNS TRIGGER AS $$ " +
				"BEGIN " +
				"NEW.year_updated_at = now();" +
				"RETURN NEW;" +
				"END;" +
				"$$ language 'plpgsql';",
		},
		// UPDATE RESULT FUNC
		{
			name: "UpdateResultFunc",
			query: "CREATE OR REPLACE FUNCTION result_timestamp_column() " +
				"RETURNS TRIGGER AS $$ " +
				"BEGIN " +
				"NEW.result_updated_at = now();" +
				"RETURN NEW;" +
				"END;" +
				"$$ language 'plpgsql';",
		},
		// TRIGGERS FOR UPDATING UPDATED_AT timestamps
		{
			name:  "AccountTableTrigger",
			query: "CREATE TRIGGER update_account_timestamp BEFORE UPDATE ON api_key FOR EACH ROW EXECUTE PROCEDURE account_timestamp_column();",
		},
		{
			name:  "KeyTableTrigger",
			query: "CREATE TRIGGER update_key_timestamp BEFORE UPDATE ON api_key FOR EACH ROW EXECUTE PROCEDURE key_timestamp_column();",
		},
		{
			name:  "EventTableTrigger",
			query: "CREATE TRIGGER update_event_timestamp BEFORE UPDATE ON api_key FOR EACH ROW EXECUTE PROCEDURE event_timestamp_column();",
		},
		{
			name:  "EventYearTableTrigger",
			query: "CREATE TRIGGER update_event_year_timestamp BEFORE UPDATE ON api_key FOR EACH ROW EXECUTE PROCEDURE event_year_timestamp_column();",
		},
		{
			name:  "ResultTableTrigger",
			query: "CREATE TRIGGER update_result_timestamp BEFORE UPDATE ON api_key FOR EACH ROW EXECUTE PROCEDURE result_timestamp_column();",
		},
	}

	// Get a context and cancel function to create our tables, defer the cancel until we're done.
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()

	for _, single := range queries {
		_, err := db.Exec(ctx, single.query)
		if err != nil {
			return nil, fmt.Errorf("error executing %s query: %v", single.name, err)
		}
	}

	o.SetSetting("version", "1")

	return &o, nil
}

func TestSetupAndGet(t *testing.T) {
	t.Log("Setting up testing database variables.")
	o := &Postgres{}
	config := getTestConfig()
	t.Log("Initializing database.")
	err := o.Setup(config)
	defer o.dropTables()
	if err != nil {
		t.Fatalf("Error initializing database. %v", err)
	}
	t.Log("Database initialized.")
	if o.db == nil {
		t.Fatalf("db variable not set")
	}
	o.db.Close()
	o.updateDB(nil)
	_, err = o.GetDatabase(config)
	if err != nil {
		t.Fatalf("error getting database with config values: %v", err)
	}
	o.db.Close()
	o.updateDB(nil)
	_, err = o.GetDB()
	if err != nil {
		t.Fatalf("error getting database without config values: %v", err)
	}
	_, err = o.GetDB()
	if err != nil {
		t.Fatalf("error getting database without config values: %v", err)
	}
	err = o.dropTables()
	if err != nil {
		t.Fatalf("error deleting database: %v", err)
	}
}

func TestCheckVersion(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	version := db.checkVersion()
	if version != database.CurrentVersion {
		t.Fatalf("version found '%v' expected '%v'", version, database.CurrentVersion)
	}
}

func TestUpgrade(t *testing.T) {
	t.Log("Setting up testing database variables.")
	t.Log("Initializing database version 1.")
	db, err := setupOld()
	defer db.dropTables()
	if err != nil {
		t.Fatalf("Error initializing database. %v", err)
	}
	t.Log("Database initialized.")
	if db == nil || db.db == nil {
		t.Fatalf("db variable not set")
	}
	// Verify version 1
	version := db.checkVersion()
	if version != 1 {
		t.Fatalf("Version set to '%v' expected '1'.", version)
	}
	// Verify update works.
	err = db.updateTables(version, 2)
	if err != nil {
		t.Fatalf("error updating database from %d to %d: %v", version, 2, err)
	}
	// Verify version 2
	version = db.checkVersion()
	if version != 2 {
		t.Fatalf("Version set to '%v' expected '2'.", version)
	}
	// Verify verison 3
	err = db.updateTables(version, 3)
	if err != nil {
		t.Fatalf("error updating database from %d to %d: %v", version, 3, err)
	}
	version = db.checkVersion()
	if version != 3 {
		t.Fatalf("Version set to '%v' expected '3'.", version)
	}
	// Verify verison 4
	err = db.updateTables(version, 4)
	if err != nil {
		t.Fatalf("error updating database from %d to %d: %v", version, 4, err)
	}
	version = db.checkVersion()
	if version != 4 {
		t.Fatalf("Version set to '%v' expected '4'.", version)
	}
	// Verify verison 5
	err = db.updateTables(version, 5)
	if err != nil {
		t.Fatalf("error updating database from %d to %d: %v", version, 5, err)
	}
	version = db.checkVersion()
	if version != 5 {
		t.Fatalf("Version set to '%v' expected '5'.", version)
	}
	// Check for error on drop tables as well. Because we can.
	err = db.dropTables()
	if err != nil {
		t.Fatalf("error deleting database: %v", err)
	}
}

func TestBadDatabase(t *testing.T) {
	db := badTestSetup(t)
	_, err := db.GetDatabase(&util.Config{})
	if err == nil {
		t.Fatal("Expected error getting database.")
	}
	db = badTestSetup(t)
	_, err = db.GetDB()
	if err == nil {
		t.Fatal("Expected error getting database.")
	}
	db = badTestSetup(t)
	err = db.Setup(&util.Config{})
	if err == nil {
		t.Fatal("Expected error in Setup.")
	}
	db = badTestSetup(t)
	err = db.dropTables()
	if err == nil {
		t.Fatal("Expected error dropping tables.")
	}
	err = db.SetSetting("", "")
	if err == nil {
		t.Fatal("Expected error setting setting.")
	}
	err = db.createTables()
	if err == nil {
		t.Fatal("Expected error creating tables.")
	}
	v := db.checkVersion()
	if v != -1 {
		t.Fatal("Expected error getting database.")
	}
	err = db.updateTables(0, 0)
	if err == nil {
		t.Fatal("Expected error updating tables.")
	}
}

func TestNoDatabase(t *testing.T) {
	db := Postgres{}
	_, err := db.GetDatabase(&util.Config{})
	if err == nil {
		t.Fatal("Expected error getting database.")
	}
	db = Postgres{}
	_, err = db.GetDB()
	if err == nil {
		t.Fatal("Expected error getting database.")
	}
	db = Postgres{}
	err = db.Setup(&util.Config{})
	if err == nil {
		t.Fatal("Expected error in Setup.")
	}
	db = Postgres{}
	err = db.dropTables()
	if err == nil {
		t.Fatal("Expected error dropping tables.")
	}
	err = db.SetSetting("", "")
	if err == nil {
		t.Fatal("Expected error setting setting.")
	}
	err = db.createTables()
	if err == nil {
		t.Fatal("Expected error creating tables.")
	}
	v := db.checkVersion()
	if v != -1 {
		t.Fatal("Expected error getting database.")
	}
	err = db.updateTables(0, 0)
	if err == nil {
		t.Fatal("Expected error updating tables.")
	}
}

func getTestConfig() *util.Config {
	return &util.Config{
		DBHost:     dbHost,
		DBPort:     dbPort,
		DBName:     dbName,
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBDriver:   dbDriver,
		AdminEmail: "admin@test.com",
		AdminName:  "tester number 1",
		AdminPass:  "password",
	}
}
