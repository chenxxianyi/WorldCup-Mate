<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import ShareCopyCard from '@/components/ai/ShareCopyCard.vue'
import { apiGetUpcomingMatches, apiListMatches } from '@/api/matches'
import { normalizeMatch, type Match } from '@/types/match'
import type { ShareCopyLength, ShareCopyPlatform, ShareCopyTone } from '@/types/ai'
import { useAIStore } from '@/stores/useAIStore'

const ai = useAIStore()
const matches = ref<Match[]>([])
const selectedMatchId = ref<number | null>(null)
const matchDropdownOpen = ref(false)
const matchSelectRef = ref<HTMLElement | null>(null)
const platform = ref<ShareCopyPlatform>('wechat')
const tone = ref<ShareCopyTone>('relaxed')
const length = ref<ShareCopyLength>('short')

const selectedMatch = computed(() => matches.value.find((item) => item.id === selectedMatchId.value) || null)
const selectedMatchLabel = computed(() => {
  if (!selectedMatch.value) return '请选择比赛'
  return `${selectedMatch.value.home_team_name} vs ${selectedMatch.value.away_team_name}`
})

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

function toggleMatchDropdown() {
  if (!matches.value.length) return
  matchDropdownOpen.value = !matchDropdownOpen.value
}

function closeMatchDropdown() {
  matchDropdownOpen.value = false
}

function selectMatch(matchId: number) {
  selectedMatchId.value = matchId
  closeMatchDropdown()
}

function moveMatchSelection(step: 1 | -1) {
  if (!matches.value.length) return
  const currentIndex = matches.value.findIndex((item) => item.id === selectedMatchId.value)
  const fallbackIndex = step > 0 ? 0 : matches.value.length - 1
  const nextIndex =
    currentIndex === -1
      ? fallbackIndex
      : (currentIndex + step + matches.value.length) % matches.value.length
  selectedMatchId.value = matches.value[nextIndex].id
  matchDropdownOpen.value = true
}

function handleDocumentPointerDown(event: PointerEvent) {
  if (!matchSelectRef.value?.contains(event.target as Node)) {
    closeMatchDropdown()
  }
}

onMounted(() => {
  loadMatches()
  document.addEventListener('pointerdown', handleDocumentPointerDown)
})

onBeforeUnmount(() => {
  document.removeEventListener('pointerdown', handleDocumentPointerDown)
})
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
      <div ref="matchSelectRef" class="match-field">
        <span>比赛</span>
        <button
          class="match-select"
          type="button"
          :class="{ open: matchDropdownOpen }"
          :disabled="!matches.length"
          aria-haspopup="listbox"
          :aria-expanded="matchDropdownOpen"
          @click="toggleMatchDropdown"
          @keydown.down.prevent="moveMatchSelection(1)"
          @keydown.up.prevent="moveMatchSelection(-1)"
          @keydown.esc.prevent="closeMatchDropdown"
        >
          <span class="match-select-text">{{ selectedMatchLabel }}</span>
          <span class="material-symbols-outlined" aria-hidden="true">expand_more</span>
        </button>

        <Transition name="match-options">
          <div
            v-if="matchDropdownOpen"
            class="match-menu"
            role="listbox"
            :aria-activedescendant="selectedMatchId ? `match-option-${selectedMatchId}` : undefined"
          >
            <button
              v-for="match in matches"
              :id="`match-option-${match.id}`"
              :key="match.id"
              class="match-option"
              type="button"
              role="option"
              :aria-selected="selectedMatchId === match.id"
              :class="{ active: selectedMatchId === match.id }"
              @click="selectMatch(match.id)"
            >
              <span class="match-pair">{{ match.home_team_name }} vs {{ match.away_team_name }}</span>
              <span class="material-symbols-outlined" aria-hidden="true">check</span>
            </button>
          </div>
        </Transition>
      </div>

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

.match-field,
.option-group {
  display: grid;
  gap: 8px;
}

.match-field > span,
.option-group > span {
  color: var(--muted);
  font-size: 13px;
  font-weight: 750;
}

.match-field {
  position: relative;
  z-index: 4;
}

.match-select {
  width: 100%;
  min-height: 48px;
  display: grid;
  grid-template-columns: minmax(0, 1fr) 24px;
  align-items: center;
  gap: 10px;
  border: 1px solid color-mix(in srgb, var(--line) 88%, var(--text));
  border-radius: 14px;
  outline: none;
  padding: 0 13px 0 15px;
  color: var(--text);
  background:
    linear-gradient(180deg, color-mix(in srgb, var(--card) 92%, var(--card-soft)), var(--card-soft));
  box-shadow: inset 0 1px 0 color-mix(in srgb, #fff 64%, transparent);
  text-align: left;
  transition: border-color 160ms ease-out, background 160ms ease-out, box-shadow 160ms ease-out;
}

.match-select:hover:not(:disabled),
.match-select.open {
  border-color: color-mix(in srgb, var(--primary) 38%, var(--line));
  background: var(--card);
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--primary) 8%, transparent);
}

.match-select:focus-visible {
  border-color: color-mix(in srgb, var(--primary) 55%, var(--line));
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--primary) 14%, transparent);
}

.match-select:disabled {
  cursor: not-allowed;
  color: var(--weak);
  background: var(--card-soft);
}

.match-select .material-symbols-outlined {
  color: var(--weak);
  font-size: 22px;
  transition: transform 160ms ease-out, color 160ms ease-out;
}

.match-select.open .material-symbols-outlined {
  color: var(--primary);
  transform: rotate(180deg);
}

.match-select-text {
  min-width: 0;
  overflow: hidden;
  font-size: 15px;
  font-weight: 750;
  line-height: 1.2;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.match-menu {
  position: absolute;
  top: calc(100% + 8px);
  left: 0;
  right: 0;
  max-height: 286px;
  overflow: auto;
  overscroll-behavior: contain;
  padding: 6px;
  border: 1px solid color-mix(in srgb, var(--line) 86%, var(--text));
  border-radius: 16px;
  background: color-mix(in srgb, var(--card) 96%, transparent);
  box-shadow: 0 18px 48px rgba(15, 23, 42, 0.18);
  backdrop-filter: blur(16px);
}

.match-menu::-webkit-scrollbar {
  width: 8px;
}

.match-menu::-webkit-scrollbar-thumb {
  border: 2px solid transparent;
  border-radius: 999px;
  background: color-mix(in srgb, var(--weak) 36%, transparent);
  background-clip: padding-box;
}

.match-option {
  width: 100%;
  min-height: 42px;
  display: grid;
  grid-template-columns: minmax(0, 1fr) 20px;
  align-items: center;
  gap: 10px;
  border: 0;
  border-radius: 11px;
  padding: 0 10px 0 12px;
  color: var(--text);
  background: transparent;
  text-align: left;
  transition: background 140ms ease-out, color 140ms ease-out;
}

.match-option:hover,
.match-option:focus-visible {
  outline: none;
  background: color-mix(in srgb, var(--primary) 7%, var(--card));
}

.match-option.active {
  color: var(--primary);
  background: color-mix(in srgb, var(--primary) 11%, var(--card));
}

.match-pair {
  min-width: 0;
  overflow: hidden;
  font-size: 14px;
  font-weight: 720;
  line-height: 1.2;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.match-option .material-symbols-outlined {
  opacity: 0;
  color: var(--primary);
  font-size: 19px;
}

.match-option.active .material-symbols-outlined {
  opacity: 1;
}

.match-options-enter-active,
.match-options-leave-active {
  transition: opacity 140ms ease-out, transform 140ms ease-out;
}

.match-options-enter-from,
.match-options-leave-to {
  opacity: 0;
  transform: translateY(-4px);
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
