package types

/*
	Responses
*/

type GetBannedPhonesResponse struct {
	Phones []string `json:"phones"`
}

type GetBannedEmailsResponse struct {
	Emails []string `json:"emails"`
}

/*
	Requests
*/

type ModifyBannedPhoneRequest struct {
	Phone string `json:"phone" validate:"numeric,required"`
}

type ModifyBannedEmailRequest struct {
	Email string `json:"email" validate:"email,required"`
}
