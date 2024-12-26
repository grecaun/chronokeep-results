package mysql

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
	dbName     = "results_test"
	dbHost     = "database.lan"
	dbUser     = "results_test"
	dbPassword = "results_test"
	dbPort     = 3306
	dbDriver   = "mysql"
)

func testHashPassword(pass string) string {
	hash, _ := auth.HashPassword(pass)
	return hash
}

func badTestSetup(t *testing.T) *MySQL {
	t.Log("Setting up bad test variables.")
	o := MySQL{}
	config := getTestConfig()
	config.DBName = "InvalidDatabaseName"
	o.GetDatabase(config)
	return &o
}

func setupTests(t *testing.T) (*MySQL, func(t *testing.T), error) {
	t.Log("Setting up testing database variables.")
	o := MySQL{}
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

func oldAddEvent(m *MySQL, event types.Event) (*types.Event, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"INSERT INTO event(event_name, slug, website, image, contact_email, account_id, access_restricted) VALUES (?, ?, ?, ?, ?, ?, ?);",
		event.Name,
		event.Slug,
		event.Website,
		event.Image,
		event.ContactEmail,
		event.AccountIdentifier,
		event.AccessRestricted,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to add event: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("unable to determine ID for event: %v", err)
	}
	return &types.Event{
		Identifier:        id,
		AccountIdentifier: event.AccountIdentifier,
		Name:              event.Name,
		Slug:              event.Slug,
		Website:           event.Website,
		Image:             event.Image,
		ContactEmail:      event.ContactEmail,
		AccessRestricted:  event.AccessRestricted,
		Type:              "",
	}, nil
}

func oldGetEvent(m *MySQL, slug string) (*types.Event, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT event_id, event_name, slug, website, image, account_id, contact_email, access_restricted, "+
			"recent_time FROM event NATURAL JOIN (SELECT e.event_id, MAX(y.date_time) AS recent_time FROM event e LEFT OUTER "+
			"JOIN (SELECT event_id, date_time FROM event_year WHERE year_deleted=FALSE) y ON e.event_id=y.event_id GROUP BY e.event_id) AS time WHERE event_deleted=FALSE and slug=?;",
		slug,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving event: %v", err)
	}
	defer res.Close()
	if res.Next() {
		var outEvent types.Event
		err := res.Scan(
			&outEvent.Identifier,
			&outEvent.Name,
			&outEvent.Slug,
			&outEvent.Website,
			&outEvent.Image,
			&outEvent.AccountIdentifier,
			&outEvent.ContactEmail,
			&outEvent.AccessRestricted,
			&outEvent.RecentTime,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting event: %v", err)
		}
		return &outEvent, nil
	}
	return nil, nil
}

func oldAddResults(m *MySQL, eventYearID int64, results []types.Result) ([]types.Result, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("unable to begin transaction to add results: %v", err)
	}
	stmt, err := tx.PrepareContext(
		ctx,
		"INSERT INTO result("+
			"event_year_id, "+
			"bib, "+
			"first, "+
			"last, "+
			"age, "+
			"gender, "+
			"age_group, "+
			"distance, "+
			"seconds, "+
			"milliseconds, "+
			"chip_seconds, "+
			"chip_milliseconds, "+
			"segment, "+
			"location, "+
			"occurence, "+
			"ranking, "+
			"age_ranking, "+
			"gender_ranking, "+
			"finish"+
			") "+
			" VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) "+
			"ON DUPLICATE KEY UPDATE "+
			"first=VALUES(first), "+
			"last=VALUES(last), "+
			"age=VALUES(age), "+
			"gender=VALUES(gender), "+
			"age_group=VALUES(age_group), "+
			"distance=VALUES(distance), "+
			"seconds=VALUES(seconds), "+
			"milliseconds=VALUES(milliseconds), "+
			"chip_seconds=VALUES(chip_seconds), "+
			"chip_milliseconds=VALUES(chip_milliseconds), "+
			"segment=VALUES(segment), "+
			"ranking=VALUES(ranking), "+
			"age_ranking=VALUES(age_ranking), "+
			"gender_ranking=VALUES(gender_ranking), "+
			"finish=VALUES(finish);",
	)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare statement for result add: %v", err)
	}
	for _, result := range results {
		_, err = stmt.ExecContext(
			ctx,
			eventYearID,
			result.Bib,
			result.First,
			result.Last,
			result.Age,
			result.Gender,
			result.AgeGroup,
			result.Distance,
			result.Seconds,
			result.Milliseconds,
			result.ChipSeconds,
			result.ChipMilliseconds,
			result.Segment,
			result.Location,
			result.Occurence,
			result.Ranking,
			result.AgeRanking,
			result.GenderRanking,
			result.Finish,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return results, nil
}

func setupOld() (*MySQL, error) {
	o := MySQL{}
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
				"account_id BIGINT NOT NULL AUTO_INCREMENT, " +
				"account_name VARCHAR(100) NOT NULL, " +
				"account_email VARCHAR(100) NOT NULL, " +
				"account_password VARCHAR(300) NOT NULL, " +
				"account_type VARCHAR(20) NOT NULL, " +
				"account_wrong_pass INT NOT NULL DEFAULT 0, " +
				"account_locked BOOL DEFAULT FALSE, " +
				"account_token VARCHAR(1000) NOT NULL DEFAULT '', " +
				"account_refresh_token VARCHAR(1000) NOT NULL DEFAULT '', " +
				"account_created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
				"account_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
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
				"valid_until DATETIME DEFAULT NULL, " +
				"key_created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
				"key_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, " +
				"key_deleted BOOL DEFAULT FALSE, " +
				"UNIQUE(key_value), " +
				"FOREIGN KEY (account_id) REFERENCES account(account_id)" +
				");",
		},
		// EVENT TABLE
		{
			name: "EventTable",
			query: "CREATE TABLE IF NOT EXISTS event(" +
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
				");",
		},
		// EVENT YEAR TABLE
		{
			name: "EventYearTable",
			query: "CREATE TABLE IF NOT EXISTS event_year(" +
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
				"result_created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
				"result_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
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

	o.SetSetting("version", "1")

	return &o, nil
}

func TestSetupAndGet(t *testing.T) {
	t.Log("Setting up testing database variables.")
	o := &MySQL{}
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
	_, _ = oldAddEvent(db, *event1)
	event1, err = oldGetEvent(db, event1.Slug)
	if err != nil {
		t.Fatalf("Error adding event: %v", err)
	}
	t.Log("Adding EventYear.")
	eventYear1 := &types.EventYear{
		EventIdentifier: event1.Identifier,
		Year:            "2021",
		DateTime:        time.Date(2021, 10, 06, 9, 6, 3, 15, time.Local),
		Live:            false,
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
	_, err = oldAddResults(db, eventYear1.Identifier, results)
	if err != nil {
		t.Fatalf("Error adding results: %v", err)
	}
	t.Log("Testing upgrades.")
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
	// Verify version 3
	err = db.updateTables(version, 3)
	if err != nil {
		t.Fatalf("error updating database from %d to %d: %v", version, 3, err)
	}
	version = db.checkVersion()
	if version != 3 {
		t.Fatalf("Version set to '%v' expected '3'.", version)
	}
	// Verify version 4
	err = db.updateTables(version, 4)
	if err != nil {
		t.Fatalf("error updating database from %d to %d: %v", version, 4, err)
	}
	version = db.checkVersion()
	if version != 4 {
		t.Fatalf("Version set to '%v' expected '4'.", version)
	}
	// Verify version 5
	err = db.updateTables(version, 5)
	if err != nil {
		t.Fatalf("error updating database from %d to %d: %v", version, 5, err)
	}
	version = db.checkVersion()
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
	// Check for error on drop tables as well. Because we can.
	err = db.dropTables()
	if err != nil {
		t.Fatalf("error deleting database: %v", err)
	}
}

func TestBadDatabase(t *testing.T) {
	db := &MySQL{}
	_, err := db.GetDatabase(&util.Config{})
	if err == nil {
		t.Fatal("Expected error getting database.")
	}
	db = &MySQL{}
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
	db := &MySQL{}
	_, err := db.GetDatabase(&util.Config{})
	if err == nil {
		t.Fatal("Expected error getting database.")
	}
	db = &MySQL{}
	_, err = db.GetDB()
	if err == nil {
		t.Fatal("Expected error getting database.")
	}
	db = &MySQL{}
	err = db.Setup(&util.Config{})
	if err == nil {
		t.Fatal("Expected error in Setup.")
	}
	db = &MySQL{}
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
