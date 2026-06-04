<script setup lang="ts">
import { ref } from 'vue'

defineProps<{
  loading?: boolean
  placeholder?: string
}>()

const emit = defineEmits<{
  (e: 'send', value: string): void
}>()

const text = ref('')

function submit() {
  const value = text.value.trim()
  if (!value) return
  emit('send', value)
  text.value = ''
}
</script>

<template>
  <form class="input-box" @submit.prevent="submit">
    <textarea
      v-model="text"
      rows="2"
      :placeholder="placeholder || '问问赛程、球队、规则或今天该看哪场'"
      :disabled="loading"
      @keydown.enter.exact.prevent="submit"
    ></textarea>
    <button class="send-btn" type="submit" :disabled="loading || !text.trim()">
      <span class="material-symbols-outlined">send</span>
    </button>
  </form>
</template>

<style scoped>
.input-box {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 10px;
  align-items: end;
  padding: 10px;
  border: 1px solid var(--line);
  border-radius: 18px;
  background: var(--card);
  box-shadow: var(--shadow);
}

textarea {
  width: 100%;
  min-height: 44px;
  max-height: 128px;
  resize: vertical;
  border: 0;
  outline: 0;
  color: var(--text);
  background: transparent;
  line-height: 1.5;
  font-size: 14px;
}

textarea::placeholder {
  color: var(--weak);
}

.send-btn {
  width: 42px;
  height: 42px;
  display: grid;
  place-items: center;
  border: 0;
  border-radius: 999px;
  color: #fff;
  background: var(--primary);
}

.send-btn:disabled {
  opacity: 0.42;
  cursor: not-allowed;
}
</style>
