<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { useBreadcrumbSchema } from '@/composables/useSchemaOrg'
import { ChevronRight, Home } from 'lucide-vue-next'

export interface BreadcrumbItem {
  label: string
  to?: string
  href?: string
}

const props = defineProps<{
  items: BreadcrumbItem[]
  showHome?: boolean
}>()

const breadcrumbsWithHome = computed(() => {
  const items = [...props.items]
  if (props.showHome !== false) {
    items.unshift({ label: 'Главная', to: '/' })
  }
  return items
})

// Add Schema.org markup
useBreadcrumbSchema(
  computed(() =>
    breadcrumbsWithHome.value
      .filter(item => item.to || item.href)
      .map(item => ({
        name: item.label,
        url: item.to || item.href || '/',
      }))
  )
)
</script>

<template>
  <nav aria-label="Breadcrumb" class="breadcrumbs">
    <ol class="breadcrumbs__list">
      <li
        v-for="(item, index) in breadcrumbsWithHome"
        :key="index"
        class="breadcrumbs__item"
      >
        <!-- Separator -->
        <ChevronRight
          v-if="index > 0"
          class="breadcrumbs__separator"
          :size="14"
        />

        <!-- Home icon -->
        <Home
          v-if="index === 0 && showHome !== false"
          class="breadcrumbs__home-icon"
          :size="14"
        />

        <!-- Link or text -->
        <RouterLink
          v-if="item.to && index < breadcrumbsWithHome.length - 1"
          :to="item.to"
          class="breadcrumbs__link"
        >
          {{ item.label }}
        </RouterLink>
        <a
          v-else-if="item.href && index < breadcrumbsWithHome.length - 1"
          :href="item.href"
          class="breadcrumbs__link"
        >
          {{ item.label }}
        </a>
        <span v-else class="breadcrumbs__current" aria-current="page">
          {{ item.label }}
        </span>
      </li>
    </ol>
  </nav>
</template>

<style scoped>
.breadcrumbs {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
}

.breadcrumbs__list {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.25rem;
  list-style: none;
  padding: 0;
  margin: 0;
}

.breadcrumbs__item {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.breadcrumbs__separator {
  color: var(--color-text-tertiary);
  flex-shrink: 0;
}

.breadcrumbs__home-icon {
  color: var(--color-text-secondary);
  margin-right: 0.125rem;
}

.breadcrumbs__link {
  color: var(--color-text-secondary);
  text-decoration: none;
  transition: color 0.15s ease;
}

.breadcrumbs__link:hover {
  color: var(--color-primary);
  text-decoration: underline;
}

.breadcrumbs__current {
  color: var(--color-text-primary);
  font-weight: 500;
}

/* Truncate long breadcrumbs on mobile */
@media (max-width: 640px) {
  .breadcrumbs__link,
  .breadcrumbs__current {
    max-width: 150px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}
</style>


