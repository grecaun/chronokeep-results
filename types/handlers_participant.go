package types

/*
	Responses
*/

// GetParticipantsResponse Struct used for the response of a Get Participants request.
type GetParticipantsResponse struct {
	Participants []Person `json:"participants"`
}

/*
	Requests
*/

// UpdateParticipantRequest Struct used to update a participant's information.
type UpdateParticipantRequest struct {
	Slug        string `json:"slug"`
	Year        string `json:"year"`
	Participant Person `json:"participant"`
}

type GetParticipantRequest struct {
	Slug       string `json:"slug"`
	Year       string `json:"year"`
	Identifier string `json:"person_id"`
}

// GetParticipantsRequest Struct used for the request to get participants for an event.
type GetParticipantsRequest struct {
	Slug string `json:"slug"`
	Year string `json:"year"`
}

// AddParticipantsRequest Struct used for the request to add/update participants for an event.
type AddParticipantsRequest struct {
	Slug         string   `json:"slug"`
	Year         string   `json:"year"`
	Participants []Person `json:"participants"`
}

// DeleteParticipantsRequest Struct used for the request to delete participants for an event.
type DeleteParticipantsRequest struct {
	Slug        string   `json:"slug"`
	Year        string   `json:"year"`
	Identifiers []string `json:"identifiers"`
}
