package entity

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type ProvidersRate struct {
	BTC  *float64 `json:"BTC"`
	USDT *float64 `json:"USDT"`
	USDC *float64 `json:"USDC"`
}

// This gets the exchange rate for crypto only from fastforex
//
// it return it to exchangeRate response struct
func NewProviderRatesFromFastForest(apiKey, pairs string) (*ProvidersRate, error) {
	url := "https://api.fastforex.io/crypto/fetch-prices?api_key=%s&pairs=%s"
	url = fmt.Sprintf(url, apiKey, pairs)

	type APIResponse struct {
		Error  string `json:"error"`
		Prices struct {
			BTC  float64 `json:"BTC/USD"`
			USDT float64 `json:"USDT/USD"`
			USDC float64 `json:"USDC/USD"`
		} `json:"prices"`
	}

	var (
		success APIResponse
		failure APIResponse
	)

	resp, err := resty.New().R().
		EnableTrace().
		SetResult(&success).
		SetError(&failure).
		Get(url)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API Provider Error: %s", failure.Error)
	}

	return &ProvidersRate{
		BTC:  &success.Prices.BTC,
		USDT: &success.Prices.USDT,
		USDC: &success.Prices.USDC,
	}, nil
}
