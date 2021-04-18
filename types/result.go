package types

// Result is a structure holding information about a specific time
// result for a specific event.
type Result struct {
	Bib          string `json:"bib"`
	First        string `json:"first"`
	Last         string `json:"last"`
	Age          int    `json:"age"`
	Gender       string `json:"gender"`
	AgeGroup     string `json:"ageGroup"`
	Distance     string `json:"distance"`
	Seconds      int    `json:"seconds"`
	Milliseconds int    `json:"milliseconds"`
	Location     string `json:"location"`
	Place        int    `json:"place"`
	AgePlace     int    `json:"agePlace"`
	GenderPlace  int    `json:"genderPlace"`
	Finish       bool   `json:"finish"`
}
