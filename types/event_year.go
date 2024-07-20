package types

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

// EventYear is a structure holding information about a specific year related to an
// Event
type EventYear struct {
	Identifier      int64     `json:"-"`
	EventIdentifier int64     `json:"-"`
	Year            string    `json:"year"`
	DateTime        time.Time `json:"date_time" validate:"datetime"`
	Live            bool      `json:"live"`
	DaysAllowed     int       `json:"days_allowed"`
}

type RequestYear struct {
	Year        string `json:"year" validate:"required"`
	DateTime    string `json:"date_time"`
	Live        bool   `json:"live"`
	DaysAllowed int    `json:"days_allowed"`
}

func (e *EventYear) Equals(other *EventYear) bool {
	return e.EventIdentifier == other.EventIdentifier &&
		e.Year == other.Year &&
		e.DateTime.Equal(other.DateTime) &&
		e.Live == other.Live &&
		e.DaysAllowed == other.DaysAllowed
}

// Validate Ensures valid data in the structure.
func (e *EventYear) Validate(validate *validator.Validate) error {
	if !validYear(e.Year) {
		return errors.New("invalid year (only numbers and - allowed)")
	}
	return validate.Struct(e)
}

// Validate Ensures valid data in the structure.
func (e *RequestYear) Validate(validate *validator.Validate) error {
	if !validYear(e.Year) {
		return errors.New("invalid year (only numbers and - allowed)")
	}
	return validate.Struct(e)
}

// ToYear Turns a RequestYear into a year object.
func (e RequestYear) ToYear() EventYear {
	out := EventYear{
		Year: e.Year,
		Live: e.Live,
	}
	d, err := time.Parse(time.RFC3339, e.DateTime)
	if err == nil {
		out.DateTime = d
		return out
	}
	d, err = time.ParseInLocation("2006/01/02 15:04:05", e.DateTime, time.Local)
	if err == nil {
		out.DateTime = d
		return out
	}
	out.DateTime = time.Now()
	return out
}

// GetDateTime Gets a time.Time object from the RequestKey
func (e RequestYear) GetDateTime() time.Time {
	d, err := time.Parse(time.RFC3339, e.DateTime)
	if err == nil {
		return d
	}
	d, err = time.ParseInLocation("2006/01/02 15:04:05", e.DateTime, time.Local)
	if err == nil {
		return d
	}
	return time.Now()
}
