import { api } from './client'
import type { 
  Article, 
  ArticleCard, 
  PaginatedResponse, 
  FeedFilters 
} from '@/types'

export interface ArticleListParams {
  sort?: 'popular' | 'new' | 'hot'
  level?: 'beginner' | 'intermediate' | 'advanced'
  contentType?: string
  categoryId?: string
  tagId?: string
  timeRange?: '24h' | '7d' | '30d' | 'all'
  page?: number
  pageSize?: number
}

export interface CreateArticleInput {
  title: string
  content: string
  lead?: string
  coverImageUrl?: string
  level: 'beginner' | 'intermediate' | 'advanced'
  contentType: 'article' | 'news' | 'post' | 'question' | 'discussion'
  categoryId: string
  tags?: string[]
  isNsfw?: boolean
  commentsEnabled?: boolean
  status?: 'draft' | 'published'
  metaTitle?: string
  metaDescription?: string
}

export interface UpdateArticleInput extends Partial<CreateArticleInput> {}

export const articlesApi = {
  // Get articles list
  list: (params?: ArticleListParams) =>
    api.get<PaginatedResponse<ArticleCard>>('/articles', { params }),

  // Get article by ID
  getById: (id: string) =>
    api.get<Article>(`/articles/${id}`),

  // Get article by slug
  getBySlug: (category: string, slug: string) =>
    api.get<Article>(`/articles/slug/${category}/${slug}`),

  // Create article
  create: (data: CreateArticleInput) =>
    api.post<Article>('/articles', data),

  // Update article
  update: (id: string, data: UpdateArticleInput) =>
    api.put<Article>(`/articles/${id}`, data),

  // Delete article
  delete: (id: string) =>
    api.delete(`/articles/${id}`),

  // Add reaction
  addReaction: (id: string, emoji: string) =>
    api.post(`/articles/${id}/reactions`, { emoji }),

  // Remove reaction
  removeReaction: (id: string) =>
    api.delete(`/articles/${id}/reactions`),

  // Add bookmark
  bookmark: (id: string) =>
    api.post(`/articles/${id}/bookmark`),

  // Remove bookmark
  removeBookmark: (id: string) =>
    api.delete(`/articles/${id}/bookmark`),
}

export default articlesApi


