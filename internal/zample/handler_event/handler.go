package handler_nats

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/otyang/go-pkg/logger"
	"github.com/otyang/yasante/config"
	"github.com/otyang/yasante/internal/zample/entity"
)

type Handler struct {
	Log    logger.Interface
	Config *config.Config
	Repo   entity.IRepository
}

func NewHandler(repo entity.IRepository, config *config.Config, Log logger.Interface) *Handler {
	return &Handler{
		Log:    Log,
		Config: config,
		Repo:   repo,
	}
}

func (h *Handler) SubscribeStocks(msg *nats.Msg) {
	fmt.Println("STOCKS: " + string(msg.Subject) + " : " + string(msg.Data))
}
