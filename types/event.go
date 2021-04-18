package types

// Event is a structure holding the information regarding an event that can span
// multiple years
type Event struct {
	Identifier       string  `json:"id"`
	Name             string  `json:"name"`
	Slug             string  `json:"slug"`
	Website          string  `json:"website"`
	Image            string  `json:"image"`
	Owner            string  `json:"owner"`
	AccessRestricted bool    `json:"accessRestricted"`
	Years            []Event `json:"years"`
}
