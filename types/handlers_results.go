package types

/*
	Responses
*/

// GetResultsResponse Struct used for the response of a GetResults request.
type GetResultsResponse struct {
	Count     int                 `json:"count"`
	Event     Event               `json:"event"`
	EventYear EventYear           `json:"event_year"`
	Years     []EventYear         `json:"years"`
	Results   map[string][]Result `json:"results"`
}

// AddResultsResponse Struct used for the response to an Add/Update/Delete Results request.
type AddResultsResponse struct {
	Count int `json:"count"`
}

// GetBibResultsResponse Struct used for the response of a GetBibResults request.
type GetBibResultsResponse struct {
	Event     Event     `json:"event"`
	EventYear EventYear `json:"year"`
	Results   []Result  `json:"results"`
	Person    *Person   `json:"person"`
}

/*
	Requests
*/

// GetResultsRequest Struct used for the request of Results for an EventYear. Also used for Delete.
type GetResultsRequest struct {
	Slug string  `json:"slug"`
	Year *string `json:"year"`
}

// AddResultsRequest Struct used to add/update Results.
type AddResultsRequest struct {
	Slug    string   `json:"slug"`
	Year    string   `json:"year"`
	Results []Result `json:"results"`
}

// GetBibResultsRequest Struct used for the request of results for a person with a bib for an event year.
type GetBibResultsRequest struct {
	Bib  string `json:"bib"`
	Slug string `json:"slug"`
	Year string `json:"year"`
}
