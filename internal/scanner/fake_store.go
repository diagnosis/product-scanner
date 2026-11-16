package scanner

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type fakeStore struct {
	products map[uuid.UUID]Product
	received Product
}

func newFakeStore() *fakeStore {
	return &fakeStore{products: make(map[uuid.UUID]Product)}
}

func (f *fakeStore) Save(ctx context.Context, product *Product) error {
	f.received = *product
	now := time.Now()
	product.ID = uuid.New()
	product.CreatedAt = now
	product.UpdatedAt = now
	f.products[product.ID] = *product
	return nil
}

var _ Store = (*fakeStore)(nil)
