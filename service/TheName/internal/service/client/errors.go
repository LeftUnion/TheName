package client

import (
	"errors"

	"github.com/LeftUnion/theName/service/TheName/internal/database"
)

var (
	ErrInternalServer     = errors.New("упс... Кто-то пролил кофе на серверную стойку...")
	ErrConflict           = errors.New("данные дублируются: ФИО")
	ErrBadRequestNoResult = errors.New("данных по такому запросу нет")
)

func handleDatabaseErrors(err error) error {
	switch err {
	case database.ErrDuplicateKey:
		return ErrConflict
	case database.ErrInternalServer:
		return ErrInternalServer
	case database.ErrNoRowsAffected:
		return ErrBadRequestNoResult
	case database.ErrNoRows:
		return ErrBadRequestNoResult
	}

	return ErrInternalServer
}
