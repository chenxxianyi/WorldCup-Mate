<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import ShareCopyCard from '@/components/ai/ShareCopyCard.vue'
import { apiGetUpcomingMatches, apiListMatches } from '@/api/matches'
import { normalizeMatch, type Match } from '@/types/match'
import type { ShareCopyLength, ShareCopyPlatform, ShareCopyTone } from '@/types/ai'
import { useAIStore } from '@/stores/useAIStore'

const ai = useAIStore()
const matches = ref<Match[]>([])
const selectedMatchId = ref<number | null>(null)
const platform = ref<ShareCopyPlatform>('wechat')
const tone = ref<ShareCopyTone>('relaxed')
const length = ref<ShareCopyLength>('short')

const selectedMatch = computed(() => matches.value.find((item) => item.id === selectedMatchId.value) || null)

const platforms: Array<{ value: ShareCopyPlatform; label: string }> = [
  { value: 'wechat', label: '朋友圈' },
  { value: 'group', label: '微信群' },
  { value: 'xiaohongshu', label: '小红书' },
  { value: 'weibo', label: '微博' },
  { value: 'general', label: '通用' },
]

const tones: Array<{ value: ShareCopyTone; label: string }> = [
  { value: 'relaxed', label: '轻松' },
  { value: 'passionate', label: '热血' },
  { value: 'professional', label: '专业' },
  { value: 'beginner', label: '小白友好' },
]

const lengths: Array<{ value: ShareCopyLength; label: string }> = [
  { value: 'short', label: '短' },
  { value: 'medium', label: '中' },
  { value: 'long', label: '长' },
]

async function loadMatches() {
  try {
    const res = await apiGetUpcomingMatches() as any[]
    matches.value = (res || []).map(normalizeMatch)
  } catch {
    try {
      const res = await apiListMatches({ page: 1, page_size: 12 }) as any
      matches.value = (res.list || res || []).map(normalizeMatch)
    } catch {
      matches.value = []
    }
  }
  selectedMatchId.value = matches.value[0]?.id || null
}

function generate() {
  if (!selectedMatchId.value) return
  ai.generateShareCopy(selectedMatchId.value, platform.value, tone.value, length.value).catch(() => {})
}

onMounted(loadMatches)
</script>

<template>
  <div class="share-page">
    <div class="section-head">
      <div>
        <h2>分享文案</h2>
        <span>选择比赛和语气，生成可以直接复制的短文案</span>
      </div>
    </div>

    <section class="card form-card">
      <label>
        <span>比赛</span>
        <select v-model.number="selectedMatchId">
          <option v-for="match in matches" :key="match.id" :value="match.id">
            {{ match.home_team_name }} vs {{ match.away_team_name }}
          </option>
        </select>
      </label>

      <div class="option-group">
        <span>平台</span>
        <div class="segmented">
          <button
            v-for="item in platforms"
            :key="item.value"
            type="button"
            :class="{ active: platform === item.value }"
            @click="platform = item.value"
          >
            {{ item.label }}
          </button>
        </div>
      </div>

      <div class="option-group">
        <span>语气</span>
        <div class="segmented">
          <button
            v-for="item in tones"
            :key="item.value"
            type="button"
            :class="{ active: tone === item.value }"
            @click="tone = item.value"
          >
            {{ item.label }}
          </button>
        </div>
      </div>

      <div class="option-group">
        <span>长度</span>
        <div class="segmented compact">
          <button
            v-for="item in lengths"
            :key="item.value"
            type="button"
            :class="{ active: length === item.value }"
            @click="length = item.value"
          >
            {{ item.label }}
          </button>
        </div>
      </div>

      <p v-if="selectedMatch" class="match-note">
        当前选择：{{ selectedMatch.home_team_name }} vs {{ selectedMatch.away_team_name }}
      </p>
    </section>

    <ShareCopyCard
      :result="ai.shareCopyResult"
      :loading="ai.shareCopyLoading"
      :error="ai.shareCopyError"
      @generate="generate"
    />
  </div>
</template>

<style scoped>
.share-page {
  display: grid;
  gap: 14px;
}

.section-head h2 {
  margin: 0;
  font-size: 20px;
}

.section-head span {
  display: block;
  margin-top: 5px;
  color: var(--muted);
  font-size: 13px;
}

.form-card {
  display: grid;
  gap: 14px;
  padding: 16px;
}

label,
.option-group {
  display: grid;
  gap: 8px;
}

label span,
.option-group > span {
  color: var(--muted);
  font-size: 13px;
  font-weight: 750;
}

select {
  min-height: 42px;
  border: 1px solid var(--line);
  border-radius: var(--radius-md);
  padding: 0 12px;
  color: var(--text);
  background: var(--card-soft);
}

.segmented {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.segmented button {
  min-height: 36px;
  padding: 0 12px;
  border: 1px solid var(--line);
  border-radius: 999px;
  color: var(--muted);
  background: var(--card);
  font-size: 13px;
  font-weight: 700;
}

.segmented button.active {
  color: #fff;
  border-color: transparent;
  background: var(--primary);
}

.segmented.compact button {
  min-width: 52px;
}

.match-note {
  margin: 0;
  color: var(--muted);
  font-size: 13px;
}
</style>
