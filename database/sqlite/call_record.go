package sqlite

import (
	"chronokeep/results/types"
	"context"
	"fmt"
	"time"
)

// GetAccountCallRecords Gets all api call records for a specific account.
func (s *SQLite) GetAccountCallRecords(email string) ([]types.CallRecord, error) {
	db, err := s.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT account_id, time, count FROM call_record NATURAL JOIN account WHERE account_email=?;",
		email,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to query for account call records: %v", err)
	}
	defer res.Close()
	var records []types.CallRecord
	for res.Next() {
		var record types.CallRecord
		err := res.Scan(
			&record.AccountIdentifier,
			&record.DateTime,
			&record.Count,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to pull record information: %v", err)
		}
		records = append(records, record)
	}
	return records, nil
}

// GetCallRecord Checks the database for a specific call record.
func (s *SQLite) GetCallRecord(email string, inTime int64) (*types.CallRecord, error) {
	db, err := s.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT account_id, time, count FROM call_record NATURAL JOIN account WHERE account_email=? AND time=?;",
		email,
		inTime,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to query for specific call record: %v", err)
	}
	defer res.Close()
	var record types.CallRecord
	if res.Next() {
		err = res.Scan(
			&record.AccountIdentifier,
			&record.DateTime,
			&record.Count,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting call record information: %v", err)
		}
	} else {
		return nil, nil
	}
	return &record, nil
}

// AddCallRecord Add a call record to the database.
func (s *SQLite) AddCallRecord(record types.CallRecord) error {
	db, err := s.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"INSERT INTO call_record(account_id, time, count) VALUES (?, ?, ?) ON CONFLICT (account_id, time) DO UPDATE SET count=?;",
		record.AccountIdentifier,
		record.DateTime,
		record.Count,
		record.Count,
	)
	if err != nil {
		return fmt.Errorf("error adding call record: %v", err)
	}
	_, err = res.LastInsertId()
	if err != nil {
		return fmt.Errorf("unable to determine ID for call record input: %v", err)
	}
	return nil
}

// AddCallRecords Add multiple call records to the database.
func (s *SQLite) AddCallRecords(records []types.CallRecord) error {
	db, err := s.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("unable to start transaction: %v", err)
	}
	stmt, err := tx.PrepareContext(
		ctx,
		"INSERT INTO call_record(account_id, time, count) VALUES (?,  ?, ?) ON CONFLICT (account_id, time) DO UPDATE SET count=?;",
	)
	if err != nil {
		return fmt.Errorf("unable to prepare statement for multiple call record inserts: %v", err)
	}
	defer stmt.Close()
	for _, record := range records {
		_, err := stmt.ExecContext(
			ctx,
			record.AccountIdentifier,
			record.DateTime,
			record.Count,
			record.Count,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error adding call record to database: %v", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error committing transaction: %v", err)
	}
	return nil
}

/*
func (s *SQLite) deleteCallRecords() (int64, error) {
	db, err := s.GetDB()
	if err != nil {
		return 0, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"DELETE FROM call_records;",
	)
	if err != nil {
		return 0, fmt.Errorf("error deleting all call records: %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("unable to determine the rows deleted: %v", err)
	}
	return count, nil
}

func (s *SQLite) deleteCallRecord(record types.CallRecord) (int64, error) {
	db, err := s.GetDB()
	if err != nil {
		return 0, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"DELETE FROM call_records WHERE account_id=? AND time=?;",
		record.AccountIdentifier,
		record.DateTime,
	)
	if err != nil {
		return 0, fmt.Errorf("error deleting all call records: %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("unable to determine the rows deleted: %v", err)
	}
	return count, nil
}
*/
