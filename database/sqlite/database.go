package sqlite

import (
	"chronokeep/results/auth"
	"chronokeep/results/database"
	"chronokeep/results/types"
	"chronokeep/results/util"

	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type SQLite struct {
	db       *sql.DB
	config   *util.Config
	validate *validator.Validate
}

// GetDatabase Used to get a database with given configuration information.
func (s *SQLite) GetDatabase(inCfg *util.Config) (*sql.DB, error) {
	if s.db != nil {
		return s.db, nil
	}
	if inCfg == nil {
		return nil, fmt.Errorf("no valid config supplied")
	}

	s.config = inCfg

	dbCon, err := sql.Open("sqlite3", inCfg.DBName+"?parseTime=true")
	if err != nil {
		return nil, fmt.Errorf("unable to open database connection: %v", err)
	}

	s.db = dbCon
	return s.db, nil
}

// GetDB Used as a general way to get a database.
func (s *SQLite) GetDB() (*sql.DB, error) {
	if s.db != nil {
		return s.db, nil
	}
	if s.config != nil {
		return s.GetDatabase(s.config)
	}
	return nil, errors.New("config file not established")
}

// Setup Automatically creates and updates tables for all of our information.
func (s *SQLite) Setup(config *util.Config) error {
	if config == nil {
		return fmt.Errorf("no valid config supplied")
	}
	// Set up Validator.
	s.validate = validator.New()
	log.Info("Setting up database.")
	// Connect to DB with database name.
	_, err := s.GetDatabase(config)
	if err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	// Check our database version.
	dbVersion := s.checkVersion()

	// Error checking version, most likely means tables are not created.
	if dbVersion < 1 {
		err = s.createTables()
		if err != nil {
			return err
		}
		// Otherwise check if our database is out of date and update if necessary.
	} else if dbVersion < database.CurrentVersion {
		log.Info(fmt.Sprintf("Updating database from version %v to %v", dbVersion, database.CurrentVersion))
		err = s.updateTables(dbVersion, database.CurrentVersion)
		if err != nil {
			return err
		}
	}

	// Check if there's an account created.
	accounts, err := s.GetAccounts()
	if err != nil {
		return fmt.Errorf("error checking for account: %v", err)
	}
	if len(accounts) < 1 {
		log.Info("Creating admin user.")
		if config.AdminName == "" || config.AdminEmail == "" || config.AdminPass == "" {
			return errors.New("admin account doesn't exist and proper credentions have not been supplied")
		}
		acc := types.Account{
			Name:     config.AdminName,
			Email:    config.AdminEmail,
			Password: config.AdminPass,
			Type:     "admin",
		}
		err = s.validate.Struct(acc)
		if err != nil {
			return fmt.Errorf("error validating base admin account on setup: %v", err)
		}
		acc.Password, err = auth.HashPassword(config.AdminPass)
		if err != nil {
			return fmt.Errorf("error hashing admin account password on setup: %v", err)
		}
		_, err = s.AddAccount(acc)
		if err != nil {
			return fmt.Errorf("error adding admin account on setup: %v", err)
		}
	}
	return nil
}

func (s *SQLite) dropTables() error {
	db, err := s.GetDatabase(s.config)
	if err != nil {
		return fmt.Errorf("error connecting to database to drop tables: %v", err)
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	_, err = db.ExecContext(
		ctx,
		"DROP TABLE distances;"+
			"DROP TABLE sms_subscriptions;"+
			"DROP TABLE linked_accounts;"+
			"DROP TABLE segments;"+
			"DROP TABLE participant;"+
			"DROP TABLE chips;"+
			"DROP TABLE banned_phones;"+
			"DROP TABLE banned_emails;"+
			"DROP TABLE call_record; "+
			"DROP TABLE result;"+
			"DROP TABLE person;"+
			"DROP TABLE event_year; "+
			"DROP TABLE event;"+
			"DROP TABLE api_key;"+
			"DROP TABLE account; "+
			"DROP TABLE settings;",
	)
	if err != nil {
		return fmt.Errorf("error dropping tables: %v", err)
	}
	return nil
}

func (s *SQLite) SetSetting(name, value string) error {
	db, err := s.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	_, err = db.ExecContext(
		ctx,
		"INSERT INTO settings(name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value=$2;",
		name,
		value,
	)
	if err != nil {
		return fmt.Errorf("error setting settings value: %v", err)
	}
	return nil
}

type myQuery struct {
	name  string
	query string
}

func (s *SQLite) createTables() error {
	log.Info("Creating database tables.")
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
				"cert_name VARCHAR(100) NOT NULL, " +
				"slug VARCHAR(50) NOT NULL, " +
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
				"days_allowed INT NOT NULL DEFAULT 1, " +
				"ranking_type VARCHAR(20) DEFAULT 'gun', " +
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
				"alternate_id VARCHAR(100) NOT NULL, " +
				"event_year_id BIGINT NOT NULL, " +
				"bib VARCHAR(100) NOT NULL, " +
				"first VARCHAR(100) NOT NULL, " +
				"last VARCHAR(100) NOT NULL, " +
				"age INT NOT NULL, " +
				"gender VARCHAR(50) NOT NULL, " +
				"age_group VARCHAR(200), " +
				"distance VARCHAR(200) NOT NULL, " +
				"anonymous SMALLINT NOT NULL DEFAULT 0, " +
				"division VARCHAR(500) NOT NULL DEFAULT '', " +
				"CONSTRAINT one_person UNIQUE (event_year_id, alternate_id), " +
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
				"local_time VARCHAR(100) NOT NULL DEFAULT '', " +
				"division_ranking INT NOT NULL DEFAULT -1, " +
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
		// BANNED PHONES TABLE
		{
			name: "CreateBannedPhones",
			query: "CREATE TABLE IF NOT EXISTS banned_phones(" +
				"banned_phone VARCHAR(20), " +
				"CONSTRAINT unique_banned_phone UNIQUE (banned_phone)" +
				");",
		},
		// BANNED EMAILS TABLE
		{
			name: "CreateBannedEmails",
			query: "CREATE TABLE IF NOT EXISTS banned_emails(" +
				"banned_email VARCHAR(200), " +
				"CONSTRAINT unique_banned_email UNIQUE (banned_email)" +
				");",
		},
		// PARTICIPANTS TABLE
		{
			name: "CreateParticipantTable",
			query: "CREATE TABLE IF NOT EXISTS participant(" +
				"participant_id INTEGER PRIMARY KEY AUTOINCREMENT, " +
				"alternate_id VARCHAR(100) NOT NULL, " +
				"event_year_id BIGINT NOT NULL, " +
				"bib VARCHAR(100) NOT NULL, " +
				"first VARCHAR(100) NOT NULL, " +
				"last VARCHAR(100) NOT NULL, " +
				"birthdate VARCHAR(15) NOT NULL, " +
				"gender VARCHAR(50) NOT NULL, " +
				"age_group VARCHAR(200), " +
				"distance VARCHAR(200) NOT NULL, " +
				"anonymous SMALLINT NOT NULL DEFAULT 0, " +
				"sms_enabled SMALLINT NOT NULL DEFAULT 0, " +
				"apparel VARCHAR(150) NOT NULL DEFAULT '', " +
				"mobile VARCHAR(15) NOT NULL DEFAULT '', " +
				"updated_at BIGINT NOT NULL DEFAULT 0, " +
				"CONSTRAINT one_participant UNIQUE (event_year_id, alternate_id), " +
				"FOREIGN KEY (event_year_id) REFERENCES event_year(event_year_id)" +
				");",
		},
		// BIBCHIP TABLE
		{
			name: "CreateChipTable",
			query: "CREATE TABLE IF NOT EXISTS chips(" +
				"chip_id INTEGER PRIMARY KEY AUTOINCREMENT, " +
				"event_year_id BIGINT NOT NULL, " +
				"bib VARCHAR(100) NOT NULL, " +
				"chip VARCHAR(100) NOT NULL, " +
				"CONSTRAINT unique_combo UNIQUE (event_year_id, chip), " +
				"FOREIGN KEY (event_year_id) REFERENCES event_year(event_year_id)" +
				");",
		},
		// SEGMENTS TABLE
		{
			name: "CreateSegmentTable",
			query: "CREATE TABLE IF NOT EXISTS segments(" +
				"segment_id INTEGER PRIMARY KEY AUTOINCREMENT, " +
				"event_year_id BIGINT NOT NULL, " +
				"location_name VARCHAR NOT NULL, " +
				"distance_name VARCHAR NOT NULL, " +
				"segment_name VARCHAR NOT NULL, " +
				"segment_distance DECIMAL(10,2) NOT NULL, " +
				"segment_distance_unit VARCHAR(12) NOT NULL, " +
				"segment_gps VARCHAR NOT NULL DEFAULT '', " +
				"segment_map_link VARCHAR NOT NULL DEFAULT '', " +
				"CONSTRAINT unique_segment UNIQUE (event_year_id, distance_name, segment_name), " +
				"FOREIGN KEY (event_year_id) REFERENCES event_year(event_year_id)" +
				");",
		},
		// DISTANCES TABLE
		{
			name: "CreateDistancesTable",
			query: "CREATE TABLE IF NOT EXISTS distances(" +
				"distance_id INTEGER PRIMARY KEY AUTOINCREMENT, " +
				"event_year_id BIGINT NOT NULL, " +
				"distance_name VARCHAR NOT NULL, " +
				"certification VARCHAR NOT NULL, " +
				"CONSTRAINT unique_distance UNIQUE (event_year_id, distance_name), " +
				"FOREIGN KEY (event_year_id) REFERENCES event_year(event_year_id)" +
				");",
		},
		// LINKED ACCOUNTS TABLE
		{
			name: "CreateLinkedAccountsTable",
			query: "CREATE TABLE IF NOT EXISTS linked_accounts(" +
				"main_account_id BIGINT NOT NULL, " +
				"sub_account_id BIGINT NOT NULL, " +
				"CONSTRAINT unique_link UNIQUE (main_account_id, sub_account_id), " +
				"FOREIGN KEY (main_account_id) REFERENCES account(account_id)," +
				"FOREIGN KEY (sub_account_id) REFERENCES account(account_id)" +
				");",
		},
		// SMS SUBSCRIPTIONS TABLE
		{
			name: "CreateSMSSubscriptionsTable",
			query: "CREATE TABLE IF NOT EXISTS sms_subscriptions(" +
				"event_year_id BIGINT NOT NULL, " +
				"bib VARCHAR(100) NOT NULL, " +
				"first VARCHAR(100) NOT NULL, " +
				"last VARCHAR(100) NOT NULL, " +
				"phone VARCHAR(15) NOT NULL, " +
				"CONSTRAINT one_subscription UNIQUE (event_year_id, bib, first, last, phone), " +
				"FOREIGN KEY (event_year_id) REFERENCES event_year(event_year_id)" +
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

	if s.db == nil {
		return fmt.Errorf("database not setup")
	}

	// Get a context and cancel function to create our tables, defer the cancel until we're done.
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("unable to start transaction: %v", err)
	}
	for _, single := range queries {
		log.Info(fmt.Sprintf("Executing query for: %s", single.name))
		_, err := tx.ExecContext(ctx, single.query)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error executing %s query: %v", single.name, err)
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("unable to commit transaction: %v", err)
	}

	s.SetSetting("version", strconv.Itoa(database.CurrentVersion))

	return nil
}

func (s *SQLite) checkVersion() int {
	log.Info("Checking database version.")
	if s.db == nil {
		return -1
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	var name string
	var version string
	err := s.db.QueryRowContext(
		ctx,
		"SELECT name, value FROM settings WHERE name='version';",
	).Scan(&name, &version)
	if err != nil {
		return -1
	}
	v, err := strconv.Atoi(version)
	if err != nil {
		return -1
	}
	return v
}

func (s *SQLite) updateTables(oldVersion, newVersion int) error {
	if s.db == nil {
		return fmt.Errorf("database not set up")
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("unable to start transaction: %v", err)
	}
	// SQLite starts at version 5.  6 will be the first update version.
	if oldVersion < 6 && newVersion >= 6 {
		log.Info("Updating to database version 6.")
		_, err := tx.ExecContext(
			ctx,
			"CREATE TABLE IF NOT EXISTS person_new("+
				"person_id INTEGER PRIMARY KEY AUTOINCREMENT, "+
				"event_year_id BIGINT NOT NULL, "+
				"bib VARCHAR(100) NOT NULL, "+
				"first VARCHAR(100) NOT NULL, "+
				"last VARCHAR(100) NOT NULL, "+
				"age INT NOT NULL, "+
				"gender VARCHAR(5) NOT NULL, "+
				"age_group VARCHAR(200), "+
				"distance VARCHAR(200) NOT NULL, "+
				"CONSTRAINT one_person UNIQUE (event_year_id, bib), "+
				"FOREIGN KEY (event_year_id) REFERENCES event_year(event_year_id)"+
				");",
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error updating from verison %d to %d: %v", oldVersion, newVersion, err)
		}
		_, err = tx.ExecContext(
			ctx,
			"INSERT INTO person_new SELECT * FROM person;",
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error updating from verison %d to %d: %v", oldVersion, newVersion, err)
		}
		_, err = tx.ExecContext(
			ctx,
			"DROP TABLE person; ALTER TABLE person_new RENAME TO person;",
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error updating from verison %d to %d: %v", oldVersion, newVersion, err)
		}
	}
	if oldVersion < 7 && newVersion >= 7 {
		log.Info("Updating to database version 7.")
		_, err := tx.ExecContext(
			ctx,
			"ALTER TABLE person "+
				"ADD COLUMN chip VARCHAR(200) DEFAULT ''; "+
				"ALTER TABLE person "+
				"ADD COLUMN anonymous SMALLINT NOT NULL DEFAULT 0;",
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error updating from version %d to %d: %v", oldVersion, newVersion, err)
		}
	}
	if oldVersion < 8 && newVersion >= 8 {
		log.Info("Updating to database version 8.")
		queries := []myQuery{
			{
				name:  "RenamePerson",
				query: "ALTER TABLE person RENAME TO person_old;",
			},
			{
				name: "CreateNewPerson",
				query: "CREATE TABLE IF NOT EXISTS person(" +
					"person_id INTEGER PRIMARY KEY AUTOINCREMENT, " +
					"event_year_id BIGINT NOT NULL, " +
					"bib VARCHAR(100) NOT NULL, " +
					"first VARCHAR(100) NOT NULL, " +
					"last VARCHAR(100) NOT NULL, " +
					"age INT NOT NULL, " +
					"gender VARCHAR(50) NOT NULL, " +
					"age_group VARCHAR(200), " +
					"distance VARCHAR(200) NOT NULL, " +
					"chip VARCHAR(200) DEFAULT '', " +
					"anonymous SMALLINT NOT NULL DEFAULT 0, " +
					"CONSTRAINT one_person UNIQUE (event_year_id, bib), " +
					"FOREIGN KEY (event_year_id) REFERENCES event_year(event_year_id)" +
					");",
			},
			{
				name: "InsertPerson",
				query: "INSERT INTO person(person_id, event_year_id, bib, first, last, age, " +
					"gender, age_group, distance, chip, anonymous) " +
					"SELECT * FROM person_old;",
			},
			{
				name:  "DeleteOldPerson",
				query: "DROP TABLE person_old;",
			},
		}
		for _, q := range queries {
			_, err := tx.ExecContext(
				ctx,
				q.query,
			)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("error updating from version %d to %d in query %s: %v", oldVersion, newVersion, q.name, err)
			}
		}
	}
	if oldVersion < 9 && newVersion >= 9 {
		log.Info("Updating to database version 9.")
		queries := []myQuery{
			{
				name:  "RenameEvent",
				query: "ALTER TABLE event RENAME TO event_old;",
			},
			{
				name: "CreateNewEvent",
				query: "CREATE TABLE IF NOT EXISTS event(" +
					"event_id INTEGER PRIMARY KEY AUTOINCREMENT, " +
					"account_id BIGINT NOT NULL, " +
					"event_name VARCHAR(100) NOT NULL, " +
					"slug VARCHAR(50) NOT NULL, " +
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
			{
				name: "InsertEvent",
				query: "INSERT INTO event(event_id, account_id, event_name, slug, website, image, " +
					"contact_email, access_restricted, event_created_at, event_updated_at, event_deleted, event_type) " +
					"SELECT * FROM event_old;",
			},
			{
				name:  "DeleteOldEvent",
				query: "DROP TABLE event_old;",
			},
		}
		for _, q := range queries {
			_, err := tx.ExecContext(
				ctx,
				q.query,
			)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("error updating from version %d to %d in query %s: %v", oldVersion, newVersion, q.name, err)
			}
		}
	}
	if oldVersion < 10 && newVersion >= 10 {
		log.Info("Updating to database version 10.")
		queries := []myQuery{
			{
				name:  "Update DNF entries.",
				query: "UPDATE result SET seconds=1000000 WHERE result_type=30 OR result_type=3;",
			},
		}
		for _, q := range queries {
			_, err := tx.ExecContext(
				ctx,
				q.query,
			)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("error updating from version %d to %d in query %s: %v", oldVersion, newVersion, q.name, err)
			}
		}
	}
	if oldVersion < 11 && newVersion >= 11 {
		log.Info("Updating to database version 11.")
		queries := []myQuery{
			{
				name: "CreateBannedPhones",
				query: "CREATE TABLE IF NOT EXISTS banned_phones(" +
					"banned_phone VARCHAR(200), " +
					"UNIQUE(banned_phone)" +
					");",
			},
			{
				name: "CreateBannedEmails",
				query: "CREATE TABLE IF NOT EXISTS banned_emails(" +
					"banned_email VARCHAR(200), " +
					"UNIQUE(banned_email)" +
					");",
			},
		}
		for _, q := range queries {
			_, err := tx.ExecContext(
				ctx,
				q.query,
			)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("error updating from version %d to %d in query %s: %v", oldVersion, newVersion, q.name, err)
			}
		}
	}
	if oldVersion < 12 && newVersion >= 12 {
		log.Info("Updating to database version 12.")
		queries := []myQuery{
			{
				name: "CreateNewPerson",
				query: "CREATE TABLE IF NOT EXISTS new_person(" +
					"person_id BIGSERIAL NOT NULL, " +
					"alternate_id VARCHAR(100) NOT NULL, " +
					"event_year_id BIGINT NOT NULL, " +
					"bib VARCHAR(100) NOT NULL, " +
					"first VARCHAR(100) NOT NULL, " +
					"last VARCHAR(100) NOT NULL, " +
					"age INT NOT NULL, " +
					"gender VARCHAR(50) NOT NULL, " +
					"age_group VARCHAR(200), " +
					"distance VARCHAR(200) NOT NULL, " +
					"chip VARCHAR(200) DEFAULT '', " +
					"anonymous SMALLINT NOT NULL DEFAULT 0, " +
					"sms_enabled SMALLINT NOT NULL DEFAULT 0, " +
					"CONSTRAINT one_person UNIQUE (event_year_id, alternate_id), " +
					"FOREIGN KEY (event_year_id) REFERENCES event_year(event_year_id), " +
					"PRIMARY KEY (person_id)" +
					");",
			},
			{
				name: "InsertPersonData",
				query: "INSERT INTO new_person(" +
					"person_id, " +
					"alternate_id, " +
					"event_year_id, " +
					"bib, " +
					"first, " +
					"last, " +
					"age, " +
					"gender, " +
					"age_group, " +
					"distance, " +
					"chip, " +
					"anonymous" +
					") SELECT person_id, bib, event_year_id, bib, first, last, age, gender, age_group, distance, chip, anonymous FROM person;",
			},
			{
				name:  "RenameNewPersonDropOld",
				query: "DROP TABLE person; ALTER TABLE new_person RENAME TO person;",
			},
		}
		for _, q := range queries {
			_, err := tx.ExecContext(
				ctx,
				q.query,
			)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("error updating from version %d to %d in query %s: %v", oldVersion, newVersion, q.name, err)
			}
		}
	}
	if oldVersion < 13 && newVersion >= 13 {
		log.Info("Updating to database version 13.")
		queries := []myQuery{
			{
				name:  "DropPersonColumns",
				query: "ALTER TABLE person DROP COLUMN chip; ALTER TABLE person DROP COLUMN sms_enabled;",
			},
			{
				name: "CreateParticipantTable",
				query: "CREATE TABLE IF NOT EXISTS participant(" +
					"participant_id INTEGER PRIMARY KEY AUTOINCREMENT, " +
					"alternate_id VARCHAR(100) NOT NULL, " +
					"event_year_id BIGINT NOT NULL, " +
					"bib VARCHAR(100) NOT NULL, " +
					"first VARCHAR(100) NOT NULL, " +
					"last VARCHAR(100) NOT NULL, " +
					"birthdate VARCHAR(15) NOT NULL, " +
					"gender VARCHAR(50) NOT NULL, " +
					"age_group VARCHAR(200), " +
					"distance VARCHAR(200) NOT NULL, " +
					"anonymous SMALLINT NOT NULL DEFAULT 0, " +
					"sms_enabled SMALLINT NOT NULL DEFAULT 0, " +
					"apparel VARCHAR(150) NOT NULL DEFAULT '', " +
					"mobile VARCHAR(15) NOT NULL DEFAULT '', " +
					"CONSTRAINT one_participant UNIQUE (event_year_id, alternate_id), " +
					"FOREIGN KEY (event_year_id) REFERENCES event_year(event_year_id)" +
					");",
			},
			{
				name: "CreateChipTable",
				query: "CREATE TABLE IF NOT EXISTS chips(" +
					"chip_id INTEGER PRIMARY KEY AUTOINCREMENT, " +
					"event_year_id BIGINT NOT NULL, " +
					"bib VARCHAR(100) NOT NULL, " +
					"chip VARCHAR(100) NOT NULL, " +
					"CONSTRAINT unique_combo UNIQUE (event_year_id, chip), " +
					"FOREIGN KEY (event_year_id) REFERENCES event_year(event_year_id)" +
					");",
			},
			{
				name: "CreateSegmentTable",
				query: "CREATE TABLE IF NOT EXISTS segments(" +
					"segment_id INTEGER PRIMARY KEY AUTOINCREMENT, " +
					"event_year_id BIGINT NOT NULL, " +
					"location_name VARCHAR NOT NULL, " +
					"distance_name VARCHAR NOT NULL, " +
					"segment_name VARCHAR NOT NULL, " +
					"segment_distance DECIMAL(10,2) NOT NULL, " +
					"segment_distance_unit VARCHAR(12) NOT NULL, " +
					"segment_gps VARCHAR NOT NULL DEFAULT '', " +
					"segment_map_link VARCHAR NOT NULL DEFAULT '', " +
					"CONSTRAINT unique_segment UNIQUE (event_year_id, distance_name, segment_name), " +
					"FOREIGN KEY (event_year_id) REFERENCES event_year(event_year_id)" +
					");",
			},
		}
		for _, q := range queries {
			_, err := tx.ExecContext(
				ctx,
				q.query,
			)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("error updating from version %d to %d in query %s: %v", oldVersion, newVersion, q.name, err)
			}
		}
	}
	if oldVersion < 14 && newVersion >= 14 {
		log.Info("Updating to database version 14.")
		queries := []myQuery{
			{
				name: "CreateLinkedAccountsTable",
				query: "CREATE TABLE IF NOT EXISTS linked_accounts(" +
					"main_account_id BIGINT NOT NULL, " +
					"sub_account_id BIGINT NOT NULL, " +
					"CONSTRAINT unique_link UNIQUE (main_account_id, sub_account_id), " +
					"FOREIGN KEY (main_account_id) REFERENCES account(account_id)," +
					"FOREIGN KEY (sub_account_id) REFERENCES account(account_id)" +
					");",
			},
		}
		for _, q := range queries {
			_, err := tx.ExecContext(
				ctx,
				q.query,
			)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("error updating from version %d to %d in query %s: %v", oldVersion, newVersion, q.name, err)
			}
		}
	}
	if oldVersion < 15 && newVersion >= 15 {
		log.Info("Updating to database version 15.")
		queries := []myQuery{
			// SMS SUBSCRIPTIONS TABLE
			{
				name: "CreateSMSSubscriptionsTable",
				query: "CREATE TABLE IF NOT EXISTS sms_subscriptions(" +
					"event_year_id BIGINT NOT NULL, " +
					"bib VARCHAR(100) NOT NULL, " +
					"first VARCHAR(100) NOT NULL, " +
					"last VARCHAR(100) NOT NULL, " +
					"phone VARCHAR(15) NOT NULL, " +
					"CONSTRAINT one_subscription UNIQUE (event_year_id, bib, first, last, phone), " +
					"FOREIGN KEY (event_year_id) REFERENCES event_year(event_year_id)" +
					");",
			},
			{
				name:  "AlterEventYearTable",
				query: "ALTER TABLE event_year ADD COLUMN days_allowed INT NOT NULL DEFAULT 1;",
			},
			{
				name:  "AlterResultTable",
				query: "ALTER TABLE result ADD COLUMN local_time VARCHAR(100) NOT NULL DEFAULT '';",
			},
		}
		for _, q := range queries {
			_, err := tx.ExecContext(
				ctx,
				q.query,
			)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("error updating from version %d to %d in query %s: %v", oldVersion, newVersion, q.name, err)
			}
		}
	}
	if oldVersion < 16 && newVersion >= 16 {
		log.Info("Updating to database version 16.")
		queries := []myQuery{
			{
				name:  "AlterEventYearTable",
				query: "ALTER TABLE event_year ADD COLUMN ranking_type VARCHAR(20) DEFAULT 'gun';",
			},
		}
		for _, q := range queries {
			_, err := tx.ExecContext(
				ctx,
				q.query,
			)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("error updating from version %d to %d in query %s: %v", oldVersion, newVersion, q.name, err)
			}
		}
	}
	if oldVersion < 17 && newVersion >= 17 {
		log.Info("Updating to database version 17.")
		queries := []myQuery{
			{
				name:  "AlterEventTable",
				query: "ALTER TABLE event ADD COLUMN cert_name VARCHAR(100) NOT NULL DEFAULT '';",
			},
			{
				name: "CreateDistancesTable",
				query: "CREATE TABLE IF NOT EXISTS distances(" +
					"distance_id INTEGER PRIMARY KEY AUTOINCREMENT, " +
					"event_year_id BIGINT NOT NULL, " +
					"distance_name VARCHAR NOT NULL, " +
					"certification VARCHAR NOT NULL, " +
					"CONSTRAINT unique_distance UNIQUE (event_year_id, distance_name), " +
					"FOREIGN KEY (event_year_id) REFERENCES event_year(event_year_id)" +
					");",
			},
		}
		for _, q := range queries {
			_, err := tx.ExecContext(
				ctx,
				q.query,
			)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("error updating from version %d to %d in query %s: %v", oldVersion, newVersion, q.name, err)
			}
		}
	}
	if oldVersion < 18 && newVersion >= 18 {
		log.Info("Updating to database version 18.")
		queries := []myQuery{
			{
				name:  "AlterResultTable",
				query: "ALTER TABLE result ADD COLUMN division_ranking INT NOT NULL DEFAULT -1;",
			},
			{
				name:  "AlterPersonTable",
				query: "ALTER TABLE person ADD COLUMN division VARCHAR(500) NOT NULL DEFAULT '';",
			},
		}
		for _, q := range queries {
			_, err := tx.ExecContext(
				ctx,
				q.query,
			)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("error updating from version %d to %d in query %s: %v", oldVersion, newVersion, q.name, err)
			}
		}
	}
	if oldVersion < 19 && newVersion >= 19 {
		log.Info("Updating to database version 19.")
		queries := []myQuery{
			{
				name:  "AlterParticipantTable",
				query: "ALTER TABLE participant ADD COLUMN updated_at BIGINT NOT NULL DEFAULT 0;",
			},
		}
		for _, q := range queries {
			_, err := tx.ExecContext(
				ctx,
				q.query,
			)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("error updating from version %d to %d in query %s: %v", oldVersion, newVersion, q.name, err)
			}
		}
	}
	_, err = tx.ExecContext(
		ctx,
		"UPDATE settings SET value=$1 WHERE name='version';",
		strconv.Itoa(newVersion),
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error updating from version %d to %d: %v", oldVersion, newVersion, err)
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error committing transaction: %v", err)
	}
	return nil
}

func (s *SQLite) updateDB(newdb *sql.DB) {
	s.db = newdb
}

// Close Closes database.
func (s *SQLite) Close() {
	s.db.Close()
}
