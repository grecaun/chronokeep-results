package postgres

import "chronokeep/results/types"

func (p *Postgres) AddSubscribedPhone(eventYearID int64, subscription types.TextSubscription) error {
	return nil
}

func (p *Postgres) RemoveSubscribedPhone(eventYearID int64, phone string) error {
	return nil
}

func (p *Postgres) GetSubscribedPhones(eventYearID int64) ([]types.TextSubscription, error) {
	return nil, nil
}
