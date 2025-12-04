<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Users, FileText, TrendingUp } from 'lucide-vue-next'
import FeedList from '@/components/feed/FeedList.vue'
import FeedFilters from '@/components/feed/FeedFilters.vue'
import Button from '@/components/ui/Button.vue'
import Badge from '@/components/ui/Badge.vue'
import { useFeedStore } from '@/stores/feed'
import { useAuthStore } from '@/stores/auth'
import { formatCompactNumber } from '@/utils/formatters'
import type { ArticleLevel, ContentType, Category } from '@/types'

const route = useRoute()
const feedStore = useFeedStore()
const authStore = useAuthStore()

// Filters
const sort = ref<'popular' | 'new'>('popular')
const level = ref<ArticleLevel | undefined>()
const contentType = ref<ContentType | undefined>()

// Category data
const category = ref<Category | null>(null)
const isSubscribed = ref(false)

// Mock category data based on route
const categoryData: Record<string, Omit<Category, 'id' | 'articleCount'>> = {
  chatbots: {
    name: 'Чат-боты',
    slug: 'chatbots',
    icon: '💬',
    description: 'Всё о ChatGPT, Claude, Gemini и других AI-ассистентах. Гайды, промпты, новости и практические советы.',
  },
  images: {
    name: 'Изображения',
    slug: 'images',
    icon: '🎨',
    description: 'Генерация изображений с помощью Midjourney, DALL-E, Stable Diffusion. Обзоры инструментов и творческие техники.',
  },
  video: {
    name: 'Видео',
    slug: 'video',
    icon: '🎬',
    description: 'AI для создания и редактирования видео. Runway, Sora, Pika и другие инструменты нового поколения.',
  },
  music: {
    name: 'Музыка',
    slug: 'music',
    icon: '🎵',
    description: 'Генерация музыки и звуков с помощью AI. Suno, Udio и другие музыкальные нейросети.',
  },
  text: {
    name: 'Текст',
    slug: 'text',
    icon: '✍️',
    description: 'Копирайтинг, рерайтинг и работа с текстом. AI-инструменты для писателей и контент-мейкеров.',
  },
  code: {
    name: 'Код',
    slug: 'code',
    icon: '💻',
    description: 'AI для программистов. GitHub Copilot, Cursor, Cody и другие инструменты для разработки.',
  },
}

const categorySlug = computed(() => route.meta.category as string || route.params.category as string)

// Load category and articles
const loadData = async () => {
  const slug = categorySlug.value
  const data = categoryData[slug]
  
  if (data) {
    category.value = {
      ...data,
      id: slug,
      articleCount: Math.floor(Math.random() * 500) + 50, // Mock
    }
  }
  
  await feedStore.fetchArticles({
    sort: sort.value,
    level: level.value,
    contentType: contentType.value,
    categoryId: slug,
  })
}

const handleSubscribe = () => {
  if (!authStore.isLoggedIn) {
    // TODO: Show login modal
    return
  }
  isSubscribed.value = !isSubscribed.value
}

// Watch for filter changes
watch([sort, level, contentType], loadData)
watch(categorySlug, loadData)

onMounted(loadData)
</script>

<template>
  <div class="min-h-screen">
    <!-- Category header -->
    <div 
      v-if="category"
      class="bg-white dark:bg-dark-secondary rounded-xl border border-border dark:border-dark-tertiary p-6 mb-6"
    >
      <div class="flex items-start gap-4">
        <!-- Icon -->
        <div class="text-5xl">{{ category.icon }}</div>
        
        <!-- Info -->
        <div class="flex-1">
          <div class="flex items-center gap-3 mb-2">
            <h1 class="text-2xl font-bold text-text-primary dark:text-white">
              {{ category.name }}
            </h1>
            <Badge variant="secondary">
              {{ formatCompactNumber(category.articleCount) }} статей
            </Badge>
          </div>
          
          <p class="text-text-secondary mb-4">
            {{ category.description }}
          </p>
          
          <!-- Stats and actions -->
          <div class="flex items-center gap-4">
            <div class="flex items-center gap-6 text-sm text-text-tertiary">
              <span class="flex items-center gap-1.5">
                <FileText class="w-4 h-4" />
                {{ formatCompactNumber(category.articleCount) }} публикаций
              </span>
              <span class="flex items-center gap-1.5">
                <Users class="w-4 h-4" />
                {{ formatCompactNumber(Math.floor(Math.random() * 10000)) }} подписчиков
              </span>
              <span class="flex items-center gap-1.5">
                <TrendingUp class="w-4 h-4" />
                +{{ Math.floor(Math.random() * 100) }} сегодня
              </span>
            </div>
            
            <div class="ml-auto">
              <Button 
                :variant="isSubscribed ? 'secondary' : 'subscribe'"
                @click="handleSubscribe"
              >
                {{ isSubscribed ? '✓ Подписка' : 'Подписаться' }}
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Filters -->
    <FeedFilters
      v-model:sort="sort"
      v-model:level="level"
      v-model:content-type="contentType"
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

