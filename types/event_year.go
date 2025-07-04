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
	RankingType     string    `json:"ranking_type"`
}

type EventYearVers1 struct {
	Year     string    `json:"year"`
	DateTime time.Time `json:"date_time"`
}

type RequestYear struct {
	Year        string `json:"year" validate:"required"`
	DateTime    string `json:"date_time"`
	Live        bool   `json:"live"`
	DaysAllowed int    `json:"days_allowed"`
	RankingType string `json:"ranking_type"`
}

type AllEventYear struct {
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Year        string    `json:"year"`
	DateTime    time.Time `json:"date_time" validate:"datetime"`
	Live        bool      `json:"live"`
	DaysAllowed int       `json:"days_allowed"`
	RankingType string    `json:"ranking_type"`
}

func (e *EventYear) ConvertToVers1() EventYearVers1 {
	return EventYearVers1{
		Year:     e.Year,
		DateTime: e.DateTime,
	}
}

func (e *EventYear) Equals(other *EventYear) bool {
	return e.Year == other.Year &&
		e.DateTime.Equal(other.DateTime) &&
		e.Live == other.Live &&
		e.DaysAllowed == other.DaysAllowed &&
		e.RankingType == other.RankingType
}

func (e *EventYear) EqualsAll(other *AllEventYear) bool {
	return e.Year == other.Year &&
		e.DateTime.Equal(other.DateTime) &&
		e.Live == other.Live &&
		e.DaysAllowed == other.DaysAllowed &&
		e.RankingType == other.RankingType
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
		Year:        e.Year,
		Live:        e.Live,
		DaysAllowed: e.DaysAllowed,
		RankingType: e.RankingType,
	}
	d, err := time.Parse(time.RFC3339, e.DateTime)
	if err == nil {
		out.DateTime = d
		return out
	}
	d, err = time.Parse("2006/01/02 15:04:05 -07:00", e.DateTime)
	if err == nil {
		out.DateTime = d
		return out
	}
	d, err = time.Parse("2006/01/02 15:04:05", e.DateTime)
	if err != nil {
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
	d, err = time.Parse("2006/01/02 15:04:05 -07:00", e.DateTime)
	if err == nil {
		return d
	}
	d, err = time.Parse("2006/01/02 15:04:05", e.DateTime)
	if err != nil {
		return d
	}
	return time.Now()
}
