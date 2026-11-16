package scanner

import "context"

type fakeModel struct {
	reply string
}

func newFakeModel() *fakeModel {
	return &fakeModel{}
}
func (f *fakeModel) Extract(ctx context.Context, img []byte) (string, error) {
	return f.reply, nil
}

var _ AIModel = (*fakeModel)(nil)
