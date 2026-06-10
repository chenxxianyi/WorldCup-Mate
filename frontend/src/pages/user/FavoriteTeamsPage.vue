<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import TeamFlag from '@/components/common/TeamFlag.vue'
import { useAuthStore } from '@/stores/useAuthStore'
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import { useTeamStore } from '@/stores/useTeamStore'

const router = useRouter()
const auth = useAuthStore()
const fav = useFavoriteStore()
const teamStore = useTeamStore()

const loading = ref(false)
const error = ref('')

const followedTeams = computed(() =>
  teamStore.teams.filter((team) => fav.isTeamFollowed(team.id)),
)

async function load() {
  if (!auth.isLoggedIn) return
  loading.value = true
  error.value = ''
  try {
    await Promise.all([
      fav.fetchFavoriteTeams(),
      teamStore.fetchTeams({ page_size: 100 }),
    ])
  } catch {
    error.value = '关注球队加载失败'
  } finally {
    loading.value = false
  }
}

function goTeam(teamId: number) {
  router.push(`/teams/${teamId}`)
}

onMounted(load)
</script>

<template>
  <div class="detail-page">
    <div class="detail-head">
      <button class="icon-back" type="button" title="返回" @click="router.back()">
        <span class="material-symbols-outlined">arrow_back</span>
      </button>
      <div>
        <h2>关注球队</h2>
        <p>{{ followedTeams.length }} 支球队</p>
      </div>
    </div>

    <div v-if="loading" class="empty-state">正在加载关注球队...</div>
    <div v-else-if="!auth.isLoggedIn" class="empty-state">
      <span>登录后可以查看关注球队</span>
      <button class="pill-btn primary" @click="router.push('/login')">去登录</button>
    </div>
    <div v-else-if="error" class="empty-state state-error">{{ error }}</div>
    <div v-else-if="!followedTeams.length" class="empty-state">
      <span>还没有关注球队</span>
      <button class="pill-btn" @click="router.push('/teams')">去关注球队</button>
    </div>

    <div v-else class="team-list">
      <article
        v-for="team in followedTeams"
        :key="team.id"
        class="surface-row team-row"
        tabindex="0"
        @click="goTeam(team.id)"
        @keydown.enter="goTeam(team.id)"
      >
        <TeamFlag :value="team.flag" :alt="team.name" :fallback="team.code" size="lg" />
        <div class="team-copy">
          <h3>{{ team.name }}</h3>
          <p>{{ team.name_en || team.code }} · {{ team.group_name || team.continent }}</p>
        </div>
        <span class="tag green">已关注</span>
        <span class="material-symbols-outlined arrow">chevron_right</span>
      </article>
    </div>
  </div>
</template>

<style scoped>
.detail-page {
  display: grid;
  gap: 14px;
}

.detail-head {
  display: flex;
  align-items: center;
  gap: 10px;
}

.icon-back {
  width: 36px;
  height: 36px;
  display: grid;
  place-items: center;
  border: 1px solid var(--line);
  border-radius: 999px;
  color: var(--text);
  background: var(--card);
}

.icon-back .material-symbols-outlined {
  font-size: 20px;
}

.detail-head h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 800;
}

.detail-head p {
  margin: 4px 0 0;
  color: var(--muted);
  font-size: 13px;
}

.team-list {
  overflow: hidden;
  border-top: 1px solid color-mix(in srgb, var(--line) 78%, transparent);
}

.team-row {
  min-width: 0;
  display: grid;
  grid-template-columns: auto minmax(0, 1fr) auto auto;
  align-items: center;
  gap: 10px;
  padding: 13px 2px;
  cursor: pointer;
  transition: border-color 160ms ease-out, background 160ms ease-out, transform 160ms ease-out;
}

.team-row:hover {
  background: color-mix(in srgb, var(--primary) 4%, transparent);
}

.team-row:active {
  transform: scale(0.99);
}

.team-row:focus-visible {
  outline: 2px solid var(--primary);
  outline-offset: 2px;
}

.team-copy {
  min-width: 0;
}

.team-copy h3 {
  overflow: hidden;
  margin: 0;
  font-size: 16px;
  font-weight: 800;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.team-copy p {
  overflow: hidden;
  margin: 5px 0 0;
  color: var(--muted);
  font-size: 12px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.arrow {
  color: var(--weak);
  font-size: 20px;
}

.state-error {
  color: var(--primary);
}

@media (max-width: 390px) {
  .team-row {
    grid-template-columns: auto minmax(0, 1fr) auto;
  }

  .team-row .tag {
    display: none;
  }
}
</style>
