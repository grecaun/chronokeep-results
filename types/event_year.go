package types

// EventYear is a structure holding information about a specific year related to an
// Event
type EventYear struct {
	Identifier      int64  `json:"-"`
	EventIdentifier int64  `json:"-"`
	Year            string `json:"year"`
	Date            string `json:"date"`
	Time            string `json:"time"`
	Live            bool   `json:"live"`
}

func (this *EventYear) Equals(other *EventYear) bool {
	return this.EventIdentifier == other.EventIdentifier &&
		this.Year == other.Year &&
		this.Date == other.Date &&
		this.Time == other.Time &&
		this.Live == other.Live
}
