package model

import (
	"time"
	"github.com/google/uuid"
)

type Container struct {
	ID         uuid.UUID `db:"id"          json:"id"`
	Name       string    `db:"name"        json:"name"`
	Image      string    `db:"image"       json:"image"`
	Status     string    `db:"status"      json:"status"`
	CPUPercent float64   `db:"cpu_percent" json:"cpu_percent"`
	MemUsage   int64     `db:"mem_usage"   json:"mem_usage"`
	MemLimit   int64     `db:"mem_limit"   json:"mem_limit"`
	CreatedAt  time.Time `db:"created_at"  json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"  json:"updated_at"`
}

type Log struct {
	ID          int64     `db:"id"           json:"id"`
	ContainerID uuid.UUID `db:"container_id" json:"container_id"`
	Level       string    `db:"level"        json:"level"`
	Message     string    `db:"message"      json:"message"`
	Service     string    `db:"service"      json:"service"`
	CreatedAt   time.Time `db:"created_at"   json:"created_at"`
}

type Event struct {
	ID          int64     `db:"id"           json:"id"`
	ContainerID uuid.UUID `db:"container_id" json:"container_id"`
	Type        string    `db:"type"         json:"type"`
	ExitCode    int       `db:"exit_code"    json:"exit_code"`
	Description string    `db:"description"  json:"description"`
	CreatedAt   time.Time `db:"created_at"   json:"created_at"`
}
