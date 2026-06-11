<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'

interface Trail {
  id: number
  x: number
  y: number
  size: number
  angle: number
}

interface GlowDot {
  id: string
  x: number
  y: number
  size: number
  opacity: number
  lowOpacity: number
  delay: number
  ring: number
  tone: number
}

const root = ref<HTMLElement | null>(null)
const active = ref(false)
const trails = ref<Trail[]>([])

const glowDots: GlowDot[] = [
  { count: 1, radius: 0 },
  { count: 8, radius: 42 },
  { count: 12, radius: 78 },
  { count: 16, radius: 116 },
  { count: 20, radius: 154 },
  { count: 24, radius: 192 },
].flatMap((ring, ringIndex) => {
  const offset = ringIndex % 2 ? Math.PI / ring.count : 0
  return Array.from({ length: ring.count }, (_, dotIndex) => {
    const angle = offset + (Math.PI * 2 * dotIndex) / ring.count
    const fade = Math.max(0.22, 1 - ringIndex * 0.13)
    return {
      id: `${ringIndex}-${dotIndex}`,
      x: Math.cos(angle) * ring.radius,
      y: Math.sin(angle) * ring.radius,
      size: Math.max(3, 7 - ringIndex * 0.65),
      opacity: fade,
      lowOpacity: fade * (0.36 + ringIndex * 0.025),
      delay: (ringIndex * 0.42 + dotIndex * 0.08) % 2.8,
      ring: ringIndex,
      tone: (dotIndex / ring.count) * 360,
    }
  })
})

let host: HTMLElement | null = null
let frame = 0
let trailId = 0
let lastTrailTime = 0
let lastX = 0
let lastY = 0
let canRun = false
const trailTimers = new Map<number, number>()

function supportsPointerGlow() {
  if (typeof window === 'undefined') return false
  const reducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
  return !reducedMotion
}

function setGlowPosition(x: number, y: number) {
  const el = root.value
  if (!el) return
  el.style.setProperty('--glow-x', `${x}px`)
  el.style.setProperty('--glow-y', `${y}px`)
}

function addTrail(x: number, y: number, distance: number, angle: number) {
  const now = performance.now()
  if (now - lastTrailTime < 28 || distance < 14) return
  lastTrailTime = now

  const id = ++trailId
  trails.value = [
    ...trails.value.slice(-11),
    {
      id,
      x,
      y,
      angle,
      size: Math.min(26, Math.max(12, distance * 0.32)),
    },
  ]

  const timer = window.setTimeout(() => {
    trails.value = trails.value.filter((trail) => trail.id !== id)
    trailTimers.delete(id)
  }, 720)
  trailTimers.set(id, timer)
}

function handlePointerMove(event: PointerEvent) {
  if (!host || !canRun || event.pointerType !== 'mouse') return
  const rect = host.getBoundingClientRect()
  const x = event.clientX - rect.left
  const y = event.clientY - rect.top
  const dx = x - lastX
  const dy = y - lastY
  const distance = Math.hypot(dx, dy)

  if (!frame) {
    frame = window.requestAnimationFrame(() => {
      frame = 0
      setGlowPosition(x, y)
    })
  }

  if (active.value) {
    addTrail(x, y, distance, Math.atan2(dy, dx))
  }

  active.value = true
  lastX = x
  lastY = y
}

function handlePointerEnter(event: PointerEvent) {
  if (!host || !canRun || event.pointerType !== 'mouse') return
  const rect = host.getBoundingClientRect()
  lastX = event.clientX - rect.left
  lastY = event.clientY - rect.top
  setGlowPosition(lastX, lastY)
  active.value = true
}

function handlePointerLeave() {
  active.value = false
}

onMounted(() => {
  canRun = supportsPointerGlow()
  if (!canRun) return

  host = root.value?.parentElement || null
  if (!host) return

  host.addEventListener('pointerenter', handlePointerEnter)
  host.addEventListener('pointermove', handlePointerMove)
  host.addEventListener('pointerleave', handlePointerLeave)
})

onBeforeUnmount(() => {
  if (host) {
    host.removeEventListener('pointerenter', handlePointerEnter)
    host.removeEventListener('pointermove', handlePointerMove)
    host.removeEventListener('pointerleave', handlePointerLeave)
  }
  if (frame) window.cancelAnimationFrame(frame)
  trailTimers.forEach((timer) => window.clearTimeout(timer))
  trailTimers.clear()
})
</script>

<template>
  <div ref="root" class="pointer-glow" :class="{ active }" aria-hidden="true">
    <div class="pointer-glow-ambient"></div>
    <div class="pointer-glow-field"></div>
    <div class="pointer-glow-dot-field">
      <span
        v-for="dot in glowDots"
        :key="dot.id"
        class="pointer-glow-dot"
        :style="{
          '--dot-x': `${dot.x}px`,
          '--dot-y': `${dot.y}px`,
          '--dot-size': `${dot.size}px`,
          '--dot-opacity': `${dot.opacity}`,
          '--dot-low-opacity': `${dot.lowOpacity}`,
          '--dot-delay': `${dot.delay}s`,
          '--dot-ring': `${dot.ring}`,
          '--dot-tone': `${dot.tone}deg`,
        }"
      ></span>
    </div>
    <span
      v-for="trail in trails"
      :key="trail.id"
      class="pointer-glow-trail"
      :style="{
        left: `${trail.x}px`,
        top: `${trail.y}px`,
        width: `${trail.size * 2.8}px`,
        height: `${trail.size}px`,
        rotate: `${trail.angle}rad`,
      }"
    ></span>
  </div>
</template>

<style scoped>
.pointer-glow {
  --glow-x: 50%;
  --glow-y: 18%;

  position: absolute;
  inset: -40px;
  z-index: 2;
  overflow: hidden;
  pointer-events: none;
  opacity: 0.72;
  contain: layout paint;
  transition: opacity 420ms ease-out;
}

.pointer-glow.active {
  opacity: 1;
}

.pointer-glow-ambient {
  position: absolute;
  inset: -14%;
  background:
    radial-gradient(circle 220px at 22% 24%, rgba(96, 165, 250, 0.18), transparent 72%),
    radial-gradient(circle 260px at 74% 20%, rgba(147, 197, 253, 0.12), transparent 70%),
    radial-gradient(circle 190px at 46% 74%, rgba(244, 114, 182, 0.1), transparent 76%);
  filter: blur(22px) saturate(1.1);
  mix-blend-mode: screen;
  opacity: 0.46;
  transform: translate3d(0, 0, 0);
  animation: pointer-glow-ambient 18s ease-in-out infinite alternate;
}

.pointer-glow-field {
  position: absolute;
  inset: 0;
  overflow: hidden;
  background:
    radial-gradient(
      circle 44px at var(--glow-x) var(--glow-y),
      rgba(255, 255, 255, 0.46),
      rgba(219, 234, 254, 0.22) 38%,
      transparent 74%
    ),
    radial-gradient(
      circle 190px at var(--glow-x) var(--glow-y),
      rgba(219, 234, 254, 0.3),
      rgba(96, 165, 250, 0.18) 32%,
      rgba(244, 114, 182, 0.1) 52%,
      transparent 74%
    ),
    radial-gradient(
      circle 310px at var(--glow-x) var(--glow-y),
      rgba(59, 130, 246, 0.08),
      transparent 72%
    );
  filter: blur(8px) saturate(1.16);
  mix-blend-mode: screen;
  animation: pointer-glow-color 16s ease-in-out infinite alternate;
}

.pointer-glow-field::before,
.pointer-glow-field::after {
  content: '';
  position: absolute;
  left: calc(var(--glow-x) - 210px);
  top: calc(var(--glow-y) - 210px);
  width: 420px;
  height: 420px;
  border-radius: 999px;
  pointer-events: none;
  opacity: 0.58;
  mix-blend-mode: screen;
  transform-origin: 50% 50%;
}

.pointer-glow-field::before {
  background:
    conic-gradient(
      from 18deg,
      transparent 0deg,
      rgba(219, 234, 254, 0.2) 24deg,
      rgba(96, 165, 250, 0.18) 54deg,
      transparent 92deg,
      rgba(244, 114, 182, 0.14) 136deg,
      transparent 186deg,
      rgba(147, 197, 253, 0.18) 246deg,
      transparent 318deg,
      transparent 360deg
    );
  filter: blur(20px);
  mask-image: radial-gradient(circle, #000 0 36%, transparent 72%);
  animation: pointer-glow-stream 12s linear infinite;
}

.pointer-glow-field::after {
  background:
    linear-gradient(
      115deg,
      transparent 26%,
      rgba(219, 234, 254, 0.2) 42%,
      rgba(96, 165, 250, 0.16) 48%,
      transparent 62%
    );
  filter: blur(12px);
  mask-image: radial-gradient(ellipse at center, transparent 0 18%, #000 34%, transparent 70%);
  animation: pointer-glow-sweep 10s ease-in-out infinite alternate;
}

.pointer-glow-dot-field {
  position: absolute;
  left: var(--glow-x);
  top: var(--glow-y);
  width: 1px;
  height: 1px;
  mix-blend-mode: screen;
  animation: pointer-glow-dots-spin 28s linear infinite;
}

.pointer-glow-dot-field::before,
.pointer-glow-dot-field::after {
  content: '';
  position: absolute;
  left: -236px;
  top: -236px;
  width: 472px;
  height: 472px;
  border-radius: 999px;
  pointer-events: none;
}

.pointer-glow-dot-field::before {
  border: 1px solid rgba(147, 197, 253, 0.13);
  box-shadow:
    inset 0 0 46px rgba(96, 165, 250, 0.08),
    0 0 34px rgba(147, 197, 253, 0.08);
  opacity: 0.72;
  scale: 0.92;
  animation: pointer-glow-orbit-ring 9s ease-in-out infinite alternate;
}

.pointer-glow-dot-field::after {
  background:
    conic-gradient(
      from 20deg,
      transparent 0deg,
      transparent 62deg,
      rgba(219, 234, 254, 0.22) 96deg,
      rgba(96, 165, 250, 0.16) 118deg,
      transparent 156deg,
      transparent 242deg,
      rgba(244, 114, 182, 0.13) 284deg,
      transparent 326deg,
      transparent 360deg
    );
  filter: blur(11px);
  mask-image: radial-gradient(circle, transparent 0 28%, #000 42%, transparent 72%);
  opacity: 0.54;
  animation: pointer-glow-orbit-sweep 14s linear infinite;
}

.pointer-glow-dot {
  position: absolute;
  left: 0;
  top: 0;
  width: var(--dot-size);
  height: var(--dot-size);
  border-radius: 999px;
  background:
    radial-gradient(circle, rgba(255, 255, 255, 0.95), rgba(147, 197, 253, 0.85) 42%, rgba(59, 130, 246, 0.32) 72%, transparent);
  box-shadow:
    0 0 10px rgba(147, 197, 253, 0.42),
    0 0 22px rgba(59, 130, 246, 0.16);
  opacity: calc(var(--dot-opacity) * 0.82);
  transform: translate(-50%, -50%) translate(var(--dot-x), var(--dot-y));
  animation:
    pointer-glow-dot-breathe 4.8s ease-in-out var(--dot-delay) infinite alternate,
    pointer-glow-dot-color 12s ease-in-out var(--dot-delay) infinite alternate;
}

.pointer-glow-dot:nth-child(3n) {
  background:
    radial-gradient(circle, rgba(255, 255, 255, 0.92), rgba(191, 219, 254, 0.78) 44%, rgba(14, 165, 233, 0.28) 74%, transparent);
}

.pointer-glow-dot:nth-child(4n) {
  box-shadow:
    0 0 12px rgba(244, 114, 182, 0.28),
    0 0 26px rgba(96, 165, 250, 0.14);
}

.pointer-glow-trail {
  position: absolute;
  border-radius: 999px;
  background:
    radial-gradient(circle at 18% 50%, rgba(255, 255, 255, 0.82), transparent 24%),
    linear-gradient(
      90deg,
      color-mix(in srgb, var(--accent) 66%, transparent),
      color-mix(in srgb, var(--primary) 48%, transparent),
      transparent
    );
  filter: blur(1px) saturate(1.24);
  mix-blend-mode: screen;
  transform: translate(-50%, -50%);
  animation: pointer-glow-trail 720ms ease-out forwards;
}

@keyframes pointer-glow-trail {
  0% {
    opacity: 0.82;
    scale: 0.72;
  }

  100% {
    opacity: 0;
    scale: 1.8;
  }
}

@keyframes pointer-glow-ambient {
  0% {
    filter: blur(18px) saturate(1.12) hue-rotate(0deg);
    opacity: 0.46;
    transform: translate3d(-3%, -2%, 0) scale(1);
  }

  42% {
    filter: blur(20px) saturate(1.26) hue-rotate(18deg);
    opacity: 0.68;
    transform: translate3d(4%, 2%, 0) scale(1.06);
  }

  100% {
    filter: blur(18px) saturate(1.18) hue-rotate(-16deg);
    opacity: 0.56;
    transform: translate3d(1%, 5%, 0) scale(1.03);
  }
}

@keyframes pointer-glow-color {
  0% {
    filter: blur(7px) saturate(1.12) hue-rotate(0deg);
  }

  50% {
    filter: blur(8px) saturate(1.28) hue-rotate(14deg);
  }

  100% {
    filter: blur(7px) saturate(1.18) hue-rotate(-12deg);
  }
}

@keyframes pointer-glow-dots-spin {
  0% {
    filter: hue-rotate(0deg) saturate(1.04);
    transform: rotate(0deg) scale(0.98);
  }

  50% {
    filter: hue-rotate(12deg) saturate(1.18);
    transform: rotate(5deg) scale(1.035);
  }

  100% {
    filter: hue-rotate(-8deg) saturate(1.08);
    transform: rotate(0deg) scale(0.98);
  }
}

@keyframes pointer-glow-dot-breathe {
  0% {
    opacity: var(--dot-low-opacity);
    scale: 0.78;
  }

  100% {
    opacity: calc(var(--dot-opacity) * 0.96);
    scale: 1.14;
  }
}

@keyframes pointer-glow-dot-color {
  0% {
    filter: hue-rotate(var(--dot-tone)) saturate(1);
  }

  100% {
    filter: hue-rotate(calc(var(--dot-tone) + 22deg)) saturate(1.3);
  }
}

@keyframes pointer-glow-orbit-ring {
  0% {
    opacity: 0.42;
    transform: scale(0.86);
  }

  100% {
    opacity: 0.76;
    transform: scale(1);
  }
}

@keyframes pointer-glow-orbit-sweep {
  0% {
    opacity: 0.34;
    transform: rotate(0deg) scale(0.94);
  }

  50% {
    opacity: 0.58;
    transform: rotate(180deg) scale(1.02);
  }

  100% {
    opacity: 0.36;
    transform: rotate(360deg) scale(0.94);
  }
}

@keyframes pointer-glow-stream {
  0% {
    opacity: 0.38;
    transform: rotate(0deg) scale(0.92);
  }

  46% {
    opacity: 0.68;
    transform: rotate(168deg) scale(1.04);
  }

  100% {
    opacity: 0.42;
    transform: rotate(360deg) scale(0.96);
  }
}

@keyframes pointer-glow-sweep {
  0% {
    opacity: 0.22;
    transform: translate3d(-18px, 14px, 0) rotate(-18deg) scale(0.88);
  }

  52% {
    opacity: 0.5;
  }

  100% {
    opacity: 0.34;
    transform: translate3d(22px, -12px, 0) rotate(24deg) scale(1.05);
  }
}

@media (prefers-reduced-motion: reduce) {
  .pointer-glow {
    display: none;
  }
}
</style>
