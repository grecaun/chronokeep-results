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

type GetDistancesResponse struct {
	Distances []Distance `json:"distances"`
}

type GetDistanceResponse struct {
	Distance *Distance `json:"distance"`
}

type DeleteDistancesResponse struct {
	Count int64 `json:"count"`
}

/*
	Requests
*/

type GetDistancesRequest struct {
	Slug     string  `json:"slug"`
	Year     *string `json:"year"`
	Distance *string `json:"distance"`
}

type AddDistancesRequest struct {
	Slug      string     `json:"slug"`
	Year      string     `json:"year"`
	Distances []Distance `json:"distances"`
}

type DeleteDistancesRequest struct {
	Slug string `json:"slug"`
	Year string `json:"year"`
}

