<script setup lang="ts">
defineProps<{
  title: string
  description?: string
  icon?: string
}>()

defineEmits<{
  (e: 'select'): void
}>()
</script>

<template>
  <button class="prompt-card" type="button" @click="$emit('select')">
    <span class="material-symbols-outlined">{{ icon || 'auto_awesome' }}</span>
    <strong>{{ title }}</strong>
    <small v-if="description">{{ description }}</small>
  </button>
</template>

<style scoped>
.prompt-card {
  min-height: 86px;
  display: grid;
  align-content: start;
  gap: 6px;
  padding: 13px;
  border: 1px solid var(--line);
  border-radius: var(--radius-lg);
  text-align: left;
  color: var(--text);
  background: var(--card);
  box-shadow: var(--shadow);
}

.prompt-card .material-symbols-outlined {
  color: var(--primary);
  font-size: 22px;
}

.prompt-card strong {
  font-size: 14px;
}

.prompt-card small {
  color: var(--muted);
  line-height: 1.4;
}
</style>
