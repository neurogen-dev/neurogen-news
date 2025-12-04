<script setup lang="ts">
import { ref, computed } from 'vue'
import { Send, X, Image, Link, Sparkles } from 'lucide-vue-next'
import Avatar from '@/components/ui/Avatar.vue'
import Button from '@/components/ui/Button.vue'
import { useAuthStore } from '@/stores/auth'

interface Props {
  placeholder?: string
  compact?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: 'Написать комментарий...',
  compact: false,
})

const emit = defineEmits<{
  submit: [content: string]
  cancel: []
}>()

const authStore = useAuthStore()
const content = ref('')
const isFocused = ref(false)
const isSubmitting = ref(false)

const isEmpty = computed(() => !content.value.trim())

const handleSubmit = async () => {
  if (isEmpty.value || isSubmitting.value) return
  
  isSubmitting.value = true
  try {
    emit('submit', content.value.trim())
    content.value = ''
    isFocused.value = false
  } finally {
    isSubmitting.value = false
  }
}

const handleCancel = () => {
  content.value = ''
  isFocused.value = false
  emit('cancel')
}
</script>

<template>
  <div 
    class="flex gap-3"
    :class="{ 'items-start': isFocused || content, 'items-center': !isFocused && !content }"
  >
    <!-- Avatar -->
    <Avatar 
      v-if="!compact"
      :src="authStore.user?.avatarUrl" 
      :alt="authStore.user?.displayName"
      :size="isFocused ? 'md' : 'sm'"
    />
    
    <!-- Input area -->
    <div 
      class="flex-1 velvet-panel overflow-hidden transition-all duration-200"
      :class="{ 'shadow-raised ring-2 ring-primary/20': isFocused }"
    >
      <textarea
        v-model="content"
        :placeholder="placeholder"
        @focus="isFocused = true"
        class="w-full px-4 py-3 bg-transparent border-none outline-none resize-none text-gray-800 dark:text-white placeholder:text-gray-400"
        :rows="isFocused || content ? 3 : 1"
      />
      
      <!-- Actions (visible when focused) -->
      <div 
        v-if="isFocused || content"
        class="flex items-center justify-between px-3 py-2.5 border-t border-gray-100 dark:border-gray-800"
      >
        <div class="flex items-center gap-1">
          <button 
            class="p-2 text-gray-400 hover:text-primary hover:bg-primary/10 rounded-lg transition-all duration-200"
            title="Добавить изображение"
          >
            <Image class="w-4 h-4" />
          </button>
          <button 
            class="p-2 text-gray-400 hover:text-primary hover:bg-primary/10 rounded-lg transition-all duration-200"
            title="Добавить ссылку"
          >
            <Link class="w-4 h-4" />
          </button>
        </div>
        
        <div class="flex items-center gap-2">
          <Button 
            v-if="compact"
            variant="ghost" 
            size="sm"
            @click="handleCancel"
          >
            <X class="w-4 h-4" />
          </Button>
          
          <Button 
            variant="primary"
            size="sm"
            :disabled="isEmpty"
            :loading="isSubmitting"
            @click="handleSubmit"
          >
            <Send class="w-4 h-4" />
            <span class="hidden sm:inline">Отправить</span>
            <Sparkles class="w-3 h-3 opacity-60 hidden sm:inline" />
          </Button>
        </div>
      </div>
    </div>
  </div>
</template>
