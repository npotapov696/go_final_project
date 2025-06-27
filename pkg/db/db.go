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

var db *sqlx.DB // инициализация обработчика базы данных.

var envDbFile = os.Getenv("TODO_DBFILE") // Получаем переменную окружения TODO_DBFILE.

// getDbFile возвращает путь к файлу базы данных scheduler.db.
// Если нет переменной среды окружения TODO_DBFILE с актуальным адресом, возвращает значение по умолчанию defaultDbFile.
func getDbFile() string {
	dbFile := DefaultDbFile
	if len(envDbFile) > 0 {
		dbFile = envDbFile
	}
	return dbFile
}

// Init проверяет наличие файла базы данных scheduler.db по актуальному пути.
// Если файл отсутствует, создает его и таблицу scheduler в нём. Устанавливает
// соединение db с этой БД.
func Init() error {
	dbFile := getDbFile()
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err = sqlx.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		_, err := db.Exec(Schema)
		if err != nil {
			return err
		}
	}
	return nil
}

// Close закрывает соединение db с базой данных scheduler.db
func Close() {
	db.Close()
}
