<script setup lang="ts">
import { Check, Star, Zap, Shield, Crown } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Badge from '@/components/ui/Badge.vue'

const features = [
  {
    icon: '🚫',
    title: 'Без рекламы',
    description: 'Читайте статьи без отвлекающих баннеров',
  },
  {
    icon: '⭐',
    title: 'Эксклюзивный контент',
    description: 'Доступ к премиум-статьям и гайдам',
  },
  {
    icon: '💬',
    title: 'Приоритетные комментарии',
    description: 'Ваши комментарии выделяются',
  },
  {
    icon: '🎨',
    title: 'Кастомизация',
    description: 'Уникальные темы и настройки профиля',
  },
  {
    icon: '📊',
    title: 'Расширенная аналитика',
    description: 'Детальная статистика ваших статей',
  },
  {
    icon: '🛡️',
    title: 'Приоритетная поддержка',
    description: 'Быстрые ответы от команды',
  },
]

const plans = [
  {
    name: 'Месяц',
    price: 299,
    period: 'месяц',
    popular: false,
  },
  {
    name: 'Год',
    price: 2399,
    period: 'год',
    popular: true,
    discount: '-33%',
  },
]
</script>

<template>
  <div class="min-h-screen">
    <!-- Hero -->
    <section class="text-center mb-12">
      <Badge variant="warning" class="mb-4">
        <Crown class="w-4 h-4 mr-1" />
        Premium
      </Badge>
      <h1 class="text-4xl font-bold text-text-primary dark:text-white mb-4">
        Neurogen <span class="text-primary">Plus</span>
      </h1>
      <p class="text-xl text-text-secondary max-w-2xl mx-auto">
        Получите максимум от платформы с премиум-подпиской
      </p>
    </section>
    
    <!-- Features -->
    <section class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-12">
      <div 
        v-for="feature in features"
        :key="feature.title"
        class="p-6 bg-white dark:bg-dark-secondary rounded-xl"
      >
        <div class="text-3xl mb-3">{{ feature.icon }}</div>
        <h3 class="font-bold text-text-primary dark:text-white mb-2">
          {{ feature.title }}
        </h3>
        <p class="text-text-secondary">
          {{ feature.description }}
        </p>
      </div>
    </section>
    
    <!-- Pricing -->
    <section class="max-w-2xl mx-auto mb-12">
      <h2 class="text-2xl font-bold text-text-primary dark:text-white text-center mb-8">
        Выберите план
      </h2>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div
          v-for="plan in plans"
          :key="plan.name"
          class="relative p-6 bg-white dark:bg-dark-secondary rounded-2xl border-2 transition-all"
          :class="[
            plan.popular 
              ? 'border-primary shadow-lg scale-105' 
              : 'border-border dark:border-dark-tertiary'
          ]"
        >
          <Badge 
            v-if="plan.popular"
            variant="primary"
            class="absolute -top-3 left-1/2 -translate-x-1/2"
          >
            Популярный
          </Badge>
          
          <div class="text-center mb-6">
            <h3 class="text-lg font-bold text-text-primary dark:text-white mb-2">
              {{ plan.name }}
            </h3>
            <div class="flex items-baseline justify-center gap-1">
              <span class="text-4xl font-bold text-text-primary dark:text-white">
                {{ plan.price }}₽
              </span>
              <span class="text-text-tertiary">/ {{ plan.period }}</span>
            </div>
            <Badge v-if="plan.discount" variant="success" class="mt-2">
              {{ plan.discount }}
            </Badge>
          </div>
          
          <ul class="space-y-3 mb-6">
            <li class="flex items-center gap-2 text-text-secondary">
              <Check class="w-5 h-5 text-success" />
              Все премиум-функции
            </li>
            <li class="flex items-center gap-2 text-text-secondary">
              <Check class="w-5 h-5 text-success" />
              Отмена в любое время
            </li>
          </ul>
          
          <Button 
            :variant="plan.popular ? 'primary' : 'secondary'"
            class="w-full"
          >
            Оформить подписку
          </Button>
        </div>
      </div>
    </section>
    
    <!-- FAQ -->
    <section class="max-w-2xl mx-auto">
      <h2 class="text-2xl font-bold text-text-primary dark:text-white text-center mb-8">
        Частые вопросы
      </h2>
      
      <div class="space-y-4">
        <div class="p-4 bg-white dark:bg-dark-secondary rounded-xl">
          <h3 class="font-bold text-text-primary dark:text-white mb-2">
            Как оформить подписку?
          </h3>
          <p class="text-text-secondary">
            Выберите план, нажмите кнопку и следуйте инструкциям. 
            Мы принимаем банковские карты.
          </p>
        </div>
        
        <div class="p-4 bg-white dark:bg-dark-secondary rounded-xl">
          <h3 class="font-bold text-text-primary dark:text-white mb-2">
            Можно ли отменить подписку?
          </h3>
          <p class="text-text-secondary">
            Да, вы можете отменить подписку в любой момент в настройках аккаунта. 
            Доступ сохранится до конца оплаченного периода.
          </p>
        </div>
        
        <div class="p-4 bg-white dark:bg-dark-secondary rounded-xl">
          <h3 class="font-bold text-text-primary dark:text-white mb-2">
            Есть ли возврат средств?
          </h3>
          <p class="text-text-secondary">
            Мы предоставляем возврат в течение 7 дней после первой оплаты, 
            если вы не удовлетворены сервисом.
          </p>
        </div>
      </div>
    </section>
  </div>
</template>

