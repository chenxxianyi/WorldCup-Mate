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
  min-width: 190px;
  min-height: 44px;
  display: inline-flex;
  align-items: center;
  gap: 7px;
  padding: 9px 11px;
  border: 1px solid var(--line);
  border-radius: var(--radius-sm);
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
  flex: 0 0 auto;
  color: var(--muted);
  font-size: 18px;
}

.prompt-card strong {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 13px;
  font-weight: 700;
  line-height: 1.45;
}

.prompt-card small {
  color: var(--muted);
  line-height: 1.4;
}
</style>
