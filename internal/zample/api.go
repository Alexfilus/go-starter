package zample

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
	"github.com/urfave/cli/v3"

	"github.com/otyang/go-starter/internal/event"
	"github.com/otyang/go-starter/internal/zample/entity"
	handlercli "github.com/otyang/go-starter/internal/zample/handler_cli"
	handlerEvent "github.com/otyang/go-starter/internal/zample/handler_event"
	handlerHttp "github.com/otyang/go-starter/internal/zample/handler_http"
	bRepo "github.com/otyang/go-starter/internal/zample/repository/bun"
	"github.com/otyang/go-starter/internal/zample/seeder"

	"github.com/otyang/go-starter/config"
	"github.com/otyang/go-starter/internal/middleware"
	"github.com/otyang/go-starter/pkg/datastore"
	loggers "github.com/otyang/go-starter/pkg/logger"
	"github.com/otyang/go-starter/pkg/pubsub"
)

func RegisterHttpHandlers(
	ctx context.Context, router *fiber.App,
	config *config.Config, log loggers.Interface, db datastore.OrmDB,
) {
	var (
		repo    = bRepo.NewRepository(db)
		handler = handlerHttp.NewHandler(repo, config, log)
	)

	group := router.Group("/zample")
	{
		group.Get("/", handler.Welcome)
		group.Get("/list-via-repo", handler.ListViaRepository)
		group.Get("/list-via-db-helpers", handler.ListViaDbHelpers(db))

		group.Get("/error-500", handler.Error500)
		group.Get("/error-returned", handler.ErrorReturned)

		group.Post(
			"/validation-via-middleware",
			middleware.MwValidateBody[entity.ProfileRequest],
			handler.Profile,
		)
		group.Post("/validation-in-handler", handler.Profile2)
	}
}

func RegisterEventsHandlers(
	ctx context.Context,
	ps pubsub.IEvent,
	config *config.Config,
	log loggers.Interface,
	db datastore.OrmDB,
) {
	var (
		repo    = bRepo.NewRepository(db)
		handler = handlerEvent.NewHandler(repo, config, log)
	)

	{
		ps.Subscribe(context.TODO(), event.SubjectStockUpdates, handler.SubscribeStocks)

		err := ps.Publish(context.TODO(), event.SubjectStockUpdates, &event.Stock{
			Symbol: "GOOG",
			Price:  200,
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func RegisterMigrationAndSeeder(
	ctx context.Context,
	config *config.Config,
	log loggers.Interface,
	db datastore.OrmDB,
) {
	a := datastore.NewDBHelper(db)

	err := a.Transactional(
		ctx,
		func(ctx context.Context, tx bun.Tx) error {
			err := a.NewWithTx(tx).Migrate(ctx, (*entity.User)(nil), (*entity.Stock)(nil))
			if err != nil {
				return err
			}

			return a.NewWithTx(tx).Insert(ctx, &seeder.Users, true)
		},
	)
	if err != nil {
		log.Fatal("seeders: error seeding database " + err.Error())
	}
}

func RegisterCLIHandlers(
	ctx context.Context,
	config *config.Config,
	log loggers.Interface,
) cli.Commands {
	handler := handlercli.NewHandler(config, log)

	return cli.Commands{
		{
			Name:  "web",
			Usage: "Start the web server.",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "name", Value: "young"},
				&cli.BoolFlag{Name: "active", Value: true},
				&cli.IntFlag{Name: "number", Value: 9000},
			},
			Action: handler.GenerateAction(),
		},
	}
}
