package postgres

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

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
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
	if inCfg == nil {
		return nil, fmt.Errorf("no valid config supplied")
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

	if !inCfg.Development {
		conString = conString + "?sslmode=require"
	}

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
	if config == nil {
		return fmt.Errorf("no valid config supplied")
	}
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
		err = p.createTables()
		if err != nil {
			return err
		}
		// Otherwise check if our database is out of date and update if necessary.
	} else if dbVersion < database.CurrentVersion {
		log.Info(fmt.Sprintf("Updating database from version %v to %v", dbVersion, database.CurrentVersion))
		err = p.updateTables(dbVersion, database.CurrentVersion)
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
		"DROP TABLE call_record, result, person, event_year, event, api_key, account, settings;",
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
				"key_name VARCHAR(100) NOT NULL DEFAULT ''," +
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
		// PERSON TABLE
		{
			name: "PersonTable",
			query: "CREATE TABLE IF NOT EXISTS person(" +
				"person_id BIGSERIAL NOT NULL, " +
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
				"result_created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP, " +
				"result_updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP," +
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

	if p.db == nil {
		return fmt.Errorf("database not setup")
	}

	// Get a context and cancel function to create our tables, defer the cancel until we're done.
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()

	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("unable to start transaction: %v", err)
	}
	for _, single := range queries {
		log.Info(fmt.Sprintf("Executing query for: %s", single.name))
		_, err := tx.Exec(ctx, single.query)
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("error executing %s query: %v", single.name, err)
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("unable to commit transaction: %v", err)
	}

	p.SetSetting("version", strconv.Itoa(database.CurrentVersion))

	return nil
}

func (p *Postgres) checkVersion() int {
	log.Info("Checking database version.")
	if p.db == nil {
		return -1
	}
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
	if p.db == nil {
		return fmt.Errorf("database not set up")
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()

	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("unable to start transaction: %v", err)
	}
	if oldVersion < 2 && newVersion >= 2 {
		log.Debug("Updating to database version 2.")
		_, err := tx.Exec(
			ctx,
			"ALTER TABLE result ADD COLUMN result_type INT DEFAULT 0;",
		)
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("error updating from version %d to %d: %v", oldVersion, newVersion, err)
		}
	}
	if oldVersion < 3 && newVersion >= 3 {
		log.Debug("Updating to database version 3.")
		queries := []myQuery{
			{
				name:  "RenameResult",
				query: "ALTER TABLE result RENAME TO result_old;",
			},
			{
				name: "CreatePerson",
				query: "CREATE TABLE IF NOT EXISTS person(" +
					"person_id BIGSERIAL NOT NULL, " +
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
					"result_created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP, " +
					"result_updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP," +
					"CONSTRAINT one_occurrence_res UNIQUE (person_id, location, occurence)," +
					"FOREIGN KEY (person_id) REFERENCES person(person_id)" +
					");",
			},
			{
				name: "InsertPerson",
				query: "INSERT INTO person (event_year_id, bib, first, last, age, gender, age_group, distance) " +
					" SELECT event_year_id, bib, first, last, age, gender, age_group, distance FROM result_old " +
					"ON CONFLICT (event_year_id, bib) DO NOTHING;",
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
			_, err := tx.Exec(
				ctx,
				q.query,
			)
			if err != nil {
				tx.Rollback(ctx)
				return fmt.Errorf("error updating from version %d to %d in query %s: %v", oldVersion, newVersion, q.name, err)
			}
		}
	}
	if oldVersion < 4 && newVersion >= 4 {
		log.Debug("Updating to database version 4.")
		_, err := tx.Exec(
			ctx,
			"ALTER TABLE event ADD COLUMN event_type VARCHAR(20) DEFAULT 'distance';",
		)
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("error updating from verison %d to %d: %v", oldVersion, newVersion, err)
		}
	}
	if oldVersion < 5 && newVersion >= 5 {
		log.Debug("Updating to database version 5.")
		_, err := tx.Exec(
			ctx,
			"ALTER TABLE api_key ADD COLUMN key_name VARCHAR(100) NOT NULL DEFAULT '';",
		)
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("error updating from verison %d to %d: %v", oldVersion, newVersion, err)
		}
	}
	if oldVersion < 6 && newVersion >= 6 {
		log.Debug("Updating to database version 6.")
		_, err := tx.Exec(
			ctx,
			"ALTER TABLE person ALTER COLUMN gender TYPE VARCHAR(5);",
		)
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("error updating from verison %d to %d: %v", oldVersion, newVersion, err)
		}
	}
	if oldVersion < 7 && newVersion >= 7 {
		log.Debug("Updating to database version 7.")
		_, err := tx.Exec(
			ctx,
			"ALTER TABLE person "+
				"ADD COLUMN chip VARCHAR(200) DEFAULT '', "+
				"ADD COLUMN anonymous SMALLINT NOT NULL DEFAULT 0;",
		)
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("error updating from version %d to %d: %v", oldVersion, newVersion, err)
		}
	}
	if oldVersion < 8 && newVersion >= 8 {
		log.Debug("Updating to database version 8.")
		_, err := tx.Exec(
			ctx,
			"ALTER TABLE person ALTER COLUMN gender TYPE VARCHAR(50) NOT NULL;",
		)
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("error updating from version %d to %d: %v", oldVersion, newVersion, err)
		}
	}
	_, err = tx.Exec(
		ctx,
		"UPDATE settings SET value=$1 WHERE name='version';",
		strconv.Itoa(newVersion),
	)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("error updating from version %d to %d: %v", oldVersion, newVersion, err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("error committing transaction: %v", err)
	}
	return nil
}

func (p *Postgres) updateDB(newdb *pgxpool.Pool) {
	p.db = newdb
}

// Close Closes database.
func (p *Postgres) Close() {
	p.db.Close()
}
