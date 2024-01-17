package names

type Names interface {
	GetHumanById(human_id int)
	GetAgeById(human_id int)
	GetSexById(human_id int)
	GetNationById(human_id int)
}
