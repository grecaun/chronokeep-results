package util

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const testSecret = "a-very-secret-key-that-is-long-enough"

func setBaseEnv(t *testing.T) {
	t.Setenv("DB_NAME", "name")
	t.Setenv("DB_HOST", "host")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "user")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_CONNECTOR", "")
	t.Setenv("RECORD_INTERVAL", "")
	t.Setenv("PORT", "")
	t.Setenv("VERSION", "")
	t.Setenv("AUTOTLS", "")
	t.Setenv("SECRET_KEY", testSecret)
	t.Setenv("REFRESH_KEY", testSecret+"-refresh")
}

func TestGetConfigDefaults(t *testing.T) {
	setBaseEnv(t)
	cfg, err := GetConfig()
	if assert.NoError(t, err) {
		assert.Equal(t, "mysql", cfg.DBDriver)
		assert.Equal(t, 300, cfg.RecordInterval)
		assert.Equal(t, 8181, cfg.Port)
		assert.True(t, cfg.Development)
		assert.False(t, cfg.AutoTLS)
		assert.Equal(t, 5432, cfg.DBPort)
	}
}

func TestGetConfigValues(t *testing.T) {
	setBaseEnv(t)
	t.Setenv("DB_CONNECTOR", "postgres")
	t.Setenv("RECORD_INTERVAL", "600")
	t.Setenv("PORT", "8080")
	t.Setenv("VERSION", "production")
	t.Setenv("AUTOTLS", "enabled")
	cfg, err := GetConfig()
	if assert.NoError(t, err) {
		assert.Equal(t, "postgres", cfg.DBDriver)
		assert.Equal(t, 600, cfg.RecordInterval)
		assert.Equal(t, 8080, cfg.Port)
		assert.False(t, cfg.Development)
		assert.True(t, cfg.AutoTLS)
	}
}

func TestGetConfigInvalidNumbersUseDefaults(t *testing.T) {
	// A valid but too small record interval used to be accepted as is.
	for _, interval := range []string{"10", "0", "-5", "abc"} {
		setBaseEnv(t)
		t.Setenv("RECORD_INTERVAL", interval)
		cfg, err := GetConfig()
		if assert.NoError(t, err) {
			assert.Equal(t, 300, cfg.RecordInterval, interval)
		}
	}
	for _, port := range []string{"0", "-1", "70000", "abc"} {
		setBaseEnv(t)
		t.Setenv("PORT", port)
		cfg, err := GetConfig()
		if assert.NoError(t, err) {
			assert.Equal(t, 8181, cfg.Port, port)
		}
	}
	setBaseEnv(t)
	t.Setenv("PORT", "80")
	cfg, err := GetConfig()
	if assert.NoError(t, err) {
		assert.Equal(t, 80, cfg.Port)
	}
}

func TestGetConfigRequiredValues(t *testing.T) {
	for name, mod := range map[string]func(){
		"no db name":       func() { t.Setenv("DB_NAME", "") },
		"no db host":       func() { t.Setenv("DB_HOST", "") },
		"no db user":       func() { t.Setenv("DB_USER", "") },
		"bad db port":      func() { t.Setenv("DB_PORT", "abc") },
		"negative db port": func() { t.Setenv("DB_PORT", "-1") },
		"no secret":        func() { t.Setenv("SECRET_KEY", "") },
		"short secret":     func() { t.Setenv("SECRET_KEY", "short") },
		"no refresh":       func() { t.Setenv("REFRESH_KEY", "") },
		"short refresh":    func() { t.Setenv("REFRESH_KEY", "short") },
	} {
		setBaseEnv(t)
		mod()
		_, err := GetConfig()
		assert.Error(t, err, name)
	}
}

func signToken(t *testing.T, claims jwt.MapClaims, key string) string {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(key))
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func TestGetTokenAndEmail(t *testing.T) {
	t.Setenv("SECRET_KEY", testSecret)
	token := signToken(t, jwt.MapClaims{"email": "j@test.com", "exp": time.Now().Add(time.Hour).Unix()}, testSecret)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	// This used to panic dereferencing a nil email pointer.
	var tok, email *string
	var err error
	assert.NotPanics(t, func() {
		tok, email, err = GetTokenAndEmail(req)
	})
	if assert.NoError(t, err) && assert.NotNil(t, email) && assert.NotNil(t, tok) {
		assert.Equal(t, "j@test.com", *email)
		assert.Equal(t, token, *tok)
	}
}

func TestGetTokenAndEmailErrors(t *testing.T) {
	t.Setenv("SECRET_KEY", testSecret)
	cases := map[string]string{
		"no header":       "",
		"no bearer":       "token",
		"too many parts":  "Bearer a b",
		"garbage token":   "Bearer garbage",
		"wrong key":       "Bearer " + signToken(t, jwt.MapClaims{"email": "j@test.com"}, "another-secret-key-that-is-long"),
		"expired":         "Bearer " + signToken(t, jwt.MapClaims{"email": "j@test.com", "exp": time.Now().Add(-time.Hour).Unix()}, testSecret),
		"no email claim":  "Bearer " + signToken(t, jwt.MapClaims{"exp": time.Now().Add(time.Hour).Unix()}, testSecret),
		"non-string mail": "Bearer " + signToken(t, jwt.MapClaims{"email": 5}, testSecret),
	}
	for name, header := range cases {
		req := httptest.NewRequest("GET", "/", nil)
		if header != "" {
			req.Header.Set("Authorization", header)
		}
		_, _, err := GetTokenAndEmail(req)
		assert.Error(t, err, name)
	}
}
