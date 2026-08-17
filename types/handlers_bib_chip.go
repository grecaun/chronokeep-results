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

// GetBibChipsResponse Struct used for the response of a GetParticipantsRequest
type GetBibChipsResponse struct {
	BibChips []BibChip `json:"bib_chips"`
}

/*
	Requests
*/

// AddBibChipsRequest Struct used to add BibChips.
type AddBibChipsRequest struct {
	Slug     string    `json:"slug"`
	Year     string    `json:"year"`
	BibChips []BibChip `json:"bib_chips"`
}

// GetBibChipsRequest Struct used to get/delete BibChips.
type GetBibChipsRequest struct {
	Slug string `json:"slug"`
	Year string `json:"year"`
}

