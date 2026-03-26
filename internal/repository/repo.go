package repository

import (
	"context"
	"math/rand"
	"sort"

	"gorm.io/gorm"

	"xhs-go-cli/internal/db"
	"xhs-go-cli/internal/model"
)

type GormSourceRepo struct {
	db *db.DB
}

func NewSourceRepo(db *db.DB) *GormSourceRepo {
	return &GormSourceRepo{db: db}
}

func (r *GormSourceRepo) Create(ctx context.Context, src *model.Source) error {
	return r.db.WithContext(ctx).Create(src).Error
}

func (r *GormSourceRepo) List(ctx context.Context, limit int) ([]model.Source, error) {
	var sources []model.Source
	if err := r.db.WithContext(ctx).Find(&sources).Error; err != nil {
		return nil, err
	}

	// 计算每个 source 的分数并排序
	type scoredSource struct {
		model.Source
		score float64
	}

	// 获取最大 query_count
	maxCount := 0
	for _, s := range sources {
		if s.QueryCount > maxCount {
			maxCount = s.QueryCount
		}
	}
	if maxCount == 0 {
		maxCount = 1
	}

	scored := make([]scoredSource, len(sources))
	for i, s := range sources {
		randomScore := rand.Float64() * 10
		score := float64(s.Priority)*5 + float64(maxCount-s.QueryCount)*3 + randomScore
		scored[i] = scoredSource{Source: s, score: score}
	}

	// 按分数降序排序
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	// 提取排序后的 sources
	result := make([]model.Source, 0, limit)
	for i, s := range scored {
		if limit > 0 && i >= limit {
			break
		}
		result = append(result, s.Source)
	}
	return result, nil
}

func (r *GormSourceRepo) GetByID(ctx context.Context, id uint) (*model.Source, error) {
	var src model.Source
	err := r.db.WithContext(ctx).First(&src, id).Error
	if err != nil {
		return nil, err
	}
	return &src, nil
}

func (r *GormSourceRepo) IncQueryCount(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.Source{}).Where("id = ?", id).UpdateColumn("query_count", gorm.Expr("query_count + ?", 1)).Error
}

type GormQueryRepo struct {
	db *db.DB
}

func NewQueryRepo(db *db.DB) *GormQueryRepo {
	return &GormQueryRepo{db: db}
}

func (r *GormQueryRepo) Create(ctx context.Context, q *model.GeneratedQuery) error {
	return r.db.WithContext(ctx).Create(q).Error
}

func (r *GormQueryRepo) List(ctx context.Context, limit int) ([]model.GeneratedQuery, error) {
	var queries []model.GeneratedQuery
	query := r.db.WithContext(ctx).Order("id ASC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&queries).Error
	return queries, err
}

func (r *GormQueryRepo) ListPending(ctx context.Context, limit int) ([]model.GeneratedQuery, error) {
	var queries []model.GeneratedQuery
	query := r.db.WithContext(ctx).Where("status = ?", "pending").Order("id ASC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&queries).Error
	return queries, err
}

func (r *GormQueryRepo) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).Model(&model.GeneratedQuery{}).Where("id = ?", id).Update("status", status).Error
}

func (r *GormQueryRepo) Exists(ctx context.Context, sourceID uint, query string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.GeneratedQuery{}).Where("source_id = ? AND query = ?", sourceID, query).Count(&count).Error
	return count > 0, err
}

type GormSearchResultRepo struct {
	db *db.DB
}

func NewSearchResultRepo(db *db.DB) *GormSearchResultRepo {
	return &GormSearchResultRepo{db: db}
}

func (r *GormSearchResultRepo) Create(ctx context.Context, result *model.SearchResult) error {
	return r.db.WithContext(ctx).Create(result).Error
}

func (r *GormSearchResultRepo) ListPending(ctx context.Context, limit int) ([]model.SearchResult, error) {
	var results []model.SearchResult
	query := r.db.WithContext(ctx).Order("id ASC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&results).Error
	return results, err
}

type GormDetailRepo struct {
	db *db.DB
}

func NewDetailRepo(db *db.DB) *GormDetailRepo {
	return &GormDetailRepo{db: db}
}

func (r *GormDetailRepo) Create(ctx context.Context, detail *model.Detail) error {
	return r.db.WithContext(ctx).Create(detail).Error
}

func (r *GormDetailRepo) ListPending(ctx context.Context, limit int) ([]model.Detail, error) {
	var details []model.Detail
	query := r.db.WithContext(ctx).Order("id ASC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&details).Error
	return details, err
}

type GormQualificationRepo struct {
	db *db.DB
}

func NewQualificationRepo(db *db.DB) *GormQualificationRepo {
	return &GormQualificationRepo{db: db}
}

func (r *GormQualificationRepo) Create(ctx context.Context, q *model.Qualification) error {
	return r.db.WithContext(ctx).Create(q).Error
}

func (r *GormQualificationRepo) List(ctx context.Context, limit int) ([]model.Qualification, error) {
	var qualifications []model.Qualification
	query := r.db.WithContext(ctx).Order("id ASC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&qualifications).Error
	return qualifications, err
}
