package repository

import (
	"database/sql"
	"go_final_project/internal/models"
)

func (r *Repository) SearchTasks(search string, searchByDate bool) ([]models.Task, error) {
	var tasks = []models.Task{}
	var rows *sql.Rows
	var query string
	if !searchByDate {
		query = `SELECT * FROM scheduler
				 WHERE LOWER(title) LIKE LOWER(:search)
				 OR LOWER(comment) LIKE LOWER(:search)
			  	 ORDER BY date LIMIT 50`
		search = "%" + search + "%"
	} else {
		query = `SELECT * FROM scheduler
				 WHERE date = :search
			     ORDER BY date LIMIT 50`
	}

	rows, err := r.db.Query(query, sql.Named("search", search))
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {
		task := models.Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return tasks, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil

}
