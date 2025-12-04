<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { 
  MessageCircle, 
  Bookmark, 
  Share2, 
  MoreHorizontal, 
  Eye, 
  Check,
  ChevronRight,
  Clock
} from 'lucide-vue-next'
import Avatar from '@/components/ui/Avatar.vue'
import Badge from '@/components/ui/Badge.vue'
import Button from '@/components/ui/Button.vue'
import ReactionPanel from '@/components/article/ReactionPanel.vue'
import CommentSection from '@/components/comments/CommentSection.vue'
import ArticleCardSkeleton from '@/components/feed/ArticleCardSkeleton.vue'
import { formatRelativeTime, formatCompactNumber, formatFullDate, formatReadingTime } from '@/utils/formatters'
import type { Article } from '@/types'

const route = useRoute()

const article = ref<Article | null>(null)
const isLoading = ref(true)
const isBookmarked = ref(false)

// Mock article data
const loadArticle = async () => {
  isLoading.value = true
  
  // Simulate API call
  await new Promise(resolve => setTimeout(resolve, 800))
  
  article.value = {
    id: '1',
    title: 'Полное руководство по ChatGPT: от основ до продвинутых техник',
    slug: route.params.slug as string,
    lead: 'Узнайте, как максимально эффективно использовать ChatGPT для работы, учёбы и творчества. Практические советы и проверенные промпты.',
    content: '',
    htmlContent: `
      <h2>Введение</h2>
      <p>ChatGPT изменил наше представление о том, как можно взаимодействовать с искусственным интеллектом. В этом руководстве мы разберём все аспекты работы с этим инструментом — от базовых команд до продвинутых техник промпт-инженерии.</p>
      
      <h2>Что такое ChatGPT?</h2>
      <p>ChatGPT — это большая языковая модель, разработанная компанией OpenAI. Она способна понимать контекст, отвечать на вопросы, писать тексты, код и многое другое.</p>
      
      <div class="prompt-block">
        <p><strong>Пример промпта:</strong></p>
        <p>Ты — опытный копирайтер. Напиши продающий текст для лендинга курсов по программированию. Целевая аудитория: новички 25-35 лет. Тон: дружелюбный и мотивирующий.</p>
      </div>
      
      <h2>Основные принципы работы</h2>
      <p>Чтобы получить качественный результат от ChatGPT, важно соблюдать несколько правил:</p>
      <ul>
        <li>Будьте конкретны в своих запросах</li>
        <li>Указывайте контекст и роль</li>
        <li>Разбивайте сложные задачи на подзадачи</li>
        <li>Используйте примеры желаемого результата</li>
      </ul>
      
      <div class="try-block">
        <p>Попробуйте сами: Откройте ChatGPT и напишите запрос с указанием роли. Сравните результат с обычным запросом без контекста.</p>
      </div>
      
      <h2>Продвинутые техники</h2>
      <p>Для более сложных задач используйте технику цепочки мыслей (Chain of Thought) — просите модель объяснять свои рассуждения пошагово.</p>
      
      <div class="warning-block">
        <p>ChatGPT может генерировать неточную информацию. Всегда проверяйте факты из важных ответов.</p>
      </div>
      
      <h2>Заключение</h2>
      <p>ChatGPT — мощный инструмент, который становится ещё полезнее, когда вы понимаете принципы его работы. Экспериментируйте с промптами и находите свой стиль общения с AI.</p>
    `,
    coverImage: {
      url: 'https://images.unsplash.com/photo-1677442136019-21780ecad995?w=1200',
      width: 1200,
      height: 630,
    },
    level: 'beginner',
    contentType: 'article',
    readingTime: 8,
    isEditorial: true,
    isPinned: false,
    status: 'published',
    commentsEnabled: true,
    isNSFW: false,
    
    author: {
      id: '1',
      username: 'neurogen',
      displayName: 'Редакция Neurogen',
      avatarUrl: undefined,
      isVerified: true,
    },
    
    category: {
      id: 'chatbots',
      name: 'Чат-боты',
      slug: 'chatbots',
      icon: '💬',
      articleCount: 150,
    },
    
    tags: [
      { id: '1', name: 'ChatGPT', slug: 'chatgpt' },
      { id: '2', name: 'промпты', slug: 'prompts' },
      { id: '3', name: 'гайд', slug: 'guide' },
    ],
    
    reactions: [
      { emoji: '👍', count: 234, isReacted: false },
      { emoji: '❤️', count: 89, isReacted: true },
      { emoji: '🔥', count: 56, isReacted: false },
    ],
    
    commentCount: 42,
    viewCount: 12500,
    bookmarkCount: 189,
    
    publishedAt: new Date(Date.now() - 1000 * 60 * 60 * 24 * 3).toISOString(),
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 24 * 5).toISOString(),
    updatedAt: new Date(Date.now() - 1000 * 60 * 60 * 24 * 2).toISOString(),
  }
  
  isLoading.value = false
}

const levelInfo = computed(() => {
  if (!article.value) return null
  const levels = {
    beginner: { emoji: '🟢', label: 'Для новичков', variant: 'beginner' as const },
    intermediate: { emoji: '🟡', label: 'Продвинутое', variant: 'intermediate' as const },
    advanced: { emoji: '🔴', label: 'Для бизнеса', variant: 'advanced' as const },
  }
  return levels[article.value.level]
})

const handleBookmark = () => {
  isBookmarked.value = !isBookmarked.value
}

const handleShare = async () => {
  if (navigator.share) {
    await navigator.share({
      title: article.value?.title,
      url: window.location.href,
    })
  } else {
    await navigator.clipboard.writeText(window.location.href)
    // TODO: Show toast
  }
}

onMounted(loadArticle)
</script>

<template>
  <div class="min-h-screen">
    <!-- Loading -->
    <ArticleCardSkeleton v-if="isLoading" />
    
    <!-- Article -->
    <article v-else-if="article" class="max-w-4xl mx-auto px-4 py-8">
      <!-- Breadcrumbs -->
      <nav class="flex items-center gap-2 text-sm text-text-tertiary mb-6" aria-label="Навигация">
        <RouterLink to="/" class="hover:text-primary transition-colors">Главная</RouterLink>
        <ChevronRight class="w-4 h-4 flex-shrink-0" />
        <RouterLink :to="`/${article.category.slug}`" class="hover:text-primary transition-colors">
          {{ article.category.icon }} {{ article.category.name }}
        </RouterLink>
        <ChevronRight class="w-4 h-4 flex-shrink-0" />
        <span class="text-text-secondary truncate">{{ article.title }}</span>
      </nav>
      
      <!-- Header -->
      <header class="mb-8">
        <!-- Meta badges -->
        <div class="flex items-center gap-2 mb-4 flex-wrap">
          <Badge v-if="levelInfo" :variant="levelInfo.variant">
            {{ levelInfo.emoji }} {{ levelInfo.label }}
          </Badge>
          <Badge v-if="article.isEditorial" variant="primary">
            ✓ Материал редакции
          </Badge>
        </div>
        
        <!-- Title -->
        <h1 class="text-4xl md:text-5xl font-bold text-text-primary dark:text-white mb-5 leading-tight font-display">
          {{ article.title }}
        </h1>
        
        <!-- Lead -->
        <p v-if="article.lead" class="text-xl text-text-secondary mb-6 leading-relaxed max-w-3xl">
          {{ article.lead }}
        </p>
        
        <!-- Author info -->
        <div class="flex items-center gap-4 mb-6 p-4 bg-bg-surface/50 dark:bg-dark-tertiary/30 rounded-xl border border-border-subtle">
          <RouterLink :to="`/@${article.author.username}`" class="flex-shrink-0">
            <Avatar 
              :src="article.author.avatarUrl" 
              :alt="article.author.displayName" 
              :size="56"
            />
          </RouterLink>
          
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 mb-1">
              <RouterLink 
                :to="`/@${article.author.username}`"
                class="font-semibold text-text-primary dark:text-white hover:text-primary transition-colors"
              >
                {{ article.author.displayName }}
              </RouterLink>
              <span 
                v-if="article.author.isVerified" 
                class="w-5 h-5 bg-primary rounded-full flex items-center justify-center flex-shrink-0"
                title="Верифицированный автор"
              >
                <Check class="w-3 h-3 text-white" />
              </span>
            </div>
            <div class="flex items-center gap-4 text-sm text-text-tertiary flex-wrap">
              <time :datetime="article.publishedAt" class="flex items-center gap-1.5">
                <Clock class="w-4 h-4" />
                {{ formatFullDate(article.publishedAt) }}
              </time>
              <span class="flex items-center gap-1.5">
                <Clock class="w-4 h-4" />
                {{ formatReadingTime(article.readingTime) }} чтения
              </span>
              <span class="flex items-center gap-1.5">
                <Eye class="w-4 h-4" />
                {{ formatCompactNumber(article.viewCount) }} просмотров
              </span>
            </div>
          </div>
          
          <Button variant="ghost" size="sm" class="flex-shrink-0">
            Подписаться
          </Button>
        </div>
      </header>
      
      <!-- Cover image -->
      <div v-if="article.coverImage" class="mb-8 rounded-2xl overflow-hidden shadow-elevated">
        <img 
          :src="article.coverImage.url"
          :alt="article.title"
          class="w-full aspect-video object-cover"
          loading="eager"
        />
      </div>
      
      <!-- Content -->
      <div 
        class="prose-article mb-10"
        v-html="article.htmlContent"
      />
      
      <!-- Tags -->
      <div class="flex flex-wrap gap-2 mb-8">
        <RouterLink
          v-for="tag in article.tags"
          :key="tag.id"
          :to="`/tag/${tag.slug}`"
          class="px-4 py-2 text-sm font-medium bg-bg-surface dark:bg-dark-tertiary rounded-full text-text-secondary hover:text-primary hover:bg-primary/10 transition-all duration-200 border border-border-subtle"
        >
          #{{ tag.name }}
        </RouterLink>
      </div>
      
      <!-- Reactions and actions -->
      <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 py-6 px-6 bg-bg-surface/30 dark:bg-dark-tertiary/20 rounded-xl border border-border-subtle mb-10">
        <ReactionPanel 
          :reactions="article.reactions" 
          :article-id="article.id"
        />
        
        <div class="flex items-center gap-2">
          <Button 
            :variant="isBookmarked ? 'primary' : 'ghost'" 
            size="sm"
            @click="handleBookmark"
          >
            <Bookmark class="w-4 h-4" :class="{ 'fill-current': isBookmarked }" />
            <span class="hidden sm:inline ml-1.5">{{ formatCompactNumber(article.bookmarkCount) }}</span>
          </Button>
          
          <Button variant="ghost" size="sm" @click="handleShare">
            <Share2 class="w-4 h-4" />
            <span class="hidden sm:inline ml-1.5">Поделиться</span>
          </Button>
          
          <Button variant="ghost" size="sm">
            <MoreHorizontal class="w-4 h-4" />
          </Button>
        </div>
      </div>
      
      <!-- Comments -->
      <section id="comments" class="mt-12">
        <h2 class="flex items-center gap-3 text-2xl font-bold text-text-primary dark:text-white mb-8 font-display">
          <MessageCircle class="w-6 h-6 text-primary" />
          Комментарии
          <span class="text-text-tertiary font-normal text-lg">({{ article.commentCount }})</span>
        </h2>
        
        <CommentSection :article-id="article.id" />
      </section>
    </article>
  </div>
</template>

