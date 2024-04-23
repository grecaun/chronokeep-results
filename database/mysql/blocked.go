package mysql

// AddBlockedPhone adds a phone number to the blocked phone numbers list
func (m *MySQL) AddBlockedPhone(phone string) error {
	return nil
}

// AddBlockedPhones adds one or more phone numbers to the blocked phone numbers list
func (m *MySQL) AddBlockedPhones(phones []string) error {
	return nil
}

// GetBlockedPhones gets the blocked phone numbers list
func (m *MySQL) GetBlockedPhones() ([]string, error) {
	return nil, nil
}

// UnblockPhone removes a phone number from the blocked phone numbers list
func (m *MySQL) UnblockPhone(phone string) error {
	return nil
}

// AddBlockedEmail adds an email address to the blocked emails list
func (m *MySQL) AddBlockedEmail(email string) error {
	return nil
}

// AddBlockedEmails adds one or more email addresses to the blocked emails list
func (m *MySQL) AddBlockedEmails(emails []string) error {
	return nil
}

// GetBlockedEmails gets the blocked phone emails list
func (m *MySQL) GetBlockedEmails() ([]string, error) {
	return nil, nil
}

// UnblockEmail removes an email address from the blocked blocked emails list
func (m *MySQL) UnblockEmail(email string) error {
	return nil
}
