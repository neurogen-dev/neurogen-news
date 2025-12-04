package service

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/model"
	"github.com/neurogen-news/backend/internal/repository"
)

type CommentService interface {
	Create(ctx context.Context, userID uuid.UUID, input CreateCommentInput) (*model.Comment, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Comment, error)
	Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, content string) (*model.Comment, error)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error

	// Lists
	GetByArticle(ctx context.Context, articleID uuid.UUID, params CommentListParams) (*CommentListResult, error)
	GetReplies(ctx context.Context, parentID uuid.UUID, limit, offset int) ([]model.Comment, error)
	GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) (*CommentListResult, error)

	// Reactions
	AddReaction(ctx context.Context, userID, commentID uuid.UUID, emoji string) error
	RemoveReaction(ctx context.Context, userID, commentID uuid.UUID) error
}

type CreateCommentInput struct {
	ArticleID uuid.UUID  `json:"articleId" validate:"required"`
	ParentID  *uuid.UUID `json:"parentId,omitempty"`
	Content   string     `json:"content" validate:"required,min=1,max=10000"`
}

type CommentListParams struct {
	Sort     string `query:"sort" validate:"omitempty,oneof=new popular old"`
	Page     int    `query:"page" validate:"min=1"`
	PageSize int    `query:"pageSize" validate:"min=1,max=100"`
}

type CommentListResult struct {
	Items    []model.Comment `json:"items"`
	Total    int             `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"pageSize"`
	HasMore  bool            `json:"hasMore"`
}

type commentService struct {
	commentRepo      repository.CommentRepository
	notificationRepo repository.NotificationRepository
	reactionRepo     repository.ReactionRepository
	redis            *repository.RedisClient
	logger           *zap.Logger
}

func NewCommentService(
	commentRepo repository.CommentRepository,
	notificationRepo repository.NotificationRepository,
	reactionRepo repository.ReactionRepository,
	redis *repository.RedisClient,
	logger *zap.Logger,
) CommentService {
	return &commentService{
		commentRepo:      commentRepo,
		notificationRepo: notificationRepo,
		reactionRepo:     reactionRepo,
		redis:            redis,
		logger:           logger,
	}
}

func (s *commentService) Create(ctx context.Context, userID uuid.UUID, input CreateCommentInput) (*model.Comment, error) {
	// Convert content to HTML (basic)
	htmlContent := convertCommentToHTML(input.Content)

	comment := &model.Comment{
		Content:     input.Content,
		HTMLContent: htmlContent,
		AuthorID:    userID,
		ArticleID:   input.ArticleID,
		ParentID:    input.ParentID,
	}

	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, err
	}

	// Get full comment with author
	fullComment, err := s.commentRepo.GetByID(ctx, comment.ID)
	if err != nil {
		return nil, err
	}

	// TODO: Send notification to article author or parent comment author

	return fullComment, nil
}

func (s *commentService) GetByID(ctx context.Context, id uuid.UUID) (*model.Comment, error) {
	comment, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get reactions
	reactions, _ := s.commentRepo.GetReactions(ctx, id, nil)
	comment.Reactions = reactions

	return comment, nil
}

func (s *commentService) Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, content string) (*model.Comment, error) {
	comment, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if comment.AuthorID != userID {
		return nil, ErrForbidden
	}

	comment.Content = content
	comment.HTMLContent = convertCommentToHTML(content)
	comment.IsEdited = true

	if err := s.commentRepo.Update(ctx, comment); err != nil {
		return nil, err
	}

	return s.commentRepo.GetByID(ctx, id)
}

func (s *commentService) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	comment, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check ownership (or moderation rights in the future)
	if comment.AuthorID != userID {
		return ErrForbidden
	}

	return s.commentRepo.Delete(ctx, id)
}

func (s *commentService) GetByArticle(ctx context.Context, articleID uuid.UUID, params CommentListParams) (*CommentListResult, error) {
	// Defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 20
	}
	if params.Sort == "" {
		params.Sort = "new"
	}

	offset := (params.Page - 1) * params.PageSize

	comments, total, err := s.commentRepo.GetByArticle(ctx, articleID, repository.CommentListParams{
		Sort:   params.Sort,
		Limit:  params.PageSize,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	// Load first 3 replies for each comment and reactions
	for i := range comments {
		if comments[i].ReplyCount > 0 {
			replies, _ := s.commentRepo.GetReplies(ctx, comments[i].ID, 3, 0)
			comments[i].Replies = replies
		}
		reactions, _ := s.commentRepo.GetReactions(ctx, comments[i].ID, nil)
		comments[i].Reactions = reactions
	}

	return &CommentListResult{
		Items:    comments,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
		HasMore:  offset+len(comments) < total,
	}, nil
}

func (s *commentService) GetReplies(ctx context.Context, parentID uuid.UUID, limit, offset int) ([]model.Comment, error) {
	replies, err := s.commentRepo.GetReplies(ctx, parentID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Load reactions for each reply
	for i := range replies {
		reactions, _ := s.commentRepo.GetReactions(ctx, replies[i].ID, nil)
		replies[i].Reactions = reactions
	}

	return replies, nil
}

func (s *commentService) GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) (*CommentListResult, error) {
	if limit < 1 || limit > 50 {
		limit = 20
	}

	page := (offset / limit) + 1

	comments, total, err := s.commentRepo.GetByUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return &CommentListResult{
		Items:    comments,
		Total:    total,
		Page:     page,
		PageSize: limit,
		HasMore:  offset+len(comments) < total,
	}, nil
}

func (s *commentService) AddReaction(ctx context.Context, userID, commentID uuid.UUID, emoji string) error {
	// Get reaction type by emoji
	reactionType, err := s.reactionRepo.GetByEmoji(ctx, emoji)
	if err != nil {
		return err
	}

	return s.commentRepo.AddReaction(ctx, commentID, userID, reactionType.ID)
}

func (s *commentService) RemoveReaction(ctx context.Context, userID, commentID uuid.UUID) error {
	return s.commentRepo.RemoveReaction(ctx, commentID, userID)
}

// convertCommentToHTML converts markdown to simple HTML
func convertCommentToHTML(content string) string {
	// Basic conversion - in production use a proper markdown parser
	html := strings.TrimSpace(content)

	// Escape HTML entities
	html = strings.ReplaceAll(html, "&", "&amp;")
	html = strings.ReplaceAll(html, "<", "&lt;")
	html = strings.ReplaceAll(html, ">", "&gt;")

	// Convert line breaks
	html = strings.ReplaceAll(html, "\n\n", "</p><p>")
	html = strings.ReplaceAll(html, "\n", "<br>")

	// Wrap in paragraph
	html = "<p>" + html + "</p>"

	return html
}

