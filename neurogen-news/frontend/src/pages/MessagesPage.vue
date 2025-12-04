<script setup lang="ts">
import { ref } from 'vue'
import { MessageCircle, Search, Send, MoreHorizontal } from 'lucide-vue-next'
import Avatar from '@/components/ui/Avatar.vue'
import Button from '@/components/ui/Button.vue'
import { formatRelativeTime } from '@/utils/formatters'

interface Conversation {
  id: string
  user: {
    username: string
    displayName: string
    avatarUrl?: string
    isOnline: boolean
  }
  lastMessage: {
    content: string
    createdAt: string
    isRead: boolean
  }
}

const conversations = ref<Conversation[]>([])
const selectedConversation = ref<string | null>(null)
const messageInput = ref('')
const searchQuery = ref('')

// Empty state for now - real implementation would fetch from API
</script>

<template>
  <div class="min-h-screen -mx-4 -my-6">
    <div class="flex h-[calc(100vh-8rem)]">
      <!-- Conversations list -->
      <aside class="w-80 border-r border-border dark:border-dark-tertiary flex flex-col">
        <!-- Header -->
        <div class="p-4 border-b border-border dark:border-dark-tertiary">
          <h1 class="text-xl font-bold text-text-primary dark:text-white flex items-center gap-2">
            <MessageCircle class="w-5 h-5" />
            Сообщения
          </h1>
        </div>
        
        <!-- Search -->
        <div class="p-3 border-b border-border dark:border-dark-tertiary">
          <div class="relative">
            <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-text-tertiary" />
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Поиск..."
              class="w-full pl-9 pr-3 py-2 bg-background-secondary dark:bg-dark-tertiary rounded-lg text-sm"
            />
          </div>
        </div>
        
        <!-- Conversations -->
        <div class="flex-1 overflow-y-auto">
          <template v-if="conversations.length > 0">
            <button
              v-for="conv in conversations"
              :key="conv.id"
              @click="selectedConversation = conv.id"
              class="w-full flex items-center gap-3 p-3 hover:bg-background-secondary dark:hover:bg-dark-tertiary transition-colors"
              :class="{ 'bg-primary/5': selectedConversation === conv.id }"
            >
              <Avatar 
                :src="conv.user.avatarUrl"
                :alt="conv.user.displayName"
                :size="48"
                :show-online="true"
                :is-online="conv.user.isOnline"
              />
              <div class="flex-1 min-w-0 text-left">
                <div class="flex items-center justify-between">
                  <span class="font-medium text-text-primary dark:text-white truncate">
                    {{ conv.user.displayName }}
                  </span>
                  <span class="text-xs text-text-tertiary">
                    {{ formatRelativeTime(conv.lastMessage.createdAt) }}
                  </span>
                </div>
                <p 
                  class="text-sm truncate"
                  :class="conv.lastMessage.isRead ? 'text-text-tertiary' : 'text-text-primary dark:text-white font-medium'"
                >
                  {{ conv.lastMessage.content }}
                </p>
              </div>
            </button>
          </template>
          
          <!-- Empty state -->
          <div v-else class="flex flex-col items-center justify-center h-full p-6 text-center">
            <MessageCircle class="w-12 h-12 text-text-tertiary mb-3" />
            <p class="text-text-secondary">
              Нет сообщений
            </p>
          </div>
        </div>
      </aside>
      
      <!-- Chat area -->
      <main class="flex-1 flex flex-col">
        <template v-if="selectedConversation">
          <!-- Chat header -->
          <div class="flex items-center justify-between p-4 border-b border-border dark:border-dark-tertiary">
            <div class="flex items-center gap-3">
              <Avatar :size="40" alt="User" />
              <div>
                <h2 class="font-medium text-text-primary dark:text-white">
                  Пользователь
                </h2>
                <span class="text-xs text-success">В сети</span>
              </div>
            </div>
            <button class="p-2 text-text-tertiary hover:text-text-primary">
              <MoreHorizontal class="w-5 h-5" />
            </button>
          </div>
          
          <!-- Messages -->
          <div class="flex-1 overflow-y-auto p-4">
            <!-- Messages would go here -->
          </div>
          
          <!-- Input -->
          <div class="p-4 border-t border-border dark:border-dark-tertiary">
            <div class="flex gap-2">
              <input
                v-model="messageInput"
                type="text"
                placeholder="Введите сообщение..."
                class="flex-1 px-4 py-2 bg-background-secondary dark:bg-dark-tertiary rounded-lg"
                @keydown.enter="() => {}"
              />
              <Button>
                <Send class="w-4 h-4" />
              </Button>
            </div>
          </div>
        </template>
        
        <!-- No conversation selected -->
        <div v-else class="flex-1 flex items-center justify-center">
          <div class="text-center">
            <MessageCircle class="w-16 h-16 text-text-tertiary mx-auto mb-4" />
            <p class="text-text-secondary">
              Выберите диалог или начните новый
            </p>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

