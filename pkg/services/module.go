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
		fx.Provide(NewCornService),
		fx.Provide(NewMailerService),
		fx.Provide(NewOTPService),
	)
}
