package types

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Person Describes a person.
type Person struct {
	Identifier int64  `json:"-"`
	Bib        string `json:"bib" validate:"required"`
	First      string `json:"first" validate:"required"`
	Last       string `json:"last" validate:"required"`
	Age        int    `json:"age" validate:"gte=0,lte=130"`
	Gender     string `json:"gender"`
	AgeGroup   string `json:"age_group"`
	Distance   string `json:"distance" validate:"required"`
}

// Validate Ensures valid data in the struct.
func (p *Person) Validate(validate *validator.Validate) error {
	p.Gender = strings.ToUpper(p.Gender)
	if p.Gender != "M" && p.Gender != "F" && p.Gender != "O" && p.Gender != "U" && p.Gender != "NB" {
		return errors.New("invalid gender (M/F/NB/O/U)")
	}
	return validate.Struct(p)
}
