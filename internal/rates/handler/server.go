package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/otyang/yasante/config"
	rates_entity "github.com/otyang/yasante/internal/rates/entity"
	"github.com/uptrace/bun"
)

type Handler struct {
	DBConnForSettings *bun.DB
	Config            *config.Config
}

func NewHandler(config *config.Config, dbConnForSettings *bun.DB) *Handler {
	return &Handler{
		DBConnForSettings: dbConnForSettings,
		Config:            config,
	}
}

func (h *Handler) RestAPI_GetRates(c *fiber.Ctx) error {
	return nil
}

func CronJob_RatesUpdate(h *Handler) error {
	var ihrs rates_entity.InhouseRates

	rfp, err := rates_entity.NewProviderRatesFromFastForest("fastForestApiKey", "BTC/USD,USDT/USD")
	if err != nil {
		return fmt.Errorf("error loading rates from provider fast forex: %s", err.Error())
	}

	ihrs.UpdateRates("BTC", rfp.BTC, rfp.BTC)
	ihrs.UpdateRates("USDT", rfp.USDT, rfp.USDT)

	return nil
}
