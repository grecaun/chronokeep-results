package handlers

import (
	db "chronokeep/results/database"
	"chronokeep/results/database/mysql"
	"chronokeep/results/database/postgres"
	"chronokeep/results/util"
	"errors"

	"github.com/go-playground/validator/v10"
)

var (
	database db.Database
	config   *util.Config
)

func Setup(inCfg *util.Config) error {
	config = inCfg
	switch config.DBDriver {
	case "mysql":
		database = &mysql.MySQL{}
		return database.Setup(config)
	case "postgres":
		database = &postgres.Postgres{}
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
