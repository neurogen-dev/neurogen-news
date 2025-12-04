<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Award, Lock, Check } from 'lucide-vue-next'
import Badge from '@/components/ui/Badge.vue'
import { formatFullDate } from '@/utils/formatters'
import type { Achievement } from '@/types'

interface UserAchievement extends Achievement {
  isUnlocked: boolean
}

const achievements = ref<UserAchievement[]>([])
const isLoading = ref(true)

const mockAchievements: UserAchievement[] = [
  {
    id: '1',
    name: 'Первые шаги',
    description: 'Зарегистрироваться на платформе',
    icon: '🎓',
    rarity: 'common',
    unlockedAt: '2024-01-15T10:00:00Z',
    isUnlocked: true,
  },
  {
    id: '2',
    name: 'Первая статья',
    description: 'Опубликовать первую статью',
    icon: '✍️',
    rarity: 'common',
    unlockedAt: '2024-02-01T12:00:00Z',
    isUnlocked: true,
  },
  {
    id: '3',
    name: 'Комментатор',
    description: 'Оставить 10 комментариев',
    icon: '💬',
    rarity: 'uncommon',
    isUnlocked: false,
    progress: { current: 7, required: 10 },
  },
  {
    id: '4',
    name: 'Популярный автор',
    description: 'Набрать 1000 просмотров',
    icon: '🔥',
    rarity: 'rare',
    isUnlocked: false,
    progress: { current: 456, required: 1000 },
  },
  {
    id: '5',
    name: 'AI Эксперт',
    description: 'Опубликовать 10 статей',
    icon: '🧠',
    rarity: 'epic',
    isUnlocked: false,
    progress: { current: 3, required: 10 },
  },
  {
    id: '6',
    name: 'Легенда',
    description: 'Набрать 10000 кармы',
    icon: '👑',
    rarity: 'legendary',
    isUnlocked: false,
    progress: { current: 1234, required: 10000 },
  },
]

const unlockedCount = computed(() => 
  achievements.value.filter(a => a.isUnlocked).length
)

const rarityColors = {
  common: 'border-gray-300 dark:border-gray-600',
  uncommon: 'border-green-500',
  rare: 'border-blue-500',
  epic: 'border-purple-500',
  legendary: 'border-yellow-500',
}

const rarityLabels = {
  common: 'Обычное',
  uncommon: 'Необычное',
  rare: 'Редкое',
  epic: 'Эпическое',
  legendary: 'Легендарное',
}

onMounted(async () => {
  await new Promise(resolve => setTimeout(resolve, 500))
  achievements.value = mockAchievements
  isLoading.value = false
})
</script>

<template>
  <div class="min-h-screen">
    <!-- Header -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-text-primary dark:text-white flex items-center gap-2">
        <Award class="w-6 h-6 text-primary" />
        Достижения
      </h1>
      <p class="text-text-secondary mt-1">
        Получено {{ unlockedCount }} из {{ achievements.length }}
      </p>
    </div>
    
    <!-- Progress bar -->
    <div class="bg-white dark:bg-dark-secondary rounded-xl p-4 mb-6">
      <div class="flex justify-between text-sm mb-2">
        <span class="text-text-secondary">Прогресс</span>
        <span class="font-medium text-text-primary dark:text-white">
          {{ Math.round((unlockedCount / achievements.length) * 100) }}%
        </span>
      </div>
      <div class="h-2 bg-background-tertiary dark:bg-dark-tertiary rounded-full overflow-hidden">
        <div 
          class="h-full bg-primary rounded-full transition-all duration-500"
          :style="{ width: `${(unlockedCount / achievements.length) * 100}%` }"
        />
      </div>
    </div>
    
    <!-- Achievements grid -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="achievement in achievements"
        :key="achievement.id"
        class="relative bg-white dark:bg-dark-secondary rounded-xl border-2 p-5 transition-all"
        :class="[
          achievement.isUnlocked 
            ? rarityColors[achievement.rarity]
            : 'border-border dark:border-dark-tertiary opacity-60'
        ]"
      >
        <!-- Unlocked badge -->
        <div 
          v-if="achievement.isUnlocked"
          class="absolute top-3 right-3"
        >
          <div class="w-6 h-6 bg-success rounded-full flex items-center justify-center">
            <Check class="w-4 h-4 text-white" />
          </div>
        </div>
        
        <!-- Icon -->
        <div 
          class="text-5xl mb-3"
          :class="{ 'grayscale': !achievement.isUnlocked }"
        >
          {{ achievement.icon }}
        </div>
        
        <!-- Info -->
        <h3 class="font-bold text-text-primary dark:text-white mb-1">
          {{ achievement.name }}
        </h3>
        <p class="text-sm text-text-secondary mb-3">
          {{ achievement.description }}
        </p>
        
        <!-- Rarity badge -->
        <Badge 
          :variant="achievement.rarity === 'legendary' ? 'warning' : 'secondary'"
          class="mb-3"
        >
          {{ rarityLabels[achievement.rarity] }}
        </Badge>
        
        <!-- Progress or unlock date -->
        <div v-if="achievement.isUnlocked" class="text-xs text-text-tertiary">
          Получено {{ formatFullDate(achievement.unlockedAt!) }}
        </div>
        <div v-else-if="achievement.progress" class="mt-2">
          <div class="flex justify-between text-xs text-text-tertiary mb-1">
            <span>Прогресс</span>
            <span>{{ achievement.progress.current }} / {{ achievement.progress.required }}</span>
          </div>
          <div class="h-1.5 bg-background-tertiary dark:bg-dark-tertiary rounded-full overflow-hidden">
            <div 
              class="h-full bg-primary rounded-full"
              :style="{ width: `${(achievement.progress.current / achievement.progress.required) * 100}%` }"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

