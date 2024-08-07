package util

import (
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

// GetConfig returns a config struct filled with values stored in local environment variables
func GetConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Info(".env file wasn't loaded.")
	}
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
		return nil, errors.New("DB_PORT not found in environment")
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

	secret_key := os.Getenv("SECRET_KEY")
	if secret_key == "" || len(secret_key) < 20 {
		return nil, errors.New("SECRET_KEY not set or under 20 characters")
	}

	refresh_key := os.Getenv("REFRESH_KEY")
	if refresh_key == "" || len(refresh_key) < 20 {
		return nil, errors.New("REFRESH_KEY not set or under 20 characters")
	}

	certificate_url := os.Getenv("CERTIFICATE_TEMPLATE_URL")
	if certificate_url == "" || len(certificate_url) < 10 {
		return nil, errors.New("CERTIFICATE_TEMPLATE_URL not set or under 10 characters")
	}

	admin_email := os.Getenv("ADMIN_EMAIL")
	admin_name := os.Getenv("ADMIN_NAME")
	admin_pass := os.Getenv("ADMIN_PASS")

	twilio_auth_token := os.Getenv("TWILIO_AUTH_TOKEN")
	twilio_response_webhook_url := os.Getenv("TWILIO_RESPONSE_WEBHOOK_URL")

	domain := os.Getenv("DOMAIN")

	return &Config{
		DBName:                   dbName,
		DBHost:                   dbHost,
		DBPort:                   dbPort,
		DBUser:                   dbUser,
		DBPassword:               dbPassword,
		DBDriver:                 dbDriver,
		RecordInterval:           recordInterval,
		Port:                     port,
		Development:              development,
		AutoTLS:                  autotls,
		SecretKey:                secret_key,
		RefreshKey:               refresh_key,
		AdminEmail:               admin_email,
		AdminName:                admin_name,
		AdminPass:                admin_pass,
		Domain:                   domain,
		CertificateURL:           certificate_url,
		TwilioAuthToken:          twilio_auth_token,
		TwilioResponseWebhookURL: twilio_response_webhook_url,
	}, nil
}

// Config is the struct that holds all of the config values for connecting to a database
type Config struct {
	DBName                   string
	DBHost                   string
	DBPort                   int
	DBUser                   string
	DBPassword               string
	DBDriver                 string
	RecordInterval           int
	Port                     int
	Development              bool
	AutoTLS                  bool
	SecretKey                string
	RefreshKey               string
	AdminEmail               string
	AdminName                string
	AdminPass                string
	Domain                   string
	CertificateURL           string
	TwilioAuthToken          string
	TwilioResponseWebhookURL string
}
