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
func Test_Scanner_BadReplies(t *testing.T) {
	f := newFakeStore()
	fm := newFakeModel()
	svc := NewService(f, fm)
	cases := []struct {
		name  string
		reply string
	}{
		{name: "null reply", reply: `null`},
		{name: "empty object", reply: `{}`},
		{name: "empty name", reply: `{"name": ""}`},
		{name: "truncated json", reply: `{"name":`},
		{name: "refusal text", reply: `Sorry, I can't read this image.`},
		{name: "empty reply", reply: ``},
	}
	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			fm.reply = cc.reply
			input := generateInput()
			_, err := svc.Scan(context.Background(), input)
			if !errors.Is(err, ErrInvalidModelReply) {
				t.Errorf("expected err %v got %v", ErrInvalidModelReply, err)
			}
			if len(f.products) != 0 {
				t.Fatalf("expected 0 got %d", len(f.products))
			}
		})
	}

}

func Test_Scanner_ModelError(t *testing.T) {
	f := newFakeStore()
	fm := newFakeModel()
	fm.err = errors.New("connection refused")
	fm.reply = `{"name":"Pepsi", "calories": 0, "category": "soda"}`
	svc := NewService(f, fm)
	_, err := svc.Scan(context.Background(), generateInput())
	if !errors.Is(err, ErrModelUnavailable) {
		t.Errorf("expected err %v got %v", ErrInvalidModelReply, err)
	}
}
