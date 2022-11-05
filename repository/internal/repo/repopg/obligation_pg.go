package repopg

import (
	"context"
	"fmt"
	"time"

	"github.com/b2b2b-pro/lib/object"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

// TODO проверить токен
func (pg *RepoPg) CreateObligation(tkn string, o object.Obligation) (int, error) {
	var err error

	zap.S().Debug("postgres CreateObligation: ", o, "\n")

	query := "INSERT INTO obligation (debtor_id, creditor_id, cost, origin_id, payment_date) VALUES ($1, $2, $3, $4, $5) RETURNING (id);"
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

// TODO проверить токен
func (pg *RepoPg) ListObligation(tkn string) ([]object.Obligation, error) {
	zap.S().Debug("postgres ListObligation\n")

	var obls []object.Obligation
	var err error
	var rows pgx.Rows
	var d interface{}

	id, lgn := pg.getAccount(tkn)

	if id < 1 && lgn != "admin" {
		return obls, fmt.Errorf("user %s not found", lgn)
	}

	if lgn == "admin" {
		rows, err = pg.db.Pool.Query(context.Background(), "SELECT id, debtor_id, creditor_id, cost, payment_date, origin_id FROM obligation;")
	} else {
		query := `select 
			obligation.id, obligation.debtor_id, obligation.creditor_id, obligation.cost, obligation.payment_date, obligation.origin_id 
		from obligation join  
		(select account_entity.account_id, account_entity.entity_id from account_entity 
			where account_entity.account_id = $1) as A 
		ON obligation.debtor_id = A.entity_id OR obligation.creditor_id = A.entity_id;`

		rows, err = pg.db.Pool.Query(context.Background(), query, id)
	}

	if err != nil {
		zap.S().Debugf("SELECT * FROM obligation; Error: %s\n", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		o := object.Obligation{}
		// TODO средствами SQL собирать табличку, заменяя INN на name
		err = rows.Scan(&o.ID, &o.DebtorID, &o.CreditorID, &o.Cost, &d, &o.OriginID)
		if err != nil {
			zap.S().Error("ListObligation rows.scan error: ", err, "\n")
			return nil, err
		}

		if d == nil {
			o.Date = object.NewPaymentDate(time.Now())
		} else {
			o.Date = object.NewPaymentDate(d.(time.Time))
		}

		zap.S().Debugf("Долг: %v\nдата: %s\n", o, &o.Date)
		obls = append(obls, o)
	}

	return obls, err
}
