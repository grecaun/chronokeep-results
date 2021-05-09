package types

import "golang.org/x/crypto/bcrypt"

// Account is a structure holding information on accounts that have access
// to this module
type Account struct {
	Identifier int64  `json:"-"`
	Password   string `json:"-" validate:"required,min=8"`
	Name       string `json:"name" validate:"required"`
	Email      string `json:"email" validate:"email,required"`
	Type       string `json:"type" validate:"required"`
}

// Equals is used to check if the fields of an Account other than the identifier are identical.
func (a *Account) Equals(other *Account) bool {
	return a.Name == other.Name &&
		a.Email == other.Email &&
		a.Type == other.Type
}

// PasswordIsHashed is used to check if a password has been hashed for insert into the database.
func (a *Account) PasswordIsHashed() bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(""))
	switch err.(type) {
	case bcrypt.InvalidHashPrefixError:
		return false
	case nil:
		// CompareHashAndPassword only returns nil if the password matches, empty string should NEVER be a correct password.
		return false
	default:
		// Check for hash too short (i.e. starts with $ but not long enough to be a hash value)
		if err.Error() == "crypto/bcrypt: hashedSecret too short to be a bcrypted password" {
			return false
		}
	}
	return true
}
