package repository

import (
	"database/sql"
	"go_final_project/internal/models"
)

func (r *Repository) UpdateTask(task models.Task, id int) error {
	query := `UPDATE scheduler 
			  SET date=:date, title=:title, comment=:comment, repeat=:repeat
			  WHERE id=:id`

	_, err := r.db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", id))
	if err != nil {
		return err
	}

	return nil
}
