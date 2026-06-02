<script setup lang="ts">
defineProps<{
  value: string | number
  label: string
}>()
</script>

<template>
  <div class="card stat-card">
    <strong>{{ value }}</strong>
    <span>{{ label }}</span>
  </div>
</template>

<style scoped>
.stat-card {
  padding: 15px;
}

.stat-card strong {
  display: block;
  font-size: 25px;
  line-height: 1;
  font-weight: 850;
}

.stat-card span {
  display: block;
  margin-top: 6px;
  color: var(--muted);
  font-size: 12px;
}
</style>
