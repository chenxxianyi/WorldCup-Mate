<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import MatchCard from '@/components/common/MatchCard.vue'
import { useAuthStore } from '@/stores/useAuthStore'
import { useFavoriteStore } from '@/stores/useFavoriteStore'

const router = useRouter()
const auth = useAuthStore()
const fav = useFavoriteStore()

const loading = ref(false)
const error = ref('')

async function load() {
  if (!auth.isLoggedIn) return
  loading.value = true
  error.value = ''
  try {
    await fav.fetchFavoriteMatches()
  } catch {
    error.value = '收藏比赛加载失败'
  } finally {
    loading.value = false
  }
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
        <h2>收藏比赛</h2>
        <p>{{ fav.favoriteMatchIds.length }} 场比赛</p>
      </div>
    </div>

    <div v-if="loading" class="card state-card">正在加载收藏比赛...</div>
    <div v-else-if="!auth.isLoggedIn" class="card state-card">
      <span>登录后可以查看收藏比赛</span>
      <button class="pill-btn primary" @click="router.push('/login')">去登录</button>
    </div>
    <div v-else-if="error" class="card state-card error">{{ error }}</div>
    <div v-else-if="!fav.favoriteMatches.length" class="card state-card">
      <span>还没有收藏比赛</span>
      <button class="pill-btn" @click="router.push('/schedule')">去赛程收藏</button>
    </div>

    <div v-else class="match-list">
      <MatchCard
        v-for="match in fav.favoriteMatches"
        :key="match.id"
        :match="match"
      />
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

.match-list {
  display: grid;
  gap: 12px;
}

.state-card {
  min-height: 120px;
  display: grid;
  place-items: center;
  gap: 12px;
  padding: 20px;
  color: var(--muted);
  text-align: center;
}

.state-card.error {
  color: var(--primary);
}
</style>
