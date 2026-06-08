<script setup lang="ts">
import { nextTick, ref, watch } from 'vue'

const props = defineProps<{
  options: string[]
  modelValue: string
}>()

defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const chipsEl = ref<HTMLElement | null>(null)

function scrollActiveIntoView() {
  const active = chipsEl.value?.querySelector('.chip.active') as HTMLElement | null
  active?.scrollIntoView({ behavior: 'smooth', block: 'nearest', inline: 'center' })
}

watch(
  () => props.modelValue,
  async () => {
    await nextTick()
    scrollActiveIntoView()
  },
)
</script>

<template>
  <div ref="chipsEl" class="chips">
    <button
      v-for="opt in options"
      :key="opt"
      class="chip"
      :class="{ active: modelValue === opt }"
      @click="$emit('update:modelValue', opt)"
    >
      {{ opt }}
    </button>
  </div>
</template>

<style scoped>
.chips {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  padding: 2px 2px 12px;
}

.chip {
  flex: 0 0 auto;
  min-height: 38px;
  padding: 0 14px;
  border: 1px solid transparent;
  border-radius: 999px;
  color: var(--muted);
  background: var(--card-soft);
  transition: all 180ms ease-out;
  font-size: 13px;
  font-weight: 600;
  white-space: nowrap;
}

.chip.active {
  color: #fff;
  background: var(--primary);
  box-shadow: 0 10px 24px color-mix(in srgb, var(--primary) 22%, transparent);
}
</style>
