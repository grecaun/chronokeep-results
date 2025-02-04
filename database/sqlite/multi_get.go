package sqlite

import (
	"chronokeep/results/types"
	"context"
	"database/sql"
	"fmt"
	"time"
)

// GetAccountAndEvent Gets an event and its corresponding account.
func (s *SQLite) GetAccountAndEvent(slug string) (*types.MultiGet, error) {
	db, err := s.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT "+
			"account_id, account_name, account_email, account_type, account_locked, "+
			"event_id, event_name, slug, website, image, contact_email, access_restricted, event_type, cert_name "+
			"FROM account NATURAL JOIN event WHERE account_deleted=FALSE AND event_deleted=FALSE and slug=?",
		slug,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting account and event from database: %v", err)
	}
	if res.Next() {
		outVal := types.MultiGet{
			Account: &types.Account{},
			Event:   &types.Event{},
		}
		err := res.Scan(
			&outVal.Account.Identifier,
			&outVal.Account.Name,
			&outVal.Account.Email,
			&outVal.Account.Type,
			&outVal.Account.Locked,
			&outVal.Event.Identifier,
			&outVal.Event.Name,
			&outVal.Event.Slug,
			&outVal.Event.Website,
			&outVal.Event.Image,
			&outVal.Event.ContactEmail,
			&outVal.Event.AccessRestricted,
			&outVal.Event.Type,
			&outVal.Event.CertificateName,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting values for account and event: %v", err)
		}
		outVal.Event.AccountIdentifier = outVal.Account.Identifier
		return &outVal, nil
	}
	return nil, nil
}

// GetAccountEventAndYear Gets an eventyear and its corresponding event and account.
func (s *SQLite) GetAccountEventAndYear(slug, year string) (*types.MultiGet, error) {
	db, err := s.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	var res *sql.Rows
	if year == "" {
		res, err = db.QueryContext(
			ctx,
			"SELECT "+
				"account_id, account_name, account_email, account_type, account_locked, "+
				"event_id, event_name, slug, website, image, contact_email, access_restricted, event_type, "+
				"event_year_id, year, date_time, live, days_allowed, cert_name, ranking_type "+
				"FROM account NATURAL JOIN event NATURAL JOIN event_year y INNER JOIN "+
				"(SELECT event_id AS e_id, MAX(date_time) AS d_time FROM event_year WHERE year_deleted=FALSE GROUP BY e_id) AS g ON g.e_id=y.event_id AND g.d_time=y.date_time "+
				"WHERE account_deleted=FALSE AND event_deleted=FALSE AND year_deleted=FALSE AND slug=?",
			slug,
		)
	} else {
		res, err = db.QueryContext(
			ctx,
			"SELECT "+
				"account_id, account_name, account_email, account_type, account_locked, "+
				"event_id, event_name, slug, website, image, contact_email, access_restricted, event_type, "+
				"event_year_id, year, date_time, live, days_allowed, cert_name, ranking_type "+
				"FROM account NATURAL JOIN event NATURAL JOIN event_year WHERE account_deleted=FALSE AND event_deleted=FALSE AND year_deleted=FALSE AND slug=? AND year=?",
			slug,
			year,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("error getting account and event from database: %v", err)
	}
	if res.Next() {
		outVal := types.MultiGet{
			Account:   &types.Account{},
			Event:     &types.Event{},
			EventYear: &types.EventYear{},
		}
		err := res.Scan(
			&outVal.Account.Identifier,
			&outVal.Account.Name,
			&outVal.Account.Email,
			&outVal.Account.Type,
			&outVal.Account.Locked,
			&outVal.Event.Identifier,
			&outVal.Event.Name,
			&outVal.Event.Slug,
			&outVal.Event.Website,
			&outVal.Event.Image,
			&outVal.Event.ContactEmail,
			&outVal.Event.AccessRestricted,
			&outVal.Event.Type,
			&outVal.EventYear.Identifier,
			&outVal.EventYear.Year,
			&outVal.EventYear.DateTime,
			&outVal.EventYear.Live,
			&outVal.EventYear.DaysAllowed,
			&outVal.Event.CertificateName,
			&outVal.EventYear.RankingType,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting values for account and event: %v", err)
		}
		outVal.Event.AccountIdentifier = outVal.Account.Identifier
		outVal.EventYear.EventIdentifier = outVal.Event.Identifier
		return &outVal, nil
	}
	return nil, nil
}

// GetEventAndYear Gets an event and eventyear.
func (s *SQLite) GetEventAndYear(slug, year string) (*types.MultiGet, error) {
	db, err := s.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	var res *sql.Rows
	var countRes *sql.Rows
	if year == "" {
		res, err = db.QueryContext(
			ctx,
			"SELECT "+
				"account_id, event_id, event_name, slug, website, image, contact_email, access_restricted, event_type, "+
				"event_year_id, year, date_time, live, days_allowed, cert_name, ranking_type "+
				"FROM event NATURAL JOIN event_year y INNER JOIN "+
				"(SELECT event_id AS e_id, MAX(date_time) AS d_time FROM event_year WHERE year_deleted=FALSE GROUP BY e_id) AS g ON g.e_id=y.event_id AND g.d_time=y.date_time "+
				" WHERE event_deleted=FALSE AND year_deleted=FALSE AND slug=?",
			slug,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting account and event from database: %v", err)
		}
		countRes, err = db.QueryContext(
			ctx,
			"SELECT COUNT(DISTINCT distance) AS dist_count FROM event NATURAL JOIN event_year y INNER JOIN "+
				"(SELECT event_id AS e_id, MAX(date_time) AS d_time FROM event_year WHERE year_deleted=FALSE GROUP BY e_id) AS g ON g.e_id=y.event_id AND g.d_time=y.date_time "+
				"NATURAL JOIN person "+
				" WHERE event_deleted=FALSE AND year_deleted=FALSE and slug=?",
			slug,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting distance count from database: %v", err)
		}
	} else {
		res, err = db.QueryContext(
			ctx,
			"SELECT "+
				"account_id, event_id, event_name, slug, website, image, contact_email, access_restricted, event_type, "+
				"event_year_id, year, date_time, live, days_allowed, cert_name, ranking_type "+
				"FROM event NATURAL JOIN event_year WHERE event_deleted=FALSE AND year_deleted=FALSE AND slug=? AND year=?",
			slug,
			year,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting account and event from database: %v", err)
		}
		countRes, err = db.QueryContext(
			ctx,
			"SELECT COUNT(DISTINCT distance) AS dist_count "+
				"FROM person NATURAL JOIN event_year "+
				"NATURAL JOIN event "+
				" WHERE event_deleted=FALSE AND year_deleted=FALSE AND slug=? AND year=?",
			slug,
			year,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting distance count from database: %v", err)
		}
	}
	var distCount = 0
	if countRes.Next() {
		err := countRes.Scan(
			&distCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting distance count from database: %v", err)
		}
	}
	if res.Next() {
		outVal := types.MultiGet{
			Event:         &types.Event{},
			EventYear:     &types.EventYear{},
			DistanceCount: &distCount,
		}
		err := res.Scan(
			&outVal.Event.AccountIdentifier,
			&outVal.Event.Identifier,
			&outVal.Event.Name,
			&outVal.Event.Slug,
			&outVal.Event.Website,
			&outVal.Event.Image,
			&outVal.Event.ContactEmail,
			&outVal.Event.AccessRestricted,
			&outVal.Event.Type,
			&outVal.EventYear.Identifier,
			&outVal.EventYear.Year,
			&outVal.EventYear.DateTime,
			&outVal.EventYear.Live,
			&outVal.EventYear.DaysAllowed,
			&outVal.Event.CertificateName,
			&outVal.EventYear.RankingType,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting values for account and event: %v", err)
		}
		outVal.EventYear.EventIdentifier = outVal.Event.Identifier
		return &outVal, nil
	}
	return nil, nil
}

// GetKeyAndAccount Gets an account and key based upon the key value.
func (s *SQLite) GetKeyAndAccount(key string) (*types.MultiKey, error) {
	db, err := s.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT "+
			"account_id, account_name, account_email, account_type, account_locked, "+
			"key_value, key_type, allowed_hosts, valid_until "+
			"FROM account NATURAL JOIN api_key WHERE account_deleted=FALSE AND key_deleted=FALSE AND key_value=?",
		key,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting account and event from database: %v", err)
	}
	if res.Next() {
		outVal := types.MultiKey{
			Key:     &types.Key{},
			Account: &types.Account{},
		}
		err := res.Scan(
			&outVal.Account.Identifier,
			&outVal.Account.Name,
			&outVal.Account.Email,
			&outVal.Account.Type,
			&outVal.Account.Locked,
			&outVal.Key.Value,
			&outVal.Key.Type,
			&outVal.Key.AllowedHosts,
			&outVal.Key.ValidUntil,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting values for account and event: %v", err)
		}
		outVal.Key.AccountIdentifier = outVal.Account.Identifier
		return &outVal, nil
	}
	return nil, nil
}
