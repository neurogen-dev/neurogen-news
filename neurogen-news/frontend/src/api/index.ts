// Re-export all API modules
export { api, apiClient, uploadFile } from './client'
export { articlesApi } from './articles'
export { commentsApi } from './comments'
export { categoriesApi } from './categories'
export { tagsApi } from './tags'
export { usersApi } from './users'
export { notificationsApi } from './notifications'
export { bookmarksApi } from './bookmarks'
export { draftsApi } from './drafts'
export { achievementsApi } from './achievements'
export { searchApi } from './search'
export { uploadApi } from './upload'

// Re-export types
export type { ArticleListParams, CreateArticleInput, UpdateArticleInput } from './articles'
export type { CommentListParams, CreateCommentInput } from './comments'
export type { CategoryArticleParams } from './categories'
export type { UserProfile, UpdateProfileInput } from './users'
export type { NotificationListResult } from './notifications'
export type { BookmarkFolder, BookmarkListParams } from './bookmarks'
export type { CreateDraftInput, UpdateDraftInput, AutoSaveInput } from './drafts'
export type { UserAchievement, AchievementProgress } from './achievements'
export type { SearchParams, SearchResult, SearchSuggestion } from './search'
export type { UploadResult, UploadType } from './upload'


