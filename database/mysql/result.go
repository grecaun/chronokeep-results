package mysql

import (
	"chronokeep/results/types"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (m *MySQL) getResultsInternal(eventYearID int64, bib *string, all bool, distance string) ([]types.Result, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	var (
		res *sql.Rows
	)
	if bib != nil {
		res, err = db.QueryContext(
			ctx,
			"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
				"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
				"gender_ranking, finish, result_type FROM result NATURAL JOIN person "+
				"WHERE event_year_id=? AND bib=? ORDER BY seconds DESC;",
			eventYearID,
			bib,
		)
	} else if distance != "" {
		res, err = db.QueryContext(
			ctx,
			"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
				"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
				"gender_ranking, finish, result_type FROM result NATURAL JOIN person "+
				"WHERE event_year_id=? AND bib=? AND distance=? ORDER BY seconds DESC;",
			eventYearID,
			bib,
			distance,
		)
	} else if all {
		res, err = db.QueryContext(
			ctx,
			"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
				"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
				"gender_ranking, finish, result_type FROM result NATURAL JOIN person "+
				"WHERE event_year_id=? ORDER BY seconds DESC;",
			eventYearID,
		)
	} else {
		res, err = db.QueryContext(
			ctx,
			"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
				"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
				"gender_ranking, finish, result_type FROM result r NATURAL JOIN person p "+
				"JOIN (SELECT bib AS mx_bib, event_year_id AS mx_event_year_id, MAX(occurence) as mx_occurence "+
				"FROM result NATURAL JOIN person GROUP BY bib, event_year_id) b "+
				"ON b.mx_bib=p.bib AND b.mx_event_year_id=p.event_year_id "+
				"AND b.mx_occurence=r.occurence WHERE event_year_id=? ORDER BY seconds DESC;",
			eventYearID,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("error retrieving results: %v", err)
	}
	defer res.Close()
	var outResults []types.Result
	for res.Next() {
		var result types.Result
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
		)
		if err != nil {
			return nil, fmt.Errorf("error getting result: %v", err)
		}
		outResults = append(outResults, result)
	}
	return outResults, nil
}

// GetResults Gets results for an event year.
func (m *MySQL) GetResults(eventYearID int64) ([]types.Result, error) {
	return m.getResultsInternal(eventYearID, nil, true, "")
}

func (m *MySQL) GetDistanceResults(eventYearID int64, distance string) ([]types.Result, error) {
	return m.getResultsInternal(eventYearID, nil, false, distance)
}

// GetLastResults Gets the last result for a bib in an event year.
func (m *MySQL) GetLastResults(eventYearID int64) ([]types.Result, error) {
	return m.getResultsInternal(eventYearID, nil, false, "")
}

// GetBibResults Gets results for an event year of a specific individual specified by their bib.
func (m *MySQL) GetBibResults(eventYearID int64, bib string) ([]types.Result, error) {
	return m.getResultsInternal(eventYearID, &bib, false, "")
}

// DeleteResults Deletes results from the database.
func (m *MySQL) DeleteResults(eventYearID int64, results []types.Result) error {
	db, err := m.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	stmt, err := db.PrepareContext(
		ctx,
		"DELETE r FROM result AS r WHERE location=? AND occurence=? AND EXISTS (SELECT * FROM person AS p WHERE event_year_id=? AND bib=? AND p.person_id=r.person_id);",
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
			return fmt.Errorf("error executing prepared delete statement: %v", err)
		}
	}
	return nil
}

// DeleteEventResults Deletes results for an event year.
func (m *MySQL) DeleteEventResults(eventYearID int64) (int64, error) {
	db, err := m.GetDB()
	if err != nil {
		return 0, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"DELETE r FROM result AS r WHERE EXISTS (SELECT * FROM person AS p WHERE p.event_year_id=? AND p.person_id=r.person_id);",
		eventYearID,
	)
	if err != nil {
		return 0, fmt.Errorf("unable to delete results for event year: %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error fetching rows affected from event year results deletion: %v", err)
	}
	return count, nil
}

// AddResults Adds results to the database. (Also updates.)
func (m *MySQL) AddResults(eventYearID int64, results []types.Result) ([]types.Result, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	personstmt, err := db.PrepareContext(
		ctx,
		"INSERT INTO person("+
			"event_year_id, "+
			"bib, "+
			"first, "+
			"last, "+
			"age, "+
			"gender, "+
			"age_group, "+
			"distance"+
			")"+
			" VALUES (?,?,?,?,?,?,?,?) "+
			"ON DUPLICATE KEY UPDATE "+
			"first=VALUES(first), "+
			"last=VALUES(last), "+
			"age=VALUES(age), "+
			"gender=VALUES(gender), "+
			"age_group=VALUES(age_group), "+
			"distance=VALUES(distance);",
	)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare statement for person add: %v", err)
	}
	defer personstmt.Close()
	stmt, err := db.PrepareContext(
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
			" VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?) "+
			"ON DUPLICATE KEY UPDATE "+
			"seconds=VALUES(seconds), "+
			"milliseconds=VALUES(milliseconds), "+
			"chip_seconds=VALUES(chip_seconds), "+
			"chip_milliseconds=VALUES(chip_milliseconds), "+
			"segment=VALUES(segment), "+
			"ranking=VALUES(ranking), "+
			"age_ranking=VALUES(age_ranking), "+
			"gender_ranking=VALUES(gender_ranking), "+
			"finish=VALUES(finish), "+
			"result_type=VALUES(result_type);",
	)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare statement for result add: %v", err)
	}
	defer stmt.Close()
	// Add people to the database and get their ids.
	for _, result := range results {
		_, err = personstmt.QueryContext(
			ctx,
			eventYearID,
			result.Bib,
			result.First,
			result.Last,
			result.Age,
			result.Gender,
			result.AgeGroup,
			result.Distance,
		)
		if err != nil {
			return nil, fmt.Errorf("error adding person to database: %v", err)
		}
	}
	// Get person ids from the database.
	bibMap := make(map[string]int64)
	res, err := db.QueryContext(
		ctx,
		"SELECT person_id, bib FROM person WHERE event_year_id=?;",
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
	var outResults []types.Result
	for _, result := range results {
		if _, ok := bibMap[result.Bib]; !ok {
			return outResults, errors.New("something went wrong and bib not found")
		}
		_, err = stmt.ExecContext(
			ctx,
			bibMap[result.Bib],
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
			return outResults, fmt.Errorf("error adding result to database: %v", err)
		}
		outResults = append(outResults, result)
	}
	return outResults, nil
}
