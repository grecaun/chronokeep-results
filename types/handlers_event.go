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
	EventYears []EventYear `json:"eventYears"`
}

// ModifyEventResponse Struct used for the response of a Add/Update Event request.
type ModifyEventResponse struct {
	Event Event `json:"event"`
}

/*
	Requests
*/

// GetAccountEventsRequest Struct used for the request to Get Events based on account.
type GetAccountEventsRequest struct {
	Key   string `json:"key"`
	Email string `json:"email"`
}

// GetEventRequest Struct used for the request to get a single Event.
type GetEventRequest struct {
	Key       string `json:"key"`
	EventSlug string `json:"eventSlug"`
}

// AddEventRequest Struct used for the request to Add an Event.
type AddEventRequest struct {
	Key          string `json:"key"`
	AccountEmail string `json:"accountEmail"`
	Event        Event  `json:"event"`
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
