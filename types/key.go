package types

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

// Key outline for data stored about an PI key
// Account should be a unique value for the account that owns the Key.
// Example types are: read (readonly), delete (read, write, delete), write (read, write)
// Allowed hosts are the hosts the calls are allowed to come from. Default of empty string is all hosts.
type Key struct {
	AccountIdentifier int64     `json:"-"`
	Value             string    `json:"value"`
	Type              string    `json:"type" validate:"required"`
	AllowedHosts      string    `json:"allowedHosts"`
	ValidUntil        time.Time `json:"validUntil" validate:"datetime,required"`
}

func (k *Key) Equal(other *Key) bool {
	return k.AccountIdentifier == other.AccountIdentifier &&
		k.Value == other.Value &&
		k.Type == other.Type &&
		k.AllowedHosts == other.AllowedHosts &&
		k.ValidUntil.Equal(other.ValidUntil)
}

// Validate Ensures valid data in the structure.
func (k *Key) Validate(validate *validator.Validate) error {
	valid := false
	switch k.Type {
	case "read":
		valid = true
	case "write":
		valid = true
	case "delete":
		valid = true
	}
	if !valid {
		return errors.New("invalid key type specified")
	}
	// TODO: validation on the allowed hosts
	return validate.Struct(k)
}
