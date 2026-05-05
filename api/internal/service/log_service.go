package service

import (
	"github.com/duybuidev/sentinel-20/api/internal/model"
	"github.com/duybuidev/sentinel-20/api/internal/repository"
)

type LogService struct {
	repo *repository.LogRepository
}

func NewLogService(repo *repository.LogRepository) *LogService {
	return &LogService{repo: repo}
}

func (s *LogService) GetAll(filter repository.LogFilter) ([]model.Log, error) {
	return s.repo.FindAll(filter)
}

func (s *LogService) Insert(l *model.Log) error {
	// Validate level
	validLevels := map[string]bool{
		"INFO": true, "WARNING": true,
		"ERROR": true, "CRITICAL": true,
	}
	if !validLevels[l.Level] {
		l.Level = "INFO"
	}
	return s.repo.Insert(l)
}

func (s *LogService) GetStats() (map[string]interface{}, error) {
	counts, err := s.repo.CountByLevel()
	if err != nil {
		return nil, err
	}

	total := 0
	for _, v := range counts {
		total += v
	}

	return map[string]interface{}{
		"total":    total,
		"by_level": counts,
	}, nil
}
