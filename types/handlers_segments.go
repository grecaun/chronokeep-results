package types

/*
	Responses
*/

type GetSegmentsResponse struct {
	Segments []Segment `json:"segments"`
}

type AddSegmentsResponse struct {
	Segments []Segment `json:"segments"`
}

type RemoveSegmentsResponse struct {
	Count int64 `json:"count"`
}

/*
	Requests
*/

type GetSegmentsRequest struct {
	Slug string  `json:"slug"`
	Year *string `json:"year"`
}

type AddSegmentsRequest struct {
	Slug     string    `json:"slug"`
	Year     string    `json:"year"`
	Segments []Segment `json:"segments"`
}

type RemoveSegmentsRequest struct {
	Slug string `json:"slug"`
	Year string `json:"year"`
}
