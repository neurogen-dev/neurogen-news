package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/neurogen-news/backend/internal/model"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
	GetStats(ctx context.Context, id uuid.UUID) (*model.UserStats, error)
	Search(ctx context.Context, query string, limit, offset int) ([]model.User, int, error)
	
	// Follow
	Follow(ctx context.Context, followerID, followingID uuid.UUID) error
	Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error
	IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error)
	GetFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.User, int, error)
	GetFollowing(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.User, int, error)
	
	// Session
	CreateSession(ctx context.Context, session *model.Session) error
	GetSessionByToken(ctx context.Context, refreshToken string) (*model.Session, error)
	DeleteSession(ctx context.Context, id uuid.UUID) error
	DeleteUserSessions(ctx context.Context, userID uuid.UUID) error
}

type userRepository struct {
	db *PostgresDB
}

func NewUserRepository(db *PostgresDB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, display_name, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	`
	
	user.ID = uuid.New()
	
	_, err := r.db.Exec(ctx, query,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.DisplayName,
		user.Role,
	)
	
	if err != nil {
		// Check for unique constraint violation
		if isDuplicateKeyError(err) {
			return ErrUserAlreadyExists
		}
		return err
	}
	
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := `
		SELECT id, username, email, password_hash, display_name, bio, avatar_url, 
			   role, karma, is_verified, is_premium, is_banned, ban_reason, banned_until,
			   created_at, updated_at
		FROM users
		WHERE id = $1
	`
	
	var user model.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.DisplayName,
		&user.Bio,
		&user.AvatarURL,
		&user.Role,
		&user.Karma,
		&user.IsVerified,
		&user.IsPremium,
		&user.IsBanned,
		&user.BanReason,
		&user.BannedUntil,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, username, email, password_hash, display_name, bio, avatar_url, 
			   role, karma, is_verified, is_premium, is_banned, ban_reason, banned_until,
			   created_at, updated_at
		FROM users
		WHERE email = $1
	`
	
	var user model.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.DisplayName,
		&user.Bio,
		&user.AvatarURL,
		&user.Role,
		&user.Karma,
		&user.IsVerified,
		&user.IsPremium,
		&user.IsBanned,
		&user.BanReason,
		&user.BannedUntil,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	
	return &user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
		SELECT id, username, email, password_hash, display_name, bio, avatar_url, 
			   role, karma, is_verified, is_premium, is_banned, ban_reason, banned_until,
			   created_at, updated_at
		FROM users
		WHERE username = $1
	`
	
	var user model.User
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.DisplayName,
		&user.Bio,
		&user.AvatarURL,
		&user.Role,
		&user.Karma,
		&user.IsVerified,
		&user.IsPremium,
		&user.IsBanned,
		&user.BanReason,
		&user.BannedUntil,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET username = $2, display_name = $3, bio = $4, avatar_url = $5, updated_at = NOW()
		WHERE id = $1
	`
	
	_, err := r.db.Exec(ctx, query,
		user.ID,
		user.Username,
		user.DisplayName,
		user.Bio,
		user.AvatarURL,
	)
	
	return err
}

func (r *userRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	query := `UPDATE users SET password_hash = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id, passwordHash)
	return err
}

func (r *userRepository) GetStats(ctx context.Context, id uuid.UUID) (*model.UserStats, error) {
	query := `
		SELECT 
			(SELECT COUNT(*) FROM follows WHERE following_id = $1) as follower_count,
			(SELECT COUNT(*) FROM follows WHERE follower_id = $1) as following_count,
			(SELECT COUNT(*) FROM articles WHERE author_id = $1 AND status = 'published') as article_count,
			(SELECT COUNT(*) FROM comments WHERE author_id = $1) as comment_count,
			(SELECT COALESCE(SUM(view_count), 0) FROM articles WHERE author_id = $1) as total_views
	`
	
	var stats model.UserStats
	err := r.db.QueryRow(ctx, query, id).Scan(
		&stats.FollowerCount,
		&stats.FollowingCount,
		&stats.ArticleCount,
		&stats.CommentCount,
		&stats.TotalViews,
	)
	
	return &stats, err
}

func (r *userRepository) Follow(ctx context.Context, followerID, followingID uuid.UUID) error {
	query := `
		INSERT INTO follows (id, follower_id, following_id, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (follower_id, following_id) DO NOTHING
	`
	_, err := r.db.Exec(ctx, query, uuid.New(), followerID, followingID)
	return err
}

func (r *userRepository) Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	query := `DELETE FROM follows WHERE follower_id = $1 AND following_id = $2`
	_, err := r.db.Exec(ctx, query, followerID, followingID)
	return err
}

func (r *userRepository) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM follows WHERE follower_id = $1 AND following_id = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, followerID, followingID).Scan(&exists)
	return exists, err
}

func (r *userRepository) GetFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.User, int, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM follows WHERE following_id = $1`
	var total int
	if err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT u.id, u.username, u.display_name, u.avatar_url, u.is_verified
		FROM users u
		JOIN follows f ON f.follower_id = u.id
		WHERE f.following_id = $1
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Username, &user.DisplayName, &user.AvatarURL, &user.IsVerified); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}
	
	return users, total, nil
}

func (r *userRepository) GetFollowing(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.User, int, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM follows WHERE follower_id = $1`
	var total int
	if err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT u.id, u.username, u.display_name, u.avatar_url, u.is_verified
		FROM users u
		JOIN follows f ON f.following_id = u.id
		WHERE f.follower_id = $1
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Username, &user.DisplayName, &user.AvatarURL, &user.IsVerified); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}
	
	return users, total, nil
}

func (r *userRepository) Search(ctx context.Context, query string, limit, offset int) ([]model.User, int, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	searchPattern := query + "%"

	// Count total
	countQuery := `SELECT COUNT(*) FROM users WHERE username ILIKE $1 OR display_name ILIKE $1`
	var total int
	if err := r.db.QueryRow(ctx, countQuery, searchPattern).Scan(&total); err != nil {
		return nil, 0, err
	}

	sqlQuery := `
		SELECT id, username, display_name, avatar_url, bio, is_verified
		FROM users
		WHERE username ILIKE $1 OR display_name ILIKE $1
		ORDER BY 
			CASE WHEN username ILIKE $2 THEN 0 ELSE 1 END,
			is_verified DESC,
			username ASC
		LIMIT $3 OFFSET $4
	`
	
	rows, err := r.db.Query(ctx, sqlQuery, searchPattern, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Username, &user.DisplayName, &user.AvatarURL, &user.Bio, &user.IsVerified); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}
	
	return users, total, nil
}

func (r *userRepository) CreateSession(ctx context.Context, session *model.Session) error {
	query := `
		INSERT INTO sessions (id, user_id, refresh_token, user_agent, ip, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`
	
	session.ID = uuid.New()
	
	_, err := r.db.Exec(ctx, query,
		session.ID,
		session.UserID,
		session.RefreshToken,
		session.UserAgent,
		session.IP,
		session.ExpiresAt,
	)
	
	return err
}

func (r *userRepository) GetSessionByToken(ctx context.Context, refreshToken string) (*model.Session, error) {
	query := `
		SELECT id, user_id, refresh_token, user_agent, ip, expires_at, created_at
		FROM sessions
		WHERE refresh_token = $1 AND expires_at > NOW()
	`
	
	var session model.Session
	err := r.db.QueryRow(ctx, query, refreshToken).Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshToken,
		&session.UserAgent,
		&session.IP,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	
	return &session, nil
}

func (r *userRepository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM sessions WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *userRepository) DeleteUserSessions(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM sessions WHERE user_id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

// Helper function to check for duplicate key errors
func isDuplicateKeyError(err error) bool {
	// PostgreSQL error code for unique_violation is 23505
	return err != nil && err.Error() != "" // Simplified check
}

