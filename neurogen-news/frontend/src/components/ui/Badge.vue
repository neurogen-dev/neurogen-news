<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/utils/cn'

interface Props {
  variant?: 'primary' | 'secondary' | 'success' | 'warning' | 'error' | 'beginner' | 'intermediate' | 'advanced'
  size?: 'sm' | 'md'
  shimmer?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'secondary',
  size: 'sm',
  shimmer: false,
})

const variantClasses = {
  primary: 'text-primary bg-primary/15 border-primary/30',
  secondary: 'text-text-secondary bg-bg-surface border-border-subtle',
  success: 'text-success bg-success/15 border-success/30',
  warning: 'text-warning bg-warning/15 border-warning/30',
  error: 'text-error bg-error/15 border-error/30',
  beginner: 'text-beginner bg-beginner/15 border-beginner/30',
  intermediate: 'text-intermediate bg-intermediate/15 border-intermediate/30',
  advanced: 'text-advanced bg-advanced/15 border-advanced/30',
}

const sizeClasses = {
  sm: 'px-2.5 py-0.5 text-xs',
  md: 'px-3 py-1 text-sm',
}

const classes = computed(() => 
  cn(
    'inline-flex items-center gap-1.5 font-medium rounded-full border',
    'transition-all duration-200',
    'hover:scale-[1.02]',
    variantClasses[props.variant],
    sizeClasses[props.size],
    props.shimmer && 'relative overflow-hidden'
  )
)
</script>

<template>
  <span :class="classes">
    <!-- Shimmer effect -->
    <span 
      v-if="shimmer"
      class="absolute inset-0 -translate-x-full animate-[shimmer_3s_ease-in-out_infinite]"
      style="background: linear-gradient(90deg, transparent, rgba(255,255,255,0.15), transparent);"
    />
    
    <span class="relative z-10 flex items-center gap-1.5">
      <slot />
    </span>
  </span>
</template>
