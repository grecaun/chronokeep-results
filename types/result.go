package types

// Result is a structure holding information about a specific time
// result for a specific event.
type Result struct {
	Bib           string `json:"bib"`
	First         string `json:"first"`
	Last          string `json:"last"`
	Age           int    `json:"age"`
	Gender        string `json:"gender"`
	AgeGroup      string `json:"ageGroup"`
	Distance      string `json:"distance"`
	Seconds       int    `json:"seconds"`
	Milliseconds  int    `json:"milliseconds"`
	Segment       string `json:"segment"`
	Location      string `json:"location"`
	Occurence     int    `json:"occurence"`
	Ranking       int    `json:"ranking"`
	AgeRanking    int    `json:"ageRanking"`
	GenderRanking int    `json:"genderRanking"`
	Finish        bool   `json:"finish"`
}
