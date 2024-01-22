package database

import "errors"

var (
	ErrNoRowsAffected = errors.New("ни одна строка не была затронута")
	ErrInternalServer = errors.New("упс... Кто-то пролил кофе на серверную стойку")
	ErrNoRows         = errors.New("по данному запросу нет данных")
	ErrDuplicateKey   = errors.New("данные с такими ФИО уже есть")
)
