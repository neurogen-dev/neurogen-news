<script setup lang="ts">
import { ref } from 'vue'
import { X, Sparkles, BookOpen, Lightbulb, Rocket, ArrowRight, GraduationCap, Wand2 } from 'lucide-vue-next'
import { RouterLink } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const isVisible = ref(true)

const dismiss = () => {
  isVisible.value = false
  authStore.markStartHereSeen()
}

const quickGuides = [
  {
    icon: Sparkles,
    title: 'Что такое нейросети?',
    description: 'Простое объяснение для новичков',
    link: '/start/what-is-ai',
    gradient: 'from-primary/20 to-secondary/10',
    iconColor: 'text-primary',
  },
  {
    icon: BookOpen,
    title: 'С чего начать?',
    description: 'Пошаговый план изучения',
    link: '/start/getting-started',
    gradient: 'from-accent/20 to-primary/10',
    iconColor: 'text-accent',
  },
  {
    icon: Lightbulb,
    title: 'Лучшие промпты',
    description: 'Готовые запросы для ChatGPT',
    link: '/prompts',
    gradient: 'from-warning/20 to-secondary/10',
    iconColor: 'text-warning',
  },
  {
    icon: Rocket,
    title: 'Каталог инструментов',
    description: 'Все AI-сервисы в одном месте',
    link: '/tools',
    gradient: 'from-secondary/20 to-accent/10',
    iconColor: 'text-secondary',
  },
]
</script>

<template>
  <Transition
    enter-active-class="transition duration-500 ease-smooth"
    enter-from-class="opacity-0 scale-95"
    enter-to-class="opacity-100 scale-100"
    leave-active-class="transition duration-300 ease-smooth"
    leave-from-class="opacity-100 scale-100"
    leave-to-class="opacity-0 scale-95"
  >
    <section 
      v-if="isVisible"
      class="relative empatra-panel overflow-hidden"
    >
      <!-- Ambient gradient background -->
      <div class="absolute inset-0 overflow-hidden pointer-events-none z-0">
        <div class="absolute -top-1/2 -left-1/4 w-full h-full bg-gradient-to-br from-primary/20 via-transparent to-transparent rounded-full blur-3xl animate-pulse"></div>
        <div class="absolute -bottom-1/2 -right-1/4 w-3/4 h-3/4 bg-gradient-to-tl from-secondary/15 via-transparent to-transparent rounded-full blur-3xl animate-float"></div>
      </div>
      
      <!-- Content -->
      <div class="relative p-8 z-10">
        <!-- Close button -->
        <button
          @click="dismiss"
          class="absolute top-4 right-4 p-2.5 text-text-tertiary hover:text-text-primary transition-all duration-200 hover:bg-bg-hover rounded-xl"
          aria-label="Скрыть секцию"
        >
          <X class="w-5 h-5" stroke-width="2" />
        </button>
        
        <!-- Header -->
        <div class="text-center mb-10">
          <div class="flex items-center justify-center gap-3 mb-4">
            <Wand2 class="w-8 h-8 text-primary" stroke-width="2" />
            <h2 class="text-3xl md:text-4xl font-bold text-text-primary font-display">
              Добро пожаловать в 
              <span class="text-gradient">Neurogen.News</span>!
            </h2>
          </div>
          <p class="text-lg text-text-secondary max-w-2xl mx-auto leading-relaxed">
            Мы поможем вам освоить нейросети с нуля. Начните с простого — выберите интересную тему или пройдите базовое обучение.
          </p>
        </div>
        
        <!-- Quick guides grid -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-10">
          <RouterLink
            v-for="(guide, index) in quickGuides"
            :key="guide.link"
            :to="guide.link"
            class="group relative empatra-card p-5 hover:scale-[1.02] transition-all duration-300"
            :style="{ animationDelay: `${index * 0.1}s` }"
          >
            <!-- Card gradient overlay -->
            <div 
              :class="`absolute inset-0 bg-gradient-to-br ${guide.gradient} opacity-0 group-hover:opacity-100 transition-opacity duration-300 rounded-[inherit]`"
            ></div>
            
            <!-- Card content -->
            <div class="relative flex flex-col gap-4">
              <div class="p-3 rounded-xl bg-bg-surface/80 backdrop-blur-sm group-hover:scale-110 transition-transform duration-300 w-fit">
                <component 
                  :is="guide.icon" 
                  :class="['w-6 h-6', guide.iconColor]" 
                  stroke-width="2"
                />
              </div>
              <div class="flex-1 min-w-0">
                <h3 class="font-bold text-lg text-text-primary group-hover:text-primary transition-colors duration-200 flex items-center gap-2 mb-2 font-display">
                  {{ guide.title }}
                  <ArrowRight class="w-4 h-4 opacity-0 -translate-x-2 group-hover:opacity-100 group-hover:translate-x-0 transition-all duration-200" stroke-width="2" />
                </h3>
                <p class="text-sm text-text-secondary leading-relaxed">
                  {{ guide.description }}
                </p>
              </div>
            </div>
          </RouterLink>
        </div>
        
        <!-- CTA buttons -->
        <div class="flex flex-col sm:flex-row items-center justify-center gap-4">
          <Button 
            variant="primary"
            size="lg"
            to="/start"
          >
            <GraduationCap class="w-5 h-5" stroke-width="2" />
            <span>Начать обучение</span>
            <Sparkles class="w-4 h-4 opacity-60" stroke-width="2" />
          </Button>
          
          <Button 
            variant="ghost"
            size="lg"
            @click="dismiss"
          >
            Показать ленту
          </Button>
        </div>
      </div>
    </section>
  </Transition>
</template>
