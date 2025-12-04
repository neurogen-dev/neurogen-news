import { api } from './client'
import type { User, ArticleCard, PaginatedResponse } from '@/types'

export interface UserProfile {
  user: User
  articleCount: number
  commentCount: number
  followerCount: number
  followingCount: number
  isFollowing: boolean
}

export interface UpdateProfileInput {
  displayName?: string
  bio?: string
  avatarUrl?: string
  coverUrl?: string
  location?: string
  website?: string
  telegram?: string
  github?: string
}

export const usersApi = {
  // Get current user profile
  getCurrentProfile: () =>
    api.get<UserProfile>('/users/me'),

  // Update current user profile
  updateProfile: (data: UpdateProfileInput) =>
    api.put<User>('/users/me', data),

  // Get user profile by username
  getProfile: (username: string) =>
    api.get<UserProfile>(`/users/${username}`),

  // Get user articles
  getArticles: (username: string, limit = 20, offset = 0) =>
    api.get<PaginatedResponse<ArticleCard>>(`/users/${username}/articles`, { 
      params: { limit, offset } 
    }),

  // Get user followers
  getFollowers: (username: string, limit = 20, offset = 0) =>
    api.get<{ items: User[]; total: number; hasMore: boolean }>(`/users/${username}/followers`, { 
      params: { limit, offset } 
    }),

  // Get user following
  getFollowing: (username: string, limit = 20, offset = 0) =>
    api.get<{ items: User[]; total: number; hasMore: boolean }>(`/users/${username}/following`, { 
      params: { limit, offset } 
    }),

  // Follow user
  follow: (username: string) =>
    api.post(`/users/${username}/follow`),

  // Unfollow user
  unfollow: (username: string) =>
    api.delete(`/users/${username}/follow`),
}

export default usersApi


