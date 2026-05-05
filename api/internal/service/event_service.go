package service

import (
	"github.com/duybuidev/sentinel-20/api/internal/model"
	"github.com/duybuidev/sentinel-20/api/internal/repository"
)

type EventService struct {
	repo *repository.EventRepository
}

func NewEventService(repo *repository.EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) GetRecent(limit int) ([]model.Event, error) {
	return s.repo.FindRecent(limit)
}

func (s *EventService) Insert(e *model.Event) error {
	validTypes := map[string]bool{
		"started": true, "stopped": true,
		"crashed": true, "restarted": true,
	}
	if !validTypes[e.Type] {
		e.Type = "stopped"
	}
	return s.repo.Insert(e)
}
