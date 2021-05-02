package types

type CallRecord struct {
	AccountIdentifier int64 `json:"-"`
	DateTime          int64 `json:"dateTime"`
	Count             int   `json:"count"`
}
