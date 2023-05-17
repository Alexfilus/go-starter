package entity

import (
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

var (
	ErrInsufficientAvailableBalance = errors.New("withdrawal error insufficient available balance")
	ErrInsufficientPendingBalance   = errors.New("withdrawal error, insufficient pending balance")
	ErrInsufficientLockedBalance    = errors.New("withdrawal error, insufficient locked balance")
	ErrFrozenWallet                 = errors.New("account frozen")
	ErrInvalidAmount                = errors.New("invalid amount: cannot be  zero or negative amount")
)

func GenerateOrGetWalletID(currencyISO, userid string) string {
	return "w_" + userid + "_" + currencyISO
}

type Wallet struct {
	ID               string          `json:"id" bun:"id,pk"`
	UserID           string          `json:"userId" bun:",notnull"`
	CurrencyCode     string          `json:"currencyCode"  bun:",notnull"`
	AvailableBalance decimal.Decimal `json:"availableBalance"  bun:"type:decimal(24,8),notnull"`
	PendingBalance   decimal.Decimal `json:"pendingBalance"  bun:"type:decimal(24,8),notnull"`
	LockedBalance    decimal.Decimal `json:"lockedBalance"  bun:"type:decimal(24,8),notnull"`
	IsFrozen         bool            `json:"isFrozen" bun:",notnull,default:'FALSE'"`
	IsFiat           bool            `json:"isFiat" bun:",notnull"`
	CreatedAt        time.Time       `json:"createdAt" bun:",notnull"`
	UpdatedAt        time.Time       `json:"updatedAt" bun:",notnull"`
}

func NewWallet(currencyISO, userID string, isFiat bool) Wallet {
	return Wallet{
		ID:               GenerateOrGetWalletID(currencyISO, userID),
		UserID:           userID,
		CurrencyCode:     currencyISO,
		AvailableBalance: decimal.NewFromFloat(0),
		PendingBalance:   decimal.NewFromFloat(0),
		LockedBalance:    decimal.NewFromFloat(0),
		IsFrozen:         false,
		IsFiat:           isFiat,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

// Freeze freezes a wallet
func (w *Wallet) Freeze() *Wallet {
	w.IsFrozen = true
	return w
}

// UnFreeze unfreeze a wallet
func (w Wallet) UnFreeze() Wallet {
	w.IsFrozen = false
	return w
}

// TotalBalance Gets the total balance
func (w *Wallet) TotalBalance() decimal.Decimal {
	return w.AvailableBalance.Add(w.PendingBalance).Add(w.LockedBalance)
}

// Add This adds to the balance
func (w *Wallet) Add(amount, fee decimal.Decimal) (*Wallet, error) {
	zero := decimal.NewFromFloat(0)

	if w.IsFrozen {
		fmt.Println("account frozen but lets allow adding money to it")
		// return nil, ErrFrozenWallet
	}

	if amount.LessThanOrEqual(zero) {
		return nil, ErrInvalidAmount // no negative amount.
	}

	if w.AvailableBalance.Add(amount).Sub(fee).LessThan(zero) {
		fmt.Println("fee is more than available balance lets allow it though")
	}
	w.AvailableBalance = w.AvailableBalance.Add(amount).Sub(fee)
	return w, nil
}

// Sub This subtracts from the balance
func (w *Wallet) Sub(amount, fee decimal.Decimal, addToPendingBalance bool) (*Wallet, error) {
	zero := decimal.NewFromFloat(0)

	if w.IsFrozen {
		return nil, ErrFrozenWallet
	}

	if amount.LessThanOrEqual(zero) {
		return w, ErrInvalidAmount
	}

	if w.AvailableBalance.Sub(amount).Sub(fee).LessThan(zero) {
		return nil, ErrInsufficientAvailableBalance
	}

	w.AvailableBalance = w.AvailableBalance.Sub(amount).Sub(fee)

	if addToPendingBalance {
		w.PendingBalance = w.PendingBalance.Add(amount)
	}
	return w, nil
}

// UnPendBalance This removes the amount from the pending balance
func (w *Wallet) UnPendBalance(amount decimal.Decimal) (*Wallet, error) {
	zero := decimal.NewFromFloat(0)

	if w.IsFrozen {
		return nil, ErrFrozenWallet
	}

	if amount.LessThanOrEqual(zero) {
		return w, ErrInvalidAmount
	}

	if w.PendingBalance.Sub(amount).LessThan(zero) {
		return nil, ErrInsufficientPendingBalance
	}
	w.PendingBalance = w.PendingBalance.Sub(amount)
	return w, nil
}

func (w Wallet) LockAmount(amount decimal.Decimal) (Wallet, error) {
	zero := decimal.NewFromFloat(0)

	if amount.LessThanOrEqual(zero) {
		return w, ErrInvalidAmount
	}

	if w.AvailableBalance.Sub(amount).LessThan(zero) {
		return w, ErrInsufficientAvailableBalance
	}

	w.AvailableBalance = w.AvailableBalance.Sub(amount)
	w.LockedBalance = w.LockedBalance.Add(amount)
	return w, nil
}

func (w Wallet) UnLockAmount(amount decimal.Decimal) (Wallet, error) {
	zero := decimal.NewFromFloat(0)

	if amount.LessThanOrEqual(zero) {
		return w, ErrInvalidAmount
	}

	if w.LockedBalance.Sub(amount).LessThan(zero) {
		return w, ErrInsufficientLockedBalance
	}
	w.LockedBalance = w.LockedBalance.Sub(amount)
	w.AvailableBalance = w.AvailableBalance.Add(amount)
	return w, nil
}
