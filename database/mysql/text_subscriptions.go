package mysql

import "chronokeep/results/types"

func (m *MySQL) AddSubscribedPhone(eventYearID int64, subscription types.TextSubscription) error {
	return nil
}

func (m *MySQL) RemoveSubscribedPhone(eventYearID int64, phone string) error {
	return nil
}

func (m *MySQL) GetSubscribedPhones(eventYearID int64) ([]types.TextSubscription, error) {
	return nil, nil
}
