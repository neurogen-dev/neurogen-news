import { api } from './client'
import type { Notification } from '@/types'

export interface NotificationListResult {
  items: Notification[]
  total: number
  unreadCount: number
  hasMore: boolean
}

export const notificationsApi = {
  // Get notifications
  list: (limit = 20, offset = 0) =>
    api.get<NotificationListResult>('/notifications', { params: { limit, offset } }),

  // Get unread count
  getUnreadCount: () =>
    api.get<{ count: number }>('/notifications/unread'),

  // Mark as read
  markAsRead: (id: string) =>
    api.put(`/notifications/${id}/read`),

  // Mark all as read
  markAllAsRead: () =>
    api.put('/notifications/read-all'),

  // Delete notification
  delete: (id: string) =>
    api.delete(`/notifications/${id}`),

  // Delete all notifications
  deleteAll: () =>
    api.delete('/notifications'),
}

export default notificationsApi


