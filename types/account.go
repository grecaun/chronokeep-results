package types

// Account is a structure holding information on accounts that have access
// to this module
type Account struct {
	Identifier int64  `json:"-"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Type       string `json:"type"`
}
