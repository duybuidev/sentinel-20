package repository

import (
	"github.com/duybuidev/sentinel-20/api/internal/model"
	"github.com/jmoiron/sqlx"
)

type EventRepository struct {
	db *sqlx.DB
}

func NewEventRepository(db *sqlx.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) FindRecent(limit int) ([]model.Event, error) {
	var events []model.Event
	if limit == 0 {
		limit = 50
	}
	query := `SELECT * FROM events ORDER BY created_at DESC LIMIT $1`
	err := r.db.Select(&events, query, limit)
	return events, err
}

func (r *EventRepository) Insert(e *model.Event) error {
	query := `
		INSERT INTO events (container_id, type, exit_code, description)
		VALUES (:container_id, :type, :exit_code, :description)`
	_, err := r.db.NamedExec(query, e)
	return err
}
