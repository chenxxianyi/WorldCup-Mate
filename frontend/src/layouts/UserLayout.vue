<script setup lang="ts">
import { onMounted, watch } from 'vue'
import CompetitionDialog from '@/components/theme/CompetitionDialog.vue'
import LeagueNavigation from '@/components/theme/LeagueNavigation.vue'
import LeagueTopbar from '@/components/theme/LeagueTopbar.vue'
import { useLeagueThemeStore } from '@/stores/useLeagueThemeStore'
import { useAuthStore } from '@/stores/useAuthStore'
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import { useReminderStore } from '@/stores/useReminderStore'
import { useNotificationStore } from '@/stores/useNotificationStore'

const theme = useLeagueThemeStore()
const auth = useAuthStore()
const favorites = useFavoriteStore()
const reminders = useReminderStore()
const notifications = useNotificationStore()

async function loadPersonalData() {
  if (!auth.isLoggedIn) return
  if (!auth.user) await auth.fetchProfile()
  if (!auth.isLoggedIn) return
  await Promise.all([
    favorites.fetchFavoriteTeams(),
    favorites.fetchFavoriteMatches(),
    reminders.fetchReminders(),
    notifications.fetchUnreadCount(),
  ])
}

onMounted(async () => {
  await theme.initialize()
  await loadPersonalData()
})

watch(() => auth.isLoggedIn, loadPersonalData)
</script>

<template>
  <div class="league-app">
    <svg class="svg-sprite" aria-hidden="true">
      <symbol id="i-ball" viewBox="0 0 24 24"><circle cx="12" cy="12" r="9" /><path d="m12 7 3 2.2-1.2 3.5h-3.6L9 9.2 12 7Zm-7.4 4 3.3 1.1 1.2 3.6-2.1 2.6M17 18.3l-2.1-2.6 1.2-3.6 3.3-1.1M9.1 15.7h5.8M7.9 12.1 9 9.2M16.1 12.1 15 9.2" /></symbol>
      <symbol id="i-home" viewBox="0 0 24 24"><path d="m3 11 9-7 9 7" /><path d="M5.5 9.5V20h13V9.5M9 20v-6h6v6" /></symbol>
      <symbol id="i-calendar" viewBox="0 0 24 24"><rect x="3" y="5" width="18" height="16" rx="2" /><path d="M7 3v4M17 3v4M3 10h18M7 14h2M12 14h2M17 14h.1M7 18h2M12 18h2" /></symbol>
      <symbol id="i-shield" viewBox="0 0 24 24"><path d="M12 3 20 6v5c0 5.3-3.3 8.3-8 10-4.7-1.7-8-4.7-8-10V6l8-3Z" /><path d="m9 12 2 2 4-5" /></symbol>
      <symbol id="i-table" viewBox="0 0 24 24"><path d="M5 4h14a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2ZM3 9h18M9 4v16M15 4v16" /></symbol>
      <symbol id="i-user" viewBox="0 0 24 24"><circle cx="12" cy="8" r="4" /><path d="M4.5 21c.8-4.2 3.3-6 7.5-6s6.7 1.8 7.5 6" /></symbol>
      <symbol id="i-sun" viewBox="0 0 24 24"><circle cx="12" cy="12" r="4" /><path d="M12 2v2M12 20v2M4.9 4.9l1.4 1.4M17.7 17.7l1.4 1.4M2 12h2M20 12h2M4.9 19.1l1.4-1.4M17.7 6.3l1.4-1.4" /></symbol>
      <symbol id="i-moon" viewBox="0 0 24 24"><path d="M20 15.6A8.5 8.5 0 0 1 8.4 4 8.5 8.5 0 1 0 20 15.6Z" /></symbol>
      <symbol id="i-chevron" viewBox="0 0 24 24"><path d="m7 10 5 5 5-5" /></symbol>
      <symbol id="i-back" viewBox="0 0 24 24"><path d="m15 18-6-6 6-6" /></symbol>
      <symbol id="i-bell" viewBox="0 0 24 24"><path d="M18 9a6 6 0 0 0-12 0c0 7-3 7-3 9h18c0-2-3-2-3-9ZM10 21h4" /></symbol>
      <symbol id="i-star" viewBox="0 0 24 24"><path d="m12 3 2.8 5.7 6.2.9-4.5 4.4 1.1 6.2-5.6-3-5.6 3 1.1-6.2L3 9.6l6.2-.9L12 3Z" /></symbol>
      <symbol id="i-search" viewBox="0 0 24 24"><circle cx="11" cy="11" r="7" /><path d="m20 20-4-4" /></symbol>
      <symbol id="i-clock" viewBox="0 0 24 24"><circle cx="12" cy="12" r="9" /><path d="M12 7v5l3 2" /></symbol>
      <symbol id="i-arrow" viewBox="0 0 24 24"><path d="M5 12h14M14 7l5 5-5 5" /></symbol>
      <symbol id="i-trophy" viewBox="0 0 24 24"><path d="M8 4h8v4c0 4-1.7 6-4 6s-4-2-4-6V4ZM9 20h6M12 14v6M8 6H4c0 4 1.4 6 5 6M16 6h4c0 4-1.4 6-5 6" /></symbol>
      <symbol id="i-mail" viewBox="0 0 24 24"><rect x="3" y="5" width="18" height="14" rx="2" /><path d="m4 7 8 6 8-6" /></symbol>
      <symbol id="i-lock" viewBox="0 0 24 24"><rect x="5" y="10" width="14" height="11" rx="2" /><path d="M8 10V7a4 4 0 0 1 8 0v3" /></symbol>
      <symbol id="i-check" viewBox="0 0 24 24"><path d="m5 12 4 4L19 6" /></symbol>
    </svg>

    <a class="skip-link" href="#main-content">跳到主要内容</a>
    <div class="ambient" aria-hidden="true" />
    <LeagueTopbar />

    <div class="app-frame">
      <LeagueNavigation />
      <main id="main-content" class="page-stage" tabindex="-1">
        <router-view v-slot="{ Component }">
          <transition name="page" mode="out-in"><component :is="Component" /></transition>
        </router-view>
      </main>
    </div>

    <CompetitionDialog />
    <div class="toast" :class="{ show: theme.toast }" role="status" aria-live="polite">{{ theme.toast }}</div>
  </div>
</template>

<style src="../styles/league-theme.css"></style>
