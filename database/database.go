package database

import (
	"chronokeep/results/util"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db     *sql.DB
	config *util.Config
)

const (
	MaxOpenConnections    = 20
	MaxIdleConnections    = 20
	MaxConnectionLifetime = time.Minute * 5
	CurrentVersion        = 1
)

// GetDatabase Used to get a database with given configuration information.
func GetDatabase(inCfg *util.Config) (*sql.DB, error) {
	if db != nil {
		return db, nil
	}

	config = inCfg
	conString := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBName,
	)

	switch config.DBDriver {
	case "postgres":
	case "mysql":
	default:
		return nil, errors.New("invalid database type given")
	}

	dbCon, err := sql.Open(config.DBDriver, conString)
	if err != nil {
		return nil, fmt.Errorf("unable to open database connection: %v", err)
	}
	dbCon.SetMaxIdleConns(MaxIdleConnections)
	dbCon.SetMaxOpenConns(MaxOpenConnections)
	dbCon.SetConnMaxLifetime(MaxConnectionLifetime)

	db = dbCon
	return db, nil
}

// GetDB Used as a general way to get a database.
func GetDB() (*sql.DB, error) {
	if db != nil {
		return db, nil
	}
	if config != nil {
		return GetDatabase(config)
	}
	return nil, errors.New("config file not established")
}

// Setup Automatically creates and updates tables for all of our information.
func Setup(config *util.Config) error {
	// Connect to DB with database name.
	_, err := GetDatabase(config)
	if err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	// Check our database version.
	dbVersion := checkVersion()

	// Error checking version, most likely means tables are not created.
	if dbVersion < 1 {
		err = createTables()
		if err != nil {
			return err
		}
		// Otherwise check if our database is out of date and update if necessary.
	} else if dbVersion < CurrentVersion {
		err = updateTables(dbVersion, CurrentVersion)
		if err != nil {
			return err
		}
	}
	return nil
}

func dropTables() error {
	db, err := GetDatabase(config)
	if err != nil {
		return fmt.Errorf("error connecting to database to drop tables: %v", err)
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	_, err = db.ExecContext(
		ctx,
		"DROP TABLE call_record, result, event_year, event, api_key, account, settings;",
	)
	if err != nil {
		return fmt.Errorf("error dropping tables: %v", err)
	}
	return nil
}

func SetSetting(name, value string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	_, err = db.ExecContext(
		ctx,
		"INSERT INTO settings(name, value) VALUES (?, ?) ON DUPLICATE KEY UPDATE value=VALUES(value);",
		name,
		value,
	)
	if err != nil {
		return fmt.Errorf("error setting settings value: %v", err)
	}
	return nil
}

func createTables() error {
	var settingsTable, accountTable, keyTable, eventTable, eventYearTable, resultTable, recordTable string
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
			"token VARCHAR(100), " +
			"refresh_token VARCHAR(100), " +
			"type VARCHAR(20) NOT NULL, " +
			"created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
			"updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
			"account_deleted BOOL DEFAULT FALSE, " +
			"UNIQUE(account_email), " +
			"PRIMARY KEY (account_id)" +
			");"

		keyTable = "CREATE TABLE IF NOT EXISTS api_key(" +
			"account_id BIGINT NOT NULL, " +
			"value CHAR(100) NOT NULL, " +
			"type VARCHAR(20) NOT NULL, " +
			"allowed_hosts TEXT, " +
			"valid_until DATETIME DEFAULT CURRENT_TIMESTAMP, " +
			"created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
			"updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, " +
			"api_key_deleted BOOL DEFAULT FALSE, " +
			"UNIQUE(value), " +
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
			"created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
			"updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
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
			"created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
			"updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
			"event_year_deleted BOOL DEFAULT FALSE, " +
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
			"created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
			"updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
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

	// Get a context and cancel function to create our tables, defer the cancel until we're done.
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()

	_, err := db.ExecContext(ctx, settingsTable)
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

	SetSetting("version", strconv.Itoa(CurrentVersion))

	return nil
}

func checkVersion() int {
	res, err := db.Query("SELECT * FROM settings WHERE name='version';")
	if err != nil {
		return -1
	}
	if res.Next() {
		var name string
		var version int
		err = res.Scan(&name, &version)
		if err != nil {
			return -1
		}
		return version
	}
	return -1
}

func updateTables(oldVersion, newVersion int) error {
	switch config.DBDriver {
	case "postgres":
		return errors.New("postgres not yet supported")
	case "mysql":
	default:
		return errors.New("invalid database type given")
	}
	return nil
}

func updateDB(newdb *sql.DB) {
	db = newdb
}

// Close Closes database.
func Close() {
	db.Close()
}
