import { useHead } from '@unhead/vue'
import { computed, type MaybeRef, unref } from 'vue'

export interface SeoMeta {
  title?: string
  description?: string
  image?: string
  url?: string
  type?: 'website' | 'article' | 'profile'
  author?: string
  publishedTime?: string
  modifiedTime?: string
  section?: string
  tags?: string[]
  canonical?: string
  robots?: string
  locale?: string
}

const SITE_NAME = 'Neurogen.News'
const DEFAULT_DESCRIPTION = 'Технологический портал о нейросетях, искусственном интеллекте и машинном обучении'
const DEFAULT_IMAGE = '/og-image.png'
const BASE_URL = import.meta.env.VITE_BASE_URL || 'https://neurogen.news'

export function useSeo(meta: MaybeRef<SeoMeta>) {
  const computedMeta = computed(() => unref(meta))

  useHead({
    title: computed(() => {
      const title = computedMeta.value.title
      return title ? `${title} — ${SITE_NAME}` : SITE_NAME
    }),
    meta: computed(() => {
      const m = computedMeta.value
      const title = m.title || SITE_NAME
      const description = m.description || DEFAULT_DESCRIPTION
      const image = m.image || DEFAULT_IMAGE
      const url = m.url || BASE_URL
      const type = m.type || 'website'
      const locale = m.locale || 'ru_RU'

      const metas: Array<{ name?: string; property?: string; content: string }> = [
        // Basic
        { name: 'description', content: description },

        // Open Graph
        { property: 'og:type', content: type },
        { property: 'og:title', content: title },
        { property: 'og:description', content: description },
        { property: 'og:image', content: image.startsWith('http') ? image : `${BASE_URL}${image}` },
        { property: 'og:url', content: url.startsWith('http') ? url : `${BASE_URL}${url}` },
        { property: 'og:site_name', content: SITE_NAME },
        { property: 'og:locale', content: locale },

        // Twitter
        { name: 'twitter:card', content: 'summary_large_image' },
        { name: 'twitter:title', content: title },
        { name: 'twitter:description', content: description },
        { name: 'twitter:image', content: image.startsWith('http') ? image : `${BASE_URL}${image}` },
      ]

      // Article-specific
      if (type === 'article') {
        if (m.author) {
          metas.push({ property: 'article:author', content: m.author })
        }
        if (m.publishedTime) {
          metas.push({ property: 'article:published_time', content: m.publishedTime })
        }
        if (m.modifiedTime) {
          metas.push({ property: 'article:modified_time', content: m.modifiedTime })
        }
        if (m.section) {
          metas.push({ property: 'article:section', content: m.section })
        }
        if (m.tags?.length) {
          m.tags.forEach(tag => {
            metas.push({ property: 'article:tag', content: tag })
          })
        }
      }

      // Robots
      if (m.robots) {
        metas.push({ name: 'robots', content: m.robots })
      }

      return metas
    }),
    link: computed(() => {
      const links: Array<{ rel: string; href: string }> = []
      const m = computedMeta.value

      // Canonical
      if (m.canonical) {
        links.push({ rel: 'canonical', href: m.canonical })
      } else if (m.url) {
        links.push({
          rel: 'canonical',
          href: m.url.startsWith('http') ? m.url : `${BASE_URL}${m.url}`,
        })
      }

      return links
    }),
  })
}

// Pre-configured SEO for article pages
export function useArticleSeo(article: MaybeRef<{
  title: string
  lead?: string
  coverImageUrl?: string
  author: { displayName: string }
  category: { name: string; slug: string }
  tags?: { name: string }[]
  publishedAt?: string
  updatedAt?: string
  slug: string
} | null>) {
  useSeo(computed(() => {
    const a = unref(article)
    if (!a) return {}

    return {
      title: a.title,
      description: a.lead,
      image: a.coverImageUrl,
      type: 'article' as const,
      author: a.author?.displayName,
      publishedTime: a.publishedAt,
      modifiedTime: a.updatedAt,
      section: a.category?.name,
      tags: a.tags?.map(t => t.name),
      url: `/${a.category?.slug}/${a.slug}`,
    }
  }))
}

// Pre-configured SEO for user profile pages
export function useProfileSeo(user: MaybeRef<{
  displayName: string
  username: string
  bio?: string
  avatarUrl?: string
} | null>) {
  useSeo(computed(() => {
    const u = unref(user)
    if (!u) return {}

    return {
      title: u.displayName,
      description: u.bio || `Профиль пользователя ${u.displayName}`,
      image: u.avatarUrl,
      type: 'profile' as const,
      url: `/u/${u.username}`,
    }
  }))
}

// Pre-configured SEO for category pages
export function useCategorySeo(category: MaybeRef<{
  name: string
  description?: string
  coverUrl?: string
  slug: string
} | null>) {
  useSeo(computed(() => {
    const c = unref(category)
    if (!c) return {}

    return {
      title: c.name,
      description: c.description,
      image: c.coverUrl,
      url: `/${c.slug}`,
    }
  }))
}

export default useSeo


