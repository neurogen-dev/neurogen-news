<script setup lang="ts">
import ArticleCard from './ArticleCard.vue'
import ArticleCardSkeleton from './ArticleCardSkeleton.vue'
import { useInfiniteScroll } from '@/composables/useInfiniteScroll'
import type { ArticleCard as ArticleCardType } from '@/types'

interface Props {
  articles: ArticleCardType[]
  isLoading?: boolean
  isLoadingMore?: boolean
  hasMore?: boolean
  compact?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  isLoading: false,
  isLoadingMore: false,
  hasMore: false,
  compact: false,
})

const emit = defineEmits<{
  loadMore: []
}>()

const { target } = useInfiniteScroll(async () => {
  if (props.hasMore && !props.isLoadingMore) {
    emit('loadMore')
  }
})
</script>

<template>
  <div class="space-y-4">
    <!-- Loading skeleton -->
    <template v-if="isLoading">
      <ArticleCardSkeleton v-for="i in 5" :key="i" />
    </template>
    
    <!-- Articles -->
    <template v-else>
      <ArticleCard
        v-for="(article, index) in articles"
        :key="article.id"
        :article="article"
        :compact="compact"
        :index="index"
      />
      
      <!-- Load more trigger -->
      <div ref="target" class="h-1" />
      
      <!-- Loading more indicator -->
      <div v-if="isLoadingMore" class="flex justify-center py-6">
        <div class="relative w-10 h-10">
          <div class="absolute inset-0 rounded-full border-2 border-primary/20"></div>
          <div class="absolute inset-0 rounded-full border-2 border-primary border-t-transparent animate-spin"></div>
        </div>
      </div>
      
      <!-- No more articles -->
      <div 
        v-else-if="!hasMore && articles.length > 0" 
        class="text-center py-8"
      >
        <p class="text-gray-400 dark:text-gray-500">
          ✨ Вы достигли конца ленты
        </p>
      </div>
      
      <!-- Empty state -->
      <div 
        v-if="articles.length === 0" 
        class="velvet-panel text-center py-16"
      >
        <div class="text-6xl mb-4">📭</div>
        <h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">
          Пока ничего нет
        </h3>
        <p class="text-gray-500 dark:text-gray-400">
          Здесь появятся статьи, когда они будут опубликованы
        </p>
      </div>
    </template>
  </div>
</template>
