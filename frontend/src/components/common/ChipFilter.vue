<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'

const props = defineProps<{
  options: string[]
  modelValue: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const chipsEl = ref<HTMLElement | null>(null)
const hasOverflow = ref(false)
const canMovePrevious = ref(false)
const canMoveNext = ref(false)
let resizeObserver: ResizeObserver | null = null

function updateNavState() {
  const el = chipsEl.value
  if (!el) {
    hasOverflow.value = false
    canMovePrevious.value = false
    canMoveNext.value = false
    return
  }

  const maxScrollLeft = el.scrollWidth - el.clientWidth
  const activeIndex = props.options.indexOf(props.modelValue)
  hasOverflow.value = maxScrollLeft > 1
  canMovePrevious.value = activeIndex > 0
  canMoveNext.value = activeIndex >= 0 && activeIndex < props.options.length - 1
}

function scrollActiveIntoView() {
  const active = chipsEl.value?.querySelector('.chip.active') as HTMLElement | null
  active?.scrollIntoView({ behavior: 'smooth', block: 'nearest', inline: 'center' })
}

function moveSelection(direction: 'previous' | 'next') {
  const activeIndex = props.options.indexOf(props.modelValue)
  if (activeIndex < 0) return

  const nextIndex = direction === 'previous' ? activeIndex - 1 : activeIndex + 1
  const nextValue = props.options[nextIndex]
  if (!nextValue) return

  emit('update:modelValue', nextValue)
}

watch(
  () => props.modelValue,
  async () => {
    await nextTick()
    scrollActiveIntoView()
    updateNavState()
  },
)

watch(
  () => props.options,
  async () => {
    await nextTick()
    updateNavState()
  },
)

onMounted(async () => {
  await nextTick()
  updateNavState()
  window.addEventListener('resize', updateNavState)
  if (chipsEl.value) {
    resizeObserver = new ResizeObserver(updateNavState)
    resizeObserver.observe(chipsEl.value)
  }
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', updateNavState)
  resizeObserver?.disconnect()
})
</script>

<template>
  <div class="chip-filter" :class="{ 'has-overflow': hasOverflow }">
    <button
      v-if="hasOverflow"
      class="chip-scroll-btn left"
      type="button"
      aria-label="切换到上一个筛选项"
      :disabled="!canMovePrevious"
      @click="moveSelection('previous')"
    >
      <span class="material-symbols-outlined" aria-hidden="true">chevron_left</span>
    </button>

    <div ref="chipsEl" class="chips" @scroll="updateNavState">
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

    <button
      v-if="hasOverflow"
      class="chip-scroll-btn right"
      type="button"
      aria-label="切换到下一个筛选项"
      :disabled="!canMoveNext"
      @click="moveSelection('next')"
    >
      <span class="material-symbols-outlined" aria-hidden="true">chevron_right</span>
    </button>
  </div>
</template>

<style scoped>
.chip-filter {
  position: relative;
  min-width: 0;
}

.chips {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  padding: 2px 2px 12px;
  scroll-behavior: smooth;
}

.chip-filter.has-overflow .chips {
  padding-inline: 44px;
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

.chip-scroll-btn {
  position: absolute;
  top: 2px;
  z-index: 2;
  width: 38px;
  height: 38px;
  display: grid;
  place-items: center;
  border: 1px solid var(--line);
  border-radius: 999px;
  color: var(--text);
  background: color-mix(in srgb, var(--card) 96%, transparent);
  box-shadow: 0 8px 20px rgba(17, 17, 17, 0.08);
  backdrop-filter: blur(10px);
  transition: color 180ms ease, background 180ms ease, opacity 180ms ease;
}

.chip-scroll-btn.left {
  left: 0;
}

.chip-scroll-btn.right {
  right: 0;
}

.chip-scroll-btn:hover:not(:disabled),
.chip-scroll-btn:focus-visible {
  color: var(--primary);
  background: var(--card);
  outline: none;
}

.chip-scroll-btn:disabled {
  cursor: default;
  opacity: 0.38;
}

.chip-scroll-btn .material-symbols-outlined {
  font-size: 23px;
}
</style>
