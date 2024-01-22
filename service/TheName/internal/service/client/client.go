package client

import (
	"errors"
	"time"

	"github.com/LeftUnion/theName/service/TheName/dtos"
	"github.com/LeftUnion/theName/service/TheName/internal/autogen"
	"github.com/LeftUnion/theName/service/TheName/internal/database"
	"github.com/LeftUnion/theName/service/TheName/internal/models"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

// Базовые URL.
const (
	ageURL    = "https://api.agify.io/"
	sexURL    = "https://api.genderize.io/"
	nationURL = "https://api.nationalize.io/"
)

// Ошибки.
var (
	// BadRequest.
	ErrNoRowsAffected = errors.New("Ни одна строка не была затронута")

	// InternalServer.
	ErrGetProbablyAge    = errors.New("Ошибка получения возможного возраста")
	ErrGetProbablySex    = errors.New("Ошибка получения возможного пола")
	ErrGetProbablyNation = errors.New("Ошибка получения возможной нации")
)

type Client struct {
	client *resty.Client
	db     database.IDatabase
}

func New(db database.IDatabase) IClient {
	return &Client{
		client: resty.New(),
		db:     db,
	}
}

func (c *Client) GetHumans(params autogen.GetHumansParams) ([]dtos.ResponseRichHuman, int, error) {
	var (
		defaultLimit  = 10
		defautlOffset = 0
	)

	// Debug.
	logrus.Debugf("Параметры для получения людей: %+v", params)

	if params.Limit == nil {
		params.Limit = &defaultLimit
	}

	if params.Offset == nil {
		params.Offset = &defautlOffset
	}

	humans, err := c.db.GetHumans(params)
	if err != nil {
		logrus.Errorf("Ошибка получения людей по фильтрам: %s", err)
		return nil, 0, handleDatabaseErrors(err)
	}

	// Debug.
	logrus.Debugf("Полученные данные: %+v", humans)

	data := make([]dtos.ResponseRichHuman, 0, *params.Limit)

	for i, human := range humans {
		var some dtos.ResponseRichHuman
		some.Id = &humans[i].Id
		some.Name = human.Name
		some.Surname = human.Surname
		some.Patronymic = human.Patronymic
		some.Age = human.Age
		some.Nation = human.Nation
		some.Sex = human.Sex
		data = append(data, some)
	}

	return data, *params.Limit + *params.Offset, nil
}

func (c *Client) UpdateHumans(humans autogen.UpdateHumansJSONBody) error {
	// Debug.
	logrus.Debugf("Данные для обновления людей: %+v", humans)
	sliceLen := len(humans)
	if sliceLen == 0 {
		logrus.Errorf("Ошибка обновления людей: пустой слайс")
		return ErrBadRequestNoResult
	}

	data := make(autogen.AddHumansJSONBody, sliceLen)
	for i, human := range humans {
		data[i].Name = human.Name
		data[i].Surname = human.Surname
		data[i].Patronymic = human.Patronymic
	}

	enrichHumans, err := c.EnrichHumans(data)
	if err != nil {
		logrus.Errorf("Ошибка обогащения информации о человеке: %s", err)
		return handleDatabaseErrors(err)
	}

	dataToUpdata := make([]models.Human, sliceLen)
	for i, human := range enrichHumans {
		dataToUpdata[i].Id = *humans[i].Id
		dataToUpdata[i].Name = human.Name
		dataToUpdata[i].Surname = human.Surname
		dataToUpdata[i].Patronymic = human.Patronymic
		dataToUpdata[i].Age = human.Age
		dataToUpdata[i].Sex = human.Sex
		dataToUpdata[i].Nation = human.Nation
		dataToUpdata[i].Last_update = time.Now()
	}

	if err := c.db.UpdateHumans(dataToUpdata); err != nil {
		logrus.Errorf("Ошибка обновления людей: %s", err)
		return handleDatabaseErrors(err)
	}

	return nil
}

func (c *Client) DeleteHumans(ids []int) error {
	if err := c.db.DeleteHumans(ids); err != nil {
		logrus.Errorf("Ошибка удаления людей по id: %s", err)
		return handleDatabaseErrors(err)
	}
	return nil
}

func (c *Client) AddHumans(poorHumans autogen.AddHumansJSONBody) error {
	enrichHumans, err := c.EnrichHumans(poorHumans)
	if err != nil {
		logrus.Errorf("Ошибка обогащения информации о человеке: %s", err)
		return handleDatabaseErrors(err)
	}

	if err := c.db.AddHumans(enrichHumans); err != nil {
		logrus.Errorf("Добавление списка людей в базу данных: %s", err)
		return handleDatabaseErrors(err)
	}

	return nil
}

func (c *Client) EnrichHumans(poorHumans autogen.AddHumansJSONBody) ([]models.Human, error) {
	sliceLength := len(poorHumans)
	enrichHumans := make([]models.Human, sliceLength, sliceLength)

	for i, poorHuman := range poorHumans {
		name := poorHuman.Name

		age, err := c.getProbablyAge(name)
		// age := 20
		// err := errors.New("")
		err = nil
		if err != nil {
			return nil, ErrGetProbablyAge
		}

		sex, err := c.getProbablySex(name)
		// sex := "male"
		// err = errors.New("")
		err = nil
		if err != nil {
			return nil, ErrGetProbablySex
		}

		nation, err := c.getProbablyNation(name)
		// nation := "RU"
		// err = errors.New("")
		err = nil
		if err != nil {
			return nil, ErrGetProbablyNation
		}

		enrichHumans[i].Name = name
		enrichHumans[i].Surname = poorHuman.Surname
		enrichHumans[i].Patronymic = poorHuman.Patronymic
		enrichHumans[i].Age = age
		enrichHumans[i].Sex = sex
		enrichHumans[i].Nation = nation
	}

	return enrichHumans, nil
}

func (c *Client) UpdateHuman(human autogen.UpdateHumanJSONBody) error {
	// Debug.
	logrus.Debugf("Информация для обновления информации человека: %+v", human)

	data := make(autogen.AddHumansJSONBody, 1)

	data[0].Name = human.Name
	data[0].Surname = human.Surname
	data[0].Patronymic = human.Patronymic

	// _, err := c.EnrichHumans(data)
	enrichHumans, err := c.EnrichHumans(data)
	if err != nil {
		logrus.Errorf("Ошибка получения обогащённой информации о человеке: %s", err)
		return ErrInternalServer
	}

	forDbHuman := models.Human{}
	forDbHuman.Id = *human.Id
	forDbHuman.Name = data[0].Name
	forDbHuman.Surname = data[0].Surname
	forDbHuman.Patronymic = data[0].Patronymic
	forDbHuman.Age = enrichHumans[0].Age
	forDbHuman.Sex = enrichHumans[0].Sex
	forDbHuman.Nation = enrichHumans[0].Nation
	forDbHuman.Last_update = time.Now()

	if err := c.db.UpdateHumanById(forDbHuman); err != nil {
		logrus.Errorf("Ошибка обновления человека по id: %s", err)
		return handleDatabaseErrors(err)
	}

	return nil
}

// Метод получения информации о человеке по id
func (c *Client) GetHuman(id int) (*dtos.ResponseRichHuman, error) {
	human, err := c.db.GetHumanById(id)

	// Debug.
	logrus.Debugf("Данные человека полученные по id: %+v", human)

	if err != nil {
		logrus.Errorf("Ошибка получения человека по id: %s", err)
		return nil, handleDatabaseErrors(err)
	}

	enrichHuman := &dtos.ResponseRichHuman{
		Id:         &human.Id,
		Name:       human.Name,
		Surname:    human.Surname,
		Patronymic: human.Patronymic,
		Age:        human.Age,
		Sex:        human.Sex,
		Nation:     human.Nation,
	}

	return enrichHuman, nil
}

func (c *Client) DeleteHuman(id int) error {

	// Debug.
	logrus.Debugf("id человека для удаления: %d", id)

	if err := c.db.DeleteHumanById(id); err != nil {
		logrus.Errorf("Ошибка удаления человека по id: %s", err)
		return handleDatabaseErrors(err)
	}

	return nil
}

// Приватный метод для получения возможного возраста по имени.
func (c *Client) getProbablyAge(name string) (int, error) {
	resp, err := c.client.SetBaseURL(ageURL).R().
		SetQueryParams(map[string]string{
			"name": name,
		}).
		SetHeader("Accept", "application/json").
		Get("/")

	if err != nil {
		logrus.Errorf("Ошибка отправки запроса на сервис обогащения возрастом: %s", err)
		return 0, ErrInternalServer
	}

	age, err := validateAge(resp.Body())
	if err != nil {
		return 0, err
	}
	return *age.Age, nil

}

func (c *Client) getProbablySex(name string) (string, error) {
	resp, err := c.client.SetBaseURL(sexURL).R().
		SetQueryParams(map[string]string{
			"name": name,
		}).
		SetHeader("Accept", "application/json").
		Get("/")

	if err != nil {
		logrus.Errorf("Ошибка отправки запроса на сервис обогащения полом: %s", err)
		return "", ErrInternalServer
	}

	sex, err := validateSex(resp.Body())
	if err != nil {
		logrus.Errorf("Ошибка валидации пола человека: %s", err)
		return "", err
	}
	return *sex.Gender, nil
}

func (c *Client) getProbablyNation(name string) (string, error) {
	resp, err := c.client.SetBaseURL(nationURL).R().
		SetQueryParams(map[string]string{
			"name": name,
		}).
		SetHeader("Accept", "application/json").
		Get("/")

	if err != nil {
		logrus.Errorf("Ошибка отправки запроса на сервис обогащения полом: %s", err)
		return "", ErrInternalServer
	}

	nation, err := validateNation(resp.Body())
	if err != nil {
		logrus.Errorf("Ошибка валидации пола человека: %s", err)
		return "", err
	}
	return nation.Country[0].CountryId, nil
}
