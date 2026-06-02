<script setup lang="ts">
import { computed, onMounted } from 'vue'
import StatCard from '@/components/common/StatCard.vue'
import { useSettingStore } from '@/stores/useSettingStore'
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import { useReminderStore } from '@/stores/useReminderStore'
import { useTeamStore } from '@/stores/useTeamStore'
import { useAuthStore } from '@/stores/useAuthStore'

const settings = useSettingStore()
const fav = useFavoriteStore()
const reminder = useReminderStore()
const teamStore = useTeamStore()
const auth = useAuthStore()

const followedTeamNames = computed(() =>
  teamStore.teams
    .filter((t) => fav.isTeamFollowed(t.id))
    .map((t) => t.name)
    .join('、')
)

onMounted(() => {
  teamStore.fetchTeams()
  if (auth.isLoggedIn) {
    fav.fetchFavoriteTeams()
    fav.fetchFavoriteMatches()
    reminder.fetchReminders()
  }
})
</script>

<template>
  <div>
    <article class="card profile-card">
      <div class="avatar">{{ auth.user?.avatar || 'U' }}</div>
      <div>
        <h2>{{ auth.user?.nickname || auth.user?.username || '未登录' }}</h2>
        <p>{{ settings.timezone }} · 已关注 {{ fav.followedTeamIds.length }} 支球队</p>
      </div>
    </article>

    <section class="section">
      <div class="stats-row">
        <StatCard :value="fav.followedTeamIds.length" label="关注球队" />
        <StatCard :value="fav.favoriteMatchIds.length" label="收藏比赛" />
        <StatCard :value="reminder.count" label="比赛提醒" />
      </div>
    </section>

    <section class="section">
      <div class="card settings-list">
        <div class="setting-item"><b>我的关注</b><span>{{ followedTeamNames }}</span></div>
        <div class="setting-item"><b>我的提醒</b><span>{{ reminder.count }} 个待发送</span></div>
        <div class="setting-item"><b>时区</b><span>{{ settings.timezone }}</span></div>
        <div class="setting-item">
          <b>深色模式</b>
          <button class="pill-btn" @click="settings.toggleTheme">切换</button>
        </div>
        <div class="setting-item"><b>语言</b><span>简体中文</span></div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.profile-card {
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 13px;
}

.avatar {
  width: 54px;
  height: 54px;
  display: grid;
  place-items: center;
  border-radius: 18px;
  color: #fff;
  font-size: 22px;
  font-weight: 800;
  background: linear-gradient(145deg, var(--primary), var(--secondary));
}

.profile-card h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 750;
}

.profile-card p {
  margin: 5px 0 0;
  color: var(--muted);
  font-size: 13px;
}

.section {
  margin-top: 18px;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
}

.settings-list {
  overflow: hidden;
}

.setting-item {
  min-height: 58px;
  padding: 0 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--line);
}

.setting-item:last-child {
  border-bottom: 0;
}

.setting-item span {
  color: var(--muted);
  font-size: 13px;
}
</style>
