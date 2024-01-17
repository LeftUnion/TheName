package models

type Human struct {
	Id         int    `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	Surname    string `db:"surname" json:"surname"`
	Patronymic string `db:"patronymic" json:"patronymic"`
	Sex        string `db:"sex" json:"sex"`
	Age        string `db:"age" json:"age"`
	Nation     string `db:"nation" json:"nation"`
}
