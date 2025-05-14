package types

/*
	Responses
*/

// GetParticipantsResponse Struct used for the response of a Get Participants request.
type GetParticipantsResponse struct {
	Event        Event         `json:"event"`
	Year         EventYear     `json:"year"`
	Participants []Participant `json:"participants"`
}

// UpdateParticipantResponse Struct used for the response of an Update Participants request.
type UpdateParticipantResponse struct {
	Participant Participant   `json:"participant"`
	Updated     []Participant `json:"updated_participants"`
}

// UpdateParticipantsResponse Struct used for the response of an Update Participants request.
type UpdateParticipantsResponse struct {
	Participants []Participant `json:"participants"`
	Updated      []Participant `json:"updated_participants"`
}

// AddResultsResponse Struct used for the response to an Add/Update/Delete Participants request.
type AddParticipantsResponse struct {
	Count   int           `json:"count"`
	Updated []Participant `json:"updated_participants"`
}

/*
	Requests
*/

// UpdateParticipantRequest Struct used to update a participant's information.
type UpdateParticipantRequest struct {
	Slug         string      `json:"slug"`
	Year         string      `json:"year"`
	Participant  Participant `json:"participant"`
	UpdatedAfter *int64      `json:"updated_after"`
}

type GetParticipantRequest struct {
	Slug       string `json:"slug"`
	Year       string `json:"year"`
	Identifier string `json:"person_id"`
}

// GetParticipantsRequest Struct used for the request to get participants for an event.
type GetParticipantsRequest struct {
	Slug         string  `json:"slug"`
	Year         *string `json:"year"`
	Limit        *int    `json:"limit"`
	Page         *int    `json:"page"`
	UpdatedAfter *int64  `json:"updated_after"`
}

// AddParticipantsRequest Struct used for the request to add/update participants for an event.
type AddParticipantsRequest struct {
	Slug         string        `json:"slug"`
	Year         string        `json:"year"`
	Participants []Participant `json:"participants"`
	UpdatedAfter *int64        `json:"updated_after"`
}

// AddParticipantRequest Struct used for the request to add/update participants for an event.
type AddParticipantRequest struct {
	Slug         string      `json:"slug"`
	Year         string      `json:"year"`
	Participant  Participant `json:"participant"`
	UpdatedAfter *int64      `json:"updated_after"`
}

// DeleteParticipantsRequest Struct used for the request to delete participants for an event.
type DeleteParticipantsRequest struct {
	Slug        string   `json:"slug"`
	Year        string   `json:"year"`
	Identifiers []string `json:"identifiers"`
}
