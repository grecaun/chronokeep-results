package database

import (
	"chronokeep/results/util"
	"context"
	"errors"
	"fmt"
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
		"%s:%s@tcp(%s)/%s",
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

func createDatabase() error {
	cfg := util.Config{
		DBDriver:   config.DBDriver,
		DBHost:     config.DBHost,
		DBPassword: config.DBHost,
		DBUser:     config.DBUser,
		DBName:     "",
	}
	db, err := GetDatabase(&cfg)
	if err != nil {
		return fmt.Errorf("error connecting to database to create database: %v", err)
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", config.DBName))
	if err != nil {
		return fmt.Errorf("error creating database: %v", err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error fetching rows on create database: %v", err)
	}
	updateDB(nil)
	return db.Close()
}

func dropTables() error {
	db, err := GetDatabase(config)
	if err != nil {
		return fmt.Errorf("error connecting to database to drop tables: %v", err)
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"DROP TABLE call_record; DROP TABLE result; DROP TABLE event_year; DROP TABLE event; DROP TABLE key; DROP TABLE account; DROP TABLE settings;",
	)
	if err != nil {
		return fmt.Errorf("error dropping tables: %v", err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error fetching rows on database drop tables: %v", err)
	}
	return nil
}

func deleteDatabase() error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, fmt.Sprintf("DELETE DATABASE %s", config.DBName))
	if err != nil {
		return fmt.Errorf("error deleting database: %v", err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error fetching rows on database delete: %v", err)
	}
	updateDB(nil)
	return db.Close()
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
			"UNIQUE (name) ON CONFLICT UPDATE);"

		accountTable = "CREATE TABLE IF NOT EXISTS account(" +
			"account_id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT, " +
			"name VARCHAR(100) NOT NULL, " +
			"email VARCHAR(100) NOT NULL, " +
			"type VARCHAR(20) NOT NULL, " +
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
			"segment VARCHAR(500), " +
			"location VARCHAR(500), " +
			"occurence INT DEFAULT -1, " +
			"ranking INT DEFAULT -1, " +
			"age_ranking INT DEFUALT -1, " +
			"gender_ranking INT DEFAULT -1, " +
			"finish BOOL DEFAULT TRUE, " +
			"created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
			"updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
			"CONSTRAINT one_occurrence UNIQUE (event_year_id, bib, location, occurence) ON CONFLICT UPDATE" +
			");"

		recordTable = "CREATE TABLE IF NOT EXISTS call_record(" +
			"account_id BIGINT FOREIGN KEY REFERENCES account(account_id), " +
			"time BIGINT NOT NULL, " +
			"count INT DEFAULT 0, " +
			"CONSTRAINT account_time UNIQUE (account_id, time) ON CONFLICT UPDATE" +
			");"
	default:
		return errors.New("invalid database type given")
	}

	settingsValue := fmt.Sprintf("INSERT INTO settings(name, value) VALUES ('version', '%v');", CurrentVersion)

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

	_, err = db.ExecContext(ctx, settingsValue)
	if err != nil {
		return fmt.Errorf("error adding settings: %v", err)
	}

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
