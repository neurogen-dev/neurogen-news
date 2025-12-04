import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { ArticleCard, FeedFilters, PaginatedResponse } from '@/types'
import { apiClient } from '@/api/client'

// Mock data for development
const MOCK_ARTICLES: ArticleCard[] = [
  {
    id: '1',
    title: 'Как использовать ChatGPT для написания кода: полное руководство 2024',
    slug: 'kak-ispolzovat-chatgpt-dlya-napisaniya-koda',
    lead: 'Подробный гайд по эффективному использованию ChatGPT и других LLM для программирования. Рассмотрим лучшие промпты, техники и реальные примеры.',
    level: 'beginner',
    contentType: 'article',
    readingTime: 12,
    viewCount: 15420,
    commentCount: 89,
    bookmarkCount: 342,
    publishedAt: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(),
    author: {
      id: 'u1',
      username: 'aiexpert',
      displayName: 'AI Expert',
      avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=aiexpert',
      isVerified: true,
    },
    category: {
      id: 'c1',
      name: 'Чат-боты',
      slug: 'chatbots',
      icon: '💬',
    },
    reactions: [
      { emoji: '🔥', count: 156, isReacted: false },
      { emoji: '👍', count: 89, isReacted: true },
      { emoji: '❤️', count: 45, isReacted: false },
    ],
    isEditorial: true,
    coverImage: { url: 'https://images.unsplash.com/photo-1677442136019-21780ecad995?w=800&h=450&fit=crop' },
  },
  {
    id: '2',
    title: 'Midjourney v6: Что нового и как использовать',
    slug: 'midjourney-v6-chto-novogo',
    lead: 'Разбираем все новые функции Midjourney v6: улучшенная работа с текстом, новые стили и более точное следование промптам.',
    level: 'intermediate',
    contentType: 'news',
    readingTime: 8,
    viewCount: 8932,
    commentCount: 45,
    bookmarkCount: 189,
    publishedAt: new Date(Date.now() - 5 * 60 * 60 * 1000).toISOString(),
    author: {
      id: 'u2',
      username: 'designer_ai',
      displayName: 'Дизайнер AI',
      avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=designer',
      isVerified: false,
    },
    category: {
      id: 'c2',
      name: 'Изображения',
      slug: 'images',
      icon: '🎨',
    },
    reactions: [
      { emoji: '😍', count: 234, isReacted: false },
      { emoji: '🔥', count: 78, isReacted: false },
    ],
    isEditorial: false,
    coverImage: { url: 'https://images.unsplash.com/photo-1686191128892-3b37add13e64?w=800&h=450&fit=crop' },
  },
  {
    id: '3',
    title: 'Создаём музыку с помощью Suno AI: пошаговый гайд',
    slug: 'sozdaem-muzyku-s-pomoschyu-suno-ai',
    lead: 'Научитесь создавать музыкальные треки любого жанра с помощью нейросети Suno. От простых мелодий до полноценных песен.',
    level: 'beginner',
    contentType: 'article',
    readingTime: 15,
    viewCount: 6234,
    commentCount: 67,
    bookmarkCount: 256,
    publishedAt: new Date(Date.now() - 12 * 60 * 60 * 1000).toISOString(),
    author: {
      id: 'u3',
      username: 'musicmaker',
      displayName: 'Music Maker',
      avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=music',
      isVerified: true,
    },
    category: {
      id: 'c3',
      name: 'Музыка',
      slug: 'music',
      icon: '🎵',
    },
    reactions: [
      { emoji: '🎵', count: 189, isReacted: true },
      { emoji: '👏', count: 67, isReacted: false },
    ],
    isEditorial: false,
    coverImage: { url: 'https://images.unsplash.com/photo-1511379938547-c1f69419868d?w=800&h=450&fit=crop' },
  },
  {
    id: '4',
    title: 'Claude 3.5 Sonnet vs GPT-4o: детальное сравнение для разработчиков',
    slug: 'claude-3-5-sonnet-vs-gpt-4o',
    lead: 'Сравниваем две топовые модели по скорости, качеству кода, следованию инструкциям и стоимости API.',
    level: 'advanced',
    contentType: 'article',
    readingTime: 20,
    viewCount: 12567,
    commentCount: 134,
    bookmarkCount: 456,
    publishedAt: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString(),
    author: {
      id: 'u4',
      username: 'techreviewer',
      displayName: 'Tech Reviewer',
      avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=tech',
      isVerified: true,
    },
    category: {
      id: 'c1',
      name: 'Чат-боты',
      slug: 'chatbots',
      icon: '💬',
    },
    reactions: [
      { emoji: '🤔', count: 89, isReacted: false },
      { emoji: '💡', count: 156, isReacted: false },
      { emoji: '👍', count: 234, isReacted: false },
    ],
    isEditorial: true,
    coverImage: { url: 'https://images.unsplash.com/photo-1555255707-c07966088b7b?w=800&h=450&fit=crop' },
  },
  {
    id: '5',
    title: 'Как я автоматизировал свой бизнес с помощью AI-агентов',
    slug: 'kak-ya-avtomatiziroval-biznes-ai-agentami',
    lead: 'Реальный кейс автоматизации рутинных задач: от обработки email до генерации отчётов. Экономия 20+ часов в неделю.',
    level: 'intermediate',
    contentType: 'post',
    readingTime: 10,
    viewCount: 4532,
    commentCount: 78,
    bookmarkCount: 167,
    publishedAt: new Date(Date.now() - 36 * 60 * 60 * 1000).toISOString(),
    author: {
      id: 'u5',
      username: 'entrepreneur',
      displayName: 'Предприниматель',
      avatarUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=entrepreneur',
      isVerified: false,
    },
    category: {
      id: 'c4',
      name: 'Код',
      slug: 'code',
      icon: '💻',
    },
    reactions: [
      { emoji: '💰', count: 123, isReacted: false },
      { emoji: '🚀', count: 89, isReacted: false },
    ],
    isEditorial: false,
  },
]

const USE_MOCK = true // Toggle for development

export const useFeedStore = defineStore('feed', () => {
  // State
  const articles = ref<ArticleCard[]>([])
  const isLoading = ref(false)
  const isLoadingMore = ref(false)
  const error = ref<string | null>(null)
  const filters = ref<FeedFilters>({
    sort: 'popular',
  })
  const pagination = ref({
    page: 1,
    pageSize: 20,
    total: 0,
  })
  
  // Getters
  const hasMore = computed(() => 
    articles.value.length < pagination.value.total
  )
  
  const isEmpty = computed(() => 
    !isLoading.value && articles.value.length === 0
  )
  
  // Actions
  async function fetchArticles(newFilters?: Partial<FeedFilters>) {
    if (newFilters) {
      filters.value = { ...filters.value, ...newFilters }
    }
    
    isLoading.value = true
    error.value = null
    
    // Use mock data in development
    if (USE_MOCK) {
      await new Promise(resolve => setTimeout(resolve, 500)) // Simulate loading
      articles.value = MOCK_ARTICLES
      pagination.value.total = MOCK_ARTICLES.length
      isLoading.value = false
      return
    }
    
    try {
      const params = new URLSearchParams({
        sort: filters.value.sort,
        page: '1',
        pageSize: String(pagination.value.pageSize),
      })
      
      if (filters.value.level) params.set('level', filters.value.level)
      if (filters.value.contentType) params.set('contentType', filters.value.contentType)
      if (filters.value.categoryId) params.set('categoryId', filters.value.categoryId)
      if (filters.value.tagId) params.set('tagId', filters.value.tagId)
      if (filters.value.timeRange) params.set('timeRange', filters.value.timeRange)
      
      const response = await apiClient.get<PaginatedResponse<ArticleCard>>(
        `/articles?${params.toString()}`
      )
      
      articles.value = response.data.items
      pagination.value.page = response.data.page
      pagination.value.total = response.data.total
    } catch (e) {
      error.value = 'Не удалось загрузить статьи'
      console.error('Failed to fetch articles:', e)
    } finally {
      isLoading.value = false
    }
  }
  
  async function loadMore() {
    if (!hasMore.value || isLoadingMore.value) return
    
    isLoadingMore.value = true
    
    try {
      const nextPage = pagination.value.page + 1
      const params = new URLSearchParams({
        sort: filters.value.sort,
        page: String(nextPage),
        pageSize: String(pagination.value.pageSize),
      })
      
      if (filters.value.level) params.set('level', filters.value.level)
      if (filters.value.contentType) params.set('contentType', filters.value.contentType)
      if (filters.value.categoryId) params.set('categoryId', filters.value.categoryId)
      if (filters.value.tagId) params.set('tagId', filters.value.tagId)
      if (filters.value.timeRange) params.set('timeRange', filters.value.timeRange)
      
      const response = await apiClient.get<PaginatedResponse<ArticleCard>>(
        `/articles?${params.toString()}`
      )
      
      articles.value.push(...response.data.items)
      pagination.value.page = response.data.page
    } catch (e) {
      console.error('Failed to load more articles:', e)
    } finally {
      isLoadingMore.value = false
    }
  }
  
  function reset() {
    articles.value = []
    isLoading.value = false
    isLoadingMore.value = false
    error.value = null
    pagination.value = {
      page: 1,
      pageSize: 20,
      total: 0,
    }
    filters.value = {
      sort: 'popular',
    }
  }
  
  return {
    // State
    articles,
    isLoading,
    isLoadingMore,
    error,
    filters,
    pagination,
    
    // Getters
    hasMore,
    isEmpty,
    
    // Actions
    fetchArticles,
    loadMore,
    reset,
  }
})
