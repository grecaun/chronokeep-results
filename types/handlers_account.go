package types

/*
	Responses
*/

// GetAccountResponse Struct used for the response of the Get Account Request.
type GetAccountResponse struct {
	Code     int     `json:"code"`
	Response string  `json:"response"`
	Account  Account `json:"account"`
	Keys     []Key   `json:"keys"`
	Events   []Event `json:"events"`
}

// GetAllAccountsResponse Struct used to get all of the accounts.
type GetAllAccountsResponse struct {
	Code     int       `json:"code"`
	Response string    `json:"response"`
	Accounts []Account `json:"accounts"`
}

// ModifyAccountResponse Struct used to respond to add/update account requests.
type ModifyAccountResponse struct {
	Code     int     `json:"code"`
	Response string  `json:"response"`
	Account  Account `json:"account"`
}

/*
	Requests
*/

// GetAccountRequest Struct used to request information about a specific account.
type GetAccountRequest struct {
	Email string `json:"email"`
}

// ModifyAccountRequest Struct used to request to add/update an account.
type ModifyAccountRequest struct {
	Account Account `json:"account"`
}

// DeleteAccountRequest Struct used to delete an account.
type DeleteAccountRequest struct {
	Email string `json:"email"`
}

// LoginRequest Struct used to login to an account.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ChangeEmailRequest Struct used to change the email on an account.
type ChangeEmailRequest struct {
	OldEmail string `json:"oldEmail"`
	NewEmail string `json:"newEmail"`
}

// ChangePasswordRequest Struct used to change the password on an account.
type ChangePasswordRequest struct {
	Email       string `json:"email"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}
