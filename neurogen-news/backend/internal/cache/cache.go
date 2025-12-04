package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/neurogen-news/backend/internal/model"
	"github.com/neurogen-news/backend/internal/repository"
)

// Cache keys
const (
	KeyArticleList     = "articles:list:%s:%s:%s:%d"    // sort:level:contentType:page
	KeyArticle         = "article:%s"                    // articleID
	KeyArticleBySlug   = "article:slug:%s:%s"           // category:slug
	KeyCategories      = "categories:all"
	KeyCategory        = "category:%s"                   // slug
	KeyPopularTags     = "tags:popular:%d"              // limit
	KeyUserProfile     = "user:profile:%s"              // userID
	KeyNotificationCnt = "notifications:count:%s"       // userID
)

// TTL durations
const (
	TTLShort  = 1 * time.Minute
	TTLMedium = 5 * time.Minute
	TTLLong   = 30 * time.Minute
	TTLHour   = 1 * time.Hour
	TTLDay    = 24 * time.Hour
)

type Cache struct {
	redis *repository.RedisClient
}

func New(redis *repository.RedisClient) *Cache {
	return &Cache{redis: redis}
}

// ============================================
// Articles cache
// ============================================

func (c *Cache) GetArticleList(ctx context.Context, sort, level, contentType string, page int) ([]model.ArticleCard, bool) {
	key := fmt.Sprintf(KeyArticleList, sort, level, contentType, page)
	
	var articles []model.ArticleCard
	err := c.redis.GetJSON(ctx, key, &articles)
	if err != nil {
		return nil, false
	}
	
	return articles, true
}

func (c *Cache) SetArticleList(ctx context.Context, sort, level, contentType string, page int, articles []model.ArticleCard) error {
	key := fmt.Sprintf(KeyArticleList, sort, level, contentType, page)
	return c.redis.SetJSON(ctx, key, articles, TTLShort)
}

func (c *Cache) GetArticle(ctx context.Context, id string) (*model.Article, bool) {
	key := fmt.Sprintf(KeyArticle, id)
	
	var article model.Article
	err := c.redis.GetJSON(ctx, key, &article)
	if err != nil {
		return nil, false
	}
	
	return &article, true
}

func (c *Cache) SetArticle(ctx context.Context, article *model.Article) error {
	key := fmt.Sprintf(KeyArticle, article.ID.String())
	return c.redis.SetJSON(ctx, key, article, TTLMedium)
}

func (c *Cache) InvalidateArticle(ctx context.Context, id string) error {
	key := fmt.Sprintf(KeyArticle, id)
	return c.redis.Del(ctx, key).Err()
}

func (c *Cache) InvalidateArticleLists(ctx context.Context) error {
	// Delete all article list caches using pattern
	keys, err := c.redis.Keys(ctx, "articles:list:*").Result()
	if err != nil {
		return err
	}
	
	if len(keys) > 0 {
		return c.redis.Del(ctx, keys...).Err()
	}
	
	return nil
}

// ============================================
// Categories cache
// ============================================

func (c *Cache) GetCategories(ctx context.Context) ([]model.Category, bool) {
	var categories []model.Category
	err := c.redis.GetJSON(ctx, KeyCategories, &categories)
	if err != nil {
		return nil, false
	}
	
	return categories, true
}

func (c *Cache) SetCategories(ctx context.Context, categories []model.Category) error {
	return c.redis.SetJSON(ctx, KeyCategories, categories, TTLLong)
}

func (c *Cache) GetCategory(ctx context.Context, slug string) (*model.Category, bool) {
	key := fmt.Sprintf(KeyCategory, slug)
	
	var category model.Category
	err := c.redis.GetJSON(ctx, key, &category)
	if err != nil {
		return nil, false
	}
	
	return &category, true
}

func (c *Cache) SetCategory(ctx context.Context, category *model.Category) error {
	key := fmt.Sprintf(KeyCategory, category.Slug)
	return c.redis.SetJSON(ctx, key, category, TTLLong)
}

func (c *Cache) InvalidateCategories(ctx context.Context) error {
	keys, err := c.redis.Keys(ctx, "categor*").Result()
	if err != nil {
		return err
	}
	
	if len(keys) > 0 {
		return c.redis.Del(ctx, keys...).Err()
	}
	
	return nil
}

// ============================================
// Tags cache
// ============================================

func (c *Cache) GetPopularTags(ctx context.Context, limit int) ([]model.Tag, bool) {
	key := fmt.Sprintf(KeyPopularTags, limit)
	
	var tags []model.Tag
	err := c.redis.GetJSON(ctx, key, &tags)
	if err != nil {
		return nil, false
	}
	
	return tags, true
}

func (c *Cache) SetPopularTags(ctx context.Context, limit int, tags []model.Tag) error {
	key := fmt.Sprintf(KeyPopularTags, limit)
	return c.redis.SetJSON(ctx, key, tags, TTLMedium)
}

// ============================================
// User cache
// ============================================

func (c *Cache) GetUserProfile(ctx context.Context, userID string) (*model.User, bool) {
	key := fmt.Sprintf(KeyUserProfile, userID)
	
	var user model.User
	err := c.redis.GetJSON(ctx, key, &user)
	if err != nil {
		return nil, false
	}
	
	return &user, true
}

func (c *Cache) SetUserProfile(ctx context.Context, user *model.User) error {
	key := fmt.Sprintf(KeyUserProfile, user.ID.String())
	return c.redis.SetJSON(ctx, key, user, TTLMedium)
}

func (c *Cache) InvalidateUserProfile(ctx context.Context, userID string) error {
	key := fmt.Sprintf(KeyUserProfile, userID)
	return c.redis.Del(ctx, key).Err()
}

// ============================================
// Notification count cache
// ============================================

func (c *Cache) GetNotificationCount(ctx context.Context, userID string) (int, bool) {
	key := fmt.Sprintf(KeyNotificationCnt, userID)
	
	count, err := c.redis.Get(ctx, key).Int()
	if err != nil {
		return 0, false
	}
	
	return count, true
}

func (c *Cache) SetNotificationCount(ctx context.Context, userID string, count int) error {
	key := fmt.Sprintf(KeyNotificationCnt, userID)
	return c.redis.Set(ctx, key, count, TTLShort).Err()
}

func (c *Cache) IncrementNotificationCount(ctx context.Context, userID string) error {
	key := fmt.Sprintf(KeyNotificationCnt, userID)
	return c.redis.Incr(ctx, key).Err()
}

func (c *Cache) InvalidateNotificationCount(ctx context.Context, userID string) error {
	key := fmt.Sprintf(KeyNotificationCnt, userID)
	return c.redis.Del(ctx, key).Err()
}

// ============================================
// Online users tracking
// ============================================

func (c *Cache) SetUserOnline(ctx context.Context, userID string) error {
	return c.redis.SAdd(ctx, "online:users", userID).Err()
}

func (c *Cache) SetUserOffline(ctx context.Context, userID string) error {
	return c.redis.SRem(ctx, "online:users", userID).Err()
}

func (c *Cache) GetOnlineUsers(ctx context.Context) ([]string, error) {
	return c.redis.SMembers(ctx, "online:users").Result()
}

func (c *Cache) GetOnlineCount(ctx context.Context) (int64, error) {
	return c.redis.SCard(ctx, "online:users").Result()
}

func (c *Cache) IsUserOnline(ctx context.Context, userID string) (bool, error) {
	return c.redis.SIsMember(ctx, "online:users", userID).Result()
}

// ============================================
// View counting with deduplication
// ============================================

func (c *Cache) RecordView(ctx context.Context, articleID, userIP string) (bool, error) {
	key := fmt.Sprintf("view:%s:%s", articleID, userIP)
	
	// Check if already viewed
	exists, err := c.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	
	if exists > 0 {
		return false, nil // Already viewed
	}
	
	// Mark as viewed for 24 hours
	if err := c.redis.Set(ctx, key, "1", TTLDay).Err(); err != nil {
		return false, err
	}
	
	return true, nil
}

// ============================================
// Rate limiting
// ============================================

func (c *Cache) CheckRateLimit(ctx context.Context, identifier, action string, limit int, window time.Duration) (bool, int, error) {
	key := fmt.Sprintf("ratelimit:%s:%s", action, identifier)
	
	count, err := c.redis.Incr(ctx, key).Result()
	if err != nil {
		return false, 0, err
	}
	
	if count == 1 {
		c.redis.Expire(ctx, key, window)
	}
	
	remaining := limit - int(count)
	if remaining < 0 {
		remaining = 0
	}
	
	return count <= int64(limit), remaining, nil
}

// ============================================
// Generic helpers
// ============================================

func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
	return c.redis.GetJSON(ctx, key, dest)
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return c.redis.SetJSON(ctx, key, value, ttl)
}

func (c *Cache) Delete(ctx context.Context, keys ...string) error {
	return c.redis.Del(ctx, keys...).Err()
}

func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.redis.Exists(ctx, key).Result()
	return result > 0, err
}

// CacheOrFetch tries to get from cache, otherwise calls fetch function and caches result
func CacheOrFetch[T any](
	ctx context.Context,
	cache *Cache,
	key string,
	ttl time.Duration,
	fetch func() (T, error),
) (T, error) {
	var result T
	
	// Try cache first
	err := cache.redis.GetJSON(ctx, key, &result)
	if err == nil {
		return result, nil
	}
	
	// Fetch from source
	result, err = fetch()
	if err != nil {
		return result, err
	}
	
	// Cache result
	_ = cache.redis.SetJSON(ctx, key, result, ttl)
	
	return result, nil
}


