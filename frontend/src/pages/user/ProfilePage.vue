<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import StatCard from '@/components/common/StatCard.vue'
import { useSettingStore } from '@/stores/useSettingStore'
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import { useReminderStore } from '@/stores/useReminderStore'
import { useTeamStore } from '@/stores/useTeamStore'
import { useAuthStore } from '@/stores/useAuthStore'
import { apiUpdateProfile } from '@/api/auth'

const settings = useSettingStore()
const fav = useFavoriteStore()
const reminder = useReminderStore()
const teamStore = useTeamStore()
const auth = useAuthStore()

// Notification email editing
const editingEmail = ref(false)
const notificationEmail = ref('')
const emailSaving = ref(false)

function startEditEmail() {
  notificationEmail.value = auth.user?.notificationEmail || ''
  editingEmail.value = true
}

function cancelEditEmail() {
  editingEmail.value = false
  notificationEmail.value = ''
}

async function saveEmail() {
  emailSaving.value = true
  try {
    await apiUpdateProfile({ notification_email: notificationEmail.value })
    if (auth.user) {
      auth.user.notificationEmail = notificationEmail.value
    }
    editingEmail.value = false
  } catch {
    alert('保存失败')
  } finally {
    emailSaving.value = false
  }
}

const followedTeamNames = computed(() =>
  teamStore.teams
    .filter((t) => fav.isTeamFollowed(t.id))
    .map((t) => t.name)
    .join('、')
)

// Avatar upload
const fileInput = ref<HTMLInputElement>()
const uploading = ref(false)

function triggerUpload() {
  fileInput.value?.click()
}

async function onFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  uploading.value = true
  try {
    await auth.uploadAvatar(file)
  } catch {
    alert('头像上传失败')
  } finally {
    uploading.value = false
    input.value = ''
  }
}

// Change password
const showPwdModal = ref(false)
const oldPwd = ref('')
const newPwd = ref('')
const confirmPwd = ref('')
const pwdLoading = ref(false)
const pwdError = ref('')

function openPwdModal() {
  oldPwd.value = ''
  newPwd.value = ''
  confirmPwd.value = ''
  pwdError.value = ''
  showPwdModal.value = true
}

async function submitPassword() {
  pwdError.value = ''
  if (!oldPwd.value || !newPwd.value) {
    pwdError.value = '请填写所有字段'
    return
  }
  if (newPwd.value.length < 6) {
    pwdError.value = '新密码至少 6 位'
    return
  }
  if (newPwd.value !== confirmPwd.value) {
    pwdError.value = '两次输入的密码不一致'
    return
  }
  pwdLoading.value = true
  try {
    await auth.changePassword(oldPwd.value, newPwd.value)
    showPwdModal.value = false
    alert('密码修改成功')
  } catch (err: any) {
    pwdError.value = err?.response?.data?.message || '修改失败，请检查旧密码'
  } finally {
    pwdLoading.value = false
  }
}

onMounted(() => {
  teamStore.fetchTeams()
  if (auth.isLoggedIn) {
    fav.fetchFavoriteTeams()
    fav.fetchFavoriteMatches()
    reminder.fetchReminders()
  }
})
</script>

<template>
  <div>
    <article class="card profile-card">
      <div class="avatar-wrap" @click="triggerUpload">
        <template v-if="auth.user?.avatar && auth.user.avatar.startsWith('/')">
          <img class="avatar-img" :src="auth.user.avatar" alt="头像" />
        </template>
        <template v-else>
          <div class="avatar-text">{{ auth.user?.avatar || 'U' }}</div>
        </template>
        <div class="avatar-overlay">
          <span v-if="uploading" class="material-symbols-outlined spinning">progress_activity</span>
          <span v-else class="material-symbols-outlined">photo_camera</span>
        </div>
      </div>
      <input ref="fileInput" type="file" accept="image/*" hidden @change="onFileChange" />
      <div>
        <h2>{{ auth.user?.nickname || auth.user?.username || '未登录' }}</h2>
        <p>{{ settings.timezone }} · 已关注 {{ fav.followedTeamIds.length }} 支球队</p>
      </div>
    </article>

    <section class="section">
      <div class="stats-row">
        <StatCard :value="fav.followedTeamIds.length" label="关注球队" />
        <StatCard :value="fav.favoriteMatchIds.length" label="收藏比赛" />
        <StatCard :value="reminder.count" label="比赛提醒" />
      </div>
    </section>

    <section class="section">
      <div class="card settings-list">
        <div class="setting-item"><b>我的关注</b><span>{{ followedTeamNames }}</span></div>
        <div class="setting-item"><b>我的提醒</b><span>{{ reminder.count }} 个待发送</span></div>
        <div class="setting-item">
          <b>默认提醒渠道</b>
          <select class="channel-select" :value="settings.defaultReminderChannel" @change="settings.setDefaultReminderChannel(($event.target as HTMLSelectElement).value)">
            <option value="site">站内通知</option>
            <option value="email">邮件通知</option>
          </select>
        </div>
        <div class="setting-item" v-if="!editingEmail">
          <b>通知邮箱</b>
          <span @click="startEditEmail" style="cursor:pointer">
            {{ auth.user?.notificationEmail || auth.user?.email || '未设置' }} <span class="muted-arrow">›</span>
          </span>
        </div>
        <div class="setting-item email-edit-row" v-else>
          <div class="email-edit-form">
            <input v-model="notificationEmail" type="email" placeholder="输入通知邮箱" class="email-input" />
            <div class="email-edit-actions">
              <button class="pill-btn" @click="cancelEditEmail">取消</button>
              <button class="pill-btn primary" :disabled="emailSaving" @click="saveEmail">
                {{ emailSaving ? '保存中...' : '保存' }}
              </button>
            </div>
          </div>
        </div>
        <div class="setting-item"><b>时区</b><span>{{ settings.timezone }}</span></div>
        <div class="setting-item">
          <b>深色模式</b>
          <button class="pill-btn" @click="settings.toggleTheme">切换</button>
        </div>
        <div class="setting-item"><b>语言</b><span>简体中文</span></div>
        <div class="setting-item" @click="openPwdModal" style="cursor:pointer">
          <b>修改密码</b><span class="muted-arrow">›</span>
        </div>
      </div>
    </section>

    <!-- Password Modal -->
    <Teleport to="body">
      <div v-if="showPwdModal" class="modal-mask" @click.self="showPwdModal = false">
        <div class="modal-box">
          <h3>修改密码</h3>
          <div class="form-group">
            <label>旧密码</label>
            <input v-model="oldPwd" type="password" placeholder="请输入旧密码" />
          </div>
          <div class="form-group">
            <label>新密码</label>
            <input v-model="newPwd" type="password" placeholder="至少 6 位" />
          </div>
          <div class="form-group">
            <label>确认新密码</label>
            <input v-model="confirmPwd" type="password" placeholder="再次输入新密码" />
          </div>
          <p v-if="pwdError" class="form-error">{{ pwdError }}</p>
          <div class="modal-actions">
            <button class="pill-btn" @click="showPwdModal = false">取消</button>
            <button class="pill-btn primary" :disabled="pwdLoading" @click="submitPassword">
              {{ pwdLoading ? '提交中...' : '确认修改' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
.profile-card {
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 13px;
}

.avatar-wrap {
  position: relative;
  width: 54px;
  height: 54px;
  flex-shrink: 0;
  cursor: pointer;
  border-radius: 18px;
  overflow: hidden;
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.avatar-text {
  width: 100%;
  height: 100%;
  display: grid;
  place-items: center;
  color: #fff;
  font-size: 22px;
  font-weight: 800;
  background: linear-gradient(145deg, var(--primary), var(--secondary));
}

.avatar-overlay {
  position: absolute;
  inset: 0;
  display: grid;
  place-items: center;
  background: rgba(0, 0, 0, 0.45);
  color: #fff;
  font-size: 22px;
  opacity: 0;
  transition: opacity 0.2s;
}

.avatar-wrap:hover .avatar-overlay {
  opacity: 1;
}

.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.profile-card h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 750;
}

.profile-card p {
  margin: 5px 0 0;
  color: var(--muted);
  font-size: 13px;
}

.section {
  margin-top: 18px;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
}

.settings-list {
  overflow: hidden;
}

.setting-item {
  min-height: 58px;
  padding: 0 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--line);
}

.setting-item:last-child {
  border-bottom: 0;
}

.setting-item span {
  color: var(--muted);
  font-size: 13px;
}

.muted-arrow {
  font-size: 20px;
  color: var(--muted);
}

/* Modal */
.modal-mask {
  position: fixed;
  inset: 0;
  z-index: 1000;
  display: grid;
  place-items: center;
  background: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(4px);
}

.modal-box {
  width: min(90vw, 380px);
  padding: 24px;
  border-radius: 20px;
  background: var(--card);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
}

.modal-box h3 {
  margin: 0 0 18px;
  font-size: 18px;
  font-weight: 800;
}

.form-group {
  margin-bottom: 14px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-size: 13px;
  font-weight: 600;
}

.form-group input {
  width: 100%;
  padding: 10px 14px;
  border: 1px solid var(--line);
  border-radius: 12px;
  font-size: 14px;
  background: var(--card-soft);
  outline: none;
  box-sizing: border-box;
}

.form-group input:focus {
  border-color: var(--primary);
}

.form-error {
  margin: 0 0 12px;
  color: #dc3545;
  font-size: 13px;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 18px;
}

.pill-btn {
  padding: 8px 18px;
  border: 1px solid var(--line);
  border-radius: 999px;
  background: var(--card);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
}

.channel-select {
  padding: 6px 12px;
  border: 1px solid var(--line);
  border-radius: 10px;
  background: var(--card-soft);
  font-size: 13px;
  color: inherit;
  outline: none;
  cursor: pointer;
}

.email-edit-row {
  flex-direction: column;
  align-items: stretch;
  padding: 12px 16px !important;
}

.email-edit-form {
  width: 100%;
}

.email-input {
  width: 100%;
  padding: 10px 14px;
  border: 1px solid var(--line);
  border-radius: 12px;
  font-size: 14px;
  background: var(--card-soft);
  outline: none;
  box-sizing: border-box;
  margin-bottom: 10px;
}

.email-input:focus {
  border-color: var(--primary);
}

.email-edit-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.pill-btn.primary {
  color: #fff;
  background: var(--primary);
  border-color: var(--primary);
}

.pill-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
