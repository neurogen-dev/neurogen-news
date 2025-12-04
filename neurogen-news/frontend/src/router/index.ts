import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  // Home/Feed
  {
    path: '/',
    name: 'home',
    component: () => import('@/pages/HomePage.vue'),
    meta: { title: 'Главная' }
  },
  {
    path: '/popular',
    name: 'popular',
    component: () => import('@/pages/HomePage.vue'),
    meta: { title: 'Популярное' }
  },
  {
    path: '/new',
    name: 'new',
    component: () => import('@/pages/HomePage.vue'),
    meta: { title: 'Свежее' }
  },
  {
    path: '/my-feed',
    name: 'my-feed',
    component: () => import('@/pages/MyFeedPage.vue'),
    meta: { title: 'Моя лента', requiresAuth: true }
  },
  
  // Auth
  {
    path: '/login',
    name: 'login',
    component: () => import('@/pages/auth/LoginPage.vue'),
    meta: { title: 'Вход', guest: true }
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('@/pages/auth/RegisterPage.vue'),
    meta: { title: 'Регистрация', guest: true }
  },
  {
    path: '/forgot-password',
    name: 'forgot-password',
    component: () => import('@/pages/auth/ForgotPasswordPage.vue'),
    meta: { title: 'Восстановление пароля', guest: true }
  },
  
  // Categories
  {
    path: '/chatbots',
    name: 'chatbots',
    component: () => import('@/pages/CategoryPage.vue'),
    meta: { title: 'Чат-боты', category: 'chatbots' }
  },
  {
    path: '/images',
    name: 'images',
    component: () => import('@/pages/CategoryPage.vue'),
    meta: { title: 'Изображения', category: 'images' }
  },
  {
    path: '/video',
    name: 'video',
    component: () => import('@/pages/CategoryPage.vue'),
    meta: { title: 'Видео', category: 'video' }
  },
  {
    path: '/music',
    name: 'music',
    component: () => import('@/pages/CategoryPage.vue'),
    meta: { title: 'Музыка', category: 'music' }
  },
  {
    path: '/text',
    name: 'text',
    component: () => import('@/pages/CategoryPage.vue'),
    meta: { title: 'Текст', category: 'text' }
  },
  {
    path: '/code',
    name: 'code',
    component: () => import('@/pages/CategoryPage.vue'),
    meta: { title: 'Код', category: 'code' }
  },
  
  // Content types
  {
    path: '/guides',
    name: 'guides',
    component: () => import('@/pages/ContentTypePage.vue'),
    meta: { title: 'Гайды', contentType: 'guides' }
  },
  {
    path: '/reviews',
    name: 'reviews',
    component: () => import('@/pages/ContentTypePage.vue'),
    meta: { title: 'Обзоры', contentType: 'reviews' }
  },
  {
    path: '/prompts',
    name: 'prompts',
    component: () => import('@/pages/PromptsPage.vue'),
    meta: { title: 'Промпты' }
  },
  {
    path: '/news',
    name: 'news',
    component: () => import('@/pages/ContentTypePage.vue'),
    meta: { title: 'Новости', contentType: 'news' }
  },
  {
    path: '/qa',
    name: 'qa',
    component: () => import('@/pages/QAPage.vue'),
    meta: { title: 'Вопросы и ответы' }
  },
  
  // Tools
  {
    path: '/tools',
    name: 'tools',
    component: () => import('@/pages/ToolsPage.vue'),
    meta: { title: 'Каталог инструментов' }
  },
  {
    path: '/tools/:slug',
    name: 'tool',
    component: () => import('@/pages/ToolPage.vue'),
    meta: { title: 'Инструмент' }
  },
  
  // Start Here (для новичков)
  {
    path: '/start',
    name: 'start',
    component: () => import('@/pages/StartPage.vue'),
    meta: { title: 'Начать здесь' }
  },
  {
    path: '/start/:slug',
    name: 'start-article',
    component: () => import('@/pages/ArticlePage.vue'),
    meta: { title: 'Обучение' }
  },
  
  // Article
  {
    path: '/:category/:slug',
    name: 'article',
    component: () => import('@/pages/ArticlePage.vue'),
    meta: { title: 'Статья' }
  },
  
  // Editor
  {
    path: '/editor/new',
    name: 'editor-new',
    component: () => import('@/pages/EditorPage.vue'),
    meta: { title: 'Новая статья', requiresAuth: true }
  },
  {
    path: '/editor/:id',
    name: 'editor-edit',
    component: () => import('@/pages/EditorPage.vue'),
    meta: { title: 'Редактирование', requiresAuth: true }
  },
  
  // User
  {
    path: '/@:username',
    name: 'user-profile',
    component: () => import('@/pages/UserProfilePage.vue'),
    meta: { title: 'Профиль' }
  },
  {
    path: '/@:username/articles',
    name: 'user-articles',
    component: () => import('@/pages/UserArticlesPage.vue'),
    meta: { title: 'Статьи пользователя' }
  },
  {
    path: '/@:username/comments',
    name: 'user-comments',
    component: () => import('@/pages/UserCommentsPage.vue'),
    meta: { title: 'Комментарии пользователя' }
  },
  
  // Settings
  {
    path: '/settings',
    name: 'settings',
    component: () => import('@/pages/settings/SettingsPage.vue'),
    meta: { title: 'Настройки', requiresAuth: true },
    children: [
      {
        path: '',
        redirect: { name: 'settings-profile' }
      },
      {
        path: 'profile',
        name: 'settings-profile',
        component: () => import('@/pages/settings/ProfileSettings.vue'),
        meta: { title: 'Профиль' }
      },
      {
        path: 'account',
        name: 'settings-account',
        component: () => import('@/pages/settings/AccountSettings.vue'),
        meta: { title: 'Аккаунт' }
      },
      {
        path: 'notifications',
        name: 'settings-notifications',
        component: () => import('@/pages/settings/NotificationSettings.vue'),
        meta: { title: 'Уведомления' }
      },
      {
        path: 'feeds',
        name: 'settings-feeds',
        component: () => import('@/pages/settings/FeedSettings.vue'),
        meta: { title: 'Ленты' }
      },
    ]
  },
  
  // Bookmarks
  {
    path: '/bookmarks',
    name: 'bookmarks',
    component: () => import('@/pages/BookmarksPage.vue'),
    meta: { title: 'Закладки', requiresAuth: true }
  },
  
  // Drafts
  {
    path: '/drafts',
    name: 'drafts',
    component: () => import('@/pages/DraftsPage.vue'),
    meta: { title: 'Черновики', requiresAuth: true }
  },
  
  // Achievements
  {
    path: '/achievements',
    name: 'achievements',
    component: () => import('@/pages/AchievementsPage.vue'),
    meta: { title: 'Достижения', requiresAuth: true }
  },
  
  // Messages
  {
    path: '/messages',
    name: 'messages',
    component: () => import('@/pages/MessagesPage.vue'),
    meta: { title: 'Сообщения', requiresAuth: true }
  },
  
  // Search
  {
    path: '/search',
    name: 'search',
    component: () => import('@/pages/SearchPage.vue'),
    meta: { title: 'Поиск' }
  },
  
  // Static pages
  {
    path: '/about',
    name: 'about',
    component: () => import('@/pages/AboutPage.vue'),
    meta: { title: 'О проекте' }
  },
  {
    path: '/rules',
    name: 'rules',
    component: () => import('@/pages/RulesPage.vue'),
    meta: { title: 'Правила' }
  },
  {
    path: '/privacy',
    name: 'privacy',
    component: () => import('@/pages/PrivacyPage.vue'),
    meta: { title: 'Политика конфиденциальности' }
  },
  {
    path: '/terms',
    name: 'terms',
    component: () => import('@/pages/TermsPage.vue'),
    meta: { title: 'Условия использования' }
  },
  
  // Plus (premium)
  {
    path: '/plus',
    name: 'plus',
    component: () => import('@/pages/PlusPage.vue'),
    meta: { title: 'Neurogen Plus' }
  },
  
  // 404
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('@/pages/NotFoundPage.vue'),
    meta: { title: '404 - Страница не найдена' }
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, _from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    }
    if (to.hash) {
      return { el: to.hash, behavior: 'smooth' }
    }
    return { top: 0 }
  },
})

// Navigation guards
router.beforeEach(async (to, _from, next) => {
  // Update document title
  const title = to.meta.title as string
  document.title = title ? `${title} — Neurogen.News` : 'Neurogen.News'
  
  // Check auth requirements
  const { useAuthStore } = await import('@/stores/auth')
  const authStore = useAuthStore()
  
  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    return next({ name: 'login', query: { redirect: to.fullPath } })
  }
  
  if (to.meta.guest && authStore.isLoggedIn) {
    return next({ name: 'home' })
  }
  
  next()
})

export default router
