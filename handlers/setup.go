package handlers

import (
	db "chronokeep/results/database"
	"chronokeep/results/database/mysql"
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
		return errors.New("postgres not supported")
	default:
		return errors.New("unknown sql driver specified")
	}
}

func Finalize() {
	database.Close()
}

func (h *Handler) Setup() {
	// Set up Validator.
	h.validate = validator.New()
}
