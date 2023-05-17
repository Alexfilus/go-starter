package config

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/otyang/go-starter/pkg/config"
	"github.com/otyang/go-starter/pkg/datastore"
	"github.com/otyang/go-starter/pkg/logger"
	"github.com/uptrace/bun"
)

var (
	Load     = config.MustLoad
	MustLoad = config.MustLoad
)

type Config struct {
	App struct {
		Name     string `toml:"Name" env:"APP_NAME" env-default:"Auth"`
		Address  string `toml:"Address" env:"APP_ADDRESS" env-default:"0.0.0.0:3000"`
		LogLevel string `toml:"LogLevel" env:"APP_LOG_LEVEL" env-default:"debug"`
	} `toml:"App"`

	DB struct {
		PoolMax              int    `toml:"PoolMax" env:"DB_POOL_MAX" env-default:"1"`
		URL                  string `toml:"URL" env:"DB_URL" env-default:"file::memory:?cache=shared"`
		Driver               string `toml:"Driver" env:"DB_DRIVER" env-default:"sqliteshim"`
		PrintQueriesToStdout bool   `toml:"PrintQueriesToStdout" env:"DB_PRINT_TO_STDOUT"  env-default:"true"`
	} `toml:"DB"`

	Redis struct {
		URL       string `toml:"URL" env:"REDIS_URI" env-default:".." `
		EnableTLS bool   `toml:"EnableTLS" env:"REDIS_ENABLE_TLS" env-default:"true"`
	} `toml:"Redis"`
}

// InitialiseSetup initiates all required data types to use through out program cycle.
// i deliberately refused placing it in a struct so everything is expressed
// in case of any changes made, as struct could hide some uninstantiated values.
// placing it here way makes it easy to set up go tests.
func InitialiseSetup(configFile *string, ctx context.Context,
) (
	*Config,
	*logger.SlogLogger,
	*bun.DB,
	*fiber.App,
) {
	cfg := &Config{}
	{
		if *configFile == "" {
			fmt.Println("Path/to/ConfigFile is required to start application")
			fmt.Println()
			flag.Usage()
			os.Exit(0)
		}
		config.MustLoad(*configFile, cfg)
	}

	log := logger.NewSlogLogger(logger.LevelDebug, logger.DisplayTypeJson, true, os.Stdout)
	logger.WithBuildInfo(log)
	log.With(cfg.App.Name, "address", "port", cfg.App.Address)

	var (
		router = fiber.New(fiber.Config{})
		db     = datastore.NewDBConnection(cfg.DB.Driver, cfg.DB.URL, cfg.DB.PoolMax, cfg.DB.PrintQueriesToStdout)
	)

	return cfg, log, db, router
}
