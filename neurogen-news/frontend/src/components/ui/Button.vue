<script setup lang="ts">
import { computed, ref } from 'vue'
import { cn } from '@/utils/cn'

interface Props {
  variant?: 'primary' | 'secondary' | 'glass' | 'soft' | 'ghost' | 'success' | 'danger' | 'subscribe'
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
  loading?: boolean
  as?: string | object
  to?: string
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  disabled: false,
  loading: false,
  as: 'button',
})

const buttonRef = ref<HTMLElement | null>(null)
const rippleStyle = ref({ left: '50%', top: '50%' })
const isRippling = ref(false)

const handleClick = (e: MouseEvent) => {
  if (!buttonRef.value || props.disabled || props.loading) return
  
  const rect = buttonRef.value.getBoundingClientRect()
  rippleStyle.value = {
    left: `${e.clientX - rect.left}px`,
    top: `${e.clientY - rect.top}px`,
  }
  
  isRippling.value = true
  setTimeout(() => {
    isRippling.value = false
  }, 600)
}

const baseClasses = `
  relative inline-flex items-center justify-center font-semibold
  transition-all duration-200 ease-smooth
  disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none
  overflow-hidden
  active:scale-[0.98]
`

const variantClasses = {
  primary: `
    bg-gradient-to-br from-primary via-primary to-primary-700
    text-white 
    shadow-raised shadow-glow-primary
    hover:shadow-floating hover:shadow-primary/40
    hover:-translate-y-0.5
  `,
  secondary: `
    bg-bg-surface
    text-text-primary
    border border-border
    hover:bg-bg-hover hover:border-primary/30
    hover:-translate-y-0.5
    shadow-raised
  `,
  glass: `
    bg-bg-surface/80
    backdrop-blur-xl
    text-text-primary
    border border-border-subtle
    shadow-raised
    hover:bg-bg-hover
    hover:border-primary/30
    hover:-translate-y-0.5
  `,
  soft: `
    bg-primary/15
    text-primary-400
    hover:bg-primary/25
    hover:-translate-y-0.5
  `,
  ghost: `
    text-text-secondary
    hover:text-text-primary
    hover:bg-bg-surface
  `,
  success: `
    bg-gradient-to-br from-success to-emerald-700
    text-white
    shadow-raised shadow-success/30
    hover:shadow-floating
    hover:-translate-y-0.5
  `,
  danger: `
    bg-gradient-to-br from-error to-rose-700
    text-white
    shadow-raised shadow-error/30
    hover:shadow-floating
    hover:-translate-y-0.5
  `,
  subscribe: `
    bg-primary/15
    text-primary-400
    border border-primary/25
    hover:bg-primary/25
    hover:border-primary/40
    hover:-translate-y-0.5
  `,
}

const sizeClasses = {
  sm: 'px-3 py-1.5 text-sm gap-1.5 rounded-lg',
  md: 'px-4 py-2 text-sm gap-2 rounded-xl',
  lg: 'px-6 py-3 text-base gap-2 rounded-xl',
}

const classes = computed(() => 
  cn(
    baseClasses,
    variantClasses[props.variant],
    sizeClasses[props.size],
    props.loading && 'opacity-70 cursor-wait'
  )
)

const component = computed(() => {
  if (props.to) return 'RouterLink'
  return props.as
})
</script>

<template>
  <component
    ref="buttonRef"
    :is="component"
    :to="to"
    :class="classes"
    :disabled="disabled || loading"
    @click="handleClick"
    v-bind="$attrs"
  >
    <!-- Highlight -->
    <span 
      v-if="['primary', 'success', 'danger'].includes(variant)"
      class="absolute top-0 left-0 right-0 h-1/2 bg-gradient-to-b from-white/20 to-transparent rounded-t-xl pointer-events-none z-0"
    />
    
    <!-- Ripple effect -->
    <span 
      v-if="isRippling"
      class="absolute w-0 h-0 bg-white/20 rounded-full transform -translate-x-1/2 -translate-y-1/2 pointer-events-none animate-ripple z-0"
      :style="rippleStyle"
    />
    
    <!-- Loading spinner -->
    <span 
      v-if="loading" 
      class="relative w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin z-10"
    />
    
    <!-- Content -->
    <span class="relative z-10 flex items-center gap-inherit">
      <slot />
    </span>
  </component>
</template>

<style scoped>
@keyframes ripple {
  to {
    width: 300%;
    padding-bottom: 300%;
    opacity: 0;
  }
}

.animate-ripple {
  animation: ripple 0.6s cubic-bezier(0.4, 0, 0.2, 1);
}
</style>
