import { api } from './client'
import type { ArticleCard, PaginatedResponse } from '@/types'

export interface BookmarkFolder {
  id: string
  name: string
  bookmarkCount: number
  createdAt: string
}

export interface BookmarkListParams {
  folderId?: string
  page?: number
  pageSize?: number
}

export const bookmarksApi = {
  // Get bookmarks
  list: (params?: BookmarkListParams) =>
    api.get<PaginatedResponse<ArticleCard>>('/bookmarks', { params }),

  // Add bookmark
  add: (articleId: string, folderId?: string) =>
    api.post('/bookmarks', { articleId, folderId }),

  // Remove bookmark
  remove: (articleId: string) =>
    api.delete(`/bookmarks/${articleId}`),

  // Check if bookmarked
  check: (articleId: string) =>
    api.get<{ isBookmarked: boolean }>(`/bookmarks/${articleId}/check`),

  // Get folders
  getFolders: () =>
    api.get<{ items: BookmarkFolder[] }>('/bookmarks/folders'),

  // Create folder
  createFolder: (name: string) =>
    api.post<BookmarkFolder>('/bookmarks/folders', { name }),

  // Update folder
  updateFolder: (id: string, name: string) =>
    api.put(`/bookmarks/folders/${id}`, { name }),

  // Delete folder
  deleteFolder: (id: string) =>
    api.delete(`/bookmarks/folders/${id}`),

  // Move bookmark to folder
  moveToFolder: (articleId: string, folderId?: string) =>
    api.put('/bookmarks/move', { articleId, folderId }),
}

export default bookmarksApi


