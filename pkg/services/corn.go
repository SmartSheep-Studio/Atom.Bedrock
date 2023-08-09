package services

import (
	"context"

	"github.com/robfig/cron/v3"
	"go.uber.org/fx"
)

type CronService struct {
	Ctx *cron.Cron
}

func NewCornService(cycle fx.Lifecycle) *CronService {
	service := &CronService{cron.New()}

	cycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			service.Ctx.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return service.Ctx.Stop().Err()
		},
	})

	return service
}
