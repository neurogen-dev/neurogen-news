<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Copy, Check, Filter, Search } from 'lucide-vue-next'
import Badge from '@/components/ui/Badge.vue'
import Button from '@/components/ui/Button.vue'

interface Prompt {
  id: string
  title: string
  description: string
  content: string
  category: string
  level: 'beginner' | 'intermediate' | 'advanced'
  copyCount: number
  tags: string[]
}

const prompts = ref<Prompt[]>([])
const isLoading = ref(true)
const copiedId = ref<string | null>(null)
const searchQuery = ref('')
const selectedCategory = ref<string | null>(null)

const categories = [
  { id: 'writing', name: 'Копирайтинг', icon: '✍️' },
  { id: 'coding', name: 'Программирование', icon: '💻' },
  { id: 'learning', name: 'Обучение', icon: '📚' },
  { id: 'creative', name: 'Креатив', icon: '🎨' },
  { id: 'business', name: 'Бизнес', icon: '💼' },
  { id: 'analysis', name: 'Анализ', icon: '📊' },
]

// Mock data
const mockPrompts: Prompt[] = [
  {
    id: '1',
    title: 'Универсальный копирайтер',
    description: 'Промпт для создания продающих текстов любой сложности',
    content: 'Ты — опытный копирайтер с 10-летним стажем. Твоя задача — написать [тип текста] для [продукт/услуга]. Целевая аудитория: [описание ЦА]. Тон коммуникации: [тон]. Текст должен содержать: [ключевые сообщения]. Объём: [количество] символов.',
    category: 'writing',
    level: 'beginner',
    copyCount: 1234,
    tags: ['копирайтинг', 'продажи', 'маркетинг'],
  },
  {
    id: '2',
    title: 'Code Review Assistant',
    description: 'Промпт для детального ревью кода с рекомендациями',
    content: 'Проанализируй следующий код как senior разработчик. Оцени: 1) Читаемость и качество кода 2) Потенциальные баги 3) Производительность 4) Безопасность 5) Соответствие best practices. Дай конкретные рекомендации по улучшению с примерами.',
    category: 'coding',
    level: 'intermediate',
    copyCount: 892,
    tags: ['код', 'ревью', 'программирование'],
  },
  {
    id: '3',
    title: 'Репетитор по любой теме',
    description: 'Превращает ChatGPT в персонального учителя',
    content: 'Ты — опытный преподаватель [предмет]. Объясни [тема] простым языком, как будто я ребёнок 10 лет. Используй аналогии из повседневной жизни. После объяснения задай мне 3 вопроса, чтобы проверить понимание. Если я ошибусь, объясни по-другому.',
    category: 'learning',
    level: 'beginner',
    copyCount: 2156,
    tags: ['обучение', 'образование', 'объяснение'],
  },
  {
    id: '4',
    title: 'Бизнес-аналитик',
    description: 'Анализ рынка и конкурентов',
    content: 'Проведи комплексный анализ рынка [отрасль] в [регион]. Включи: 1) Объём и динамику рынка 2) Ключевых игроков и их доли 3) Тренды и драйверы роста 4) SWOT-анализ для нового игрока 5) Рекомендации по позиционированию. Источники данных укажи в конце.',
    category: 'business',
    level: 'advanced',
    copyCount: 543,
    tags: ['бизнес', 'анализ', 'маркетинг'],
  },
]

const copyPrompt = async (prompt: Prompt) => {
  await navigator.clipboard.writeText(prompt.content)
  copiedId.value = prompt.id
  setTimeout(() => {
    copiedId.value = null
  }, 2000)
}

const filteredPrompts = () => {
  return prompts.value.filter(prompt => {
    if (selectedCategory.value && prompt.category !== selectedCategory.value) {
      return false
    }
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase()
      return (
        prompt.title.toLowerCase().includes(query) ||
        prompt.description.toLowerCase().includes(query) ||
        prompt.tags.some(tag => tag.toLowerCase().includes(query))
      )
    }
    return true
  })
}

const levelInfo = (level: string) => {
  const levels = {
    beginner: { emoji: '🟢', label: 'Простой', variant: 'beginner' as const },
    intermediate: { emoji: '🟡', label: 'Средний', variant: 'intermediate' as const },
    advanced: { emoji: '🔴', label: 'Продвинутый', variant: 'advanced' as const },
  }
  return levels[level as keyof typeof levels]
}

onMounted(() => {
  // Simulate loading
  setTimeout(() => {
    prompts.value = mockPrompts
    isLoading.value = false
  }, 500)
})
</script>

<template>
  <div class="min-h-screen">
    <!-- Header -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-text-primary dark:text-white mb-2">
        💡 Библиотека промптов
      </h1>
      <p class="text-text-secondary">
        Готовые промпты для ChatGPT, Claude и других AI — копируй и используй
      </p>
    </div>
    
    <!-- Filters -->
    <div class="flex flex-col sm:flex-row gap-4 mb-6">
      <!-- Search -->
      <div class="relative flex-1">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-text-tertiary" />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Поиск промптов..."
          class="w-full pl-10 pr-4 py-2 bg-white dark:bg-dark-secondary border border-border dark:border-dark-tertiary rounded-lg text-text-primary dark:text-white"
        />
      </div>
      
      <!-- Category filter -->
      <div class="flex gap-2 overflow-x-auto scrollbar-hide">
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
    </div>
    
    <!-- Prompts grid -->
    <div class="grid gap-4">
      <article
        v-for="prompt in filteredPrompts()"
        :key="prompt.id"
        class="bg-white dark:bg-dark-secondary rounded-xl border border-border dark:border-dark-tertiary p-5"
      >
        <!-- Header -->
        <div class="flex items-start justify-between gap-4 mb-3">
          <div>
            <h2 class="text-lg font-bold text-text-primary dark:text-white mb-1">
              {{ prompt.title }}
            </h2>
            <p class="text-text-secondary text-sm">
              {{ prompt.description }}
            </p>
          </div>
          <Badge :variant="levelInfo(prompt.level).variant">
            {{ levelInfo(prompt.level).emoji }} {{ levelInfo(prompt.level).label }}
          </Badge>
        </div>
        
        <!-- Prompt content -->
        <div class="relative bg-background-secondary dark:bg-dark-tertiary rounded-lg p-4 mb-4">
          <pre class="text-sm text-text-primary dark:text-gray-200 whitespace-pre-wrap font-mono">{{ prompt.content }}</pre>
        </div>
        
        <!-- Footer -->
        <div class="flex items-center justify-between">
          <div class="flex flex-wrap gap-2">
            <Badge 
              v-for="tag in prompt.tags" 
              :key="tag"
              variant="secondary"
            >
              #{{ tag }}
            </Badge>
          </div>
          
          <div class="flex items-center gap-3">
            <span class="text-sm text-text-tertiary">
              Скопировано {{ prompt.copyCount }} раз
            </span>
            <Button
              @click="copyPrompt(prompt)"
              :variant="copiedId === prompt.id ? 'subscribe' : 'secondary'"
              size="sm"
            >
              <Check v-if="copiedId === prompt.id" class="w-4 h-4 mr-1" />
              <Copy v-else class="w-4 h-4 mr-1" />
              {{ copiedId === prompt.id ? 'Скопировано!' : 'Копировать' }}
            </Button>
          </div>
        </div>
      </article>
      
      <!-- Empty state -->
      <div 
        v-if="filteredPrompts().length === 0 && !isLoading"
        class="text-center py-12"
      >
        <div class="text-5xl mb-4">🔍</div>
        <p class="text-text-secondary">
          Промпты не найдены. Попробуйте изменить фильтры.
        </p>
      </div>
    </div>
  </div>
</template>

