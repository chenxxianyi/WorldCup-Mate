<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ThemeIcon from '@/components/theme/ThemeIcon.vue'
import ThemeMatchCard from '@/components/theme/ThemeMatchCard.vue'
import { matchToThemeMatch } from '@/data/themeAdapters'
import { useAuthStore } from '@/stores/useAuthStore'
import { useFavoriteStore } from '@/stores/useFavoriteStore'
import { useLeagueThemeStore } from '@/stores/useLeagueThemeStore'
import { useNotificationStore } from '@/stores/useNotificationStore'
import { useReminderStore } from '@/stores/useReminderStore'

const route = useRoute()
const router = useRouter()
const theme = useLeagueThemeStore()
const auth = useAuthStore()
const favorites = useFavoriteStore()
const reminders = useReminderStore()
const notifications = useNotificationStore()
const profileForm = reactive({ timezone: 'Asia/Shanghai', language: 'zh-CN', notification_email: '' })
const oldPassword = ref('')
const newPassword = ref('')
const savingProfile = ref(false)
const changingPassword = ref(false)
const uploadingAvatar = ref(false)
const recentFavorite = computed(() => favorites.favoriteMatches[0] ? matchToThemeMatch(favorites.favoriteMatches[0]) : null)

onMounted(() => {
  if (auth.isLoggedIn) notifications.fetchNotifications()
})

watch(() => auth.user, (user) => {
  if (!user) return
  profileForm.timezone = user.timezone || 'Asia/Shanghai'
  profileForm.language = user.language || 'zh-CN'
  profileForm.notification_email = user.notificationEmail || user.email
}, { immediate: true })

function login() {
  router.push({ path: '/login', query: { redirect: route.fullPath } })
}

function logout() {
  auth.logout()
  theme.showToast('已退出登录')
  router.push('/')
}

async function saveProfile() {
  savingProfile.value = true
  try {
    await auth.updateProfile(profileForm)
    theme.showToast('个人设置已保存')
  } catch (reason) {
    theme.showToast(reason instanceof Error ? reason.message : '保存失败，请稍后重试')
  } finally {
    savingProfile.value = false
  }
}

async function uploadAvatar(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  // Client-side guard (server enforces the same limit, SEC-06).
  if (file.size > 5 * 1024 * 1024) {
    theme.showToast('文件大小不能超过 5MB')
    input.value = ''
    return
  }
  uploadingAvatar.value = true
  try {
    await auth.uploadAvatar(file)
    theme.showToast('头像已更新')
  } catch (reason) {
    theme.showToast(reason instanceof Error ? reason.message : '头像上传失败')
  } finally {
    uploadingAvatar.value = false
    input.value = ''
  }
}

async function changePassword() {
  // SEC-03: client-side hint; the server-side policy is authoritative.
  const pwd = newPassword.value
  if (pwd.length < 8 || !/[A-Za-z]/.test(pwd) || !/[0-9]/.test(pwd)) {
    theme.showToast('密码至少 8 位，且需同时包含字母和数字')
    return
  }
  changingPassword.value = true
  try {
    await auth.changePassword(oldPassword.value, newPassword.value)
    oldPassword.value = ''
    newPassword.value = ''
    theme.showToast('密码修改成功')
  } catch (reason) {
    theme.showToast(reason instanceof Error ? reason.message : '密码修改失败')
  } finally {
    changingPassword.value = false
  }
}
</script>

<template>
  <div class="page-view">
    <header class="page-heading">
      <div><p class="eyebrow">MY MATCHDAY</p><h1>我的</h1></div>
      <button v-if="auth.isLoggedIn" class="secondary-button" type="button" @click="logout">退出登录</button>
      <button v-else class="primary-button" type="button" @click="login">登录账号</button>
    </header>

    <article v-if="!auth.isLoggedIn" class="login-required-card">
      <span class="profile-avatar">M</span>
      <div><p class="eyebrow">PERSONAL MATCHDAY</p><h2>登录后保存你的比赛世界</h2><p>跨设备同步关注球队、收藏比赛、开球提醒和站内通知。</p></div>
      <button class="primary-button" type="button" @click="login">登录并继续 <ThemeIcon name="arrow" /></button>
    </article>

    <div v-else class="profile-layout">
      <div class="stack">
        <article class="card profile-card">
          <img v-if="auth.user?.avatar?.startsWith('/')" class="profile-avatar image" :src="auth.user.avatar" alt="用户头像" />
          <span v-else class="profile-avatar">{{ auth.user?.avatar || auth.user?.username?.charAt(0).toUpperCase() || 'M' }}</span>
          <h2>{{ auth.user?.nickname || auth.user?.username }}</h2><p>{{ auth.user?.email }}</p>
          <label class="avatar-upload secondary-button">{{ uploadingAvatar ? '上传中…' : '更换头像' }}<input type="file" accept="image/png,image/jpeg,image/gif,image/webp" :disabled="uploadingAvatar" @change="uploadAvatar" /></label>
          <div class="profile-stats"><span class="profile-stat"><strong>{{ favorites.followedTeamIds.length }}</strong><span>关注球队</span></span><span class="profile-stat"><strong>{{ favorites.favoriteMatchIds.length }}</strong><span>收藏比赛</span></span><span class="profile-stat"><strong>{{ reminders.count }}</strong><span>比赛提醒</span></span></div>
        </article>

        <details class="card account-panel">
          <summary>账号安全 <span>修改登录密码</span></summary>
          <form class="account-form" @submit.prevent="changePassword">
            <label class="field-label">当前密码<div class="field-control"><ThemeIcon name="lock" /><input v-model="oldPassword" type="password" autocomplete="current-password" required /></div></label>
            <label class="field-label">新密码<div class="field-control"><ThemeIcon name="lock" /><input v-model="newPassword" type="password" autocomplete="new-password" minlength="8" required /></div><small class="field-hint">至少 8 位，需包含字母和数字</small></label>
            <button class="primary-button full-button" type="submit" :disabled="changingPassword">{{ changingPassword ? '提交中…' : '修改密码' }}</button>
          </form>
        </details>
      </div>

      <div class="stack">
        <form class="card profile-settings-form" @submit.prevent="saveProfile">
          <div class="setting-line"><span class="setting-label"><strong>当前赛事</strong><span>首页与数据范围</span></span><button class="setting-value text-button" type="button" @click="theme.competitionDialogOpen = true">{{ theme.current.name }} · {{ theme.current.season }}</button></div>
          <label class="setting-line"><span class="setting-label"><strong>默认时区</strong><span>比赛时间自动换算</span></span><select v-model="profileForm.timezone" class="setting-select"><option value="Asia/Shanghai">北京时间</option><option value="Europe/London">伦敦时间</option><option value="Europe/Madrid">马德里时间</option><option value="America/New_York">纽约时间</option></select></label>
          <label class="setting-line"><span class="setting-label"><strong>界面语言</strong><span>账号首选语言</span></span><select v-model="profileForm.language" class="setting-select"><option value="zh-CN">简体中文</option><option value="en-US">English</option></select></label>
          <label class="setting-line setting-email"><span class="setting-label"><strong>通知邮箱</strong><span>邮件提醒接收地址</span></span><input v-model="profileForm.notification_email" class="setting-input" type="email" required /></label>
          <div class="setting-line"><span class="setting-label"><strong>外观模式</strong><span>与赛事主题自由组合</span></span><button class="text-button" type="button" @click="theme.toggleTheme">{{ theme.settings.theme === 'dark' ? '深色' : '浅色' }}</button></div>
          <button class="primary-button settings-submit" type="submit" :disabled="savingProfile">{{ savingProfile ? '保存中…' : '保存个人设置' }}</button>
        </form>

        <section><div class="section-heading"><div><p class="eyebrow">SAVED</p><h2>最近收藏</h2></div><button class="text-link" type="button" @click="router.push('/schedule')">查看赛程</button></div><ThemeMatchCard v-if="recentFavorite" :match="recentFavorite" /><article v-else class="card empty-mini">还没有收藏比赛</article></section>

        <section><div class="section-heading"><div><p class="eyebrow">NOTIFICATIONS</p><h2>最近通知</h2></div><button v-if="notifications.unreadCount" class="text-link" type="button" @click="notifications.markAllRead">全部已读</button></div>
          <div v-if="notifications.notifications.length" class="card notification-list"><button v-for="item in notifications.notifications.slice(0, 5)" :key="item.id" class="notification-item" :class="{ unread: !item.is_read }" type="button" @click="notifications.markRead(item.id)"><i /><span><strong>{{ item.title }}</strong><small>{{ item.content }}</small></span><time>{{ new Date(item.created_at).toLocaleDateString('zh-CN') }}</time></button></div>
          <article v-else class="card empty-mini">暂无通知</article>
        </section>
      </div>
    </div>
  </div>
</template>
