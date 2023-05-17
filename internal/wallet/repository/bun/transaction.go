package bunrepo

import (
	"context"
	"strings"
	"time"

	"github.com/otyang/yasante/config"
	"github.com/otyang/yasante/internal/wallet/entity"
	"github.com/uptrace/bun"
)

var _ entity.ITransactionRepository = (*TransactionRepository)(nil)

type TransactionRepository struct {
	db     bun.IDB
	config *config.Config
}

func NewTransactionRepository(db *bun.DB, config *config.Config) *TransactionRepository {
	return &TransactionRepository{
		db:     db,
		config: config,
	}
}

// NewWithTx returns a clone of TransactionRepository, HOWEVER OVERRIDING the dbConnection with a db-Transaction conn
// as the new dbConnection
func (r *TransactionRepository) NewWithTx(tx bun.Tx) entity.ITransactionRepository {
	return &TransactionRepository{
		db:     tx,
		config: r.config,
	}
}

func (r *TransactionRepository) GetByID(ctx context.Context, txID string) (*entity.Transaction, error) {
	transaction := entity.Transaction{}
	err := r.db.
		NewSelect().
		Model(&transaction).
		Where("id = ?", txID).
		Limit(1).
		Scan(ctx)
	return &transaction, err
}

func (r *TransactionRepository) Update(ctx context.Context, txn entity.Transaction) (*entity.Transaction, error) {
	txn.UpdatedAt = time.Now()
	_, err := r.db.
		NewUpdate().
		Model(&txn).
		WherePK().
		Exec(ctx)
	return &txn, err
}

func (r *TransactionRepository) Create(ctx context.Context, txn entity.Transaction) (*entity.Transaction, error) {
	txn.CreatedAt = time.Now()
	txn.UpdatedAt = time.Now()
	_, err := r.db.NewInsert().Model(&txn).Exec(ctx)
	return &txn, err
}

func (r *TransactionRepository) List(
	ctx context.Context, arg entity.ListTransactionsParams,
) (*string, []*entity.Transaction, error) {
	var txns []*entity.Transaction

	query := r.db.NewSelect().
		Model(&txns).
		Where("wallet_id = ?", arg.UserOrWalletID).
		WhereOr("user_id = ?", arg.UserOrWalletID)

	// if transaction status given lets add to query
	if arg.Status != nil {
		query = query.Where("lower(status) = ?", strings.ToLower(*arg.Status))
	}

	// if transaction type given lets add to query
	if arg.Type != nil {
		query = query.Where("lower(type) = ?", strings.ToLower(*arg.Type))
	}

	query = query.Where("created_at >= ? AND created_at <= ?", arg.StartDate, arg.EndDate)

	if arg.AscOrder {
		query = query.OrderExpr("created_at ASC")
	} else {
		query = query.OrderExpr("created_at DESC")
	}

	err := query.Limit(arg.Limit + 1).Scan(ctx)
	if err != nil {
		return nil, nil, err
	}
	if len(txns) < (arg.Limit + 1) {
		return nil, txns, nil
	}

	txnNextCursor := txns[len(txns)-1]          // last result alone
	txnsWithoutNextCursor := txns[:len(txns)-1] // all without last results
	return &txnNextCursor.ID, txnsWithoutNextCursor, nil
}
