package types

// Account is a structure holding information on accounts that have access
// to this module
type Account struct {
	Identifier int64  `json:"-"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Type       string `json:"type"`
}

// Equals is used to check if the fields of an Account other than the identifier are identical.
func (a *Account) Equals(other *Account) bool {
	return a.Name == other.Name &&
		a.Email == other.Email &&
		a.Type == other.Type
}
