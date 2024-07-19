package sqlite

import (
	"chronokeep/results/types"
	"context"
	"fmt"
	"time"
)

func (s *SQLite) AddSubscribedPhone(eventYearID int64, subscription types.SmsSubscription) error {
	db, err := s.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("unable to start transaction: %v", err)
	}
	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO sms_subscriptions(event_year_id, bib, first, last, phone) VALUES (?,?,?,?,?)",
		eventYearID,
		subscription.Bib,
		subscription.First,
		subscription.Last,
		subscription.Phone,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("unable to add sms subscription: %v", err)
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error committing transaction: %v", err)
	}
	return nil
}

func (s *SQLite) RemoveSubscribedPhone(eventYearID int64, phone string) error {
	db, err := s.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("unable to start transaction: %v", err)
	}
	_, err = tx.ExecContext(
		ctx,
		"DELETE FROM sms_subscriptions WHERE event_year_id=? AND phone=?;",
		eventYearID,
		phone,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("unable to remove sms subscription: %v", err)
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error committing transaction: %v", err)
	}
	return nil
}

func (s *SQLite) GetSubscribedPhones(eventYearID int64) ([]types.SmsSubscription, error) {
	db, err := s.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT bib, first, last, phone FROM sms_subscriptions WHERE event_year_id=?;",
		eventYearID,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving sms subscriptions: %v", err)
	}
	defer res.Close()
	var outSubs []types.SmsSubscription
	for res.Next() {
		var sub types.SmsSubscription
		err := res.Scan(
			&sub.Bib,
			&sub.First,
			&sub.Last,
			&sub.Phone,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting subscription: %v", err)
		}
		outSubs = append(outSubs, sub)
	}
	return outSubs, nil
}
