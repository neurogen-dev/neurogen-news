<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { Settings, UserPlus } from 'lucide-vue-next'
import FeedList from '@/components/feed/FeedList.vue'
import Button from '@/components/ui/Button.vue'
import { useFeedStore } from '@/stores/feed'
import { useAuthStore } from '@/stores/auth'
import type { ArticleCard } from '@/types'

const feedStore = useFeedStore()
const authStore = useAuthStore()

const articles = ref<ArticleCard[]>([])
const isLoading = ref(true)
const hasSubscriptions = ref(true)

// Load personalized feed
const loadFeed = async () => {
  isLoading.value = true
  try {
    await feedStore.fetchArticles({ sort: 'new' })
    articles.value = feedStore.articles
    // Check if user has subscriptions
    hasSubscriptions.value = authStore.user ? authStore.user.followingCount > 0 : false
  } finally {
    isLoading.value = false
  }
}

const handleLoadMore = async () => {
  await feedStore.loadMore()
}

onMounted(loadFeed)
</script>

<template>
  <div class="min-h-screen">
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-text-primary dark:text-white">
        Моя лента
      </h1>
      <RouterLink 
        to="/settings/feeds"
        class="p-2 text-text-tertiary hover:text-text-primary transition-colors"
        title="Настроить ленту"
      >
        <Settings class="w-5 h-5" />
      </RouterLink>
    </div>
    
    <!-- Empty state - no subscriptions -->
    <template v-if="!hasSubscriptions && !isLoading">
      <div class="text-center py-12 bg-white dark:bg-dark-secondary rounded-xl border border-border dark:border-dark-tertiary">
        <div class="text-5xl mb-4">👋</div>
        <h2 class="text-xl font-bold text-text-primary dark:text-white mb-2">
          Ваша лента пуста
        </h2>
        <p class="text-text-secondary mb-6 max-w-md mx-auto">
          Подпишитесь на авторов и подсайты, чтобы видеть их публикации в своей ленте
        </p>
        
        <div class="flex flex-col sm:flex-row items-center justify-center gap-3">
          <Button as="RouterLink" to="/popular">
            <UserPlus class="w-4 h-4 mr-2" />
            Найти авторов
          </Button>
          <Button as="RouterLink" to="/" variant="secondary">
            Смотреть общую ленту
          </Button>
        </div>
        
        <!-- Suggested authors -->
        <div class="mt-8 pt-8 border-t border-border dark:border-dark-tertiary">
          <h3 class="text-sm font-medium text-text-tertiary mb-4">
            Рекомендуемые авторы
          </h3>
          <div class="flex flex-wrap justify-center gap-2">
            <RouterLink 
              to="/@neurogen"
              class="px-3 py-1.5 text-sm bg-background-secondary dark:bg-dark-tertiary rounded-full text-text-secondary hover:text-primary transition-colors"
            >
              @neurogen
            </RouterLink>
            <RouterLink 
              to="/@ai_expert"
              class="px-3 py-1.5 text-sm bg-background-secondary dark:bg-dark-tertiary rounded-full text-text-secondary hover:text-primary transition-colors"
            >
              @ai_expert
            </RouterLink>
            <RouterLink 
              to="/@promptmaster"
              class="px-3 py-1.5 text-sm bg-background-secondary dark:bg-dark-tertiary rounded-full text-text-secondary hover:text-primary transition-colors"
            >
              @promptmaster
            </RouterLink>
          </div>
        </div>
      </div>
    </template>
    
    <!-- Feed -->
    <template v-else>
      <FeedList
        :articles="articles"
        :is-loading="isLoading"
        :is-loading-more="feedStore.isLoadingMore"
        :has-more="feedStore.hasMore"
        @load-more="handleLoadMore"
      />
    </template>
  </div>
</template>

