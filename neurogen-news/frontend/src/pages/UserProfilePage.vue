<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { 
  Calendar, 
  MapPin, 
  Link as LinkIcon, 
  Twitter,
  Users,
  FileText,
  MessageCircle,
  Settings,
  Award,
  MoreHorizontal
} from 'lucide-vue-next'
import Avatar from '@/components/ui/Avatar.vue'
import Badge from '@/components/ui/Badge.vue'
import Button from '@/components/ui/Button.vue'
import FeedList from '@/components/feed/FeedList.vue'
import { useAuthStore } from '@/stores/auth'
import { formatFullDate, formatCompactNumber, formatKarma } from '@/utils/formatters'
import type { User, ArticleCard } from '@/types'

const route = useRoute()
const authStore = useAuthStore()

const user = ref<User | null>(null)
const articles = ref<ArticleCard[]>([])
const isLoading = ref(true)
const isFollowing = ref(false)
const activeTab = ref<'articles' | 'comments'>('articles')

const isOwnProfile = computed(() => 
  authStore.user?.username === route.params.username
)

const loadProfile = async () => {
  isLoading.value = true
  
  // Simulate API call
  await new Promise(resolve => setTimeout(resolve, 500))
  
  user.value = {
    id: '1',
    username: route.params.username as string,
    displayName: 'Пользователь Neurogen',
    email: 'user@example.com',
    bio: 'Энтузиаст AI и машинного обучения. Пишу о нейросетях простым языком.',
    avatarUrl: undefined,
    role: 'AUTHOR',
    karma: 1234,
    isVerified: true,
    isPremium: false,
    createdAt: '2024-01-15T10:00:00Z',
    followerCount: 567,
    followingCount: 123,
    articleCount: 42,
  }
  
  isLoading.value = false
}

const handleFollow = () => {
  if (!authStore.isLoggedIn) return
  isFollowing.value = !isFollowing.value
}

onMounted(loadProfile)
</script>

<template>
  <div class="min-h-screen">
    <!-- Loading -->
    <div v-if="isLoading" class="animate-pulse space-y-4">
      <div class="h-32 bg-background-tertiary dark:bg-dark-tertiary rounded-xl" />
      <div class="flex gap-4">
        <div class="w-24 h-24 rounded-full bg-background-tertiary dark:bg-dark-tertiary" />
        <div class="flex-1 space-y-2">
          <div class="h-6 w-48 bg-background-tertiary dark:bg-dark-tertiary rounded" />
          <div class="h-4 w-32 bg-background-tertiary dark:bg-dark-tertiary rounded" />
        </div>
      </div>
    </div>
    
    <!-- Profile -->
    <template v-else-if="user">
      <!-- Cover -->
      <div class="h-32 bg-gradient-to-r from-primary/20 via-primary/10 to-purple-500/20 rounded-xl mb-4" />
      
      <!-- User info -->
      <div class="relative px-4 -mt-16 mb-6">
        <div class="flex flex-col sm:flex-row sm:items-end gap-4">
          <Avatar 
            :src="user.avatarUrl" 
            :alt="user.displayName" 
            :size="96"
            class="ring-4 ring-white dark:ring-dark-secondary"
          />
          
          <div class="flex-1">
            <div class="flex items-center gap-2 flex-wrap">
              <h1 class="text-2xl font-bold text-text-primary dark:text-white">
                {{ user.displayName }}
              </h1>
              <Badge v-if="user.isVerified" variant="primary">✓</Badge>
              <Badge v-if="user.isPremium" variant="warning">💎 Plus</Badge>
              <Badge :variant="user.role === 'ADMIN' ? 'error' : 'secondary'">
                {{ user.role }}
              </Badge>
            </div>
            <p class="text-text-tertiary">@{{ user.username }}</p>
          </div>
          
          <div class="flex gap-2">
            <template v-if="isOwnProfile">
              <Button as="RouterLink" to="/settings" variant="secondary">
                <Settings class="w-4 h-4 mr-1" />
                Настройки
              </Button>
            </template>
            <template v-else>
              <Button 
                :variant="isFollowing ? 'secondary' : 'primary'"
                @click="handleFollow"
              >
                {{ isFollowing ? '✓ Подписка' : 'Подписаться' }}
              </Button>
              <Button variant="ghost">
                <MoreHorizontal class="w-4 h-4" />
              </Button>
            </template>
          </div>
        </div>
        
        <!-- Bio -->
        <p v-if="user.bio" class="mt-4 text-text-secondary">
          {{ user.bio }}
        </p>
        
        <!-- Meta -->
        <div class="flex flex-wrap items-center gap-4 mt-4 text-sm text-text-tertiary">
          <span class="flex items-center gap-1">
            <Calendar class="w-4 h-4" />
            На сайте с {{ formatFullDate(user.createdAt) }}
          </span>
          <span 
            class="font-medium"
            :class="{
              'text-success': user.karma > 0,
              'text-error': user.karma < 0,
            }"
          >
            Карма: {{ formatKarma(user.karma) }}
          </span>
        </div>
        
        <!-- Stats -->
        <div class="flex gap-6 mt-4">
          <RouterLink 
            :to="`/@${user.username}/followers`"
            class="hover:text-primary"
          >
            <span class="font-bold text-text-primary dark:text-white">
              {{ formatCompactNumber(user.followerCount) }}
            </span>
            <span class="text-text-tertiary ml-1">подписчиков</span>
          </RouterLink>
          <RouterLink 
            :to="`/@${user.username}/following`"
            class="hover:text-primary"
          >
            <span class="font-bold text-text-primary dark:text-white">
              {{ formatCompactNumber(user.followingCount) }}
            </span>
            <span class="text-text-tertiary ml-1">подписок</span>
          </RouterLink>
          <span>
            <span class="font-bold text-text-primary dark:text-white">
              {{ formatCompactNumber(user.articleCount) }}
            </span>
            <span class="text-text-tertiary ml-1">статей</span>
          </span>
        </div>
      </div>
      
      <!-- Tabs -->
      <div class="flex gap-1 border-b border-border dark:border-dark-tertiary mb-6">
        <button
          @click="activeTab = 'articles'"
          class="px-4 py-3 text-sm font-medium border-b-2 -mb-px transition-colors"
          :class="[
            activeTab === 'articles'
              ? 'border-primary text-primary'
              : 'border-transparent text-text-secondary hover:text-text-primary'
          ]"
        >
          <FileText class="w-4 h-4 inline mr-1" />
          Статьи
        </button>
        <button
          @click="activeTab = 'comments'"
          class="px-4 py-3 text-sm font-medium border-b-2 -mb-px transition-colors"
          :class="[
            activeTab === 'comments'
              ? 'border-primary text-primary'
              : 'border-transparent text-text-secondary hover:text-text-primary'
          ]"
        >
          <MessageCircle class="w-4 h-4 inline mr-1" />
          Комментарии
        </button>
      </div>
      
      <!-- Content -->
      <FeedList 
        :articles="articles" 
        :is-loading="false" 
        :has-more="false"
      />
    </template>
  </div>
</template>

