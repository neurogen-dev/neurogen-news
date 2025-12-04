<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterView } from 'vue-router'
import AppHeader from '@/components/layout/AppHeader.vue'
import AppSidebar from '@/components/layout/AppSidebar.vue'
import MobileNav from '@/components/layout/MobileNav.vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const isLoaded = ref(false)

onMounted(async () => {
  // Set dark theme by default
  document.documentElement.classList.add('dark')
  
  // Try to fetch current user if token exists
  if (authStore.accessToken) {
    await authStore.fetchCurrentUser()
  }
  
  // Trigger entrance animation
  setTimeout(() => {
    isLoaded.value = true
  }, 100)
})
</script>

<template>
  <div 
    class="min-h-screen bg-bg-base text-text-primary transition-all duration-700"
    :class="isLoaded ? 'opacity-100' : 'opacity-0'"
  >
    <!-- Animated Gradient Orbs - Empatra Style -->
    <div class="fixed inset-0 pointer-events-none z-0 overflow-hidden">
      <div 
        class="absolute top-[-15%] left-[-10%] w-[500px] h-[500px] bg-primary/20 rounded-full blur-[120px] animate-float opacity-40"
      ></div>
      <div 
        class="absolute bottom-[10%] right-[-10%] w-[400px] h-[400px] bg-secondary/15 rounded-full blur-[100px] animate-float opacity-30" 
        style="animation-delay: -10s;"
      ></div>
      <div 
        class="absolute top-[40%] left-[20%] w-[300px] h-[300px] bg-accent/10 rounded-full blur-[80px] animate-float opacity-25" 
        style="animation-delay: -18s;"
      ></div>
    </div>
    
    <!-- Header -->
    <AppHeader />
    
    <!-- Main layout - Unique structure (not like DTF) -->
    <div class="relative z-10 flex pt-16">
      <!-- Sidebar (desktop) - Collapsible, not fixed width like DTF -->
      <AppSidebar class="hidden lg:block" />
      
      <!-- Main content - Centered, wider content area -->
      <main class="flex-1 min-h-[calc(100vh-4rem)] lg:ml-64 pb-24 lg:pb-8">
        <div class="max-w-4xl mx-auto px-4 sm:px-6 py-8">
          <RouterView v-slot="{ Component, route }">
            <Transition 
              name="page" 
              mode="out-in"
              appear
            >
              <component 
                :is="Component" 
                :key="route.path"
              />
            </Transition>
          </RouterView>
        </div>
      </main>
    </div>
    
    <!-- Mobile navigation -->
    <MobileNav class="lg:hidden" />
  </div>
</template>

<style>
/* Page transitions - Empatra style */
.page-enter-active,
.page-leave-active {
  transition: 
    opacity 0.3s cubic-bezier(0.4, 0, 0.2, 1), 
    transform 0.3s cubic-bezier(0.4, 0, 0.2, 1),
    filter 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  will-change: opacity, transform, filter;
}

.page-enter-from {
  opacity: 0;
  transform: translateY(20px) scale(0.98);
  filter: blur(6px);
}

.page-leave-to {
  opacity: 0;
  transform: translateY(-10px) scale(0.99);
  filter: blur(3px);
}

.page-enter-active,
.page-leave-active {
  transform-style: preserve-3d;
  backface-visibility: hidden;
}
</style>
