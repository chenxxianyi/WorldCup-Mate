<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/useAuthStore'
import { LOGOUT_QUERY_KEY, LOGOUT_QUERY_VALUE, clearAuthStorage } from '@/utils/logout'

const router = useRouter()
const auth = useAuthStore()
const route = useRoute()

const redirectTo = (route.query.redirect as string) || '/'

const SAVED_LOGIN_KEY = 'wm-saved-login'
const AUTO_LOGIN_KEY = 'wm-auto-login'

type SavedLogin = {
  account?: string
  email?: string
  password?: string
}

function getSavedLogin(): SavedLogin {
  try {
    return JSON.parse(localStorage.getItem(SAVED_LOGIN_KEY) || '{}')
  } catch {
    return {}
  }
}

const savedLogin = getSavedLogin()
const isRegister = ref(false)
const username = ref('')
const account = ref(savedLogin.account || savedLogin.email || '')
const email = ref('')
const password = ref(savedLogin.password || '')
const confirmPassword = ref('')
const showPassword = ref(false)
const rememberPassword = ref(Boolean(savedLogin.account && savedLogin.password))
const autoLogin = ref(localStorage.getItem(AUTO_LOGIN_KEY) === '1')
const error = ref('')
const loading = ref(false)
const showTestAccount = import.meta.env.DEV
const passwordMismatch = computed(() => isRegister.value && confirmPassword.value.length > 0 && password.value !== confirmPassword.value)

watch(autoLogin, (enabled) => {
  if (enabled) rememberPassword.value = true
})

watch(isRegister, () => {
  error.value = ''
  confirmPassword.value = ''
})

onMounted(() => {
  if (route.query[LOGOUT_QUERY_KEY] !== LOGOUT_QUERY_VALUE) return

  clearAuthStorage()
  auth.logout()
  router.replace({ name: 'login' }).catch(() => {})
})

async function handleSubmit() {
  error.value = ''
  if (isRegister.value && password.value !== confirmPassword.value) {
    error.value = '两次输入的密码不一致'
    return
  }

  loading.value = true
  try {
    if (isRegister.value) {
      await auth.register(username.value, email.value, password.value, confirmPassword.value)
    } else {
      await auth.login(account.value, password.value, autoLogin.value)
      if (rememberPassword.value) {
        localStorage.setItem(SAVED_LOGIN_KEY, JSON.stringify({ account: account.value, password: password.value }))
      } else {
        localStorage.removeItem(SAVED_LOGIN_KEY)
      }
      if (autoLogin.value) {
        localStorage.setItem(AUTO_LOGIN_KEY, '1')
      } else {
        localStorage.removeItem(AUTO_LOGIN_KEY)
      }
    }
    router.push(redirectTo)
  } catch (e: any) {
    error.value = e.message || '操作失败'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-page">
    <article class="card login-card">
      <h2>{{ isRegister ? '注册' : '登录' }} WorldCup Mate</h2>
      <form @submit.prevent="handleSubmit">
        <div v-if="isRegister" class="field">
          <label>用户名</label>
          <input v-model="username" placeholder="请输入用户名" required />
        </div>
        <div class="field">
          <label>{{ isRegister ? '邮箱' : '邮箱 / 用户名' }}</label>
          <input
            v-if="isRegister"
            v-model="email"
            type="email"
            placeholder="请输入邮箱"
            required
          />
          <input
            v-else
            v-model="account"
            type="text"
            placeholder="请输入邮箱或用户名"
            autocomplete="username"
            required
          />
        </div>
        <div class="field">
          <label>密码</label>
          <div class="password-input">
            <input
              v-model="password"
              :type="showPassword ? 'text' : 'password'"
              placeholder="请输入密码"
              :minlength="isRegister ? 6 : undefined"
              :autocomplete="isRegister ? 'new-password' : 'current-password'"
              required
            />
            <button
              class="password-toggle"
              type="button"
              :title="showPassword ? '隐藏密码' : '显示密码'"
              :aria-label="showPassword ? '隐藏密码' : '显示密码'"
              :aria-pressed="showPassword"
              @click="showPassword = !showPassword"
            >
              <span class="material-symbols-outlined" aria-hidden="true">{{ showPassword ? 'visibility_off' : 'visibility' }}</span>
            </button>
          </div>
        </div>
        <div v-if="isRegister" class="field">
          <label>确认密码</label>
          <input
            v-model="confirmPassword"
            type="password"
            placeholder="请再次输入密码"
            minlength="6"
            autocomplete="new-password"
            required
            :aria-invalid="passwordMismatch"
            :class="{ invalid: passwordMismatch }"
          />
          <p v-if="passwordMismatch" class="field-error">两次输入的密码不一致</p>
        </div>
        <div v-if="!isRegister" class="login-options" aria-label="登录选项">
          <label class="option-row">
            <input v-model="rememberPassword" type="checkbox" />
            <span>记住密码</span>
          </label>
          <label class="option-row">
            <input v-model="autoLogin" type="checkbox" />
            <span>自动登录</span>
          </label>
        </div>
        <p v-if="error" class="error">{{ error }}</p>
        <button class="pill-btn primary submit-btn" type="submit" :disabled="loading">
          {{ loading ? '处理中...' : (isRegister ? '注册' : '登录') }}
        </button>
      </form>
      <p class="toggle">
        {{ isRegister ? '已有账号？' : '没有账号？' }}
        <a href="#" @click.prevent="isRegister = !isRegister">{{ isRegister ? '去登录' : '去注册' }}</a>
      </p>
      <p v-if="showTestAccount" class="hint">测试账号：admin@worldcup.local / admin123456</p>
    </article>
  </div>
</template>

<style scoped>
.login-page {
  display: grid;
  place-items: center;
  min-height: 60vh;
}

.login-card {
  padding: 28px 24px;
  max-width: 400px;
  width: 100%;
}

.login-card h2 {
  margin: 0 0 20px;
  font-size: 20px;
  font-weight: 750;
  text-align: center;
}

.field {
  margin-bottom: 14px;
}

.field label {
  display: block;
  margin-bottom: 6px;
  font-size: 13px;
  font-weight: 600;
  color: var(--muted);
}

.field input {
  width: 100%;
  height: 44px;
  border: 1px solid var(--line);
  border-radius: 12px;
  padding: 0 14px;
  font-size: 14px;
  color: var(--text);
  background: var(--card-soft);
  outline: none;
}

.field input:focus {
  border-color: var(--primary);
}

.password-input {
  position: relative;
}

.password-input input {
  padding-right: 52px;
}

.password-toggle {
  position: absolute;
  top: 0;
  right: 0;
  display: grid;
  place-items: center;
  width: 48px;
  height: 44px;
  border: 0;
  color: var(--muted);
  background: transparent;
  cursor: pointer;
}

.password-toggle:focus-visible {
  outline: 2px solid var(--primary);
  outline-offset: -4px;
  border-radius: 10px;
}

.password-toggle .material-symbols-outlined {
  font-size: 21px;
  line-height: 1;
}

.field input.invalid {
  border-color: var(--live);
}

.field-error {
  margin: 6px 0 0;
  color: var(--live);
  font-size: 12px;
}

.login-options {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin: 0 0 12px;
}

.option-row {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  min-height: 28px;
  color: var(--muted);
  font-size: 13px;
  font-weight: 650;
  line-height: 1;
  white-space: nowrap;
  cursor: pointer;
}

.option-row input {
  width: 16px;
  height: 16px;
  margin: 0;
  accent-color: var(--primary);
}

.error {
  color: var(--live);
  font-size: 13px;
  margin: 8px 0;
}

.submit-btn {
  width: 100%;
  margin-top: 8px;
  border: 0;
  cursor: pointer;
  font-size: 15px;
  font-weight: 700;
}

.toggle {
  text-align: center;
  margin-top: 16px;
  font-size: 13px;
  color: var(--muted);
}

.toggle a {
  color: var(--primary);
  text-decoration: none;
  font-weight: 600;
}

.hint {
  text-align: center;
  margin-top: 12px;
  font-size: 12px;
  color: var(--weak);
}
</style>
