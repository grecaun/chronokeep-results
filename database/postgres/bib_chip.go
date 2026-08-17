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

func (p *Postgres) AddBibChips(eventYearID int64, bibChips []types.BibChip) ([]types.BibChip, error) {
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
	for _, bc := range bibChips {
		_, err = tx.Exec(
			ctx,
			"INSERT INTO chips("+
				"event_year_id, "+
				"bib, "+
				"chip"+
				") VALUES ($1,$2,$3) "+
				"ON CONFLICT (event_year_id, chip) DO UPDATE SET "+
				"bib=$2"+
				";",
			eventYearID,
			bc.Bib,
			bc.Chip,
		)
		if err != nil {
			tx.Rollback(ctx)
			return nil, fmt.Errorf("error adding bibchip to database: %v", err)
		}
	}
	output := make([]types.BibChip, 0)
	for _, bc := range bibChips {
		output = append(output, types.BibChip{
			Bib:  bc.Bib,
			Chip: bc.Chip,
		})
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("error commiting transaction: %v", err)
	}
	return output, nil
}

func (p *Postgres) GetBibChips(eventYearID int64) ([]types.BibChip, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Query(
		ctx,
		"SELECT chip_id, bib, chip FROM chips WHERE event_year_id=$1;",
		eventYearID,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving bibchips: %v", err)
	}
	defer res.Close()
	output := make([]types.BibChip, 0)
	for res.Next() {
		var bc types.BibChip
		err = res.Scan(
			&bc.Identifier,
			&bc.Bib,
			&bc.Chip,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting bibchip: %v", err)
		}
		output = append(output, bc)
	}
	return output, nil
}

func (p *Postgres) DeleteBibChips(eventYearID int64) (int64, error) {
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
	var count int64
	res, err := tx.Exec(
		ctx,
		"DELETE FROM chips WHERE event_year_id=$1;",
		eventYearID,
	)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("error deleting bibchips: %v", err)
	}
	count = res.RowsAffected()
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("error committing transaction: %v", err)
	}
	return count, nil
}

