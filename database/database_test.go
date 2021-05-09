package database

import (
	"chronokeep/results/auth"
	"chronokeep/results/util"
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"
)

const (
	dbName     = "results_test"
	dbHost     = "localhost"
	dbUser     = "results_test"
	dbPassword = "results_test"
	dbPort     = 3306
	dbDriver   = "mysql"
)

func testHashPassword(pass string) string {
	hash, _ := auth.HashPassword(pass)
	return hash
}

func setupTests(t *testing.T) (func(t *testing.T), error) {
	t.Log("Setting up testing database variables.")
	config = getTestConfig()
	t.Log("Initializing database.")
	// Connect to DB with database name.
	_, err := GetDatabase(config)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	// Check our database version.
	dbVersion := checkVersion()

	// Error checking version, most likely means tables are not created.
	if dbVersion < 1 {
		err = createTables()
		if err != nil {
			return nil, err
		}
		// Otherwise check if our database is out of date and update if necessary.
	} else if dbVersion < CurrentVersion {
		err = updateTables(dbVersion, CurrentVersion)
		if err != nil {
			return nil, err
		}
	}
	t.Log("Database initialized.")
	return func(t *testing.T) {
		t.Log("Deleting old database.")
		err = dropTables()
		if err != nil {
			t.Fatalf("Error deleting database. %v", err)
			return
		}
		t.Log("Database successfully deleted.")
	}, nil
}

func setupOld(version int) error {
	// Connect to DB with database name.
	db, err := GetDatabase(config)
	if err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}
	var settingsTable, accountTable, keyTable, eventTable, eventYearTable, resultTable, recordTable string
	switch version {
	case 1:
		switch config.DBDriver {
		case "postgres":
			return errors.New("postgres not yet supported")
		case "mysql":

			settingsTable = "CREATE TABLE IF NOT EXISTS settings(" +
				"name VARCHAR(200) NOT NULL, " +
				"value VARCHAR(200) NOT NULL, " +
				"UNIQUE (name));"

			accountTable = "CREATE TABLE IF NOT EXISTS account(" +
				"account_id BIGINT NOT NULL AUTO_INCREMENT, " +
				"account_name VARCHAR(100) NOT NULL, " +
				"account_email VARCHAR(100) NOT NULL, " +
				"account_password VARCHAR(300) NOT NULL, " +
				"account_type VARCHAR(20) NOT NULL, " +
				"account_wrong_pass INT NOT NULL DEFAULT 0, " +
				"account_locked BOOL DEFAULT FALSE, " +
				"account_created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
				"account_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
				"account_deleted BOOL DEFAULT FALSE, " +
				"UNIQUE(account_email), " +
				"PRIMARY KEY (account_id)" +
				");"

			keyTable = "CREATE TABLE IF NOT EXISTS api_key(" +
				"account_id BIGINT NOT NULL, " +
				"key_value CHAR(100) NOT NULL, " +
				"key_type VARCHAR(20) NOT NULL, " +
				"allowed_hosts TEXT, " +
				"valid_until DATETIME DEFAULT CURRENT_TIMESTAMP, " +
				"key_created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
				"key_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, " +
				"key_deleted BOOL DEFAULT FALSE, " +
				"UNIQUE(key_value), " +
				"FOREIGN KEY (account_id) REFERENCES account(account_id)" +
				");"

			eventTable = "CREATE TABLE IF NOT EXISTS event(" +
				"event_id BIGINT NOT NULL AUTO_INCREMENT, " +
				"account_id BIGINT NOT NULL, " +
				"event_name VARCHAR(100) NOT NULL, " +
				"slug VARCHAR(20) NOT NULL, " +
				"website VARCHAR(200), " +
				"image VARCHAR(200), " +
				"contact_email VARCHAR(100), " +
				"access_restricted BOOL DEFAULT FALSE, " +
				"event_created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
				"event_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
				"event_deleted BOOL DEFAULT FALSE, " +
				"UNIQUE(event_name), " +
				"UNIQUE(slug)," +
				"FOREIGN KEY (account_id) REFERENCES account(account_id)," +
				"PRIMARY KEY (event_id)" +
				");"

			eventYearTable = "CREATE TABLE IF NOT EXISTS event_year(" +
				"event_year_id BIGINT NOT NULL AUTO_INCREMENT, " +
				"event_id BIGINT NOT NULL, " +
				"year VARCHAR(20) NOT NULL, " +
				"date_time DATETIME NOT NULL, " +
				"live BOOL DEFAULT FALSE, " +
				"year_created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
				"year_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
				"year_deleted BOOL DEFAULT FALSE, " +
				"CONSTRAINT year_slug UNIQUE (event_id, year)," +
				"FOREIGN KEY (event_id) REFERENCES event(event_id)," +
				"PRIMARY KEY (event_year_id)" +
				");"

			resultTable = "CREATE TABLE IF NOT EXISTS result(" +
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
				"segment VARCHAR(500), " +
				"location VARCHAR(500), " +
				"occurence INT DEFAULT -1, " +
				"ranking INT DEFAULT -1, " +
				"age_ranking INT DEFAULT -1, " +
				"gender_ranking INT DEFAULT -1, " +
				"finish BOOL DEFAULT TRUE, " +
				"result_created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
				"result_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
				"CONSTRAINT one_occurrence UNIQUE (event_year_id, bib, location, occurence)," +
				"FOREIGN KEY (event_year_id) REFERENCES event_year(event_year_id)" +
				");"

			recordTable = "CREATE TABLE IF NOT EXISTS call_record(" +
				"account_id BIGINT NOT NULL, " +
				"time BIGINT NOT NULL, " +
				"count INT DEFAULT 0, " +
				"CONSTRAINT account_time UNIQUE (account_id, time)," +
				"FOREIGN KEY (account_id) REFERENCES account(account_id)" +
				");"
		default:
			return errors.New("invalid database type given")
		}
	default:
		return errors.New("invalid version specified")
	}

	// Get a context and cancel function to create our tables, defer the cancel until we're done.
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()

	_, err = db.ExecContext(ctx, settingsTable)
	if err != nil {
		return fmt.Errorf("error creating settings table: %v", err)
	}

	_, err = db.ExecContext(ctx, accountTable)
	if err != nil {
		return fmt.Errorf("error creating account table: %v", err)
	}

	_, err = db.ExecContext(ctx, keyTable)
	if err != nil {
		return fmt.Errorf("error creating key table: %v", err)
	}

	_, err = db.ExecContext(ctx, eventTable)
	if err != nil {
		return fmt.Errorf("error creating event table: %v", err)
	}

	_, err = db.ExecContext(ctx, eventYearTable)
	if err != nil {
		return fmt.Errorf("error creating event year table: %v", err)
	}

	_, err = db.ExecContext(ctx, resultTable)
	if err != nil {
		return fmt.Errorf("error creating result table: %v", err)
	}

	_, err = db.ExecContext(ctx, recordTable)
	if err != nil {
		return fmt.Errorf("error creating record table: %v", err)
	}

	SetSetting("version", strconv.Itoa(version))

	return nil
}

func TestSetupAndGet(t *testing.T) {
	t.Log("Setting up testing database variables.")
	config = getTestConfig()
	t.Log("Initializing database.")
	err := Setup(config)
	defer dropTables()
	if err != nil {
		t.Fatalf("Error initializing database. %v", err)
	}
	t.Log("Database initialized.")
	if db == nil {
		t.Fatalf("db variable not set")
	}
	db.Close()
	updateDB(nil)
	_, err = GetDatabase(config)
	if err != nil {
		t.Fatalf("error getting database with config values: %v", err)
	}
	db.Close()
	updateDB(nil)
	_, err = GetDB()
	if err != nil {
		t.Fatalf("error getting database without config values: %v", err)
	}
	_, err = GetDB()
	if err != nil {
		t.Fatalf("error getting database without config values: %v", err)
	}
	err = dropTables()
	if err != nil {
		t.Fatalf("error deleting database: %v", err)
	}
}

func TestCheckVersion(t *testing.T) {
	finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer finalize(t)
	version := checkVersion()
	if version != CurrentVersion {
		t.Fatalf("version found '%v' expected '%v'", version, CurrentVersion)
	}
}

func TestUpgrade(t *testing.T) {
	t.Log("Setting up testing database variables.")
	config = getTestConfig()
	t.Log("Initializing database version 1.")
	err := setupOld(1)
	defer dropTables()
	if err != nil {
		t.Fatalf("Error initializing database. %v", err)
	}
	t.Log("Database initialized.")
	if db == nil {
		t.Fatalf("db variable not set")
	}
	// Verify version 1
	version := checkVersion()
	if version != 1 {
		t.Fatalf("Version set to '%v' expected '1'.", version)
	}
	// In the future this will verify updates work properly.
	// Check for error on drop tables as well. Because we can.
	err = dropTables()
	if err != nil {
		t.Fatalf("error deleting database: %v", err)
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
