<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/useAuthStore'

const router = useRouter()
const auth = useAuthStore()

const isRegister = ref(false)
const username = ref('')
const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function handleSubmit() {
  error.value = ''
  loading.value = true
  try {
    if (isRegister.value) {
      await auth.register(username.value, email.value, password.value)
    } else {
      await auth.login(email.value, password.value)
    }
    router.push('/')
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
          <label>邮箱</label>
          <input v-model="email" type="email" placeholder="请输入邮箱" required />
        </div>
        <div class="field">
          <label>密码</label>
          <input v-model="password" type="password" placeholder="请输入密码" required />
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
      <p class="hint">测试账号：admin@worldcup.local / admin123456</p>
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
