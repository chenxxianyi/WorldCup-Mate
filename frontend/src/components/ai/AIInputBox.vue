<script setup lang="ts">
import { ref } from 'vue'

defineProps<{
  loading?: boolean
  placeholder?: string
}>()

const emit = defineEmits<{
  (e: 'send', value: string): void
  (e: 'stop'): void
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
    <button
      v-if="loading"
      class="send-btn stop-btn"
      type="button"
      title="暂停生成"
      aria-label="暂停生成"
      @click="emit('stop')"
    >
      <span class="material-symbols-outlined">pause</span>
    </button>
    <button v-else class="send-btn" type="submit" :disabled="!text.trim()">
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
  border-radius: var(--radius-md);
  background: var(--card);
  box-shadow: none;
}

.input-box:focus-within {
  border-color: color-mix(in srgb, var(--blue) 36%, var(--line));
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
  width: 44px;
  height: 44px;
  display: grid;
  place-items: center;
  border: 0;
  border-radius: var(--radius-sm);
  color: #fff;
  background: var(--blue);
}

.send-btn:focus-visible {
  outline: 2px solid color-mix(in srgb, var(--blue) 55%, transparent);
  outline-offset: 2px;
}

.send-btn:disabled {
  opacity: 0.42;
  cursor: not-allowed;
}

</style>
