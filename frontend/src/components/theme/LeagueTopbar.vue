<script setup lang="ts">
import { useRouter } from 'vue-router'
import { type CompetitionCode } from '@/data/leagueTheme'
import { useLeagueThemeStore } from '@/stores/useLeagueThemeStore'
import { useAuthStore } from '@/stores/useAuthStore'
import { useNotificationStore } from '@/stores/useNotificationStore'
import ThemeIcon from './ThemeIcon.vue'

const router = useRouter()
const theme = useLeagueThemeStore()
const auth = useAuthStore()
const notifications = useNotificationStore()
</script>

<template>
  <header class="topbar">
    <div class="topbar-inner">
      <button
        class="brand"
        type="button"
        aria-label="返回首页"
        @click="router.push('/')"
      >
        <span class="brand-mark"><ThemeIcon name="ball" /></span>
        <span class="brand-copy">
          <strong>WorldCup Mate</strong>
          <small>ONE APP · SIX MATCHDAY WORLDS</small>
        </span>
      </button>

      <nav
        class="desktop-leagues"
        aria-label="赛事切换"
      >
        <button
          v-for="code in theme.competitionCodes"
          :key="code"
          class="league-tab"
          :class="{ active: code === theme.currentCode }"
          type="button"
          :aria-pressed="code === theme.currentCode"
          @click="theme.setCompetition(code as CompetitionCode)"
        >
          <span class="league-tab-mark">{{ theme.themeFor(code).mark }}</span>
          <span class="league-tab-name">{{ theme.themeFor(code).name }}</span>
        </button>
      </nav>

      <div class="top-actions">
        <button
          class="icon-button"
          type="button"
          :aria-label="theme.settings.theme === 'dark' ? '切换浅色模式' : '切换深色模式'"
          @click="theme.toggleTheme"
        >
          <ThemeIcon :name="theme.settings.theme === 'dark' ? 'sun' : 'moon'" />
        </button>
        <button
          class="avatar-button"
          type="button"
          aria-label="打开个人中心"
          @click="router.push(auth.isLoggedIn ? '/profile' : '/login')"
        >
          <img
            v-if="auth.user?.avatar?.startsWith('/')"
            :src="auth.user.avatar"
            alt=""
          >
          <span v-else>{{ auth.user?.avatar || 'M' }}</span>
          <i
            v-if="notifications.unreadCount"
            class="avatar-notice"
          >{{ notifications.unreadCount > 9 ? '9+' : notifications.unreadCount }}</i>
        </button>
      </div>

      <button
        class="mobile-competition-trigger"
        type="button"
        aria-haspopup="dialog"
        aria-controls="competition-dialog"
        @click="theme.competitionDialogOpen = true"
      >
        <span class="mobile-trigger-copy">
          <span class="league-tab-mark">{{ theme.current.mark }}</span>
          <strong>{{ theme.current.name }} · {{ theme.current.en }}</strong>
          <span>{{ theme.current.season }}</span>
        </span>
        <ThemeIcon name="chevron" />
      </button>
    </div>
  </header>
</template>
