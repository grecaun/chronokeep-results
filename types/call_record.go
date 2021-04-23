package types

type CallRecord struct {
	AccountIdentifier int    `json:"-"`
	DateTime          string `json:"dateTime"`
	Count             int    `json:"count"`
}
