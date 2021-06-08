package types

// Person Describes a person.
type Person struct {
	Identifier int64  `json:"-"`
	Bib        string `json:"bib"`
	First      string `json:"first"`
	Last       string `json:"last"`
	Age        int    `json:"age"`
	Gender     string `json:"gender"`
	AgeGroup   string `json:"age_group"`
	Distance   string `json:"distance"`
}
