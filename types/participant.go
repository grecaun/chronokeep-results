/*
Chronokeep Desktop - Race Scoring Software
Copyright (C) 2026 James Sentinella

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

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
	UpdatedAt   int64  `json:"updated_at"`
}

// ResultParticipant Describes the information we want publicly available.
type ResultParticipant struct {
	Bib      string `json:"bib"`
	First    string `json:"first"`
	Last     string `json:"last"`
	AgeGroup string `json:"age_group"`
	Gender   string `json:"gender"`
	Distance string `json:"distance"`
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
		one.Apparel == two.Apparel &&
		one.UpdatedAt == two.UpdatedAt
}

func (one *Participant) TestEquals(two *Participant) bool {
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

func (p *Participant) ToResultParticipant() ResultParticipant {
	return ResultParticipant{
		Bib:      p.Bib,
		First:    p.First,
		Last:     p.Last,
		AgeGroup: p.AgeGroup,
		Gender:   p.Gender,
		Distance: p.Distance,
	}
}

