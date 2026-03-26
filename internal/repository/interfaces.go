package repository

import (
	"context"

	"xhs-go-cli/internal/model"
)

type SourceRepository interface {
	Create(ctx context.Context, src *model.Source) error
	List(ctx context.Context, limit int) ([]model.Source, error)
	GetByID(ctx context.Context, id uint) (*model.Source, error)
	IncQueryCount(ctx context.Context, id uint) error
}

type QueryRepository interface {
	Create(ctx context.Context, q *model.GeneratedQuery) error
	List(ctx context.Context, limit int) ([]model.GeneratedQuery, error)
	ListPending(ctx context.Context, limit int) ([]model.GeneratedQuery, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
	Exists(ctx context.Context, sourceID uint, query string) (bool, error)
}

type SearchResultRepository interface {
	Create(ctx context.Context, r *model.SearchResult) error
	ListPending(ctx context.Context, limit int) ([]model.SearchResult, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
}

type DetailRepository interface {
	Create(ctx context.Context, d *model.Detail) error
	ListPending(ctx context.Context, limit int) ([]model.Detail, error)
}

type QualificationRepository interface {
	Create(ctx context.Context, q *model.Qualification) error
	List(ctx context.Context, limit int) ([]model.Qualification, error)
}
