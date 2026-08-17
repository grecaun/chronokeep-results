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

/*
	Responses
*/

// EventYearResponse Struct used for the response of a Get/Add/Update EventYear request.
type EventYearResponse struct {
	Event     Event     `json:"event"`
	EventYear EventYear `json:"event_year"`
}

// EventYearsResponse Struct used for the response of a Get Event Years request.
type EventYearsResponse struct {
	EventYears []EventYear `json:"years"`
}

// AllEventYearsResponse Struct used for the response of a Get All Event Years request.
type AllEventYearsResponse struct {
	EventYears []AllEventYear `json:"years"`
}

/*
	Requests
*/

// GetEventYearRequest Struct used for the request of an Event Year.
type GetEventYearRequest struct {
	Slug string `json:"slug"`
	Year string `json:"year"`
}

// ModifyEventYearRequest Struct used to add an Event Year.
type ModifyEventYearRequest struct {
	Slug      string      `json:"slug"`
	EventYear RequestYear `json:"event_year"`
}

// DeleteEventYearRequest Struct used to delete an Event Year.
type DeleteEventYearRequest struct {
	Slug string `json:"slug"`
	Year string `json:"year"`
}

