package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/neurogen-news/backend/internal/model"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *model.Notification) error
	GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Notification, int, error)
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error)
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteAll(ctx context.Context, userID uuid.UUID) error
}

type notificationRepository struct {
	db *PostgresDB
}

func NewNotificationRepository(db *PostgresDB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, notification *model.Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, type, title, message, link, actor_id, article_id, comment_id, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, false, NOW())
	`

	notification.ID = uuid.New()

	_, err := r.db.Exec(ctx, query,
		notification.ID,
		notification.UserID,
		notification.Type,
		notification.Title,
		notification.Message,
		notification.Link,
		notification.ActorID,
		notification.ArticleID,
		notification.CommentID,
	)

	return err
}

func (r *notificationRepository) GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Notification, int, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	// Count total
	countQuery := `SELECT COUNT(*) FROM notifications WHERE user_id = $1`
	var total int
	if err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get notifications with actor info
	query := `
		SELECT 
			n.id, n.user_id, n.type, n.title, n.message, n.link,
			n.actor_id, n.article_id, n.comment_id, n.is_read, n.created_at,
			u.username, u.display_name, u.avatar_url
		FROM notifications n
		LEFT JOIN users u ON u.id = n.actor_id
		WHERE n.user_id = $1
		ORDER BY n.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var notifications []model.Notification
	for rows.Next() {
		var n model.Notification
		var actorUsername, actorDisplayName, actorAvatarURL *string

		err := rows.Scan(
			&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.Link,
			&n.ActorID, &n.ArticleID, &n.CommentID, &n.IsRead, &n.CreatedAt,
			&actorUsername, &actorDisplayName, &actorAvatarURL,
		)
		if err != nil {
			return nil, 0, err
		}

		if n.ActorID != nil {
			n.Actor = &model.NotificationActor{
				ID:          *n.ActorID,
				Username:    *actorUsername,
				DisplayName: *actorDisplayName,
				AvatarURL:   actorAvatarURL,
			}
		}

		notifications = append(notifications, n)
	}

	return notifications, total, nil
}

func (r *notificationRepository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false`
	var count int
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE notifications SET is_read = true WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE notifications SET is_read = true WHERE user_id = $1 AND is_read = false`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *notificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM notifications WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *notificationRepository) DeleteAll(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM notifications WHERE user_id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}
