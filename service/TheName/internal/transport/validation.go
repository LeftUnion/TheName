package transport

import (
	"errors"

	"github.com/LeftUnion/theName/service/TheName/internal/autogen"
)

var (
	ErrEmptyName    = errors.New("Поле 'name' не может быть пустым")
	ErrEmptySurname = errors.New("Поле 'surname' не может быть пустым")
	ErrNoContent    = errors.New("Нет данных")
	ErrNoId         = errors.New("Поле id неверно")
)

func validateUpdateHuman(human autogen.UpdateHumanJSONBody) error {

	if human.Id == nil && *human.Id <= 0 {
		return ErrNoId
	}

	if len(human.Name) == 0 {
		return ErrEmptyName
	}

	if len(human.Surname) == 0 {
		return ErrEmptySurname
	}

	return nil
}

func validateUpdateHumans(humans autogen.UpdateHumansJSONBody) error {
	if len(humans) == 0 {
		return ErrNoContent
	}

	for _, human := range humans {
		if human.Id == nil && *human.Id <= 0 {
			return ErrNoId
		}

		if len(human.Name) == 0 {
			return ErrEmptyName
		}

		if len(human.Surname) == 0 {
			return ErrEmptySurname
		}
	}

	return nil
}

func validateHumans(humans autogen.AddHumansJSONBody) error {
	if len(humans) == 0 {
		return ErrNoContent
	}

	for _, human := range humans {
		if len(human.Name) == 0 {
			return ErrEmptyName
		}

		if len(human.Surname) == 0 {
			return ErrEmptySurname
		}
	}

	return nil
}
