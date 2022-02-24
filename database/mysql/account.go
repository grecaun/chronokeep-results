package mysql

import (
	"chronokeep/results/types"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const (
	MaxLoginAttempts = 4
)

func (m *MySQL) getAccountInternal(unique, key *string, id *int64) (*types.Account, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	var res *sql.Rows
	if unique != nil {
		res, err = db.QueryContext(
			ctx,
			"SELECT account_id, account_unique, account_type FROM account WHERE account_deleted=FALSE AND account_unique=?;",
			unique,
		)
	} else if key != nil {
		res, err = db.QueryContext(
			ctx,
			"SELECT account_id, account_unique, account_type FROM account NATURAL JOIN api_key WHERE account_deleted=FALSE AND key_deleted=FALSE AND key_value=?;",
			key,
		)
	} else if id != nil {
		res, err = db.QueryContext(
			ctx,
			"SELECT account_id, account_unique, account_type FROM account WHERE account_deleted=FALSE AND account_id=?;",
			id,
		)
	} else {
		return nil, errors.New("no valid identifying value provided to internal method")
	}
	if err != nil {
		return nil, fmt.Errorf("error retrieving account: %v", err)
	}
	defer res.Close()
	var outAccount types.Account
	if res.Next() {
		err := res.Scan(
			&outAccount.Identifier,
			&outAccount.Unique,
			&outAccount.Type,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting account information: %v", err)
		}
	} else {
		return nil, nil
	}
	return &outAccount, nil
}

// GetAccount Gets an account based on the unique id provided.
func (m *MySQL) GetAccount(unique string) (*types.Account, error) {
	return m.getAccountInternal(&unique, nil, nil)
}

// GetAccountByKey Gets an account based upon an API key provided.
func (m *MySQL) GetAccountByKey(key string) (*types.Account, error) {
	return m.getAccountInternal(nil, &key, nil)
}

// GetAccoutByID Gets an account based upon the Account ID.
func (m *MySQL) GetAccountByID(id int64) (*types.Account, error) {
	return m.getAccountInternal(nil, nil, &id)
}

// GetAccounts Get all accounts that have not been deleted.
func (m *MySQL) GetAccounts() ([]types.Account, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT account_id, account_unique, account_type FROM account WHERE account_deleted=FALSE;",
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving accounts: %v", err)
	}
	defer res.Close()
	var outAccounts []types.Account
	for res.Next() {
		var account types.Account
		err := res.Scan(
			&account.Identifier,
			&account.Unique,
			&account.Type,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting account information: %v", err)
		}
		outAccounts = append(outAccounts, account)
	}
	return outAccounts, nil
}

// AddAccount Adds an account to the database.
func (m *MySQL) AddAccount(account types.Account) (*types.Account, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"INSERT INTO account(account_unique, account_type) VALUES (?, ?)",
		account.Unique,
		account.Type,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to add account: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("unable to determine ID for account: %v", err)
	}
	return &types.Account{
		Identifier: id,
		Unique:     account.Unique,
		Type:       account.Type,
	}, nil
}

// DeleteAccount Deletes an account from view, does not permanently delete from database.
// This does not delete events associated with this account, but does set keys to deleted.
func (m *MySQL) DeleteAccount(id int64) error {
	db, err := m.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"UPDATE account SET account_deleted=TRUE WHERE account_id=?",
		id,
	)
	if err != nil {
		return fmt.Errorf("error deleting account: %v", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected on delete account: %v", err)
	}
	if rows != 1 {
		return fmt.Errorf("error deleting account, rows affected: %v", rows)
	}
	_, err = db.ExecContext(
		ctx,
		"UPDATE api_key SET key_deleted=TRUE WHERE account_id=?",
		id,
	)
	if err != nil {
		return fmt.Errorf("error deleting keys attached to account: %v", err)
	}
	return nil
}

// ResurrectAccount Brings an account out of the deleted state.
func (m *MySQL) ResurrectAccount(unique string) error {
	db, err := m.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"UPDATE account SET account_deleted=FALSE WHERE account_unique=?",
		unique,
	)
	if err != nil {
		return fmt.Errorf("error resurrecting account: %v", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected on resurrect account: %v", err)
	}
	if rows != 1 {
		return fmt.Errorf("error resurrecting account, rows affected: %v", rows)
	}
	return nil
}

// GetDeletedAccount Returns a deleted account.
func (m *MySQL) GetDeletedAccount(unique string) (*types.Account, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT account_id, account_unique, account_type FROM account WHERE account_deleted=TRUE AND account_unique=?;",
		unique,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving account: %v", err)
	}
	defer res.Close()
	var outAccount types.Account
	if res.Next() {
		err := res.Scan(
			&outAccount.Identifier,
			&outAccount.Unique,
			&outAccount.Type,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting account information: %v", err)
		}
	} else {
		return nil, nil
	}
	return &outAccount, nil
}

// UpdateAccount Updates account information in the database.
func (m *MySQL) UpdateAccount(account types.Account) error {
	db, err := m.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"UPDATE account SET account_type=? WHERE account_deleted=FALSE AND account_unique=?",
		account.Type,
		account.Unique,
	)
	if err != nil {
		return fmt.Errorf("error updating account: %v", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected on update account: %v", err)
	}
	if rows != 1 {
		return fmt.Errorf("error updating account, rows affected: %v", rows)
	}
	return nil
}
