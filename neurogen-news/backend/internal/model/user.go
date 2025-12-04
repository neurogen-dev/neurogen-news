package model

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleUser      UserRole = "USER"
	RoleAuthor    UserRole = "AUTHOR"
	RoleModerator UserRole = "MODERATOR"
	RoleEditor    UserRole = "EDITOR"
	RoleAdmin     UserRole = "ADMIN"
)

type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Username     string     `json:"username" db:"username"`
	Email        string     `json:"email,omitempty" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"`
	DisplayName  string     `json:"displayName" db:"display_name"`
	Bio          *string    `json:"bio,omitempty" db:"bio"`
	AvatarURL    *string    `json:"avatarUrl,omitempty" db:"avatar_url"`
	CoverURL     *string    `json:"coverUrl,omitempty" db:"cover_url"`
	Location     *string    `json:"location,omitempty" db:"location"`
	Website      *string    `json:"website,omitempty" db:"website"`
	Telegram     *string    `json:"telegram,omitempty" db:"telegram"`
	Github       *string    `json:"github,omitempty" db:"github"`
	Role         UserRole   `json:"role" db:"role"`
	Karma        int        `json:"karma" db:"karma"`
	IsVerified   bool       `json:"isVerified" db:"is_verified"`
	IsPremium    bool       `json:"isPremium" db:"is_premium"`
	IsBanned     bool       `json:"-" db:"is_banned"`
	BanReason    *string    `json:"-" db:"ban_reason"`
	BannedUntil  *time.Time `json:"-" db:"banned_until"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time  `json:"updatedAt" db:"updated_at"`

	// Computed fields (not in DB)
	FollowerCount  int `json:"followerCount,omitempty"`
	FollowingCount int `json:"followingCount,omitempty"`
	ArticleCount   int `json:"articleCount,omitempty"`
}

type UserStats struct {
	FollowerCount  int `json:"followerCount" db:"follower_count"`
	FollowingCount int `json:"followingCount" db:"following_count"`
	ArticleCount   int `json:"articleCount" db:"article_count"`
	CommentCount   int `json:"commentCount" db:"comment_count"`
	TotalViews     int `json:"totalViews" db:"total_views"`
}

type CurrentUser struct {
	User
	Email               string `json:"email"`
	UnreadNotifications int    `json:"unreadNotifications"`
	DraftCount          int    `json:"draftCount"`
	BookmarkCount       int    `json:"bookmarkCount"`
}

type Follow struct {
	ID          uuid.UUID `json:"id" db:"id"`
	FollowerID  uuid.UUID `json:"followerId" db:"follower_id"`
	FollowingID uuid.UUID `json:"followingId" db:"following_id"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

type Session struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"userId" db:"user_id"`
	RefreshToken string    `json:"-" db:"refresh_token"`
	UserAgent    string    `json:"userAgent" db:"user_agent"`
	IP           string    `json:"ip" db:"ip"`
	ExpiresAt    time.Time `json:"expiresAt" db:"expires_at"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}

type AuthProvider struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"userId" db:"user_id"`
	Provider     string    `json:"provider" db:"provider"` // google, vk, telegram
	ProviderID   string    `json:"providerId" db:"provider_id"`
	AccessToken  string    `json:"-" db:"access_token"`
	RefreshToken *string   `json:"-" db:"refresh_token"`
	ExpiresAt    time.Time `json:"expiresAt" db:"expires_at"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}

