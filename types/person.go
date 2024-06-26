package types

import (
	"github.com/go-playground/validator/v10"
)

// Person Describes a person.
type Person struct {
	Identifier  int64  `json:"-"`
	AlternateId string `json:"id"`
	Bib         string `json:"bib" validate:"required"`
	First       string `json:"first"`
	Last        string `json:"last"`
	Age         int    `json:"age" validate:"gte=0,lte=130"`
	Gender      string `json:"gender"`
	AgeGroup    string `json:"age_group"`
	Distance    string `json:"distance" validate:"required"`
	Anonymous   bool   `json:"anonymous"`
}

// Validate Ensures valid data in the struct.
func (p *Person) Validate(validate *validator.Validate) error {
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
		one.Anonymous == two.Anonymous &&
		one.AlternateId == two.AlternateId
}

func (p *Person) AnonyInt() int {
	if p.Anonymous {
		return 1
	}
	return 0
}
