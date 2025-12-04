<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { Search, PenSquare, Bell, Menu, X, Sparkles, Zap } from 'lucide-vue-next'
import Avatar from '@/components/ui/Avatar.vue'
import Button from '@/components/ui/Button.vue'
import UserMenu from './UserMenu.vue'
import NotificationPanel from './NotificationPanel.vue'
import SearchModal from './SearchModal.vue'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const showSearch = ref(false)
const showUserMenu = ref(false)
const showNotifications = ref(false)
const showMobileMenu = ref(false)
const scrolled = ref(false)

// Track scroll for glass effect intensity
const handleScroll = () => {
  scrolled.value = window.scrollY > 20
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll, { passive: true })
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})

const handleWrite = () => {
  if (authStore.isLoggedIn) {
    router.push({ name: 'editor-new' })
  } else {
    router.push({ name: 'login', query: { redirect: '/editor/new' } })
  }
}

// Keyboard shortcut for search
onMounted(() => {
  const handleKeydown = (e: KeyboardEvent) => {
    if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
      e.preventDefault()
      showSearch.value = true
    }
  }
  window.addEventListener('keydown', handleKeydown)
  onUnmounted(() => window.removeEventListener('keydown', handleKeydown))
})
</script>

<template>
  <header 
    class="fixed top-0 left-0 right-0 z-50 h-16 transition-all duration-300"
    :class="[
      scrolled 
        ? 'bg-bg-base/90 backdrop-blur-xl border-b border-border-subtle' 
        : 'bg-transparent'
    ]"
  >
    <div class="relative h-full flex items-center justify-between px-4 lg:px-6 max-w-screen-2xl mx-auto">
      <!-- Left section -->
      <div class="flex items-center gap-4">
        <!-- Mobile menu toggle -->
        <button 
          @click="showMobileMenu = !showMobileMenu"
          class="lg:hidden p-2 -ml-2 text-text-secondary hover:text-text-primary transition-colors rounded-xl hover:bg-bg-surface"
        >
          <Menu v-if="!showMobileMenu" class="w-6 h-6" stroke-width="2" />
          <X v-else class="w-6 h-6" stroke-width="2" />
        </button>
        
        <!-- Logo - Empatra Style -->
        <RouterLink 
          to="/" 
          class="flex items-center gap-3 font-bold text-xl group"
        >
          <!-- Animated Logo Mark -->
          <div class="relative w-9 h-9">
            <div class="absolute inset-0 bg-gradient-to-br from-primary via-secondary to-accent rounded-xl rotate-3 group-hover:rotate-6 transition-transform duration-300 opacity-80"></div>
            <div class="absolute inset-0.5 bg-bg-base rounded-[10px] flex items-center justify-center">
              <Zap class="w-5 h-5 text-primary group-hover:scale-110 transition-transform" fill="currentColor" fill-opacity="0.2" stroke-width="2.5" />
            </div>
          </div>
          <span class="hidden sm:flex items-baseline gap-1">
            <span class="text-gradient font-display text-xl tracking-tight">Neurogen</span>
            <span class="text-text-tertiary font-normal text-lg">.News</span>
          </span>
        </RouterLink>
      </div>
      
      <!-- Center section - Search -->
      <button
        @click="showSearch = true"
        class="hidden md:flex items-center gap-3 px-4 py-2.5 empatra-panel !rounded-full text-text-secondary hover:text-text-primary transition-all duration-200 w-80 max-w-md group"
      >
        <Search class="w-4 h-4 group-hover:text-primary transition-colors" stroke-width="2.5" />
        <span class="text-sm">Поиск статей, инструментов...</span>
        <kbd class="ml-auto text-xs bg-bg-surface/80 px-2 py-0.5 rounded-md font-mono border border-border-subtle">⌘K</kbd>
      </button>
      
      <!-- Right section -->
      <div class="flex items-center gap-2">
        <!-- Search (mobile) -->
        <button
          @click="showSearch = true"
          class="md:hidden p-2.5 text-text-secondary hover:text-text-primary transition-colors rounded-xl hover:bg-bg-surface"
        >
          <Search class="w-5 h-5" stroke-width="2.5" />
        </button>
        
        <!-- Write button -->
        <Button 
          @click="handleWrite"
          variant="primary"
          size="sm"
          class="hidden sm:flex"
        >
          <PenSquare class="w-4 h-4" stroke-width="2.5" />
          <span>Написать</span>
          <Sparkles class="w-3.5 h-3.5 opacity-60" stroke-width="2.5" />
        </Button>
        
        <template v-if="authStore.isLoggedIn">
          <!-- Notifications -->
          <div class="relative">
            <button
              @click="showNotifications = !showNotifications"
              class="relative p-2.5 text-text-secondary hover:text-text-primary transition-all duration-200 rounded-xl hover:bg-bg-surface"
            >
              <Bell class="w-5 h-5" stroke-width="2.5" />
              <span 
                v-if="authStore.user?.unreadNotifications"
                class="absolute top-1.5 right-1.5 w-2.5 h-2.5 bg-gradient-to-br from-secondary to-error rounded-full ring-2 ring-bg-base animate-pulse"
              />
            </button>
            
            <Transition
              enter-active-class="transition duration-200 ease-spring"
              enter-from-class="opacity-0 scale-95 -translate-y-2"
              enter-to-class="opacity-100 scale-100 translate-y-0"
              leave-active-class="transition duration-150 ease-smooth"
              leave-from-class="opacity-100 scale-100 translate-y-0"
              leave-to-class="opacity-0 scale-95 -translate-y-2"
            >
              <NotificationPanel 
                v-if="showNotifications" 
                @close="showNotifications = false"
              />
            </Transition>
          </div>
          
          <!-- User avatar -->
          <div class="relative">
            <button
              @click="showUserMenu = !showUserMenu"
              class="p-1 rounded-full ring-2 ring-transparent hover:ring-primary/30 transition-all"
            >
              <Avatar 
                :src="authStore.user?.avatarUrl" 
                :alt="authStore.user?.displayName"
                size="sm"
              />
            </button>
            
            <Transition
              enter-active-class="transition duration-200 ease-spring"
              enter-from-class="opacity-0 scale-95 -translate-y-2"
              enter-to-class="opacity-100 scale-100 translate-y-0"
              leave-active-class="transition duration-150 ease-smooth"
              leave-from-class="opacity-100 scale-100 translate-y-0"
              leave-to-class="opacity-0 scale-95 -translate-y-2"
            >
              <UserMenu 
                v-if="showUserMenu" 
                @close="showUserMenu = false"
              />
            </Transition>
          </div>
        </template>
        
        <template v-else>
          <Button
            variant="ghost"
            size="sm"
            to="/login"
          >
            Войти
          </Button>
        </template>
      </div>
    </div>
    
    <!-- Search modal -->
    <Transition
      enter-active-class="transition duration-200 ease-smooth"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-150 ease-smooth"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <SearchModal 
        v-if="showSearch" 
        @close="showSearch = false"
      />
    </Transition>
  </header>
</template>
