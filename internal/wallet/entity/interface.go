package entity

import (
	"context"
	"errors"

	"github.com/otyang/go-pkg/datastore"
	"github.com/uptrace/bun"
)

var ErrIrreversibleTransaction = errors.New("sorry swap or transfer transactions can't be reversed")

func IsWalletError(err error) bool {
	switch {
	case
		errors.Is(err, ErrIrreversibleTransaction),
		errors.Is(err, ErrInsufficientAvailableBalance),
		errors.Is(err, ErrInsufficientPendingBalance),
		errors.Is(err, ErrInsufficientLockedBalance),
		errors.Is(err, ErrFrozenWallet),
		errors.Is(err, ErrInvalidAmount):
		return true
	default:
		return false
	}
}

type (
	OrmDB   = datastore.OrmDB
	OrmDBTx = datastore.OrmDbTx

	IWalletRepository interface {
		NewWithTx(tx bun.Tx) IWalletRepository
		List(ctx context.Context, userID string) ([]Wallet, error)
		Update(ctx context.Context, wt Wallet) (*Wallet, error)
		GetByID(ctx context.Context, walletID string) (*Wallet, error)
		Create(ctx context.Context, userID, currencyISO string, isFiat bool) error
		Credit(ctx context.Context, ic InitiateCreditParams) (*Wallet, error)
		Debit(ctx context.Context, id InitiateDebitParams) (*Wallet, error)
		Swap(ctx context.Context, is InitiateSwapParams) (dbtAcc, crdtAcc *Wallet, err error)
		Transfer(ctx context.Context, it InitiateTransferParams) (sender, reciepient *Wallet, err error)
	}

	ITransactionRepository interface {
		NewWithTx(tx bun.Tx) ITransactionRepository
		Create(ctx context.Context, txn Transaction) (*Transaction, error)
		GetByID(ctx context.Context, txID string) (*Transaction, error)
		Update(ctx context.Context, txn Transaction) (*Transaction, error)
		List(ctx context.Context, arg ListTransactionsParams) (*string, []*Transaction, error)
	}

	IWalletAndTransactionRepository interface {
		DebitAndEntry(ctx context.Context, id InitiateDebitParams) (Transaction, error)
		CreditAndEntry(ctx context.Context, ic InitiateCreditParams) (Transaction, error)
		ReverseTransactionAndEntry(ctx context.Context, txID string) (*Transaction, error)
		SwapAndEntry(ctx context.Context, is InitiateSwapParams) (debit, credit Transaction, err error)
		TransferAndEntry(ctx context.Context, it InitiateTransferParams) (sender, reciepient Transaction, err error)
	}
)
