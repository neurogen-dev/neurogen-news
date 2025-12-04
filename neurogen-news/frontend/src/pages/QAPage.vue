<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { MessageCircle, Check, Clock, Filter, PenSquare, ChevronUp } from 'lucide-vue-next'
import Avatar from '@/components/ui/Avatar.vue'
import Badge from '@/components/ui/Badge.vue'
import Button from '@/components/ui/Button.vue'
import { formatRelativeTime } from '@/utils/formatters'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

interface Question {
  id: string
  title: string
  content: string
  author: {
    username: string
    displayName: string
    avatarUrl?: string
  }
  answerCount: number
  viewCount: number
  voteCount: number
  hasAcceptedAnswer: boolean
  tags: string[]
  createdAt: string
}

const questions = ref<Question[]>([])
const isLoading = ref(true)
const sortBy = ref<'new' | 'popular' | 'unanswered'>('new')

// Mock data
const mockQuestions: Question[] = [
  {
    id: '1',
    title: 'Как правильно формулировать промпты для генерации кода в ChatGPT?',
    content: 'Пытаюсь использовать ChatGPT для написания кода, но результаты не всегда соответствуют ожиданиям...',
    author: {
      username: 'newbie_dev',
      displayName: 'Начинающий разработчик',
    },
    answerCount: 5,
    viewCount: 234,
    voteCount: 12,
    hasAcceptedAnswer: true,
    tags: ['chatgpt', 'промпты', 'код'],
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 2).toISOString(),
  },
  {
    id: '2',
    title: 'Midjourney vs DALL-E 3 — что лучше для коммерческих проектов?',
    content: 'Нужно выбрать инструмент для генерации изображений для коммерческого использования...',
    author: {
      username: 'designer_pro',
      displayName: 'Дизайнер',
    },
    answerCount: 8,
    viewCount: 456,
    voteCount: 24,
    hasAcceptedAnswer: true,
    tags: ['midjourney', 'dall-e', 'сравнение'],
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 24).toISOString(),
  },
  {
    id: '3',
    title: 'Почему Claude отказывается выполнять некоторые запросы?',
    content: 'Заметил, что Claude иногда отказывается помогать с определёнными задачами, хотя они кажутся безобидными...',
    author: {
      username: 'curious_user',
      displayName: 'Любопытный пользователь',
    },
    answerCount: 3,
    viewCount: 189,
    voteCount: 8,
    hasAcceptedAnswer: false,
    tags: ['claude', 'ограничения', 'ai'],
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 48).toISOString(),
  },
  {
    id: '4',
    title: 'Как обучить свою нейросеть на собственных данных?',
    content: 'Хочу создать модель, которая будет специализироваться на моей предметной области...',
    author: {
      username: 'ml_enthusiast',
      displayName: 'ML Энтузиаст',
    },
    answerCount: 0,
    viewCount: 67,
    voteCount: 5,
    hasAcceptedAnswer: false,
    tags: ['машинное обучение', 'fine-tuning', 'данные'],
    createdAt: new Date(Date.now() - 1000 * 60 * 30).toISOString(),
  },
]

onMounted(() => {
  setTimeout(() => {
    questions.value = mockQuestions
    isLoading.value = false
  }, 500)
})
</script>

<template>
  <div class="min-h-screen">
    <!-- Header -->
    <div class="flex items-start justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-text-primary dark:text-white mb-2">
          ❓ Вопросы и ответы
        </h1>
        <p class="text-text-secondary">
          Задавайте вопросы о нейросетях и получайте ответы от сообщества
        </p>
      </div>
      
      <Button 
        v-if="authStore.isLoggedIn"
        as="RouterLink" 
        to="/editor/new?type=question"
      >
        <PenSquare class="w-4 h-4 mr-1" />
        Задать вопрос
      </Button>
    </div>
    
    <!-- Filters -->
    <div class="flex gap-2 mb-6">
      <button
        @click="sortBy = 'new'"
        class="px-4 py-2 text-sm rounded-lg transition-colors"
        :class="[
          sortBy === 'new'
            ? 'bg-primary text-white'
            : 'bg-background-secondary dark:bg-dark-tertiary text-text-secondary hover:text-text-primary'
        ]"
      >
        <Clock class="w-4 h-4 inline mr-1" />
        Новые
      </button>
      <button
        @click="sortBy = 'popular'"
        class="px-4 py-2 text-sm rounded-lg transition-colors"
        :class="[
          sortBy === 'popular'
            ? 'bg-primary text-white'
            : 'bg-background-secondary dark:bg-dark-tertiary text-text-secondary hover:text-text-primary'
        ]"
      >
        <ChevronUp class="w-4 h-4 inline mr-1" />
        Популярные
      </button>
      <button
        @click="sortBy = 'unanswered'"
        class="px-4 py-2 text-sm rounded-lg transition-colors"
        :class="[
          sortBy === 'unanswered'
            ? 'bg-primary text-white'
            : 'bg-background-secondary dark:bg-dark-tertiary text-text-secondary hover:text-text-primary'
        ]"
      >
        <MessageCircle class="w-4 h-4 inline mr-1" />
        Без ответа
      </button>
    </div>
    
    <!-- Questions list -->
    <div class="space-y-4">
      <article
        v-for="question in questions"
        :key="question.id"
        class="bg-white dark:bg-dark-secondary rounded-xl border border-border dark:border-dark-tertiary p-5"
      >
        <div class="flex gap-4">
          <!-- Stats -->
          <div class="hidden sm:flex flex-col items-center gap-2 text-center min-w-[80px]">
            <div class="text-lg font-bold text-text-primary dark:text-white">
              {{ question.voteCount }}
            </div>
            <div class="text-xs text-text-tertiary">голосов</div>
            
            <div 
              class="px-2 py-1 rounded text-sm"
              :class="[
                question.hasAcceptedAnswer
                  ? 'bg-success/10 text-success'
                  : question.answerCount > 0
                    ? 'bg-background-tertiary dark:bg-dark-tertiary text-text-secondary'
                    : 'border border-border dark:border-dark-tertiary text-text-tertiary'
              ]"
            >
              <Check v-if="question.hasAcceptedAnswer" class="w-4 h-4 inline" />
              {{ question.answerCount }}
              <span class="block text-xs">ответов</span>
            </div>
          </div>
          
          <!-- Content -->
          <div class="flex-1 min-w-0">
            <RouterLink 
              :to="`/qa/${question.id}`"
              class="text-lg font-bold text-text-primary dark:text-white hover:text-primary transition-colors line-clamp-2"
            >
              {{ question.title }}
            </RouterLink>
            
            <p class="text-text-secondary text-sm mt-2 line-clamp-2">
              {{ question.content }}
            </p>
            
            <!-- Tags and meta -->
            <div class="flex items-center justify-between mt-4 flex-wrap gap-2">
              <div class="flex flex-wrap gap-2">
                <Badge 
                  v-for="tag in question.tags" 
                  :key="tag"
                  variant="secondary"
                >
                  {{ tag }}
                </Badge>
              </div>
              
              <div class="flex items-center gap-3 text-sm text-text-tertiary">
                <span>{{ question.viewCount }} просмотров</span>
                <RouterLink 
                  :to="`/@${question.author.username}`"
                  class="flex items-center gap-2 hover:text-primary"
                >
                  <Avatar 
                    :src="question.author.avatarUrl"
                    :alt="question.author.displayName"
                    :size="20"
                  />
                  {{ question.author.displayName }}
                </RouterLink>
                <time>{{ formatRelativeTime(question.createdAt) }}</time>
              </div>
            </div>
          </div>
        </div>
      </article>
      
      <!-- Empty state -->
      <div 
        v-if="questions.length === 0 && !isLoading"
        class="text-center py-12"
      >
        <MessageCircle class="w-16 h-16 text-text-tertiary mx-auto mb-4" />
        <h2 class="text-xl font-medium text-text-primary dark:text-white mb-2">
          Пока нет вопросов
        </h2>
        <p class="text-text-secondary mb-6">
          Будьте первым, кто задаст вопрос!
        </p>
        <Button as="RouterLink" to="/editor/new?type=question">
          Задать вопрос
        </Button>
      </div>
    </div>
  </div>
</template>

