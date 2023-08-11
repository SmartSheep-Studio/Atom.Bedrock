package main

import (
	"code.smartsheep.studio/atom/bedrock/pkg/server/config"
	"code.smartsheep.studio/atom/bedrock/pkg/server/datasource"
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext"
	"code.smartsheep.studio/atom/bedrock/pkg/server/logger"
	"code.smartsheep.studio/atom/bedrock/pkg/server/services"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli"
	"go.uber.org/fx"
	"os"
)

func main() {
	instance := &cli.App{
		Name:    "Bedrock Server",
		Usage:   "Project Atom's Core Component Server",
		Version: "1.0-SNAPSHOT",
		Commands: []cli.Command{
			{
				Name:  "serve",
				Usage: "Start server",
				Action: func(c *cli.Context) error {
					log.Info().Msgf("You are running Bedrock %s!", "SNAPSHOT-1.0")
					fx.New(
						logger.Module(),
						fx.WithLogger(logger.NewEventLogger),

						config.Module(),
						datasource.Module(),
						services.Module(),
						hypertext.Module(),
					).Run()
					return nil
				},
			},
		},
	}

	if err := instance.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("Failed to run bedrock server.")
	}
}
