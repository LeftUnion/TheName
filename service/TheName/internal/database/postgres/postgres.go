package postgres

import (
	"database/sql"
	"fmt"

	configs "github.com/LeftUnion/theName/service/TheName/internal/config"
	"github.com/LeftUnion/theName/service/TheName/internal/database"
	"github.com/LeftUnion/theName/service/TheName/internal/database/postgres/names"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const driver = "pgx"

type Postgres struct {
	names.INames
}

func NewPostgres(cfg configs.DataBaseConfig) database.IDatabase {
	return Postgres{
		INames: names.New(sqlx.NewDb(newPostgresDB(cfg), driver)),
	}
}

func newPostgresDB(cfg configs.DataBaseConfig) *sql.DB {
	// Строка подкючения к базе данных.
	pattern := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Name, cfg.Password, cfg.SSLmode)

	// Debug.
	logrus.Debugf("Строка подключения к базе данных: %s", pattern)

	config, err := pgx.ParseConfig(pattern)
	if err != nil {
		logrus.Panicf("Invalid postgres config: %s", err)
	}
	config.Logger = &pgxLogger{}

	connection := stdlib.RegisterConnConfig(config)
	db, err := sql.Open(driver, connection)
	if err != nil {
		logrus.Panicf("Invalid postgres connect: %s", err)
	}

	if err := db.Ping(); err != nil {
		logrus.Panicf("Invalid postgres ping: %s", err)
	}

	return db
}
