package postgres

import (
	"chronokeep/results/types"
	"context"
	"fmt"
	"time"
)

func (p *Postgres) AddSegments(eventYearID int64, segments []types.Segment) ([]types.Segment, error) {
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
	for _, seg := range segments {
		_, err = tx.Exec(
			ctx,
			"INSERT INTO segments("+
				"event_year_id, "+
				"location_name, "+
				"distance_name, "+
				"segment_name, "+
				"segment_distance, "+
				"segment_distance_unit, "+
				"segment_gps, "+
				"segment_map_link"+
				") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) "+
				"ON CONFLICT (event_year_id, distance_name, segment_name) DO UPDATE SET "+
				"location_name=$2, "+
				"segment_distance=$5, "+
				"segment_distance_unit=$6, "+
				"segment_gps=$7, "+
				"segment_map_link=$8"+
				";",
			eventYearID,
			seg.Location,
			seg.DistanceName,
			seg.Name,
			seg.DistanceValue,
			seg.DistanceUnit,
			seg.GPS,
			seg.MapLink,
		)
		if err != nil {
			tx.Rollback(ctx)
			return nil, fmt.Errorf("error adding segment to database: %v", err)
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}
	output := make([]types.Segment, 0)
	for _, seg := range segments {
		output = append(output, types.Segment{
			Location:      seg.Location,
			DistanceName:  seg.DistanceName,
			Name:          seg.Name,
			DistanceValue: seg.DistanceValue,
			DistanceUnit:  seg.DistanceUnit,
			GPS:           seg.GPS,
			MapLink:       seg.MapLink,
		})
	}
	return output, nil
}

func (s *Postgres) GetSegments(eventYearID int64) ([]types.Segment, error) {
	db, err := s.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Query(
		ctx,
		"SELECT segment_id, location_name, distance_name, segment_name, "+
			"segment_distance, segment_distance_unit, segment_gps, "+
			"segment_map_link FROM segments WHERE event_year_id=$1;",
		eventYearID,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving segments: %v", err)
	}
	defer res.Close()
	output := make([]types.Segment, 0)
	for res.Next() {
		var seg types.Segment
		err := res.Scan(
			&seg.Identifier,
			&seg.Location,
			&seg.DistanceName,
			&seg.Name,
			&seg.DistanceValue,
			&seg.DistanceUnit,
			&seg.GPS,
			&seg.MapLink,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting segment: %v", err)
		}
		output = append(output, seg)
	}
	return output, nil
}

func (s *Postgres) DeleteSegments(eventYearID int64) (int64, error) {
	db, err := s.GetDB()
	if err != nil {
		return 0, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("error starting transaction: %v", err)
	}
	res, err := tx.Exec(
		ctx,
		"DELETE FROM segments WHERE event_year_id=$1",
		eventYearID,
	)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("error deleting segments: %v", err)
	}
	count := res.RowsAffected()
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("error committing transaction: %v", err)
	}
	return count, nil
}
