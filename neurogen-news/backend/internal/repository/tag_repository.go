package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/neurogen-news/backend/internal/model"
)

var (
	ErrTagNotFound = errors.New("tag not found")
)

type TagRepository interface {
	GetAll(ctx context.Context, limit int) ([]model.Tag, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Tag, error)
	GetBySlug(ctx context.Context, slug string) (*model.Tag, error)
	GetByName(ctx context.Context, name string) (*model.Tag, error)
	Create(ctx context.Context, tag *model.Tag) error
	Update(ctx context.Context, tag *model.Tag) error
	Delete(ctx context.Context, id uuid.UUID) error
	
	GetPopular(ctx context.Context, limit int) ([]model.Tag, error)
	Search(ctx context.Context, query string, limit int) ([]model.Tag, error)
}

type tagRepository struct {
	db *PostgresDB
}

func NewTagRepository(db *PostgresDB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) GetAll(ctx context.Context, limit int) ([]model.Tag, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	
	query := `
		SELECT id, name, slug, article_count, created_at
		FROM tags
		ORDER BY article_count DESC
		LIMIT $1
	`
	
	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var tags []model.Tag
	for rows.Next() {
		var tag model.Tag
		if err := rows.Scan(&tag.ID, &tag.Name, &tag.Slug, &tag.ArticleCount, &tag.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	
	return tags, nil
}

func (r *tagRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Tag, error) {
	query := `SELECT id, name, slug, article_count, created_at FROM tags WHERE id = $1`
	
	var tag model.Tag
	err := r.db.QueryRow(ctx, query, id).Scan(&tag.ID, &tag.Name, &tag.Slug, &tag.ArticleCount, &tag.CreatedAt)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTagNotFound
		}
		return nil, err
	}
	
	return &tag, nil
}

func (r *tagRepository) GetBySlug(ctx context.Context, slug string) (*model.Tag, error) {
	query := `SELECT id, name, slug, article_count, created_at FROM tags WHERE slug = $1`
	
	var tag model.Tag
	err := r.db.QueryRow(ctx, query, slug).Scan(&tag.ID, &tag.Name, &tag.Slug, &tag.ArticleCount, &tag.CreatedAt)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTagNotFound
		}
		return nil, err
	}
	
	return &tag, nil
}

func (r *tagRepository) GetByName(ctx context.Context, name string) (*model.Tag, error) {
	query := `SELECT id, name, slug, article_count, created_at FROM tags WHERE LOWER(name) = LOWER($1)`
	
	var tag model.Tag
	err := r.db.QueryRow(ctx, query, name).Scan(&tag.ID, &tag.Name, &tag.Slug, &tag.ArticleCount, &tag.CreatedAt)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTagNotFound
		}
		return nil, err
	}
	
	return &tag, nil
}

func (r *tagRepository) Create(ctx context.Context, tag *model.Tag) error {
	query := `
		INSERT INTO tags (id, name, slug, created_at)
		VALUES ($1, $2, $3, NOW())
	`
	
	if tag.ID == uuid.Nil {
		tag.ID = uuid.New()
	}
	
	_, err := r.db.Exec(ctx, query, tag.ID, tag.Name, tag.Slug)
	return err
}

func (r *tagRepository) Update(ctx context.Context, tag *model.Tag) error {
	query := `UPDATE tags SET name = $2, slug = $3 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, tag.ID, tag.Name, tag.Slug)
	return err
}

func (r *tagRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tags WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *tagRepository) GetPopular(ctx context.Context, limit int) ([]model.Tag, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	
	query := `
		SELECT id, name, slug, article_count, created_at
		FROM tags
		WHERE article_count > 0
		ORDER BY article_count DESC
		LIMIT $1
	`
	
	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var tags []model.Tag
	for rows.Next() {
		var tag model.Tag
		if err := rows.Scan(&tag.ID, &tag.Name, &tag.Slug, &tag.ArticleCount, &tag.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	
	return tags, nil
}

func (r *tagRepository) Search(ctx context.Context, query string, limit int) ([]model.Tag, error) {
	if limit <= 0 || limit > 20 {
		limit = 10
	}
	
	sqlQuery := `
		SELECT id, name, slug, article_count, created_at
		FROM tags
		WHERE name ILIKE $1
		ORDER BY article_count DESC
		LIMIT $2
	`
	
	rows, err := r.db.Query(ctx, sqlQuery, query+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var tags []model.Tag
	for rows.Next() {
		var tag model.Tag
		if err := rows.Scan(&tag.ID, &tag.Name, &tag.Slug, &tag.ArticleCount, &tag.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	
	return tags, nil
}

