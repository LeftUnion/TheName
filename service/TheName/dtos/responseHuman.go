package dtos

type ResponseRichHuman struct {
	Id         *int    `db:"id" json:"id"`
	Name       string  `db:"name" json:"name"`
	Surname    string  `db:"surname" json:"surname"`
	Patronymic *string `db:"patronymic" json:"patronymic"`
	Age        int     `db:"age" json:"age"`
	Sex        string  `db:"sex" json:"sex"`
	Nation     string  `db:"nation" json:"nation"`
}
