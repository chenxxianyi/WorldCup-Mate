<script setup lang="ts">
import { useRouter } from 'vue-router'
import type { TodayRecommendations } from '@/types/ai'
import AIThinking from './AIThinking.vue'

defineProps<{
  data: TodayRecommendations | null
  loading?: boolean
  error?: string
}>()

const emit = defineEmits<{
  (e: 'generate'): void
  (e: 'refresh'): void
}>()

const router = useRouter()
</script>

<template>
  <article class="card recommend-card">
    <div class="card-head">
      <div>
        <span class="tag gold">今日推荐</span>
        <h2>今天看什么</h2>
      </div>
      <button class="pill-btn" type="button" :disabled="loading" @click="data ? emit('refresh') : emit('generate')">
        {{ data ? '重新生成' : '生成推荐' }}
      </button>
    </div>

    <div v-if="loading" class="state">
      <AIThinking />
      <span>正在筛选值得看的比赛...</span>
    </div>
    <div v-else-if="error" class="state error">
      <span>{{ error }}</span>
      <button class="pill-btn primary" type="button" @click="emit('generate')">重试</button>
    </div>
    <div v-else-if="!data" class="state">
      <span>让 AI 根据赛程和关注球队给你一份观赛清单。</span>
    </div>
    <div v-else class="content">
      <div v-if="data.only_one_match" class="only-one">
        <span>只看一场</span>
        <strong>{{ data.only_one_match.title }}</strong>
        <p>{{ data.only_one_match.reason }}</p>
        <button class="pill-btn primary" type="button" @click="router.push(`/matches/${data.only_one_match?.match_id}`)">
          查看比赛
        </button>
      </div>

      <button
        v-for="item in data.recommendations"
        :key="item.match_id"
        class="rec-row"
        type="button"
        @click="router.push(`/matches/${item.match_id}`)"
      >
        <div>
          <strong>{{ item.title }}</strong>
          <span>{{ item.kickoff_time || '时间待定' }}</span>
          <p>{{ item.reason }}</p>
        </div>
        <b v-if="item.rating">{{ item.rating }}/5</b>
      </button>

      <p v-if="data.note" class="note">{{ data.note }}</p>
    </div>
  </article>
</template>

<style scoped>
.recommend-card {
  padding: 16px;
}

.card-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.card-head h2 {
  margin: 9px 0 0;
  font-size: 18px;
}

.state {
  min-height: 112px;
  display: grid;
  place-items: center;
  gap: 10px;
  text-align: center;
  color: var(--muted);
}

.state.error {
  color: var(--primary);
}

.content {
  display: grid;
  gap: 10px;
  margin-top: 14px;
}

.only-one {
  display: grid;
  gap: 8px;
  padding: 13px;
  border-radius: var(--radius-lg);
  background: color-mix(in srgb, var(--primary) 9%, var(--card-soft));
}

.only-one span {
  color: var(--primary);
  font-size: 12px;
  font-weight: 800;
}

.only-one strong,
.rec-row strong {
  font-size: 15px;
}

.only-one p,
.rec-row p,
.note {
  margin: 0;
  color: var(--muted);
  line-height: 1.55;
  font-size: 13px;
}

.rec-row {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 12px;
  border: 1px solid var(--line);
  border-radius: var(--radius-md);
  text-align: left;
  color: var(--text);
  background: var(--card);
}

.rec-row span {
  display: block;
  margin: 4px 0;
  color: var(--weak);
  font-size: 12px;
}

.rec-row b {
  min-width: 44px;
  color: var(--primary);
  text-align: right;
}
</style>
