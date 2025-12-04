package repository

type Repositories struct {
	User         UserRepository
	Article      ArticleRepository
	Comment      CommentRepository
	Category     CategoryRepository
	Tag          TagRepository
	Notification NotificationRepository
	Achievement  AchievementRepository
	Bookmark     BookmarkRepository
	Draft        DraftRepository
	Reaction     ReactionRepository
}

func NewRepositories(db *PostgresDB) *Repositories {
	return &Repositories{
		User:         NewUserRepository(db),
		Article:      NewArticleRepository(db),
		Comment:      NewCommentRepository(db),
		Category:     NewCategoryRepository(db),
		Tag:          NewTagRepository(db),
		Notification: NewNotificationRepository(db),
		Achievement:  NewAchievementRepository(db),
		Bookmark:     NewBookmarkRepository(db),
		Draft:        NewDraftRepository(db),
		Reaction:     NewReactionRepository(db),
	}
}

