package repository

import (
	"fmt"

	"github.com/duybuidev/sentinel-20/api/internal/model"
	"github.com/jmoiron/sqlx"
)

type LogRepository struct {
	db *sqlx.DB
}

func NewLogRepository(db *sqlx.DB) *LogRepository {
	return &LogRepository{db: db}
}

type LogFilter struct {
	Level   string
	Service string
	Limit   int
}

func (r *LogRepository) FindAll(f LogFilter) ([]model.Log, error) {
	var logs []model.Log
	query := `SELECT * FROM logs WHERE 1=1`
	if f.Level != "" {
		query += fmt.Sprintf(" AND level = '%s'", f.Level)
	}
	if f.Service != "" {
		query += fmt.Sprintf(" AND service = '%s'", f.Service)
	}
	if f.Limit == 0 {
		f.Limit = 100
	}
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT %d", f.Limit)
	err := r.db.Select(&logs, query)
	return logs, err
}

func (r *LogRepository) Insert(l *model.Log) error {
	query := `
		INSERT INTO logs (container_id, level, message, service)
		VALUES (:container_id, :level, :message, :service)`
	_, err := r.db.NamedExec(query, l)
	return err
}

func (r *LogRepository) CountByLevel() (map[string]int, error) {
	rows, err := r.db.Queryx(`SELECT level, COUNT(*) as count FROM logs GROUP BY level`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make(map[string]int)
	for rows.Next() {
		var level string
		var count int
		rows.Scan(&level, &count)
		result[level] = count
	}
	return result, nil
}
