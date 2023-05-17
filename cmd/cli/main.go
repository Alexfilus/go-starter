package main

import (
	"context"
	"flag"
	"os"

	"github.com/otyang/yasante/config"
	"github.com/otyang/yasante/internal/zample"
	"github.com/urfave/cli/v3"
)

// Just a simple command line app. dont want to add anything much

func main() {
	configFile := flag.String("configFile", "/config.toml", "full path to config file")
	flag.Parse()

	var (
		ctx               = context.Background()
		cfg, log, _, _, _ = config.InitialiseSetup(configFile, ctx)
	)

	app := &cli.App{}
	{
		app.Name = "Hello"
		app.Usage = "how to use my cli ..."
		app.UsageText = "Command help"
		app.Version = "1.0.0"
		app.Copyright = "© Young Otutu"
		app.Commands = zample.RegisterCLIHandlers(ctx, cfg, log)

	}

	// start cmd app
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err.Error())
	}
}
