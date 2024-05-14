package mysql

import (
	"chronokeep/results/types"
	"context"
	"fmt"
	"time"
)

func (m *MySQL) GetPerson(slug, year, bib string) (*types.Person, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT person_id, bib, first, last, age, gender, age_group, distance, anonymous, alternate_id "+
			"FROM person NATURAL JOIN event_year NATURAL JOIN event WHERE slug=? AND year=? AND bib=?",
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
		err = res.Scan(
			&outPerson.Identifier,
			&outPerson.Bib,
			&outPerson.First,
			&outPerson.Last,
			&outPerson.Age,
			&outPerson.Gender,
			&outPerson.AgeGroup,
			&outPerson.Distance,
			&anonymous,
			&outPerson.AlternateId,
		)
		outPerson.Anonymous = anonymous != 0
		if err != nil {
			return nil, fmt.Errorf("error getting person: %v", err)
		}
		return &outPerson, nil
	}
	return nil, nil
}

func (m *MySQL) GetPeople(slug, year string) ([]types.Person, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT person_id, bib, first, last, age, gender, age_group, distance, anonymous, alternate_id "+
			"FROM person NATURAL JOIN event_year NATURAL JOIN event WHERE slug=? AND year=?;",
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
		err := res.Scan(
			&person.Identifier,
			&person.Bib,
			&person.First,
			&person.Last,
			&person.Age,
			&person.Gender,
			&person.AgeGroup,
			&person.Distance,
			&anonymous,
			&person.AlternateId,
		)
		person.Anonymous = anonymous != 0
		if err != nil {
			return nil, fmt.Errorf("error getting person: %v", err)
		}
		output = append(output, person)
	}
	return output, nil
}

func (m *MySQL) AddPerson(eventYearID int64, person types.Person) (*types.Person, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction: %v", err)
	}
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
			"distance, "+
			"anonymous, "+
			"alternate_id "+
			")"+
			" VALUES (?,?,?,?,?,?,?,?,?,?) "+
			"ON DUPLICATE KEY UPDATE "+
			"first=VALUES(first), "+
			"last=VALUES(last), "+
			"age=VALUES(age), "+
			"gender=VALUES(gender), "+
			"age_group=VALUES(age_group), "+
			"distance=VALUES(distance), "+
			"anonymous=VALUES(anonymous), "+
			"alternate_id=VALUES(alternate_id);",
		eventYearID,
		person.Bib,
		person.First,
		person.Last,
		person.Age,
		person.Gender,
		person.AgeGroup,
		person.Distance,
		person.AnonyInt(),
		person.AlternateId,
	)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error adding person to database: %v", err)
	}
	res, err := tx.QueryContext(
		ctx,
		"SELECT person_id FROM person WHERE event_year_id=? AND alternate_id=?;",
		eventYearID,
		person.AlternateId,
	)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error retrieving person_id: %v", err)
	}
	defer res.Close()
	output := types.Person{
		Bib:         person.Bib,
		First:       person.First,
		Last:        person.Last,
		Age:         person.Age,
		Gender:      person.Gender,
		AgeGroup:    person.AgeGroup,
		Distance:    person.Distance,
		Anonymous:   person.Anonymous,
		AlternateId: person.AlternateId,
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

func (m *MySQL) AddPeople(eventYearID int64, people []types.Person) ([]types.Person, error) {
	db, err := m.GetDB()
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
			"distance, "+
			"anonymous, "+
			"alternate_id "+
			")"+
			" VALUES (?,?,?,?,?,?,?,?,?,?) "+
			"ON DUPLICATE KEY UPDATE "+
			"first=VALUES(first), "+
			"last=VALUES(last), "+
			"age=VALUES(age), "+
			"gender=VALUES(gender), "+
			"age_group=VALUES(age_group), "+
			"distance=VALUES(distance),"+
			"anonymous=VALUES(anonymous), "+
			"alternate_id=VALUES(alternate_id);",
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
			person.AnonyInt(),
			person.AlternateId,
		)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error adding person to database: %v", err)
		}
	}
	// get a list of all the bibs and persons for the event year
	altMap := make(map[string]int64)
	res, err := tx.QueryContext(
		ctx,
		"SELECT person_id, alternate_id FROM person WHERE event_year_id=?;",
		eventYearID,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying database for person ids: %v", err)
	}
	defer res.Close()
	for res.Next() {
		var id int64
		var altId string
		err = res.Scan(
			&id,
			&altId,
		)
		if err != nil {
			return nil, fmt.Errorf("error retrieving person ids: %v", err)
		}
		altMap[altId] = id
	}
	output := make([]types.Person, 0)
	for _, person := range people {
		if id, ok := altMap[person.AlternateId]; ok {
			output = append(output, types.Person{
				Identifier:  id,
				Bib:         person.Bib,
				First:       person.First,
				Last:        person.Last,
				Age:         person.Age,
				Gender:      person.Gender,
				AgeGroup:    person.AgeGroup,
				Distance:    person.Distance,
				Anonymous:   person.Anonymous,
				AlternateId: person.AlternateId,
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

func (m *MySQL) DeletePeople(eventYearId int64, alternateIds []string) (int64, error) {
	db, err := m.GetDB()
	if err != nil {
		return 0, err
	}
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("error starting transaction: %v", err)
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	var count int64
	if len(alternateIds) > 0 {
		count = 0
		for _, altId := range alternateIds {
			count++
			_, err = tx.ExecContext(
				ctx,
				"DELETE FROM person WHERE event_year_id=? AND alternate_id=?;",
				eventYearId,
				altId,
			)
			if err != nil {
				tx.Rollback()
				return 0, fmt.Errorf("error deleting participant: %v", err)
			}
		}
	} else {
		res, err := tx.ExecContext(
			ctx,
			"DELETE FROM person WHERE event_year_id=?;",
			eventYearId,
		)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("error deleting participants: %v", err)
		}
		count, err = res.RowsAffected()
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("error fetching rows affected from person deletion: %v", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error committing transaction: %v", err)
	}
	return count, nil
}

func (m *MySQL) UpdatePerson(eventYearID int64, person types.Person) (*types.Person, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction: %v", err)
	}
	_, err = tx.ExecContext(
		ctx,
		"UPDATE person SET "+
			"bib=?, "+
			"first=?, "+
			"last=?, "+
			"age=?, "+
			"gender=?, "+
			"age_group=?, "+
			"distance=?, "+
			"anonymous=? "+
			"WHERE event_year_id=? AND alternate_id=?;",
		person.Bib,
		person.First,
		person.Last,
		person.Age,
		person.Gender,
		person.AgeGroup,
		person.Distance,
		person.AnonyInt(),
		eventYearID,
		person.AlternateId,
	)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error adding person to database: %v", err)
	}
	res, err := tx.QueryContext(
		ctx,
		"SELECT person_id FROM person WHERE event_year_id=? AND alternate_id=?;",
		eventYearID,
		person.AlternateId,
	)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error retrieving person_id: %v", err)
	}
	defer res.Close()
	output := types.Person{
		Bib:         person.Bib,
		First:       person.First,
		Last:        person.Last,
		Age:         person.Age,
		Gender:      person.Gender,
		AgeGroup:    person.AgeGroup,
		Distance:    person.Distance,
		Anonymous:   person.Anonymous,
		AlternateId: person.AlternateId,
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
