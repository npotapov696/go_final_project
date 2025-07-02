package db

import (
	"database/sql"
	"fmt"
	"time"
)

// Task соответствует полям таблицы scheduler базы данных scheduler.db.
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет в таблицу scheduler базы данных scheduler.db задачу из task.
// Возвращает id добавленной задачи и возможную ошибку.
func AddTask(task *Task) (int64, error) {
	var id int64

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`

	res, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

// Tasks возвращает массив задач и возможную ошибку из таблицы scheduler базы данных scheduler.db.
// В search передается строка для поиска.
// Если строка пуста, возвращаются все задачи.
// Если строка в формате "02.01.2006", возвращаются все задачи с указанной датой.
// в остальных случаях возвращаются задачи, в полях title и/или comment которых присутствует эта строка.
// Количество возвращаемых задач ограничено количеством, переданным в maxEntries.
func Tasks(maxEntries int, search string) ([]*Task, error) {
	var err error
	var tasks []*Task
	var query string
	var date time.Time
	if search == "" {
		query = `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit`
	} else {
		date, err = time.Parse("02.01.2006", search)
		if err == nil {
			query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date LIMIT :limit`

		} else {
			search = "%" + search + "%"
			query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit`

		}
	}
	rows, err := db.Query(query,
		sql.Named("limit", maxEntries),
		sql.Named("search", search),
		sql.Named("date", date.Format(DateString)))
	if err != nil {
		return tasks, err
	}
	defer rows.Close()
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, &task)
	}
	if err := rows.Err(); err != nil {
		return tasks, err
	}

	if tasks == nil {
		tasks = make([]*Task, 0)
	}

	return tasks, nil
}

// GetTask возвращает задачу и возможную ошибку из таблицы scheduler базы данных scheduler.db.
// На вход получает id задачи.
func GetTask(id string) (*Task, error) {
	var task Task
	if id == "" {
		return &task, fmt.Errorf("не указан идентификатор")
	}
	var err error

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id`

	row := db.QueryRow(query, sql.Named("id", id))
	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return &task, fmt.Errorf("задача не найдена")
	}
	return &task, nil
}

// UpdateTask обновляет поля задачи таблицы scheduler базы данных scheduler.db полями задачи task.
// Поиск экземпляра задачи в базе данных в соответствии с id задачи task. Возвращает возможную ошибку.
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET
	date = :date,
	title = :title,
	comment = :comment,
	repeat = :repeat
	WHERE id = :id`

	res, err := db.Exec(query,
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("неверный id для обновления задачи")
	}
	return nil
}

// DeleteTask удаляет задачу таблицы scheduler базы данных scheduler.db по указанному id.
// Возвращает возможную ошибку
func DeleteTask(id string) error {
	if id == "" {
		return fmt.Errorf("не указан идентификатор")
	}

	query := `DELETE FROM scheduler WHERE id = :id`

	_, err := db.Exec(query, sql.Named("id", id))
	if err != nil {
		return err
	}
	return nil
}
