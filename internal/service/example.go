package service

import (
	"context"

	"go.uber.org/zap"

	"go-skeleton/internal/errcode"
	"go-skeleton/internal/model"
	applog "go-skeleton/pkg/log"
)

// ExampleRepository is the persistence dependency used by ExampleService.
type ExampleRepository interface {
	Create(ctx context.Context, example *model.Example) error
	List(ctx context.Context, limit, offset int) ([]model.Example, int64, error)
}

// ExampleService handles the example application flow.
type ExampleService struct {
	repo ExampleRepository
}

// NewExampleService creates an ExampleService with the given repository.
func NewExampleService(repo ExampleRepository) *ExampleService {
	return &ExampleService{repo: repo}
}

// CreateExampleReq is the request body for creating an example.
type CreateExampleReq struct {
	Name string `json:"name" binding:"required"`
}

// Create creates a new example.
func (s *ExampleService) Create(ctx context.Context, req *CreateExampleReq) (*model.Example, error) {
	example := model.Example{Name: req.Name}
	if err := s.repo.Create(ctx, &example); err != nil {
		applog.FromContext(ctx).Error("failed to create example", zap.Error(err))
		return nil, errcode.DatabaseError
	}
	return &example, nil
}

// ListExamplesReq is the request query for listing examples.
type ListExamplesReq struct {
	Limit  int `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset int `form:"offset" binding:"omitempty,min=0"`
}

// ListExamplesRes is the response for listing examples.
type ListExamplesRes struct {
	Examples []model.Example `json:"examples"`
	Total    int64           `json:"total"`
}

// List returns a paginated list of examples.
func (s *ExampleService) List(ctx context.Context, req *ListExamplesReq) (*ListExamplesRes, error) {
	if req.Limit == 0 {
		req.Limit = 20
	}
	examples, total, err := s.repo.List(ctx, req.Limit, req.Offset)
	if err != nil {
		applog.FromContext(ctx).Error("failed to list examples", zap.Error(err))
		return nil, errcode.DatabaseError
	}

	return &ListExamplesRes{Examples: examples, Total: total}, nil
}
