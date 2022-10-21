/*
repopg реализует методы хранилища на postgresql
*/
package repopg

import (
	"github.com/b2b2b-pro/prototype/repository/config"
	"github.com/b2b2b-pro/prototype/repository/pkg/pgdb"
	"go.uber.org/zap"
)

type RepoPg struct {
	db *pgdb.PgDB
}

// CRUDL для фирм (postgres)

func New(cfg *config.Config) (*RepoPg, error) {
	var err error

	conn, err := pgdb.New(cfg)
	if err != nil {
		zap.S().Fatal("Не получилось создать коннект к postgresql: ", err, "\n")
	}
	pg := &RepoPg{db: conn}

	err = pg.migrate()

	return pg, err
}

func (pg *RepoPg) Stop() {
	zap.S().Info("postgres Stop\n")
	pg.db.Close()
}
