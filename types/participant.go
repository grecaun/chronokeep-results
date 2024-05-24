package types

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

// Participant Describes a person.
type Participant struct {
	Identifier  int64  `json:"-"`
	AlternateId string `json:"id"`
	Bib         string `json:"bib"`
	First       string `json:"first"`
	Last        string `json:"last"`
	Birthdate   string `json:"birthdate"`
	Gender      string `json:"gender"`
	AgeGroup    string `json:"age_group"`
	Distance    string `json:"distance" validate:"required"`
	Anonymous   bool   `json:"anonymous"`
	SMSEnabled  bool   `json:"sms_enabled"`
	Mobile      string `json:"mobile"`
	Apparel     string `json:"apparel"`
}

// Validate Ensures valid data in the struct.
func (p *Participant) Validate(validate *validator.Validate) error {
	us_layout := "1/2/2006"
	iso_layout := "2006/1/2"
	t, err := time.Parse(us_layout, p.Birthdate)
	if err != nil {
		t, err := time.Parse(iso_layout, p.Birthdate)
		if err != nil || t.After(time.Now()) {
			return fmt.Errorf("invalid birthdate")
		}
	}
	if t.After(time.Now()) {
		return fmt.Errorf("invalid birthdate")
	}
	return validate.Struct(p)
}

func (one *Participant) Equals(two *Participant) bool {
	return one.Bib == two.Bib &&
		one.First == two.First &&
		one.Last == two.Last &&
		one.Birthdate == two.Birthdate &&
		one.Gender == two.Gender &&
		one.AgeGroup == two.AgeGroup &&
		one.Distance == two.Distance &&
		one.Anonymous == two.Anonymous &&
		one.SMSEnabled == two.SMSEnabled &&
		one.Mobile == two.Mobile &&
		one.Apparel == two.Apparel
}

func (p *Participant) AnonyInt() int {
	if p.Anonymous {
		return 1
	}
	return 0
}

func (p *Participant) SMSInt() int {
	if p.SMSEnabled {
		return 1
	}
	return 0
}
