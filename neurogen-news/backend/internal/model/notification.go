package model

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationNewComment      NotificationType = "new_comment"
	NotificationCommentReply    NotificationType = "comment_reply"
	NotificationReaction        NotificationType = "reaction"
	NotificationNewFollower     NotificationType = "new_follower"
	NotificationArticlePublished NotificationType = "article_published"
	NotificationMention         NotificationType = "mention"
	NotificationSystem          NotificationType = "system"
)

type Notification struct {
	ID        uuid.UUID        `json:"id" db:"id"`
	UserID    uuid.UUID        `json:"userId" db:"user_id"`
	Type      NotificationType `json:"type" db:"type"`
	Title     string           `json:"title" db:"title"`
	Message   string           `json:"message" db:"message"`
	Link      *string          `json:"link,omitempty" db:"link"`
	ImageURL  *string          `json:"imageUrl,omitempty" db:"image_url"`
	ActorID   *uuid.UUID       `json:"actorId,omitempty" db:"actor_id"`
	ArticleID *uuid.UUID       `json:"articleId,omitempty" db:"article_id"`
	CommentID *uuid.UUID       `json:"commentId,omitempty" db:"comment_id"`
	IsRead    bool             `json:"isRead" db:"is_read"`
	CreatedAt time.Time        `json:"createdAt" db:"created_at"`

	// Populated separately
	Actor *NotificationActor `json:"actor,omitempty"`
}

type NotificationActor struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Username    string    `json:"username" db:"username"`
	DisplayName string    `json:"displayName" db:"display_name"`
	AvatarURL   *string   `json:"avatarUrl,omitempty" db:"avatar_url"`
}

type Achievement struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Description string    `json:"description" db:"description"`
	Icon        string    `json:"icon" db:"icon"`
	Points      int       `json:"points" db:"points"`
	Category    string    `json:"category" db:"category"` // author, commentator, social, special
	IsHidden    bool      `json:"isHidden" db:"is_hidden"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

type UserAchievement struct {
	ID            uuid.UUID `json:"id" db:"id"`
	UserID        uuid.UUID `json:"userId" db:"user_id"`
	AchievementID uuid.UUID `json:"achievementId" db:"achievement_id"`
	AwardedAt     time.Time `json:"awardedAt" db:"awarded_at"`

	// Populated separately
	Achievement *Achievement `json:"achievement,omitempty"`
}

type AchievementProgress struct {
	UserID           uuid.UUID `json:"userId" db:"user_id"`
	ArticleCount     int       `json:"articleCount" db:"article_count"`
	CommentCount     int       `json:"commentCount" db:"comment_count"`
	TotalViews       int       `json:"totalViews" db:"total_views"`
	FollowerCount    int       `json:"followerCount" db:"follower_count"`
	AchievementCount int       `json:"achievementCount" db:"achievement_count"`
	TotalPoints      int       `json:"totalPoints" db:"total_points"`
}

type Report struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ReporterID  uuid.UUID `json:"reporterId" db:"reporter_id"`
	TargetType  string    `json:"targetType" db:"target_type"` // article, comment, user
	TargetID    uuid.UUID `json:"targetId" db:"target_id"`
	Reason      string    `json:"reason" db:"reason"`
	Description *string   `json:"description,omitempty" db:"description"`
	Status      string    `json:"status" db:"status"` // pending, resolved, dismissed
	ResolvedBy  *uuid.UUID `json:"resolvedBy,omitempty" db:"resolved_by"`
	ResolvedAt  *time.Time `json:"resolvedAt,omitempty" db:"resolved_at"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
}

