package database

import (
	"chronokeep/results/types"
	"context"
	"fmt"
	"time"
)

// GetAccount Gets an account based on the email address provided.
func GetAccount(email string) (*types.Account, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT account_id, account_name, account_email, type FROM account WHERE account_deleted=FALSE AND account_email=?;",
		email,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving account: %v", err)
	}
	defer res.Close()
	var outAccount types.Account
	if res.Next() {
		err := res.Scan(
			&outAccount.Identifier,
			&outAccount.Name,
			&outAccount.Email,
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

// GetAccounts Get all accounts that have not been deleted.
func GetAccounts() ([]types.Account, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT account_id, account_name, account_email, type FROM account WHERE account_deleted=FALSE;",
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving accounts: %v", err)
	}
	defer res.Close()
	var outAccounts []types.Account
	for res.Next() {
		var account types.Account
		err := res.Scan(&account.Identifier, &account.Name, &account.Email, &account.Type)
		if err != nil {
			return nil, fmt.Errorf("error getting account information: %v", err)
		}
		outAccounts = append(outAccounts, account)
	}
	return outAccounts, nil
}

// AddAccount Adds an account to the database.
func AddAccount(account types.Account) (*types.Account, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"INSERT INTO account(account_name, account_email, type) VALUES (?, ?, ?)",
		account.Name,
		account.Email,
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
		Name:       account.Name,
		Email:      account.Email,
		Type:       account.Type,
	}, nil
}

// DeleteAccount Deletes an account from view, does not permanently delete from database.
// This does not cascade down.  Must be done manually.
func DeleteAccount(account types.Account) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"UPDATE account SET account_deleted=TRUE WHERE account_id=?",
		account.Identifier,
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
	return nil
}

// ResurrectAccount Brings an account out of the deleted state.
func ResurrectAccount(email string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"UPDATE account SET account_deleted=FALSE WHERE account_email=?",
		email,
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
func GetDeletedAccount(email string) (*types.Account, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT account_id, account_name, account_email, type FROM account WHERE account_deleted=TRUE AND account_email=?;",
		email,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving account: %v", err)
	}
	defer res.Close()
	var outAccount types.Account
	if res.Next() {
		err := res.Scan(
			&outAccount.Identifier,
			&outAccount.Name,
			&outAccount.Email,
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
func UpdateAccount(account types.Account) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"UPDATE account SET account_name=?, account_email=?, type=? WHERE account_id=?",
		account.Name,
		account.Email,
		account.Type,
		account.Identifier,
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
