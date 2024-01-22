package models

import "time"

type Human struct {
	Id          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Surname     string    `db:"surname" json:"surname"`
	Patronymic  *string   `db:"patronymic" json:"patronymic"`
	Age         int       `db:"age" json:"age"`
	Sex         string    `db:"sex" json:"sex"`
	Nation      string    `db:"nation" json:"nation"`
	Created_at  time.Time `db:"created_at" json:"created_at"`
	Last_update time.Time `db:"last_update" json:"last_update"`
}
