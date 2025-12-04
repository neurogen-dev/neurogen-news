import { api } from './client'
import type { ArticleCard, User, Tag } from '@/types'

export interface SearchParams {
  q: string
  type?: 'all' | 'articles' | 'users' | 'tags'
  sort?: 'relevance' | 'new' | 'popular'
  level?: string
  category?: string
  limit?: number
  offset?: number
}

export interface ArticleSearchResult {
  items: ArticleCard[]
  total: number
  hasMore: boolean
}

export interface UserSearchResult {
  items: User[]
  total: number
  hasMore: boolean
}

export interface SearchResult {
  articles?: ArticleSearchResult
  users?: UserSearchResult
  tags?: Tag[]
  total: number
}

export interface SearchSuggestion {
  type: 'article' | 'user' | 'tag'
  text: string
  slug?: string
  image?: string
}

export const searchApi = {
  // Search
  search: (params: SearchParams) =>
    api.get<SearchResult>('/search', { params }),

  // Get search suggestions
  getSuggestions: (query: string, limit = 5) =>
    api.get<{ items: SearchSuggestion[] }>('/search/suggestions', { 
      params: { q: query, limit } 
    }),
}

export default searchApi


