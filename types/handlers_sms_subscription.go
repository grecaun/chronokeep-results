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

// GetSmsSubscriptionsResponse Struct used for the response of a GetSmsSubscriptions request.
type GetSmsSubscriptionsResponse struct {
	Subscriptions []SmsSubscription `json:"subscriptions"`
}

/*
	Requests
*/

// GetSmsSubscriptionsRequest Struct used to get the list of subscriptions requested.
type GetSmsSubscriptionsRequest struct {
	Slug string  `json:"slug"`
	Year *string `json:"year"`
}

// AddSmsSubscriptionRequest Struct used for the request to add a phone number to be alerted when a specific bib/person is seen.
type AddSmsSubscriptionRequest struct {
	Slug  string  `json:"slug"`
	Year  *string `json:"year"`
	Bib   *string `json:"bib"`
	First *string `json:"first"`
	Last  *string `json:"last"`
	Phone string  `json:"phone"`
}

// RemoveSmsSubscriptionRequest Struct used for the request to remove a phone number from the subscribed list.
type RemoveSmsSubscriptionRequest struct {
	Slug  string  `json:"slug"`
	Year  *string `json:"year"`
	Phone string  `json:"phone"`
}

