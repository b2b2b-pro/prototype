package repopg

import (
	"context"

	"go.uber.org/zap"
)

func (pg *RepoPg) migrate() error {
	var err error

	err = pg.tableAccount()
	if err != nil {
		return err
	}

	err = pg.tableEntity()
	if err != nil {
		return err
	}

	err = pg.tableAccountEntity()
	if err != nil {
		return err
	}

	err = pg.tableOrigin()
	if err != nil {
		return err
	}

	err = pg.tableStatus()
	if err != nil {
		return err
	}

	err = pg.tableObligation()
	if err != nil {
		return err
	}

	return nil
}

func (pg *RepoPg) tableAccount() error {
	var err error
	query1 := `CREATE TABLE IF NOT EXISTS account ( 
		id SERIAL PRIMARY KEY,
		login varchar(12)
		);`

	zap.S().Debug("postgres Create Table Account\n")

	row1, err := pg.db.Pool.Query(context.Background(), query1)
	if err != nil {
		zap.S().Debugf("tableAccount Create Table error: %v\n", err)
	}

	defer row1.Close()

	// пусть всегда будет admin
	query2 := `INSERT INTO account (login)
	SELECT 'admin'
	WHERE NOT EXISTS (
	 SELECT login FROM account WHERE  login = 'admin'
	);`

	zap.S().Debug("postgres INSERT admin into Account\n")

	row2, err := pg.db.Pool.Query(context.Background(), query2)
	if err != nil {
		zap.S().Debugf("tableOrigin Create Table error: %v\n", err)
	}

	defer row2.Close()

	return err
}

func (pg *RepoPg) tableEntity() error {
	var err error
	query := `CREATE TABLE IF NOT EXISTS entity ( 
				ID SERIAL PRIMARY KEY, 
				INN varchar(12), 
				KPP varchar(9), 
				shortname varchar(80), 
				fullname varchar(160), 
				UNIQUE (INN, KPP)
			);`

	zap.S().Debug("postgres Create Table Entity\n")

	row, err := pg.db.Pool.Query(context.Background(), query)
	if err != nil {
		zap.S().Debugf("tableEntity Create Table error: %v\n", err)
	}

	defer row.Close()

	return err
}

func (pg *RepoPg) tableAccountEntity() error {
	var err error
	query := `CREATE TABLE IF NOT EXISTS account_entity (
		id SERIAL PRIMARY KEY,
		account_id int REFERENCES account (id),
		entity_id int REFERENCES entity (id),
		UNIQUE (account_id, entity_id)
		);`

	zap.S().Debug("postgres Create Table AccountEntity\n")

	row, err := pg.db.Pool.Query(context.Background(), query)
	if err != nil {
		zap.S().Debugf("tableAccountEntity Create Table error: %v\n", err)
	}

	defer row.Close()

	return err
}

func (pg *RepoPg) tableOrigin() error {
	var err error
	query1 := `CREATE TABLE IF NOT EXISTS origin (
		ID SERIAL PRIMARY KEY,
		description varchar(160)
	);`

	zap.S().Debug("postgres Create Table Origin\n")

	row1, err := pg.db.Pool.Query(context.Background(), query1)
	defer row1.Close()

	if err != nil {
		zap.S().Debugf("tableOrigin Create Table error: %v\n", err)
		return err
	}

	query2 := `INSERT INTO origin (description)
	SELECT 'тестовый ввод' 
	WHERE NOT EXISTS (
	 SELECT description FROM origin WHERE  description = 'тестовый ввод'
	);`

	zap.S().Debug("postgres INSERT тестовый ввод Origin\n")

	row2, err := pg.db.Pool.Query(context.Background(), query2)
	if err != nil {
		zap.S().Debugf("tableOrigin Create Table error: %v\n", err)
	}

	defer row2.Close()

	return err
}

func (pg *RepoPg) tableStatus() error {
	var err error
	query1 := `CREATE TABLE IF NOT EXISTS obl_status (
		id SERIAL PRIMARY KEY,
		status varchar(24)
	);`

	zap.S().Debug("postgres Create Table Obl_Status\n")

	row1, err := pg.db.Pool.Query(context.Background(), query1)
	defer row1.Close()

	if err != nil {
		zap.S().Debugf("tableStatus Create Table error: %v\n", err)
		return err
	}

	query2 := `INSERT INTO obl_status (status)
	SELECT 'new'
	WHERE NOT EXISTS (
	SELECT status FROM obl_status WHERE  status = 'new'
	);`

	zap.S().Debug("postgres INSERT new status\n")

	row2, err := pg.db.Pool.Query(context.Background(), query2)
	if err != nil {
		zap.S().Debugf("tableOrigin Create Table error: %v\n", err)
	}

	defer row2.Close()

	return err
}

func (pg *RepoPg) tableObligation() error {
	var err error
	query := `CREATE TABLE IF NOT EXISTS obligation (
		ID SERIAL PRIMARY KEY,
		debtor_id int REFERENCES entity (id),
		creditor_id int REFERENCES entity (id),
		cost bigint,
		payment_date date,
		origin_id int REFERENCES origin (ID),
		status_id int REFERENCES obl_status (ID) NULL
	);`

	zap.S().Debug("postgres Create Table Obligation\n")

	row, err := pg.db.Pool.Query(context.Background(), query)
	if err != nil {
		zap.S().Debugf("tableObligation Create Table error: %v\n", err)
	}

	defer row.Close()

	return err
}
