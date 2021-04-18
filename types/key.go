package types

import "time"

// Key outline for data stored about an PI key
// Account should be a unique value for the account that owns the Key.
// Example types are: default, readonly, nodelete
// Allowed hosts are the hosts the calls are allowed to come from. Default of empty string is all hosts.
type Key struct {
	Key          string     `json:"key"`
	Type         string     `json:"type"`
	AllowedHosts string     `json:"allowedHosts"`
	Email        string     `json:"email"`
	ValidUntil   *time.Time `json:"validUntil"`
}
