<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const SHOW_SCROLL_Y = 420
const IDLE_HIDE_DELAY_MS = 1500

const showBackTop = ref(false)
let hideTimer: number | undefined

function clearHideTimer() {
  if (!hideTimer) return
  window.clearTimeout(hideTimer)
  hideTimer = undefined
}

function scheduleIdleHide() {
  clearHideTimer()
  if (window.scrollY <= SHOW_SCROLL_Y) return

  hideTimer = window.setTimeout(() => {
    showBackTop.value = false
    hideTimer = undefined
  }, IDLE_HIDE_DELAY_MS)
}

function updateBackTopVisibility() {
  const shouldShow = window.scrollY > SHOW_SCROLL_Y
  showBackTop.value = shouldShow
  if (shouldShow) scheduleIdleHide()
  else clearHideTimer()
}

function backToTop() {
  clearHideTimer()
  showBackTop.value = false
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  window.addEventListener('scroll', updateBackTopVisibility, { passive: true })
  updateBackTopVisibility()
})

onBeforeUnmount(() => {
  clearHideTimer()
  window.removeEventListener('scroll', updateBackTopVisibility)
})

watch(
  () => route.fullPath,
  async () => {
    await nextTick()
    updateBackTopVisibility()
  },
)
</script>

<template>
  <Transition name="back-top">
    <button
      v-if="showBackTop"
      class="back-top-btn"
      type="button"
      aria-label="返回顶部"
      title="返回顶部"
      @click="backToTop"
    >
      <span class="material-symbols-outlined">keyboard_arrow_up</span>
    </button>
  </Transition>
</template>
