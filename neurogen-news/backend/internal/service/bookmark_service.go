package service

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/model"
	"github.com/neurogen-news/backend/internal/repository"
)

type BookmarkService interface {
	Add(ctx context.Context, userID, articleID uuid.UUID, folderID *uuid.UUID) error
	Remove(ctx context.Context, userID, articleID uuid.UUID) error
	IsBookmarked(ctx context.Context, userID, articleID uuid.UUID) (bool, error)
	GetByUser(ctx context.Context, userID uuid.UUID, params BookmarkListParams) (*BookmarkListResult, error)

	// Folders
	GetFolders(ctx context.Context, userID uuid.UUID) ([]model.BookmarkFolder, error)
	CreateFolder(ctx context.Context, userID uuid.UUID, name string) (*model.BookmarkFolder, error)
	UpdateFolder(ctx context.Context, userID uuid.UUID, folderID uuid.UUID, name string) error
	DeleteFolder(ctx context.Context, userID uuid.UUID, folderID uuid.UUID) error
	MoveToFolder(ctx context.Context, userID, articleID uuid.UUID, folderID *uuid.UUID) error
}

type BookmarkListParams struct {
	FolderID *uuid.UUID
	Page     int
	PageSize int
}

type BookmarkListResult struct {
	Items    []model.ArticleCard `json:"items"`
	Total    int                 `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"pageSize"`
	HasMore  bool                `json:"hasMore"`
}

type bookmarkService struct {
	bookmarkRepo repository.BookmarkRepository
	logger       *zap.Logger
}

func NewBookmarkService(
	bookmarkRepo repository.BookmarkRepository,
	logger *zap.Logger,
) BookmarkService {
	return &bookmarkService{
		bookmarkRepo: bookmarkRepo,
		logger:       logger,
	}
}

func (s *bookmarkService) Add(ctx context.Context, userID, articleID uuid.UUID, folderID *uuid.UUID) error {
	return s.bookmarkRepo.Add(ctx, userID, articleID, folderID)
}

func (s *bookmarkService) Remove(ctx context.Context, userID, articleID uuid.UUID) error {
	return s.bookmarkRepo.Remove(ctx, userID, articleID)
}

func (s *bookmarkService) IsBookmarked(ctx context.Context, userID, articleID uuid.UUID) (bool, error) {
	return s.bookmarkRepo.IsBookmarked(ctx, userID, articleID)
}

func (s *bookmarkService) GetByUser(ctx context.Context, userID uuid.UUID, params BookmarkListParams) (*BookmarkListResult, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 50 {
		params.PageSize = 20
	}

	offset := (params.Page - 1) * params.PageSize

	articles, total, err := s.bookmarkRepo.GetByUser(ctx, userID, repository.BookmarkListParams{
		FolderID: params.FolderID,
		Limit:    params.PageSize,
		Offset:   offset,
	})
	if err != nil {
		return nil, err
	}

	return &BookmarkListResult{
		Items:    articles,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
		HasMore:  offset+len(articles) < total,
	}, nil
}

func (s *bookmarkService) GetFolders(ctx context.Context, userID uuid.UUID) ([]model.BookmarkFolder, error) {
	return s.bookmarkRepo.GetFolders(ctx, userID)
}

func (s *bookmarkService) CreateFolder(ctx context.Context, userID uuid.UUID, name string) (*model.BookmarkFolder, error) {
	folder := &model.BookmarkFolder{
		Name:   name,
		UserID: userID,
	}

	if err := s.bookmarkRepo.CreateFolder(ctx, folder); err != nil {
		return nil, err
	}

	return folder, nil
}

func (s *bookmarkService) UpdateFolder(ctx context.Context, userID uuid.UUID, folderID uuid.UUID, name string) error {
	folder := &model.BookmarkFolder{
		ID:     folderID,
		Name:   name,
		UserID: userID,
	}
	return s.bookmarkRepo.UpdateFolder(ctx, folder)
}

func (s *bookmarkService) DeleteFolder(ctx context.Context, userID uuid.UUID, folderID uuid.UUID) error {
	return s.bookmarkRepo.DeleteFolder(ctx, folderID)
}

func (s *bookmarkService) MoveToFolder(ctx context.Context, userID, articleID uuid.UUID, folderID *uuid.UUID) error {
	return s.bookmarkRepo.MoveToFolder(ctx, userID, articleID, folderID)
}

