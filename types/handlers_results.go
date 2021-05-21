package types

/*
	Responses
*/

// GetResultsResponse Struct used for the response of a GetResults request.
type GetResultsResponse struct {
	Count     int       `json:"count"`
	Event     Event     `json:"event"`
	EventYear EventYear `json:"event_year"`
	Results   []Result  `json:"results"`
}

// AddResultsResponse Struct used for the response to an Add/Update/Delete Results request.
type AddResultsResponse struct {
	Count int `json:"count"`
}

/*
	Requests
*/

// GetResultsRequest Struct used for the request of Results for an EventYear. Also used for Delete.
type GetResultsRequest struct {
	Key  string `json:"key"`
	Slug string `json:"slug"`
	Year string `json:"year"`
}

// AddResultsRequest Struct used to add/update Results.
type AddResultsRequest struct {
	Key     string   `json:"key"`
	Slug    string   `json:"slug"`
	Year    string   `json:"year"`
	Results []Result `json:"results"`
}
