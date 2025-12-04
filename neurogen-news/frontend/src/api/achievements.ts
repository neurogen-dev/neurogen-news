import { api } from './client'
import type { Achievement } from '@/types'

export interface UserAchievement {
  id: string
  achievementId: string
  awardedAt: string
  achievement: Achievement
}

export interface AchievementProgress {
  articleCount: number
  commentCount: number
  totalViews: number
  followerCount: number
  achievementCount: number
  totalPoints: number
}

export const achievementsApi = {
  // Get all achievements
  list: () =>
    api.get<{ items: Achievement[] }>('/achievements'),

  // Get achievement by ID
  getById: (id: string) =>
    api.get<Achievement>(`/achievements/${id}`),

  // Get current user achievements
  getMyAchievements: () =>
    api.get<{ items: UserAchievement[] }>('/achievements/my'),

  // Get user achievements
  getUserAchievements: (username: string) =>
    api.get<{ items: UserAchievement[] }>(`/achievements/user/${username}`),

  // Get progress
  getProgress: () =>
    api.get<AchievementProgress>('/achievements/progress'),

  // Check for new achievements
  checkAchievements: () =>
    api.post<{ newAchievements: Achievement[] }>('/achievements/check'),
}

export default achievementsApi


