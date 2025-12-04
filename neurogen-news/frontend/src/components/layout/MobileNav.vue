<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { Home, TrendingUp, Search, PenSquare, User } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const navItems = computed(() => [
  { name: 'Главная', icon: Home, to: '/', exact: true },
  { name: 'Популярное', icon: TrendingUp, to: '/popular' },
  { name: 'Поиск', icon: Search, to: '/search' },
  { name: 'Написать', icon: PenSquare, to: '/editor/new', requiresAuth: true },
  { 
    name: 'Профиль', 
    icon: User, 
    to: authStore.isLoggedIn ? `/@${authStore.user?.username}` : '/login'
  },
])

const isActive = (item: typeof navItems.value[0]) => {
  if (item.exact) return route.path === item.to
  return route.path.startsWith(item.to)
}

const handleNavClick = (item: typeof navItems.value[0], event: Event) => {
  if (item.requiresAuth && !authStore.isLoggedIn) {
    event.preventDefault()
    router.push({ name: 'login', query: { redirect: item.to } })
  }
}
</script>

<template>
  <nav class="fixed bottom-0 left-0 right-0 z-50 bg-bg-elevated/95 backdrop-blur-xl border-t border-border-subtle safe-area-bottom">
    <!-- Ambient glow -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute -bottom-10 left-1/2 -translate-x-1/2 w-64 h-20 bg-primary/15 rounded-full blur-3xl"></div>
    </div>
    
    <div class="relative flex items-center justify-around h-16 px-2">
      <RouterLink
        v-for="item in navItems"
        :key="item.to"
        :to="item.to"
        @click="(e: Event) => handleNavClick(item, e)"
        class="flex flex-col items-center justify-center gap-1 px-3 py-2 rounded-xl transition-all duration-200 min-w-[60px]"
        :class="[
          isActive(item)
            ? 'text-primary'
            : 'text-text-tertiary hover:text-text-secondary'
        ]"
      >
        <div 
          class="p-1.5 rounded-lg transition-all duration-200"
          :class="isActive(item) ? 'bg-primary/15' : ''"
        >
          <component 
            :is="item.icon" 
            class="w-5 h-5 transition-transform duration-200"
            :class="isActive(item) ? 'scale-110' : ''"
            stroke-width="2"
          />
        </div>
        <span 
          class="text-[10px] font-medium transition-all duration-200"
          :class="isActive(item) ? 'opacity-100' : 'opacity-60'"
        >
          {{ item.name }}
        </span>
      </RouterLink>
    </div>
  </nav>
</template>
