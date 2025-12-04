<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { 
  ExternalLink, 
  Star, 
  Globe, 
  DollarSign, 
  ChevronLeft,
  Share2,
  Bookmark,
  MessageCircle
} from 'lucide-vue-next'
import Badge from '@/components/ui/Badge.vue'
import Button from '@/components/ui/Button.vue'

const route = useRoute()

const tool = ref({
  id: '1',
  name: 'ChatGPT',
  description: 'Универсальный AI-ассистент от OpenAI для генерации текста, кода и ответов на вопросы',
  fullDescription: `
ChatGPT — это передовая языковая модель от OpenAI, которая может помочь с широким спектром задач:

• Написание и редактирование текстов
• Генерация и отладка кода
• Анализ данных и документов
• Перевод текстов
• Ответы на вопросы
• Брейншторминг идей

ChatGPT доступен в бесплатной версии (GPT-3.5) и платной (GPT-4) с расширенными возможностями.
  `,
  icon: '🤖',
  url: 'https://chat.openai.com',
  category: 'Чат-боты',
  tags: ['текст', 'код', 'анализ', 'перевод'],
  isPremium: false,
  isFeatured: true,
  rating: 4.8,
  reviewCount: 1523,
  pricing: 'Freemium',
  platforms: ['Web', 'iOS', 'Android'],
})

const isLoading = ref(true)

onMounted(async () => {
  await new Promise(resolve => setTimeout(resolve, 500))
  isLoading.value = false
})
</script>

<template>
  <div class="min-h-screen">
    <!-- Back link -->
    <RouterLink 
      to="/tools"
      class="inline-flex items-center gap-1 text-text-secondary hover:text-primary mb-6"
    >
      <ChevronLeft class="w-4 h-4" />
      Каталог инструментов
    </RouterLink>
    
    <!-- Main content -->
    <div class="bg-white dark:bg-dark-secondary rounded-xl border border-border dark:border-dark-tertiary">
      <!-- Header -->
      <div class="p-6 border-b border-border dark:border-dark-tertiary">
        <div class="flex items-start gap-6">
          <div class="text-6xl">{{ tool.icon }}</div>
          
          <div class="flex-1">
            <div class="flex items-center gap-3 flex-wrap mb-2">
              <h1 class="text-2xl font-bold text-text-primary dark:text-white">
                {{ tool.name }}
              </h1>
              <Badge v-if="tool.isPremium" variant="warning">💎 Premium</Badge>
              <Badge v-if="tool.isFeatured" variant="primary">⭐ Рекомендуем</Badge>
            </div>
            
            <p class="text-text-secondary mb-4">
              {{ tool.description }}
            </p>
            
            <div class="flex items-center gap-4 flex-wrap">
              <div class="flex items-center gap-1">
                <Star class="w-5 h-5 text-yellow-500 fill-yellow-500" />
                <span class="font-bold text-text-primary dark:text-white">
                  {{ tool.rating }}
                </span>
                <span class="text-text-tertiary">
                  ({{ tool.reviewCount }} отзывов)
                </span>
              </div>
              
              <div class="flex items-center gap-1 text-text-secondary">
                <DollarSign class="w-4 h-4" />
                {{ tool.pricing }}
              </div>
              
              <div class="flex items-center gap-1 text-text-secondary">
                <Globe class="w-4 h-4" />
                {{ tool.platforms.join(', ') }}
              </div>
            </div>
          </div>
          
          <div class="flex flex-col gap-2">
            <Button as="a" :href="tool.url" target="_blank" rel="noopener">
              <ExternalLink class="w-4 h-4 mr-1" />
              Открыть
            </Button>
            <div class="flex gap-2">
              <Button variant="ghost" size="sm">
                <Bookmark class="w-4 h-4" />
              </Button>
              <Button variant="ghost" size="sm">
                <Share2 class="w-4 h-4" />
              </Button>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Tags -->
      <div class="px-6 py-4 border-b border-border dark:border-dark-tertiary">
        <div class="flex items-center gap-2 flex-wrap">
          <span class="text-sm text-text-tertiary">Теги:</span>
          <Badge 
            v-for="tag in tool.tags" 
            :key="tag"
            variant="secondary"
          >
            {{ tag }}
          </Badge>
        </div>
      </div>
      
      <!-- Description -->
      <div class="p-6">
        <h2 class="text-lg font-bold text-text-primary dark:text-white mb-4">
          Описание
        </h2>
        <div class="prose-article whitespace-pre-line text-text-secondary">
          {{ tool.fullDescription }}
        </div>
      </div>
    </div>
    
    <!-- Related articles -->
    <section class="mt-8">
      <h2 class="text-lg font-bold text-text-primary dark:text-white mb-4">
        Статьи о {{ tool.name }}
      </h2>
      <div class="text-center py-8 text-text-tertiary">
        Статьи скоро появятся
      </div>
    </section>
  </div>
</template>

