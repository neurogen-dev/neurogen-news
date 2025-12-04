import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { ArticleCard } from '@/types'
import { bookmarksApi, type BookmarkFolder, type BookmarkListParams } from '@/api'

export const useBookmarksStore = defineStore('bookmarks', () => {
  // State
  const bookmarks = ref<ArticleCard[]>([])
  const folders = ref<BookmarkFolder[]>([])
  const currentFolderId = ref<string | null>(null)
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(20)
  const hasMore = ref(false)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // Set to track bookmarked article IDs for quick lookup
  const bookmarkedIds = ref<Set<string>>(new Set())

  // Getters
  const isEmpty = computed(() => bookmarks.value.length === 0 && !isLoading.value)
  const currentFolder = computed(() => 
    currentFolderId.value 
      ? folders.value.find(f => f.id === currentFolderId.value) 
      : null
  )

  // Actions
  async function fetchBookmarks(params?: BookmarkListParams, append = false) {
    if (isLoading.value) return

    isLoading.value = true
    error.value = null

    try {
      const response = await bookmarksApi.list({
        folderId: currentFolderId.value || undefined,
        page: append ? page.value + 1 : 1,
        pageSize: pageSize.value,
        ...params,
      })

      if (append) {
        bookmarks.value = [...bookmarks.value, ...response.data.items]
        page.value++
      } else {
        bookmarks.value = response.data.items
        page.value = 1
      }

      total.value = response.data.total
      hasMore.value = response.data.hasMore

      // Update bookmarked IDs set
      response.data.items.forEach(item => bookmarkedIds.value.add(item.id))
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Не удалось загрузить закладки'
      throw e
    } finally {
      isLoading.value = false
    }
  }

  async function fetchMore() {
    if (!hasMore.value || isLoading.value) return
    await fetchBookmarks({}, true)
  }

  async function fetchFolders() {
    try {
      const response = await bookmarksApi.getFolders()
      folders.value = response.data.items
    } catch (e) {
      console.error('Failed to fetch folders', e)
    }
  }

  async function addBookmark(articleId: string, folderId?: string) {
    try {
      await bookmarksApi.add(articleId, folderId)
      bookmarkedIds.value.add(articleId)
      
      // Refresh bookmarks if viewing the same folder
      if (folderId === currentFolderId.value || (!folderId && !currentFolderId.value)) {
        await fetchBookmarks()
      }
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Не удалось добавить в закладки'
      throw e
    }
  }

  async function removeBookmark(articleId: string) {
    try {
      await bookmarksApi.remove(articleId)
      bookmarkedIds.value.delete(articleId)
      bookmarks.value = bookmarks.value.filter(b => b.id !== articleId)
      total.value--
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Не удалось удалить из закладок'
      throw e
    }
  }

  async function checkBookmark(articleId: string): Promise<boolean> {
    if (bookmarkedIds.value.has(articleId)) {
      return true
    }

    try {
      const response = await bookmarksApi.check(articleId)
      if (response.data.isBookmarked) {
        bookmarkedIds.value.add(articleId)
      }
      return response.data.isBookmarked
    } catch (e) {
      return false
    }
  }

  function isBookmarked(articleId: string): boolean {
    return bookmarkedIds.value.has(articleId)
  }

  async function toggleBookmark(articleId: string, folderId?: string) {
    if (isBookmarked(articleId)) {
      await removeBookmark(articleId)
      return false
    } else {
      await addBookmark(articleId, folderId)
      return true
    }
  }

  async function createFolder(name: string) {
    try {
      const response = await bookmarksApi.createFolder(name)
      folders.value.push(response.data)
      return response.data
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Не удалось создать папку'
      throw e
    }
  }

  async function updateFolder(id: string, name: string) {
    try {
      await bookmarksApi.updateFolder(id, name)
      const index = folders.value.findIndex(f => f.id === id)
      if (index !== -1) {
        folders.value[index] = { ...folders.value[index], name }
      }
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Не удалось обновить папку'
      throw e
    }
  }

  async function deleteFolder(id: string) {
    try {
      await bookmarksApi.deleteFolder(id)
      folders.value = folders.value.filter(f => f.id !== id)
      
      // If viewing deleted folder, switch to all bookmarks
      if (currentFolderId.value === id) {
        currentFolderId.value = null
        await fetchBookmarks()
      }
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Не удалось удалить папку'
      throw e
    }
  }

  async function moveToFolder(articleId: string, folderId?: string) {
    try {
      await bookmarksApi.moveToFolder(articleId, folderId)
      
      // Refresh if moved out of current folder
      if (currentFolderId.value && currentFolderId.value !== folderId) {
        bookmarks.value = bookmarks.value.filter(b => b.id !== articleId)
        total.value--
      }
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Не удалось переместить закладку'
      throw e
    }
  }

  function setCurrentFolder(folderId: string | null) {
    currentFolderId.value = folderId
  }

  function reset() {
    bookmarks.value = []
    folders.value = []
    currentFolderId.value = null
    total.value = 0
    page.value = 1
    hasMore.value = false
    error.value = null
    bookmarkedIds.value.clear()
  }

  return {
    // State
    bookmarks,
    folders,
    currentFolderId,
    total,
    page,
    pageSize,
    hasMore,
    isLoading,
    error,
    // Getters
    isEmpty,
    currentFolder,
    // Actions
    fetchBookmarks,
    fetchMore,
    fetchFolders,
    addBookmark,
    removeBookmark,
    checkBookmark,
    isBookmarked,
    toggleBookmark,
    createFolder,
    updateFolder,
    deleteFolder,
    moveToFolder,
    setCurrentFolder,
    reset,
  }
})


