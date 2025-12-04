package handler

import (
	"github.com/neurogen-news/backend/internal/service"
	"go.uber.org/zap"
)

type Handlers struct {
	Auth         *AuthHandler
	User         *UserHandler
	Article      *ArticleHandler
	Comment      *CommentHandler
	Category     *CategoryHandler
	Tag          *TagHandler
	Search       *SearchHandler
	Notification *NotificationHandler
	Achievement  *AchievementHandler
	Bookmark     *BookmarkHandler
	Draft        *DraftHandler
	Upload       *UploadHandler
	Admin        *AdminHandler
}

func NewHandlers(services *service.Services, logger *zap.Logger) *Handlers {
	return &Handlers{
		Auth:         NewAuthHandler(services.Auth, logger),
		User:         NewUserHandler(services.User, logger),
		Article:      NewArticleHandler(services.Article, logger),
		Comment:      NewCommentHandler(services.Comment, logger),
		Category:     NewCategoryHandler(services.Category, services.Article, logger),
		Tag:          NewTagHandler(services.Tag, services.Article, logger),
		Search:       NewSearchHandler(services.Search, logger),
		Notification: NewNotificationHandler(services.Notification, logger),
		Achievement:  NewAchievementHandler(services.Achievement, logger),
		Bookmark:     NewBookmarkHandler(services.Bookmark, logger),
		Draft:        NewDraftHandler(services.Draft, logger),
		Upload:       NewUploadHandler(services.Upload, logger),
		Admin:        NewAdminHandler(services, logger),
	}
}

