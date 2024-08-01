package db

import (
	"github.com/desulaidovich/balance/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func New(c *config.Data) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", c.ConnectionString)

	if err != nil {
		return nil, err
	}

	return db, nil
}
