package types

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

// Account is a structure holding information on accounts that have access
// to this module
type Account struct {
	Identifier int64  `json:"-"`
	Unique     string `json:"identifier"`
	Name       string `json:"name" validate:"required"`
	Type       string `json:"type" validate:"required"`
	Locked     bool   `json:"locked"`
	ResultsAPI bool   `json:"results_api"`
	RemoteAPI  bool   `json:"remote_api"`
}

// Equals is used to check if the fields of an Account other than the identifier are identical.
func (a *Account) Equals(other *Account) bool {
	if other == nil {
		return false
	}
	return a.Unique == other.Unique &&
		a.Type == other.Type
}

// Validate Ensures that the struct is viable for entry into the database.
func (a *Account) Validate(validate *validator.Validate) error {
	valid := false
	switch a.Type {
	case "admin":
		valid = true
	case "free":
		valid = true
	case "paid":
		valid = true
	}
	if !valid {
		return errors.New("invalid account type specified")
	}
	return validate.Struct(a)
}
