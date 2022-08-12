package sqlite

import (
	"chronokeep/results/types"
	"context"
	"fmt"
	"time"
)

func (s *SQLite) GetPerson(slug, year, bib string) (*types.Person, error) {
	db, err := s.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT person_id, bib, first, last, age, gender, age_group, distance FROM person NATURAL JOIN event_year NATURAL JOIN event WHERE slug=? AND year=? AND bib=?",
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

func (s *SQLite) GetPeople(slug, year string) ([]types.Person, error) {
	db, err := s.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT person_id, bib, first, last, age, gender, age_group, distance FROM person NATURAL JOIN event_year NATURAL JOIN event WHERE slug=? AND year=?;",
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

func (s *SQLite) AddPerson(eventYearID int64, person types.Person) (*types.Person, error) {
	db, err := s.GetDB()
	if err != nil {
		return nil, err
	}
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction: %v", err)
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	_, err = tx.ExecContext(
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
			"distance=$8;",
		eventYearID,
		person.Bib,
		person.First,
		person.Last,
		person.Age,
		person.Gender,
		person.AgeGroup,
		person.Distance,
	)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error adding person to database: %v", err)
	}
	res, err := tx.QueryContext(
		ctx,
		"SELECT person_id FROM person WHERE event_year_id=$1 AND bib=$2;",
		eventYearID,
		person.Bib,
	)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error retrieving person_id: %v", err)
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
	if res.Next() {
		res.Scan(
			&output.Identifier,
		)
	} else {
		tx.Rollback()
		return nil, fmt.Errorf("person not found after add: %v", err)
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}
	return &output, nil
}

func (s *SQLite) AddPeople(eventYearID int64, people []types.Person) ([]types.Person, error) {
	db, err := s.GetDB()
	if err != nil {
		return nil, err
	}
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction: %v", err)
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
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
			"distance"+
			") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) "+
			"ON CONFLICT (event_year_id, bib) DO UPDATE SET "+
			"first=$3, "+
			"last=$4, "+
			"age=$5, "+
			"gender=$6, "+
			"age_group=$7, "+
			"distance=$8;",
	)
	if err != nil {
		return nil, fmt.Errorf("error preparing statement for adding people: %v", err)
	}
	for _, person := range people {
		_, err = personstmt.ExecContext(
			ctx,
			eventYearID,
			person.Bib,
			person.First,
			person.Last,
			person.Age,
			person.Gender,
			person.AgeGroup,
			person.Distance,
		)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error adding person to database: %v", err)
		}
	}
	// get a list of all the bibs and persons for the event year
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
	output := make([]types.Person, 0)
	for _, person := range people {
		if id, ok := bibMap[person.Bib]; ok {
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
		} else {
			tx.Rollback()
			return nil, fmt.Errorf("person not found after add: %v", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}
	return output, nil
}

func (s *SQLite) DeletePeople(eventYearId int64, bibs []string) error {
	db, err := s.GetDB()
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	if len(bibs) > 0 {
		for _, bib := range bibs {
			_, err = tx.ExecContext(
				ctx,
				"DELETE FROM person WHERE event_year_id=$1 AND bib=$2;",
				eventYearId,
				bib,
			)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("error deleting participant: %v", err)
			}
		}
	} else {
		_, err = tx.ExecContext(
			ctx,
			"DELETE FROM person WHERE event_year_id=$1;",
			eventYearId,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error deleting participants: %v", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error committing transaction: %v", err)
	}
	return nil
}
