package postgre

import (
	"github.com/jmoiron/sqlx"
)

type Database struct {
}

func NewDatabasePostgres(db *sqlx.DB) *Database {
	return &Database{}
}
