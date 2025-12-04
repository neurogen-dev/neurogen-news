package service

import (
	"github.com/neurogen-news/backend/internal/repository"
	"go.uber.org/zap"
)

type Services struct {
	Auth         AuthService
	User         UserService
	Article      ArticleService
	Comment      CommentService
	Category     CategoryService
	Tag          TagService
	Notification NotificationService
	Achievement  AchievementService
	Bookmark     BookmarkService
	Draft        DraftService
	Search       SearchService
	Upload       UploadService
}

type Deps struct {
	Repos     *repository.Repositories
	Redis     *repository.RedisClient
	JWTSecret string
	Logger    *zap.Logger
}

func NewServices(deps Deps) *Services {
	return &Services{
		Auth:         NewAuthService(deps.Repos.User, deps.Redis, deps.JWTSecret, deps.Logger),
		User:         NewUserService(deps.Repos.User, deps.Repos.Article, deps.Redis, deps.Logger),
		Article:      NewArticleService(deps.Repos.Article, deps.Repos.Tag, deps.Redis, deps.Logger),
		Comment:      NewCommentService(deps.Repos.Comment, deps.Repos.Notification, deps.Repos.Reaction, deps.Redis, deps.Logger),
		Category:     NewCategoryService(deps.Repos.Category, deps.Redis, deps.Logger),
		Tag:          NewTagService(deps.Repos.Tag, deps.Redis, deps.Logger),
		Notification: NewNotificationService(deps.Repos.Notification, deps.Repos.User, deps.Redis, deps.Logger),
		Achievement:  NewAchievementService(deps.Repos.Achievement, deps.Logger),
		Bookmark:     NewBookmarkService(deps.Repos.Bookmark, deps.Logger),
		Draft:        NewDraftService(deps.Repos.Draft, deps.Logger),
		Search:       NewSearchService(deps.Repos.Article, deps.Repos.User, deps.Logger),
		Upload:       NewUploadService(deps.Logger),
	}
}

