package sqlite

import (
	"chronokeep/results/types"
	"context"
	"fmt"
	"time"
)

func (s *SQLite) AddSegments(eventYearID int64, segments []types.Segment) ([]types.Segment, error) {
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
	segmentstmt, err := tx.PrepareContext(
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
	)
	if err != nil {
		return nil, fmt.Errorf("error preparing statement for adding segments: %v", err)
	}
	for _, seg := range segments {
		_, err = segmentstmt.ExecContext(
			ctx,
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
			tx.Rollback()
			return nil, fmt.Errorf("error adding segment to database: %v", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
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

func (s *SQLite) GetSegments(eventYearID int64) ([]types.Segment, error) {
	db, err := s.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
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

func (s *SQLite) DeleteSegments(eventYearID int64) (int64, error) {
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
	res, err := tx.ExecContext(
		ctx,
		"DELETE FROM segments WHERE event_year_id=$1",
		eventYearID,
	)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error deleting segments: %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error fetching rows affected from segments deletion: %v", err)
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error committing transaction: %v", err)
	}
	return count, nil
}
