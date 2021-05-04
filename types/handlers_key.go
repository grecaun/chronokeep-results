package types

/*
	Responses
*/

// ModifyKeyResponse Struct used to respond to a Add/Update Key Request.
type ModifyKeyResponse struct {
	Code     int    `json:"code"`
	Response string `json:"response"`
	Key      Key    `json:"key"`
}

// GetKeysResponse Struct used to respond to the requets for account keys.
type GetKeysResponse struct {
	Code     int    `json:"code"`
	Response string `json:"response"`
	Keys     []Key  `json:"keys"`
}

/*
	Requests
*/

// DeleteKeyRequest Struct used for the Delete Key request.
type DeleteKeyRequest struct {
	Key string `json:"key"`
}

// AddKeyRequest Struct used for the Add Key request.
type AddKeyRequest struct {
	Email string `json:"email"`
	Key   Key    `json:"key"`
}

// UpdateKeyRequest Struct used for the Update Key request.
type UpdateKeyRequest struct {
	Key Key `json:"key"`
}

// GetKeysRequest Struct used for the Get Keys request.
type GetKeysRequest struct {
	Email string `json:"email"`
}
