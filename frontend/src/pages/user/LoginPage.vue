<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ThemeIcon from '@/components/theme/ThemeIcon.vue'
import { useAuthStore } from '@/stores/useAuthStore'
import { useLeagueThemeStore } from '@/stores/useLeagueThemeStore'
import { safeRedirectPath } from '@/router'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const theme = useLeagueThemeStore()
const email = ref(import.meta.env.DEV ? 'admin@worldcup.local' : '')
const password = ref(import.meta.env.DEV ? 'admin123456' : '')
const username = ref('')
const isRegister = ref(false)
const remember = ref(true)
const loading = ref(false)
const error = ref('')

async function submit() {
  error.value = ''
  if (isRegister.value) {
    // SEC-03: client-side hint; the server-side policy is authoritative.
    const pwd = password.value
    if (pwd.length < 8 || !/[A-Za-z]/.test(pwd) || !/[0-9]/.test(pwd)) {
      error.value = '密码至少 8 位，且需同时包含字母和数字'
      return
    }
  }
  loading.value = true
  try {
    if (isRegister.value) await auth.register(username.value, email.value, password.value, remember.value)
    else await auth.login(email.value, password.value, remember.value)
    theme.showToast(isRegister.value ? '注册成功，欢迎加入' : '登录成功，欢迎回来')
    // SEC-05: only in-app absolute paths are allowed as redirect targets.
    router.push(safeRedirectPath(route.query.redirect, '/profile'))
  } catch (reason: unknown) {
    error.value = reason instanceof Error ? reason.message : '登录失败，请检查账号或服务状态'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="page-view login-wrap">
    <section class="login-panel">
      <div class="login-visual"><span class="hero-mark">{{ theme.current.mark }}</span><div><h1>你的比赛日，应该更简单。</h1><p>登录后同步关注球队、收藏比赛与开球提醒，在六个赛事世界之间无缝切换。</p></div><small>{{ theme.current.name }} · {{ theme.current.season }}</small></div>
      <form class="login-form" @submit.prevent="submit">
        <p class="eyebrow">{{ isRegister ? 'JOIN MATCHDAY' : 'WELCOME BACK' }}</p><h2>{{ isRegister ? '注册' : '登录' }} WorldCup Mate</h2><p>登录后同步你的关注、收藏和提醒设置。</p>
        <div class="form-stack">
          <label v-if="isRegister" class="field-label">用户名<div class="field-control"><ThemeIcon name="user" /><input v-model="username" type="text" autocomplete="username" minlength="2" required /></div></label>
          <label class="field-label">邮箱<div class="field-control"><ThemeIcon name="mail" /><input v-model="email" type="email" autocomplete="email" required /></div></label>
          <label class="field-label">密码<div class="field-control"><ThemeIcon name="lock" /><input v-model="password" type="password" :autocomplete="isRegister ? 'new-password' : 'current-password'" required /></div><small v-if="isRegister" class="field-hint">至少 8 位，需包含字母和数字</small></label>
          <div class="form-row"><label><input v-model="remember" type="checkbox" /> 记住我</label><button class="text-button" type="button" @click="theme.showToast('重置密码邮件功能待接入')">忘记密码？</button></div>
          <p v-if="error" class="login-error">{{ error }}</p>
          <button class="primary-button full-button" type="submit" :disabled="loading">{{ loading ? '处理中…' : isRegister ? '注册并继续' : '登录并继续' }} <ThemeIcon name="arrow" /></button>
          <button class="text-button auth-switch" type="button" @click="isRegister = !isRegister; error = ''">{{ isRegister ? '已有账号？返回登录' : '没有账号？创建账号' }}</button>
          <button class="secondary-button full-button" type="button" @click="router.push('/profile')">返回个人中心</button>
        </div>
      </form>
    </section>
  </div>
</template>

<style scoped>
.login-error { margin: 0; color: var(--status-danger); font-size: 13px; }
button:disabled { cursor: wait; opacity: .65; }
</style>
