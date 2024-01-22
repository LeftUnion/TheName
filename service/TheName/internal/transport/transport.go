package transport

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/LeftUnion/theName/service/TheName/internal/autogen"
	"github.com/LeftUnion/theName/service/TheName/internal/service"
	"github.com/LeftUnion/theName/service/TheName/internal/service/client"
	"github.com/sirupsen/logrus"

	jsoniter "github.com/json-iterator/go"
)

const (
	messageField = "message"
)

const (
	BaseUrl = "http://localhost:3000/api/humans?"
)

var (
	ErrInternalServer = errors.New("Упс... Кажется кто-то пролил кофе на серверную стойку")
)

var (
	errorBadRequstWithExplain = "Bad Request: %s"
)

type Transport struct {
	srv service.IService
}

func NewTransport(srv service.IService) autogen.ServerInterface {
	return &Transport{
		srv: srv,
	}
}

// Удаление информации о людях.
// (DELETE /humans)
func (t *Transport) DeleteHumans(w http.ResponseWriter, r *http.Request, params autogen.DeleteHumansParams) {
	// Debug.
	logrus.Debugf("id для удаления людей: %+v", params)

	rMessage := make(map[string]string, 1)

	// Проверка id.
	for _, id := range params.Ids {
		if id <= 0 {
			rMessage["message"] = "id не может быть <= 0"
			writeJSONResponse(w, rMessage, http.StatusBadRequest)
			return
		}
	}

	if err := t.srv.DeleteHumans(params.Ids); err != nil {
		rMessage["message"] = err.Error()
		switch err {
		case client.ErrInternalServer:
			writeJSONResponse(w, rMessage, http.StatusInternalServerError)
		case client.ErrConflict:
			writeJSONResponse(w, rMessage, http.StatusBadRequest)
		case client.ErrBadRequestNoResult:
			writeJSONResponse(w, rMessage, http.StatusBadRequest)
		default:
			writeJSONResponse(w, rMessage, http.StatusInternalServerError)
		}
		return
	}

	writeOk(w)
}

// Получение информации о людях.
// (GET /humans)
func (t *Transport) GetHumans(w http.ResponseWriter, r *http.Request, params autogen.GetHumansParams) {
	var (
		ascSort          = "ASC"
		descSort         = "DESC"
		defaultSortField = "id"
		defaultSortScale = "ASC"
	)

	// Debug.
	logrus.Debugf("Параметры сортировки и пагинации: %+v", params)

	rMessage := make(map[string]string, 1)

	if params.SortBy == nil {
		params.SortBy = &defaultSortField
	}

	if params.SortScale == nil {
		params.SortScale = &defaultSortScale
	} else {
		upperSortScale := strings.ToUpper(*params.SortScale)
		if upperSortScale != ascSort && upperSortScale != descSort {
			rMessage["message"] = fmt.Sprintf(errorBadRequstWithExplain, "неверные параметры сортировки")
			writeJSONResponse(w, rMessage, http.StatusBadRequest)
		}
	}

	humans, next, err := t.srv.GetHumans(params)
	if err != nil {
		rMessage["message"] = err.Error()
		switch err {
		case client.ErrInternalServer:
			writeJSONResponse(w, rMessage, http.StatusInternalServerError)
		case client.ErrConflict:
			writeJSONResponse(w, rMessage, http.StatusBadRequest)
		case client.ErrBadRequestNoResult:
			writeJSONResponse(w, rMessage, http.StatusBadRequest)
		default:
			writeJSONResponse(w, rMessage, http.StatusInternalServerError)
		}
		return
	}

	nextURL := BaseUrl

	if params.Id != nil {
		nextURL = fmt.Sprintf("%s&id=%d", nextURL, *params.Id)
	}
	if params.Name != nil {
		nextURL = fmt.Sprintf("%s&name=%s", nextURL, *params.Name)
	}
	if params.Surname != nil {
		nextURL = fmt.Sprintf("%s&surname=%s", nextURL, *params.Surname)
	}
	if params.Age != nil {
		nextURL = fmt.Sprintf("%s&age=%d", nextURL, *params.Age)
	}
	if params.Sex != nil {
		nextURL = fmt.Sprintf("%s&sex=%s", nextURL, *params.Sex)
	}
	if params.Nation != nil {
		nextURL = fmt.Sprintf("%s&nation=%s", nextURL, *params.Nation)
	}
	if params.Limit != nil {
		nextURL = fmt.Sprintf("%s&limit=%d", nextURL, *params.Limit)
	}
	if params.Offset != nil {
		nextURL = fmt.Sprintf("%s&offset=%d", nextURL, next)
	}
	if params.SortBy != nil {
		nextURL = fmt.Sprintf("%s&sort_by=%s", nextURL, *params.SortBy)
	}
	if params.SortScale != nil {
		nextURL = fmt.Sprintf("%s&sort_scale=%s", nextURL, *params.SortScale)
	}
	nextUrl := fmt.Sprintf("%s;rel=next", nextURL)
	writeLinkHeader(w, nextUrl)
	writeJSONResponse(w, humans, http.StatusOK)
}

// Создание информации о людях.
// (POST /humans)
func (t *Transport) AddHumans(w http.ResponseWriter, r *http.Request) {
	rMessage := make(map[string]string, 1)

	// Чтение тела запроса.
	data, err := io.ReadAll(r.Body)
	if err != nil {
		rMessage["message"] = ErrInternalServer.Error()
		writeJSONResponse(w, rMessage, http.StatusInternalServerError)
		return
	}

	// Парсинг JSON.
	var humans autogen.AddHumansJSONBody
	if err := jsoniter.Unmarshal(data, &humans); err != nil {
		rMessage["message"] = fmt.Sprintf(errorBadRequstWithExplain, err) // TODO: handle err right
		writeJSONResponse(w, rMessage, http.StatusBadRequest)
		return
	}

	// Валидация массива людей.
	if err := validateHumans(humans); err != nil {
		rMessage["message"] = fmt.Sprintf(errorBadRequstWithExplain, err)
		writeJSONResponse(w, err, http.StatusBadRequest)
		return
	}

	if err := t.srv.AddHumans(humans); err != nil {
		rMessage["message"] = err.Error()
		switch err {
		case client.ErrInternalServer:
			writeJSONResponse(w, rMessage, http.StatusInternalServerError)
		case client.ErrConflict:
			writeJSONResponse(w, rMessage, http.StatusBadRequest)
		case client.ErrBadRequestNoResult:
			writeJSONResponse(w, rMessage, http.StatusBadRequest)
		default:
			writeJSONResponse(w, rMessage, http.StatusInternalServerError)
		}
		return
	}

	writeOk(w)
}

// Обновление информации о людях.
// (PUT /humans)
func (t *Transport) UpdateHumans(w http.ResponseWriter, r *http.Request) {
	rMessage := make(map[string]string, 1)

	// Чтение тела запроса.
	data, err := io.ReadAll(r.Body)
	if err != nil {
		rMessage["message"] = ErrInternalServer.Error()
		writeJSONResponse(w, rMessage, http.StatusInternalServerError)
		return
	}

	// Парсинг JSON.
	var humans autogen.UpdateHumansJSONBody
	if err := jsoniter.Unmarshal(data, &humans); err != nil {
		rMessage["message"] = fmt.Sprintf(errorBadRequstWithExplain, err) // TODO: handle err right
		writeJSONResponse(w, rMessage, http.StatusBadRequest)
		return
	}

	// Валидация массива людей.
	if err := validateUpdateHumans(humans); err != nil {
		rMessage["message"] = fmt.Sprintf(errorBadRequstWithExplain, err)
		writeJSONResponse(w, err, http.StatusBadRequest)
		return
	}

	if err := t.srv.UpdateHumans(humans); err != nil {
		rMessage["message"] = err.Error()
		switch err {
		case client.ErrInternalServer:
			writeJSONResponse(w, rMessage, http.StatusInternalServerError)
		case client.ErrConflict:
			writeJSONResponse(w, rMessage, http.StatusBadRequest)
		case client.ErrBadRequestNoResult:
			writeJSONResponse(w, rMessage, http.StatusBadRequest)
		default:
			writeJSONResponse(w, rMessage, http.StatusInternalServerError)
		}
		return
	}

	writeOk(w)
}

// Удаление информации о человеке.
// (DELETE /humans/{humanId})
func (t *Transport) DeleteHuman(w http.ResponseWriter, r *http.Request, humanId int) {
	// Debug.
	logrus.Debugf("Id пользователя для удаления: %d", humanId)

	rMessage := make(map[string]string, 1)

	// Проверка id.
	if humanId <= 0 {
		rMessage["message"] = "id не может быть <= 0"
		writeJSONResponse(w, rMessage, http.StatusBadRequest)
		return
	}

	// Удаление информации о человеке по id.
	if err := t.srv.DeleteHuman(humanId); err != nil {
		rMessage["message"] = err.Error()
		switch err {
		case client.ErrInternalServer:
			writeJSONResponse(w, rMessage, http.StatusInternalServerError)
		case client.ErrConflict:
			writeJSONResponse(w, rMessage, http.StatusBadRequest)
		case client.ErrBadRequestNoResult:
			writeJSONResponse(w, rMessage, http.StatusBadRequest)
		default:
			writeJSONResponse(w, rMessage, http.StatusInternalServerError)
		}
		return
	}

	writeOk(w)
}

// Получение информации о человеке.
// (GET /humans/{humanId})
func (t *Transport) GetHuman(w http.ResponseWriter, r *http.Request, humanId int) {
	rMessage := make(map[string]string, 1)

	// Проверка id.
	if humanId <= 0 {
		rMessage["message"] = "id не может быть <= 0."
		writeJSONResponse(w, rMessage, http.StatusBadRequest)
		return
	}

	human, err := t.srv.GetHuman(humanId)
	if err != nil {
		rMessage["message"] = err.Error()
		switch err {
		case client.ErrInternalServer:
			writeJSONResponse(w, rMessage, http.StatusInternalServerError)
		case client.ErrConflict:
			writeJSONResponse(w, rMessage, http.StatusBadRequest)
		case client.ErrBadRequestNoResult:
			writeJSONResponse(w, rMessage, http.StatusBadRequest)
		default:
			writeJSONResponse(w, rMessage, http.StatusInternalServerError)
		}
		return
	}

	writeJSONResponse(w, human, http.StatusOK)
}

// Обновление информации о человеке.
// (PUT /humans/{humanId})
func (t *Transport) UpdateHuman(w http.ResponseWriter, r *http.Request, humanId int) {
	rMessage := make(map[string]string, 1)

	// Чтение тела запроса.
	data, err := io.ReadAll(r.Body)
	if err != nil {
		rMessage["message"] = ErrInternalServer.Error()
		writeJSONResponse(w, rMessage, http.StatusInternalServerError)
		return
	}

	// Проверка id.
	if humanId <= 0 {
		rMessage["message"] = "id не может быть <= 0."
		writeJSONResponse(w, rMessage, http.StatusBadRequest)
		return
	}

	// Парсинг JSON.
	var human autogen.UpdateHumanJSONBody
	if err := jsoniter.Unmarshal(data, &human); err != nil {
		rMessage["message"] = fmt.Sprintf(errorBadRequstWithExplain, err) // TODO: handle err right
		writeJSONResponse(w, rMessage, http.StatusBadRequest)
		return
	}
	human.Id = &humanId

	// Валидация массива людей.
	if err := validateUpdateHuman(human); err != nil {
		rMessage["message"] = fmt.Sprintf(errorBadRequstWithExplain, err)
		writeJSONResponse(w, err, http.StatusBadRequest)
		return
	}

	if err := t.srv.UpdateHuman(human); err != nil {
		rMessage["message"] = err.Error()
		switch err {
		case client.ErrInternalServer:
			writeJSONResponse(w, rMessage, http.StatusInternalServerError)
		case client.ErrConflict:
			writeJSONResponse(w, rMessage, http.StatusBadRequest)
		case client.ErrBadRequestNoResult:
			writeJSONResponse(w, rMessage, http.StatusBadRequest)
		default:
			writeJSONResponse(w, rMessage, http.StatusInternalServerError)
		}
		return
	}

	writeOk(w)
}
