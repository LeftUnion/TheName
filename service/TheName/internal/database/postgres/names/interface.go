package names

import (
	"github.com/LeftUnion/theName/service/TheName/internal/autogen"
	"github.com/LeftUnion/theName/service/TheName/internal/models"
)

type INames interface {
	//Human
	DeleteHumanById(id int) error
	UpdateHumanById(humans models.Human) error
	GetHumanById(id int) (*models.Human, error)

	//Humans
	AddHumans(enrichHumans []models.Human) error
	DeleteHumans([]int) error
	UpdateHumans(humans []models.Human) error
	GetHumans(params autogen.GetHumansParams) ([]models.Human, error)
}
