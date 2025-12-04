<script setup lang="ts">
import { computed } from 'vue'
import { 
  Sparkles, 
  TrendingUp, 
  Clock, 
  BookOpen, 
  Newspaper, 
  HelpCircle,
  Circle,
  CheckCircle2
} from 'lucide-vue-next'
import { cn } from '@/utils/cn'
import type { ArticleLevel, ContentType } from '@/types'

interface Props {
  sort?: 'popular' | 'new'
  level?: ArticleLevel
  contentType?: ContentType
}

const props = withDefaults(defineProps<Props>(), {
  sort: 'popular',
})

const emit = defineEmits<{
  'update:sort': [value: 'popular' | 'new']
  'update:level': [value: ArticleLevel | undefined]
  'update:contentType': [value: ContentType | undefined]
}>()

const sortOptions = [
  { value: 'popular' as const, label: 'Популярное', icon: TrendingUp },
  { value: 'new' as const, label: 'Свежее', icon: Clock },
]

const levelOptions: { 
  value: ArticleLevel | undefined
  label: string
  icon: typeof Circle
  colorClass: string
}[] = [
  { value: undefined, label: 'Все', icon: Circle, colorClass: 'text-text-secondary' },
  { value: 'beginner', label: 'Для новичков', icon: CheckCircle2, colorClass: 'text-beginner' },
  { value: 'intermediate', label: 'Продвинутое', icon: CheckCircle2, colorClass: 'text-intermediate' },
  { value: 'advanced', label: 'Для бизнеса', icon: CheckCircle2, colorClass: 'text-advanced' },
]

const contentTypeOptions: { 
  value: ContentType | undefined
  label: string
  icon: typeof BookOpen
}[] = [
  { value: undefined, label: 'Все', icon: Circle },
  { value: 'article', label: 'Статьи', icon: BookOpen },
  { value: 'news', label: 'Новости', icon: Newspaper },
  { value: 'question', label: 'Вопросы', icon: HelpCircle },
]

const currentSort = computed({
  get: () => props.sort,
  set: (value) => emit('update:sort', value),
})

const handleLevelChange = (value: ArticleLevel | undefined) => {
  emit('update:level', value)
}

const handleContentTypeChange = (value: ContentType | undefined) => {
  emit('update:contentType', value)
}

const getLevelBadgeClasses = (option: typeof levelOptions[0]) => {
  if (props.level !== option.value) {
    return 'hover:bg-bg-hover'
  }
  
  if (option.value === 'beginner') {
    return 'bg-beginner/15 text-beginner border-beginner/30'
  }
  if (option.value === 'intermediate') {
    return 'bg-intermediate/15 text-intermediate border-intermediate/30'
  }
  if (option.value === 'advanced') {
    return 'bg-advanced/15 text-advanced border-advanced/30'
  }
  return 'bg-primary/15 text-primary border-primary/30'
}
</script>

<template>
  <div class="empatra-panel p-5 space-y-5">
    <!-- Sort tabs -->
    <div class="empatra-tabs w-fit">
      <button
        v-for="option in sortOptions"
        :key="option.value"
        @click="currentSort = option.value"
        class="empatra-tab"
        :class="{ 'active': currentSort === option.value }"
      >
        <span class="relative z-10 flex items-center gap-2">
          <component 
            :is="option.icon" 
            class="w-4 h-4"
            stroke-width="2"
          />
          <span class="font-medium">{{ option.label }}</span>
        </span>
      </button>
    </div>
    
    <!-- Level filters -->
    <div class="flex flex-wrap gap-2">
      <button
        v-for="option in levelOptions"
        :key="String(option.value)"
        @click="handleLevelChange(option.value)"
        :class="cn(
          'inline-flex items-center gap-2 px-3 py-1.5 text-sm font-medium rounded-full border border-border-subtle transition-all duration-200',
          getLevelBadgeClasses(option)
        )"
      >
        <component 
          :is="option.icon" 
          class="w-3.5 h-3.5"
          :class="props.level === option.value ? option.colorClass : 'text-text-tertiary'"
          stroke-width="2.5"
          :fill="props.level === option.value ? 'currentColor' : 'none'"
          :fill-opacity="props.level === option.value ? 0.2 : 0"
        />
        <span>{{ option.label }}</span>
        <Sparkles 
          v-if="props.level === option.value" 
          class="w-3 h-3 opacity-60 animate-pulse" 
          stroke-width="2.5"
        />
      </button>
    </div>
    
    <!-- Content type filters -->
    <div class="flex flex-wrap gap-2">
      <button
        v-for="option in contentTypeOptions"
        :key="String(option.value)"
        @click="handleContentTypeChange(option.value)"
        :class="cn(
          'inline-flex items-center gap-2 px-3 py-1.5 text-sm font-medium rounded-full border border-border-subtle transition-all duration-200',
          contentType === option.value
            ? 'bg-primary/15 text-primary border-primary/30'
            : 'hover:bg-bg-hover'
        )"
      >
        <component 
          :is="option.icon" 
          class="w-3.5 h-3.5"
          :class="contentType === option.value ? 'text-primary' : 'text-text-tertiary'"
          stroke-width="2.5"
        />
        <span>{{ option.label }}</span>
        <Sparkles 
          v-if="contentType === option.value" 
          class="w-3 h-3 opacity-60 animate-pulse" 
          stroke-width="2.5"
        />
      </button>
    </div>
  </div>
</template>
