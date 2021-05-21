package postgres

import (
	"chronokeep/results/auth"
	"chronokeep/results/types"
	"chronokeep/results/util"

	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

const (
	MaxOpenConnections    = 20
	MaxIdleConnections    = 20
	MaxConnectionLifetime = time.Minute * 5
	CurrentVersion        = 1
)

type Postgres struct {
	db       *pgxpool.Pool
	config   *util.Config
	validate *validator.Validate
}

// GetDatabase Used to get a database with given configuration information.
func (p *Postgres) GetDatabase(inCfg *util.Config) (*pgxpool.Pool, error) {
	if p.db != nil {
		return p.db, nil
	}

	p.config = inCfg
	conString := fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s",
		p.config.DBDriver,
		p.config.DBUser,
		p.config.DBPassword,
		p.config.DBHost,
		p.config.DBPort,
		p.config.DBName,
	)

	dbCon, err := pgxpool.Connect(context.Background(), conString)
	if err != nil {
		return nil, fmt.Errorf("unable to open database connection: %v", err)
	}

	p.db = dbCon
	return p.db, nil
}

// GetDB Used as a general way to get a database.
func (p *Postgres) GetDB() (*pgxpool.Pool, error) {
	if p.db != nil {
		return p.db, nil
	}
	if p.config != nil {
		return p.GetDatabase(p.config)
	}
	return nil, errors.New("config file not established")
}

// Setup Automatically creates and updates tables for all of our information.
func (p *Postgres) Setup(config *util.Config) error {
	// Set up Validator.
	p.validate = validator.New()
	log.Info("Setting up database.")
	// Connect to DB with database name.
	_, err := p.GetDatabase(config)
	if err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	// Check our database version.
	dbVersion := p.checkVersion()

	// Error checking version, most likely means tables are not created.
	if dbVersion < 1 {
		log.Info("Creating database tables.")
		err = p.createTables()
		if err != nil {
			return err
		}
		// Otherwise check if our database is out of date and update if necessary.
	} else if dbVersion < CurrentVersion {
		log.Info(fmt.Sprintf("Updating database from version %v to %v", dbVersion, CurrentVersion))
		err = p.updateTables(dbVersion, CurrentVersion)
		if err != nil {
			return err
		}
	}

	// Check if there's an account created.
	accounts, err := p.GetAccounts()
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
		err = p.validate.Struct(acc)
		if err != nil {
			return fmt.Errorf("error validating base admin account on setup: %v", err)
		}
		acc.Password, err = auth.HashPassword(config.AdminPass)
		if err != nil {
			return fmt.Errorf("error hashing admin account password on setup: %v", err)
		}
		_, err = p.AddAccount(acc)
		if err != nil {
			return fmt.Errorf("error adding admin account on setup: %v", err)
		}
	}
	return nil
}

func (p *Postgres) dropTables() error {
	db, err := p.GetDatabase(p.config)
	if err != nil {
		return fmt.Errorf("error connecting to database to drop tables: %v", err)
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	_, err = db.Exec(
		ctx,
		"DROP TABLE call_record, result, event_year, event, api_key, account, settings;",
	)
	if err != nil {
		return fmt.Errorf("error dropping tables: %v", err)
	}
	return nil
}

func (p *Postgres) SetSetting(name, value string) error {
	db, err := p.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	_, err = db.Exec(
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

func (p *Postgres) createTables() error {
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
			query: "CREATE TRIGGER update_account_timestamp BEFORE UPDATE ON account FOR EACH ROW EXECUTE PROCEDURE account_timestamp_column();",
		},
		{
			name:  "KeyTableTrigger",
			query: "CREATE TRIGGER update_key_timestamp BEFORE UPDATE ON api_key FOR EACH ROW EXECUTE PROCEDURE key_timestamp_column();",
		},
		{
			name:  "EventTableTrigger",
			query: "CREATE TRIGGER update_event_timestamp BEFORE UPDATE ON event FOR EACH ROW EXECUTE PROCEDURE event_timestamp_column();",
		},
		{
			name:  "EventYearTableTrigger",
			query: "CREATE TRIGGER update_event_year_timestamp BEFORE UPDATE ON event_year FOR EACH ROW EXECUTE PROCEDURE event_year_timestamp_column();",
		},
		{
			name:  "ResultTableTrigger",
			query: "CREATE TRIGGER update_result_timestamp BEFORE UPDATE ON result FOR EACH ROW EXECUTE PROCEDURE result_timestamp_column();",
		},
	}

	// Get a context and cancel function to create our tables, defer the cancel until we're done.
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()

	for _, single := range queries {
		_, err := p.db.Exec(ctx, single.query)
		if err != nil {
			return fmt.Errorf("error executing %s query: %v", single.name, err)
		}
	}

	p.SetSetting("version", strconv.Itoa(CurrentVersion))

	return nil
}

func (p *Postgres) checkVersion() int {
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	var name string
	var version string
	err := p.db.QueryRow(
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

func (p *Postgres) updateTables(oldVersion, newVersion int) error {
	return nil
}

func (p *Postgres) updateDB(newdb *pgxpool.Pool) {
	p.db = newdb
}

// Close Closes database.
func (p *Postgres) Close() {
	p.db.Close()
}
