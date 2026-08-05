<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const tabs = [
  { name: '首页', icon: 'home', route: '/' },
  { name: '赛程', icon: 'calendar_month', route: '/schedule' },
  { name: '球队', icon: 'flag', route: '/teams' },
  { name: '积分榜', icon: 'trophy', route: '/standings' },
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
  <nav
    class="bottom-nav"
    aria-label="底部导航"
  >
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
  bottom: calc(18px + env(safe-area-inset-bottom));
  transform: translateX(-50%);
  width: min(calc(100% - 40px), 430px);
  height: 68px;
  display: flex;
  align-items: center;
  justify-content: space-around;
  padding: 0 10px;
  border: 1px solid rgba(255, 255, 255, 0.4);
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.85);
  box-shadow:
    0 8px 32px rgba(0, 0, 0, 0.05),
    inset 0 1px 0 rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
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
  transition: all 300ms ease;
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
  filter:
    drop-shadow(0 0 10px rgba(227, 29, 36, 0.5))
    drop-shadow(0 0 3px rgba(227, 29, 36, 0.3));
}

.nav-item > span:not(.material-symbols-outlined) {
  font-size: 12px;
  line-height: 1;
  font-weight: 600;
  letter-spacing: 0.5px;
  transition: all 300ms ease;
}

[data-theme='dark'] .bottom-nav {
  border-color: rgba(255, 255, 255, 0.08);
  background: rgba(17, 17, 17, 0.78);
  box-shadow:
    0 8px 32px rgba(0, 0, 0, 0.36),
    inset 0 1px 0 rgba(255, 255, 255, 0.08);
}

@media (min-width: 768px) {
  .bottom-nav {
    display: none;
  }
}
</style>
