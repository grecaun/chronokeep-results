package postgres

import (
	"chronokeep/results/types"
	"context"
	"fmt"
	"time"
)

func (p *Postgres) GetPerson(slug, year, bib string) (*types.Person, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Query(
		ctx,
		"SELECT person_id, bib, first, last, age, gender, age_group, distance FROM person NATURAL JOIN event_year NATURAL JOIN event WHERE slug=$1 AND year=$2 AND bib=$3",
		slug,
		year,
		bib,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving person: %v", err)
	}
	defer res.Close()
	if res.Next() {
		var outPerson types.Person
		err = res.Scan(
			&outPerson.Identifier,
			&outPerson.Bib,
			&outPerson.First,
			&outPerson.Last,
			&outPerson.Age,
			&outPerson.Gender,
			&outPerson.AgeGroup,
			&outPerson.Distance,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting person: %v", err)
		}
		return &outPerson, nil
	}
	return nil, nil
}

func (p *Postgres) GetPeople(slug, year string) ([]types.Person, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Query(
		ctx,
		"SELECT person_id, bib, first, last, age, gender, age_group, distance FROM person NATURAL JOIN event_year NATURAL JOIN event WHERE slug=$1 AND year=$2;",
		slug,
		year,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving people: %v", err)
	}
	defer res.Close()
	output := make([]types.Person, 0)
	for res.Next() {
		var person types.Person
		err := res.Scan(
			&person.Identifier,
			&person.Bib,
			&person.First,
			&person.Last,
			&person.Age,
			&person.Gender,
			&person.AgeGroup,
			&person.Distance,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting person: %v", err)
		}
		output = append(output, person)
	}
	return output, nil
}

func (p *Postgres) AddPerson(eventYearID int64, person types.Person) (*types.Person, error) {
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
	output := types.Person{
		Bib:      person.Bib,
		First:    person.First,
		Last:     person.Last,
		Age:      person.Age,
		Gender:   person.Gender,
		AgeGroup: person.AgeGroup,
		Distance: person.Distance,
	}
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
		person.Bib,
		person.First,
		person.Last,
		person.Age,
		person.Gender,
		person.AgeGroup,
		person.Distance,
	).Scan(&output.Identifier)
	if err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("error adding person to database: %v", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}
	return &output, nil
}

func (p *Postgres) AddPeople(eventYearID int64, people []types.Person) ([]types.Person, error) {
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
	output := make([]types.Person, 0)
	for _, person := range people {
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
			person.Bib,
			person.First,
			person.Last,
			person.Age,
			person.Gender,
			person.AgeGroup,
			person.Distance,
		).Scan(&id)
		if err != nil {
			tx.Rollback(ctx)
			return nil, fmt.Errorf("error adding person to database: %v", err)
		}
		output = append(output, types.Person{
			Identifier: id,
			Bib:        person.Bib,
			First:      person.First,
			Last:       person.Last,
			Age:        person.Age,
			Gender:     person.Gender,
			AgeGroup:   person.AgeGroup,
			Distance:   person.Distance,
		})
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}
	return output, nil
}

func (p *Postgres) DeletePeople(eventYearId int64, bibs []string) error {
	db, err := p.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	if len(bibs) > 0 {
		for _, bib := range bibs {
			_, err = tx.Exec(
				ctx,
				"DELETE FROM person WHERE event_year_id=$1 AND bib=$2;",
				eventYearId,
				bib,
			)
			if err != nil {
				tx.Rollback(ctx)
				return fmt.Errorf("error deleting participant: %v", err)
			}
		}
	} else {
		_, err = tx.Exec(
			ctx,
			"DELETE FROM person WHERE event_year_id=$1;",
			eventYearId,
		)
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("error deleting participants: %v", err)
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("error committing transaction: %v", err)
	}
	return nil
}
