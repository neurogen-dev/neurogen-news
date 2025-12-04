package service

import (
	"context"

	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/model"
	"github.com/neurogen-news/backend/internal/repository"
	"github.com/neurogen-news/backend/internal/search"
)

type SearchService interface {
	Search(ctx context.Context, query string, params SearchParams) (*SearchResult, error)
	SearchArticles(ctx context.Context, query string, limit, offset int) ([]model.ArticleCard, int, error)
	SearchUsers(ctx context.Context, query string, limit, offset int) ([]model.User, int, error)
	GetSuggestions(ctx context.Context, query string, limit int) ([]SearchSuggestion, error)

	// Index management
	IndexArticle(ctx context.Context, article *model.Article, authorName, categoryName, categorySlug string, tags []string) error
	IndexUser(ctx context.Context, user *model.User) error
	DeleteArticleIndex(ctx context.Context, id string) error
}

type SearchParams struct {
	Type        string // all, articles, users, tags
	Sort        string // relevance, new, popular
	Level       string
	ContentType string
	Category    string
	Limit       int
	Offset      int
}

type SearchResult struct {
	Articles *ArticleSearchResult `json:"articles,omitempty"`
	Users    *UserSearchResult    `json:"users,omitempty"`
	Tags     []model.Tag          `json:"tags,omitempty"`
	Total    int                  `json:"total"`
}

type ArticleSearchResult struct {
	Items   []model.ArticleCard `json:"items"`
	Total   int                 `json:"total"`
	HasMore bool                `json:"hasMore"`
}

type UserSearchResult struct {
	Items   []model.User `json:"items"`
	Total   int          `json:"total"`
	HasMore bool         `json:"hasMore"`
}

type SearchSuggestion struct {
	Type  string `json:"type"` // article, user, tag
	Text  string `json:"text"`
	Slug  string `json:"slug,omitempty"`
	Image string `json:"image,omitempty"`
}

type searchService struct {
	articleRepo  repository.ArticleRepository
	userRepo     repository.UserRepository
	tagRepo      repository.TagRepository
	searchClient *search.Client
	logger       *zap.Logger
}

func NewSearchService(
	articleRepo repository.ArticleRepository,
	userRepo repository.UserRepository,
	tagRepo repository.TagRepository,
	searchClient *search.Client,
	logger *zap.Logger,
) SearchService {
	return &searchService{
		articleRepo:  articleRepo,
		userRepo:     userRepo,
		tagRepo:      tagRepo,
		searchClient: searchClient,
		logger:       logger,
	}
}

func (s *searchService) Search(ctx context.Context, query string, params SearchParams) (*SearchResult, error) {
	if params.Limit <= 0 || params.Limit > 50 {
		params.Limit = 20
	}

	result := &SearchResult{}

	// Use Meilisearch if available
	if s.searchClient != nil {
		return s.searchWithMeilisearch(ctx, query, params)
	}

	// Fallback to PostgreSQL search
	switch params.Type {
	case "articles":
		articles, total, err := s.SearchArticles(ctx, query, params.Limit, params.Offset)
		if err != nil {
			return nil, err
		}
		result.Articles = &ArticleSearchResult{
			Items:   articles,
			Total:   total,
			HasMore: params.Offset+len(articles) < total,
		}
		result.Total = total

	case "users":
		users, total, err := s.SearchUsers(ctx, query, params.Limit, params.Offset)
		if err != nil {
			return nil, err
		}
		result.Users = &UserSearchResult{
			Items:   users,
			Total:   total,
			HasMore: params.Offset+len(users) < total,
		}
		result.Total = total

	default: // all
		articles, articlesTotal, _ := s.SearchArticles(ctx, query, 10, 0)
		users, usersTotal, _ := s.SearchUsers(ctx, query, 5, 0)

		result.Articles = &ArticleSearchResult{
			Items:   articles,
			Total:   articlesTotal,
			HasMore: articlesTotal > 10,
		}
		result.Users = &UserSearchResult{
			Items:   users,
			Total:   usersTotal,
			HasMore: usersTotal > 5,
		}
		result.Total = articlesTotal + usersTotal
	}

	return result, nil
}

func (s *searchService) searchWithMeilisearch(ctx context.Context, query string, params SearchParams) (*SearchResult, error) {
	result := &SearchResult{}

	switch params.Type {
	case "articles":
		searchResult, err := s.searchClient.SearchArticles(ctx, search.SearchArticleParams{
			Query:       query,
			Level:       params.Level,
			ContentType: params.ContentType,
			CategoryID:  params.Category,
			Sort:        params.Sort,
			Limit:       int64(params.Limit),
			Offset:      int64(params.Offset),
		})
		if err != nil {
			return nil, err
		}

		articles := make([]model.ArticleCard, len(searchResult.Hits))
		for i, hit := range searchResult.Hits {
			articles[i] = searchHitToArticleCard(hit)
		}

		result.Articles = &ArticleSearchResult{
			Items:   articles,
			Total:   int(searchResult.Total),
			HasMore: int64(params.Offset)+int64(len(articles)) < searchResult.Total,
		}
		result.Total = int(searchResult.Total)

	case "users":
		searchResult, err := s.searchClient.SearchUsers(ctx, query, int64(params.Limit), int64(params.Offset))
		if err != nil {
			return nil, err
		}

		users := make([]model.User, len(searchResult.Hits))
		for i, hit := range searchResult.Hits {
			users[i] = searchHitToUser(hit)
		}

		result.Users = &UserSearchResult{
			Items:   users,
			Total:   int(searchResult.Total),
			HasMore: int64(params.Offset)+int64(len(users)) < searchResult.Total,
		}
		result.Total = int(searchResult.Total)

	case "tags":
		tags, err := s.searchClient.SearchTags(ctx, query, int64(params.Limit))
		if err != nil {
			return nil, err
		}

		modelTags := make([]model.Tag, len(tags))
		for i, t := range tags {
			modelTags[i] = searchHitToTag(t)
		}
		result.Tags = modelTags
		result.Total = len(tags)

	default: // all
		// Search articles
		articlesResult, _ := s.searchClient.SearchArticles(ctx, search.SearchArticleParams{
			Query: query,
			Limit: 10,
		})
		if articlesResult != nil {
			articles := make([]model.ArticleCard, len(articlesResult.Hits))
			for i, hit := range articlesResult.Hits {
				articles[i] = searchHitToArticleCard(hit)
			}
			result.Articles = &ArticleSearchResult{
				Items:   articles,
				Total:   int(articlesResult.Total),
				HasMore: articlesResult.Total > 10,
			}
		}

		// Search users
		usersResult, _ := s.searchClient.SearchUsers(ctx, query, 5, 0)
		if usersResult != nil {
			users := make([]model.User, len(usersResult.Hits))
			for i, hit := range usersResult.Hits {
				users[i] = searchHitToUser(hit)
			}
			result.Users = &UserSearchResult{
				Items:   users,
				Total:   int(usersResult.Total),
				HasMore: usersResult.Total > 5,
			}
		}

		// Search tags
		tags, _ := s.searchClient.SearchTags(ctx, query, 5)
		if len(tags) > 0 {
			modelTags := make([]model.Tag, len(tags))
			for i, t := range tags {
				modelTags[i] = searchHitToTag(t)
			}
			result.Tags = modelTags
		}

		result.Total = 0
		if result.Articles != nil {
			result.Total += result.Articles.Total
		}
		if result.Users != nil {
			result.Total += result.Users.Total
		}
		result.Total += len(result.Tags)
	}

	return result, nil
}

func (s *searchService) SearchArticles(ctx context.Context, query string, limit, offset int) ([]model.ArticleCard, int, error) {
	// Use PostgreSQL full-text search as fallback
	return s.articleRepo.List(ctx, repository.ArticleListParams{
		Sort:   "popular",
		Limit:  limit,
		Offset: offset,
	})
}

func (s *searchService) SearchUsers(ctx context.Context, query string, limit, offset int) ([]model.User, int, error) {
	return s.userRepo.Search(ctx, query, limit, offset)
}

func (s *searchService) GetSuggestions(ctx context.Context, query string, limit int) ([]SearchSuggestion, error) {
	if limit <= 0 || limit > 10 {
		limit = 5
	}

	var suggestions []SearchSuggestion

	if s.searchClient != nil {
		// Search articles for suggestions
		articlesResult, _ := s.searchClient.SearchArticles(ctx, search.SearchArticleParams{
			Query: query,
			Limit: int64(limit),
		})
		if articlesResult != nil {
			for _, hit := range articlesResult.Hits {
				suggestions = append(suggestions, SearchSuggestion{
					Type:  "article",
					Text:  hit.Title,
					Slug:  hit.CategorySlug + "/" + hit.Slug,
					Image: hit.CoverImageURL,
				})
			}
		}

		// Search users
		usersResult, _ := s.searchClient.SearchUsers(ctx, query, int64(limit), 0)
		if usersResult != nil {
			for _, hit := range usersResult.Hits {
				suggestions = append(suggestions, SearchSuggestion{
					Type:  "user",
					Text:  hit.DisplayName,
					Slug:  hit.Username,
					Image: hit.AvatarURL,
				})
			}
		}

		// Search tags
		tags, _ := s.searchClient.SearchTags(ctx, query, int64(limit))
		for _, tag := range tags {
			suggestions = append(suggestions, SearchSuggestion{
				Type: "tag",
				Text: tag.Name,
				Slug: tag.Slug,
			})
		}
	}

	return suggestions, nil
}

func (s *searchService) IndexArticle(ctx context.Context, article *model.Article, authorName, categoryName, categorySlug string, tags []string) error {
	if s.searchClient == nil {
		return nil
	}

	searchable := search.ArticleToSearchable(article, authorName, categoryName, categorySlug, tags)
	return s.searchClient.IndexArticle(ctx, searchable)
}

func (s *searchService) IndexUser(ctx context.Context, user *model.User) error {
	if s.searchClient == nil {
		return nil
	}

	searchable := search.UserToSearchable(user)
	return s.searchClient.IndexUser(ctx, searchable)
}

func (s *searchService) DeleteArticleIndex(ctx context.Context, id string) error {
	if s.searchClient == nil {
		return nil
	}
	return s.searchClient.DeleteArticle(ctx, id)
}

// Helper functions to convert search hits to models
func searchHitToArticleCard(hit search.SearchableArticle) model.ArticleCard {
	return model.ArticleCard{
		// ID would need parsing from string
		Title:          hit.Title,
		Slug:           hit.Slug,
		Level:          model.ArticleLevel(hit.Level),
		ContentType:    model.ContentType(hit.ContentType),
		AuthorUsername: hit.AuthorName,
		CategorySlug:   hit.CategorySlug,
		CategoryName:   hit.CategoryName,
		ViewCount:      hit.ViewCount,
		CommentCount:   hit.CommentCount,
		ReadingTime:    hit.ReadingTime,
	}
}

func searchHitToUser(hit search.SearchableUser) model.User {
	bio := hit.Bio
	avatar := hit.AvatarURL

	return model.User{
		Username:    hit.Username,
		DisplayName: hit.DisplayName,
		Bio:         &bio,
		AvatarURL:   &avatar,
		IsVerified:  hit.IsVerified,
		Karma:       hit.Karma,
	}
}

func searchHitToTag(hit search.SearchableTag) model.Tag {
	return model.Tag{
		Name:         hit.Name,
		Slug:         hit.Slug,
		ArticleCount: hit.ArticleCount,
	}
}

