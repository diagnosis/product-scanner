package scanner

import (
	"context"
	"errors"
	"testing"
)

func generateInput() []byte {
	return make([]byte, 256)
}

func Test_Scanner_Success(t *testing.T) {
	f := newFakeStore()
	fm := newFakeModel()
	fm.reply = `{"name":"Pepsi", "calories": 0, "category": "soda"}`
	svc := NewService(f, fm)
	input := generateInput()
	p, err := svc.Scan(context.Background(), input)
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}
	if len(f.products) != 1 {
		t.Fatalf("expected product len is 1 got %d", len(f.products))
	}
	if p.Name != "Pepsi" {
		t.Errorf("expected name is %s got %s", "Pepsi", p.Name)
	}
	if f.products[p.ID].Name != "Pepsi" {
		t.Errorf("expected name is %s got %s", "Pepsi", f.products[p.ID].Name)
	}

}
func Test_Scanner_ReplyNull(t *testing.T) {
	f := newFakeStore()
	fm := newFakeModel()
	fm.reply = "null"
	svc := NewService(f, fm)
	input := generateInput()
	_, err := svc.Scan(context.Background(), input)
	if !errors.Is(err, ErrInvalidModelReply) {
		t.Errorf("expected err %v got %v", ErrInvalidModelReply, err)
	}

}
