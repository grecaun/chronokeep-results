package util

import (
	"os"
	"strconv"

	"github.com/pkg/errors"
)

// GetConfig returns a config struct filled with values stored in local environment variables
func GetConfig() (*Config, error) {
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		return nil, errors.New("DB_NAME not found in environment")
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		return nil, errors.New("DB_HOST not found in environment")
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil || dbPort < 0 {
		dbPort = 3306
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		return nil, errors.New("DB_USER not found in environment")
	}

	dbPassword := os.Getenv("DB_PASSWORD")

	dbDriver := os.Getenv("DB_CONNECTOR")
	if dbDriver == "" {
		dbDriver = "mysql"
	}

	recordInterval, err := strconv.Atoi(os.Getenv("RECORD_INTERVAL"))
	if err != nil && recordInterval < 60 {
		recordInterval = 300
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil && port < 200 {
		port = 8181
	}

	development := os.Getenv("VERSION") != "production"

	autotls := os.Getenv("AUTOTLS") == "enabled"

	session_auth := os.Getenv("SESSION_AUTH_KEY")
	if session_auth == "" || (len(session_auth) != 32 && len(session_auth) != 64) {
		return nil, errors.New("SESSION_AUTH_KEY not set or not 32/64 bytes")
	}

	session_encr := os.Getenv("SESSION_ENCRY_KEY")
	if session_encr == "" || (len(session_encr) != 32 && len(session_encr) != 64) {
		return nil, errors.New("SESSION_ENCRY_KEY not set or not 32/64 bytes")
	}

	admin_email := os.Getenv("ADMIN_EMAIL")
	admin_name := os.Getenv("ADMIN_NAME")
	admin_pass := os.Getenv("ADMIN_PASS")

	return &Config{
		DBName:         dbName,
		DBHost:         dbHost,
		DBPort:         dbPort,
		DBUser:         dbUser,
		DBPassword:     dbPassword,
		DBDriver:       dbDriver,
		RecordInterval: recordInterval,
		Port:           port,
		Development:    development,
		AutoTLS:        autotls,
		SessionAuthKey: session_auth,
		SessionEncrKey: session_encr,
		AdminEmail:     admin_email,
		AdminName:      admin_name,
		AdminPass:      admin_pass,
	}, nil
}

// Config is the struct that holds all of the config values for connecting to a database
type Config struct {
	DBName         string
	DBHost         string
	DBPort         int
	DBUser         string
	DBPassword     string
	DBDriver       string
	RecordInterval int
	Port           int
	Development    bool
	AutoTLS        bool
	SessionAuthKey string
	SessionEncrKey string
	AdminEmail     string
	AdminName      string
	AdminPass      string
}
