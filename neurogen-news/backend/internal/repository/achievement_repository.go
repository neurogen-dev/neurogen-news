package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/neurogen-news/backend/internal/model"
)

var (
	ErrAchievementNotFound = errors.New("achievement not found")
)

type AchievementRepository interface {
	GetAll(ctx context.Context) ([]model.Achievement, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Achievement, error)
	GetBySlug(ctx context.Context, slug string) (*model.Achievement, error)
	
	// User achievements
	GetUserAchievements(ctx context.Context, userID uuid.UUID) ([]model.UserAchievement, error)
	Award(ctx context.Context, userID, achievementID uuid.UUID) error
	HasAchievement(ctx context.Context, userID, achievementID uuid.UUID) (bool, error)
	
	// Progress tracking
	GetProgress(ctx context.Context, userID uuid.UUID) (*model.AchievementProgress, error)
	UpdateProgress(ctx context.Context, userID uuid.UUID, field string, increment int) error
}

type achievementRepository struct {
	db *PostgresDB
}

func NewAchievementRepository(db *PostgresDB) AchievementRepository {
	return &achievementRepository{db: db}
}

func (r *achievementRepository) GetAll(ctx context.Context) ([]model.Achievement, error) {
	query := `
		SELECT id, name, slug, description, icon, points, category, is_hidden, created_at
		FROM achievements
		WHERE is_hidden = false
		ORDER BY category, points ASC
	`
	
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var achievements []model.Achievement
	for rows.Next() {
		var a model.Achievement
		err := rows.Scan(
			&a.ID, &a.Name, &a.Slug, &a.Description, &a.Icon, &a.Points, &a.Category, &a.IsHidden, &a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		achievements = append(achievements, a)
	}
	
	return achievements, nil
}

func (r *achievementRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Achievement, error) {
	query := `
		SELECT id, name, slug, description, icon, points, category, is_hidden, created_at
		FROM achievements
		WHERE id = $1
	`
	
	var a model.Achievement
	err := r.db.QueryRow(ctx, query, id).Scan(
		&a.ID, &a.Name, &a.Slug, &a.Description, &a.Icon, &a.Points, &a.Category, &a.IsHidden, &a.CreatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAchievementNotFound
		}
		return nil, err
	}
	
	return &a, nil
}

func (r *achievementRepository) GetBySlug(ctx context.Context, slug string) (*model.Achievement, error) {
	query := `
		SELECT id, name, slug, description, icon, points, category, is_hidden, created_at
		FROM achievements
		WHERE slug = $1
	`
	
	var a model.Achievement
	err := r.db.QueryRow(ctx, query, slug).Scan(
		&a.ID, &a.Name, &a.Slug, &a.Description, &a.Icon, &a.Points, &a.Category, &a.IsHidden, &a.CreatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAchievementNotFound
		}
		return nil, err
	}
	
	return &a, nil
}

func (r *achievementRepository) GetUserAchievements(ctx context.Context, userID uuid.UUID) ([]model.UserAchievement, error) {
	query := `
		SELECT ua.id, ua.user_id, ua.achievement_id, ua.awarded_at,
		       a.name, a.slug, a.description, a.icon, a.points, a.category
		FROM user_achievements ua
		JOIN achievements a ON a.id = ua.achievement_id
		WHERE ua.user_id = $1
		ORDER BY ua.awarded_at DESC
	`
	
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var userAchievements []model.UserAchievement
	for rows.Next() {
		var ua model.UserAchievement
		ua.Achievement = &model.Achievement{}
		
		err := rows.Scan(
			&ua.ID, &ua.UserID, &ua.AchievementID, &ua.AwardedAt,
			&ua.Achievement.Name, &ua.Achievement.Slug, &ua.Achievement.Description,
			&ua.Achievement.Icon, &ua.Achievement.Points, &ua.Achievement.Category,
		)
		if err != nil {
			return nil, err
		}
		ua.Achievement.ID = ua.AchievementID
		userAchievements = append(userAchievements, ua)
	}
	
	return userAchievements, nil
}

func (r *achievementRepository) Award(ctx context.Context, userID, achievementID uuid.UUID) error {
	query := `
		INSERT INTO user_achievements (id, user_id, achievement_id, awarded_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id, achievement_id) DO NOTHING
	`
	_, err := r.db.Exec(ctx, query, uuid.New(), userID, achievementID)
	return err
}

func (r *achievementRepository) HasAchievement(ctx context.Context, userID, achievementID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_achievements WHERE user_id = $1 AND achievement_id = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, achievementID).Scan(&exists)
	return exists, err
}

func (r *achievementRepository) GetProgress(ctx context.Context, userID uuid.UUID) (*model.AchievementProgress, error) {
	// Get user stats for achievement progress
	query := `
		SELECT 
			(SELECT COUNT(*) FROM articles WHERE author_id = $1 AND status = 'published') as article_count,
			(SELECT COUNT(*) FROM comments WHERE author_id = $1 AND is_deleted = false) as comment_count,
			(SELECT COALESCE(SUM(view_count), 0) FROM articles WHERE author_id = $1) as total_views,
			(SELECT COUNT(*) FROM follows WHERE following_id = $1) as follower_count,
			(SELECT COUNT(*) FROM user_achievements WHERE user_id = $1) as achievement_count,
			(SELECT COALESCE(SUM(a.points), 0) FROM user_achievements ua JOIN achievements a ON a.id = ua.achievement_id WHERE ua.user_id = $1) as total_points
	`
	
	var progress model.AchievementProgress
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&progress.ArticleCount,
		&progress.CommentCount,
		&progress.TotalViews,
		&progress.FollowerCount,
		&progress.AchievementCount,
		&progress.TotalPoints,
	)
	
	if err != nil {
		return nil, err
	}
	
	progress.UserID = userID
	
	return &progress, nil
}

func (r *achievementRepository) UpdateProgress(ctx context.Context, userID uuid.UUID, field string, increment int) error {
	// This would be used for tracking specific progress
	// For now, we calculate progress dynamically in GetProgress
	return nil
}

