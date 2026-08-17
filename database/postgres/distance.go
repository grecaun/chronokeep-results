/*
Chronokeep Desktop - Race Scoring Software
Copyright (C) 2026 James Sentinella

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package postgres

import (
	"chronokeep/results/types"
	"context"
	"fmt"
	"time"
)

func (p *Postgres) AddDistances(eventYearID int64, distances []types.Distance) ([]types.Distance, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction: %v", err)
	}
	for _, dist := range distances {
		_, err = tx.Exec(
			ctx,
			"INSERT INTO distances("+
				"event_year_id, "+
				"distance_name, "+
				"certification"+
				") VALUES ($1,$2,$3) "+
				"ON CONFLICT (event_year_id, distance_name) DO UPDATE SET "+
				"certification=$3"+
				";",
			eventYearID,
			dist.Name,
			dist.Certification,
		)
		if err != nil {
			tx.Rollback(ctx)
			return nil, fmt.Errorf("error adding distance to database: %v", err)
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
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

func (p *Postgres) GetDistance(eventYearID int64, dist_name string) (*types.Distance, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Query(
		ctx,
		"SELECT distance_id, distance_name, certification "+
			"FROM distances WHERE event_year_id=$1 AND distance_name=$2;",
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

func (p *Postgres) GetDistances(eventYearID int64) ([]types.Distance, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Query(
		ctx,
		"SELECT distance_id, distance_name, certification "+
			"FROM distances WHERE event_year_id=$1;",
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

func (p *Postgres) DeleteDistances(eventYearID int64) (int64, error) {
	db, err := p.GetDB()
	if err != nil {
		return 0, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("error starting transaction: %v", err)
	}
	res, err := tx.Exec(
		ctx,
		"DELETE FROM distances WHERE event_year_id=$1",
		eventYearID,
	)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("error deleting distances: %v", err)
	}
	count := res.RowsAffected()
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("error committing transaction: %v", err)
	}
	return count, nil
}

