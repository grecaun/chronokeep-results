package types

import (
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// Event is a structure holding the information regarding an event that can span
// multiple years
type Event struct {
	AccountIdentifier int64      `json:"-"`
	Identifier        int64      `json:"-"`
	Name              string     `json:"name" validate:"required"`
	CertificateName   string     `json:"cert_name"`
	Slug              string     `json:"slug" validate:"required"`
	Website           string     `json:"website"`
	Image             string     `json:"image"`
	ContactEmail      string     `json:"contact_email" validate:"email"`
	AccessRestricted  bool       `json:"access_restricted"`
	Type              string     `json:"type"`
	RecentTime        *time.Time `json:"recent_time"`
}

// Equals Returns true if all fields other than the Identifier fields are equal.
func (e *Event) Equals(other *Event) bool {
	return e.Name == other.Name &&
		e.CertificateName == other.CertificateName &&
		e.Slug == other.Slug &&
		e.Website == other.Website &&
		e.Image == other.Image &&
		e.ContactEmail == other.ContactEmail &&
		e.AccessRestricted == other.AccessRestricted &&
		e.Type == other.Type
}

// Validate Ensures valid information in the structure.
func (e *Event) Validate(validate *validator.Validate) error {
	e.Slug = strings.ToLower(e.Slug)
	if !validSlug(e.Slug) {
		return errors.New("invalid slug (only letters, numbers, and - character allowed)")
	}
	if !validEventName(e.Name) {
		return errors.New("invalid event name (only letters, numbers, ', /, and spaces allowed)")
	}
	valid := false
	switch e.Type {
	case "distance":
		valid = true
	case "time":
		valid = true
	case "backyardultra":
		valid = true
	}
	if !valid {
		return errors.New("invalid event type specified")
	}
	err := validate.Var(e.Website, "url")
	if e.Website != "" && err != nil {
		return errors.New("invalid website url")
	}
	err = validate.Var(e.Image, "url")
	if e.Image != "" && err != nil {
		return errors.New("invalid image url")
	}
	return validate.Struct(e)
}
