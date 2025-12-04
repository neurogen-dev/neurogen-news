<script setup lang="ts">
import Skeleton from '@/components/ui/Skeleton.vue'
</script>

<template>
  <article class="velvet-panel p-5 space-y-4">
    <!-- Header -->
    <div class="flex items-center gap-3">
      <Skeleton width="32px" height="32px" rounded="full" />
      <div class="flex-1 space-y-2">
        <Skeleton width="150px" height="14px" />
        <Skeleton width="100px" height="12px" />
      </div>
      <Skeleton width="90px" height="32px" rounded="lg" />
    </div>
    
    <!-- Badges -->
    <div class="flex gap-2">
      <Skeleton width="100px" height="24px" rounded="full" />
      <Skeleton width="80px" height="24px" rounded="full" />
      <Skeleton width="60px" height="24px" rounded="full" />
    </div>
    
    <!-- Title -->
    <Skeleton width="100%" height="24px" />
    
    <!-- Lead text -->
    <div class="space-y-2">
      <Skeleton width="100%" height="16px" />
      <Skeleton width="80%" height="16px" />
    </div>
    
    <!-- Cover image -->
    <Skeleton width="100%" height="200px" rounded="xl" />
    
    <!-- Actions -->
    <div class="pt-4 border-t border-gray-100 dark:border-gray-800 flex justify-between">
      <div class="flex gap-3">
        <Skeleton width="60px" height="28px" rounded="lg" />
        <Skeleton width="60px" height="28px" rounded="lg" />
        <Skeleton width="40px" height="28px" rounded="lg" />
      </div>
      <Skeleton width="50px" height="28px" rounded="lg" />
    </div>
  </article>
</template>
