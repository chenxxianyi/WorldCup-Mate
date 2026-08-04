<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { statusLabel } from '@/data/leagueTheme'
import { formatMatchday, type ThemeMatch } from '@/data/themeAdapters'
import { useLeagueThemeStore } from '@/stores/useLeagueThemeStore'
import { useAuthStore } from '@/stores/useAuthStore'
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import { useReminderStore } from '@/stores/useReminderStore'
import TeamBadge from './TeamBadge.vue'
import ThemeIcon from './ThemeIcon.vue'

const props = defineProps<{ match: ThemeMatch; compact?: boolean }>()
const route = useRoute()
const router = useRouter()
const theme = useLeagueThemeStore()
const auth = useAuthStore()
const favorites = useFavoriteStore()
const reminders = useReminderStore()
const pendingAction = ref<'favorite' | 'reminder' | ''>('')

function openMatch() {
  router.push(`/matches/${props.match.id}`)
}

function requireLogin() {
  if (auth.isLoggedIn) return true
  theme.showToast('登录后才能使用收藏和提醒')
  router.push({ path: '/login', query: { redirect: route.fullPath } })
  return false
}

async function toggleFavorite() {
  if (!requireLogin() || pendingAction.value) return
  pendingAction.value = 'favorite'
  const wasFavorite = favorites.isMatchFavorite(props.match.id)
  try {
    const saved = await favorites.toggleMatchFavorite(props.match.id)
    theme.showToast(saved ? (wasFavorite ? '已取消收藏' : '比赛已加入收藏') : '收藏操作失败，请稍后重试')
  } finally {
    pendingAction.value = ''
  }
}

async function toggleReminder() {
  if (!requireLogin() || pendingAction.value) return
  pendingAction.value = 'reminder'
  const hadReminder = reminders.hasReminder(props.match.id)
  try {
    const saved = await reminders.toggleReminder(props.match.id, 15)
    theme.showToast(saved ? (hadReminder ? '开球提醒已关闭' : '开球提醒已开启') : '提醒操作失败，请稍后重试')
  } finally {
    pendingAction.value = ''
  }
}
</script>

<template>
  <article class="card match-card clickable-card" :class="{ featured: match.featured }" @click="openMatch">
    <div class="match-card-top">
      <span class="label">{{ formatMatchday(match.source, theme.current.stage) }}</span>
      <span class="status-pill" :class="match.status">
        <i v-if="match.status === 'live'" class="live-dot" />
        {{ statusLabel(match.status) }}
      </span>
    </div>
    <div class="match-teams">
      <div class="match-team">
        <TeamBadge :team="match.home" size="small" />
        <span class="match-team-copy"><strong>{{ match.home[0] }}</strong><small>{{ match.home[1] }}</small></span>
      </div>
      <button class="match-score text-button" type="button" :aria-label="`查看 ${match.home[0]} 对 ${match.away[0]} 比赛详情`" @click.stop="openMatch">
        {{ match.score }}
      </button>
      <div class="match-team away">
        <span class="match-team-copy"><strong>{{ match.away[0] }}</strong><small>{{ match.away[1] }}</small></span>
        <TeamBadge :team="match.away" size="small" />
      </div>
    </div>
    <div class="match-card-bottom">
      <span class="match-meta">
        <strong>{{ match.date }} · {{ match.time }}</strong>
        <span>{{ compact ? theme.current.name : match.venue }}</span>
      </span>
      <span class="card-actions">
        <button
          class="icon-action"
          :class="{ active: reminders.hasReminder(match.id) }"
          type="button"
          aria-label="设置比赛提醒"
          @click.stop="toggleReminder"
        ><ThemeIcon name="bell" /></button>
        <button
          class="icon-action"
          :class="{ active: favorites.isMatchFavorite(match.id) }"
          type="button"
          aria-label="收藏比赛"
          @click.stop="toggleFavorite"
        ><ThemeIcon name="star" /></button>
      </span>
    </div>
  </article>
</template>
