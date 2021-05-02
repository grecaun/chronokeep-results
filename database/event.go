package database

import (
	"chronokeep/results/types"
	"context"
	"fmt"
	"time"
)

// GetEvent Gets an event with a slug.
func GetEvent(slug string) (*types.Event, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT event_id, name, slug, website, image, account_id, contact_email, access_restricted FROM event WHERE deleted=FALSE and slug=?;",
		slug,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving event: %v", err)
	}
	defer res.Close()
	var outEvent types.Event
	if res.Next() {
		err := res.Scan(
			&outEvent.Identifier,
			&outEvent.Name,
			&outEvent.Slug,
			&outEvent.Website,
			&outEvent.Image,
			&outEvent.AccountIdentifier,
			&outEvent.ContactEmail,
			&outEvent.AccessRestricted,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting event: %v", err)
		}
	} else {
		return nil, nil
	}
	return &outEvent, nil
}

// GetEvents Gets all events.
func GetEvents() ([]types.Event, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT event_id, name, slug, website, image, account_id, contact_email, access_restricted FROM event WHERE deleted=FALSE;",
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving event: %v", err)
	}
	defer res.Close()
	var outEvents []types.Event
	for res.Next() {
		var event types.Event
		err := res.Scan(
			&event.Identifier,
			&event.Name,
			&event.Slug,
			&event.Website,
			&event.Image,
			&event.AccountIdentifier,
			&event.ContactEmail,
			&event.AccessRestricted,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting event: %v", err)
		}
		outEvents = append(outEvents, event)
	}
	return outEvents, nil
}

// GetAccountsEvents Gets all events associated with an account.
func GetAccountEvents(email string) ([]types.Event, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT event_id, name, slug, website, image, account_id, contact_email, access_restricted FROM event NATURAL JOIN account WHERE deleted=FALSE AND email=?;",
		email,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving event: %v", err)
	}
	defer res.Close()
	var outEvents []types.Event
	for res.Next() {
		var event types.Event
		err := res.Scan(
			&event.Identifier,
			&event.Name,
			&event.Slug,
			&event.Website,
			&event.Image,
			&event.AccountIdentifier,
			&event.ContactEmail,
			&event.AccessRestricted,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting event: %v", err)
		}
		outEvents = append(outEvents, event)
	}
	return outEvents, nil
}

// AddEvent Adds an event to the database.
func AddEvent(event types.Event) (*types.Event, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"INSERT INTO event(name, slug, website, image, contact_email, account_id, access_restricted) VALUES (?, ?, ?, ?, ?, ?, ?);",
		event.Name,
		event.Slug,
		event.Website,
		event.Image,
		event.ContactEmail,
		event.AccountIdentifier,
		event.AccessRestricted,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to add event: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("unable to determine ID for event: %v", err)
	}
	return &types.Event{
		Identifier:        id,
		AccountIdentifier: event.AccountIdentifier,
		Name:              event.Name,
		Slug:              event.Slug,
		Website:           event.Website,
		Image:             event.Image,
		ContactEmail:      event.ContactEmail,
		AccessRestricted:  event.AccessRestricted,
	}, nil
}

// DeleteEvent Deletes and event from view, does not permanently delete from database.
// This does not cascade down. Must be done manually.
func DeleteEvent(event types.Event) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"UPDATE event SET deleted=TRUE WHERE event_id=?",
		event.Identifier,
	)
	if err != nil {
		return fmt.Errorf("error deleting event: %v", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected on delete event: %v", err)
	}
	if rows != 1 {
		return fmt.Errorf("error deleting event, rows affected: %v", rows)
	}
	return nil
}

// UpdateEvent Updates an Event in the database. Name and Slug cannot be changed once set.
func UpdateEvent(event types.Event) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"UPDATE event SET website=?, image=?, contact_email=?, access_restricted=? WHERE event_id=?",
		event.Website,
		event.Image,
		event.ContactEmail,
		event.AccessRestricted,
		event.Identifier,
	)
	if err != nil {
		return fmt.Errorf("error updating event: %v", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected on update event: %v", err)
	}
	if rows != 1 {
		return fmt.Errorf("error updating event, rows affected: %v", rows)
	}
	return nil
}
