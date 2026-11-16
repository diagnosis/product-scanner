package scanner

import "context"

type fakeModel struct {
	reply string
	err   error
}

func newFakeModel() *fakeModel {
	return &fakeModel{}
}
func (f *fakeModel) Extract(ctx context.Context, img []byte) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	return f.reply, nil
}

var _ AIModel = (*fakeModel)(nil)
