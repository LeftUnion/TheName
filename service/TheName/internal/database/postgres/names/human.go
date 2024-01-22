package names

import (
	"database/sql"

	"github.com/LeftUnion/theName/service/TheName/internal/database"
	"github.com/LeftUnion/theName/service/TheName/internal/models"
	"github.com/jackc/pgconn"
	"github.com/sirupsen/logrus"
)

// Метод добавления обогащённой информации о людях в бд.
func (n *Names) UpdateHumanById(human models.Human) error {
	// Debug.
	logrus.Debugf("Информация о человеке для обновления: %+v", human)

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
		:name,
		:surname,
		:patronymic,
		:age,
		:sex,
		:nation,
		:last_update
	)
	WHERE id = :id;`

	result, err := n.db.NamedExec(query, human)
	if err != nil {
		logrus.Errorf("Invalid UPDATE the_name_service.humans %s", err)
		if err.(*pgconn.PgError).Code == "23505" { // TODO: is Ok postgres err
			return database.ErrDuplicateKey
		}
		return database.ErrInternalServer
	}

	// Debug.
	logrus.Debugf("Результат: %+v", result)

	affected, err := result.RowsAffected()
	if err != nil {
		logrus.Errorf("Invalid UPDATE the_name_service.humans %s", err)
		return database.ErrInternalServer
	}
	if affected == 0 {
		logrus.Errorf("Invalid UPDATE the_name_service.humans %s", err)
		return database.ErrNoRowsAffected
	}

	// Debug.
	logrus.Debugf("Строк затронуто: %d", affected)

	return nil
}

// Метод получения обогащённой информации о человеке по id.
func (n *Names) GetHumanById(id int) (*models.Human, error) {
	// Debug.
	logrus.Debugf("Id человека для получения: %d", id)

	const query = `
	SELECT * FROM the_name_service.humans
	WHERE id = $1;`

	human := &models.Human{}

	if err := n.db.Get(human, query, id); err != nil {
		logrus.Errorf("Ошибка получения человека по id: %s", err)
		if err == sql.ErrNoRows {
			return nil, database.ErrNoRows
		} else {
			return nil, database.ErrInternalServer
		}
	}

	// Debug.
	logrus.Debugf("Человек с id(%d) получен: %+v", id, *human)

	return human, nil
}

// Метод принимающий id человека для удаления
func (n *Names) DeleteHumanById(id int) error {
	// Debug.
	logrus.Debugf("Id человека для удаления: %d", id)

	const query = `
		DELETE FROM the_name_service.humans
		WHERE id = $1;`

	result, err := n.db.Exec(query, id)
	if err != nil {
		logrus.Errorf("Invalid DELETE the_name_service.humans %s", err)
		return database.ErrInternalServer
	}

	// Debug.
	logrus.Debugf("Результат: %+v", result)

	affected, err := result.RowsAffected()
	if err != nil {
		logrus.Errorf("Invalid DELETE the_name_service.humans %s", err)
		return database.ErrInternalServer
	}
	if affected == 0 {
		logrus.Errorf("Invalid DELETE the_name_service.humans %s", err)
		return database.ErrNoRowsAffected
	}

	// Debug.
	logrus.Debugf("Строк затронуто: %+v", affected)

	return nil
}
