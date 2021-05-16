package handlers

import (
	"chronokeep/results/types"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h Handler) GetKeys(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	var request types.GetKeysRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	if claims["type"].(string) != "admin" && claims["email"].(string) != request.Email {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	account, err := database.GetAccount(request.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if account == nil {
		return getAPIError(c, http.StatusNotFound, "Account Not Found", nil)
	}
	keys, err := database.GetAccountKeys(account.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	return c.JSON(http.StatusOK, types.GetKeysResponse{
		Keys: keys,
	})
}

func (h Handler) AddKey(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	var request types.AddKeyRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	if claims["type"].(string) != "admin" && claims["email"].(string) != request.Email {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	account, err := database.GetAccount(request.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if account == nil {
		return getAPIError(c, http.StatusNotFound, "Account Not Found", nil)
	}
	// Create new API Key for our key to add.
	newKey, err := uuid.NewRandom()
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Key Generation Error", err)
	}
	key, err := database.AddKey(types.Key{
		AccountIdentifier: account.Identifier,
		Value:             newKey.String(),
		Type:              request.Key.Type,
		AllowedHosts:      request.Key.AllowedHosts,
		ValidUntil:        request.Key.ValidUntil,
	})
	if err != nil || key == nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	return c.JSON(http.StatusOK, types.ModifyKeyResponse{
		Key: *key,
	})
}

func (h Handler) DeleteKey(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	var request types.DeleteKeyRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	// Get Key to be deleted.
	key, err := database.GetKey(request.Key)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if key == nil {
		return getAPIError(c, http.StatusNotFound, "Key Not Found", nil)
	}
	// Get Account associated with this key
	account, err := database.GetAccountByID(key.AccountIdentifier)
	if err != nil || account == nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	// Deny access to non admins who do not own the key
	if claims["type"].(string) != "admin" && claims["email"].(string) != account.Email {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	err = database.DeleteKey(*key)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	return c.NoContent(http.StatusOK)
}

func (h Handler) UpdateKey(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	var request types.UpdateKeyRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	// Get Account associated with this key
	account, err := database.GetAccountByKey(request.Key.Value)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if account == nil {
		return getAPIError(c, http.StatusNotFound, "Key Not Found", nil)
	}
	// Deny access to non admins who do not own the key
	if claims["type"].(string) != "admin" && claims["email"].(string) != account.Email {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	err = database.UpdateKey(request.Key)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	key, err := database.GetKey(request.Key.Value)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	return c.JSON(http.StatusOK, types.ModifyKeyResponse{
		Key: *key,
	})
}
