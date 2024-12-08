package repository

import (
	"database/sql"
	"strconv"
)

func (r *Repository) AddTask(date, title, comment, repeat string) (string, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) 
			  VALUES (:date, :title, :comment, :repeat)`

	res, err := r.db.Exec(query,
		sql.Named("date", date),
		sql.Named("title", title),
		sql.Named("comment", comment),
		sql.Named("repeat", repeat))
	if err != nil {
		return "", err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}

	return strconv.Itoa(int(id)), nil
}
