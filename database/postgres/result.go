package postgres

import (
	"chronokeep/results/types"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
)

type ResultType int

const (
	All ResultType = iota
	Finish
	Last
)

func (p *Postgres) getResultsInternal(eventYearID int64, bib *string, rtype ResultType, distance string, limit, page int) ([]types.Result, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	var (
		res pgx.Rows
	)
	if bib != nil {
		if limit > 0 {
			res, err = db.Query(
				ctx,
				"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
					"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
					"gender_ranking, finish, result_type FROM result NATURAL JOIN person WHERE event_year_id=$1 AND bib=$2 ORDER BY seconds ASC LIMIT $3 OFFSET $4;",
				eventYearID,
				bib,
				limit,
				page*limit,
			)
		} else {
			res, err = db.Query(
				ctx,
				"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
					"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
					"gender_ranking, finish, result_type FROM result NATURAL JOIN person WHERE event_year_id=$1 AND bib=$2 ORDER BY seconds ASC;",
				eventYearID,
				bib,
			)
		}
	} else if distance != "" {
		if rtype == Finish {
			if limit > 0 {
				res, err = db.Query(
					ctx,
					"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
						"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
						"gender_ranking, finish, result_type FROM result NATURAL JOIN person WHERE finish=TRUE AND event_year_id=$1 AND distance=$2 ORDER BY seconds ASC LIMIT $3 OFFSET $4;",
					eventYearID,
					distance,
					limit,
					page*limit,
				)
			} else {
				res, err = db.Query(
					ctx,
					"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
						"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
						"gender_ranking, finish, result_type FROM result NATURAL JOIN person WHERE finish=TRUE AND event_year_id=$1 AND distance=$2 ORDER BY seconds ASC;",
					eventYearID,
					distance,
				)
			}
		} else {
			if limit > 0 {
				res, err = db.Query(
					ctx,
					"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
						"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
						"gender_ranking, finish, result_type FROM result r NATURAL JOIN person p "+
						"JOIN (SELECT bib AS mx_bib, event_year_id AS mx_event_year_id, MAX(occurence) as mx_occurence "+
						"FROM result NATURAL JOIN person GROUP BY bib, event_year_id) b ON b.mx_bib=p.bib AND b.mx_event_year_id=p.event_year_id AND b.mx_occurence=r.occurence "+
						"WHERE event_year_id=$1 AND distance=$2 ORDER BY seconds ASC LIMIT $3 OFFSET $4;",
					eventYearID,
					distance,
					limit,
					page*limit,
				)
			} else {
				res, err = db.Query(
					ctx,
					"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
						"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
						"gender_ranking, finish, result_type FROM result r NATURAL JOIN person p "+
						"JOIN (SELECT bib AS mx_bib, event_year_id AS mx_event_year_id, MAX(occurence) as mx_occurence "+
						"FROM result NATURAL JOIN person GROUP BY bib, event_year_id) b ON b.mx_bib=p.bib AND b.mx_event_year_id=p.event_year_id AND b.mx_occurence=r.occurence "+
						"WHERE event_year_id=$1 AND distance=$2 ORDER BY seconds ASC;",
					eventYearID,
					distance,
				)
			}
		}
	} else if rtype == All {
		if limit > 0 {
			res, err = db.Query(
				ctx,
				"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
					"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
					"gender_ranking, finish, result_type FROM result NATURAL JOIN person WHERE event_year_id=$1 ORDER BY seconds ASC LIMIT $2 OFFSET $3;",
				eventYearID,
				limit,
				page*limit,
			)
		} else {
			res, err = db.Query(
				ctx,
				"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
					"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
					"gender_ranking, finish, result_type FROM result NATURAL JOIN person WHERE event_year_id=$1 ORDER BY seconds ASC;",
				eventYearID,
			)
		}
	} else if rtype == Finish {
		if limit > 0 {
			res, err = db.Query(
				ctx,
				"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
					"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
					"gender_ranking, finish, result_type FROM result NATURAL JOIN person "+
					"WHERE finish=TRUE AND event_year_id=$1 ORDER BY seconds ASC LIMIT $2 OFFSET $3;",
				eventYearID,
				limit,
				page*limit,
			)
		} else {
			res, err = db.Query(
				ctx,
				"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
					"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
					"gender_ranking, finish, result_type FROM result NATURAL JOIN person "+
					"WHERE finish=TRUE AND event_year_id=$1 ORDER BY seconds ASC;",
				eventYearID,
			)
		}
	} else {
		if limit > 0 {
			res, err = db.Query(
				ctx,
				"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
					"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
					"gender_ranking, finish, result_type FROM result r NATURAL JOIN person p "+
					"JOIN (SELECT bib AS mx_bib, event_year_id AS mx_event_year_id, MAX(occurence) as mx_occurence "+
					"FROM result NATURAL JOIN person GROUP BY bib, event_year_id) b ON b.mx_bib=p.bib AND b.mx_event_year_id=p.event_year_id "+
					"AND b.mx_occurence=r.occurence WHERE event_year_id=$1 ORDER BY seconds ASC LIMIT $2 OFFSET $3;",
				eventYearID,
				limit,
				page*limit,
			)
		} else {
			res, err = db.Query(
				ctx,
				"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, "+
					"chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, "+
					"gender_ranking, finish, result_type FROM result r NATURAL JOIN person p "+
					"JOIN (SELECT bib AS mx_bib, event_year_id AS mx_event_year_id, MAX(occurence) as mx_occurence "+
					"FROM result NATURAL JOIN person GROUP BY bib, event_year_id) b ON b.mx_bib=p.bib AND b.mx_event_year_id=p.event_year_id "+
					"AND b.mx_occurence=r.occurence WHERE event_year_id=$1 ORDER BY seconds ASC;",
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
func (p *Postgres) GetResults(eventYearID int64, limit, page int) ([]types.Result, error) {
	return p.getResultsInternal(eventYearID, nil, All, "", limit, page)
}

// GetLastResults Gets the last result for a bib in an event year.
func (p *Postgres) GetLastResults(eventYearID int64, limit, page int) ([]types.Result, error) {
	return p.getResultsInternal(eventYearID, nil, Last, "", limit, page)
}

// GetDistanceResults Gets the distance results (last only) for a distance.
func (m *Postgres) GetDistanceResults(eventYearID int64, distance string, limit, page int) ([]types.Result, error) {
	return m.getResultsInternal(eventYearID, nil, Last, distance, limit, page)
}

// GetFinishResults Gets the finish results for an entire event (empty distance) or just a distance.
func (m *Postgres) GetFinishResults(eventYearID int64, distance string, limit, page int) ([]types.Result, error) {
	return m.getResultsInternal(eventYearID, nil, Finish, distance, limit, page)
}

// GetBibResults Gets results for an event year of a specific individual specified by their bib.
func (p *Postgres) GetBibResults(eventYearID int64, bib string) ([]types.Result, error) {
	return p.getResultsInternal(eventYearID, &bib, All, "", 0, 0)
}

// DeleteResults Deletes results from the database.
func (p *Postgres) DeleteResults(eventYearID int64, results []types.Result) error {
	db, err := p.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("unable to begin transaction to delete results: %v", err)
	}
	for _, result := range results {
		_, err = tx.Exec(
			ctx,
			"DELETE FROM result r WHERE location=$3 AND occurence=$4 AND EXISTS (SELECT * FROM person p WHERE event_year_id=$1 AND bib=$2 AND r.person_id=p.person_id);",
			eventYearID,
			result.Bib,
			result.Location,
			result.Occurence,
		)
		if err != nil {
			return fmt.Errorf("error executing delete query: %v", err)
		}
	}
	return tx.Commit(ctx)
}

// DeleteEventResults Deletes results for an event year.
func (p *Postgres) DeleteEventResults(eventYearID int64) (int64, error) {
	db, err := p.GetDB()
	if err != nil {
		return 0, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Exec(
		ctx,
		"DELETE FROM result r WHERE EXISTS (SELECT * FROM person p WHERE event_year_id=$1 AND p.person_id=r.person_id);",
		eventYearID,
	)
	if err != nil {
		return 0, fmt.Errorf("unable to delete results for event year: %v", err)
	}
	_, err = db.Exec(
		ctx,
		"DELETE FROM person WHERE event_year_id=$1;",
		eventYearID,
	)
	if err != nil {
		return 0, fmt.Errorf("unable to delete persons for event year: %v", err)
	}
	return res.RowsAffected(), nil
}

// AddResults Adds results to the database. (Also updates.)
func (p *Postgres) AddResults(eventYearID int64, results []types.Result) ([]types.Result, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to begin transaction to add results: %v", err)
	}
	for _, result := range results {
		var id int64
		err = tx.QueryRow(
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
				") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) "+
				"ON CONFLICT (event_year_id, bib) DO UPDATE SET "+
				"first=$3, "+
				"last=$4, "+
				"age=$5, "+
				"gender=$6, "+
				"age_group=$7, "+
				"distance=$8 RETURNING (person_id);",
			eventYearID,
			result.Bib,
			result.First,
			result.Last,
			result.Age,
			result.Gender,
			result.AgeGroup,
			result.Distance,
		).Scan(&id)
		if err != nil {
			return nil, err
		}
		if id == 0 {
			return nil, errors.New("id value set to 0")
		}
		_, err = tx.Exec(
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
			return nil, err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}
	return results, nil
}
