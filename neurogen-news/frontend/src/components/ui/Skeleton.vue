<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/utils/cn'

interface Props {
  width?: string | number
  height?: string | number
  rounded?: 'sm' | 'md' | 'lg' | 'xl' | '2xl' | 'full'
  variant?: 'default' | 'text' | 'circular' | 'rectangular'
}

const props = withDefaults(defineProps<Props>(), {
  rounded: 'lg',
  variant: 'default',
})

const roundedClasses = {
  sm: 'rounded-sm',
  md: 'rounded-md',
  lg: 'rounded-lg',
  xl: 'rounded-xl',
  '2xl': 'rounded-2xl',
  full: 'rounded-full',
}

const variantClasses = {
  default: '',
  text: 'h-4',
  circular: 'rounded-full',
  rectangular: 'rounded-none',
}

const style = computed(() => {
  const s: Record<string, string> = {}
  
  if (props.width) {
    s.width = typeof props.width === 'number' ? `${props.width}px` : props.width
  }
  
  if (props.height) {
    s.height = typeof props.height === 'number' ? `${props.height}px` : props.height
  }
  
  return s
})

const classes = computed(() => cn(
  'animate-pulse',
  'bg-gradient-to-r from-gray-200 via-gray-100 to-gray-200',
  'dark:from-slate-700 dark:via-slate-600 dark:to-slate-700',
  'bg-[length:200%_100%]',
  'animate-shimmer',
  roundedClasses[props.rounded],
  variantClasses[props.variant],
  props.variant === 'circular' && 'aspect-square'
))
</script>

<template>
  <div :class="classes" :style="style">
    <slot />
  </div>
</template>

<style scoped>
.animate-shimmer {
  animation: shimmer 1.5s ease-in-out infinite;
}

@keyframes shimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}
</style>
