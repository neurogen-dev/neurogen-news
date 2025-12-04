// Re-export all composables
export { useSeo, useArticleSeo, useProfileSeo, useCategorySeo } from './useSeo'
export {
  useOrganizationSchema,
  useWebsiteSchema,
  useArticleSchema,
  useHowToSchema,
  useFAQSchema,
  useSoftwareSchema,
  useBreadcrumbSchema,
  usePersonSchema,
} from './useSchemaOrg'
export type {
  ArticleSchemaData,
  HowToSchemaData,
  HowToStep,
  FAQItem,
  SoftwareSchemaData,
  BreadcrumbItem,
  PersonSchemaData,
} from './useSchemaOrg'


