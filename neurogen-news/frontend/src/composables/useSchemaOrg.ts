import { useHead } from '@unhead/vue'
import { computed, type MaybeRef, unref } from 'vue'

const SITE_NAME = 'Neurogen.News'
const BASE_URL = import.meta.env.VITE_BASE_URL || 'https://neurogen.news'
const LOGO_URL = `${BASE_URL}/logo.png`

// Organization schema (for site-wide)
export function useOrganizationSchema() {
  useHead({
    script: [
      {
        type: 'application/ld+json',
        innerHTML: JSON.stringify({
          '@context': 'https://schema.org',
          '@type': 'Organization',
          name: SITE_NAME,
          url: BASE_URL,
          logo: LOGO_URL,
          sameAs: [
            'https://t.me/neurogen_news',
            'https://vk.com/neurogen_news',
          ],
        }),
      },
    ],
  })
}

// Website schema (for homepage)
export function useWebsiteSchema() {
  useHead({
    script: [
      {
        type: 'application/ld+json',
        innerHTML: JSON.stringify({
          '@context': 'https://schema.org',
          '@type': 'WebSite',
          name: SITE_NAME,
          url: BASE_URL,
          potentialAction: {
            '@type': 'SearchAction',
            target: {
              '@type': 'EntryPoint',
              urlTemplate: `${BASE_URL}/search?q={search_term_string}`,
            },
            'query-input': 'required name=search_term_string',
          },
        }),
      },
    ],
  })
}

// Article schema
export interface ArticleSchemaData {
  title: string
  description?: string
  imageUrl?: string
  authorName: string
  authorUrl?: string
  publishedAt: string
  modifiedAt?: string
  url: string
  section?: string
  tags?: string[]
  wordCount?: number
}

export function useArticleSchema(article: MaybeRef<ArticleSchemaData | null>) {
  useHead({
    script: computed(() => {
      const a = unref(article)
      if (!a) return []

      return [
        {
          type: 'application/ld+json',
          innerHTML: JSON.stringify({
            '@context': 'https://schema.org',
            '@type': 'Article',
            headline: a.title,
            description: a.description,
            image: a.imageUrl,
            author: {
              '@type': 'Person',
              name: a.authorName,
              url: a.authorUrl,
            },
            publisher: {
              '@type': 'Organization',
              name: SITE_NAME,
              logo: {
                '@type': 'ImageObject',
                url: LOGO_URL,
              },
            },
            datePublished: a.publishedAt,
            dateModified: a.modifiedAt || a.publishedAt,
            mainEntityOfPage: {
              '@type': 'WebPage',
              '@id': a.url.startsWith('http') ? a.url : `${BASE_URL}${a.url}`,
            },
            articleSection: a.section,
            keywords: a.tags?.join(', '),
            wordCount: a.wordCount,
          }),
        },
      ]
    }),
  })
}

// HowTo schema for tutorials
export interface HowToStep {
  name: string
  text: string
  imageUrl?: string
}

export interface HowToSchemaData {
  name: string
  description: string
  imageUrl?: string
  totalTime?: string // ISO 8601 duration, e.g., "PT30M"
  steps: HowToStep[]
}

export function useHowToSchema(howTo: MaybeRef<HowToSchemaData | null>) {
  useHead({
    script: computed(() => {
      const h = unref(howTo)
      if (!h) return []

      return [
        {
          type: 'application/ld+json',
          innerHTML: JSON.stringify({
            '@context': 'https://schema.org',
            '@type': 'HowTo',
            name: h.name,
            description: h.description,
            image: h.imageUrl,
            totalTime: h.totalTime,
            step: h.steps.map((step, index) => ({
              '@type': 'HowToStep',
              position: index + 1,
              name: step.name,
              text: step.text,
              image: step.imageUrl,
            })),
          }),
        },
      ]
    }),
  })
}

// FAQ schema
export interface FAQItem {
  question: string
  answer: string
}

export function useFAQSchema(items: MaybeRef<FAQItem[]>) {
  useHead({
    script: computed(() => {
      const faqs = unref(items)
      if (!faqs.length) return []

      return [
        {
          type: 'application/ld+json',
          innerHTML: JSON.stringify({
            '@context': 'https://schema.org',
            '@type': 'FAQPage',
            mainEntity: faqs.map(faq => ({
              '@type': 'Question',
              name: faq.question,
              acceptedAnswer: {
                '@type': 'Answer',
                text: faq.answer,
              },
            })),
          }),
        },
      ]
    }),
  })
}

// Software Application schema (for tools)
export interface SoftwareSchemaData {
  name: string
  description: string
  imageUrl?: string
  operatingSystem?: string
  applicationCategory?: string
  aggregateRating?: {
    ratingValue: number
    ratingCount: number
  }
  offers?: {
    price: number | string
    priceCurrency: string
  }
}

export function useSoftwareSchema(software: MaybeRef<SoftwareSchemaData | null>) {
  useHead({
    script: computed(() => {
      const s = unref(software)
      if (!s) return []

      const schema: Record<string, unknown> = {
        '@context': 'https://schema.org',
        '@type': 'SoftwareApplication',
        name: s.name,
        description: s.description,
        image: s.imageUrl,
        operatingSystem: s.operatingSystem || 'Web',
        applicationCategory: s.applicationCategory || 'UtilitiesApplication',
      }

      if (s.aggregateRating) {
        schema.aggregateRating = {
          '@type': 'AggregateRating',
          ratingValue: s.aggregateRating.ratingValue,
          ratingCount: s.aggregateRating.ratingCount,
        }
      }

      if (s.offers) {
        schema.offers = {
          '@type': 'Offer',
          price: s.offers.price,
          priceCurrency: s.offers.priceCurrency,
        }
      }

      return [
        {
          type: 'application/ld+json',
          innerHTML: JSON.stringify(schema),
        },
      ]
    }),
  })
}

// Breadcrumb schema
export interface BreadcrumbItem {
  name: string
  url: string
}

export function useBreadcrumbSchema(items: MaybeRef<BreadcrumbItem[]>) {
  useHead({
    script: computed(() => {
      const breadcrumbs = unref(items)
      if (!breadcrumbs.length) return []

      return [
        {
          type: 'application/ld+json',
          innerHTML: JSON.stringify({
            '@context': 'https://schema.org',
            '@type': 'BreadcrumbList',
            itemListElement: breadcrumbs.map((item, index) => ({
              '@type': 'ListItem',
              position: index + 1,
              name: item.name,
              item: item.url.startsWith('http') ? item.url : `${BASE_URL}${item.url}`,
            })),
          }),
        },
      ]
    }),
  })
}

// Person schema (for author/profile pages)
export interface PersonSchemaData {
  name: string
  url: string
  imageUrl?: string
  description?: string
  sameAs?: string[]
}

export function usePersonSchema(person: MaybeRef<PersonSchemaData | null>) {
  useHead({
    script: computed(() => {
      const p = unref(person)
      if (!p) return []

      return [
        {
          type: 'application/ld+json',
          innerHTML: JSON.stringify({
            '@context': 'https://schema.org',
            '@type': 'Person',
            name: p.name,
            url: p.url.startsWith('http') ? p.url : `${BASE_URL}${p.url}`,
            image: p.imageUrl,
            description: p.description,
            sameAs: p.sameAs,
          }),
        },
      ]
    }),
  })
}

export default {
  useOrganizationSchema,
  useWebsiteSchema,
  useArticleSchema,
  useHowToSchema,
  useFAQSchema,
  useSoftwareSchema,
  useBreadcrumbSchema,
  usePersonSchema,
}


