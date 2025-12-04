package service

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/model"
	"github.com/neurogen-news/backend/internal/repository"
)

type CategoryService interface {
	GetAll(ctx context.Context) ([]model.Category, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Category, error)
	GetBySlug(ctx context.Context, slug string) (*model.Category, error)
	
	// Subscriptions
	Subscribe(ctx context.Context, userID, categoryID uuid.UUID) error
	Unsubscribe(ctx context.Context, userID, categoryID uuid.UUID) error
	IsSubscribed(ctx context.Context, userID, categoryID uuid.UUID) (bool, error)
	GetUserSubscriptions(ctx context.Context, userID uuid.UUID) ([]model.Category, error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
	redis        *repository.RedisClient
	logger       *zap.Logger
}

func NewCategoryService(
	categoryRepo repository.CategoryRepository,
	redis *repository.RedisClient,
	logger *zap.Logger,
) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
		redis:        redis,
		logger:       logger,
	}
}

func (s *categoryService) GetAll(ctx context.Context) ([]model.Category, error) {
	// TODO: Add caching
	return s.categoryRepo.GetAll(ctx)
}

func (s *categoryService) GetByID(ctx context.Context, id uuid.UUID) (*model.Category, error) {
	return s.categoryRepo.GetByID(ctx, id)
}

func (s *categoryService) GetBySlug(ctx context.Context, slug string) (*model.Category, error) {
	return s.categoryRepo.GetBySlug(ctx, slug)
}

func (s *categoryService) Subscribe(ctx context.Context, userID, categoryID uuid.UUID) error {
	return s.categoryRepo.Subscribe(ctx, userID, categoryID)
}

func (s *categoryService) Unsubscribe(ctx context.Context, userID, categoryID uuid.UUID) error {
	return s.categoryRepo.Unsubscribe(ctx, userID, categoryID)
}

func (s *categoryService) IsSubscribed(ctx context.Context, userID, categoryID uuid.UUID) (bool, error) {
	return s.categoryRepo.IsSubscribed(ctx, userID, categoryID)
}

func (s *categoryService) GetUserSubscriptions(ctx context.Context, userID uuid.UUID) ([]model.Category, error) {
	return s.categoryRepo.GetUserSubscriptions(ctx, userID)
}

