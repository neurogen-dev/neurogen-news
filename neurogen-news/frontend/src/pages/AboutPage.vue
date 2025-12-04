<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { Users, Target, Heart, Zap, Globe, Shield } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'

const stats = [
  { value: '50K+', label: 'Читателей' },
  { value: '2K+', label: 'Статей' },
  { value: '100+', label: 'Авторов' },
  { value: '500+', label: 'AI инструментов' },
]

const values = [
  {
    icon: Target,
    title: 'Доступность',
    description: 'Делаем сложные технологии понятными для каждого',
  },
  {
    icon: Heart,
    title: 'Практичность',
    description: 'Фокус на реальном применении AI в жизни и работе',
  },
  {
    icon: Zap,
    title: 'Актуальность',
    description: 'Следим за новинками и быстро о них рассказываем',
  },
  {
    icon: Shield,
    title: 'Честность',
    description: 'Объективные обзоры без скрытой рекламы',
  },
]

const team = [
  { name: 'Команда редакции', role: 'Создание контента', avatar: '👨‍💻' },
  { name: 'AI-эксперты', role: 'Консультации', avatar: '🧠' },
  { name: 'Сообщество', role: 'Идеи и поддержка', avatar: '👥' },
]
</script>

<template>
  <div class="min-h-screen">
    <!-- Hero -->
    <section class="text-center mb-12">
      <div class="text-6xl mb-4">🧠</div>
      <h1 class="text-4xl font-bold text-text-primary dark:text-white mb-4">
        О <span class="text-primary">Neurogen</span><span class="text-text-secondary">.News</span>
      </h1>
      <p class="text-xl text-text-secondary max-w-2xl mx-auto">
        Делаем нейросети понятными и полезными для каждого. 
        Гайды, обзоры, новости и инструменты — всё в одном месте.
      </p>
    </section>
    
    <!-- Stats -->
    <section class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-12">
      <div 
        v-for="stat in stats" 
        :key="stat.label"
        class="text-center p-6 bg-white dark:bg-dark-secondary rounded-xl"
      >
        <div class="text-3xl font-bold text-primary mb-1">
          {{ stat.value }}
        </div>
        <div class="text-text-secondary">{{ stat.label }}</div>
      </div>
    </section>
    
    <!-- Mission -->
    <section class="bg-gradient-to-r from-primary/5 via-primary/10 to-purple-500/5 rounded-2xl p-8 mb-12">
      <div class="flex items-center gap-3 mb-4">
        <Globe class="w-6 h-6 text-primary" />
        <h2 class="text-2xl font-bold text-text-primary dark:text-white">
          Наша миссия
        </h2>
      </div>
      <p class="text-lg text-text-secondary leading-relaxed">
        Искусственный интеллект меняет мир прямо сейчас. Мы верим, что каждый должен 
        иметь возможность использовать эти технологии — независимо от технического 
        бэкграунда. Neurogen.News помогает разобраться в мире AI: от базовых понятий 
        до продвинутых техник, от бесплатных инструментов до профессиональных решений.
      </p>
    </section>
    
    <!-- Values -->
    <section class="mb-12">
      <h2 class="text-2xl font-bold text-text-primary dark:text-white text-center mb-8">
        Наши принципы
      </h2>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div 
          v-for="value in values"
          :key="value.title"
          class="flex gap-4 p-6 bg-white dark:bg-dark-secondary rounded-xl"
        >
          <div class="w-12 h-12 bg-primary/10 rounded-xl flex items-center justify-center shrink-0">
            <component :is="value.icon" class="w-6 h-6 text-primary" />
          </div>
          <div>
            <h3 class="font-bold text-text-primary dark:text-white mb-1">
              {{ value.title }}
            </h3>
            <p class="text-text-secondary">
              {{ value.description }}
            </p>
          </div>
        </div>
      </div>
    </section>
    
    <!-- Team -->
    <section class="mb-12">
      <h2 class="text-2xl font-bold text-text-primary dark:text-white text-center mb-8">
        Кто мы
      </h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div 
          v-for="member in team"
          :key="member.name"
          class="text-center p-6 bg-white dark:bg-dark-secondary rounded-xl"
        >
          <div class="text-5xl mb-3">{{ member.avatar }}</div>
          <h3 class="font-bold text-text-primary dark:text-white">
            {{ member.name }}
          </h3>
          <p class="text-text-secondary">{{ member.role }}</p>
        </div>
      </div>
    </section>
    
    <!-- CTA -->
    <section class="text-center py-8 px-6 bg-white dark:bg-dark-secondary rounded-2xl">
      <Users class="w-10 h-10 text-primary mx-auto mb-4" />
      <h2 class="text-xl font-bold text-text-primary dark:text-white mb-2">
        Присоединяйтесь к сообществу
      </h2>
      <p class="text-text-secondary mb-6">
        Станьте частью растущего сообщества AI-энтузиастов
      </p>
      <div class="flex justify-center gap-3">
        <Button as="RouterLink" to="/register">
          Создать аккаунт
        </Button>
        <Button as="a" href="https://t.me/neurogen_news" target="_blank" variant="secondary">
          Telegram канал
        </Button>
      </div>
    </section>
  </div>
</template>

