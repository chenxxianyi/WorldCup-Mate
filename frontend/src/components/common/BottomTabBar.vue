<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const tabs = [
  { name: '首页', icon: 'home', route: '/' },
  { name: '赛程', icon: 'calendar_month', route: '/schedule' },
  { name: '球队', icon: 'flag', route: '/teams' },
  { name: 'AI', icon: 'auto_awesome', route: '/ai' },
  { name: '我的', icon: 'person', route: '/profile' },
]

function isActive(tabRoute: string) {
  if (tabRoute === '/') return route.path === '/'
  return route.path.startsWith(tabRoute)
}

function navigate(tabRoute: string) {
  router.push(tabRoute)
}
</script>

<template>
  <nav class="bottom-nav" aria-label="底部导航">
    <button
      v-for="tab in tabs"
      :key="tab.route"
      class="nav-item"
      :class="{ active: isActive(tab.route) }"
      @click="navigate(tab.route)"
    >
      <div class="icon-wrapper">
        <span class="material-symbols-outlined">{{ tab.icon }}</span>
      </div>
      <span>{{ tab.name }}</span>
    </button>
  </nav>
</template>

<style scoped>
.bottom-nav {
  position: fixed;
  z-index: 20;
  left: 50%;
  bottom: calc(14px + env(safe-area-inset-bottom));
  transform: translateX(-50%);
  width: min(calc(100% - 32px), 430px);
  height: 62px;
  display: flex;
  align-items: center;
  justify-content: space-around;
  padding: 0 8px;
  border: 1px solid var(--line);
  border-radius: 999px;
  background: color-mix(in srgb, var(--card) 94%, transparent);
  box-shadow: 0 10px 28px rgba(17, 17, 17, 0.07);
  backdrop-filter: blur(14px);
  -webkit-backdrop-filter: blur(14px);
}

.nav-item {
  border: 0;
  flex: 1;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 3px;
  color: #8e8e93;
  background: transparent;
  text-decoration: none;
  position: relative;
  transition: color 180ms ease-out;
}

.nav-item:active {
  transform: scale(0.97);
}

.nav-item.active {
  color: var(--primary);
}

.icon-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.nav-item.active .material-symbols-outlined {
  filter: none;
}

.nav-item > span:not(.material-symbols-outlined) {
  font-size: 12px;
  line-height: 1;
  font-weight: 650;
  letter-spacing: 0;
  transition: color 180ms ease-out;
}

[data-theme='dark'] .bottom-nav {
  border-color: var(--line);
  background: color-mix(in srgb, var(--card) 94%, transparent);
  box-shadow: 0 10px 28px rgba(0, 0, 0, 0.28);
}

@media (min-width: 768px) {
  .bottom-nav {
    display: none;
  }
}
</style>
