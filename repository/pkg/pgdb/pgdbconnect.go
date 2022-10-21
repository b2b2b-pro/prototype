/*
pgdb обеспечивает подключение к СУБД
*/
package pgdb

import (
	"context"
	"fmt"
	"time"

	"github.com/b2b2b-pro/prototype/repository/config"
	"github.com/jackc/pgx/v4/pgxpool"

	"go.uber.org/zap"
)

type PgDB struct {
	Pool *pgxpool.Pool // TODO: по возможности инкапсулировать Pool
}

func New(cfg *config.Config) (*PgDB, error) {
	zap.S().Debug("Инициализирется репозиторий в Postgres\n")

	s := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=public", cfg.DBUser, cfg.DBPass, cfg.HostPG, cfg.PortPG, cfg.DBName)

	pool, err := pgxpool.Connect(context.Background(), s)
	if err != nil {
		// TODO затычка переделать
		zap.S().Debugf("pgxpool.Connect error: %v\n", err)
		time.Sleep(time.Second * 10)
		pool, err = pgxpool.Connect(context.Background(), s)
		if err != nil {
			zap.S().Fatal("Unable to connect to database: ", err, "\n")
		}
	}

	pg := &PgDB{}
	pg.Pool = pool

	return pg, err
}

func (pg *PgDB) Close() {
	pg.Pool.Close()
}
