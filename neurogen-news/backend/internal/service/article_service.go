package service

import (
	"context"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/model"
	"github.com/neurogen-news/backend/internal/repository"
)

type ArticleService interface {
	Create(ctx context.Context, userID uuid.UUID, input CreateArticleInput) (*model.Article, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Article, error)
	GetBySlug(ctx context.Context, categorySlug, articleSlug string) (*model.Article, error)
	Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, input UpdateArticleInput) (*model.Article, error)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	List(ctx context.Context, params ArticleListParams) (*ArticleListResult, error)
	
	// Reactions
	AddReaction(ctx context.Context, userID, articleID uuid.UUID, emoji string) error
	RemoveReaction(ctx context.Context, userID, articleID uuid.UUID) error
	
	// View count
	RecordView(ctx context.Context, articleID uuid.UUID, userIP string) error
}

type CreateArticleInput struct {
	Title           string             `json:"title" validate:"required,min=5,max=200"`
	Content         string             `json:"content" validate:"required,min=100"`
	Lead            *string            `json:"lead,omitempty" validate:"omitempty,max=500"`
	CoverImageURL   *string            `json:"coverImageUrl,omitempty" validate:"omitempty,url"`
	Level           model.ArticleLevel `json:"level" validate:"required,oneof=beginner intermediate advanced"`
	ContentType     model.ContentType  `json:"contentType" validate:"required,oneof=article news post question discussion"`
	CategoryID      uuid.UUID          `json:"categoryId" validate:"required"`
	Tags            []string           `json:"tags,omitempty" validate:"max=10"`
	IsNSFW          bool               `json:"isNsfw"`
	CommentsEnabled bool               `json:"commentsEnabled"`
	Status          model.ArticleStatus `json:"status" validate:"oneof=draft published"`
	
	// SEO
	MetaTitle       *string `json:"metaTitle,omitempty" validate:"omitempty,max=60"`
	MetaDescription *string `json:"metaDescription,omitempty" validate:"omitempty,max=160"`
}

type UpdateArticleInput struct {
	Title           *string             `json:"title,omitempty" validate:"omitempty,min=5,max=200"`
	Content         *string             `json:"content,omitempty" validate:"omitempty,min=100"`
	Lead            *string             `json:"lead,omitempty" validate:"omitempty,max=500"`
	CoverImageURL   *string             `json:"coverImageUrl,omitempty" validate:"omitempty,url"`
	Level           *model.ArticleLevel `json:"level,omitempty" validate:"omitempty,oneof=beginner intermediate advanced"`
	ContentType     *model.ContentType  `json:"contentType,omitempty" validate:"omitempty,oneof=article news post question discussion"`
	CategoryID      *uuid.UUID          `json:"categoryId,omitempty"`
	Tags            []string            `json:"tags,omitempty" validate:"max=10"`
	IsNSFW          *bool               `json:"isNsfw,omitempty"`
	CommentsEnabled *bool               `json:"commentsEnabled,omitempty"`
	Status          *model.ArticleStatus `json:"status,omitempty" validate:"omitempty,oneof=draft published archived"`
	MetaTitle       *string             `json:"metaTitle,omitempty" validate:"omitempty,max=60"`
	MetaDescription *string             `json:"metaDescription,omitempty" validate:"omitempty,max=160"`
}

type ArticleListParams struct {
	Sort        string `query:"sort" validate:"omitempty,oneof=popular new hot"`
	Level       string `query:"level" validate:"omitempty,oneof=beginner intermediate advanced"`
	ContentType string `query:"contentType" validate:"omitempty,oneof=article news post question discussion"`
	CategoryID  string `query:"categoryId" validate:"omitempty,uuid"`
	TagID       string `query:"tagId" validate:"omitempty,uuid"`
	TimeRange   string `query:"timeRange" validate:"omitempty,oneof=24h 7d 30d all"`
	Page        int    `query:"page" validate:"min=1"`
	PageSize    int    `query:"pageSize" validate:"min=1,max=50"`
}

type ArticleListResult struct {
	Items    []model.ArticleCard `json:"items"`
	Total    int                 `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"pageSize"`
	HasMore  bool                `json:"hasMore"`
}

type articleService struct {
	articleRepo repository.ArticleRepository
	tagRepo     repository.TagRepository
	redis       *repository.RedisClient
	logger      *zap.Logger
}

func NewArticleService(
	articleRepo repository.ArticleRepository,
	tagRepo repository.TagRepository,
	redis *repository.RedisClient,
	logger *zap.Logger,
) ArticleService {
	return &articleService{
		articleRepo: articleRepo,
		tagRepo:     tagRepo,
		redis:       redis,
		logger:      logger,
	}
}

func (s *articleService) Create(ctx context.Context, userID uuid.UUID, input CreateArticleInput) (*model.Article, error) {
	// Generate slug
	slug := generateSlug(input.Title)
	
	// Calculate reading time (rough estimate: 200 words per minute)
	wordCount := len(strings.Fields(input.Content))
	readingTime := (wordCount + 199) / 200
	if readingTime < 1 {
		readingTime = 1
	}
	
	// Convert markdown to HTML (simplified)
	htmlContent := convertToHTML(input.Content)
	
	article := &model.Article{
		Title:           input.Title,
		Slug:            slug,
		Lead:            input.Lead,
		Content:         input.Content,
		HTMLContent:     htmlContent,
		CoverImageURL:   input.CoverImageURL,
		Level:           input.Level,
		ContentType:     input.ContentType,
		Status:          input.Status,
		ReadingTime:     readingTime,
		IsNSFW:          input.IsNSFW,
		CommentsEnabled: input.CommentsEnabled,
		AuthorID:        userID,
		CategoryID:      input.CategoryID,
		MetaTitle:       input.MetaTitle,
		MetaDescription: input.MetaDescription,
	}
	
	if input.Status == model.StatusPublished {
		now := time.Now()
		article.PublishedAt = &now
	}
	
	if err := s.articleRepo.Create(ctx, article); err != nil {
		return nil, err
	}
	
	// Add tags
	if len(input.Tags) > 0 {
		tagIDs, err := s.getOrCreateTags(ctx, input.Tags)
		if err != nil {
			s.logger.Error("Failed to create tags", zap.Error(err))
		} else {
			_ = s.articleRepo.AddTags(ctx, article.ID, tagIDs)
		}
	}
	
	// Invalidate cache
	s.invalidateCache(ctx)
	
	return article, nil
}

func (s *articleService) GetByID(ctx context.Context, id uuid.UUID) (*model.Article, error) {
	article, err := s.articleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Get tags
	tags, _ := s.articleRepo.GetTags(ctx, id)
	article.Tags = tags
	
	return article, nil
}

func (s *articleService) GetBySlug(ctx context.Context, categorySlug, articleSlug string) (*model.Article, error) {
	article, err := s.articleRepo.GetBySlug(ctx, categorySlug, articleSlug)
	if err != nil {
		return nil, err
	}
	
	// Get tags
	tags, _ := s.articleRepo.GetTags(ctx, article.ID)
	article.Tags = tags
	
	return article, nil
}

func (s *articleService) Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, input UpdateArticleInput) (*model.Article, error) {
	article, err := s.articleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Check ownership (or admin)
	if article.AuthorID != userID {
		return nil, ErrForbidden
	}
	
	// Update fields
	if input.Title != nil {
		article.Title = *input.Title
		article.Slug = generateSlug(*input.Title)
	}
	if input.Content != nil {
		article.Content = *input.Content
		article.HTMLContent = convertToHTML(*input.Content)
		wordCount := len(strings.Fields(*input.Content))
		article.ReadingTime = (wordCount + 199) / 200
	}
	if input.Lead != nil {
		article.Lead = input.Lead
	}
	if input.CoverImageURL != nil {
		article.CoverImageURL = input.CoverImageURL
	}
	if input.Level != nil {
		article.Level = *input.Level
	}
	if input.ContentType != nil {
		article.ContentType = *input.ContentType
	}
	if input.CategoryID != nil {
		article.CategoryID = *input.CategoryID
	}
	if input.IsNSFW != nil {
		article.IsNSFW = *input.IsNSFW
	}
	if input.CommentsEnabled != nil {
		article.CommentsEnabled = *input.CommentsEnabled
	}
	if input.Status != nil {
		article.Status = *input.Status
		if *input.Status == model.StatusPublished && article.PublishedAt == nil {
			now := time.Now()
			article.PublishedAt = &now
		}
	}
	if input.MetaTitle != nil {
		article.MetaTitle = input.MetaTitle
	}
	if input.MetaDescription != nil {
		article.MetaDescription = input.MetaDescription
	}
	
	if err := s.articleRepo.Update(ctx, article); err != nil {
		return nil, err
	}
	
	// Update tags
	if input.Tags != nil {
		_ = s.articleRepo.RemoveTags(ctx, article.ID)
		if len(input.Tags) > 0 {
			tagIDs, err := s.getOrCreateTags(ctx, input.Tags)
			if err == nil {
				_ = s.articleRepo.AddTags(ctx, article.ID, tagIDs)
			}
		}
	}
	
	s.invalidateCache(ctx)
	
	return article, nil
}

func (s *articleService) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	article, err := s.articleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	if article.AuthorID != userID {
		return ErrForbidden
	}
	
	if err := s.articleRepo.Delete(ctx, id); err != nil {
		return err
	}
	
	s.invalidateCache(ctx)
	
	return nil
}

func (s *articleService) List(ctx context.Context, params ArticleListParams) (*ArticleListResult, error) {
	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 50 {
		params.PageSize = 20
	}
	if params.Sort == "" {
		params.Sort = "popular"
	}
	
	offset := (params.Page - 1) * params.PageSize
	
	// Convert params
	repoParams := repository.ArticleListParams{
		Sort:        params.Sort,
		Level:       params.Level,
		ContentType: params.ContentType,
		TimeRange:   params.TimeRange,
		Limit:       params.PageSize,
		Offset:      offset,
	}
	
	if params.CategoryID != "" {
		id, _ := uuid.Parse(params.CategoryID)
		repoParams.CategoryID = &id
	}
	if params.TagID != "" {
		id, _ := uuid.Parse(params.TagID)
		repoParams.TagID = &id
	}
	
	articles, total, err := s.articleRepo.List(ctx, repoParams)
	if err != nil {
		return nil, err
	}
	
	return &ArticleListResult{
		Items:    articles,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
		HasMore:  offset+len(articles) < total,
	}, nil
}

func (s *articleService) AddReaction(ctx context.Context, userID, articleID uuid.UUID, emoji string) error {
	// TODO: Implement reaction logic
	return nil
}

func (s *articleService) RemoveReaction(ctx context.Context, userID, articleID uuid.UUID) error {
	// TODO: Implement reaction removal
	return nil
}

func (s *articleService) RecordView(ctx context.Context, articleID uuid.UUID, userIP string) error {
	// Use Redis to deduplicate views (1 view per IP per 24 hours)
	key := "view:" + articleID.String() + ":" + userIP
	exists, err := s.redis.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	
	if exists == 0 {
		// New view
		if err := s.redis.Set(ctx, key, "1", 24*time.Hour).Err(); err != nil {
			return err
		}
		return s.articleRepo.IncrementViewCount(ctx, articleID)
	}
	
	return nil
}

func (s *articleService) getOrCreateTags(ctx context.Context, tagNames []string) ([]uuid.UUID, error) {
	var tagIDs []uuid.UUID
	
	for _, name := range tagNames {
		name = strings.TrimSpace(strings.ToLower(name))
		if name == "" {
			continue
		}
		
		tag, err := s.tagRepo.GetByName(ctx, name)
		if err != nil {
			// Create new tag
			tag = &model.Tag{
				ID:   uuid.New(),
				Name: name,
				Slug: generateSlug(name),
			}
			if err := s.tagRepo.Create(ctx, tag); err != nil {
				continue
			}
		}
		tagIDs = append(tagIDs, tag.ID)
	}
	
	return tagIDs, nil
}

func (s *articleService) invalidateCache(ctx context.Context) {
	// Delete cached article lists
	s.redis.Del(ctx, "articles:popular", "articles:new", "articles:hot")
}

// Helper functions

func generateSlug(title string) string {
	// Transliteration map for Cyrillic
	translit := map[rune]string{
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo",
		'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m",
		'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
		'ф': "f", 'х': "h", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "sch",
		'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
	}
	
	var result strings.Builder
	title = strings.ToLower(title)
	
	for _, r := range title {
		if trans, ok := translit[r]; ok {
			result.WriteString(trans)
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
		} else if r == ' ' || r == '-' || r == '_' {
			result.WriteRune('-')
		}
	}
	
	// Clean up multiple dashes
	slug := result.String()
	re := regexp.MustCompile(`-+`)
	slug = re.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	
	// Limit length
	if len(slug) > 100 {
		slug = slug[:100]
	}
	
	// Add random suffix to ensure uniqueness
	slug = slug + "-" + uuid.New().String()[:8]
	
	return slug
}

func convertToHTML(markdown string) string {
	// Simplified conversion - in production, use a proper markdown parser
	// This is just a placeholder
	html := markdown
	
	// Basic conversions
	html = strings.ReplaceAll(html, "\n\n", "</p><p>")
	html = "<p>" + html + "</p>"
	
	return html
}

// Errors
var ErrForbidden = &AppError{Code: "FORBIDDEN", Message: "You don't have permission to perform this action"}

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

