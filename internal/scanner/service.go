package scanner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Store interface {
	Save(ctx context.Context, product *Product) error
}

type AIModel interface {
	Extract(ctx context.Context, img []byte) (string, error)
}

type Service struct {
	store Store
	model AIModel
}

func NewService(store Store, model AIModel) *Service {
	return &Service{store: store, model: model}
}

var ErrInvalidModelReply = errors.New("scanner: invalid model reply")
var ErrModelUnavailable = errors.New("scanner: model is not available")

func (s *Service) Scan(ctx context.Context, input []byte) (*Product, error) {
	reply, err := s.model.Extract(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrModelUnavailable, err)
	}
	var p *Product
	if err = json.Unmarshal([]byte(reply), &p); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidModelReply, err)
	}
	if p == nil || strings.TrimSpace(p.Name) == "" {
		return nil, ErrInvalidModelReply
	}
	if err = s.store.Save(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}
