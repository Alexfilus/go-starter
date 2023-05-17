package bunrepo

import (
	"context"
	"time"

	"github.com/otyang/yasante/config"
	"github.com/otyang/yasante/internal/wallet/entity"
	"github.com/uptrace/bun"
)

var _ entity.IWalletRepository = (*WalletRepository)(nil)

type WalletRepository struct {
	db     bun.IDB
	config *config.Config
}

func NewWalletRepository(db *bun.DB, config *config.Config) *WalletRepository {
	return &WalletRepository{
		db:     db,
		config: config,
	}
}

// NewWithTx returns a clone of WalletRepository, HOWEVER OVERRIDING the dbConnection with a db-Transaction conn
// as the new dbConnection
func (r *WalletRepository) NewWithTx(tx bun.Tx) entity.IWalletRepository {
	return &WalletRepository{
		db:     tx,
		config: r.config,
	}
}

func (r *WalletRepository) GetByID(ctx context.Context, walletID string) (*entity.Wallet, error) {
	wallet := entity.Wallet{}
	err := r.db.
		NewSelect().
		Model(&wallet).
		Where("id = ?", walletID).
		Limit(1).
		Scan(ctx)
	return &wallet, err
}

func (r *WalletRepository) Update(ctx context.Context, wt entity.Wallet) (*entity.Wallet, error) {
	wt.UpdatedAt = time.Now()
	_, err := r.db.NewUpdate().Model(&wt).WherePK().Exec(ctx)
	return &wt, err
}

func (r *WalletRepository) List(ctx context.Context, userID string) ([]entity.Wallet, error) {
	var wallets []entity.Wallet

	err := r.db.NewSelect().
		Model(&wallets).
		Where("user_id = ?", userID).
		OrderExpr("currency ASC").
		Scan(ctx)
	return wallets, err
}

func (r *WalletRepository) Create(ctx context.Context, userID, currencyISO string, isFiat bool) error {
	wallet := entity.NewWallet(currencyISO, userID, isFiat)
	_, err := r.db.NewInsert().Model(&wallet).Ignore().Exec(ctx)
	return err
}

func (r *WalletRepository) Credit(ctx context.Context, ic entity.InitiateCreditParams) (*entity.Wallet, error) {
	wallet, err := r.GetByID(ctx, ic.GetWalletID())
	if err != nil {
		return nil, err
	}

	wlt, err := wallet.Add(ic.Amount, ic.Fee)
	if err != nil {
		return nil, err
	}

	wt, err := r.Update(ctx, *wlt)
	return wt, err
}

func (r *WalletRepository) Debit(ctx context.Context, id entity.InitiateDebitParams) (*entity.Wallet, error) {
	wallet, err := r.GetByID(ctx, id.GetWalletID())
	if err != nil {
		return nil, err
	}

	wlt, err := wallet.Sub(id.Amount, id.Fee, false)
	if err != nil {
		return nil, err
	}

	wt, err := r.Update(ctx, *wlt)
	return wt, err
}

func (r *WalletRepository) Transfer(
	ctx context.Context,
	it entity.InitiateTransferParams,
) (sender, reciepient *entity.Wallet, err error) {
	debit, err := r.Debit(ctx, entity.InitiateDebitParams{
		UserID:      it.FromUserID,
		CurrencyISO: it.CurrencyISO,
		Amount:      it.Amount,
		Fee:         it.Fee,
	})
	if err != nil {
		return nil, nil, err
	}

	credit, err := r.Credit(ctx, entity.InitiateCreditParams{
		UserID:      it.ToUserID,
		CurrencyISO: it.CurrencyISO,
		Amount:      it.Amount,
		Fee:         entity.Zero(),
	})
	return debit, credit, err
}

func (r *WalletRepository) Swap(
	ctx context.Context, is entity.InitiateSwapParams,
) (debitW, creditW *entity.Wallet, err error) {
	debit, err := r.Debit(ctx, entity.InitiateDebitParams{
		UserID:      is.UserID,
		CurrencyISO: is.FromCurrencyISO,
		Amount:      is.FromAmount,
		Fee:         is.Fee,
	})
	if err != nil {
		return nil, nil, err
	}

	credit, err := r.Credit(ctx, entity.InitiateCreditParams{
		UserID:      is.UserID,
		CurrencyISO: is.ToCurrencyISO,
		Amount:      is.ToAmount,
		Fee:         entity.Zero(),
	})
	return debit, credit, err
}
