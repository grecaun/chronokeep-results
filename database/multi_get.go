package database

import (
	"chronokeep/results/types"
	"context"
	"fmt"
	"time"
)

// GetAccountAndEvent Gets an event and its corresponding account.
func GetAccountAndEvent(slug string) (*MultiGet, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT "+
			"account_id, account_name, account_email, account_type, account_locked, "+
			"event_id, event_name, slug, website, image, contact_email, access_restricted "+
			"FROM account NATURAL JOIN event WHERE account_deleted=FALSE AND event_deleted=FALSE and slug=?",
		slug,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting account and event from database: %v", err)
	}
	if res.Next() {
		outVal := MultiGet{
			Account: &types.Account{},
			Event:   &types.Event{},
		}
		err := res.Scan(
			&outVal.Account.Identifier,
			&outVal.Account.Name,
			&outVal.Account.Email,
			&outVal.Account.Type,
			&outVal.Account.Locked,
			&outVal.Event.Identifier,
			&outVal.Event.Name,
			&outVal.Event.Slug,
			&outVal.Event.Website,
			&outVal.Event.Image,
			&outVal.Event.ContactEmail,
			&outVal.Event.AccessRestricted,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting values for account and event: %v", err)
		}
		outVal.Event.AccountIdentifier = outVal.Account.Identifier
		return &outVal, nil
	}
	return nil, nil
}

// GetAccountEventAndYear Gets an eventyear and its corresponding event and account.
func GetAccountEventAndYear(slug, year string) (*MultiGet, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT "+
			"account_id, account_name, account_email, account_type, account_locked, "+
			"event_id, event_name, slug, website, image, contact_email, access_restricted, "+
			"event_year_id, year, date_time, live "+
			"FROM account NATURAL JOIN event NATURAL JOIN event_year WHERE account_deleted=FALSE AND event_deleted=FALSE AND year_deleted=FALSE AND slug=? AND year=?",
		slug,
		year,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting account and event from database: %v", err)
	}
	if res.Next() {
		outVal := MultiGet{
			Account:   &types.Account{},
			Event:     &types.Event{},
			EventYear: &types.EventYear{},
		}
		err := res.Scan(
			&outVal.Account.Identifier,
			&outVal.Account.Name,
			&outVal.Account.Email,
			&outVal.Account.Type,
			&outVal.Account.Locked,
			&outVal.Event.Identifier,
			&outVal.Event.Name,
			&outVal.Event.Slug,
			&outVal.Event.Website,
			&outVal.Event.Image,
			&outVal.Event.ContactEmail,
			&outVal.Event.AccessRestricted,
			&outVal.EventYear.Identifier,
			&outVal.EventYear.Year,
			&outVal.EventYear.DateTime,
			&outVal.EventYear.Live,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting values for account and event: %v", err)
		}
		outVal.Event.AccountIdentifier = outVal.Account.Identifier
		outVal.EventYear.EventIdentifier = outVal.Event.Identifier
		return &outVal, nil
	}
	return nil, nil
}

type MultiGet struct {
	Account   *types.Account
	Event     *types.Event
	EventYear *types.EventYear
}
