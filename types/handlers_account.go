package types

/*
	Responses
*/

// GetAccountResponse Struct used for the response of the Get Account Request.
type GetAccountResponse struct {
	Account Account `json:"account"`
	Keys    []Key   `json:"keys"`
	Events  []Event `json:"events"`
}

// GetAllAccountsResponse Struct used to get all of the accounts.
type GetAllAccountsResponse struct {
	Accounts []Account `json:"accounts"`
}

// ModifyAccountResponse Struct used to respond to add/update account requests.
type ModifyAccountResponse struct {
	Account Account `json:"account"`
}

/*
	Requests
*/

// GetAccountRequest Struct used to request information about a specific account.
type GetAccountRequest struct {
	Ident *string `json:"identifier"`
}

// ModifyAccountRequest Struct used to request to update an account.
type UpdateAccountRequest struct {
	Account Account `json:"account"`
}

// AddAccountRequest Struct used to add an account.
type AddAccountRequest struct {
	Account  Account `json:"account"`
	Password string  `json:"password"`
}

// DeleteAccountRequest Struct used to delete an account.
type DeleteAccountRequest struct {
	Ident string `json:"identifier"`
}

/*
	Accounts API Responses
*/

// GetAccountResponse Struct used for the response of the Get/Modify/Add Account Requests.
type AccountsSingleResponse struct {
	Account     *Account `json:"account"`
	UserAccount *Account `json:"user"`
	Message     *string  `json:"message"`
}

/*
	Accounts API Requests
*/

// GetAccountRequest Struct used to request information about a specific account.
type AccountsGetAccountRequest struct {
	Token string  `json:"token"`
	Ident *string `json:"identifier"`
}

// ModifyAccountRequest Struct used to request to update an account.
type AccountsUpdateAccountRequest struct {
	Token   string  `json:"token"`
	Account Account `json:"account"`
}
