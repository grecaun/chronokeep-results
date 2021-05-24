package types

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Result is a structure holding information about a specific time
// result for a specific event.
type Result struct {
	Bib              string `json:"bib" validate:"required"`
	First            string `json:"first" validate:"required"`
	Last             string `json:"last" validate:"required"`
	Age              int    `json:"age" validate:"gte=0,lte=130"`
	Gender           string `json:"gender"`
	AgeGroup         string `json:"age_group"`
	Distance         string `json:"distance"`
	Seconds          int    `json:"seconds" validate:"gte=0"`
	Milliseconds     int    `json:"milliseconds"`
	ChipSeconds      int    `json:"chip_seconds" validate:"gte=0"`
	ChipMilliseconds int    `json:"chip_milliseconds"`
	Segment          string `json:"segment"`
	Location         string `json:"location" validate:"required"`
	Occurence        int    `json:"occurence" validate:"gte=0"`
	Ranking          int    `json:"ranking"`
	AgeRanking       int    `json:"age_ranking"`
	GenderRanking    int    `json:"gender_ranking"`
	Finish           bool   `json:"finish"`
}

// Validate Ensures valid data in the struct.
func (r *Result) Validate(validate *validator.Validate) error {
	r.Gender = strings.ToUpper(r.Gender)
	if r.Gender != "M" && r.Gender != "F" && r.Gender != "O" && r.Gender != "U" {
		return errors.New("invalid gender (M/F/O/U)")
	}
	return validate.Struct(r)
}
