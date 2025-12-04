<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { 
  MessageCircle, 
  Bookmark, 
  Share2, 
  MoreHorizontal, 
  Eye, 
  Check, 
  Heart,
  BookOpen,
  Newspaper,
  FileText,
  HelpCircle,
  MessageSquare,
  CheckCircle2
} from 'lucide-vue-next'
import Avatar from '@/components/ui/Avatar.vue'
import Badge from '@/components/ui/Badge.vue'
import Button from '@/components/ui/Button.vue'
import ReactionPanel from '@/components/article/ReactionPanel.vue'
import { formatRelativeTime, formatCompactNumber, formatReadingTime } from '@/utils/formatters'
import type { ArticleCard as ArticleCardType } from '@/types'

interface Props {
  article: ArticleCardType
  compact?: boolean
  index?: number
}

const props = withDefaults(defineProps<Props>(), {
  compact: false,
  index: 0,
})

const cardRef = ref<HTMLElement | null>(null)
const isVisible = ref(false)
const mousePos = ref({ x: 50, y: 50 })
const isLiked = ref(false)
const isBookmarked = ref(false)

const levelInfo = computed(() => {
  const levels = {
    beginner: { 
      icon: CheckCircle2, 
      label: 'Для новичков', 
      variant: 'beginner' as const,
      color: 'beginner'
    },
    intermediate: { 
      icon: CheckCircle2, 
      label: 'Продвинутое', 
      variant: 'intermediate' as const,
      color: 'intermediate'
    },
    advanced: { 
      icon: CheckCircle2, 
      label: 'Для бизнеса', 
      variant: 'advanced' as const,
      color: 'advanced'
    },
  }
  return levels[props.article.level]
})

const contentTypeInfo = computed(() => {
  const types = {
    article: { icon: BookOpen, label: 'Статья' },
    news: { icon: Newspaper, label: 'Новость' },
    post: { icon: FileText, label: 'Пост' },
    question: { icon: HelpCircle, label: 'Вопрос' },
    discussion: { icon: MessageSquare, label: 'Обсуждение' },
  }
  return types[props.article.contentType] || types.article
})

const articleUrl = computed(() => 
  `/${props.article.category.slug}/${props.article.slug}`
)

// Mouse tracking for glow effect
const handleMouseMove = (e: MouseEvent) => {
  if (!cardRef.value) return
  const rect = cardRef.value.getBoundingClientRect()
  mousePos.value = {
    x: ((e.clientX - rect.left) / rect.width) * 100,
    y: ((e.clientY - rect.top) / rect.height) * 100,
  }
}

// Intersection observer for entrance animation
onMounted(() => {
  if (!cardRef.value) return
  
  const observer = new IntersectionObserver(
    ([entry]) => {
      if (entry?.isIntersecting) {
        isVisible.value = true
        observer.disconnect()
      }
    },
    { threshold: 0.1 }
  )
  
  observer.observe(cardRef.value)
})

const toggleLike = () => {
  isLiked.value = !isLiked.value
}

const toggleBookmark = () => {
  isBookmarked.value = !isBookmarked.value
}
</script>

<template>
  <article 
    ref="cardRef"
    @mousemove="handleMouseMove"
    class="empatra-card group relative"
    :class="[
      { 'px-5 py-5': compact, 'p-6': !compact },
      isVisible ? 'animate-slide-up' : 'opacity-0'
    ]"
    :style="{
      '--mouse-x': `${mousePos.x}%`,
      '--mouse-y': `${mousePos.y}%`,
      animationDelay: `${index * 0.05}s`
    }"
  >
    <!-- Header -->
    <header class="flex items-center gap-3 mb-4">
      <RouterLink :to="`/@${article.author.username}`" class="scale-hover">
        <Avatar 
          :src="article.author.avatarUrl" 
          :alt="article.author.displayName" 
          size="sm"
        />
      </RouterLink>
      
      <div class="flex-1 min-w-0">
        <div class="flex items-center gap-2 text-sm flex-wrap">
          <RouterLink 
            :to="`/@${article.author.username}`"
            class="font-semibold text-text-primary hover:text-primary transition-colors duration-200"
          >
            {{ article.author.displayName }}
          </RouterLink>
          
          <span 
            v-if="article.author.isVerified" 
            class="w-4 h-4 bg-gradient-to-br from-primary to-secondary rounded-full flex items-center justify-center"
            title="Верифицированный аккаунт"
          >
            <Check class="w-2.5 h-2.5 text-white" stroke-width="3" />
          </span>
          
          <span class="text-text-tertiary">•</span>
          
          <RouterLink 
            :to="`/${article.category.slug}`"
            class="text-text-secondary hover:text-primary transition-colors duration-200 flex items-center gap-1.5"
          >
            <component 
              :is="contentTypeInfo.icon" 
              class="w-3.5 h-3.5"
              stroke-width="2"
            />
            <span>{{ article.category.name }}</span>
          </RouterLink>
          
          <span class="text-text-tertiary">•</span>
          
          <time 
            :datetime="article.publishedAt"
            class="text-text-tertiary"
          >
            {{ formatRelativeTime(article.publishedAt) }}
          </time>
        </div>
      </div>
      
      <Button 
        variant="subscribe" 
        size="sm" 
        class="hidden sm:flex"
      >
        Подписаться
      </Button>
    </header>
    
    <!-- Meta badges -->
    <div class="flex items-center gap-2.5 mb-4 flex-wrap">
      <Badge 
        :variant="levelInfo.variant" 
        shimmer
        class="flex items-center gap-1.5"
      >
        <component 
          :is="levelInfo.icon" 
          class="w-3.5 h-3.5"
          :class="`text-${levelInfo.color}`"
          stroke-width="2.5"
          fill="currentColor"
          fill-opacity="0.2"
        />
        <span>{{ levelInfo.label }}</span>
      </Badge>
      <Badge 
        variant="secondary"
        class="flex items-center gap-1.5"
      >
        <component 
          :is="contentTypeInfo.icon" 
          class="w-3.5 h-3.5"
          stroke-width="2"
        />
        <span>{{ contentTypeInfo.label }}</span>
      </Badge>
      <span class="text-sm text-text-tertiary flex items-center gap-1">
        <Eye class="w-3.5 h-3.5" stroke-width="2" />
        {{ formatReadingTime(article.readingTime) }}
      </span>
    </div>
    
    <!-- Content -->
    <RouterLink :to="articleUrl" class="block group/content">
      <h2 class="text-xl font-bold mb-3 text-text-primary group-hover/content:text-primary transition-colors duration-300 leading-tight font-display">
        {{ article.title }}
        <Badge v-if="article.isEditorial" class="ml-2 align-middle">
          <Check class="w-3 h-3 mr-1" stroke-width="2.5" />
          Материал редакции
        </Badge>
      </h2>
      
      <p 
        v-if="!compact && article.lead" 
        class="text-text-secondary mb-4 line-clamp-2 leading-relaxed"
      >
        {{ article.lead }}
      </p>
      
      <div 
        v-if="article.coverImage && !compact"
        class="relative overflow-hidden rounded-xl"
      >
        <img 
          :src="article.coverImage.url"
          :alt="article.title"
          class="w-full aspect-video object-cover transition-transform duration-500 group-hover/content:scale-105"
          loading="lazy"
        />
        <!-- Gradient overlay -->
        <div class="absolute inset-0 bg-gradient-to-t from-bg-base/40 via-transparent to-transparent opacity-0 group-hover/content:opacity-100 transition-opacity duration-300"></div>
      </div>
    </RouterLink>
    
    <!-- Reactions -->
    <div class="mt-5 pt-4 border-t border-border-subtle">
      <ReactionPanel 
        :reactions="article.reactions" 
        :article-id="article.id"
        compact
      />
    </div>
    
    <!-- Actions -->
    <div class="flex items-center justify-between mt-4 text-sm">
      <div class="flex items-center gap-1">
        <RouterLink 
          :to="`${articleUrl}#comments`"
          class="flex items-center gap-1.5 px-3 py-2 hover:text-primary hover:bg-primary/10 rounded-lg transition-all duration-200"
        >
          <MessageCircle class="w-4 h-4" stroke-width="2" />
          <span class="font-medium">{{ formatCompactNumber(article.commentCount) }}</span>
        </RouterLink>
        
        <button 
          @click="toggleBookmark"
          class="flex items-center gap-1.5 px-3 py-2 rounded-lg transition-all duration-200"
          :class="isBookmarked ? 'text-warning bg-warning/10' : 'text-text-secondary hover:text-warning hover:bg-warning/10'"
        >
          <Bookmark class="w-4 h-4" :class="isBookmarked ? 'fill-current' : ''" stroke-width="2" />
          <span class="font-medium">{{ formatCompactNumber(article.bookmarkCount) }}</span>
        </button>
        
        <button class="flex items-center gap-1.5 px-3 py-2 text-text-secondary hover:text-accent hover:bg-accent/10 rounded-lg transition-all duration-200">
          <Share2 class="w-4 h-4" stroke-width="2" />
        </button>
        
        <button class="flex items-center gap-1.5 px-3 py-2 text-text-secondary hover:text-text-primary hover:bg-bg-hover rounded-lg transition-all duration-200">
          <MoreHorizontal class="w-4 h-4" stroke-width="2" />
        </button>
      </div>
      
      <div class="flex items-center gap-3">
        <button 
          @click="toggleLike"
          class="flex items-center gap-1.5 px-3 py-2 rounded-lg transition-all duration-200"
          :class="isLiked ? 'text-secondary bg-secondary/10' : 'text-text-secondary hover:text-secondary hover:bg-secondary/10'"
        >
          <Heart class="w-4 h-4 transition-transform duration-200" :class="isLiked ? 'fill-current scale-110' : ''" stroke-width="2" />
        </button>
        
        <div class="flex items-center gap-1.5 text-text-tertiary">
          <Eye class="w-4 h-4" stroke-width="2" />
          <span class="font-medium">{{ formatCompactNumber(article.viewCount) }}</span>
        </div>
      </div>
    </div>
  </article>
</template>
