package types

/*
	Responses
*/

// GetEventYearResponse Struct used for the response of a Get EventYear request.
type GetEventYearResponse struct {
	Code      int       `json:"code"`
	Response  string    `json:"response"`
	Event     Event     `json:"event"`
	EventYear EventYear `json:"eventYear"`
}

// ModifyEventYearResponse Struct used for the response of an Add/Update EventYear Request.
type ModifyEventYearResponse struct {
	Code      int       `json:"code"`
	Response  string    `json:"response"`
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
