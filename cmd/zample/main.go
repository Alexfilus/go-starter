package main

import (
	"context"
	"flag"

	"github.com/otyang/go-starter/config"
	"github.com/otyang/go-starter/internal/middleware"
	"github.com/otyang/go-starter/internal/zample"
)

func main() {
	configFile := flag.String("configFile", "/config.toml", "full path to config file")
	flag.Parse()

	var (
		ctx                  = context.Background()
		cfg, log, db, router = config.InitialiseSetup(configFile, ctx)
	)

	defer db.Close()

	{
		router.Use(
			middleware.TimeTakenToProcessEndpoint(),
			middleware.Cors(),
			middleware.RequestID(),
			middleware.ErrorHandler(log),
		)

		zample.RegisterHttpHandlers(ctx, router, cfg, log, db)
		zample.RegisterMigrationAndSeeder(ctx, cfg, log, db)
	}

	// start go http app
	if err := router.Listen(cfg.App.Address); err != nil {
		log.Fatal("Error starting server: " + err.Error())
	}
}
