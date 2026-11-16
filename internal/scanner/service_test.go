package scanner

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
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

func Test_Scanner_Field_Validation(t *testing.T) {

	cases := []struct {
		name  string
		reply string
		want  *Product
	}{
		{name: "Zero Calories", reply: `{"name": "Pepsi", "calories": 0,  "category": "soda", "weight":"500 mL"}`, want: &Product{
			Category: "soda",
			Name:     "Pepsi",
			Weight:   ptr("500 mL"),
			Calories: ptr(0),
		}},
		{name: "nil calories", reply: `{"name": "Dassani", "calories": null, "category": "water", "weight": "330 mL"}`, want: &Product{
			Category: "water",
			Name:     "Dassani",
			Weight:   ptr("330 mL"),
			Calories: nil,
		}},
		{name: "calories key absent", reply: `{"name":"Four Loko", "category": "soda", "weight": "1 pint"}`, want: &Product{
			Category: "soda",
			Name:     "Four Loko",
			Weight:   ptr("1 pint"),
			Calories: nil,
		}},
		{name: "weight ket absent", reply: `{"name":"Coca Cola", "category": "soda","calories":240}`, want: &Product{
			Category: "soda",
			Calories: ptr(240),
			Name:     "Coca Cola",
			Weight:   nil,
		}},
	}
	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			f := newFakeStore()
			fm := newFakeModel()
			fm.reply = cc.reply
			svc := NewService(f, fm)
			p, err := svc.Scan(context.Background(), generateInput())
			if err != nil {
				t.Fatalf("expected no error got %v", err)
			}
			if !samePtr(p.Calories, cc.want.Calories) {
				t.Errorf("expected %v got %v", show(cc.want.Calories), show(p.Calories))
			}
			if !samePtr(p.Weight, cc.want.Weight) {
				t.Errorf("expected %v got %v", show(cc.want.Weight), show(p.Weight))
			}
			if cc.want.Name != p.Name {
				t.Errorf("expected %s got %s", cc.want.Name, p.Name)
			}
			if cc.want.Category != p.Category {
				t.Errorf("expected %s got %s", cc.want.Category, p.Category)
			}
		})
	}
}

func Test_Scanner_ID(t *testing.T) {
	f := newFakeStore()
	fm := newFakeModel()
	svc := NewService(f, fm)
	fm.reply = `{"name":"Pepsi", "id":"11111111-1111-1111-1111-111111111111"}`
	_, err := svc.Scan(context.Background(), generateInput())
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}
	if f.received.ID != uuid.Nil {
		t.Errorf("expected uuid to be nil")
	}
	if !f.received.CreatedAt.IsZero() {
		t.Errorf("expected created at to be zero got %v", f.received.CreatedAt)
	}
	if !f.received.UpdatedAt.IsZero() {
		t.Errorf("expected updated at to be zero got %v", f.received.UpdatedAt)
	}

}

func ptr[T any](v T) *T {
	return &v
}
func samePtr[T comparable](a, b *T) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
func show[T any](p *T) string {
	if p == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", *p)
}
