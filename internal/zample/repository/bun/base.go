package bun

import (
	"context"

	"github.com/otyang/go-pkg/datastore"
	"github.com/otyang/yasante/internal/zample/entity"
)

var _ entity.IRepository = (entity.IRepository)(nil)

type Repository struct {
	db datastore.OrmDB
}

func NewRepository(db datastore.OrmDB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUsers(ctx context.Context) ([]*entity.User, error) {
	users := []*entity.User{} // ensures empty result returns empty array []
	err := r.db.NewSelect().Model(&users).Scan(ctx)
	return users, err
}
