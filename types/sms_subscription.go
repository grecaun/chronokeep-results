package types

// TextSubscription holds the information regarding a text subscription.
type TextSubscription struct {
	Bib   string `json:"bib"`
	First string `json:"first"`
	Last  string `json:"last"`
	Phone string `json:"phone"`
}
