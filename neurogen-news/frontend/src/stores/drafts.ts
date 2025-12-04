import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Draft, ContentType } from '@/types'
import { draftsApi, type CreateDraftInput, type UpdateDraftInput, type AutoSaveInput } from '@/api'

export const useDraftsStore = defineStore('drafts', () => {
  // State
  const drafts = ref<Draft[]>([])
  const currentDraft = ref<Draft | null>(null)
  const autoSaveDraft = ref<Draft | null>(null)
  const total = ref(0)
  const page = ref(1)
  const hasMore = ref(false)
  const isLoading = ref(false)
  const isSaving = ref(false)
  const error = ref<string | null>(null)
  const lastAutoSave = ref<Date | null>(null)

  // Getters
  const isEmpty = computed(() => drafts.value.length === 0 && !isLoading.value)
  const hasAutoSave = computed(() => autoSaveDraft.value !== null)

  // Actions
  async function fetchDrafts(limit = 20, offset = 0, append = false) {
    if (isLoading.value) return

    isLoading.value = true
    error.value = null

    try {
      const response = await draftsApi.list(limit, offset)

      if (append) {
        drafts.value = [...drafts.value, ...response.data.items]
      } else {
        drafts.value = response.data.items
        page.value = 1
      }

      total.value = response.data.total
      hasMore.value = response.data.hasMore
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Не удалось загрузить черновики'
      throw e
    } finally {
      isLoading.value = false
    }
  }

  async function fetchMore() {
    if (!hasMore.value || isLoading.value) return
    page.value++
    await fetchDrafts(20, (page.value - 1) * 20, true)
  }

  async function fetchDraftById(id: string) {
    isLoading.value = true
    error.value = null

    try {
      const response = await draftsApi.getById(id)
      currentDraft.value = response.data
      return response.data
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Черновик не найден'
      throw e
    } finally {
      isLoading.value = false
    }
  }

  async function createDraft(data: CreateDraftInput) {
    isSaving.value = true
    error.value = null

    try {
      const response = await draftsApi.create(data)
      drafts.value = [response.data, ...drafts.value]
      total.value++
      currentDraft.value = response.data
      return response.data
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Не удалось создать черновик'
      throw e
    } finally {
      isSaving.value = false
    }
  }

  async function updateDraft(id: string, data: UpdateDraftInput) {
    isSaving.value = true
    error.value = null

    try {
      const response = await draftsApi.update(id, data)

      // Update in list
      const index = drafts.value.findIndex(d => d.id === id)
      if (index !== -1) {
        drafts.value[index] = response.data
      }

      if (currentDraft.value?.id === id) {
        currentDraft.value = response.data
      }

      return response.data
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Не удалось обновить черновик'
      throw e
    } finally {
      isSaving.value = false
    }
  }

  async function deleteDraft(id: string) {
    try {
      await draftsApi.delete(id)
      drafts.value = drafts.value.filter(d => d.id !== id)
      total.value--

      if (currentDraft.value?.id === id) {
        currentDraft.value = null
      }
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Не удалось удалить черновик'
      throw e
    }
  }

  async function autoSave(data: AutoSaveInput) {
    try {
      await draftsApi.autoSave(data)
      lastAutoSave.value = new Date()
    } catch (e) {
      console.error('Auto-save failed', e)
      // Don't throw - auto-save failures shouldn't interrupt user
    }
  }

  async function fetchAutoSave(articleId?: string) {
    try {
      const response = await draftsApi.getAutoSave(articleId)
      autoSaveDraft.value = response.data.draft
      return response.data.draft
    } catch (e) {
      console.error('Failed to fetch auto-save', e)
      return null
    }
  }

  function clearAutoSave() {
    autoSaveDraft.value = null
  }

  function clearCurrentDraft() {
    currentDraft.value = null
  }

  function reset() {
    drafts.value = []
    currentDraft.value = null
    autoSaveDraft.value = null
    total.value = 0
    page.value = 1
    hasMore.value = false
    error.value = null
    lastAutoSave.value = null
  }

  return {
    // State
    drafts,
    currentDraft,
    autoSaveDraft,
    total,
    page,
    hasMore,
    isLoading,
    isSaving,
    error,
    lastAutoSave,
    // Getters
    isEmpty,
    hasAutoSave,
    // Actions
    fetchDrafts,
    fetchMore,
    fetchDraftById,
    createDraft,
    updateDraft,
    deleteDraft,
    autoSave,
    fetchAutoSave,
    clearAutoSave,
    clearCurrentDraft,
    reset,
  }
})


