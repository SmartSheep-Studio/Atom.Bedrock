package datasource

import (
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewDb(conf *viper.Viper) *gorm.DB {
	dialector := postgres.Open(conf.GetString("datasource.master.dsn"))
	conn, err := gorm.Open(dialector, &gorm.Config{NamingStrategy: schema.NamingStrategy{
		TablePrefix: conf.GetString("datasource.master.table_prefix"),
	}, Logger: logger.New(&log.Logger, logger.Config{
		Colorful:                  true,
		IgnoreRecordNotFoundError: true,
		LogLevel:                  lo.Ternary(conf.GetBool("debug"), logger.Info, logger.Silent),
	})})

	if err != nil {
		log.Fatal().Err(err).Msg("An error occurred when connect master db")
		return nil
	} else {
		return conn
	}
}
