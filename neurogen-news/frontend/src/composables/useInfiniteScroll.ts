import { ref, onMounted, onUnmounted } from 'vue'

interface UseInfiniteScrollOptions {
  threshold?: number
  rootMargin?: string
}

export function useInfiniteScroll(
  callback: () => Promise<void> | void,
  options: UseInfiniteScrollOptions = {}
) {
  const { threshold = 0.1, rootMargin = '200px' } = options
  
  const target = ref<HTMLElement | null>(null)
  const isLoading = ref(false)
  
  let observer: IntersectionObserver | null = null
  
  const handleIntersect = async (entries: IntersectionObserverEntry[]) => {
    const [entry] = entries
    
    if (entry.isIntersecting && !isLoading.value) {
      isLoading.value = true
      try {
        await callback()
      } finally {
        isLoading.value = false
      }
    }
  }
  
  onMounted(() => {
    observer = new IntersectionObserver(handleIntersect, {
      threshold,
      rootMargin,
    })
    
    if (target.value) {
      observer.observe(target.value)
    }
  })
  
  onUnmounted(() => {
    if (observer) {
      observer.disconnect()
    }
  })
  
  return {
    target,
    isLoading,
  }
}
