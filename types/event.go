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
func (this *Event) Equals(other *Event) bool {
	return this.AccountIdentifier == other.AccountIdentifier &&
		this.Name == other.Name &&
		this.Slug == other.Slug &&
		this.Website == other.Website &&
		this.Image == other.Website &&
		this.ContactEmail == other.ContactEmail &&
		this.AccessRestricted == other.AccessRestricted
}
