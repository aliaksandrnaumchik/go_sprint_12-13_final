package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL,
    comment TEXT,
    repeat VARCHAR(128)
);

CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	var errOpen error
	DB, errOpen = sql.Open("sqlite", dbFile)
	if errOpen != nil {
		return fmt.Errorf("ошибка открытия БД: %w", errOpen)
	}

	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	if install {
		_, err = DB.Exec(schema)
		if err != nil {
			return fmt.Errorf("ошибка создания схемы БД: %w", err)
		}
		log.Println("База данных инициализирована, создана таблица scheduler")
	} else {
		log.Println("Подключено к существующей БД scheduler.db")
	}

	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
