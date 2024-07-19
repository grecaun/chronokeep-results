package sqlite

import "chronokeep/results/types"

func (s *SQLite) AddSubscribedPhone(eventYearID int64, subscription types.TextSubscription) error {
	return nil
}

func (s *SQLite) RemoveSubscribedPhone(eventYearID int64, phone string) error {
	return nil
}

func (s *SQLite) GetSubscribedPhones(eventYearID int64) ([]types.TextSubscription, error) {
	return nil, nil
}
