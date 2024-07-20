package types

import "github.com/go-playground/validator/v10"

type Segment struct {
	Identifier    int64   `json:"-"`
	Location      string  `json:"location" validate:"required"`
	DistanceName  string  `json:"distance_name" validate:"required"`
	Name          string  `json:"name" validate:"required"`
	DistanceValue float64 `json:"distance_value" validate:"required"`
	DistanceUnit  string  `json:"distance_unit" validate:"required"`
	GPS           string  `json:"gps"`
	MapLink       string  `json:"map_link"`
}

func (s *Segment) Validate(validate *validator.Validate) error {
	return validate.Struct(s)
}

func (s Segment) Equals(other Segment) bool {
	return s.Location == other.Location &&
		s.DistanceName == other.DistanceName &&
		s.Name == other.Name &&
		s.DistanceValue == other.DistanceValue &&
		s.DistanceUnit == other.DistanceUnit &&
		s.GPS == other.GPS &&
		s.MapLink == other.MapLink
}
