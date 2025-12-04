<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink } from 'vue-router'
import { MessageCircle, MoreHorizontal, Check, ChevronDown, ChevronUp } from 'lucide-vue-next'
import Avatar from '@/components/ui/Avatar.vue'
import ReactionPanel from '@/components/article/ReactionPanel.vue'
import CommentInput from './CommentInput.vue'
import { formatRelativeTime } from '@/utils/formatters'
import type { Comment } from '@/types'

interface Props {
  comment: Comment
}

const props = defineProps<Props>()

const showReplyInput = ref(false)
const showReplies = ref(true)

const handleReply = async (content: string) => {
  // TODO: Call API to create reply
  console.log('Reply to', props.comment.id, ':', content)
  showReplyInput.value = false
}
</script>

<template>
  <div 
    class="relative"
    :class="{ 'ml-6 sm:ml-12': comment.depth > 0 }"
  >
    <!-- Thread line for nested comments -->
    <div 
      v-if="comment.depth > 0"
      class="absolute -left-4 sm:-left-6 top-0 bottom-0 w-px bg-gradient-to-b from-gray-200 via-gray-200 to-transparent dark:from-gray-700 dark:via-gray-700"
    />
    
    <div class="flex gap-3 group/comment">
      <!-- Avatar -->
      <RouterLink :to="`/@${comment.author.username}`" class="shrink-0">
        <Avatar 
          :src="comment.author.avatarUrl" 
          :alt="comment.author.displayName" 
          :size="comment.depth > 0 ? 'sm' : 'md'"
        />
      </RouterLink>
      
      <!-- Content -->
      <div class="flex-1 min-w-0">
        <!-- Header -->
        <div class="flex items-center gap-2 mb-1">
          <RouterLink 
            :to="`/@${comment.author.username}`"
            class="font-medium text-gray-900 dark:text-white hover:text-primary text-sm transition-colors"
          >
            {{ comment.author.displayName }}
          </RouterLink>
          
          <span 
            v-if="comment.author.isVerified" 
            class="w-4 h-4 bg-gradient-to-br from-primary to-primary-600 rounded-full flex items-center justify-center"
          >
            <Check class="w-2.5 h-2.5 text-white" />
          </span>
          
          <span class="text-gray-400 text-sm">
            {{ formatRelativeTime(comment.createdAt) }}
          </span>
          
          <span v-if="comment.isEdited" class="text-gray-400 text-xs">
            (изменено)
          </span>
        </div>
        
        <!-- Text -->
        <div 
          class="text-gray-700 dark:text-gray-200 text-sm mb-2 prose prose-sm dark:prose-invert max-w-none"
          v-html="comment.htmlContent"
        />
        
        <!-- Actions -->
        <div class="flex items-center gap-3">
          <ReactionPanel 
            :reactions="comment.reactions" 
            :article-id="comment.id"
            compact
          />
          
          <button
            @click="showReplyInput = !showReplyInput"
            class="flex items-center gap-1.5 text-sm text-gray-400 hover:text-primary transition-all duration-200"
          >
            <MessageCircle class="w-4 h-4" />
            Ответить
          </button>
          
          <button class="p-1.5 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg opacity-0 group-hover/comment:opacity-100 transition-all duration-200">
            <MoreHorizontal class="w-4 h-4" />
          </button>
        </div>
        
        <!-- Reply input -->
        <div v-if="showReplyInput" class="mt-3">
          <CommentInput 
            @submit="handleReply"
            @cancel="showReplyInput = false"
            placeholder="Написать ответ..."
            compact
          />
        </div>
        
        <!-- Replies -->
        <div v-if="comment.replies && comment.replies.length > 0" class="mt-4">
          <!-- Toggle replies -->
          <button
            @click="showReplies = !showReplies"
            class="flex items-center gap-1.5 text-sm text-primary hover:text-primary-600 transition-colors mb-3"
          >
            <component :is="showReplies ? ChevronUp : ChevronDown" class="w-4 h-4" />
            {{ showReplies ? 'Скрыть' : 'Показать' }} ответы ({{ comment.replies.length }})
          </button>
          
          <!-- Nested comments -->
          <div v-if="showReplies" class="space-y-4">
            <CommentItem
              v-for="reply in comment.replies"
              :key="reply.id"
              :comment="reply"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
