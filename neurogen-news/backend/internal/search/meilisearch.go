package search

import (
	"context"
	"fmt"

	"github.com/meilisearch/meilisearch-go"
	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/model"
)

// Index names
const (
	IndexArticles = "articles"
	IndexUsers    = "users"
	IndexTags     = "tags"
)

// SearchableArticle represents article document for search
type SearchableArticle struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Slug          string   `json:"slug"`
	Lead          string   `json:"lead"`
	Content       string   `json:"content"`
	Level         string   `json:"level"`
	ContentType   string   `json:"contentType"`
	AuthorID      string   `json:"authorId"`
	AuthorName    string   `json:"authorName"`
	CategoryID    string   `json:"categoryId"`
	CategoryName  string   `json:"categoryName"`
	CategorySlug  string   `json:"categorySlug"`
	Tags          []string `json:"tags"`
	ViewCount     int      `json:"viewCount"`
	CommentCount  int      `json:"commentCount"`
	PublishedAt   int64    `json:"publishedAt"` // Unix timestamp for sorting
	CoverImageURL string   `json:"coverImageUrl,omitempty"`
	ReadingTime   int      `json:"readingTime"`
}

// SearchableUser represents user document for search
type SearchableUser struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Bio         string `json:"bio"`
	AvatarURL   string `json:"avatarUrl,omitempty"`
	IsVerified  bool   `json:"isVerified"`
	Karma       int    `json:"karma"`
}

// SearchableTag represents tag document for search
type SearchableTag struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	ArticleCount int    `json:"articleCount"`
}

// Client wraps Meilisearch client
type Client struct {
	client *meilisearch.Client
	logger *zap.Logger
}

// NewClient creates new Meilisearch client
func NewClient(host, apiKey string, logger *zap.Logger) (*Client, error) {
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   host,
		APIKey: apiKey,
	})

	// Check connection
	if _, err := client.Health(); err != nil {
		return nil, fmt.Errorf("failed to connect to Meilisearch: %w", err)
	}

	c := &Client{
		client: client,
		logger: logger,
	}

	// Setup indexes
	if err := c.setupIndexes(); err != nil {
		return nil, err
	}

	return c, nil
}

// setupIndexes creates and configures search indexes
func (c *Client) setupIndexes() error {
	// Articles index
	articlesIndex := c.client.Index(IndexArticles)

	// Set filterable attributes
	_, err := articlesIndex.UpdateFilterableAttributes(&[]string{
		"level",
		"contentType",
		"categoryId",
		"categorySlug",
		"authorId",
		"tags",
		"publishedAt",
	})
	if err != nil {
		return fmt.Errorf("failed to set filterable attributes for articles: %w", err)
	}

	// Set sortable attributes
	_, err = articlesIndex.UpdateSortableAttributes(&[]string{
		"publishedAt",
		"viewCount",
		"commentCount",
	})
	if err != nil {
		return fmt.Errorf("failed to set sortable attributes for articles: %w", err)
	}

	// Set searchable attributes with priority
	_, err = articlesIndex.UpdateSearchableAttributes(&[]string{
		"title",
		"lead",
		"content",
		"authorName",
		"categoryName",
		"tags",
	})
	if err != nil {
		return fmt.Errorf("failed to set searchable attributes for articles: %w", err)
	}

	// Set ranking rules
	_, err = articlesIndex.UpdateRankingRules(&[]string{
		"words",
		"typo",
		"proximity",
		"attribute",
		"sort",
		"exactness",
		"viewCount:desc",
	})
	if err != nil {
		return fmt.Errorf("failed to set ranking rules for articles: %w", err)
	}

	// Users index
	usersIndex := c.client.Index(IndexUsers)

	_, err = usersIndex.UpdateFilterableAttributes(&[]string{
		"isVerified",
	})
	if err != nil {
		return fmt.Errorf("failed to set filterable attributes for users: %w", err)
	}

	_, err = usersIndex.UpdateSortableAttributes(&[]string{
		"karma",
	})
	if err != nil {
		return fmt.Errorf("failed to set sortable attributes for users: %w", err)
	}

	_, err = usersIndex.UpdateSearchableAttributes(&[]string{
		"username",
		"displayName",
		"bio",
	})
	if err != nil {
		return fmt.Errorf("failed to set searchable attributes for users: %w", err)
	}

	// Tags index
	tagsIndex := c.client.Index(IndexTags)

	_, err = tagsIndex.UpdateSortableAttributes(&[]string{
		"articleCount",
	})
	if err != nil {
		return fmt.Errorf("failed to set sortable attributes for tags: %w", err)
	}

	_, err = tagsIndex.UpdateSearchableAttributes(&[]string{
		"name",
	})
	if err != nil {
		return fmt.Errorf("failed to set searchable attributes for tags: %w", err)
	}

	c.logger.Info("Meilisearch indexes configured successfully")
	return nil
}

// ============================================
// Articles
// ============================================

// IndexArticle adds or updates article in search index
func (c *Client) IndexArticle(ctx context.Context, article *SearchableArticle) error {
	_, err := c.client.Index(IndexArticles).AddDocuments([]SearchableArticle{*article}, "id")
	if err != nil {
		c.logger.Error("Failed to index article", zap.String("id", article.ID), zap.Error(err))
		return err
	}
	return nil
}

// IndexArticles adds multiple articles to search index
func (c *Client) IndexArticles(ctx context.Context, articles []SearchableArticle) error {
	if len(articles) == 0 {
		return nil
	}
	_, err := c.client.Index(IndexArticles).AddDocuments(articles, "id")
	if err != nil {
		c.logger.Error("Failed to index articles", zap.Error(err))
		return err
	}
	return nil
}

// DeleteArticle removes article from search index
func (c *Client) DeleteArticle(ctx context.Context, id string) error {
	_, err := c.client.Index(IndexArticles).DeleteDocument(id)
	if err != nil {
		c.logger.Error("Failed to delete article from index", zap.String("id", id), zap.Error(err))
		return err
	}
	return nil
}

// SearchArticleParams represents article search parameters
type SearchArticleParams struct {
	Query       string
	Level       string
	ContentType string
	CategoryID  string
	Tags        []string
	Sort        string // "relevance", "new", "popular"
	Limit       int64
	Offset      int64
}

// SearchArticleResult represents search result
type SearchArticleResult struct {
	Hits             []SearchableArticle `json:"hits"`
	Total            int64               `json:"total"`
	ProcessingTimeMs int64               `json:"processingTimeMs"`
}

// SearchArticles searches for articles
func (c *Client) SearchArticles(ctx context.Context, params SearchArticleParams) (*SearchArticleResult, error) {
	// Build filter
	var filters []string

	if params.Level != "" {
		filters = append(filters, fmt.Sprintf("level = '%s'", params.Level))
	}
	if params.ContentType != "" {
		filters = append(filters, fmt.Sprintf("contentType = '%s'", params.ContentType))
	}
	if params.CategoryID != "" {
		filters = append(filters, fmt.Sprintf("categoryId = '%s'", params.CategoryID))
	}
	for _, tag := range params.Tags {
		filters = append(filters, fmt.Sprintf("tags = '%s'", tag))
	}

	// Build sort
	var sort []string
	switch params.Sort {
	case "new":
		sort = []string{"publishedAt:desc"}
	case "popular":
		sort = []string{"viewCount:desc"}
	default:
		// relevance - no explicit sort
	}

	// Set defaults
	if params.Limit <= 0 || params.Limit > 100 {
		params.Limit = 20
	}

	searchReq := &meilisearch.SearchRequest{
		Limit:  params.Limit,
		Offset: params.Offset,
	}

	if len(filters) > 0 {
		filterStr := filters[0]
		for i := 1; i < len(filters); i++ {
			filterStr += " AND " + filters[i]
		}
		searchReq.Filter = filterStr
	}

	if len(sort) > 0 {
		searchReq.Sort = sort
	}

	resp, err := c.client.Index(IndexArticles).Search(params.Query, searchReq)
	if err != nil {
		c.logger.Error("Article search failed", zap.Error(err))
		return nil, err
	}

	// Convert hits
	var hits []SearchableArticle
	for _, hit := range resp.Hits {
		if hitMap, ok := hit.(map[string]interface{}); ok {
			article := mapToSearchableArticle(hitMap)
			hits = append(hits, article)
		}
	}

	return &SearchArticleResult{
		Hits:             hits,
		Total:            resp.EstimatedTotalHits,
		ProcessingTimeMs: resp.ProcessingTimeMs,
	}, nil
}

// ============================================
// Users
// ============================================

// IndexUser adds or updates user in search index
func (c *Client) IndexUser(ctx context.Context, user *SearchableUser) error {
	_, err := c.client.Index(IndexUsers).AddDocuments([]SearchableUser{*user}, "id")
	if err != nil {
		c.logger.Error("Failed to index user", zap.String("id", user.ID), zap.Error(err))
		return err
	}
	return nil
}

// DeleteUser removes user from search index
func (c *Client) DeleteUser(ctx context.Context, id string) error {
	_, err := c.client.Index(IndexUsers).DeleteDocument(id)
	return err
}

// SearchUserResult represents user search result
type SearchUserResult struct {
	Hits             []SearchableUser `json:"hits"`
	Total            int64            `json:"total"`
	ProcessingTimeMs int64            `json:"processingTimeMs"`
}

// SearchUsers searches for users
func (c *Client) SearchUsers(ctx context.Context, query string, limit, offset int64) (*SearchUserResult, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	resp, err := c.client.Index(IndexUsers).Search(query, &meilisearch.SearchRequest{
		Limit:  limit,
		Offset: offset,
		Sort:   []string{"karma:desc"},
	})
	if err != nil {
		c.logger.Error("User search failed", zap.Error(err))
		return nil, err
	}

	var hits []SearchableUser
	for _, hit := range resp.Hits {
		if hitMap, ok := hit.(map[string]interface{}); ok {
			user := mapToSearchableUser(hitMap)
			hits = append(hits, user)
		}
	}

	return &SearchUserResult{
		Hits:             hits,
		Total:            resp.EstimatedTotalHits,
		ProcessingTimeMs: resp.ProcessingTimeMs,
	}, nil
}

// ============================================
// Tags
// ============================================

// IndexTag adds or updates tag in search index
func (c *Client) IndexTag(ctx context.Context, tag *SearchableTag) error {
	_, err := c.client.Index(IndexTags).AddDocuments([]SearchableTag{*tag}, "id")
	return err
}

// SearchTags searches for tags (autocomplete)
func (c *Client) SearchTags(ctx context.Context, query string, limit int64) ([]SearchableTag, error) {
	if limit <= 0 || limit > 20 {
		limit = 10
	}

	resp, err := c.client.Index(IndexTags).Search(query, &meilisearch.SearchRequest{
		Limit: limit,
		Sort:  []string{"articleCount:desc"},
	})
	if err != nil {
		return nil, err
	}

	var tags []SearchableTag
	for _, hit := range resp.Hits {
		if hitMap, ok := hit.(map[string]interface{}); ok {
			tag := mapToSearchableTag(hitMap)
			tags = append(tags, tag)
		}
	}

	return tags, nil
}

// ============================================
// Helpers
// ============================================

func mapToSearchableArticle(m map[string]interface{}) SearchableArticle {
	article := SearchableArticle{}

	if v, ok := m["id"].(string); ok {
		article.ID = v
	}
	if v, ok := m["title"].(string); ok {
		article.Title = v
	}
	if v, ok := m["slug"].(string); ok {
		article.Slug = v
	}
	if v, ok := m["lead"].(string); ok {
		article.Lead = v
	}
	if v, ok := m["level"].(string); ok {
		article.Level = v
	}
	if v, ok := m["contentType"].(string); ok {
		article.ContentType = v
	}
	if v, ok := m["authorId"].(string); ok {
		article.AuthorID = v
	}
	if v, ok := m["authorName"].(string); ok {
		article.AuthorName = v
	}
	if v, ok := m["categoryId"].(string); ok {
		article.CategoryID = v
	}
	if v, ok := m["categoryName"].(string); ok {
		article.CategoryName = v
	}
	if v, ok := m["categorySlug"].(string); ok {
		article.CategorySlug = v
	}
	if v, ok := m["coverImageUrl"].(string); ok {
		article.CoverImageURL = v
	}
	if v, ok := m["viewCount"].(float64); ok {
		article.ViewCount = int(v)
	}
	if v, ok := m["commentCount"].(float64); ok {
		article.CommentCount = int(v)
	}
	if v, ok := m["readingTime"].(float64); ok {
		article.ReadingTime = int(v)
	}
	if v, ok := m["publishedAt"].(float64); ok {
		article.PublishedAt = int64(v)
	}
	if v, ok := m["tags"].([]interface{}); ok {
		for _, t := range v {
			if tag, ok := t.(string); ok {
				article.Tags = append(article.Tags, tag)
			}
		}
	}

	return article
}

func mapToSearchableUser(m map[string]interface{}) SearchableUser {
	user := SearchableUser{}

	if v, ok := m["id"].(string); ok {
		user.ID = v
	}
	if v, ok := m["username"].(string); ok {
		user.Username = v
	}
	if v, ok := m["displayName"].(string); ok {
		user.DisplayName = v
	}
	if v, ok := m["bio"].(string); ok {
		user.Bio = v
	}
	if v, ok := m["avatarUrl"].(string); ok {
		user.AvatarURL = v
	}
	if v, ok := m["isVerified"].(bool); ok {
		user.IsVerified = v
	}
	if v, ok := m["karma"].(float64); ok {
		user.Karma = int(v)
	}

	return user
}

func mapToSearchableTag(m map[string]interface{}) SearchableTag {
	tag := SearchableTag{}

	if v, ok := m["id"].(string); ok {
		tag.ID = v
	}
	if v, ok := m["name"].(string); ok {
		tag.Name = v
	}
	if v, ok := m["slug"].(string); ok {
		tag.Slug = v
	}
	if v, ok := m["articleCount"].(float64); ok {
		tag.ArticleCount = int(v)
	}

	return tag
}

// ArticleToSearchable converts model.Article to SearchableArticle
func ArticleToSearchable(article *model.Article, authorName, categoryName, categorySlug string, tags []string) *SearchableArticle {
	var publishedAt int64
	if article.PublishedAt != nil {
		publishedAt = article.PublishedAt.Unix()
	}

	var lead string
	if article.Lead != nil {
		lead = *article.Lead
	}

	var coverURL string
	if article.CoverImageURL != nil {
		coverURL = *article.CoverImageURL
	}

	return &SearchableArticle{
		ID:            article.ID.String(),
		Title:         article.Title,
		Slug:          article.Slug,
		Lead:          lead,
		Content:       article.Content,
		Level:         string(article.Level),
		ContentType:   string(article.ContentType),
		AuthorID:      article.AuthorID.String(),
		AuthorName:    authorName,
		CategoryID:    article.CategoryID.String(),
		CategoryName:  categoryName,
		CategorySlug:  categorySlug,
		Tags:          tags,
		ViewCount:     article.ViewCount,
		CommentCount:  article.CommentCount,
		PublishedAt:   publishedAt,
		CoverImageURL: coverURL,
		ReadingTime:   article.ReadingTime,
	}
}

// UserToSearchable converts model.User to SearchableUser
func UserToSearchable(user *model.User) *SearchableUser {
	var bio, avatarURL string
	if user.Bio != nil {
		bio = *user.Bio
	}
	if user.AvatarURL != nil {
		avatarURL = *user.AvatarURL
	}

	return &SearchableUser{
		ID:          user.ID.String(),
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Bio:         bio,
		AvatarURL:   avatarURL,
		IsVerified:  user.IsVerified,
		Karma:       user.Karma,
	}
}

// TagToSearchable converts model.Tag to SearchableTag
func TagToSearchable(tag *model.Tag) *SearchableTag {
	return &SearchableTag{
		ID:           tag.ID.String(),
		Name:         tag.Name,
		Slug:         tag.Slug,
		ArticleCount: tag.ArticleCount,
	}
}


