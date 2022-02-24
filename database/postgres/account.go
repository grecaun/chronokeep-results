package postgres

import (
	"chronokeep/results/types"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
)

func (p *Postgres) getAccountInternal(unique, key *string, id *int64) (*types.Account, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	var res pgx.Rows
	if unique != nil {
		res, err = db.Query(
			ctx,
			"SELECT account_id, account_unique, account_type FROM account WHERE account_deleted=FALSE AND account_unique=$1;",
			unique,
		)
	} else if key != nil {
		res, err = db.Query(
			ctx,
			"SELECT account_id, account_unique, account_type FROM account NATURAL JOIN api_key WHERE account_deleted=FALSE AND key_deleted=FALSE AND key_value=$1;",
			key,
		)
	} else if id != nil {
		res, err = db.Query(
			ctx,
			"SELECT account_id, account_unique, account_type FROM account WHERE account_deleted=FALSE AND account_id=$1;",
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

// GetAccount Gets an account based on the unique id address provided.
func (p *Postgres) GetAccount(unique string) (*types.Account, error) {
	return p.getAccountInternal(&unique, nil, nil)
}

// GetAccountByKey Gets an account based upon an API key provided.
func (p *Postgres) GetAccountByKey(key string) (*types.Account, error) {
	return p.getAccountInternal(nil, &key, nil)
}

// GetAccoutByID Gets an account based upon the Account ID.
func (p *Postgres) GetAccountByID(id int64) (*types.Account, error) {
	return p.getAccountInternal(nil, nil, &id)
}

// GetAccounts Get all accounts that have not been deleted.
func (p *Postgres) GetAccounts() ([]types.Account, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Query(
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
func (p *Postgres) AddAccount(account types.Account) (*types.Account, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	var id int64
	err = db.QueryRow(
		ctx,
		"INSERT INTO account(account_unique, account_type) VALUES ($1, $2) RETURNING (account_id);",
		account.Unique,
		account.Type,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("unable to add account: %v", err)
	}
	return &types.Account{
		Identifier: id,
		Unique:     account.Unique,
		Type:       account.Type,
	}, nil
}

// DeleteAccount Deletes an account from view, does not permanently delete from database.
// This does not delete events associated with this account, but does set keys to deleted.
func (p *Postgres) DeleteAccount(id int64) error {
	db, err := p.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Exec(
		ctx,
		"UPDATE account SET account_deleted=TRUE WHERE account_id=$1",
		id,
	)
	if err != nil {
		return fmt.Errorf("error deleting account: %v", err)
	}
	if res.RowsAffected() != 1 {
		return fmt.Errorf("error deleting account, rows affected: %v", res.RowsAffected())
	}
	_, err = db.Exec(
		ctx,
		"UPDATE api_key SET key_deleted=TRUE WHERE account_id=$1",
		id,
	)
	if err != nil {
		return fmt.Errorf("error deleting keys attached to account: %v", err)
	}
	return nil
}

// ResurrectAccount Brings an account out of the deleted state.
func (p *Postgres) ResurrectAccount(unique string) error {
	db, err := p.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Exec(
		ctx,
		"UPDATE account SET account_deleted=FALSE WHERE account_unique=$1",
		unique,
	)
	if err != nil {
		return fmt.Errorf("error resurrecting account: %v", err)
	}
	if res.RowsAffected() != 1 {
		return fmt.Errorf("error resurrecting account, rows affected: %v", res.RowsAffected())
	}
	return nil
}

// GetDeletedAccount Returns a deleted account.
func (p *Postgres) GetDeletedAccount(unique string) (*types.Account, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Query(
		ctx,
		"SELECT account_id, account_unique, account_type FROM account WHERE account_deleted=TRUE AND account_unique=$1;",
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
func (p *Postgres) UpdateAccount(account types.Account) error {
	db, err := p.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Exec(
		ctx,
		"UPDATE account SET account_type=$1 WHERE account_deleted=FALSE AND account_unique=$2",
		account.Type,
		account.Unique,
	)
	if err != nil {
		return fmt.Errorf("error updating account: %v", err)
	}
	if res.RowsAffected() != 1 {
		return fmt.Errorf("error updating account, rows affected: %v", res.RowsAffected())
	}
	return nil
}
