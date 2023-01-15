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
	Chip       string `json:"chip"`
	Anonymous  bool   `json:"anonymous"`
}

// Validate Ensures valid data in the struct.
func (p *Person) Validate(validate *validator.Validate) error {
	p.Gender = strings.ToUpper(p.Gender)
	if p.Gender != "M" && p.Gender != "F" && p.Gender != "O" && p.Gender != "U" && p.Gender != "NB" {
		return errors.New("invalid gender (M/F/NB/O/U)")
	}
	return validate.Struct(p)
}

func (one *Person) Equals(two *Person) bool {
	return one.Bib == two.Bib &&
		one.First == two.First &&
		one.Last == two.Last &&
		one.Age == two.Age &&
		one.Gender == two.Gender &&
		one.AgeGroup == two.AgeGroup &&
		one.Distance == two.Distance &&
		one.Chip == two.Chip &&
		one.Anonymous == two.Anonymous
}

func (p *Person) AnonyInt() int {
	if p.Anonymous {
		return 1
	}
	return 0
}
