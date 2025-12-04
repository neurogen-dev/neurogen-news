<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { 
  Save, 
  Send, 
  Eye, 
  Settings, 
  ChevronLeft,
  Image as ImageIcon,
  X
} from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import Badge from '@/components/ui/Badge.vue'
import TipTapEditor from '@/components/editor/TipTapEditor.vue'
import { useAuthStore } from '@/stores/auth'
import type { ArticleLevel, ContentType, Category } from '@/types'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const isEditMode = computed(() => !!route.params.id)

// Form state
const title = ref('')
const content = ref('')
const lead = ref('')
const coverImage = ref<string | null>(null)
const level = ref<ArticleLevel>('beginner')
const contentType = ref<ContentType>('article')
const categoryId = ref('')
const tags = ref<string[]>([])
const tagInput = ref('')
const commentsEnabled = ref(true)
const isNSFW = ref(false)

// UI state
const isSaving = ref(false)
const isPublishing = ref(false)
const showSettings = ref(false)
const lastSaved = ref<Date | null>(null)

// Categories (mock - would come from API)
const categories: Category[] = [
  { id: '1', name: 'Чат-боты', slug: 'chatbots', icon: '💬', articleCount: 0 },
  { id: '2', name: 'Изображения', slug: 'images', icon: '🎨', articleCount: 0 },
  { id: '3', name: 'Видео', slug: 'video', icon: '🎬', articleCount: 0 },
  { id: '4', name: 'Музыка', slug: 'music', icon: '🎵', articleCount: 0 },
  { id: '5', name: 'Текст', slug: 'text', icon: '✍️', articleCount: 0 },
  { id: '6', name: 'Код', slug: 'code', icon: '💻', articleCount: 0 },
]

const levels = [
  { value: 'beginner', label: '🟢 Для новичков' },
  { value: 'intermediate', label: '🟡 Продвинутое' },
  { value: 'advanced', label: '🔴 Для бизнеса' },
]

const contentTypes = [
  { value: 'article', label: '📖 Статья' },
  { value: 'news', label: '📰 Новость' },
  { value: 'post', label: '✏️ Пост' },
  { value: 'question', label: '❓ Вопрос' },
]

// Auto-save
let autoSaveTimer: ReturnType<typeof setTimeout>

const autoSave = () => {
  clearTimeout(autoSaveTimer)
  autoSaveTimer = setTimeout(async () => {
    await saveDraft()
  }, 30000) // Auto-save every 30 seconds
}

watch([title, content, lead, level, contentType, categoryId, tags], autoSave)

const saveDraft = async () => {
  if (!title.value && !content.value) return
  
  isSaving.value = true
  try {
    // TODO: Call API to save draft
    console.log('Saving draft...', {
      title: title.value,
      content: content.value,
      level: level.value,
      contentType: contentType.value,
    })
    
    lastSaved.value = new Date()
  } finally {
    isSaving.value = false
  }
}

const publish = async () => {
  if (!title.value || !content.value || !categoryId.value) {
    // TODO: Show validation errors
    return
  }
  
  isPublishing.value = true
  try {
    // TODO: Call API to publish article
    console.log('Publishing...', {
      title: title.value,
      content: content.value,
      lead: lead.value,
      coverImage: coverImage.value,
      level: level.value,
      contentType: contentType.value,
      categoryId: categoryId.value,
      tags: tags.value,
      commentsEnabled: commentsEnabled.value,
      isNSFW: isNSFW.value,
    })
    
    // Redirect to the published article
    router.push('/')
  } finally {
    isPublishing.value = false
  }
}

const addTag = () => {
  const tag = tagInput.value.trim().toLowerCase()
  if (tag && !tags.value.includes(tag) && tags.value.length < 10) {
    tags.value.push(tag)
  }
  tagInput.value = ''
}

const removeTag = (index: number) => {
  tags.value.splice(index, 1)
}

const handleCoverUpload = () => {
  // TODO: Implement file upload
  const url = window.prompt('URL изображения для обложки')
  if (url) {
    coverImage.value = url
  }
}

// Load draft if editing
onMounted(async () => {
  if (isEditMode.value) {
    // TODO: Load article data
  }
  
  // Handle content from quick post
  const queryContent = route.query.content as string
  if (queryContent) {
    content.value = queryContent
  }
})

onBeforeUnmount(() => {
  clearTimeout(autoSaveTimer)
})
</script>

<template>
  <div class="min-h-screen -mx-4 -my-6">
    <!-- Header -->
    <header class="sticky top-16 z-40 bg-white dark:bg-dark-secondary border-b border-border dark:border-dark-tertiary">
      <div class="flex items-center justify-between px-4 py-3">
        <div class="flex items-center gap-4">
          <button
            @click="router.back()"
            class="p-2 text-text-tertiary hover:text-text-primary transition-colors"
          >
            <ChevronLeft class="w-5 h-5" />
          </button>
          
          <div class="text-sm text-text-tertiary">
            <span v-if="lastSaved">
              Сохранено {{ lastSaved.toLocaleTimeString() }}
            </span>
            <span v-else-if="isSaving">
              Сохранение...
            </span>
            <span v-else>
              Черновик
            </span>
          </div>
        </div>
        
        <div class="flex items-center gap-2">
          <Button variant="ghost" size="sm" @click="saveDraft" :loading="isSaving">
            <Save class="w-4 h-4 mr-1" />
            Сохранить
          </Button>
          
          <Button variant="secondary" size="sm">
            <Eye class="w-4 h-4 mr-1" />
            Предпросмотр
          </Button>
          
          <Button variant="ghost" size="sm" @click="showSettings = !showSettings">
            <Settings class="w-4 h-4" />
          </Button>
          
          <Button @click="publish" :loading="isPublishing">
            <Send class="w-4 h-4 mr-1" />
            Опубликовать
          </Button>
        </div>
      </div>
    </header>
    
    <div class="flex">
      <!-- Main editor -->
      <main class="flex-1 p-6 max-w-4xl mx-auto">
        <!-- Cover image -->
        <div 
          v-if="coverImage"
          class="relative aspect-video rounded-xl overflow-hidden mb-6"
        >
          <img :src="coverImage" alt="Cover" class="w-full h-full object-cover" />
          <button
            @click="coverImage = null"
            class="absolute top-2 right-2 p-1 bg-black/50 text-white rounded-full hover:bg-black/70 transition-colors"
          >
            <X class="w-4 h-4" />
          </button>
        </div>
        
        <button
          v-else
          @click="handleCoverUpload"
          class="flex items-center justify-center gap-2 w-full aspect-video rounded-xl border-2 border-dashed border-border dark:border-dark-tertiary hover:border-primary hover:bg-primary/5 transition-colors mb-6"
        >
          <ImageIcon class="w-6 h-6 text-text-tertiary" />
          <span class="text-text-tertiary">Добавить обложку</span>
        </button>
        
        <!-- Title -->
        <input
          v-model="title"
          type="text"
          placeholder="Заголовок статьи"
          class="w-full text-3xl font-bold bg-transparent border-none outline-none text-text-primary dark:text-white placeholder:text-text-tertiary mb-4"
        />
        
        <!-- Lead -->
        <textarea
          v-model="lead"
          placeholder="Краткое описание (опционально)"
          rows="2"
          class="w-full text-lg bg-transparent border-none outline-none resize-none text-text-secondary placeholder:text-text-tertiary mb-6"
        />
        
        <!-- Content editor -->
        <TipTapEditor 
          v-model="content"
          placeholder="Начните писать..."
        />
      </main>
      
      <!-- Settings sidebar -->
      <aside 
        v-if="showSettings"
        class="w-80 border-l border-border dark:border-dark-tertiary bg-white dark:bg-dark-secondary p-4 space-y-6"
      >
        <h3 class="font-semibold text-text-primary dark:text-white">
          Настройки публикации
        </h3>
        
        <!-- Category -->
        <div>
          <label class="block text-sm text-text-tertiary mb-2">Категория *</label>
          <select
            v-model="categoryId"
            class="w-full px-3 py-2 bg-background-secondary dark:bg-dark-tertiary border border-border dark:border-dark-tertiary rounded-lg text-text-primary dark:text-white"
          >
            <option value="">Выберите категорию</option>
            <option 
              v-for="category in categories" 
              :key="category.id" 
              :value="category.id"
            >
              {{ category.icon }} {{ category.name }}
            </option>
          </select>
        </div>
        
        <!-- Level -->
        <div>
          <label class="block text-sm text-text-tertiary mb-2">Уровень сложности</label>
          <select
            v-model="level"
            class="w-full px-3 py-2 bg-background-secondary dark:bg-dark-tertiary border border-border dark:border-dark-tertiary rounded-lg text-text-primary dark:text-white"
          >
            <option 
              v-for="l in levels" 
              :key="l.value" 
              :value="l.value"
            >
              {{ l.label }}
            </option>
          </select>
        </div>
        
        <!-- Content type -->
        <div>
          <label class="block text-sm text-text-tertiary mb-2">Тип контента</label>
          <select
            v-model="contentType"
            class="w-full px-3 py-2 bg-background-secondary dark:bg-dark-tertiary border border-border dark:border-dark-tertiary rounded-lg text-text-primary dark:text-white"
          >
            <option 
              v-for="type in contentTypes" 
              :key="type.value" 
              :value="type.value"
            >
              {{ type.label }}
            </option>
          </select>
        </div>
        
        <!-- Tags -->
        <div>
          <label class="block text-sm text-text-tertiary mb-2">
            Теги (до 10)
          </label>
          <div class="flex flex-wrap gap-2 mb-2">
            <Badge 
              v-for="(tag, index) in tags" 
              :key="tag"
              variant="primary"
              class="cursor-pointer"
              @click="removeTag(index)"
            >
              {{ tag }} ×
            </Badge>
          </div>
          <input
            v-model="tagInput"
            @keydown.enter.prevent="addTag"
            type="text"
            placeholder="Добавить тег и нажать Enter"
            class="w-full px-3 py-2 bg-background-secondary dark:bg-dark-tertiary border border-border dark:border-dark-tertiary rounded-lg text-text-primary dark:text-white text-sm"
            :disabled="tags.length >= 10"
          />
        </div>
        
        <!-- Options -->
        <div class="space-y-3">
          <label class="flex items-center gap-2 cursor-pointer">
            <input 
              v-model="commentsEnabled" 
              type="checkbox"
              class="w-4 h-4 rounded text-primary"
            />
            <span class="text-sm text-text-secondary">Разрешить комментарии</span>
          </label>
          
          <label class="flex items-center gap-2 cursor-pointer">
            <input 
              v-model="isNSFW" 
              type="checkbox"
              class="w-4 h-4 rounded text-primary"
            />
            <span class="text-sm text-text-secondary">18+ контент</span>
          </label>
        </div>
      </aside>
    </div>
  </div>
</template>

