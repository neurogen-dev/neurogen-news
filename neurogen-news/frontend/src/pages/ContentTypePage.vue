<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import FeedList from '@/components/feed/FeedList.vue'
import FeedFilters from '@/components/feed/FeedFilters.vue'
import { useFeedStore } from '@/stores/feed'
import type { ArticleLevel, ContentType } from '@/types'

const route = useRoute()
const feedStore = useFeedStore()

// Filters
const sort = ref<'popular' | 'new'>('popular')
const level = ref<ArticleLevel | undefined>()

// Content type from route
const contentTypeSlug = computed(() => route.meta.contentType as string)

const contentTypeInfo: Record<string, { title: string; description: string; emoji: string }> = {
  guides: {
    title: 'Гайды',
    description: 'Подробные инструкции и руководства по работе с нейросетями',
    emoji: '📚',
  },
  reviews: {
    title: 'Обзоры',
    description: 'Обзоры и сравнения AI-инструментов и сервисов',
    emoji: '🔍',
  },
  news: {
    title: 'Новости',
    description: 'Последние новости мира искусственного интеллекта',
    emoji: '📰',
  },
}

const currentInfo = computed(() => contentTypeInfo[contentTypeSlug.value])

// Load articles
const loadData = async () => {
  await feedStore.fetchArticles({
    sort: sort.value,
    level: level.value,
    contentType: contentTypeSlug.value as ContentType,
  })
}

// Watch for filter changes
watch([sort, level], loadData)
watch(contentTypeSlug, loadData)

onMounted(loadData)
</script>

<template>
  <div class="min-h-screen">
    <!-- Header -->
    <div class="mb-6">
      <div class="flex items-center gap-3 mb-2">
        <span class="text-3xl">{{ currentInfo?.emoji }}</span>
        <h1 class="text-2xl font-bold text-text-primary dark:text-white">
          {{ currentInfo?.title }}
        </h1>
      </div>
      <p class="text-text-secondary">
        {{ currentInfo?.description }}
      </p>
    </div>
    
    <!-- Filters -->
    <FeedFilters
      v-model:sort="sort"
      v-model:level="level"
      class="mb-6"
    />
    
    <!-- Articles -->
    <FeedList
      :articles="feedStore.articles"
      :is-loading="feedStore.isLoading"
      :is-loading-more="feedStore.isLoadingMore"
      :has-more="feedStore.hasMore"
      @load-more="feedStore.loadMore"
    />
  </div>
</template>

