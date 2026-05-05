package repository

import (
	"github.com/duybuidev/sentinel-20/api/internal/model"
	"github.com/jmoiron/sqlx"
)

type ContainerRepository struct {
	db *sqlx.DB
}

func NewContainerRepository(db *sqlx.DB) *ContainerRepository {
	return &ContainerRepository{db: db}
}

func (r *ContainerRepository) FindAll() ([]model.Container, error) {
	var containers []model.Container
	query := `SELECT * FROM containers ORDER BY updated_at DESC`
	err := r.db.Select(&containers, query)
	return containers, err
}

func (r *ContainerRepository) FindByID(id string) (*model.Container, error) {
	var container model.Container
	query := `SELECT * FROM containers WHERE id = $1`
	err := r.db.Get(&container, query, id)
	return &container, err
}

func (r *ContainerRepository) Upsert(c *model.Container) error {
	query := `
		INSERT INTO containers (name, image, status, cpu_percent, mem_usage, mem_limit)
		VALUES (:name, :image, :status, :cpu_percent, :mem_usage, :mem_limit)
		ON CONFLICT (name) DO UPDATE SET
			status      = EXCLUDED.status,
			cpu_percent = EXCLUDED.cpu_percent,
			mem_usage   = EXCLUDED.mem_usage,
			mem_limit   = EXCLUDED.mem_limit,
			updated_at  = NOW()
		RETURNING id`
	_, err := r.db.NamedExec(query, c)
	return err
}

func (r *ContainerRepository) UpdateStatus(name, status string) error {
	query := `UPDATE containers SET status = $1, updated_at = NOW() WHERE name = $2`
	_, err := r.db.Exec(query, status, name)
	return err
}
func (r *ContainerRepository) FindByName(name string) (*model.Container, error) {
	var container model.Container
	query := `SELECT * FROM containers WHERE name = $1`
	err := r.db.Get(&container, query, name)
	return &container, err
}
