package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

func InitDB(dbFile string) (*sql.DB, error) {
	var install bool

	if _, err := os.Open(dbFile); err != nil { // если нет БД файла
		file, err := os.Create(dbFile) // создаем файл
		if err != nil {
			return nil, fmt.Errorf("failed to create db file: %v", err)
		}
		defer file.Close()

		install = true
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open the database: %v", err)
	}

	if install {
		_, err := db.Exec(`
			CREATE TABLE scheduler (
				id INTEGER PRIMARY KEY,
				date VARCHAR(255),
				title VARCHAR(255) NOT NULL,
				comment TEXT,
				repeat VARCHAR(128)
			);
			CREATE INDEX idx_date ON scheduler(date);
		`)
		if err != nil {
			if err := os.Remove(dbFile); err != nil {
				return nil, fmt.Errorf("failed to delete .db file: %v", err)
			}
			return nil, fmt.Errorf("failed to create table: %v", err)
		}
	}
	return db, nil
}
