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
    :class="{ 'is-clickable': clickable }"
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
  position: relative;
  min-height: 78px;
  padding: 14px 14px 12px;
  overflow: hidden;
  transition: border-color 180ms ease-out, background 180ms ease-out, box-shadow 180ms ease-out, transform 180ms ease-out;
}

.stat-card.is-clickable {
  cursor: pointer;
}

.stat-card.is-clickable:hover {
  border-color: color-mix(in srgb, var(--blue) 26%, var(--line));
  background: color-mix(in srgb, var(--blue) 4%, var(--card));
  box-shadow: 0 12px 28px rgba(15, 23, 42, 0.08);
}

.stat-card.is-clickable:active {
  transform: translateY(1px);
}

.stat-card.is-clickable:focus-visible {
  outline: 2px solid color-mix(in srgb, var(--blue) 58%, transparent);
  outline-offset: 2px;
}

.stat-card strong {
  display: block;
  color: var(--text);
  font-size: 24px;
  line-height: 1;
  font-weight: 800;
  letter-spacing: 0;
  text-decoration: none;
}

.stat-card span {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 4px;
  margin-top: 6px;
  color: var(--muted);
  font-size: 12px;
  font-weight: 600;
  line-height: 1.2;
  text-decoration: none;
}

.stat-card i {
  display: inline-grid;
  width: 20px;
  height: 20px;
  flex: 0 0 auto;
  place-items: center;
  border-radius: 999px;
  color: var(--blue);
  background: color-mix(in srgb, var(--blue) 8%, transparent);
  font-size: 15px;
  line-height: 1;
  transition: background 180ms ease-out, transform 180ms ease-out;
}

.stat-card.is-clickable:hover i {
  background: color-mix(in srgb, var(--blue) 13%, transparent);
  transform: translateX(1px);
}
</style>
