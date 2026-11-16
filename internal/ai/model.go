package ai

import (
	"context"

	"github.com/diagnosis/product-scanner/internal/scanner"
)

type AIModel struct {
}

func (A AIModel) Extract(ctx context.Context, img []byte) (string, error) {
	//TODO implement me
	panic("implement me")
}

func NewAIModel() *AIModel {
	return &AIModel{}
}

var _ scanner.AIModel = (*AIModel)(nil)
