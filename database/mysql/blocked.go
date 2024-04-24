package mysql

import (
	"context"
	"fmt"
	"time"
)

// AddBlockedPhone adds a phone number to the blocked phone numbers list
func (m *MySQL) AddBlockedPhone(phone string) error {
	db, err := m.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"INSERT IGNORE INTO banned_phones(banned_phone) VALUES (?);",
		phone,
	)
	if err != nil {
		return fmt.Errorf("error adding blocked phone: %v", err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to determine number of rows on add blocked phone: %v", err)
	}
	return nil
}

// AddBlockedPhones adds one or more phone numbers to the blocked phone numbers list
func (m *MySQL) AddBlockedPhones(phones []string) error {
	db, err := m.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("unable to start transaction: %v", err)
	}
	stmt, err := tx.PrepareContext(
		ctx,
		"INSERT IGNORE INTO banned_phones(banned_phone) VALUES (?);",
	)
	if err != nil {
		return fmt.Errorf("unable to prepare statement for multiple blocked phone adds: %v", err)
	}
	defer stmt.Close()
	for _, phone := range phones {
		_, err := stmt.ExecContext(
			ctx,
			phone,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error adding blocked phone: %v", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}
	return nil
}

// GetBlockedPhones gets the blocked phone numbers list
func (m *MySQL) GetBlockedPhones() ([]string, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT * FROM banned_phones;",
	)
	if err != nil {
		return nil, fmt.Errorf("unable to query for blocked phone numbers: %v", err)
	}
	defer res.Close()
	var phones []string
	for res.Next() {
		var phone string
		err := res.Scan(&phone)
		if err != nil {
			return nil, fmt.Errorf("unable to pull phone number: %v", err)
		}
		phones = append(phones, phone)
	}
	return phones, nil
}

// UnblockPhone removes a phone number from the blocked phone numbers list
func (m *MySQL) UnblockPhone(phone string) error {
	db, err := m.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"DELETE FROM banned_phones WHERE banned_phone=?;",
		phone,
	)
	if err != nil {
		return fmt.Errorf("error unblocking phone number: %v", err)
	}
	num, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error determining number of rows on unblock phone: %v", err)
	}
	if num < 1 {
		return fmt.Errorf("phone number not found")
	}
	return nil
}

// AddBlockedEmail adds an email address to the blocked emails list
func (m *MySQL) AddBlockedEmail(email string) error {
	db, err := m.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"INSERT IGNORE INTO banned_emails(banned_email) VALUES (?);",
		email,
	)
	if err != nil {
		return fmt.Errorf("error adding blocked email: %v", err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to determine number of rows on add blocked email: %v", err)
	}
	return nil
}

// AddBlockedEmails adds one or more email addresses to the blocked emails list
func (m *MySQL) AddBlockedEmails(emails []string) error {
	db, err := m.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("unable to start transaction: %v", err)
	}
	stmt, err := tx.PrepareContext(
		ctx,
		"INSERT IGNORE INTO banned_emails(banned_email) VALUES (?);",
	)
	if err != nil {
		return fmt.Errorf("unable to prepare statement for multiple blocked email adds: %v", err)
	}
	defer stmt.Close()
	for _, email := range emails {
		_, err := stmt.ExecContext(
			ctx,
			email,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error adding blocked email: %v", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}
	return nil
}

// GetBlockedEmails gets the blocked phone emails list
func (m *MySQL) GetBlockedEmails() ([]string, error) {
	db, err := m.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.QueryContext(
		ctx,
		"SELECT * FROM banned_emails;",
	)
	if err != nil {
		return nil, fmt.Errorf("unable to query for blocked emails: %v", err)
	}
	defer res.Close()
	var emails []string
	for res.Next() {
		var email string
		err := res.Scan(&email)
		if err != nil {
			return nil, fmt.Errorf("unable to pull email: %v", err)
		}
		emails = append(emails, email)
	}
	return emails, nil
}

// UnblockEmail removes an email address from the blocked blocked emails list
func (m *MySQL) UnblockEmail(email string) error {
	db, err := m.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.ExecContext(
		ctx,
		"DELETE FROM banned_emails WHERE banned_email=?;",
		email,
	)
	if err != nil {
		return fmt.Errorf("error unblocking email: %v", err)
	}
	num, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error determining number of rows on unblock email: %v", err)
	}
	if num < 1 {
		return fmt.Errorf("email not found")
	}
	return nil
}
