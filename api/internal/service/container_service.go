package service

import (
	"github.com/duybuidev/sentinel-20/api/internal/model"
	"github.com/duybuidev/sentinel-20/api/internal/repository"
)

type ContainerService struct {
	repo *repository.ContainerRepository
}

func NewContainerService(repo *repository.ContainerRepository) *ContainerService {
	return &ContainerService{repo: repo}
}

func (s *ContainerService) GetAll() ([]model.Container, error) {
	return s.repo.FindAll()
}

func (s *ContainerService) GetByID(id string) (*model.Container, error) {
	return s.repo.FindByID(id)
}

func (s *ContainerService) Upsert(c *model.Container) error {
	return s.repo.Upsert(c)
}

func (s *ContainerService) UpdateStatus(name, status string) error {
	return s.repo.UpdateStatus(name, status)
}

func (s *ContainerService) GetSummary() (map[string]interface{}, error) {
	containers, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	total, running, stopped, dead := 0, 0, 0, 0
	for _, c := range containers {
		total++
		switch c.Status {
		case "running":
			running++
		case "exited":
			stopped++
		case "dead":
			dead++
		}
	}

	return map[string]interface{}{
		"total":   total,
		"running": running,
		"stopped": stopped,
		"dead":    dead,
	}, nil
}
