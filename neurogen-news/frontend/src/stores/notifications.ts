import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Notification } from '@/types'
import { notificationsApi } from '@/api'
import { useWebSocket } from '@/api/websocket'

export const useNotificationsStore = defineStore('notifications', () => {
  // State
  const notifications = ref<Notification[]>([])
  const unreadCount = ref(0)
  const total = ref(0)
  const hasMore = ref(false)
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  const isOpen = ref(false)

  const ws = useWebSocket()

  // Getters
  const isEmpty = computed(() => notifications.value.length === 0 && !isLoading.value)
  const unreadNotifications = computed(() => notifications.value.filter(n => !n.isRead))
  const hasUnread = computed(() => unreadCount.value > 0)

  // Actions
  async function fetchNotifications(limit = 20, offset = 0, append = false) {
    if (isLoading.value && !append) return

    isLoading.value = true
    error.value = null

    try {
      const response = await notificationsApi.list(limit, offset)

      if (append) {
        notifications.value = [...notifications.value, ...response.data.items]
      } else {
        notifications.value = response.data.items
      }

      total.value = response.data.total
      unreadCount.value = response.data.unreadCount
      hasMore.value = response.data.hasMore
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Не удалось загрузить уведомления'
      throw e
    } finally {
      isLoading.value = false
    }
  }

  async function fetchMore() {
    if (!hasMore.value || isLoading.value) return
    await fetchNotifications(20, notifications.value.length, true)
  }

  async function fetchUnreadCount() {
    try {
      const response = await notificationsApi.getUnreadCount()
      unreadCount.value = response.data.count
    } catch (e) {
      console.error('Failed to fetch unread count', e)
    }
  }

  async function markAsRead(id: string) {
    try {
      await notificationsApi.markAsRead(id)

      const notification = notifications.value.find(n => n.id === id)
      if (notification && !notification.isRead) {
        notification.isRead = true
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      }
    } catch (e) {
      console.error('Failed to mark as read', e)
      throw e
    }
  }

  async function markAllAsRead() {
    try {
      await notificationsApi.markAllAsRead()

      notifications.value = notifications.value.map(n => ({ ...n, isRead: true }))
      unreadCount.value = 0
    } catch (e) {
      console.error('Failed to mark all as read', e)
      throw e
    }
  }

  async function deleteNotification(id: string) {
    try {
      await notificationsApi.delete(id)

      const notification = notifications.value.find(n => n.id === id)
      if (notification && !notification.isRead) {
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      }

      notifications.value = notifications.value.filter(n => n.id !== id)
      total.value--
    } catch (e) {
      console.error('Failed to delete notification', e)
      throw e
    }
  }

  async function deleteAll() {
    try {
      await notificationsApi.deleteAll()
      notifications.value = []
      unreadCount.value = 0
      total.value = 0
    } catch (e) {
      console.error('Failed to delete all', e)
      throw e
    }
  }

  function addNotification(notification: Notification) {
    notifications.value = [notification, ...notifications.value]
    if (!notification.isRead) {
      unreadCount.value++
    }
    total.value++
  }

  // Subscribe to WebSocket notifications
  function subscribeToNotifications() {
    ws.on<Notification>('notification', (payload) => {
      addNotification(payload as Notification)
    })
  }

  function unsubscribeFromNotifications() {
    // Clear notification handler
  }

  function togglePanel() {
    isOpen.value = !isOpen.value
  }

  function openPanel() {
    isOpen.value = true
  }

  function closePanel() {
    isOpen.value = false
  }

  function reset() {
    notifications.value = []
    unreadCount.value = 0
    total.value = 0
    hasMore.value = false
    error.value = null
    isOpen.value = false
  }

  return {
    // State
    notifications,
    unreadCount,
    total,
    hasMore,
    isLoading,
    error,
    isOpen,
    // Getters
    isEmpty,
    unreadNotifications,
    hasUnread,
    // Actions
    fetchNotifications,
    fetchMore,
    fetchUnreadCount,
    markAsRead,
    markAllAsRead,
    deleteNotification,
    deleteAll,
    addNotification,
    subscribeToNotifications,
    unsubscribeFromNotifications,
    togglePanel,
    openPanel,
    closePanel,
    reset,
  }
})


