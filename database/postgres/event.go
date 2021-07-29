package postgres

import (
	"chronokeep/results/types"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
)

// GetEvent Gets an event with a slug.
func (p *Postgres) GetEvent(slug string) (*types.Event, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Query(
		ctx,
		"SELECT event_id, event_name, slug, website, image, account_id, contact_email, access_restricted, event_type, "+
			"recent_time FROM event NATURAL JOIN (SELECT e.event_id, MAX(y.date_time) AS recent_time FROM event e LEFT OUTER "+
			"JOIN event_year y ON e.event_id=y.event_id GROUP BY e.event_id) AS time WHERE event_deleted=FALSE and slug=$1;",
		slug,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving event: %v", err)
	}
	defer res.Close()
	if res.Next() {
		var outEvent types.Event
		err := res.Scan(
			&outEvent.Identifier,
			&outEvent.Name,
			&outEvent.Slug,
			&outEvent.Website,
			&outEvent.Image,
			&outEvent.AccountIdentifier,
			&outEvent.ContactEmail,
			&outEvent.AccessRestricted,
			&outEvent.Type,
			&outEvent.RecentTime,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting event: %v", err)
		}
		return &outEvent, nil
	}
	return nil, nil
}

func (p *Postgres) getEventsInternal(email *string) ([]types.Event, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	var res pgx.Rows
	if email == nil {
		res, err = db.Query(
			ctx,
			"SELECT event_id, event_name, slug, website, image, account_id, contact_email, access_restricted, event_type, "+
				"recent_time FROM event NATURAL JOIN (SELECT e.event_id, MAX(y.date_time) AS recent_time FROM event e LEFT OUTER "+
				"JOIN event_year y ON e.event_id=y.event_id GROUP BY e.event_id) AS time WHERE event_deleted=FALSE AND access_restricted=FALSE;",
		)
	} else {
		res, err = db.Query(
			ctx,
			"SELECT event_id, event_name, slug, website, image, account_id, contact_email, access_restricted, event_type, "+
				"recent_time FROM event NATURAL JOIN account NATURAL JOIN (SELECT e.event_id, MAX(y.date_time) AS recent_time FROM event e LEFT OUTER "+
				"JOIN event_year y ON e.event_id=y.event_id GROUP BY e.event_id) AS time WHERE event_deleted=FALSE AND account_email=$1;",
			email,
		)
	}
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
			&event.Type,
			&event.RecentTime,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting event: %v", err)
		}
		outEvents = append(outEvents, event)
	}
	return outEvents, nil
}

// GetEvents Gets all events.
func (p *Postgres) GetEvents() ([]types.Event, error) {
	return p.getEventsInternal(nil)
}

// GetAccountsEvents Gets all events associated with an account.
func (p *Postgres) GetAccountEvents(email string) ([]types.Event, error) {
	return p.getEventsInternal(&email)
}

// AddEvent Adds an event to the database.
func (p *Postgres) AddEvent(event types.Event) (*types.Event, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	var id int64
	err = db.QueryRow(
		ctx,
		"INSERT INTO event(event_name, slug, website, image, contact_email, account_id, access_restricted, event_type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING (event_id);",
		event.Name,
		event.Slug,
		event.Website,
		event.Image,
		event.ContactEmail,
		event.AccountIdentifier,
		event.AccessRestricted,
		event.Type,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("unable to add event: %v", err)
	}
	if id == 0 {
		return nil, errors.New("id value set to 0")
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
		Type:              event.Type,
	}, nil
}

// DeleteEvent Deletes and event from view, does not permanently delete from database.
// This does not cascade down. Must be done manually.
func (p *Postgres) DeleteEvent(event types.Event) error {
	db, err := p.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Exec(
		ctx,
		"UPDATE event SET event_deleted=TRUE WHERE event_id=$1;",
		event.Identifier,
	)
	if err != nil {
		return fmt.Errorf("error deleting event: %v", err)
	}
	if res.RowsAffected() != 1 {
		return fmt.Errorf("error deleting event, rows affected: %v", res.RowsAffected())
	}
	_, err = db.Exec(
		ctx,
		"UPDATE event_year SET year_deleted=TRUE where event_id=$1;",
		event.Identifier,
	)
	if err != nil {
		return fmt.Errorf("error deleting event years attached to event: %v", err)
	}
	return nil
}

// RealDeleteEvent Really deletes an event from the database.
func (p *Postgres) RealDeleteEvent(event types.Event) error {
	db, err := p.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	_, err = db.Exec(
		ctx,
		"DELETE FROM result r WHERE EXISTS (SELECT * FROM person p NATURAL JOIN event_year y WHERE r.person_id=p.person_id AND y.event_id=$1);",
		event.Identifier,
	)
	if err != nil {
		return fmt.Errorf("error deleting event results: %v", err)
	}
	_, err = db.Exec(
		ctx,
		"DELETE FROM person p WHERE EXISTS (SELECT * FROM event_year y WHERE p.event_year_id=y.event_year_id AND y.event_id=$1);",
		event.Identifier,
	)
	if err != nil {
		return fmt.Errorf("error deleting event people: %v", err)
	}
	_, err = db.Exec(
		ctx,
		"DELETE FROM event_year WHERE event_id=$1;",
		event.Identifier,
	)
	if err != nil {
		return fmt.Errorf("error deleting event years: %v", err)
	}
	_, err = db.Exec(
		ctx,
		"DELETE FROM event WHERE event_id=$1;",
		event.Identifier,
	)
	if err != nil {
		return fmt.Errorf("error deleting event: %v", err)
	}
	return nil
}

// UpdateEvent Updates an Event in the database. Name and Slug cannot be changed once set.
func (p *Postgres) UpdateEvent(event types.Event) error {
	db, err := p.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Exec(
		ctx,
		"UPDATE event SET event_name=$1, website=$2, image=$3, contact_email=$4, access_restricted=$5, event_type=$7 WHERE event_id=$6;",
		event.Name,
		event.Website,
		event.Image,
		event.ContactEmail,
		event.AccessRestricted,
		event.Identifier,
		event.Type,
	)
	if err != nil {
		return fmt.Errorf("error updating event: %v", err)
	}
	if res.RowsAffected() < 1 {
		return errors.New("no rows affected")
	}
	return nil
}
