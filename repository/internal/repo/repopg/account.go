package repopg

import (
	"context"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

// возвращает id и login учётной записи
func (pg *RepoPg) getAccount(tkn string) (int, string) {
	// TODO нужно ли проверить не просрочен ли токен?
	t, err := pg.tp.DecryptToken(tkn)
	if err != nil {
		zap.S().Debugf("DecryptToken: %v\nError: %v\n", tkn, err)
		return -1, ""
	} else {
		zap.S().Debugf("DecryptToken User: %v\n", t.Credential)
	}

	zap.S().Debugf("postgres find account: %s\n", t.Credential)

	query := "SELECT id FROM account WHERE login = $1;"
	// query := "SELECT id FROM account WHERE login = 'user01';"
	var row pgx.Rows
	var id int

	row, err = pg.db.Pool.Query(context.Background(), query, t.Credential)
	// row, err = pg.db.Pool.Query(context.Background(), query)
	if err != nil {
		zap.S().Errorf("%s. Error: %v\n", query, err)
		return -1, t.Credential
	}
	defer row.Close()

	if row.Next() {
		err = row.Scan(&id)
		if err != nil {
			zap.S().Error("SELECT id FROM entity WHERE ... scan; Error: ", err, "\n")
			return -1, t.Credential
		}
	}

	/*
		err = row.Scan(&id)
		if err != nil {
			zap.S().Errorf("row.Scan %s. Error: %v\n", query, err)
			return -1, t.Credential
		}
	*/

	return id, t.Credential
}
