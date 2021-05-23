package types

/*
	Responses
*/

// GetEventsResponse Struct used for the response of a Get Events request.
type GetEventsResponse struct {
	Events []Event `json:"events"`
}

// GetEventResponse Struct used for the response of a Get Event request.
type GetEventResponse struct {
	Event      Event       `json:"event"`
	EventYears []EventYear `json:"event_years"`
	Year       *EventYear  `json:"year"`
	Results    []Result    `json:"results"`
}

// ModifyEventResponse Struct used for the response of a Add/Update Event request.
type ModifyEventResponse struct {
	Event Event `json:"event"`
}

/*
	Requests
*/

// GetEventRequest Struct used for the request to get a single Event.
type GetEventRequest struct {
	Key  string `json:"key"`
	Slug string `json:"slug"`
}

// GetEventsRequest Struct used for the request for multiple events.  If an account email is specified it pulls all events
// associated with that account, otherwise it pulls all non-restricted events. Used only for JWT restricted endpoints.
type GetREventsRequest struct {
	Email *string `json:"account_email,omitempty"`
}

// AddEventRequest Struct used for the request to Add an Event.
type AddEventRequest struct {
	Key   string  `json:"key"`
	Email *string `json:"account_email,omitempty"`
	Event Event   `json:"event"`
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
