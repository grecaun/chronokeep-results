package types

/*
	Responses
*/

// GetResultsResponse Struct used for the response of a GetResults request.
type GetResultsResponse struct {
	Count        int                 `json:"count"`
	Event        Event               `json:"event"`
	EventYear    EventYear           `json:"event_year"`
	Years        []EventYear         `json:"years"`
	Results      map[string][]Result `json:"results"`
	Participants []ResultParticipant `json:"participants"`
	Distances    []Distance          `json:"distances"`
}

// GetResultsResponse Struct used for the response of a GetResults request.
type GetMultiResultsResponse struct {
	Event   Event                          `json:"event"`
	Results map[string]map[string][]Result `json:"results"`
}

type GetResultsResponseVers1 struct {
	Count   int                      `json:"count"`
	Event   EventVers1               `json:"event"`
	Years   []EventYearVers1         `json:"years"`
	Results map[string][]ResultVers1 `json:"results"`
}

// AddResultsResponse Struct used for the response to an Add/Update/Delete Results request.
type AddResultsResponse struct {
	Count int `json:"count"`
}

// GetBibResultsResponse Struct used for the response of a GetBibResults request.
type GetBibResultsResponse struct {
	Event          Event     `json:"event"`
	EventYear      EventYear `json:"year"`
	Results        []Result  `json:"results"`
	Person         *Person   `json:"person"`
	SingleDistance bool      `json:"single_distance"`
	Segments       []Segment `json:"segments"`
	Distance       *Distance `json:"distance"`
}

/*
	Requests
*/

// GetResultsRequest Struct used for the request of Results for an EventYear. Also used for Delete.
type GetResultsRequest struct {
	Slug     string  `json:"slug"`
	Year     *string `json:"year"`
	Distance *string `json:"distance"`
	Limit    *int    `json:"limit"`
	Page     *int    `json:"page"`
	Version  *int    `json:"version"`
}

// GetMultiResultsRequest Struct used for the request of results for many years for an Event.
type GetMultiResultsRequest struct {
	Slug  string   `json:"slug"`
	Years []string `json:"years"`
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
