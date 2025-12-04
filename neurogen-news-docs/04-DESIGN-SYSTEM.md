# Дизайн-система Neurogen.News

## Философия дизайна

Neurogen.News использует современный, чистый дизайн с акцентом на читабельность контента. 
Основная цель — создать платформу, где контент находится в центре внимания, а UI ненавязчиво 
поддерживает пользовательский опыт.

**Референсы:** VC.ru (основной) + DTF.ru (дополнительный)

---

## Цветовая палитра

### Светлая тема (по умолчанию)

```css
:root {
  /* Основные цвета */
  --color-background: #FFFFFF;
  --color-background-secondary: #F9F9F9;
  --color-background-tertiary: #F3F3F3;
  --color-background-pink: #FFF0F0;  /* Акцентный баннер сверху */
  
  /* Текст */
  --color-text-primary: #1A1A1A;
  --color-text-secondary: #666666;
  --color-text-tertiary: #999999;
  --color-text-link: #346EE0;
  
  /* Акценты */
  --color-primary: #346EE0;           /* Синий - основной акцент */
  --color-primary-hover: #2856B3;
  --color-success: #34C759;
  --color-warning: #FF9500;
  --color-error: #FF3B30;
  
  /* Границы и разделители */
  --color-border: #E5E5E5;
  --color-border-light: #F0F0F0;
  
  /* Карточки и элементы */
  --color-card-background: #FFFFFF;
  --color-card-hover: #FAFAFA;
  
  /* Специальные */
  --color-verified-badge: #346EE0;    /* Синяя галочка */
  --color-plus-badge: #4FC3F7;        /* Голубой ромбик Plus */
  --color-editorial-badge: #346EE0;   /* Материал редакции */
}
```

### Тёмная тема

```css
:root[data-theme="dark"] {
  --color-background: #1A1A1A;
  --color-background-secondary: #242424;
  --color-background-tertiary: #2D2D2D;
  --color-background-pink: #2D1F1F;
  
  --color-text-primary: #FFFFFF;
  --color-text-secondary: #B3B3B3;
  --color-text-tertiary: #808080;
  --color-text-link: #5B8DEF;
  
  --color-border: #333333;
  --color-border-light: #2D2D2D;
  
  --color-card-background: #242424;
  --color-card-hover: #2D2D2D;
}
```

---

## Типографика

### Шрифты

```css
:root {
  /* Основной шрифт - системный стек */
  --font-family-base: -apple-system, BlinkMacSystemFont, 'Segoe UI', 
                      Roboto, Oxygen, Ubuntu, Cantarell, 
                      'Open Sans', 'Helvetica Neue', sans-serif;
  
  /* Моноширинный - для кода */
  --font-family-mono: 'SF Mono', 'Monaco', 'Inconsolata', 
                      'Roboto Mono', 'Source Code Pro', monospace;
}
```

### Размеры текста

| Элемент | Desktop | Mobile | Line Height |
|---------|---------|--------|-------------|
| H1 (Заголовок статьи) | 28px | 24px | 1.3 |
| H2 (Заголовок в статье) | 22px | 20px | 1.35 |
| H3 (Подзаголовок) | 18px | 17px | 1.4 |
| Body (Основной текст) | 17px | 16px | 1.6 |
| Small (Метаданные) | 14px | 13px | 1.4 |
| XSmall (Счётчики) | 12px | 12px | 1.3 |

### Веса шрифтов

```css
--font-weight-regular: 400;
--font-weight-medium: 500;
--font-weight-semibold: 600;
--font-weight-bold: 700;
```

---

## Отступы и сетка

### Spacing Scale

```css
--space-1: 4px;
--space-2: 8px;
--space-3: 12px;
--space-4: 16px;
--space-5: 20px;
--space-6: 24px;
--space-8: 32px;
--space-10: 40px;
--space-12: 48px;
--space-16: 64px;
```

### Grid System

```
Desktop (1200px+):
┌─────────────────────────────────────────────────────────────────┐
│ [Сайдбар 240px] │ [Контент max 640px] │ [Пустое пространство]  │
└─────────────────────────────────────────────────────────────────┘

Tablet (768px - 1199px):
┌─────────────────────────────────────────────────────────────────┐
│     [Сайдбар 200px]     │     [Контент 100%]                   │
└─────────────────────────────────────────────────────────────────┘

Mobile (<768px):
┌─────────────────────────────────────────────────────────────────┐
│                      [Контент 100%]                             │
│              [Нижняя навигация fixed]                           │
└─────────────────────────────────────────────────────────────────┘
```

### Максимальные ширины

```css
--max-width-content: 640px;      /* Основной контент */
--max-width-article: 720px;      /* Статья с медиа */
--max-width-sidebar: 240px;      /* Левый сайдбар */
--max-width-container: 1200px;   /* Общий контейнер */
```

---

## Компоненты

### Кнопки

#### Primary Button
```css
.button-primary {
  background: var(--color-primary);
  color: white;
  padding: 10px 20px;
  border-radius: 8px;
  font-weight: 500;
  font-size: 14px;
  transition: background 0.2s;
}

.button-primary:hover {
  background: var(--color-primary-hover);
}
```

#### Secondary Button (Ghost)
```css
.button-secondary {
  background: transparent;
  color: var(--color-text-primary);
  padding: 10px 20px;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  font-weight: 500;
}

.button-secondary:hover {
  background: var(--color-background-secondary);
}
```

#### Subscribe Button
```css
.button-subscribe {
  background: var(--color-primary);
  color: white;
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
}

.button-subscribe[data-subscribed="true"] {
  background: transparent;
  color: var(--color-text-secondary);
  border: 1px solid var(--color-border);
}
```

### Карточки

#### Article Card в ленте

```
┌─────────────────────────────────────────────────────────────────┐
│  [32px Avatar]  Author Name    • Category • 2h    [Subscribe]   │
│                                                                 │
│  ЗАГОЛОВОК СТАТЬИ [✓ badge]                                     │
│                                                                 │
│  Лид-текст статьи с ограничением в 2-3 строки...               │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    Cover Image 16:9                      │   │
│  │                      (max 640px)                         │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                            [Показать полностью] │
│                                                                 │
│  [😊 42] [❤️ 15] [🔥 8] [+]                                     │
│                                                                 │
│  [💬 128]  [🔖 15]  [↗️ Share]  [···]    [👁 12.5K]              │
│                                                                 │
│  "Лучший комментарий из обсуждения..."                         │
└─────────────────────────────────────────────────────────────────┘
```

**Стили:**
```css
.article-card {
  background: var(--color-card-background);
  border-radius: 12px;
  padding: var(--space-5);
  margin-bottom: var(--space-4);
  transition: background 0.2s;
}

.article-card:hover {
  background: var(--color-card-hover);
}

.article-card__header {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-bottom: var(--space-4);
}

.article-card__avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
}

.article-card__title {
  font-size: 20px;
  font-weight: 600;
  line-height: 1.3;
  margin-bottom: var(--space-3);
}

.article-card__cover {
  border-radius: 12px;
  overflow: hidden;
  aspect-ratio: 16 / 9;
  margin: var(--space-4) 0;
}
```

### Аватары

| Размер | Использование |
|--------|---------------|
| 24px | Компактный (комментарии) |
| 32px | Карточка статьи в ленте |
| 40px | Шапка статьи |
| 48px | Профиль в сайдбаре |
| 96px | Страница профиля |

```css
.avatar {
  border-radius: 50%;
  object-fit: cover;
}

.avatar--24 { width: 24px; height: 24px; }
.avatar--32 { width: 32px; height: 32px; }
.avatar--40 { width: 40px; height: 40px; }
.avatar--48 { width: 48px; height: 48px; }
.avatar--96 { width: 96px; height: 96px; }
```

### Бейджи

#### Верифицированный аккаунт (галочка)
```css
.badge-verified {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  background: var(--color-verified-badge);
  border-radius: 50%;
  color: white;
}

.badge-verified svg {
  width: 10px;
  height: 10px;
}
```

#### Plus подписка (ромбик)
```css
.badge-plus {
  display: inline-flex;
  width: 14px;
  height: 14px;
  background: var(--color-plus-badge);
  transform: rotate(45deg);
  border-radius: 2px;
}
```

#### Материал редакции
```css
.badge-editorial {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  color: var(--color-editorial-badge);
  font-size: 13px;
}
```

### Реакции

```css
.reaction {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  border-radius: 16px;
  background: transparent;
  transition: background 0.2s;
  cursor: pointer;
}

.reaction:hover {
  background: var(--color-background-secondary);
}

.reaction--active {
  background: rgba(52, 110, 224, 0.1);
}

.reaction__emoji {
  font-size: 16px;
}

.reaction__count {
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text-secondary);
}
```

**Доступные реакции:**
| Эмодзи | Название | Код |
|--------|----------|-----|
| 👍 | Нравится | like |
| ❤️ | Супер | love |
| 😂 | Смешно | haha |
| 🤔 | Хмм | think |
| 😢 | Грустно | sad |
| 😡 | Злюсь | angry |
| 🔥 | Огонь | fire |
| 🎉 | Праздник | party |

---

## Навигация

### Левый сайдбар

```
┌───────────────────────────────┐
│  [Logo]                        │
│                               │
│  ○ Популярное                 │
│  ○ Свежее                     │
│  ● Моя лента  ← активный      │
│  ○ Сообщения                  │
│  ○ Рейтинг                    │
│  ○ Курсы                      │
│                               │
│  ─────────────────────────────│
│  Темы                         │
│  🤖 Нейросети                  │
│  🎨 Генерация                  │
│  💼 AI бизнес                  │
│  📚 Обучение                   │
│  ...                          │
│  ∨ Показать все               │
│                               │
│  ─────────────────────────────│
│  neurogen.news                │
│  ℹ О проекте                  │
│  📋 Правила                   │
│  📢 Реклама                   │
│  📱 Приложения                │
└───────────────────────────────┘
```

**Стили меню:**
```css
.sidebar-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-4);
  border-radius: 8px;
  color: var(--color-text-primary);
  transition: background 0.2s;
}

.sidebar-item:hover {
  background: var(--color-background-secondary);
}

.sidebar-item--active {
  background: var(--color-background-secondary);
  font-weight: 500;
}

.sidebar-item__icon {
  width: 24px;
  height: 24px;
  opacity: 0.7;
}

.sidebar-item--active .sidebar-item__icon {
  opacity: 1;
  color: var(--color-primary);
}
```

### Шапка (Header)

```
Desktop:
┌─────────────────────────────────────────────────────────────────┐
│  [Logo]              [Search] [🔔 Notifications] [✏️] [Avatar▼]│
└─────────────────────────────────────────────────────────────────┘

Mobile:
┌─────────────────────────────────────────────────────────────────┐
│  [☰ Menu]            [Logo]                    [🔍] [Avatar]   │
└─────────────────────────────────────────────────────────────────┘
```

### Меню пользователя (Dropdown)

```
┌───────────────────────────────────┐
│  Мой профиль                      │
│  ┌─────────────────────────────┐ │
│  │ [Avatar] ИМЯ 💎             │ │
│  │          @username          │ │
│  └─────────────────────────────┘ │
│  ─────────────────────────────── │
│  ✏️ Черновики                    │
│  🔖 Закладки                     │
│  🏆 Ачивки                       │
│  💰 Донаты          [Подключить] │
│  ⚙️ Настройки              [◐]  │ ← переключатель темы
│  ─────────────────────────────── │
│  💎 Подписка Plus                │
│  ─────────────────────────────── │
│  ↪️ Выйти                        │
└───────────────────────────────────┘
```

---

## Панель уведомлений

```
┌───────────────────────────────────┐
│  Уведомления                [···] │
│  ─────────────────────────────── │
│  ● [Avatar] Иван ответил на...   │
│    "Текст ответа..."             │
│    2 минуты назад                │
│  ─────────────────────────────── │
│  ○ [Avatar] Мария подписалась    │
│    1 час назад                   │
│  ─────────────────────────────── │
│  ○ [Avatar] Реакция на статью    │
│    3 часа назад                  │
└───────────────────────────────────┘
```

**Стили:**
```css
.notification-item {
  display: flex;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-4);
  transition: background 0.2s;
}

.notification-item:hover {
  background: var(--color-background-secondary);
}

.notification-item--unread::before {
  content: '';
  width: 8px;
  height: 8px;
  background: var(--color-primary);
  border-radius: 50%;
  flex-shrink: 0;
}
```

---

## Иконки

### Набор иконок

Используем Lucide Icons (или Phosphor Icons) для единообразия.

**Основные иконки:**
| Иконка | Название | Использование |
|--------|----------|---------------|
| Home | Главная | Навигация |
| Clock | Свежее | Навигация |
| Newspaper | Моя лента | Навигация |
| MessageSquare | Сообщения | Навигация |
| Trophy | Рейтинг | Навигация |
| GraduationCap | Курсы | Навигация |
| Search | Поиск | Header |
| Bell | Уведомления | Header |
| PenSquare | Написать | Header |
| Bookmark | Закладки | Действия |
| Share | Поделиться | Действия |
| MoreHorizontal | Меню | Действия |
| Eye | Просмотры | Метрики |
| MessageCircle | Комментарии | Метрики |
| Settings | Настройки | Профиль |

---

## Анимации

### Transitions

```css
:root {
  --transition-fast: 150ms ease;
  --transition-base: 200ms ease;
  --transition-slow: 300ms ease;
}
```

### Hover Effects

```css
/* Подъём карточки */
.card-lift {
  transition: transform var(--transition-base), 
              box-shadow var(--transition-base);
}

.card-lift:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}
```

### Loading States

```css
/* Skeleton */
.skeleton {
  background: linear-gradient(
    90deg,
    var(--color-background-secondary) 25%,
    var(--color-background-tertiary) 50%,
    var(--color-background-secondary) 75%
  );
  background-size: 200% 100%;
  animation: skeleton-loading 1.5s infinite;
  border-radius: 4px;
}

@keyframes skeleton-loading {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}
```

---

## Формы

### Input Field

```css
.input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  font-size: 16px;
  transition: border-color var(--transition-base);
}

.input:focus {
  outline: none;
  border-color: var(--color-primary);
}

.input::placeholder {
  color: var(--color-text-tertiary);
}
```

### Comment Input

```css
.comment-input {
  display: flex;
  align-items: flex-start;
  gap: var(--space-3);
  padding: var(--space-4);
  background: var(--color-background-secondary);
  border-radius: 12px;
}

.comment-input__textarea {
  flex: 1;
  min-height: 40px;
  max-height: 200px;
  resize: none;
  border: none;
  background: transparent;
  font-size: 15px;
}

.comment-input__actions {
  display: flex;
  gap: var(--space-2);
}

.comment-input__action {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  color: var(--color-text-tertiary);
}
```

---

## Адаптивность

### Breakpoints

```css
/* Mobile first approach */
--breakpoint-sm: 576px;
--breakpoint-md: 768px;
--breakpoint-lg: 992px;
--breakpoint-xl: 1200px;
--breakpoint-xxl: 1400px;
```

### Mobile Navigation (Bottom Bar)

```
┌─────────────────────────────────────────────────────────────────┐
│  [🏠 Главная] [📰 Лента] [✏️ Пост] [🔍 Поиск] [👤 Профиль]      │
└─────────────────────────────────────────────────────────────────┘
```

```css
.mobile-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 56px;
  background: var(--color-background);
  border-top: 1px solid var(--color-border);
  display: flex;
  justify-content: space-around;
  align-items: center;
  padding-bottom: env(safe-area-inset-bottom);
}

@media (min-width: 768px) {
  .mobile-nav {
    display: none;
  }
}
```

---

## Тёмная тема

### Переключение темы

```javascript
// Автоматическое определение + ручное переключение
const getPreferredTheme = () => {
  const stored = localStorage.getItem('theme');
  if (stored) return stored;
  return window.matchMedia('(prefers-color-scheme: dark)').matches 
    ? 'dark' 
    : 'light';
};

const setTheme = (theme) => {
  document.documentElement.setAttribute('data-theme', theme);
  localStorage.setItem('theme', theme);
};
```

### Изображения в тёмной теме

```css
/* Уменьшаем яркость изображений в тёмной теме */
[data-theme="dark"] .article-cover img,
[data-theme="dark"] .article-content img {
  filter: brightness(0.9);
}
```


