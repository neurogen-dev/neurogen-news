package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/neurogen-news/backend/internal/model"
)

var (
	ErrArticleNotFound = errors.New("article not found")
)

type ArticleRepository interface {
	Create(ctx context.Context, article *model.Article) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Article, error)
	GetBySlug(ctx context.Context, categorySlug, articleSlug string) (*model.Article, error)
	Update(ctx context.Context, article *model.Article) error
	Delete(ctx context.Context, id uuid.UUID) error
	
	List(ctx context.Context, params ArticleListParams) ([]model.ArticleCard, int, error)
	GetByAuthor(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]model.ArticleCard, int, error)
	GetByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]model.ArticleCard, int, error)
	GetByTag(ctx context.Context, tagID uuid.UUID, limit, offset int) ([]model.ArticleCard, int, error)
	
	IncrementViewCount(ctx context.Context, id uuid.UUID) error
	
	// Tags
	AddTags(ctx context.Context, articleID uuid.UUID, tagIDs []uuid.UUID) error
	RemoveTags(ctx context.Context, articleID uuid.UUID) error
	GetTags(ctx context.Context, articleID uuid.UUID) ([]model.Tag, error)
}

type ArticleListParams struct {
	Sort        string // popular, new, hot
	Level       string
	ContentType string
	CategoryID  *uuid.UUID
	TagID       *uuid.UUID
	AuthorID    *uuid.UUID
	TimeRange   string // 24h, 7d, 30d, all
	Limit       int
	Offset      int
}

type articleRepository struct {
	db *PostgresDB
}

func NewArticleRepository(db *PostgresDB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) Create(ctx context.Context, article *model.Article) error {
	query := `
		INSERT INTO articles (
			id, title, slug, lead, content, html_content, cover_image_url,
			level, content_type, status, reading_time, is_editorial, is_pinned,
			is_nsfw, comments_enabled, author_id, category_id,
			meta_title, meta_description, canonical_url,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, NOW(), NOW()
		)
	`
	
	article.ID = uuid.New()
	
	_, err := r.db.Exec(ctx, query,
		article.ID,
		article.Title,
		article.Slug,
		article.Lead,
		article.Content,
		article.HTMLContent,
		article.CoverImageURL,
		article.Level,
		article.ContentType,
		article.Status,
		article.ReadingTime,
		article.IsEditorial,
		article.IsPinned,
		article.IsNSFW,
		article.CommentsEnabled,
		article.AuthorID,
		article.CategoryID,
		article.MetaTitle,
		article.MetaDescription,
		article.CanonicalURL,
	)
	
	return err
}

func (r *articleRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Article, error) {
	query := `
		SELECT 
			a.id, a.title, a.slug, a.lead, a.content, a.html_content, a.cover_image_url,
			a.level, a.content_type, a.status, a.reading_time, a.is_editorial, a.is_pinned,
			a.is_nsfw, a.comments_enabled, a.author_id, a.category_id,
			a.meta_title, a.meta_description, a.canonical_url,
			a.view_count, a.comment_count, a.bookmark_count,
			a.published_at, a.created_at, a.updated_at
		FROM articles a
		WHERE a.id = $1
	`
	
	var article model.Article
	err := r.db.QueryRow(ctx, query, id).Scan(
		&article.ID,
		&article.Title,
		&article.Slug,
		&article.Lead,
		&article.Content,
		&article.HTMLContent,
		&article.CoverImageURL,
		&article.Level,
		&article.ContentType,
		&article.Status,
		&article.ReadingTime,
		&article.IsEditorial,
		&article.IsPinned,
		&article.IsNSFW,
		&article.CommentsEnabled,
		&article.AuthorID,
		&article.CategoryID,
		&article.MetaTitle,
		&article.MetaDescription,
		&article.CanonicalURL,
		&article.ViewCount,
		&article.CommentCount,
		&article.BookmarkCount,
		&article.PublishedAt,
		&article.CreatedAt,
		&article.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}
	
	return &article, nil
}

func (r *articleRepository) GetBySlug(ctx context.Context, categorySlug, articleSlug string) (*model.Article, error) {
	query := `
		SELECT 
			a.id, a.title, a.slug, a.lead, a.content, a.html_content, a.cover_image_url,
			a.level, a.content_type, a.status, a.reading_time, a.is_editorial, a.is_pinned,
			a.is_nsfw, a.comments_enabled, a.author_id, a.category_id,
			a.meta_title, a.meta_description, a.canonical_url,
			a.view_count, a.comment_count, a.bookmark_count,
			a.published_at, a.created_at, a.updated_at
		FROM articles a
		JOIN categories c ON c.id = a.category_id
		WHERE a.slug = $1 AND c.slug = $2 AND a.status = 'published'
	`
	
	var article model.Article
	err := r.db.QueryRow(ctx, query, articleSlug, categorySlug).Scan(
		&article.ID,
		&article.Title,
		&article.Slug,
		&article.Lead,
		&article.Content,
		&article.HTMLContent,
		&article.CoverImageURL,
		&article.Level,
		&article.ContentType,
		&article.Status,
		&article.ReadingTime,
		&article.IsEditorial,
		&article.IsPinned,
		&article.IsNSFW,
		&article.CommentsEnabled,
		&article.AuthorID,
		&article.CategoryID,
		&article.MetaTitle,
		&article.MetaDescription,
		&article.CanonicalURL,
		&article.ViewCount,
		&article.CommentCount,
		&article.BookmarkCount,
		&article.PublishedAt,
		&article.CreatedAt,
		&article.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}
	
	return &article, nil
}

func (r *articleRepository) Update(ctx context.Context, article *model.Article) error {
	query := `
		UPDATE articles SET
			title = $2, slug = $3, lead = $4, content = $5, html_content = $6,
			cover_image_url = $7, level = $8, content_type = $9, status = $10,
			reading_time = $11, is_editorial = $12, is_pinned = $13, is_nsfw = $14,
			comments_enabled = $15, category_id = $16,
			meta_title = $17, meta_description = $18, canonical_url = $19,
			published_at = $20, updated_at = NOW()
		WHERE id = $1
	`
	
	_, err := r.db.Exec(ctx, query,
		article.ID,
		article.Title,
		article.Slug,
		article.Lead,
		article.Content,
		article.HTMLContent,
		article.CoverImageURL,
		article.Level,
		article.ContentType,
		article.Status,
		article.ReadingTime,
		article.IsEditorial,
		article.IsPinned,
		article.IsNSFW,
		article.CommentsEnabled,
		article.CategoryID,
		article.MetaTitle,
		article.MetaDescription,
		article.CanonicalURL,
		article.PublishedAt,
	)
	
	return err
}

func (r *articleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM articles WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *articleRepository) List(ctx context.Context, params ArticleListParams) ([]model.ArticleCard, int, error) {
	var conditions []string
	var args []interface{}
	argNum := 1
	
	// Base condition - only published
	conditions = append(conditions, "a.status = 'published'")
	
	// Filters
	if params.Level != "" {
		conditions = append(conditions, fmt.Sprintf("a.level = $%d", argNum))
		args = append(args, params.Level)
		argNum++
	}
	
	if params.ContentType != "" {
		conditions = append(conditions, fmt.Sprintf("a.content_type = $%d", argNum))
		args = append(args, params.ContentType)
		argNum++
	}
	
	if params.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("a.category_id = $%d", argNum))
		args = append(args, *params.CategoryID)
		argNum++
	}
	
	if params.AuthorID != nil {
		conditions = append(conditions, fmt.Sprintf("a.author_id = $%d", argNum))
		args = append(args, *params.AuthorID)
		argNum++
	}
	
	if params.TimeRange != "" && params.TimeRange != "all" {
		var interval string
		switch params.TimeRange {
		case "24h":
			interval = "1 day"
		case "7d":
			interval = "7 days"
		case "30d":
			interval = "30 days"
		}
		conditions = append(conditions, fmt.Sprintf("a.published_at > NOW() - INTERVAL '%s'", interval))
	}
	
	whereClause := strings.Join(conditions, " AND ")
	
	// Sort
	var orderBy string
	switch params.Sort {
	case "new":
		orderBy = "a.published_at DESC"
	case "hot":
		orderBy = "(a.view_count + a.comment_count * 10) DESC, a.published_at DESC"
	default: // popular
		orderBy = "a.view_count DESC, a.published_at DESC"
	}
	
	// Count query
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM articles a WHERE %s`, whereClause)
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	
	// Main query
	query := fmt.Sprintf(`
		SELECT 
			a.id, a.title, a.slug, a.lead, a.cover_image_url,
			a.level, a.content_type, a.reading_time, a.is_editorial, a.is_pinned,
			a.author_id, u.username, u.display_name, u.avatar_url, u.is_verified,
			a.category_id, c.slug, c.name, c.icon,
			a.view_count, a.comment_count, a.bookmark_count,
			a.published_at
		FROM articles a
		JOIN users u ON u.id = a.author_id
		JOIN categories c ON c.id = a.category_id
		WHERE %s
		ORDER BY a.is_pinned DESC, %s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderBy, argNum, argNum+1)
	
	args = append(args, params.Limit, params.Offset)
	
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var articles []model.ArticleCard
	for rows.Next() {
		var article model.ArticleCard
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Slug,
			&article.Lead,
			&article.CoverImageURL,
			&article.Level,
			&article.ContentType,
			&article.ReadingTime,
			&article.IsEditorial,
			&article.IsPinned,
			&article.AuthorID,
			&article.AuthorUsername,
			&article.AuthorName,
			&article.AuthorAvatar,
			&article.AuthorVerified,
			&article.CategoryID,
			&article.CategorySlug,
			&article.CategoryName,
			&article.CategoryIcon,
			&article.ViewCount,
			&article.CommentCount,
			&article.BookmarkCount,
			&article.PublishedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		articles = append(articles, article)
	}
	
	return articles, total, nil
}

func (r *articleRepository) GetByAuthor(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]model.ArticleCard, int, error) {
	return r.List(ctx, ArticleListParams{
		AuthorID: &authorID,
		Sort:     "new",
		Limit:    limit,
		Offset:   offset,
	})
}

func (r *articleRepository) GetByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]model.ArticleCard, int, error) {
	return r.List(ctx, ArticleListParams{
		CategoryID: &categoryID,
		Sort:       "popular",
		Limit:      limit,
		Offset:     offset,
	})
}

func (r *articleRepository) GetByTag(ctx context.Context, tagID uuid.UUID, limit, offset int) ([]model.ArticleCard, int, error) {
	// This would need a JOIN with article_tags
	// For now, use the List method with a tag filter
	return r.List(ctx, ArticleListParams{
		TagID:  &tagID,
		Sort:   "popular",
		Limit:  limit,
		Offset: offset,
	})
}

func (r *articleRepository) IncrementViewCount(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE articles SET view_count = view_count + 1 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *articleRepository) AddTags(ctx context.Context, articleID uuid.UUID, tagIDs []uuid.UUID) error {
	if len(tagIDs) == 0 {
		return nil
	}
	
	query := `INSERT INTO article_tags (article_id, tag_id) VALUES `
	values := make([]string, len(tagIDs))
	args := make([]interface{}, len(tagIDs)*2)
	
	for i, tagID := range tagIDs {
		values[i] = fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2)
		args[i*2] = articleID
		args[i*2+1] = tagID
	}
	
	query += strings.Join(values, ", ") + " ON CONFLICT DO NOTHING"
	
	_, err := r.db.Exec(ctx, query, args...)
	return err
}

func (r *articleRepository) RemoveTags(ctx context.Context, articleID uuid.UUID) error {
	query := `DELETE FROM article_tags WHERE article_id = $1`
	_, err := r.db.Exec(ctx, query, articleID)
	return err
}

func (r *articleRepository) GetTags(ctx context.Context, articleID uuid.UUID) ([]model.Tag, error) {
	query := `
		SELECT t.id, t.name, t.slug
		FROM tags t
		JOIN article_tags at ON at.tag_id = t.id
		WHERE at.article_id = $1
	`
	
	rows, err := r.db.Query(ctx, query, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var tags []model.Tag
	for rows.Next() {
		var tag model.Tag
		if err := rows.Scan(&tag.ID, &tag.Name, &tag.Slug); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	
	return tags, nil
}

