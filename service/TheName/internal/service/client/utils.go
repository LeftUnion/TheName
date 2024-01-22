package client

import (
	"errors"

	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

var (
	ErrZeroCountOfNames = "Количество людей с таким именем 0"
	ErrNoAge            = "У данного имени нет возможного возраста"
	// ErrInternalServer   = " Внутренняя ошибка сервера"
	ErrNoSex         = "У данного имени нет возможного пола"
	ErrNoProbability = "Шанс равен 0"
	ErrNoNation      = "У этого имени нет нации"
)

func validateAge(rawBody []byte) (*Age, error) {
	age := &Age{}

	if err := jsoniter.Unmarshal(rawBody, age); err != nil {
		logrus.Errorf("Невалидное тело ответа сервиса обогащения имени возрастом: %s", err)
		return nil, ErrInternalServer
	}

	// Debug.
	logrus.Debugf("Ответ от сервиса обогащения age: %+v", *age)

	if age.Count == 0 {
		return nil, errors.New(ErrZeroCountOfNames)
	}

	if age.Age == nil {
		return nil, errors.New(ErrNoAge)
	}

	// Debug.
	logrus.Debugf("Возможный возраст %d", *(*age).Age)

	return age, nil
}

func validateSex(rawBody []byte) (*Sex, error) {
	sex := &Sex{}

	if err := jsoniter.Unmarshal(rawBody, sex); err != nil {
		logrus.Errorf("Невалидное тело ответа сервиса обогащения имени полом: %s", err)
		return nil, ErrInternalServer
	}

	// Debug.
	logrus.Debugf("Ответ от сервиса обогащения sex: %+v", *sex)

	if sex.Count == 0 {
		return nil, errors.New(ErrZeroCountOfNames)
	}

	if sex.Probability == 0 {
		return nil, errors.New(ErrNoProbability)
	}

	if sex.Gender == nil {
		return nil, errors.New(ErrNoSex)
	}

	// Debug.
	logrus.Debugf("Возможный пол %s", *(*sex).Gender)

	return sex, nil
}

func validateNation(rawBody []byte) (*Nation, error) {
	nation := &Nation{}

	if err := jsoniter.Unmarshal(rawBody, nation); err != nil {
		return nil, ErrInternalServer
	}

	// Debug.
	logrus.Debugf("Ответ от сервиса обогащения nation: %+v", nation)

	if nation.Count == 0 {
		return nil, errors.New(ErrZeroCountOfNames)
	}

	if len(nation.Country) == 0 {
		return nil, errors.New(ErrNoNation)
	}

	// Debug.
	logrus.Debugf("Возможная национальность %s", (*nation).Country[0].CountryId)

	return nation, nil
}
