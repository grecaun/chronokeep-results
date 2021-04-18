package types

// GenderalRequest A generalized request struct when only a key is required for the call.
type GenderalRequest struct {
	Key string `json:"key"`
}

/*
	Handlers related to Events
*/

// GetEventsResponse Struct used for the response of a Get Events request.
type GetEventsResponse struct {
	Count  int     `json:"count"`
	Events []Event `json:"events"`
}

// GetAllEventsRequest Struct used for the request to Get Events based on account.
type GetEventsRequest struct {
	Key     string `json:"key"`
	Account string `json:"account"`
}

// GetEventResponse Struct used for the response of a Get Event request.
type GetEventResponse struct {
	Event Event `json:"event"`
}

// GetEventRequest Struct used for the request to get a single Event.
type GetEventRequest struct {
	Key             string `json:"key"`
	EventIdentifier string `json:"eventID"`
	EventSlug       string `json:"eventSlug"`
}

// ModifyEventResponse Struct used for the response of a Add/Update Event request.
type ModifyEventResponse struct {
	Comment string `json:"comment"`
	Event   Event  `json:"event"`
}

// AddEventRequest Struct used for the request to Add an Event.
type AddEventRequest struct {
	Key   string `json:"key"`
	Event Event  `json:"event"`
}

// UpdateEventRequest Struct used for the request to Update an Event.
type UpdateEventRequest struct {
	Key   string `json:"key"`
	Event Event  `json:"event"`
}

// DeleteEventRequest Struct used for the request to Delete an Event.
type DeleteEventRequest struct {
	Key             string `json:"key"`
	EventIdentifier string `json:"eventID"`
}

/*
	Handlers related to an EventYear
*/

// GetEventYearResponse Struct used for the response of a Get EventYear request.
type GetEventYearResponse struct {
	Event     Event     `json:"event"`
	EventYear EventYear `json:"eventYear"`
}

// GetEventYearRequest Struct used for the request of an Event Year.
type GetEventYearRequest struct {
	Key             string `json:"key"`
	EventIdentifier string `json:"eventID"`
	Year            string `json:"year"`
}

// ModifyEventYearResponse Struct used for the response of an Add/Update EventYear Request.
type ModifyEventYearResponse struct {
	Comment   string    `json:"comment"`
	Event     Event     `json:"event"`
	EventYear EventYear `json:"eventYear"`
}

// AddEventYearRequest Struct used to add an Event Year.
type AddEventYearRequest struct {
	Key             string    `json:"key"`
	EventIdentifier string    `json:"eventID"`
	EventYear       EventYear `json:"eventYear"`
}

// UpdateEventYearRequest Struct used to update an Event Year.
type UpdateEventYearRequest struct {
	Key       string    `json:"key"`
	EventYear EventYear `json:"eventYear"`
}

// DeleteEventYearRequest Struct used to delete an Event Year.
type DeleteEventYearRequest struct {
	Key                 string `json:"key"`
	EventYearIdentifier string `json:"eventYearID"`
}

/*
	Handlers realted to results.
*/

// GetResultsResponse Struct used for the response of a GetResults request.
type GetResultsResponse struct {
	Results []Result `json:"results"`
}

// GetResultsRequest Struct used for the request of Results for an EventYear
type GetResultsRequest struct {
	Key                 string `json:"key"`
	EventYearIdentifier string `json:"eventYearID"`
}

// ModifyResultsResponse Struct used for the response to an Add/Update/Delete Results request.
type ModifyResultsResponse struct {
	Count   int    `json:"count"`
	Comment string `json:"comment"`
}

// ModifyResultsRequest Struct used to add/update Results.
type ModifyResultsRequest struct {
	Key                 string   `json:"key"`
	EventYearIdentifier string   `json:"eventYearID"`
	Results             []Result `json:"results"`
}

// DeleteResultsRequest Struct used to delete results for an EventYear.
type DeleteResultsRequest struct {
	Key                 string `json:"key"`
	EventYearIdentifier string `json:"eventYearID"`
}

/*
	Handlers related to api keys
*/

// ModifyKeyResponse Struct used to respond to a Add/Update Key Request.
type ModifyKeyResponse struct {
	Comment string `json:"comment"`
	Key     Key    `json:"key"`
}

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

/*
	Handlers related to accounts.
*/

// GetAccountResponse Struct used for the response of the Get Account Request.
type GetAccountResponse struct {
	Account Account `json:"account"`
	Keys    []Key   `json:"keys"`
	Events  []Event `json:"events"`
}

// GetAccountRequest Struct used to request information about a specific account.
type GetAccountRequest struct {
	AccountIdentifier string `json:"accountID"`
}

// GetAllAccountsResponse Struct used to get all of the accounts.
type GetAllAccountsResponse struct {
	Accounts []Account `json:"accounts"`
}

// ModifyAccountResponse Struct used to respond to add/update account requests.
type ModifyAccountResponse struct {
	Comment string  `json:"comment"`
	Account Account `json:"account"`
}

// ModifyAccountRequest Struct used to request to add/update an account.
type ModifyAccountRequest struct {
	Account Account `json:"account"`
}

// DeleteAccountRequest Struct used to delete an account.
type DeleteAccountRequest struct {
	AccountIdentifier string `json:"accountID"`
}
