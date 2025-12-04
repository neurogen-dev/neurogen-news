import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Comment } from '@/types'
import { commentsApi, type CommentListParams, type CreateCommentInput } from '@/api'

export const useCommentsStore = defineStore('comments', () => {
  // State
  const comments = ref<Comment[]>([])
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(20)
  const hasMore = ref(false)
  const isLoading = ref(false)
  const isSubmitting = ref(false)
  const error = ref<string | null>(null)
  const sort = ref<'new' | 'popular' | 'old'>('new')

  // Getters
  const isEmpty = computed(() => comments.value.length === 0 && !isLoading.value)
  const rootComments = computed(() => comments.value.filter(c => !c.parentId))

  // Actions
  async function fetchComments(articleId: string, params?: CommentListParams, append = false) {
    if (isLoading.value) return

    isLoading.value = true
    error.value = null

    try {
      const response = await commentsApi.getByArticle(articleId, {
        sort: sort.value,
        page: append ? page.value + 1 : 1,
        pageSize: pageSize.value,
        ...params,
      })

      if (append) {
        comments.value = [...comments.value, ...response.data.items]
        page.value++
      } else {
        comments.value = response.data.items
        page.value = 1
      }

      total.value = response.data.total
      hasMore.value = response.data.hasMore
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Не удалось загрузить комментарии'
      throw e
    } finally {
      isLoading.value = false
    }
  }

  async function fetchMore(articleId: string) {
    if (!hasMore.value || isLoading.value) return
    await fetchComments(articleId, { sort: sort.value }, true)
  }

  async function fetchReplies(commentId: string, limit = 10, offset = 0) {
    try {
      const response = await commentsApi.getReplies(commentId, limit, offset)
      
      // Update comment in the list with loaded replies
      const commentIndex = comments.value.findIndex(c => c.id === commentId)
      if (commentIndex !== -1) {
        comments.value[commentIndex] = {
          ...comments.value[commentIndex],
          replies: offset === 0 
            ? response.data.items 
            : [...(comments.value[commentIndex].replies || []), ...response.data.items],
        }
      }

      return response.data.items
    } catch (e) {
      console.error('Failed to fetch replies', e)
      throw e
    }
  }

  async function createComment(data: CreateCommentInput) {
    isSubmitting.value = true
    error.value = null

    try {
      const response = await commentsApi.create(data)
      const newComment = response.data

      if (data.parentId) {
        // Add reply to parent comment
        const parentIndex = comments.value.findIndex(c => c.id === data.parentId)
        if (parentIndex !== -1) {
          const parent = comments.value[parentIndex]
          comments.value[parentIndex] = {
            ...parent,
            replyCount: (parent.replyCount || 0) + 1,
            replies: [...(parent.replies || []), newComment],
          }
        }
      } else {
        // Add root comment
        comments.value = [newComment, ...comments.value]
        total.value++
      }

      return newComment
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Не удалось отправить комментарий'
      throw e
    } finally {
      isSubmitting.value = false
    }
  }

  async function updateComment(id: string, content: string) {
    try {
      const response = await commentsApi.update(id, content)
      
      // Update in list
      const index = comments.value.findIndex(c => c.id === id)
      if (index !== -1) {
        comments.value[index] = response.data
      }

      // Update in replies
      comments.value = comments.value.map(comment => ({
        ...comment,
        replies: comment.replies?.map(reply => 
          reply.id === id ? response.data : reply
        ),
      }))

      return response.data
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Не удалось обновить комментарий'
      throw e
    }
  }

  async function deleteComment(id: string) {
    try {
      await commentsApi.delete(id)

      // Mark as deleted in list
      const index = comments.value.findIndex(c => c.id === id)
      if (index !== -1) {
        comments.value[index] = {
          ...comments.value[index],
          isDeleted: true,
          content: '[Комментарий удалён]',
        }
      }

      // Mark in replies
      comments.value = comments.value.map(comment => ({
        ...comment,
        replies: comment.replies?.map(reply => 
          reply.id === id 
            ? { ...reply, isDeleted: true, content: '[Комментарий удалён]' }
            : reply
        ),
      }))
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Не удалось удалить комментарий'
      throw e
    }
  }

  async function addReaction(commentId: string, emoji: string) {
    try {
      await commentsApi.addReaction(commentId, emoji)
    } catch (e) {
      console.error('Failed to add reaction', e)
      throw e
    }
  }

  async function removeReaction(commentId: string) {
    try {
      await commentsApi.removeReaction(commentId)
    } catch (e) {
      console.error('Failed to remove reaction', e)
      throw e
    }
  }

  function setSort(newSort: 'new' | 'popular' | 'old') {
    sort.value = newSort
  }

  function reset() {
    comments.value = []
    total.value = 0
    page.value = 1
    hasMore.value = false
    error.value = null
    sort.value = 'new'
  }

  return {
    // State
    comments,
    total,
    page,
    pageSize,
    hasMore,
    isLoading,
    isSubmitting,
    error,
    sort,
    // Getters
    isEmpty,
    rootComments,
    // Actions
    fetchComments,
    fetchMore,
    fetchReplies,
    createComment,
    updateComment,
    deleteComment,
    addReaction,
    removeReaction,
    setSort,
    reset,
  }
})


