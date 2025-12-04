<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/utils/cn'

interface Props {
  src?: string | null
  alt?: string
  size?: number | 'sm' | 'md' | 'lg' | 'xl'
  fallback?: string
  ring?: boolean
  status?: 'online' | 'offline' | 'away' | 'busy'
}

const props = withDefaults(defineProps<Props>(), {
  size: 40,
  ring: false,
})

const sizeMap = {
  sm: 32,
  md: 40,
  lg: 56,
  xl: 80,
}

const resolvedSize = computed(() => 
  typeof props.size === 'string' ? sizeMap[props.size] : props.size
)

const initials = computed(() => {
  if (props.fallback) return props.fallback
  if (!props.alt) return '?'
  return props.alt
    .split(' ')
    .map(word => word[0])
    .join('')
    .slice(0, 2)
    .toUpperCase()
})

const statusColors = {
  online: 'bg-success',
  offline: 'bg-gray-400',
  away: 'bg-warning',
  busy: 'bg-error',
}

const statusSize = computed(() => Math.max(8, resolvedSize.value / 5))
</script>

<template>
  <div 
    class="relative inline-flex shrink-0"
    :style="{ width: `${resolvedSize}px`, height: `${resolvedSize}px` }"
  >
    <!-- Avatar image or fallback -->
    <div
      :class="cn(
        'w-full h-full rounded-full overflow-hidden',
        'bg-gradient-to-br from-primary/20 via-secondary/20 to-accent/20',
        'flex items-center justify-center',
        'transition-all duration-200',
        ring && 'ring-2 ring-primary/30 hover:ring-primary/50'
      )"
    >
      <img
        v-if="src"
        :src="src"
        :alt="alt"
        class="w-full h-full object-cover"
        loading="lazy"
      />
      <span 
        v-else
        class="text-primary font-semibold"
        :style="{ fontSize: `${resolvedSize / 2.5}px` }"
      >
        {{ initials }}
      </span>
    </div>
    
    <!-- Status indicator -->
    <span
      v-if="status"
      :class="cn(
        'absolute bottom-0 right-0 rounded-full',
        'ring-2 ring-white dark:ring-slate-900',
        'animate-pulse',
        statusColors[status]
      )"
      :style="{ 
        width: `${statusSize}px`, 
        height: `${statusSize}px`,
      }"
    />
  </div>
</template>
