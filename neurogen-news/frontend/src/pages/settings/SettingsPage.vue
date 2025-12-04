<script setup lang="ts">
import { RouterLink, RouterView, useRoute } from 'vue-router'
import { User, Lock, Bell, Rss, ChevronLeft } from 'lucide-vue-next'

const route = useRoute()

const menuItems = [
  { name: 'Профиль', icon: User, to: '/settings/profile' },
  { name: 'Аккаунт', icon: Lock, to: '/settings/account' },
  { name: 'Уведомления', icon: Bell, to: '/settings/notifications' },
  { name: 'Ленты', icon: Rss, to: '/settings/feeds' },
]

const isActive = (path: string) => route.path === path
</script>

<template>
  <div class="min-h-screen">
    <!-- Header -->
    <div class="flex items-center gap-4 mb-6">
      <RouterLink 
        to="/"
        class="p-2 text-text-tertiary hover:text-text-primary transition-colors lg:hidden"
      >
        <ChevronLeft class="w-5 h-5" />
      </RouterLink>
      <h1 class="text-2xl font-bold text-text-primary dark:text-white">
        Настройки
      </h1>
    </div>
    
    <div class="flex flex-col lg:flex-row gap-6">
      <!-- Sidebar -->
      <aside class="lg:w-64 shrink-0">
        <nav class="space-y-1">
          <RouterLink
            v-for="item in menuItems"
            :key="item.to"
            :to="item.to"
            class="flex items-center gap-3 px-4 py-3 rounded-lg transition-colors"
            :class="[
              isActive(item.to)
                ? 'bg-primary/10 text-primary'
                : 'text-text-secondary hover:text-text-primary hover:bg-background-secondary dark:hover:bg-dark-tertiary'
            ]"
          >
            <component :is="item.icon" class="w-5 h-5" />
            {{ item.name }}
          </RouterLink>
        </nav>
      </aside>
      
      <!-- Content -->
      <main class="flex-1 bg-white dark:bg-dark-secondary rounded-xl border border-border dark:border-dark-tertiary p-6">
        <RouterView />
      </main>
    </div>
  </div>
</template>

