package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go-skeleton/internal/errcode"
	"go-skeleton/internal/model"
	"go-skeleton/internal/service"
	applog "go-skeleton/pkg/log"
	"go-skeleton/pkg/response"
	"go-skeleton/pkg/validator"
)

func init() {
	applog.SetLogger(zap.NewNop())
	validator.InitValidator()
}

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

func setupRouter(repo service.ExampleRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("trace_id", "test-trace")
		c.Next()
	})

	svc := service.NewExampleService(repo)
	h := NewExampleHandler(svc)
	r.POST("/examples", h.Create)
	r.GET("/examples", h.List)
	return r
}

func TestCreateExampleSuccess(t *testing.T) {
	repo := &mockExampleRepo{
		createFunc: func(_ context.Context, example *model.Example) error {
			example.ID = 1
			return nil
		},
	}
	router := setupRouter(repo)

	req := httptest.NewRequest(http.MethodPost, "/examples", strings.NewReader(`{"name":"test-example"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp response.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Code != 0 {
		t.Fatalf("expected code 0, got %d", resp.Code)
	}
}

func TestCreateExampleValidationError(t *testing.T) {
	router := setupRouter(&mockExampleRepo{})

	req := httptest.NewRequest(http.MethodPost, "/examples", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var resp response.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Code != errcode.InvalidParams.Code() {
		t.Fatalf("expected code %d, got %d", errcode.InvalidParams.Code(), resp.Code)
	}
}

func TestCreateExampleDatabaseError(t *testing.T) {
	repo := &mockExampleRepo{
		createFunc: func(_ context.Context, _ *model.Example) error {
			return errors.New("db down")
		},
	}
	router := setupRouter(repo)

	req := httptest.NewRequest(http.MethodPost, "/examples", strings.NewReader(`{"name":"test"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var resp response.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Code != errcode.DatabaseError.Code() {
		t.Fatalf("expected code %d, got %d", errcode.DatabaseError.Code(), resp.Code)
	}
}

func TestListExamplesSuccess(t *testing.T) {
	repo := &mockExampleRepo{
		listFunc: func(_ context.Context, _, _ int) ([]model.Example, int64, error) {
			return []model.Example{{ID: 1, Name: "example1"}}, 1, nil
		},
	}
	router := setupRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/examples?limit=10&offset=0", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var resp response.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Code != 0 {
		t.Fatalf("expected code 0, got %d", resp.Code)
	}
}

func TestListExamplesInvalidQuery(t *testing.T) {
	router := setupRouter(&mockExampleRepo{})

	req := httptest.NewRequest(http.MethodGet, "/examples?limit=-1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var resp response.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Code != errcode.InvalidParams.Code() {
		t.Fatalf("expected validation error code, got %d", resp.Code)
	}
}
