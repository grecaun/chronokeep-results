package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// APIError holds information on an error from the API
type APIError struct {
	Message string `json:"message,omitempty"`
}

func getAPIError(c echo.Context, code int, message string, err error) error {
	log.WithFields(log.Fields{
		"message": message,
		"error":   err,
		"code":    code,
	}).Error("API Error.")
	return c.JSON(code, APIError{Message: message})
}

func retrieveKey(r *http.Request) (*string, error) {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) != 2 {
		return nil, errors.New("unknown authorization header")
	}
	return &strArr[1], nil
}

func createTokens(email string) (*string, *string, error) {
	// Create token
	claims := jwt.MapClaims{}
	claims["email"] = email
	claims["authorized"] = true
	claims["exp"] = time.Now().Add(expirationWindow).Unix()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString([]byte(config.SecretKey))
	if err != nil {
		return nil, nil, err
	}
	// Create refresh token
	claims = jwt.MapClaims{}
	claims["email"] = email
	claims["exp"] = time.Now().Add(refreshWindow).Unix()
	r := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refresh, err := r.SignedString([]byte(config.RefreshKey))
	if err != nil {
		return nil, nil, err
	}
	return &token, &refresh, nil
}
