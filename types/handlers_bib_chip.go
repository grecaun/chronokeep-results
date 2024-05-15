package types

/*
	Responses
*/

// GetBibChipsResponse Struct used for the response of a GetParticipantsRequest
type GetBibChipsResponse struct {
	BibChips []BibChip `json:"bib_chips"`
}

/*
	Requests
*/

// AddBibChipsRequest Struct used to add BibChips.
type AddBibChipsRequest struct {
	Slug     string    `json:"slug"`
	Year     string    `json:"year"`
	BibChips []BibChip `json:"bib_chips"`
}

// GetBibChipsRequest Struct used to get/delete BibChips.
type GetBibChipsRequest struct {
	Slug string `json:"slug"`
	Year string `json:"year"`
}
