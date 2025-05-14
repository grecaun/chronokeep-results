package sqlite

import (
	"chronokeep/results/types"
	"context"
	"database/sql"
	"fmt"
	"time"
)

func (s *SQLite) AddParticipants(eventYearID int64, participants []types.Participant) ([]types.Participant, error) {
	db, err := s.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction: %v", err)
	}
	participantstmt, err := tx.PrepareContext(
		ctx,
		"INSERT INTO participant("+
			"event_year_id, "+
			"bib, "+
			"first, "+
			"last, "+
			"birthdate, "+
			"gender, "+
			"age_group, "+
			"distance, "+
			"anonymous, "+
			"alternate_id, "+
			"apparel, "+
			"sms_enabled, "+
			"mobile, "+
			"updated_at"+
			") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) "+
			"ON CONFLICT (event_year_id, alternate_id) DO UPDATE SET "+
			"first=$3, "+
			"last=$4, "+
			"birthdate=$5, "+
			"gender=$6, "+
			"age_group=$7, "+
			"distance=$8, "+
			"anonymous=$9, "+
			"alternate_id=$10, "+
			"apparel=$11, "+
			"sms_enabled=$12, "+
			"mobile=$13, "+
			"updated_at=$14"+
			";",
	)
	if err != nil {
		return nil, fmt.Errorf("error preparing statement for adding participants: %v", err)
	}
	for _, participant := range participants {
		_, err = participantstmt.ExecContext(
			ctx,
			eventYearID,
			participant.Bib,
			participant.First,
			participant.Last,
			participant.Birthdate,
			participant.Gender,
			participant.AgeGroup,
			participant.Distance,
			participant.AnonyInt(),
			participant.AlternateId,
			participant.Apparel,
			participant.SMSInt(),
			participant.Mobile,
			participant.UpdatedAt,
		)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error adding participant to database: %v", err)
		}
	}
	// get a list of all the bibs and participants for the event year
	altMap := make(map[string]int64)
	res, err := tx.QueryContext(
		ctx,
		"SELECT participant_id, alternate_id FROM participant WHERE event_year_id=$1;",
		eventYearID,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying database for participant ids: %v", err)
	}
	for res.Next() {
		var id int64
		var altId string
		err = res.Scan(
			&id,
			&altId,
		)
		if err != nil {
			return nil, fmt.Errorf("error retrieving participant ids: %v", err)
		}
		altMap[altId] = id
	}
	output := make([]types.Participant, 0)
	for _, participant := range participants {
		if id, ok := altMap[participant.AlternateId]; ok {
			output = append(output, types.Participant{
				Identifier:  id,
				Bib:         participant.Bib,
				First:       participant.First,
				Last:        participant.Last,
				Birthdate:   participant.Birthdate,
				Gender:      participant.Gender,
				AgeGroup:    participant.AgeGroup,
				Distance:    participant.Distance,
				Anonymous:   participant.Anonymous,
				AlternateId: participant.AlternateId,
				Apparel:     participant.Apparel,
				SMSEnabled:  participant.SMSEnabled,
				Mobile:      participant.Mobile,
				UpdatedAt:   participant.UpdatedAt,
			})
		} else {
			tx.Rollback()
			return nil, fmt.Errorf("participant not found after add: %v", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}
	return output, nil
}

func (s *SQLite) GetParticipants(eventYearID int64, limit, page int, updated_after *int64) ([]types.Participant, error) {
	db, err := s.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	var (
		res *sql.Rows
	)
	if limit > 0 {
		if updated_after != nil && *updated_after >= 0 {
			res, err = db.QueryContext(
				ctx,
				"SELECT participant_id, bib, first, last, birthdate, gender, age_group, distance, anonymous, alternate_id, apparel, sms_enabled, mobile, updated_at "+
					"FROM participant WHERE event_year_id=$1 AND updated_at>=$2 ORDER BY distance ASC, last ASC, first ASC LIMIT $3 OFFSET $4;",
				eventYearID,
				*updated_after,
				limit,
				page*limit,
			)
		} else {
			res, err = db.QueryContext(
				ctx,
				"SELECT participant_id, bib, first, last, birthdate, gender, age_group, distance, anonymous, alternate_id, apparel, sms_enabled, mobile, updated_at "+
					"FROM participant WHERE event_year_id=$1 ORDER BY distance ASC, last ASC, first ASC LIMIT $2 OFFSET $3;",
				eventYearID,
				limit,
				page*limit,
			)
		}
	} else {
		if updated_after != nil && *updated_after >= 0 {
			res, err = db.QueryContext(
				ctx,
				"SELECT participant_id, bib, first, last, birthdate, gender, age_group, distance, anonymous, alternate_id, apparel, sms_enabled, mobile, updated_at "+
					"FROM participant WHERE event_year_id=$1 AND updated_at>=$2;",
				eventYearID,
				*updated_after,
			)
		} else {
			res, err = db.QueryContext(
				ctx,
				"SELECT participant_id, bib, first, last, birthdate, gender, age_group, distance, anonymous, alternate_id, apparel, sms_enabled, mobile, updated_at "+
					"FROM participant WHERE event_year_id=$1;",
				eventYearID,
			)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("error retrieving participants: %v", err)
	}
	defer res.Close()
	output := make([]types.Participant, 0)
	for res.Next() {
		var participant types.Participant
		var anonymous int
		var sms int
		err := res.Scan(
			&participant.Identifier,
			&participant.Bib,
			&participant.First,
			&participant.Last,
			&participant.Birthdate,
			&participant.Gender,
			&participant.AgeGroup,
			&participant.Distance,
			&anonymous,
			&participant.AlternateId,
			&participant.Apparel,
			&sms,
			&participant.Mobile,
			&participant.UpdatedAt,
		)
		participant.Anonymous = anonymous != 0
		participant.SMSEnabled = sms != 0
		if err != nil {
			return nil, fmt.Errorf("error getting participant: %v", err)
		}
		output = append(output, participant)
	}
	return output, nil
}

func (s *SQLite) DeleteParticipants(eventYearID int64, alternateIds []string) (int64, error) {
	db, err := s.GetDB()
	if err != nil {
		return 0, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("error starting transaction: %v", err)
	}
	var count int64
	if len(alternateIds) > 0 {
		count = 0
		for _, altId := range alternateIds {
			count++
			_, err = tx.ExecContext(
				ctx,
				"DELETE FROM participant WHERE event_year_id=$1 AND alternate_id=$2;",
				eventYearID,
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
			"DELETE FROM participant WHERE event_year_id=$1;",
			eventYearID,
		)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("error deleting participants: %v", err)
		}
		count, err = res.RowsAffected()
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("error fetching rows affected from participants deletion: %v", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error committing transaction: %v", err)
	}
	return count, nil
}

func (s *SQLite) UpdateParticipant(eventYearID int64, participant types.Participant) (*types.Participant, error) {
	db, err := s.GetDB()
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
		"UPDATE participant SET "+
			"bib=$1, "+
			"first=$2, "+
			"last=$3, "+
			"birthdate=$4, "+
			"gender=$5, "+
			"age_group=$6, "+
			"distance=$7, "+
			"anonymous=$8, "+
			"apparel=$9, "+
			"sms_enabled=$10, "+
			"mobile=$11, "+
			"updated_at=$12 "+
			"WHERE event_year_id=$13 AND alternate_id=$14;",
		participant.Bib,
		participant.First,
		participant.Last,
		participant.Birthdate,
		participant.Gender,
		participant.AgeGroup,
		participant.Distance,
		participant.AnonyInt(),
		participant.Apparel,
		participant.SMSInt(),
		participant.Mobile,
		participant.UpdatedAt,
		eventYearID,
		participant.AlternateId,
	)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error adding participant to database: %v", err)
	}
	res, err := tx.QueryContext(
		ctx,
		"SELECT participant_id FROM participant WHERE event_year_id=$1 AND alternate_id=$2;",
		eventYearID,
		participant.AlternateId,
	)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error retrieving participant_id: %v", err)
	}
	defer res.Close()
	output := types.Participant{
		Bib:         participant.Bib,
		First:       participant.First,
		Last:        participant.Last,
		Birthdate:   participant.Birthdate,
		Gender:      participant.Gender,
		AgeGroup:    participant.AgeGroup,
		Distance:    participant.Distance,
		Anonymous:   participant.Anonymous,
		AlternateId: participant.AlternateId,
		SMSEnabled:  participant.SMSEnabled,
		Apparel:     participant.Apparel,
		Mobile:      participant.Mobile,
		UpdatedAt:   participant.UpdatedAt,
	}
	if res.Next() {
		res.Scan(
			&output.Identifier,
		)
	} else {
		tx.Rollback()
		return nil, fmt.Errorf("participant not found after add: %v", err)
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}
	return &output, nil
}

func (s *SQLite) UpdateParticipants(eventYearID int64, participants []types.Participant) ([]types.Participant, error) {
	db, err := s.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction: %v", err)
	}
	output := make([]types.Participant, 0)
	for _, participant := range participants {
		_, err = tx.ExecContext(
			ctx,
			"UPDATE participant SET "+
				"bib=$1, "+
				"first=$2, "+
				"last=$3, "+
				"birthdate=$4, "+
				"gender=$5, "+
				"age_group=$6, "+
				"distance=$7, "+
				"anonymous=$8, "+
				"apparel=$9, "+
				"sms_enabled=$10, "+
				"mobile=$11, "+
				"updated_at=$12 "+
				"WHERE event_year_id=$13 AND alternate_id=$14;",
			participant.Bib,
			participant.First,
			participant.Last,
			participant.Birthdate,
			participant.Gender,
			participant.AgeGroup,
			participant.Distance,
			participant.AnonyInt(),
			participant.Apparel,
			participant.SMSInt(),
			participant.Mobile,
			participant.UpdatedAt,
			eventYearID,
			participant.AlternateId,
		)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error adding participant to database: %v", err)
		}
		res, err := tx.QueryContext(
			ctx,
			"SELECT participant_id FROM participant WHERE event_year_id=$1 AND alternate_id=$2;",
			eventYearID,
			participant.AlternateId,
		)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error retrieving participant_id: %v", err)
		}
		defer res.Close()
		tmp := types.Participant{
			Bib:         participant.Bib,
			First:       participant.First,
			Last:        participant.Last,
			Birthdate:   participant.Birthdate,
			Gender:      participant.Gender,
			AgeGroup:    participant.AgeGroup,
			Distance:    participant.Distance,
			Anonymous:   participant.Anonymous,
			AlternateId: participant.AlternateId,
			SMSEnabled:  participant.SMSEnabled,
			Apparel:     participant.Apparel,
			Mobile:      participant.Mobile,
			UpdatedAt:   participant.UpdatedAt,
		}
		if res.Next() {
			res.Scan(
				&tmp.Identifier,
			)
			output = append(output, tmp)
		} else {
			tx.Rollback()
			return nil, fmt.Errorf("participant not found after add: %v", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}
	return output, nil
}
