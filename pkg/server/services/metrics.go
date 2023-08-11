package services

import (
	"time"
)

type MetricsService struct {
	IsReady   bool          `json:"is_ready"`
	StartAt   time.Time     `json:"start_at"`
	StartCost time.Duration `json:"start_cost"`
}

func NewMetricsService() *MetricsService {
	return &MetricsService{
		IsReady: false,
		StartAt: time.Now(),
	}
}
