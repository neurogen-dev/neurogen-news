<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { Home, ArrowLeft, Search } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
</script>

<template>
  <div class="min-h-[60vh] flex items-center justify-center">
    <div class="text-center max-w-md">
      <!-- 404 illustration -->
      <div class="text-8xl mb-6">🤖</div>
      
      <h1 class="text-4xl font-bold text-text-primary dark:text-white mb-4">
        404
      </h1>
      
      <h2 class="text-xl font-medium text-text-secondary mb-2">
        Страница не найдена
      </h2>
      
      <p class="text-text-tertiary mb-8">
        Похоже, нейросети не смогли найти эту страницу. 
        Возможно, она была удалена или вы перешли по неверной ссылке.
      </p>
      
      <div class="flex flex-col sm:flex-row items-center justify-center gap-3">
        <Button as="RouterLink" to="/">
          <Home class="w-4 h-4 mr-2" />
          На главную
        </Button>
        
        <Button 
          variant="secondary" 
          @click="$router.back()"
        >
          <ArrowLeft class="w-4 h-4 mr-2" />
          Назад
        </Button>
        
        <Button 
          as="RouterLink" 
          to="/search"
          variant="ghost"
        >
          <Search class="w-4 h-4 mr-2" />
          Поиск
        </Button>
      </div>
      
      <!-- Popular links -->
      <div class="mt-12 pt-8 border-t border-border dark:border-dark-tertiary">
        <p class="text-sm text-text-tertiary mb-4">
          Возможно, вас заинтересует:
        </p>
        <div class="flex flex-wrap justify-center gap-2">
          <RouterLink 
            to="/start"
            class="px-3 py-1.5 text-sm bg-background-secondary dark:bg-dark-tertiary rounded-full text-text-secondary hover:text-primary transition-colors"
          >
            🎓 Начать здесь
          </RouterLink>
          <RouterLink 
            to="/popular"
            class="px-3 py-1.5 text-sm bg-background-secondary dark:bg-dark-tertiary rounded-full text-text-secondary hover:text-primary transition-colors"
          >
            🔥 Популярное
          </RouterLink>
          <RouterLink 
            to="/tools"
            class="px-3 py-1.5 text-sm bg-background-secondary dark:bg-dark-tertiary rounded-full text-text-secondary hover:text-primary transition-colors"
          >
            🛠️ Инструменты
          </RouterLink>
          <RouterLink 
            to="/prompts"
            class="px-3 py-1.5 text-sm bg-background-secondary dark:bg-dark-tertiary rounded-full text-text-secondary hover:text-primary transition-colors"
          >
            💡 Промпты
          </RouterLink>
        </div>
      </div>
    </div>
  </div>
</template>

