package types

// Account is a structure holding information on accounts that have access
// to this module
type Account struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Identifier string `json:"id"`
	Type       string `json:"type"`
}
