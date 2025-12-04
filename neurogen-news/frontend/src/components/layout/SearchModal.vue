<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Search, X, Clock, TrendingUp, FileText, User, Wrench, ArrowRight } from 'lucide-vue-next'
import type { SearchResult } from '@/types'

const emit = defineEmits<{
  close: []
}>()

const router = useRouter()
const query = ref('')
const inputRef = ref<HTMLInputElement | null>(null)
const isLoading = ref(false)
const results = ref<SearchResult[]>([])
const selectedIndex = ref(0)

// Recent searches (from localStorage)
const recentSearches = ref<string[]>(
  JSON.parse(localStorage.getItem('recentSearches') || '[]')
)

// Popular searches
const popularSearches = [
  'ChatGPT промпты',
  'Midjourney гайд',
  'Как создать нейросеть',
  'Лучшие AI инструменты',
  'Генерация изображений',
]

const showSuggestions = computed(() => 
  !query.value && (recentSearches.value.length > 0 || popularSearches.length > 0)
)

const typeIcons = {
  article: FileText,
  user: User,
  tool: Wrench,
  category: TrendingUp,
}

// Debounced search
let searchTimeout: ReturnType<typeof setTimeout>
watch(query, (newQuery) => {
  clearTimeout(searchTimeout)
  if (newQuery.length < 2) {
    results.value = []
    return
  }
  
  isLoading.value = true
  searchTimeout = setTimeout(async () => {
    // TODO: Call search API
    // For now, mock results
    results.value = [
      {
        type: 'article',
        id: '1',
        title: `${newQuery} - Полное руководство`,
        description: 'Подробный гайд по использованию нейросетей...',
        url: '/chatbots/full-guide',
      },
      {
        type: 'tool',
        id: '2',
        title: `AI Tool - ${newQuery}`,
        description: 'Инструмент для генерации контента',
        url: '/tools/ai-tool',
      },
    ]
    isLoading.value = false
    selectedIndex.value = 0
  }, 300)
})

const handleSearch = () => {
  if (!query.value.trim()) return
  
  // Save to recent searches
  const searches = recentSearches.value.filter(s => s !== query.value)
  searches.unshift(query.value)
  recentSearches.value = searches.slice(0, 5)
  localStorage.setItem('recentSearches', JSON.stringify(recentSearches.value))
  
  // Navigate to search page
  router.push({ name: 'search', query: { q: query.value } })
  emit('close')
}

const handleResultClick = (result: SearchResult) => {
  router.push(result.url)
  emit('close')
}

const handleSuggestionClick = (suggestion: string) => {
  query.value = suggestion
  handleSearch()
}

const clearRecentSearches = () => {
  recentSearches.value = []
  localStorage.removeItem('recentSearches')
}

// Keyboard navigation
const handleKeydown = (event: KeyboardEvent) => {
  if (event.key === 'Escape') {
    emit('close')
  } else if (event.key === 'ArrowDown') {
    event.preventDefault()
    selectedIndex.value = Math.min(selectedIndex.value + 1, results.value.length - 1)
  } else if (event.key === 'ArrowUp') {
    event.preventDefault()
    selectedIndex.value = Math.max(selectedIndex.value - 1, 0)
  } else if (event.key === 'Enter' && results.value[selectedIndex.value]) {
    handleResultClick(results.value[selectedIndex.value])
  }
}

// Global keyboard shortcut
const handleGlobalKeydown = (event: KeyboardEvent) => {
  if ((event.metaKey || event.ctrlKey) && event.key === 'k') {
    event.preventDefault()
    inputRef.value?.focus()
  }
}

onMounted(() => {
  inputRef.value?.focus()
  document.addEventListener('keydown', handleGlobalKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleGlobalKeydown)
  clearTimeout(searchTimeout)
})
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-start justify-center pt-20 px-4">
    <!-- Backdrop -->
    <div 
      class="absolute inset-0 bg-black/40 backdrop-blur-sm"
      @click="emit('close')"
    />
    
    <!-- Modal -->
    <div 
      class="relative w-full max-w-2xl velvet-panel shadow-floating overflow-hidden animate-scale-in"
      @keydown="handleKeydown"
    >
      <!-- Search input -->
      <div class="flex items-center gap-3 px-4 py-4 border-b border-gray-100 dark:border-gray-800">
        <Search class="w-5 h-5 text-gray-400 shrink-0" />
        <input
          ref="inputRef"
          v-model="query"
          type="text"
          placeholder="Поиск статей, инструментов, пользователей..."
          class="flex-1 bg-transparent border-none outline-none text-gray-800 dark:text-white placeholder:text-gray-400"
          @keydown.enter="handleSearch"
        />
        <button
          @click="emit('close')"
          class="p-1.5 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-all duration-200"
        >
          <X class="w-5 h-5" />
        </button>
      </div>
      
      <!-- Content -->
      <div class="max-h-96 overflow-y-auto scrollbar-thin">
        <!-- Loading -->
        <div v-if="isLoading" class="p-8 flex justify-center">
          <div class="relative w-8 h-8">
            <div class="absolute inset-0 rounded-full border-2 border-primary/20"></div>
            <div class="absolute inset-0 rounded-full border-2 border-primary border-t-transparent animate-spin"></div>
          </div>
        </div>
        
        <!-- Results -->
        <template v-else-if="results.length > 0">
          <div class="py-2">
            <button
              v-for="(result, index) in results"
              :key="result.id"
              @click="handleResultClick(result)"
              class="flex items-center gap-3 w-full px-4 py-3 text-left transition-all duration-200"
              :class="[
                index === selectedIndex
                  ? 'bg-primary/10 text-primary'
                  : 'hover:bg-gray-50 dark:hover:bg-gray-800/50'
              ]"
            >
              <component 
                :is="typeIcons[result.type]" 
                class="w-5 h-5 text-gray-400 shrink-0" 
              />
              <div class="flex-1 min-w-0">
                <div class="font-medium text-gray-800 dark:text-white truncate">
                  {{ result.title }}
                </div>
                <div 
                  v-if="result.description" 
                  class="text-sm text-gray-500 truncate"
                >
                  {{ result.description }}
                </div>
              </div>
              <ArrowRight class="w-4 h-4 text-gray-400 shrink-0" />
            </button>
          </div>
        </template>
        
        <!-- Suggestions -->
        <template v-else-if="showSuggestions">
          <!-- Recent searches -->
          <div v-if="recentSearches.length > 0" class="py-3">
            <div class="flex items-center justify-between px-4 mb-2">
              <h4 class="text-xs font-semibold text-gray-400 uppercase tracking-wider">
                Недавние поиски
              </h4>
              <button
                @click="clearRecentSearches"
                class="text-xs text-gray-400 hover:text-primary transition-colors"
              >
                Очистить
              </button>
            </div>
            <button
              v-for="search in recentSearches"
              :key="search"
              @click="handleSuggestionClick(search)"
              class="flex items-center gap-3 w-full px-4 py-2.5 text-left text-gray-600 dark:text-gray-300 hover:text-gray-800 dark:hover:text-white hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors"
            >
              <Clock class="w-4 h-4 text-gray-400" />
              {{ search }}
            </button>
          </div>
          
          <!-- Popular searches -->
          <div class="py-3 border-t border-gray-100 dark:border-gray-800">
            <h4 class="px-4 mb-2 text-xs font-semibold text-gray-400 uppercase tracking-wider">
              Популярные запросы
            </h4>
            <button
              v-for="search in popularSearches"
              :key="search"
              @click="handleSuggestionClick(search)"
              class="flex items-center gap-3 w-full px-4 py-2.5 text-left text-gray-600 dark:text-gray-300 hover:text-gray-800 dark:hover:text-white hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors"
            >
              <TrendingUp class="w-4 h-4 text-gray-400" />
              {{ search }}
            </button>
          </div>
        </template>
        
        <!-- No results -->
        <div 
          v-else-if="query.length >= 2" 
          class="p-8 text-center text-gray-400"
        >
          Ничего не найдено по запросу "{{ query }}"
        </div>
      </div>
      
      <!-- Footer -->
      <div class="flex items-center justify-between px-4 py-3 border-t border-gray-100 dark:border-gray-800 text-xs text-gray-400">
        <div class="flex items-center gap-4">
          <span class="flex items-center gap-1">
            <kbd class="px-1.5 py-0.5 bg-gray-100 dark:bg-gray-800 rounded font-mono">↵</kbd>
            для поиска
          </span>
          <span class="flex items-center gap-1">
            <kbd class="px-1.5 py-0.5 bg-gray-100 dark:bg-gray-800 rounded font-mono">↑↓</kbd>
            для навигации
          </span>
        </div>
        <span class="flex items-center gap-1">
          <kbd class="px-1.5 py-0.5 bg-gray-100 dark:bg-gray-800 rounded font-mono">Esc</kbd>
          для закрытия
        </span>
      </div>
    </div>
  </div>
</template>
