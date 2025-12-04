<script setup lang="ts">
import { computed, ref } from 'vue'
import { cn } from '@/utils/cn'
import { Eye, EyeOff } from 'lucide-vue-next'

interface Props {
  modelValue?: string | number
  type?: string
  placeholder?: string
  label?: string
  error?: string
  hint?: string
  disabled?: boolean
  required?: boolean
  id?: string
  icon?: object
}

const props = withDefaults(defineProps<Props>(), {
  type: 'text',
})

const emit = defineEmits<{
  'update:modelValue': [value: string | number]
}>()

const inputRef = ref<HTMLInputElement | null>(null)
const showPassword = ref(false)
const isFocused = ref(false)

const inputType = computed(() => {
  if (props.type === 'password' && showPassword.value) return 'text'
  return props.type
})

const inputId = computed(() => 
  props.id || `input-${Math.random().toString(36).slice(2, 9)}`
)

const handleInput = (event: Event) => {
  const target = event.target as HTMLInputElement
  emit('update:modelValue', target.value)
}

const focus = () => inputRef.value?.focus()

defineExpose({ focus, inputRef })
</script>

<template>
  <div class="w-full">
    <!-- Label -->
    <label 
      v-if="label" 
      :for="inputId"
      class="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2"
    >
      {{ label }}
      <span v-if="required" class="text-error ml-0.5">*</span>
    </label>
    
    <!-- Input wrapper -->
    <div 
      class="relative group"
      :class="{ 'velvet-input': true }"
    >
      <!-- Icon -->
      <span 
        v-if="icon" 
        class="input-icon transition-colors duration-200"
        :class="isFocused ? 'text-primary' : 'text-gray-400'"
      >
        <component :is="icon" class="w-[18px] h-[18px]" />
      </span>
      
      <input
        :id="inputId"
        ref="inputRef"
        :type="inputType"
        :value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        :required="required"
        @input="handleInput"
        @focus="isFocused = true"
        @blur="isFocused = false"
        :class="cn(
          'w-full px-4 py-3 bg-white/60 dark:bg-slate-800/60',
          'border-2 border-transparent rounded-xl',
          'text-gray-800 dark:text-gray-100 placeholder:text-gray-400',
          'transition-all duration-200 ease-smooth',
          'shadow-inset-soft',
          'focus:bg-white dark:focus:bg-slate-800',
          'focus:border-primary',
          'focus:shadow-[0_0_0_4px_rgba(99,102,241,0.15),var(--shadow-raised)]',
          'disabled:opacity-50 disabled:cursor-not-allowed',
          error 
            ? 'border-error focus:border-error focus:shadow-[0_0_0_4px_rgba(239,68,68,0.15)]'
            : '',
          icon && 'pl-11',
          type === 'password' && 'pr-11'
        )"
      />
      
      <!-- Glow effect -->
      <div 
        class="absolute inset-0 rounded-xl pointer-events-none opacity-0 transition-opacity duration-300"
        :class="{ 'opacity-50': isFocused && !error }"
        style="background: radial-gradient(circle at center, rgba(99, 102, 241, 0.15), transparent 70%);"
      />
      
      <!-- Password toggle -->
      <button
        v-if="type === 'password'"
        type="button"
        @click="showPassword = !showPassword"
        class="absolute right-3 top-1/2 -translate-y-1/2 p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors rounded-lg hover:bg-black/[0.04] dark:hover:bg-white/[0.06]"
      >
        <EyeOff v-if="showPassword" class="w-5 h-5" />
        <Eye v-else class="w-5 h-5" />
      </button>
    </div>
    
    <!-- Error message -->
    <p 
      v-if="error" 
      class="mt-2 text-sm text-error flex items-center gap-1.5"
    >
      <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="10"/>
        <line x1="12" y1="8" x2="12" y2="12"/>
        <line x1="12" y1="16" x2="12.01" y2="16"/>
      </svg>
      {{ error }}
    </p>
    
    <!-- Hint -->
    <p v-else-if="hint" class="mt-2 text-sm text-gray-500 dark:text-gray-400">
      {{ hint }}
    </p>
  </div>
</template>

<style scoped>
.velvet-input .input-icon {
  position: absolute;
  left: 0.875rem;
  top: 50%;
  transform: translateY(-50%);
  pointer-events: none;
  z-index: 1;
}
</style>
