package hypertext

import (
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/controllers"
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/middlewares"
	"code.smartsheep.studio/atom/bedrock/pkg/server/services"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"

	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("hypertext",
		fx.Provide(NewHttpServer),

		middlewares.Module(),
		controllers.Module(),

		fx.Invoke(fx.Annotate(MapControllers, fx.ParamTags(`group:"hypertext_controllers"`))),

		fx.Invoke(func(metrics *services.MetricsService, conf *viper.Viper, log zerolog.Logger) {
			metrics.IsReady = true
			metrics.StartCost = time.Since(metrics.StartAt)

			log.Info().Msgf("Your Bedrock instance is ready on: %s", conf.GetString("base_url"))
		}),
	)
}
