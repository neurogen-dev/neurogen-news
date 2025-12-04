<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/utils/cn'

interface Props {
  icon?: any
  size?: 'sm' | 'md' | 'lg' | 'xl'
  variant?: 'primary' | 'secondary' | 'accent' | 'success' | 'warning' | 'danger' | 'glass'
  gradient?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  variant: 'primary',
  gradient: true,
})

const sizeClasses = {
  sm: 'w-8 h-8 p-1.5 rounded-lg',
  md: 'w-10 h-10 p-2 rounded-xl',
  lg: 'w-12 h-12 p-2.5 rounded-2xl',
  xl: 'w-16 h-16 p-4 rounded-3xl',
}

const iconSizes = {
  sm: 'w-4 h-4',
  md: 'w-5 h-5',
  lg: 'w-6 h-6',
  xl: 'w-8 h-8',
}

const variantStyles = {
  primary: 'text-primary-600 dark:text-primary-400 bg-primary-50 dark:bg-primary-900/20 border-primary-100 dark:border-primary-800',
  secondary: 'text-secondary-600 dark:text-secondary-400 bg-secondary-50 dark:bg-secondary-900/20 border-secondary-100 dark:border-secondary-800',
  accent: 'text-accent-600 dark:text-accent-400 bg-accent-50 dark:bg-accent-900/20 border-accent-100 dark:border-accent-800',
  success: 'text-emerald-600 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-900/20 border-emerald-100 dark:border-emerald-800',
  warning: 'text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-900/20 border-amber-100 dark:border-amber-800',
  danger: 'text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border-red-100 dark:border-red-800',
  glass: 'text-gray-700 dark:text-gray-200 bg-white/60 dark:bg-white/5 border-white/40 dark:border-white/10 backdrop-blur-md',
}

const gradientStyles = {
  primary: 'from-primary-100 to-primary-50/50 dark:from-primary-900/40 dark:to-primary-900/10',
  secondary: 'from-secondary-100 to-secondary-50/50 dark:from-secondary-900/40 dark:to-secondary-900/10',
  accent: 'from-accent-100 to-accent-50/50 dark:from-accent-900/40 dark:to-accent-900/10',
  success: 'from-emerald-100 to-emerald-50/50 dark:from-emerald-900/40 dark:to-emerald-900/10',
  warning: 'from-amber-100 to-amber-50/50 dark:from-amber-900/40 dark:to-amber-900/10',
  danger: 'from-red-100 to-red-50/50 dark:from-red-900/40 dark:to-red-900/10',
  glass: 'from-white/80 to-white/40 dark:from-white/10 dark:to-white/5',
}

const containerClasses = computed(() => cn(
  'relative flex items-center justify-center border shadow-sm transition-all duration-300 group-hover:scale-110 group-hover:shadow-md',
  sizeClasses[props.size],
  variantStyles[props.variant],
  props.gradient && 'bg-gradient-to-br',
  props.gradient && gradientStyles[props.variant]
))
</script>

<template>
  <div :class="containerClasses">
    <!-- Inner glow/shine effect -->
    <div class="absolute inset-0 rounded-[inherit] bg-gradient-to-br from-white/40 to-transparent opacity-50 dark:opacity-20 pointer-events-none"></div>
    
    <!-- Icon -->
    <component 
      :is="icon" 
      :class="cn('relative z-10 stroke-[1.5]', iconSizes[props.size])" 
    />
    
    <slot />
  </div>
</template>

