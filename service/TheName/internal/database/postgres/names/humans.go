package names

import (
	"database/sql"
	"fmt"

	"github.com/LeftUnion/theName/service/TheName/internal/autogen"
	"github.com/LeftUnion/theName/service/TheName/internal/database"
	"github.com/LeftUnion/theName/service/TheName/internal/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/sirupsen/logrus"
)

/*
Метод добавления обогащённой информации о людях в бд.
get: слайс models.Human;
return: nil;
err: ErrInternalServer, ErrDuplicateKey
*/
func (db *Names) AddHumans(humans []models.Human) error {
	// Debug.
	logrus.Debugf("Слайс обогащённой информации о людях: %+v", humans)

	const query = `
	INSERT INTO the_name_service.humans (
		name,
		surname,
		patronymic,
		age,
		nation,
		sex
	) VALUES (
		:name,
		:surname,
		:patronymic,
		:age,
		:nation,
		:sex
	);`

	result, err := db.db.NamedExec(query, humans)
	if err != nil {
		logrus.Errorf("Ошибка вставки в таблицу: %s", err)
		if err.(*pgconn.PgError).Code == "23505" {
			return database.ErrDuplicateKey
		}
		return database.ErrInternalServer
	}

	// Debug.
	logrus.Debugf("Результат: %+v", result)

	affected, err := result.RowsAffected()
	if err != nil {
		logrus.Errorf("Ошибка получения результата: %s", err)
		return database.ErrInternalServer
	}

	// Debug.
	logrus.Debugf("Строк затронуто: %+v", affected)

	return nil
}

/*
 */
func (db *Names) DeleteHumans(ids []int) error {
	// Debug.
	logrus.Debugf("Слайс id для удаления: %+v", ids)

	const query = `
	DELETE FROM the_name_service.humans
	WHERE id = $1;
	`

	tx, err := db.db.Begin()
	if err != nil {
		logrus.Errorf("Ошибка начала транзакции: %s", err)
		return database.ErrInternalServer
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		logrus.Errorf("Ошибка подготовки запроса: %s", err)
		return database.ErrInternalServer
	}

	for _, id := range ids {
		if _, err := stmt.Exec(id); err != nil {
			logrus.Errorf("Ошибка удаления человека по id: %s", err)
			return database.ErrNoRows
		}
	}

	if err := tx.Commit(); err != nil {
		logrus.Errorf("Ошибка подтверждения изменений: %s", err)
		return database.ErrInternalServer
	}

	return nil
}

func (db *Names) UpdateHumans(humans []models.Human) error {
	// Debug.
	logrus.Debugf("Слайс id для обновления: %+v", humans)

	const query = `
	UPDATE the_name_service.humans
	SET (
		name,
		surname,
		patronymic,
		age,
		sex,
		nation,
		last_update
	) = (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7
	)
	WHERE id = $8;`

	tx, err := db.db.Begin()
	if err != nil {
		logrus.Errorf("Ошибка начала транзакции: %s", err)
		return database.ErrInternalServer
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		logrus.Errorf("Ошибка подготовки запроса: %s", err)
		return database.ErrInternalServer
	}

	for _, human := range humans {
		if _, err := stmt.Exec(human.Name, human.Surname, *human.Patronymic, human.Age, human.Sex, human.Nation, human.Last_update, human.Id); err != nil {
			logrus.Errorf("Ошибка обновления человека: %s", err)
			return database.ErrNoRows
		}
	}

	if err := tx.Commit(); err != nil {
		logrus.Errorf("Ошибка подтверждения изменений: %s", err)
		return database.ErrInternalServer
	}

	return nil
}

// TODO: move to other file

type SortOptions struct {
	Field string
	Order string
}

type paramsions struct {
	Id         *int
	Name       *string
	Surname    *string
	Patronymic *string
	Age        *int
	Sex        *string
	Nation     *string
	Limit      *int
	Offset     *int
}

func (db *Names) GetHumans(params autogen.GetHumansParams) ([]models.Human, error) {
	// Debug.
	logrus.Debugf("Фильтры для полуения людей: %+v", params)

	qb := sq.Select("id", "name", "surname", "patronymic", "age", "nation", "sex").From("the_name_service.humans").PlaceholderFormat(sq.Dollar)
	if params.SortBy != nil && params.SortScale != nil {
		qb = qb.OrderBy(fmt.Sprintf("%s %s", *params.SortBy, *params.SortScale))
	}

	// фильтрация
	if params.Id != nil {
		qb = qb.Where(sq.Eq{
			"id": *params.Id,
		})
	}
	if params.Name != nil {
		qb = qb.Where(sq.Eq{
			"name": *params.Name,
		})
	}
	if params.Surname != nil {
		qb = qb.Where(sq.Eq{
			"surname": *params.Surname,
		})
	}
	if params.Sex != nil {
		qb = qb.Where(sq.Eq{
			"sex": *params.Sex,
		})
	}
	if params.Age != nil {
		qb = qb.Where(sq.Eq{
			"age": *params.Age,
		})
	}
	if params.Nation != nil {
		qb = qb.Where(sq.Eq{
			"nation": *params.Nation,
		})
	}
	if params.Limit != nil {

		qb = qb.Limit(uint64(*params.Limit))
	}
	if params.Offset != nil {

		qb = qb.Offset((uint64(*params.Offset)))
	}

	query, args, err := qb.ToSql()
	if err != nil {
		logrus.Errorf("Ошибка получения строки запроса: %s", err)
		return nil, database.ErrInternalServer
	}
	// Debug.
	logrus.Debugf("Строка запроса к бд: %s", err)

	rows, err := db.db.Query(query, args...)
	defer rows.Close()
	if err != nil {
		logrus.Errorf("Ошибка получения информации о людях по фильтрам: %s", err)
		if err == sql.ErrNoRows {
			return nil, database.ErrNoRows
		}
		return nil, database.ErrInternalServer
	}

	result := make([]models.Human, 0, *params.Limit)
	for rows.Next() {
		var human models.Human
		err := rows.Scan(
			&human.Id,
			&human.Name,
			&human.Surname,
			&human.Patronymic,
			&human.Age,
			&human.Nation,
			&human.Sex,
		)
		if err != nil {
			logrus.Errorf("Ошибка чтения человека: %s", err)
			return nil, database.ErrInternalServer
		}
		result = append(result, human)
	}

	// Debug.
	logrus.Debugf("Результат: %+v", result)

	return result, nil
}
