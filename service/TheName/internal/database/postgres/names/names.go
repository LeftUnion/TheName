package names

import (
	"github.com/jmoiron/sqlx"
)

type Names struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) INames {
	return &Names{db: db}
}
