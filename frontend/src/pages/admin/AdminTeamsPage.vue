<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  apiAdminListTeams,
  apiAdminCreateTeam,
  apiAdminUpdateTeam,
  apiAdminDeleteTeam,
} from '@/api/admin'
import type { TeamInput } from '@/api/admin'
import type { ApiTeam as RawApiTeam } from '@/types/team'

const search = ref('')
const teamType = ref('')
const teams = ref<RawApiTeam[]>([])
const loading = ref(false)

// Dialog state for create/edit
const dialogVisible = ref(false)
const editing = ref(false)
const editingId = ref<number | null>(null)
const form = ref<TeamInput>({
  name: '',
  name_en: '',
  fifa_code: '',
  external_code: '',
  team_type: 'national',
  flag_url: '',
  continent: '欧洲',
  country: '',
  venue: '',
  group_id: undefined,
  coach: '',
  description: '',
})

const filteredTeams = computed(() => {
  const q = search.value.trim().toLowerCase()
  if (!q) return teams.value
  return teams.value.filter((team) =>
    [team.name, team.name_en, team.fifa_code, team.continent, team.group?.name]
      .filter(Boolean)
      .some((value) => String(value).toLowerCase().includes(q)),
  )
})

async function loadTeams() {
  loading.value = true
  try {
    const params: Record<string, any> = { page: 1, page_size: 100 }
    if (teamType.value) params.teamType = teamType.value
    const res = await apiAdminListTeams(params)
    teams.value = res.list
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editing.value = false
  editingId.value = null
  form.value = {
    name: '',
    name_en: '',
    fifa_code: '',
    external_code: '',
    team_type: 'national',
    flag_url: '',
    continent: '欧洲',
    country: '',
    venue: '',
    group_id: undefined,
    coach: '',
    description: '',
  }
  dialogVisible.value = true
}

function openEdit(team: RawApiTeam) {
  editing.value = true
  editingId.value = team.id
  form.value = {
    name: team.name,
    name_en: team.name_en || '',
    fifa_code: team.fifa_code || '',
    external_code: team.external_code || '',
    team_type: team.team_type || 'national',
    flag_url: team.flag_url || '',
    continent: team.continent || '欧洲',
    country: team.country || '',
    venue: team.venue || '',
    group_id: team.group_id ?? undefined,
    coach: team.coach || '',
    description: team.description || '',
  }
  dialogVisible.value = true
}

async function save() {
  try {
    if (editing.value && editingId.value) {
      await apiAdminUpdateTeam(editingId.value, form.value)
      ElMessage.success('球队已更新')
    } else {
      await apiAdminCreateTeam(form.value)
      ElMessage.success('球队已创建')
    }
    dialogVisible.value = false
    await loadTeams()
  } catch (e: any) {
    ElMessage.error(e?.response?.data?.message || '保存失败')
  }
}

async function remove(team: RawApiTeam) {
  try {
    await ElMessageBox.confirm(
      `确定删除球队「${team.name}」吗？关联的比赛或积分将阻止删除。`,
      '删除确认',
      { type: 'warning' },
    )
    await apiAdminDeleteTeam(team.id)
    ElMessage.success('已删除')
    await loadTeams()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e?.response?.data?.message || '删除失败')
    }
  }
}

onMounted(loadTeams)
watch(teamType, loadTeams)
watch(search, () => {})
</script>

<template>
  <div class="admin-page">
    <div class="admin-head">
      <div>
        <h2>球队管理</h2>
        <span>{{ filteredTeams.length }} 支球队</span>
      </div>
      <div class="tools">
        <select
          v-model="teamType"
          class="admin-select"
          aria-label="球队类型"
        >
          <option value="">
            全部类型
          </option>
          <option value="national">
            国家队
          </option>
          <option value="club">
            俱乐部
          </option>
        </select>
        <input
          v-model="search"
          class="admin-search"
          placeholder="搜索球队 / 小组 / 大洲"
        >
        <el-button type="primary" @click="openCreate">
          新建球队
        </el-button>
      </div>
    </div>

    <div class="card table-card table-scroll">
      <table class="admin-table">
        <thead>
          <tr>
            <th>球队</th><th>代码</th><th>英文名</th><th>类型</th><th>大洲</th><th>小组</th><th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="team in filteredTeams"
            :key="team.id"
          >
            <td><b>{{ team.name }}</b></td>
            <td>{{ team.fifa_code }}</td>
            <td>{{ team.name_en }}</td>
            <td>{{ team.team_type }}</td>
            <td>{{ team.continent }}</td>
            <td>{{ team.group?.name || '-' }}</td>
            <td>
              <el-button size="small" @click="openEdit(team)">编辑</el-button>
              <el-button size="small" type="danger" @click="remove(team)">删除</el-button>
            </td>
          </tr>
          <tr v-if="!loading && !filteredTeams.length">
            <td
              colspan="7"
              class="empty-row"
            >
              暂无球队数据
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create / Edit dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="editing ? '编辑球队' : '新建球队'"
      width="520px"
      destroy-on-close
    >
      <el-form :model="form" label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="英文名">
          <el-input v-model="form.name_en" />
        </el-form-item>
        <el-form-item label="FIFA 代码">
          <el-input v-model="form.fifa_code" />
        </el-form-item>
        <el-form-item label="外部代码">
          <el-input v-model="form.external_code" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="form.team_type">
            <el-option label="国家队" value="national" />
            <el-option label="俱乐部" value="club" />
          </el-select>
        </el-form-item>
        <el-form-item label="大洲">
          <el-select v-model="form.continent">
            <el-option label="亚洲" value="亚洲" />
            <el-option label="欧洲" value="欧洲" />
            <el-option label="南美洲" value="南美洲" />
            <el-option label="北美洲" value="北美洲" />
            <el-option label="非洲" value="非洲" />
            <el-option label="大洋洲" value="大洋洲" />
          </el-select>
        </el-form-item>
        <el-form-item label="国家">
          <el-input v-model="form.country" />
        </el-form-item>
        <el-form-item label="球场">
          <el-input v-model="form.venue" />
        </el-form-item>
        <el-form-item label="教练">
          <el-input v-model="form.coach" />
        </el-form-item>
        <el-form-item label="简介">
          <el-input v-model="form.description" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.admin-page {
  display: grid;
  gap: 14px;
}

.admin-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.tools {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
  flex-wrap: wrap;
}

.admin-head h2 {
  margin: 0;
  font-size: 20px;
}

.admin-head span {
  color: var(--muted);
  font-size: 13px;
}

.admin-search {
  width: min(320px, 100%);
  min-height: 40px;
  padding: 0 14px;
  border: 1px solid var(--line);
  border-radius: 8px;
  color: var(--text);
  background: var(--card);
}

.table-card {
  overflow: hidden;
}

.table-scroll {
  overflow-x: auto;
}

.admin-table {
  width: 100%;
  min-width: 640px;
  border-collapse: collapse;
}

th,
td {
  padding: 12px 10px;
  border-bottom: 1px solid var(--line);
  text-align: left;
  font-size: 13px;
}

th {
  color: var(--muted);
  background: var(--card-soft);
}

.empty-row {
  color: var(--muted);
  text-align: center;
}
</style>
