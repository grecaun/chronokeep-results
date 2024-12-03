package mysql

import (
	"chronokeep/results/types"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// GetEvent Gets an event with a slug.
func (m *MySQL) GetEvent(slug string) (*types.Event, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT event_id, event_name, cert_name, slug, website, image, account_id, contact_email, access_restricted, event_type, "+
			"recent_time FROM event NATURAL JOIN (SELECT e.event_id, MAX(y.date_time) AS recent_time FROM event e LEFT OUTER "+
			"JOIN event_year y ON e.event_id=y.event_id GROUP BY e.event_id) AS time WHERE event_deleted=FALSE and slug=?;",
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
			&outEvent.CertificateName,
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

func (m *MySQL) getEventsInternal(email *string, restricted ...bool) ([]types.Event, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	var res *sql.Rows
	if email == nil {
		if len(restricted) > 0 && restricted[0] {
			res, err = db.QueryContext(
				ctx,
				"SELECT event_id, event_name, cert_name, slug, website, image, account_id, contact_email, access_restricted, event_type, "+
					"recent_time FROM event NATURAL JOIN (SELECT e.event_id, MAX(y.date_time) AS recent_time FROM event e LEFT OUTER "+
					"JOIN event_year y ON e.event_id=y.event_id GROUP BY e.event_id) AS time WHERE event_deleted=FALSE;",
			)
		} else {
			res, err = db.QueryContext(
				ctx,
				"SELECT event_id, event_name, cert_name, slug, website, image, account_id, contact_email, access_restricted, event_type, "+
					"recent_time FROM event NATURAL JOIN (SELECT e.event_id, MAX(y.date_time) AS recent_time FROM event e LEFT OUTER "+
					"JOIN event_year y ON e.event_id=y.event_id GROUP BY e.event_id) AS time WHERE event_deleted=FALSE AND access_restricted=FALSE;",
			)
		}
	} else {
		res, err = db.QueryContext(
			ctx,
			"SELECT event_id, event_name, cert_name, slug, website, image, account_id, contact_email, access_restricted, event_type, "+
				"recent_time FROM event NATURAL JOIN account a NATURAL JOIN (SELECT e.event_id, MAX(y.date_time) AS recent_time FROM event e LEFT OUTER "+
				"JOIN event_year y ON e.event_id=y.event_id GROUP BY e.event_id) AS time WHERE event_deleted=FALSE AND (account_email=? "+
				"OR EXISTS (SELECT sub_account_id FROM linked_accounts JOIN account b ON b.account_id=sub_account_id WHERE "+
				"a.account_id=main_account_id AND b.account_email=?));",
			email,
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
			&event.CertificateName,
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

// GetEvents Gets all events not including restricted.
func (m *MySQL) GetEvents() ([]types.Event, error) {
	return m.getEventsInternal(nil)
}

// GetEvents Gets all events.
func (m *MySQL) GetAllEvents() ([]types.Event, error) {
	return m.getEventsInternal(nil, true)
}

// GetAccountsEvents Gets all events associated with an account.
func (m *MySQL) GetAccountEvents(email string) ([]types.Event, error) {
	return m.getEventsInternal(&email)
}

// AddEvent Adds an event to the database.
func (m *MySQL) AddEvent(event types.Event) (*types.Event, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"INSERT INTO event(event_name, cert_name, slug, website, image, contact_email, account_id, access_restricted, event_type) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);",
		event.Name,
		event.CertificateName,
		event.Slug,
		event.Website,
		event.Image,
		event.ContactEmail,
		event.AccountIdentifier,
		event.AccessRestricted,
		event.Type,
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
		CertificateName:   event.CertificateName,
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
func (m *MySQL) DeleteEvent(event types.Event) error {
	db, err := m.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("unable to start transaction: %v", err)
	}
	res, err := tx.ExecContext(
		ctx,
		"UPDATE event SET event_deleted=TRUE WHERE event_id=?;",
		event.Identifier,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting event: %v", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error checking rows affected on delete event: %v", err)
	}
	if rows != 1 {
		tx.Rollback()
		return fmt.Errorf("error deleting event, rows affected: %v", rows)
	}
	_, err = tx.ExecContext(
		ctx,
		"UPDATE event_year SET year_deleted=TRUE where event_id=?;",
		event.Identifier,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting event: %v", err)
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error committing transaction: %v", err)
	}
	return nil
}

// RealDeleteEvent Really deletes an event from the database.
func (m *MySQL) RealDeleteEvent(event types.Event) error {
	db, err := m.GetDB()
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
		"DELETE FROM result r WHERE EXISTS (SELECT * FROM person p NATURAL JOIN event_year y WHERE r.person_id=p.person_id AND y.event_id=?);",
		event.Identifier,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting event results: %v", err)
	}
	_, err = tx.ExecContext(
		ctx,
		"DELETE FROM person p WHERE EXISTS (SELECT * FROM event_year y WHERE p.event_year_id=y.event_year_id AND y.event_id=?);",
		event.Identifier,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting event people: %v", err)
	}
	_, err = tx.ExecContext(
		ctx,
		"DELETE FROM event_year WHERE event_id=?;",
		event.Identifier,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting event years: %v", err)
	}
	_, err = tx.ExecContext(
		ctx,
		"DELETE FROM event WHERE event_id=?;",
		event.Identifier,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting event: %v", err)
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error committing transaction: %v", err)
	}
	return nil
}

// UpdateEvent Updates an Event in the database. Name and Slug cannot be changed once set.
func (m *MySQL) UpdateEvent(event types.Event) error {
	db, err := m.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"UPDATE event SET event_name=?, cert_name=?, website=?, image=?, contact_email=?, access_restricted=?, event_type=? WHERE event_id=?;",
		event.Name,
		event.CertificateName,
		event.Website,
		event.Image,
		event.ContactEmail,
		event.AccessRestricted,
		event.Type,
		event.Identifier,
	)
	if err != nil {
		return fmt.Errorf("error updating event: %v", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected on update event: %v", err)
	}
	if rows < 1 {
		return errors.New("no rows affected")
	}
	return nil
}
