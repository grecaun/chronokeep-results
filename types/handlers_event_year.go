package types

/*
	Responses
*/

// EventYearResponse Struct used for the response of a Get/Add/Update EventYear request.
type EventYearResponse struct {
	Event     Event     `json:"event"`
	EventYear EventYear `json:"event_year"`
}

// EventYearsResponse Struct used for the response of a Get Event Years request.
type EventYearsResponse struct {
	EventYears []EventYear `json:"years"`
}

/*
	Requests
*/

// GetEventYearRequest Struct used for the request of an Event Year.
type GetEventYearRequest struct {
	Slug string `json:"slug"`
	Year string `json:"year"`
}

// ModifyEventYearRequest Struct used to add an Event Year.
type ModifyEventYearRequest struct {
	Slug      string      `json:"slug"`
	EventYear RequestYear `json:"event_year"`
}

// DeleteEventYearRequest Struct used to delete an Event Year.
type DeleteEventYearRequest struct {
	Slug string `json:"slug"`
	Year string `json:"year"`
}
