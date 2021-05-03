package types

import "time"

// Key outline for data stored about an PI key
// Account should be a unique value for the account that owns the Key.
// Example types are: default (readonly), delete (read, write, delete), write (read, write)
// Allowed hosts are the hosts the calls are allowed to come from. Default of empty string is all hosts.
type Key struct {
	AccountIdentifier int64     `json:"-"`
	Value             string    `json:"value"`
	Type              string    `json:"type"`
	AllowedHosts      string    `json:"allowedHosts"`
	ValidUntil        time.Time `json:"validUntil"`
}

func (k *Key) Equal(other *Key) bool {
	return k.AccountIdentifier == other.AccountIdentifier &&
		k.Value == other.Value &&
		k.Type == other.Type &&
		k.AllowedHosts == other.AllowedHosts &&
		k.ValidUntil.Equal(other.ValidUntil)
}
