<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'

const props = defineProps<{
  hours?: number
  minutes?: number
  seconds?: number
  targetTime?: string
}>()

const h = ref('00')
const m = ref('00')
const s = ref('00')
const expired = ref(false)

let timer: ReturnType<typeof setInterval>

function calcRemaining(): number {
  if (props.targetTime) {
    const target = new Date(props.targetTime).getTime()
    const now = Date.now()
    return Math.max(0, Math.floor((target - now) / 1000))
  }
  return (props.hours || 0) * 3600 + (props.minutes || 0) * 60 + (props.seconds || 0)
}

function updateDisplay(total: number) {
  h.value = String(Math.floor(total / 3600)).padStart(2, '0')
  m.value = String(Math.floor((total % 3600) / 60)).padStart(2, '0')
  s.value = String(total % 60).padStart(2, '0')
  expired.value = total === 0
}

onMounted(() => {
  let total = calcRemaining()
  updateDisplay(total)
  timer = setInterval(() => {
    total = Math.max(0, total - 1)
    updateDisplay(total)
    if (total === 0) clearInterval(timer)
  }, 1000)
})

onUnmounted(() => {
  clearInterval(timer)
})

watch(() => props.targetTime, () => {
  clearInterval(timer)
  let total = calcRemaining()
  updateDisplay(total)
  timer = setInterval(() => {
    total = Math.max(0, total - 1)
    updateDisplay(total)
    if (total === 0) clearInterval(timer)
  }, 1000)
})
</script>

<template>
  <article class="card countdown-card">
    <div class="countdown-inner">
      <slot />
      <div class="countdown-time">
        <div class="time-box"><strong>{{ h }}</strong><span>小时</span></div>
        <div class="time-box"><strong>{{ m }}</strong><span>分钟</span></div>
        <div class="time-box"><strong>{{ s }}</strong><span>秒</span></div>
      </div>
    </div>
  </article>
</template>

<style scoped>
.countdown-card {
  overflow: hidden;
  position: relative;
  padding: 20px;
  color: #fff;
  border-color: rgba(255, 255, 255, 0.14);
  background:
    radial-gradient(circle at 100% 0%, rgba(206, 17, 38, 0.26), transparent 32%),
    radial-gradient(circle at 8% 92%, rgba(0, 104, 71, 0.42), transparent 34%),
    linear-gradient(135deg, rgba(7, 7, 7, 0.96), rgba(0, 40, 104, 0.9) 52%, rgba(22, 18, 10, 0.95)),
    #070707;
}

.countdown-card::before {
  content: '';
  position: absolute;
  width: 220px;
  height: 220px;
  right: -92px;
  top: -80px;
  border: 1px solid rgba(255, 255, 255, 0.16);
  border-radius: 50%;
}

.countdown-card::after {
  content: '';
  position: absolute;
  inset: auto -20px -60px 30%;
  height: 160px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(56, 189, 248, 0.24), transparent 68%);
}

.countdown-inner {
  position: relative;
  z-index: 1;
}

.countdown-time {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
  margin-top: 22px;
}

.time-box {
  padding: 13px 8px 11px;
  border: 1px solid rgba(255, 255, 255, 0.11);
  border-radius: 16px;
  text-align: center;
  background: rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(16px);
}

.time-box strong {
  display: block;
  font-family: var(--font-display);
  font-size: clamp(30px, 11vw, 44px);
  line-height: 1;
  font-weight: 800;
}

.time-box span {
  display: block;
  margin-top: 6px;
  color: rgba(255, 255, 255, 0.64);
  font-size: 12px;
}
</style>
