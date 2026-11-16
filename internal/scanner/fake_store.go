package scanner

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type fakeStore struct {
	products map[uuid.UUID]Product
}

func newFakeStore() *fakeStore {
	return &fakeStore{make(map[uuid.UUID]Product)}
}

func (f *fakeStore) Save(ctx context.Context, product *Product) error {
	now := time.Now()
	product.ID = uuid.New()
	product.CreatedAt = now
	product.UpdatedAt = now
	f.products[product.ID] = *product
	return nil
}

var _ Store = (*fakeStore)(nil)
