<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { 
  ArrowRight, 
  BookOpen, 
  Lightbulb, 
  Wrench, 
  MessageCircle,
  GraduationCap,
  Rocket,
  Star
} from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'

const learningPaths = [
  {
    id: 'basics',
    title: 'Основы AI',
    description: 'Что такое нейросети и как они работают',
    icon: '🧠',
    lessons: 5,
    duration: '30 мин',
    level: 'beginner',
  },
  {
    id: 'chatgpt',
    title: 'Мастер ChatGPT',
    description: 'Научитесь получать максимум от AI-ассистентов',
    icon: '💬',
    lessons: 8,
    duration: '1 час',
    level: 'beginner',
  },
  {
    id: 'images',
    title: 'Генерация изображений',
    description: 'Midjourney, DALL-E и Stable Diffusion',
    icon: '🎨',
    lessons: 6,
    duration: '45 мин',
    level: 'beginner',
  },
  {
    id: 'prompts',
    title: 'Промпт-инженерия',
    description: 'Продвинутые техники составления запросов',
    icon: '✨',
    lessons: 10,
    duration: '1.5 часа',
    level: 'intermediate',
  },
]

const quickGuides = [
  {
    title: 'Первые шаги с ChatGPT',
    description: 'Создайте аккаунт и напишите первый запрос',
    link: '/start/chatgpt-first-steps',
    emoji: '👋',
  },
  {
    title: '10 лучших промптов',
    description: 'Готовые запросы для повседневных задач',
    link: '/prompts',
    emoji: '💡',
  },
  {
    title: 'Сравнение AI-моделей',
    description: 'ChatGPT vs Claude vs Gemini',
    link: '/reviews/ai-models-comparison',
    emoji: '⚖️',
  },
  {
    title: 'AI для работы',
    description: 'Как использовать нейросети в профессии',
    link: '/guides/ai-for-work',
    emoji: '💼',
  },
]
</script>

<template>
  <div class="min-h-screen">
    <!-- Hero section -->
    <div class="text-center mb-12">
      <div class="text-6xl mb-4">🎓</div>
      <h1 class="text-3xl font-bold text-text-primary dark:text-white mb-4">
        Добро пожаловать в мир AI!
      </h1>
      <p class="text-lg text-text-secondary max-w-2xl mx-auto">
        Мы поможем вам освоить нейросети с нуля. Выберите путь обучения или начните с быстрых гайдов.
      </p>
    </div>
    
    <!-- Learning paths -->
    <section class="mb-12">
      <h2 class="text-xl font-bold text-text-primary dark:text-white mb-6">
        📚 Пути обучения
      </h2>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <RouterLink
          v-for="path in learningPaths"
          :key="path.id"
          :to="`/start/${path.id}`"
          class="group flex items-start gap-4 p-5 bg-white dark:bg-dark-secondary rounded-xl border border-border dark:border-dark-tertiary hover:border-primary hover:shadow-card-hover transition-all"
        >
          <div class="text-4xl">{{ path.icon }}</div>
          
          <div class="flex-1">
            <h3 class="font-bold text-text-primary dark:text-white group-hover:text-primary transition-colors mb-1">
              {{ path.title }}
            </h3>
            <p class="text-text-secondary text-sm mb-3">
              {{ path.description }}
            </p>
            <div class="flex items-center gap-4 text-sm text-text-tertiary">
              <span>{{ path.lessons }} уроков</span>
              <span>{{ path.duration }}</span>
              <span 
                class="px-2 py-0.5 rounded-full"
                :class="{
                  'bg-beginner/10 text-beginner': path.level === 'beginner',
                  'bg-intermediate/10 text-intermediate': path.level === 'intermediate',
                }"
              >
                {{ path.level === 'beginner' ? '🟢 Базовый' : '🟡 Средний' }}
              </span>
            </div>
          </div>
          
          <ArrowRight class="w-5 h-5 text-text-tertiary group-hover:text-primary group-hover:translate-x-1 transition-all" />
        </RouterLink>
      </div>
    </section>
    
    <!-- Quick guides -->
    <section class="mb-12">
      <h2 class="text-xl font-bold text-text-primary dark:text-white mb-6">
        ⚡ Быстрый старт
      </h2>
      
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <RouterLink
          v-for="guide in quickGuides"
          :key="guide.title"
          :to="guide.link"
          class="group flex items-center gap-4 p-4 bg-background-secondary dark:bg-dark-tertiary rounded-xl hover:bg-primary/5 dark:hover:bg-primary/10 transition-colors"
        >
          <div class="text-2xl">{{ guide.emoji }}</div>
          <div class="flex-1">
            <h3 class="font-medium text-text-primary dark:text-white group-hover:text-primary transition-colors">
              {{ guide.title }}
            </h3>
            <p class="text-sm text-text-tertiary">
              {{ guide.description }}
            </p>
          </div>
          <ArrowRight class="w-4 h-4 text-text-tertiary group-hover:text-primary" />
        </RouterLink>
      </div>
    </section>
    
    <!-- Resources -->
    <section class="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-12">
      <RouterLink
        to="/prompts"
        class="p-6 bg-gradient-to-br from-amber-500/10 to-orange-500/10 rounded-xl border border-amber-500/20 hover:border-amber-500/40 transition-colors text-center"
      >
        <Lightbulb class="w-8 h-8 text-amber-500 mx-auto mb-3" />
        <h3 class="font-bold text-text-primary dark:text-white mb-1">
          Библиотека промптов
        </h3>
        <p class="text-sm text-text-secondary">
          Готовые запросы для любых задач
        </p>
      </RouterLink>
      
      <RouterLink
        to="/tools"
        class="p-6 bg-gradient-to-br from-blue-500/10 to-cyan-500/10 rounded-xl border border-blue-500/20 hover:border-blue-500/40 transition-colors text-center"
      >
        <Wrench class="w-8 h-8 text-blue-500 mx-auto mb-3" />
        <h3 class="font-bold text-text-primary dark:text-white mb-1">
          Каталог инструментов
        </h3>
        <p class="text-sm text-text-secondary">
          Все AI-сервисы в одном месте
        </p>
      </RouterLink>
      
      <RouterLink
        to="/qa"
        class="p-6 bg-gradient-to-br from-purple-500/10 to-pink-500/10 rounded-xl border border-purple-500/20 hover:border-purple-500/40 transition-colors text-center"
      >
        <MessageCircle class="w-8 h-8 text-purple-500 mx-auto mb-3" />
        <h3 class="font-bold text-text-primary dark:text-white mb-1">
          Вопросы и ответы
        </h3>
        <p class="text-sm text-text-secondary">
          Задайте вопрос сообществу
        </p>
      </RouterLink>
    </section>
    
    <!-- CTA -->
    <div class="text-center py-8 px-6 bg-gradient-to-r from-primary/5 via-primary/10 to-primary/5 rounded-2xl border border-primary/20">
      <Star class="w-10 h-10 text-primary mx-auto mb-4" />
      <h2 class="text-xl font-bold text-text-primary dark:text-white mb-2">
        Готовы начать?
      </h2>
      <p class="text-text-secondary mb-6">
        Создайте аккаунт, чтобы сохранять прогресс и получать персональные рекомендации
      </p>
      <div class="flex justify-center gap-3">
        <Button as="RouterLink" to="/register" size="lg">
          Создать аккаунт
        </Button>
        <Button as="RouterLink" to="/popular" variant="secondary" size="lg">
          Смотреть контент
        </Button>
      </div>
    </div>
  </div>
</template>

