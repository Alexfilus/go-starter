package entity

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

const (
	TransactionStatusNew     TransactionStatus = "new"
	TransactionStatusPending TransactionStatus = "pending"
	TransactionStatusFailed  TransactionStatus = "failed"
	TransactionStatusSuccess TransactionStatus = "success"
)

type TransactionStatus string

func (ts TransactionStatus) IsValid() bool {
	switch ts {
	case TransactionStatusNew,
		TransactionStatusPending,
		TransactionStatusFailed,
		TransactionStatusSuccess:
		return true
	}
	return false
}

func (ts *TransactionStatus) Scan(src any) error {
	var status string
	switch src := src.(type) {
	case []byte:
		status = string(src)
	case string:
		status = src
	default:
		return fmt.Errorf("unsupported data type: %T", src)
	}

	*ts = TransactionStatus(status)
	return nil
}

func (ts TransactionStatus) Value() (driver.Value, error) {
	if !ts.IsValid() {
		return nil, errors.New("invalid transaction status enum type")
	}
	return string(ts), nil
}

const (
	TransactionTypeSwap     TransactionType = "swap"
	TransactionTypeTransfer TransactionType = "transfer"
)

type TransactionType string

func (tt TransactionType) IsValid() bool {
	switch tt {
	case TransactionTypeSwap, TransactionTypeTransfer:
		return true
	}
	return false
}

func (tt *TransactionType) Scan(src any) error {
	switch src := src.(type) {
	case []byte:
		*tt = TransactionType(string(src))
		return nil
	case string:
		*tt = TransactionType(src)
		return nil
	default:
		return fmt.Errorf("unsupported data type: %T", src)
	}
}

func (tt TransactionType) Value() (driver.Value, error) {
	if !tt.IsValid() {
		return nil, errors.New("tnvalid transaction enum type")
	}
	return string(tt), nil
}

///////

var (
	ErrInvalidTransactionStatus = errors.New("invalid transaction status")
	ErrInvalidTransactionType   = errors.New("invalid transaction type")
)

type InitiateDebitParams struct {
	UserID      string
	CurrencyISO string
	Amount      decimal.Decimal
	Fee         decimal.Decimal
	Type        TransactionType
	Status      TransactionStatus
}

func (i InitiateDebitParams) GetWalletID() string {
	return GenerateOrGetWalletID(i.CurrencyISO, i.UserID)
}

func (i InitiateDebitParams) Validate() error {
	if !i.Type.IsValid() {
		return ErrInvalidTransactionType
	}
	if !i.Status.IsValid() {
		return ErrInvalidTransactionStatus
	}
	return nil
}

type InitiateCreditParams struct {
	UserID      string
	CurrencyISO string
	Amount      decimal.Decimal
	Fee         decimal.Decimal
	Type        TransactionType
	Status      TransactionStatus
}

func (i InitiateCreditParams) GetWalletID() string {
	return GenerateOrGetWalletID(i.CurrencyISO, i.UserID)
}

func (i InitiateCreditParams) Validate() error {
	if !i.Type.IsValid() {
		return ErrInvalidTransactionType
	}
	if !i.Status.IsValid() {
		return ErrInvalidTransactionStatus
	}
	return nil
}

type InitiateSwapParams struct {
	UserID          string
	FromCurrencyISO string
	ToCurrencyISO   string
	FromAmount      decimal.Decimal
	ToAmount        decimal.Decimal
	Fee             decimal.Decimal
}

func (i InitiateSwapParams) GetWalletIDs() (from, to string) {
	return GenerateOrGetWalletID(i.FromCurrencyISO, i.UserID),
		GenerateOrGetWalletID(i.ToCurrencyISO, i.UserID)
}

type InitiateTransferParams struct {
	FromUserID  string
	ToUserID    string
	CurrencyISO string
	Amount      decimal.Decimal
	Fee         decimal.Decimal
}

func (i InitiateTransferParams) GetWalletIDs() (from, to string) {
	return GenerateOrGetWalletID(i.CurrencyISO, i.FromUserID),
		GenerateOrGetWalletID(i.CurrencyISO, i.ToUserID)
}

// /
type ListTransactionsParams struct {
	UserOrWalletID string
	StartDate      time.Time
	EndDate        time.Time
	Limit          int
	Status         *string
	Type           *string
	AscOrder       bool
}

func Zero() decimal.Decimal {
	return decimal.NewFromFloat(0)
}
