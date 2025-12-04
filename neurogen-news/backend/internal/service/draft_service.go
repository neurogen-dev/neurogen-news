package service

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/model"
	"github.com/neurogen-news/backend/internal/repository"
)

type DraftService interface {
	Create(ctx context.Context, userID uuid.UUID, input CreateDraftInput) (*model.Draft, error)
	GetByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.Draft, error)
	Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, input UpdateDraftInput) (*model.Draft, error)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) (*DraftListResult, error)

	// Auto-save
	AutoSave(ctx context.Context, userID uuid.UUID, input AutoSaveInput) error
	GetLatestAutoSave(ctx context.Context, userID uuid.UUID, articleID *uuid.UUID) (*model.Draft, error)
}

type CreateDraftInput struct {
	Title         string             `json:"title"`
	Content       string             `json:"content"`
	CoverImageURL *string            `json:"coverImageUrl,omitempty"`
	ContentType   model.ContentType  `json:"contentType"`
	CategoryID    *uuid.UUID         `json:"categoryId,omitempty"`
	Tags          []string           `json:"tags,omitempty"`
}

type UpdateDraftInput struct {
	Title         *string            `json:"title,omitempty"`
	Content       *string            `json:"content,omitempty"`
	CoverImageURL *string            `json:"coverImageUrl,omitempty"`
	ContentType   *model.ContentType `json:"contentType,omitempty"`
	CategoryID    *uuid.UUID         `json:"categoryId,omitempty"`
	Tags          []string           `json:"tags,omitempty"`
}

type AutoSaveInput struct {
	ArticleID     *uuid.UUID         `json:"articleId,omitempty"`
	Title         string             `json:"title"`
	Content       string             `json:"content"`
	CoverImageURL *string            `json:"coverImageUrl,omitempty"`
	ContentType   model.ContentType  `json:"contentType"`
	CategoryID    *uuid.UUID         `json:"categoryId,omitempty"`
	Tags          []string           `json:"tags,omitempty"`
}

type DraftListResult struct {
	Items    []model.Draft `json:"items"`
	Total    int           `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"pageSize"`
	HasMore  bool          `json:"hasMore"`
}

type draftService struct {
	draftRepo repository.DraftRepository
	logger    *zap.Logger
}

func NewDraftService(
	draftRepo repository.DraftRepository,
	logger *zap.Logger,
) DraftService {
	return &draftService{
		draftRepo: draftRepo,
		logger:    logger,
	}
}

func (s *draftService) Create(ctx context.Context, userID uuid.UUID, input CreateDraftInput) (*model.Draft, error) {
	draft := &model.Draft{
		UserID:        userID,
		Title:         input.Title,
		Content:       input.Content,
		CoverImageURL: input.CoverImageURL,
		ContentType:   input.ContentType,
		CategoryID:    input.CategoryID,
		Tags:          input.Tags,
		IsAutoSave:    false,
	}

	if err := s.draftRepo.Create(ctx, draft); err != nil {
		return nil, err
	}

	return s.draftRepo.GetByID(ctx, draft.ID)
}

func (s *draftService) GetByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.Draft, error) {
	draft, err := s.draftRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if draft.UserID != userID {
		return nil, ErrForbidden
	}

	return draft, nil
}

func (s *draftService) Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, input UpdateDraftInput) (*model.Draft, error) {
	draft, err := s.draftRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if draft.UserID != userID {
		return nil, ErrForbidden
	}

	// Update fields
	if input.Title != nil {
		draft.Title = *input.Title
	}
	if input.Content != nil {
		draft.Content = *input.Content
	}
	if input.CoverImageURL != nil {
		draft.CoverImageURL = input.CoverImageURL
	}
	if input.ContentType != nil {
		draft.ContentType = *input.ContentType
	}
	if input.CategoryID != nil {
		draft.CategoryID = input.CategoryID
	}
	if input.Tags != nil {
		draft.Tags = input.Tags
	}

	if err := s.draftRepo.Update(ctx, draft); err != nil {
		return nil, err
	}

	return s.draftRepo.GetByID(ctx, id)
}

func (s *draftService) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	draft, err := s.draftRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check ownership
	if draft.UserID != userID {
		return ErrForbidden
	}

	return s.draftRepo.Delete(ctx, id)
}

func (s *draftService) GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) (*DraftListResult, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	page := (offset / limit) + 1

	drafts, total, err := s.draftRepo.GetByUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return &DraftListResult{
		Items:    drafts,
		Total:    total,
		Page:     page,
		PageSize: limit,
		HasMore:  offset+len(drafts) < total,
	}, nil
}

func (s *draftService) AutoSave(ctx context.Context, userID uuid.UUID, input AutoSaveInput) error {
	draft := &model.Draft{
		UserID:        userID,
		ArticleID:     input.ArticleID,
		Title:         input.Title,
		Content:       input.Content,
		CoverImageURL: input.CoverImageURL,
		ContentType:   input.ContentType,
		CategoryID:    input.CategoryID,
		Tags:          input.Tags,
		IsAutoSave:    true,
	}

	return s.draftRepo.AutoSave(ctx, draft)
}

func (s *draftService) GetLatestAutoSave(ctx context.Context, userID uuid.UUID, articleID *uuid.UUID) (*model.Draft, error) {
	return s.draftRepo.GetLatestAutoSave(ctx, userID, articleID)
}

