package service

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/model"
	"github.com/neurogen-news/backend/internal/repository"
)

type AchievementService interface {
	GetAll(ctx context.Context) ([]model.Achievement, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Achievement, error)
	GetUserAchievements(ctx context.Context, userID uuid.UUID) ([]model.UserAchievement, error)
	GetProgress(ctx context.Context, userID uuid.UUID) (*model.AchievementProgress, error)

	// Check and award achievements
	CheckAndAward(ctx context.Context, userID uuid.UUID) ([]model.Achievement, error)
}

type achievementService struct {
	achievementRepo repository.AchievementRepository
	logger          *zap.Logger
}

func NewAchievementService(
	achievementRepo repository.AchievementRepository,
	logger *zap.Logger,
) AchievementService {
	return &achievementService{
		achievementRepo: achievementRepo,
		logger:          logger,
	}
}

func (s *achievementService) GetAll(ctx context.Context) ([]model.Achievement, error) {
	return s.achievementRepo.GetAll(ctx)
}

func (s *achievementService) GetByID(ctx context.Context, id uuid.UUID) (*model.Achievement, error) {
	return s.achievementRepo.GetByID(ctx, id)
}

func (s *achievementService) GetUserAchievements(ctx context.Context, userID uuid.UUID) ([]model.UserAchievement, error) {
	return s.achievementRepo.GetUserAchievements(ctx, userID)
}

func (s *achievementService) GetProgress(ctx context.Context, userID uuid.UUID) (*model.AchievementProgress, error) {
	return s.achievementRepo.GetProgress(ctx, userID)
}

func (s *achievementService) CheckAndAward(ctx context.Context, userID uuid.UUID) ([]model.Achievement, error) {
	progress, err := s.achievementRepo.GetProgress(ctx, userID)
	if err != nil {
		return nil, err
	}

	allAchievements, err := s.achievementRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var newAchievements []model.Achievement

	for _, achievement := range allAchievements {
		// Check if user already has this achievement
		has, _ := s.achievementRepo.HasAchievement(ctx, userID, achievement.ID)
		if has {
			continue
		}

		// Check conditions based on achievement slug
		shouldAward := false

		switch achievement.Slug {
		case "first-article":
			shouldAward = progress.ArticleCount >= 1
		case "prolific-writer":
			shouldAward = progress.ArticleCount >= 10
		case "author-100":
			shouldAward = progress.ArticleCount >= 100
		case "first-comment":
			shouldAward = progress.CommentCount >= 1
		case "active-commenter":
			shouldAward = progress.CommentCount >= 50
		case "commentator-100":
			shouldAward = progress.CommentCount >= 100
		case "popular-author":
			shouldAward = progress.TotalViews >= 1000
		case "viral-author":
			shouldAward = progress.TotalViews >= 10000
		case "first-follower":
			shouldAward = progress.FollowerCount >= 1
		case "influencer":
			shouldAward = progress.FollowerCount >= 100
		case "celebrity":
			shouldAward = progress.FollowerCount >= 1000
		}

		if shouldAward {
			if err := s.achievementRepo.Award(ctx, userID, achievement.ID); err != nil {
				s.logger.Error("Failed to award achievement",
					zap.String("achievement", achievement.Slug),
					zap.Error(err))
				continue
			}
			newAchievements = append(newAchievements, achievement)
		}
	}

	return newAchievements, nil
}

