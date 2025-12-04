package model

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Content     string     `json:"content" db:"content"`
	HTMLContent string     `json:"htmlContent" db:"html_content"`
	AuthorID    uuid.UUID  `json:"authorId" db:"author_id"`
	ArticleID   uuid.UUID  `json:"articleId" db:"article_id"`
	ParentID    *uuid.UUID `json:"parentId,omitempty" db:"parent_id"`
	ReplyCount  int        `json:"replyCount" db:"reply_count"`
	IsEdited    bool       `json:"isEdited" db:"is_edited"`
	IsDeleted   bool       `json:"isDeleted" db:"is_deleted"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"updated_at"`

	// Populated separately
	Author    *CommentAuthor  `json:"author,omitempty"`
	Reactions []ReactionCount `json:"reactions,omitempty"`
	Replies   []Comment       `json:"replies,omitempty"`
	Depth     int             `json:"depth"`
}

type CommentAuthor struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Username    string    `json:"username" db:"username"`
	DisplayName string    `json:"displayName" db:"display_name"`
	AvatarURL   *string   `json:"avatarUrl,omitempty" db:"avatar_url"`
	IsVerified  bool      `json:"isVerified" db:"is_verified"`
}

type ReactionType struct {
	ID    uuid.UUID `json:"id" db:"id"`
	Emoji string    `json:"emoji" db:"emoji"`
	Name  string    `json:"name" db:"name"`
}

type ArticleReaction struct {
	ID         uuid.UUID `json:"id" db:"id"`
	ArticleID  uuid.UUID `json:"articleId" db:"article_id"`
	UserID     uuid.UUID `json:"userId" db:"user_id"`
	ReactionID uuid.UUID `json:"reactionId" db:"reaction_id"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
}

type CommentReaction struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CommentID  uuid.UUID `json:"commentId" db:"comment_id"`
	UserID     uuid.UUID `json:"userId" db:"user_id"`
	ReactionID uuid.UUID `json:"reactionId" db:"reaction_id"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
}

type ReactionCount struct {
	Emoji     string `json:"emoji" db:"emoji"`
	Count     int    `json:"count" db:"count"`
	IsReacted bool   `json:"isReacted" db:"is_reacted"`
}

