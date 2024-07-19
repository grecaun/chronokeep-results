package postgres

import (
	"chronokeep/results/types"
	"context"
	"fmt"
	"time"
)

func (p *Postgres) AddSubscribedPhone(eventYearID int64, subscription types.SmsSubscription) error {
	db, err := p.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("unable to start transaction: %v", err)
	}
	_, err = tx.Exec(
		ctx,
		"INSERT INTO sms_subscriptions(event_year_id, bib, first, last, phone) VALUES ($1,$2,$3,$4,$5)",
		eventYearID,
		subscription.Bib,
		subscription.First,
		subscription.Last,
		subscription.Phone,
	)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("unable to add sms subscription: %v", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("error committing transaction: %v", err)
	}
	return nil
}

func (p *Postgres) RemoveSubscribedPhone(eventYearID int64, phone string) error {
	db, err := p.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("unable to start transaction: %v", err)
	}
	_, err = tx.Exec(
		ctx,
		"DELETE FROM sms_subscriptions WHERE event_year_id=$1 AND phone=$2;",
		eventYearID,
		phone,
	)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("unable to remove sms subscription: %v", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("error committing transaction: %v", err)
	}
	return nil
}

func (p *Postgres) GetSubscribedPhones(eventYearID int64) ([]types.SmsSubscription, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Query(
		ctx,
		"SELECT bib, first, last, phone FROM sms_subscriptions WHERE event_year_id=$1;",
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
