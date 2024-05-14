package mysql

import (
	"chronokeep/results/types"
	"context"
	"fmt"
	"time"
)

func (m *MySQL) AddBibChips(eventYearID int64, bibChips []types.BibChip) ([]types.BibChip, error) {
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
	chipstmt, err := tx.PrepareContext(
		ctx,
		"INSERT INTO chips("+
			"event_year_id, "+
			"bib, "+
			"chip"+
			") VALUES (?,?,?) "+
			"ON DUPLICATE KEY UPDATE "+
			"bib=VALUES(bib)"+
			";",
	)
	if err != nil {
		return nil, fmt.Errorf("error preparing statement for adding bibchips: %v", err)
	}
	for _, bc := range bibChips {
		_, err = chipstmt.Exec(
			eventYearID,
			bc.Bib,
			bc.Chip,
		)
		if err != nil {
			tx.Rollback()
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
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error commiting transaction: %v", err)
	}
	return output, nil
}

func (m *MySQL) GetBibChips(eventYearID int64) ([]types.BibChip, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT chip_id, bib, chip FROM chips WHERE event_year_id=?;",
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

func (m *MySQL) DeleteBibChips(eventYearID int64) (int64, error) {
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
	var count int64
	res, err := tx.ExecContext(
		ctx,
		"DELETE FROM chips WHERE event_year_id=?;",
		eventYearID,
	)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error deleting bibchips: %v", err)
	}
	count, err = res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error fetching rows affected from bibchips deletion: %v", err)
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error committing transaction: %v", err)
	}
	return count, nil
}
