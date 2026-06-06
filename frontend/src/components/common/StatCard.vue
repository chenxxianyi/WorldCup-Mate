<script setup lang="ts">
defineProps<{
  value: string | number
  label: string
  clickable?: boolean
}>()

const emit = defineEmits<{
  (e: 'click'): void
}>()
</script>

<template>
  <div
    class="card stat-card"
    :class="{ clickable }"
    :tabindex="clickable ? 0 : undefined"
    :role="clickable ? 'button' : undefined"
    @click="clickable && emit('click')"
    @keydown.enter="clickable && emit('click')"
    @keydown.space.prevent="clickable && emit('click')"
  >
    <strong>{{ value }}</strong>
    <span>
      {{ label }}
      <i v-if="clickable" class="material-symbols-outlined">chevron_right</i>
    </span>
  </div>
</template>

<style scoped>
.stat-card {
  padding: 15px;
  transition: border-color 160ms ease-out, background 160ms ease-out, transform 160ms ease-out;
}

.stat-card.clickable {
  cursor: pointer;
}

.stat-card.clickable:hover {
  border-color: color-mix(in srgb, var(--primary) 28%, var(--line));
  background: color-mix(in srgb, var(--primary) 4%, var(--card));
}

.stat-card.clickable:active {
  transform: scale(0.98);
}

.stat-card.clickable:focus-visible {
  outline: 2px solid var(--primary);
  outline-offset: 2px;
}

.stat-card strong {
  display: block;
  font-size: 25px;
  line-height: 1;
  font-weight: 850;
}

.stat-card span {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 4px;
  margin-top: 6px;
  color: var(--muted);
  font-size: 12px;
}

.stat-card i {
  font-size: 15px;
  line-height: 1;
}
</style>
