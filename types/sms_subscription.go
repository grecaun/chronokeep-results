package types

// SmsSubscription holds the information regarding a text subscription.
type SmsSubscription struct {
	Bib   string `json:"bib"`
	First string `json:"first"`
	Last  string `json:"last"`
	Phone string `json:"phone"`
}

func (s *SmsSubscription) Equals(o *SmsSubscription) bool {
	return s.Bib == o.Bib &&
		s.First == o.First &&
		s.Last == o.Last &&
		s.Phone == o.Phone
}
