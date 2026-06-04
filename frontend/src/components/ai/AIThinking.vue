<template>
  <div class="thinking">
    <span></span>
    <span></span>
    <span></span>
  </div>
</template>

<style scoped>
.thinking {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  min-height: 20px;
}

.thinking span {
  width: 6px;
  height: 6px;
  border-radius: 999px;
  background: var(--primary);
  opacity: 0.35;
  animation: bounce 1.1s infinite ease-in-out;
}

.thinking span:nth-child(2) {
  animation-delay: 0.15s;
}

.thinking span:nth-child(3) {
  animation-delay: 0.3s;
}

@keyframes bounce {
  0%, 80%, 100% {
    transform: translateY(0);
    opacity: 0.35;
  }

  40% {
    transform: translateY(-5px);
    opacity: 1;
  }
}
</style>
