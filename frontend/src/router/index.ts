import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/useAuthStore'

// Safe in-app redirect target: absolute path only; reject protocol-relative
// and external URLs (SEC-05).
export function safeRedirectPath(target: unknown, fallback = '/') {
  if (typeof target !== 'string') return fallback
  if (!target.startsWith('/') || target.startsWith('//') || target.includes('://')) return fallback
  return target
}

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: () => import('@/layouts/UserLayout.vue'),
      children: [
        { path: '', name: 'home', component: () => import('@/pages/user/HomePage.vue') },
        { path: 'schedule', name: 'schedule', component: () => import('@/pages/user/SchedulePage.vue') },
        { path: 'matches/:id', name: 'match-detail', component: () => import('@/pages/user/MatchDetailPage.vue') },
        { path: 'teams', name: 'teams', component: () => import('@/pages/user/TeamsPage.vue') },
        { path: 'teams/:id', name: 'team-detail', component: () => import('@/pages/user/TeamDetailPage.vue') },
        { path: 'standings', name: 'standings', component: () => import('@/pages/user/StandingsPage.vue') },
        { path: 'profile', name: 'profile', component: () => import('@/pages/user/ProfilePage.vue'), meta: { requiresAuth: true } },
        { path: 'login', name: 'login', component: () => import('@/pages/user/LoginPage.vue') },
      ],
    },
    {
      path: '/admin/login',
      name: 'admin-login',
      component: () => import('@/pages/admin/AdminLoginPage.vue'),
    },
    {
      path: '/admin',
      component: () => import('@/layouts/AdminLayout.vue'),
      meta: { requiresAdmin: true },
      children: [
        { path: '', name: 'admin', component: () => import('@/pages/admin/AdminDashboard.vue') },
        { path: 'competitions', name: 'admin-competitions', component: () => import('@/pages/admin/AdminCompetitionsPage.vue') },
        { path: 'featured', name: 'admin-featured', component: () => import('@/pages/admin/AdminFeaturedPage.vue') },
        { path: 'teams', name: 'admin-teams', component: () => import('@/pages/admin/AdminTeamsPage.vue') },
        { path: 'matches', name: 'admin-matches', component: () => import('@/pages/admin/AdminMatchesPage.vue') },
        { path: 'standings', name: 'admin-standings', component: () => import('@/pages/admin/AdminStandingsPage.vue') },
        { path: 'sync', name: 'admin-sync', component: () => import('@/pages/admin/AdminSyncPage.vue') },
      ],
    },
  ],
  scrollBehavior() {
    return { top: 0 }
  },
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()

  if (to.meta.requiresAdmin) {
    if (!auth.isLoggedIn) {
      return { name: 'login', query: { redirect: to.fullPath } }
    }
    if (!auth.user) {
      await auth.fetchProfile() // reload role after a page refresh
    }
    if (!auth.isLoggedIn) {
      return { name: 'login', query: { redirect: to.fullPath } }
    }
    if (auth.user?.role !== 'admin') {
      return { name: 'home' } // normal users get bounced to the home page
    }
    return true
  }

  if (to.meta.requiresAuth) {
    if (!auth.isLoggedIn) {
      return { name: 'login', query: { redirect: to.fullPath } }
    }
    return true
  }

  if (to.name === 'login' && auth.isLoggedIn) {
    return { name: 'home' }
  }

  if (to.name === 'admin-login' && auth.isLoggedIn) {
    return { name: 'admin' }
  }

  return true
})

export default router
