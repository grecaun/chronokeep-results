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
	Code     int     `json:"code"`
	Response string  `json:"response"`
	Count    int     `json:"count"`
	Events   []Event `json:"events"`
}

// GetAccountEventsRequest Struct used for the request to Get Events based on account.
type GetAccountEventsRequest struct {
	Key     string `json:"key"`
	Account string `json:"account"`
}

// GetEventResponse Struct used for the response of a Get Event request.
type GetEventResponse struct {
	Code     int    `json:"code"`
	Response string `json:"response"`
	Event    Event  `json:"event"`
}

// GetEventRequest Struct used for the request to get a single Event.
type GetEventRequest struct {
	Key       string `json:"key"`
	EventSlug string `json:"eventSlug"`
}

// ModifyEventResponse Struct used for the response of a Add/Update Event request.
type ModifyEventResponse struct {
	Code     int    `json:"code"`
	Response string `json:"response"`
	Comment  string `json:"comment"`
	Event    Event  `json:"event"`
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
	Key  string `json:"key"`
	Slug string `json:"slug"`
}

/*
	Handlers related to an EventYear
*/

// GetEventYearResponse Struct used for the response of a Get EventYear request.
type GetEventYearResponse struct {
	Code      int       `json:"code"`
	Response  string    `json:"response"`
	Event     Event     `json:"event"`
	EventYear EventYear `json:"eventYear"`
}

// GetEventYearRequest Struct used for the request of an Event Year.
type GetEventYearRequest struct {
	Key  string `json:"key"`
	Slug string `json:"slug"`
	Year string `json:"year"`
}

// ModifyEventYearResponse Struct used for the response of an Add/Update EventYear Request.
type ModifyEventYearResponse struct {
	Code      int       `json:"code"`
	Response  string    `json:"response"`
	Event     Event     `json:"event"`
	EventYear EventYear `json:"eventYear"`
}

// AddEventYearRequest Struct used to add an Event Year.
type AddEventYearRequest struct {
	Key       string    `json:"key"`
	Slug      string    `json:"slug"`
	EventYear EventYear `json:"eventYear"`
}

// UpdateEventYearRequest Struct used to update an Event Year.
type UpdateEventYearRequest struct {
	Key       string    `json:"key"`
	EventYear EventYear `json:"eventYear"`
}

// DeleteEventYearRequest Struct used to delete an Event Year.
type DeleteEventYearRequest struct {
	Key  string `json:"key"`
	Slug string `json:"slug"`
	Year string `json:"year"`
}

/*
	Handlers realted to results.
*/

// GetResultsResponse Struct used for the response of a GetResults request.
type GetResultsResponse struct {
	Code     int      `json:"code"`
	Response string   `json:"response"`
	Count    int      `json:"count"`
	Results  []Result `json:"results"`
}

// GetResultsRequest Struct used for the request of Results for an EventYear
type GetResultsRequest struct {
	Key  string `json:"key"`
	Slug string `json:"slug"`
	Year string `json:"year"`
}

// ModifyResultsResponse Struct used for the response to an Add/Update/Delete Results request.
type ModifyResultsResponse struct {
	Code     int    `json:"code"`
	Response string `json:"response"`
	Count    int    `json:"count"`
}

// ModifyResultsRequest Struct used to add/update Results.
type ModifyResultsRequest struct {
	Key     string   `json:"key"`
	Slug    string   `json:"slug"`
	Year    string   `json:"year"`
	Results []Result `json:"results"`
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
	Code     int    `json:"code"`
	Response string `json:"response"`
	Key      Key    `json:"key"`
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
	Code     int     `json:"code"`
	Response string  `json:"response"`
	Account  Account `json:"account"`
	Keys     []Key   `json:"keys"`
	Events   []Event `json:"events"`
}

// GetAccountRequest Struct used to request information about a specific account.
type GetAccountRequest struct {
	AccountIdentifier string `json:"accountID"`
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

// ModifyAccountRequest Struct used to request to add/update an account.
type ModifyAccountRequest struct {
	Account Account `json:"account"`
}

// DeleteAccountRequest Struct used to delete an account.
type DeleteAccountRequest struct {
	Email string `json:"email"`
}
