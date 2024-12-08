package repository

import (
	"go_final_project/internal/models"
)

func (r *Repository) SelectAllTasks() ([]models.Task, error) {
	var tasks = []models.Task{}
	rows, err := r.db.Query("SELECT * FROM scheduler ORDER BY date LIMIT 50")
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
