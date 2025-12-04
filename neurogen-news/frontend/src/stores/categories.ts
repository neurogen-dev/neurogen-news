import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Category, ArticleCard } from '@/types'
import { categoriesApi, type CategoryArticleParams } from '@/api'

export const useCategoriesStore = defineStore('categories', () => {
  // State
  const categories = ref<Category[]>([])
  const currentCategory = ref<Category | null>(null)
  const categoryArticles = ref<ArticleCard[]>([])
  const subscriptions = ref<Category[]>([])
  const isSubscribed = ref(false)
  const articlesTotal = ref(0)
  const articlesPage = ref(1)
  const articlesHasMore = ref(false)
  const isLoading = ref(false)
  const isLoadingArticles = ref(false)
  const error = ref<string | null>(null)

  // Getters
  const isEmpty = computed(() => categories.value.length === 0 && !isLoading.value)
  const officialCategories = computed(() => categories.value.filter(c => c.isOfficial))
  const communityCategories = computed(() => categories.value.filter(c => !c.isOfficial))

  // Actions
  async function fetchCategories() {
    if (isLoading.value) return

    isLoading.value = true
    error.value = null

    try {
      const response = await categoriesApi.list()
      categories.value = response.data.items
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Не удалось загрузить категории'
      throw e
    } finally {
      isLoading.value = false
    }
  }

  async function fetchCategoryBySlug(slug: string) {
    isLoading.value = true
    error.value = null

    try {
      const response = await categoriesApi.getBySlug(slug)
      currentCategory.value = response.data.category
      isSubscribed.value = response.data.isSubscribed
      return response.data.category
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Категория не найдена'
      throw e
    } finally {
      isLoading.value = false
    }
  }

  async function fetchCategoryArticles(slug: string, params?: CategoryArticleParams, append = false) {
    if (isLoadingArticles.value) return

    isLoadingArticles.value = true
    error.value = null

    try {
      const response = await categoriesApi.getArticles(slug, {
        page: append ? articlesPage.value + 1 : 1,
        pageSize: 20,
        ...params,
      })

      if (append) {
        categoryArticles.value = [...categoryArticles.value, ...response.data.items]
        articlesPage.value++
      } else {
        categoryArticles.value = response.data.items
        articlesPage.value = 1
      }

      articlesTotal.value = response.data.total
      articlesHasMore.value = response.data.hasMore
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Не удалось загрузить статьи'
      throw e
    } finally {
      isLoadingArticles.value = false
    }
  }

  async function fetchMoreArticles(slug: string, params?: CategoryArticleParams) {
    if (!articlesHasMore.value || isLoadingArticles.value) return
    await fetchCategoryArticles(slug, params, true)
  }

  async function subscribe(slug: string) {
    try {
      await categoriesApi.subscribe(slug)
      isSubscribed.value = true

      // Update subscriber count
      if (currentCategory.value) {
        currentCategory.value.subscriberCount++
      }

      // Add to subscriptions list
      if (currentCategory.value && !subscriptions.value.find(c => c.slug === slug)) {
        subscriptions.value.push(currentCategory.value)
      }
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Не удалось подписаться'
      throw e
    }
  }

  async function unsubscribe(slug: string) {
    try {
      await categoriesApi.unsubscribe(slug)
      isSubscribed.value = false

      // Update subscriber count
      if (currentCategory.value) {
        currentCategory.value.subscriberCount = Math.max(0, currentCategory.value.subscriberCount - 1)
      }

      // Remove from subscriptions list
      subscriptions.value = subscriptions.value.filter(c => c.slug !== slug)
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Не удалось отписаться'
      throw e
    }
  }

  async function toggleSubscription(slug: string) {
    if (isSubscribed.value) {
      await unsubscribe(slug)
    } else {
      await subscribe(slug)
    }
  }

  async function fetchSubscriptions() {
    try {
      const response = await categoriesApi.getSubscriptions()
      subscriptions.value = response.data.items
    } catch (e) {
      console.error('Failed to fetch subscriptions', e)
    }
  }

  function getCategoryBySlug(slug: string): Category | undefined {
    return categories.value.find(c => c.slug === slug)
  }

  function clearCurrentCategory() {
    currentCategory.value = null
    categoryArticles.value = []
    isSubscribed.value = false
    articlesTotal.value = 0
    articlesPage.value = 1
    articlesHasMore.value = false
  }

  function reset() {
    categories.value = []
    currentCategory.value = null
    categoryArticles.value = []
    subscriptions.value = []
    isSubscribed.value = false
    articlesTotal.value = 0
    articlesPage.value = 1
    articlesHasMore.value = false
    error.value = null
  }

  return {
    // State
    categories,
    currentCategory,
    categoryArticles,
    subscriptions,
    isSubscribed,
    articlesTotal,
    articlesPage,
    articlesHasMore,
    isLoading,
    isLoadingArticles,
    error,
    // Getters
    isEmpty,
    officialCategories,
    communityCategories,
    // Actions
    fetchCategories,
    fetchCategoryBySlug,
    fetchCategoryArticles,
    fetchMoreArticles,
    subscribe,
    unsubscribe,
    toggleSubscription,
    fetchSubscriptions,
    getCategoryBySlug,
    clearCurrentCategory,
    reset,
  }
})


