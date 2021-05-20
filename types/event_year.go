package types

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// EventYear is a structure holding information about a specific year related to an
// Event
type EventYear struct {
	Identifier      int64     `json:"-"`
	EventIdentifier int64     `json:"-"`
	Year            string    `json:"year"`
	DateTime        time.Time `json:"dateTime" validate:"datetime"`
	Live            bool      `json:"live"`
}

func (e *EventYear) Equals(other *EventYear) bool {
	return e.EventIdentifier == other.EventIdentifier &&
		e.Year == other.Year &&
		e.DateTime.Equal(other.DateTime) &&
		e.Live == other.Live
}

// Validate Ensures valid data ion the structure.
func (e *EventYear) Validate(validate *validator.Validate) error {
	return validate.Struct(e)
}
