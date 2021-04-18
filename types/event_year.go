package types

// EventYear is a structure holding information about a specific year related to an
// Event
type EventYear struct {
	Identifier string `json:"id"`
	Year       string `json:"year"`
	Date       string `json:"date"`
	Time       string `json:"time"`
	Live       bool   `json:"live"`
}
