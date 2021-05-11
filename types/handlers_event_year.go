package types

/*
	Responses
*/

// EventYearResponse Struct used for the response of a Get/Add/Update EventYear request.
type EventYearResponse struct {
	Event     Event     `json:"event"`
	EventYear EventYear `json:"eventYear"`
}

/*
	Requests
*/

// GetEventYearRequest Struct used for the request of an Event Year.
type GetEventYearRequest struct {
	Key  string `json:"key"`
	Slug string `json:"slug"`
	Year string `json:"year"`
}

// ModifyEventYearRequest Struct used to add an Event Year.
type ModifyEventYearRequest struct {
	Key       string    `json:"key"`
	Slug      string    `json:"slug"`
	EventYear EventYear `json:"eventYear"`
}

// DeleteEventYearRequest Struct used to delete an Event Year.
type DeleteEventYearRequest struct {
	Key  string `json:"key"`
	Slug string `json:"slug"`
	Year string `json:"year"`
}
