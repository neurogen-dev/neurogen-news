package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/neurogen-news/backend/internal/model"
)

var (
	ErrReactionTypeNotFound = errors.New("reaction type not found")
)

type ReactionRepository interface {
	// Reaction types
	GetAllTypes(ctx context.Context) ([]model.ReactionType, error)
	GetByEmoji(ctx context.Context, emoji string) (*model.ReactionType, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.ReactionType, error)
	
	// Article reactions
	AddArticleReaction(ctx context.Context, articleID, userID, reactionID uuid.UUID) error
	RemoveArticleReaction(ctx context.Context, articleID, userID uuid.UUID) error
	GetArticleReactions(ctx context.Context, articleID uuid.UUID, userID *uuid.UUID) ([]model.ReactionCount, error)
	GetUserArticleReaction(ctx context.Context, articleID, userID uuid.UUID) (*model.ReactionType, error)
}

type reactionRepository struct {
	db *PostgresDB
}

func NewReactionRepository(db *PostgresDB) ReactionRepository {
	return &reactionRepository{db: db}
}

func (r *reactionRepository) GetAllTypes(ctx context.Context) ([]model.ReactionType, error) {
	query := `SELECT id, emoji, name FROM reaction_types ORDER BY name`
	
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var types []model.ReactionType
	for rows.Next() {
		var rt model.ReactionType
		if err := rows.Scan(&rt.ID, &rt.Emoji, &rt.Name); err != nil {
			return nil, err
		}
		types = append(types, rt)
	}
	
	return types, nil
}

func (r *reactionRepository) GetByEmoji(ctx context.Context, emoji string) (*model.ReactionType, error) {
	query := `SELECT id, emoji, name FROM reaction_types WHERE emoji = $1`
	
	var rt model.ReactionType
	err := r.db.QueryRow(ctx, query, emoji).Scan(&rt.ID, &rt.Emoji, &rt.Name)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrReactionTypeNotFound
		}
		return nil, err
	}
	
	return &rt, nil
}

func (r *reactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.ReactionType, error) {
	query := `SELECT id, emoji, name FROM reaction_types WHERE id = $1`
	
	var rt model.ReactionType
	err := r.db.QueryRow(ctx, query, id).Scan(&rt.ID, &rt.Emoji, &rt.Name)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrReactionTypeNotFound
		}
		return nil, err
	}
	
	return &rt, nil
}

func (r *reactionRepository) AddArticleReaction(ctx context.Context, articleID, userID, reactionID uuid.UUID) error {
	query := `
		INSERT INTO article_reactions (id, article_id, user_id, reaction_id, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (article_id, user_id) DO UPDATE SET reaction_id = $4, created_at = NOW()
	`
	
	_, err := r.db.Exec(ctx, query, uuid.New(), articleID, userID, reactionID)
	return err
}

func (r *reactionRepository) RemoveArticleReaction(ctx context.Context, articleID, userID uuid.UUID) error {
	query := `DELETE FROM article_reactions WHERE article_id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, articleID, userID)
	return err
}

func (r *reactionRepository) GetArticleReactions(ctx context.Context, articleID uuid.UUID, userID *uuid.UUID) ([]model.ReactionCount, error) {
	query := `
		SELECT 
			rt.emoji,
			COUNT(*) as count,
			COALESCE(bool_or(ar.user_id = $2), false) as is_reacted
		FROM article_reactions ar
		JOIN reaction_types rt ON rt.id = ar.reaction_id
		WHERE ar.article_id = $1
		GROUP BY rt.emoji
		ORDER BY count DESC
	`
	
	var uid interface{} = nil
	if userID != nil {
		uid = *userID
	}
	
	rows, err := r.db.Query(ctx, query, articleID, uid)
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

func (r *reactionRepository) GetUserArticleReaction(ctx context.Context, articleID, userID uuid.UUID) (*model.ReactionType, error) {
	query := `
		SELECT rt.id, rt.emoji, rt.name
		FROM article_reactions ar
		JOIN reaction_types rt ON rt.id = ar.reaction_id
		WHERE ar.article_id = $1 AND ar.user_id = $2
	`
	
	var rt model.ReactionType
	err := r.db.QueryRow(ctx, query, articleID, userID).Scan(&rt.ID, &rt.Emoji, &rt.Name)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // No reaction is not an error
		}
		return nil, err
	}
	
	return &rt, nil
}

