package mysql

import (
	"chronokeep/results/auth"
	"chronokeep/results/database"
	"chronokeep/results/types"
	"chronokeep/results/util"

	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"database/sql"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

type MySQL struct {
	db       *sql.DB
	config   *util.Config
	validate *validator.Validate
}

// GetDatabase Used to get a database with given configuration information.
func (m *MySQL) GetDatabase(inCfg *util.Config) (*sql.DB, error) {
	if m.db != nil {
		return m.db, nil
	}
	if inCfg == nil {
		return nil, fmt.Errorf("no valid config supplied")
	}

	m.config = inCfg
	conString := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		m.config.DBUser,
		m.config.DBPassword,
		m.config.DBHost,
		m.config.DBPort,
		m.config.DBName,
	)

	dbCon, err := sql.Open(m.config.DBDriver, conString)
	if err != nil {
		return nil, fmt.Errorf("unable to open database connection: %v", err)
	}
	dbCon.SetMaxIdleConns(database.MaxIdleConnections)
	dbCon.SetMaxOpenConns(database.MaxOpenConnections)
	dbCon.SetConnMaxLifetime(database.MaxConnectionLifetime)

	m.db = dbCon
	return m.db, nil
}

// GetDB Used as a general way to get a database.
func (m *MySQL) GetDB() (*sql.DB, error) {
	if m.db != nil {
		return m.db, nil
	}
	if m.config != nil {
		return m.GetDatabase(m.config)
	}
	return nil, errors.New("config file not established")
}

// Setup Automatically creates and updates tables for all of our information.
func (m *MySQL) Setup(config *util.Config) error {
	if config == nil {
		return fmt.Errorf("no valid config supplied")
	}
	// Set up Validator.
	m.validate = validator.New()
	log.Info("Setting up database.")
	// Connect to DB with database name.
	_, err := m.GetDatabase(config)
	if err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	// Check our database version.
	dbVersion := m.checkVersion()

	// Error checking version, most likely means tables are not created.
	if dbVersion < 1 {
		err = m.createTables()
		if err != nil {
			return err
		}
		// Otherwise check if our database is out of date and update if necessary.
	} else if dbVersion < database.CurrentVersion {
		log.Info(fmt.Sprintf("Updating database from version %v to %v", dbVersion, database.CurrentVersion))
		err = m.updateTables(dbVersion, database.CurrentVersion)
		if err != nil {
			return err
		}
	}

	// Check if there's an account created.
	accounts, err := m.GetAccounts()
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
		err = m.validate.Struct(acc)
		if err != nil {
			return fmt.Errorf("error validating base admin account on setup: %v", err)
		}
		acc.Password, err = auth.HashPassword(config.AdminPass)
		if err != nil {
			return fmt.Errorf("error hashing admin account password on setup: %v", err)
		}
		_, err = m.AddAccount(acc)
		if err != nil {
			return fmt.Errorf("error adding admin account on setup: %v", err)
		}
	}
	return nil
}

func (m *MySQL) dropTables() error {
	db, err := m.GetDatabase(m.config)
	if err != nil {
		return fmt.Errorf("error connecting to database to drop tables: %v", err)
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	_, err = db.ExecContext(
		ctx,
		"DROP TABLE "+
			"distances, "+
			"sms_subscriptions, "+
			"linked_accounts, "+
			"segments, "+
			"participant, "+
			"chips, "+
			"banned_phones, "+
			"banned_emails, "+
			"call_record, "+
			"result, "+
			"person, "+
			"event_year, "+
			"event, "+
			"api_key, "+
			"account, "+
			"settings;",
	)
	if err != nil {
		return fmt.Errorf("error dropping tables: %v", err)
	}
	return nil
}

func (m *MySQL) SetSetting(name, value string) error {
	db, err := m.GetDB()
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

type myQuery struct {
	name  string
	query string
}

func (m *MySQL) createTables() error {
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
				"key_name VARCHAR(100) NOT NULL DEFAULT ''," +
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
				"cert_name VARCHAR(100) NOT NULL, " +
				"slug VARCHAR(50) NOT NULL, " +
				"website VARCHAR(200), " +
				"image VARCHAR(200), " +
				"contact_email VARCHAR(100), " +
				"access_restricted BOOL DEFAULT FALSE, " +
				"event_created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
				"event_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
				"event_deleted BOOL DEFAULT FALSE, " +
				"event_type VARCHAR(20) DEFAULT 'distance', " +
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
				"days_allowed INT NOT NULL DEFAULT 1, " +
				"ranking_type VARCHAR(20) DEFAULT 'gun', " +
				"year_created_at DATETIME DEFAULT CURRENT_TIMESTAMP, " +
				"year_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
				"year_deleted BOOL DEFAULT FALSE, " +
				"CONSTRAINT year_slug UNIQUE (event_id, year)," +
				"FOREIGN KEY (event_id) REFERENCES event(event_id)," +
				"PRIMARY KEY (event_year_id)" +
				");",
		},
		// PERSON TABLE
		{
			name: "PersonTable",
			query: "CREATE TABLE IF NOT EXISTS person(" +
				"person_id BIGINT NOT NULL AUTO_INCREMENT, " +
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
				"FOREIGN KEY (event_year_id) REFERENCES event_year(event_year_id), " +
				"PRIMARY KEY (person_id)" +
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
				"result_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
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
				"CONSTRAINT unique_banned_phone UNIQUE(banned_phone)" +
				");",
		},
		// BANNED EMAILS TABLE
		{
			name: "CreateBannedEmails",
			query: "CREATE TABLE IF NOT EXISTS banned_emails(" +
				"banned_email VARCHAR(200), " +
				"CONSTRAINT unique_banned_email UNIQUE(banned_email)" +
				");",
		},
		// PARTICIPANTS TABLE
		{
			name: "CreateParticipantTable",
			query: "CREATE TABLE IF NOT EXISTS participant(" +
				"participant_id BIGINT NOT NULL AUTO_INCREMENT, " +
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
				"FOREIGN KEY (event_year_id) REFERENCES event_year(event_year_id), " +
				"PRIMARY KEY (participant_id)" +
				");",
		},
		// BIBCHIP TABLE
		{
			name: "CreateChipTable",
			query: "CREATE TABLE IF NOT EXISTS chips(" +
				"chip_id BIGINT NOT NULL AUTO_INCREMENT, " +
				"event_year_id BIGINT NOT NULL, " +
				"bib VARCHAR(100) NOT NULL, " +
				"chip VARCHAR(100) NOT NULL, " +
				"CONSTRAINT unique_combo UNIQUE (event_year_id, chip), " +
				"FOREIGN KEY (event_year_id) REFERENCES event_year(event_year_id), " +
				"PRIMARY KEY (chip_id)" +
				");",
		},
		// SEGMENTS TABLE
		{
			name: "CreateSegmentTable",
			query: "CREATE TABLE IF NOT EXISTS segments(" +
				"segment_id BIGINT NOT NULL AUTO_INCREMENT, " +
				"event_year_id BIGINT NOT NULL, " +
				"location_name VARCHAR(100) NOT NULL, " +
				"distance_name VARCHAR(100) NOT NULL, " +
				"segment_name VARCHAR(100) NOT NULL, " +
				"segment_distance DECIMAL(10,2) NOT NULL DEFAULT 0.0, " +
				"segment_distance_unit VARCHAR(12) NOT NULL DEFAULT '', " +
				"segment_gps VARCHAR(500) NOT NULL DEFAULT '', " +
				"segment_map_link VARCHAR(500) NOT NULL DEFAULT '', " +
				"CONSTRAINT unique_segment UNIQUE (event_year_id, distance_name, segment_name), " +
				"FOREIGN KEY (event_year_id) REFERENCES event_year(event_year_id), " +
				"PRIMARY KEY (segment_id)" +
				");",
		},
		// DISTANCES TABLE
		{
			name: "CreateDistancesTable",
			query: "CREATE TABLE IF NOT EXISTS distances(" +
				"distance_id BIGINT NOT NULL AUTO_INCREMENT, " +
				"event_year_id BIGINT NOT NULL, " +
				"distance_name VARCHAR(100) NOT NULL, " +
				"certification VARCHAR(150) NOT NULL, " +
				"CONSTRAINT unique_distance UNIQUE (event_year_id, distance_name), " +
				"FOREIGN KEY (event_year_id) REFERENCES event_year(event_year_id), " +
				"PRIMARY KEY (distance_id)" +
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
	}

	if m.db == nil {
		return fmt.Errorf("database not setup")
	}

	// Get a context and cancel function to create our tables, defer the cancel until we're done.
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()

	tx, err := m.db.Begin()
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
		return fmt.Errorf("error committing transaction: %v", err)
	}

	m.SetSetting("version", strconv.Itoa(database.CurrentVersion))

	return nil
}

func (m *MySQL) checkVersion() int {
	log.Info("Checking database version.")
	if m.db == nil {
		return -1
	}
	res, err := m.db.Query("SELECT * FROM settings WHERE name='version';")
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

func (m *MySQL) updateTables(oldVersion, newVersion int) error {
	if m.db == nil {
		return fmt.Errorf("database not set up")
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	if oldVersion < 2 && newVersion >= 2 {
		log.Debug("Updating to database version 2.")
		_, err := m.db.ExecContext(
			ctx,
			"ALTER TABLE result ADD COLUMN result_type INT DEFAULT 0;",
		)
		if err != nil {
			return fmt.Errorf("error updating from version %d to %d: %v", oldVersion, newVersion, err)
		}
	}
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("unable to start transaction: %v", err)
	}
	if oldVersion < 3 && newVersion >= 3 {
		log.Info("Updating to database version 3.")
		queries := []myQuery{
			{
				name:  "RenameResult",
				query: "ALTER TABLE result RENAME TO result_old;",
			},
			{
				name: "CreatePerson",
				query: "CREATE TABLE IF NOT EXISTS person(" +
					"person_id BIGINT NOT NULL AUTO_INCREMENT, " +
					"event_year_id BIGINT NOT NULL, " +
					"bib VARCHAR(100) NOT NULL, " +
					"first VARCHAR(100) NOT NULL, " +
					"last VARCHAR(100) NOT NULL, " +
					"age INT NOT NULL, " +
					"gender CHAR(1) NOT NULL, " +
					"age_group VARCHAR(200), " +
					"distance VARCHAR(200) NOT NULL, " +
					"CONSTRAINT one_person UNIQUE (event_year_id, bib), " +
					"FOREIGN KEY (event_year_id) REFERENCES event_year(event_year_id), " +
					"PRIMARY KEY (person_id)" +
					");",
			},
			{
				name: "CreateNewResult",
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
					"result_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
					"CONSTRAINT one_occurrence_res UNIQUE (person_id, location, occurence)," +
					"FOREIGN KEY (person_id) REFERENCES person(person_id)" +
					");",
			},
			{
				name: "InsertPerson",
				query: "INSERT INTO person (event_year_id, bib, first, last, age, gender, age_group, distance) " +
					" SELECT event_year_id, bib, first, last, age, gender, age_group, distance FROM result_old " +
					"ON DUPLICATE KEY UPDATE first=VALUES(first), last=VALUES(last), age=VALUES(age), gender=VALUES(gender), " +
					"age_group=VALUES(age_group), distance=VALUES(distance);",
			},
			{
				name: "InsertResult",
				query: "INSERT INTO result (person_id, seconds, milliseconds, chip_seconds, chip_milliseconds, " +
					"segment, location, occurence, ranking, age_ranking, gender_ranking, finish, result_type, " +
					"result_created_at, result_updated_at" +
					") SELECT person_id, seconds, milliseconds, chip_seconds, chip_milliseconds, " +
					"segment, location, occurence, ranking, age_ranking, gender_ranking, finish, result_type, " +
					"result_created_at, result_updated_at FROM result_old r JOIN person p ON (r.bib = p.bib " +
					"AND r.event_year_id=p.event_year_id);",
			},
			{
				name:  "DeleteResult",
				query: "DROP TABLE result_old;",
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
	if oldVersion < 4 && newVersion >= 4 {
		log.Info("Updating to database version 4.")
		_, err := tx.ExecContext(
			ctx,
			"ALTER TABLE event ADD COLUMN event_type VARCHAR(20) DEFAULT 'distance';",
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error updating from version %d to %d: %v", oldVersion, newVersion, err)
		}
	}
	if oldVersion < 5 && newVersion >= 5 {
		log.Info("Updating to database version 5.")
		_, err := tx.ExecContext(
			ctx,
			"ALTER TABLE api_key ADD COLUMN key_name VARCHAR(100) NOT NULL DEFAULT '';",
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error updating from verison %d to %d: %v", oldVersion, newVersion, err)
		}
	}
	if oldVersion < 6 && newVersion >= 6 {
		log.Info("Updating to database version 6.")
		_, err := tx.ExecContext(
			ctx,
			"ALTER TABLE person MODIFY gender VARCHAR(5);",
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
				"ADD COLUMN chip VARCHAR(200) DEFAULT '', "+
				"ADD COLUMN anonymous SMALLINT NOT NULL DEFAULT 0;",
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error updating from version %d to %d: %v", oldVersion, newVersion, err)
		}
	}
	if oldVersion < 8 && newVersion >= 8 {
		log.Info("Updating to database version 8.")
		_, err := tx.ExecContext(
			ctx,
			"ALTER TABLE person MODIFY gender VARCHAR(50);",
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error updating from version %d to %d: %v", oldVersion, newVersion, err)
		}
	}
	if oldVersion < 9 && newVersion >= 9 {
		log.Info("Updating to database version 9.")
		_, err := tx.ExecContext(
			ctx,
			"ALTER TABLE event MODIFY slug VARCHAR(50);",
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error updating from version %d to %d: %v", oldVersion, newVersion, err)
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
					"banned_phone VARCHAR(20), " +
					"CONSTRAINT unique_banned_phone UNIQUE(banned_phone)" +
					");",
			},
			{
				name: "CreateBannedEmails",
				query: "CREATE TABLE IF NOT EXISTS banned_emails(" +
					"banned_email VARCHAR(200), " +
					"CONSTRAINT unique_banned_email UNIQUE(banned_email)" +
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
				name: "CreateNewPersonTable",
				query: "CREATE TABLE IF NOT EXISTS new_person(" +
					"person_id BIGINT NOT NULL AUTO_INCREMENT, " +
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
				name: "ResultTable",
				query: "CREATE TABLE IF NOT EXISTS new_result(" +
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
					"result_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
					"CONSTRAINT one_occurrence_res UNIQUE (person_id, location, occurence)," +
					"FOREIGN KEY (person_id) REFERENCES new_person(person_id)" +
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
				name:  "InsertPersonData",
				query: "INSERT INTO new_result SELECT * FROM result;",
			},
			{
				name:  "DropOld",
				query: "DROP TABLE result, person;",
			},
			{
				name:  "RenamePerson",
				query: "ALTER TABLE new_person RENAME TO person;",
			},
			{
				name:  "RenameResult",
				query: "ALTER TABLE new_result RENAME TO result;",
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
				name:  "DropPersonColumns1",
				query: "ALTER TABLE person DROP COLUMN chip;",
			},
			{
				name:  "DropPersonColumns2",
				query: "ALTER TABLE person DROP COLUMN sms_enabled;",
			},
			{
				name: "CreateParticipantTable",
				query: "CREATE TABLE IF NOT EXISTS participant(" +
					"participant_id BIGINT NOT NULL AUTO_INCREMENT, " +
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
					"FOREIGN KEY (event_year_id) REFERENCES event_year(event_year_id), " +
					"PRIMARY KEY (participant_id)" +
					");",
			},
			{
				name: "CreateChipTable",
				query: "CREATE TABLE IF NOT EXISTS chips(" +
					"chip_id BIGINT NOT NULL AUTO_INCREMENT, " +
					"event_year_id BIGINT NOT NULL, " +
					"bib VARCHAR(100) NOT NULL, " +
					"chip VARCHAR(100) NOT NULL, " +
					"CONSTRAINT unique_combo UNIQUE (event_year_id, chip), " +
					"FOREIGN KEY (event_year_id) REFERENCES event_year(event_year_id), " +
					"PRIMARY KEY (chip_id)" +
					");",
			},
			{
				name: "CreateSegmentTable",
				query: "CREATE TABLE IF NOT EXISTS segments(" +
					"segment_id BIGINT NOT NULL AUTO_INCREMENT, " +
					"event_year_id BIGINT NOT NULL, " +
					"location_name VARCHAR(100) NOT NULL, " +
					"distance_name VARCHAR(100) NOT NULL, " +
					"segment_name VARCHAR(100) NOT NULL, " +
					"segment_distance DECIMAL(10,2) NOT NULL DEFAULT 0.0, " +
					"segment_distance_unit VARCHAR(12) NOT NULL DEFAULT '', " +
					"segment_gps VARCHAR(500) NOT NULL DEFAULT '', " +
					"segment_map_link VARCHAR(500) NOT NULL DEFAULT '', " +
					"CONSTRAINT unique_segment UNIQUE (event_year_id, distance_name, segment_name), " +
					"FOREIGN KEY (event_year_id) REFERENCES event_year(event_year_id), " +
					"PRIMARY KEY (segment_id)" +
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
					"distance_id BIGINT NOT NULL AUTO_INCREMENT, " +
					"event_year_id BIGINT NOT NULL, " +
					"distance_name VARCHAR(100) NOT NULL, " +
					"certification VARCHAR(150) NOT NULL, " +
					"CONSTRAINT unique_distance UNIQUE (event_year_id, distance_name), " +
					"FOREIGN KEY (event_year_id) REFERENCES event_year(event_year_id), " +
					"PRIMARY KEY (distance_id)" +
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
		"UPDATE settings SET value=? WHERE name='version';",
		newVersion,
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

func (m *MySQL) updateDB(newdb *sql.DB) {
	m.db = newdb
}

// Close Closes database.
func (m *MySQL) Close() {
	m.db.Close()
}
