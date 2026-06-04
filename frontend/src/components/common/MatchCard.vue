<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import { useReminderStore } from '@/stores/useReminderStore'
import { useSettingStore } from '@/stores/useSettingStore'
import type { Match } from '@/types/match'
import TeamFlag from '@/components/common/TeamFlag.vue'

const props = defineProps<{
  match: Match
  featured?: boolean
}>()

const router = useRouter()
const fav = useFavoriteStore()
const reminder = useReminderStore()
const settings = useSettingStore()

const showChannelPicker = ref(false)

function stageLabel(stage: Match['stage']) {
  const labels: Record<Match['stage'], string> = {
    group: '小组赛',
    group_stage: '小组赛',
    round_of_32: '32强赛',
    round_of_16: '16强赛',
    quarter_final: '四分之一决赛',
    semi_final: '半决赛',
    third_place: '季军赛',
    final: '决赛',
  }
  return labels[stage] || '世界杯比赛'
}

function stageTagText(match: Match) {
  const label = stageLabel(match.stage)
  if (match.group_name) {
    return `${label} · ${match.group_name}`
  }
  if (match.stage === 'group' || match.stage === 'group_stage') {
    return label
  }
  return label
}

function onBellClick() {
  if (reminder.hasReminder(props.match.id)) {
    reminder.toggleReminder(props.match.id)
  } else {
    showChannelPicker.value = !showChannelPicker.value
  }
}

function createWithChannel(channel: string) {
  reminder.toggleReminder(props.match.id, 30, channel)
  showChannelPicker.value = false
}

function goDetail() {
  router.push(`/matches/${props.match.id}`)
}
</script>

<template>
  <article
    class="card match-card"
    :class="{ featured: featured || match.is_featured }"
    @click="showChannelPicker = false"
  >
    <div v-if="featured || match.is_featured" class="hot-ribbon">热门比赛</div>
    <div class="match-top">
      <span class="tag blue">
        {{ stageTagText(match) }}
      </span>
      <span v-if="match.status === 'live'" class="tag live">
        <i class="live-dot" style="background: #fff"></i> {{ match.minute }}' LIVE
      </span>
      <span v-else-if="match.status === 'finished'" class="tag green">已结束</span>
      <span v-else class="tag">未开始</span>
    </div>
    <div class="teams-line" @click="goDetail">
      <div class="team-side">
        <TeamFlag :value="match.home_flag" :alt="match.home_team_name" :fallback="match.home_team_code" />
        <div class="team-copy">
          <span class="team-name">{{ match.home_team_name }}</span>
          <span class="team-meta">{{ match.home_team_code }}</span>
        </div>
      </div>
      <div v-if="match.status === 'live' || match.status === 'finished'" class="score">
        {{ match.home_score }}-{{ match.away_score }}
      </div>
      <div v-else class="vs">VS</div>
      <div class="team-side away">
        <div class="team-copy">
          <span class="team-name">{{ match.away_team_name }}</span>
          <span class="team-meta">{{ match.away_team_code }}</span>
        </div>
        <TeamFlag :value="match.away_flag" :alt="match.away_team_name" :fallback="match.away_team_code" />
      </div>
    </div>
    <div class="match-bottom">
      <div class="where">
        <template v-if="match.status === 'live'">
          进行中 · {{ match.minute }}′<br />{{ match.stadium }}
        </template>
        <template v-else>
          {{ match.local_kickoff_time.includes(' ') ? match.local_kickoff_time.split(' ')[1] : match.local_kickoff_time || 'TBD' }} · {{ match.city || 'TBD' }}<br />{{ match.stadium || 'TBD' }}
        </template>
      </div>
      <div class="actions">
        <div class="bell-wrap">
          <button
            class="icon-action bell"
            :class="{ active: reminder.hasReminder(match.id) }"
            title="提醒我"
            @click.stop="onBellClick"
          >
            <span class="material-symbols-outlined" :style="reminder.hasReminder(match.id) ? 'font-variation-settings: \'FILL\' 1' : ''">notifications</span>
          </button>
          <div v-if="showChannelPicker && !reminder.hasReminder(match.id)" class="channel-popover" @click.stop>
            <div class="channel-title">提醒方式</div>
            <button class="channel-option" @click="createWithChannel('site')">
              <span class="material-symbols-outlined">notifications</span>
              <span>站内通知</span>
            </button>
            <button class="channel-option" @click="createWithChannel('email')">
              <span class="material-symbols-outlined">mail</span>
              <span>邮件通知</span>
            </button>
          </div>
        </div>
        <button
          class="icon-action star"
          :class="{ active: fav.isMatchFavorite(match.id), ghost: !fav.isMatchFavorite(match.id) }"
          title="收藏"
          @click.stop="fav.toggleMatchFavorite(match.id)"
        >
          <span class="material-symbols-outlined" :style="fav.isMatchFavorite(match.id) ? 'font-variation-settings: \'FILL\' 1' : ''">star</span>
        </button>
      </div>
    </div>
  </article>
</template>

<style scoped>
.match-card {
  position: relative;
  overflow: hidden;
  padding: 15px;
}

.match-card.featured {
  border-color: color-mix(in srgb, var(--hot) 30%, var(--line));
  padding-top: 28px;
}

.hot-ribbon {
  position: absolute;
  top: 0;
  right: 0;
  padding: 6px 12px;
  border-radius: 0 0 0 12px;
  color: #fff;
  font-size: 10px;
  font-weight: 800;
  background: var(--hot);
}

.match-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.match-card.featured .match-top {
  padding-right: 0;
}

.teams-line {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto minmax(0, 1fr);
  align-items: center;
  gap: 14px;
  margin: 18px 0 16px;
  cursor: pointer;
}

.team-side {
  min-width: 0;
  min-height: 54px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.team-side.away {
  justify-content: flex-end;
  text-align: right;
}

.team-copy {
  min-width: 0;
  display: grid;
  gap: 5px;
}

.team-name {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-weight: 750;
  font-size: 16px;
}

.team-meta {
  display: block;
  color: var(--weak);
  font-size: 12px;
}

.score {
  min-width: 70px;
  text-align: center;
  font-weight: 850;
  font-size: 30px;
  line-height: 1;
}

.vs {
  min-width: 64px;
  height: 42px;
  display: grid;
  place-items: center;
  border-radius: 999px;
  color: var(--muted);
  font-size: 13px;
  font-weight: 800;
  background: var(--card-soft);
}

.match-bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding-top: 12px;
  border-top: 1px solid var(--line);
}

.where {
  min-width: 0;
  color: var(--muted);
  font-size: 12px;
  line-height: 1.5;
}

.actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.bell-wrap {
  position: relative;
}

.channel-popover {
  position: absolute;
  bottom: calc(100% + 8px);
  right: 0;
  z-index: 100;
  min-width: 140px;
  padding: 8px;
  border-radius: 14px;
  background: var(--card);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  border: 1px solid var(--line);
}

.channel-title {
  padding: 4px 8px 8px;
  font-size: 11px;
  font-weight: 700;
  color: var(--weak);
}

.channel-option {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 8px;
  border: none;
  border-radius: 10px;
  background: transparent;
  font-size: 13px;
  cursor: pointer;
  color: inherit;
}

.channel-option:hover {
  background: var(--card-soft);
}

.channel-option .material-symbols-outlined {
  font-size: 18px;
  color: var(--primary);
}
</style>
