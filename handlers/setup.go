package handlers

import (
	db "chronokeep/results/database"
	"chronokeep/results/database/mysql"
	"chronokeep/results/util"
	"errors"
)

var (
	database db.Database
)

func Setup(config *util.Config) error {
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
