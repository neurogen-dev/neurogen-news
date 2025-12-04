package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/neurogen-news/backend/internal/model"
)

var (
	ErrDraftNotFound = errors.New("draft not found")
)

type DraftRepository interface {
	Create(ctx context.Context, draft *model.Draft) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Draft, error)
	Update(ctx context.Context, draft *model.Draft) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Draft, int, error)
	
	// Auto-save
	AutoSave(ctx context.Context, draft *model.Draft) error
	GetLatestAutoSave(ctx context.Context, userID uuid.UUID, articleID *uuid.UUID) (*model.Draft, error)
}

type draftRepository struct {
	db *PostgresDB
}

func NewDraftRepository(db *PostgresDB) DraftRepository {
	return &draftRepository{db: db}
}

func (r *draftRepository) Create(ctx context.Context, draft *model.Draft) error {
	query := `
		INSERT INTO drafts (id, user_id, article_id, title, content, cover_image_url, 
		                    content_type, category_id, tags, is_auto_save, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
	`
	
	draft.ID = uuid.New()
	
	_, err := r.db.Exec(ctx, query,
		draft.ID,
		draft.UserID,
		draft.ArticleID,
		draft.Title,
		draft.Content,
		draft.CoverImageURL,
		draft.ContentType,
		draft.CategoryID,
		draft.Tags,
		draft.IsAutoSave,
	)
	
	return err
}

func (r *draftRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Draft, error) {
	query := `
		SELECT id, user_id, article_id, title, content, cover_image_url, 
		       content_type, category_id, tags, is_auto_save, created_at, updated_at
		FROM drafts
		WHERE id = $1
	`
	
	var draft model.Draft
	err := r.db.QueryRow(ctx, query, id).Scan(
		&draft.ID,
		&draft.UserID,
		&draft.ArticleID,
		&draft.Title,
		&draft.Content,
		&draft.CoverImageURL,
		&draft.ContentType,
		&draft.CategoryID,
		&draft.Tags,
		&draft.IsAutoSave,
		&draft.CreatedAt,
		&draft.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDraftNotFound
		}
		return nil, err
	}
	
	return &draft, nil
}

func (r *draftRepository) Update(ctx context.Context, draft *model.Draft) error {
	query := `
		UPDATE drafts SET
			title = $2, content = $3, cover_image_url = $4, content_type = $5,
			category_id = $6, tags = $7, updated_at = NOW()
		WHERE id = $1
	`
	
	_, err := r.db.Exec(ctx, query,
		draft.ID,
		draft.Title,
		draft.Content,
		draft.CoverImageURL,
		draft.ContentType,
		draft.CategoryID,
		draft.Tags,
	)
	
	return err
}

func (r *draftRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM drafts WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *draftRepository) GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Draft, int, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	
	// Count total (excluding auto-saves)
	countQuery := `SELECT COUNT(*) FROM drafts WHERE user_id = $1 AND is_auto_save = false`
	var total int
	if err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}
	
	// Get drafts with category info
	query := `
		SELECT d.id, d.user_id, d.article_id, d.title, d.content, d.cover_image_url,
		       d.content_type, d.category_id, d.tags, d.is_auto_save, d.created_at, d.updated_at,
		       c.name, c.slug
		FROM drafts d
		LEFT JOIN categories c ON c.id = d.category_id
		WHERE d.user_id = $1 AND d.is_auto_save = false
		ORDER BY d.updated_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var drafts []model.Draft
	for rows.Next() {
		var draft model.Draft
		var categoryName, categorySlug *string
		
		err := rows.Scan(
			&draft.ID,
			&draft.UserID,
			&draft.ArticleID,
			&draft.Title,
			&draft.Content,
			&draft.CoverImageURL,
			&draft.ContentType,
			&draft.CategoryID,
			&draft.Tags,
			&draft.IsAutoSave,
			&draft.CreatedAt,
			&draft.UpdatedAt,
			&categoryName,
			&categorySlug,
		)
		if err != nil {
			return nil, 0, err
		}
		
		if categoryName != nil {
			draft.CategoryName = categoryName
			draft.CategorySlug = categorySlug
		}
		
		drafts = append(drafts, draft)
	}
	
	return drafts, total, nil
}

func (r *draftRepository) AutoSave(ctx context.Context, draft *model.Draft) error {
	// Upsert auto-save draft
	query := `
		INSERT INTO drafts (id, user_id, article_id, title, content, cover_image_url,
		                    content_type, category_id, tags, is_auto_save, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, true, NOW(), NOW())
		ON CONFLICT (user_id, article_id) WHERE is_auto_save = true
		DO UPDATE SET
			title = EXCLUDED.title,
			content = EXCLUDED.content,
			cover_image_url = EXCLUDED.cover_image_url,
			content_type = EXCLUDED.content_type,
			category_id = EXCLUDED.category_id,
			tags = EXCLUDED.tags,
			updated_at = NOW()
	`
	
	if draft.ID == uuid.Nil {
		draft.ID = uuid.New()
	}
	
	_, err := r.db.Exec(ctx, query,
		draft.ID,
		draft.UserID,
		draft.ArticleID,
		draft.Title,
		draft.Content,
		draft.CoverImageURL,
		draft.ContentType,
		draft.CategoryID,
		draft.Tags,
	)
	
	return err
}

func (r *draftRepository) GetLatestAutoSave(ctx context.Context, userID uuid.UUID, articleID *uuid.UUID) (*model.Draft, error) {
	query := `
		SELECT id, user_id, article_id, title, content, cover_image_url,
		       content_type, category_id, tags, is_auto_save, created_at, updated_at
		FROM drafts
		WHERE user_id = $1 AND is_auto_save = true
	`
	
	args := []interface{}{userID}
	if articleID != nil {
		query += ` AND article_id = $2`
		args = append(args, *articleID)
	} else {
		query += ` AND article_id IS NULL`
	}
	
	query += ` ORDER BY updated_at DESC LIMIT 1`
	
	var draft model.Draft
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&draft.ID,
		&draft.UserID,
		&draft.ArticleID,
		&draft.Title,
		&draft.Content,
		&draft.CoverImageURL,
		&draft.ContentType,
		&draft.CategoryID,
		&draft.Tags,
		&draft.IsAutoSave,
		&draft.CreatedAt,
		&draft.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDraftNotFound
		}
		return nil, err
	}
	
	return &draft, nil
}

