import { createRouter, createWebHistory } from 'vue-router'

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
        { path: 'ai', name: 'ai-home', component: () => import('@/pages/user/AIHomePage.vue') },
        { path: 'ai/chat', name: 'ai-chat', component: () => import('@/pages/user/AIChatPage.vue') },
        { path: 'ai/match/:id', name: 'ai-match-insight', component: () => import('@/pages/user/AIMatchInsightPage.vue') },
        { path: 'ai/share-copy', name: 'ai-share-copy', component: () => import('@/pages/user/AIShareCopyPage.vue') },
        { path: 'profile', name: 'profile', component: () => import('@/pages/user/ProfilePage.vue') },
        { path: 'login', name: 'login', component: () => import('@/pages/user/LoginPage.vue') },
      ],
    },
    {
      path: '/admin',
      component: () => import('@/layouts/AdminLayout.vue'),
      children: [
        { path: '', name: 'admin', component: () => import('@/pages/admin/AdminDashboard.vue') },
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

export default router
