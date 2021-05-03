package database

import (
	"chronokeep/results/types"
	"context"
	"fmt"
	"time"
)

// GetResults Gets results for an event year.
func GetResults(eventYearID int64) ([]types.Result, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, segment, location, occurence, ranking, age_ranking, gender_ranking, finish FROM result WHERE event_year_id=?;",
		eventYearID,
	)
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
			&result.Segment,
			&result.Location,
			&result.Occurence,
			&result.Ranking,
			&result.AgeRanking,
			&result.GenderRanking,
			&result.Finish,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting result: %v", err)
		}
		outResults = append(outResults, result)
	}
	return outResults, nil
}

// DeleteResults Deletes results from the database.
func DeleteResults(eventYearID int64, results []types.Result) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	stmt, err := db.PrepareContext(
		ctx,
		"DELETE FROM result WHERE event_year_id=? AND bib=? AND location=? AND occurence=?;",
	)
	if err != nil {
		return fmt.Errorf("unable to get prepared statement for result deletion: %v", err)
	}
	defer stmt.Close()
	for _, result := range results {
		_, err := stmt.ExecContext(
			ctx,
			eventYearID,
			result.Bib,
			result.Location,
			result.Occurence,
		)
		if err != nil {
			return fmt.Errorf("error executing prepared delete statement: %v", err)
		}
	}
	return nil
}

// DeleteEventResults Deletes results for an event year.
func DeleteEventResults(eventYearID int64) (int64, error) {
	db, err := GetDB()
	if err != nil {
		return 0, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"DELETE FROM result WHERE event_year_id=?;",
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
func AddResults(eventYearID int64, results []types.Result) ([]types.Result, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	stmt, err := db.PrepareContext(
		ctx,
		"INSERT INTO result("+
			"event_year_id, "+
			"bib, "+
			"first, "+
			"last, "+
			"age, "+
			"gender, "+
			"age_group, "+
			"distance, "+
			"seconds, "+
			"milliseconds, "+
			"segment, "+
			"location, "+
			"occurence, "+
			"ranking, "+
			"age_ranking, "+
			"gender_ranking, "+
			"finish) "+
			" VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) "+
			"ON DUPLICATE KEY UPDATE "+
			"first=VALUES(first), "+
			"last=VALUES(last), "+
			"age=VALUES(age), "+
			"gender=VALUES(gender), "+
			"age_group=VALUES(age_group), "+
			"distance=VALUES(distance), "+
			"seconds=VALUES(seconds), "+
			"milliseconds=VALUES(milliseconds), "+
			"segment=VALUES(segment), "+
			"ranking=VALUES(ranking), "+
			"age_ranking=VALUES(age_ranking), "+
			"gender_ranking=VALUES(gender_ranking), "+
			"finish=VALUES(finish);",
	)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare statement for result add: %v", err)
	}
	defer stmt.Close()
	var outResults []types.Result
	for _, result := range results {
		_, err := stmt.ExecContext(
			ctx,
			eventYearID,
			result.Bib,
			result.First,
			result.Last,
			result.Age,
			result.Gender,
			result.AgeGroup,
			result.Distance,
			result.Seconds,
			result.Milliseconds,
			result.Segment,
			result.Location,
			result.Occurence,
			result.Ranking,
			result.AgeRanking,
			result.GenderRanking,
			result.Finish,
		)
		if err != nil {
			return outResults, fmt.Errorf("error adding result to database: %v", err)
		}
		outResults = append(outResults, result)
	}
	return outResults, nil
}
