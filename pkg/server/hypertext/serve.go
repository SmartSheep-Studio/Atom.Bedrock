package hypertext

import (
	view "code.smartsheep.studio/atom/bedrock/packages/bedrock-web"
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/controllers"
	"code.smartsheep.studio/atom/bedrock/pkg/server/services"

	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	flog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var server *fiber.App

func NewHttpServer(cycle fx.Lifecycle, conf *viper.Viper, metrics *services.MetricsService, cop *services.HeLiCoPtErService) *fiber.App {
	// Create app
	server = fiber.New(fiber.Config{
		Prefork:               viper.GetBool("hypertext.advanced.prefork"),
		ProxyHeader:           fiber.HeaderXForwardedFor,
		CaseSensitive:         false,
		StrictRouting:         false,
		DisableStartupMessage: true,
		EnableIPValidation:    true,
		ServerHeader:          "Neutron",
		AppName:               "Neutron v2.0",
		BodyLimit:             viper.GetInt("hypertext.max_body_size"),
	})

	// Apply global middlewares
	server.Use(recover.New())
	server.Use(idempotency.New())
	server.Use(requestid.New())
	server.Use(etag.New())
	server.Use(limiter.New(limiter.Config{
		Max:               20,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))
	server.Use(flog.New(flog.Config{
		Format: "${status} | ${latency} | ${method} ${path}\n",
		Output: log.Logger,
	}))
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodOptions,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))

	cycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info().Msgf("Hypertext transfer protocol server prefork is %s", lo.Ternary(conf.GetBool("hypertext.advanced.prefork"), "enabled", "disabled"))
			log.Info().Msg("Hypertext transfer protocol server is starting...")

			go func() {
				err := server.Listen(conf.GetString("hypertext.bind_addr"))
				if err != nil {
					log.Fatal().Err(err).Msg("An error occurred when start http server.")
				}
			}()

			if conf.GetBool("helicopter.enabled") {
				if err := cop.StartAll(); err != nil {
					log.Err(err).Msg("HeLiCoPtEr start failed...")
				}
			} else {
				log.Info().Msg("HeLiCoPtEr was disabled.")
			}

			metrics.IsReady = true
			url := conf.GetString("base_url")
			log.Info().Msgf("Hypertext transfer protocol server is ready on %s", url)

			return nil
		},
	})

	return server
}

func MapControllers(controllers []controllers.HypertextController, server *fiber.App) {
	for _, controller := range controllers {
		controller.Map(server)
	}

	// Handle APIs not found
	server.Get("/api/*", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotFound, "not found")
	})
	// Handle CGIs not found
	server.Get("/cgi/*", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotFound, "not found")
	})

	// Serve static files
	server.Use("/", cache.New(cache.Config{
		Expiration:   24 * time.Hour,
		CacheControl: true,
	}), filesystem.New(filesystem.Config{
		Root:         view.GetHttpFS(),
		Index:        "index.html",
		NotFoundFile: "index.html",
	}))
}
