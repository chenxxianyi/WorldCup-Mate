<script setup lang="ts">
import { computed, onMounted, onUnmounted, reactive, watch } from 'vue'

const props = defineProps<{
  hours?: number
  minutes?: number
  seconds?: number
  targetTime?: string
}>()

type UnitKey = 'hours' | 'minutes' | 'seconds'

interface TimePart {
  key: UnitKey
  label: string
  value: string
  previous: string
  flipping: boolean
  flipKey: number
}

const unitKeys: UnitKey[] = ['hours', 'minutes', 'seconds']
const timeParts = reactive<Record<UnitKey, TimePart>>({
  hours: { key: 'hours', label: '小时', value: '00', previous: '00', flipping: false, flipKey: 0 },
  minutes: { key: 'minutes', label: '分钟', value: '00', previous: '00', flipping: false, flipKey: 0 },
  seconds: { key: 'seconds', label: '秒', value: '00', previous: '00', flipping: false, flipKey: 0 },
})

const displayParts = computed(() => unitKeys.map((key) => timeParts[key]))

let timer: number | undefined
const flipTimers: Partial<Record<UnitKey, number>> = {}

function calcRemaining(): number {
  if (props.targetTime) {
    const target = new Date(props.targetTime).getTime()
    if (!Number.isFinite(target)) return 0
    const now = Date.now()
    return Math.max(0, Math.floor((target - now) / 1000))
  }
  return (props.hours || 0) * 3600 + (props.minutes || 0) * 60 + (props.seconds || 0)
}

function formatTime(total: number) {
  return {
    hours: String(Math.floor(total / 3600)).padStart(2, '0'),
    minutes: String(Math.floor((total % 3600) / 60)).padStart(2, '0'),
    seconds: String(total % 60).padStart(2, '0'),
  }
}

function setPart(key: UnitKey, nextValue: string, animate: boolean) {
  const part = timeParts[key]
  if (part.value === nextValue) return

  part.previous = part.value
  part.value = nextValue

  if (!animate) {
    part.flipping = false
    return
  }

  part.flipKey += 1
  part.flipping = true
  if (flipTimers[key] !== undefined) window.clearTimeout(flipTimers[key])
  flipTimers[key] = window.setTimeout(() => {
    part.flipping = false
  }, 920)
}

function updateDisplay(total: number, animate = true) {
  const next = formatTime(total)
  setPart('hours', next.hours, animate)
  setPart('minutes', next.minutes, animate)
  setPart('seconds', next.seconds, animate)
}

function clearCountdownTimer() {
  if (timer === undefined) return
  window.clearInterval(timer)
  timer = undefined
}

function clearFlipTimers() {
  for (const key of unitKeys) {
    if (flipTimers[key] !== undefined) window.clearTimeout(flipTimers[key])
    flipTimers[key] = undefined
    timeParts[key].flipping = false
  }
}

function startCountdown() {
  clearCountdownTimer()
  clearFlipTimers()
  let total = calcRemaining()
  updateDisplay(total, false)
  if (total === 0) return

  timer = window.setInterval(() => {
    total = props.targetTime ? calcRemaining() : Math.max(0, total - 1)
    updateDisplay(total)
    if (total === 0) clearCountdownTimer()
  }, 1000)
}

function digitClass(value: string) {
  if (value.length >= 4) return 'digits-tight'
  if (value.length === 3) return 'digits-medium'
  return 'digits-normal'
}

onMounted(startCountdown)

onUnmounted(() => {
  clearCountdownTimer()
  clearFlipTimers()
})

watch(() => [props.targetTime, props.hours, props.minutes, props.seconds], startCountdown)
</script>

<template>
  <article class="card countdown-card">
    <div class="countdown-inner">
      <slot />
      <div class="countdown-time">
        <div
          v-for="part in displayParts"
          :key="part.key"
          class="time-box"
          :class="digitClass(part.value)"
        >
          <div class="flip-stack" :class="{ 'is-flipping': part.flipping }" aria-hidden="true">
            <div class="flip-panel flip-panel-top">
              <span class="flip-value">{{ part.value }}</span>
            </div>
            <div class="flip-panel flip-panel-bottom">
              <span class="flip-value">{{ part.value }}</span>
            </div>
            <template v-if="part.flipping">
              <div :key="`${part.flipKey}-top`" class="score-flip score-flip-top">
                <span class="flip-value">{{ part.previous }}</span>
              </div>
              <div :key="`${part.flipKey}-bottom`" class="score-flip score-flip-bottom">
                <span class="flip-value">{{ part.value }}</span>
              </div>
            </template>
            <span class="flip-crease"></span>
          </div>
          <span class="sr-only">{{ part.value }} {{ part.label }}</span>
          <span class="time-label" aria-hidden="true">{{ part.label }}</span>
        </div>
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
  gap: 10px;
  margin-top: 22px;
}

.time-box {
  min-width: 0;
  display: grid;
  gap: 8px;
  text-align: center;
}

.flip-stack {
  position: relative;
  height: 92px;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.13);
  border-radius: 14px;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.17), rgba(255, 255, 255, 0.05)),
    rgba(10, 25, 44, 0.62);
  box-shadow:
    0 18px 34px rgba(0, 0, 0, 0.22),
    inset 0 1px 0 rgba(255, 255, 255, 0.18),
    inset 0 -18px 30px rgba(0, 0, 0, 0.2);
  backdrop-filter: blur(16px);
  perspective: 1000px;
  transform-style: preserve-3d;
}

.flip-stack::before {
  content: '';
  position: absolute;
  left: 18px;
  right: 18px;
  top: 9px;
  z-index: 4;
  height: 2px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.18);
}

.flip-panel,
.score-flip {
  position: absolute;
  left: 0;
  right: 0;
  height: 50%;
  display: grid;
  place-items: center;
  overflow: hidden;
  border-radius: inherit;
  background:
    linear-gradient(145deg, rgba(42, 72, 116, 0.92), rgba(10, 22, 39, 0.95));
}

.flip-panel-top,
.score-flip-top {
  top: 0;
  align-items: end;
  border-radius: 14px 14px 0 0;
  z-index: 1;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.2), rgba(255, 255, 255, 0.08)),
    linear-gradient(145deg, rgba(72, 98, 135, 0.92), rgba(31, 50, 76, 0.95));
}

.flip-panel-bottom,
.score-flip-bottom {
  bottom: 0;
  align-items: start;
  border-radius: 0 0 14px 14px;
  background:
    linear-gradient(180deg, rgba(2, 8, 23, 0.34), rgba(255, 255, 255, 0.06)),
    linear-gradient(145deg, rgba(20, 38, 61, 0.95), rgba(8, 19, 34, 0.96));
}

.flip-panel {
  z-index: 1;
}

.score-flip {
  z-index: 3;
  transform-style: preserve-3d;
  backface-visibility: hidden;
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.16),
    inset 0 -1px 0 rgba(0, 0, 0, 0.22);
}

.score-flip-top {
  transform-origin: 50% 100%;
  animation: scoreboard-top-flip 420ms cubic-bezier(0.45, 0, 0.2, 1) forwards;
}

.score-flip-bottom {
  z-index: 2;
  transform-origin: 50% 0%;
  transform: rotateX(90deg);
  animation: scoreboard-bottom-flip 460ms 360ms cubic-bezier(0.16, 0.85, 0.22, 1) forwards;
}

.score-flip::after {
  content: '';
  position: absolute;
  inset: 0;
  pointer-events: none;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.08), rgba(0, 0, 0, 0.18));
  opacity: 0;
}

.score-flip-top::after {
  animation: scoreboard-top-shade 420ms cubic-bezier(0.45, 0, 0.2, 1) forwards;
}

.score-flip-bottom::after {
  animation: scoreboard-bottom-shade 460ms 360ms cubic-bezier(0.16, 0.85, 0.22, 1) forwards;
}

.flip-crease {
  position: absolute;
  left: 0;
  right: 0;
  top: 50%;
  z-index: 5;
  height: 1px;
  background: rgba(0, 0, 0, 0.42);
  box-shadow:
    0 -1px 0 rgba(255, 255, 255, 0.1),
    0 1px 0 rgba(255, 255, 255, 0.06),
    0 7px 16px rgba(0, 0, 0, 0.2);
}

.flip-value {
  position: relative;
  z-index: 1;
  color: #fff;
  font-size: 44px;
  line-height: 1;
  font-weight: 850;
  letter-spacing: 0;
  text-shadow: 0 8px 18px rgba(0, 0, 0, 0.28);
}

.flip-panel-top .flip-value,
.score-flip-top .flip-value {
  transform: translateY(50%);
}

.flip-panel-bottom .flip-value,
.score-flip-bottom .flip-value {
  transform: translateY(-50%);
}

.digits-medium .flip-value {
  font-size: 36px;
}

.digits-tight .flip-value {
  font-size: 30px;
}

.time-label {
  display: block;
  color: rgba(255, 255, 255, 0.64);
  font-size: 12px;
  font-weight: 700;
}

.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  clip-path: inset(50%);
}

@keyframes scoreboard-top-flip {
  0% {
    transform: rotateX(0deg);
  }

  72% {
    transform: rotateX(-88deg);
  }

  100% {
    transform: rotateX(-92deg);
  }
}

@keyframes scoreboard-bottom-flip {
  0% {
    transform: rotateX(92deg);
  }

  100% {
    transform: rotateX(0deg);
  }
}

@keyframes scoreboard-top-shade {
  0% {
    opacity: 0;
  }

  100% {
    opacity: 0.46;
  }
}

@keyframes scoreboard-bottom-shade {
  0% {
    opacity: 0.42;
  }

  100% {
    opacity: 0;
  }
}

@media (min-width: 768px) {
  .flip-stack {
    height: 108px;
  }

  .flip-value {
    font-size: 56px;
  }

  .digits-medium .flip-value {
    font-size: 46px;
  }

  .digits-tight .flip-value {
    font-size: 38px;
  }
}

@media (prefers-reduced-motion: reduce) {
  .score-flip {
    display: none;
    animation: none;
  }
}
</style>
