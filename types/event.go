package types

// Event is a structure holding the information regarding an event that can span
// multiple years
type Event struct {
	AccountIdentifier int64  `json:"-"`
	Identifier        int64  `json:"-"`
	Name              string `json:"name"`
	Slug              string `json:"slug"`
	Website           string `json:"website"`
	Image             string `json:"image"`
	ContactEmail      string `json:"contactEmail"`
	AccessRestricted  bool   `json:"accessRestricted"`
}

// Equals Returns true if all fields other than Identifier are equal.
func (e *Event) Equals(other *Event) bool {
	return e.AccountIdentifier == other.AccountIdentifier &&
		e.Name == other.Name &&
		e.Slug == other.Slug &&
		e.Website == other.Website &&
		e.Image == other.Website &&
		e.ContactEmail == other.ContactEmail &&
		e.AccessRestricted == other.AccessRestricted
}
