package repopg

import (
	"context"

	"github.com/b2b2b-pro/lib/object"
	"go.uber.org/zap"
)

func (pg *RepoPg) CreateObligation(o object.Obligation) (int, error) {
	var err error

	zap.S().Debug("postgres CreateObligation: ", o, "\n")

	query := "INSERT INTO obligation (debtor_id, creditor_id, cost, origin_id, date) VALUES ($1, $2, $3, $4, $5) RETURNING (id);"
	// TODO: разобраться как это чинится средствами SQL
	if o.OriginID == 0 {
		o.OriginID = 1 // если не установлено происхождение, использовать "тестовый ввод"
	}

	row, err := pg.db.Pool.Query(context.Background(), query, o.DebtorID, o.CreditorID, o.Cost, o.OriginID, o.Date.String())
	if err != nil {
		zap.S().Error("INSERT INTO obligation Error: ", err, "\n")
		return 0, err
	}

	defer row.Close()

	if row.Next() {
		if err = row.Scan(&o.ID); err != nil {
			zap.S().Debugf("%s, Error: %v\n", err)
		}
	}

	zap.S().Debug("Create Obligation ID: ", o.ID, "\n")

	return o.ID, err
}

func (pg *RepoPg) ListObligation() ([]object.Obligation, error) {
	zap.S().Debug("postgres ListObligation\n")

	rows, err := pg.db.Pool.Query(context.Background(), "SELECT * FROM obligation;")
	if err != nil {
		zap.S().Error("SELECT * FROM obligation; Error: ", err, "\n")
		return nil, err
	}

	defer rows.Close()

	var obls []object.Obligation

	for rows.Next() {
		o := object.Obligation{}
		// TODO научиться средствами SQL собирать табличку, заменяя INN на name

		err = rows.Scan(&o.ID, &o.DebtorID, &o.CreditorID, &o.Cost, &o.OriginID)
		if err != nil {
			zap.S().Error("ListObligation rows.scan error: ", err, "\n")
			return nil, err
		}

		obls = append(obls, o)
	}

	return obls, err
}
