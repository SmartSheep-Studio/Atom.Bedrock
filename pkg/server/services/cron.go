package services

import (
	"context"
	"time"

	"code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type CronService struct {
	Ctx *cron.Cron

	db     *gorm.DB
	logger zerolog.Logger
}

func NewCronService(cycle fx.Lifecycle, db *gorm.DB, logger zerolog.Logger) *CronService {
	service := &CronService{
		cron.New(cron.WithLogger(cron.VerbosePrintfLogger(&logger))),

		db,
		logger,
	}

	service.Ctx.AddFunc("@every 15m", service.CleanDatabase)

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

func (v *CronService) CleanDatabase() {
	targets := []any{
		&models.Lock{},
		&models.User{},
		&models.UserGroup{},
		&models.Notification{},
		&models.UserContact{},
		&models.UserSession{},
		&models.OauthClient{},
		&models.OTP{},
		&models.StorageFile{},
	}

	v.logger.Info().Msg("Now starting cleanning database...")

	start := time.Now()
	for _, target := range targets {
		if stat := v.db.Unscoped().Where("deleted_at IS NOT NULL").Delete(&target); stat.Error != nil {
			v.logger.Error().Err(stat.Error).Msgf("Clean up table %s an error occurred.", stat.Statement.Table)
		} else {
			v.logger.Debug().Msgf("Successfully cleaned table %s.", stat.Statement.Table)
		}
	}

	v.logger.Info().Msgf("Database has been cleaned, took %dms.", time.Since(start))
}
