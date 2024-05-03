package types

type MultiGet struct {
	Account       *Account
	Event         *Event
	EventYear     *EventYear
	DistanceCount *int
}

type MultiKey struct {
	Key     *Key
	Account *Account
}
