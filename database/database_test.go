package database

import (
	"chronokeep/results/util"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

const (
	dbName     = "results_test"
	dbHost     = "localhost"
	dbUser     = "user"
	dbPassword = "password"
)

func setupTests(t *testing.T) (func(t *testing.T), error) {
	t.Log("Setting up testing database variables.")
	config := &util.Config{
		DBHost:     dbHost,
		DBName:     dbName,
		DBUser:     dbUser,
		DBPassword: dbPassword,
	}
	t.Log("Creating database.")
	err := createDatabase()
	if err != nil {
		return nil, err
	}
	t.Log("Initializing database.")
	err = Setup(config)
	if err != nil {
		return nil, err
	}
	t.Log("Database initialized.")
	return func(t *testing.T) {
		t.Log("Deleting old database.")
		err = deleteDatabase()
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
				"UNIQUE (name) ON CONFLICT UPDATE);"

			accountTable = "CREATE TABLE IF NOT EXISTS account(" +
				"account_id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT, " +
				"name VARCHAR(100), " +
				"email VARCHAR(100), " +
				"type VARCHAR(20), " +
				"created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
				"updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
				"deleted BOOL DEFAULT FALSE, " +
				"UNIQUE(email));"

			keyTable = "CREATE TABLE IF NOT EXISTS key(" +
				"key_id BIGINT PRIMARY KEY AUTO_INCREMENT, " +
				"account_id BIGINT FOREIGN KEY REFERENCES account(account_id), " +
				"value CHAR(100) NOT NULL, " +
				"type VARCHAR(20) NOT NULL, " +
				"allowed_hosts TEXT, " +
				"valid_until DATETIME DEFAULT CURRENT_TIMESTAMP," +
				"created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
				"updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
				"deleted BOOL DEFAULT FALSE," +
				"UNIQUE(value));"

			eventTable = "CREATE TABLE IF NOT EXISTS event(" +
				"event_id BIGINT PRIMARY KEY AUTO_INCREMENT, " +
				"name VARCHAR(100) NOT NULL, " +
				"slug VARCHAR(20) NOT NULL, " +
				"website VARCHAR(200), " +
				"image VARCHAR(200), " +
				"contact_email VARCHAR(100), " +
				"account_id BIGINT FOREIGN KEY REFERENCES account(account_id), " +
				"access_restricted BOOL DEFAULT FALSE, " +
				"created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
				"updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
				"deleted BOOL DEFAULT FALSE, " +
				"UNIQUE(name), " +
				"UNIQUE(slug)" +
				");"

			eventYearTable = "CREATE TABLE IF NOT EXISTS event_year(" +
				"event_year_id BIGINT PRIMARY KEY AUTO_INCREMENT, " +
				"event_id BIGINT FOREIGN KEY REFERENCES event(event_id), " +
				"year VARCHAR(20) NOT NULL, " +
				"date DATE NOT NULL, " +
				"time TIME NOT NULL, " +
				"live BOOL DEFAULT FALSE, " +
				"created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
				"updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
				"deleted BOOL DEFAULT FALSE, " +
				"CONSTRAINT year_slug UNIQUE (event_id, year)" +
				");"

			resultTable = "CREATE TABLE IF NOT EXISTS result(" +
				"event_year_id BIGINT FOREIGN KEY REFERENCES event_year(event_year_id), " +
				"bib VARCHAR(100) NOT NULL, " +
				"first VARCHAR(100) NOT NULL, " +
				"last VARCHAR(100) NOT NULL, " +
				"age INT NOT NULL, " +
				"gender CHAR(1) NOT NULL, " +
				"age_group VARCHAR(200), " +
				"distance VARCHAR(200) NOT NULL, " +
				"seconds INT DEFAULT 0, " +
				"milliseconds INT DEFAULT 0, " +
				"location VARCHAR(500), " +
				"occurence INT DEFAULT -1, " +
				"ranking INT DEFAULT -1, " +
				"age_ranking INT DEFUALT -1, " +
				"gender_ranking INT DEFAULT -1, " +
				"finish BOOL DEFAULT TRUE, " +
				"created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
				"updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
				"CONSTRAINT one_finish UNIQUE (event_year_id, bib, finish) ON CONFLICT UPDATE, " +
				"CONSTRAINT one_occurrence UNIQUE (event_year_id, bib, location, occurence) ON CONFLICT UPDATE" +
				");"

			recordTable = "CREATE TABLE IF NOT EXISTS call_record(" +
				"account_id BIGINT FOREIGN KEY REFERENCES account(account_id), " +
				"time DATETIME NOT NULL, " +
				"count INT DEFAULT 0, " +
				"CONSTRAINT account_time UNIQUE (account_id, time) ON CONFLICT UPDATE" +
				");"
		default:
			return errors.New("invalid database type given")
		}
	default:
		return errors.New("invalid version specified")
	}

	settingsValue := fmt.Sprintf("INSERT INTO settings(name, value) VALUES ('version', '%v');", version)

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

	_, err = db.ExecContext(ctx, settingsValue)
	if err != nil {
		return fmt.Errorf("error adding settings: %v", err)
	}

	return nil
}

func TestSetupAndGet(t *testing.T) {
	t.Log("Setting up testing database variables.")
	config = &util.Config{
		DBHost:     dbHost,
		DBName:     dbName,
		DBUser:     dbUser,
		DBPassword: dbPassword,
	}
	t.Log("Creating database.")
	err := createDatabase()
	if err != nil {
		t.Fatalf("Error creating database. %v", err)
	}
	t.Log("Initializing database.")
	err = Setup(config)
	if err != nil {
		t.Fatalf("Error initializing database. %v", err)
	}
	t.Log("Database initialized.")
	if db == nil {
		t.Fatalf("db variable not set")
	}
	defer deleteDatabase()
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
	err = deleteDatabase()
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
	config = &util.Config{
		DBHost:     dbHost,
		DBName:     dbName,
		DBUser:     dbUser,
		DBPassword: dbPassword,
	}
	err := createDatabase()
	if err != nil {
		t.Fatalf("Error creating database. %v", err)
	}
	defer deleteDatabase()
	t.Log("Initializing database version 1.")
	err = setupOld(1)
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
}
