package postgres

import (
	"chronokeep/results/types"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
)

func (p *Postgres) getResultsInternal(eventYearID int64, bib *string) ([]types.Result, error) {
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
		res, err = db.Query(
			ctx,
			"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, gender_ranking, finish FROM result WHERE event_year_id=$1 AND bib=$2;",
			eventYearID,
			bib,
		)
	} else {
		res, err = db.Query(
			ctx,
			"SELECT bib, first, last, age, gender, age_group, distance, seconds, milliseconds, chip_seconds, chip_milliseconds, segment, location, occurence, ranking, age_ranking, gender_ranking, finish FROM result WHERE event_year_id=$1;",
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
		)
		if err != nil {
			return nil, fmt.Errorf("error getting result: %v", err)
		}
		outResults = append(outResults, result)
	}
	return outResults, nil
}

// GetResults Gets results for an event year.
func (p *Postgres) GetResults(eventYearID int64) ([]types.Result, error) {
	return p.getResultsInternal(eventYearID, nil)
}

// GetBibResults Gets results for an event year of a specific individual specified by their bib.
func (p *Postgres) GetBibResults(eventYearID int64, bib string) ([]types.Result, error) {
	return p.getResultsInternal(eventYearID, &bib)
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
			"DELETE FROM result WHERE event_year_id=$1 AND bib=$2 AND location=$3 AND occurence=$4;",
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
		"DELETE FROM result WHERE event_year_id=$1;",
		eventYearID,
	)
	if err != nil {
		return 0, fmt.Errorf("unable to delete results for event year: %v", err)
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
		_, err = tx.Exec(
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
				"chip_seconds, "+
				"chip_milliseconds, "+
				"segment, "+
				"location, "+
				"occurence, "+
				"ranking, "+
				"age_ranking, "+
				"gender_ranking, "+
				"finish) "+
				" VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17) "+
				"ON CONFLICT (event_year_id, bib, location, occurence) DO UPDATE SET "+
				"first=$3, "+
				"last=$4, "+
				"age=$5, "+
				"gender=$6, "+
				"age_group=$7, "+
				"distance=$8, "+
				"seconds=$9, "+
				"milliseconds=$10, "+
				"seconds=$11, "+
				"milliseconds=$12, "+
				"segment=$13, "+
				"ranking=$16, "+
				"age_ranking=$17, "+
				"gender_ranking=$18, "+
				"finish=$19;",
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
			result.ChipSeconds,
			result.ChipMilliseconds,
			result.Segment,
			result.Location,
			result.Occurence,
			result.Ranking,
			result.AgeRanking,
			result.GenderRanking,
			result.Finish,
		)
		if err != nil {
			return nil, err
		}
	}
	return results, tx.Commit(ctx)
}
