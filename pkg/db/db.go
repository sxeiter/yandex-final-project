package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func Get() *sql.DB {
	return db
}

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
`

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)
	conn, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка при открытии БД: %w", err)
	}
	if err := conn.Ping(); err != nil {
		conn.Close()
		return fmt.Errorf("ошибка при подключении к базе данных: %w", err)
	}

	if install {
		_, err := conn.Exec(schema)
		if err != nil {
			conn.Close()
			return fmt.Errorf("ошибка при создании таблиц: %w", err)
		}
		fmt.Println(" База данных успешно создана и инициализирована.")
	} else {
		fmt.Println("Подключение к существующей базе данных установлено.")
	}

	db = conn
	return nil
}
