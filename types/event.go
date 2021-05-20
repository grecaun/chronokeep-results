package types

import (
	"github.com/go-playground/validator/v10"
)

// Event is a structure holding the information regarding an event that can span
// multiple years
type Event struct {
	AccountIdentifier int64  `json:"-"`
	Identifier        int64  `json:"-"`
	Name              string `json:"name" validate:"required"`
	Slug              string `json:"slug" validate:"required"`
	Website           string `json:"website" validate:"url"`
	Image             string `json:"image" validate:"url"`
	ContactEmail      string `json:"contactEmail" validate:"email,required"`
	AccessRestricted  bool   `json:"accessRestricted"`
}

// Equals Returns true if all fields other than Identifier are equal.
func (e *Event) Equals(other *Event) bool {
	return e.AccountIdentifier == other.AccountIdentifier &&
		e.Name == other.Name &&
		e.Slug == other.Slug &&
		e.Website == other.Website &&
		e.Image == other.Website &&
		e.ContactEmail == other.ContactEmail &&
		e.AccessRestricted == other.AccessRestricted
}

// Validate Ensures valid information in the structure.
func (e *Event) Validate(validate *validator.Validate) error {
	return validate.Struct(e)
}
