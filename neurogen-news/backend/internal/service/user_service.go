package service

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/model"
	"github.com/neurogen-news/backend/internal/repository"
)

type UserService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*UserProfile, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, input UpdateProfileInput) (*model.User, error)
	UpdateSettings(ctx context.Context, userID uuid.UUID, input UpdateSettingsInput) error

	// Follows
	Follow(ctx context.Context, followerID, followingID uuid.UUID) error
	Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error
	IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error)
	GetFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.User, int, error)
	GetFollowing(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.User, int, error)

	// Articles
	GetUserArticles(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.ArticleCard, int, error)
}

type UserProfile struct {
	User           *model.User `json:"user"`
	ArticleCount   int         `json:"articleCount"`
	CommentCount   int         `json:"commentCount"`
	FollowerCount  int         `json:"followerCount"`
	FollowingCount int         `json:"followingCount"`
	IsFollowing    bool        `json:"isFollowing"`
}

type UpdateProfileInput struct {
	DisplayName *string `json:"displayName,omitempty" validate:"omitempty,min=2,max=100"`
	Bio         *string `json:"bio,omitempty" validate:"omitempty,max=500"`
	AvatarURL   *string `json:"avatarUrl,omitempty"`
	CoverURL    *string `json:"coverUrl,omitempty"`
	Location    *string `json:"location,omitempty" validate:"omitempty,max=100"`
	Website     *string `json:"website,omitempty" validate:"omitempty,url"`
	Telegram    *string `json:"telegram,omitempty" validate:"omitempty,max=50"`
	Github      *string `json:"github,omitempty" validate:"omitempty,max=50"`
}

type UpdateSettingsInput struct {
	EmailNotifications    *bool   `json:"emailNotifications,omitempty"`
	PushNotifications     *bool   `json:"pushNotifications,omitempty"`
	Theme                 *string `json:"theme,omitempty" validate:"omitempty,oneof=light dark system"`
	Language              *string `json:"language,omitempty" validate:"omitempty,oneof=ru en"`
	ShowEmail             *bool   `json:"showEmail,omitempty"`
	AllowDirectMessages   *bool   `json:"allowDirectMessages,omitempty"`
}

type userService struct {
	userRepo    repository.UserRepository
	articleRepo repository.ArticleRepository
	redis       *repository.RedisClient
	logger      *zap.Logger
}

func NewUserService(
	userRepo repository.UserRepository,
	articleRepo repository.ArticleRepository,
	redis *repository.RedisClient,
	logger *zap.Logger,
) UserService {
	return &userService{
		userRepo:    userRepo,
		articleRepo: articleRepo,
		redis:       redis,
		logger:      logger,
	}
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	return s.userRepo.GetByUsername(ctx, username)
}

func (s *userService) GetProfile(ctx context.Context, userID uuid.UUID) (*UserProfile, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	stats, err := s.userRepo.GetStats(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user stats", zap.Error(err))
		stats = &model.UserStats{}
	}

	return &UserProfile{
		User:           user,
		ArticleCount:   stats.ArticleCount,
		CommentCount:   stats.CommentCount,
		FollowerCount:  stats.FollowerCount,
		FollowingCount: stats.FollowingCount,
	}, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID uuid.UUID, input UpdateProfileInput) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if input.DisplayName != nil {
		user.DisplayName = *input.DisplayName
	}
	if input.Bio != nil {
		user.Bio = input.Bio
	}
	if input.AvatarURL != nil {
		user.AvatarURL = input.AvatarURL
	}
	if input.CoverURL != nil {
		user.CoverURL = input.CoverURL
	}
	if input.Location != nil {
		user.Location = input.Location
	}
	if input.Website != nil {
		user.Website = input.Website
	}
	if input.Telegram != nil {
		user.Telegram = input.Telegram
	}
	if input.Github != nil {
		user.Github = input.Github
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateSettings(ctx context.Context, userID uuid.UUID, input UpdateSettingsInput) error {
	// TODO: Implement user settings update
	return nil
}

func (s *userService) Follow(ctx context.Context, followerID, followingID uuid.UUID) error {
	if followerID == followingID {
		return &AppError{Code: "SELF_FOLLOW", Message: "Cannot follow yourself"}
	}
	return s.userRepo.Follow(ctx, followerID, followingID)
}

func (s *userService) Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	return s.userRepo.Unfollow(ctx, followerID, followingID)
}

func (s *userService) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	return s.userRepo.IsFollowing(ctx, followerID, followingID)
}

func (s *userService) GetFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.User, int, error) {
	return s.userRepo.GetFollowers(ctx, userID, limit, offset)
}

func (s *userService) GetFollowing(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.User, int, error) {
	return s.userRepo.GetFollowing(ctx, userID, limit, offset)
}

func (s *userService) GetUserArticles(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.ArticleCard, int, error) {
	return s.articleRepo.GetByAuthor(ctx, userID, limit, offset)
}

