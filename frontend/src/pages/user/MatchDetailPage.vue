<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useMatchStore } from '@/stores/useMatchStore'
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import { useReminderStore } from '@/stores/useReminderStore'
import { useAuthStore } from '@/stores/useAuthStore'
import { toLocalTime, toUTCTime } from '@/utils/timezone'
import { useSettingStore } from '@/stores/useSettingStore'
import { apiGetGroupStandings } from '@/api/standings'
import { normalizeStanding, type Standing } from '@/types/standing'
import TeamFlag from '@/components/common/TeamFlag.vue'

const route = useRoute()
const matchStore = useMatchStore()
const fav = useFavoriteStore()
const reminder = useReminderStore()
const auth = useAuthStore()
const settings = useSettingStore()
const match = ref<any>(null)
const groupStandings = ref<Standing[]>([])

onMounted(async () => {
  const id = Number(route.params.id)
  match.value = await matchStore.fetchMatchDetail(id)
  if (auth.isLoggedIn) {
    fav.fetchFavoriteTeams()
    fav.fetchFavoriteMatches()
    reminder.fetchReminders()
  }
  if (match.value?.group_name) {
    const groupNum = match.value.group_name.charCodeAt(6) - 64
    try {
      const res = await apiGetGroupStandings(groupNum) as any[]
      groupStandings.value = res.map(normalizeStanding)
    } catch {}
  }
})

function rowClass(status: string) {
  if (status === '晋级') return 'rank-ok'
  if (status === '待定') return 'rank-mid'
  return 'rank-out'
}
</script>

<template>
  <div v-if="match">
    <article class="card detail-hero">
      <div class="match-top">
        <span class="tag gold">焦点大战 · {{ match.group_name }}</span>
        <span v-if="match.status === 'live'" class="tag live">
          <i class="live-dot" style="background: #fff"></i> {{ match.minute }}' LIVE
        </span>
        <span v-else-if="match.status === 'finished'" class="tag green">已结束</span>
        <span v-else class="tag">未开始</span>
      </div>
      <div class="detail-score">
        <div class="detail-team">
          <TeamFlag :value="match.home_flag" :alt="match.home_team_name" :fallback="match.home_team_code" size="lg" />
          <h2>{{ match.home_team_name }}</h2>
          <p class="team-meta">{{ match.home_team_code }}</p>
        </div>
        <div class="big-vs">
          <template v-if="match.status === 'live' || match.status === 'finished'">
            {{ match.home_score }} - {{ match.away_score }}
          </template>
          <template v-else>VS</template>
        </div>
        <div class="detail-team">
          <TeamFlag :value="match.away_flag" :alt="match.away_team_name" :fallback="match.away_team_code" size="lg" />
          <h2>{{ match.away_team_name }}</h2>
          <p class="team-meta">{{ match.away_team_code }}</p>
        </div>
      </div>
      <div class="info-grid">
        <div class="info-cell"><span>本地时间</span><strong>{{ toLocalTime(match.kickoff_time_utc, settings.timezone) }}</strong></div>
        <div class="info-cell"><span>UTC 时间</span><strong>{{ toUTCTime(match.kickoff_time_utc) }}</strong></div>
        <div class="info-cell"><span>城市</span><strong>{{ match.city }}</strong></div>
        <div class="info-cell"><span>球场</span><strong>{{ match.stadium }}</strong></div>
      </div>
      <div class="actions">
        <button
          class="pill-btn"
          :class="{ active: fav.isMatchFavorite(match.id) }"
          @click="fav.toggleMatchFavorite(match.id)"
        >
          <span class="material-symbols-outlined" style="font-size: 18px; vertical-align: -4px" :style="fav.isMatchFavorite(match.id) ? 'font-variation-settings: \'FILL\' 1' : ''">star</span>
          {{ fav.isMatchFavorite(match.id) ? '已收藏' : '收藏' }}
        </button>
        <button
          class="pill-btn primary"
          @click="reminder.toggleReminder(match.id)"
        >
          <span class="material-symbols-outlined" style="font-size: 18px; vertical-align: -4px">notifications</span>
          {{ reminder.hasReminder(match.id) ? '取消提醒' : '设置提醒' }}
        </button>
      </div>
    </article>

    <section class="section">
      <div class="section-head">
        <h2>所在小组积分</h2>
        <span>{{ match.group_name }}</span>
      </div>
      <div class="card table-card table-scroll">
        <table class="standing-table">
          <thead>
            <tr>
              <th>排名</th><th>球队</th><th>场次</th><th>胜</th><th>平</th><th>负</th><th>净胜球</th><th>积分</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(s, i) in groupStandings" :key="s.team_id" :class="rowClass(s.status)">
              <td>{{ i + 1 }}</td>
              <td class="team-cell">
                <TeamFlag :value="s.flag" :alt="s.team_name" :fallback="s.team_code" size="sm" />
                <span>{{ s.team_name }}</span>
              </td>
              <td>{{ s.played }}</td>
              <td>{{ s.won }}</td>
              <td>{{ s.drawn }}</td>
              <td>{{ s.lost }}</td>
              <td>{{ s.goal_difference > 0 ? '+' : '' }}{{ s.goal_difference }}</td>
              <td><b>{{ s.points }}</b></td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </div>
</template>

<style scoped>
.detail-hero {
  padding: 18px;
}

.match-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.detail-score {
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  gap: 12px;
  align-items: center;
  margin: 22px 0;
}

.detail-team {
  text-align: center;
}

.detail-team h2 {
  margin: 8px 0 0;
  font-size: 18px;
  font-weight: 750;
}

.team-meta {
  color: var(--weak);
  font-size: 12px;
}

.big-vs {
  display: grid;
  place-items: center;
  min-width: 88px;
  min-height: 68px;
  border-radius: 18px;
  color: var(--primary);
  font-size: 24px;
  font-weight: 850;
  background: color-mix(in srgb, var(--primary) 10%, transparent);
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.info-cell {
  padding: 13px;
  border-radius: 14px;
  background: var(--card-soft);
}

.info-cell span {
  display: block;
  color: var(--muted);
  font-size: 12px;
}

.info-cell strong {
  display: block;
  margin-top: 5px;
  overflow-wrap: anywhere;
  font-size: 14px;
}

.actions {
  display: flex;
  gap: 8px;
  margin-top: 16px;
}

.section {
  margin-top: 18px;
}

.section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
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

.table-card {
  overflow: hidden;
}

.table-scroll {
  overflow-x: auto;
}

.standing-table {
  width: 100%;
  border-collapse: collapse;
  min-width: 520px;
}

th, td {
  padding: 12px 10px;
  border-bottom: 1px solid var(--line);
  text-align: left;
  font-size: 13px;
}

th {
  color: var(--muted);
  font-weight: 750;
  background: var(--card-soft);
}

td {
  color: var(--text);
}

.team-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 150px;
}

tr:last-child td {
  border-bottom: 0;
}

.rank-ok td:first-child {
  border-left: 4px solid var(--success);
}

.rank-mid td:first-child {
  border-left: 4px solid var(--gold);
}

.rank-out {
  color: var(--weak);
  opacity: 0.66;
}
</style>
