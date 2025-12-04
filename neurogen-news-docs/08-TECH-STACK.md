# Технологический стек Neurogen.News

## Обзор архитектуры

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              КЛИЕНТЫ                                     │
│                    (Браузеры, Мобильные приложения)                      │
└─────────────────────────────────────┬───────────────────────────────────┘
                                      │
┌─────────────────────────────────────▼───────────────────────────────────┐
│                            CDN (Cloudflare)                              │
│                    Статика, Edge Caching, DDoS Protection                │
└─────────────────────────────────────┬───────────────────────────────────┘
                                      │
┌─────────────────────────────────────▼───────────────────────────────────┐
│                         ЕДИНЫЙ БИНАРНИК GO                               │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                     Vue 3 SPA (embed)                            │   │
│  │              Статика встроена через go:embed                     │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                      Go HTTP Server                              │   │
│  │  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐             │   │
│  │  │  API Routes  │ │  SSR Engine  │ │  WebSocket   │             │   │
│  │  │  (REST/JSON) │ │  (Prerender) │ │  (Real-time) │             │   │
│  │  └──────────────┘ └──────────────┘ └──────────────┘             │   │
│  └─────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────┬───────────────────────────────────┘
                                      │
         ┌────────────────────────────┼────────────────────────────┐
         │                            │                            │
┌────────▼────────┐         ┌─────────▼─────────┐       ┌──────────▼──────────┐
│   PostgreSQL    │         │      Redis        │       │    Meilisearch      │
│   (Primary DB)  │         │  (Cache/Pub-Sub)  │       │  (Full-text Search) │
└─────────────────┘         └───────────────────┘       └─────────────────────┘
         │
         │ Репликация
         ▼
┌─────────────────┐
│ PostgreSQL      │
│ (Read Replica)  │
└─────────────────┘
```

---

## 1. Frontend: Vue 3

### 1.1 Основные технологии

| Технология | Версия | Назначение |
|------------|--------|------------|
| **Vue.js** | 3.5+ | UI фреймворк |
| **TypeScript** | 5.6+ | Типизация |
| **Vite** | 6.0+ | Сборщик |
| **Vue Router** | 4.4+ | Маршрутизация |
| **Pinia** | 2.2+ | State Management |
| **VueUse** | 11.0+ | Composition Utilities |

### 1.2 UI и стилизация

| Технология | Версия | Назначение |
|------------|--------|------------|
| **Tailwind CSS** | 3.4+ | Utility-first CSS |
| **Headless UI** | 1.7+ | Accessible компоненты |
| **Radix Vue** | 1.9+ | Primitive компоненты |
| **Lucide Icons** | 0.460+ | Иконки |
| **VueUse Motion** | 2.2+ | Анимации |

### 1.3 Редактор контента

| Технология | Версия | Назначение |
|------------|--------|------------|
| **TipTap** | 2.9+ | Rich Text Editor |
| **Shiki** | 1.22+ | Подсветка синтаксиса |
| **Marked** | 15.0+ | Markdown парсинг |

### 1.4 Структура frontend

```
frontend/
├── src/
│   ├── assets/                    # Статические ресурсы
│   │   ├── styles/
│   │   │   ├── main.css          # Tailwind + глобальные стили
│   │   │   └── variables.css     # CSS переменные (темы)
│   │   └── fonts/
│   │
│   ├── components/                # Vue компоненты
│   │   ├── ui/                   # Базовые UI компоненты
│   │   │   ├── Button.vue
│   │   │   ├── Input.vue
│   │   │   ├── Modal.vue
│   │   │   ├── Dropdown.vue
│   │   │   ├── Avatar.vue
│   │   │   ├── Badge.vue
│   │   │   ├── Skeleton.vue
│   │   │   └── Toast.vue
│   │   │
│   │   ├── layout/               # Лейаут компоненты
│   │   │   ├── AppHeader.vue
│   │   │   ├── AppSidebar.vue
│   │   │   ├── AppFooter.vue
│   │   │   └── UserMenu.vue
│   │   │
│   │   ├── feed/                 # Компоненты ленты
│   │   │   ├── ArticleCard.vue
│   │   │   ├── ArticleCardCompact.vue
│   │   │   ├── FeedList.vue
│   │   │   ├── FeedFilters.vue
│   │   │   └── QuickEditor.vue
│   │   │
│   │   ├── article/              # Компоненты статьи
│   │   │   ├── ArticleContent.vue
│   │   │   ├── ArticleHeader.vue
│   │   │   ├── ReactionPanel.vue
│   │   │   ├── ActionBar.vue
│   │   │   ├── RelatedArticles.vue
│   │   │   └── ShareButton.vue
│   │   │
│   │   ├── comments/             # Комментарии
│   │   │   ├── CommentList.vue
│   │   │   ├── CommentItem.vue
│   │   │   ├── CommentForm.vue
│   │   │   └── CommentThread.vue
│   │   │
│   │   ├── editor/               # Редактор
│   │   │   ├── RichTextEditor.vue
│   │   │   ├── EditorToolbar.vue
│   │   │   ├── BlockPicker.vue
│   │   │   ├── ImageUploader.vue
│   │   │   ├── PromptBlock.vue
│   │   │   └── ChecklistBlock.vue
│   │   │
│   │   ├── user/                 # Пользователь
│   │   │   ├── UserProfile.vue
│   │   │   ├── UserCard.vue
│   │   │   ├── AchievementBadge.vue
│   │   │   └── SettingsForm.vue
│   │   │
│   │   ├── tools/                # Каталог инструментов
│   │   │   ├── ToolCard.vue
│   │   │   ├── ToolGrid.vue
│   │   │   ├── ToolFilters.vue
│   │   │   └── ToolComparison.vue
│   │   │
│   │   ├── seo/                  # SEO компоненты
│   │   │   ├── SeoHead.vue
│   │   │   ├── SchemaOrg.vue
│   │   │   └── Breadcrumbs.vue
│   │   │
│   │   └── common/               # Общие компоненты
│   │       ├── CopyableBlock.vue
│   │       ├── LevelBadge.vue
│   │       ├── ReadingTime.vue
│   │       └── InfiniteScroll.vue
│   │
│   ├── composables/              # Vue Composables
│   │   ├── useAuth.ts
│   │   ├── useApi.ts
│   │   ├── useSeo.ts
│   │   ├── useNotifications.ts
│   │   ├── useDrafts.ts
│   │   ├── useBookmarks.ts
│   │   ├── useWebSocket.ts
│   │   ├── useInfiniteScroll.ts
│   │   └── useTheme.ts
│   │
│   ├── stores/                   # Pinia Stores
│   │   ├── auth.ts
│   │   ├── user.ts
│   │   ├── feed.ts
│   │   ├── article.ts
│   │   ├── notifications.ts
│   │   ├── drafts.ts
│   │   └── ui.ts
│   │
│   ├── pages/                    # Страницы (Vue Router)
│   │   ├── index.vue             # Главная (популярное)
│   │   ├── new.vue               # Свежее
│   │   ├── my.vue                # Моя лента
│   │   ├── start.vue             # С чего начать
│   │   ├── search.vue            # Поиск
│   │   ├── rating.vue            # Рейтинг
│   │   │
│   │   ├── [category]/
│   │   │   ├── index.vue         # Страница категории
│   │   │   └── [slug].vue        # Страница статьи
│   │   │
│   │   ├── tools/
│   │   │   ├── index.vue         # Каталог инструментов
│   │   │   └── [slug].vue        # Карточка инструмента
│   │   │
│   │   ├── questions/
│   │   │   ├── index.vue         # Q&A лента
│   │   │   └── [id].vue          # Страница вопроса
│   │   │
│   │   ├── user/
│   │   │   ├── [username].vue    # Профиль пользователя
│   │   │   ├── drafts.vue        # Черновики
│   │   │   ├── bookmarks.vue     # Закладки
│   │   │   ├── achievements.vue  # Ачивки
│   │   │   └── settings/
│   │   │       ├── blog.vue
│   │   │       ├── feeds.vue
│   │   │       ├── general.vue
│   │   │       └── notifications.vue
│   │   │
│   │   ├── editor/
│   │   │   ├── index.vue         # Новая статья
│   │   │   └── [draftId].vue     # Редактирование черновика
│   │   │
│   │   ├── auth/
│   │   │   ├── login.vue
│   │   │   └── register.vue
│   │   │
│   │   └── static/
│   │       ├── about.vue
│   │       ├── rules.vue
│   │       └── advertising.vue
│   │
│   ├── router/
│   │   └── index.ts              # Vue Router конфигурация
│   │
│   ├── api/                      # API клиент
│   │   ├── client.ts             # Axios/Fetch wrapper
│   │   ├── articles.ts
│   │   ├── comments.ts
│   │   ├── users.ts
│   │   ├── drafts.ts
│   │   ├── tools.ts
│   │   └── types.ts              # API типы
│   │
│   ├── utils/                    # Утилиты
│   │   ├── formatters.ts         # Форматирование дат, чисел
│   │   ├── validators.ts         # Валидация форм
│   │   ├── slug.ts               # Генерация slug
│   │   └── seo.ts                # SEO хелперы
│   │
│   ├── types/                    # TypeScript типы
│   │   ├── article.ts
│   │   ├── user.ts
│   │   ├── comment.ts
│   │   └── index.ts
│   │
│   ├── App.vue
│   └── main.ts
│
├── public/
│   ├── favicon.ico
│   ├── robots.txt
│   └── og-default.jpg
│
├── index.html
├── vite.config.ts
├── tailwind.config.ts
├── tsconfig.json
└── package.json
```

### 1.5 Пример компонента

```vue
<!-- src/components/feed/ArticleCard.vue -->
<script setup lang="ts">
import { computed } from 'vue'
import type { Article } from '@/types'
import { formatRelativeTime, formatNumber } from '@/utils/formatters'

const props = defineProps<{
  article: Article
  compact?: boolean
}>()

const levelBadge = computed(() => {
  const levels = {
    beginner: { emoji: '🟢', label: 'Для новичков' },
    intermediate: { emoji: '🟡', label: 'Продвинутое' },
    advanced: { emoji: '🔴', label: 'Для бизнеса' },
  }
  return levels[props.article.level]
})

const contentTypeBadge = computed(() => {
  const types = {
    guide: '📖 Гайд',
    review: '🔍 Обзор',
    prompts: '💡 Промпты',
    news: '📰 Новость',
    question: '❓ Вопрос',
  }
  return types[props.article.contentType]
})
</script>

<template>
  <article 
    class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 overflow-hidden hover:shadow-lg transition-shadow"
    :class="{ 'p-4': compact, 'p-6': !compact }"
  >
    <!-- Header -->
    <header class="flex items-center gap-3 mb-4">
      <RouterLink :to="`/@${article.author.username}`">
        <Avatar :src="article.author.avatarUrl" :alt="article.author.displayName" size="sm" />
      </RouterLink>
      
      <div class="flex-1 min-w-0">
        <div class="flex items-center gap-2 text-sm">
          <RouterLink 
            :to="`/@${article.author.username}`"
            class="font-medium hover:text-blue-600 truncate"
          >
            {{ article.author.displayName }}
          </RouterLink>
          <span class="text-gray-400">•</span>
          <RouterLink 
            :to="`/${article.category.slug}`"
            class="text-gray-500 hover:text-blue-600"
          >
            {{ article.category.name }}
          </RouterLink>
          <span class="text-gray-400">•</span>
          <time class="text-gray-400" :datetime="article.publishedAt">
            {{ formatRelativeTime(article.publishedAt) }}
          </time>
        </div>
      </div>
      
      <Button variant="ghost" size="sm">
        Подписаться
      </Button>
    </header>
    
    <!-- Meta badges -->
    <div class="flex items-center gap-2 mb-3 text-sm">
      <Badge :variant="article.level">
        {{ levelBadge.emoji }} {{ levelBadge.label }}
      </Badge>
      <Badge variant="secondary">
        {{ contentTypeBadge }}
      </Badge>
      <span class="text-gray-400">
        {{ article.readingTime }} мин
      </span>
    </div>
    
    <!-- Content -->
    <RouterLink :to="`/${article.category.slug}/${article.slug}`" class="block group">
      <h2 class="text-xl font-bold mb-2 group-hover:text-blue-600 transition-colors">
        {{ article.title }}
        <Badge v-if="article.isEditorial" variant="primary" class="ml-2">
          ✓ Материал редакции
        </Badge>
      </h2>
      
      <p v-if="!compact" class="text-gray-600 dark:text-gray-400 mb-4 line-clamp-3">
        {{ article.lead }}
      </p>
      
      <img 
        v-if="article.coverImage && !compact"
        :src="article.coverImage.url"
        :alt="article.title"
        class="w-full aspect-video object-cover rounded-lg"
        loading="lazy"
      />
    </RouterLink>
    
    <!-- Reactions -->
    <div class="flex items-center gap-2 mt-4 pt-4 border-t border-gray-100 dark:border-gray-800">
      <ReactionPanel :reactions="article.reactions" :articleId="article.id" />
    </div>
    
    <!-- Actions -->
    <div class="flex items-center justify-between mt-3 text-sm text-gray-500">
      <div class="flex items-center gap-4">
        <RouterLink 
          :to="`/${article.category.slug}/${article.slug}#comments`"
          class="flex items-center gap-1 hover:text-blue-600"
        >
          💬 {{ formatNumber(article.commentCount) }}
        </RouterLink>
        
        <button class="flex items-center gap-1 hover:text-blue-600">
          🔖 {{ formatNumber(article.bookmarkCount) }}
        </button>
        
        <ShareButton :url="`/${article.category.slug}/${article.slug}`" :title="article.title" />
      </div>
      
      <span class="text-gray-400">
        👁 {{ formatNumber(article.viewCount) }}
      </span>
    </div>
  </article>
</template>
```

### 1.6 Vite конфигурация

```typescript
// vite.config.ts
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import { compression } from 'vite-plugin-compression2'

export default defineConfig({
  plugins: [
    vue(),
    // Brotli и Gzip сжатие для production
    compression({
      algorithm: 'brotliCompress',
      exclude: [/\.(br)$/, /\.(gz)$/],
    }),
    compression({
      algorithm: 'gzip',
      exclude: [/\.(br)$/, /\.(gz)$/],
    }),
  ],
  
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  
  build: {
    target: 'esnext',
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,
        drop_debugger: true,
      },
    },
    rollupOptions: {
      output: {
        manualChunks: {
          // Разделение vendor chunks для кэширования
          'vue-vendor': ['vue', 'vue-router', 'pinia'],
          'editor': ['@tiptap/vue-3', '@tiptap/starter-kit'],
          'ui': ['@headlessui/vue', '@vueuse/core'],
        },
      },
    },
    // Для встраивания в Go binary
    assetsDir: 'assets',
    cssCodeSplit: true,
    sourcemap: false,
  },
  
  // SSR prerender для SEO-критичных страниц
  // Используем vite-ssg или prerender-spa-plugin
})
```

---

## 2. Backend: Go

### 2.1 Основные технологии

| Технология | Версия | Назначение |
|------------|--------|------------|
| **Go** | 1.23+ | Язык |
| **Fiber** | v2.52+ / v3 | HTTP фреймворк |
| **sqlc** | 1.27+ | Type-safe SQL |
| **pgx** | v5.7+ | PostgreSQL драйвер |
| **go-redis** | v9.7+ | Redis клиент |
| **validator** | v10 | Валидация |
| **zap** | 1.27+ | Логирование |
| **viper** | 1.19+ | Конфигурация |

### 2.2 Структура backend

```
backend/
├── cmd/
│   └── server/
│       └── main.go                # Точка входа
│
├── internal/
│   ├── config/
│   │   └── config.go             # Конфигурация приложения
│   │
│   ├── server/
│   │   ├── server.go             # HTTP сервер
│   │   ├── routes.go             # Маршрутизация
│   │   └── middleware.go         # Middleware
│   │
│   ├── handler/                  # HTTP handlers
│   │   ├── article.go
│   │   ├── comment.go
│   │   ├── user.go
│   │   ├── auth.go
│   │   ├── draft.go
│   │   ├── bookmark.go
│   │   ├── notification.go
│   │   ├── tool.go
│   │   ├── search.go
│   │   ├── feed.go
│   │   └── seo.go                # SEO данные (sitemap, etc.)
│   │
│   ├── service/                  # Бизнес-логика
│   │   ├── article.go
│   │   ├── comment.go
│   │   ├── user.go
│   │   ├── auth.go
│   │   ├── draft.go
│   │   ├── bookmark.go
│   │   ├── notification.go
│   │   ├── achievement.go
│   │   ├── search.go
│   │   ├── feed.go
│   │   └── seo.go
│   │
│   ├── repository/               # Работа с БД (sqlc generated)
│   │   ├── queries/              # SQL запросы
│   │   │   ├── articles.sql
│   │   │   ├── comments.sql
│   │   │   ├── users.sql
│   │   │   └── ...
│   │   ├── db.go                 # Сгенерированный код
│   │   ├── articles.sql.go
│   │   ├── comments.sql.go
│   │   └── ...
│   │
│   ├── cache/                    # Кэширование (Redis)
│   │   ├── cache.go
│   │   ├── article.go
│   │   ├── feed.go
│   │   └── user.go
│   │
│   ├── search/                   # Полнотекстовый поиск
│   │   ├── meilisearch.go
│   │   └── indexer.go
│   │
│   ├── websocket/                # Real-time
│   │   ├── hub.go
│   │   ├── client.go
│   │   └── events.go
│   │
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── ratelimit.go
│   │   ├── cors.go
│   │   ├── compress.go
│   │   ├── cache.go
│   │   └── seo.go
│   │
│   ├── model/                    # Доменные модели
│   │   ├── article.go
│   │   ├── comment.go
│   │   ├── user.go
│   │   ├── reaction.go
│   │   ├── notification.go
│   │   └── achievement.go
│   │
│   ├── dto/                      # Data Transfer Objects
│   │   ├── request/
│   │   │   ├── article.go
│   │   │   ├── comment.go
│   │   │   └── user.go
│   │   └── response/
│   │       ├── article.go
│   │       ├── comment.go
│   │       └── user.go
│   │
│   └── pkg/                      # Внутренние пакеты
│       ├── validator/
│       ├── slug/
│       ├── jwt/
│       ├── hash/
│       └── errors/
│
├── pkg/                          # Публичные пакеты
│   ├── logger/
│   └── utils/
│
├── migrations/                   # SQL миграции
│   ├── 001_init.up.sql
│   ├── 001_init.down.sql
│   └── ...
│
├── web/                          # Встроенные статические файлы
│   └── dist/                     # Собранный Vue SPA
│       └── ... (go:embed)
│
├── scripts/
│   ├── build.sh
│   └── migrate.sh
│
├── Makefile
├── Dockerfile
├── go.mod
└── go.sum
```

### 2.3 Встраивание Vue SPA в Go бинарник

```go
// cmd/server/main.go
package main

import (
    "embed"
    "io/fs"
    "net/http"
    
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/filesystem"
    
    "neurogen-news/internal/config"
    "neurogen-news/internal/server"
)

//go:embed web/dist/*
var embeddedFiles embed.FS

func main() {
    cfg := config.Load()
    
    app := fiber.New(fiber.Config{
        Prefork:       cfg.Prefork,
        ServerHeader:  "Neurogen",
        StrictRouting: true,
        CaseSensitive: true,
        // Оптимизации для производительности
        ReadBufferSize:  8192,
        WriteBufferSize: 8192,
        BodyLimit:       10 * 1024 * 1024, // 10MB
    })
    
    // Извлекаем статику из embedded FS
    distFS, _ := fs.Sub(embeddedFiles, "web/dist")
    
    // Middleware
    app.Use(middleware.Compress())
    app.Use(middleware.CORS(cfg.CORS))
    app.Use(middleware.RateLimit(cfg.RateLimit))
    app.Use(middleware.Logger())
    
    // API routes
    api := app.Group("/api/v1")
    server.SetupAPIRoutes(api, cfg)
    
    // SSR для SEO-критичных страниц (статьи, категории)
    app.Get("/sitemap.xml", handler.Sitemap)
    app.Get("/robots.txt", handler.Robots)
    
    // SPA fallback — отдаём Vue приложение
    app.Use("/", filesystem.New(filesystem.Config{
        Root:         http.FS(distFS),
        Index:        "index.html",
        NotFoundFile: "index.html", // SPA fallback
        MaxAge:       86400,         // 1 day cache
    }))
    
    app.Listen(cfg.ServerAddr)
}
```

### 2.4 Пример handler'а

```go
// internal/handler/article.go
package handler

import (
    "github.com/gofiber/fiber/v2"
    
    "neurogen-news/internal/dto/request"
    "neurogen-news/internal/dto/response"
    "neurogen-news/internal/service"
)

type ArticleHandler struct {
    articleService *service.ArticleService
    cacheService   *service.CacheService
    seoService     *service.SEOService
}

func NewArticleHandler(
    articleService *service.ArticleService,
    cacheService *service.CacheService,
    seoService *service.SEOService,
) *ArticleHandler {
    return &ArticleHandler{
        articleService: articleService,
        cacheService:   cacheService,
        seoService:     seoService,
    }
}

// GetArticle godoc
// @Summary      Get article by slug
// @Description  Returns article with SEO data for SSR
// @Tags         articles
// @Accept       json
// @Produce      json
// @Param        category  path     string  true  "Category slug"
// @Param        slug      path     string  true  "Article slug"
// @Success      200       {object} response.ArticleWithSEO
// @Failure      404       {object} response.Error
// @Router       /api/v1/articles/{category}/{slug} [get]
func (h *ArticleHandler) GetArticle(c *fiber.Ctx) error {
    categorySlug := c.Params("category")
    articleSlug := c.Params("slug")
    
    // Проверяем кэш
    cacheKey := "article:" + categorySlug + ":" + articleSlug
    if cached, err := h.cacheService.Get(c.Context(), cacheKey); err == nil {
        return c.JSON(cached)
    }
    
    // Получаем статью
    article, err := h.articleService.GetBySlug(c.Context(), categorySlug, articleSlug)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(response.Error{
            Message: "Article not found",
        })
    }
    
    // Инкремент просмотров (асинхронно)
    go h.articleService.IncrementViews(article.ID)
    
    // Генерируем SEO данные
    seoData := h.seoService.GenerateArticleSEO(article)
    
    // Формируем ответ
    resp := response.ArticleWithSEO{
        Article: response.ArticleFromModel(article),
        SEO:     seoData,
    }
    
    // Кэшируем на 10 минут
    h.cacheService.Set(c.Context(), cacheKey, resp, 10*time.Minute)
    
    return c.JSON(resp)
}

// GetFeed godoc
// @Summary      Get articles feed
// @Description  Returns paginated feed with filters
// @Tags         articles
// @Accept       json
// @Produce      json
// @Param        sort      query    string  false  "Sort: popular, new"  default(popular)
// @Param        category  query    string  false  "Category slug"
// @Param        level     query    string  false  "Level: beginner, intermediate, advanced"
// @Param        type      query    string  false  "Content type: guide, review, prompts, news"
// @Param        cursor    query    string  false  "Pagination cursor"
// @Param        limit     query    int     false  "Items per page"  default(20)
// @Success      200       {object} response.ArticleFeed
// @Router       /api/v1/feed [get]
func (h *ArticleHandler) GetFeed(c *fiber.Ctx) error {
    req := request.FeedRequest{
        Sort:     c.Query("sort", "popular"),
        Category: c.Query("category"),
        Level:    c.Query("level"),
        Type:     c.Query("type"),
        Cursor:   c.Query("cursor"),
        Limit:    c.QueryInt("limit", 20),
    }
    
    // Валидация
    if err := validate.Struct(req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ValidationError(err))
    }
    
    // Получаем ленту
    feed, err := h.articleService.GetFeed(c.Context(), req)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.Error{
            Message: "Failed to load feed",
        })
    }
    
    return c.JSON(feed)
}

// CreateArticle godoc
// @Summary      Create new article
// @Description  Creates article from draft
// @Tags         articles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        article  body     request.CreateArticle  true  "Article data"
// @Success      201      {object} response.Article
// @Failure      400      {object} response.Error
// @Failure      401      {object} response.Error
// @Router       /api/v1/articles [post]
func (h *ArticleHandler) CreateArticle(c *fiber.Ctx) error {
    userID := c.Locals("userID").(string)
    
    var req request.CreateArticle
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.Error{
            Message: "Invalid request body",
        })
    }
    
    // Валидация
    if err := validate.Struct(req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ValidationError(err))
    }
    
    // Создаём статью
    article, err := h.articleService.Create(c.Context(), userID, req)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.Error{
            Message: "Failed to create article",
        })
    }
    
    // Инвалидируем кэш ленты
    h.cacheService.InvalidatePattern(c.Context(), "feed:*")
    
    // Индексируем для поиска
    go h.articleService.IndexForSearch(article)
    
    return c.Status(fiber.StatusCreated).JSON(response.ArticleFromModel(article))
}
```

### 2.5 SQL с sqlc

```sql
-- internal/repository/queries/articles.sql

-- name: GetArticleBySlug :one
SELECT 
    a.id,
    a.title,
    a.slug,
    a.lead,
    a.content,
    a.cover_image,
    a.level,
    a.content_type,
    a.reading_time,
    a.is_editorial,
    a.view_count,
    a.comment_count,
    a.bookmark_count,
    a.published_at,
    a.created_at,
    a.updated_at,
    -- Author
    u.id as author_id,
    u.username as author_username,
    u.display_name as author_display_name,
    u.avatar_url as author_avatar_url,
    u.is_verified as author_is_verified,
    -- Category
    c.id as category_id,
    c.slug as category_slug,
    c.name as category_name,
    c.icon as category_icon
FROM articles a
JOIN users u ON a.author_id = u.id
JOIN categories c ON a.category_id = c.id
WHERE c.slug = $1 
  AND a.slug = $2 
  AND a.status = 'published'
  AND a.deleted_at IS NULL;

-- name: GetFeedPopular :many
SELECT 
    a.id,
    a.title,
    a.slug,
    a.lead,
    a.cover_image,
    a.level,
    a.content_type,
    a.reading_time,
    a.is_editorial,
    a.view_count,
    a.comment_count,
    a.bookmark_count,
    a.published_at,
    u.id as author_id,
    u.username as author_username,
    u.display_name as author_display_name,
    u.avatar_url as author_avatar_url,
    c.id as category_id,
    c.slug as category_slug,
    c.name as category_name
FROM articles a
JOIN users u ON a.author_id = u.id
JOIN categories c ON a.category_id = c.id
WHERE a.status = 'published'
  AND a.deleted_at IS NULL
  AND (sqlc.narg('category_slug')::text IS NULL OR c.slug = sqlc.narg('category_slug'))
  AND (sqlc.narg('level')::text IS NULL OR a.level = sqlc.narg('level'))
  AND (sqlc.narg('content_type')::text IS NULL OR a.content_type = sqlc.narg('content_type'))
  AND (sqlc.narg('cursor_score')::float IS NULL OR 
       (a.popularity_score, a.id) < (sqlc.narg('cursor_score'), sqlc.narg('cursor_id')::uuid))
ORDER BY a.popularity_score DESC, a.id DESC
LIMIT sqlc.arg('limit');

-- name: IncrementArticleViews :exec
UPDATE articles 
SET view_count = view_count + 1
WHERE id = $1;

-- name: CreateArticle :one
INSERT INTO articles (
    id,
    author_id,
    category_id,
    title,
    slug,
    lead,
    content,
    cover_image,
    level,
    content_type,
    reading_time,
    is_editorial,
    status,
    published_at
) VALUES (
    gen_random_uuid(),
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
    CASE WHEN $12 = 'published' THEN NOW() ELSE NULL END
)
RETURNING *;
```

### 2.6 Rate Limiting и кэширование

```go
// internal/middleware/ratelimit.go
package middleware

import (
    "time"
    
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/limiter"
)

func RateLimit(cfg config.RateLimitConfig) fiber.Handler {
    return limiter.New(limiter.Config{
        Max:        cfg.RequestsPerMinute,
        Expiration: time.Minute,
        KeyGenerator: func(c *fiber.Ctx) string {
            // Используем IP + User ID для авторизованных
            if userID := c.Locals("userID"); userID != nil {
                return userID.(string)
            }
            return c.IP()
        },
        LimitReached: func(c *fiber.Ctx) error {
            return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
                "error": "Too many requests, please try again later",
            })
        },
        // Используем Redis для распределённого rate limiting
        Storage: redisStorage,
    })
}
```

```go
// internal/cache/cache.go
package cache

import (
    "context"
    "encoding/json"
    "time"
    
    "github.com/redis/go-redis/v9"
)

type Cache struct {
    client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
    return &Cache{client: client}
}

func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
    val, err := c.client.Get(ctx, key).Result()
    if err != nil {
        return err
    }
    return json.Unmarshal([]byte(val), dest)
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *Cache) InvalidatePattern(ctx context.Context, pattern string) error {
    var cursor uint64
    for {
        keys, nextCursor, err := c.client.Scan(ctx, cursor, pattern, 100).Result()
        if err != nil {
            return err
        }
        
        if len(keys) > 0 {
            c.client.Del(ctx, keys...)
        }
        
        cursor = nextCursor
        if cursor == 0 {
            break
        }
    }
    return nil
}

// Pub/Sub для инвалидации кэша между инстансами
func (c *Cache) SubscribeInvalidation(ctx context.Context) {
    pubsub := c.client.Subscribe(ctx, "cache:invalidate")
    
    for msg := range pubsub.Channel() {
        c.client.Del(ctx, msg.Payload)
    }
}
```

---

## 3. База данных: PostgreSQL 17

### 3.1 Конфигурация для высокой нагрузки

```sql
-- postgresql.conf оптимизации

-- Память
shared_buffers = '4GB'              -- 25% RAM
effective_cache_size = '12GB'       -- 75% RAM
work_mem = '256MB'
maintenance_work_mem = '1GB'

-- WAL и checkpoint
wal_buffers = '64MB'
checkpoint_completion_target = 0.9
max_wal_size = '4GB'

-- Параллельные запросы
max_parallel_workers_per_gather = 4
max_parallel_workers = 8
max_worker_processes = 16

-- Connections
max_connections = 200
```

### 3.2 Индексы для производительности

```sql
-- Индексы для быстрой фильтрации ленты
CREATE INDEX CONCURRENTLY idx_articles_feed_popular 
ON articles (popularity_score DESC, id DESC) 
WHERE status = 'published' AND deleted_at IS NULL;

CREATE INDEX CONCURRENTLY idx_articles_feed_new 
ON articles (published_at DESC, id DESC) 
WHERE status = 'published' AND deleted_at IS NULL;

-- Составные индексы для фильтров
CREATE INDEX CONCURRENTLY idx_articles_category_level 
ON articles (category_id, level, popularity_score DESC) 
WHERE status = 'published' AND deleted_at IS NULL;

-- Полнотекстовый поиск (если не используем Meilisearch)
CREATE INDEX CONCURRENTLY idx_articles_search 
ON articles USING gin(to_tsvector('russian', title || ' ' || lead));

-- Индексы для комментариев
CREATE INDEX CONCURRENTLY idx_comments_article 
ON comments (article_id, created_at DESC) 
WHERE deleted_at IS NULL;

-- Индексы для уведомлений
CREATE INDEX CONCURRENTLY idx_notifications_user_unread 
ON notifications (user_id, created_at DESC) 
WHERE is_read = false;
```

---

## 4. Кэширование: Redis 7

### 4.1 Стратегия кэширования

```go
// Ключи и TTL
const (
    // Горячие данные — короткий TTL
    CacheFeedPopular    = "feed:popular:%s"      // 5 min
    CacheFeedNew        = "feed:new:%s"          // 2 min
    CacheArticle        = "article:%s:%s"        // 10 min
    
    // Профили и пользователи
    CacheUser           = "user:%s"              // 15 min
    CacheUserStats      = "user:stats:%s"        // 5 min
    
    // Агрегаты — дольше
    CacheCategories     = "categories:all"       // 1 hour
    CacheToolsCatalog   = "tools:catalog"        // 30 min
    
    // Сессии
    CacheSession        = "session:%s"           // 24 hours
)
```

### 4.2 Pub/Sub для real-time

```go
// internal/websocket/events.go
package websocket

const (
    EventNewComment     = "comment:new"
    EventNewReaction    = "reaction:new"
    EventOnlineCount    = "online:count"
    EventNotification   = "notification"
    EventAchievement    = "achievement:unlocked"
)

type Event struct {
    Type    string      `json:"type"`
    Payload interface{} `json:"payload"`
}

// Публикация события
func (h *Hub) PublishEvent(ctx context.Context, channel string, event Event) error {
    data, _ := json.Marshal(event)
    return h.redis.Publish(ctx, channel, data).Err()
}

// Подписка на события
func (h *Hub) SubscribeChannel(ctx context.Context, channel string) {
    pubsub := h.redis.Subscribe(ctx, channel)
    
    for msg := range pubsub.Channel() {
        var event Event
        json.Unmarshal([]byte(msg.Payload), &event)
        
        // Рассылаем подключенным клиентам
        h.BroadcastToChannel(channel, event)
    }
}
```

---

## 5. Полнотекстовый поиск: Meilisearch

### 5.1 Настройка индексов

```go
// internal/search/meilisearch.go
package search

import (
    "github.com/meilisearch/meilisearch-go"
)

type SearchService struct {
    client *meilisearch.Client
}

func NewSearchService(host, apiKey string) *SearchService {
    client := meilisearch.NewClient(meilisearch.ClientConfig{
        Host:   host,
        APIKey: apiKey,
    })
    
    return &SearchService{client: client}
}

func (s *SearchService) SetupIndexes() error {
    // Индекс статей
    articlesIndex := s.client.Index("articles")
    
    // Настройка searchable атрибутов (приоритет)
    articlesIndex.UpdateSearchableAttributes(&[]string{
        "title",      // Высший приоритет
        "lead",
        "content",
        "tags",
        "author_name",
    })
    
    // Фильтруемые атрибуты
    articlesIndex.UpdateFilterableAttributes(&[]string{
        "category_slug",
        "level",
        "content_type",
        "is_editorial",
        "published_at",
    })
    
    // Сортируемые атрибуты
    articlesIndex.UpdateSortableAttributes(&[]string{
        "published_at",
        "popularity_score",
        "view_count",
    })
    
    // Настройка ranking rules
    articlesIndex.UpdateRankingRules(&[]string{
        "words",
        "typo",
        "proximity",
        "attribute",
        "sort",
        "exactness",
        "published_at:desc", // Свежие выше
    })
    
    // Настройка для русского языка
    articlesIndex.UpdateTypoTolerance(&meilisearch.TypoTolerance{
        Enabled: true,
        MinWordSizeForTypos: meilisearch.MinWordSizeForTypos{
            OneTypo:  4,
            TwoTypos: 8,
        },
    })
    
    return nil
}

// Поиск
func (s *SearchService) Search(ctx context.Context, query string, filters SearchFilters) (*SearchResult, error) {
    index := s.client.Index("articles")
    
    searchParams := &meilisearch.SearchRequest{
        Query:  query,
        Limit:  filters.Limit,
        Offset: filters.Offset,
        Filter: buildFilterString(filters),
        Sort:   []string{filters.Sort},
        AttributesToHighlight: []string{"title", "lead"},
        HighlightPreTag:  "<mark>",
        HighlightPostTag: "</mark>",
    }
    
    result, err := index.Search(query, searchParams)
    if err != nil {
        return nil, err
    }
    
    return mapSearchResult(result), nil
}
```

---

## 6. Сборка и деплой

### 6.1 Makefile

```makefile
# Makefile

.PHONY: all build build-frontend build-backend run dev test clean

# Версии
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -s -w"

# Paths
FRONTEND_DIR := frontend
BACKEND_DIR := backend
DIST_DIR := $(BACKEND_DIR)/web/dist
OUTPUT := neurogen-news

all: build

# Сборка frontend
build-frontend:
	cd $(FRONTEND_DIR) && npm ci && npm run build
	rm -rf $(DIST_DIR)
	cp -r $(FRONTEND_DIR)/dist $(DIST_DIR)

# Сборка backend с embedded frontend
build-backend:
	cd $(BACKEND_DIR) && CGO_ENABLED=0 go build $(LDFLAGS) -o ../$(OUTPUT) ./cmd/server

# Полная сборка
build: build-frontend build-backend
	@echo "Built $(OUTPUT) $(VERSION)"
	@ls -lh $(OUTPUT)

# Сборка для разных платформ
build-linux:
	GOOS=linux GOARCH=amd64 $(MAKE) build-backend OUTPUT=neurogen-news-linux-amd64

build-darwin:
	GOOS=darwin GOARCH=arm64 $(MAKE) build-backend OUTPUT=neurogen-news-darwin-arm64

build-all: build-frontend build-linux build-darwin

# Разработка
dev:
	@echo "Starting dev servers..."
	cd $(FRONTEND_DIR) && npm run dev &
	cd $(BACKEND_DIR) && air

# Запуск production
run:
	./$(OUTPUT)

# Тесты
test:
	cd $(BACKEND_DIR) && go test -v -race ./...
	cd $(FRONTEND_DIR) && npm run test

# Линтеры
lint:
	cd $(BACKEND_DIR) && golangci-lint run
	cd $(FRONTEND_DIR) && npm run lint

# Миграции
migrate-up:
	cd $(BACKEND_DIR) && goose -dir migrations postgres "$(DATABASE_URL)" up

migrate-down:
	cd $(BACKEND_DIR) && goose -dir migrations postgres "$(DATABASE_URL)" down

# Генерация sqlc
sqlc:
	cd $(BACKEND_DIR) && sqlc generate

# Очистка
clean:
	rm -f $(OUTPUT) neurogen-news-*
	rm -rf $(DIST_DIR)
	cd $(FRONTEND_DIR) && rm -rf dist node_modules/.cache
```

### 6.2 Dockerfile (multi-stage)

```dockerfile
# Dockerfile

# Stage 1: Build frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: Build backend
FROM golang:1.23-alpine AS backend-builder
WORKDIR /app

# Устанавливаем зависимости
RUN apk add --no-cache git

# Копируем go.mod и go.sum
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Копируем исходники
COPY backend/ ./

# Копируем собранный frontend
COPY --from=frontend-builder /app/frontend/dist ./web/dist

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /neurogen-news \
    ./cmd/server

# Stage 3: Production image
FROM alpine:3.19

# Добавляем ca-certificates для HTTPS
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Копируем бинарник
COPY --from=backend-builder /neurogen-news .

# Создаём непривилегированного пользователя
RUN adduser -D -g '' appuser
USER appuser

# Порт
EXPOSE 8080

# Healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Запуск
CMD ["./neurogen-news"]
```

### 6.3 Docker Compose для разработки

```yaml
# docker-compose.yml
version: '3.8'

services:
  postgres:
    image: postgres:17-alpine
    environment:
      POSTGRES_DB: neurogen
      POSTGRES_USER: neurogen
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:-dev_password}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backend/migrations:/docker-entrypoint-initdb.d
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U neurogen"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  meilisearch:
    image: getmeili/meilisearch:v1.11
    environment:
      MEILI_MASTER_KEY: ${MEILI_MASTER_KEY:-dev_master_key}
      MEILI_ENV: development
    volumes:
      - meilisearch_data:/meili_data
    ports:
      - "7700:7700"

  app:
    build:
      context: .
      dockerfile: Dockerfile
    environment:
      - DATABASE_URL=postgres://neurogen:${POSTGRES_PASSWORD:-dev_password}@postgres:5432/neurogen?sslmode=disable
      - REDIS_URL=redis://redis:6379
      - MEILISEARCH_URL=http://meilisearch:7700
      - MEILISEARCH_KEY=${MEILI_MASTER_KEY:-dev_master_key}
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
      meilisearch:
        condition: service_started

volumes:
  postgres_data:
  redis_data:
  meilisearch_data:
```

---

## 7. Производительность и масштабирование

### 7.1 Метрики производительности

| Метрика | Целевое значение |
|---------|------------------|
| Время ответа API (p95) | < 100ms |
| Время ответа API (p99) | < 200ms |
| RPS на инстанс | 5,000+ |
| Память на инстанс | < 512MB |
| Горячий старт | < 100ms |
| Размер бинарника | < 50MB |

### 7.2 Оптимизации Go

```go
// Пул объектов для уменьшения аллокаций
var responsePool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

// Connection pooling для PostgreSQL
func NewDBPool(connString string) (*pgxpool.Pool, error) {
    config, err := pgxpool.ParseConfig(connString)
    if err != nil {
        return nil, err
    }
    
    config.MaxConns = 50
    config.MinConns = 10
    config.MaxConnLifetime = time.Hour
    config.MaxConnIdleTime = 30 * time.Minute
    config.HealthCheckPeriod = time.Minute
    
    return pgxpool.NewWithConfig(context.Background(), config)
}

// Graceful shutdown
func (s *Server) Shutdown(ctx context.Context) error {
    // Прекращаем принимать новые запросы
    s.app.Shutdown()
    
    // Ждём завершения текущих запросов
    shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()
    
    // Закрываем соединения
    s.db.Close()
    s.redis.Close()
    
    return nil
}
```

### 7.3 Горизонтальное масштабирование

```yaml
# kubernetes/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: neurogen-news
spec:
  replicas: 3
  selector:
    matchLabels:
      app: neurogen-news
  template:
    metadata:
      labels:
        app: neurogen-news
    spec:
      containers:
      - name: app
        image: neurogen/news:latest
        ports:
        - containerPort: 8080
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: url
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: neurogen-news-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: neurogen-news
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

---

## 8. Мониторинг и observability

### 8.1 Метрики (Prometheus)

```go
// internal/metrics/metrics.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    HTTPRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
    
    HTTPRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
        },
        []string{"method", "path"},
    )
    
    DBQueryDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "db_query_duration_seconds",
            Help:    "Database query duration in seconds",
            Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5},
        },
        []string{"query"},
    )
    
    CacheHits = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "cache_hits_total",
            Help: "Total number of cache hits",
        },
        []string{"cache"},
    )
    
    CacheMisses = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "cache_misses_total",
            Help: "Total number of cache misses",
        },
        []string{"cache"},
    )
    
    ActiveWebSockets = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_websocket_connections",
            Help: "Number of active WebSocket connections",
        },
    )
)
```

### 8.2 Логирование (Zap)

```go
// internal/pkg/logger/logger.go
package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func NewLogger(env string) (*zap.Logger, error) {
    var config zap.Config
    
    if env == "production" {
        config = zap.NewProductionConfig()
        config.EncoderConfig.TimeKey = "timestamp"
        config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    } else {
        config = zap.NewDevelopmentConfig()
        config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
    }
    
    return config.Build()
}

// Structured logging middleware
func LoggingMiddleware(logger *zap.Logger) fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()
        
        err := c.Next()
        
        logger.Info("request",
            zap.String("method", c.Method()),
            zap.String("path", c.Path()),
            zap.Int("status", c.Response().StatusCode()),
            zap.Duration("latency", time.Since(start)),
            zap.String("ip", c.IP()),
            zap.String("user_agent", c.Get("User-Agent")),
        )
        
        return err
    }
}
```

---

## 9. Безопасность

### 9.1 Аутентификация (JWT)

```go
// internal/pkg/jwt/jwt.go
package jwt

import (
    "time"
    
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID   string `json:"uid"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

type JWTService struct {
    secretKey     []byte
    accessExpiry  time.Duration
    refreshExpiry time.Duration
}

func NewJWTService(secretKey string) *JWTService {
    return &JWTService{
        secretKey:     []byte(secretKey),
        accessExpiry:  15 * time.Minute,
        refreshExpiry: 7 * 24 * time.Hour,
    }
}

func (s *JWTService) GenerateAccessToken(user *model.User) (string, error) {
    claims := Claims{
        UserID:   user.ID,
        Username: user.Username,
        Role:     string(user.Role),
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessExpiry)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "neurogen.news",
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.secretKey)
}

func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return s.secretKey, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, fmt.Errorf("invalid token")
}
```

### 9.2 Security Headers

```go
// internal/middleware/security.go
package middleware

import "github.com/gofiber/fiber/v2"

func SecurityHeaders() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Prevent XSS
        c.Set("X-XSS-Protection", "1; mode=block")
        c.Set("X-Content-Type-Options", "nosniff")
        
        // Prevent clickjacking
        c.Set("X-Frame-Options", "SAMEORIGIN")
        
        // HSTS
        c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        
        // CSP
        c.Set("Content-Security-Policy", 
            "default-src 'self'; "+
            "script-src 'self' 'unsafe-inline'; "+
            "style-src 'self' 'unsafe-inline'; "+
            "img-src 'self' data: https:; "+
            "font-src 'self' data:; "+
            "connect-src 'self' wss:")
        
        // Referrer Policy
        c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
        
        // Permissions Policy
        c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
        
        return c.Next()
    }
}
```

---

## 10. Документация для команды

### 10.1 Структура документации

```
docs/
├── README.md                 # Обзор проекта
├── CONTRIBUTING.md           # Как контрибьютить
├── ARCHITECTURE.md           # Архитектура системы
│
├── api/
│   ├── openapi.yaml         # OpenAPI спецификация
│   └── README.md            # Как использовать API
│
├── frontend/
│   ├── README.md            # Настройка frontend
│   ├── COMPONENTS.md        # Документация компонентов
│   └── STYLING.md           # Гайд по стилям
│
├── backend/
│   ├── README.md            # Настройка backend
│   ├── DATABASE.md          # Схема БД, миграции
│   └── CACHING.md           # Стратегия кэширования
│
├── deployment/
│   ├── README.md            # Деплой инструкции
│   ├── DOCKER.md            # Docker конфигурация
│   └── KUBERNETES.md        # K8s манифесты
│
└── guides/
    ├── QUICKSTART.md        # Быстрый старт
    ├── DEVELOPMENT.md       # Разработка
    └── TESTING.md           # Тестирование
```

### 10.2 Быстрый старт (QUICKSTART.md)

```markdown
# Быстрый старт

## Требования

- Go 1.23+
- Node.js 20+
- Docker & Docker Compose
- Make

## Запуск в dev-режиме

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/neurogen/news.git
   cd news
   ```

2. Скопируйте .env файл:
   ```bash
   cp .env.example .env
   ```

3. Запустите инфраструктуру:
   ```bash
   docker-compose up -d postgres redis meilisearch
   ```

4. Примените миграции:
   ```bash
   make migrate-up
   ```

5. Запустите dev-серверы:
   ```bash
   make dev
   ```

6. Откройте http://localhost:5173 (frontend) или http://localhost:8080 (backend)

## Сборка production-бинарника

```bash
make build
./neurogen-news
```

## Структура проекта

[Ссылка на ARCHITECTURE.md]
```

---

## Сводка технологий

| Категория | Технология | Версия |
|-----------|------------|--------|
| **Frontend** | Vue.js | 3.5+ |
| | TypeScript | 5.6+ |
| | Vite | 6.0+ |
| | Tailwind CSS | 3.4+ |
| | TipTap | 2.9+ |
| **Backend** | Go | 1.23+ |
| | Fiber | v2.52+ |
| | sqlc | 1.27+ |
| | pgx | v5.7+ |
| **Database** | PostgreSQL | 17 |
| | Redis | 7.4+ |
| **Search** | Meilisearch | 1.11+ |
| **Infra** | Docker | 27+ |
| | Cloudflare | CDN + WAF |


