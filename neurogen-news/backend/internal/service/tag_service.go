package service

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/model"
	"github.com/neurogen-news/backend/internal/repository"
)

type TagService interface {
	GetAll(ctx context.Context, limit int) ([]model.Tag, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Tag, error)
	GetBySlug(ctx context.Context, slug string) (*model.Tag, error)
	GetPopular(ctx context.Context, limit int) ([]model.Tag, error)
	Search(ctx context.Context, query string, limit int) ([]model.Tag, error)
}

type tagService struct {
	tagRepo repository.TagRepository
	redis   *repository.RedisClient
	logger  *zap.Logger
}

func NewTagService(
	tagRepo repository.TagRepository,
	redis *repository.RedisClient,
	logger *zap.Logger,
) TagService {
	return &tagService{
		tagRepo: tagRepo,
		redis:   redis,
		logger:  logger,
	}
}

func (s *tagService) GetAll(ctx context.Context, limit int) ([]model.Tag, error) {
	return s.tagRepo.GetAll(ctx, limit)
}

func (s *tagService) GetByID(ctx context.Context, id uuid.UUID) (*model.Tag, error) {
	return s.tagRepo.GetByID(ctx, id)
}

func (s *tagService) GetBySlug(ctx context.Context, slug string) (*model.Tag, error) {
	return s.tagRepo.GetBySlug(ctx, slug)
}

func (s *tagService) GetPopular(ctx context.Context, limit int) ([]model.Tag, error) {
	return s.tagRepo.GetPopular(ctx, limit)
}

func (s *tagService) Search(ctx context.Context, query string, limit int) ([]model.Tag, error) {
	return s.tagRepo.Search(ctx, query, limit)
}

