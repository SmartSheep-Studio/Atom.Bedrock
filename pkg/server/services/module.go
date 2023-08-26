package services

import (
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("services",
		fx.Provide(NewMetricsService),
		fx.Provide(NewHeLiCoPtErService),
		fx.Provide(NewUserService),
		fx.Provide(NewAuthService),
		fx.Provide(NewStorageService),
		fx.Provide(NewCronService),
		fx.Provide(NewMailerService),
		fx.Provide(NewOTPService),

		fx.Invoke(func(cron *CronService) {
			// Cleanup at the startup
			cron.CleanDatabase()
		}),
	)
}
