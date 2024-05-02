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
		"SELECT person_id, bib, first, last, age, gender, age_group, distance, chip, anonymous, alternate_id, sms_enabled "+
			"FROM person NATURAL JOIN event_year NATURAL JOIN event WHERE slug=$1 AND year=$2 AND bib=$3",
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
		var anonymous int
		var sms int
		err = res.Scan(
			&outPerson.Identifier,
			&outPerson.Bib,
			&outPerson.First,
			&outPerson.Last,
			&outPerson.Age,
			&outPerson.Gender,
			&outPerson.AgeGroup,
			&outPerson.Distance,
			&outPerson.Chip,
			&anonymous,
			&outPerson.AlternateId,
			&sms,
		)
		outPerson.Anonymous = anonymous != 0
		outPerson.SMSEnabled = sms != 0
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
		"SELECT person_id, bib, first, last, age, gender, age_group, distance, chip, anonymous, alternate_id, sms_enabled "+
			"FROM person NATURAL JOIN event_year NATURAL JOIN event WHERE slug=$1 AND year=$2;",
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
		var anonymous int
		var sms int
		err := res.Scan(
			&person.Identifier,
			&person.Bib,
			&person.First,
			&person.Last,
			&person.Age,
			&person.Gender,
			&person.AgeGroup,
			&person.Distance,
			&person.Chip,
			&anonymous,
			&person.AlternateId,
			&sms,
		)
		person.Anonymous = anonymous != 0
		person.SMSEnabled = sms != 0
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
		Bib:         person.Bib,
		First:       person.First,
		Last:        person.Last,
		Age:         person.Age,
		Gender:      person.Gender,
		AgeGroup:    person.AgeGroup,
		Distance:    person.Distance,
		Chip:        person.Chip,
		Anonymous:   person.Anonymous,
		AlternateId: person.AlternateId,
		SMSEnabled:  person.SMSEnabled,
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
			"distance, "+
			"chip, "+
			"anonymous, "+
			"alternate_id, "+
			"sms_enabled"+
			") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) "+
			"ON CONFLICT (event_year_id, alternate_id) DO UPDATE SET "+
			"first=$3, "+
			"last=$4, "+
			"age=$5, "+
			"gender=$6, "+
			"age_group=$7, "+
			"distance=$8, "+
			"chip=$9, "+
			"anonymous=$10, "+
			"alternate_id=$11, "+
			"sms_enabled=$12 "+
			"RETURNING (person_id);",
		eventYearID,
		person.Bib,
		person.First,
		person.Last,
		person.Age,
		person.Gender,
		person.AgeGroup,
		person.Distance,
		person.Chip,
		person.AnonyInt(),
		person.AlternateId,
		person.SMSInt(),
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
				"distance, "+
				"chip, "+
				"anonymous, "+
				"alternate_id, "+
				"sms_enabled"+
				") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) "+
				"ON CONFLICT (event_year_id, alternate_id) DO UPDATE SET "+
				"first=$3, "+
				"last=$4, "+
				"age=$5, "+
				"gender=$6, "+
				"age_group=$7, "+
				"distance=$8, "+
				"chip=$9, "+
				"anonymous=$10, "+
				"alternate_id=$11, "+
				"sms_enabled=$12 "+
				"RETURNING (person_id);",
			eventYearID,
			person.Bib,
			person.First,
			person.Last,
			person.Age,
			person.Gender,
			person.AgeGroup,
			person.Distance,
			person.Chip,
			person.AnonyInt(),
			person.AlternateId,
			person.SMSInt(),
		).Scan(&id)
		if err != nil {
			tx.Rollback(ctx)
			return nil, fmt.Errorf("error adding person to database: %v", err)
		}
		output = append(output, types.Person{
			Identifier:  id,
			Bib:         person.Bib,
			First:       person.First,
			Last:        person.Last,
			Age:         person.Age,
			Gender:      person.Gender,
			AgeGroup:    person.AgeGroup,
			Distance:    person.Distance,
			Chip:        person.Chip,
			Anonymous:   person.Anonymous,
			AlternateId: person.AlternateId,
			SMSEnabled:  person.SMSEnabled,
		})
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}
	return output, nil
}

func (p *Postgres) DeletePeople(eventYearId int64, alternateIds []string) (int64, error) {
	db, err := p.GetDB()
	if err != nil {
		return 0, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("error starting transaction: %v", err)
	}
	var count int64
	if len(alternateIds) > 0 {
		count = 0
		for _, altId := range alternateIds {
			count++
			_, err = tx.Exec(
				ctx,
				"DELETE FROM person WHERE event_year_id=$1 AND alternate_id=$2;",
				eventYearId,
				altId,
			)
			if err != nil {
				tx.Rollback(ctx)
				return 0, fmt.Errorf("error deleting participant: %v", err)
			}
		}
	} else {
		res, err := tx.Exec(
			ctx,
			"DELETE FROM person WHERE event_year_id=$1;",
			eventYearId,
		)
		if err != nil {
			tx.Rollback(ctx)
			return 0, fmt.Errorf("error deleting participants: %v", err)
		}
		count = res.RowsAffected()
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("error committing transaction: %v", err)
	}
	return count, nil
}

func (p *Postgres) UpdatePerson(eventYearID int64, person types.Person) (*types.Person, error) {
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
		Bib:         person.Bib,
		First:       person.First,
		Last:        person.Last,
		Age:         person.Age,
		Gender:      person.Gender,
		AgeGroup:    person.AgeGroup,
		Distance:    person.Distance,
		Chip:        person.Chip,
		Anonymous:   person.Anonymous,
		AlternateId: person.AlternateId,
		SMSEnabled:  person.SMSEnabled,
	}
	err = tx.QueryRow(
		ctx,
		"UPDATE person SET "+
			"bib=$1, "+
			"first=$2, "+
			"last=$3, "+
			"age=$4, "+
			"gender=$5, "+
			"age_group=$6, "+
			"distance=$7, "+
			"chip=$8, "+
			"anonymous=$9, "+
			"sms_enabled=$12 "+
			"WHERE event_year_id=$10 AND alternate_id=$11 RETURNING (person_id);",
		person.Bib,
		person.First,
		person.Last,
		person.Age,
		person.Gender,
		person.AgeGroup,
		person.Distance,
		person.Chip,
		person.AnonyInt(),
		eventYearID,
		person.AlternateId,
		person.SMSInt(),
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
