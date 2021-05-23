package types

// GenderalRequest A generalized request struct when only a key is required for the call.
type GeneralRequest struct {
	Key string `json:"key"`
}
