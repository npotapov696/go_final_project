// Пакет db реализует подключение к базе данных и содержит функции для взаимодействия с ней.
package db

import (
	"os"

	"github.com/jmoiron/sqlx"

	_ "modernc.org/sqlite"
)

// Schema содержит команду инициализации таблицы scheduler.
const Schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(256) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX scheduler_date ON scheduler (date);
`

// DefaultDbFile содержит путь по умолчанию к базе данных scheduler.db.
var DefaultDbFile = "scheduler.db"

// DateString содержит строковый формат представления даты
var DateString = "20060102"

// Dbinstance является обёрткой над экземпляром обработчика базы данных.
type Dbinstance struct {
	Db *sqlx.DB
}

var DB Dbinstance

// getDbFile возвращает путь к файлу базы данных scheduler.db.
// Если нет переменной среды окружения TODO_DBFILE с актуальным адресом, возвращает значение по умолчанию defaultDbFile.
func getDbFile() string {
	dbFile := DefaultDbFile
	envDbFile := os.Getenv("TODO_DBFILE")
	if len(envDbFile) > 0 {
		dbFile = envDbFile
	}
	return dbFile
}

// Init проверяет наличие файла базы данных scheduler.db по актуальному пути.
// Если файл отсутствует, создает его и таблицу scheduler в нём.
func Init() error {
	dbFile := getDbFile()
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	DB.Db, err = sqlx.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	defer DB.Db.Close()

	if install {
		_, err := DB.Db.Exec(Schema)
		if err != nil {
			return err
		}
	}
	return nil
}
