package sqlite

import (
	"chronokeep/results/types"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type ResultType int

const (
	All ResultType = iota
	Finish
	Last
)

func (s *SQLite) getResultsInternal(eventYearID int64, bib *string, rtype ResultType, distance string, limit, page int) ([]types.Result, error) {
	db, err := s.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	var (
		res *sql.Rows
	)
	if bib != nil {
		if limit > 0 {
			res, err = db.QueryContext(
				ctx,
				"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
					"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
					"gender_ranking, finish, result_type, chip, anonymous, alternate_id FROM result NATURAL JOIN person "+
					"WHERE event_year_id=? AND bib=? ORDER BY seconds ASC LIMIT ? OFFSET ?;",
				eventYearID,
				bib,
				limit,
				page*limit,
			)
		} else {
			res, err = db.QueryContext(
				ctx,
				"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
					"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
					"gender_ranking, finish, result_type, chip, anonymous, alternate_id FROM result NATURAL JOIN person "+
					"WHERE event_year_id=? AND bib=? ORDER BY seconds ASC;",
				eventYearID,
				bib,
			)
		}
	} else if distance != "" {
		if rtype == Finish {
			if limit > 0 {
				res, err = db.QueryContext(
					ctx,
					"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
						"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
						"gender_ranking, finish, result_type, chip, anonymous, alternate_id FROM result NATURAL JOIN person "+
						"WHERE finish=TRUE AND event_year_id=? AND distance=? ORDER BY seconds ASC LIMIT ? OFFSET ?;",
					eventYearID,
					distance,
					limit,
					page*limit,
				)
			} else {
				res, err = db.QueryContext(
					ctx,
					"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
						"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
						"gender_ranking, finish, result_type, chip, anonymous, alternate_id FROM result NATURAL JOIN person "+
						"WHERE finish=TRUE AND event_year_id=? AND distance=? ORDER BY seconds ASC;",
					eventYearID,
					distance,
				)
			}
		} else if rtype == All {
			if limit > 0 {
				res, err = db.QueryContext(
					ctx,
					"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
						"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
						"gender_ranking, finish, result_type, chip, anonymous, alternate_id FROM result NATURAL JOIN person "+
						"WHERE event_year_id=? AND distance=? ORDER BY seconds ASC LIMIT ? OFFSET ?;",
					eventYearID,
					distance,
					limit,
					page*limit,
				)
			} else {
				res, err = db.QueryContext(
					ctx,
					"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
						"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
						"gender_ranking, finish, result_type, chip, anonymous, alternate_id FROM result NATURAL JOIN person "+
						"WHERE event_year_id=? AND distance=? ORDER BY seconds ASC;",
					eventYearID,
					distance,
				)
			}
		} else {
			if limit > 0 {
				res, err = db.QueryContext(
					ctx,
					"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
						"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
						"gender_ranking, finish, result_type, chip, anonymous, alternate_id FROM result r NATURAL JOIN person p "+
						"JOIN (SELECT bib AS mx_bib, event_year_id AS mx_event_year_id, MAX(seconds) as mx_seconds "+
						"FROM result NATURAL JOIN person GROUP BY bib, event_year_id, segment) b "+
						"ON b.mx_bib=p.bib AND b.mx_event_year_id=p.event_year_id "+
						"AND b.mx_seconds=r.seconds "+
						"WHERE event_year_id=? AND distance=? ORDER BY seconds ASC LIMIT ? OFFSET ?;",
					eventYearID,
					distance,
					limit,
					page*limit,
				)
			} else {
				res, err = db.QueryContext(
					ctx,
					"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
						"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
						"gender_ranking, finish, result_type, chip, anonymous, alternate_id FROM result r NATURAL JOIN person p "+
						"JOIN (SELECT bib AS mx_bib, event_year_id AS mx_event_year_id, MAX(seconds) as mx_seconds "+
						"FROM result NATURAL JOIN person GROUP BY bib, event_year_id, segment) b "+
						"ON b.mx_bib=p.bib AND b.mx_event_year_id=p.event_year_id "+
						"AND b.mx_seconds=r.seconds "+
						"WHERE event_year_id=? AND distance=? ORDER BY seconds ASC;",
					eventYearID,
					distance,
				)
			}
		}
	} else if rtype == All {
		if limit > 0 {
			res, err = db.QueryContext(
				ctx,
				"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
					"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
					"gender_ranking, finish, result_type, chip, anonymous, alternate_id FROM result NATURAL JOIN person "+
					"WHERE event_year_id=? ORDER BY seconds ASC LIMIT ? OFFSET ?;",
				eventYearID,
				limit,
				page*limit,
			)
		} else {
			res, err = db.QueryContext(
				ctx,
				"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
					"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
					"gender_ranking, finish, result_type, chip, anonymous, alternate_id FROM result NATURAL JOIN person "+
					"WHERE event_year_id=? ORDER BY seconds ASC;",
				eventYearID,
			)
		}
	} else if rtype == Finish {
		if limit > 0 {
			res, err = db.QueryContext(
				ctx,
				"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
					"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
					"gender_ranking, finish, result_type, chip, anonymous, alternate_id FROM result NATURAL JOIN person "+
					"WHERE finish=TRUE AND event_year_id=? ORDER BY seconds ASC LIMIT ? OFFSET ?;",
				eventYearID,
				limit,
				page*limit,
			)
		} else {
			res, err = db.QueryContext(
				ctx,
				"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
					"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
					"gender_ranking, finish, result_type, chip, anonymous, alternate_id FROM result NATURAL JOIN person "+
					"WHERE finish=TRUE AND event_year_id=? ORDER BY seconds ASC;",
				eventYearID,
			)
		}
	} else {
		if limit > 0 {
			res, err = db.QueryContext(
				ctx,
				"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
					"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
					"gender_ranking, finish, result_type, chip, anonymous, alternate_id FROM result r NATURAL JOIN person p "+
					"JOIN (SELECT bib AS mx_bib, event_year_id AS mx_event_year_id, MAX(seconds) as mx_seconds "+
					"FROM result NATURAL JOIN person GROUP BY bib, event_year_id, segment) b "+
					"ON b.mx_bib=p.bib AND b.mx_event_year_id=p.event_year_id "+
					"AND b.mx_seconds=r.seconds "+
					"WHERE event_year_id=? ORDER BY seconds ASC LIMIT ? OFFSET ?;",
				eventYearID,
				limit,
				page*limit,
			)
		} else {
			res, err = db.QueryContext(
				ctx,
				"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
					"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
					"gender_ranking, finish, result_type, chip, anonymous, alternate_id FROM result r NATURAL JOIN person p "+
					"JOIN (SELECT bib AS mx_bib, event_year_id AS mx_event_year_id, MAX(seconds) as mx_seconds "+
					"FROM result NATURAL JOIN person GROUP BY bib, event_year_id, segment) b "+
					"ON b.mx_bib=p.bib AND b.mx_event_year_id=p.event_year_id "+
					"AND b.mx_seconds=r.seconds "+
					"WHERE event_year_id=? ORDER BY seconds ASC;",
				eventYearID,
			)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("error retrieving results: %v", err)
	}
	defer res.Close()
	var outResults []types.Result
	for res.Next() {
		var result types.Result
		var anonymous int
		err := res.Scan(
			&result.Bib,
			&result.First,
			&result.Last,
			&result.Age,
			&result.Gender,
			&result.AgeGroup,
			&result.Distance,
			&result.Seconds,
			&result.Milliseconds,
			&result.ChipSeconds,
			&result.ChipMilliseconds,
			&result.Segment,
			&result.Location,
			&result.Occurence,
			&result.Ranking,
			&result.AgeRanking,
			&result.GenderRanking,
			&result.Finish,
			&result.Type,
			&result.Chip,
			&anonymous,
			&result.AlternateId,
		)
		result.Anonymous = anonymous != 0
		if err != nil {
			return nil, fmt.Errorf("error getting result: %v", err)
		}
		outResults = append(outResults, result)
	}
	return outResults, nil
}

// GetResults Gets results for an event year.
func (s *SQLite) GetResults(eventYearID int64, limit, page int) ([]types.Result, error) {
	return s.getResultsInternal(eventYearID, nil, All, "", limit, page)
}

// GetLastResults Gets the last result for a bib in an event year.
func (s *SQLite) GetLastResults(eventYearID int64, limit, page int) ([]types.Result, error) {
	return s.getResultsInternal(eventYearID, nil, Last, "", limit, page)
}

// GetDistanceResults Gets the distance results (last only) for a distance.
func (s *SQLite) GetDistanceResults(eventYearID int64, distance string, limit, page int) ([]types.Result, error) {
	return s.getResultsInternal(eventYearID, nil, Last, distance, limit, page)
}

// GetAllDistanceResults Gets the distance results (all) for a distance.
func (s *SQLite) GetAllDistanceResults(eventYearId int64, distance string, limit, page int) ([]types.Result, error) {
	return s.getResultsInternal(eventYearId, nil, All, distance, limit, page)
}

// GetFinishResults Gets the finish results for an entire event (empty distance) or just a distance.
func (s *SQLite) GetFinishResults(eventYearID int64, distance string, limit, page int) ([]types.Result, error) {
	return s.getResultsInternal(eventYearID, nil, Finish, distance, limit, page)
}

// GetBibResults Gets results for an event year of a specific individual specified by their bib.
func (s *SQLite) GetBibResults(eventYearID int64, bib string) ([]types.Result, error) {
	return s.getResultsInternal(eventYearID, &bib, All, "", 0, 0)
}

// DeleteResults Deletes results from the database.
func (s *SQLite) DeleteResults(eventYearID int64, results []types.Result) error {
	db, err := s.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("unable to begin transaction to delete results: %v", err)
	}
	stmt, err := tx.PrepareContext(
		ctx,
		"DELETE FROM result AS r WHERE location=$1 AND occurence=$2 AND EXISTS (SELECT * FROM person AS p WHERE event_year_id=$3 AND bib=$4 AND p.person_id=r.person_id);",
	)
	if err != nil {
		return fmt.Errorf("unable to get prepared statement for result deletion: %v", err)
	}
	defer stmt.Close()
	for _, result := range results {
		_, err := stmt.ExecContext(
			ctx,
			result.Location,
			result.Occurence,
			eventYearID,
			result.Bib,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error executing prepared delete statement: %v", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error committing transaction: %v", err)
	}
	return nil
}

// DeleteEventResults Deletes results for an event year.
func (s *SQLite) DeleteEventResults(eventYearID int64) (int64, error) {
	db, err := s.GetDB()
	if err != nil {
		return 0, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("unable to start transaction: %v", err)
	}
	res, err := tx.ExecContext(
		ctx,
		"DELETE FROM result AS r WHERE EXISTS (SELECT * FROM person AS p WHERE p.event_year_id=$1 AND p.person_id=r.person_id);",
		eventYearID,
	)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("unable to delete results for event year: %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error fetching rows affected from event year results deletion: %v", err)
	}
	_, err = tx.ExecContext(
		ctx,
		"DELETE FROM person WHERE event_year_id=?;",
		eventYearID,
	)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("unable to delete persons for event year: %v", err)
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error committing transaction: %v", err)
	}
	return count, nil
}

// AddResults Adds results to the database. (Also updates.)
func (s *SQLite) AddResults(eventYearID int64, results []types.Result) ([]types.Result, error) {
	db, err := s.GetDB()
	if err != nil {
		return nil, err
	}
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction: %v", err)
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*15)
	defer cancelfunc()
	personstmt, err := tx.PrepareContext(
		ctx,
		"INSERT INTO person("+
			"event_year_id, "+
			"bib, "+
			"first, "+
			"last, "+
			"age, "+
			"gender, "+
			"age_group, "+
			"distance, "+
			"chip, "+
			"anonymous, "+
			"alternate_id"+
			") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) "+
			"ON CONFLICT (event_year_id, alternate_id) DO UPDATE SET "+
			"first=$3, "+
			"last=$4, "+
			"age=$5, "+
			"gender=$6, "+
			"age_group=$7, "+
			"distance=$8, "+
			"chip=$9, "+
			"anonymous=$10"+
			";",
	)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare statement for person add: %v", err)
	}
	defer personstmt.Close()
	stmt, err := tx.PrepareContext(
		ctx,
		"INSERT INTO result("+
			"person_id, "+
			"seconds, "+
			"milliseconds, "+
			"chip_seconds, "+
			"chip_milliseconds, "+
			"segment, "+
			"location, "+
			"occurence, "+
			"ranking, "+
			"age_ranking, "+
			"gender_ranking, "+
			"finish, "+
			"result_type"+
			") "+
			" VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) "+
			"ON CONFLICT (person_id, location, occurence) DO UPDATE SET "+
			"seconds=$2, "+
			"milliseconds=$3, "+
			"chip_seconds=$4, "+
			"chip_milliseconds=$5, "+
			"segment=$6, "+
			"ranking=$9, "+
			"age_ranking=$10, "+
			"gender_ranking=$11, "+
			"finish=$12, "+
			"result_type=$13;",
	)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare statement for result add: %v", err)
	}
	defer stmt.Close()
	// Get person ids from the database.
	bibMap := make(map[string]int64)
	res, err := tx.QueryContext(
		ctx,
		"SELECT person_id, bib FROM person WHERE event_year_id=$1;",
		eventYearID,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying database for person ids: %v", err)
	}
	for res.Next() {
		var id int64
		var bib string
		err = res.Scan(
			&id,
			&bib,
		)
		if err != nil {
			return nil, fmt.Errorf("error retrieving person ids: %v", err)
		}
		bibMap[bib] = id
	}
	// Add people to the database and get their ids.
	var outResults []types.Result
	for _, result := range results {
		res, err := personstmt.ExecContext(
			ctx,
			eventYearID,
			result.Bib,
			result.First,
			result.Last,
			result.Age,
			result.Gender,
			result.AgeGroup,
			result.Distance,
			result.Chip,
			result.AnonyInt(),
			result.AlternateId,
		)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error adding person to database: %v", err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error adding person to database: %v", err)
		}
		// LastInsertId for SQLite returns an autoincremented ID even if nothing was inserted into the database.
		// Thus, check if we already know the correct person ID for someone before trying to use the insert ID.
		if bib, ok := bibMap[result.Bib]; ok {
			id = bib
		} else if id > 0 {
			bibMap[result.Bib] = id
		}
		_, err = stmt.ExecContext(
			ctx,
			id,
			result.Seconds,
			result.Milliseconds,
			result.ChipSeconds,
			result.ChipMilliseconds,
			result.Segment,
			result.Location,
			result.Occurence,
			result.Ranking,
			result.AgeRanking,
			result.GenderRanking,
			result.Finish,
			result.Type,
		)
		if err != nil {
			tx.Rollback()
			return outResults, fmt.Errorf("error adding result to database: %v", err)
		}
		outResults = append(outResults, result)
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("unable to commit transaction: %v", err)
	}
	return outResults, nil
}
