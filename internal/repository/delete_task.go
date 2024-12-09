package repository

import (
	"database/sql"
)

func (r *Repository) DeleteTask(id int) error {
	query := `DELETE FROM scheduler WHERE id=:id`
	_, err := r.db.Exec(query, sql.Named("id", id))
	if err != nil {
		return err
	}
	return nil
}
