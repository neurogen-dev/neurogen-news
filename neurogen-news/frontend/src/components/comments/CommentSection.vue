<script setup lang="ts">
import { ref, onMounted } from 'vue'
import CommentItem from './CommentItem.vue'
import CommentInput from './CommentInput.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import { useAuthStore } from '@/stores/auth'
import type { Comment } from '@/types'

interface Props {
  articleId: string
}

const props = defineProps<Props>()

const authStore = useAuthStore()
const comments = ref<Comment[]>([])
const isLoading = ref(true)
const sortBy = ref<'best' | 'new'>('best')

// Mock comments
const loadComments = async () => {
  isLoading.value = true
  
  // Simulate API call
  await new Promise(resolve => setTimeout(resolve, 600))
  
  comments.value = [
    {
      id: '1',
      content: 'Отличная статья! Особенно понравились примеры промптов. Уже попробовал — работает намного лучше, чем мои обычные запросы.',
      htmlContent: '<p>Отличная статья! Особенно понравились примеры промптов. Уже попробовал — работает намного лучше, чем мои обычные запросы.</p>',
      author: {
        id: '2',
        username: 'alexey',
        displayName: 'Алексей Иванов',
        avatarUrl: undefined,
        isVerified: false,
      },
      articleId: props.articleId,
      reactions: [
        { emoji: '👍', count: 12, isReacted: false },
        { emoji: '❤️', count: 3, isReacted: false },
      ],
      replyCount: 2,
      isEdited: false,
      createdAt: new Date(Date.now() - 1000 * 60 * 60 * 5).toISOString(),
      depth: 0,
      replies: [
        {
          id: '2',
          content: 'Согласен! Техника с указанием роли реально работает. Я теперь всегда начинаю с "Ты — эксперт в..."',
          htmlContent: '<p>Согласен! Техника с указанием роли реально работает. Я теперь всегда начинаю с "Ты — эксперт в..."</p>',
          author: {
            id: '3',
            username: 'maria',
            displayName: 'Мария',
            avatarUrl: undefined,
            isVerified: false,
          },
          articleId: props.articleId,
          parentId: '1',
          reactions: [
            { emoji: '👍', count: 5, isReacted: true },
          ],
          replyCount: 0,
          isEdited: false,
          createdAt: new Date(Date.now() - 1000 * 60 * 60 * 3).toISOString(),
          depth: 1,
        },
      ],
    },
    {
      id: '3',
      content: 'А есть подобный гайд для Claude? Хотелось бы сравнить подходы к разным моделям.',
      htmlContent: '<p>А есть подобный гайд для Claude? Хотелось бы сравнить подходы к разным моделям.</p>',
      author: {
        id: '4',
        username: 'ivan',
        displayName: 'Иван',
        avatarUrl: undefined,
        isVerified: false,
      },
      articleId: props.articleId,
      reactions: [
        { emoji: '🤔', count: 4, isReacted: false },
      ],
      replyCount: 1,
      isEdited: false,
      createdAt: new Date(Date.now() - 1000 * 60 * 60 * 2).toISOString(),
      depth: 0,
      replies: [
        {
          id: '4',
          content: 'Планируем опубликовать на следующей неделе! Подписывайтесь, чтобы не пропустить.',
          htmlContent: '<p>Планируем опубликовать на следующей неделе! Подписывайтесь, чтобы не пропустить.</p>',
          author: {
            id: '1',
            username: 'neurogen',
            displayName: 'Редакция Neurogen',
            avatarUrl: undefined,
            isVerified: true,
          },
          articleId: props.articleId,
          parentId: '3',
          reactions: [
            { emoji: '👍', count: 8, isReacted: false },
            { emoji: '🎉', count: 2, isReacted: false },
          ],
          replyCount: 0,
          isEdited: false,
          createdAt: new Date(Date.now() - 1000 * 60 * 60).toISOString(),
          depth: 1,
        },
      ],
    },
  ]
  
  isLoading.value = false
}

const handleNewComment = async (content: string) => {
  // TODO: Call API to create comment
  console.log('New comment:', content)
  
  // Optimistically add comment
  const newComment: Comment = {
    id: String(Date.now()),
    content,
    htmlContent: `<p>${content}</p>`,
    author: {
      id: authStore.user?.id || '',
      username: authStore.user?.username || '',
      displayName: authStore.user?.displayName || '',
      avatarUrl: authStore.user?.avatarUrl,
      isVerified: authStore.user?.isVerified || false,
    },
    articleId: props.articleId,
    reactions: [],
    replyCount: 0,
    isEdited: false,
    createdAt: new Date().toISOString(),
    depth: 0,
  }
  
  comments.value.unshift(newComment)
}

onMounted(loadComments)
</script>

<template>
  <div>
    <!-- Comment input -->
    <CommentInput 
      v-if="authStore.isLoggedIn"
      @submit="handleNewComment"
      class="mb-6"
    />
    
    <!-- Login prompt -->
    <div 
      v-else
      class="p-4 bg-background-secondary dark:bg-dark-tertiary rounded-xl text-center mb-6"
    >
      <p class="text-text-secondary mb-2">
        Войдите, чтобы оставить комментарий
      </p>
      <RouterLink 
        to="/login"
        class="text-primary hover:underline"
      >
        Войти
      </RouterLink>
    </div>
    
    <!-- Sort -->
    <div class="flex items-center gap-2 mb-4">
      <span class="text-sm text-text-tertiary">Сортировка:</span>
      <button
        @click="sortBy = 'best'"
        class="px-3 py-1 text-sm rounded-full transition-colors"
        :class="sortBy === 'best' 
          ? 'bg-primary text-white' 
          : 'bg-background-secondary dark:bg-dark-tertiary text-text-secondary hover:text-text-primary'"
      >
        Лучшие
      </button>
      <button
        @click="sortBy = 'new'"
        class="px-3 py-1 text-sm rounded-full transition-colors"
        :class="sortBy === 'new' 
          ? 'bg-primary text-white' 
          : 'bg-background-secondary dark:bg-dark-tertiary text-text-secondary hover:text-text-primary'"
      >
        Новые
      </button>
    </div>
    
    <!-- Loading -->
    <div v-if="isLoading" class="space-y-4">
      <div v-for="i in 3" :key="i" class="space-y-2">
        <div class="flex items-center gap-3">
          <Skeleton class="w-10 h-10 rounded-full" />
          <Skeleton class="h-4 w-32" />
        </div>
        <Skeleton class="h-16 ml-13" />
      </div>
    </div>
    
    <!-- Comments list -->
    <div v-else class="space-y-4">
      <CommentItem
        v-for="comment in comments"
        :key="comment.id"
        :comment="comment"
      />
      
      <!-- Empty state -->
      <div 
        v-if="comments.length === 0"
        class="text-center py-8 text-text-tertiary"
      >
        Пока нет комментариев. Будьте первым!
      </div>
    </div>
  </div>
</template>

