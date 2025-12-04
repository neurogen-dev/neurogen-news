<script setup lang="ts">
import { ref } from 'vue'
import { Image, Link, FileText, Send, Sparkles } from 'lucide-vue-next'
import { useRouter } from 'vue-router'
import Avatar from '@/components/ui/Avatar.vue'
import Button from '@/components/ui/Button.vue'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const content = ref('')
const isFocused = ref(false)

const handleSubmit = async () => {
  if (!content.value.trim()) return
  
  // If content is short, create a quick post
  // Otherwise, open full editor
  if (content.value.length > 200) {
    router.push({
      name: 'editor-new',
      query: { content: content.value }
    })
  } else {
    // TODO: Create quick post via API
    console.log('Quick post:', content.value)
    content.value = ''
  }
}

const openFullEditor = () => {
  router.push({
    name: 'editor-new',
    query: content.value ? { content: content.value } : undefined
  })
}
</script>

<template>
  <div class="velvet-panel p-4">
    <div class="flex gap-3">
      <!-- User avatar -->
      <Avatar 
        :src="authStore.user?.avatarUrl" 
        :alt="authStore.user?.displayName" 
        size="md"
        ring
      />
      
      <div class="flex-1">
        <!-- Text input -->
        <div class="velvet-input">
          <textarea
            v-model="content"
            @focus="isFocused = true"
            placeholder="Поделитесь мыслями, задайте вопрос или напишите пост..."
            class="w-full px-4 py-3 bg-white/60 dark:bg-slate-800/60 border-2 border-transparent rounded-xl text-gray-800 dark:text-gray-100 placeholder:text-gray-400 transition-all duration-200 ease-smooth shadow-inset-soft focus:bg-white dark:focus:bg-slate-800 focus:border-primary focus:shadow-[0_0_0_4px_rgba(99,102,241,0.15),var(--shadow-raised)] resize-none"
            :rows="isFocused ? 3 : 1"
          />
        </div>
        
        <!-- Actions (visible when focused) -->
        <div 
          v-if="isFocused || content"
          class="flex items-center justify-between mt-3 pt-3 border-t border-gray-100 dark:border-gray-800"
        >
          <div class="flex items-center gap-1">
            <button 
              class="p-2 text-gray-400 hover:text-primary hover:bg-primary/10 rounded-lg transition-all duration-200"
              title="Добавить изображение"
            >
              <Image class="w-5 h-5" />
            </button>
            <button 
              class="p-2 text-gray-400 hover:text-primary hover:bg-primary/10 rounded-lg transition-all duration-200"
              title="Добавить ссылку"
            >
              <Link class="w-5 h-5" />
            </button>
            <button 
              @click="openFullEditor"
              class="p-2 text-gray-400 hover:text-primary hover:bg-primary/10 rounded-lg transition-all duration-200"
              title="Открыть полный редактор"
            >
              <FileText class="w-5 h-5" />
            </button>
          </div>
          
          <div class="flex items-center gap-3">
            <span 
              v-if="content.length > 0" 
              class="text-sm transition-colors"
              :class="[
                content.length > 200 ? 'text-error' : 
                content.length > 180 ? 'text-warning' : 'text-gray-400'
              ]"
            >
              {{ content.length }}/200
            </span>
            <Button 
              @click="handleSubmit"
              :disabled="!content.trim()"
              variant="primary"
              size="sm"
            >
              <Send class="w-4 h-4" />
              Опубликовать
              <Sparkles class="w-3 h-3 opacity-60" />
            </Button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
