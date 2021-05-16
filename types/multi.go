package types

type MultiGet struct {
	Account   *Account
	Event     *Event
	EventYear *EventYear
}

type MultiKey struct {
	Key     *Key
	Account *Account
}
