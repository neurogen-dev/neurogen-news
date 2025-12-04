package service

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/model"
	"github.com/neurogen-news/backend/internal/repository"
)

type NotificationService interface {
	Create(ctx context.Context, notification *model.Notification) error
	GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) (*NotificationListResult, error)
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error)
	MarkAsRead(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	DeleteAll(ctx context.Context, userID uuid.UUID) error

	// Notification creation helpers
	NotifyNewComment(ctx context.Context, articleAuthorID, commentAuthorID, articleID, commentID uuid.UUID, articleTitle string) error
	NotifyCommentReply(ctx context.Context, parentAuthorID, replyAuthorID, articleID, commentID uuid.UUID) error
	NotifyNewFollower(ctx context.Context, userID, followerID uuid.UUID) error
	NotifyReaction(ctx context.Context, authorID, reactorID, articleID uuid.UUID, emoji string) error
}

type NotificationListResult struct {
	Items       []model.Notification `json:"items"`
	Total       int                  `json:"total"`
	UnreadCount int                  `json:"unreadCount"`
	HasMore     bool                 `json:"hasMore"`
}

type notificationService struct {
	notificationRepo repository.NotificationRepository
	userRepo         repository.UserRepository
	redis            *repository.RedisClient
	logger           *zap.Logger
}

func NewNotificationService(
	notificationRepo repository.NotificationRepository,
	userRepo repository.UserRepository,
	redis *repository.RedisClient,
	logger *zap.Logger,
) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
		redis:            redis,
		logger:           logger,
	}
}

func (s *notificationService) Create(ctx context.Context, notification *model.Notification) error {
	return s.notificationRepo.Create(ctx, notification)
}

func (s *notificationService) GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) (*NotificationListResult, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	notifications, total, err := s.notificationRepo.GetByUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	unreadCount, _ := s.notificationRepo.GetUnreadCount(ctx, userID)

	return &NotificationListResult{
		Items:       notifications,
		Total:       total,
		UnreadCount: unreadCount,
		HasMore:     offset+len(notifications) < total,
	}, nil
}

func (s *notificationService) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	return s.notificationRepo.GetUnreadCount(ctx, userID)
}

func (s *notificationService) MarkAsRead(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	// TODO: Verify notification belongs to user
	return s.notificationRepo.MarkAsRead(ctx, id)
}

func (s *notificationService) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	return s.notificationRepo.MarkAllAsRead(ctx, userID)
}

func (s *notificationService) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	// TODO: Verify notification belongs to user
	return s.notificationRepo.Delete(ctx, id)
}

func (s *notificationService) DeleteAll(ctx context.Context, userID uuid.UUID) error {
	return s.notificationRepo.DeleteAll(ctx, userID)
}

func (s *notificationService) NotifyNewComment(ctx context.Context, articleAuthorID, commentAuthorID, articleID, commentID uuid.UUID, articleTitle string) error {
	// Don't notify if user comments on their own article
	if articleAuthorID == commentAuthorID {
		return nil
	}

	notification := &model.Notification{
		UserID:    articleAuthorID,
		Type:      model.NotificationNewComment,
		Title:     "Новый комментарий",
		Message:   "Новый комментарий к вашей статье \"" + articleTitle + "\"",
		ActorID:   &commentAuthorID,
		ArticleID: &articleID,
		CommentID: &commentID,
	}

	link := "/article/" + articleID.String() + "#comment-" + commentID.String()
	notification.Link = &link

	return s.notificationRepo.Create(ctx, notification)
}

func (s *notificationService) NotifyCommentReply(ctx context.Context, parentAuthorID, replyAuthorID, articleID, commentID uuid.UUID) error {
	// Don't notify if user replies to their own comment
	if parentAuthorID == replyAuthorID {
		return nil
	}

	notification := &model.Notification{
		UserID:    parentAuthorID,
		Type:      model.NotificationCommentReply,
		Title:     "Ответ на комментарий",
		Message:   "Кто-то ответил на ваш комментарий",
		ActorID:   &replyAuthorID,
		ArticleID: &articleID,
		CommentID: &commentID,
	}

	link := "/article/" + articleID.String() + "#comment-" + commentID.String()
	notification.Link = &link

	return s.notificationRepo.Create(ctx, notification)
}

func (s *notificationService) NotifyNewFollower(ctx context.Context, userID, followerID uuid.UUID) error {
	notification := &model.Notification{
		UserID:  userID,
		Type:    model.NotificationNewFollower,
		Title:   "Новый подписчик",
		Message: "У вас новый подписчик",
		ActorID: &followerID,
	}

	link := "/user/" + followerID.String()
	notification.Link = &link

	return s.notificationRepo.Create(ctx, notification)
}

func (s *notificationService) NotifyReaction(ctx context.Context, authorID, reactorID, articleID uuid.UUID, emoji string) error {
	// Don't notify if user reacts to their own article
	if authorID == reactorID {
		return nil
	}

	notification := &model.Notification{
		UserID:    authorID,
		Type:      model.NotificationReaction,
		Title:     "Новая реакция",
		Message:   "Кто-то поставил " + emoji + " вашей статье",
		ActorID:   &reactorID,
		ArticleID: &articleID,
	}

	link := "/article/" + articleID.String()
	notification.Link = &link

	return s.notificationRepo.Create(ctx, notification)
}

