import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Article, ArticleCard } from '@/types'
import { articlesApi, type ArticleListParams, type CreateArticleInput, type UpdateArticleInput } from '@/api'

export const useArticlesStore = defineStore('articles', () => {
  // State
  const articles = ref<ArticleCard[]>([])
  const currentArticle = ref<Article | null>(null)
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(20)
  const hasMore = ref(false)
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  const filters = ref<ArticleListParams>({
    sort: 'popular',
    timeRange: 'all',
  })

  // Getters
  const isEmpty = computed(() => articles.value.length === 0 && !isLoading.value)

  // Actions
  async function fetchArticles(params?: ArticleListParams, append = false) {
    if (isLoading.value) return

    isLoading.value = true
    error.value = null

    try {
      const mergedParams = { ...filters.value, ...params, page: append ? page.value + 1 : 1, pageSize: pageSize.value }
      const response = await articlesApi.list(mergedParams)
      
      if (append) {
        articles.value = [...articles.value, ...response.data.items]
        page.value++
      } else {
        articles.value = response.data.items
        page.value = 1
      }
      
      total.value = response.data.total
      hasMore.value = response.data.hasMore
      filters.value = mergedParams
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Не удалось загрузить статьи'
      throw e
    } finally {
      isLoading.value = false
    }
  }

  async function fetchMore() {
    if (!hasMore.value || isLoading.value) return
    await fetchArticles(filters.value, true)
  }

  async function fetchArticleById(id: string) {
    isLoading.value = true
    error.value = null

    try {
      const response = await articlesApi.getById(id)
      currentArticle.value = response.data
      return response.data
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Статья не найдена'
      throw e
    } finally {
      isLoading.value = false
    }
  }

  async function fetchArticleBySlug(category: string, slug: string) {
    isLoading.value = true
    error.value = null

    try {
      const response = await articlesApi.getBySlug(category, slug)
      currentArticle.value = response.data
      return response.data
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Статья не найдена'
      throw e
    } finally {
      isLoading.value = false
    }
  }

  async function createArticle(data: CreateArticleInput) {
    isLoading.value = true
    error.value = null

    try {
      const response = await articlesApi.create(data)
      return response.data
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Не удалось создать статью'
      throw e
    } finally {
      isLoading.value = false
    }
  }

  async function updateArticle(id: string, data: UpdateArticleInput) {
    isLoading.value = true
    error.value = null

    try {
      const response = await articlesApi.update(id, data)
      if (currentArticle.value?.id === id) {
        currentArticle.value = response.data
      }
      return response.data
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Не удалось обновить статью'
      throw e
    } finally {
      isLoading.value = false
    }
  }

  async function deleteArticle(id: string) {
    try {
      await articlesApi.delete(id)
      articles.value = articles.value.filter(a => a.id !== id)
      if (currentArticle.value?.id === id) {
        currentArticle.value = null
      }
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Не удалось удалить статью'
      throw e
    }
  }

  async function addReaction(id: string, emoji: string) {
    try {
      await articlesApi.addReaction(id, emoji)
      // Update local state if needed
    } catch (e) {
      console.error('Failed to add reaction', e)
      throw e
    }
  }

  async function removeReaction(id: string) {
    try {
      await articlesApi.removeReaction(id)
    } catch (e) {
      console.error('Failed to remove reaction', e)
      throw e
    }
  }

  function setFilters(newFilters: ArticleListParams) {
    filters.value = { ...filters.value, ...newFilters }
  }

  function clearCurrentArticle() {
    currentArticle.value = null
  }

  function reset() {
    articles.value = []
    currentArticle.value = null
    total.value = 0
    page.value = 1
    hasMore.value = false
    error.value = null
    filters.value = { sort: 'popular', timeRange: 'all' }
  }

  return {
    // State
    articles,
    currentArticle,
    total,
    page,
    pageSize,
    hasMore,
    isLoading,
    error,
    filters,
    // Getters
    isEmpty,
    // Actions
    fetchArticles,
    fetchMore,
    fetchArticleById,
    fetchArticleBySlug,
    createArticle,
    updateArticle,
    deleteArticle,
    addReaction,
    removeReaction,
    setFilters,
    clearCurrentArticle,
    reset,
  }
})


