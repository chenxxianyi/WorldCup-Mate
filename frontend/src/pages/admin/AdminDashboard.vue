<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import StatCard from '@/components/common/StatCard.vue'
import SyncStatusBadge from '@/components/common/SyncStatusBadge.vue'
import { apiAdminDashboard } from '@/api/admin'
import { useMatchStore } from '@/stores/useMatchStore'

const router = useRouter()
const matchStore = useMatchStore()

const dashboard = ref({
  total_matches: 0,
  total_groups: 0,
  total_teams: 0,
  total_users: 0,
  total_competitions: 0,
  total_reminders: 0,
})
const loading = ref(true)
const error = ref('')

async function loadDashboard() {
  loading.value = true
  error.value = ''
  try {
    dashboard.value = await apiAdminDashboard()
  } catch (err: any) {
    error.value = err?.message || '看板数据加载失败'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadDashboard()
  matchStore.fetchMatches({ page: 1, page_size: 20 })
})
</script>

<template>
  <div class="admin-shell">
    <article class="card admin-hero">
      <div
        class="section-head"
        style="margin-bottom: 0"
      >
        <div>
          <h2>后台数据看板</h2>
          <span>WorldCup Mate Admin</span>
        </div>
        <button
          class="pill-btn primary"
          @click="router.push('/admin/matches')"
        >
          新增比赛
        </button>
      </div>

      <div
        v-if="loading"
        class="admin-hint"
        aria-busy="true"
      >
        看板数据加载中...
      </div>
      <div
        v-else-if="error"
        class="admin-hint error"
        role="alert"
      >
        {{ error }}
      </div>
      <div
        v-else
        class="admin-kpis"
      >
        <StatCard
          :value="dashboard.total_matches"
          label="总比赛数"
        />
        <StatCard
          :value="dashboard.total_groups"
          label="小组数量"
        />
        <StatCard
          :value="dashboard.total_teams"
          label="球队数量"
        />
        <StatCard
          :value="dashboard.total_competitions"
          label="赛事数量"
        />
        <StatCard
          :value="dashboard.total_users"
          label="用户数量"
        />
        <StatCard
          :value="dashboard.total_reminders"
          label="提醒总数"
        />
      </div>
      <SyncStatusBadge mode="card" />
    </article>

    <div class="card table-card table-scroll">
      <table class="admin-table">
        <thead>
          <tr>
            <th>比赛</th><th>阶段</th><th>时间</th><th>状态</th><th>推荐</th><th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="m in matchStore.matches"
            :key="m.id"
          >
            <td>{{ m.home_team_name }} vs {{ m.away_team_name }}</td>
            <td>{{ m.group_name || (m.matchday ? `第 ${m.matchday} 轮` : '-') }}</td>
            <td>{{ m.local_kickoff_time }}</td>
            <td>
              <span
                v-if="m.status === 'live'"
                class="tag live"
              >
                <i
                  class="live-dot"
                  style="background: #fff"
                /> LIVE
              </span>
              <span
                v-else-if="m.status === 'finished'"
                class="tag green"
              >已结束</span>
              <span
                v-else
                class="tag"
              >未开始</span>
            </td>
            <td>
              <span
                v-if="m.is_featured"
                class="tag gold"
              >热门</span>
              <span
                v-else-if="m.importance_level >= 2"
                class="tag blue"
              >焦点</span>
              <span
                v-else
                class="tag"
              >-</span>
            </td>
            <td>
              <button
                class="pill-btn"
                @click="router.push('/admin/matches')"
              >
                编辑
              </button>
            </td>
          </tr>
          <tr v-if="!matchStore.matches.length">
            <td
              colspan="6"
              class="empty-row"
            >
              暂无比赛数据
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
.admin-shell {
  display: grid;
  gap: 14px;
}

.admin-hero {
  padding: 18px;
  display: grid;
  gap: 14px;
}

.section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.section-head h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 700;
}

.section-head span {
  color: var(--muted);
  font-size: 13px;
}

.admin-kpis {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.admin-hint {
  padding: 16px 0;
  color: var(--muted);
  font-size: 14px;
}

.admin-hint.error {
  color: var(--primary);
}

.table-card {
  overflow: hidden;
}

.table-scroll {
  overflow-x: auto;
}

.admin-table {
  width: 100%;
  min-width: 620px;
  border-collapse: collapse;
}

th,
td {
  padding: 10px 12px;
  border-bottom: 1px solid var(--line);
  text-align: left;
  font-size: 13px;
}

th {
  color: var(--muted);
  font-weight: 600;
  background: var(--card-soft);
}

.empty-row {
  color: var(--muted);
  text-align: center;
}

@media (min-width: 768px) {
  .admin-kpis {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (min-width: 1024px) {
  .admin-kpis {
    grid-template-columns: repeat(6, minmax(0, 1fr));
  }
}
</style>
