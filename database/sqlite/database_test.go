package sqlite

import (
	"chronokeep/results/auth"
	"chronokeep/results/database"
	"chronokeep/results/types"
	"chronokeep/results/util"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

const (
	dbName     = "./results_test.sqlite"
	dbHost     = ""
	dbUser     = ""
	dbPassword = ""
	dbPort     = 0
	dbDriver   = "sqlite3"
)

func testHashPassword(pass string) string {
	hash, _ := auth.HashPassword(pass)
	return hash
}

func badTestSetup(t *testing.T) *SQLite {
	t.Log("Setting up bad test variables.")
	o := SQLite{}
	config := getTestConfig()
	config.DBName = "InvalidDatabaseName.sqlite"
	o.GetDatabase(config)
	return &o
}

func setupTests(t *testing.T) (*SQLite, func(t *testing.T), error) {
	t.Log("Setting up testing database variables.")
	o := SQLite{}
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

func setupOld() (*SQLite, error) {
	o := SQLite{}
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
				"account_id INTEGER PRIMARY KEY AUTOINCREMENT, " +
				"account_name VARCHAR(100) NOT NULL, " +
				"account_email VARCHAR(100) NOT NULL, " +
				"account_password VARCHAR(300) NOT NULL, " +
				"account_type VARCHAR(20) NOT NULL, " +
				"account_wrong_pass INT NOT NULL DEFAULT 0, " +
				"account_locked BOOL DEFAULT FALSE, " +
				"account_token VARCHAR(1000) NOT NULL DEFAULT '', " +
				"account_refresh_token VARCHAR(1000) NOT NULL DEFAULT '', " +
				"account_created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
				"account_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP," +
				"account_deleted BOOL DEFAULT FALSE, " +
				"UNIQUE(account_email)" +
				");",
		},
		// KEY TABLE
		{
			name: "KeyTable",
			query: "CREATE TABLE IF NOT EXISTS api_key(" +
				"account_id INTEGER NOT NULL, " +
				"key_name VARCHAR(100) NOT NULL DEFAULT ''," +
				"key_value VARCHAR(100) NOT NULL, " +
				"key_type VARCHAR(20) NOT NULL, " +
				"allowed_hosts TEXT, " +
				"valid_until DATETIME DEFAULT NULL, " +
				"key_created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
				"key_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
				"key_deleted BOOL DEFAULT FALSE, " +
				"UNIQUE(key_value), " +
				"FOREIGN KEY (account_id) REFERENCES account(account_id)" +
				");",
		},
		// EVENT TABLE
		{
			name: "EventTable",
			query: "CREATE TABLE IF NOT EXISTS event(" +
				"event_id INTEGER PRIMARY KEY AUTOINCREMENT, " +
				"account_id BIGINT NOT NULL, " +
				"event_name VARCHAR(100) NOT NULL, " +
				"slug VARCHAR(20) NOT NULL, " +
				"website VARCHAR(200), " +
				"image VARCHAR(200), " +
				"contact_email VARCHAR(100), " +
				"access_restricted BOOL DEFAULT FALSE, " +
				"event_created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
				"event_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP," +
				"event_deleted BOOL DEFAULT FALSE, " +
				"event_type VARCHAR(20) DEFAULT 'distance', " +
				"UNIQUE(event_name), " +
				"UNIQUE(slug)," +
				"FOREIGN KEY (account_id) REFERENCES account(account_id)" +
				");",
		},
		// EVENT YEAR TABLE
		{
			name: "EventYearTable",
			query: "CREATE TABLE IF NOT EXISTS event_year(" +
				"event_year_id INTEGER PRIMARY KEY AUTOINCREMENT, " +
				"event_id BIGINT NOT NULL, " +
				"year VARCHAR(20) NOT NULL, " +
				"date_time DATETIME NOT NULL, " +
				"live BOOL DEFAULT FALSE, " +
				"year_created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
				"year_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP," +
				"year_deleted BOOL DEFAULT FALSE, " +
				"CONSTRAINT year_slug UNIQUE (event_id, year)," +
				"FOREIGN KEY (event_id) REFERENCES event(event_id)" +
				");",
		},
		// PERSON TABLE
		{
			name: "PersonTable",
			query: "CREATE TABLE IF NOT EXISTS person(" +
				"person_id INTEGER PRIMARY KEY AUTOINCREMENT, " +
				"event_year_id BIGINT NOT NULL, " +
				"bib VARCHAR(100) NOT NULL, " +
				"first VARCHAR(100) NOT NULL, " +
				"last VARCHAR(100) NOT NULL, " +
				"age INT NOT NULL, " +
				"gender CHAR(1) NOT NULL, " +
				"age_group VARCHAR(200), " +
				"distance VARCHAR(200) NOT NULL, " +
				"CONSTRAINT one_person UNIQUE (event_year_id, bib), " +
				"FOREIGN KEY (event_year_id) REFERENCES event_year(event_year_id)" +
				");",
		},
		// RESULT TABLE
		{
			name: "ResultTable",
			query: "CREATE TABLE IF NOT EXISTS result(" +
				"person_id BIGINT NOT NULL, " +
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
				"result_type INT DEFAULT 0, " +
				"result_created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
				"result_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP," +
				"CONSTRAINT one_occurrence_res UNIQUE (person_id, location, occurence)," +
				"FOREIGN KEY (person_id) REFERENCES person(person_id)" +
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
			query: "CREATE TRIGGER UpdateAccountTime UPDATE OF account_name, account_email, account_password, " +
				"account_type, account_wrong_pass, account_locked, account_deleted ON account " +
				"BEGIN" +
				"    UPDATE account SET account_updated_at=CURRENT_TIMESTAMP WHERE account_id=account_id;" +
				"END;",
		},
		// UPDATE KEY FUNC
		{
			name: "UpdateKeyFunc",
			query: "CREATE TRIGGER UpdateKeyTime UPDATE OF account_id, key_name, key_type, allowed_hosts, " +
				"valid_until, key_deleted ON api_key " +
				"BEGIN" +
				"    UPDATE api_key SET key_updated_at=CURRENT_TIMESTAMP WHERE key_value=key_value;" +
				"END;",
		},
		// UPDATE EVENT FUNC
		{
			name: "UpdateEventFunc",
			query: "CREATE TRIGGER UpdateEventTime UPDATE OF account_id, event_name, slug, website, image, " +
				"contact_email, access_restricted, event_deleted, event_type ON event " +
				"BEGIN" +
				"    UPDATE event SET event_updated_at=CURRENT_TIMESTAMP WHERE event_id=event_id;" +
				"END;",
		},
		// UPDATE EVENT YEAR FUNC
		{
			name: "UpdateEventYearFunc",
			query: "CREATE TRIGGER UpdateEventYearTime UPDATE OF year, date_time, live, year_deleted ON event_year " +
				"BEGIN" +
				"    UPDATE event_year SET year_updated_at=CURRENT_TIMESTAMP WHERE event_year_id=event_year_id;" +
				"END;",
		},
		// UPDATE RESULT FUNC
		{
			name: "UpdateResultFunc",
			query: "CREATE TRIGGER UpdateResultTime UPDATE OF person_id, seconds, milliseconds, chip_seconds, " +
				"chip_milliseconds, segment, location, occurence, ranking, age_ranking, gender_ranking, finish, " +
				"result_type ON result " +
				"BEGIN" +
				"    UPDATE result SET result_updated_at=CURRENT_TIMESTAMP WHERE person_id=person_id AND location=location AND occurence=occurence;" +
				"END;",
		},
	}

	// Get a context and cancel function to create our tables, defer the cancel until we're done.
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()

	for _, single := range queries {
		_, err := db.ExecContext(ctx, single.query)
		if err != nil {
			return nil, fmt.Errorf("error executing %s query: %v", single.name, err)
		}
	}

	o.SetSetting("version", "5")

	return &o, nil
}

func TestSetupAndGet(t *testing.T) {
	t.Log("Setting up testing database variables.")
	o := &SQLite{}
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
	// Set up some basic information in the database to ensure we can
	// upgrade with existing data.
	t.Log("Adding account.")
	account1 := &types.Account{
		Name:     "John Smith",
		Email:    "j@test.com",
		Type:     "admin",
		Password: testHashPassword("password"),
	}
	_, _ = db.AddAccount(*account1)
	account1, err = db.GetAccount(account1.Email)
	if err != nil {
		t.Fatalf("Error adding account: %v", err)
	}
	t.Log("Adding Event.")
	event1 := &types.Event{
		AccountIdentifier: account1.Identifier,
		Name:              "Event 1",
		Slug:              "event1",
		ContactEmail:      "event1@test.com",
		AccessRestricted:  false,
	}
	_, _ = db.oldAddEvent(*event1)
	event1, err = db.oldGetEvent(event1.Slug)
	if err != nil {
		t.Fatalf("Error adding event: %v", err)
	}
	t.Log("Adding EventYear.")
	eventYear1 := &types.EventYear{
		EventIdentifier: event1.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 10, 06, 9, 6, 3, 15, time.Local),
		Live:            false,
		DaysAllowed:     1,
		RankingType:     "chip",
	}
	_, _ = db.oldAddEventYear(*eventYear1)
	eventYear1, err = db.oldGetEventYear(event1.Slug, eventYear1.Year)
	if err != nil {
		t.Fatalf("Error adding event year: %v", err)
	}
	t.Log("Adding results.")
	results := []types.Result{
		{
			Bib:           "100",
			First:         "John",
			Last:          "Smith",
			Age:           24,
			Gender:        "M",
			AgeGroup:      "20-29",
			Distance:      "1 Mile",
			Seconds:       377,
			Milliseconds:  0,
			Segment:       "",
			Location:      "Start/Finish",
			Occurence:     1,
			Ranking:       1,
			AgeRanking:    1,
			GenderRanking: 1,
			Finish:        false,
			Anonymous:     true,
		},
	}
	_, _ = db.AddResults(eventYear1.Identifier, results)
	// Verify version 5
	version := db.checkVersion()
	if version != 5 {
		t.Fatalf("Version set to '%v' expected '5'.", version)
	}
	// Verify version 6
	err = db.updateTables(version, 6)
	if err != nil {
		t.Fatalf("error updating database from %d to %d: %v", version, 6, err)
	}
	version = db.checkVersion()
	if version != 6 {
		t.Fatalf("Version set to '%v' expected '6'.", version)
	}
	// Verify version 7
	err = db.updateTables(version, 7)
	if err != nil {
		t.Fatalf("error updating database from %d to %d: %v", version, 7, err)
	}
	version = db.checkVersion()
	if version != 7 {
		t.Fatalf("Version set to '%v' expected '7'.", version)
	}
	// Verify version 8
	err = db.updateTables(version, 8)
	if err != nil {
		t.Fatalf("error updating database from %d to %d: %v", version, 8, err)
	}
	version = db.checkVersion()
	if version != 8 {
		t.Fatalf("Version set to '%v' expected '8'.", version)
	}
	// Verify version 9
	err = db.updateTables(version, 9)
	if err != nil {
		t.Fatalf("error updating database from %d to %d: %v", version, 9, err)
	}
	version = db.checkVersion()
	if version != 9 {
		t.Fatalf("Version set to '%v' expected '9'.", version)
	}
	// Verify version 10
	err = db.updateTables(version, 10)
	if err != nil {
		t.Fatalf("error updating database from %d to %d: %v", version, 10, err)
	}
	version = db.checkVersion()
	if version != 10 {
		t.Fatalf("Version set to '%v' expected '10'.", version)
	}
	// Verify version 11
	err = db.updateTables(version, 11)
	if err != nil {
		t.Fatalf("error updating database from %d to %d: %v", version, 11, err)
	}
	version = db.checkVersion()
	if version != 11 {
		t.Fatalf("Version set to '%v' expected '11'.", version)
	}
	// Verify version 12
	err = db.updateTables(version, 12)
	if err != nil {
		t.Fatalf("error updating database from %d to %d: %v", version, 12, err)
	}
	version = db.checkVersion()
	if version != 12 {
		t.Fatalf("Version set to '%v' expected '12'.", version)
	}
	// Verify version 13
	err = db.updateTables(version, 13)
	if err != nil {
		t.Fatalf("error updating database from %d to %d: %v", version, 13, err)
	}
	version = db.checkVersion()
	if version != 13 {
		t.Fatalf("Version set to '%v' expected '13'.", version)
	}
	// Verify version 14
	err = db.updateTables(version, 14)
	if err != nil {
		t.Fatalf("error updating database from %d to %d: %v", version, 14, err)
	}
	version = db.checkVersion()
	if version != 14 {
		t.Fatalf("Version set to '%v' expected '14'.", version)
	}
	// Verify version 15
	err = db.updateTables(version, 15)
	if err != nil {
		t.Fatalf("error updating database from %d to %d: %v", version, 15, err)
	}
	version = db.checkVersion()
	if version != 15 {
		t.Fatalf("Version set to '%v' expected '15'.", version)
	}
	// Verify version 16
	err = db.updateTables(version, 16)
	if err != nil {
		t.Fatalf("error updating database from %d to %d: %v", version, 16, err)
	}
	version = db.checkVersion()
	if version != 16 {
		t.Fatalf("Version set to '%v' expected '16'.", version)
	}
	// Verify version 17
	err = db.updateTables(version, 17)
	if err != nil {
		t.Fatalf("error updating database from %d to %d: %v", version, 17, err)
	}
	version = db.checkVersion()
	if version != 17 {
		t.Fatalf("Version set to '%v' expected '17'.", version)
	}
	// Verify version 18
	err = db.updateTables(version, 18)
	if err != nil {
		t.Fatalf("error updating database from %d to %d: %v", version, 18, err)
	}
	version = db.checkVersion()
	if version != 18 {
		t.Fatalf("Version set to '%v' expected '18'.", version)
	}
	// Verify version 19
	err = db.updateTables(version, 19)
	if err != nil {
		t.Fatalf("error updating database from %d to %d: %v", version, 19, err)
	}
	version = db.checkVersion()
	if version != 19 {
		t.Fatalf("Version set to '%v' expected '19'.", version)
	}
	// Check for error on drop tables as well. Because we can.
	err = db.dropTables()
	if err != nil {
		t.Fatalf("error deleting database: %v", err)
	}
}

func TestNoDatabase(t *testing.T) {
	db := SQLite{}
	_, err := db.GetDatabase(nil)
	if err == nil {
		t.Fatal("Expected error getting database.")
	}
	db = SQLite{}
	_, err = db.GetDB()
	if err == nil {
		t.Fatal("Expected error getting database.")
	}
	db = SQLite{}
	err = db.Setup(&util.Config{})
	if err == nil {
		t.Fatal("Expected error in Setup.")
	}
	db = SQLite{}
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
