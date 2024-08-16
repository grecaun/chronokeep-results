package postgres

import (
	"chronokeep/results/types"
	"context"
	"errors"
	"fmt"
	"time"
)

// GetEventYear Gets an event year for an event with a slug and a specific year.
func (p *Postgres) GetEventYear(event_slug, year string) (*types.EventYear, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Query(
		ctx,
		"SELECT event_year_id, event_id, year, date_time, live, days_allowed FROM event_year NATURAL JOIN event WHERE slug=$1 AND year=$2 AND year_deleted=FALSE;",
		event_slug,
		year,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving event year: %v", err)
	}
	defer res.Close()
	var outEventYear types.EventYear
	if res.Next() {
		err := res.Scan(
			&outEventYear.Identifier,
			&outEventYear.EventIdentifier,
			&outEventYear.Year,
			&outEventYear.DateTime,
			&outEventYear.Live,
			&outEventYear.DaysAllowed,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting event year: %v", err)
		}
	} else {
		return nil, nil
	}
	return &outEventYear, nil
}

// oldGetEventYear Gets an event year for an event with a slug and a specific year.
func (p *Postgres) oldGetEventYear(event_slug, year string) (*types.EventYear, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Query(
		ctx,
		"SELECT event_year_id, event_id, year, date_time, live FROM event_year NATURAL JOIN event WHERE slug=$1 AND year=$2 AND year_deleted=FALSE;",
		event_slug,
		year,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving event year: %v", err)
	}
	defer res.Close()
	var outEventYear types.EventYear
	if res.Next() {
		err := res.Scan(
			&outEventYear.Identifier,
			&outEventYear.EventIdentifier,
			&outEventYear.Year,
			&outEventYear.DateTime,
			&outEventYear.Live,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting event year: %v", err)
		}
	} else {
		return nil, nil
	}
	return &outEventYear, nil
}

// GetEventYears Gets all event years for a specific event based on the slug.
func (p *Postgres) GetEventYears(event_slug string) ([]types.EventYear, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Query(
		ctx,
		"SELECT event_year_id, event_id, year, date_time, live, days_allowed FROM event_year NATURAL JOIN event WHERE slug=$1 AND year_deleted=FALSE;",
		event_slug,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving event years: %v", err)
	}
	defer res.Close()
	var outEventYears []types.EventYear
	for res.Next() {
		var year types.EventYear
		err := res.Scan(
			&year.Identifier,
			&year.EventIdentifier,
			&year.Year,
			&year.DateTime,
			&year.Live,
			&year.DaysAllowed,
		)
		if err != nil {
			return nil, nil
		}
		outEventYears = append(outEventYears, year)
	}
	return outEventYears, nil
}

// GetAllEventYears Gets all event years.
func (p *Postgres) GetAllEventYears() ([]types.EventYear, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Query(
		ctx,
		"SELECT event_year_id, event_id, year, date_time, live, days_allowed FROM event_year NATURAL JOIN event WHERE year_deleted=FALSE;",
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving event years: %v", err)
	}
	defer res.Close()
	var outEventYears []types.EventYear
	for res.Next() {
		var year types.EventYear
		err := res.Scan(
			&year.Identifier,
			&year.EventIdentifier,
			&year.Year,
			&year.DateTime,
			&year.Live,
			&year.DaysAllowed,
		)
		if err != nil {
			return nil, nil
		}
		outEventYears = append(outEventYears, year)
	}
	return outEventYears, nil
}

// AddEventYear Adds an event year to the database.
func (p *Postgres) AddEventYear(year types.EventYear) (*types.EventYear, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	var id int64
	err = db.QueryRow(
		ctx,
		"INSERT INTO event_year(event_id, year, date_time, live, days_allowed) VALUES ($1, $2, $3, $4, $5) RETURNING (event_year_id);",
		year.EventIdentifier,
		year.Year,
		year.DateTime,
		year.Live,
		year.DaysAllowed,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("unable to add event year: %v", err)
	}
	if id == 0 {
		return nil, errors.New("id value set to 0")
	}
	return &types.EventYear{
		Identifier:      id,
		EventIdentifier: year.EventIdentifier,
		Year:            year.Year,
		DateTime:        year.DateTime,
		Live:            year.Live,
		DaysAllowed:     year.DaysAllowed,
	}, nil
}

// AddEventYear Adds an event year to the database.
func (p *Postgres) oldAddEventYear(year types.EventYear) (*types.EventYear, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	var id int64
	err = db.QueryRow(
		ctx,
		"INSERT INTO event_year(event_id, year, date_time, live) VALUES ($1, $2, $3, $4) RETURNING (event_year_id);",
		year.EventIdentifier,
		year.Year,
		year.DateTime,
		year.Live,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("unable to add event year: %v", err)
	}
	if id == 0 {
		return nil, errors.New("id value set to 0")
	}
	return &types.EventYear{
		Identifier:      id,
		EventIdentifier: year.EventIdentifier,
		Year:            year.Year,
		DateTime:        year.DateTime,
		Live:            year.Live,
	}, nil
}

// DeleteEventYear Deletes an EventYear from view, does not permanently delete from database.
// This does not cascade down.  Must be done manually.
func (p *Postgres) DeleteEventYear(year types.EventYear) error {
	db, err := p.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	_, err = db.Exec(
		ctx,
		"UPDATE event_year SET year_deleted=TRUE WHERE event_year_id=$1",
		year.Identifier,
	)
	if err != nil {
		return fmt.Errorf("error deleting event year: %v", err)
	}
	return nil
}

// UpdateEventYear Updates an Event Year information in the database. Year cannot be changed once set.
func (p *Postgres) UpdateEventYear(year types.EventYear) error {
	db, err := p.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	_, err = db.Exec(
		ctx,
		"UPDATE event_year SET date_time=$1, live=$2, days_allowed=$4 WHERE event_year_id=$3",
		year.DateTime,
		year.Live,
		year.Identifier,
		year.DaysAllowed,
	)
	if err != nil {
		return fmt.Errorf("error updating event year: %v", err)
	}
	return nil
}
