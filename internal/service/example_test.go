package service

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/zap"

	"go-skeleton/internal/errcode"
	"go-skeleton/internal/model"
	applog "go-skeleton/pkg/log"
)

type mockExampleRepo struct {
	createFunc func(ctx context.Context, example *model.Example) error
	listFunc   func(ctx context.Context, limit, offset int) ([]model.Example, int64, error)
}

func (m *mockExampleRepo) Create(ctx context.Context, example *model.Example) error {
	return m.createFunc(ctx, example)
}

func (m *mockExampleRepo) List(ctx context.Context, limit, offset int) ([]model.Example, int64, error) {
	return m.listFunc(ctx, limit, offset)
}

func init() {
	applog.SetLogger(zap.NewNop())
}

func TestCreateSuccess(t *testing.T) {
	repo := &mockExampleRepo{
		createFunc: func(_ context.Context, example *model.Example) error {
			example.ID = 1
			return nil
		},
	}
	svc := NewExampleService(repo)

	example, err := svc.Create(context.Background(), &CreateExampleReq{Name: "test"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if example.ID != 1 {
		t.Fatalf("expected ID 1, got %d", example.ID)
	}
	if example.Name != "test" {
		t.Fatalf("expected name test, got %q", example.Name)
	}
}

func TestCreateDatabaseError(t *testing.T) {
	repo := &mockExampleRepo{
		createFunc: func(_ context.Context, _ *model.Example) error {
			return errors.New("connection refused")
		},
	}
	svc := NewExampleService(repo)

	_, err := svc.Create(context.Background(), &CreateExampleReq{Name: "test"})
	if err == nil {
		t.Fatal("expected error")
	}

	var ec errcode.Error
	if !errors.As(err, &ec) {
		t.Fatalf("expected errcode.Error, got %T", err)
	}
	if ec.Code() != errcode.DatabaseError.Code() {
		t.Fatalf("expected code %d, got %d", errcode.DatabaseError.Code(), ec.Code())
	}
}

func TestListSuccess(t *testing.T) {
	examples := []model.Example{
		{ID: 1, Name: "a"},
		{ID: 2, Name: "b"},
	}
	repo := &mockExampleRepo{
		listFunc: func(_ context.Context, _, _ int) ([]model.Example, int64, error) {
			return examples, 2, nil
		},
	}
	svc := NewExampleService(repo)

	res, err := svc.List(context.Background(), &ListExamplesReq{Limit: 10, Offset: 0})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if res.Total != 2 {
		t.Fatalf("expected total 2, got %d", res.Total)
	}
	if len(res.Examples) != 2 {
		t.Fatalf("expected 2 examples, got %d", len(res.Examples))
	}
}

func TestListDefaultLimit(t *testing.T) {
	var capturedLimit int
	repo := &mockExampleRepo{
		listFunc: func(_ context.Context, limit, _ int) ([]model.Example, int64, error) {
			capturedLimit = limit
			return nil, 0, nil
		},
	}
	svc := NewExampleService(repo)

	_, err := svc.List(context.Background(), &ListExamplesReq{Limit: 0})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if capturedLimit != 20 {
		t.Fatalf("expected default limit 20, got %d", capturedLimit)
	}
}

func TestListDatabaseError(t *testing.T) {
	repo := &mockExampleRepo{
		listFunc: func(_ context.Context, _, _ int) ([]model.Example, int64, error) {
			return nil, 0, errors.New("timeout")
		},
	}
	svc := NewExampleService(repo)

	_, err := svc.List(context.Background(), &ListExamplesReq{Limit: 10})
	if err == nil {
		t.Fatal("expected error")
	}

	var ec errcode.Error
	if !errors.As(err, &ec) {
		t.Fatalf("expected errcode.Error, got %T", err)
	}
	if ec.Code() != errcode.DatabaseError.Code() {
		t.Fatalf("expected database error code, got %d", ec.Code())
	}
}
