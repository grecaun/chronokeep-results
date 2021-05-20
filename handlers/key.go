package handlers

import (
	"chronokeep/results/types"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h Handler) GetKeys(c echo.Context) error {
	var request types.GetKeysRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	account, err := verifyToken(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized Token", err)
	}
	if account.Type != "admin" && account.Email != request.Email {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	keyAccount, err := database.GetAccount(request.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if keyAccount == nil {
		return getAPIError(c, http.StatusNotFound, "Account Not Found", nil)
	}
	keys, err := database.GetAccountKeys(keyAccount.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	return c.JSON(http.StatusOK, types.GetKeysResponse{
		Keys: keys,
	})
}

func (h Handler) AddKey(c echo.Context) error {
	var request types.AddKeyRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	account, err := verifyToken(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized Token", err)
	}
	if err := request.Key.Validate(h.validate); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Field(s)", err)
	}
	if account.Type != "admin" && account.Email != request.Email {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	keyAccount, err := database.GetAccount(request.Email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if keyAccount == nil {
		return getAPIError(c, http.StatusNotFound, "Account Not Found", nil)
	}
	// Create new API Key for our key to add.
	newKey, err := uuid.NewRandom()
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Key Generation Error", err)
	}
	key, err := database.AddKey(types.Key{
		AccountIdentifier: keyAccount.Identifier,
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
	var request types.DeleteKeyRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	account, err := verifyToken(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized Token", err)
	}
	// Get Key to be deleted.
	multiKey, err := database.GetKeyAndAccount(request.Key)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if multiKey == nil || multiKey.Key == nil || multiKey.Account == nil {
		return getAPIError(c, http.StatusNotFound, "Key Not Found", nil)
	}
	// Deny access to non admins who do not own the key
	if account.Type != "admin" && account.Email != multiKey.Account.Email {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	err = database.DeleteKey(*multiKey.Key)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	return c.NoContent(http.StatusOK)
}

func (h Handler) UpdateKey(c echo.Context) error {
	var request types.UpdateKeyRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	account, err := verifyToken(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized Token", err)
	}
	if err := request.Key.Validate(h.validate); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Field(s)", err)
	}
	// Get Account associated with this key
	keyAccount, err := database.GetAccountByKey(request.Key.Value)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	if keyAccount == nil {
		return getAPIError(c, http.StatusNotFound, "Key Not Found", nil)
	}
	// Deny access to non admins who do not own the key
	if account.Type != "admin" && account.Email != keyAccount.Email {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	err = database.UpdateKey(request.Key)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Unable To Update Key", err)
	}
	key, err := database.GetKey(request.Key.Value)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Database Error", err)
	}
	return c.JSON(http.StatusOK, types.ModifyKeyResponse{
		Key: *key,
	})
}
