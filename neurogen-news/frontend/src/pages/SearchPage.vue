<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Search, Filter, X } from 'lucide-vue-next'
import FeedList from '@/components/feed/FeedList.vue'
import Button from '@/components/ui/Button.vue'
import Badge from '@/components/ui/Badge.vue'
import { useFeedStore } from '@/stores/feed'
import type { ArticleLevel, ContentType } from '@/types'

const route = useRoute()
const router = useRouter()
const feedStore = useFeedStore()

const query = ref((route.query.q as string) || '')
const showFilters = ref(false)

// Filters
const level = ref<ArticleLevel | undefined>()
const contentType = ref<ContentType | undefined>()
const timeRange = ref<'24h' | '7d' | '30d' | 'all'>('all')

const hasFilters = computed(() => 
  level.value || contentType.value || timeRange.value !== 'all'
)

const activeFilterCount = computed(() => {
  let count = 0
  if (level.value) count++
  if (contentType.value) count++
  if (timeRange.value !== 'all') count++
  return count
})

// Perform search
const performSearch = async () => {
  if (!query.value.trim()) {
    feedStore.reset()
    return
  }
  
  // Update URL
  router.replace({ query: { q: query.value } })
  
  // TODO: Call search API instead of regular feed
  await feedStore.fetchArticles({
    sort: 'popular',
    level: level.value,
    contentType: contentType.value,
    timeRange: timeRange.value,
  })
}

const clearFilters = () => {
  level.value = undefined
  contentType.value = undefined
  timeRange.value = 'all'
}

// Watch for query changes
watch(query, performSearch)
watch([level, contentType, timeRange], performSearch)

// Watch route query
watch(
  () => route.query.q,
  (newQuery) => {
    if (newQuery !== query.value) {
      query.value = (newQuery as string) || ''
    }
  }
)

onMounted(() => {
  if (query.value) {
    performSearch()
  }
})
</script>

<template>
  <div class="min-h-screen">
    <!-- Search header -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-text-primary dark:text-white mb-4">
        Поиск
      </h1>
      
      <!-- Search input -->
      <div class="flex gap-2">
        <div class="relative flex-1">
          <Search class="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-text-tertiary" />
          <input
            v-model="query"
            type="text"
            placeholder="Поиск статей, инструментов..."
            class="w-full pl-12 pr-4 py-3 bg-white dark:bg-dark-secondary border border-border dark:border-dark-tertiary rounded-xl text-text-primary dark:text-white placeholder:text-text-tertiary focus:outline-none focus:ring-2 focus:ring-primary/50"
            @keydown.enter="performSearch"
          />
          <button
            v-if="query"
            @click="query = ''"
            class="absolute right-4 top-1/2 -translate-y-1/2 text-text-tertiary hover:text-text-primary"
          >
            <X class="w-5 h-5" />
          </button>
        </div>
        
        <Button
          :variant="showFilters ? 'primary' : 'secondary'"
          @click="showFilters = !showFilters"
          class="relative"
        >
          <Filter class="w-5 h-5" />
          <span 
            v-if="activeFilterCount > 0"
            class="absolute -top-1 -right-1 w-5 h-5 bg-primary text-white text-xs rounded-full flex items-center justify-center"
          >
            {{ activeFilterCount }}
          </span>
        </Button>
      </div>
      
      <!-- Filters panel -->
      <div 
        v-if="showFilters"
        class="mt-4 p-4 bg-white dark:bg-dark-secondary rounded-xl border border-border dark:border-dark-tertiary"
      >
        <div class="flex items-center justify-between mb-4">
          <h3 class="font-medium text-text-primary dark:text-white">
            Фильтры
          </h3>
          <button
            v-if="hasFilters"
            @click="clearFilters"
            class="text-sm text-primary hover:underline"
          >
            Сбросить
          </button>
        </div>
        
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <!-- Level filter -->
          <div>
            <label class="block text-sm text-text-tertiary mb-2">Уровень</label>
            <select
              v-model="level"
              class="w-full px-3 py-2 bg-background-secondary dark:bg-dark-tertiary border border-border dark:border-dark-tertiary rounded-lg text-text-primary dark:text-white"
            >
              <option :value="undefined">Все</option>
              <option value="beginner">🟢 Для новичков</option>
              <option value="intermediate">🟡 Продвинутое</option>
              <option value="advanced">🔴 Для бизнеса</option>
            </select>
          </div>
          
          <!-- Content type filter -->
          <div>
            <label class="block text-sm text-text-tertiary mb-2">Тип контента</label>
            <select
              v-model="contentType"
              class="w-full px-3 py-2 bg-background-secondary dark:bg-dark-tertiary border border-border dark:border-dark-tertiary rounded-lg text-text-primary dark:text-white"
            >
              <option :value="undefined">Все</option>
              <option value="article">📖 Статьи</option>
              <option value="news">📰 Новости</option>
              <option value="question">❓ Вопросы</option>
            </select>
          </div>
          
          <!-- Time range filter -->
          <div>
            <label class="block text-sm text-text-tertiary mb-2">Период</label>
            <select
              v-model="timeRange"
              class="w-full px-3 py-2 bg-background-secondary dark:bg-dark-tertiary border border-border dark:border-dark-tertiary rounded-lg text-text-primary dark:text-white"
            >
              <option value="all">За всё время</option>
              <option value="24h">За 24 часа</option>
              <option value="7d">За неделю</option>
              <option value="30d">За месяц</option>
            </select>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Search results -->
    <template v-if="query">
      <!-- Results count -->
      <div class="mb-4 text-text-secondary">
        <template v-if="!feedStore.isLoading">
          Найдено: <span class="font-medium text-text-primary dark:text-white">{{ feedStore.pagination.total }}</span> результатов
          <span v-if="hasFilters">
            с фильтрами
            <Badge v-if="level" variant="primary" class="ml-1">
              {{ level }}
            </Badge>
            <Badge v-if="contentType" variant="primary" class="ml-1">
              {{ contentType }}
            </Badge>
          </span>
        </template>
      </div>
      
      <FeedList
        :articles="feedStore.articles"
        :is-loading="feedStore.isLoading"
        :is-loading-more="feedStore.isLoadingMore"
        :has-more="feedStore.hasMore"
        @load-more="feedStore.loadMore"
      />
    </template>
    
    <!-- Empty state -->
    <template v-else>
      <div class="text-center py-12">
        <Search class="w-16 h-16 text-text-tertiary mx-auto mb-4" />
        <h2 class="text-xl font-medium text-text-primary dark:text-white mb-2">
          Введите поисковый запрос
        </h2>
        <p class="text-text-secondary">
          Найдите статьи, инструменты и авторов
        </p>
      </div>
    </template>
  </div>
</template>

