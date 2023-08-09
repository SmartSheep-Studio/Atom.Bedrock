package main

import (
	"code.smartsheep.studio/atom/bedrock/pkg/config"
	"code.smartsheep.studio/atom/bedrock/pkg/datasource"
	"code.smartsheep.studio/atom/bedrock/pkg/hypertext"
	"code.smartsheep.studio/atom/bedrock/pkg/logger"
	"code.smartsheep.studio/atom/bedrock/pkg/services"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		logger.Module(),
		fx.WithLogger(logger.NewEventLogger),

		config.Module(),
		datasource.Module(),
		services.Module(),
		hypertext.Module(),
	).Run()
}
