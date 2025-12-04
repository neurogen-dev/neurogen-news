# SEO-стратегия Neurogen.News

## Цель

Вывести сайт в топ-10 поисковых систем (Яндекс, Google) по ключевым запросам о нейросетях и ИИ за **6-12 месяцев** с нуля.

---

## 1. Техническое SEO

### 1.1 Серверный рендеринг (SSR/SSG)

**Критически важно для индексации!**

```
┌─────────────────────────────────────────────────────────────────┐
│  Стратегия рендеринга                                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  SSG (Static Site Generation):                                  │
│  ├── Главная страница                                           │
│  ├── Страницы подсайтов (категорий)                            │
│  ├── Статические страницы (О проекте, Правила)                 │
│  └── Каталог инструментов                                      │
│                                                                 │
│  SSR (Server-Side Rendering):                                   │
│  ├── Страницы статей (динамический контент)                    │
│  ├── Профили пользователей                                      │
│  ├── Результаты поиска                                          │
│  └── Ленты (популярное, свежее)                                │
│                                                                 │
│  SPA (Client-Side):                                             │
│  ├── Редактор статей                                            │
│  ├── Настройки пользователя                                     │
│  └── Личный кабинет                                             │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 1.2 Core Web Vitals

| Метрика | Целевое значение | Как достигнуть |
|---------|------------------|----------------|
| **LCP** (Largest Contentful Paint) | < 2.5s | Оптимизация изображений, preload критичных ресурсов, CDN |
| **FID** (First Input Delay) | < 100ms | Code splitting, отложенная загрузка JS |
| **CLS** (Cumulative Layout Shift) | < 0.1 | Резервирование места для изображений, skeleton loaders |
| **TTFB** (Time to First Byte) | < 200ms | Edge caching, оптимизация бэкенда |
| **INP** (Interaction to Next Paint) | < 200ms | Оптимизация обработчиков событий |

### 1.3 Техническая оптимизация

```go
// Пример middleware для SEO-заголовков на Go
func SEOMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Canonical URL
        w.Header().Set("Link", fmt.Sprintf("<%s>; rel=\"canonical\"", r.URL.String()))
        
        // Cache headers для статики
        if isStaticAsset(r.URL.Path) {
            w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
        }
        
        // Security headers (влияют на доверие)
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "SAMEORIGIN")
        w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        
        next.ServeHTTP(w, r)
    })
}
```

#### Чек-лист технического SEO

- [ ] **HTTPS** — обязательно, с HSTS
- [ ] **Robots.txt** — правильная настройка
- [ ] **Sitemap.xml** — динамическая генерация
- [ ] **Canonical URLs** — на каждой странице
- [ ] **Hreflang** — если будет многоязычность
- [ ] **Mobile-first** — адаптивный дизайн
- [ ] **Скорость загрузки** — Core Web Vitals в зелёной зоне
- [ ] **Structured Data** — Schema.org разметка
- [ ] **Open Graph / Twitter Cards** — для социальных сетей
- [ ] **AMP** (опционально) — для новостных статей

---

## 2. Структурированные данные (Schema.org)

### 2.1 Типы разметки

```json
// Article (для статей)
{
  "@context": "https://schema.org",
  "@type": "Article",
  "headline": "Как начать пользоваться ChatGPT за 5 минут",
  "description": "Пошаговый гайд для новичков...",
  "image": "https://neurogen.news/images/article-cover.jpg",
  "author": {
    "@type": "Person",
    "name": "Иван Петров",
    "url": "https://neurogen.news/@ivan"
  },
  "publisher": {
    "@type": "Organization",
    "name": "Neurogen.News",
    "logo": {
      "@type": "ImageObject",
      "url": "https://neurogen.news/logo.png"
    }
  },
  "datePublished": "2025-12-01T10:00:00+03:00",
  "dateModified": "2025-12-02T15:30:00+03:00",
  "mainEntityOfPage": {
    "@type": "WebPage",
    "@id": "https://neurogen.news/guides/chatgpt-quickstart"
  }
}
```

```json
// HowTo (для гайдов)
{
  "@context": "https://schema.org",
  "@type": "HowTo",
  "name": "Как создать изображение в Midjourney",
  "description": "Пошаговая инструкция создания картинки с помощью Midjourney",
  "totalTime": "PT10M",
  "estimatedCost": {
    "@type": "MonetaryAmount",
    "currency": "USD",
    "value": "10"
  },
  "step": [
    {
      "@type": "HowToStep",
      "name": "Зайти в Discord",
      "text": "Откройте Discord и найдите сервер Midjourney",
      "image": "https://neurogen.news/images/step1.jpg"
    },
    {
      "@type": "HowToStep",
      "name": "Написать промпт",
      "text": "Введите /imagine и ваш промпт",
      "image": "https://neurogen.news/images/step2.jpg"
    }
  ]
}
```

```json
// FAQPage (для вопросов)
{
  "@context": "https://schema.org",
  "@type": "FAQPage",
  "mainEntity": [
    {
      "@type": "Question",
      "name": "Что такое ChatGPT?",
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "ChatGPT — это ИИ-ассистент от OpenAI..."
      }
    }
  ]
}
```

```json
// SoftwareApplication (для каталога инструментов)
{
  "@context": "https://schema.org",
  "@type": "SoftwareApplication",
  "name": "ChatGPT",
  "applicationCategory": "AI Assistant",
  "operatingSystem": "Web, iOS, Android",
  "offers": {
    "@type": "Offer",
    "price": "0",
    "priceCurrency": "USD"
  },
  "aggregateRating": {
    "@type": "AggregateRating",
    "ratingValue": "4.8",
    "ratingCount": "12345"
  }
}
```

### 2.2 Breadcrumbs (хлебные крошки)

```json
{
  "@context": "https://schema.org",
  "@type": "BreadcrumbList",
  "itemListElement": [
    {
      "@type": "ListItem",
      "position": 1,
      "name": "Главная",
      "item": "https://neurogen.news"
    },
    {
      "@type": "ListItem",
      "position": 2,
      "name": "Гайды",
      "item": "https://neurogen.news/guides"
    },
    {
      "@type": "ListItem",
      "position": 3,
      "name": "ChatGPT за 5 минут",
      "item": "https://neurogen.news/guides/chatgpt-quickstart"
    }
  ]
}
```

---

## 3. Семантическое ядро и контент-стратегия

### 3.1 Кластеры ключевых слов

```
┌─────────────────────────────────────────────────────────────────┐
│  КЛАСТЕР: "ChatGPT" (Высокий приоритет)                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Информационные (для новичков):                                 │
│  ├── "что такое chatgpt" — 50K запросов/мес                    │
│  ├── "как пользоваться chatgpt" — 30K запросов/мес             │
│  ├── "chatgpt бесплатно" — 25K запросов/мес                    │
│  ├── "chatgpt на русском" — 20K запросов/мес                   │
│  └── "chatgpt регистрация" — 15K запросов/мес                  │
│                                                                 │
│  Практические (для активных):                                   │
│  ├── "промпты для chatgpt" — 10K запросов/мес                  │
│  ├── "chatgpt для работы" — 8K запросов/мес                    │
│  ├── "chatgpt api" — 5K запросов/мес                           │
│  └── "chatgpt vs claude" — 3K запросов/мес                     │
│                                                                 │
│  Бизнес (для профессионалов):                                   │
│  ├── "chatgpt для бизнеса" — 2K запросов/мес                   │
│  ├── "chatgpt enterprise" — 1K запросов/мес                    │
│  └── "интеграция chatgpt" — 500 запросов/мес                   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 3.2 Типы контента под SEO

| Тип контента | Цель SEO | Примеры запросов |
|--------------|----------|------------------|
| **Гайды для новичков** | Высокочастотные информационные | "как начать", "что такое" |
| **Сравнения** | Среднечастотные коммерческие | "X vs Y", "лучший" |
| **Обзоры инструментов** | Среднечастотные брендовые | "обзор Midjourney" |
| **Подборки** | Высокочастотные списковые | "топ нейросетей 2025" |
| **Новости** | Быстрый трафик, Google Discover | Текущие события |
| **Глоссарий** | Низкочастотные, longtail | "что значит LLM" |

### 3.3 Контент-план (первые 3 месяца)

**Месяц 1 — Базовые гайды:**
- 10 статей "Как начать с [сервис]" (ChatGPT, Midjourney, Claude, etc.)
- 5 статей "Что такое [термин]" (нейросеть, промпт, LLM, etc.)
- 3 топ-подборки ("Лучшие бесплатные нейросети 2025")

**Месяц 2 — Расширение:**
- 10 сравнительных статей ("ChatGPT vs Claude")
- 10 практических гайдов ("Промпты для копирайтинга")
- 5 обзоров инструментов

**Месяц 3 — Углубление:**
- 10 продвинутых гайдов (API, автоматизация)
- 5 кейсов применения
- 10 ответов на частые вопросы (FAQ контент)

---

## 4. URL-структура

### 4.1 Иерархия URL

```
neurogen.news/
├── /                                  # Главная
├── /popular                           # Популярное
├── /new                               # Свежее
├── /my                                # Моя лента
│
├── /chatbots/                         # Подсайт: Чат-боты
│   ├── /chatbots/chatgpt-quickstart  # Статья
│   └── /chatbots/claude-vs-gpt       # Статья
│
├── /images/                           # Подсайт: Изображения
│   └── /images/midjourney-guide      # Статья
│
├── /guides/                           # Все гайды
├── /reviews/                          # Все обзоры
├── /prompts/                          # Библиотека промптов
├── /tools/                            # Каталог инструментов
│   └── /tools/chatgpt                # Карточка инструмента
│
├── /questions/                        # Q&A
│   └── /questions/123-slug           # Вопрос
│
├── /@username/                        # Профиль пользователя
│   ├── /@username/posts              # Посты пользователя
│   └── /@username/comments           # Комментарии
│
├── /search?q=                         # Поиск
├── /tag/prompt-engineering           # Страница тега
│
└── /about, /rules, /advertising      # Статические страницы
```

### 4.2 Правила формирования URL

```go
// Пример генерации SEO-friendly slug
func GenerateSlug(title string) string {
    // Транслитерация кириллицы
    slug := transliterate(title)
    
    // Замена пробелов на дефисы
    slug = strings.ReplaceAll(slug, " ", "-")
    
    // Удаление спецсимволов
    slug = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(strings.ToLower(slug), "")
    
    // Удаление множественных дефисов
    slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")
    
    // Ограничение длины (60 символов оптимально)
    if len(slug) > 60 {
        slug = slug[:60]
    }
    
    return strings.Trim(slug, "-")
}
```

---

## 5. Мета-теги и заголовки

### 5.1 Шаблоны title и description

```html
<!-- Главная страница -->
<title>Neurogen.News — Нейросети на практике: гайды, обзоры, промпты</title>
<meta name="description" content="Научим пользоваться нейросетями: ChatGPT, Midjourney, Claude и другие. Простые гайды для новичков и продвинутые техники для профи.">

<!-- Статья -->
<title>Как начать с ChatGPT за 5 минут — пошаговый гайд | Neurogen.News</title>
<meta name="description" content="Простая инструкция для новичков: регистрация, первый запрос, полезные команды. Начните использовать ChatGPT прямо сейчас!">

<!-- Подсайт (категория) -->
<title>Чат-боты — ChatGPT, Claude, Gemini: гайды и обзоры | Neurogen.News</title>
<meta name="description" content="Всё о текстовых нейросетях: как начать, лучшие промпты, сравнения. Практические гайды для работы и жизни.">

<!-- Каталог инструментов -->
<title>Каталог нейросетей 2025 — 100+ сервисов с обзорами | Neurogen.News</title>
<meta name="description" content="Выберите нейросеть для своих задач: фильтры по цене, категории, уровню. Честные обзоры и рейтинги от пользователей.">

<!-- Страница инструмента -->
<title>ChatGPT — обзор, цены, альтернативы | Neurogen.News</title>
<meta name="description" content="Всё о ChatGPT: возможности, тарифы (бесплатный и Plus за $20), плюсы и минусы. Рейтинг 4.8 от 12K пользователей.">
```

### 5.2 Open Graph и Twitter Cards

```html
<!-- Open Graph -->
<meta property="og:type" content="article">
<meta property="og:title" content="Как начать с ChatGPT за 5 минут">
<meta property="og:description" content="Простая инструкция для новичков...">
<meta property="og:image" content="https://neurogen.news/og/chatgpt-guide.jpg">
<meta property="og:image:width" content="1200">
<meta property="og:image:height" content="630">
<meta property="og:url" content="https://neurogen.news/guides/chatgpt-quickstart">
<meta property="og:site_name" content="Neurogen.News">

<!-- Twitter Cards -->
<meta name="twitter:card" content="summary_large_image">
<meta name="twitter:title" content="Как начать с ChatGPT за 5 минут">
<meta name="twitter:description" content="Простая инструкция для новичков...">
<meta name="twitter:image" content="https://neurogen.news/og/chatgpt-guide.jpg">
```

---

## 6. Внутренняя перелинковка

### 6.1 Стратегия перелинковки

```
┌─────────────────────────────────────────────────────────────────┐
│  СТРУКТУРА ПЕРЕЛИНКОВКИ                                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Пилларные страницы (hub pages):                               │
│  ├── /chatbots/ — ссылается на все статьи о чат-ботах         │
│  ├── /images/ — ссылается на все статьи о генерации           │
│  └── /guides/ — ссылается на все гайды                         │
│                                                                 │
│  Кластерные статьи:                                             │
│  ├── Ссылаются на пилларную страницу (1-2 ссылки)              │
│  ├── Ссылаются на связанные статьи (3-5 ссылок)                │
│  └── Ссылаются на инструменты из каталога                      │
│                                                                 │
│  Автоматическая перелинковка:                                   │
│  ├── "Читайте также" — 3 похожие статьи                        │
│  ├── "Следующий шаг" — рекомендация после прочтения            │
│  └── Упоминания инструментов → карточки в каталоге             │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 6.2 Автоматические блоки

```vue
<!-- Компонент "Читайте также" -->
<template>
  <section class="related-articles">
    <h3>Читайте также</h3>
    <div class="articles-grid">
      <ArticleCard 
        v-for="article in relatedArticles" 
        :key="article.id"
        :article="article"
        compact
      />
    </div>
  </section>
</template>

<script setup>
// Алгоритм подбора:
// 1. Та же категория
// 2. Общие теги
// 3. Похожий уровень сложности
// 4. Не читал ранее (для авторизованных)
const relatedArticles = await getRelatedArticles(currentArticle.id, 3)
</script>
```

---

## 7. Sitemap и индексация

### 7.1 Динамическая генерация Sitemap

```go
// Генератор sitemap на Go
func GenerateSitemap(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/xml")
    
    var urls []SitemapURL
    
    // Статические страницы
    urls = append(urls, SitemapURL{
        Loc:        "https://neurogen.news/",
        ChangeFreq: "hourly",
        Priority:   1.0,
    })
    
    // Категории
    categories := db.GetAllCategories()
    for _, cat := range categories {
        urls = append(urls, SitemapURL{
            Loc:        fmt.Sprintf("https://neurogen.news/%s/", cat.Slug),
            ChangeFreq: "hourly",
            Priority:   0.9,
        })
    }
    
    // Статьи (только опубликованные)
    articles := db.GetPublishedArticles()
    for _, article := range articles {
        urls = append(urls, SitemapURL{
            Loc:        fmt.Sprintf("https://neurogen.news/%s/%s", article.Category.Slug, article.Slug),
            LastMod:    article.UpdatedAt,
            ChangeFreq: "weekly",
            Priority:   0.8,
        })
    }
    
    // Инструменты
    tools := db.GetAllTools()
    for _, tool := range tools {
        urls = append(urls, SitemapURL{
            Loc:        fmt.Sprintf("https://neurogen.news/tools/%s", tool.Slug),
            ChangeFreq: "monthly",
            Priority:   0.7,
        })
    }
    
    renderSitemap(w, urls)
}
```

### 7.2 Robots.txt

```
User-agent: *
Allow: /

# Исключаем служебные страницы
Disallow: /api/
Disallow: /admin/
Disallow: /settings/
Disallow: /drafts/
Disallow: /editor/
Disallow: /_nuxt/
Disallow: /search?

# Sitemap
Sitemap: https://neurogen.news/sitemap.xml
Sitemap: https://neurogen.news/sitemap-articles.xml
Sitemap: https://neurogen.news/sitemap-tools.xml
Sitemap: https://neurogen.news/sitemap-users.xml

# Crawl-delay (опционально)
Crawl-delay: 1
```

---

## 8. Ссылочная стратегия

### 8.1 Источники ссылок

| Источник | Приоритет | Как получить |
|----------|-----------|--------------|
| **Гостевые статьи** | Высокий | VC.ru, Habr, DTF |
| **Упоминания в медиа** | Высокий | Пресс-релизы, экспертные комментарии |
| **Агрегаторы** | Средний | Яндекс.Дзен, Telegram-каналы |
| **Форумы и сообщества** | Средний | Полезные ответы со ссылкой |
| **Социальные сети** | Низкий (для сигналов) | Шаринг контента |

### 8.2 Линкбейт-контент

Контент, который естественно привлекает ссылки:

1. **Исследования** — "Мы проанализировали 1000 промптов и вот что узнали"
2. **Инструменты** — Бесплатные генераторы промптов, калькуляторы
3. **Инфографика** — Визуальные сравнения нейросетей
4. **Топ-подборки** — "100 лучших промптов для ChatGPT"
5. **Калькуляторы** — "Сколько стоит использовать GPT-4 API"

---

## 9. Мониторинг и аналитика

### 9.1 Инструменты

| Инструмент | Назначение |
|------------|------------|
| **Google Search Console** | Индексация, позиции, ошибки |
| **Яндекс.Вебмастер** | Индексация в Яндексе |
| **Google Analytics 4** | Трафик, поведение |
| **Яндекс.Метрика** | Вебвизор, карта кликов |
| **Ahrefs / Semrush** | Позиции, конкуренты, ссылки |

### 9.2 KPI для отслеживания

| Метрика | Цель (6 мес) | Цель (12 мес) |
|---------|--------------|---------------|
| Органический трафик | 50K/мес | 200K/мес |
| Проиндексировано страниц | 500 | 2000 |
| Средняя позиция | Топ-30 | Топ-10 |
| CTR из поиска | 3% | 5% |
| Ссылающиеся домены | 100 | 500 |
| Core Web Vitals | Все зелёные | Все зелёные |

---

## 10. Локальное SEO и региональность

### 10.1 Настройка для РФ

```html
<!-- Региональные мета-теги -->
<meta name="geo.region" content="RU">
<meta name="geo.placename" content="Russia">
<meta property="og:locale" content="ru_RU">

<!-- Yandex-специфичные -->
<meta name="yandex-verification" content="XXXXXXXX">
```

### 10.2 Яндекс-оптимизация

- Регистрация в Яндекс.Вебмастере
- Настройка региона (Россия)
- Турбо-страницы для мобильных
- ИКС (Индекс качества сайта)
- Яндекс.Дзен — дистрибуция контента

---

## 11. Техническая реализация SEO в коде

### 11.1 Vue-компонент для SEO

```vue
<!-- composables/useSeo.ts -->
<script setup lang="ts">
interface SeoProps {
  title: string
  description: string
  image?: string
  type?: 'article' | 'website'
  publishedAt?: string
  author?: string
  tags?: string[]
}

export function useSeo(props: SeoProps) {
  const route = useRoute()
  const config = useRuntimeConfig()
  
  const fullTitle = `${props.title} | Neurogen.News`
  const canonicalUrl = `${config.public.siteUrl}${route.path}`
  const ogImage = props.image || `${config.public.siteUrl}/og-default.jpg`
  
  useHead({
    title: fullTitle,
    meta: [
      { name: 'description', content: props.description },
      { name: 'robots', content: 'index, follow' },
      
      // Open Graph
      { property: 'og:title', content: fullTitle },
      { property: 'og:description', content: props.description },
      { property: 'og:image', content: ogImage },
      { property: 'og:url', content: canonicalUrl },
      { property: 'og:type', content: props.type || 'website' },
      { property: 'og:locale', content: 'ru_RU' },
      
      // Twitter
      { name: 'twitter:card', content: 'summary_large_image' },
      { name: 'twitter:title', content: fullTitle },
      { name: 'twitter:description', content: props.description },
      { name: 'twitter:image', content: ogImage },
      
      // Article-specific
      ...(props.publishedAt ? [
        { property: 'article:published_time', content: props.publishedAt }
      ] : []),
      ...(props.author ? [
        { property: 'article:author', content: props.author }
      ] : []),
      ...(props.tags?.map(tag => ({
        property: 'article:tag', content: tag
      })) || []),
    ],
    link: [
      { rel: 'canonical', href: canonicalUrl }
    ]
  })
  
  // JSON-LD
  useSchemaOrg([
    defineWebPage({
      name: fullTitle,
      description: props.description,
    }),
    ...(props.type === 'article' ? [
      defineArticle({
        headline: props.title,
        description: props.description,
        image: ogImage,
        datePublished: props.publishedAt,
        author: props.author,
      })
    ] : [])
  ])
}
</script>
```

### 11.2 Серверный SEO-контроллер на Go

```go
// internal/seo/handler.go
package seo

type SEOData struct {
    Title       string   `json:"title"`
    Description string   `json:"description"`
    Canonical   string   `json:"canonical"`
    OGImage     string   `json:"og_image"`
    Type        string   `json:"type"`
    Schema      string   `json:"schema"` // JSON-LD
    PublishedAt string   `json:"published_at,omitempty"`
    Author      string   `json:"author,omitempty"`
    Tags        []string `json:"tags,omitempty"`
}

func (h *Handler) GetArticleSEO(ctx context.Context, slug string) (*SEOData, error) {
    article, err := h.articleRepo.GetBySlug(ctx, slug)
    if err != nil {
        return nil, err
    }
    
    schema := generateArticleSchema(article)
    
    return &SEOData{
        Title:       article.Title,
        Description: truncate(article.Lead, 160),
        Canonical:   fmt.Sprintf("https://neurogen.news/%s/%s", article.Category.Slug, article.Slug),
        OGImage:     article.CoverImage.URL,
        Type:        "article",
        Schema:      schema,
        PublishedAt: article.PublishedAt.Format(time.RFC3339),
        Author:      article.Author.DisplayName,
        Tags:        article.TagNames(),
    }, nil
}

func generateArticleSchema(article *Article) string {
    schema := map[string]interface{}{
        "@context": "https://schema.org",
        "@type":    "Article",
        "headline": article.Title,
        "description": article.Lead,
        "image": article.CoverImage.URL,
        "datePublished": article.PublishedAt.Format(time.RFC3339),
        "dateModified": article.UpdatedAt.Format(time.RFC3339),
        "author": map[string]interface{}{
            "@type": "Person",
            "name":  article.Author.DisplayName,
            "url":   fmt.Sprintf("https://neurogen.news/@%s", article.Author.Username),
        },
        "publisher": map[string]interface{}{
            "@type": "Organization",
            "name":  "Neurogen.News",
            "logo": map[string]interface{}{
                "@type": "ImageObject",
                "url":   "https://neurogen.news/logo.png",
            },
        },
    }
    
    jsonBytes, _ := json.Marshal(schema)
    return string(jsonBytes)
}
```

---

## 12. Чек-лист запуска

### До запуска
- [ ] Настроен HTTPS с редиректом
- [ ] Robots.txt корректен
- [ ] Sitemap.xml генерируется
- [ ] Google Search Console подключен
- [ ] Яндекс.Вебмастер подключен
- [ ] Analytics/Метрика установлены
- [ ] Core Web Vitals в зелёной зоне
- [ ] Structured Data проверен через Google Rich Results Test
- [ ] Open Graph проверен через Facebook Debugger

### После запуска (первая неделя)
- [ ] Проверить индексацию главной
- [ ] Отправить sitemap в поисковики
- [ ] Проверить ошибки в Search Console
- [ ] Опубликовать 10+ статей
- [ ] Настроить автоматические отчёты

### Регулярно (еженедельно)
- [ ] Проверять ошибки индексации
- [ ] Анализировать позиции по ключам
- [ ] Отслеживать Core Web Vitals
- [ ] Публиковать новый контент
- [ ] Обновлять устаревший контент


