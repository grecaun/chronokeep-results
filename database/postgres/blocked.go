package postgres

import (
	"context"
	"fmt"
	"time"
)

// AddBlockedPhone adds a phone number to the blocked phone numbers list
func (p *Postgres) AddBlockedPhone(phone string) error {
	db, err := p.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	_, err = db.Exec(
		ctx,
		"INSERT INTO banned_phones(banned_phone) VALUES ($1) ON CONFLICT DO NOTHING;",
		phone,
	)
	if err != nil {
		return fmt.Errorf("error adding blocked phone: %v", err)
	}
	return nil
}

// AddBlockedPhones adds one or more phone numbers to the blocked phone numbers list
func (p *Postgres) AddBlockedPhones(phones []string) error {
	db, err := p.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("unable to start transaction: %v", err)
	}
	for _, phone := range phones {
		_, err := tx.Exec(
			ctx,
			"INSERT INTO banned_phones(banned_phone) VALUES ($1) ON CONFLICT DO NOTHING;",
			phone,
		)
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("error adding blocked phone: %v", err)
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}
	return nil
}

// GetBlockedPhones gets the blocked phone numbers list
func (p *Postgres) GetBlockedPhones() ([]string, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Query(
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
func (p *Postgres) UnblockPhone(phone string) error {
	db, err := p.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	_, err = db.Exec(
		ctx,
		"DELETE FROM banned_phones WHERE banned_phone=$1;",
		phone,
	)
	if err != nil {
		return fmt.Errorf("error unblocking phone number: %v", err)
	}
	return nil
}

// AddBlockedEmail adds an email address to the blocked emails list
func (p *Postgres) AddBlockedEmail(email string) error {
	db, err := p.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	_, err = db.Exec(
		ctx,
		"INSERT INTO banned_emails(banned_email) VALUES ($1) ON CONFLICT DO NOTHING;",
		email,
	)
	if err != nil {
		return fmt.Errorf("error adding blocked email: %v", err)
	}
	return nil
}

// AddBlockedEmails adds one or more email addresses to the blocked emails list
func (p *Postgres) AddBlockedEmails(emails []string) error {
	db, err := p.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("unable to start transaction: %v", err)
	}
	for _, email := range emails {
		_, err := tx.Exec(
			ctx,
			"INSERT INTO banned_emails(banned_email) VALUES ($1) ON CONFLICT DO NOTHING;",
			email,
		)
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("error adding blocked email: %v", err)
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}
	return nil
}

// GetBlockedEmails gets the blocked phone emails list
func (p *Postgres) GetBlockedEmails() ([]string, error) {
	db, err := p.GetDB()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	res, err := db.Query(
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
func (p *Postgres) UnblockEmail(email string) error {
	db, err := p.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	_, err = db.Exec(
		ctx,
		"DELETE FROM banned_emails WHERE banned_email=$1;",
		email,
	)
	if err != nil {
		return fmt.Errorf("error unblocking email: %v", err)
	}
	return nil
}
