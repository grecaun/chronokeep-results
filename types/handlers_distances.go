package types

/*
	Responses
*/

type GetDistancesResponse struct {
	Distances []Distance `json:"distances"`
}

type GetDistanceResponse struct {
	Distance *Distance `json:"distance"`
}

type DeleteDistancesResponse struct {
	Count int64 `json:"count"`
}

/*
	Requests
*/

type GetDistancesRequest struct {
	Slug     string  `json:"slug"`
	Year     *string `json:"year"`
	Distance *string `json:"distance"`
}

type AddDistancesRequest struct {
	Slug      string     `json:"slug"`
	Year      string     `json:"year"`
	Distances []Distance `json:"distances"`
}

type DeleteDistancesRequest struct {
	Slug string `json:"slug"`
	Year string `json:"year"`
}
