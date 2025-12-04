<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { RouterLink } from 'vue-router'
import { Bell, Check, Settings } from 'lucide-vue-next'
import Avatar from '@/components/ui/Avatar.vue'
import { formatRelativeTime } from '@/utils/formatters'
import type { Notification } from '@/types'

const emit = defineEmits<{
  close: []
}>()

const panelRef = ref<HTMLElement | null>(null)
const notifications = ref<Notification[]>([])
const isLoading = ref(true)

// TODO: Fetch from API
const mockNotifications: Notification[] = [
  {
    id: '1',
    type: 'new_comment',
    title: 'Новый комментарий',
    message: 'Алексей прокомментировал вашу статью "Как использовать ChatGPT для работы"',
    link: '/chatbots/how-to-use-chatgpt#comments',
    isRead: false,
    createdAt: new Date(Date.now() - 1000 * 60 * 30).toISOString(),
    actor: {
      id: '2',
      username: 'alexey',
      displayName: 'Алексей',
      avatarUrl: undefined
    }
  },
  {
    id: '2',
    type: 'reaction',
    title: 'Реакция',
    message: 'Марии понравилась ваша статья',
    link: '/images/best-ai-art-generators',
    isRead: false,
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 2).toISOString(),
    actor: {
      id: '3',
      username: 'maria',
      displayName: 'Мария',
      avatarUrl: undefined
    }
  },
  {
    id: '3',
    type: 'new_follower',
    title: 'Новый подписчик',
    message: 'Иван подписался на вас',
    link: '/@ivan',
    isRead: true,
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 24).toISOString(),
    actor: {
      id: '4',
      username: 'ivan',
      displayName: 'Иван',
      avatarUrl: undefined
    }
  }
]

const handleClickOutside = (event: MouseEvent) => {
  if (panelRef.value && !panelRef.value.contains(event.target as Node)) {
    emit('close')
  }
}

const markAllAsRead = () => {
  notifications.value = notifications.value.map(n => ({ ...n, isRead: true }))
  // TODO: Call API to mark all as read
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  // Simulate loading
  setTimeout(() => {
    notifications.value = mockNotifications
    isLoading.value = false
  }, 500)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <div 
    ref="panelRef"
    class="absolute right-0 top-full mt-2 w-96 max-h-[32rem] velvet-panel shadow-floating overflow-hidden animate-slide-up"
  >
    <!-- Header -->
    <div class="flex items-center justify-between px-4 py-3 border-b border-gray-100 dark:border-gray-800">
      <h3 class="font-semibold text-gray-900 dark:text-white">
        Уведомления
      </h3>
      <div class="flex items-center gap-1">
        <button
          @click="markAllAsRead"
          class="p-2 text-gray-400 hover:text-primary hover:bg-primary/10 rounded-lg transition-all duration-200"
          title="Отметить все как прочитанные"
        >
          <Check class="w-4 h-4" />
        </button>
        <RouterLink
          to="/settings/notifications"
          @click="emit('close')"
          class="p-2 text-gray-400 hover:text-primary hover:bg-primary/10 rounded-lg transition-all duration-200"
          title="Настройки уведомлений"
        >
          <Settings class="w-4 h-4" />
        </RouterLink>
      </div>
    </div>
    
    <!-- Notifications list -->
    <div class="overflow-y-auto max-h-80 scrollbar-thin">
      <!-- Loading -->
      <div v-if="isLoading" class="p-8 flex justify-center">
        <div class="relative w-8 h-8">
          <div class="absolute inset-0 rounded-full border-2 border-primary/20"></div>
          <div class="absolute inset-0 rounded-full border-2 border-primary border-t-transparent animate-spin"></div>
        </div>
      </div>
      
      <!-- Empty state -->
      <div v-else-if="notifications.length === 0" class="p-8 text-center">
        <div class="w-14 h-14 mx-auto mb-4 rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center">
          <Bell class="w-7 h-7 text-gray-400" />
        </div>
        <p class="text-gray-500">Нет новых уведомлений</p>
      </div>
      
      <!-- Notification items -->
      <template v-else>
        <RouterLink
          v-for="notification in notifications"
          :key="notification.id"
          :to="notification.link || '#'"
          @click="emit('close')"
          class="flex gap-3 px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-all duration-200"
          :class="{ 'bg-primary/5': !notification.isRead }"
        >
          <!-- Actor avatar or icon -->
          <Avatar 
            v-if="notification.actor"
            :src="notification.actor.avatarUrl" 
            :alt="notification.actor.displayName"
            size="md"
          />
          <div 
            v-else 
            class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center shrink-0"
          >
            <Bell class="w-5 h-5 text-primary" />
          </div>
          
          <!-- Content -->
          <div class="flex-1 min-w-0">
            <p class="text-sm text-gray-800 dark:text-white line-clamp-2">
              {{ notification.message }}
            </p>
            <time class="text-xs text-gray-400 mt-1 block">
              {{ formatRelativeTime(notification.createdAt) }}
            </time>
          </div>
          
          <!-- Unread indicator -->
          <div 
            v-if="!notification.isRead"
            class="w-2 h-2 bg-primary rounded-full shrink-0 mt-2 glow-pulse"
          />
        </RouterLink>
      </template>
    </div>
    
    <!-- Footer -->
    <div class="px-4 py-3 border-t border-gray-100 dark:border-gray-800 text-center">
      <RouterLink
        to="/notifications"
        @click="emit('close')"
        class="text-sm text-primary hover:text-primary-600 transition-colors font-medium"
      >
        Все уведомления
      </RouterLink>
    </div>
  </div>
</template>
