package types

import "github.com/go-playground/validator/v10"

// RegistrationPerson Describes a person.
type RegistrationPerson struct {
	Identifier int64  `json:"id"`
	Bib        string `json:"bib" validate:"required"`
	First      string `json:"first"`
	Last       string `json:"last"`
	Age        int    `json:"age" validate:"gte=0,lte=130"`
	Gender     string `json:"gender"`
	AgeGroup   string `json:"age_group"`
	Distance   string `json:"distance" validate:"required"`
	Chip       string `json:"chip"`
	Anonymous  bool   `json:"anonymous"`
	SMSEnabled bool   `json:"sms"`
}

// Validate Ensures valid data in the struct.
func (p *RegistrationPerson) Validate(validate *validator.Validate) error {
	return validate.Struct(p)
}

func (p *RegistrationPerson) AnonyInt() int {
	if p.Anonymous {
		return 1
	}
	return 0
}

func (p *RegistrationPerson) SMSInt() int {
	if p.SMSEnabled {
		return 1
	}
	return 0
}
