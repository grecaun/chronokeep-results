package types

// Event is a structure holding the information regarding an event that can span
// multiple years
type Distance struct {
	Identifier    int64  `json:"-"`
	Name          string `json:"name"`
	Certification string `json:"certification"`
}
