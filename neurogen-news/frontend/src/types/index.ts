// User types
export type UserRole = 'USER' | 'AUTHOR' | 'MODERATOR' | 'EDITOR' | 'ADMIN'

export interface User {
  id: string
  username: string
  displayName: string
  email?: string
  avatarUrl?: string
  bio?: string
  role: UserRole
  karma: number
  isVerified: boolean
  isPremium: boolean
  createdAt: string
  
  // Stats
  followerCount: number
  followingCount: number
  articleCount: number
}

export interface CurrentUser extends User {
  email: string
  unreadNotifications: number
  draftCount: number
  bookmarkCount: number
}

// Article types
export type ArticleLevel = 'beginner' | 'intermediate' | 'advanced'
export type ContentType = 'article' | 'news' | 'post' | 'question' | 'discussion'
export type ArticleStatus = 'draft' | 'pending' | 'published' | 'archived'

export interface Category {
  id: string
  name: string
  slug: string
  icon: string
  description?: string
  articleCount: number
}

export interface Tag {
  id: string
  name: string
  slug: string
}

export interface ReactionCount {
  emoji: string
  count: number
  isReacted: boolean
}

export interface CoverImage {
  url: string
  altText?: string
  width: number
  height: number
}

export interface ArticleCard {
  id: string
  title: string
  slug: string
  lead?: string
  coverImage?: CoverImage
  level: ArticleLevel
  contentType: ContentType
  readingTime: number
  isEditorial: boolean
  isPinned: boolean
  
  author: {
    id: string
    username: string
    displayName: string
    avatarUrl?: string
    isVerified: boolean
  }
  
  category: Category
  tags: Tag[]
  
  reactions: ReactionCount[]
  commentCount: number
  viewCount: number
  bookmarkCount: number
  
  publishedAt: string
}

export interface Article extends ArticleCard {
  content: string
  htmlContent: string
  status: ArticleStatus
  
  createdAt: string
  updatedAt: string
  
  // SEO
  metaTitle?: string
  metaDescription?: string
  canonicalUrl?: string
  
  // Settings
  commentsEnabled: boolean
  isNSFW: boolean
}

// Comment types
export interface Comment {
  id: string
  content: string
  htmlContent: string
  
  author: {
    id: string
    username: string
    displayName: string
    avatarUrl?: string
    isVerified: boolean
  }
  
  articleId: string
  parentId?: string
  
  reactions: ReactionCount[]
  replyCount: number
  
  isEdited: boolean
  createdAt: string
  updatedAt?: string
  
  // Tree structure
  depth: number
  replies?: Comment[]
}

// Notification types
export type NotificationType = 
  | 'new_comment'
  | 'comment_reply'
  | 'reaction'
  | 'new_follower'
  | 'article_published'
  | 'mention'
  | 'system'

export interface Notification {
  id: string
  type: NotificationType
  title: string
  message: string
  link?: string
  imageUrl?: string
  isRead: boolean
  createdAt: string
  
  actor?: {
    id: string
    username: string
    displayName: string
    avatarUrl?: string
  }
}

// Achievement types
export interface Achievement {
  id: string
  name: string
  description: string
  icon: string
  rarity: 'common' | 'uncommon' | 'rare' | 'epic' | 'legendary'
  unlockedAt?: string
  progress?: {
    current: number
    required: number
  }
}

// Draft types
export interface Draft {
  id: string
  title: string
  content: string
  level?: ArticleLevel
  contentType?: ContentType
  categoryId?: string
  tags: string[]
  coverImage?: CoverImage
  createdAt: string
  updatedAt: string
  autoSavedAt?: string
}

// Tool types (for AI tools catalog)
export interface Tool {
  id: string
  name: string
  description: string
  shortDescription: string
  icon: string
  url: string
  category: string
  tags: Tag[]
  isPremium: boolean
  isFeatured: boolean
  rating: number
  reviewCount: number
  addedAt: string
}

// Theme types
export type Theme = 'light' | 'dark' | 'system'

// Pagination
export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  pageSize: number
  hasMore: boolean
}

// API types
export interface ApiError {
  code: string
  message: string
  details?: Record<string, string[]>
}

export interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: ApiError
}

// Auth types
export interface LoginRequest {
  email: string
  password: string
  remember?: boolean
}

export interface RegisterRequest {
  email: string
  username: string
  password: string
}

export interface AuthResponse {
  user: CurrentUser
  accessToken: string
  refreshToken: string
  expiresAt: number
}

// Feed filters
export interface FeedFilters {
  sort: 'popular' | 'new' | 'hot'
  level?: ArticleLevel
  contentType?: ContentType
  categoryId?: string
  tagId?: string
  timeRange?: '24h' | '7d' | '30d' | 'all'
}

// Search
export interface SearchResult {
  type: 'article' | 'user' | 'category' | 'tool'
  id: string
  title: string
  description?: string
  url: string
  imageUrl?: string
  highlights?: string[]
}

export interface SearchResponse {
  results: SearchResult[]
  total: number
  query: string
  took: number
}

// Subscription/Bookmark
export interface Bookmark {
  id: string
  articleId: string
  article: ArticleCard
  createdAt: string
  folder?: string
}

export interface Follow {
  id: string
  userId: string
  user: User
  createdAt: string
}

// Report
export type ReportReason = 
  | 'spam'
  | 'harassment'
  | 'inappropriate'
  | 'misinformation'
  | 'copyright'
  | 'other'

export interface Report {
  id: string
  reason: ReportReason
  description?: string
  targetType: 'article' | 'comment' | 'user'
  targetId: string
  status: 'pending' | 'resolved' | 'dismissed'
  createdAt: string
}

// Settings
export interface UserSettings {
  theme: Theme
  language: string
  emailNotifications: {
    newFollower: boolean
    newComment: boolean
    mentions: boolean
    newsletter: boolean
    marketing: boolean
  }
  privacy: {
    showEmail: boolean
    showOnline: boolean
    allowMessages: boolean
  }
  feed: {
    defaultSort: 'popular' | 'new'
    showNSFW: boolean
  }
}
