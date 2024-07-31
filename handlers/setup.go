package handlers

import (
	db "chronokeep/results/database"
	"chronokeep/results/database/mysql"
	"chronokeep/results/database/postgres"
	"chronokeep/results/database/sqlite"
	"chronokeep/results/util"
	"errors"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"github.com/twilio/twilio-go/client"
)

var (
	database               db.Database
	config                 *util.Config
	twilioRequestValidator client.RequestValidator
)

func Setup(inCfg *util.Config) error {
	config = inCfg
	twilioRequestValidator = client.NewRequestValidator(config.TwilioAuthToken)
	switch config.DBDriver {
	case "mysql":
		log.Info("Database set to MySQL")
		database = &mysql.MySQL{}
		return database.Setup(config)
	case "postgres":
		log.Info("Database set to Postgresql")
		database = &postgres.Postgres{}
		return database.Setup(config)
	case "sqlite3":
		log.Info("Database set to SQLite")
		database = &sqlite.SQLite{}
		return database.Setup(config)
	default:
		return errors.New("unknown database driver specified")
	}
}

func Finalize() {
	database.Close()
}

func (h *Handler) Setup() {
	// Set up Validator.
	h.validate = validator.New()
}
