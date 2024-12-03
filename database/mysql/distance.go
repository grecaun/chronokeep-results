package mysql

import (
	"chronokeep/results/types"
	"context"
	"fmt"
	"time"
)

func (m *MySQL) AddDistances(eventYearID int64, distances []types.Distance) ([]types.Distance, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction: %v", err)
	}
	distancestmt, err := tx.PrepareContext(
		ctx,
		"INSERT INTO distances("+
			"event_year_id, "+
			"distance_name, "+
			"certification"+
			") VALUES (?,?,?) "+
			"ON DUPLICATE KEY UPDATE "+
			"certification=VALUES(certification)"+
			";",
	)
	if err != nil {
		return nil, fmt.Errorf("error preparing statement for adding distances: %v", err)
	}
	for _, dist := range distances {
		_, err = distancestmt.ExecContext(
			ctx,
			eventYearID,
			dist.Name,
			dist.Certification,
		)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error adding distance to database: %v", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}
	output := make([]types.Distance, 0)
	for _, dist := range distances {
		output = append(output, types.Distance{
			Name:          dist.Name,
			Certification: dist.Certification,
		})
	}
	return output, nil
}

func (m *MySQL) GetDistance(eventYearID int64, dist_name string) (*types.Distance, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT distance_id, distance_name, certification "+
			"FROM distances WHERE event_year_id=? AND distance_name=?;",
		eventYearID,
		dist_name,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving distances: %v", err)
	}
	defer res.Close()
	if res.Next() {
		var dist types.Distance
		err := res.Scan(
			&dist.Identifier,
			&dist.Name,
			&dist.Certification,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting distance: %v", err)
		}
		return &dist, nil
	}
	return nil, nil
}

func (m *MySQL) GetDistances(eventYearID int64) ([]types.Distance, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT distance_id, distance_name, certification "+
			"FROM distances WHERE event_year_id=?;",
		eventYearID,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving distances: %v", err)
	}
	defer res.Close()
	output := make([]types.Distance, 0)
	for res.Next() {
		var dist types.Distance
		err := res.Scan(
			&dist.Identifier,
			&dist.Name,
			&dist.Certification,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting distance: %v", err)
		}
		output = append(output, dist)
	}
	return output, nil
}

func (m *MySQL) DeleteDistances(eventYearID int64) (int64, error) {
	db, err := m.GetDB()
	if err != nil {
		return 0, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("error starting transaction: %v", err)
	}
	res, err := tx.ExecContext(
		ctx,
		"DELETE FROM distances WHERE event_year_id=?",
		eventYearID,
	)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error deleting distances: %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error fetching rows affected from distances deletion: %v", err)
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error committing transaction: %v", err)
	}
	return count, nil
}
