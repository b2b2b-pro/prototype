package repopg

import (
	"context"

	"github.com/b2b2b-pro/lib/object"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

// CRUDL для фирм (postgres)

func (pg *RepoPg) CreateEntity(frm object.Entity) (int, error) {
	zap.S().Debug("postgres CreateEntity: ", frm, "\n")
	var row1, row2 pgx.Rows
	var err error
	qselect := "SELECT id FROM entity WHERE inn = $1 AND kpp = $2;"
	row1, err = pg.db.Pool.Query(context.Background(), qselect, frm.INN, frm.KPP)
	if err != nil {
		zap.S().Error("SELECT id FROM entity WHERE inn = $1 AND kpp = $2; Error: ", err, "\n")
		return 0, err
	}
	defer row1.Close()

	if row1.Next() {
		var id int
		err = row1.Scan(&id)
		if err != nil {
			zap.S().Error("SELECT id FROM entity WHERE ... scan; Error: ", err, "\n")
			return 0, err
		}
		if id != 0 {
			frm.ID = id
			return frm.ID, nil
		}
	}

	qinsert := "INSERT INTO entity (inn, kpp, shortname, fullname) VALUES ($1, $2, $3, $4) RETURNING (id);"
	row2, err = pg.db.Pool.Query(context.Background(), qinsert, frm.INN, frm.KPP, frm.ShortName, frm.FullName)
	if err != nil {
		zap.S().Error("INSERT INTO entity Error: ", err, "\n")
		return 0, err
	}
	defer row2.Close()
	// проверять результат
	if row2.Next() {
		var id int
		err = row2.Scan(&id)
		if err != nil {
			zap.S().Errorf("%s scan; Error: %v\n", qinsert, err)
			return 0, err
		}
		if id != 0 {
			frm.ID = id
			return frm.ID, nil
		}
	}
	// TDOD разве мы должны сюда попадать?
	zap.S().Debug("postgres return: ", frm.ID, "\n")
	return frm.ID, err
}

func (pg *RepoPg) ListEntity() ([]object.Entity, error) {
	zap.S().Debug("postgres ListEntity\n")

	query := "SELECT entity.id, entity.inn, entity.kpp, entity.shortname, entity.fullname FROM entity;"
	var rows pgx.Rows
	var err error
	rows, err = pg.db.Pool.Query(context.Background(), query)
	if err != nil {
		zap.S().Errorf("%s. Error: %v\n", query, err)
		return nil, err
	}
	defer rows.Close()

	var frms []object.Entity
	for rows.Next() {
		f := object.Entity{}
		err = rows.Scan(&f.ID, &f.INN, &f.KPP, &f.ShortName, &f.FullName)
		if err != nil {
			zap.S().Error("ListEntity rows.scan error: ", err, "\n")
			return nil, err
		}
		frms = append(frms, f)
	}

	return frms, err
}
