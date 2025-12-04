<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { 
  User, 
  FileText, 
  Bookmark, 
  Award, 
  Settings, 
  LogOut,
  Moon,
  Sun,
  Monitor
} from 'lucide-vue-next'
import Avatar from '@/components/ui/Avatar.vue'
import { useAuthStore } from '@/stores/auth'
import { formatKarma } from '@/utils/formatters'
import type { Theme } from '@/types'

const emit = defineEmits<{
  close: []
}>()

const router = useRouter()
const authStore = useAuthStore()
const menuRef = ref<HTMLElement | null>(null)

const menuItems = [
  { name: 'Мой профиль', icon: User, to: `/@${authStore.user?.username}` },
  { name: 'Черновики', icon: FileText, to: '/drafts', badge: authStore.user?.draftCount },
  { name: 'Закладки', icon: Bookmark, to: '/bookmarks', badge: authStore.user?.bookmarkCount },
  { name: 'Достижения', icon: Award, to: '/achievements' },
  { name: 'Настройки', icon: Settings, to: '/settings' },
]

const themes = [
  { value: 'light' as Theme, icon: Sun, label: 'Светлая' },
  { value: 'dark' as Theme, icon: Moon, label: 'Тёмная' },
  { value: 'system' as Theme, icon: Monitor, label: 'Система' },
]

const handleLogout = async () => {
  await authStore.logout()
  emit('close')
  router.push('/')
}

const handleThemeChange = (theme: Theme) => {
  authStore.setTheme(theme)
}

const handleClickOutside = (event: MouseEvent) => {
  if (menuRef.value && !menuRef.value.contains(event.target as Node)) {
    emit('close')
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <div 
    ref="menuRef"
    class="absolute right-0 top-full mt-2 w-72 velvet-panel shadow-floating overflow-hidden animate-slide-up"
  >
    <!-- User info -->
    <div class="p-4 border-b border-gray-100 dark:border-gray-800">
      <div class="flex items-center gap-3">
        <Avatar 
          :src="authStore.user?.avatarUrl" 
          :alt="authStore.user?.displayName"
          size="lg"
          ring
        />
        <div class="flex-1 min-w-0">
          <div class="font-semibold text-gray-900 dark:text-white truncate">
            {{ authStore.user?.displayName }}
          </div>
          <div class="text-sm text-gray-500 truncate">
            @{{ authStore.user?.username }}
          </div>
        </div>
      </div>
      
      <!-- Karma -->
      <div class="mt-3 flex items-center gap-4 text-sm">
        <div>
          <span class="text-gray-500">Карма:</span>
          <span 
            class="ml-1 font-medium"
            :class="{
              'text-success': (authStore.user?.karma ?? 0) > 0,
              'text-error': (authStore.user?.karma ?? 0) < 0,
              'text-gray-400': authStore.user?.karma === 0
            }"
          >
            {{ formatKarma(authStore.user?.karma ?? 0) }}
          </span>
        </div>
        <div>
          <span class="text-gray-500">Подписчики:</span>
          <span class="ml-1 font-medium text-gray-800 dark:text-white">
            {{ authStore.user?.followerCount }}
          </span>
        </div>
      </div>
    </div>
    
    <!-- Menu items -->
    <div class="py-2">
      <RouterLink
        v-for="item in menuItems"
        :key="item.to"
        :to="item.to"
        @click="emit('close')"
        class="flex items-center gap-3 px-4 py-2.5 text-sm text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-all duration-200"
      >
        <component :is="item.icon" class="w-4 h-4" />
        <span>{{ item.name }}</span>
        <span 
          v-if="item.badge" 
          class="ml-auto velvet-badge !bg-primary/10 !text-primary !text-xs"
        >
          {{ item.badge }}
        </span>
      </RouterLink>
    </div>
    
    <!-- Theme switcher -->
    <div class="px-4 py-3 border-t border-gray-100 dark:border-gray-800">
      <div class="text-xs font-medium text-gray-400 mb-2 uppercase tracking-wider">Тема оформления</div>
      <div class="velvet-tabs">
        <button
          v-for="theme in themes"
          :key="theme.value"
          @click="handleThemeChange(theme.value)"
          class="velvet-tab flex-1"
          :class="{ 'active': authStore.theme === theme.value }"
        >
          <span class="velvet-tab-bg"></span>
          <span class="relative z-10 flex items-center justify-center gap-1.5">
            <component :is="theme.icon" class="w-4 h-4" />
            <span class="hidden sm:inline">{{ theme.label }}</span>
          </span>
        </button>
      </div>
    </div>
    
    <!-- Logout -->
    <div class="py-2 border-t border-gray-100 dark:border-gray-800">
      <button
        @click="handleLogout"
        class="flex items-center gap-3 w-full px-4 py-2.5 text-sm text-error hover:bg-error/5 transition-colors"
      >
        <LogOut class="w-4 h-4" />
        Выйти
      </button>
    </div>
  </div>
</template>
