# Схема базы данных Neurogen.News

## Обзор

База данных Neurogen.News построена на PostgreSQL с использованием Prisma ORM.
Структура оптимизирована для UGC-платформы с высокой нагрузкой на чтение.

---

## ER-диаграмма (основные сущности)

```
┌──────────────┐       ┌──────────────┐       ┌──────────────┐
│    User      │───────│   Article    │───────│   Comment    │
│              │       │              │       │              │
│ id           │       │ id           │       │ id           │
│ email        │       │ title        │       │ content      │
│ username     │───┐   │ slug         │   ┌───│ authorId     │
│ displayName  │   │   │ content      │   │   │ articleId    │
│ avatarUrl    │   │   │ authorId ────│───┘   │ parentId     │
│ ...          │   │   │ subsiteId    │       │ ...          │
└──────────────┘   │   │ ...          │       └──────────────┘
                   │   └──────────────┘
                   │          │
                   │          │
                   │   ┌──────────────┐
                   │   │   Subsite    │
                   │   │              │
                   └───│ id           │
                       │ slug         │
                       │ name         │
                       │ ...          │
                       └──────────────┘
```

---

## Prisma Schema

### User (Пользователь)

```prisma
model User {
  id            String    @id @default(cuid())
  email         String    @unique
  username      String    @unique
  displayName   String
  avatarUrl     String?
  coverUrl      String?
  bio           String?
  
  // Статус
  role          UserRole  @default(USER)
  isVerified    Boolean   @default(false)
  isPlusActive  Boolean   @default(false)
  plusExpiresAt DateTime?
  karma         Int       @default(0)
  
  // Счётчики (денормализация для производительности)
  followerCount  Int      @default(0)
  followingCount Int      @default(0)
  articleCount   Int      @default(0)
  commentCount   Int      @default(0)
  
  // Настройки
  settings      Json      @default("{}")
  
  // Связи
  articles      Article[]
  comments      Comment[]
  reactions     Reaction[]
  bookmarks     Bookmark[]
  notifications Notification[]
  achievements  UserAchievement[]
  
  // Подписки
  followers     Follow[]  @relation("followers")
  following     Follow[]  @relation("following")
  subsiteSubscriptions SubsiteSubscription[]
  
  // Метки времени
  createdAt     DateTime  @default(now())
  updatedAt     DateTime  @updatedAt
  lastActiveAt  DateTime  @default(now())
  
  @@index([username])
  @@index([email])
  @@index([createdAt])
  @@index([karma])
}

enum UserRole {
  USER
  AUTHOR
  MODERATOR
  EDITOR
  ADMIN
}
```

### Article (Статья)

```prisma
model Article {
  id            String        @id @default(cuid())
  
  // Контент
  title         String
  slug          String        @unique
  lead          String?       // Подзаголовок/лид
  content       Json          // TipTap JSON
  coverUrl      String?
  
  // Связи
  author        User          @relation(fields: [authorId], references: [id])
  authorId      String
  subsite       Subsite       @relation(fields: [subsiteId], references: [id])
  subsiteId     String
  
  // Тип и статус
  type          ContentType   @default(ARTICLE)
  status        ArticleStatus @default(DRAFT)
  isEditorial   Boolean       @default(false) // Материал редакции
  isPinned      Boolean       @default(false)
  
  // Счётчики
  viewCount     Int           @default(0)
  commentCount  Int           @default(0)
  reactionCount Int           @default(0)
  bookmarkCount Int           @default(0)
  repostCount   Int           @default(0)
  
  // SEO
  metaTitle       String?
  metaDescription String?
  
  // Связи
  comments      Comment[]
  reactions     ArticleReaction[]
  bookmarks     Bookmark[]
  tags          ArticleTag[]
  reposts       Repost[]
  
  // Метки времени
  createdAt     DateTime      @default(now())
  updatedAt     DateTime      @updatedAt
  publishedAt   DateTime?
  
  @@index([authorId])
  @@index([subsiteId])
  @@index([status, publishedAt])
  @@index([createdAt])
  @@index([viewCount])
  @@index([slug])
}

enum ContentType {
  ARTICLE
  NEWS
  POST
  QUESTION
  DISCUSSION
}

enum ArticleStatus {
  DRAFT
  PENDING_REVIEW
  PUBLISHED
  REJECTED
  ARCHIVED
}
```

### Comment (Комментарий)

```prisma
model Comment {
  id            String    @id @default(cuid())
  content       String    @db.Text
  
  // Связи
  author        User      @relation(fields: [authorId], references: [id])
  authorId      String
  article       Article   @relation(fields: [articleId], references: [id], onDelete: Cascade)
  articleId     String
  
  // Вложенность
  parent        Comment?  @relation("CommentReplies", fields: [parentId], references: [id])
  parentId      String?
  replies       Comment[] @relation("CommentReplies")
  
  // Счётчики
  reactionCount Int       @default(0)
  replyCount    Int       @default(0)
  
  // Статус
  isDeleted     Boolean   @default(false)
  isEdited      Boolean   @default(false)
  isPinned      Boolean   @default(false)
  
  // Связи
  reactions     CommentReaction[]
  
  // Метки времени
  createdAt     DateTime  @default(now())
  updatedAt     DateTime  @updatedAt
  
  @@index([articleId, createdAt])
  @@index([authorId])
  @@index([parentId])
}
```

### Subsite (Категория/Подсайт)

```prisma
model Subsite {
  id            String    @id @default(cuid())
  slug          String    @unique
  name          String
  description   String?
  icon          String?   // Emoji или URL
  coverUrl      String?
  
  // Настройки
  isOfficial    Boolean   @default(false)
  requiresModeration Boolean @default(false)
  isHidden      Boolean   @default(false)
  
  // Счётчики
  subscriberCount Int     @default(0)
  articleCount    Int     @default(0)
  
  // Связи
  articles      Article[]
  subscriptions SubsiteSubscription[]
  
  // Метки времени
  createdAt     DateTime  @default(now())
  updatedAt     DateTime  @updatedAt
  
  @@index([slug])
  @@index([subscriberCount])
}
```

### Reactions (Реакции)

```prisma
model ReactionType {
  id      String  @id @default(cuid())
  emoji   String  @unique
  label   String
  weight  Int     @default(1) // Влияние на karma
  
  articleReactions ArticleReaction[]
  commentReactions CommentReaction[]
}

model ArticleReaction {
  id            String       @id @default(cuid())
  
  user          User         @relation(fields: [userId], references: [id])
  userId        String
  article       Article      @relation(fields: [articleId], references: [id], onDelete: Cascade)
  articleId     String
  reactionType  ReactionType @relation(fields: [reactionTypeId], references: [id])
  reactionTypeId String
  
  createdAt     DateTime     @default(now())
  
  @@unique([userId, articleId, reactionTypeId])
  @@index([articleId])
}

model CommentReaction {
  id            String       @id @default(cuid())
  
  user          User         @relation(fields: [userId], references: [id])
  userId        String
  comment       Comment      @relation(fields: [commentId], references: [id], onDelete: Cascade)
  commentId     String
  reactionType  ReactionType @relation(fields: [reactionTypeId], references: [id])
  reactionTypeId String
  
  createdAt     DateTime     @default(now())
  
  @@unique([userId, commentId, reactionTypeId])
  @@index([commentId])
}
```

### Follow (Подписки)

```prisma
model Follow {
  id            String   @id @default(cuid())
  
  follower      User     @relation("following", fields: [followerId], references: [id])
  followerId    String
  following     User     @relation("followers", fields: [followingId], references: [id])
  followingId   String
  
  // Настройки уведомлений
  notifyEmail   Boolean  @default(true)
  notifyPush    Boolean  @default(true)
  
  createdAt     DateTime @default(now())
  
  @@unique([followerId, followingId])
  @@index([followerId])
  @@index([followingId])
}

model SubsiteSubscription {
  id            String   @id @default(cuid())
  
  user          User     @relation(fields: [userId], references: [id])
  userId        String
  subsite       Subsite  @relation(fields: [subsiteId], references: [id])
  subsiteId     String
  
  // Настройки уведомлений
  notifyEmail   Boolean  @default(false)
  notifyPush    Boolean  @default(true)
  
  createdAt     DateTime @default(now())
  
  @@unique([userId, subsiteId])
  @@index([userId])
  @@index([subsiteId])
}
```

### Bookmark (Закладки)

```prisma
model Bookmark {
  id            String       @id @default(cuid())
  
  user          User         @relation(fields: [userId], references: [id])
  userId        String
  article       Article?     @relation(fields: [articleId], references: [id], onDelete: Cascade)
  articleId     String?
  
  // Можно расширить для комментариев
  targetType    BookmarkType @default(ARTICLE)
  
  createdAt     DateTime     @default(now())
  
  @@unique([userId, articleId])
  @@index([userId, createdAt])
}

enum BookmarkType {
  ARTICLE
  COMMENT
}
```

### Tag (Теги)

```prisma
model Tag {
  id            String       @id @default(cuid())
  slug          String       @unique
  name          String
  
  // Счётчики
  articleCount  Int          @default(0)
  
  articles      ArticleTag[]
  
  createdAt     DateTime     @default(now())
  
  @@index([slug])
  @@index([articleCount])
}

model ArticleTag {
  id            String   @id @default(cuid())
  
  article       Article  @relation(fields: [articleId], references: [id], onDelete: Cascade)
  articleId     String
  tag           Tag      @relation(fields: [tagId], references: [id])
  tagId         String
  
  @@unique([articleId, tagId])
  @@index([tagId])
}
```

### Notification (Уведомления)

```prisma
model Notification {
  id            String           @id @default(cuid())
  
  user          User             @relation(fields: [userId], references: [id])
  userId        String
  
  type          NotificationType
  title         String
  message       String?
  data          Json?            // Дополнительные данные
  
  // Ссылка на объект
  targetType    String?          // article, comment, user
  targetId      String?
  
  // Статус
  isRead        Boolean          @default(false)
  isEmailSent   Boolean          @default(false)
  isPushSent    Boolean          @default(false)
  
  createdAt     DateTime         @default(now())
  readAt        DateTime?
  
  @@index([userId, isRead, createdAt])
  @@index([userId, createdAt])
}

enum NotificationType {
  NEW_FOLLOWER
  NEW_COMMENT
  COMMENT_REPLY
  MENTION
  REACTION
  ARTICLE_PUBLISHED
  SYSTEM
  ACHIEVEMENT
}
```

### Achievement (Достижения)

```prisma
model Achievement {
  id            String            @id @default(cuid())
  slug          String            @unique
  name          String
  description   String
  icon          String            // URL или emoji
  
  // Условия получения
  category      AchievementCategory
  threshold     Int               // Порог для получения
  
  users         UserAchievement[]
  
  createdAt     DateTime          @default(now())
  
  @@index([category])
}

enum AchievementCategory {
  ARTICLES      // Количество статей
  COMMENTS      // Количество комментариев
  FOLLOWERS     // Количество подписчиков
  REACTIONS     // Количество полученных реакций
  SPECIAL       // Особые достижения
}

model UserAchievement {
  id            String      @id @default(cuid())
  
  user          User        @relation(fields: [userId], references: [id])
  userId        String
  achievement   Achievement @relation(fields: [achievementId], references: [id])
  achievementId String
  
  // Прогресс (для незавершённых)
  progress      Int         @default(0)
  isCompleted   Boolean     @default(false)
  
  completedAt   DateTime?
  createdAt     DateTime    @default(now())
  
  @@unique([userId, achievementId])
  @@index([userId, isCompleted])
}
```

### Session & Auth (Сессии и авторизация)

```prisma
model Session {
  id            String   @id @default(cuid())
  
  userId        String
  token         String   @unique
  
  userAgent     String?
  ipAddress     String?
  
  expiresAt     DateTime
  createdAt     DateTime @default(now())
  lastUsedAt    DateTime @default(now())
  
  @@index([userId])
  @@index([token])
  @@index([expiresAt])
}

model AuthProvider {
  id            String   @id @default(cuid())
  
  userId        String
  provider      String   // google, vk, telegram
  providerId    String
  
  createdAt     DateTime @default(now())
  
  @@unique([provider, providerId])
  @@index([userId])
}
```

### Draft (Черновики)

```prisma
model Draft {
  id            String    @id @default(cuid())
  
  author        User      @relation(fields: [authorId], references: [id])
  authorId      String
  
  title         String?
  content       Json?     // TipTap JSON
  coverUrl      String?
  subsiteId     String?
  
  // Автосохранение
  lastSavedAt   DateTime  @default(now())
  
  createdAt     DateTime  @default(now())
  updatedAt     DateTime  @updatedAt
  
  @@index([authorId, updatedAt])
}
```

---

## Индексы для производительности

### Составные индексы

```sql
-- Лента популярного
CREATE INDEX idx_articles_popular ON articles (
  status, 
  published_at DESC, 
  (view_count + comment_count * 10 + reaction_count * 5) DESC
) WHERE status = 'PUBLISHED';

-- Поиск по подсайту
CREATE INDEX idx_articles_subsite_date ON articles (
  subsite_id, 
  status, 
  published_at DESC
) WHERE status = 'PUBLISHED';

-- Комментарии к статье
CREATE INDEX idx_comments_article_tree ON comments (
  article_id, 
  parent_id NULLS FIRST, 
  created_at DESC
);

-- Уведомления пользователя
CREATE INDEX idx_notifications_user_unread ON notifications (
  user_id, 
  created_at DESC
) WHERE is_read = false;
```

---

## Миграции (примеры)

### Начальная миграция

```sql
-- CreateEnum UserRole
CREATE TYPE "UserRole" AS ENUM ('USER', 'AUTHOR', 'MODERATOR', 'EDITOR', 'ADMIN');

-- CreateTable User
CREATE TABLE "users" (
    "id" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "username" TEXT NOT NULL,
    "display_name" TEXT NOT NULL,
    "avatar_url" TEXT,
    "cover_url" TEXT,
    "bio" TEXT,
    "role" "UserRole" NOT NULL DEFAULT 'USER',
    "is_verified" BOOLEAN NOT NULL DEFAULT false,
    "is_plus_active" BOOLEAN NOT NULL DEFAULT false,
    "karma" INTEGER NOT NULL DEFAULT 0,
    "follower_count" INTEGER NOT NULL DEFAULT 0,
    "following_count" INTEGER NOT NULL DEFAULT 0,
    "article_count" INTEGER NOT NULL DEFAULT 0,
    "comment_count" INTEGER NOT NULL DEFAULT 0,
    "settings" JSONB NOT NULL DEFAULT '{}',
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,
    "last_active_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "users_pkey" PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "users_email_key" ON "users"("email");
CREATE UNIQUE INDEX "users_username_key" ON "users"("username");
CREATE INDEX "users_username_idx" ON "users"("username");
CREATE INDEX "users_created_at_idx" ON "users"("created_at");
```

---

## Seed данные

### Подсайты (категории)

```typescript
const subsites = [
  { slug: 'neural-networks', name: 'Нейросети', icon: '🤖' },
  { slug: 'generation', name: 'Генерация', icon: '🎨' },
  { slug: 'ai-business', name: 'AI для бизнеса', icon: '💼' },
  { slug: 'research', name: 'Исследования', icon: '🧪' },
  { slug: 'tools', name: 'Инструменты', icon: '🛠' },
  { slug: 'education', name: 'Обучение', icon: '📚' },
  { slug: 'investments', name: 'Инвестиции', icon: '💰' },
  { slug: 'opinions', name: 'Мнения', icon: '🗣' },
  { slug: 'regulation', name: 'Регулирование', icon: '⚖️' },
  { slug: 'future', name: 'Будущее', icon: '🔮' },
];
```

### Типы реакций

```typescript
const reactionTypes = [
  { emoji: '👍', label: 'Нравится', weight: 1 },
  { emoji: '❤️', label: 'Супер', weight: 2 },
  { emoji: '😂', label: 'Смешно', weight: 1 },
  { emoji: '🤔', label: 'Хмм', weight: 0 },
  { emoji: '😢', label: 'Грустно', weight: 0 },
  { emoji: '😡', label: 'Злюсь', weight: -1 },
  { emoji: '🔥', label: 'Огонь', weight: 2 },
  { emoji: '🎉', label: 'Праздник', weight: 1 },
];
```

### Достижения

```typescript
const achievements = [
  // Статьи
  { slug: 'first-article', name: 'Первый пост', category: 'ARTICLES', threshold: 1 },
  { slug: '10-articles', name: '10 постов', category: 'ARTICLES', threshold: 10 },
  { slug: '100-articles', name: '100 постов', category: 'ARTICLES', threshold: 100 },
  
  // Подписчики
  { slug: '10-followers', name: '10 подписчиков', category: 'FOLLOWERS', threshold: 10 },
  { slug: '100-followers', name: '100 подписчиков', category: 'FOLLOWERS', threshold: 100 },
  { slug: '1k-followers', name: '1K подписчиков', category: 'FOLLOWERS', threshold: 1000 },
  
  // Реакции
  { slug: '100-reactions', name: '100 лайков', category: 'REACTIONS', threshold: 100 },
  { slug: '1k-reactions', name: '1K лайков', category: 'REACTIONS', threshold: 1000 },
  
  // Особые
  { slug: 'early-user', name: 'Ранний пользователь', category: 'SPECIAL', threshold: 0 },
  { slug: 'plus-subscriber', name: 'Plus подписка', category: 'SPECIAL', threshold: 0 },
];
```

---

## Модерация и администрирование

> Полное описание системы ролей и модерации: [09-ADMIN-MODERATION.md](./09-ADMIN-MODERATION.md)

### ModerationAction (Действия модерации)

```prisma
model ModerationAction {
  id            String           @id @default(cuid())
  
  // Кто выполнил
  moderator     User             @relation("moderator", fields: [moderatorId], references: [id])
  moderatorId   String
  
  // Над кем/чем
  targetType    ModerationTarget // USER, ARTICLE, COMMENT
  targetId      String
  
  // Что сделано
  action        ModActionType
  reason        String?
  details       Json?
  
  // Сроки (для банов)
  expiresAt     DateTime?
  
  // Связь с жалобой
  reportId      String?
  report        Report?          @relation(fields: [reportId], references: [id])
  
  createdAt     DateTime         @default(now())
  
  @@index([moderatorId])
  @@index([targetType, targetId])
  @@index([createdAt])
}

enum ModerationTarget {
  USER
  ARTICLE
  COMMENT
}

enum ModActionType {
  // Контент
  APPROVE
  REJECT
  HIDE
  RESTORE
  MOVE
  PIN
  UNPIN
  MARK_NSFW
  EDIT
  
  // Пользователи
  WARNING
  MUTE
  COMMENT_BAN
  TEMP_BAN
  PERM_BAN
  SHADOW_BAN
  UNBAN
  
  // Роли
  ROLE_CHANGE
}
```

### Report (Жалобы)

```prisma
model Report {
  id            String        @id @default(cuid())
  
  // Кто подал
  reporter      User          @relation("reporter", fields: [reporterId], references: [id])
  reporterId    String
  
  // На что жалоба
  targetType    ReportTarget
  targetId      String
  
  // Причина
  reason        ReportReason
  description   String?
  
  // Статус
  status        ReportStatus  @default(PENDING)
  
  // Обработка
  moderator     User?         @relation("reportModerator", fields: [moderatorId], references: [id])
  moderatorId   String?
  resolution    String?
  
  // Связанные действия
  actions       ModerationAction[]
  
  createdAt     DateTime      @default(now())
  resolvedAt    DateTime?
  
  @@index([status])
  @@index([targetType, targetId])
  @@index([reporterId])
}

enum ReportTarget {
  ARTICLE
  COMMENT
  USER
  MESSAGE
}

enum ReportReason {
  SPAM
  OFFENSIVE
  FRAUD
  COPYRIGHT
  OUTDATED
  OFF_TOPIC
  NSFW_UNMARKED
  PERSONAL_DATA
  OTHER
}

enum ReportStatus {
  PENDING
  IN_REVIEW
  RESOLVED
  DISMISSED
  ESCALATED
}
```

### UserBan (Баны пользователей)

```prisma
model UserBan {
  id            String      @id @default(cuid())
  
  user          User        @relation(fields: [userId], references: [id])
  userId        String
  
  type          BanType
  reason        String
  
  bannedBy      User        @relation("bannedBy", fields: [bannedById], references: [id])
  bannedById    String
  
  startsAt      DateTime    @default(now())
  expiresAt     DateTime?   // null = постоянный
  
  isActive      Boolean     @default(true)
  
  // Апелляция
  appealText    String?
  appealStatus  AppealStatus?
  appealedAt    DateTime?
  
  createdAt     DateTime    @default(now())
  
  @@index([userId, isActive])
  @@index([expiresAt])
}

enum BanType {
  COMMENT_BAN   // Запрет комментировать
  POST_BAN      // Запрет публиковать
  TEMP_BAN      // Временный полный бан
  PERM_BAN      // Постоянный бан
  SHADOW_BAN    // Скрытый бан
  SUBSITE_BAN   // Бан в конкретном подсайте
}

enum AppealStatus {
  PENDING
  APPROVED
  REJECTED
}
```

### ArticleVersion (История версий)

```prisma
model ArticleVersion {
  id            String   @id @default(cuid())
  
  article       Article  @relation(fields: [articleId], references: [id])
  articleId     String
  
  version       Int
  title         String
  content       Json
  
  editedBy      User     @relation(fields: [editedById], references: [id])
  editedById    String
  changeReason  String?
  
  createdAt     DateTime @default(now())
  
  @@unique([articleId, version])
  @@index([articleId])
}
```

### SubsiteTeamMember (Команда подсайта)

```prisma
model SubsiteTeamMember {
  id            String         @id @default(cuid())
  
  subsite       Subsite        @relation(fields: [subsiteId], references: [id])
  subsiteId     String
  
  user          User           @relation(fields: [userId], references: [id])
  userId        String
  
  role          SubsiteRole
  
  addedBy       User           @relation("addedBy", fields: [addedById], references: [id])
  addedById     String
  
  createdAt     DateTime       @default(now())
  
  @@unique([subsiteId, userId])
  @@index([subsiteId])
  @@index([userId])
}

enum SubsiteRole {
  OWNER         // Владелец (полные права)
  EDITOR        // Редактор (редактирование контента)
  MODERATOR     // Модератор (модерация контента)
}
```

### SubsiteBan (Бан в подсайте)

```prisma
model SubsiteBan {
  id            String    @id @default(cuid())
  
  subsite       Subsite   @relation(fields: [subsiteId], references: [id])
  subsiteId     String
  
  user          User      @relation(fields: [userId], references: [id])
  userId        String
  
  reason        String
  
  bannedBy      User      @relation("subsiteBannedBy", fields: [bannedById], references: [id])
  bannedById    String
  
  expiresAt     DateTime?
  isActive      Boolean   @default(true)
  
  createdAt     DateTime  @default(now())
  
  @@unique([subsiteId, userId])
  @@index([subsiteId, isActive])
}
```

### Расширения существующих моделей

#### Subsite (дополнительные поля)

```prisma
// Добавить к существующей модели Subsite:

  // Тип подсайта
  type          SubsiteType @default(OFFICIAL)
  
  // Владелец (для сообществ)
  owner         User?     @relation("subsiteOwner", fields: [ownerId], references: [id])
  ownerId       String?
  
  // Команда
  team          SubsiteTeamMember[]
  
  // Модерация
  moderationMode ModerationMode @default(PREMOD_NEW)
  rules         String?   // Markdown с правилами
  
  // Ограничения
  allowedContentTypes ContentType[]
  whitelistedDomains  String[]
  
  // Флаги
  isArchived    Boolean   @default(false)
  requiresInvite Boolean  @default(false)
  
  // Баны
  bans          SubsiteBan[]

enum SubsiteType {
  OFFICIAL      // Официальный подсайт
  COMMUNITY     // Сообщество
  COMPANY       // Корпоративный блог
}

enum ModerationMode {
  FREE          // Без модерации
  PREMOD_NEW    // Премодерация для новых
  PREMOD_ALL    // Полная премодерация
  INVITE_ONLY   // По приглашению
}
```

#### Article (дополнительные поля)

```prisma
// Добавить к существующей модели Article:

  // Модерация
  moderationStatus ModerationStatus @default(PENDING)
  moderatedBy      User?            @relation("moderatedBy", fields: [moderatedById], references: [id])
  moderatedById    String?
  moderatedAt      DateTime?
  rejectionReason  String?
  
  // Версионирование
  version          Int              @default(1)
  versions         ArticleVersion[]
  
  // Флаги
  isNSFW          Boolean          @default(false)
  isSponsored     Boolean          @default(false)
  isLocked        Boolean          @default(false) // Комментарии закрыты
  
  // История модерации
  moderationHistory ModerationAction[]

enum ModerationStatus {
  PENDING       // Ожидает модерации
  APPROVED      // Одобрено
  REJECTED      // Отклонено
  AUTO_APPROVED // Авто-одобрено
}
```

