<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Zap } from 'lucide-vue-next'
import FeedList from '@/components/feed/FeedList.vue'
import FeedFilters from '@/components/feed/FeedFilters.vue'
import StartHereSection from '@/components/common/StartHereSection.vue'
import QuickPost from '@/components/editor/QuickPost.vue'
import { useFeedStore } from '@/stores/feed'
import { useAuthStore } from '@/stores/auth'
import type { ArticleLevel, ContentType } from '@/types'

const route = useRoute()
const feedStore = useFeedStore()
const authStore = useAuthStore()

// Filters
const sort = ref<'popular' | 'new'>('popular')
const level = ref<ArticleLevel | undefined>()
const contentType = ref<ContentType | undefined>()

// Update sort from route
watch(
  () => route.name,
  (name) => {
    if (name === 'popular') sort.value = 'popular'
    else if (name === 'new') sort.value = 'new'
  },
  { immediate: true }
)

// Load articles
const loadArticles = async () => {
  await feedStore.fetchArticles({
    sort: sort.value,
    level: level.value,
    contentType: contentType.value,
  })
}

// Handle filter changes
watch([sort, level, contentType], loadArticles)

// Handle load more
const handleLoadMore = async () => {
  await feedStore.loadMore()
}

onMounted(loadArticles)
</script>

<template>
  <div class="min-h-screen">
    <!-- Welcome Hero - Empatra Style -->
    <div 
      v-if="!authStore.isLoggedIn"
      class="empatra-panel p-8 mb-8 text-center animate-slide-up overflow-hidden relative"
    >
      <!-- Gradient orbs -->
      <div class="absolute -top-20 -right-20 w-48 h-48 bg-gradient-to-br from-primary/30 to-secondary/20 rounded-full blur-3xl"></div>
      <div class="absolute -bottom-20 -left-20 w-40 h-40 bg-gradient-to-tr from-accent/20 to-primary/10 rounded-full blur-3xl"></div>
      
      <div class="relative">
        <h1 class="text-3xl md:text-4xl font-bold text-text-primary mb-4 font-display">
          Медиа об <span class="text-gradient">искусственном интеллекте</span>
        </h1>
        <p class="text-text-secondary max-w-xl mx-auto text-lg leading-relaxed">
          Гайды, обзоры инструментов, промпты и новости из мира AI. 
          Присоединяйтесь к сообществу энтузиастов нейросетей!
        </p>
        
        <!-- Empatra branding -->
        <div class="mt-6 flex items-center justify-center gap-2 text-text-tertiary text-sm">
          <Zap class="w-4 h-4 text-primary" />
          <span>Продукт</span>
          <a href="https://empatra.ai" target="_blank" rel="noopener noreferrer" class="text-primary hover:text-primary-400 transition-colors font-medium">
            Empatra AI
          </a>
        </div>
      </div>
    </div>
    
    <!-- Start Here Section for beginners -->
    <StartHereSection v-if="!authStore.hasSeenStartHere" class="mb-8" />
    
    <!-- Quick post (for logged in users) -->
    <QuickPost v-if="authStore.isLoggedIn" class="mb-8" />
    
    <!-- Feed filters -->
    <FeedFilters
      v-model:sort="sort"
      v-model:level="level"
      v-model:content-type="contentType"
      class="mb-8"
    />
    
    <!-- Article feed -->
    <FeedList
      :articles="feedStore.articles"
      :is-loading="feedStore.isLoading"
      :is-loading-more="feedStore.isLoadingMore"
      :has-more="feedStore.hasMore"
      @load-more="handleLoadMore"
    />
  </div>
</template>
