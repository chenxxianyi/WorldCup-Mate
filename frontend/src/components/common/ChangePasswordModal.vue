<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/useAuthStore'

const props = defineProps<{ modelValue: boolean }>()
const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}>()

const auth = useAuthStore()
const router = useRouter()

const oldPwd = ref('')
const newPwd = ref('')
const confirmPwd = ref('')
const loading = ref(false)
const error = ref('')

function resetForm() {
  oldPwd.value = ''
  newPwd.value = ''
  confirmPwd.value = ''
  error.value = ''
  loading.value = false
}

function close() {
  emit('update:modelValue', false)
}

// Reset the form each time the modal opens so it never carries stale input.
watch(
  () => props.modelValue,
  (open) => {
    if (open) resetForm()
  },
)

async function submit() {
  error.value = ''
  if (!auth.isLoggedIn) {
    close()
    router.replace('/login').catch(() => {})
    return
  }
  if (!oldPwd.value || !newPwd.value) {
    error.value = '请填写所有字段'
    return
  }
  if (newPwd.value.length < 6) {
    error.value = '新密码至少 6 位'
    return
  }
  if (newPwd.value !== confirmPwd.value) {
    error.value = '两次输入的密码不一致'
    return
  }
  loading.value = true
  try {
    await auth.changePassword(oldPwd.value, newPwd.value)
    close()
    emit('success')
    alert('密码修改成功')
  } catch (err: any) {
    error.value = err?.response?.data?.message || '修改失败，请检查旧密码'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div v-if="modelValue" class="modal-mask" @click.self="close">
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
      <p v-if="error" class="form-error">{{ error }}</p>
      <div class="modal-actions">
        <button class="pill-btn" @click="close">取消</button>
        <button class="pill-btn primary" :disabled="loading" @click="submit">
          {{ loading ? '提交中...' : '确认修改' }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
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
