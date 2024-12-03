package types

import "github.com/go-playground/validator/v10"

// Event is a structure holding the information regarding an event that can span
// multiple years
type Distance struct {
	Identifier    int64  `json:"-"`
	Name          string `json:"name" validate:"required"`
	Certification string `json:"certification" validate:"required"`
}

func (d *Distance) Validate(validate *validator.Validate) error {
	return validate.Struct(d)
}

func (d Distance) Equals(other Distance) bool {
	return d.Name == other.Name &&
		d.Certification == other.Certification
}
