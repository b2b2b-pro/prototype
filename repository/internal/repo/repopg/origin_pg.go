package repopg

import (
	"context"

	"github.com/b2b2b-pro/lib/object"
	"go.uber.org/zap"
)

// TODO проверить токен
func (pg *RepoPg) CreateOrigin(tkn string, orn object.Origin) (int, error) {
	var err error
	query := "INSERT INTO origin (description) VALUES ($1) RETURNING (id);"

	zap.S().Debug("postgres CreateOrigin: ", orn, "\n")

	row, err := pg.db.Pool.Query(context.Background(), query, orn.Description)
	if err != nil {
		zap.S().Error("INSERT INTO origin Error: ", err, "\n")
		return 0, err
	}

	defer row.Close()

	row.Next()
	err = row.Scan(&orn.ID)

	return orn.ID, err
}
