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
  width: 100%;
  min-height: 58px;
  display: grid;
  grid-template-columns: 22px minmax(0, 1fr);
  align-items: center;
  column-gap: 10px;
  row-gap: 2px;
  padding: 11px 12px;
  border: 1px solid var(--line);
  border-radius: var(--radius-md);
  text-align: left;
  color: var(--text);
  background: var(--card);
  box-shadow: none;
  transition: border-color 160ms ease-out, color 160ms ease-out, background 160ms ease-out;
}

.prompt-card:focus-visible {
  outline: 2px solid color-mix(in srgb, var(--blue) 55%, transparent);
  outline-offset: 2px;
}

.prompt-card:hover {
  color: var(--blue);
  border-color: color-mix(in srgb, var(--blue) 28%, var(--line));
  background: color-mix(in srgb, var(--blue) 5%, var(--card));
}

.prompt-card .material-symbols-outlined {
  grid-row: 1 / span 2;
  flex: 0 0 auto;
  color: var(--muted);
  font-size: 20px;
}

.prompt-card strong {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 14px;
  font-weight: 700;
  line-height: 1.35;
}

.prompt-card small {
  color: var(--muted);
  line-height: 1.4;
  font-size: 12px;
}
</style>
