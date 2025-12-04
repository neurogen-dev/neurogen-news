<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { 
  Home, 
  TrendingUp, 
  Clock, 
  Rss,
  MessageSquare,
  Image as ImageIcon,
  Video,
  Music,
  FileText,
  Code,
  BookOpen,
  Search as SearchIcon,
  Lightbulb,
  Newspaper,
  HelpCircle,
  Wrench,
  GraduationCap,
  ExternalLink,
  Send,
  Heart,
  Sparkles
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const authStore = useAuthStore()

const mainNav = [
  { name: 'Главная', icon: Home, to: '/' },
  { name: 'Популярное', icon: TrendingUp, to: '/popular' },
  { name: 'Свежее', icon: Clock, to: '/new' },
]

const personalNav = computed(() => 
  authStore.isLoggedIn ? [
    { name: 'Моя лента', icon: Rss, to: '/my-feed' },
  ] : []
)

const categories = [
  { name: 'Чат-боты', icon: MessageSquare, to: '/chatbots' },
  { name: 'Изображения', icon: ImageIcon, to: '/images' },
  { name: 'Видео', icon: Video, to: '/video' },
  { name: 'Музыка', icon: Music, to: '/music' },
  { name: 'Текст', icon: FileText, to: '/text' },
  { name: 'Код', icon: Code, to: '/code' },
]

const contentTypes = [
  { name: 'Гайды', icon: BookOpen, to: '/guides' },
  { name: 'Обзоры', icon: SearchIcon, to: '/reviews' },
  { name: 'Промпты', icon: Lightbulb, to: '/prompts' },
  { name: 'Новости', icon: Newspaper, to: '/news' },
  { name: 'Вопросы', icon: HelpCircle, to: '/qa' },
]

const tools = [
  { name: 'Каталог AI', icon: Wrench, to: '/tools' },
  { name: 'Начать здесь', icon: GraduationCap, to: '/start' },
]

const externalLinks = [
  { name: 'Telegram', url: 'https://t.me/neurogen_news', icon: Send },
  { name: 'Empatra', url: 'https://empatra.ai', icon: Heart },
]

const isActive = (path: string) => {
  if (path === '/') return route.path === '/'
  return route.path.startsWith(path)
}
</script>

<template>
  <aside class="fixed top-16 left-0 w-64 h-[calc(100vh-4rem)] bg-bg-elevated/50 backdrop-blur-xl border-r border-border-subtle overflow-y-auto scrollbar-thin">
    <!-- Ambient glow -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none z-0">
      <div class="absolute top-20 -left-20 w-48 h-48 bg-primary/10 rounded-full blur-3xl"></div>
      <div class="absolute bottom-40 -left-12 w-40 h-40 bg-secondary/8 rounded-full blur-3xl"></div>
    </div>
    
    <nav class="relative p-4 space-y-6 z-10">
      <!-- Main navigation -->
      <div class="space-y-1">
        <RouterLink
          v-for="item in mainNav"
          :key="item.to"
          :to="item.to"
          class="empatra-nav-item"
          :class="[
            isActive(item.to)
              ? 'active'
              : ''
          ]"
        >
          <component 
            :is="item.icon" 
            class="w-5 h-5 transition-transform duration-200 group-hover:scale-110" 
            stroke-width="2"
          />
          <span class="flex-1">{{ item.name }}</span>
          <Sparkles 
            v-if="isActive(item.to)" 
            class="w-3.5 h-3.5 opacity-60 animate-pulse" 
            stroke-width="2.5"
          />
        </RouterLink>
        
        <RouterLink
          v-for="item in personalNav"
          :key="item.to"
          :to="item.to"
          class="empatra-nav-item"
          :class="[isActive(item.to) ? 'active' : '']"
        >
          <component 
            :is="item.icon" 
            class="w-5 h-5 transition-transform duration-200"
            stroke-width="2"
          />
          <span class="flex-1">{{ item.name }}</span>
        </RouterLink>
      </div>
      
      <!-- Categories -->
      <div>
        <h3 class="px-3 mb-2 text-xs font-semibold text-text-tertiary uppercase tracking-widest flex items-center gap-2">
          <span class="w-8 h-px bg-gradient-to-r from-primary/50 to-transparent"></span>
          Категории
        </h3>
        <div class="space-y-0.5">
          <RouterLink
            v-for="item in categories"
            :key="item.to"
            :to="item.to"
            class="empatra-nav-item !py-2"
            :class="[isActive(item.to) ? 'active' : '']"
          >
            <component 
              :is="item.icon" 
              class="w-4 h-4"
              stroke-width="2"
            />
            <span>{{ item.name }}</span>
          </RouterLink>
        </div>
      </div>
      
      <!-- Content types -->
      <div>
        <h3 class="px-3 mb-2 text-xs font-semibold text-text-tertiary uppercase tracking-widest flex items-center gap-2">
          <span class="w-8 h-px bg-gradient-to-r from-accent/50 to-transparent"></span>
          Контент
        </h3>
        <div class="space-y-0.5">
          <RouterLink
            v-for="item in contentTypes"
            :key="item.to"
            :to="item.to"
            class="empatra-nav-item !py-2"
            :class="[isActive(item.to) ? 'active' : '']"
          >
            <component 
              :is="item.icon" 
              class="w-4 h-4"
              stroke-width="2"
            />
            <span>{{ item.name }}</span>
          </RouterLink>
        </div>
      </div>
      
      <!-- Tools -->
      <div>
        <h3 class="px-3 mb-2 text-xs font-semibold text-text-tertiary uppercase tracking-widest flex items-center gap-2">
          <span class="w-8 h-px bg-gradient-to-r from-warning/50 to-transparent"></span>
          Инструменты
        </h3>
        <div class="space-y-0.5">
          <RouterLink
            v-for="item in tools"
            :key="item.to"
            :to="item.to"
            class="empatra-nav-item !py-2"
            :class="[isActive(item.to) ? 'active' : '']"
          >
            <component 
              :is="item.icon" 
              class="w-4 h-4"
              stroke-width="2"
            />
            <span>{{ item.name }}</span>
          </RouterLink>
        </div>
      </div>
      
      <!-- External links -->
      <div class="pt-4 border-t border-border-subtle">
        <a
          v-for="link in externalLinks"
          :key="link.url"
          :href="link.url"
          target="_blank"
          rel="noopener noreferrer"
          class="empatra-nav-item !py-2 group"
        >
          <component 
            :is="link.icon" 
            class="w-4 h-4"
            stroke-width="2"
          />
          <span class="flex-1">{{ link.name }}</span>
          <ExternalLink class="w-3.5 h-3.5 opacity-40 group-hover:opacity-100 transition-opacity" stroke-width="2" />
        </a>
      </div>
      
      <!-- Empatra Branding -->
      <div class="pt-4">
        <a
          href="https://empatra.ai"
          target="_blank"
          rel="noopener noreferrer"
          class="block p-3 rounded-xl bg-gradient-to-r from-primary/10 via-secondary/10 to-accent/10 border border-primary/20 hover:border-primary/40 transition-all group"
        >
          <div class="flex items-center gap-2 mb-1">
            <Heart class="w-4 h-4 text-secondary" fill="currentColor" stroke-width="2" />
            <span class="text-xs font-semibold text-text-primary">Продукт Empatra AI</span>
          </div>
          <p class="text-xs text-text-tertiary leading-relaxed">
            От холодных данных к теплым отношениям
          </p>
        </a>
      </div>
    </nav>
  </aside>
</template>
