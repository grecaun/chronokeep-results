package handlers

import (
	"chronokeep/results/types"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func (h Handler) GetKeys(c echo.Context) error {
	var request types.GetKeysRequest
	_ = c.Bind(&request)
	account, err := verifyToken(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized Token", err)
	}
	if account.Type != "admin" && request.Email != nil && account.Email != *request.Email {
		return getAPIError(c, http.StatusUnauthorized, "Not an Admin / Ownership Error", nil)
	}
	email := account.Email
	if request.Email != nil {
		keyAccount, err := database.GetAccount(*request.Email)
		if err != nil {
			return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Key Account", err)
		}
		if keyAccount == nil {
			return getAPIError(c, http.StatusNotFound, "Account Not Found", nil)
		}
		email = keyAccount.Email
	}
	keys, err := database.GetAccountKeys(email)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Keys", err)
	}
	return c.JSON(http.StatusOK, types.GetKeysResponse{
		Keys: keys,
	})
}

func (h Handler) AddKey(c echo.Context) error {
	log.Info("Adding new Key; Binding request body to key struct.")
	var request types.AddKeyRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	log.Info("Verifying token.")
	account, err := verifyToken(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Unauthorized Token", err)
	}
	log.Info("Validating key.")
	if err := request.Key.Validate(h.validate); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Field(s)", err)
	}
	log.Info("Checking for admin or ownership.")
	if account.Type != "admin" && request.Email != nil && account.Email != *request.Email {
		return getAPIError(c, http.StatusUnauthorized, "Not an Admin / Ownership Error", nil)
	}
	// If email is set we add a key to that account, otherwise add it to the calling person's account.
	accountid := account.Identifier
	if request.Email != nil {
		log.Info("Getting key account holder for id value.")
		keyAccount, err := database.GetAccount(*request.Email)
		if err != nil {
			return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Key Account", err)
		}
		if keyAccount == nil {
			return getAPIError(c, http.StatusNotFound, "Account Not Found", nil)
		}
		accountid = keyAccount.Identifier
	}
	log.Info("Adding key to database.")
	// Create new API Key for our key to add.
	newKey, err := uuid.NewRandom()
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Key Generation Error", err)
	}
	key, err := database.AddKey(types.Key{
		AccountIdentifier: accountid,
		Name:              request.Key.Name,
		Value:             newKey.String(),
		Type:              request.Key.Type,
		AllowedHosts:      strings.TrimSpace(request.Key.AllowedHosts),
		ValidUntil:        request.Key.GetValidUntil(),
	})
	if err != nil || key == nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Adding Key", err)
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
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Key/Account", err)
	}
	if multiKey == nil || multiKey.Key == nil || multiKey.Account == nil {
		return getAPIError(c, http.StatusNotFound, "Key Not Found", nil)
	}
	// Deny access to non admins who do not own the key
	if account.Type != "admin" && account.Email != multiKey.Account.Email {
		return getAPIError(c, http.StatusUnauthorized, "Not an Admin / Ownership Error", nil)
	}
	err = database.DeleteKey(*multiKey.Key)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Deleting Key", err)
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
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Key Account", err)
	}
	if keyAccount == nil {
		return getAPIError(c, http.StatusNotFound, "Key Not Found", nil)
	}
	// Deny access to non admins who do not own the key
	if account.Type != "admin" && account.Email != keyAccount.Email {
		return getAPIError(c, http.StatusUnauthorized, "Not an Admin / Ownership Error", nil)
	}
	err = database.UpdateKey(request.Key.ToKey())
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Unable To Update Key", err)
	}
	key, err := database.GetKey(request.Key.Value)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Updated Key", err)
	}
	if key == nil {
		return getAPIError(c, http.StatusNotFound, "Key Not Found After Update", nil)
	}
	return c.JSON(http.StatusOK, types.ModifyKeyResponse{
		Key: *key,
	})
}
