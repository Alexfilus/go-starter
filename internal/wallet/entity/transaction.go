package entity

import (
	"time"

	"github.com/otyang/go-pkg/utils"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type LinkedTxnDetail struct {
	bun.BaseModel `bun:"table:transaction"`
	Transaction
}

type Transaction struct {
	bun.BaseModel   `bun:"table:transaction"`
	ID              string            `json:"id" bun:"id,pk"`
	UserID          string            `json:"userId" bun:",notnull"`
	WalletID        string            `json:"waletId" bun:",notnull"`
	Currency        string            `json:"currency" bun:",notnull"`
	TotalAmount     decimal.Decimal   `json:"amount" bun:"type:decimal(24,8),notnull"`
	Fee             decimal.Decimal   `json:"fee" bun:"type:decimal(24,8),notnull"`
	Credit          decimal.Decimal   `json:"credit" bun:"type:decimal(24,8),notnull"`
	Debit           decimal.Decimal   `json:"debit" bun:"type:decimal(24,8),notnull"`
	BalanceAfter    decimal.Decimal   `json:"balanceAfter" bun:"type:decimal(24,8),notnull"`
	Type            TransactionType   `json:"type" bun:",notnull"`
	Status          TransactionStatus `json:"status" bun:",notnull"`
	IsThisAReversal bool              `json:"isThisAReversal"  bun:"is_this_a_reversal,notnull,default:FALSE"`
	LinkedTxnID     *string           `json:"linkedTxnID"`
	LinkedTxnDetail *LinkedTxnDetail  `json:"linkedTxnDetail" bun:"rel:has-one,join:linked_txn_id=id"`
	CreatedAt       time.Time         `json:"createdAt" bun:",notnull"`
	UpdatedAt       time.Time         `json:"updatedAt" bun:",notnull"`
}

func newTransactionID() string {
	return "txn_" + utils.RandomID(15)
}

func _sum(no1, no2 decimal.Decimal) decimal.Decimal {
	return no1.Add(no2)
}

func zero() decimal.Decimal {
	return decimal.NewFromFloat(0)
}

func NewTransaction(t Transaction) *Transaction {
	return &Transaction{
		ID:              newTransactionID(),
		UserID:          t.UserID,
		WalletID:        t.WalletID,
		Currency:        t.Currency,
		TotalAmount:     t.TotalAmount,
		Fee:             t.Fee,
		Credit:          t.Credit,
		Debit:           t.Debit,
		BalanceAfter:    t.BalanceAfter,
		Type:            t.Type,
		Status:          t.Status,
		IsThisAReversal: false,
		LinkedTxnID:     t.LinkedTxnID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func (t *Transaction) SetMeta(isThisAReversal bool) *Transaction {
	t.IsThisAReversal = isThisAReversal
	return t
}

func NewTransactionCreditEntry(w Wallet, ic InitiateCreditParams) *Transaction {
	return &Transaction{
		ID:              newTransactionID(),
		UserID:          w.UserID,
		WalletID:        w.ID,
		Currency:        w.CurrencyCode,
		TotalAmount:     _sum(ic.Amount, ic.Fee),
		Fee:             ic.Fee,
		Credit:          ic.Amount,
		Debit:           zero(),
		BalanceAfter:    w.AvailableBalance,
		Type:            ic.Type,
		Status:          ic.Status,
		IsThisAReversal: false,
		LinkedTxnID:     nil,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func NewTransactionDebitEntry(w Wallet, id InitiateDebitParams) *Transaction {
	return &Transaction{
		ID:              newTransactionID(),
		UserID:          w.UserID,
		WalletID:        w.ID,
		Currency:        w.CurrencyCode,
		TotalAmount:     _sum(id.Amount, id.Fee),
		Fee:             id.Fee,
		Credit:          zero(),
		Debit:           id.Amount,
		BalanceAfter:    w.AvailableBalance,
		Type:            id.Type,
		Status:          id.Status,
		IsThisAReversal: false,
		LinkedTxnID:     nil,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func NewTransactionTransferEntry(it InitiateTransferParams, debitedWallet, creditedWallet Wallet) (from, to *Transaction) {
	debitWID, creditWID := it.GetWalletIDs()

	_dID := newTransactionID()
	_cID := newTransactionID()

	debitEntry := &Transaction{
		ID:              _dID,
		UserID:          debitedWallet.UserID,
		WalletID:        debitWID,
		Currency:        debitedWallet.CurrencyCode,
		TotalAmount:     _sum(it.Amount, it.Fee),
		Fee:             it.Fee,
		Credit:          zero(),
		Debit:           it.Amount,
		BalanceAfter:    debitedWallet.AvailableBalance,
		Type:            TransactionTypeTransfer,
		Status:          TransactionStatusSuccess, // always succeful
		IsThisAReversal: false,                    //  cant be reversed.
		LinkedTxnID:     &_cID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	creditEntry := &Transaction{
		ID:              _cID,
		UserID:          creditedWallet.UserID,
		WalletID:        creditWID,
		Currency:        creditedWallet.CurrencyCode,
		TotalAmount:     _sum(it.Amount, it.Fee),
		Fee:             it.Fee,
		Credit:          it.Amount,
		Debit:           zero(),
		BalanceAfter:    creditedWallet.AvailableBalance,
		Type:            TransactionTypeTransfer,
		Status:          TransactionStatusSuccess, // always successful
		IsThisAReversal: false,                    // cant be reversed
		LinkedTxnID:     &_dID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	return debitEntry, creditEntry
}

func NewTransactionSwapEntry(is InitiateSwapParams, debitedWallet, creditedWallet Wallet) (from, to *Transaction) {
	debitWID, creditWID := is.GetWalletIDs()

	_dID := newTransactionID()
	_cID := newTransactionID()

	debitEntry := &Transaction{
		ID:              newTransactionID(),
		UserID:          debitedWallet.UserID,
		WalletID:        debitWID,
		Currency:        debitedWallet.CurrencyCode,
		TotalAmount:     _sum(is.FromAmount, is.Fee),
		Fee:             is.Fee,
		Credit:          zero(),
		Debit:           is.FromAmount,
		BalanceAfter:    debitedWallet.AvailableBalance,
		Type:            TransactionTypeSwap,
		Status:          TransactionStatusSuccess, // always successful
		IsThisAReversal: false,                    // cant be reversed
		LinkedTxnID:     &_cID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	creditEntry := &Transaction{
		ID:              newTransactionID(),
		UserID:          creditedWallet.UserID,
		WalletID:        creditWID,
		Currency:        creditedWallet.CurrencyCode,
		TotalAmount:     _sum(is.ToAmount, is.Fee),
		Fee:             is.Fee,
		Credit:          is.ToAmount,
		Debit:           zero(),
		BalanceAfter:    creditedWallet.AvailableBalance,
		Type:            TransactionTypeSwap,
		Status:          TransactionStatusSuccess, // always successful
		IsThisAReversal: false,                    // cant be reversed
		LinkedTxnID:     &_dID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	return debitEntry, creditEntry
}
