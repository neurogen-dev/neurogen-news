import { api } from './client'
import type { Draft, PaginatedResponse, ContentType } from '@/types'

export interface CreateDraftInput {
  title: string
  content: string
  coverImageUrl?: string
  contentType: ContentType
  categoryId?: string
  tags?: string[]
}

export interface UpdateDraftInput extends Partial<CreateDraftInput> {}

export interface AutoSaveInput {
  articleId?: string
  title: string
  content: string
  coverImageUrl?: string
  contentType: ContentType
  categoryId?: string
  tags?: string[]
}

export const draftsApi = {
  // Get drafts list
  list: (limit = 20, offset = 0) =>
    api.get<PaginatedResponse<Draft>>('/drafts', { params: { limit, offset } }),

  // Get draft by ID
  getById: (id: string) =>
    api.get<Draft>(`/drafts/${id}`),

  // Create draft
  create: (data: CreateDraftInput) =>
    api.post<Draft>('/drafts', data),

  // Update draft
  update: (id: string, data: UpdateDraftInput) =>
    api.put<Draft>(`/drafts/${id}`, data),

  // Delete draft
  delete: (id: string) =>
    api.delete(`/drafts/${id}`),

  // Auto-save
  autoSave: (data: AutoSaveInput) =>
    api.post('/drafts/autosave', data),

  // Get latest auto-save
  getAutoSave: (articleId?: string) =>
    api.get<{ draft: Draft | null }>('/drafts/autosave', { 
      params: articleId ? { articleId } : {} 
    }),
}

export default draftsApi


