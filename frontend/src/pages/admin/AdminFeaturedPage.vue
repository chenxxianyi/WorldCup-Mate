<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import {
  apiAdminListFeatured,
  apiAdminFeaturedMatches,
  apiAdminUpdateFeatured,
  type FeaturedConfig,
  type FeaturedMatchPick,
} from '@/api/featured'
import { apiAdminListCompetitions, type AdminCompetition } from '@/api/admin'

const competitions = ref<AdminCompetition[]>([])
const configs = ref<FeaturedConfig[]>([])
const selectedCode = ref('')
const matchOptions = ref<FeaturedMatchPick[]>([])
const loading = ref(false)
const saving = ref(false)
const error = ref('')
const notice = ref('')

const form = ref({ match_id: null as number | null, tagline: '', description: '', stage_label: '', enabled: true })

const orphanMatch = computed(() => {
  if (form.value.match_id == null) return null
  return matchOptions.value.find((m) => m.id === form.value.match_id) ?? null
})

const selectedComp = computed(() => competitions.value.find((c) => c.code === selectedCode.value))

async function load() {
  loading.value = true
  try {
    competitions.value = await apiAdminListCompetitions()
    configs.value = await apiAdminListFeatured()
    if (!selectedCode.value && competitions.value.length) {
      selectedCode.value = competitions.value[0].code
    }
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : '加载失败'
  } finally {
    loading.value = false
  }
}

async function loadMatches() {
  if (!selectedCode.value) return
  try {
    matchOptions.value = await apiAdminFeaturedMatches(selectedCode.value)
  } catch {
    matchOptions.value = []
  }
}

watch(selectedCode, async (code) => {
  if (!code) return
  const cfg = configs.value.find((c) => competitions.value.find((x) => x.id === c.competition_id)?.code === code)
  form.value = {
    match_id: cfg?.match_id ?? null,
    tagline: cfg?.tagline ?? '',
    description: cfg?.description ?? '',
    stage_label: cfg?.stage_label ?? '',
    enabled: cfg?.enabled ?? true,
  }
  await loadMatches()
})

async function save() {
  if (!selectedCode.value) return
  saving.value = true
  error.value = ''
  try {
    await apiAdminUpdateFeatured(selectedCode.value, {
      match_id: form.value.match_id,
      tagline: form.value.tagline.trim(),
      description: form.value.description.trim(),
      stage_label: form.value.stage_label.trim(),
      enabled: form.value.enabled,
    })
    notice.value = `「${selectedComp.value?.name ?? selectedCode.value}」焦点配置已保存，前台首页立即生效`
    await load()
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : '保存失败'
  } finally {
    saving.value = false
  }
}

function matchLabel(m: FeaturedMatchPick) {
  const k = m.kickoff_time_utc ? new Date(m.kickoff_time_utc).toLocaleString('zh-CN', { timeZone: 'Asia/Shanghai', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' }) : '时间待定'
  const statusText = m.status === 'live' ? '直播中' : m.status === 'finished' ? '已结束' : '未开始'
  return `${m.home} vs ${m.away}（${k} · ${statusText}）`
}

onMounted(load)
</script>

<template>
  <div class="admin-page">
    <div class="admin-head">
      <h2>焦点管理</h2>
      <p class="admin-sub">
        配置前台首页焦点区：指定焦点比赛与宣传文字。未配置的赛事沿用默认文案，焦点比赛自动选取下一场。
      </p>
    </div>

    <p v-if="notice" class="notice">{{ notice }}</p>
    <p v-if="error" class="form-error">{{ error }}</p>

    <div class="admin-card">
      <div class="form-row">
        <label>
          选择赛事
          <select v-model="selectedCode">
            <option v-for="c in competitions" :key="c.code" :value="c.code">{{ c.name }}（{{ c.code }}）</option>
          </select>
        </label>
        <label class="toggle-label">
          <input v-model="form.enabled" type="checkbox" />
          启用焦点配置
        </label>
      </div>

      <div class="form-row">
        <label>
          焦点比赛
          <select v-model="form.match_id">
            <option :value="null">自动（下一场未结束的比赛）</option>
            <option
              v-if="form.match_id != null && !orphanMatch"
              :value="form.match_id"
              disabled
            >
              已选比赛（跨赛季/不在当前列表，保存将保留原选择）
            </option>
            <option v-for="m in matchOptions" :key="m.id" :value="m.id">{{ matchLabel(m) }}</option>
          </select>
        </label>
      </div>

      <div class="form-row">
        <label>
          主宣传语
          <input v-model="form.tagline" placeholder="每一轮，都有新的主角" maxlength="60" />
        </label>
      </div>
      <div class="form-row">
        <label>
          副文案
          <textarea v-model="form.description" rows="3" placeholder="一句话介绍该赛事看点…" maxlength="200" />
        </label>
      </div>
      <div class="form-row">
        <label>
          阶段标签
          <input v-model="form.stage_label" placeholder="MATCHWEEK / 第 X 轮" maxlength="30" />
        </label>
      </div>

      <div class="modal-actions">
        <button class="btn primary" type="button" :disabled="saving || !selectedCode" @click="save">
          {{ saving ? '保存中…' : '保存配置' }}
        </button>
      </div>
    </div>

    <div v-if="selectedComp" class="admin-card preview-card">
      <h3>前台预览</h3>
      <div class="hero-preview">
        <div class="hero-preview-mark">{{ selectedComp.code.slice(0, 2) }}</div>
        <div>
          <strong>{{ selectedComp.name }} · {{ selectedComp.season }} 赛季</strong>
          <h4>{{ form.tagline || '（默认宣传语）' }}</h4>
          <p>{{ form.description || '（默认描述）' }}</p>
          <small>{{ form.stage_label || '（默认标签）' }} · 焦点比赛：{{ form.match_id ? matchOptions.find((m) => m.id === form.match_id)?.home + ' vs ' + matchOptions.find((m) => m.id === form.match_id)?.away : '自动' }}</small>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.admin-sub { color: var(--text-secondary, #8a8fa3); font-size: 13px; margin: 4px 0 16px; }
.notice { color: #16a34a; background: rgba(22, 163, 74, 0.1); padding: 8px 12px; border-radius: 8px; margin-bottom: 12px; font-size: 13px; }
.form-error { color: #dc2626; font-size: 13px; margin-bottom: 8px; }
.form-row { display: flex; flex-direction: column; gap: 6px; margin-bottom: 14px; }
.form-row label { font-size: 13px; display: flex; flex-direction: column; gap: 4px; }
.toggle-label { flex-direction: row !important; align-items: center; gap: 8px !important; }
.form-row input, .form-row select, .form-row textarea {
  padding: 8px 10px; border-radius: 8px; border: 1px solid rgba(255, 255, 255, 0.12);
  background: rgba(255, 255, 255, 0.05); color: inherit; font: inherit;
}
.preview-card { margin-top: 16px; }
.hero-preview { display: flex; gap: 14px; align-items: flex-start; padding: 16px; border-radius: 12px;
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.18), rgba(168, 85, 247, 0.12)); }
.hero-preview-mark { width: 44px; height: 44px; border-radius: 50%; background: rgba(255, 255, 255, 0.14);
  display: flex; align-items: center; justify-content: center; font-weight: 800; }
.hero-preview h4 { margin: 6px 0; font-size: 18px; }
.hero-preview p { margin: 0 0 6px; font-size: 13px; color: var(--text-secondary, #8a8fa3); }
.hero-preview small { color: var(--text-secondary, #8a8fa3); }
.modal-actions { display: flex; justify-content: flex-end; margin-top: 8px; }
</style>
