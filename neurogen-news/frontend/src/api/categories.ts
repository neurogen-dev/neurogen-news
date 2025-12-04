import { api } from './client'
import type { Category, ArticleCard, PaginatedResponse } from '@/types'

export interface CategoryArticleParams {
  sort?: 'popular' | 'new' | 'hot'
  level?: string
  contentType?: string
  timeRange?: string
  page?: number
  pageSize?: number
}

export const categoriesApi = {
  // Get all categories
  list: () =>
    api.get<{ items: Category[] }>('/categories'),

  // Get category by slug
  getBySlug: (slug: string) =>
    api.get<{ category: Category; isSubscribed: boolean }>(`/categories/${slug}`),

  // Get articles for category
  getArticles: (slug: string, params?: CategoryArticleParams) =>
    api.get<PaginatedResponse<ArticleCard>>(`/categories/${slug}/articles`, { params }),

  // Subscribe to category
  subscribe: (slug: string) =>
    api.post(`/categories/${slug}/subscribe`),

  // Unsubscribe from category
  unsubscribe: (slug: string) =>
    api.delete(`/categories/${slug}/subscribe`),

  // Get user subscriptions
  getSubscriptions: () =>
    api.get<{ items: Category[] }>('/categories/subscriptions'),
}

export default categoriesApi


