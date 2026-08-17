/*
Chronokeep Desktop - Race Scoring Software
Copyright (C) 2026 James Sentinella

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

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

