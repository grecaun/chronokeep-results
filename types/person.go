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
	"github.com/go-playground/validator/v10"
)

// Person Describes a person.
type Person struct {
	Identifier  int64  `json:"-"`
	AlternateId string `json:"id"`
	Bib         string `json:"bib" validate:"required"`
	First       string `json:"first"`
	Last        string `json:"last"`
	Age         int    `json:"age" validate:"gte=0,lte=130"`
	Gender      string `json:"gender"`
	AgeGroup    string `json:"age_group"`
	Distance    string `json:"distance" validate:"required"`
	Anonymous   bool   `json:"anonymous"`
}

// Validate Ensures valid data in the struct.
func (p *Person) Validate(validate *validator.Validate) error {
	return validate.Struct(p)
}

func (one *Person) Equals(two *Person) bool {
	return one.Bib == two.Bib &&
		one.First == two.First &&
		one.Last == two.Last &&
		one.Age == two.Age &&
		one.Gender == two.Gender &&
		one.AgeGroup == two.AgeGroup &&
		one.Distance == two.Distance &&
		one.Anonymous == two.Anonymous &&
		one.AlternateId == two.AlternateId
}

func (p *Person) AnonyInt() int {
	if p.Anonymous {
		return 1
	}
	return 0
}

