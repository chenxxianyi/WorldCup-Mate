<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'

const startY = ref(0)
const pullDistance = ref(0)
const touching = ref(false)
const refreshing = ref(false)

const maxPull = 86
const triggerPull = 58

const clampedPull = computed(() => Math.min(pullDistance.value, maxPull))
const progress = computed(() => Math.min(clampedPull.value / triggerPull, 1))
const ready = computed(() => progress.value >= 1)
const visible = computed(() => clampedPull.value > 2 || refreshing.value)

const indicatorStyle = computed(() => ({
  opacity: visible.value ? 1 : 0,
  transform: `translate3d(-50%, ${refreshing.value ? 18 : clampedPull.value * 0.55}px, 0) scale(${0.86 + progress.value * 0.14})`,
}))

const refreshIconStyle = computed(() => ({
  transform: refreshing.value ? undefined : `rotate(${progress.value * 180}deg)`,
}))

function canPull() {
  return window.scrollY <= 0 && document.documentElement.scrollTop <= 0 && document.body.scrollTop <= 0
}

function onTouchStart(event: TouchEvent) {
  if (refreshing.value || event.touches.length !== 1 || !canPull()) return
  touching.value = true
  startY.value = event.touches[0].clientY
  pullDistance.value = 0
}

function onTouchMove(event: TouchEvent) {
  if (!touching.value || refreshing.value || event.touches.length !== 1) return

  const delta = event.touches[0].clientY - startY.value
  if (delta <= 0) {
    pullDistance.value = 0
    return
  }

  if (!canPull()) {
    touching.value = false
    pullDistance.value = 0
    return
  }

  event.preventDefault()
  pullDistance.value = Math.min(delta * 0.58, maxPull)
}

function onTouchEnd() {
  if (!touching.value || refreshing.value) return
  touching.value = false

  if (pullDistance.value >= triggerPull) {
    refreshing.value = true
    pullDistance.value = triggerPull
    window.setTimeout(() => {
      window.location.reload()
    }, 260)
    return
  }

  pullDistance.value = 0
}

onMounted(() => {
  window.addEventListener('touchstart', onTouchStart, { passive: true })
  window.addEventListener('touchmove', onTouchMove, { passive: false })
  window.addEventListener('touchend', onTouchEnd, { passive: true })
  window.addEventListener('touchcancel', onTouchEnd, { passive: true })
})

onBeforeUnmount(() => {
  window.removeEventListener('touchstart', onTouchStart)
  window.removeEventListener('touchmove', onTouchMove)
  window.removeEventListener('touchend', onTouchEnd)
  window.removeEventListener('touchcancel', onTouchEnd)
})
</script>

<template>
  <div class="pull-refresh" :class="{ visible, ready, refreshing }" :style="indicatorStyle" aria-live="polite">
    <span class="material-symbols-outlined refresh-icon" :style="refreshIconStyle" aria-hidden="true">refresh</span>
    <span class="refresh-text">{{ refreshing ? '正在刷新' : ready ? '松开刷新' : '下拉刷新' }}</span>
  </div>
</template>

<style scoped>
.pull-refresh {
  position: fixed;
  top: 0;
  left: 50%;
  z-index: 80;
  min-width: 96px;
  height: 34px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 5px;
  padding: 0 10px;
  border: 1px solid var(--line);
  border-radius: 999px;
  color: var(--muted);
  background: color-mix(in srgb, var(--card) 94%, transparent);
  pointer-events: none;
  transition: opacity 140ms ease-out, transform 180ms cubic-bezier(0.2, 0.8, 0.2, 1), color 160ms ease-out, border-color 160ms ease-out, background 160ms ease-out;
  will-change: transform, opacity;
}

.pull-refresh.ready,
.pull-refresh.refreshing {
  color: var(--primary);
  border-color: color-mix(in srgb, var(--primary) 24%, var(--line));
  background: color-mix(in srgb, var(--primary) 6%, var(--card));
}

.refresh-icon {
  font-size: 18px;
}

.refresh-text {
  font-size: 12px;
  font-weight: 750;
  line-height: 1;
  white-space: nowrap;
}

.pull-refresh.refreshing .refresh-icon {
  animation: refreshSpin 760ms linear infinite;
}

@keyframes refreshSpin {
  to {
    transform: rotate(360deg);
  }
}

@media (prefers-reduced-motion: reduce) {
  .pull-refresh {
    transition: opacity 120ms ease-out;
  }

  .pull-refresh.refreshing .refresh-icon {
    animation: none;
  }
}
</style>
