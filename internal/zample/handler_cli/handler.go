package handlercli

import (
	"fmt"

	"github.com/otyang/go-pkg/logger"
	"github.com/otyang/yasante/config"
	"github.com/urfave/cli/v3"
)

type Handler struct {
	Log    logger.Interface
	Config *config.Config
}

func NewHandler(config *config.Config, Log logger.Interface) *Handler {
	return &Handler{
		Log:    Log,
		Config: config,
	}
}

func (h *Handler) GenerateAction() cli.ActionFunc {
	return func(c *cli.Context) error {
		fmt.Println("Starting server")
		fmt.Println(c.String("name"))
		fmt.Println(c.Bool("active"))
		fmt.Println(c.Int("number"))
		fmt.Println("	Hello World		")
		return nil
	}
}
