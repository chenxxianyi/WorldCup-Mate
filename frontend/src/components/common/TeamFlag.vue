<script setup lang="ts">
import { computed, ref, watch } from 'vue'

const props = withDefaults(defineProps<{
  value?: string | null
  alt?: string
  fallback?: string
  size?: 'sm' | 'md' | 'lg'
}>(), {
  value: '',
  alt: '',
  fallback: '',
  size: 'md',
})

const failed = ref(false)

watch(() => props.value, () => {
  failed.value = false
})

const source = computed(() => (props.value || '').trim())
const fallbackText = computed(() => {
  const text = (props.fallback || props.alt || '').trim()
  return text ? text.slice(0, 3).toUpperCase() : 'TBD'
})
const isImage = computed(() => /^(https?:\/\/|\/)/i.test(source.value))
</script>

<template>
  <span class="team-flag" :class="size" aria-hidden="true">
    <img
      v-if="isImage && !failed"
      :src="source"
      :alt="alt"
      loading="lazy"
      decoding="async"
      @error="failed = true"
    />
    <span v-else-if="source && !isImage" class="flag-text">{{ source }}</span>
    <span v-else class="flag-fallback">{{ fallbackText }}</span>
  </span>
</template>

<style scoped>
.team-flag {
  --flag-size: 34px;
  width: var(--flag-size);
  height: var(--flag-size);
  display: inline-grid;
  place-items: center;
  flex: 0 0 auto;
  overflow: hidden;
  border-radius: 999px;
  background: var(--card-soft);
  color: var(--text);
  vertical-align: middle;
}

.team-flag.sm {
  --flag-size: 24px;
}

.team-flag.lg {
  --flag-size: 54px;
}

.team-flag img {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: contain;
  padding: 5px;
}

.flag-text {
  max-width: 100%;
  overflow: hidden;
  font-size: calc(var(--flag-size) * 0.62);
  line-height: 1;
  text-align: center;
}

.flag-fallback {
  max-width: 100%;
  padding: 0 3px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--muted);
  font-size: calc(var(--flag-size) * 0.28);
  font-weight: 800;
  letter-spacing: 0;
}
</style>
