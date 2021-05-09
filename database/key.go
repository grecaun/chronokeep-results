package database

import (
	"chronokeep/results/types"
	"context"
	"fmt"
	"time"
)

// GetAccountKeys Gets all keys associated with an account.
func GetAccountKeys(email string) ([]types.Key, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT account_id, key_value, key_type, allowed_hosts, valid_until FROM api_key NATURAL JOIN account WHERE key_deleted=FALSE AND account_email=?;",
		email,
	)
	if err != nil {
		return nil, fmt.Errorf("erorr retrieving key: %v", err)
	}
	defer res.Close()
	var outKeys []types.Key
	for res.Next() {
		var key types.Key
		err := res.Scan(
			&key.AccountIdentifier,
			&key.Value,
			&key.Type,
			&key.AllowedHosts,
			&key.ValidUntil,
		)
		if err != nil {
			return nil, fmt.Errorf("unable getting key: %v", err)
		}
		outKeys = append(outKeys, key)
	}
	return outKeys, nil
}

// GetKey Gets a key based on its key value.
func GetKey(key string) (*types.Key, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT account_id, key_value, key_type, allowed_hosts, valid_until FROM api_key WHERE key_deleted=FALSE AND key_value=?;",
		key,
	)
	if err != nil {
		return nil, fmt.Errorf("erorr retrieving key: %v", err)
	}
	defer res.Close()
	var outKey types.Key
	if res.Next() {
		err := res.Scan(
			&outKey.AccountIdentifier,
			&outKey.Value,
			&outKey.Type,
			&outKey.AllowedHosts,
			&outKey.ValidUntil,
		)
		if err != nil {
			return nil, fmt.Errorf("unable getting key: %v", err)
		}
	} else {
		return nil, nil
	}
	return &outKey, nil
}

// AddKey Adds a key to the database.
func AddKey(key types.Key) (*types.Key, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"INSERT INTO api_key(account_id, key_value, key_type, allowed_hosts, valid_until) VALUES (?, ?, ?, ?, ?);",
		key.AccountIdentifier,
		key.Value,
		key.Type,
		key.AllowedHosts,
		key.ValidUntil,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to add key: %v", err)
	}
	_, err = res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("unable to determine ID for key: %v", err)
	}
	return &types.Key{
		AccountIdentifier: key.AccountIdentifier,
		Value:             key.Value,
		Type:              key.Type,
		AllowedHosts:      key.AllowedHosts,
		ValidUntil:        key.ValidUntil,
	}, nil
}

// DeleteKey Deletes a key from view. Does not permanently delete a key.
func DeleteKey(key types.Key) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"UPDATE api_key SET key_deleted=TRUE WHERE key_value=?;",
		key.Value,
	)
	if err != nil {
		return fmt.Errorf("error deleting key: %v", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected on delete key: %v", err)
	}
	if rows != 1 {
		return fmt.Errorf("error deleting key, rows affected: %v", err)
	}
	return nil
}

// UpdateKey Updates information for a key.  Cannot change account_id or value.
func UpdateKey(key types.Key) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"UPDATE api_key SET key_type=?, allowed_hosts=?, valid_until=? WHERE key_value=?;",
		key.Type,
		key.AllowedHosts,
		key.ValidUntil,
		key.Value,
	)
	if err != nil {
		return fmt.Errorf("error updating key: %v", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected on update key: %v", err)
	}
	if rows != 1 {
		return fmt.Errorf("error updating key, rows affected: %v", err)
	}
	return nil
}
