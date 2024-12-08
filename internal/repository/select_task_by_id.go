package repository

import (
	"database/sql"
	"go_final_project/internal/models"
)

func (r *Repository) SelectTaskById(id int) (models.Task, error) {
	task := models.Task{}
	row := r.db.QueryRow("SELECT * FROM scheduler WHERE id=:id", sql.Named("id", id))

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return task, err
	}

	return task, nil
}
