<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/useAuthStore'

const auth = useAuthStore()
const router = useRouter()
const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function submit() {
  if (!email.value || !password.value) {
    error.value = '请输入邮箱和密码'
    return
  }
  loading.value = true
  error.value = ''
  try {
    await auth.adminLogin(email.value, password.value)
    router.push('/admin')
  } catch (err: any) {
    error.value = err?.message || '登录失败'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="admin-login-page">
    <div class="admin-login-card">
      <div class="brand">
        <div class="brand-mark">
          WM
        </div>
        <span>WorldCup Mate Admin</span>
      </div>
      <input
        v-model="email"
        class="admin-login-input"
        type="email"
        placeholder="管理员邮箱"
        autocomplete="username"
      >
      <input
        v-model="password"
        class="admin-login-input"
        type="password"
        placeholder="密码"
        autocomplete="current-password"
        @keyup.enter="submit"
      >
      <p
        v-if="error"
        class="login-error"
        role="alert"
      >
        {{ error }}
      </p>
      <button
        class="pill-btn primary admin-login-btn"
        :disabled="loading"
        @click="submit"
      >
        {{ loading ? '登录中...' : '管理员登录' }}
      </button>
      <a
        class="back-link"
        @click="router.push('/')"
      >← 返回前台</a>
    </div>
  </div>
</template>

<style scoped>
.admin-login-page {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 24px;
  background: var(--bg);
}

.admin-login-card {
  width: 100%;
  max-width: 380px;
  display: grid;
  gap: 12px;
  padding: 28px;
  border-radius: 20px;
  background: var(--card);
  box-shadow: 0 24px 48px rgba(0, 0, 0, 0.1);
}

.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
  font-size: 16px;
  font-weight: 800;
}

.brand-mark {
  width: 38px;
  height: 38px;
  display: grid;
  place-items: center;
  border-radius: 12px;
  color: #fff;
  font-weight: 800;
  font-size: 14px;
  background: linear-gradient(145deg, var(--primary), var(--secondary));
}

.admin-login-input {
  width: 100%;
  height: 44px;
  padding: 0 14px;
  border: 1px solid var(--line);
  border-radius: 12px;
  color: var(--text);
  background: var(--card-soft);
  font-size: 14px;
  outline: none;
}

.login-error {
  margin: 0;
  color: var(--primary);
  font-size: 13px;
}

.admin-login-btn {
  width: 100%;
}

.back-link {
  text-align: center;
  color: var(--muted);
  font-size: 13px;
  cursor: pointer;
}
</style>
