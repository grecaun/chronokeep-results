package types

import (
	"github.com/go-playground/validator/v10"
)

// Result is a structure holding information about a specific time
// result for a specific event.
type Result struct {
	Bib           string `json:"bib" validate:""`
	First         string `json:"first" validate:""`
	Last          string `json:"last" validate:""`
	Age           int    `json:"age" validate:"gte=0,lte=130"`
	Gender        string `json:"gender" validate:"required"`
	AgeGroup      string `json:"ageGroup"`
	Distance      string `json:"distance"`
	Seconds       int    `json:"seconds" validate:"gte=1"`
	Milliseconds  int    `json:"milliseconds"`
	Segment       string `json:"segment"`
	Location      string `json:"location" validate:"required"`
	Occurence     int    `json:"occurence" validate:"gte=1"`
	Ranking       int    `json:"ranking"`
	AgeRanking    int    `json:"age_ranking"`
	GenderRanking int    `json:"gender_ranking"`
	Finish        bool   `json:"finish"`
}

// Validate Ensures valid data in the struct.
func (r *Result) Validate(validate *validator.Validate) error {
	return validate.Struct(r)
}
