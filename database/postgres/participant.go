package postgres

import (
	"chronokeep/results/types"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
)

func (p *Postgres) AddParticipants(eventYearID int64, participants []types.Participant) ([]types.Participant, error) {
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
	output := make([]types.Participant, 0)
	for _, participant := range participants {
		var id int64
		err = tx.QueryRow(
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
				"updated_at=$14 "+
				"RETURNING (participant_id);",
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
		).Scan(&id)
		if err != nil {
			tx.Rollback(ctx)
			return nil, fmt.Errorf("error adding participant to database: %v", err)
		}
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
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}
	return output, nil
}

func (p *Postgres) GetParticipants(eventYearID int64, limit, page int, updated_after *int64) ([]types.Participant, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	var (
		res pgx.Rows
	)
	if limit > 0 {
		if updated_after != nil && *updated_after >= 0 {
			res, err = db.Query(
				ctx,
				"SELECT participant_id, bib, first, last, birthdate, gender, age_group, distance, anonymous, alternate_id, apparel, sms_enabled, mobile, updated_at "+
					"FROM participant WHERE event_year_id=$1 AND updated_at>=$2 ORDER BY distance ASC, last ASC, first ASC LIMIT $3 OFFSET $4;",
				eventYearID,
				*updated_after,
				limit,
				page*limit,
			)
		} else {
			res, err = db.Query(
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
			res, err = db.Query(
				ctx,
				"SELECT participant_id, bib, first, last, birthdate, gender, age_group, distance, anonymous, alternate_id, apparel, sms_enabled, mobile, updated_at "+
					"FROM participant WHERE event_year_id=$1 AND updated_at>=$2;",
				eventYearID,
				*updated_after,
			)
		} else {
			res, err = db.Query(
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

func (p *Postgres) DeleteParticipants(eventYearID int64, alternateIds []string) (int64, error) {
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
				"DELETE FROM participant WHERE event_year_id=$1 AND alternate_id=$2;",
				eventYearID,
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
			"DELETE FROM participant WHERE event_year_id=$1;",
			eventYearID,
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

func (p *Postgres) UpdateParticipant(eventYearID int64, participant types.Participant) (*types.Participant, error) {
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
	err = tx.QueryRow(
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
			"WHERE event_year_id=$13 AND alternate_id=$14 RETURNING (participant_id);",
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
	).Scan(&output.Identifier)
	if err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("error adding participant to database: %v", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}
	return &output, nil
}

func (p *Postgres) UpdateParticipants(eventYearID int64, participants []types.Participant) ([]types.Participant, error) {
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
	output := make([]types.Participant, 0)
	for _, participant := range participants {
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
		err = tx.QueryRow(
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
				"WHERE event_year_id=$13 AND alternate_id=$14 RETURNING (participant_id);",
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
		).Scan(&tmp.Identifier)
		if err != nil {
			tx.Rollback(ctx)
			return nil, fmt.Errorf("error adding participant to database: %v", err)
		}
		output = append(output, tmp)
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}
	return output, nil
}
