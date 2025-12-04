<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { Search, Star, ExternalLink, Filter } from 'lucide-vue-next'
import Badge from '@/components/ui/Badge.vue'
import Button from '@/components/ui/Button.vue'

interface Tool {
  id: string
  name: string
  description: string
  icon: string
  url: string
  category: string
  tags: string[]
  isPremium: boolean
  isFeatured: boolean
  rating: number
  reviewCount: number
}

const tools = ref<Tool[]>([])
const isLoading = ref(true)
const searchQuery = ref('')
const selectedCategory = ref<string | null>(null)

const categories = [
  { id: 'chatbots', name: 'Чат-боты', icon: '💬' },
  { id: 'images', name: 'Изображения', icon: '🎨' },
  { id: 'video', name: 'Видео', icon: '🎬' },
  { id: 'audio', name: 'Аудио', icon: '🎵' },
  { id: 'writing', name: 'Тексты', icon: '✍️' },
  { id: 'coding', name: 'Код', icon: '💻' },
  { id: 'productivity', name: 'Продуктивность', icon: '⚡' },
]

// Mock tools data
const mockTools: Tool[] = [
  {
    id: '1',
    name: 'ChatGPT',
    description: 'Универсальный AI-ассистент от OpenAI для генерации текста, кода и ответов на вопросы',
    icon: '🤖',
    url: 'https://chat.openai.com',
    category: 'chatbots',
    tags: ['текст', 'код', 'анализ'],
    isPremium: false,
    isFeatured: true,
    rating: 4.8,
    reviewCount: 1523,
  },
  {
    id: '2',
    name: 'Midjourney',
    description: 'Генерация высококачественных изображений по текстовому описанию',
    icon: '🎨',
    url: 'https://midjourney.com',
    category: 'images',
    tags: ['изображения', 'арт', 'дизайн'],
    isPremium: true,
    isFeatured: true,
    rating: 4.9,
    reviewCount: 987,
  },
  {
    id: '3',
    name: 'Claude',
    description: 'AI-ассистент от Anthropic с улучшенным пониманием контекста',
    icon: '🧠',
    url: 'https://claude.ai',
    category: 'chatbots',
    tags: ['текст', 'анализ', 'код'],
    isPremium: false,
    isFeatured: true,
    rating: 4.7,
    reviewCount: 654,
  },
  {
    id: '4',
    name: 'Runway',
    description: 'Платформа для создания и редактирования видео с помощью AI',
    icon: '🎬',
    url: 'https://runway.ml',
    category: 'video',
    tags: ['видео', 'монтаж', 'эффекты'],
    isPremium: true,
    isFeatured: false,
    rating: 4.5,
    reviewCount: 342,
  },
  {
    id: '5',
    name: 'Suno',
    description: 'Генерация музыки и песен по текстовому описанию',
    icon: '🎵',
    url: 'https://suno.ai',
    category: 'audio',
    tags: ['музыка', 'аудио', 'песни'],
    isPremium: false,
    isFeatured: true,
    rating: 4.6,
    reviewCount: 521,
  },
  {
    id: '6',
    name: 'GitHub Copilot',
    description: 'AI-помощник для написания кода от GitHub и OpenAI',
    icon: '💻',
    url: 'https://github.com/features/copilot',
    category: 'coding',
    tags: ['код', 'программирование', 'IDE'],
    isPremium: true,
    isFeatured: true,
    rating: 4.7,
    reviewCount: 1102,
  },
]

const filteredTools = computed(() => {
  return tools.value.filter(tool => {
    if (selectedCategory.value && tool.category !== selectedCategory.value) {
      return false
    }
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase()
      return (
        tool.name.toLowerCase().includes(query) ||
        tool.description.toLowerCase().includes(query) ||
        tool.tags.some(tag => tag.toLowerCase().includes(query))
      )
    }
    return true
  })
})

const featuredTools = computed(() => 
  tools.value.filter(tool => tool.isFeatured).slice(0, 4)
)

onMounted(() => {
  setTimeout(() => {
    tools.value = mockTools
    isLoading.value = false
  }, 500)
})
</script>

<template>
  <div class="min-h-screen">
    <!-- Header -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-text-primary dark:text-white mb-2">
        🛠️ Каталог AI-инструментов
      </h1>
      <p class="text-text-secondary">
        Полная коллекция нейросетей и AI-сервисов для любых задач
      </p>
    </div>
    
    <!-- Featured tools -->
    <div v-if="!selectedCategory && !searchQuery" class="mb-8">
      <h2 class="text-lg font-semibold text-text-primary dark:text-white mb-4">
        ⭐ Рекомендуемые
      </h2>
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <a
          v-for="tool in featuredTools"
          :key="tool.id"
          :href="tool.url"
          target="_blank"
          rel="noopener noreferrer"
          class="flex items-start gap-4 p-4 bg-gradient-to-br from-primary/5 to-primary/10 dark:from-primary/10 dark:to-primary/5 rounded-xl border border-primary/20 hover:border-primary/40 transition-colors"
        >
          <div class="text-3xl">{{ tool.icon }}</div>
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2">
              <h3 class="font-bold text-text-primary dark:text-white">
                {{ tool.name }}
              </h3>
              <Badge v-if="tool.isPremium" variant="warning">💎 Premium</Badge>
            </div>
            <p class="text-sm text-text-secondary line-clamp-2 mt-1">
              {{ tool.description }}
            </p>
            <div class="flex items-center gap-2 mt-2">
              <div class="flex items-center gap-1 text-sm">
                <Star class="w-4 h-4 text-yellow-500 fill-yellow-500" />
                <span class="text-text-primary dark:text-white font-medium">
                  {{ tool.rating }}
                </span>
              </div>
              <span class="text-text-tertiary text-sm">
                ({{ tool.reviewCount }} отзывов)
              </span>
            </div>
          </div>
          <ExternalLink class="w-4 h-4 text-text-tertiary shrink-0" />
        </a>
      </div>
    </div>
    
    <!-- Search and filters -->
    <div class="flex flex-col sm:flex-row gap-4 mb-6">
      <div class="relative flex-1">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-text-tertiary" />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Поиск инструментов..."
          class="w-full pl-10 pr-4 py-2 bg-white dark:bg-dark-secondary border border-border dark:border-dark-tertiary rounded-lg text-text-primary dark:text-white"
        />
      </div>
    </div>
    
    <!-- Categories -->
    <div class="flex gap-2 overflow-x-auto scrollbar-hide mb-6">
      <button
        @click="selectedCategory = null"
        class="px-3 py-2 text-sm rounded-lg whitespace-nowrap transition-colors"
        :class="[
          !selectedCategory
            ? 'bg-primary text-white'
            : 'bg-background-secondary dark:bg-dark-tertiary text-text-secondary hover:text-text-primary'
        ]"
      >
        Все
      </button>
      <button
        v-for="category in categories"
        :key="category.id"
        @click="selectedCategory = category.id"
        class="px-3 py-2 text-sm rounded-lg whitespace-nowrap transition-colors"
        :class="[
          selectedCategory === category.id
            ? 'bg-primary text-white'
            : 'bg-background-secondary dark:bg-dark-tertiary text-text-secondary hover:text-text-primary'
        ]"
      >
        {{ category.icon }} {{ category.name }}
      </button>
    </div>
    
    <!-- Tools grid -->
    <div class="grid gap-4">
      <a
        v-for="tool in filteredTools"
        :key="tool.id"
        :href="tool.url"
        target="_blank"
        rel="noopener noreferrer"
        class="flex items-start gap-4 p-5 bg-white dark:bg-dark-secondary rounded-xl border border-border dark:border-dark-tertiary hover:shadow-card-hover transition-shadow"
      >
        <div class="text-4xl">{{ tool.icon }}</div>
        
        <div class="flex-1 min-w-0">
          <div class="flex items-center gap-2 flex-wrap mb-1">
            <h3 class="text-lg font-bold text-text-primary dark:text-white">
              {{ tool.name }}
            </h3>
            <Badge v-if="tool.isPremium" variant="warning">💎 Premium</Badge>
            <Badge v-if="tool.isFeatured" variant="primary">⭐ Рекомендуем</Badge>
          </div>
          
          <p class="text-text-secondary mb-3">
            {{ tool.description }}
          </p>
          
          <div class="flex items-center gap-4 flex-wrap">
            <div class="flex items-center gap-1">
              <Star class="w-4 h-4 text-yellow-500 fill-yellow-500" />
              <span class="text-text-primary dark:text-white font-medium">
                {{ tool.rating }}
              </span>
              <span class="text-text-tertiary text-sm">
                ({{ tool.reviewCount }})
              </span>
            </div>
            
            <div class="flex flex-wrap gap-2">
              <Badge 
                v-for="tag in tool.tags" 
                :key="tag"
                variant="secondary"
              >
                {{ tag }}
              </Badge>
            </div>
          </div>
        </div>
        
        <ExternalLink class="w-5 h-5 text-text-tertiary shrink-0" />
      </a>
      
      <!-- Empty state -->
      <div 
        v-if="filteredTools.length === 0 && !isLoading"
        class="text-center py-12"
      >
        <div class="text-5xl mb-4">🔍</div>
        <p class="text-text-secondary">
          Инструменты не найдены. Попробуйте изменить фильтры.
        </p>
      </div>
    </div>
  </div>
</template>

