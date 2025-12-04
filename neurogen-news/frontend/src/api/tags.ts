import { api } from './client'
import type { Tag, ArticleCard, PaginatedResponse } from '@/types'

export const tagsApi = {
  // Get all tags
  list: (limit = 50) =>
    api.get<{ items: Tag[] }>('/tags', { params: { limit } }),

  // Get popular tags
  getPopular: (limit = 20) =>
    api.get<{ items: Tag[] }>('/tags/popular', { params: { limit } }),

  // Search tags
  search: (query: string, limit = 10) =>
    api.get<{ items: Tag[] }>('/tags/search', { params: { q: query, limit } }),

  // Get tag by slug
  getBySlug: (slug: string) =>
    api.get<Tag>(`/tags/${slug}`),

  // Get articles for tag
  getArticles: (slug: string, params?: { page?: number; pageSize?: number }) =>
    api.get<PaginatedResponse<ArticleCard>>(`/tags/${slug}/articles`, { params }),
}

export default tagsApi


