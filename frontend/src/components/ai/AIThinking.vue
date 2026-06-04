<template>
  <div class="thinking" role="status" aria-label="AI 正在回复">
    <span class="dot blue"></span>
    <span class="dot red"></span>
    <span class="dot yellow"></span>
    <span class="dot green"></span>
  </div>
</template>

<style scoped>
.thinking {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  min-height: 24px;
  padding: 2px 0;
}

.dot {
  width: 7px;
  height: 7px;
  border-radius: 999px;
  opacity: 0.36;
  transform: scale(0.72);
  animation: gemini-pulse 1.35s infinite cubic-bezier(0.4, 0, 0.2, 1);
}

.dot.blue {
  background: #4285f4;
}

.dot.red {
  background: #ea4335;
  animation-delay: 0.14s;
}

.dot.yellow {
  background: #fbbc04;
  animation-delay: 0.28s;
}

.dot.green {
  background: #34a853;
  animation-delay: 0.42s;
}

@keyframes gemini-pulse {
  0%,
  70%,
  100% {
    opacity: 0.32;
    transform: translateY(0) scale(0.72);
  }

  28% {
    opacity: 1;
    transform: translateY(-2px) scale(1);
  }
}

@media (prefers-reduced-motion: reduce) {
  .dot {
    animation: none;
    opacity: 0.8;
    transform: none;
  }
}
</style>
