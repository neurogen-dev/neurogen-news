package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/neurogen-news/backend/internal/model"
)

var (
	ErrCommentNotFound = errors.New("comment not found")
)

type CommentRepository interface {
	Create(ctx context.Context, comment *model.Comment) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Comment, error)
	Update(ctx context.Context, comment *model.Comment) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Lists
	GetByArticle(ctx context.Context, articleID uuid.UUID, params CommentListParams) ([]model.Comment, int, error)
	GetReplies(ctx context.Context, parentID uuid.UUID, limit, offset int) ([]model.Comment, error)
	GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Comment, int, error)

	// Reactions
	AddReaction(ctx context.Context, commentID, userID, reactionID uuid.UUID) error
	RemoveReaction(ctx context.Context, commentID, userID uuid.UUID) error
	GetReactions(ctx context.Context, commentID uuid.UUID, userID *uuid.UUID) ([]model.ReactionCount, error)
}

type CommentListParams struct {
	Sort   string // new, popular, old
	Limit  int
	Offset int
}

type commentRepository struct {
	db *PostgresDB
}

func NewCommentRepository(db *PostgresDB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, comment *model.Comment) error {
	query := `
		INSERT INTO comments (id, content, html_content, author_id, article_id, parent_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	`

	comment.ID = uuid.New()

	_, err := r.db.Exec(ctx, query,
		comment.ID,
		comment.Content,
		comment.HTMLContent,
		comment.AuthorID,
		comment.ArticleID,
		comment.ParentID,
	)

	if err != nil {
		return err
	}

	// Increment reply count if this is a reply
	if comment.ParentID != nil {
		_, _ = r.db.Exec(ctx, `UPDATE comments SET reply_count = reply_count + 1 WHERE id = $1`, *comment.ParentID)
	}

	return nil
}

func (r *commentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Comment, error) {
	query := `
		SELECT 
			c.id, c.content, c.html_content, c.author_id, c.article_id, c.parent_id,
			c.reply_count, c.is_edited, c.is_deleted, c.created_at, c.updated_at,
			u.id, u.username, u.display_name, u.avatar_url, u.is_verified
		FROM comments c
		JOIN users u ON u.id = c.author_id
		WHERE c.id = $1
	`

	var comment model.Comment
	comment.Author = &model.CommentAuthor{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&comment.ID,
		&comment.Content,
		&comment.HTMLContent,
		&comment.AuthorID,
		&comment.ArticleID,
		&comment.ParentID,
		&comment.ReplyCount,
		&comment.IsEdited,
		&comment.IsDeleted,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.Author.ID,
		&comment.Author.Username,
		&comment.Author.DisplayName,
		&comment.Author.AvatarURL,
		&comment.Author.IsVerified,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCommentNotFound
		}
		return nil, err
	}

	return &comment, nil
}

func (r *commentRepository) Update(ctx context.Context, comment *model.Comment) error {
	query := `
		UPDATE comments SET
			content = $2,
			html_content = $3,
			is_edited = true,
			updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query,
		comment.ID,
		comment.Content,
		comment.HTMLContent,
	)

	return err
}

func (r *commentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Soft delete - mark as deleted instead of removing
	query := `
		UPDATE comments SET
			is_deleted = true,
			content = '[Комментарий удалён]',
			html_content = '<p>[Комментарий удалён]</p>',
			updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *commentRepository) GetByArticle(ctx context.Context, articleID uuid.UUID, params CommentListParams) ([]model.Comment, int, error) {
	// Set defaults
	if params.Limit <= 0 || params.Limit > 100 {
		params.Limit = 20
	}

	// Sort order
	var orderBy string
	switch params.Sort {
	case "popular":
		orderBy = "c.reply_count DESC, c.created_at DESC"
	case "old":
		orderBy = "c.created_at ASC"
	default: // new
		orderBy = "c.created_at DESC"
	}

	// Count total root comments
	countQuery := `SELECT COUNT(*) FROM comments WHERE article_id = $1 AND parent_id IS NULL`
	var total int
	if err := r.db.QueryRow(ctx, countQuery, articleID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get root comments with authors
	query := fmt.Sprintf(`
		SELECT 
			c.id, c.content, c.html_content, c.author_id, c.article_id, c.parent_id,
			c.reply_count, c.is_edited, c.is_deleted, c.created_at, c.updated_at,
			u.id, u.username, u.display_name, u.avatar_url, u.is_verified
		FROM comments c
		JOIN users u ON u.id = c.author_id
		WHERE c.article_id = $1 AND c.parent_id IS NULL
		ORDER BY %s
		LIMIT $2 OFFSET $3
	`, orderBy)

	rows, err := r.db.Query(ctx, query, articleID, params.Limit, params.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		comment.Author = &model.CommentAuthor{}

		err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.HTMLContent,
			&comment.AuthorID,
			&comment.ArticleID,
			&comment.ParentID,
			&comment.ReplyCount,
			&comment.IsEdited,
			&comment.IsDeleted,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.Author.ID,
			&comment.Author.Username,
			&comment.Author.DisplayName,
			&comment.Author.AvatarURL,
			&comment.Author.IsVerified,
		)
		if err != nil {
			return nil, 0, err
		}
		comment.Depth = 0
		comments = append(comments, comment)
	}

	return comments, total, nil
}

func (r *commentRepository) GetReplies(ctx context.Context, parentID uuid.UUID, limit, offset int) ([]model.Comment, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	query := `
		SELECT 
			c.id, c.content, c.html_content, c.author_id, c.article_id, c.parent_id,
			c.reply_count, c.is_edited, c.is_deleted, c.created_at, c.updated_at,
			u.id, u.username, u.display_name, u.avatar_url, u.is_verified
		FROM comments c
		JOIN users u ON u.id = c.author_id
		WHERE c.parent_id = $1
		ORDER BY c.created_at ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, parentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []model.Comment
	for rows.Next() {
		var comment model.Comment
		comment.Author = &model.CommentAuthor{}

		err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.HTMLContent,
			&comment.AuthorID,
			&comment.ArticleID,
			&comment.ParentID,
			&comment.ReplyCount,
			&comment.IsEdited,
			&comment.IsDeleted,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.Author.ID,
			&comment.Author.Username,
			&comment.Author.DisplayName,
			&comment.Author.AvatarURL,
			&comment.Author.IsVerified,
		)
		if err != nil {
			return nil, err
		}
		comment.Depth = 1
		replies = append(replies, comment)
	}

	return replies, nil
}

func (r *commentRepository) GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Comment, int, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	countQuery := `SELECT COUNT(*) FROM comments WHERE author_id = $1 AND is_deleted = false`
	var total int
	if err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			c.id, c.content, c.html_content, c.author_id, c.article_id, c.parent_id,
			c.reply_count, c.is_edited, c.is_deleted, c.created_at, c.updated_at,
			u.id, u.username, u.display_name, u.avatar_url, u.is_verified
		FROM comments c
		JOIN users u ON u.id = c.author_id
		WHERE c.author_id = $1 AND c.is_deleted = false
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		comment.Author = &model.CommentAuthor{}

		err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.HTMLContent,
			&comment.AuthorID,
			&comment.ArticleID,
			&comment.ParentID,
			&comment.ReplyCount,
			&comment.IsEdited,
			&comment.IsDeleted,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.Author.ID,
			&comment.Author.Username,
			&comment.Author.DisplayName,
			&comment.Author.AvatarURL,
			&comment.Author.IsVerified,
		)
		if err != nil {
			return nil, 0, err
		}
		comments = append(comments, comment)
	}

	return comments, total, nil
}

func (r *commentRepository) AddReaction(ctx context.Context, commentID, userID, reactionID uuid.UUID) error {
	query := `
		INSERT INTO comment_reactions (id, comment_id, user_id, reaction_id, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (comment_id, user_id) DO UPDATE SET reaction_id = $4, created_at = NOW()
	`

	_, err := r.db.Exec(ctx, query, uuid.New(), commentID, userID, reactionID)
	return err
}

func (r *commentRepository) RemoveReaction(ctx context.Context, commentID, userID uuid.UUID) error {
	query := `DELETE FROM comment_reactions WHERE comment_id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, commentID, userID)
	return err
}

func (r *commentRepository) GetReactions(ctx context.Context, commentID uuid.UUID, userID *uuid.UUID) ([]model.ReactionCount, error) {
	query := `
		SELECT 
			rt.emoji,
			COUNT(*) as count,
			COALESCE(bool_or(cr.user_id = $2), false) as is_reacted
		FROM comment_reactions cr
		JOIN reaction_types rt ON rt.id = cr.reaction_id
		WHERE cr.comment_id = $1
		GROUP BY rt.emoji
		ORDER BY count DESC
	`

	var uid interface{} = nil
	if userID != nil {
		uid = *userID
	}

	rows, err := r.db.Query(ctx, query, commentID, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reactions []model.ReactionCount
	for rows.Next() {
		var rc model.ReactionCount
		if err := rows.Scan(&rc.Emoji, &rc.Count, &rc.IsReacted); err != nil {
			return nil, err
		}
		reactions = append(reactions, rc)
	}

	return reactions, nil
}
