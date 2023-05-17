package bunrepo

import (
	"context"
	"errors"

	"github.com/otyang/yasante/config"
	"github.com/otyang/yasante/internal/wallet/entity"
	"github.com/uptrace/bun"
)

var _ entity.IWalletAndTransactionRepository = (*WalletTransactionRepository)(nil)

type WalletTransactionRepository struct {
	db         *bun.DB
	config     *config.Config
	walletRepo entity.IWalletRepository
	txnRepo    entity.ITransactionRepository
}

func NewWalletTransactionRepo(db *bun.DB, config *config.Config) *WalletTransactionRepository {
	return &WalletTransactionRepository{
		db:         db,
		config:     config,
		walletRepo: NewWalletRepository(db, config),
		txnRepo:    NewTransactionRepository(db, config),
	}
}

func (r *WalletTransactionRepository) CreditAndEntry(
	ctx context.Context, ic entity.InitiateCreditParams,
) (entity.Transaction, error) {
	var txn *entity.Transaction

	if err := ic.Validate(); err != nil {
		return *txn, err
	}

	err := r.db.RunInTx(
		ctx,
		nil,
		func(ctx context.Context, tx bun.Tx) error {
			wallet, err := r.walletRepo.NewWithTx(tx).Credit(ctx, ic)
			if err != nil {
				return err
			}

			txn = entity.NewTransactionCreditEntry(*wallet, ic)

			txn, err = r.txnRepo.NewWithTx(tx).Create(ctx, *txn)
			return err
		})
	return *txn, err
}

func (r *WalletTransactionRepository) DebitAndEntry(
	ctx context.Context, id entity.InitiateDebitParams,
) (entity.Transaction, error) {
	var txn *entity.Transaction

	if err := id.Validate(); err != nil {
		return *txn, err
	}

	err := r.db.RunInTx(
		ctx,
		nil,
		func(ctx context.Context, tx bun.Tx) error {
			wallet, err := r.walletRepo.NewWithTx(tx).Debit(ctx, id)
			if err != nil {
				return err
			}

			txn = entity.NewTransactionDebitEntry(*wallet, id)

			txn, err = r.txnRepo.NewWithTx(tx).Create(ctx, *txn)
			return err
		})
	return *txn, err
}

func (r *WalletTransactionRepository) TransferAndEntry(
	ctx context.Context, it entity.InitiateTransferParams,
) (senderTxn, reciepientTxn entity.Transaction, err error) {
	var fromtxn *entity.Transaction
	var totxn *entity.Transaction

	err = r.db.RunInTx(
		ctx,
		nil,
		func(ctx context.Context, tx bun.Tx) error {
			debitedWallet, creditedWallet, err := r.walletRepo.NewWithTx(tx).Transfer(ctx, it)
			if err != nil {
				return err
			}

			fromtxn, totxn = entity.NewTransactionTransferEntry(it, *debitedWallet, *creditedWallet)
			fromtxn, err = r.txnRepo.NewWithTx(tx).Create(ctx, *fromtxn)

			if err != nil {
				return err
			}

			totxn, err = r.txnRepo.NewWithTx(tx).Create(ctx, *fromtxn)
			return err
		})
	return *fromtxn, *totxn, err
}

func (r *WalletTransactionRepository) SwapAndEntry(
	ctx context.Context, is entity.InitiateSwapParams,
) (debit, credit entity.Transaction, err error) {
	var debitTxn *entity.Transaction
	var creditTxn *entity.Transaction

	err = r.db.RunInTx(
		ctx,
		nil,
		func(ctx context.Context, tx bun.Tx) error {
			debitedWallet, creditedWallet, err := r.walletRepo.NewWithTx(tx).Swap(ctx, is)
			if err != nil {
				return err
			}

			debitTxn, creditTxn = entity.NewTransactionSwapEntry(is, *debitedWallet, *creditedWallet)
			debitTxn, err = r.txnRepo.NewWithTx(tx).Create(ctx, *debitTxn)

			if err != nil {
				return err
			}

			creditTxn, err = r.txnRepo.NewWithTx(tx).Create(ctx, *creditTxn)
			return err
		})

	return *debitTxn, *creditTxn, err
}

func (r *WalletTransactionRepository) ReverseTransactionAndEntry(ctx context.Context, txID string) (*entity.Transaction, error) {
	newTxn := &entity.Transaction{}
	err := r.db.RunInTx(
		ctx,
		nil,
		func(ctx context.Context, tx bun.Tx) error {
			txn2reverse, err := r.txnRepo.NewWithTx(tx).GetByID(ctx, txID)
			if err != nil {
				return err
			}

			switch txn2reverse.Type {
			case entity.TransactionTypeSwap, entity.TransactionTypeTransfer:
				return entity.ErrIrreversibleTransaction
			}

			var txn entity.Transaction

			// DebitTransaction: Lets credit
			if txn2reverse.Debit.GreaterThan(txn2reverse.Credit) {

				ic := entity.InitiateCreditParams{
					UserID:      txn2reverse.UserID,
					CurrencyISO: txn2reverse.Currency,
					Amount:      txn.TotalAmount,
					Fee:         entity.Zero(),
				}
				wallet, err := r.walletRepo.NewWithTx(tx).Credit(ctx, ic)
				if err != nil {
					return err
				}
				txn = *entity.NewTransactionCreditEntry(*wallet, ic).SetMeta(true)
			}

			// CreditTransaction: Lets debit
			if txn2reverse.Credit.GreaterThan(txn2reverse.Debit) {

				id := entity.InitiateDebitParams{
					UserID:      txn2reverse.UserID,
					CurrencyISO: txn2reverse.Currency,
					Amount:      txn.TotalAmount,
					Fee:         entity.Zero(),
				}
				wallet, err := r.walletRepo.NewWithTx(tx).Debit(ctx, id)
				if err != nil {
					return err
				}
				txn = *entity.NewTransactionDebitEntry(*wallet, id).SetMeta(true)
			}

			if txn2reverse.Credit.Equal(txn2reverse.Debit) {
				return errors.New("reversal: odd case credit should never be same as debit")
			}

			// reverse txn
			txn2reverse.Status = entity.TransactionStatusFailed
			if _, err = r.txnRepo.NewWithTx(tx).Update(ctx, *txn2reverse); err != nil {
				return err
			}

			newTxn, err = r.txnRepo.NewWithTx(tx).Update(ctx, txn)
			return err
		})

	return newTxn, err
}
