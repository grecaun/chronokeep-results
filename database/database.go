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
	DatabaseVersion       = 1
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
	if config != nil {
		return GetDatabase(config)
	}
	return nil, errors.New("config file not established")
}

// Setup Automatically creates and updates tables for all of our information.
func Setup(inCfg *util.Config) error {
	// Connect to the database software without a database name.
	config = inCfg
	dbName := config.DBName
	config.DBName = ""

	// Create the database if it doesn't exists.
	err := createDatabase(dbName)
	if err != nil {
		return err
	}

	// Connect to DB with database name.
	config.DBName = dbName
	db, err = GetDatabase(config)
	if err != nil {
		return fmt.Errorf("error connecting to database after establishing database exists: %v", err)
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
	} else if dbVersion < DatabaseVersion {
		err = updateTables()
		if err != nil {
			return err
		}
	}
	return nil
}

func createDatabase(dbName string) error {
	db, err := GetDatabase(config)
	if err != nil {
		return fmt.Errorf("error connecting to database with no database name: %v", err)
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	result, err := db.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", dbName))
	if err != nil {
		return fmt.Errorf("error creating database: %v", err)
	}
	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error fetching rows: %v", err)
	}
	return db.Close()
}

func createTables() error {
	settingsTable := "CREATE TABLE IF NOT EXISTS settings(" +
		"name VARCHAR(200) NOT NULL, " +
		"value VARCHAR(200) NOT NULL, " +
		"UNIQUE (name));"

	accountTable := "CREATE TABLE IF NOT EXISTS account(" +
		"account_id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT, " +
		"name VARCHAR(100), " +
		"email VARCHAR(100), " +
		"type VARCHAR(20), " +
		"created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
		"updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
		"deleted BOOL DEFAULT FALSE, " +
		"UNIQUE(email));"

	keyTable := "CREATE TABLE IF NOT EXISTS key(" +
		"key_id BIGINT PRIMARY KEY AUTO_INCREMENT, " +
		"account_id BIGINT FOREIGN KEY REFERENCES account(account_id), " +
		"value CHAR(100) NOT NULL, " +
		"type VARCHAR(20) NOT NULL, " +
		"allowed_hosts TEXT, " +
		"valid_until DATETIME DEFAULT CURRENT_TIMESTAMP," +
		"created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
		"updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
		"deleted BOOL DEFAULT FALSE);"

	eventTable := "CREATE TABLE IF NOT EXISTS event(" +
		"event_id BIGINT PRIMARY KEY AUTO_INCREMENT, " +
		"name VARCHAR(100) NOT NULL, " +
		"slug VARCHAR(20) NOT NULL, " +
		"website VARCHAR(200), " +
		"image VARCHAR(200), " +
		"account_id BIGINT FOREIGN KEY REFERENCES account(account_id), " +
		"access_restricted BOOL DEFAULT FALSE, " +
		"created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
		"updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
		"deleted BOOL DEFAULT FALSE, " +
		"UNIQUE(name), " +
		"UNIQUE(slug)" +
		");"

	eventYearTable := "CREATE TABLE IF NOT EXISTS event_year(" +
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

	resultTable := "CREATE TABLE IF NOT EXISTS result(" +
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
		"CONSTRAINT one_finish UNIQUE (event_year_id, bib, finish), " +
		"CONSTRAINT one_occurrence UNIQUE (event_year_id, bib, location, occurence)" +
		");"

	recordTable := "CREATE TABLE IF NOT EXISTS call_record(" +
		"account_id BIGINT FOREIGN KEY REFERENCES account(account_id), " +
		"time DATETIME NOT NULL, " +
		"count INT DEFAULT 0, " +
		"CONSTRAINT account_time UNIQUE (account_id, time)" +
		");"

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

func updateTables() error {
	return nil
}

func Close() {
	db.Close()
}
