<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import {
  apiAdminListCompetitions,
  apiAdminCreateCompetition,
  apiAdminUpdateCompetition,
  type AdminCompetition,
  type CompetitionInput,
} from '@/api/admin'
import { seasonLabel } from '@/types/competition'

const items = ref<AdminCompetition[]>([])
const loading = ref(false)
const editing = ref<AdminCompetition | null>(null)
const creating = ref(false)
const form = ref<CompetitionInput>({ name: '', name_en: '', country: '', format: 'league', season: 2026, status: 'active', sort_order: 0 })
const saving = ref(false)
const error = ref('')
const notice = ref('')

const sorted = computed(() => [...items.value].sort((a, b) => a.sort_order - b.sort_order))

async function load() {
  loading.value = true
  try {
    items.value = await apiAdminListCompetitions()
  } finally {
    loading.value = false
  }
}

function openCreate() {
  form.value = { name: '', name_en: '', country: '', format: 'league', season: 2026, status: 'active', sort_order: sorted.value.length }
  creating.value = true
  error.value = ''
}

function openEdit(item: AdminCompetition) {
  editing.value = item
  form.value = {
    name: item.name, name_en: item.name_en, country: item.country,
    logo_url: item.logo_url, format: item.format, season: item.season,
    status: item.status, sort_order: item.sort_order,
  }
  error.value = ''
}

function close() {
  creating.value = false
  editing.value = null
}

async function save() {
  error.value = ''
  if (!form.value.name?.trim()) {
    error.value = '请填写赛事名称'
    return
  }
  saving.value = true
  try {
    if (creating.value) {
      if (!form.value.code?.trim()) {
        error.value = '请填写赛事代码（如 PL、PD）'
        return
      }
      await apiAdminCreateCompetition({ ...form.value, code: form.value.code.trim().toUpperCase() })
      notice.value = '赛事已创建，前台切换列表将自动更新'
    } else if (editing.value) {
      await apiAdminUpdateCompetition(editing.value.id, form.value)
      notice.value = '赛事已更新'
    }
    close()
    await load()
  } catch (e) {
    error.value = e instanceof Error ? e.message : '保存失败'
  } finally {
    saving.value = false
  }
}

async function toggleStatus(item: AdminCompetition) {
  const next = item.status === 'active' ? 'inactive' : 'active'
  try {
    await apiAdminUpdateCompetition(item.id, { status: next })
    item.status = next
    notice.value = next === 'active' ? `「${item.name}」已启用` : `「${item.name}」已停用（前台不再显示）`
  } catch (e) {
    error.value = e instanceof Error ? e.message : '操作失败'
  }
}

onMounted(load)
</script>

<template>
  <div class="admin-page">
    <div class="admin-head">
      <h2>赛事管理</h2>
      <p class="admin-sub">
        配置前台"切换赛事"选项：启停、显示名称、赛季与排序。停用的赛事数据保留，仅不在前台展示。
      </p>
      <button class="btn primary" type="button" @click="openCreate">+ 新增赛事</button>
    </div>

    <p v-if="notice" class="notice">{{ notice }}</p>

    <div class="admin-card">
      <div class="table-wrap">
        <table class="admin-table">
          <thead>
            <tr>
              <th>排序</th>
              <th>名称</th>
              <th>代码</th>
              <th>赛季</th>
              <th>国家</th>
              <th>状态</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in sorted" :key="item.id">
              <td><span class="order-badge">{{ item.sort_order }}</span></td>
              <td>
                <div class="team-cell">
                  <img v-if="item.logo_url" :src="item.logo_url" class="comp-logo" alt="" loading="lazy" />
                  <span v-else class="comp-logo comp-logo-text">{{ item.code.slice(0, 2) }}</span>
                  <div>
                    <strong>{{ item.name }}</strong>
                    <small>{{ item.name_en }}</small>
                  </div>
                </div>
              </td>
              <td><code>{{ item.code }}</code></td>
              <td>{{ seasonLabel(item.season) || '—' }}</td>
              <td>{{ item.country || '—' }}</td>
              <td>
                <button
                  class="status-toggle"
                  :class="item.status"
                  type="button"
                  @click="toggleStatus(item)"
                >
                  {{ item.status === 'active' ? '启用中' : '已停用' }}
                </button>
              </td>
              <td>
                <button class="btn sm" type="button" @click="openEdit(item)">编辑</button>
              </td>
            </tr>
            <tr v-if="!loading && sorted.length === 0">
              <td colspan="7" class="empty-cell">暂无赛事，点击"新增赛事"创建第一个联赛配置</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="admin-tip">
      💡 世界杯与五大联赛一样由本页控制：停用后前台切换列表即隐藏（数据保留）；全部停用仅剩世界杯时，前台仍以世界杯作为默认视图。
    </div>

    <div v-if="creating || editing" class="modal-mask" @click.self="close">
      <div class="modal">
        <h3>{{ creating ? '新增赛事' : '编辑赛事' }}</h3>
        <p v-if="error" class="form-error">{{ error }}</p>
        <div class="form-grid">
          <label>
            代码（如 PL）
            <input v-model="form.code" :disabled="!creating" placeholder="PL" />
          </label>
          <label>
            名称
            <input v-model="form.name" placeholder="英超" />
          </label>
          <label>
            英文名
            <input v-model="form.name_en" placeholder="PREMIER LEAGUE" />
          </label>
          <label>
            赛季起始年
            <input v-model.number="form.season" type="number" />
          </label>
          <label>
            国家
            <input v-model="form.country" placeholder="England" />
          </label>
          <label>
            排序（小的在前）
            <input v-model.number="form.sort_order" type="number" />
          </label>
          <label class="full">
            Logo 地址
            <input v-model="form.logo_url" placeholder="https://…" />
          </label>
          <label>
            状态
            <select v-model="form.status">
              <option value="active">启用（前台显示）</option>
              <option value="inactive">停用（前台隐藏）</option>
            </select>
          </label>
        </div>
        <div class="modal-actions">
          <button class="btn" type="button" @click="close">取消</button>
          <button class="btn primary" type="button" :disabled="saving" @click="save">
            {{ saving ? '保存中…' : '保存' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.admin-sub { color: var(--text-secondary, #8a8fa3); font-size: 13px; margin: 4px 0 16px; }
.notice { color: #16a34a; background: rgba(22, 163, 74, 0.1); padding: 8px 12px; border-radius: 8px; margin-bottom: 12px; font-size: 13px; }
.form-error { color: #dc2626; font-size: 13px; margin-bottom: 8px; }
.team-cell { display: flex; align-items: center; gap: 10px; }
.team-cell small { display: block; color: var(--text-secondary, #8a8fa3); }
.comp-logo { width: 30px; height: 30px; border-radius: 50%; object-fit: contain; background: rgba(255, 255, 255, 0.06); display: inline-flex; align-items: center; justify-content: center; }
.comp-logo-text { font-size: 11px; font-weight: 700; color: var(--text-secondary, #8a8fa3); }
.order-badge { background: rgba(255, 255, 255, 0.08); border-radius: 6px; padding: 2px 8px; font-size: 12px; }
.status-toggle { border: none; border-radius: 999px; padding: 4px 12px; font-size: 12px; cursor: pointer; }
.status-toggle.active { background: rgba(22, 163, 74, 0.15); color: #4ade80; }
.status-toggle.inactive { background: rgba(100, 116, 139, 0.15); color: #94a3b8; }
.admin-tip { margin-top: 12px; font-size: 13px; color: var(--text-secondary, #8a8fa3); }
.empty-cell { text-align: center; padding: 24px; color: var(--text-secondary, #8a8fa3); }
.form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.form-grid .full { grid-column: 1 / -1; }
.form-grid label { display: flex; flex-direction: column; gap: 4px; font-size: 13px; }
.form-grid input, .form-grid select { padding: 8px 10px; border-radius: 8px; border: 1px solid rgba(255, 255, 255, 0.12); background: rgba(255, 255, 255, 0.05); color: inherit; }
</style>
