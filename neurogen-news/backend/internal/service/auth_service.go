package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/neurogen-news/backend/internal/model"
	"github.com/neurogen-news/backend/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrUserBanned         = errors.New("user is banned")
)

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (*AuthResult, error)
	Login(ctx context.Context, input LoginInput) (*AuthResult, error)
	RefreshToken(ctx context.Context, refreshToken string) (*AuthResult, error)
	Logout(ctx context.Context, userID uuid.UUID, refreshToken string) error
	ValidateToken(ctx context.Context, token string) (*TokenClaims, error)
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
}

type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Remember bool   `json:"remember"`
}

type AuthResult struct {
	User         *model.CurrentUser `json:"user"`
	AccessToken  string             `json:"accessToken"`
	RefreshToken string             `json:"refreshToken"`
	ExpiresAt    int64              `json:"expiresAt"`
}

type TokenClaims struct {
	UserID   uuid.UUID      `json:"userId"`
	Username string         `json:"username"`
	Role     model.UserRole `json:"role"`
	jwt.RegisteredClaims
}

type authService struct {
	userRepo  repository.UserRepository
	redis     *repository.RedisClient
	jwtSecret []byte
	logger    *zap.Logger
}

func NewAuthService(
	userRepo repository.UserRepository,
	redis *repository.RedisClient,
	jwtSecret string,
	logger *zap.Logger,
) AuthService {
	return &authService{
		userRepo:  userRepo,
		redis:     redis,
		jwtSecret: []byte(jwtSecret),
		logger:    logger,
	}
}

func (s *authService) Register(ctx context.Context, input RegisterInput) (*AuthResult, error) {
	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &model.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(passwordHash),
		DisplayName:  input.Username,
		Role:         model.RoleUser,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate tokens
	return s.generateTokens(ctx, user, "", "")
}

func (s *authService) Login(ctx context.Context, input LoginInput) (*AuthResult, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Check if user is banned
	if user.IsBanned {
		return nil, ErrUserBanned
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.generateTokens(ctx, user, "", "")
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*AuthResult, error) {
	// Get session by refresh token
	session, err := s.userRepo.GetSessionByToken(ctx, refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	if user.IsBanned {
		return nil, ErrUserBanned
	}

	// Delete old session
	_ = s.userRepo.DeleteSession(ctx, session.ID)

	return s.generateTokens(ctx, user, session.UserAgent, session.IP)
}

func (s *authService) Logout(ctx context.Context, userID uuid.UUID, refreshToken string) error {
	// Delete session by refresh token
	session, err := s.userRepo.GetSessionByToken(ctx, refreshToken)
	if err != nil {
		return nil
	}

	if session.UserID != userID {
		return nil
	}

	return s.userRepo.DeleteSession(ctx, session.ID)
}

func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (s *authService) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// Don't reveal if user exists
		return nil
	}

	// Generate reset token
	resetToken := uuid.New().String()
	
	// Store in Redis with 1 hour expiry
	key := "password_reset:" + resetToken
	if err := s.redis.Set(ctx, key, user.ID.String(), time.Hour).Err(); err != nil {
		return err
	}

	// TODO: Send email with reset link
	s.logger.Info("Password reset requested",
		zap.String("email", email),
		zap.String("token", resetToken),
	)

	return nil
}

func (s *authService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Get user ID from Redis
	key := "password_reset:" + token
	userIDStr, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		return ErrInvalidToken
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ErrInvalidToken
	}

	// Hash new password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	if err := s.userRepo.UpdatePassword(ctx, userID, string(passwordHash)); err != nil {
		return err
	}

	// Delete reset token
	s.redis.Del(ctx, key)

	// Invalidate all sessions
	_ = s.userRepo.DeleteUserSessions(ctx, userID)

	return nil
}

func (s *authService) generateTokens(ctx context.Context, user *model.User, userAgent, ip string) (*AuthResult, error) {
	// Access token - 15 minutes
	accessExpiry := time.Now().Add(15 * time.Minute)
	accessClaims := TokenClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	// Refresh token - 7 days
	refreshExpiry := time.Now().Add(7 * 24 * time.Hour)
	refreshToken := uuid.New().String()

	// Save session
	session := &model.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		IP:           ip,
		ExpiresAt:    refreshExpiry,
	}

	if err := s.userRepo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	// Get user stats for CurrentUser
	stats, _ := s.userRepo.GetStats(ctx, user.ID)

	currentUser := &model.CurrentUser{
		User:                *user,
		Email:               user.Email,
		UnreadNotifications: 0, // TODO: Get from notifications
		DraftCount:          0, // TODO: Get from drafts
		BookmarkCount:       0, // TODO: Get from bookmarks
	}

	if stats != nil {
		currentUser.FollowerCount = stats.FollowerCount
		currentUser.FollowingCount = stats.FollowingCount
		currentUser.ArticleCount = stats.ArticleCount
	}

	return &AuthResult{
		User:         currentUser,
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiry.Unix(),
	}, nil
}

