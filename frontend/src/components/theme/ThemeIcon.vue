<script setup lang="ts">
defineProps<{ name: string; label?: string }>()
</script>

<template>
  <svg
    :aria-hidden="label ? undefined : 'true'"
    :aria-label="label"
  >
    <use :href="`#i-${name}`" />
  </svg>
</template>
