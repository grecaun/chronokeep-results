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

import "github.com/go-playground/validator/v10"

// Event is a structure holding the information regarding an event that can span
// multiple years
type Distance struct {
	Identifier    int64  `json:"-"`
	Name          string `json:"name" validate:"required"`
	Certification string `json:"certification" validate:"required"`
}

func (d *Distance) Validate(validate *validator.Validate) error {
	return validate.Struct(d)
}

func (d Distance) Equals(other Distance) bool {
	return d.Name == other.Name &&
		d.Certification == other.Certification
}

