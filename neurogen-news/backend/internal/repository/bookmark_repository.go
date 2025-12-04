package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/neurogen-news/backend/internal/model"
)

type BookmarkRepository interface {
	Add(ctx context.Context, userID, articleID uuid.UUID, folderID *uuid.UUID) error
	Remove(ctx context.Context, userID, articleID uuid.UUID) error
	IsBookmarked(ctx context.Context, userID, articleID uuid.UUID) (bool, error)
	GetByUser(ctx context.Context, userID uuid.UUID, params BookmarkListParams) ([]model.ArticleCard, int, error)
	
	// Folders
	GetFolders(ctx context.Context, userID uuid.UUID) ([]model.BookmarkFolder, error)
	CreateFolder(ctx context.Context, folder *model.BookmarkFolder) error
	UpdateFolder(ctx context.Context, folder *model.BookmarkFolder) error
	DeleteFolder(ctx context.Context, id uuid.UUID) error
	MoveToFolder(ctx context.Context, userID, articleID uuid.UUID, folderID *uuid.UUID) error
}

type BookmarkListParams struct {
	FolderID *uuid.UUID
	Limit    int
	Offset   int
}

type bookmarkRepository struct {
	db *PostgresDB
}

func NewBookmarkRepository(db *PostgresDB) BookmarkRepository {
	return &bookmarkRepository{db: db}
}

func (r *bookmarkRepository) Add(ctx context.Context, userID, articleID uuid.UUID, folderID *uuid.UUID) error {
	query := `
		INSERT INTO bookmarks (id, user_id, article_id, folder_id, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (user_id, article_id) DO UPDATE SET folder_id = $4
	`
	_, err := r.db.Exec(ctx, query, uuid.New(), userID, articleID, folderID)
	return err
}

func (r *bookmarkRepository) Remove(ctx context.Context, userID, articleID uuid.UUID) error {
	query := `DELETE FROM bookmarks WHERE user_id = $1 AND article_id = $2`
	_, err := r.db.Exec(ctx, query, userID, articleID)
	return err
}

func (r *bookmarkRepository) IsBookmarked(ctx context.Context, userID, articleID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM bookmarks WHERE user_id = $1 AND article_id = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, articleID).Scan(&exists)
	return exists, err
}

func (r *bookmarkRepository) GetByUser(ctx context.Context, userID uuid.UUID, params BookmarkListParams) ([]model.ArticleCard, int, error) {
	if params.Limit <= 0 || params.Limit > 50 {
		params.Limit = 20
	}
	
	// Build WHERE clause
	whereClause := "b.user_id = $1"
	args := []interface{}{userID}
	argNum := 2
	
	if params.FolderID != nil {
		whereClause += " AND b.folder_id = $2"
		args = append(args, *params.FolderID)
		argNum++
	}
	
	// Count total
	countQuery := "SELECT COUNT(*) FROM bookmarks b WHERE " + whereClause
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	
	// Get bookmarked articles
	query := `
		SELECT 
			a.id, a.title, a.slug, a.lead, a.cover_image_url,
			a.level, a.content_type, a.reading_time, a.is_editorial, a.is_pinned,
			a.author_id, u.username, u.display_name, u.avatar_url, u.is_verified,
			a.category_id, c.slug, c.name, c.icon,
			a.view_count, a.comment_count, a.bookmark_count,
			a.published_at
		FROM bookmarks b
		JOIN articles a ON a.id = b.article_id
		JOIN users u ON u.id = a.author_id
		JOIN categories c ON c.id = a.category_id
		WHERE ` + whereClause + `
		ORDER BY b.created_at DESC
		LIMIT $` + string(rune('0'+argNum)) + ` OFFSET $` + string(rune('0'+argNum+1))
	
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
			&article.ID, &article.Title, &article.Slug, &article.Lead, &article.CoverImageURL,
			&article.Level, &article.ContentType, &article.ReadingTime, &article.IsEditorial, &article.IsPinned,
			&article.AuthorID, &article.AuthorUsername, &article.AuthorName, &article.AuthorAvatar, &article.AuthorVerified,
			&article.CategoryID, &article.CategorySlug, &article.CategoryName, &article.CategoryIcon,
			&article.ViewCount, &article.CommentCount, &article.BookmarkCount,
			&article.PublishedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		articles = append(articles, article)
	}
	
	return articles, total, nil
}

func (r *bookmarkRepository) GetFolders(ctx context.Context, userID uuid.UUID) ([]model.BookmarkFolder, error) {
	query := `
		SELECT bf.id, bf.name, bf.user_id, bf.created_at,
		       (SELECT COUNT(*) FROM bookmarks b WHERE b.folder_id = bf.id) as bookmark_count
		FROM bookmark_folders bf
		WHERE bf.user_id = $1
		ORDER BY bf.name ASC
	`
	
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []model.BookmarkFolder{}, nil
		}
		return nil, err
	}
	defer rows.Close()
	
	var folders []model.BookmarkFolder
	for rows.Next() {
		var folder model.BookmarkFolder
		if err := rows.Scan(&folder.ID, &folder.Name, &folder.UserID, &folder.CreatedAt, &folder.BookmarkCount); err != nil {
			return nil, err
		}
		folders = append(folders, folder)
	}
	
	return folders, nil
}

func (r *bookmarkRepository) CreateFolder(ctx context.Context, folder *model.BookmarkFolder) error {
	query := `INSERT INTO bookmark_folders (id, name, user_id, created_at) VALUES ($1, $2, $3, NOW())`
	folder.ID = uuid.New()
	_, err := r.db.Exec(ctx, query, folder.ID, folder.Name, folder.UserID)
	return err
}

func (r *bookmarkRepository) UpdateFolder(ctx context.Context, folder *model.BookmarkFolder) error {
	query := `UPDATE bookmark_folders SET name = $2 WHERE id = $1 AND user_id = $3`
	_, err := r.db.Exec(ctx, query, folder.ID, folder.Name, folder.UserID)
	return err
}

func (r *bookmarkRepository) DeleteFolder(ctx context.Context, id uuid.UUID) error {
	// Move bookmarks to "no folder" before deleting
	_, _ = r.db.Exec(ctx, `UPDATE bookmarks SET folder_id = NULL WHERE folder_id = $1`, id)
	query := `DELETE FROM bookmark_folders WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *bookmarkRepository) MoveToFolder(ctx context.Context, userID, articleID uuid.UUID, folderID *uuid.UUID) error {
	query := `UPDATE bookmarks SET folder_id = $3 WHERE user_id = $1 AND article_id = $2`
	_, err := r.db.Exec(ctx, query, userID, articleID, folderID)
	return err
}

