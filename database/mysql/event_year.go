package mysql

import (
	"chronokeep/results/types"
	"context"
	"fmt"
	"time"
)

// oldGetEventYear The old getter for an event year.
func (m *MySQL) oldGetEventYear(event_slug, year string) (*types.EventYear, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT event_year_id, event_id, year, date_time, live FROM event_year NATURAL JOIN event WHERE slug=? AND year=? AND year_deleted=FALSE;",
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

// GetEventYear Gets an event year for an event with a slug and a specific year.
func (m *MySQL) GetEventYear(event_slug, year string) (*types.EventYear, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT event_year_id, event_id, year, date_time, live, days_allowed, ranking_type FROM event_year NATURAL JOIN event WHERE slug=? AND year=? AND year_deleted=FALSE;",
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
			&outEventYear.RankingType,
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
func (m *MySQL) GetEventYears(event_slug string) ([]types.EventYear, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT event_year_id, event_id, year, date_time, live, days_allowed, ranking_type FROM event_year NATURAL JOIN event WHERE slug=? AND year_deleted=FALSE;",
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
			&year.RankingType,
		)
		if err != nil {
			return nil, nil
		}
		outEventYears = append(outEventYears, year)
	}
	return outEventYears, nil
}

// GetAllEventYears Gets all event years.
func (m *MySQL) GetAllEventYears() ([]types.EventYear, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT event_year_id, event_id, year, date_time, live, days_allowed, ranking_type FROM event_year NATURAL JOIN event WHERE year_deleted=FALSE;",
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
			&year.RankingType,
		)
		if err != nil {
			return nil, nil
		}
		outEventYears = append(outEventYears, year)
	}
	return outEventYears, nil
}

// oldAddEventYear Old function that adds an event year to the database
func (m *MySQL) oldAddEventYear(year types.EventYear) (*types.EventYear, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"INSERT INTO event_year(event_id, year, date_time, live) VALUES (?, ?, ?, ?);",
		year.EventIdentifier,
		year.Year,
		year.DateTime,
		year.Live,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to add event year: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("unable to determine ID for event year: %v", err)
	}
	return &types.EventYear{
		Identifier:      id,
		EventIdentifier: year.EventIdentifier,
		Year:            year.Year,
		DateTime:        year.DateTime,
		Live:            year.Live,
		DaysAllowed:     year.DaysAllowed,
		RankingType:     year.RankingType,
	}, nil
}

// AddEventYear Adds an event year to the database.
func (m *MySQL) AddEventYear(year types.EventYear) (*types.EventYear, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"INSERT INTO event_year(event_id, year, date_time, live, days_allowed, ranking_type) VALUES (?, ?, ?, ?, ?, ?);",
		year.EventIdentifier,
		year.Year,
		year.DateTime,
		year.Live,
		year.DaysAllowed,
		year.RankingType,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to add event year: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("unable to determine ID for event year: %v", err)
	}
	return &types.EventYear{
		Identifier:      id,
		EventIdentifier: year.EventIdentifier,
		Year:            year.Year,
		DateTime:        year.DateTime,
		Live:            year.Live,
		DaysAllowed:     year.DaysAllowed,
		RankingType:     year.RankingType,
	}, nil
}

// DeleteEventYear Deletes an EventYear from view, does not permanently delete from database.
// This does not cascade down.  Must be done manually.
func (m *MySQL) DeleteEventYear(year types.EventYear) error {
	db, err := m.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	_, err = db.ExecContext(
		ctx,
		"UPDATE event_year SET year_deleted=TRUE WHERE event_year_id=?",
		year.Identifier,
	)
	if err != nil {
		return fmt.Errorf("error deleting event year: %v", err)
	}
	return nil
}

// UpdateEventYear Updates an Event Year information in the database. Year cannot be changed once set.
func (m *MySQL) UpdateEventYear(year types.EventYear) error {
	db, err := m.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	_, err = db.ExecContext(
		ctx,
		"UPDATE event_year SET date_time=?, live=?, days_allowed=?, ranking_type=? WHERE event_year_id=?",
		year.DateTime,
		year.Live,
		year.DaysAllowed,
		year.RankingType,
		year.Identifier,
	)
	if err != nil {
		return fmt.Errorf("error updating event year: %v", err)
	}
	return nil
}
