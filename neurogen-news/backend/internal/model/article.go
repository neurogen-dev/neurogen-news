package model

import (
	"time"

	"github.com/google/uuid"
)

type ArticleLevel string

const (
	LevelBeginner     ArticleLevel = "beginner"
	LevelIntermediate ArticleLevel = "intermediate"
	LevelAdvanced     ArticleLevel = "advanced"
)

type ContentType string

const (
	ContentTypeArticle    ContentType = "article"
	ContentTypeNews       ContentType = "news"
	ContentTypePost       ContentType = "post"
	ContentTypeQuestion   ContentType = "question"
	ContentTypeDiscussion ContentType = "discussion"
)

type ArticleStatus string

const (
	StatusDraft     ArticleStatus = "draft"
	StatusPending   ArticleStatus = "pending"
	StatusPublished ArticleStatus = "published"
	StatusArchived  ArticleStatus = "archived"
)

type Article struct {
	ID              uuid.UUID     `json:"id" db:"id"`
	Title           string        `json:"title" db:"title"`
	Slug            string        `json:"slug" db:"slug"`
	Lead            *string       `json:"lead,omitempty" db:"lead"`
	Content         string        `json:"content" db:"content"`
	HTMLContent     string        `json:"htmlContent" db:"html_content"`
	CoverImageURL   *string       `json:"coverImageUrl,omitempty" db:"cover_image_url"`
	Level           ArticleLevel  `json:"level" db:"level"`
	ContentType     ContentType   `json:"contentType" db:"content_type"`
	Status          ArticleStatus `json:"status" db:"status"`
	ReadingTime     int           `json:"readingTime" db:"reading_time"`
	IsEditorial     bool          `json:"isEditorial" db:"is_editorial"`
	IsPinned        bool          `json:"isPinned" db:"is_pinned"`
	IsNSFW          bool          `json:"isNsfw" db:"is_nsfw"`
	CommentsEnabled bool          `json:"commentsEnabled" db:"comments_enabled"`
	AuthorID        uuid.UUID     `json:"authorId" db:"author_id"`
	CategoryID      uuid.UUID     `json:"categoryId" db:"category_id"`

	// SEO fields
	MetaTitle       *string `json:"metaTitle,omitempty" db:"meta_title"`
	MetaDescription *string `json:"metaDescription,omitempty" db:"meta_description"`
	CanonicalURL    *string `json:"canonicalUrl,omitempty" db:"canonical_url"`

	// Stats
	ViewCount     int `json:"viewCount" db:"view_count"`
	CommentCount  int `json:"commentCount" db:"comment_count"`
	BookmarkCount int `json:"bookmarkCount" db:"bookmark_count"`

	// Timestamps
	PublishedAt *time.Time `json:"publishedAt,omitempty" db:"published_at"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"updated_at"`

	// Relations (populated separately)
	Author    *User          `json:"author,omitempty"`
	Category  *Category      `json:"category,omitempty"`
	Tags      []Tag          `json:"tags,omitempty"`
	Reactions []ReactionCount `json:"reactions,omitempty"`
}

type ArticleCard struct {
	ID            uuid.UUID    `json:"id" db:"id"`
	Title         string       `json:"title" db:"title"`
	Slug          string       `json:"slug" db:"slug"`
	Lead          *string      `json:"lead,omitempty" db:"lead"`
	CoverImageURL *string      `json:"coverImageUrl,omitempty" db:"cover_image_url"`
	Level         ArticleLevel `json:"level" db:"level"`
	ContentType   ContentType  `json:"contentType" db:"content_type"`
	ReadingTime   int          `json:"readingTime" db:"reading_time"`
	IsEditorial   bool         `json:"isEditorial" db:"is_editorial"`
	IsPinned      bool         `json:"isPinned" db:"is_pinned"`

	AuthorID       uuid.UUID `json:"authorId" db:"author_id"`
	AuthorUsername string    `json:"authorUsername" db:"author_username"`
	AuthorName     string    `json:"authorName" db:"author_display_name"`
	AuthorAvatar   *string   `json:"authorAvatar,omitempty" db:"author_avatar_url"`
	AuthorVerified bool      `json:"authorVerified" db:"author_is_verified"`

	CategoryID   uuid.UUID `json:"categoryId" db:"category_id"`
	CategorySlug string    `json:"categorySlug" db:"category_slug"`
	CategoryName string    `json:"categoryName" db:"category_name"`
	CategoryIcon string    `json:"categoryIcon" db:"category_icon"`

	ViewCount     int `json:"viewCount" db:"view_count"`
	CommentCount  int `json:"commentCount" db:"comment_count"`
	BookmarkCount int `json:"bookmarkCount" db:"bookmark_count"`

	PublishedAt time.Time `json:"publishedAt" db:"published_at"`

	// Populated separately
	Tags      []Tag          `json:"tags,omitempty"`
	Reactions []ReactionCount `json:"reactions,omitempty"`
}

type Category struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	Name            string     `json:"name" db:"name"`
	Slug            string     `json:"slug" db:"slug"`
	Description     *string    `json:"description,omitempty" db:"description"`
	Icon            string     `json:"icon" db:"icon"`
	Color           *string    `json:"color,omitempty" db:"color"`
	IsOfficial      bool       `json:"isOfficial" db:"is_official"`
	ParentID        *uuid.UUID `json:"parentId,omitempty" db:"parent_id"`
	ArticleCount    int        `json:"articleCount" db:"article_count"`
	SubscriberCount int        `json:"subscriberCount" db:"subscriber_count"`
	CreatedAt       time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time  `json:"updatedAt" db:"updated_at"`
}

type Tag struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Slug         string    `json:"slug" db:"slug"`
	ArticleCount int       `json:"articleCount,omitempty" db:"article_count"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}

type ArticleTag struct {
	ArticleID uuid.UUID `json:"articleId" db:"article_id"`
	TagID     uuid.UUID `json:"tagId" db:"tag_id"`
}

type Bookmark struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"userId" db:"user_id"`
	ArticleID uuid.UUID `json:"articleId" db:"article_id"`
	Folder    *string   `json:"folder,omitempty" db:"folder"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`

	// Populated separately
	Article *ArticleCard `json:"article,omitempty"`
}

type Draft struct {
	ID            uuid.UUID   `json:"id" db:"id"`
	UserID        uuid.UUID   `json:"userId" db:"user_id"`
	ArticleID     *uuid.UUID  `json:"articleId,omitempty" db:"article_id"`
	Title         string      `json:"title" db:"title"`
	Content       string      `json:"content" db:"content"`
	CoverImageURL *string     `json:"coverImageUrl,omitempty" db:"cover_image_url"`
	ContentType   ContentType `json:"contentType" db:"content_type"`
	CategoryID    *uuid.UUID  `json:"categoryId,omitempty" db:"category_id"`
	Tags          []string    `json:"tags" db:"tags"`
	IsAutoSave    bool        `json:"isAutoSave" db:"is_auto_save"`
	CreatedAt     time.Time   `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time   `json:"updatedAt" db:"updated_at"`

	// Populated separately
	CategoryName *string `json:"categoryName,omitempty"`
	CategorySlug *string `json:"categorySlug,omitempty"`
}

type BookmarkFolder struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	UserID        uuid.UUID `json:"userId" db:"user_id"`
	BookmarkCount int       `json:"bookmarkCount" db:"bookmark_count"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
}

