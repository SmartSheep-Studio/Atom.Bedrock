package config

import (
	"go.uber.org/fx"
	"os"

	"github.com/spf13/viper"
)

func Module() fx.Option {
	return fx.Module("config",
		fx.Provide(NewConfigProvider),
		fx.Invoke(AutoMkdir),
	)
}

func AutoMkdir(conf *viper.Viper) error {
	folders := []string{
		conf.GetString("paths.user_contents"),
	}

	for _, folder := range folders {
		stat, err := os.Stat(folder)
		if os.IsNotExist(err) || !stat.IsDir() {
			if err := os.MkdirAll(folder, 0755); err != nil {
				return err
			}
		}
	}

	return nil
}
