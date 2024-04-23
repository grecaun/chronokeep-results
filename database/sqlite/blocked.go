package sqlite

// AddBlockedPhone adds a phone number to the blocked phone numbers list
func (s *SQLite) AddBlockedPhone(phone string) error {
	return nil
}

// AddBlockedPhones adds one or more phone numbers to the blocked phone numbers list
func (s *SQLite) AddBlockedPhones(phones []string) error {
	return nil
}

// GetBlockedPhones gets the blocked phone numbers list
func (s *SQLite) GetBlockedPhones() ([]string, error) {
	return nil, nil
}

// UnblockPhone removes a phone number from the blocked phone numbers list
func (s *SQLite) UnblockPhone(phone string) error {
	return nil
}

// AddBlockedEmail adds an email address to the blocked emails list
func (s *SQLite) AddBlockedEmail(email string) error {
	return nil
}

// AddBlockedEmails adds one or more email addresses to the blocked emails list
func (s *SQLite) AddBlockedEmails(emails []string) error {
	return nil
}

// GetBlockedEmails gets the blocked phone emails list
func (s *SQLite) GetBlockedEmails() ([]string, error) {
	return nil, nil
}

// UnblockEmail removes an email address from the blocked blocked emails list
func (s *SQLite) UnblockEmail(email string) error {
	return nil
}
