<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import NotificationList from '@/components/common/NotificationList.vue'
import StatCard from '@/components/common/StatCard.vue'
import ChangePasswordModal from '@/components/common/ChangePasswordModal.vue'
import { useSettingStore } from '@/stores/useSettingStore'
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import { useReminderStore } from '@/stores/useReminderStore'
import { useTeamStore } from '@/stores/useTeamStore'
import { useAuthStore } from '@/stores/useAuthStore'
import { apiUpdateProfile } from '@/api/auth'

const settings = useSettingStore()
const router = useRouter()
const fav = useFavoriteStore()
const reminder = useReminderStore()
const teamStore = useTeamStore()
const auth = useAuthStore()

// Notification email editing
const editingEmail = ref(false)
const notificationEmail = ref('')
const emailSaving = ref(false)

function startEditEmail() {
  if (!auth.isLoggedIn) {
    router.replace('/login').catch(() => {})
    return
  }
  notificationEmail.value = auth.user?.notificationEmail || ''
  editingEmail.value = true
}

function cancelEditEmail() {
  editingEmail.value = false
  notificationEmail.value = ''
}

async function saveEmail() {
  if (!auth.isLoggedIn) {
    cancelEditEmail()
    router.replace('/login').catch(() => {})
    return
  }
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

function goProfileDetail(path: string) {
  router.push(auth.isLoggedIn ? path : '/login')
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
const avatarTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp']
const maxAvatarSize = 5 * 1024 * 1024

function triggerUpload() {
  if (!auth.isLoggedIn) {
    router.replace('/login').catch(() => {})
    return
  }
  fileInput.value?.click()
}

async function onFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  if (!auth.isLoggedIn) {
    input.value = ''
    router.replace('/login').catch(() => {})
    return
  }
  if (!avatarTypes.includes(file.type)) {
    alert('仅支持 JPG、PNG、GIF、WebP 图片')
    input.value = ''
    return
  }
  if (file.size > maxAvatarSize) {
    alert('头像不能超过 5MB')
    input.value = ''
    return
  }
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

function openPwdModal() {
  if (!auth.isLoggedIn) {
    router.replace('/login').catch(() => {})
    return
  }
  showPwdModal.value = true
}

onMounted(() => {
  teamStore.fetchTeams()
  if (auth.isLoggedIn) {
    fav.fetchFavoriteTeams()
    fav.fetchFavoriteMatches()
    reminder.fetchReminders()
  }
})

watch(
  () => auth.isLoggedIn,
  (loggedIn) => {
    if (loggedIn) return
    editingEmail.value = false
    emailSaving.value = false
    uploading.value = false
    showPwdModal.value = false
  },
)
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
      <input ref="fileInput" type="file" accept="image/jpeg,image/png,image/gif,image/webp" hidden @change="onFileChange" />
      <div>
        <h2 :class="{ clickable: !auth.isLoggedIn }" @click="!auth.isLoggedIn && router.push('/login')">{{ auth.isLoggedIn ? (auth.user?.nickname || auth.user?.username) : '未登录 · 点击登录' }}</h2>
        <p>{{ settings.timezone }} · 已关注 {{ fav.followedTeamIds.length }} 支球队</p>
      </div>
    </article>

    <section class="section">
      <div class="stats-row">
        <StatCard
          :value="fav.followedTeamIds.length"
          label="关注球队"
          clickable
          @click="goProfileDetail('/profile/favorite-teams')"
        />
        <StatCard
          :value="fav.favoriteMatchIds.length"
          label="收藏比赛"
          clickable
          @click="goProfileDetail('/profile/favorite-matches')"
        />
        <StatCard
          :value="reminder.count"
          label="比赛提醒"
          clickable
          @click="goProfileDetail('/profile/reminders')"
        />
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

    <NotificationList v-if="auth.isLoggedIn" />

    <ChangePasswordModal v-model="showPwdModal" />
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

.profile-card h2.clickable {
  color: var(--primary);
  cursor: pointer;
  text-decoration: underline;
  text-underline-offset: 4px;
}
</style>
