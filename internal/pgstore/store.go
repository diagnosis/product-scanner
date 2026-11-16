package pgstore

import (
	"context"

	"github.com/diagnosis/product-scanner/internal/scanner"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGStore struct {
	pool *pgxpool.Pool
}

func NewPGStore(pool *pgxpool.Pool) *PGStore {
	return &PGStore{pool: pool}
}

func (P PGStore) Save(ctx context.Context, product *scanner.Product) error {
	//TODO implement me
	panic("implement me")
}

var _ scanner.Store = (*PGStore)(nil)
