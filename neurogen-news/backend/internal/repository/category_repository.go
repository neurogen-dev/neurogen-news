package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/neurogen-news/backend/internal/model"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
)

type CategoryRepository interface {
	GetAll(ctx context.Context) ([]model.Category, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Category, error)
	GetBySlug(ctx context.Context, slug string) (*model.Category, error)
	Create(ctx context.Context, category *model.Category) error
	Update(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	
	// Subscriptions
	Subscribe(ctx context.Context, userID, categoryID uuid.UUID) error
	Unsubscribe(ctx context.Context, userID, categoryID uuid.UUID) error
	IsSubscribed(ctx context.Context, userID, categoryID uuid.UUID) (bool, error)
	GetUserSubscriptions(ctx context.Context, userID uuid.UUID) ([]model.Category, error)
}

type categoryRepository struct {
	db *PostgresDB
}

func NewCategoryRepository(db *PostgresDB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) GetAll(ctx context.Context) ([]model.Category, error) {
	query := `
		SELECT id, name, slug, description, icon, color, is_official, parent_id,
		       article_count, subscriber_count, created_at, updated_at
		FROM categories
		WHERE parent_id IS NULL
		ORDER BY article_count DESC, name ASC
	`
	
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var categories []model.Category
	for rows.Next() {
		var cat model.Category
		err := rows.Scan(
			&cat.ID, &cat.Name, &cat.Slug, &cat.Description, &cat.Icon, &cat.Color,
			&cat.IsOfficial, &cat.ParentID, &cat.ArticleCount, &cat.SubscriberCount,
			&cat.CreatedAt, &cat.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	
	return categories, nil
}

func (r *categoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Category, error) {
	query := `
		SELECT id, name, slug, description, icon, color, is_official, parent_id,
		       article_count, subscriber_count, created_at, updated_at
		FROM categories
		WHERE id = $1
	`
	
	var cat model.Category
	err := r.db.QueryRow(ctx, query, id).Scan(
		&cat.ID, &cat.Name, &cat.Slug, &cat.Description, &cat.Icon, &cat.Color,
		&cat.IsOfficial, &cat.ParentID, &cat.ArticleCount, &cat.SubscriberCount,
		&cat.CreatedAt, &cat.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	
	return &cat, nil
}

func (r *categoryRepository) GetBySlug(ctx context.Context, slug string) (*model.Category, error) {
	query := `
		SELECT id, name, slug, description, icon, color, is_official, parent_id,
		       article_count, subscriber_count, created_at, updated_at
		FROM categories
		WHERE slug = $1
	`
	
	var cat model.Category
	err := r.db.QueryRow(ctx, query, slug).Scan(
		&cat.ID, &cat.Name, &cat.Slug, &cat.Description, &cat.Icon, &cat.Color,
		&cat.IsOfficial, &cat.ParentID, &cat.ArticleCount, &cat.SubscriberCount,
		&cat.CreatedAt, &cat.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	
	return &cat, nil
}

func (r *categoryRepository) Create(ctx context.Context, category *model.Category) error {
	query := `
		INSERT INTO categories (id, name, slug, description, icon, color, is_official, parent_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
	`
	
	category.ID = uuid.New()
	
	_, err := r.db.Exec(ctx, query,
		category.ID, category.Name, category.Slug, category.Description,
		category.Icon, category.Color, category.IsOfficial, category.ParentID,
	)
	
	return err
}

func (r *categoryRepository) Update(ctx context.Context, category *model.Category) error {
	query := `
		UPDATE categories SET
			name = $2, slug = $3, description = $4, icon = $5, color = $6,
			is_official = $7, updated_at = NOW()
		WHERE id = $1
	`
	
	_, err := r.db.Exec(ctx, query,
		category.ID, category.Name, category.Slug, category.Description,
		category.Icon, category.Color, category.IsOfficial,
	)
	
	return err
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *categoryRepository) Subscribe(ctx context.Context, userID, categoryID uuid.UUID) error {
	query := `
		INSERT INTO category_subscriptions (id, user_id, category_id, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id, category_id) DO NOTHING
	`
	_, err := r.db.Exec(ctx, query, uuid.New(), userID, categoryID)
	return err
}

func (r *categoryRepository) Unsubscribe(ctx context.Context, userID, categoryID uuid.UUID) error {
	query := `DELETE FROM category_subscriptions WHERE user_id = $1 AND category_id = $2`
	_, err := r.db.Exec(ctx, query, userID, categoryID)
	return err
}

func (r *categoryRepository) IsSubscribed(ctx context.Context, userID, categoryID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM category_subscriptions WHERE user_id = $1 AND category_id = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, categoryID).Scan(&exists)
	return exists, err
}

func (r *categoryRepository) GetUserSubscriptions(ctx context.Context, userID uuid.UUID) ([]model.Category, error) {
	query := `
		SELECT c.id, c.name, c.slug, c.description, c.icon, c.color, c.is_official, c.parent_id,
		       c.article_count, c.subscriber_count, c.created_at, c.updated_at
		FROM categories c
		JOIN category_subscriptions cs ON cs.category_id = c.id
		WHERE cs.user_id = $1
		ORDER BY c.name ASC
	`
	
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var categories []model.Category
	for rows.Next() {
		var cat model.Category
		err := rows.Scan(
			&cat.ID, &cat.Name, &cat.Slug, &cat.Description, &cat.Icon, &cat.Color,
			&cat.IsOfficial, &cat.ParentID, &cat.ArticleCount, &cat.SubscriberCount,
			&cat.CreatedAt, &cat.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	
	return categories, nil
}

