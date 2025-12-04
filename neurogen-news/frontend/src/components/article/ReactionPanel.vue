<script setup lang="ts">
import { ref } from 'vue'
import { Plus } from 'lucide-vue-next'
import { cn } from '@/utils/cn'
import type { ReactionCount } from '@/types'

interface Props {
  reactions: ReactionCount[]
  articleId: string
  compact?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  compact: false,
})

const showPicker = ref(false)

// Available reactions
const availableReactions = [
  { emoji: '👍', label: 'Нравится' },
  { emoji: '❤️', label: 'Супер' },
  { emoji: '😂', label: 'Смешно' },
  { emoji: '🤔', label: 'Хмм' },
  { emoji: '😢', label: 'Грустно' },
  { emoji: '😡', label: 'Злюсь' },
  { emoji: '🔥', label: 'Огонь' },
  { emoji: '🎉', label: 'Праздник' },
]

const handleReaction = (emoji: string) => {
  // TODO: Send reaction to API
  console.log('React with', emoji, 'on article', props.articleId)
  showPicker.value = false
}

const toggleReaction = (emoji: string) => {
  // TODO: Toggle existing reaction
  console.log('Toggle reaction', emoji)
}
</script>

<template>
  <div class="flex items-center gap-2 flex-wrap">
    <!-- Existing reactions -->
    <button
      v-for="reaction in reactions"
      :key="reaction.emoji"
      @click="toggleReaction(reaction.emoji)"
      :class="cn(
        'inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full transition-all duration-200',
        reaction.isReacted
          ? 'bg-primary/15 text-primary shadow-sm shadow-primary/10'
          : 'bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-gray-700'
      )"
    >
      <span class="text-base">{{ reaction.emoji }}</span>
      <span class="text-sm font-medium">{{ reaction.count }}</span>
    </button>
    
    <!-- Add reaction button -->
    <div class="relative">
      <button
        @click="showPicker = !showPicker"
        class="inline-flex items-center justify-center w-8 h-8 rounded-full bg-gray-100 dark:bg-gray-800 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-700 transition-all duration-200"
      >
        <Plus class="w-4 h-4" />
      </button>
      
      <!-- Reaction picker -->
      <Transition
        enter-active-class="transition duration-150 ease-spring"
        enter-from-class="opacity-0 scale-95 translate-y-1"
        enter-to-class="opacity-100 scale-100 translate-y-0"
        leave-active-class="transition duration-100 ease-smooth"
        leave-from-class="opacity-100 scale-100"
        leave-to-class="opacity-0 scale-95"
      >
        <div 
          v-if="showPicker"
          class="absolute bottom-full left-0 mb-2 p-2 velvet-panel shadow-floating flex gap-1 z-50"
        >
          <button
            v-for="reaction in availableReactions"
            :key="reaction.emoji"
            @click="handleReaction(reaction.emoji)"
            class="p-2 text-xl hover:bg-primary/10 rounded-lg transition-all duration-200 hover:scale-110"
            :title="reaction.label"
          >
            {{ reaction.emoji }}
          </button>
        </div>
      </Transition>
    </div>
  </div>
</template>
