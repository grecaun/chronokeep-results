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

// GetAccountCallRecords Gets all api call records for a specific account.
func (p *Postgres) GetAccountCallRecords(email string) ([]types.CallRecord, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Query(
		ctx,
		"SELECT account_id, time, count FROM call_record NATURAL JOIN account WHERE account_email=$1;",
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
func (p *Postgres) GetCallRecord(email string, inTime int64) (*types.CallRecord, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Query(
		ctx,
		"SELECT account_id, time, count FROM call_record NATURAL JOIN account WHERE account_email=$1 AND time=$2;",
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
func (p *Postgres) AddCallRecord(record types.CallRecord) error {
	db, err := p.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Exec(
		ctx,
		"INSERT INTO call_record(account_id, time, count) VALUES ($1, $2, $3) ON CONFLICT (account_id, time) DO UPDATE SET count=$3;",
		record.AccountIdentifier,
		record.DateTime,
		record.Count,
	)
	if err != nil {
		return fmt.Errorf("error adding call record: %v", err)
	}
	if res.RowsAffected() < 1 {
		return fmt.Errorf("call record insert rows affected: %v", res.RowsAffected())
	}
	return nil
}

// AddCallRecords Add multiple call records to the database.
func (p *Postgres) AddCallRecords(records []types.CallRecord) error {
	db, err := p.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("unable to start transaction: %v", err)
	}
	for _, record := range records {
		_, err = tx.Exec(
			ctx,
			"INSERT INTO call_record(account_id, time, count) VALUES ($1, $2, $3) ON CONFLICT (account_id, time) DO UPDATE SET count=$3;",
			record.AccountIdentifier,
			record.DateTime,
			record.Count,
		)
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("error adding call record to database: %v", err)
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("error committing transaction: %v", err)
	}
	return nil
}

/*
func (p *Postgres) deleteCallRecords() (int64, error) {
	db, err := p.GetDB()
	if err != nil {
		return 0, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Exec(
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

func (p *Postgres) deleteCallRecord(record types.CallRecord) (int64, error) {
	db, err := p.GetDB()
	if err != nil {
		return 0, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Exec(
		ctx,
		"DELETE FROM call_records WHERE account_id=$1 AND time=$2;",
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

