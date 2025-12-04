# Архитектура Neurogen.News

> **Стек:** Vue 3 + Go | Сборка в единый бинарник | Рассчитан на высокую нагрузку

## Общая архитектура системы

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
                               Load Balancer
                                      │
         ┌────────────────────────────┼────────────────────────────┐
         │                            │                            │
┌────────▼──────────────────┐ ┌───────▼───────────────────┐ ┌──────▼────────────────────┐
│   ЕДИНЫЙ GO БИНАРНИК #1   │ │   ЕДИНЫЙ GO БИНАРНИК #2   │ │   ЕДИНЫЙ GO БИНАРНИК #N   │
│  ┌──────────────────────┐ │ │  ┌──────────────────────┐ │ │  ┌──────────────────────┐ │
│  │  Vue 3 SPA (embed)   │ │ │  │  Vue 3 SPA (embed)   │ │ │  │  Vue 3 SPA (embed)   │ │
│  └──────────────────────┘ │ │  └──────────────────────┘ │ │  └──────────────────────┘ │
│  ┌──────────────────────┐ │ │  ┌──────────────────────┐ │ │  ┌──────────────────────┐ │
│  │  Go HTTP Server      │ │ │  │  Go HTTP Server      │ │ │  │  Go HTTP Server      │ │
│  │  (Fiber v2.52+)      │ │ │  │  (Fiber v2.52+)      │ │ │  │  (Fiber v2.52+)      │ │
│  └──────────────────────┘ │ │  └──────────────────────┘ │ │  └──────────────────────┘ │
└────────────┬──────────────┘ └───────────┬───────────────┘ └──────────────┬───────────┘
             │                            │                                │
             └────────────────────────────┼────────────────────────────────┘
                                          │
         ┌────────────────────────────────┼────────────────────────────────┐
         │                                │                                │
┌────────▼────────┐            ┌──────────▼──────────┐          ┌──────────▼──────────┐
│   PostgreSQL    │            │       Redis         │          │    Meilisearch      │
│   17 (Primary)  │            │  7.4+ (Cache/PubSub)│          │  1.11+ (Search)     │
└────────┬────────┘            └─────────────────────┘          └─────────────────────┘
         │
         │ Репликация (Read Replicas)
         ▼
┌─────────────────┐
│ PostgreSQL      │
│ (Read Replica)  │
└─────────────────┘
```

### Ключевые принципы

| Принцип | Реализация |
|---------|------------|
| **Единый бинарник** | Vue SPA встроен в Go через `go:embed` |
| **Stateless** | Вся сессионная информация в Redis |
| **Горизонтальное масштабирование** | Добавление инстансов без изменения кода |
| **Чистая архитектура** | handler → service → repository |
| **Type-safe SQL** | sqlc генерирует Go-код из SQL |

---

## Структура проекта

```
neurogen-news/
├── frontend/                   # Vue 3 приложение
│   ├── src/
│   │   ├── components/         # Vue компоненты
│   │   │   ├── ui/             # Базовые UI (Button, Input, Modal...)
│   │   │   ├── layout/         # AppHeader, AppSidebar, UserMenu...
│   │   │   ├── feed/           # ArticleCard, FeedList, FeedFilters...
│   │   │   ├── article/        # ArticleContent, ReactionPanel...
│   │   │   ├── comments/       # CommentList, CommentForm, CommentThread...
│   │   │   ├── editor/         # RichTextEditor, Toolbar, BlockPicker...
│   │   │   ├── user/           # UserProfile, AchievementBadge...
│   │   │   ├── tools/          # ToolCard, ToolGrid, ToolComparison...
│   │   │   └── seo/            # SeoHead, SchemaOrg, Breadcrumbs...
│   │   │
│   │   ├── pages/              # Vue Router страницы
│   │   │   ├── index.vue       # Главная (популярное)
│   │   │   ├── new.vue         # Свежее
│   │   │   ├── my.vue          # Моя лента
│   │   │   ├── start.vue       # С чего начать
│   │   │   ├── [category]/
│   │   │   │   ├── index.vue   # Страница категории
│   │   │   │   └── [slug].vue  # Страница статьи
│   │   │   ├── tools/
│   │   │   │   ├── index.vue   # Каталог инструментов
│   │   │   │   └── [slug].vue  # Карточка инструмента
│   │   │   ├── questions/      # Q&A раздел
│   │   │   ├── user/           # Настройки, черновики, закладки
│   │   │   └── editor/         # Редактор статей
│   │   │
│   │   ├── composables/        # Vue Composables
│   │   │   ├── useAuth.ts
│   │   │   ├── useApi.ts
│   │   │   ├── useSeo.ts
│   │   │   ├── useWebSocket.ts
│   │   │   └── ...
│   │   │
│   │   ├── stores/             # Pinia Stores
│   │   │   ├── auth.ts
│   │   │   ├── user.ts
│   │   │   ├── feed.ts
│   │   │   └── ...
│   │   │
│   │   ├── api/                # API клиент
│   │   ├── types/              # TypeScript типы
│   │   └── utils/              # Утилиты
│   │
│   ├── vite.config.ts
│   ├── tailwind.config.ts
│   └── package.json
│
├── backend/                    # Go приложение
│   ├── cmd/
│   │   └── server/
│   │       └── main.go         # Точка входа
│   │
│   ├── internal/
│   │   ├── config/             # Конфигурация
│   │   ├── server/             # HTTP сервер, роуты
│   │   ├── handler/            # HTTP handlers
│   │   │   ├── article.go
│   │   │   ├── comment.go
│   │   │   ├── user.go
│   │   │   ├── auth.go
│   │   │   ├── draft.go
│   │   │   ├── search.go
│   │   │   └── seo.go
│   │   │
│   │   ├── service/            # Бизнес-логика
│   │   │   ├── article.go
│   │   │   ├── comment.go
│   │   │   ├── user.go
│   │   │   ├── achievement.go
│   │   │   └── ...
│   │   │
│   │   ├── repository/         # Работа с БД (sqlc)
│   │   │   ├── queries/        # SQL запросы
│   │   │   └── *.sql.go        # Сгенерированный код
│   │   │
│   │   ├── cache/              # Redis кэширование
│   │   ├── search/             # Meilisearch
│   │   ├── websocket/          # Real-time
│   │   ├── middleware/         # Auth, RateLimit, CORS...
│   │   ├── model/              # Доменные модели
│   │   └── dto/                # Request/Response объекты
│   │
│   ├── web/
│   │   └── dist/               # Собранный Vue SPA (go:embed)
│   │
│   ├── migrations/             # SQL миграции
│   ├── go.mod
│   └── go.sum
│
├── docs/                       # Документация
├── scripts/                    # Скрипты сборки
├── docker-compose.yml
├── Dockerfile
└── Makefile
```

## Модули системы

### 1. Система контента

#### Типы контента
```typescript
enum ContentType {
  ARTICLE = 'article',      // Полноценная статья
  NEWS = 'news',            // Короткая новость
  POST = 'post',            // Пост пользователя (UGC)
  QUESTION = 'question',    // Вопрос сообществу
  DISCUSSION = 'discussion' // Дискуссия
}

interface Article {
  id: string;
  type: ContentType;
  title: string;
  slug: string;
  lead: string;              // Подзаголовок/лид
  content: JSONContent;      // TipTap JSON
  coverImage?: Media;
  author: User;
  subsite: Subsite;
  tags: Tag[];
  isEditorialContent: boolean; // Материал редакции
  
  // Статистика
  viewCount: number;
  commentCount: number;
  reactionCount: number;
  bookmarkCount: number;
  
  // Даты
  createdAt: Date;
  updatedAt: Date;
  publishedAt?: Date;
  
  // Модерация
  status: ArticleStatus;
  isBlocked: boolean;
}

enum ArticleStatus {
  DRAFT = 'draft',
  PENDING_REVIEW = 'pending_review',
  PUBLISHED = 'published',
  ARCHIVED = 'archived'
}
```

#### Система черновиков
```typescript
interface Draft {
  id: string;
  userId: string;
  title?: string;
  content: JSONContent;
  coverImage?: Media;
  subsiteId?: string;
  tags: string[];
  
  // Версионирование
  version: number;
  versions: DraftVersion[];
  
  // Автосохранение
  lastAutoSaveAt: Date;
  
  createdAt: Date;
  updatedAt: Date;
}

interface DraftVersion {
  id: string;
  draftId: string;
  version: number;
  content: JSONContent;
  createdAt: Date;
}
```

#### Подсайты (категории)
```typescript
interface Subsite {
  id: string;
  slug: string;            // 'neural-networks', 'generation', etc.
  name: string;
  description: string;
  icon: string;            // Emoji или URL иконки
  coverImage?: string;
  subscriberCount: number;
  
  // Настройки
  isOfficial: boolean;     // Официальный подсайт
  requiresModeration: boolean;
  allowedContentTypes: ContentType[];
}
```

### 2. Система пользователей

```typescript
interface User {
  id: string;
  username: string;
  displayName: string;
  email: string;
  avatarUrl?: string;
  coverUrl?: string;
  bio?: string;
  
  // Статус
  role: UserRole;
  isVerified: boolean;      // Синяя галочка
  isPro: boolean;           // Премиум подписка
  karma: number;
  
  // Статистика
  followerCount: number;
  followingCount: number;
  articleCount: number;
  commentCount: number;
  
  // Связи
  links: SocialLink[];
  badges: Badge[];
  achievements: Achievement[];
  
  // Настройки
  settings: UserSettings;
  
  createdAt: Date;
  lastActiveAt: Date;
}

enum UserRole {
  USER = 'user',
  AUTHOR = 'author',        // Может публиковать
  MODERATOR = 'moderator',
  EDITOR = 'editor',        // Редакция
  ADMIN = 'admin'
}

interface UserSettings {
  // Настройки блога
  blog: {
    name: string;
    description: string;
    slug: string;
    telegramLink?: string;
    avatar?: string;
    cover?: string;
  };
  
  // Настройки лент
  feeds: {
    preferredSubsites: string[];
    hiddenSubsites: string[];
    contentFilter: 'all' | 'editorial' | 'community';
  };
  
  // Основные настройки
  general: {
    timezone: string;
    interests: string[];
    showInJobSearch: boolean;  // Чекбокс для карьеры
    language: 'ru' | 'en';
  };
  
  // Настройки уведомлений
  notifications: {
    email: {
      comments: boolean;
      replies: boolean;
      mentions: boolean;
      followers: boolean;
      digest: 'daily' | 'weekly' | 'never';
    };
    push: {
      comments: boolean;
      replies: boolean;
      mentions: boolean;
      followers: boolean;
    };
  };
}
```

### 3. Система достижений (Ачивки)

```typescript
interface Achievement {
  id: string;
  slug: string;
  name: string;
  description: string;
  icon: string;
  category: AchievementCategory;
  
  // Условия получения
  condition: AchievementCondition;
  threshold: number;
  
  // Награды
  points: number;
  badge?: Badge;
  
  isHidden: boolean;  // Секретные достижения
}

enum AchievementCategory {
  WRITING = 'writing',       // За статьи
  COMMENTING = 'commenting', // За комментарии
  SOCIAL = 'social',         // За подписчиков
  ENGAGEMENT = 'engagement', // За реакции
  SPECIAL = 'special'        // Особые
}

interface UserAchievement {
  userId: string;
  achievementId: string;
  unlockedAt: Date;
  progress: number;          // Прогресс до достижения (0-100)
}

// Примеры достижений (по аналогии с VC.ru)
const ACHIEVEMENTS = [
  { name: 'Первая статья', condition: 'articles_count >= 1', points: 10 },
  { name: 'Автор', condition: 'articles_count >= 10', points: 50 },
  { name: 'Плодовитый автор', condition: 'articles_count >= 50', points: 200 },
  { name: 'Топ-комментатор', condition: 'comments_count >= 100', points: 100 },
  { name: 'Популярный', condition: 'followers_count >= 100', points: 150 },
  { name: 'Звезда', condition: 'followers_count >= 1000', points: 500 },
  { name: 'Ранний пользователь', condition: 'registered_before_launch', points: 100 },
];
```

### 4. Система закладок

```typescript
interface Bookmark {
  id: string;
  userId: string;
  articleId: string;
  
  // Организация
  folderId?: string;
  note?: string;
  
  createdAt: Date;
}

interface BookmarkFolder {
  id: string;
  userId: string;
  name: string;
  isPrivate: boolean;
  
  createdAt: Date;
}
```

### 5. Система уведомлений

```typescript
enum NotificationType {
  COMMENT = 'comment',           // Комментарий к статье
  REPLY = 'reply',               // Ответ на комментарий
  MENTION = 'mention',           // Упоминание @username
  REACTION = 'reaction',         // Реакция на статью/комментарий
  FOLLOWER = 'follower',         // Новый подписчик
  PUBLICATION = 'publication',   // Новая статья от подписки
  MODERATION = 'moderation',     // Действия модерации
  ACHIEVEMENT = 'achievement',   // Получено достижение
  SYSTEM = 'system'              // Системное уведомление
}

interface Notification {
  id: string;
  userId: string;
  type: NotificationType;
  
  // Контент
  title: string;
  message: string;
  link?: string;
  
  // Связанные сущности
  actorId?: string;              // Кто вызвал уведомление
  articleId?: string;
  commentId?: string;
  
  // Статус
  isRead: boolean;
  readAt?: Date;
  
  createdAt: Date;
}
```

### 6. Система реакций

По аналогии с DTF/VC — расширенная система эмодзи реакций:

```typescript
interface Reaction {
  id: string;
  emoji: string;           // '👍', '❤️', '😂', '🤔', '😢', '😡'
  label: string;           // 'Нравится', 'Супер', 'Смешно', etc.
  weight: number;          // Влияние на karma (+1, -1, etc.)
}

interface UserReaction {
  userId: string;
  targetId: string;
  targetType: 'article' | 'comment';
  reactionId: string;
  createdAt: Date;
}
```

### 7. Система комментариев

```typescript
interface Comment {
  id: string;
  content: string;
  authorId: string;
  articleId: string;
  parentId?: string;        // Для вложенных комментариев
  
  // Реакции
  reactionCounts: Record<string, number>;
  
  // Модерация
  isDeleted: boolean;
  isEdited: boolean;
  editedAt?: Date;
  
  createdAt: Date;
  updatedAt: Date;
}
```

### 8. Система подписок

```typescript
interface Subscription {
  id: string;
  userId: string;
  
  // Можно подписаться на:
  targetType: 'user' | 'subsite' | 'tag';
  targetId: string;
  
  // Настройки уведомлений
  notifyEmail: boolean;
  notifyPush: boolean;
  
  createdAt: Date;
}
```

---

## Редактор статей

### Архитектура редактора

```typescript
// Конфигурация редактора (TipTap)
const EDITOR_EXTENSIONS = [
  // Базовые
  StarterKit,
  
  // Форматирование текста
  Bold,
  Italic,
  Underline,
  Strike,
  
  // Структура
  Heading.configure({ levels: [2, 3] }),
  BulletList,
  OrderedList,
  Blockquote,
  CodeBlock,
  
  // Медиа
  Image,
  Video,
  Gallery,
  
  // Специальные блоки
  Spoiler,
  Embed,
  ArticleCard,
  
  // Утилиты
  Link,
  Placeholder,
  CharacterCount,
];

interface EditorState {
  // Контент
  content: JSONContent;
  
  // Мета
  title: string;
  subsiteId: string;
  tags: string[];
  coverImage?: Media;
  
  // Статус
  isDirty: boolean;
  isSaving: boolean;
  lastSavedAt?: Date;
  
  // Версии
  currentVersion: number;
  versions: DraftVersion[];
}
```

### Блоки редактора

```typescript
enum EditorBlockType {
  // Текстовые
  PARAGRAPH = 'paragraph',
  HEADING = 'heading',
  BLOCKQUOTE = 'blockquote',
  CODE_BLOCK = 'codeBlock',
  LIST = 'list',
  
  // Медиа
  IMAGE = 'image',
  GALLERY = 'gallery',
  VIDEO = 'video',
  
  // Embed
  TWITTER = 'twitter',
  TELEGRAM = 'telegram',
  YOUTUBE = 'youtube',
  
  // Специальные
  SPOILER = 'spoiler',
  DIVIDER = 'divider',
  ARTICLE_CARD = 'articleCard',
}
```

---

## Кэширование

### Стратегия кэширования

```typescript
// Уровни кэша
const CACHE_LAYERS = {
  // Edge (CDN) — статика, изображения
  edge: {
    static: '1y',
    images: '1d',
    api: '1m'  // stale-while-revalidate
  },
  
  // Redis — горячие данные
  redis: {
    feed: '5m',           // Лента популярного
    article: '10m',       // Контент статьи
    userProfile: '5m',    // Профиль пользователя
    subsiteList: '1h',    // Список подсайтов
    reactionCounts: '1m', // Счётчики реакций
    achievements: '10m',  // Достижения пользователя
    notifications: '1m',  // Уведомления (инвалидируются часто)
    drafts: '1m',         // Черновики
  },
  
  // In-memory — сессии, rate limiting
  memory: {
    session: '24h',
    rateLimit: '1m'
  }
};
```

### Инвалидация кэша

```typescript
// События инвалидации
const INVALIDATION_EVENTS = {
  'article.publish': ['feed:*', 'subsite:{subsiteId}'],
  'article.update': ['article:{id}'],
  'article.delete': ['feed:*', 'article:{id}'],
  'comment.create': ['article:{articleId}:comments'],
  'reaction.add': ['article:{articleId}:reactions'],
  'draft.save': ['drafts:{userId}'],
  'bookmark.add': ['bookmarks:{userId}'],
  'notification.create': ['notifications:{userId}'],
  'achievement.unlock': ['achievements:{userId}', 'user:{userId}'],
  'settings.update': ['settings:{userId}', 'user:{userId}'],
};
```

---

## Real-time функционал

### WebSocket события

```typescript
// Подключение к каналам
const CHANNELS = {
  global: 'global',                    // Глобальные события
  article: 'article:{id}',             // События статьи
  subsite: 'subsite:{slug}',           // События подсайта
  user: 'user:{id}',                   // Личные уведомления
};

// События
interface WSEvent {
  type: 
    | 'comment.new'
    | 'reaction.update'
    | 'article.update'
    | 'online.count'
    | 'notification.new'
    | 'achievement.unlocked';
  payload: unknown;
  timestamp: Date;
}
```

---

## Безопасность

### Rate Limiting

```typescript
const RATE_LIMITS = {
  anonymous: {
    requests: 100,
    window: '1m'
  },
  authenticated: {
    requests: 1000,
    window: '1m'
  },
  article: {
    create: { count: 5, window: '1h' },
    edit: { count: 20, window: '1h' }
  },
  comment: {
    create: { count: 30, window: '1h' }
  },
  reaction: {
    create: { count: 100, window: '1h' }
  },
  draft: {
    autoSave: { count: 60, window: '1h' }
  }
};
```

### Модерация контента

```typescript
enum ModerationStatus {
  PENDING = 'pending',
  APPROVED = 'approved',
  REJECTED = 'rejected',
  NEEDS_REVIEW = 'needs_review'
}

// Автомодерация
const AUTO_MODERATION = {
  spam: 'akismet',
  profanity: 'dictionary',
  imageNSFW: 'ml-model',
  duplicateContent: 'hash-check'
};
```

---

## Масштабирование

### Горизонтальное масштабирование

```
                    ┌─────────────┐
                    │   Clients   │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │ Load Balancer│
                    └──────┬──────┘
                           │
        ┌─────────────┬────┴────┬─────────────┐
        │             │         │             │
   ┌────▼────┐  ┌─────▼────┐ ┌──▼────┐  ┌────▼────┐
   │ App 1   │  │  App 2   │ │App 3  │  │ App N   │
   └────┬────┘  └─────┬────┘ └───┬───┘  └────┬────┘
        │             │          │            │
        └─────────────┴────┬─────┴────────────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
        ┌─────▼────┐ ┌─────▼────┐ ┌─────▼────┐
        │ PG Master│ │PG Replica│ │PG Replica│
        └──────────┘ └──────────┘ └──────────┘
```

### Шардинг данных

```typescript
// Шардинг по подсайтам для изоляции нагрузки
const SHARDING_STRATEGY = {
  articles: {
    key: 'subsite_id',
    shards: 8
  },
  comments: {
    key: 'article_id',
    shards: 16
  },
  notifications: {
    key: 'user_id',
    shards: 8
  }
};
```

---

## Мониторинг

### Метрики

```typescript
const METRICS = {
  // Бизнес метрики
  business: [
    'articles.published.count',
    'comments.created.count',
    'users.registered.count',
    'reactions.total.count',
    'feed.views.count',
    'drafts.created.count',
    'bookmarks.added.count',
    'achievements.unlocked.count',
    'notifications.sent.count'
  ],
  
  // Технические метрики
  technical: [
    'response.time.p95',
    'error.rate',
    'cache.hit.rate',
    'db.query.time',
    'ws.connections.active',
    'editor.autosave.latency'
  ]
};
```

### Алерты

```typescript
const ALERTS = {
  critical: {
    'error.rate > 5%': 'pagerduty',
    'response.time.p95 > 3s': 'pagerduty',
    'db.connections > 90%': 'pagerduty'
  },
  warning: {
    'cache.hit.rate < 80%': 'slack',
    'queue.lag > 1000': 'slack',
    'notifications.delivery.failed > 1%': 'slack'
  }
};
```
