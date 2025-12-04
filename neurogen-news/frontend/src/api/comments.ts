import { api } from './client'
import type { Comment, PaginatedResponse } from '@/types'

export interface CommentListParams {
  sort?: 'new' | 'popular' | 'old'
  page?: number
  pageSize?: number
}

export interface CreateCommentInput {
  articleId: string
  parentId?: string
  content: string
}

export const commentsApi = {
  // Get comments for article
  getByArticle: (articleId: string, params?: CommentListParams) =>
    api.get<PaginatedResponse<Comment>>(`/comments/article/${articleId}`, { params }),

  // Get comment by ID
  getById: (id: string) =>
    api.get<Comment>(`/comments/${id}`),

  // Get replies for comment
  getReplies: (id: string, limit = 10, offset = 0) =>
    api.get<{ items: Comment[] }>(`/comments/${id}/replies`, { params: { limit, offset } }),

  // Create comment
  create: (data: CreateCommentInput) =>
    api.post<Comment>('/comments', data),

  // Update comment
  update: (id: string, content: string) =>
    api.put<Comment>(`/comments/${id}`, { content }),

  // Delete comment
  delete: (id: string) =>
    api.delete(`/comments/${id}`),

  // Add reaction
  addReaction: (id: string, emoji: string) =>
    api.post(`/comments/${id}/reactions`, { emoji }),

  // Remove reaction
  removeReaction: (id: string) =>
    api.delete(`/comments/${id}/reactions`),
}

export default commentsApi


