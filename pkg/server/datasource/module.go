package datasource

import (
	"code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("datasource",
		fx.Provide(NewDb),
		fx.Invoke(models.RunMigration),
	)
}
