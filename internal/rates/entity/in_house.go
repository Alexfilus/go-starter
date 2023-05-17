package entity

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

type Currency struct {
	ISOCode           string    `json:"isoCode" bun:",pk"`
	Name              string    `json:"name"`
	Symbol            string    `json:"symbol"`
	IsFiat            bool      `json:"isFiat" bun:",notnull"`
	Precision         int       `json:"precision"`
	LogoURL           string    `json:"logoUrl"`
	BuyRate           float64   `json:"buyRate"`
	SellRate          float64   `json:"sellRate"`
	DepositFee        float64   `json:"depositFee"`
	DepositEnabled    bool      `json:"depositEnabled"`
	WithdrawalFee     float64   `json:"withdrawalFee"`
	WithdrawalEnabled bool      `json:"withdrawalEnabled"`
	MinDepositAmount  float64   `json:"minDepositAmount"`
	MinWithdrawAmount float64   `json:"minWithdrawAmount"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type InhouseRate struct {
	ISOCode   string    `json:"isoCode"`
	Symbol    string    `json:"symbol"`
	Name      string    `json:"name"`
	BuyRate   float64   `json:"buyRate"`
	SellRate  float64   `json:"sellRate"`
	Precision int       `json:"precision"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type InhouseRates []InhouseRate

func (ihr InhouseRates) GetInfo(currencyCode string) (InhouseRate, error) {
	for _, value := range ihr {
		if strings.EqualFold(currencyCode, value.ISOCode) {
			return value, nil
		}
	}
	return InhouseRate{}, fmt.Errorf("unknown currency: %s", currencyCode)
}

func (ihr InhouseRates) UpdateRates(currencyCode string, buyRate, sellRate *float64) error {
	found := false
	for index, value := range ihr {
		if strings.EqualFold(currencyCode, value.ISOCode) {
			found = true

			if buyRate != nil {
				ihr[index].BuyRate = *buyRate
			}

			if sellRate != nil {
				ihr[index].SellRate = *sellRate
			}
		}
	}

	if !found {
		return errors.New("currencycode does not exist")
	}
	return nil
}

func (ihr InhouseRates) calculateFxRate(from, to, base string) (float64, error) {
	// SameCurrency Conversion :-  AmountToConvert * 1
	if strings.EqualFold(from, to) {
		return 1, nil
	}

	// FromBaseCurrency to OthersConversion :
	// - AmountToConvert * SellRate of OTHERCurrency
	if strings.EqualFold(from, base) && !strings.EqualFold(to, base) {
		toCurrency, err := ihr.GetInfo(to)
		return toCurrency.SellRate, err
	}

	// FromOtherCurrency to BaseConversion
	// :- (1 / buyRateOfOtherCurrency)
	if !strings.EqualFold(from, base) && strings.EqualFold(to, base) {
		fromCurrency, err := ihr.GetInfo(from)
		if err != nil {
			return 0, err
		}
		return (1.00 / fromCurrency.BuyRate), err
	}

	// CrossRateCurrenct Conversion :-
	// (1 / buyrate of fromcurrency) * (sellrate of tocurrency)
	if !strings.EqualFold(from, base) && !strings.EqualFold(to, base) {
		fromCurrency, err := ihr.GetInfo(from)
		if err != nil {
			return 0, err
		}
		toCurrency, err := ihr.GetInfo(to)
		return (1 / fromCurrency.BuyRate) * toCurrency.SellRate, err
	}

	return 0, errors.New("odd conversion case: invalid operation")
}

// ToFixed: Rounds up float64 numbers to precision
func (ihr InhouseRates) toFixed(num float64, precision int) float64 {
	round := func(num float64) int {
		return int(num + math.Copysign(0.5, num))
	}

	output := math.Pow(10, float64(precision))

	return float64(round(num*output)) / output
}

type Quote struct {
	From struct {
		Amount   float64
		Currency string
	}
	To struct {
		Amount   float64
		Currency string
	}
	Fee                 float64
	TotalAmountToDeduct float64
	Rate                float64
	Date                time.Time
}

func (ihr InhouseRates) GenerateQuote(from, to string, amount, fee float64, baseCurrency string) (*Quote, error) {
	rate, err := ihr.calculateFxRate(from, to, baseCurrency)
	if err != nil {
		return nil, err
	}

	// Get To details
	infoTo, err := ihr.GetInfo(to)
	if err != nil {
		return nil, err
	}

	// Get from details
	infoFrom, err := ihr.GetInfo(from)
	if err != nil {
		return nil, err
	}

	quote := new(Quote)
	quote.From.Amount = amount
	quote.From.Currency = from
	quote.To.Amount = ihr.toFixed(amount*rate, infoTo.Precision)
	quote.To.Currency = to
	quote.Fee = fee
	quote.TotalAmountToDeduct = ihr.toFixed(amount+fee, infoFrom.Precision)
	quote.Rate = rate
	quote.Date = time.Now()

	return quote, nil
}
