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
