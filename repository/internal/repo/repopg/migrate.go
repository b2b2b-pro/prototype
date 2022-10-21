package repopg

import (
	"context"

	"go.uber.org/zap"
)

func (pg *RepoPg) migrate() error {
	var err error

	err = pg.tableEntity()
	if err != nil {
		zap.S().Errorf("postgres Create Table Entity error: %v\n", err)
		return err
	}

	err = pg.tableOrigin()
	if err != nil {
		zap.S().Errorf("postgres Create Table Origin error: %v\n", err)
		return err
	}

	err = pg.tableObligation()
	if err != nil {
		zap.S().Errorf("postgres Create Table Obligation error: %v\n", err)
		return err
	}

	return nil
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
	defer row.Close()

	zap.S().Debugf("tableEntity Create Table error: %v\n", err)

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
		return err
	}

	query2 := `INSERT INTO origin (description)
	SELECT 'тестовый ввод' 
	WHERE NOT EXISTS (
	 SELECT description FROM origin WHERE  description = 'тестовый ввод'
	);`

	zap.S().Debug("postgres INSERT тестовый ввод Origin\n")

	row2, err := pg.db.Pool.Query(context.Background(), query2)
	defer row2.Close()

	zap.S().Debugf("tableOrigin Create Table error: %v\n", err)

	return err
}

func (pg *RepoPg) tableObligation() error {
	var err error
	query := `CREATE TABLE IF NOT EXISTS obligation (
		ID SERIAL PRIMARY KEY,
		debtor_id int REFERENCES entity (id),
		creditor_id int REFERENCES entity (id),
		cost money,
		date date,
		origin_id int REFERENCES origin (ID),
		status varchar(10)
	);`

	zap.S().Debug("postgres Create Table Obligation\n")

	row, err := pg.db.Pool.Query(context.Background(), query)
	defer row.Close()

	zap.S().Debugf("tableObligation Create Table error: %v\n", err)

	return err
}
